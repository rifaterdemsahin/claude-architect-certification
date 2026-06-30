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
	"image"
	"image/color"
	"image/jpeg"
	_ "image/png"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"
)

// ── Config ────────────────────────────────────────────────────────────────────

type config struct {
	supabaseURL       string
	supabaseAnon      string
	axiomToken        string
	axiomDataset      string
	axiomAPIURL       string
	axiomQueryURL     string
	port              string
	azureAccountName  string
	azureAccountKey   string
	azureKeyVaultName string
	azureTenantID     string
	azureClientID     string
	azureClientSecret string
	googleClientID    string
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
	// Resolve the Google OAuth client ID once at startup so pages can read it from
	// /api/config instead of requiring it to be set in the browser (Folder Creator).
	cfg.googleClientID = firstNonEmpty(
		cfg.getSecret("claude-architect-GOOGLE-CLIENT-ID"),
		cfg.getSecret("google-oauth-client-id"),
		cfg.getSecret("GOOGLE_CLIENT_ID"),
	)
	return cfg
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
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

// isAdminRequest returns true when the caller is trusted: localhost origin OR a
// valid `admin_logged_in` cookie (set by /api/admin/login). Every privileged
// handler uses this so the rule is identical everywhere. UIs mirror it via
// /api/admin/status and hide destructive buttons when it is false.
func isAdminRequest(r *http.Request) bool {
	if c, err := r.Cookie("admin_logged_in"); err == nil && c.Value == "true" {
		return true
	}
	return false
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

// supabaseUpsert POSTs a row with the merge-duplicates Prefer header so it
// inserts-or-updates keyed on the given onConflict column(s). Requires the
// table to expose the onConflict set via a unique constraint (e.g.
// sentence_animations(sentence_id, animation_type)).
func supabaseUpsert(ctx context.Context, cfg config, table string, body any, onConflict string) error {
	u := fmt.Sprintf("%s/rest/v1/%s?select=", cfg.supabaseURL, table)
	if onConflict != "" {
		u += "&on_conflict=" + url.QueryEscape(onConflict)
	}
	b, _ := json.Marshal(body)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("apikey", cfg.supabaseAnon)
	req.Header.Set("Authorization", "Bearer "+cfg.supabaseAnon)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	// merge-duplicates = upsert; return=minimal keeps the payload tiny.
	req.Header.Set("Prefer", "resolution=merge-duplicates,return=minimal")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()
	if resp.StatusCode >= 400 {
		return fmt.Errorf("supabase UPSERT %s: %s", table, resp.Status)
	}
	return nil
}

func supabasePatch(ctx context.Context, cfg config, table, query string, body any) error {
	resp, err := supabaseReq(ctx, cfg, http.MethodPatch, table, query, body)
	if err != nil {
		return err
	}
	resp.Body.Close()
	if resp.StatusCode >= 400 {
		return fmt.Errorf("supabase PATCH %s: %s", table, resp.Status)
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
								if v, ok := m.Data["level"].(string); ok {
									ev.Level = v
								}
								if v, ok := m.Data["method"].(string); ok {
									ev.Method = v
								}
								if v, ok := m.Data["path"].(string); ok {
									ev.Path = v
								}
								if v, ok := m.Data["status"].(float64); ok {
									ev.Status = int(v)
								}
								if v, ok := m.Data["duration_ms"].(float64); ok {
									ev.Duration = int64(v)
								}
								if v, ok := m.Data["err"].(string); ok {
									ev.Err = v
								}
								if v, ok := m.Data["panic"].(string); ok {
									ev.Panic = v
								}
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

// axiomLogsHandler returns the latest Axiom events as JSON for the debug
// panel's "📊 Show Axiom Logs" button. Admin-gated (AXIOM_TOKEN is sensitive
// and logs may carry request details); reuses the same APL query as the
// /admin/errors HTML page.
func axiomLogsHandler(cfg config) http.HandlerFunc {
	type axiomLogEvent struct {
		Time     string `json:"_time"`
		Level    string `json:"level"`
		Stage    string `json:"stage"`
		Method   string `json:"method"`
		Path     string `json:"path"`
		Status   int    `json:"status"`
		Duration int64  `json:"duration_ms"`
		Err      string `json:"err"`
		Message  string `json:"message"`
		URL      string `json:"url"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		setCORS(w, r)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if !isAdminRequest(r) {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]any{"error": "unauthorized — sign in as admin", "events": []any{}})
			return
		}
		limit := 50
		if l := r.URL.Query().Get("limit"); l != "" {
			if n, err := strconv.Atoi(l); err == nil && n > 0 && n <= 500 {
				limit = n
			}
		}
		if cfg.axiomToken == "" {
			json.NewEncoder(w).Encode(map[string]any{"error": "AXIOM_TOKEN not set on server", "events": []any{}})
			return
		}
		apl := fmt.Sprintf(`['%s'] | sort by _time desc | limit %d`, cfg.axiomDataset, limit)
		body, _ := json.Marshal(map[string]any{
			"apl":       apl,
			"startTime": "now-24h",
		})
		queryURL := fmt.Sprintf("%s/v1/datasets/%s/query", cfg.axiomQueryURL, cfg.axiomDataset)
		log.Printf("axiom logs -> %s (apl: %s)", queryURL, apl)
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, queryURL, bytes.NewReader(body))
		if err != nil {
			json.NewEncoder(w).Encode(map[string]any{"error": err.Error(), "events": []any{}})
			return
		}
		req.Header.Set("Authorization", "Bearer "+cfg.axiomToken)
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]any{"error": err.Error(), "events": []any{}})
			return
		}
		defer resp.Body.Close()
		respBody, _ := io.ReadAll(resp.Body)
		if resp.StatusCode >= 300 {
			json.NewEncoder(w).Encode(map[string]any{"error": fmt.Sprintf("Axiom API %d: %s", resp.StatusCode, string(respBody)), "events": []any{}})
			return
		}
		var qr axiomQueryResp
		if err := json.Unmarshal(respBody, &qr); err != nil {
			json.NewEncoder(w).Encode(map[string]any{"error": "decode: " + err.Error(), "events": []any{}})
			return
		}
		events := []axiomLogEvent{}
		for _, m := range qr.Matches {
			ev := axiomLogEvent{Time: m.Time}
			if v, ok := m.Data["level"].(string); ok {
				ev.Level = v
			}
			if v, ok := m.Data["stage"].(string); ok {
				ev.Stage = v
			}
			if v, ok := m.Data["method"].(string); ok {
				ev.Method = v
			}
			if v, ok := m.Data["path"].(string); ok {
				ev.Path = v
			}
			if v, ok := m.Data["status"].(float64); ok {
				ev.Status = int(v)
			}
			if v, ok := m.Data["duration_ms"].(float64); ok {
				ev.Duration = int64(v)
			}
			if v, ok := m.Data["err"].(string); ok {
				ev.Err = v
			}
			if v, ok := m.Data["message"].(string); ok {
				ev.Message = v
			}
			if v, ok := m.Data["url"].(string); ok {
				ev.URL = v
			}
			events = append(events, ev)
		}
		log.Printf("axiom logs query returned %d events", len(events))
		// Axiom's /query endpoint ignores the APL `| limit N` clause (it caps at
		// 1000 matches), so enforce the requested limit here. Events are already
		// newest-first from the APL `sort by _time desc`.
		if len(events) > limit {
			events = events[:limit]
		}
		json.NewEncoder(w).Encode(map[string]any{"events": events, "count": len(events)})
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
		GoogleClientID   string `json:"googleClientId"`
	}
	payload, _ := json.Marshal(configResp{
		SupabaseURL:      cfg.supabaseURL,
		SupabaseAnon:     cfg.supabaseAnon,
		AzureAccountName: cfg.azureAccountName,
		GoogleClientID:   cfg.googleClientID,
	})
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "public, max-age=300")
		w.Write(payload)
	}
}

// maskSecret returns a non-reversible preview of a secret so the env page can
// confirm a value is loaded without exposing it (e.g. "ab••••… (40 chars)").
func maskSecret(v string) string {
	if v == "" {
		return ""
	}
	n := len(v)
	if n <= 8 {
		return fmt.Sprintf("•••• (%d chars)", n)
	}
	return fmt.Sprintf("%s••••%s (%d chars)", v[:2], v[n-2:], n)
}

// envStatusHandler reports which environment values the server loaded at
// startup. Gated to localhost or an admin session (same policy as the gdrive
// credentials endpoint). Secret values are masked — only presence + a short
// preview is returned, never the raw secret.
func envStatusHandler(cfg config) http.HandlerFunc {
	type envEntry struct {
		Key    string `json:"key"`
		Group  string `json:"group"`
		Set    bool   `json:"set"`
		Secret bool   `json:"secret"`
		Value  string `json:"value"`
		Note   string `json:"note,omitempty"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if !isAdminRequest(r) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error":"Unauthorized — sign in as admin or open from localhost."}`))
			return
		}

		show := func(key, group, val string, secret bool, note string) envEntry {
			e := envEntry{Key: key, Group: group, Set: val != "", Secret: secret, Note: note}
			if val == "" {
				e.Value = ""
			} else if secret {
				e.Value = maskSecret(val)
			} else {
				e.Value = val
			}
			return e
		}

		entries := []envEntry{
			show("SUPABASE_URL", "Supabase", cfg.supabaseURL, false, ""),
			show("SUPABASE_ANON_KEY", "Supabase", cfg.supabaseAnon, true, "Public anon key (RLS-protected)"),
			show("AXIOM_DATASET", "Axiom", cfg.axiomDataset, false, ""),
			show("AXIOM_API_URL", "Axiom", cfg.axiomAPIURL, false, ""),
			show("AXIOM_QUERY_URL", "Axiom", cfg.axiomQueryURL, false, ""),
			show("AXIOM_TOKEN", "Axiom", cfg.axiomToken, true, ""),
			show("PORT", "Server", cfg.port, false, ""),
			show("AZURE_KEYVAULT_NAME", "Azure", cfg.azureKeyVaultName, false, ""),
			show("AZURE_TENANT_ID", "Azure", cfg.azureTenantID, true, ""),
			show("AZURE_CLIENT_ID", "Azure", cfg.azureClientID, true, ""),
			show("AZURE_CLIENT_SECRET", "Azure", cfg.azureClientSecret, true, ""),
			show("AZURE_STORAGE (account name)", "Azure", cfg.azureAccountName, false, "Parsed from AZURE_STORAGE_CONN_STR"),
			show("AZURE_STORAGE (account key)", "Azure", cfg.azureAccountKey, true, "Parsed from AZURE_STORAGE_CONN_STR"),
			show("GOOGLE_CLIENT_ID", "Google", cfg.googleClientID, false, "OAuth client ID (public; served via /api/config)"),
			show("RUNPOD_API_KEY", "RunPod", cfg.getSecret("RUNPOD_API_KEY"), true, "Key Vault 'runpod-api-key' or env; powers the Animation Generator render"),
			show("RUNPOD_ENDPOINT_ID", "RunPod", os.Getenv("RUNPOD_ENDPOINT_ID"), false, "Serverless endpoint id for the Remotion worker"),
			show("REMOTION_SERVE_URL", "RunPod", os.Getenv("REMOTION_SERVE_URL"), false, "Deployed Remotion bundle URL the worker renders"),
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-store")
		json.NewEncoder(w).Encode(map[string]any{
			"keyVaultActive": cfg.azureKeyVaultName != "" && cfg.azureTenantID != "" && cfg.azureClientID != "" && cfg.azureClientSecret != "",
			"entries":        entries,
		})
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
	"research-images":     true,
	"research-audio":      true,
	"research-videos":     true,
	"research-notes":      true,
	"research-animations": true, // Remotion MP4s from the Animation Generator
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

// researchFileProxyURL builds the same-origin proxy path the pages use to load
// a blob. The research-images container is private, so raw blob URLs 404 — blobs
// must be served through /api/research/file, which signs a short-lived SAS.
func researchFileProxyURL(container, name string) string {
	return "/api/research/file?container=" + container + "&name=" + url.QueryEscape(name)
}

// uploadBlobToAzure PUTs a block blob into the given container. Used for both
// originals and thumbnails so the upload logic lives in one place.
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

// downloadBlobFromAzure fetches a blob's bytes (and content type) via a
// short-lived read SAS. Used by the rename flow, which re-uploads the same
// bytes under a new name (Azure blobs cannot be renamed in place).
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

// deleteBlobFromAzure removes a blob via a short-lived delete SAS. Best-effort
// callers ignore the error (e.g. a missing thumbnail).
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

// thumbPrefix matches the convention used by the image gallery pages
// (research/images.html, scripts/index.html): a thumbnail for blob "X" is the
// blob "__thumb__X" in the same container. Keeping the original extension in
// the name lets those pages discover thumbnails by listing the container.
const thumbPrefix = "__thumb__"

// thumbBlobName derives the thumbnail blob name from an original blob name,
// e.g. "m1_v2_123.png" → "__thumb__m1_v2_123.png".
func thumbBlobName(original string) string {
	return thumbPrefix + original
}

// generateThumbnail decodes a PNG/JPEG image and produces a downscaled JPEG
// thumbnail whose longest edge is capped at maxEdge px, preserving aspect
// ratio. Uses a box-average downscale — stdlib only, no external deps.
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

		// Read the upload into memory so we can both store the original and
		// derive a thumbnail from the same bytes.
		data, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "read upload: "+err.Error(), http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		// ── Upload the original ──
		if err := uploadBlobToAzure(ctx, cfg, container, blobName, contentType, data); err != nil {
			http.Error(w, "azure upload: "+err.Error(), http.StatusBadGateway)
			return
		}

		// ── Generate + upload a thumbnail so gallery grids load small images,
		//    not full-resolution originals (best-effort; original still saved). ──
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
			// 🛡️ Destructive: require sign-in (localhost or admin cookie). Without
			// this an unsigned-in visitor on the deployed site could delete blobs.
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

// researchRenameHandler renames a blob (and its __thumb__ companion) by copying
// the bytes to the new name and deleting the old — Azure has no in-place rename.
// It also repoints research_relationships.item_name so existing image↔video and
// image↔sentence links survive the rename. The file extension (suffix) is
// preserved: the new name must keep the original extension.
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
		// Disallow path separators and thumbnail-internal names.
		if strings.ContainsAny(to, "/\\") || strings.HasPrefix(to, thumbPrefix) || strings.HasPrefix(from, thumbPrefix) {
			http.Error(w, "invalid name", http.StatusBadRequest)
			return
		}
		// The suffix must not change: the new name keeps the original extension.
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

		// Copy the original blob to the new name, then delete the old.
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

		// Best-effort: move the thumbnail companion too.
		oldThumb, newThumb := thumbBlobName(from), thumbBlobName(to)
		if tData, tct, terr := downloadBlobFromAzure(ctx, cfg, container, oldThumb); terr == nil {
			if uerr := uploadBlobToAzure(ctx, cfg, container, newThumb, tct, tData); uerr != nil {
				log.Printf("research rename: thumbnail copy failed for %s: %v", newThumb, uerr)
			} else if derr := deleteBlobFromAzure(ctx, cfg, container, oldThumb); derr != nil {
				log.Printf("research rename: thumbnail delete failed for %s: %v", oldThumb, derr)
			}
		}

		// Repoint any relationship rows so existing links keep resolving.
		if err := supabasePatch(ctx, cfg, "research_relationships",
			"item_name=eq."+url.QueryEscape(from),
			map[string]string{"item_name": to}); err != nil {
			log.Printf("research rename: relationship update failed %s→%s: %v", from, to, err)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"ok": "true", "name": to, "thumbnail": newThumb})
	}
}

// ── Explanations ──────────────────────────────────────────────────────────────

type Explanation struct {
	ID              int    `json:"id,omitempty"`
	EntityType      string `json:"entity_type"`
	EntityID        string `json:"entity_id"`
	ExplanationText string `json:"explanation_text"`
	GeneratedBy     string `json:"generated_by"`
	CreatedAt       string `json:"created_at,omitempty"`
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

		ctx := r.Context()
		blobName := fmt.Sprintf("m%d_v%d_%d.png", req.ModuleNumber, req.VideoNumber, time.Now().Unix())
		container := "research-images"

		// ── Upload the original ──
		if err := uploadBlobToAzure(ctx, cfg, container, blobName, contentType, imageData); err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"azure upload failed: %s"}`, err.Error()), http.StatusBadGateway)
			return
		}

		// ── Generate + upload the thumbnail (best-effort; original still saved if this fails) ──
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

		// Store loadable same-origin proxy URLs (the container is private, so
		// raw blob URLs 404 in the browser). Blob names remain the canonical ref.
		imageURL := researchFileProxyURL(container, blobName)
		var thumbURL string
		if thumbName != "" {
			thumbURL = researchFileProxyURL(container, thumbName)
		}

		// ── Save both references to Supabase ──
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

// DrawingSaveRequest is the payload from drawing_generator.html when the user
// hits 💾 Save Drawing: the editable Excalidraw scene plus a flattened PNG
// (data URL) of the same drawing, tied to a single sentence.
type DrawingSaveRequest struct {
	SentenceID     int             `json:"sentence_id"`
	ModuleNumber   int             `json:"module_number"`
	VideoNumber    int             `json:"video_number"`
	ExcalidrawJSON json.RawMessage `json:"excalidraw_json"`
	PNG            string          `json:"png"` // data URL (image/png;base64,...)
	Prompt         string          `json:"prompt"`
}

// drawingSaveHandler persists one per-sentence Excalidraw drawing: it uploads
// the flattened PNG to Azure Blob Storage (research-images container) and
// upserts the editable scene + blob reference into the sentence_drawings table,
// keyed on sentence_id. Mirrors imageSaveHandler's Azure path.
// POST /api/drawings/save
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

		// PNG is optional — a scene can be saved without a rendered image — but
		// when present we push it to Azure so the drawing has a flat thumbnail.
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

		// Upsert the editable scene + blob reference keyed on sentence_id.
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

// generatedImageRow is a minimal view of generated_images for the backfill job.
type generatedImageRow struct {
	ID            int64  `json:"id"`
	AzureBlobName string `json:"azure_blob_name"`
	ImageURL      string `json:"image_url"`
	ThumbnailURL  string `json:"thumbnail_url"`
}

// imageThumbnailBackfillHandler walks every generated_images row that has an
// original blob but no thumbnail, downloads the original from Azure, builds a
// thumbnail, uploads it, and patches the row with thumbnail_url/thumbnail_blob_name.
// POST /api/images/backfill-thumbnails
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
		// Rows missing a thumbnail but having an original blob.
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
			// Download the original via a short-lived read SAS.
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
			"status":     "ok",
			"message":    "Gemini API is reachable and key is valid",
			"key_set":    true,
			"model":      "gemini-2.5-flash",
			"ping_reply": reply,
			"finishReason": func() string {
				if len(gResp.Candidates) > 0 {
					return gResp.Candidates[0].FinishReason
				}
				return ""
			}(),
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
		if !isAdminRequest(r) {
			http.Error(w, `{"error":"Unauthorized — sign in as admin to save infographics."}`, http.StatusUnauthorized)
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
			"module_number":   req.ModuleNumber,
			"video_number":    req.VideoNumber,
			"sentence_id":     req.SentenceID,
			"topic":           req.Topic,
			"style":           req.Style,
			"layout_json":     req.LayoutJSON,
			"azure_blob_name": blobName,
			"status":          "saved_to_azure",
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
		ModuleNumber int  `json:"module_number"`
		VideoNumber  int  `json:"video_number"`
		Force        bool `json:"force"`
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

		ctx, cancel := context.WithTimeout(r.Context(), 40*time.Second)
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

		// Check existing lower thirds for this video
		var existingLTs []struct {
			ID        int    `json:"id"`
			MainText  string `json:"main_text"`
			SubText   string `json:"sub_text"`
			Rationale string `json:"rationale"`
			SortOrder int    `json:"sort_order"`
		}
		ltQ := fmt.Sprintf("select=id,main_text,sub_text,rationale,sort_order&module_number=eq.%d&video_number=eq.%d&order=sort_order.asc", req.ModuleNumber, req.VideoNumber)
		_ = supabaseGet(ctx, cfg, "lower_thirds", ltQ, &existingLTs)

		if len(existingLTs) > 0 && !req.Force {
			var suggestions []map[string]any
			for _, lt := range existingLTs {
				suggestions = append(suggestions, map[string]any{
					"main":       lt.MainText,
					"sub":        lt.SubText,
					"rationale":  lt.Rationale,
					"db_id":      lt.ID,
					"sort_order": lt.SortOrder,
					"from_cache": true,
				})
			}

			var scripts []struct {
				ScriptText string `json:"script_text"`
			}
			scriptQ := fmt.Sprintf("select=script_text&video_id=eq.%d&limit=1", videoID)
			scriptContent := ""
			_ = supabaseGet(ctx, cfg, "scripts", scriptQ, &scripts)
			if len(scripts) > 0 {
				scriptContent = scripts[0].ScriptText
			}

			json.NewEncoder(w).Encode(map[string]any{
				"suggestions":   suggestions,
				"script_text":   scriptContent,
				"video_title":   videoTitle,
				"module_number": req.ModuleNumber,
				"video_number":  req.VideoNumber,
				"from_cache":    true,
			})
			return
		}

		// Fetch script for Gemini context
		var scripts []struct {
			ScriptText string `json:"script_text"`
		}
		scriptQ := fmt.Sprintf("select=script_text&video_id=eq.%d&limit=1", videoID)
		_ = supabaseGet(ctx, cfg, "scripts", scriptQ, &scripts)

		scriptContent := ""
		if len(scripts) > 0 {
			scriptContent = scripts[0].ScriptText
		}

		geminiKey := cfg.getSecret("GEMINI_API_KEY")
		if geminiKey == "" {
			http.Error(w, `{"error":"GEMINI_API_KEY missing"}`, http.StatusServiceUnavailable)
			return
		}

		// Build Gemini prompt
		prompt := fmt.Sprintf(`You are an expert video production assistant for the "Claude AI Certification for Architects" course.

Module %d, Video %d: "%s"

%s

Generate 3 lower third overlay suggestions for this video. For each suggestion provide:
1. Main text (short, punchy, max 40 chars)
2. Sub text (descriptive, max 60 chars)
3. Rationale — why this lower third is important for the audience to learn (1-2 sentences)

Return a JSON array ONLY:
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

		type Suggestion struct {
			Main      string `json:"main"`
			Sub       string `json:"sub"`
			Rationale string `json:"rationale"`
		}

		var suggestions []Suggestion
		if len(gResp.Candidates) > 0 && len(gResp.Candidates[0].Content.Parts) > 0 {
			raw := gResp.Candidates[0].Content.Parts[0].Text
			if err := json.Unmarshal([]byte(raw), &suggestions); err != nil {
				// Try wrapping if it's a single object
				var single Suggestion
				if json.Unmarshal([]byte(raw), &single) == nil {
					suggestions = []Suggestion{single}
				}
			}
		}

		if suggestions == nil {
			suggestions = []Suggestion{}
		}

		// Save each suggestion to lower_thirds table (insert, conflicts are ignored)
		for i, s := range suggestions {
			body := map[string]any{
				"module_number": req.ModuleNumber,
				"video_number":  req.VideoNumber,
				"module_id":     modules[0].ID,
				"video_id":      videoID,
				"main_text":     s.Main,
				"sub_text":      s.Sub,
				"rationale":     s.Rationale,
				"sort_order":    i + 1,
			}
			_ = supabasePost(ctx, cfg, "lower_thirds", body)
		}

		// Convert to response format
		var respSuggestions []map[string]any
		for i, s := range suggestions {
			respSuggestions = append(respSuggestions, map[string]any{
				"main":       s.Main,
				"sub":        s.Sub,
				"rationale":  s.Rationale,
				"sort_order": i + 1,
				"from_cache": false,
			})
		}

		json.NewEncoder(w).Encode(map[string]any{
			"suggestions":   respSuggestions,
			"script_text":   scriptContent,
			"video_title":   videoTitle,
			"module_number": req.ModuleNumber,
			"video_number":  req.VideoNumber,
			"from_cache":    false,
		})
	}
}

func openRouterGenerateHandler(cfg config) http.HandlerFunc {
	type ORRequest struct {
		Script              string `json:"script"`
		CourseName          string `json:"courseName"`
		ModuleName          string `json:"moduleName"`
		VideoName           string `json:"videoName"`
		Presenter           string `json:"presenter"`
		ExistingLowerThirds string `json:"existingLowerThirds"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodPost {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		var req ORRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
			return
		}

		apiKey := cfg.getSecret("OPENROUTER_API_KEY")
		if apiKey == "" {
			log.Printf("OPENROUTER_API_KEY unavailable: not in Key Vault and not in env")
			http.Error(w, `{"error":"OPENROUTER_API_KEY missing from server configuration"}`, http.StatusServiceUnavailable)
			return
		}

		prompt := fmt.Sprintf(`You are the lower-thirds editor for a technical teaching video.
Your goal is to help a learner follow along: name the speaker/module,
put a label + one-line gloss on screen the moment a domain-specific term
is first spoken, reinforce each learning objective, and mark the key takeaway.

Course: %s
Module: %s
Video: %s
Presenter: %s   (use only for the opening speaker_id; skip if empty)

Below is the scene-by-scene script. Walk it IN ORDER, scene by scene.
For each scene, decide whether an on-screen lower third would aid
comprehension. Emit one when it helps; emit none for filler or pure
transition scenes. Do NOT pad to a fixed number and do NOT cap the count —
let the content decide. As a guide, expect roughly one lower third per
distinct idea, term, or section, and skip anything a viewer wouldn't
need spelled out.

When a domain-specific or technical term is first spoken, emit a
"term_definition" so the viewer sees the term plus a plain-language gloss.

Lower-third types to choose from:
- speaker_id      (who's talking / module banner — opening only)
- section_title   (entering a new part of the video)
- term_definition (a technical term just appeared — label + define it)
- key_point       (a claim worth reinforcing on screen)
- takeaway        (the one thing to remember)

Constraints: "main" ≤ 40 chars, "sub" ≤ 60 chars.

DO NOT GENERATE these existing lower thirds (use this as negative prompt):
%s

SCRIPT (scene-anchored — use this, never the prose version):
%s

Return ONLY a JSON array, in scene order, with no markdown fences:
[
  {"scene": <int>, "type": "<one of the types above>", "main": "...", "sub": "...", "rationale": "..."}
]`, req.CourseName, req.ModuleName, req.VideoName, req.Presenter, req.ExistingLowerThirds, req.Script)

		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()

		orURL := "https://openrouter.ai/api/v1/chat/completions"
		reqBody := map[string]any{
			"model": "google/gemini-2.5-flash",
			"messages": []map[string]string{
				{"role": "user", "content": prompt},
			},
			"response_format": map[string]map[string]string{
				"type": {"type": "json_object"},
			},
		}

		b, _ := json.Marshal(reqBody)
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
			http.Error(w, `{"error":"OpenRouter API call failed"}`, http.StatusBadGateway)
			return
		}
		defer hresp.Body.Close()

		if hresp.StatusCode >= 400 {
			body, _ := io.ReadAll(hresp.Body)
			http.Error(w, fmt.Sprintf(`{"error":"OpenRouter API returned HTTP %d: %s"}`, hresp.StatusCode, body), http.StatusBadGateway)
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

		content := ""
		if len(orResp.Choices) > 0 {
			content = orResp.Choices[0].Message.Content
		}

		// Return the exact prompt that was sent so the UI can expose the
		// "actual executed prompt" for feedback/debugging (no client rebuild drift).
		json.NewEncoder(w).Encode(map[string]string{
			"content": content,
			"prompt":  prompt,
			"model":   "google/gemini-2.5-flash",
		})
	}
}

func drawingGenerateHandler(cfg config) http.HandlerFunc {
	type DrawingRequest struct {
		Sentences  string `json:"sentences"`
		CourseName string `json:"courseName"`
		ModuleName string `json:"moduleName"`
		VideoName  string `json:"videoName"`
	}

	const model = "anthropic/claude-sonnet-4.6"

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodPost {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		var req DrawingRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
			return
		}

		apiKey := cfg.getSecret("OPENROUTER_API_KEY")
		if apiKey == "" {
			http.Error(w, `{"error":"OPENROUTER_API_KEY missing from server configuration"}`, http.StatusServiceUnavailable)
			return
		}

		prompt := fmt.Sprintf(`You are a technical architect drawing an Excalidraw diagram.
Create a valid Excalidraw JSON object that visually represents the core architecture, data flow, or concepts in the following script.

Course: %s
Module: %s
Video: %s

Script:
%s

# Output format
Return ONLY valid JSON (no markdown fences, no commentary). The JSON must have this exact structure (do not deviate):
{
  "type": "excalidraw",
  "version": 2,
  "source": "https://excalidraw.com",
  "elements": [
    { "type": "rectangle", "x": 100, "y": 100, "width": 200, "height": 100, "strokeColor": "#000000", "backgroundColor": "transparent", "fillStyle": "hachure", "strokeWidth": 1, "strokeStyle": "solid", "roughness": 1, "opacity": 100, "id": "rect1" },
    { "type": "text", "x": 120, "y": 120, "text": "Example Text", "fontSize": 20, "fontFamily": 1, "textAlign": "left", "verticalAlign": "top", "strokeColor": "#000000", "id": "text1" }
  ]
}

Make sure x and y coordinates are laid out logically so elements do not overlap.
Use sensible colors and ensure text fits within rectangles.
`, req.CourseName, req.ModuleName, req.VideoName, req.Sentences)

		ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
		defer cancel()

		orURL := "https://openrouter.ai/api/v1/chat/completions"
		reqBody := map[string]any{
			"model": model,
			"messages": []map[string]string{
				{"role": "user", "content": prompt},
			},
			"max_tokens": 4000,
		}

		b, _ := json.Marshal(reqBody)
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
			http.Error(w, `{"error":"OpenRouter API call failed"}`, http.StatusBadGateway)
			return
		}
		defer hresp.Body.Close()

		if hresp.StatusCode >= 400 {
			body, _ := io.ReadAll(hresp.Body)
			http.Error(w, fmt.Sprintf(`{"error":"OpenRouter API returned HTTP %d: %s"}`, hresp.StatusCode, body), http.StatusBadGateway)
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

		content := ""
		if len(orResp.Choices) > 0 {
			content = strings.TrimSpace(orResp.Choices[0].Message.Content)
			content = strings.TrimPrefix(content, "```json")
			content = strings.TrimPrefix(content, "```")
			content = strings.TrimSuffix(content, "```")
			content = strings.TrimSpace(content)
		}

		json.NewEncoder(w).Encode(map[string]string{
			"content": content,
			"prompt":  prompt,
			"model":   model,
		})
	}
}

// slideGenerateHandler powers the "Generate Slides Complex" button on the Marp
// Slide Generator. It builds a napkin.ai-style visual-slide prompt from the
// selected script sentences and proxies to OpenRouter. The model is chosen for
// design/structure reasoning quality (see below), keeping the API key server-side.
func slideGenerateHandler(cfg config) http.HandlerFunc {
	type SlideRequest struct {
		Sentences        string `json:"sentences"`
		CourseName       string `json:"courseName"`
		ModuleName       string `json:"moduleName"`
		VideoName        string `json:"videoName"`
		BrandFrontmatter string `json:"brandFrontmatter"`
	}

	// Best model for the job: napkin.ai-style slides need strong layout/visual
	// reasoning and reliable structured-markdown output. Anthropic's Claude
	// Sonnet 4.6 leads on design-structured generation at a sensible cost, so we
	// use it here (vs. the cheaper gemini-2.5-flash used for bulk lower-thirds).
	const slideModel = "anthropic/claude-sonnet-4.6"

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodPost {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		var req SlideRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
			return
		}

		apiKey := cfg.getSecret("OPENROUTER_API_KEY")
		if apiKey == "" {
			log.Printf("OPENROUTER_API_KEY unavailable: not in Key Vault and not in env")
			http.Error(w, `{"error":"OPENROUTER_API_KEY missing from server configuration"}`, http.StatusServiceUnavailable)
			return
		}

		brand := req.BrandFrontmatter
		if strings.TrimSpace(brand) == "" {
			brand = `---
marp: true
theme: default
class: lead
backgroundColor: #030712
color: #f3f4f6
style: |
  section { font-family: 'Plus Jakarta Sans', sans-serif; display: flex; flex-direction: column; justify-content: center; align-items: center; text-align: center; }
  h1 { font-family: 'Outfit', sans-serif; color: #8b5cf6; font-size: 3.4rem; font-weight: 800; margin: 0.4rem 0; line-height: 1.1; }
  h2 { font-size: 4rem; margin: 0; }
  ul { font-size: 1.5rem; text-align: left; line-height: 1.5; }
---`
		}

		prompt := fmt.Sprintf(`You are a world-class visual presentation designer in the style of napkin.ai:
you turn raw script text into clean, visual, scannable slides where each idea is
expressed as a tiny "visual" — a bold icon, a punchy headline, and a few
supporting points or a simple left-to-right flow (e.g. "Idea -> Build -> Ship").

Course: %s
Module: %s
Video: %s

# Design rules (napkin.ai style)
1. One core idea per slide. Lead with a single large, highly relevant emoji icon
   that visually represents the idea.
2. Headline = the punchy core message in 3-8 words. No full sentences.
3. Add 2-4 short supporting bullet points (max ~6 words each) OR a simple flow
   line using arrows (->) when the idea is a process or sequence.
4. Group the sentences intelligently: merge closely-related sentences into one
   visual slide instead of one slide per sentence. Aim for high signal, low noise.
5. Open with a title slide and close with a single takeaway slide.
6. Keep it visual and minimal — think keynote, not a document.

# Output format
Return ONLY raw Marp markdown (no markdown fences, no commentary).
Start with EXACTLY this frontmatter block (do not alter it):

%s

Then one slide per idea, separated by a line containing only ---, each formatted as:

## [single big emoji]

# [Punchy headline]

- [supporting point]
- [supporting point]

# Input script sentences
%s
`, req.CourseName, req.ModuleName, req.VideoName, brand, req.Sentences)

		ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
		defer cancel()

		orURL := "https://openrouter.ai/api/v1/chat/completions"
		reqBody := map[string]any{
			"model": slideModel,
			"messages": []map[string]string{
				{"role": "user", "content": prompt},
			},
			// A slide deck's markdown is small; cap output so the request stays
			// well within account credits (an unbounded ceiling triggers HTTP 402).
			"max_tokens": 4000,
		}

		b, _ := json.Marshal(reqBody)
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
			http.Error(w, `{"error":"OpenRouter API call failed"}`, http.StatusBadGateway)
			return
		}
		defer hresp.Body.Close()

		if hresp.StatusCode >= 400 {
			body, _ := io.ReadAll(hresp.Body)
			http.Error(w, fmt.Sprintf(`{"error":"OpenRouter API returned HTTP %d: %s"}`, hresp.StatusCode, body), http.StatusBadGateway)
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

		content := ""
		if len(orResp.Choices) > 0 {
			content = orResp.Choices[0].Message.Content
		}

		json.NewEncoder(w).Encode(map[string]string{
			"content": content,
			"prompt":  prompt,
			"model":   slideModel,
		})
	}
}

