package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func infographicGenerateHandler(cfg config) http.HandlerFunc {
	type InfoGenRequest struct {
		Topic string `json:"topic"`
		Style string `json:"style"`
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

		blobName := fmt.Sprintf("infographic_m%d_v%d_%d.json", req.ModuleNumber, req.VideoNumber, time.Now().Unix())
		container := "research-notes"

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
					blobName = ""
				}
			} else {
				log.Printf("sas error for infographic: %v", err)
				blobName = ""
			}
		}

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
