package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type ImageGenRequest struct {
	Prompt       string   `json:"prompt"`
	ModuleNumber int      `json:"module_number"`
	VideoNumber  int      `json:"video_number"`
	AssetTypes   []string `json:"asset_types"`
}

type geminiContentResp struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text       string `json:"text"`
				InlineData struct {
					MimeType string `json:"mimeType"`
					Data     string `json:"data"`
				} `json:"inlineData"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
	UsageMetadata struct {
		PromptTokenCount     int `json:"promptTokenCount"`
		CandidatesTokenCount int `json:"candidatesTokenCount"`
		TotalTokenCount      int `json:"totalTokenCount"`
	} `json:"usageMetadata"`
}

var assetTypeStyles = map[string]string{
	"explain":      "Annotated diagram style with clear callouts, labels, and explanatory arrows — suitable for a slide or textbook",
	"infographic":  "Structured infographic layout with data points, icons, section dividers, statistics, and a clear visual hierarchy",
	"graphic":      "Full illustrative scene with cinematic lighting, rich detail, suitable for a video thumbnail or hero background, 16:9",
	"diagram":      "Clean technical diagram with labeled components, directional arrows, minimalist color palette, and precise geometry",
	"code":         "Stylised code snippet or terminal window with syntax highlighting, dark background, monospace font, suitable for commands, configs, YAML, or manifests",
	"comparison":   "Side-by-side or before/after split layout showing contrast between two items, with labels on each side and a clear divider",
	"stepbystep":   "Numbered sequence card with clear progression arrows, step indicators, or checklist design walking through a process",
	"thumbnail":    "Bold high-contrast 16:9 YouTube-optimised layout with prominent text zone, strong focal point, and space for overlay graphics",
	"architecture": "System architecture diagram with service boxes, trust zones, data flow arrows, and cloud service icons in a clean technical layout",
	"callout":      "Clean quote-card or highlight-card design emphasising one key insight, statistic, or takeaway with minimal surrounding detail",
	"timeline":     "Horizontal timeline with milestone nodes, dates, labels, and a clear progression arc from left to right",
	"table":        "Structured grid or matrix layout with labelled rows and columns, clear cell hierarchy, suitable for feature comparisons or tiered pricing",
	"titlecard":    "Clean transition or chapter-divider card with title, subtitle, decorative accent line, minimal background detail",
	"analogy":      "Illustrative metaphor image linking a technical concept to a familiar real-world scene or object, with subtle labelling",
	"background":   "Ambient environmental background with soft focal depth, abstract textures, cinematic lighting, and wide-angle perspective suitable for video backdrops",
	"transparent":  "Asset with no background, isolated subject on transparent canvas, clean edges, suitable for compositing in video editing",
	"icon":         "Single clean glyph or badge, scalable, minimal detail, high contrast, suitable for repeatable visual language across the course",
}

func buildAssetTypeInstruction(types []string) string {
	if len(types) == 0 {
		return "Style: Minimalist, dark corporate, glassmorphism, tech-focused, professional, 16:9."
	}
	var lines []string
	for _, t := range types {
		if style, ok := assetTypeStyles[t]; ok {
			lines = append(lines, fmt.Sprintf("  - %s", style))
		}
	}
	if len(lines) == 0 {
		return "Style: Minimalist, dark corporate, glassmorphism, tech-focused, professional, 16:9."
	}
	return fmt.Sprintf("Generate a prompt for the following asset type(s):\n%s\n\nThe visual style should be: dark corporate, glassmorphism, tech-focused, professional.", strings.Join(lines, "\n"))
}

func imageEnhancePromptHandler(cfg config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodPost {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			Prompt     string   `json:"prompt"`
			AssetTypes []string `json:"asset_types"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
			return
		}

		geminiKey := cfg.getSecret("GEMINI_API_KEY")
		if geminiKey == "" {
			http.Error(w, `{"error":"GEMINI_API_KEY missing"}`, http.StatusServiceUnavailable)
			return
		}

		typeInstr := buildAssetTypeInstruction(req.AssetTypes)
		refinePrompt := fmt.Sprintf(`You are an expert AI Image Prompt Engineer.
Refine the following user prompt into a high-quality, descriptive prompt for a professional image generator (like Midjourney or Imagen).
%s
Return ONLY the refined prompt text, no preamble or extra commentary.
User Prompt: %s`, typeInstr, req.Prompt)

		geminiURL := "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash:generateContent?key=" + geminiKey
		geminiReqBody := map[string]any{
			"contents": []any{map[string]any{"parts": []any{map[string]any{"text": refinePrompt}}}},
		}

		b, _ := json.Marshal(geminiReqBody)
		hresp, err := http.Post(geminiURL, "application/json", bytes.NewReader(b))
		if err != nil || hresp.StatusCode >= 400 {
			http.Error(w, `{"error":"Gemini refinement failed"}`, http.StatusBadGateway)
			return
		}
		defer hresp.Body.Close()

		var gResp struct {
			Candidates []struct {
				Content struct {
					Parts []struct {
						Text string `json:"text"`
					} `json:"parts"`
				} `json:"content"`
			} `json:"candidates"`
		}
		json.NewDecoder(hresp.Body).Decode(&gResp)

		refined := ""
		if len(gResp.Candidates) > 0 && len(gResp.Candidates[0].Content.Parts) > 0 {
			refined = gResp.Candidates[0].Content.Parts[0].Text
		}

		json.NewEncoder(w).Encode(map[string]any{
			"refined_prompt": refined,
		})
	}
}

