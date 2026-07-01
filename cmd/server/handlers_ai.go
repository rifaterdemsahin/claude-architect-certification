package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

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

func sanityCheckHandler(cfg config) http.HandlerFunc {
	type SanityRequest struct {
		ItemName string `json:"item_name"`
		ItemDesc string `json:"item_desc"`
		UserNote string `json:"user_note"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
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