// ── Animation Generator (RunPod serverless + Remotion) ───────────────────────
//
// Three handlers back 5_Symbols/production/postprod/animation_generator.html:
//  1. POST /api/animations/generate-prompt  — builds the Remotion composition
//     prompt + inputProps + RunPod payload for a (sentence, animationType).
//     Open (no cost, no secret leaked — only the prompt text is returned).
//  2. POST /api/animations/runpod/run       — admin-gated. Submits the job to
//     RunPod serverless and upserts a `sentence_animations` row (status=generating).
//  3. GET  /api/animations/runpod/status    — admin-gated. Polls RunPod; on
//     COMPLETED downloads the MP4, uploads it to Azure `research-animations`,
//     and patches the row (status=completed + animation_url). Idempotent.

// animationTypeMeta describes one of the 10 course-content animation types:
// its slug, display label, emoji, and the Remotion composition prompt builder.
// The prompt is a complete spec an LLM or developer can turn straight into a
// Remotion <Composition> — it encodes the animation choreography + the sentence
type AnimationTypeMeta struct {
	Slug        string
	Label       string
	Emoji       string
	BuildPrompt func(sentence, sentenceType, moduleTitle, videoTitle string) string
}

func clampLen(s string, n int) string {
	s = strings.TrimSpace(s)
	if len(s) > n {
		return s[:n]
	}
	return s
}

