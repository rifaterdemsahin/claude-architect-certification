package main

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	_ "image/png"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

func sasPerms(raw string) string {
	const order = "racwdxltfmoepiyq"
	var b strings.Builder
	for _, c := range order {
		if strings.ContainsRune(raw, c) {
			b.WriteRune(c)
		}
	}
	return b.String()
}

func generateContainerSAS(accountName, accountKey, container, permissions string, expiry time.Time) (string, error) {
	permissions = sasPerms(permissions)
	const version = "2026-04-06"
	expiryStr := expiry.UTC().Format("2006-01-02T15:04:05Z")
	canonResource := "/blob/" + accountName + "/" + container
	stringToSign := strings.Join([]string{
		permissions, "", expiryStr, canonResource,
		"", "", "https", version, "c",
		"", "", "", "", "", "", "",
	}, "\n")
	keyBytes, err := base64.StdEncoding.DecodeString(accountKey)
	if err != nil {
		return "", fmt.Errorf("decode account key: %w", err)
	}
	mac := hmac.New(sha256.New, keyBytes)
	mac.Write([]byte(stringToSign))
	sig := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	params := url.Values{
		"sv": {version}, "se": {expiryStr}, "sr": {"c"},
		"sp": {permissions}, "spr": {"https"}, "sig": {sig},
	}
	return params.Encode(), nil
}

var allowedResearchContainers = map[string]bool{
	"research-images":       true,
	"research-audio":        true,
	"research-videos":       true,
	"research-notes":        true,
	"research-animations":   true,
	"research-infographics": true,
}

type blobInfo struct {
	Name         string `json:"name"`
	Size         int64  `json:"size"`
	ContentType  string `json:"contentType"`
	LastModified string `json:"lastModified"`
}

type blobListXML struct {
	XMLName xml.Name `xml:"EnumerationResults"`
	Blobs   []struct {
		Name  string `xml:"Name"`
		Props struct {
			LastModified  string `xml:"Last-Modified"`
			ContentLength int64  `xml:"Content-Length"`
			ContentType   string `xml:"Content-Type"`
		} `xml:"Properties"`
	} `xml:"Blobs>Blob"`
}

func blobURL(accountName, container, name, sasQuery string) string {
	base := fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, container)
	if name != "" {
		base += "/" + url.PathEscape(name)
	}
	if sasQuery == "" {
		return base
	}
	return base + "?" + sasQuery
}

func researchFileProxyURL(container, name string) string {
	return "/api/research/file?container=" + container + "&name=" + url.QueryEscape(name)
}

func uploadBlobToAzure(ctx context.Context, cfg config, container, blobName, contentType string, data []byte) error {
	expiry := time.Now().UTC().Add(10 * time.Minute)
	sasQuery, err := generateContainerSAS(cfg.azureAccountName, cfg.azureAccountKey, container, "rcwl", expiry)
	if err != nil {
		return fmt.Errorf("sas: %w", err)
	}
	putURL := blobURL(cfg.azureAccountName, container, blobName, sasQuery)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set("x-ms-blob-type", "BlockBlob")
	req.Header.Set("Content-Type", contentType)
	req.ContentLength = int64(len(data))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("azure upload %d: %s", resp.StatusCode, strings.TrimSpace(string(b)))
	}
	return nil
}

func downloadBlobFromAzure(ctx context.Context, cfg config, container, blobName string) ([]byte, string, error) {
	expiry := time.Now().UTC().Add(10 * time.Minute)
	sasQuery, err := generateContainerSAS(cfg.azureAccountName, cfg.azureAccountKey, container, "rl", expiry)
	if err != nil {
		return nil, "", fmt.Errorf("sas: %w", err)
	}
	getURL := blobURL(cfg.azureAccountName, container, blobName, sasQuery)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, "", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return nil, "", fmt.Errorf("azure get %d", resp.StatusCode)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	ct := resp.Header.Get("Content-Type")
	if ct == "" {
		ct = "application/octet-stream"
	}
	return data, ct, nil
}

