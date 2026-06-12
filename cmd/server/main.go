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
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// ── Config ────────────────────────────────────────────────────────────────────

type config struct {
	supabaseURL        string
	supabaseAnon       string
	axiomToken         string
	axiomDataset       string
	axiomAPIURL        string
	axiomQueryURL      string
	port               string
	azureAccountName   string
	azureAccountKey    string
	azureKeyVaultName  string
	azureTenantID      string
	azureClientID      string
	azureClientSecret  string
}

func loadConfig() config {
	cfg := config{
		supabaseURL:       mustEnv("SUPABASE_URL"),
		supabaseAnon:      mustEnv("SUPABASE_ANON_KEY"),
		axiomToken:        os.Getenv("AXIOM_TOKEN"),
		axiomDataset:      mustEnv("AXIOM_DATASET"),
		axiomAPIURL:       envOr("AXIOM_API_URL", "https://api.axiom.co"),
		axiomQueryURL:     envOr("AXIOM_QUERY_URL", "https://api.axiom.co"),
		port:              envOr("PORT", "8080"),
		azureKeyVaultName: os.Getenv("AZURE_KEYVAULT_NAME"),
		azureTenantID:     os.Getenv("AZURE_TENANT_ID"),
		azureClientID:     os.Getenv("AZURE_CLIENT_ID"),
		azureClientSecret: os.Getenv("AZURE_CLIENT_SECRET"),
	}
	if connStr := os.Getenv("AZURE_STORAGE_CONN_STR"); connStr != "" {
		cfg.azureAccountName, cfg.azureAccountKey = parseStorageConnStr(connStr)
	}
	return cfg
}

func getSecretFromKeyVault(vaultName, tenantID, clientID, clientSecret, secretName string) (string, error) {
	secretName = strings.ReplaceAll(strings.ToLower(secretName), "_", "-")

	tokenURL := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", tenantID)
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("scope", "https://vault.azure.net/.default")

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("oauth token error (%d): %s", resp.StatusCode, b)
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", err
	}

	secretURL := fmt.Sprintf("https://%s.vault.azure.net/secrets/%s?api-version=7.4", vaultName, secretName)
	reqSecret, err := http.NewRequest("GET", secretURL, nil)
	if err != nil {
		return "", err
	}
	reqSecret.Header.Set("Authorization", "Bearer "+tokenResp.AccessToken)

	respSecret, err := http.DefaultClient.Do(reqSecret)
	if err != nil {
		return "", err
	}
	defer respSecret.Body.Close()

	if respSecret.StatusCode >= 400 {
		b, _ := io.ReadAll(respSecret.Body)
		return "", fmt.Errorf("keyvault get secret error (%d): %s", respSecret.StatusCode, b)
	}

	var secretResp struct {
		Value string `json:"value"`
	}
	if err := json.NewDecoder(respSecret.Body).Decode(&secretResp); err != nil {
		return "", err
	}

	return secretResp.Value, nil
}

func (c config) getSecret(secretName string) string {
	if c.azureKeyVaultName != "" && c.azureTenantID != "" && c.azureClientID != "" && c.azureClientSecret != "" {
		val, err := getSecretFromKeyVault(c.azureKeyVaultName, c.azureTenantID, c.azureClientID, c.azureClientSecret, secretName)
		if err == nil && val != "" {
			log.Printf("Successfully loaded secret '%s' from Key Vault '%s'", secretName, c.azureKeyVaultName)
			return val
		}
		log.Printf("Keyvault getSecret failed for %s, falling back to env: %v", secretName, err)
	}
	return os.Getenv(secretName)
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("required env var %s is not set", key)
	}
	return v
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// ── Axiom ─────────────────────────────────────────────────────────────────────

