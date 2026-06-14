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

func loadDotEnv(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		val = strings.Trim(val, "'\"")
		if os.Getenv(key) == "" {
			os.Setenv(key, val)
		}
	}
}

func loadConfig() config {
	loadDotEnv(".env")
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

		// Cross-origin support for every route: the static GitHub Pages site
		// (display-only) diverts all /api calls to this Fly.io backend. Allow
		// *.github.io and local dev, and answer CORS preflights here.
		setCORS(sw, r)
		if r.Method == http.MethodOptions {
			sw.WriteHeader(http.StatusNoContent)
			return
		}

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
		SupabaseURL      string `json:"supabaseUrl"`
		SupabaseAnon     string `json:"supabaseAnon"`
		AzureAccountName string `json:"azureAccountName"`
	}
	payload, _ := json.Marshal(configResp{
		SupabaseURL:      cfg.supabaseURL,
		SupabaseAnon:     cfg.supabaseAnon,
		AzureAccountName: cfg.azureAccountName,
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
		geminiURL := "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash:generateContent?key=" + geminiKey
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

// ── Image Generation ─────────────────────────────────────────────────────────

type ImageGenRequest struct {
	Prompt       string   `json:"prompt"`
	ModuleNumber int      `json:"module_number"`
	VideoNumber  int      `json:"video_number"`
	AssetTypes   []string `json:"asset_types"`
}

// geminiContentResp covers both text and inline-image (base64) responses from
// the generateContent endpoint, plus token usage for cost calculation.
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

		// ── Step 1: Refine the user prompt with the text model ──
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

		// ── Step 2: Generate the actual image with the Gemini image model ──
		// gemini-2.5-flash-image ("Nano Banana") is the same model the Gemini
		// web app uses; it returns the image as inline base64 data.
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

		// Extract the inline image data (base64) → build a data URL.
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

		// ── Step 3: Cost — refinement (gemini-2.5-flash) + image (gemini-2.5-flash-image) ──
		// Pricing per 1M tokens (USD): flash text in $0.30 / out $2.50;
		// flash-image text in $0.30 / image out $30.00 (1290 tokens per image ≈ $0.039).
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

		var req ImageSaveRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
			return
		}

		if cfg.azureAccountName == "" {
			http.Error(w, `{"error":"Azure Storage not configured"}`, http.StatusServiceUnavailable)
			return
		}

		// Obtain the image bytes — either from an inline data URL (Gemini
		// returns base64) or by downloading a remote URL.
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

		blobName := fmt.Sprintf("m%d_v%d_%d.png", req.ModuleNumber, req.VideoNumber, time.Now().Unix())
		container := "research-images"

		expiry := time.Now().UTC().Add(10 * time.Minute)
		sasQuery, err := generateContainerSAS(cfg.azureAccountName, cfg.azureAccountKey, container, "rcwl", expiry)
		if err != nil {
			http.Error(w, `{"error":"sas error"}`, http.StatusInternalServerError)
			return
		}

		putURL := blobURL(cfg.azureAccountName, container, blobName, sasQuery)
		ureq, _ := http.NewRequest(http.MethodPut, putURL, bytes.NewReader(imageData))
		ureq.Header.Set("x-ms-blob-type", "BlockBlob")
		ureq.Header.Set("Content-Type", contentType)
		ureq.ContentLength = int64(len(imageData))

		uresp, err := http.DefaultClient.Do(ureq)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"azure upload failed: %s"}`, err.Error()), http.StatusBadGateway)
			return
		}
		defer uresp.Body.Close()
		if uresp.StatusCode >= 400 {
			b, _ := io.ReadAll(uresp.Body)
			http.Error(w, fmt.Sprintf(`{"error":"azure upload %d: %s"}`, uresp.StatusCode, strings.TrimSpace(string(b))), http.StatusBadGateway)
			return
		}

		// Save to Supabase
		ctx := r.Context()
		dbEntry := map[string]any{
			"module_number":   req.ModuleNumber,
			"video_number":    req.VideoNumber,
			"prompt":          req.Prompt,
			"azure_blob_name": blobName,
			"status":          "saved_to_azure",
			"image_url":       blobURL(cfg.azureAccountName, container, blobName, ""), // URL without SAS for DB
		}
		
		if err := supabasePost(ctx, cfg, "generated_images", dbEntry); err != nil {
			log.Printf("supabase save generated image error: %v", err)
		}

		json.NewEncoder(w).Encode(map[string]any{
			"ok":        true,
			"blob_name": blobName,
			"url":       blobURL(cfg.azureAccountName, container, blobName, ""),
		})
	}
}

// ── Gemini Connection Test ────────────────────────────────────────────────────

func imageTestGeminiHandler(cfg config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodGet {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		geminiKey := cfg.getSecret("GEMINI_API_KEY")
		if geminiKey == "" {
			json.NewEncoder(w).Encode(map[string]any{
				"status":  "error",
				"message": "GEMINI_API_KEY not found in env or Key Vault",
				"key_set": false,
			})
			return
		}

		// Try a minimal ping to Gemini
		geminiURL := "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash:generateContent?key=" + geminiKey
		body := map[string]any{
			"contents": []any{map[string]any{"parts": []any{map[string]any{"text": "Say OK"}}}},
		}
		b, _ := json.Marshal(body)
		hresp, err := http.Post(geminiURL, "application/json", bytes.NewReader(b))
		if err != nil {
			json.NewEncoder(w).Encode(map[string]any{
				"status":  "error",
				"message": fmt.Sprintf("Network error: %v", err),
				"key_set": true,
				"detail":  err.Error(),
			})
			return
		}
		defer hresp.Body.Close()

		respBody, _ := io.ReadAll(hresp.Body)

		if hresp.StatusCode >= 400 {
			json.NewEncoder(w).Encode(map[string]any{
				"status":     "error",
				"key_set":    true,
				"message":    fmt.Sprintf("Gemini API returned HTTP %d", hresp.StatusCode),
				"detail":     string(respBody),
				"statusCode": hresp.StatusCode,
			})
			return
		}

		var gResp struct {
			Candidates []struct {
				Content struct {
					Parts []struct {
						Text string `json:"text"`
					} `json:"parts"`
				} `json:"content"`
				FinishReason string `json:"finishReason"`
			} `json:"candidates"`
		}
		json.Unmarshal(respBody, &gResp)

		reply := ""
		if len(gResp.Candidates) > 0 && len(gResp.Candidates[0].Content.Parts) > 0 {
			reply = gResp.Candidates[0].Content.Parts[0].Text
		}

		json.NewEncoder(w).Encode(map[string]any{
			"status":       "ok",
			"message":      "Gemini API is reachable and key is valid",
			"key_set":      true,
			"model":        "gemini-2.5-flash",
			"ping_reply":   reply,
			"finishReason": func() string { if len(gResp.Candidates) > 0 { return gResp.Candidates[0].FinishReason }; return "" }(),
		})
	}
}

// ── Infographic Generation ───────────────────────────────────────────────────

func infographicGenerateHandler(cfg config) http.HandlerFunc {
	type InfoGenRequest struct {
		Topic string `json:"topic"`
		Style string `json:"style"` // "modern", "minimalist", "technical"
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodPost {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		var req InfoGenRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
			return
		}

		geminiKey := cfg.getSecret("GEMINI_API_KEY")
		if geminiKey == "" {
			http.Error(w, `{"error":"GEMINI_API_KEY missing"}`, http.StatusServiceUnavailable)
			return
		}

		prompt := fmt.Sprintf(`You are an expert Infographic Designer and Cloud Architect.
