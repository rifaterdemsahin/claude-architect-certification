package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type AnalogyProps struct {
	AnalogyTheme     string   `json:"analogyTheme"`
	MainTitle        string   `json:"mainTitle"`
	ProblemTitle     string   `json:"problemTitle"`
	ProblemConcept   string   `json:"problemConcept"`
	ProblemAnalogy   string   `json:"problemAnalogy"`
	ProblemCallouts  []string `json:"problemCallouts"`
	ProblemMood      string   `json:"problemMood"`
	SolutionTitle    string   `json:"solutionTitle"`
	SolutionConcept  string   `json:"solutionConcept"`
	SolutionAnalogy  string   `json:"solutionAnalogy"`
	SolutionCallouts []string `json:"solutionCallouts"`
	SolutionMood     string   `json:"solutionMood"`
	FooterTakeaway   string   `json:"footerTakeaway"`
}

func cleanAnalogyThemeSlug(theme string) string {
	t := strings.ToLower(strings.TrimSpace(theme))
	switch {
	case strings.Contains(t, "rac"):
		return "racing"
	case strings.Contains(t, "cook"):
		return "cooking"
	case strings.Contains(t, "const"):
		return "construction"
	case strings.Contains(t, "sail"):
		return "sailing"
	case strings.Contains(t, "aviat") || strings.Contains(t, "plane") || strings.Contains(t, "flight"):
		return "aviation"
	case strings.Contains(t, "space") || strings.Contains(t, "rocket"):
		return "space"
	case strings.Contains(t, "medic") || strings.Contains(t, "surg") || strings.Contains(t, "hosp"):
		return "medical"
	case strings.Contains(t, "traff") || strings.Contains(t, "highw") || strings.Contains(t, "road"):
		return "traffic"
	case strings.Contains(t, "fact") || strings.Contains(t, "assemb"):
		return "factory"
	case strings.Contains(t, "gard") || strings.Contains(t, "plant") || strings.Contains(t, "grow"):
		return "gardening"
	default:
		return "custom"
	}
}