// animationTypeMetas is the single source of truth for the 10 animation types
// and their Remotion prompt templates. The page mirrors this list for its
// <select> options; the prompt handler renders the chosen template.
func animationTypeMetas() []AnimationTypeMeta {
	brand := `BRAND KIT: background #030712 → deep navy; primary #8b5cf6 (violet); secondary #3b82f6 (blue); success #10b981; text #f3f4f6; muted #9ca3af. Fonts: 'Outfit' (headings, 800/900), 'Plus Jakarta Sans' (body). 1920x1080, 30fps. Easing: Remotion spring() for entrances, interpolate() with Easing.inOut for exits. Always export h264 MP4.`
	mk := func(name, goal, choreography string) string {
		return fmt.Sprintf(`# Remotion composition: %s

GOAL: %s

%s

CONSTRAINTS:
- React + Remotion only (@remotion/player, useCurrentFrame, useVideoConfig, spring, interpolate, AbsoluteFill, Sequence, Img, Audio). No external chart libs.
- Read every value from `+"`inputProps`"+` (props.title, props.subtitle, props.items[], props.metric, props.brandColor…) — NEVER hard-code the sentence text in the component; the same <Composition> renders every sentence via props.
- 1920x1080, props.fps || 30, props.durationInFrames || 150.
- One focal point at a time; text must be legible at 1280px wide.
- %s

OUTPUT: a single self-contained React+Remotion <Composition id="Main"> whose defaultProps match the inputProps below. Return ONLY the component source + defaultProps, no prose.`, name, goal, choreography, brand)
	}

	return []AnimationTypeMeta{
		{
			Slug: "architecture_diagram", Label: "Architecture Diagram", Emoji: "🏛️",
			BuildPrompt: func(s, st, m, v string) string {
				return mk("ArchitectureDiagram", "Animate a system architecture building up node-by-node so the viewer sees how the pieces connect (e.g. Claude + MCP + Supabase + Azure Key Vault).",
					fmt.Sprintf(`CHOREOGRAPHY (150 frames):
- f0–30: title "%s" fades in top-center (spring scale 0.8→1).
- f30–135: each node (props.nodes[] = {id,label,x,y}) pops in with spring() staggered 12 frames apart, connectors (SVG <path>) draw via pathLength interpolate after both endpoints exist.
- f135–150: the active node (props.activeNode) gets a violet glow pulse (boxShadow animate) + its label scales 1.1.
INPUT PROPS: title="%s" (from module %q / video %q), subtitle=derived from sentence, nodes[]=from sentence entities, sentenceType=%q.
SENTENCE (authoritative content): %q`, "System Architecture", clampLen(s, 90), m, v, st, s))
			},
		},
		{
			Slug: "data_flow", Label: "Data Flow", Emoji: "🌊",
			BuildPrompt: func(s, st, m, v string) string {
				return mk("DataFlow", "Show data flowing left-to-right through pipeline stages (e.g. Request → Auth → LLM → Storage → Response).",
					fmt.Sprintf(`CHOREOGRAPHY (150 frames):
- f0–20: stage boxes (props.stages[] = string[]) slide up in a row, staggered 10 frames.
- f20–140: a packet (small rounded square, props.color) travels stage→stage via interpolate on x, looping every 30 frames; each stage glows violet while the packet is inside it.
- f140–150: final stage gets a green check + success flash.
INPUT PROPS: title="Data Flow", stages[]=stages parsed from the sentence, sentence="%s", sentenceType=%q, module=%q, video=%q.
SENTENCE: %q`, clampLen(s, 90), st, m, v, s))
			},
		},
		{
			Slug: "code_typing", Label: "Code Typing", Emoji: "💻",
			BuildPrompt: func(s, st, m, v string) string {
				return mk("CodeTyping", "Reveal a code snippet with a typewriter effect + a blinking caret + simple syntax highlighting (keyword/string/comment colors).",
					fmt.Sprintf(`CHOREOGRAPHY (180 frames):
- Use a monospaced font (Fira Code / monospace). Show one extra character per frame (Math.min(frame, code.length)); a 2px violet caret blinks at the cursor position (visible when frame %% 30 < 15).
- Keyword tokens (#3b82f6), strings (#10b981), comments (#9ca3af), default (#f3f4f6).
- f150–180: after the snippet is fully typed, a violet underline sweeps under it + a caption (props.caption) fades in.
INPUT PROPS: code=props.code (the snippet to type), language=props.language, caption="%s", sentence="%s".
SENTENCE (context for caption): %q`, clampLen(s, 60), clampLen(s, 90), s))
			},
		},
		{
			Slug: "concept_reveal", Label: "Concept Reveal", Emoji: "💡",
			BuildPrompt: func(s, st, m, v string) string {
				return mk("ConceptReveal", "Kinetic typography: a single big idea punches in with scale + fade + blur-out exit (Steve Jobs keynote style).",
					fmt.Sprintf(`CHOREOGRAPHY (120 frames):
- f0–25: headline "%s" scales 1.4→1 with spring(bounce) + opacity 0→1 + blur 20px→0.
- f25–95: hold, a violet accent bar wipes left-to-right under the text (interpolate width 0→100%%).
- f95–120: blur 0→16px + scale 1→1.08 + opacity 1→0.
- One concept only; max 8 words in the headline. Background subtle radial violet glow.
INPUT PROPS: headline="%s", subheadline=optional, sentence="%s".
SENTENCE: %q`, clampLen(s, 64), clampLen(s, 64), clampLen(s, 90), s))
			},
		},
		{
			Slug: "timeline", Label: "Timeline", Emoji: "📅",
			BuildPrompt: func(s, st, m, v string) string {
				return mk("Timeline", "Animate a horizontal milestone timeline building left-to-right, one milestone at a time.",
					fmt.Sprintf(`CHOREOGRAPHY (150 frames):
- f0–20: a horizontal axis line draws across the screen (interpolate width 0→100%%).
- f20–140: each milestone (props.milestones[] = {at,label}) pops: a dot drops via spring, a label fades in above, staggered 18 frames.
- f140–150: the last milestone gets a violet ring pulse.
INPUT PROPS: title="%s", milestones[]=from sentence, sentence="%s", module=%q.
SENTENCE: %q`, clampLen(s, 60), clampLen(s, 90), m, s))
			},
		},
		{
			Slug: "comparison", Label: "Comparison", Emoji: "⚖️",
			BuildPrompt: func(s, st, m, v string) string {
				return mk("Comparison", "Split-screen comparison: two panels slide in from left/right and the differences highlight (e.g. Option A vs Option B, before vs after).",
					fmt.Sprintf(`CHOREOGRAPHY (150 frames):
- f0–40: left panel (props.left={title,points[]}) slides from -50%% + right panel (props.right) from +50%%; a center divider draws vertically.
- f40–130: points in each panel fade in staggered, one per side alternately; matching rows get a subtle violet link line.
- f130–150: the winning side (props.winner) lifts up + green border.
INPUT PROPS: left={title,points[]}, right={title,points[]}, winner="left"|"right"|null, sentence="%s".
SENTENCE (the comparison being made): %q`, clampLen(s, 90), s))
			},
		},
		{
			Slug: "process_steps", Label: "Process Steps", Emoji: "👣",
			BuildPrompt: func(s, st, m, v string) string {
				return mk("ProcessSteps", "Sequential numbered steps that animate in one-by-one with green checkmarks completing each.",
					fmt.Sprintf(`CHOREOGRAPHY (props.steps.length*24 + 30 frames):
- Each step card (props.steps[] = {n,title}) slides in from the right staggered 24 frames; a large step number (violet) counts up.
- 18 frames after a step appears, a green ✓ check scales in (spring) to mark it done.
- A progress bar at the bottom fills proportionally to completed steps.
INPUT PROPS: title="%s", steps[]=from sentence, sentence="%s".
SENTENCE (the procedure described): %q`, clampLen(s, 60), clampLen(s, 90), s))
			},
		},
		{
			Slug: "metric_counter", Label: "Metric Counter", Emoji: "🔢",
			BuildPrompt: func(s, st, m, v string) string {
				return mk("MetricCounter", "Animate a number counting up to a target value with a big reveal + supporting caption (cost, %, count, latency…).",
					fmt.Sprintf(`CHOREOGRAPHY (120 frames):
- f0–90: the number interpolates 0→props.target with an ease-out (use interpolate(frame,[0,90],[0,target],{extrapolateRight:'clamp',easing:Easing.out(Easing.cubic)})); format with props.prefix/suffix and props.decimals.
- f90–120: the value pulses once (scale 1→1.06→1) + a violet ring expands outward + caption "%s" fades in below.
INPUT PROPS: target=props.target (number), prefix, suffix, decimals, caption="%s", sentence="%s".
SENTENCE (the metric being stated): %q`, clampLen(s, 80), clampLen(s, 80), clampLen(s, 90), s))
			},
		},
		{
			Slug: "flowchart", Label: "Flowchart", Emoji: "🔀",
			BuildPrompt: func(s, st, m, v string) string {
				return mk("Flowchart", "Decision flowchart: rectangular process nodes + diamond decision nodes, branches revealing conditionally.",
					fmt.Sprintf(`CHOREOGRAPHY (180 frames):
- Build the graph top-down from props.nodes[] = {id,type:'process'|'decision'|'terminal',label,y} and props.edges[] = {from,to,label}.
- Nodes appear via spring staggered 18 frames; edges (SVG paths) draw after both endpoints exist (pathLength interpolate).
- At decision diamonds, the YES branch highlights green and NO branch highlights muted, 20 frames after the diamond appears.
INPUT PROPS: title="%s", nodes[]=parsed from sentence, edges[], sentence="%s".
SENTENCE (the logic/decision being mapped): %q`, clampLen(s, 60), clampLen(s, 90), s))
			},
		},
		{
			Slug: "callout_zoom", Label: "Callout Zoom", Emoji: "🔍",
			BuildPrompt: func(s, st, m, v string) string {
				return mk("CalloutZoom", "Zoom into a region of a screenshot/diagram and a callout label points at the highlighted element.",
					fmt.Sprintf(`CHOREOGRAPHY (120 frames):
- f0–40: full image (props.image) scales 1→props.zoom (e.g. 2.2) toward props.focusPoint {x,y} (transform-origin set to focus point), ease-in-out.
- f40–100: a violet rounded rectangle border pulses around the focus region; an SVG leader line draws from it to a callout pill (props.callout) that fades in.
- f100–120: hold; callout text typed/fully visible.
INPUT PROPS: image=props.image (URL), focusPoint={x,y} in 0–1, zoom=props.zoom||2, callout="%s", sentence="%s".
SENTENCE (what is being pointed out): %q`, clampLen(s, 70), clampLen(s, 90), s))
			},
		},
	}
}