func deleteBlobFromAzure(ctx context.Context, cfg config, container, blobName string) error {
	expiry := time.Now().UTC().Add(10 * time.Minute)
	sasQuery, err := generateContainerSAS(cfg.azureAccountName, cfg.azureAccountKey, container, "rdl", expiry)
	if err != nil {
		return fmt.Errorf("sas: %w", err)
	}
	delURL := blobURL(cfg.azureAccountName, container, blobName, sasQuery)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, delURL, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()
	if resp.StatusCode >= 400 {
		return fmt.Errorf("azure delete %d", resp.StatusCode)
	}
	return nil
}

const thumbPrefix = "__thumb__"

func thumbBlobName(original string) string {
	return thumbPrefix + original
}

func generateThumbnail(data []byte, maxEdge int) ([]byte, error) {
	src, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}
	b := src.Bounds()
	sw, sh := b.Dx(), b.Dy()
	if sw == 0 || sh == 0 {
		return nil, fmt.Errorf("empty image")
	}
	tw, th := sw, sh
	if sw > maxEdge || sh > maxEdge {
		if sw >= sh {
			tw, th = maxEdge, int(float64(sh)*float64(maxEdge)/float64(sw))
		} else {
			th, tw = maxEdge, int(float64(sw)*float64(maxEdge)/float64(sh))
		}
	}
	if tw < 1 {
		tw = 1
	}
	if th < 1 {
		th = 1
	}
	dst := image.NewRGBA(image.Rect(0, 0, tw, th))
	for y := 0; y < th; y++ {
		sy0 := b.Min.Y + y*sh/th
		sy1 := b.Min.Y + (y+1)*sh/th
		if sy1 <= sy0 {
			sy1 = sy0 + 1
		}
		for x := 0; x < tw; x++ {
			sx0 := b.Min.X + x*sw/tw
			sx1 := b.Min.X + (x+1)*sw/tw
			if sx1 <= sx0 {
				sx1 = sx0 + 1
			}
			var rs, gs, bs, as, n uint64
			for sy := sy0; sy < sy1; sy++ {
				for sx := sx0; sx < sx1; sx++ {
					r, g, bb, a := src.At(sx, sy).RGBA()
					rs += uint64(r)
					gs += uint64(g)
					bs += uint64(bb)
					as += uint64(a)
					n++
				}
			}
			if n == 0 {
				n = 1
			}
			dst.Set(x, y, color.RGBA{
				R: uint8((rs / n) >> 8),
				G: uint8((gs / n) >> 8),
				B: uint8((bs / n) >> 8),
				A: uint8((as / n) >> 8),
			})
		}
	}
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, dst, &jpeg.Options{Quality: 80}); err != nil {
		return nil, fmt.Errorf("encode: %w", err)
	}
	return buf.Bytes(), nil
}

func researchUploadHandler(cfg config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if !isAdminRequest(r) {
			http.Error(w, "Unauthorized — sign in as admin to upload.", http.StatusUnauthorized)
			return
		}
		container := r.URL.Query().Get("container")
		if !allowedResearchContainers[container] {
			http.Error(w, "invalid container", http.StatusBadRequest)
			return
		}
		if cfg.azureAccountName == "" {
			http.Error(w, "Azure Storage not configured", http.StatusServiceUnavailable)
			return
		}
		if err := r.ParseMultipartForm(100 << 20); err != nil {
			http.Error(w, "parse form: "+err.Error(), http.StatusBadRequest)
			return
		}
		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "missing file field", http.StatusBadRequest)
			return
		}
		defer file.Close()

		blobName := strings.ReplaceAll(header.Filename, "/", "_")
		blobName = strings.ReplaceAll(blobName, "\\", "_")
		if blobName == "" || blobName == "." {
			blobName = fmt.Sprintf("file-%d", time.Now().UnixMilli())
		}
		contentType := header.Header.Get("Content-Type")
		if contentType == "" {
			contentType = "application/octet-stream"
		}

		data, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "read upload: "+err.Error(), http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		if err := uploadBlobToAzure(ctx, cfg, container, blobName, contentType, data); err != nil {
			http.Error(w, "azure upload: "+err.Error(), http.StatusBadGateway)
			return
		}

		var thumbName string
		if thumbData, terr := generateThumbnail(data, 320); terr != nil {
			log.Printf("research thumbnail generation skipped for %s: %v", blobName, terr)
		} else {
			tName := thumbBlobName(blobName)
			if uerr := uploadBlobToAzure(ctx, cfg, container, tName, "image/jpeg", thumbData); uerr != nil {
				log.Printf("research thumbnail upload failed for %s: %v", tName, uerr)
			} else {
				thumbName = tName
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"ok": "true", "name": blobName, "thumbnail": thumbName})
	}
}