func imageGenerateHandler(cfg config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodPost {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		var req ImageGenRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
			return
		}

		geminiKey := cfg.getSecret("GEMINI_API_KEY")
		if geminiKey == "" {
			http.Error(w, `{"error":"GEMINI_API_KEY missing"}`, http.StatusServiceUnavailable)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 120*time.Second)
		defer cancel()

		typeInstr := buildAssetTypeInstruction(req.AssetTypes)
		refinePrompt := fmt.Sprintf(`You are an expert AI Image Prompt Engineer.
Refine the following user prompt into a single high-quality, descriptive image-generation prompt.
%s
Return only the refined prompt text, no preamble.
User Prompt: %s`, typeInstr, req.Prompt)

		refineURL := "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash:generateContent?key=" + geminiKey
		refineBody, _ := json.Marshal(map[string]any{
			"contents": []any{map[string]any{"parts": []any{map[string]any{"text": refinePrompt}}}},
		})
		refineReq, _ := http.NewRequestWithContext(ctx, "POST", refineURL, bytes.NewReader(refineBody))
		refineReq.Header.Set("Content-Type", "application/json")
		refineResp, err := http.DefaultClient.Do(refineReq)
		if err != nil || refineResp.StatusCode >= 400 {
			http.Error(w, `{"error":"Gemini refinement failed"}`, http.StatusBadGateway)
			return
		}
		defer refineResp.Body.Close()

		var refineParsed geminiContentResp
		json.NewDecoder(refineResp.Body).Decode(&refineParsed)

		refined := req.Prompt
		if len(refineParsed.Candidates) > 0 && len(refineParsed.Candidates[0].Content.Parts) > 0 {
			if t := strings.TrimSpace(refineParsed.Candidates[0].Content.Parts[0].Text); t != "" {
				refined = t
			}
		}

		const imageModel = "gemini-2.5-flash-image"
		imageURL := "https://generativelanguage.googleapis.com/v1beta/models/" + imageModel + ":generateContent?key=" + geminiKey
		imageBody, _ := json.Marshal(map[string]any{
			"contents": []any{map[string]any{"parts": []any{map[string]any{"text": refined}}}},
			"generationConfig": map[string]any{
				"responseModalities": []string{"IMAGE"},
			},
		})
		imageReq, _ := http.NewRequestWithContext(ctx, "POST", imageURL, bytes.NewReader(imageBody))
		imageReq.Header.Set("Content-Type", "application/json")
		imageResp, err := http.DefaultClient.Do(imageReq)
		if err != nil {
			http.Error(w, `{"error":"Gemini image generation failed"}`, http.StatusBadGateway)
			return
		}
		defer imageResp.Body.Close()
		if imageResp.StatusCode >= 400 {
			body, _ := io.ReadAll(imageResp.Body)
			http.Error(w, fmt.Sprintf(`{"error":"Gemini image generation %d: %s"}`, imageResp.StatusCode, strings.TrimSpace(string(body))), http.StatusBadGateway)
			return
		}

		var imgParsed geminiContentResp
		if err := json.NewDecoder(imageResp.Body).Decode(&imgParsed); err != nil {
			http.Error(w, `{"error":"failed to decode image response"}`, http.StatusBadGateway)
			return
		}

		dataURL := ""
		if len(imgParsed.Candidates) > 0 {
			for _, p := range imgParsed.Candidates[0].Content.Parts {
				if p.InlineData.Data != "" {
					mime := p.InlineData.MimeType
					if mime == "" {
						mime = "image/png"
					}
					dataURL = "data:" + mime + ";base64," + p.InlineData.Data
					break
				}
			}
		}
		if dataURL == "" {
			http.Error(w, `{"error":"Gemini returned no image data"}`, http.StatusBadGateway)
			return
		}

		refineCost := float64(refineParsed.UsageMetadata.PromptTokenCount)*0.30/1e6 +
			float64(refineParsed.UsageMetadata.CandidatesTokenCount)*2.50/1e6
		imageCost := float64(imgParsed.UsageMetadata.PromptTokenCount)*0.30/1e6 +
			float64(imgParsed.UsageMetadata.CandidatesTokenCount)*30.0/1e6
		totalCost := refineCost + imageCost

		json.NewEncoder(w).Encode(map[string]any{
			"original_prompt": req.Prompt,
			"refined_prompt":  refined,
			"image_url":       dataURL,
			"module_number":   req.ModuleNumber,
			"video_number":    req.VideoNumber,
			"model":           imageModel,
			"cost_usd":        totalCost,
			"tokens": map[string]any{
				"refine_prompt":     refineParsed.UsageMetadata.PromptTokenCount,
				"refine_candidates": refineParsed.UsageMetadata.CandidatesTokenCount,
				"image_prompt":      imgParsed.UsageMetadata.PromptTokenCount,
				"image_candidates":  imgParsed.UsageMetadata.CandidatesTokenCount,
			},
		})
	}
}