func shipToAxiom(cfg config, events []map[string]any) {
	if cfg.axiomToken == "" {
		return
	}
	body, err := json.Marshal(events)
	if err != nil {
		log.Printf("axiom marshal: %v", err)
		return
	}
	var url string
	if strings.Contains(cfg.axiomAPIURL, ".edge.axiom.co") {
		url = fmt.Sprintf("%s/v1/ingest/%s", cfg.axiomAPIURL, cfg.axiomDataset)
	} else {
		url = fmt.Sprintf("%s/v1/datasets/%s/ingest", cfg.axiomAPIURL, cfg.axiomDataset)
	}
	log.Printf("axiom ingest -> %s (%d events)", url, len(events))
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		log.Printf("axiom request build: %v", err)
		return
	}
	req.Header.Set("Authorization", "Bearer "+cfg.axiomToken)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("axiom send: %v", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		log.Printf("axiom ingest %d: %s", resp.StatusCode, b)
	}
}

// ── observe middleware ────────────────────────────────────────────────────────

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (sw *statusWriter) WriteHeader(code int) {
	sw.status = code
	sw.ResponseWriter.WriteHeader(code)
}

func observe(cfg config, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sw := &statusWriter{ResponseWriter: w, status: http.StatusOK}

		defer func() {
			if rec := recover(); rec != nil {
				http.Error(sw, "internal server error", http.StatusInternalServerError)
				go shipToAxiom(cfg, []map[string]any{{
					"_time":  time.Now().UTC().Format(time.RFC3339),
					"level":  "error",
					"method": r.Method,
					"path":   r.URL.Path,
					"panic":  fmt.Sprintf("%v", rec),
				}})
				log.Printf("PANIC %s %s: %v", r.Method, r.URL.Path, rec)
			}
		}()

		next.ServeHTTP(sw, r)

		dur := time.Since(start)
		go shipToAxiom(cfg, []map[string]any{{
			"_time":       time.Now().UTC().Format(time.RFC3339),
			"level":       "info",
			"method":      r.Method,
			"path":        r.URL.Path,
			"status":      sw.status,
			"duration_ms": dur.Milliseconds(),
		}})
		log.Printf("%s %s %d %dms", r.Method, r.URL.Path, sw.status, dur.Milliseconds())
	})
}

// ── Supabase models ───────────────────────────────────────────────────────────

type Course struct {
	ID                    string   `json:"id"`
	CourseTitle           string   `json:"course_title"`
	Instructor            string   `json:"instructor"`
	TargetAudience        string   `json:"target_audience"`
	TotalDuration         string   `json:"total_duration"`
	DifficultyLevel       string   `json:"difficulty_level"`
	LearningObjectives    []string `json:"learning_objectives"`
	KeyTakeaways          string   `json:"key_takeaways"`
	RealWorldConnections  string   `json:"real_world_connections"`
	RecommendedBackground string   `json:"recommended_background"`
	ProofOfLearning       string   `json:"proof_of_learning"`
	CourseDescription     string   `json:"course_description"`
	Skills                []string `json:"skills"`
}

type Tool struct {
	ToolName       string `json:"tool_name"`
	Purpose        string `json:"purpose"`
	FreeOrPaid     string `json:"free_or_paid"`
	TrialAvailable bool   `json:"trial_available"`
	TrialDuration  string `json:"trial_duration"`
	ToolValidity   string `json:"tool_validity"`
}

// ── Supabase helpers ──────────────────────────────────────────────────────────

func supabaseReq(ctx context.Context, cfg config, method, table, query string, body any) (*http.Response, error) {
	u := fmt.Sprintf("%s/rest/v1/%s", cfg.supabaseURL, table)
	if query != "" {
		u += "?" + query
	}
	var bodyReader io.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		bodyReader = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, method, u, bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("apikey", cfg.supabaseAnon)
	req.Header.Set("Authorization", "Bearer "+cfg.supabaseAnon)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return http.DefaultClient.Do(req)
}

