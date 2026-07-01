package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

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

		var modules []struct {
			ID int `json:"id"`
		}
		modQ := fmt.Sprintf("select=id&module_number=eq.%d&limit=1", req.ModuleNumber)
		if err := supabaseGet(ctx, cfg, "modules", modQ, &modules); err != nil || len(modules) == 0 {
			http.Error(w, `{"error":"module not found"}`, http.StatusNotFound)
			return
		}

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
				var single Suggestion
				if json.Unmarshal([]byte(raw), &single) == nil {
					suggestions = []Suggestion{single}
				}
			}
		}

		if suggestions == nil {
			suggestions = []Suggestion{}
		}

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

func slideGenerateHandler(cfg config) http.HandlerFunc {
	type SlideRequest struct {
		Sentences        string `json:"sentences"`
		CourseName       string `json:"courseName"`
		ModuleName       string `json:"moduleName"`
		VideoName        string `json:"videoName"`
		BrandFrontmatter string `json:"brandFrontmatter"`
	}

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