func analogyDefaultProps(theme, sentence, sentenceType, moduleName, videoName string) AnalogyProps {
	cleanSentence := strings.TrimSpace(sentence)
	if cleanSentence == "" {
		cleanSentence = "Unindexed raw data loading vs. semantic vector search"
	}
	shortConcept := cleanSentence
	if len(shortConcept) > 80 {
		shortConcept = shortConcept[:77] + "..."
	}

	props := AnalogyProps{
		AnalogyTheme: theme,
		MainTitle:    fmt.Sprintf("%s vs. Optimized Architecture", shortConcept),
		ProblemTitle: "The Inefficient Approach",
		ProblemConcept:   cleanSentence,
		ProblemCallouts:  []string{"High Latency", "Slower Initialization", "High Memory Footprint", "Resource Bottleneck"},
		SolutionTitle:    "The Optimized Solution",
		SolutionConcept:  "Streamlined cloud-native architecture with dynamic semantic retrieval and caching",
		SolutionCallouts: []string{"Instant Response", "Sub-second Startup", "Low Memory Footprint", "Linear Scalability"},
		FooterTakeaway:   "Adopt modern semantic indexing and caching to maximize cloud application performance and efficiency.",
	}

	switch strings.ToLower(theme) {
	case "racing", "racing 🏁":
		props.AnalogyTheme = "Racing 🏁"
		props.ProblemAnalogy = "A race car driving with the handbrake engaged, towing a heavy trailer through thick smoke and warning flags"
		props.ProblemMood = "Frustrated driver, chaotic pit stop, overheating engine indicators"
		props.SolutionAnalogy = "A sleek aerodynamic Formula 1 car accelerating effortlessly down a clean track with rocket boosters and speed trails"
		props.SolutionMood = "Confident driver, high-tech streamlined cockpit, glowing green velocity meters"
	case "cooking", "cooking 🍳":
		props.AnalogyTheme = "Cooking 🍳"
		props.ProblemAnalogy = "A chef chopping every ingredient from scratch in a cluttered, smokey kitchen during peak rush hour"
		props.ProblemMood = "Stressed chef, overflowing sink, burning timers and chaotic order slips"
		props.SolutionAnalogy = "An automated Michelin-star kitchen with pre-prepped mise-en-place, robotic precision, and clean induction stations"
		props.SolutionMood = "Calm master chef, spotless stainless steel surfaces, instant dish assembly"
	case "construction", "construction 🏗️":
		props.AnalogyTheme = "Construction 🏗️"
		props.ProblemAnalogy = "Workers carrying individual bricks up wooden ladders by hand on a muddy, unorganized building site"
		props.ProblemMood = "Exhausted workers, collapsing scaffolding, hazard signs and delayed timelines"
		props.SolutionAnalogy = "Automated tower cranes and modular prefabricated steel beams slotting together rapidly on a pristine digital site"
		props.SolutionMood = "Efficient engineers with tablets, glowing laser blueprints, rapid vertical progress"
	case "sailing", "sailing ⛵":
		props.AnalogyTheme = "Sailing ⛵"
		props.ProblemAnalogy = "A heavy wooden galleon battling against stormy headwinds with torn sails and a jammed anchor dragging in the mud"
		props.ProblemMood = "Panicking crew, crashing waves, dark storm clouds and creaking timber"
		props.SolutionAnalogy = "A state-of-the-art hydrofoil yacht skimming effortlessly above crystal-blue waters powered by steady tailwinds"
		props.SolutionMood = "Relaxed captain at a digital helm, sunlit horizon, smooth high-speed gliding"
	case "aviation", "aviation ✈️":
		props.AnalogyTheme = "Aviation ✈️"
		props.ProblemAnalogy = "An overloaded propeller plane taxiing on a bumpy dirt runway in dense fog with flashing engine warnings"
		props.ProblemMood = "Anxious pilot, rattling dashboard, low visibility and delayed departure alarms"
		props.SolutionAnalogy = "A supersonic jet breaking the sound barrier in clear stratosphere with streamlined swept wings and glowing afterburners"
		props.SolutionMood = "Calm flight crew, head-up holographic display, clear skies and optimal cruising altitude"
	case "space", "space exploration 🚀":
		props.AnalogyTheme = "Space Exploration 🚀"
		props.ProblemAnalogy = "A clunky rocket struggling against gravity with heavy chemical fuel tanks and sputtering ignition nozzles"
		props.ProblemMood = "Tense mission control, red siren lights, vibration warnings and high fuel consumption"
		props.SolutionAnalogy = "An advanced ion-propulsion spacecraft gliding smoothly through orbital trajectories with glowing energy shields"
		props.SolutionMood = "Serene astronauts, blue holographic navigation consoles, effortless stellar velocity"
	case "medical", "medical / surgery 🏥":
		props.AnalogyTheme = "Medical / Surgery 🏥"
		props.ProblemAnalogy = "A surgeon searching through messy paper patient charts in a dimly lit, crowded emergency room"
		props.ProblemMood = "Overwhelmed staff, beeping monitors, chaotic piles of folders and critical delays"
		props.SolutionAnalogy = "A robotic surgical suite with real-time AI biometrics, holographic patient overlays, and instantaneous diagnosis"
		props.SolutionMood = "Focused specialist, sterile glowing environment, precise and rapid intervention"
	case "traffic", "traffic / highway 🛣️":
		props.AnalogyTheme = "Traffic / Highway 🛣️"
		props.ProblemAnalogy = "A gridlocked single-lane city street blocked by delivery trucks, potholes, and red traffic lights at every intersection"
		props.ProblemMood = "Honking cars, frustrated drivers, exhaust fumes and standstill delays"
		props.SolutionAnalogy = "An automated multi-level express smart highway with green synchronized lights and high-speed autonomous transit lanes"
		props.SolutionMood = "Seamless traffic flow, glowing LED guidance lines, zero friction travel"
	case "factory", "factory / assembly line 🏭":
		props.AnalogyTheme = "Factory / Assembly Line 🏭"
		props.ProblemAnalogy = "Workers manually passing heavy parts along a jammed conveyor belt with frequent mechanical breakdowns and sparks"
		props.ProblemMood = "Stressed supervisors, piling inventory, alarm bells and bottlenecked chutes"
		props.SolutionAnalogy = "A high-speed robotic assembly line with synchronized robotic arms, automated quality vision sensors, and zero waste"
		props.SolutionMood = "Clean tech floor, real-time analytics dashboards, rhythmic flawless throughput"
	case "gardening", "gardening 🌱":
		props.AnalogyTheme = "Gardening 🌱"
		props.ProblemAnalogy = "Watering a dry, rocky field bucket by bucket from a distant well while weeds choke the crops"
		props.ProblemMood = "Wilted plants, cracked earth, exhausting manual labor and poor yield"
		props.SolutionAnalogy = "An automated hydroponic greenhouse with nutrient drip irrigation, LED grow lights, and lush flourishing vegetation"
		props.SolutionMood = "Vibrant green foliage, automated sensors, effortless growth and bountiful harvests"
	default:
		props.AnalogyTheme = "Custom Analogy"
		props.ProblemAnalogy = "An outdated, friction-heavy mechanism struggling under heavy load with manual bottlenecks"
		props.ProblemMood = "Chaotic, slow, resource-intensive and error-prone"
		props.SolutionAnalogy = "A streamlined, automated, cloud-native architecture operating with zero friction and high velocity"
		props.SolutionMood = "Clean, responsive, high-speed and effortless"
	}
	return props
}