// findAnimationType resolves a slug to its meta (case-insensitive). Returns nil
// when the slug is not one of the 10 known types.
func findAnimationType(slug string) *AnimationTypeMeta {
	slug = strings.ToLower(strings.TrimSpace(slug))
	for i := range animationTypeMetas() {
		m := animationTypeMetas()[i]
		if m.Slug == slug {
			return &m
		}
	}
	return nil
}

// suggestAnimationType maps a script sentence_type to a sensible default
// animation type so the page can pre-select the <select> for the user.
func suggestAnimationType(sentenceType string) string {
	switch strings.ToLower(strings.TrimSpace(sentenceType)) {
	case "hook", "heading", "title":
		return "concept_reveal"
	case "objective", "takeaway", "insight":
		return "callout_zoom"
	case "step", "transition":
		return "process_steps"
	case "cue":
		return "code_typing"
	default:
		return "concept_reveal"
	}
}

// animationDefaultProps returns the inputProps the Remotion <Composition>
// consumes for the given (type, sentence). This is what gets baked into the
// RunPod serverless render payload AND stored in remotion_props.
func animationDefaultProps(animType, sentence, sentenceType, moduleTitle, videoTitle string) map[string]any {
	base := map[string]any{
		"title":        clampLen(sentence, 64),
		"subtitle":     clampLen(moduleTitle+" · "+videoTitle, 80),
		"sentence":     sentence,
		"sentenceType": sentenceType,
		"module":       moduleTitle,
		"video":        videoTitle,
		"fps":           30,
		"durationInFrames": 150,
		"brandColor":   "#8b5cf6",
		"secondaryColor": "#3b82f6",
		"bgColor":      "#030712",
	}
	switch animType {
	case "code_typing":
		base["code"] = "// paste the snippet to type here\nconst result = await call({ model: 'claude' });"
		base["language"] = "javascript"
		base["durationInFrames"] = 180
	case "comparison":
		base["left"] = map[string]any{"title": "Option A", "points": []string{"point one", "point two"}}
		base["right"] = map[string]any{"title": "Option B", "points": []string{"point one", "point two"}}
		base["winner"] = "left"
	case "process_steps":
		base["steps"] = []map[string]any{{"n": 1, "title": "First step"}, {"n": 2, "title": "Second step"}, {"n": 3, "title": "Third step"}}
		base["durationInFrames"] = 102
	case "metric_counter":
		base["target"] = 100
		base["prefix"] = ""
		base["suffix"] = "%"
		base["decimals"] = 0
	case "timeline":
		base["milestones"] = []map[string]any{{"at": 0, "label": "Start"}, {"at": 50, "label": "Middle"}, {"at": 100, "label": "Now"}}
	case "data_flow":
		base["stages"] = []string{"Input", "Process", "Store", "Output"}
	case "architecture_diagram":
		base["nodes"] = []map[string]any{
			{"id": "a", "label": "Client", "x": 200, "y": 540},
			{"id": "b", "label": "API", "x": 700, "y": 540},
			{"id": "c", "label": "Database", "x": 1200, "y": 540},
		}
		base["activeNode"] = "b"
	case "flowchart":
		base["nodes"] = []map[string]any{
			{"id": "s", "type": "terminal", "label": "Start", "y": 200},
			{"id": "d", "type": "decision", "label": "Valid?", "y": 480},
			{"id": "e", "type": "process", "label": "Handle", "y": 760},
		}
		base["edges"] = []map[string]any{{"from": "s", "to": "d"}, {"from": "d", "to": "e", "label": "YES"}}
	case "callout_zoom":
		base["image"] = ""
		base["focusPoint"] = map[string]any{"x": 0.5, "y": 0.5}
		base["zoom"] = 2
	}
	return base
}