func supabaseGet(ctx context.Context, cfg config, table, query string, out any) error {
	resp, err := supabaseReq(ctx, cfg, http.MethodGet, table, query, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("supabase GET %s %s: %s", table, resp.Status, b)
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

func supabasePost(ctx context.Context, cfg config, table string, body any) error {
	resp, err := supabaseReq(ctx, cfg, http.MethodPost, table, "select=", body)
	if err != nil {
		return err
	}
	resp.Body.Close()
	if resp.StatusCode >= 400 {
		return fmt.Errorf("supabase POST %s: %s", table, resp.Status)
	}
	return nil
}

func supabaseDelete(ctx context.Context, cfg config, table, query string) error {
	resp, err := supabaseReq(ctx, cfg, http.MethodDelete, table, query, nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	if resp.StatusCode >= 400 {
		return fmt.Errorf("supabase DELETE %s: %s", table, resp.Status)
	}
	return nil
}

// ── Models ────────────────────────────────────────────────────────────────────

type NavFav struct {
	URL   string `json:"url"`
	Label string `json:"label"`
}

// ── Template page data ────────────────────────────────────────────────────────

type indexData struct {
	Course        *Course
	Tools         []Tool
	FetchErr      string
	NavFavsJSON   template.JS // pre-serialised for window.__NAV_FAVS__
	NavConfigJSON template.JS // pre-serialised for window.__NAV_CONFIG__
}

// ── Handlers ──────────────────────────────────────────────────────────────────

func homeHandler(tmpl *template.Template, cfg config, navConfigJS template.JS) http.HandlerFunc {
	static := http.FileServer(http.Dir("."))
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" && r.URL.Path != "/index.html" {
			static.ServeHTTP(w, r) // serve unported routes straight from disk
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		data := indexData{}

		var courses []Course
		if err := supabaseGet(ctx, cfg, "course_metadata", "select=*&limit=1", &courses); err != nil {
			data.FetchErr = err.Error()
			go shipToAxiom(cfg, []map[string]any{{
				"_time": time.Now().UTC().Format(time.RFC3339),
				"level": "error", "path": "/", "err": err.Error(),
			}})
		} else if len(courses) > 0 {
			c := courses[0]
			data.Course = &c

			var tools []Tool
			q := fmt.Sprintf("select=*&course_id=eq.%s&order=display_order.asc", c.ID)
			if err := supabaseGet(ctx, cfg, "course_tools", q, &tools); err != nil {
				log.Printf("course_tools fetch: %v", err)
			} else {
				data.Tools = tools
			}
		}

		var favs []NavFav
		if err := supabaseGet(ctx, cfg, "nav_favorites", "select=url,label&order=created_at.asc", &favs); err != nil {
			log.Printf("nav_favorites fetch: %v", err)
		}
		favsJSON, _ := json.Marshal(favs)
		data.NavFavsJSON = template.JS(favsJSON)
		data.NavConfigJSON = navConfigJS

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tmpl.ExecuteTemplate(w, "index.html", data); err != nil {
			log.Printf("template execute: %v", err)
		}
	}
}

// ── Axiom errors admin ───────────────────────────────────────────────────────

type axiomEvent struct {
	Time     string `json:"_time"`
	Level    string `json:"level"`
	Method   string `json:"method"`
	Path     string `json:"path"`
	Status   int    `json:"status"`
	Duration int64  `json:"duration_ms"`
	Err      string `json:"err"`
	Panic    string `json:"panic"`
}

type axiomQueryResp struct {
	Matches []struct {
		Time string         `json:"_time"`
		Data map[string]any `json:"data"`
	} `json:"matches"`
}

type axiomErrorsData struct {
	Events        []axiomEvent
	FetchErr      string
	QueryURL      string
	APL           string
	NavFavsJSON   template.JS
	NavConfigJSON template.JS
}

func axiomErrorsHandler(tmpl *template.Template, cfg config, navConfigJS template.JS) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := axiomErrorsData{NavConfigJSON: navConfigJS}

		var favs []NavFav
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()
		if err := supabaseGet(ctx, cfg, "nav_favorites", "select=url,label&order=created_at.asc", &favs); err != nil {
			log.Printf("nav_favorites fetch: %v", err)
		}
		favsJSON, _ := json.Marshal(favs)
		data.NavFavsJSON = template.JS(favsJSON)

		if cfg.axiomToken == "" {
			data.FetchErr = "AXIOM_TOKEN not set — configure it via Fly.io secrets"
		} else {
			apl := fmt.Sprintf(`['%s'] | sort by _time desc | limit 100`, cfg.axiomDataset)
			data.APL = apl
			data.QueryURL = fmt.Sprintf("%s/v1/datasets/%s/query", cfg.axiomQueryURL, cfg.axiomDataset)
			body, _ := json.Marshal(map[string]any{
				"apl":       apl,
				"startTime": "now-24h",
			})
			queryURL := fmt.Sprintf("%s/v1/datasets/%s/query", cfg.axiomQueryURL, cfg.axiomDataset)
			log.Printf("axiom query -> %s (apl: %s)", queryURL, apl)

			req, err := http.NewRequestWithContext(ctx, http.MethodPost, queryURL, bytes.NewReader(body))
			if err != nil {
				data.FetchErr = err.Error()
			} else {
				req.Header.Set("Authorization", "Bearer "+cfg.axiomToken)
				req.Header.Set("Content-Type", "application/json")
				resp, err := http.DefaultClient.Do(req)
				if err != nil {
					data.FetchErr = err.Error()
				} else {
					defer resp.Body.Close()
					respBody, _ := io.ReadAll(resp.Body)
					log.Printf("axiom query response %d: %s", resp.StatusCode, string(respBody[:min(len(respBody), 1000)]))

					if resp.StatusCode >= 300 {
						data.FetchErr = fmt.Sprintf("Axiom API %d: %s", resp.StatusCode, string(respBody))
					} else {
						var qr axiomQueryResp
						if err := json.Unmarshal(respBody, &qr); err != nil {
							data.FetchErr = "decode: " + err.Error()
						} else {
							for _, m := range qr.Matches {
								ev := axiomEvent{Time: m.Time}
								if v, ok := m.Data["level"].(string); ok { ev.Level = v }
								if v, ok := m.Data["method"].(string); ok { ev.Method = v }
								if v, ok := m.Data["path"].(string); ok { ev.Path = v }
								if v, ok := m.Data["status"].(float64); ok { ev.Status = int(v) }
								if v, ok := m.Data["duration_ms"].(float64); ok { ev.Duration = int64(v) }
								if v, ok := m.Data["err"].(string); ok { ev.Err = v }
								if v, ok := m.Data["panic"].(string); ok { ev.Panic = v }
								data.Events = append(data.Events, ev)
							}
							log.Printf("axiom query returned %d events", len(data.Events))
						}
					}
				}
			}
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tmpl.ExecuteTemplate(w, "axiom_errors.html", data); err != nil {
			log.Printf("axiom_errors template: %v", err)
		}
	}
}

// ── Public config endpoint ────────────────────────────────────────────────────
// Returns the Supabase anon key (public by design) so static pages can
// auto-connect without hardcoding it in HTML.

func configHandler(cfg config) http.HandlerFunc {
	type configResp struct {
		SupabaseURL  string `json:"supabaseUrl"`
		SupabaseAnon string `json:"supabaseAnon"`
	}
	payload, _ := json.Marshal(configResp{
		SupabaseURL:  cfg.supabaseURL,
		SupabaseAnon: cfg.supabaseAnon,
	})
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "public, max-age=300")
		w.Write(payload)
	}
}