func formatAnalogyPrompt(props AnalogyProps) string {
	return fmt.Sprintf(`Create a high-quality 1080p infographic that uses a clear, side-by-side split screen comparison to explain a technical concept through a stark analogy.

**The Analogy Theme:** %s

**Left Side: The Inefficient/Problem Approach**
- **Concept:** %s
- **Visual Analogy:** %s
- **Key Callouts/Labels:** Include 3-4 text bubbles pointing out the technical consequences: %s.
- **Tone/Mood:** %s

**Right Side: The Optimized/Solution Approach**
- **Concept:** %s
- **Visual Analogy:** %s
- **Key Callouts/Labels:** Include 3-4 text bubbles pointing out the technical benefits: %s.
- **Tone/Mood:** %s

**Overall Design & Layout:**
- **Structure:** Split vertically or clearly divided down the center.
- **Top Header:** A prominent title banner reading: "%s"
- **Section Headers:** Clear subtitles for both sides: "%s" vs. "%s"
- **Bottom Footer:** A summary banner with a concluding takeaway: "%s"
- **Style:** Bright, clean, and vibrant vector illustration/infographic style. Text labels must be highly legible, integrated cleanly into the graphic, and directly mapped to the visual elements.`,
		props.AnalogyTheme,
		props.ProblemConcept,
		props.ProblemAnalogy,
		strings.Join(props.ProblemCallouts, ", "),
		props.ProblemMood,
		props.SolutionConcept,
		props.SolutionAnalogy,
		strings.Join(props.SolutionCallouts, ", "),
		props.SolutionMood,
		props.MainTitle,
		props.ProblemTitle,
		props.SolutionTitle,
		props.FooterTakeaway,
	)
}