Create a structured JSON layout for an infographic about: "%s".
The style should be: %s.

Return ONLY a JSON object with the following structure:
{
  "title": "Clear catchy title",
  "subtitle": "Informative subtitle",
  "sections": [
    {
      "icon": "Emoji representative",
      "heading": "Section Heading",
      "content": "Short concise bullet points or description (max 30 words)"
    }
  ],
  "visual_cue": "A description for an AI image generator to create a background or supporting visual for this infographic"
}

Keep it professional, architect-focused, and high-signal.`, req.Topic, req.Style)

		geminiURL := "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash:generateContent?key=" + geminiKey
		geminiReqBody := map[string]any{
			"contents": []any{map[string]any{"parts": []any{map[string]any{"text": prompt}}}},
			"generationConfig": map[string]any{
				"responseMimeType": "application/json",
			},
		}

		b, _ := json.Marshal(geminiReqBody)
		hresp, err := http.Post(geminiURL, "application/json", bytes.NewReader(b))
		if err != nil || hresp.StatusCode >= 400 {
			http.Error(w, `{"error":"Gemini generation failed"}`, http.StatusBadGateway)
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

		if len(gResp.Candidates) > 0 && len(gResp.Candidates[0].Content.Parts) > 0 {
			w.Write([]byte(gResp.Candidates[0].Content.Parts[0].Text))
		} else {
			http.Error(w, `{"error":"Empty response from AI"}`, http.StatusInternalServerError)
		}
	}
}

func infographicSaveHandler(cfg config) http.HandlerFunc {
	type InfoSaveRequest struct {
		ModuleNumber int            `json:"module_number"`
		VideoNumber  int            `json:"video_number"`
		SentenceID   *int64         `json:"sentence_id"`
		Topic        string         `json:"topic"`
		Style        string         `json:"style"`
		LayoutJSON   map[string]any `json:"layout_json"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodPost {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		var req InfoSaveRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
			return
		}

		// 1. Save to Azure Blob Storage first
		blobName := fmt.Sprintf("infographic_m%d_v%d_%d.json", req.ModuleNumber, req.VideoNumber, time.Now().Unix())
		container := "research-notes" // Storing layout JSON in notes container

		if cfg.azureAccountName != "" {
			expiry := time.Now().UTC().Add(10 * time.Minute)
			sasQuery, err := generateContainerSAS(cfg.azureAccountName, cfg.azureAccountKey, container, "rcwl", expiry)
			if err == nil {
				putURL := blobURL(cfg.azureAccountName, container, blobName, sasQuery)
				layoutBytes, _ := json.Marshal(req.LayoutJSON)
				
				ureq, _ := http.NewRequest(http.MethodPut, putURL, bytes.NewReader(layoutBytes))
				ureq.Header.Set("x-ms-blob-type", "BlockBlob")
				ureq.Header.Set("Content-Type", "application/json")

				uresp, err := http.DefaultClient.Do(ureq)
				if err != nil || uresp.StatusCode >= 400 {
					log.Printf("azure infographic upload failed: %v", err)
					blobName = "" // mark as failed for DB
				}
			} else {
				log.Printf("sas error for infographic: %v", err)
				blobName = ""
			}
		}

		// 2. Save to Supabase
		ctx := r.Context()
		dbEntry := map[string]any{
			"module_number":  req.ModuleNumber,
			"video_number":   req.VideoNumber,
			"sentence_id":    req.SentenceID,
			"topic":          req.Topic,
			"style":          req.Style,
			"layout_json":    req.LayoutJSON,
			"azure_blob_name": blobName,
			"status":         "saved_to_azure",
		}

		if err := supabasePost(ctx, cfg, "infographics", dbEntry); err != nil {
			log.Printf("supabase save infographic error: %v", err)
			http.Error(w, `{"error":"failed to save to database"}`, http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]any{
			"ok":        true,
			"blob_name": blobName,
		})
	}
}