// animationPromptHandler builds the Remotion composition prompt + inputProps +
// the full RunPod serverless payload for a (sentence, animationType). Open: it
// returns only prompt text + non-secret config, never the API key.
func animationPromptHandler(cfg config) http.HandlerFunc {
	type Req struct {
		Sentence      string `json:"sentence"`
		SentenceType  string `json:"sentenceType"`
		AnimationType string `json:"animationType"`
		ModuleName    string `json:"moduleName"`
		VideoName     string `json:"videoName"`
		CustomPrompt  string `json:"customPrompt"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		setCORS(w, r)
		w.Header().Set("Content-Type", "application/json")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		if r.Method != http.MethodPost {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}
		var req Req
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
			return
		}
		meta := findAnimationType(req.AnimationType)
		if meta == nil {
			http.Error(w, `{"error":"unknown animation type"}`, http.StatusBadRequest)
			return
		}
		sentence := strings.TrimSpace(req.Sentence)
		if sentence == "" {
			sentence = "(no sentence provided)"
		}
		prompt := strings.TrimSpace(req.CustomPrompt)
		if prompt == "" {
			prompt = meta.BuildPrompt(sentence, req.SentenceType, req.ModuleName, req.VideoName)
		}
		props := animationDefaultProps(meta.Slug, sentence, req.SentenceType, req.ModuleName, req.VideoName)

		// Build the RunPod serverless payload (the exact body POSTed to
		// /v2/{endpoint}/run). serveUrl is filled at submit time from env so the
		// secret-free preview is still useful to copy/paste into RunPod manually.
		serveURL := os.Getenv("REMOTION_SERVE_URL")
		dur, _ := props["durationInFrames"].(int)
		if dur == 0 {
			dur = 150
		}
		runpodPayload := map[string]any{
			"input": map[string]any{
				"serveUrl":         serveURL,
				"composition":      "Main",
				"codec":            "h264",
				"imageFormat":      "jpeg",
				"crf":              18,
				"inputProps":       props,
				"durationInFrames": dur,
				"fps":              30,
				"width":            1920,
				"height":           1080,
			},
		}

		json.NewEncoder(w).Encode(map[string]any{
			"animationType":  meta.Slug,
			"label":          meta.Label,
			"emoji":          meta.Emoji,
			"prompt":         prompt,
			"inputProps":     props,
			"runpodPayload":  runpodPayload,
			"serveUrl":       serveURL,
			"runpodConfigured": cfg.getSecret("RUNPOD_API_KEY") != "" && os.Getenv("RUNPOD_ENDPOINT_ID") != "" && serveURL != "",
		})
	}
}

func scriptGenerateHandler(cfg config) http.HandlerFunc {
	type ORRequest struct {
		Prompt   string `json:"prompt"`
		JsonMode bool   `json:"jsonMode"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodPost {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		var req ORRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
			return
		}

		apiKey := cfg.getSecret("OPENROUTER_API_KEY")
		if apiKey == "" {
			http.Error(w, `{"error":"OPENROUTER_API_KEY missing from server configuration"}`, http.StatusServiceUnavailable)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()

		orURL := "https://openrouter.ai/api/v1/chat/completions"
		reqBody := map[string]any{
			"model": "google/gemini-2.5-flash",
			"messages": []map[string]string{
				{"role": "user", "content": req.Prompt},
			},
		}
		if req.JsonMode {
			reqBody["response_format"] = map[string]map[string]string{
				"type": {"type": "json_object"},
			}
		}

		b, _ := json.Marshal(reqBody)
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
			http.Error(w, `{"error":"OpenRouter API call failed"}`, http.StatusBadGateway)
			return
		}
		defer hresp.Body.Close()

		if hresp.StatusCode >= 400 {
			body, _ := io.ReadAll(hresp.Body)
			http.Error(w, fmt.Sprintf(`{"error":"OpenRouter API returned HTTP %d: %s"}`, hresp.StatusCode, body), http.StatusBadGateway)
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

		content := ""
		if len(orResp.Choices) > 0 {
			content = orResp.Choices[0].Message.Content
		}

		json.NewEncoder(w).Encode(map[string]string{
			"content": content,
		})
	}
}

// animationOpenRouterPrepareHandler calls OpenRouter to intelligently populate
// the inputProps for a given sentence and animation type, returning JSON.
func animationOpenRouterPrepareHandler(cfg config) http.HandlerFunc {
	type Req struct {
		AnimationType string `json:"animationType"`
		Sentence      string `json:"sentence"`
		SentenceType  string `json:"sentenceType"`
		ModuleName    string `json:"moduleName"`
		VideoName     string `json:"videoName"`
	}
	const model = "anthropic/claude-sonnet-4.6"
	return func(w http.ResponseWriter, r *http.Request) {
		setCORS(w, r)
		w.Header().Set("Content-Type", "application/json")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		if !isAdminRequest(r) {
			http.Error(w, `{"error":"Unauthorized"}`, http.StatusUnauthorized)
			return
		}
		if r.Method != http.MethodPost {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		apiKey := cfg.getSecret("OPENROUTER_API_KEY")
		if apiKey == "" {
			http.Error(w, `{"error":"OPENROUTER_API_KEY missing from server configuration"}`, http.StatusServiceUnavailable)
			return
		}

		var req Req
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
			return
		}

		meta := findAnimationType(req.AnimationType)
		if meta == nil {
			http.Error(w, `{"error":"unknown animation type"}`, http.StatusBadRequest)
			return
		}

		defaultProps := animationDefaultProps(req.AnimationType, req.Sentence, req.SentenceType, req.ModuleName, req.VideoName)
		defaultPropsJSON, _ := json.MarshalIndent(defaultProps, "", "  ")

		prompt := fmt.Sprintf(`You are an expert data extractor for video animations.
I will give you a sentence from a technical course script and an animation type.
I need you to output the JSON properties ("inputProps") that will be passed to the Remotion video renderer.

Animation Type: %s (%s)
Sentence: %q

Instructions:
1. Return ONLY a valid JSON object. No markdown formatting, no prose, no code blocks (do not wrap in ` + "`" + `json` + "`" + `).
2. The JSON schema must exactly match the keys and data types of this default example:
%s
3. Intelligently update the values (e.g. titles, points, steps, nodes, code snippets) to match the meaning of the Sentence.
4. Keep strings concise. Do not change colors, fps, or durationInFrames unless absolutely necessary for the content to fit.`, meta.Label, meta.Slug, req.Sentence, string(defaultPropsJSON))

		payload := map[string]any{
			"model":       model,
			"messages":    []map[string]string{{"role": "user", "content": prompt}},
			"temperature": 0.2,
		}
		body, _ := json.Marshal(payload)
		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()
		hreq, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://openrouter.ai/api/v1/chat/completions", bytes.NewReader(body))
		if err != nil {
			http.Error(w, `{"error":"request build failed"}`, http.StatusInternalServerError)
			return
		}
		hreq.Header.Set("Content-Type", "application/json")
		hreq.Header.Set("Authorization", "Bearer "+apiKey)
		hreq.Header.Set("HTTP-Referer", "https://github.com/rifaterdemsahin/claude-architect-certification")
		hreq.Header.Set("X-Title", "Claude Architect Certification")

		hresp, err := http.DefaultClient.Do(hreq)
		if err != nil {
			http.Error(w, `{"error":"OpenRouter call failed"}`, http.StatusBadGateway)
			return
		}
		defer hresp.Body.Close()
		b, _ := io.ReadAll(hresp.Body)
		if hresp.StatusCode != http.StatusOK {
			http.Error(w, fmt.Sprintf(`{"error":"OpenRouter HTTP %d: %s"}`, hresp.StatusCode, b), http.StatusBadGateway)
			return
		}

		var orResp struct {
			Choices []struct {
				Message struct {
					Content string `json:"content"`
				} `json:"message"`
			} `json:"choices"`
		}
		if err := json.Unmarshal(b, &orResp); err != nil || len(orResp.Choices) == 0 {
			http.Error(w, `{"error":"invalid OpenRouter response"}`, http.StatusBadGateway)
			return
		}
		
		content := strings.TrimSpace(orResp.Choices[0].Message.Content)
		content = strings.TrimPrefix(content, "```json")
		content = strings.TrimPrefix(content, "```")
		content = strings.TrimSuffix(content, "```")
		content = strings.TrimSpace(content)

		var parsed map[string]any
		if err := json.Unmarshal([]byte(content), &parsed); err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"OpenRouter did not return valid JSON. Content was: %s"}`, content), http.StatusBadGateway)
			return
		}

		json.NewEncoder(w).Encode(parsed)
	}
}

// animationRunpodRunHandler submits a render job to RunPod serverless and
// records a generating row in sentence_animations. Admin-gated (paid render).
func animationRunpodRunHandler(cfg config) http.HandlerFunc {
	type Req struct {
		SentenceID    int    `json:"sentenceId"`
		ModuleNumber  int    `json:"moduleNumber"`
		VideoNumber   int    `json:"videoNumber"`
		ScriptID      int    `json:"scriptId"`
		AnimationType string `json:"animationType"`
		Prompt        string `json:"prompt"`
		InputProps    map[string]any `json:"inputProps"`
		CustomPrompt  string `json:"customPrompt"`
		Duration      int    `json:"durationInFrames"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		setCORS(w, r)
		w.Header().Set("Content-Type", "application/json")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		if !isAdminRequest(r) {
			http.Error(w, `{"error":"Unauthorized — sign in as admin or open from localhost."}`, http.StatusUnauthorized)
			return
		}
		if r.Method != http.MethodPost {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}
		apiKey := cfg.getSecret("RUNPOD_API_KEY")
		endpointID := os.Getenv("RUNPOD_ENDPOINT_ID")
		serveURL := os.Getenv("REMOTION_SERVE_URL")
		if apiKey == "" || endpointID == "" || serveURL == "" {
			http.Error(w, `{"error":"RunPod not configured. Set RUNPOD_API_KEY (Key Vault/Azure secret 'runpod-api-key' or env), RUNPOD_ENDPOINT_ID, and REMOTION_SERVE_URL on the server."}`, http.StatusServiceUnavailable)
			return
		}
		var req Req
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
			return
		}
		if req.SentenceID == 0 {
			http.Error(w, `{"error":"sentenceId is required"}`, http.StatusBadRequest)
			return
		}
		if findAnimationType(req.AnimationType) == nil {
			http.Error(w, `{"error":"unknown animation type"}`, http.StatusBadRequest)
			return
		}
		if req.InputProps == nil {
			req.InputProps = animationDefaultProps(req.AnimationType, "", "", "", "")
		}
		dur := req.Duration
		if dur == 0 {
			if v, ok := req.InputProps["durationInFrames"].(float64); ok {
				dur = int(v)
			}
			if v, ok := req.InputProps["durationInFrames"].(int); ok {
				dur = v
			}
		}
		if dur == 0 {
			dur = 150
		}
		req.InputProps["durationInFrames"] = dur

		payload := map[string]any{
			"input": map[string]any{
				"serveUrl":         serveURL,
				"composition":      "Main",
				"codec":            "h264",
				"imageFormat":      "jpeg",
				"crf":              18,
				"inputProps":       req.InputProps,
				"durationInFrames": dur,
				"fps":              30,
				"width":            1920,
				"height":           1080,
			},
		}
		body, _ := json.Marshal(payload)

		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()
		runURL := "https://api.runpod.ai/v2/" + endpointID + "/run"
		hreq, err := http.NewRequestWithContext(ctx, http.MethodPost, runURL, bytes.NewReader(body))
		if err != nil {
			http.Error(w, `{"error":"failed to build runpod request"}`, http.StatusInternalServerError)
			return
		}
		hreq.Header.Set("Content-Type", "application/json")
		hreq.Header.Set("Authorization", "Bearer "+apiKey)

		hresp, err := http.DefaultClient.Do(hreq)
		if err != nil {
			http.Error(w, `{"error":"RunPod submit failed"}`, http.StatusBadGateway)
			return
		}
		defer hresp.Body.Close()
		rb, _ := io.ReadAll(hresp.Body)
		if hresp.StatusCode >= 400 {
			http.Error(w, fmt.Sprintf(`{"error":"RunPod submit HTTP %d: %s"}`, hresp.StatusCode, strings.TrimSpace(string(rb))), http.StatusBadGateway)
			return
		}
		var rpResp struct {
			ID     string `json:"id"`
			Status string `json:"status"`
		}
		if err := json.Unmarshal(rb, &rpResp); err != nil || rpResp.ID == "" {
			http.Error(w, `{"error":"RunPod response missing job id"}`, http.StatusBadGateway)
			return
		}

		// Upsert a generating row keyed on (sentence_id, animation_type).
		row := map[string]any{
			"sentence_id":     req.SentenceID,
			"module_number":   req.ModuleNumber,
			"video_number":    req.VideoNumber,
			"script_id":       req.ScriptID,
			"animation_type":  req.AnimationType,
			"generation_status": "generating",
			"prompt_used":     strings.TrimSpace(req.CustomPrompt),
			"remotion_props":  req.InputProps,
			"runpod_job_id":   rpResp.ID,
			"runpod_status":   rpResp.Status,
			"codec":           "h264",
			"duration_in_frames": dur,
			"fps":             30,
			"width":           1920,
			"height":          1080,
			"error_message":   "",
		}
		if strings.TrimSpace(req.CustomPrompt) == "" {
			row["prompt_used"] = req.Prompt
		}
		// Upsert via the unique (sentence_id, animation_type) index.
		if err := supabaseUpsert(ctx, cfg, "sentence_animations", row, "sentence_id,animation_type"); err != nil {
			log.Printf("sentence_animations upsert: %v", err)
		}

		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(map[string]any{
			"jobId":     rpResp.ID,
			"status":    rpResp.Status,
			"endpoint":  endpointID,
		})
	}
}

