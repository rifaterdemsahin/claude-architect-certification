package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

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