type ImageSaveRequest struct {
	ImageURL     string `json:"image_url"`
	ModuleNumber int    `json:"module_number"`
	VideoNumber  int    `json:"video_number"`
	Prompt       string `json:"prompt"`
}

func imageSaveHandler(cfg config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodPost {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}
		if !isAdminRequest(r) {
			http.Error(w, `{"error":"Unauthorized — sign in as admin to save images."}`, http.StatusUnauthorized)
			return
		}

		var req ImageSaveRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
			return
		}

		if cfg.azureAccountName == "" {
			http.Error(w, `{"error":"Azure Storage not configured"}`, http.StatusServiceUnavailable)
			return
		}

		var imageData []byte
		contentType := "image/png"
		if strings.HasPrefix(req.ImageURL, "data:") {
			meta, b64, found := strings.Cut(strings.TrimPrefix(req.ImageURL, "data:"), ",")
			if !found {
				http.Error(w, `{"error":"invalid data URL"}`, http.StatusBadRequest)
				return
			}
			if mime, _, ok := strings.Cut(meta, ";"); ok && mime != "" {
				contentType = mime
			} else if meta != "" {
				contentType = meta
			}
			decoded, err := base64.StdEncoding.DecodeString(b64)
			if err != nil {
				http.Error(w, `{"error":"failed to decode image data"}`, http.StatusBadRequest)
				return
			}
			imageData = decoded
		} else {
			resp, err := http.Get(req.ImageURL)
			if err != nil {
				http.Error(w, `{"error":"failed to download image"}`, http.StatusBadGateway)
				return
			}
			defer resp.Body.Close()
			imageData, err = io.ReadAll(resp.Body)
			if err != nil {
				http.Error(w, `{"error":"failed to read image data"}`, http.StatusBadGateway)
				return
			}
			if ct := resp.Header.Get("Content-Type"); ct != "" {
				contentType = ct
			}
		}

		ctx := r.Context()
		blobName := fmt.Sprintf("m%d_v%d_%d.png", req.ModuleNumber, req.VideoNumber, time.Now().Unix())
		container := "research-images"

		if err := uploadBlobToAzure(ctx, cfg, container, blobName, contentType, imageData); err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"azure upload failed: %s"}`, err.Error()), http.StatusBadGateway)
			return
		}

		var thumbName string
		if thumbData, terr := generateThumbnail(imageData, 320); terr != nil {
			log.Printf("thumbnail generation failed for %s: %v", blobName, terr)
		} else {
			thumbName = thumbBlobName(blobName)
			if uerr := uploadBlobToAzure(ctx, cfg, container, thumbName, "image/jpeg", thumbData); uerr != nil {
				log.Printf("thumbnail upload failed for %s: %v", thumbName, uerr)
				thumbName = ""
			}
		}

		imageURL := researchFileProxyURL(container, blobName)
		var thumbURL string
		if thumbName != "" {
			thumbURL = researchFileProxyURL(container, thumbName)
		}

		dbEntry := map[string]any{
			"module_number":       req.ModuleNumber,
			"video_number":        req.VideoNumber,
			"prompt":              req.Prompt,
			"azure_blob_name":     blobName,
			"image_url":           imageURL,
			"thumbnail_blob_name": thumbName,
			"thumbnail_url":       thumbURL,
			"status":              "saved_to_azure",
		}
		if err := supabasePost(ctx, cfg, "generated_images", dbEntry); err != nil {
			log.Printf("supabase save generated image error: %v", err)
		}

		json.NewEncoder(w).Encode(map[string]any{
			"ok":            true,
			"blob_name":     blobName,
			"url":           imageURL,
			"thumbnail":     thumbName,
			"thumbnail_url": thumbURL,
		})
	}
}

type DrawingSaveRequest struct {
	SentenceID     int             `json:"sentence_id"`
	ModuleNumber   int             `json:"module_number"`
	VideoNumber    int             `json:"video_number"`
	ExcalidrawJSON json.RawMessage `json:"excalidraw_json"`
	PNG            string          `json:"png"`
	Prompt         string          `json:"prompt"`
}

func drawingSaveHandler(cfg config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodPost {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}
		if !isAdminRequest(r) {
			http.Error(w, `{"error":"Unauthorized — sign in as admin to save drawings."}`, http.StatusUnauthorized)
			return
		}

		var req DrawingSaveRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
			return
		}
		if req.SentenceID == 0 {
			http.Error(w, `{"error":"sentence_id is required"}`, http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		var blobName, imageURL string

		if strings.HasPrefix(req.PNG, "data:") {
			if cfg.azureAccountName == "" {
				http.Error(w, `{"error":"Azure Storage not configured"}`, http.StatusServiceUnavailable)
				return
			}
			meta, b64, found := strings.Cut(strings.TrimPrefix(req.PNG, "data:"), ",")
			if !found {
				http.Error(w, `{"error":"invalid data URL"}`, http.StatusBadRequest)
				return
			}
			contentType := "image/png"
			if mime, _, ok := strings.Cut(meta, ";"); ok && mime != "" {
				contentType = mime
			}
			imageData, err := base64.StdEncoding.DecodeString(b64)
			if err != nil {
				http.Error(w, `{"error":"failed to decode image data"}`, http.StatusBadRequest)
				return
			}
			blobName = fmt.Sprintf("drawing_s%d_%d.png", req.SentenceID, time.Now().Unix())
			if err := uploadBlobToAzure(ctx, cfg, "research-images", blobName, contentType, imageData); err != nil {
				http.Error(w, fmt.Sprintf(`{"error":"azure upload failed: %s"}`, err.Error()), http.StatusBadGateway)
				return
			}
			imageURL = researchFileProxyURL("research-images", blobName)
		}

		row := map[string]any{
			"sentence_id":   req.SentenceID,
			"module_number": req.ModuleNumber,
			"video_number":  req.VideoNumber,
			"prompt_used":   req.Prompt,
			"updated_at":    time.Now().UTC().Format(time.RFC3339),
		}
		if len(req.ExcalidrawJSON) > 0 {
			row["excalidraw_json"] = req.ExcalidrawJSON
		}
		if blobName != "" {
			row["azure_blob_name"] = blobName
			row["image_url"] = imageURL
		}
		if err := supabaseUpsert(ctx, cfg, "sentence_drawings", row, "sentence_id"); err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"supabase save failed: %s"}`, err.Error()), http.StatusBadGateway)
			return
		}

		json.NewEncoder(w).Encode(map[string]any{
			"ok":        true,
			"blob_name": blobName,
			"image_url": imageURL,
		})
	}
}