// animationRunpodStatusHandler polls a RunPod job. On COMPLETED it downloads
// the rendered MP4, uploads it to Azure `research-animations`, and patches the
// sentence_animations row. Idempotent: once a row is completed, re-polling is a
// no-op. Admin-gated (it performs the Azure write).
func animationRunpodStatusHandler(cfg config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setCORS(w, r)
		w.Header().Set("Content-Type", "application/json")
		if !isAdminRequest(r) {
			http.Error(w, `{"error":"Unauthorized — sign in as admin or open from localhost."}`, http.StatusUnauthorized)
			return
		}
		jobID := r.URL.Query().Get("id")
		sentenceID := r.URL.Query().Get("sentence_id")
		if jobID == "" || sentenceID == "" {
			http.Error(w, `{"error":"id and sentence_id are required"}`, http.StatusBadRequest)
			return
		}
		apiKey := cfg.getSecret("RUNPOD_API_KEY")
		endpointID := os.Getenv("RUNPOD_ENDPOINT_ID")
		if apiKey == "" || endpointID == "" {
			http.Error(w, `{"error":"RunPod not configured"}`, http.StatusServiceUnavailable)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()
		statusURL := "https://api.runpod.ai/v2/" + endpointID + "/status/" + jobID
		hreq, _ := http.NewRequestWithContext(ctx, http.MethodGet, statusURL, nil)
		hreq.Header.Set("Authorization", "Bearer "+apiKey)
		hresp, err := http.DefaultClient.Do(hreq)
		if err != nil {
			http.Error(w, `{"error":"RunPod status poll failed"}`, http.StatusBadGateway)
			return
		}
		defer hresp.Body.Close()
		rb, _ := io.ReadAll(hresp.Body)
		var rpResp struct {
			Status string `json:"status"`
			Output any    `json:"output"`
			Error  string `json:"error"`
		}
		_ = json.Unmarshal(rb, &rpResp)

		status := strings.ToUpper(strings.TrimSpace(rpResp.Status))
		out := map[string]any{"status": status, "runpodStatus": status}

		switch status {
		case "COMPLETED":
			videoURL := runPodOutputVideoURL(rpResp.Output)
			if videoURL == "" {
				out["status"] = "COMPLETED_NO_URL"
				out["error"] = "render completed but no output video URL was returned"
				if e := supabasePatch(ctx, cfg, "sentence_animations",
					"runpod_job_id=eq."+jobID,
					map[string]any{"generation_status": "failed", "error_message": out["error"].(string), "runpod_status": status}); e != nil {
					log.Printf("sentence_animations patch failed(no-url): %v", e)
				}
				json.NewEncoder(w).Encode(out)
				return
			}
			// Download the rendered MP4 from RunPod.
			dl, derr := http.Get(videoURL)
			if derr != nil {
				out["error"] = "failed to download render: " + derr.Error()
				json.NewEncoder(w).Encode(out)
				return
			}
			videoBytes, _ := io.ReadAll(dl.Body)
			dl.Body.Close()
			if len(videoBytes) == 0 {
				out["error"] = "rendered video was empty"
				json.NewEncoder(w).Encode(out)
				return
			}
			// Upload to Azure `research-animations`.
			blobName := fmt.Sprintf("m%d_v%d_%d.mp4", queryInt(r, "module_number"), queryInt(r, "video_number"), time.Now().Unix())
			if uerr := uploadBlobToAzure(ctx, cfg, "research-animations", blobName, "video/mp4", videoBytes); uerr != nil {
				log.Printf("animation azure upload failed: %v", uerr)
				out["error"] = "azure upload failed: " + uerr.Error()
				json.NewEncoder(w).Encode(out)
				return
			}
			proxyURL := researchFileProxyURL("research-animations", blobName)
			if e := supabasePatch(ctx, cfg, "sentence_animations",
				"runpod_job_id=eq."+jobID,
				map[string]any{
					"generation_status": "completed",
					"runpod_status":    status,
					"azure_blob_name":  blobName,
					"animation_url":    proxyURL,
					"error_message":    "",
				}); e != nil {
				log.Printf("sentence_animations patch failed(completed): %v", e)
			}
			out["url"] = proxyURL
			out["blobName"] = blobName
		case "FAILED", "CANCELLED", "TIMED_OUT":
			errMsg := rpResp.Error
			if errMsg == "" {
				errMsg = "render " + strings.ToLower(status)
			}
			if e := supabasePatch(ctx, cfg, "sentence_animations",
				"runpod_job_id=eq."+jobID,
				map[string]any{"generation_status": "failed", "runpod_status": status, "error_message": errMsg}); e != nil {
				log.Printf("sentence_animations patch failed(%s): %v", status, e)
			}
			out["error"] = errMsg
		}
		json.NewEncoder(w).Encode(out)
	}
}

