package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

type indexData struct {
	Course        *Course
	Tools         []Tool
	FetchErr      string
	NavFavsJSON   template.JS
	NavConfigJSON template.JS
}

func homeHandler(tmpl *template.Template, cfg config, navConfigJS template.JS) http.HandlerFunc {
	static := http.FileServer(http.Dir("."))
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" && r.URL.Path != "/index.html" {
			static.ServeHTTP(w, r)
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
		if err := supabaseGet(ctx, cfg, "nav_favorites", "select=url,label&order=updated_at.desc", &favs); err != nil {
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

func navFavsHandler(cfg config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		if r.Method == http.MethodDelete {
			u := r.URL.Query().Get("url")
			if u == "" {
				http.Error(w, "bad request: missing url", http.StatusBadRequest)
				return
			}
			_ = supabaseDelete(ctx, cfg, "nav_favorites", "url=eq."+url.QueryEscape(u))
			fmt.Fprint(w, `{"favorited":false}`)
			return
		}

		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req NavFav
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.URL == "" {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		_ = supabaseDelete(ctx, cfg, "nav_favorites", "url=eq."+url.QueryEscape(req.URL))
		_ = supabasePost(ctx, cfg, "nav_favorites", req)
		fmt.Fprint(w, `{"favorited":true}`)
	}
}