// ── Lower Thirds Generation ──────────────────────────────────────────────────────

func lowerThirdGenerateHandler(cfg config) http.HandlerFunc {
	type LTRequest struct {
		ModuleNumber int `json:"module_number"`
		VideoNumber  int `json:"video_number"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodPost {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		var req LTRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
			return
		}

		geminiKey := cfg.getSecret("GEMINI_API_KEY")
		if geminiKey == "" {
			http.Error(w, `{"error":"GEMINI_API_KEY missing"}`, http.StatusServiceUnavailable)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()

		// Fetch module ID first
		var modules []struct {
			ID int `json:"id"`
		}
		modQ := fmt.Sprintf("select=id&module_number=eq.%d&limit=1", req.ModuleNumber)
		if err := supabaseGet(ctx, cfg, "modules", modQ, &modules); err != nil || len(modules) == 0 {
			http.Error(w, `{"error":"module not found"}`, http.StatusNotFound)
			return
		}

		// Fetch the video and script
		var videos []struct {
			ID    int    `json:"id"`
			Title string `json:"title"`
		}
		vidQ := fmt.Sprintf("select=id,title&module_id=eq.%d&video_number=eq.%d&limit=1", modules[0].ID, req.VideoNumber)
		if err := supabaseGet(ctx, cfg, "videos", vidQ, &videos); err != nil || len(videos) == 0 {
			http.Error(w, `{"error":"video not found"}`, http.StatusNotFound)
			return
		}
		videoID := videos[0].ID
		videoTitle := videos[0].Title

		var scripts []struct {
			ScriptText string `json:"script_text"`
		}
		scriptQ := fmt.Sprintf("select=script_text&video_id=eq.%d&limit=1", videoID)
		_ = supabaseGet(ctx, cfg, "scripts", scriptQ, &scripts)

		scriptContent := ""
		if len(scripts) > 0 {
			scriptContent = scripts[0].ScriptText
		}

		// Build Gemini prompt
		prompt := fmt.Sprintf(`You are an expert video production assistant for the "Claude AI Certification for Architects" course.

Module %d, Video %d: "%s"

%s

Generate 3 lower third overlay suggestions for this video. For each suggestion provide:
1. Main text (short, punchy, max 40 chars)
2. Sub text (descriptive, max 60 chars)
3. Brief rationale (why this fits the content)

Return JSON array:
[{"main":"...","sub":"...","rationale":"..."}]

Focus on professional, certification-quality overlays. Use the module/video theme.`, req.ModuleNumber, req.VideoNumber, videoTitle,
			func() string {
				if scriptContent != "" {
					return fmt.Sprintf("Here is the video script for context:\n\n%s\n\n---", scriptContent)
				}
				return "No script available. Generate based on the module and video title alone."
			}())

		geminiURL := "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash:generateContent?key=" + geminiKey
		geminiReqBody := map[string]any{
			"contents": []any{map[string]any{"parts": []any{map[string]any{"text": prompt}}}},
			"generationConfig": map[string]any{
				"responseMimeType": "application/json",
				"temperature":      0.7,
			},
		}

		b, _ := json.Marshal(geminiReqBody)
		hresp, err := http.Post(geminiURL, "application/json", bytes.NewReader(b))
		if err != nil || hresp.StatusCode >= 400 {
			http.Error(w, `{"error":"Gemini generation failed"}`, http.StatusBadGateway)
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

		responseText := "[]"
		if len(gResp.Candidates) > 0 && len(gResp.Candidates[0].Content.Parts) > 0 {
			responseText = gResp.Candidates[0].Content.Parts[0].Text
		}

		json.NewEncoder(w).Encode(map[string]any{
			"suggestions":   json.RawMessage(responseText),
			"script_text":   scriptContent,
			"video_title":   videoTitle,
			"module_number": req.ModuleNumber,
			"video_number":  req.VideoNumber,
		})
	}
}

// setCORS allows the static GitHub Pages site (and local dev) to call this
// backend cross-origin. The GitHub Pages site diverts /api calls here because
// Pages cannot run the Go backend that proxies to Gemini.
func setCORS(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	if strings.HasSuffix(origin, ".github.io") || strings.HasPrefix(origin, "http://localhost") {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
	w.Header().Set("Vary", "Origin")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, apikey")
}

func sanityCheckHandler(cfg config) http.HandlerFunc {
	type SanityRequest struct {
		ItemName string `json:"item_name"`
		ItemDesc string `json:"item_desc"`
		UserNote string `json:"user_note"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// Allow the static GitHub Pages site to call this backend cross-origin.
		setCORS(w, r)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodPost {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		var req SanityRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
			return
		}

		geminiKey := cfg.getSecret("GEMINI_API_KEY")
		if geminiKey == "" {
			http.Error(w, `{"error":"GEMINI_API_KEY missing"}`, http.StatusServiceUnavailable)
			return
		}

		prompt := fmt.Sprintf(`You are an expert Startup Consultant and Product Strategist. 
Analyze the following user finding/note for a specific Customer Discovery task.

**Task Name:** %s
**Task Description:** %s
**User Finding:** %s

Provide a "Sanity Check" feedback in Markdown format.
Include the following sections:
- **✅ Pros:** What is good about this finding or approach?
- **❌ Cons:** What are the potential risks or flaws?
- **🕵️ Gaps:** What is missing? What questions should the user ask next?

Be concise, critical, and constructive.`, req.ItemName, req.ItemDesc, req.UserNote)

		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()

		geminiURL := "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash:generateContent?key=" + geminiKey
		geminiReqBody := map[string]any{
			"contents": []any{map[string]any{"parts": []any{map[string]any{"text": prompt}}}},
		}

		b, _ := json.Marshal(geminiReqBody)
		hreq, err := http.NewRequestWithContext(ctx, "POST", geminiURL, bytes.NewReader(b))
		if err != nil {
			http.Error(w, `{"error":"failed to build request"}`, http.StatusInternalServerError)
			return
		}
		hreq.Header.Set("Content-Type", "application/json")

		hresp, err := http.DefaultClient.Do(hreq)
		if err != nil {
			http.Error(w, `{"error":"Gemini API call failed"}`, http.StatusBadGateway)
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
		if err := json.NewDecoder(hresp.Body).Decode(&gResp); err != nil {
			http.Error(w, `{"error":"failed to decode response"}`, http.StatusInternalServerError)
			return
		}

		feedback := ""
		if len(gResp.Candidates) > 0 && len(gResp.Candidates[0].Content.Parts) > 0 {
			feedback = gResp.Candidates[0].Content.Parts[0].Text
		}

		json.NewEncoder(w).Encode(map[string]string{
			"feedback": feedback,
		})
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
	mux.Handle("/api/images/generate", observe(cfg, imageGenerateHandler(cfg)))
	mux.Handle("/api/images/enhance-prompt", observe(cfg, imageEnhancePromptHandler(cfg)))
	mux.Handle("/api/images/save", observe(cfg, imageSaveHandler(cfg)))
	mux.Handle("/api/images/test-gemini", observe(cfg, imageTestGeminiHandler(cfg)))
	mux.Handle("/api/infographics/generate", observe(cfg, infographicGenerateHandler(cfg)))
	mux.Handle("/api/infographics/save", observe(cfg, infographicSaveHandler(cfg)))
	mux.Handle("/api/lowerthirds/generate", observe(cfg, lowerThirdGenerateHandler(cfg)))
	mux.Handle("/api/ai/sanity-check", observe(cfg, sanityCheckHandler(cfg)))
	mux.Handle("/admin/errors", observe(cfg, axiomErrorsHandler(tmpl, cfg, navConfigJS)))
	mux.Handle("/", observe(cfg, homeHandler(tmpl, cfg, navConfigJS)))

	addr := ":" + cfg.port
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