// queryInt reads an int query param, 0 when missing/invalid.
func queryInt(r *http.Request, key string) int {
	v, _ := strconv.Atoi(r.URL.Query().Get(key))
	return v
}

// runPodOutputVideoURL tolerates the output shapes different RunPod Remotion
// workers return: {url}, {output}, {video}, {url, ...}, or a bare string.
func runPodOutputVideoURL(output any) string {
	switch o := output.(type) {
	case string:
		return o
	case map[string]any:
		for _, k := range []string{"url", "output", "video", "file"} {
			if s, ok := o[k].(string); ok && s != "" {
				return s
			}
		}
	}
	return ""
}

// runPodLLMText extracts the generated text from a RunPod serverless LLM
// endpoint (vLLM/TGI, OpenAI-compatible) completion response. The endpoint can
// return the OpenAI shape (output is an array of {choices:[{message:{content}}]})
// OR a single object. We tolerate both, plus a bare string and the legacy
// {choices:[{text}]} format.
func runPodLLMText(output any) string {
	switch o := output.(type) {
	case string:
		return o
	case map[string]any:
		if arr, ok := o["choices"].([]any); ok && len(arr) > 0 {
			return firstChoiceText(arr[0])
		}
		if data, ok := o["data"].([]any); ok && len(data) > 0 {
			return firstChoiceText(data[0])
		}
		if s, ok := o["content"].(string); ok {
			return s
		}
	case []any:
		if len(o) > 0 {
			return firstChoiceText(o[0])
		}
	}
	return ""
}

func firstChoiceText(v any) string {
	m, ok := v.(map[string]any)
	if !ok {
		return ""
	}
	if msg, ok := m["message"].(map[string]any); ok {
		if s, ok := msg["content"].(string); ok {
			return s
		}
	}
	if s, ok := m["text"].(string); ok {
		return s
	}
	if arr, ok := m["choices"].([]any); ok && len(arr) > 0 {
		return firstChoiceText(arr[0])
	}
	return ""
}

// animationRunpodGenerateCodeHandler uses the RunPod serverless LLM endpoint
// (the existing endpoint rsplhkl473fnsa is an OpenAI-compatible text model) to
// turn a (sentence, animationType) into a complete, deployable Remotion
// <Composition> source. The server builds the Remotion prompt (reusing the
// same builder the preview uses), sends it to the LLM as a chat message,
// polls until COMPLETED, and returns the generated React/Remotion code.
//
// Why this path: the RunPod key in the Azure Key Vault (runpod-api-key) backs
// an LLM endpoint today, not a Remotion render farm. Generating the Remotion
// component code is the productive, real use of that key for the Animation
// Generator — the code it returns is exactly what a Remotion serve URL would
// bundle. Admin-gated (paid LLM call). The key never reaches the browser.
func animationRunpodGenerateCodeHandler(cfg config) http.HandlerFunc {
	type Req struct {
		AnimationType string `json:"animationType"`
		Sentence      string `json:"sentence"`
		SentenceType  string `json:"sentenceType"`
		ModuleName    string `json:"moduleName"`
		VideoName     string `json:"videoName"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		setCORS(w, r)
		w.Header().Set("Content-Type", "application/json")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		if !isAdminRequest(r) {
			http.Error(w, `{"error":"Unauthorized — sign in as admin or open from localhost."}`, http.StatusUnauthorized)
			return
		}
		if r.Method != http.MethodPost {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}
		apiKey := cfg.getSecret("RUNPOD_API_KEY")
		endpointID := os.Getenv("RUNPOD_ENDPOINT_ID")
		if apiKey == "" || endpointID == "" {
			http.Error(w, `{"error":"RunPod not configured. Set RUNPOD_API_KEY (Key Vault 'runpod-api-key' or env) and RUNPOD_ENDPOINT_ID on the server."}`, http.StatusServiceUnavailable)
			return
		}
		var req Req
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
			return
		}
		meta := findAnimationType(req.AnimationType)
		if meta == nil {
			http.Error(w, `{"error":"unknown animation type"}`, http.StatusBadRequest)
			return
		}
		sentence := strings.TrimSpace(req.Sentence)
		if sentence == "" {
			sentence = meta.Label + " demo for the Claude Architect Certification course."
		}
		prompt := meta.BuildPrompt(sentence, req.SentenceType, req.ModuleName, req.VideoName)
		props := animationDefaultProps(meta.Slug, sentence, req.SentenceType, req.ModuleName, req.VideoName)

		// Submit the prompt to the RunPod serverless LLM endpoint in OpenAI
		// chat format (the endpoint validated `messages` as required).
		payload := map[string]any{
			"input": map[string]any{
				"messages": []map[string]string{
					{"role": "user", "content": prompt},
				},
				"max_tokens":  4000,
				"temperature": 0.3,
			},
		}
		body, _ := json.Marshal(payload)
		submitCtx, submitCancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer submitCancel()
		runURL := "https://api.runpod.ai/v2/" + endpointID + "/run"
		hreq, err := http.NewRequestWithContext(submitCtx, http.MethodPost, runURL, bytes.NewReader(body))
		if err != nil {
			http.Error(w, `{"error":"failed to build runpod request"}`, http.StatusInternalServerError)
			return
		}
		hreq.Header.Set("Content-Type", "application/json")
		hreq.Header.Set("Authorization", "Bearer "+apiKey)
		hresp, err := http.DefaultClient.Do(hreq)
		if err != nil {
			http.Error(w, `{"error":"RunPod submit failed"}`, http.StatusBadGateway)
			return
		}
		rb, _ := io.ReadAll(hresp.Body)
		hresp.Body.Close()
		if hresp.StatusCode >= 400 {
			http.Error(w, fmt.Sprintf(`{"error":"RunPod submit HTTP %d: %s"}`, hresp.StatusCode, strings.TrimSpace(string(rb))), http.StatusBadGateway)
			return
		}
		var rpRun struct {
			ID     string `json:"id"`
			Status string `json:"status"`
		}
		if err := json.Unmarshal(rb, &rpRun); err != nil || rpRun.ID == "" {
			http.Error(w, `{"error":"RunPod response missing job id"}`, http.StatusBadGateway)
			return
		}

		// Poll for completion. LLM generation of a full composition can take
		// 20-60s on a cold worker; allow up to ~120s total.
		jobID := rpRun.ID
		deadline := time.Now().Add(120 * time.Second)
		for time.Now().Before(deadline) {
			pollCtx, pollCancel := context.WithTimeout(r.Context(), 30*time.Second)
			statusURL := "https://api.runpod.ai/v2/" + endpointID + "/status/" + jobID
			preq, _ := http.NewRequestWithContext(pollCtx, http.MethodGet, statusURL, nil)
			preq.Header.Set("Authorization", "Bearer "+apiKey)
			presp, perr := http.DefaultClient.Do(preq)
			if perr != nil {
				pollCancel()
				time.Sleep(3 * time.Second)
				continue
			}
			pb, _ := io.ReadAll(presp.Body)
			presp.Body.Close()
			pollCancel()
			var st struct {
				Status string `json:"status"`
				Output any    `json:"output"`
				Error  string `json:"error"`
			}
			_ = json.Unmarshal(pb, &st)
			status := strings.ToUpper(strings.TrimSpace(st.Status))
			switch status {
			case "COMPLETED":
				code := strings.TrimSpace(runPodLLMText(st.Output))
				if code == "" {
					http.Error(w, `{"error":"RunPod completed but returned no generated text"}`, http.StatusBadGateway)
					return
				}
				json.NewEncoder(w).Encode(map[string]any{
					"animationType": meta.Slug,
					"label":         meta.Label,
					"emoji":         meta.Emoji,
					"prompt":        prompt,
					"inputProps":    props,
					"code":          code,
					"jobId":         jobID,
					"endpoint":      endpointID,
				})
				return
			case "FAILED", "CANCELLED", "TIMED_OUT":
				errMsg := st.Error
				if errMsg == "" {
					errMsg = "RunPod LLM job " + strings.ToLower(status)
				}
				http.Error(w, fmt.Sprintf(`{"error":"RunPod LLM failed: %s"}`, errMsg), http.StatusBadGateway)
				return
			}
			time.Sleep(3 * time.Second)
		}
		http.Error(w, `{"error":"RunPod LLM timed out (120s)"}`, http.StatusGatewayTimeout)
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

func fixGrammarHandler(cfg config) http.HandlerFunc {
	type FixRequest struct {
		Text string `json:"text"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodPost {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		var req FixRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
			return
		}

		if strings.TrimSpace(req.Text) == "" {
			json.NewEncoder(w).Encode(map[string]string{"fixed_text": ""})
			return
		}

		geminiKey := cfg.getSecret("GEMINI_API_KEY")
		if geminiKey == "" {
			http.Error(w, `{"error":"GEMINI_API_KEY missing"}`, http.StatusServiceUnavailable)
			return
		}

		prompt := fmt.Sprintf(`You are an expert editor and technical writer. 
Fix any spelling, grammar, and punctuation errors in the following text. 
Maintain the original meaning, tone, and formatting (bullet points, etc.). 
Return ONLY the corrected text, no preamble or extra commentary.

TEXT TO FIX:
%s`, req.Text)

		ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
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

		fixed := ""
		if len(gResp.Candidates) > 0 && len(gResp.Candidates[0].Content.Parts) > 0 {
			fixed = gResp.Candidates[0].Content.Parts[0].Text
		}

		json.NewEncoder(w).Encode(map[string]string{
			"fixed_text": strings.TrimSpace(fixed),
		})
	}
}

