package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

// ── Config ────────────────────────────────────────────────────────────────────

type config struct {
	supabaseURL  string
	supabaseAnon string
	axiomToken   string
	axiomDataset string
	axiomAPIURL  string
	port         string
}

func loadConfig() config {
	return config{
		supabaseURL:  mustEnv("SUPABASE_URL"),
		supabaseAnon: mustEnv("SUPABASE_ANON_KEY"),
		axiomToken:   os.Getenv("AXIOM_TOKEN"), // optional — skips Axiom if unset
		axiomDataset: mustEnv("AXIOM_DATASET"),
		axiomAPIURL:  envOr("AXIOM_API_URL", "https://api.axiom.co"),
		port:         envOr("PORT", "8080"),
	}
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
	url := fmt.Sprintf("%s/v1/datasets/%s/ingest", cfg.axiomAPIURL, cfg.axiomDataset)
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
	resp.Body.Close()
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
	Events  []axiomEvent
	FetchErr string
	NavFavsJSON   template.JS
	NavConfigJSON template.JS
}

func axiomErrorsHandler(tmpl *template.Template, cfg config, navConfigJS template.JS) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := axiomErrorsData{NavConfigJSON: navConfigJS}

		var favs []NavFav
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
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
			body, _ := json.Marshal(map[string]string{"apl": apl})
			url := fmt.Sprintf("%s/v1/datasets/%s/query", cfg.axiomAPIURL, cfg.axiomDataset)
			req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
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
					var qr axiomQueryResp
					if err := json.NewDecoder(resp.Body).Decode(&qr); err != nil {
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
			"templates/index.html",
			"templates/axiom_errors.html",
		),
	)

	fs := http.FileServer(http.Dir("."))
	mux := http.NewServeMux()
	mux.Handle("/shared/", observe(cfg, fs))
	mux.Handle("/navigation_config.json", observe(cfg, fs))
	mux.Handle("/api/nav/favs", observe(cfg, navFavsHandler(cfg)))
	mux.Handle("/admin/errors", observe(cfg, axiomErrorsHandler(tmpl, cfg, navConfigJS)))
	mux.Handle("/", observe(cfg, homeHandler(tmpl, cfg, navConfigJS)))

	addr := ":" + cfg.port
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