func analogyPromptHandler(cfg config) http.HandlerFunc {
	type Req struct {
		Sentence     string `json:"sentence"`
		SentenceType string `json:"sentenceType"`
		AnalogyTheme string `json:"analogyTheme"`
		ModuleName   string `json:"moduleName"`
		VideoName    string `json:"videoName"`
		CustomPrompt string `json:"customPrompt"`
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

		props := analogyDefaultProps(req.AnalogyTheme, req.Sentence, req.SentenceType, req.ModuleName, req.VideoName)

		apiKey := cfg.getSecret("OPENROUTER_API_KEY")
		if apiKey != "" && strings.TrimSpace(req.Sentence) != "" {
			defaultPropsJSON, _ := json.MarshalIndent(props, "", "  ")
			aiPrompt := fmt.Sprintf(`You are an expert technical educational designer and infographic creator.
I will provide a sentence from a cloud architecture course script and an analogy theme.
Output structured JSON properties for a side-by-side split screen comparison infographic.

Theme: %s
Sentence: %q

Instructions:
1. Return ONLY valid JSON matching this schema:
%s
2. Intelligently adapt titles, concepts, analogies, callouts (3-4 strings per array), and mood to match the exact meaning of the Sentence.
3. No markdown wrapper or prose, just the JSON object.`, req.AnalogyTheme, req.Sentence, string(defaultPropsJSON))

			payload := map[string]any{
				"model":       "anthropic/claude-sonnet-4.6",
				"messages":    []map[string]string{{"role": "user", "content": aiPrompt}},
				"temperature": 0.3,
			}
			body, _ := json.Marshal(payload)
			ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
			defer cancel()
			hreq, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://openrouter.ai/api/v1/chat/completions", bytes.NewReader(body))
			if err == nil {
				hreq.Header.Set("Content-Type", "application/json")
				hreq.Header.Set("Authorization", "Bearer "+apiKey)
				hresp, err := http.DefaultClient.Do(hreq)
				if err == nil && hresp.StatusCode == http.StatusOK {
					defer hresp.Body.Close()
					var orResp struct {
						Choices []struct {
							Message struct {
								Content string `json:"content"`
							} `json:"message"`
						} `json:"choices"`
					}
					if json.NewDecoder(hresp.Body).Decode(&orResp) == nil && len(orResp.Choices) > 0 {
						cleanContent := strings.TrimSpace(orResp.Choices[0].Message.Content)
						cleanContent = strings.TrimPrefix(cleanContent, "```json")
						cleanContent = strings.TrimPrefix(cleanContent, "```")
						cleanContent = strings.TrimSuffix(cleanContent, "```")
						var aiProps AnalogyProps
						if json.Unmarshal([]byte(cleanContent), &aiProps) == nil && aiProps.MainTitle != "" {
							props = aiProps
						}
					}
				}
			}
		}

		promptStr := strings.TrimSpace(req.CustomPrompt)
		if promptStr == "" {
			promptStr = formatAnalogyPrompt(props)
		}

		json.NewEncoder(w).Encode(map[string]any{
			"ok":           true,
			"prompt":       promptStr,
			"analogyProps": props,
		})
	}
}