// ── Admin Handlers ────────────────────────────────────────────────────────────

func adminBackupSupabaseHandler(cfg config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !isAdminRequest(r) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		dbURL := cfg.getSecret("SUPABASE_DB_URL")
		if dbURL == "" {
			dbURL = os.Getenv("SUPABASE_DB_URL")
		}
		if dbURL == "" {
			http.Error(w, "SUPABASE_DB_URL not configured in Azure Key Vault or Env", http.StatusInternalServerError)
			return
		}

		cmd := exec.Command("pg_dump", dbURL, "--no-owner", "--no-acl", "--schema=public")
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", "attachment; filename=\"supabase_backup.sql\"")
		
		cmd.Stdout = w
		cmd.Stderr = os.Stderr
		
		if err := cmd.Run(); err != nil {
			log.Printf("pg_dump failed: %v", err)
		}
	}
}

func adminBackupAzureHandler(cfg config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !isAdminRequest(r) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		secrets := map[string]string{
			"SUPABASE_DB_URL": cfg.getSecret("SUPABASE_DB_URL"),
			"OPENROUTER_API_KEY": cfg.getSecret("OPENROUTER_API_KEY"),
			"GEMINI_API_KEY": cfg.getSecret("GEMINI_API_KEY"),
			"claude-architect-GOOGLE-CLIENT-ID": cfg.getSecret("claude-architect-GOOGLE-CLIENT-ID"),
			"GOOGLE-IMAGEN-API-KEY": cfg.getSecret("GOOGLE-IMAGEN-API-KEY"),
			"AXIOM_TOKEN": cfg.getSecret("AXIOM_TOKEN"),
			"AXIOM_ORG_ID": cfg.getSecret("AXIOM_ORG_ID"),
			"AdminPassword": cfg.getSecret("AdminPassword"),
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Disposition", "attachment; filename=\"azure_secrets_backup.json\"")
		json.NewEncoder(w).Encode(secrets)
	}
}

// adminStatusHandler lets any page ask "am I signed in?" so destructive UI
// (delete/upload buttons) can be hidden when the visitor is not trusted. The
// rule mirrors every privileged handler: localhost origin OR admin cookie.
func adminStatusHandler(cfg config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-store")
		json.NewEncoder(w).Encode(map[string]any{"admin": isAdminRequest(r)})
	}
}

func adminLoginHandler(cfg config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req struct {
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		adminPass := cfg.getSecret("ClaudeCertificateSiteAdminPassword")
		if adminPass == "" {
			adminPass = cfg.getSecret("AdminPassword")
		}
		if adminPass == "" {
			adminPass = "admin"
		}

		if req.Password != adminPass {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"success":false,"error":"Invalid password"}`))
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "admin_logged_in",
			Value:    "true",
			Path:     "/",
			MaxAge:   3600 * 24,
			HttpOnly: false,
			SameSite: http.SameSiteLaxMode,
		})

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"success":true}`))
	}
}

// adminLogoutHandler clears the `admin_logged_in` cookie so the visitor is no
// longer trusted: /api/admin/status then returns {"admin":false} and every
// destructive button hides again (shared/nav.js gate). Mirrors the login
// handler's cookie shape so the browser drops exactly the same cookie.
func adminLogoutHandler(cfg config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "admin_logged_in",
			Value:    "",
			Path:     "/",
			MaxAge:   -1, // delete now
			Expires:  time.Unix(0, 0),
			HttpOnly: false,
			SameSite: http.SameSiteLaxMode,
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"success":true}`))
	}
}

func adminGDriveCredentialsHandler(cfg config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !isAdminRequest(r) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error":"Unauthorized"}`))
			return
		}

		clientID := cfg.getSecret("claude-architect-GOOGLE-CLIENT-ID")
		if clientID == "" {
			clientID = cfg.getSecret("google-oauth-client-id")
		}
		if clientID == "" {
			clientID = cfg.getSecret("GOOGLE_CLIENT_ID")
		}

		apiKey := cfg.getSecret("GOOGLE-IMAGEN-API-KEY")
		if apiKey == "" {
			apiKey = cfg.getSecret("GOOGLE-SEARCH-API-KEY")
		}
		if apiKey == "" {
			apiKey = cfg.getSecret("GOOGLE_API_KEY")
		}

		// .env is only ever written from a trusted (localhost) origin; the handler
		// already requires isAdminRequest, so save unconditionally here.
		saveCredentialsToDotEnv(clientID, apiKey)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"clientId": clientID,
			"apiKey":   apiKey,
		})
	}
}

func saveCredentialsToDotEnv(clientId, apiKey string) {
	if clientId == "" || apiKey == "" {
		return
	}
	content, err := os.ReadFile(".env")
	if err != nil {
		content = []byte{}
	}
	lines := strings.Split(string(content), "\n")
	hasClientID := false
	hasAPIKey := false
	for i, line := range lines {
		if strings.HasPrefix(line, "GOOGLE_CLIENT_ID=") {
			lines[i] = "GOOGLE_CLIENT_ID=" + clientId
			hasClientID = true
		}
		if strings.HasPrefix(line, "GOOGLE_API_KEY=") {
			lines[i] = "GOOGLE_API_KEY=" + apiKey
			hasAPIKey = true
		}
	}
	if !hasClientID {
		lines = append(lines, "GOOGLE_CLIENT_ID="+clientId)
	}
	if !hasAPIKey {
		lines = append(lines, "GOOGLE_API_KEY="+apiKey)
	}
	os.WriteFile(".env", []byte(strings.Join(lines, "\n")), 0644)
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
	mux.Handle("/api/env-status", observe(cfg, envStatusHandler(cfg)))
	mux.Handle("/api/nav/favs", observe(cfg, navFavsHandler(cfg)))
	mux.Handle("/api/errors", observe(cfg, clientErrorsHandler(cfg)))
	mux.Handle("/api/research/upload", observe(cfg, researchUploadHandler(cfg)))
	mux.Handle("/api/research/files", observe(cfg, researchFilesHandler(cfg)))
	mux.Handle("/api/research/file", observe(cfg, researchFileHandler(cfg)))
	mux.Handle("/api/research/rename", observe(cfg, researchRenameHandler(cfg)))
	mux.Handle("/api/explanations/generate", observe(cfg, generateExplanationHandler(cfg)))
	mux.Handle("/api/explanations", observe(cfg, explanationsHandler(cfg)))
	mux.Handle("/api/images/generate", observe(cfg, imageGenerateHandler(cfg)))
	mux.Handle("/api/images/enhance-prompt", observe(cfg, imageEnhancePromptHandler(cfg)))
	mux.Handle("/api/images/save", observe(cfg, imageSaveHandler(cfg)))
	mux.Handle("/api/images/backfill-thumbnails", observe(cfg, imageThumbnailBackfillHandler(cfg)))
	mux.Handle("/api/images/test-gemini", observe(cfg, imageTestGeminiHandler(cfg)))
	mux.Handle("/api/infographics/generate", observe(cfg, infographicGenerateHandler(cfg)))
	mux.Handle("/api/infographics/save", observe(cfg, infographicSaveHandler(cfg)))
	mux.Handle("/api/scripts/openrouter", observe(cfg, scriptGenerateHandler(cfg)))
	mux.Handle("/api/lowerthirds/generate", observe(cfg, lowerThirdGenerateHandler(cfg)))
	mux.Handle("/api/lowerthirds/openrouter", observe(cfg, openRouterGenerateHandler(cfg)))
	mux.Handle("/api/drawings/openrouter", observe(cfg, drawingGenerateHandler(cfg)))
	mux.Handle("/api/drawings/save", observe(cfg, drawingSaveHandler(cfg)))
	mux.Handle("/api/slides/openrouter", observe(cfg, slideGenerateHandler(cfg)))
	mux.Handle("/api/animations/generate-prompt", observe(cfg, animationPromptHandler(cfg)))
	mux.Handle("/api/animations/openrouter-prepare", observe(cfg, animationOpenRouterPrepareHandler(cfg)))
	mux.Handle("/api/animations/runpod/run", observe(cfg, animationRunpodRunHandler(cfg)))
	mux.Handle("/api/animations/runpod/status", observe(cfg, animationRunpodStatusHandler(cfg)))
	mux.Handle("/api/animations/runpod/generate-code", observe(cfg, animationRunpodGenerateCodeHandler(cfg)))
	mux.Handle("/api/ai/sanity-check", observe(cfg, sanityCheckHandler(cfg)))
	mux.Handle("/api/ai/fix-grammar", observe(cfg, fixGrammarHandler(cfg)))
	mux.Handle("/api/admin/login", observe(cfg, adminLoginHandler(cfg)))
	mux.Handle("/api/admin/logout", observe(cfg, adminLogoutHandler(cfg)))
	mux.Handle("/api/admin/status", observe(cfg, adminStatusHandler(cfg)))
	mux.Handle("/api/admin/backup/supabase", observe(cfg, adminBackupSupabaseHandler(cfg)))
	mux.Handle("/api/admin/backup/azure", observe(cfg, adminBackupAzureHandler(cfg)))
	mux.Handle("/api/axiom/logs", observe(cfg, axiomLogsHandler(cfg)))
	mux.Handle("/api/admin/gdrive-credentials", observe(cfg, adminGDriveCredentialsHandler(cfg)))
	mux.Handle("/admin/errors", observe(cfg, axiomErrorsHandler(tmpl, cfg, navConfigJS)))
	mux.Handle("/", observe(cfg, homeHandler(tmpl, cfg, navConfigJS)))

	addr := ":" + cfg.port
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