func researchFilesHandler(cfg config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		container := r.URL.Query().Get("container")
		if !allowedResearchContainers[container] {
			http.Error(w, "invalid container", http.StatusBadRequest)
			return
		}
		if cfg.azureAccountName == "" {
			http.Error(w, "Azure Storage not configured", http.StatusServiceUnavailable)
			return
		}
		expiry := time.Now().UTC().Add(5 * time.Minute)
		sasQuery, err := generateContainerSAS(cfg.azureAccountName, cfg.azureAccountKey, container, "rl", expiry)
		if err != nil {
			http.Error(w, "sas error", http.StatusInternalServerError)
			return
		}
		listURL := fmt.Sprintf("https://%s.blob.core.windows.net/%s?restype=container&comp=list&%s",
			cfg.azureAccountName, container, sasQuery)
		req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, listURL, nil)
		if err != nil {
			http.Error(w, "request build", http.StatusInternalServerError)
			return
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			http.Error(w, "azure list: "+err.Error(), http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 400 {
			b, _ := io.ReadAll(resp.Body)
			http.Error(w, fmt.Sprintf("azure %d: %s", resp.StatusCode, b), http.StatusBadGateway)
			return
		}
		var listResult blobListXML
		if err := xml.NewDecoder(resp.Body).Decode(&listResult); err != nil {
			http.Error(w, "decode: "+err.Error(), http.StatusInternalServerError)
			return
		}
		blobs := make([]blobInfo, 0, len(listResult.Blobs))
		for _, b := range listResult.Blobs {
			blobs = append(blobs, blobInfo{
				Name:         b.Name,
				Size:         b.Props.ContentLength,
				ContentType:  b.Props.ContentType,
				LastModified: b.Props.LastModified,
			})
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(blobs)
	}
}

func researchFileHandler(cfg config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		container := r.URL.Query().Get("container")
		name := r.URL.Query().Get("name")
		if !allowedResearchContainers[container] || name == "" {
			http.Error(w, "invalid params", http.StatusBadRequest)
			return
		}
		if cfg.azureAccountName == "" {
			http.Error(w, "Azure Storage not configured", http.StatusServiceUnavailable)
			return
		}
		expiry := time.Now().UTC().Add(5 * time.Minute)
		switch r.Method {
		case http.MethodGet:
			sasQuery, err := generateContainerSAS(cfg.azureAccountName, cfg.azureAccountKey, container, "rl", expiry)
			if err != nil {
				http.Error(w, "sas error", http.StatusInternalServerError)
				return
			}
			getURL := blobURL(cfg.azureAccountName, container, name, sasQuery)
			req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, getURL, nil)
			if err != nil {
				http.Error(w, "request build", http.StatusInternalServerError)
				return
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				http.Error(w, "azure get: "+err.Error(), http.StatusBadGateway)
				return
			}
			defer resp.Body.Close()
			if resp.StatusCode >= 400 {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			if ct := resp.Header.Get("Content-Type"); ct != "" {
				w.Header().Set("Content-Type", ct)
			}
			if cl := resp.Header.Get("Content-Length"); cl != "" {
				w.Header().Set("Content-Length", cl)
			}
			w.Header().Set("Cache-Control", "public, max-age=300")
			io.Copy(w, resp.Body)

		case http.MethodDelete:
			if !isAdminRequest(r) {
				http.Error(w, "Unauthorized — sign in as admin to delete.", http.StatusUnauthorized)
				return
			}
			sasQuery, err := generateContainerSAS(cfg.azureAccountName, cfg.azureAccountKey, container, "rdl", expiry)
			if err != nil {
				http.Error(w, "sas error", http.StatusInternalServerError)
				return
			}
			delURL := blobURL(cfg.azureAccountName, container, name, sasQuery)
			req, err := http.NewRequestWithContext(r.Context(), http.MethodDelete, delURL, nil)
			if err != nil {
				http.Error(w, "request build", http.StatusInternalServerError)
				return
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				http.Error(w, "azure delete: "+err.Error(), http.StatusBadGateway)
				return
			}
			resp.Body.Close()
			if resp.StatusCode >= 400 {
				http.Error(w, "delete failed", http.StatusBadGateway)
				return
			}
			w.WriteHeader(http.StatusNoContent)

		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func researchRenameHandler(cfg config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if !isAdminRequest(r) {
			http.Error(w, "Unauthorized — sign in as admin to rename.", http.StatusUnauthorized)
			return
		}
		var body struct {
			Container string `json:"container"`
			From      string `json:"from"`
			To        string `json:"to"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		container := body.Container
		from := strings.TrimSpace(body.From)
		to := strings.TrimSpace(body.To)
		if !allowedResearchContainers[container] {
			http.Error(w, "invalid container", http.StatusBadRequest)
			return
		}
		if from == "" || to == "" {
			http.Error(w, "from and to are required", http.StatusBadRequest)
			return
		}
		if strings.ContainsAny(to, "/\\") || strings.HasPrefix(to, thumbPrefix) || strings.HasPrefix(from, thumbPrefix) {
			http.Error(w, "invalid name", http.StatusBadRequest)
			return
		}
		if !strings.EqualFold(path.Ext(from), path.Ext(to)) {
			http.Error(w, "file extension (suffix) must not change", http.StatusBadRequest)
			return
		}
		if to == from {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"ok": "true", "name": to})
			return
		}
		if cfg.azureAccountName == "" {
			http.Error(w, "Azure Storage not configured", http.StatusServiceUnavailable)
			return
		}
		ctx := r.Context()

		data, ct, err := downloadBlobFromAzure(ctx, cfg, container, from)
		if err != nil {
			http.Error(w, "source not found: "+err.Error(), http.StatusNotFound)
			return
		}
		if err := uploadBlobToAzure(ctx, cfg, container, to, ct, data); err != nil {
			http.Error(w, "azure copy: "+err.Error(), http.StatusBadGateway)
			return
		}
		if err := deleteBlobFromAzure(ctx, cfg, container, from); err != nil {
			log.Printf("research rename: original delete failed for %s: %v", from, err)
		}

		oldThumb, newThumb := thumbBlobName(from), thumbBlobName(to)
		if tData, tct, terr := downloadBlobFromAzure(ctx, cfg, container, oldThumb); terr == nil {
			if uerr := uploadBlobToAzure(ctx, cfg, container, newThumb, tct, tData); uerr != nil {
				log.Printf("research rename: thumbnail copy failed for %s: %v", newThumb, uerr)
			} else if derr := deleteBlobFromAzure(ctx, cfg, container, oldThumb); derr != nil {
				log.Printf("research rename: thumbnail delete failed for %s: %v", oldThumb, derr)
			}
		}

		if err := supabasePatch(ctx, cfg, "research_relationships",
			"item_name=eq."+url.QueryEscape(from),
			map[string]string{"item_name": to}); err != nil {
			log.Printf("research rename: relationship update failed %s→%s: %v", from, to, err)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"ok": "true", "name": to, "thumbnail": newThumb})
	}
}

func researchOcrHandler(cfg config) http.HandlerFunc {
	type OCRRequest struct {
		Container string `json:"container"`
		Name      string `json:"name"`
	}
	type OCRResponse struct {
		Text string `json:"text"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodPost {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		var req OCRRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
			return
		}

		if !allowedResearchContainers[req.Container] || req.Name == "" {
			http.Error(w, `{"error":"invalid params"}`, http.StatusBadRequest)
			return
		}

		if cfg.azureAccountName == "" {
			http.Error(w, `{"error":"Azure Storage not configured"}`, http.StatusServiceUnavailable)
			return
		}

		apiKey := cfg.getSecret("OPENROUTER_API_KEY")
		if apiKey == "" {
			http.Error(w, `{"error":"OPENROUTER_API_KEY missing from server configuration"}`, http.StatusServiceUnavailable)
			return
		}

		// 1. Download the image data from Azure Blob Storage
		data, contentType, err := downloadBlobFromAzure(r.Context(), cfg, req.Container, req.Name)
		if err != nil {
			log.Printf("Failed to download blob for OCR: %v", err)
			http.Error(w, fmt.Sprintf(`{"error":"failed to download image: %v"}`, err), http.StatusBadGateway)
			return
		}

		// Check if it's an image
		if !strings.HasPrefix(contentType, "image/") {
			ext := strings.ToLower(path.Ext(req.Name))
			isImgExt := ext == ".png" || ext == ".jpg" || ext == ".jpeg" || ext == ".webp" || ext == ".gif" || ext == ".svg"
			if !isImgExt {
				http.Error(w, `{"error":"file is not an image"}`, http.StatusBadRequest)
				return
			}
			if ext == ".png" {
				contentType = "image/png"
			} else if ext == ".webp" {
				contentType = "image/webp"
			} else if ext == ".gif" {
				contentType = "image/gif"
			} else {
				contentType = "image/jpeg"
			}
		}

		base64Image := base64.StdEncoding.EncodeToString(data)
		dataURL := fmt.Sprintf("data:%s;base64,%s", contentType, base64Image)

		// 2. Query OpenRouter with the image
		prompt := "Analyze this image and explain it in proper, clean Markdown format. Include a title representing the image content, a short descriptive overview explaining what the image depicts (e.g. system architecture, data flow diagram, code screenshot, or slide), and a structured list of all visible text organized cleanly under headings, bullet points, tables, or code blocks as appropriate. Do not include raw conversational greetings or sign-offs."

		ctx, cancel := context.WithTimeout(r.Context(), 45*time.Second)
		defer cancel()

		orURL := "https://openrouter.ai/api/v1/chat/completions"
		reqBody := map[string]any{
			"model": "google/gemini-2.5-flash",
			"messages": []map[string]any{
				{
					"role": "user",
					"content": []map[string]any{
						{
							"type": "text",
							"text": prompt,
						},
						{
							"type": "image_url",
							"image_url": map[string]string{
								"url": dataURL,
							},
						},
					},
				},
			},
		}

		b, err := json.Marshal(reqBody)
		if err != nil {
			http.Error(w, `{"error":"failed to marshal request"}`, http.StatusInternalServerError)
			return
		}

		hreq, err := http.NewRequestWithContext(ctx, "POST", orURL, bytes.NewReader(b))
		if err != nil {
			http.Error(w, `{"error":"failed to build request"}`, http.StatusInternalServerError)
			return
		}

		hreq.Header.Set("Content-Type", "application/json")
		hreq.Header.Set("Authorization", "Bearer "+apiKey)
		hreq.Header.Set("HTTP-Referer", "https://claude-architect-certification.com")
		hreq.Header.Set("X-Title", "Claude Architect Certification")

		hresp, err := http.DefaultClient.Do(hreq)
		if err != nil {
			log.Printf("OpenRouter OCR API call failed: %v", err)
			http.Error(w, `{"error":"OpenRouter API call failed"}`, http.StatusBadGateway)
			return
		}
		defer hresp.Body.Close()

		if hresp.StatusCode >= 400 {
			respBody, _ := io.ReadAll(hresp.Body)
			log.Printf("OpenRouter API returned error %d: %s", hresp.StatusCode, respBody)
			http.Error(w, fmt.Sprintf(`{"error":"OpenRouter API returned HTTP %d"}`, hresp.StatusCode), http.StatusBadGateway)
			return
		}

		var orResp struct {
			Choices []struct {
				Message struct {
					Content string `json:"content"`
				} `json:"message"`
			} `json:"choices"`
		}

		if err := json.NewDecoder(hresp.Body).Decode(&orResp); err != nil {
			http.Error(w, `{"error":"failed to decode response"}`, http.StatusInternalServerError)
			return
		}

		extractedText := ""
		if len(orResp.Choices) > 0 {
			extractedText = strings.TrimSpace(orResp.Choices[0].Message.Content)
		}

		resp := OCRResponse{
			Text: extractedText,
		}
		json.NewEncoder(w).Encode(resp)
	}
}