// ── Client-side error ingestion ───────────────────────────────────────────────

func clientErrorsHandler(cfg config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		payload["_time"] = time.Now().UTC().Format(time.RFC3339)
		if _, ok := payload["level"]; !ok {
			payload["level"] = "error"
		}
		payload["source"] = "client"
		go shipToAxiom(cfg, []map[string]any{payload})
		w.WriteHeader(http.StatusNoContent)
	}
}

// ── Nav favourites toggle ─────────────────────────────────────────────────────

func navFavsHandler(cfg config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req NavFav
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.URL == "" {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		// Check if already stored
		var existing []NavFav
		q := "url=eq." + url.QueryEscape(req.URL) + "&select=url"
		_ = supabaseGet(ctx, cfg, "nav_favorites", q, &existing)

		w.Header().Set("Content-Type", "application/json")
		if len(existing) > 0 {
			_ = supabaseDelete(ctx, cfg, "nav_favorites", "url=eq."+url.QueryEscape(req.URL))
			fmt.Fprint(w, `{"favorited":false}`)
		} else {
			_ = supabasePost(ctx, cfg, "nav_favorites", req)
			fmt.Fprint(w, `{"favorited":true}`)
		}
	}
}

// ── Azure Blob Storage ────────────────────────────────────────────────────────

