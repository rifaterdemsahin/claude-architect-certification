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
	"strconv"
	"strings"
	"time"
)

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
		if len(events) > limit {
			events = events[:limit]
		}
		json.NewEncoder(w).Encode(map[string]any{"events": events, "count": len(events)})
	}
}

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
