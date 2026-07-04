package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var (
	reTitle      = regexp.MustCompile(`<title[^>]*>([^<]+)</title>`)
	reTag        = regexp.MustCompile(`<[^>]+>`)
	reWhitespace = regexp.MustCompile(`\s+`)
	reScript     = regexp.MustCompile(`(?is)<script[^>]*>.*?</script>`)
	reStyle      = regexp.MustCompile(`(?is)<style[^>]*>.*?</style>`)
	reComment    = regexp.MustCompile(`(?is)<!--.*?-->`)
)

type searchResult struct {
	ID          int     `json:"id"`
	URL         string  `json:"url"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	MenuLabel   string  `json:"menu_label"`
	Category    string  `json:"category"`
	Excerpt     string  `json:"excerpt"`
	Rank        float64 `json:"rank"`
}

type searchResponse struct {
	Query   string         `json:"query"`
	Results []searchResult `json:"results"`
	Count   int            `json:"count"`
}

type reindexResponse struct {
	Success     bool     `json:"success"`
	PagesAdded  int      `json:"pages_added"`
	Errors      []string `json:"errors,omitempty"`
	Elapsed     string   `json:"elapsed"`
}

// ── Public search endpoint ────────────────────────────────────────────────────
// GET /api/search?q=query
func searchHandler(cfg config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		q := strings.TrimSpace(r.URL.Query().Get("q"))
		if q == "" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(searchResponse{Query: "", Results: []searchResult{}, Count: 0})
			return
		}

		ctx := r.Context()
		body := map[string]any{
			"query_text":  q,
			"max_results": 30,
		}
		var buf bytes.Buffer
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		u := cfg.supabaseURL + "/rest/v1/rpc/search_pages"
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, &buf)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		req.Header.Set("apikey", cfg.supabaseAnon)
		req.Header.Set("Authorization", "Bearer "+cfg.supabaseAnon)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 400 {
			b, _ := io.ReadAll(resp.Body)
			http.Error(w, fmt.Sprintf("search_pages RPC failed: %s %s", resp.Status, b), http.StatusInternalServerError)
			return
		}

		var results []searchResult
		if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if results == nil {
			results = []searchResult{}
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-store")
		json.NewEncoder(w).Encode(searchResponse{
			Query:   q,
			Results: results,
			Count:   len(results),
		})
	}
}

// ── Admin reindex endpoint ────────────────────────────────────────────────────
// POST /api/admin/reindex-search
func adminReindexSearchHandler(cfg config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !isAdminRequest(r) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		start := time.Now()
		var errors []string
		pagesAdded := 0

		// Index HTML pages from filesystem
		err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return nil // skip inaccessible
			}
			if d.IsDir() {
				skip := d.Name() == ".git" || d.Name() == "node_modules" || d.Name() == "_site" ||
					d.Name() == "server" || d.Name() == "vendor"
				if skip {
					return filepath.SkipDir
				}
				return nil
			}
			if !strings.HasSuffix(path, ".html") {
				return nil
			}

			// Skip template shell files
			if strings.Contains(path, "/templates/") {
				return nil
			}

			if err := indexHTMLPage(cfg, r, path); err != nil {
				errors = append(errors, fmt.Sprintf("%s: %v", path, err))
			} else {
				pagesAdded++
			}
			return nil
		})

		if err != nil {
			errors = append(errors, fmt.Sprintf("walk error: %v", err))
		}

		// Also re-index from navigation_menus
		if err := indexFromNavMenus(cfg, r); err != nil {
			errors = append(errors, fmt.Sprintf("nav_menus: %v", err))
		}

		elapsed := time.Since(start).Round(time.Millisecond).String()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(reindexResponse{
			Success:    true,
			PagesAdded: pagesAdded,
			Errors:     errors,
			Elapsed:    elapsed,
		})

		log.Printf("[reindex] %d pages indexed in %s (%d errors)", pagesAdded, elapsed, len(errors))
	}
}

func indexHTMLPage(cfg config, r *http.Request, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	html := string(data)
	title, bodyText := extractPageContent(html)

	// Normalize URL path
	urlPath := "/" + path

	// Determine category from path
	category := inferCategory(path)

	entry := map[string]any{
		"url":         urlPath,
		"title":       title,
		"content":     bodyText,
		"description": "",
		"menu_label":  "",
		"category":    category,
	}

	// Upsert using on_conflict=url
	return supabaseUpsert(r.Context(), cfg, "page_search_index", entry, "url")
}

func extractPageContent(html string) (title string, bodyText string) {
	// Remove comments, scripts, styles first
	html = reComment.ReplaceAllString(html, "")
	html = reScript.ReplaceAllString(html, " ")
	html = reStyle.ReplaceAllString(html, " ")

	// Extract title
	if m := reTitle.FindStringSubmatch(html); len(m) > 1 {
		title = strings.TrimSpace(m[1])
	}

	// Strip HTML tags
	bodyText = reTag.ReplaceAllString(html, " ")
	bodyText = reWhitespace.ReplaceAllString(bodyText, " ")
	bodyText = strings.TrimSpace(bodyText)

	// Limit content size (first 10KB of text)
	if len(bodyText) > 10000 {
		bodyText = bodyText[:10000]
	}

	return title, bodyText
}

func inferCategory(path string) string {
	path = strings.TrimPrefix(path, "5_Symbols/")
	if strings.Contains(path, "preprod") {
		return "preprod"
	}
	if strings.Contains(path, "prod/") {
		return "production"
	}
	if strings.Contains(path, "postprod") {
		return "postprod"
	}
	if strings.Contains(path, "publish") {
		return "publish"
	}
	if strings.HasPrefix(path, "index.html") || strings.HasPrefix(path, "/index.html") {
		return "home"
	}
	if strings.HasPrefix(path, "1_Real_Unknown") {
		return "planning"
	}
	if strings.HasPrefix(path, "2_Environment") {
		return "environment"
	}
	if strings.HasPrefix(path, "3_Simulation") {
		return "simulation"
	}
	if strings.HasPrefix(path, "4_Formula") {
		return "formula"
	}
	if strings.HasPrefix(path, "6_Semblance") {
		return "semblance"
	}
	if strings.HasPrefix(path, "7_Testing_Known") {
		return "testing"
	}
	return "other"
}

func indexFromNavMenus(cfg config, r *http.Request) error {
	ctx := r.Context()
	u := cfg.supabaseURL + "/rest/v1/navigation_menus?select=url%2Clabel%2Cdescription%2Cmenu_type&url=not.is.null"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return err
	}
	req.Header.Set("apikey", cfg.supabaseAnon)
	req.Header.Set("Authorization", "Bearer "+cfg.supabaseAnon)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var items []struct {
		URL         string `json:"url"`
		Label       string `json:"label"`
		Description string `json:"description"`
		MenuType    string `json:"menu_type"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return err
	}

	for _, item := range items {
		entry := map[string]any{
			"url":         item.URL,
			"title":       item.Label,
			"description": item.Description,
			"menu_label":  item.Label,
			"category":    item.MenuType,
		}
		if err := supabaseUpsert(ctx, cfg, "page_search_index", entry, "url"); err != nil {
			log.Printf("[reindex] nav menu upsert error for %s: %v", item.URL, err)
		}
	}

	return nil
}