func parseStorageConnStr(connStr string) (accountName, accountKey string) {
	for _, part := range strings.Split(connStr, ";") {
		k, v, ok := strings.Cut(part, "=")
		if !ok {
			continue
		}
		switch k {
		case "AccountName":
			accountName = v
		case "AccountKey":
			// strings.Cut splits on the first "=" so base64 padding in the key is preserved
			accountKey = v
		}
	}
	return
}

// sasPerms puts permission chars in Azure's required canonical order: r a c w d x l t f m e o p i y q
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

// generateContainerSAS creates a service SAS query string for a blob container.
// Uses 16-field string-to-sign (sv=2026-04-06, includes signedEncryptionScope as field 11).
// Format and permissions order verified against az storage container generate-sas output.
func generateContainerSAS(accountName, accountKey, container, permissions string, expiry time.Time) (string, error) {
	permissions = sasPerms(permissions) // enforce canonical order; Azure rejects sp= if out of order
	const version = "2026-04-06"
	expiryStr := expiry.UTC().Format("2006-01-02T15:04:05Z")
	canonResource := "/blob/" + accountName + "/" + container
	// 16 fields: perm, start, expiry, canon, id, ip, proto, ver, resource(c),
	//            snapshotTime, encryptionScope, rscc, rscd, rsce, rscl, rsct
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
	"research-images": true,
	"research-audio":  true,
	"research-videos": true,
	"research-notes":  true,
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
	return base + "?" + sasQuery
}