type generatedImageRow struct {
	ID            int64  `json:"id"`
	AzureBlobName string `json:"azure_blob_name"`
	ImageURL      string `json:"image_url"`
	ThumbnailURL  string `json:"thumbnail_url"`
}

func imageThumbnailBackfillHandler(cfg config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodPost {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}
		if cfg.azureAccountName == "" {
			http.Error(w, `{"error":"Azure Storage not configured"}`, http.StatusServiceUnavailable)
			return
		}
		ctx := r.Context()
		container := "research-images"

		var rows []generatedImageRow
		q := "select=id,azure_blob_name,image_url,thumbnail_url&thumbnail_url=is.null&azure_blob_name=not.is.null&limit=500"
		if err := supabaseGet(ctx, cfg, "generated_images", q, &rows); err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"supabase read: %s"}`, err.Error()), http.StatusBadGateway)
			return
		}

		processed, failed := 0, 0
		var errs []string
		for _, row := range rows {
			if row.AzureBlobName == "" {
				continue
			}
			expiry := time.Now().UTC().Add(10 * time.Minute)
			sasQuery, err := generateContainerSAS(cfg.azureAccountName, cfg.azureAccountKey, container, "r", expiry)
			if err != nil {
				failed++
				errs = append(errs, fmt.Sprintf("%d: sas %v", row.ID, err))
				continue
			}
			getURL := blobURL(cfg.azureAccountName, container, row.AzureBlobName, sasQuery)
			dresp, err := http.Get(getURL)
			if err != nil || dresp.StatusCode >= 400 {
				failed++
				code := 0
				if dresp != nil {
					code = dresp.StatusCode
					dresp.Body.Close()
				}
				errs = append(errs, fmt.Sprintf("%d: download (%d) %v", row.ID, code, err))
				continue
			}
			orig, rerr := io.ReadAll(dresp.Body)
			dresp.Body.Close()
			if rerr != nil {
				failed++
				errs = append(errs, fmt.Sprintf("%d: read %v", row.ID, rerr))
				continue
			}
			thumbData, terr := generateThumbnail(orig, 320)
			if terr != nil {
				failed++
				errs = append(errs, fmt.Sprintf("%d: thumb %v", row.ID, terr))
				continue
			}
			tName := thumbBlobName(row.AzureBlobName)
			if uerr := uploadBlobToAzure(ctx, cfg, container, tName, "image/jpeg", thumbData); uerr != nil {
				failed++
				errs = append(errs, fmt.Sprintf("%d: upload %v", row.ID, uerr))
				continue
			}
			patch := map[string]any{
				"thumbnail_blob_name": tName,
				"thumbnail_url":       researchFileProxyURL(container, tName),
			}
			if perr := supabasePatch(ctx, cfg, "generated_images", fmt.Sprintf("id=eq.%d", row.ID), patch); perr != nil {
				failed++
				errs = append(errs, fmt.Sprintf("%d: patch %v", row.ID, perr))
				continue
			}
			processed++
		}

		if len(errs) > 5 {
			errs = errs[:5]
		}
		json.NewEncoder(w).Encode(map[string]any{
			"ok":         true,
			"candidates": len(rows),
			"processed":  processed,
			"failed":     failed,
			"errors":     errs,
		})
	}
}