func analogyGenerateHandler(cfg config) http.HandlerFunc {
	type GenReq struct {
		SentenceID   int          `json:"sentence_id"`
		ModuleNumber int          `json:"module_number"`
		VideoNumber  int          `json:"video_number"`
		ScriptID     int          `json:"script_id"`
		AnalogyTheme string       `json:"analogy_theme"`
		Prompt       string       `json:"prompt"`
		AnalogyProps AnalogyProps `json:"analogy_props"`
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
		if !isAdminRequest(r) {
			http.Error(w, `{"error":"Unauthorized — sign in as admin to generate analogy infographics."}`, http.StatusUnauthorized)
			return
		}
		var req GenReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
			return
		}

		geminiKey := cfg.getSecret("GEMINI_API_KEY")
		if geminiKey == "" {
			http.Error(w, `{"error":"GEMINI_API_KEY missing — configure in admin Environment or .env"}`, http.StatusServiceUnavailable)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 120*time.Second)
		defer cancel()

		const imageModel = "gemini-2.5-flash-image"
		imageURL := "https://generativelanguage.googleapis.com/v1beta/models/" + imageModel + ":generateContent?key=" + geminiKey
		imageBody, _ := json.Marshal(map[string]any{
			"contents": []any{map[string]any{"parts": []any{map[string]any{"text": req.Prompt}}}},
			"generationConfig": map[string]any{
				"responseModalities": []string{"IMAGE"},
			},
		})
		imageReq, _ := http.NewRequestWithContext(ctx, "POST", imageURL, bytes.NewReader(imageBody))
		imageReq.Header.Set("Content-Type", "application/json")
		imageResp, err := http.DefaultClient.Do(imageReq)
		if err != nil {
			http.Error(w, `{"error":"Gemini image generation request failed"}`, http.StatusBadGateway)
			return
		}
		defer imageResp.Body.Close()
		if imageResp.StatusCode >= 400 {
			body, _ := io.ReadAll(imageResp.Body)
			http.Error(w, fmt.Sprintf(`{"error":"Gemini image model HTTP %d: %s"}`, imageResp.StatusCode, strings.TrimSpace(string(body))), http.StatusBadGateway)
			return
		}

		var imgParsed geminiContentResp
		if err := json.NewDecoder(imageResp.Body).Decode(&imgParsed); err != nil {
			http.Error(w, `{"error":"failed to decode image generation response"}`, http.StatusBadGateway)
			return
		}

		var imageData []byte
		contentType := "image/png"
		if len(imgParsed.Candidates) > 0 {
			for _, p := range imgParsed.Candidates[0].Content.Parts {
				if p.InlineData.Data != "" {
					if p.InlineData.MimeType != "" {
						contentType = p.InlineData.MimeType
					}
					decoded, derr := base64.StdEncoding.DecodeString(p.InlineData.Data)
					if derr == nil && len(decoded) > 0 {
						imageData = decoded
						break
					}
				}
			}
		}
		if len(imageData) == 0 {
			http.Error(w, `{"error":"Gemini returned no valid image data"}`, http.StatusBadGateway)
			return
		}

		cleanTheme := cleanAnalogyThemeSlug(req.AnalogyTheme)
		blobName := fmt.Sprintf("analogy_m%d_v%d_s%d_%s_%d.png", req.ModuleNumber, req.VideoNumber, req.SentenceID, cleanTheme, time.Now().Unix())
		container := "research-infographics"

		var proxyURL string
		if cfg.azureAccountName != "" {
			if uerr := uploadBlobToAzure(ctx, cfg, container, blobName, contentType, imageData); uerr != nil {
				log.Printf("azure analogy upload error: %v", uerr)
			} else {
				proxyURL = researchFileProxyURL(container, blobName)
			}
		}

		dbRow := map[string]any{
			"sentence_id":       req.SentenceID,
			"module_number":     req.ModuleNumber,
			"video_number":      req.VideoNumber,
			"script_id":         req.ScriptID,
			"analogy_theme":     cleanTheme,
			"generation_status": "completed",
			"prompt_used":       req.Prompt,
			"analogy_props":     req.AnalogyProps,
			"azure_blob_name":   blobName,
			"image_url":         proxyURL,
		}
		if err := supabaseUpsert(ctx, cfg, "sentence_analogies", dbRow, "sentence_id,analogy_theme"); err != nil {
			log.Printf("supabase save sentence_analogies error: %v", err)
		}

		json.NewEncoder(w).Encode(map[string]any{
			"ok":              true,
			"image_url":       proxyURL,
			"azure_blob_name": blobName,
			"analogy_theme":   cleanTheme,
		})
	}
}

func analogySaveHandler(cfg config) http.HandlerFunc {
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
		if !isAdminRequest(r) {
			http.Error(w, `{"error":"Unauthorized — sign in as admin to save analogy rows."}`, http.StatusUnauthorized)
			return
		}
		var row map[string]any
		if err := json.NewDecoder(r.Body).Decode(&row); err != nil {
			http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
			return
		}
		if err := supabaseUpsert(r.Context(), cfg, "sentence_analogies", row, "sentence_id,analogy_theme"); err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"database upsert failed: %s"}`, err.Error()), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]any{"ok": true})
	}
}