func researchUploadHandler(cfg config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
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

		expiry := time.Now().UTC().Add(5 * time.Minute)
		sasQuery, err := generateContainerSAS(cfg.azureAccountName, cfg.azureAccountKey, container, "rcwl", expiry)
		if err != nil {
			http.Error(w, "sas error", http.StatusInternalServerError)
			return
		}
		putURL := blobURL(cfg.azureAccountName, container, blobName, sasQuery)
		req, err := http.NewRequestWithContext(r.Context(), http.MethodPut, putURL, file)
		if err != nil {
			http.Error(w, "request build", http.StatusInternalServerError)
			return
		}
		req.Header.Set("x-ms-blob-type", "BlockBlob")
		req.Header.Set("Content-Type", contentType)
		req.ContentLength = header.Size

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			http.Error(w, "azure upload: "+err.Error(), http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 400 {
			b, _ := io.ReadAll(resp.Body)
			http.Error(w, fmt.Sprintf("azure %d: %s", resp.StatusCode, b), http.StatusBadGateway)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"ok": "true", "name": blobName})
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

// ── Explanations ──────────────────────────────────────────────────────────────

type Explanation struct {
	ID              int       `json:"id,omitempty"`
	EntityType      string    `json:"entity_type"`
	EntityID        string    `json:"entity_id"`
	ExplanationText string    `json:"explanation_text"`
	GeneratedBy     string    `json:"generated_by"`
	CreatedAt       string    `json:"created_at,omitempty"`
}

func fetchEntityContent(ctx context.Context, cfg config, entityType, entityID string) (string, error) {
	switch entityType {
	case "sentence":
		var items []struct {
			SentenceText string `json:"sentence_text"`
			SentenceType string `json:"sentence_type"`
		}
		if err := supabaseGet(ctx, cfg, "sentences", "select=sentence_text,sentence_type&id=eq."+entityID, &items); err == nil && len(items) > 0 {
			return fmt.Sprintf("Sentence [%s]: %s", items[0].SentenceType, items[0].SentenceText), nil
		}
	case "outline":
		var items []struct {
			Content     string `json:"content"`
			ContentType string `json:"content_type"`
		}
		if err := supabaseGet(ctx, cfg, "outline", "select=content,content_type&id=eq."+entityID, &items); err == nil && len(items) > 0 {
			return fmt.Sprintf("Outline Node [%s]: %s", items[0].ContentType, items[0].Content), nil
		}
	case "research":
		return fmt.Sprintf("Research Asset: %s", entityID), nil
	case "problem":
		var items []struct {
			Title    string `json:"title"`
			Headline string `json:"headline"`
		}
		if err := supabaseGet(ctx, cfg, "problem_pages", "select=title,headline&id=eq."+entityID, &items); err == nil && len(items) > 0 {
			return fmt.Sprintf("Problem Page: %s - Headline: %s", items[0].Title, items[0].Headline), nil
		}
	}
	return "Entity ID: " + entityID, nil
}

func generateExplanationHandler(cfg config) http.HandlerFunc {
	type GenRequest struct {
		EntityType string `json:"entity_type"`
		EntityID   string `json:"entity_id"`
		Prompt     string `json:"prompt"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodPost {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		var req GenRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
			return
		}

		if req.EntityType == "" || req.EntityID == "" {
			http.Error(w, `{"error":"entity_type and entity_id are required"}`, http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
		defer cancel()

		entityContent, err := fetchEntityContent(ctx, cfg, req.EntityType, req.EntityID)
		if err != nil {
			log.Printf("failed to fetch entity details: %v", err)
		}

		geminiKey := cfg.getSecret("GEMINI_API_KEY")
		if geminiKey == "" {
			http.Error(w, `{"error":"GEMINI_API_KEY is not configured in Azure Key Vault or environment"}`, http.StatusServiceUnavailable)
			return
		}

		prompt := fmt.Sprintf("You are an expert AI Cloud & Software Architect. Explain this entity in detail for the Claude AI Certification program:\n\n%s", entityContent)
		if req.Prompt != "" {
			prompt = fmt.Sprintf("%s\n\nUser request: %s", prompt, req.Prompt)
		}

		// Prepare Gemini API request
		geminiURL := "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent?key=" + geminiKey
		geminiReqBody := map[string]any{
			"contents": []any{
				map[string]any{
					"parts": []any{
						map[string]any{
							"text": prompt,
						},
					},
				},
			},
		}

		geminiReqBytes, err := json.Marshal(geminiReqBody)
		if err != nil {
			http.Error(w, `{"error":"failed to marshal request"}`, http.StatusInternalServerError)
			return
		}

		hreq, err := http.NewRequestWithContext(ctx, "POST", geminiURL, bytes.NewReader(geminiReqBytes))
		if err != nil {
			http.Error(w, `{"error":"failed to build request"}`, http.StatusInternalServerError)
			return
		}
		hreq.Header.Set("Content-Type", "application/json")

		hresp, err := http.DefaultClient.Do(hreq)
		if err != nil {
			http.Error(w, `{"error":"Gemini API call failed: `+err.Error()+`"}`, http.StatusBadGateway)
			return
		}
		defer hresp.Body.Close()

		if hresp.StatusCode >= 400 {
			b, _ := io.ReadAll(hresp.Body)
			http.Error(w, fmt.Sprintf(`{"error":"Gemini API returned error (%d): %s"}`, hresp.StatusCode, string(b)), http.StatusBadGateway)
			return
		}

		var geminiResp struct {
			Candidates []struct {
				Content struct {
					Parts []struct {
						Text string `json:"text"`
					} `json:"parts"`
				} `json:"content"`
			} `json:"candidates"`
		}

		if err := json.NewDecoder(hresp.Body).Decode(&geminiResp); err != nil {
			http.Error(w, `{"error":"failed to decode Gemini response"}`, http.StatusInternalServerError)
			return
		}

		if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
			http.Error(w, `{"error":"no response candidate from Gemini"}`, http.StatusInternalServerError)
			return
		}

		responseText := geminiResp.Candidates[0].Content.Parts[0].Text

		json.NewEncoder(w).Encode(map[string]string{
			"explanation": responseText,
		})
	}
}

func explanationsHandler(cfg config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		if r.Method == http.MethodGet {
			entityType := r.URL.Query().Get("entity_type")
			entityID := r.URL.Query().Get("entity_id")

			var query string
			if entityType != "" && entityID != "" {
				query = fmt.Sprintf("entity_type=eq.%s&entity_id=eq.%s&order=created_at.desc", url.QueryEscape(entityType), url.QueryEscape(entityID))
			} else if entityType != "" {
				query = fmt.Sprintf("entity_type=eq.%s&order=created_at.desc", url.QueryEscape(entityType))
			} else {
				query = "order=created_at.desc"
			}

			var exps []Explanation
			if err := supabaseGet(ctx, cfg, "explanations", query, &exps); err != nil {
				http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(exps)
			return
		}

		if r.Method == http.MethodPost {
			var exp Explanation
			if err := json.NewDecoder(r.Body).Decode(&exp); err != nil {
				http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
				return
			}
			if exp.EntityType == "" || exp.EntityID == "" || exp.ExplanationText == "" {
				http.Error(w, `{"error":"entity_type, entity_id, and explanation_text are required"}`, http.StatusBadRequest)
				return
			}
			if exp.GeneratedBy == "" {
				exp.GeneratedBy = "user"
			}

			if err := supabasePost(ctx, cfg, "explanations", exp); err != nil {
				http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
				return
			}

			json.NewEncoder(w).Encode(map[string]string{"ok": "true"})
			return
		}

		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

// ── Main ──────────────────────────────────────────────────────────────────────

func main() {
	cfg := loadConfig()

	// Read nav config once at startup so the template never needs a client-side fallback.
	navConfigRaw, err := os.ReadFile("navigation_config.json")
	if err != nil {
		log.Fatalf("cannot read navigation_config.json: %v", err)
	}
	navConfigJS := template.JS(navConfigRaw)

	funcs := template.FuncMap{
		"inc": func(i int) int { return i + 1 },
	}
	tmpl := template.Must(
		template.New("").Funcs(funcs).ParseFiles(
			"5_Symbols/templates/index.html",
			"5_Symbols/templates/axiom_errors.html",
		),
	)

	fs := http.FileServer(http.Dir("."))
	mux := http.NewServeMux()
	mux.Handle("/shared/", observe(cfg, fs))
	mux.Handle("/navigation_config.json", observe(cfg, fs))
	mux.Handle("/api/config", observe(cfg, configHandler(cfg)))
	mux.Handle("/api/nav/favs", observe(cfg, navFavsHandler(cfg)))
	mux.Handle("/api/errors", observe(cfg, clientErrorsHandler(cfg)))
	mux.Handle("/api/research/upload", observe(cfg, researchUploadHandler(cfg)))
	mux.Handle("/api/research/files", observe(cfg, researchFilesHandler(cfg)))
	mux.Handle("/api/research/file", observe(cfg, researchFileHandler(cfg)))
	mux.Handle("/api/explanations/generate", observe(cfg, generateExplanationHandler(cfg)))
	mux.Handle("/api/explanations", observe(cfg, explanationsHandler(cfg)))
	mux.Handle("/admin/errors", observe(cfg, axiomErrorsHandler(tmpl, cfg, navConfigJS)))
	mux.Handle("/", observe(cfg, homeHandler(tmpl, cfg, navConfigJS)))

	addr := ":" + cfg.port
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
