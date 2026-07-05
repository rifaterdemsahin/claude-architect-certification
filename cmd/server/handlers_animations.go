package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

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

func normalizeModelForOpenRouter(displayName string) string {
	t := strings.ToLower(displayName)
	if t == "" || strings.Contains(t, "auto-detect") || strings.Contains(t, "executing ai") {
		return "anthropic/claude-sonnet-4.6"
	}
	if strings.Contains(t, "deepseek v4 pro") || strings.Contains(t, "deepseek-v4-pro") {
		return "deepseek/deepseek-v4-pro"
	}
	if strings.Contains(t, "deepseek v4 flash") || strings.Contains(t, "deepseek-v4-flash") {
		return "deepseek/deepseek-v4-flash"
	}
	if strings.Contains(t, "deepseek r1") || strings.Contains(t, "deepseek v3") {
		return "deepseek/deepseek-r1"
	}
	if strings.Contains(t, "gemini 3.1") || strings.Contains(t, "gemini 2.5 pro") || strings.Contains(t, "gemini-2.5") {
		return "google/gemini-2.5-pro-preview-05-06"
	}
	if strings.Contains(t, "gemini 2.0 flash") || strings.Contains(t, "gemini-2.0-flash") {
		return "google/gemini-2.0-flash-001"
	}
	if strings.Contains(t, "claude sonnet 4.6") || strings.Contains(t, "claude-sonnet-4.6") {
		return "anthropic/claude-sonnet-4.6"
	}
	if strings.Contains(t, "claude 3.7") {
		return "anthropic/claude-3.7-sonnet"
	}
	if strings.Contains(t, "gpt-4o") || strings.Contains(t, "openai") {
		return "openai/gpt-4o"
	}
	if strings.Contains(t, "deepseek") {
		return "deepseek/deepseek-chat"
	}
	if strings.Contains(t, "gemini") {
		return "google/gemini-2.5-pro-preview-05-06"
	}
	if strings.Contains(t, "claude") {
		return "anthropic/claude-sonnet-4.6"
	}
	return displayName
}

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

func animationDefaultProps(animType, sentence, sentenceType, moduleTitle, videoTitle string) map[string]any {
	base := map[string]any{
		"title":              clampLen(sentence, 64),
		"subtitle":           clampLen(moduleTitle+" · "+videoTitle, 80),
		"sentence":           sentence,
		"sentenceType":       sentenceType,
		"module":             moduleTitle,
		"video":              videoTitle,
		"fps":                30,
		"durationInFrames":   150,
		"brandColor":         "#8b5cf6",
		"secondaryColor":     "#3b82f6",
		"bgColor":            "#030712",
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
			"animationType":    meta.Slug,
			"label":            meta.Label,
			"emoji":            meta.Emoji,
			"prompt":           prompt,
			"inputProps":       props,
			"runpodPayload":    runpodPayload,
			"serveUrl":         serveURL,
			"runpodConfigured": cfg.getSecret("RUNPOD_API_KEY") != "" && os.Getenv("RUNPOD_ENDPOINT_ID") != "" && serveURL != "",
		})
	}
}

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
1. Return ONLY a valid JSON object. No markdown formatting, no prose, no code blocks (do not wrap in `+"`"+`json`+"`"+`).
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
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadGateway)
			json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("OpenRouter HTTP %d", hresp.StatusCode), "detail": string(b)})
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
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadGateway)
			json.NewEncoder(w).Encode(map[string]string{"error": "OpenRouter did not return valid JSON", "detail": content})
			return
		}

		json.NewEncoder(w).Encode(parsed)
	}
}

func animationRunpodRunHandler(cfg config) http.HandlerFunc {
	type Req struct {
		SentenceID    int            `json:"sentenceId"`
		ModuleNumber  int            `json:"moduleNumber"`
		VideoNumber   int            `json:"videoNumber"`
		ScriptID      int            `json:"scriptId"`
		AnimationType string         `json:"animationType"`
		Prompt        string         `json:"prompt"`
		InputProps    map[string]any `json:"inputProps"`
		CustomPrompt  string         `json:"customPrompt"`
		Duration      int            `json:"durationInFrames"`
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

		row := map[string]any{
			"sentence_id":       req.SentenceID,
			"module_number":     req.ModuleNumber,
			"video_number":      req.VideoNumber,
			"script_id":         req.ScriptID,
			"animation_type":    req.AnimationType,
			"generation_status": "generating",
			"prompt_used":       strings.TrimSpace(req.CustomPrompt),
			"remotion_props":    req.InputProps,
			"runpod_job_id":     rpResp.ID,
			"runpod_status":     rpResp.Status,
			"codec":             "h264",
			"duration_in_frames": dur,
			"fps":               30,
			"width":             1920,
			"height":            1080,
			"error_message":     "",
		}
		if strings.TrimSpace(req.CustomPrompt) == "" {
			row["prompt_used"] = req.Prompt
		}
		if err := supabaseUpsert(ctx, cfg, "sentence_animations", row, "sentence_id,animation_type"); err != nil {
			log.Printf("sentence_animations upsert: %v", err)
		}

		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(map[string]any{
			"jobId":    rpResp.ID,
			"status":   rpResp.Status,
			"endpoint": endpointID,
		})
	}
}

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
					"runpod_status":     status,
					"azure_blob_name":   blobName,
					"animation_url":     proxyURL,
					"error_message":     "",
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

func queryInt(r *http.Request, key string) int {
	v, _ := strconv.Atoi(r.URL.Query().Get(key))
	return v
}

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

		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(map[string]any{
			"animationType": meta.Slug,
			"label":         meta.Label,
			"emoji":         meta.Emoji,
			"prompt":        prompt,
			"inputProps":    props,
			"jobId":         rpRun.ID,
			"endpoint":      endpointID,
			"status":        rpRun.Status,
		})
	}
}

func animationRemotionInstructionsHandler(cfg config) http.HandlerFunc {
	type Req struct {
		PromptText string `json:"promptText"`
		Model      string `json:"model"`
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

		apiKey := cfg.getSecret("OPENROUTER_API_KEY")
		if apiKey == "" {
			http.Error(w, `{"error":"OPENROUTER_API_KEY missing"}`, http.StatusServiceUnavailable)
			return
		}

		var req Req
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
			return
		}
		if req.PromptText == "" {
			http.Error(w, `{"error":"promptText is required"}`, http.StatusBadRequest)
			return
		}

		model := req.Model
		if model == "" {
			model = "anthropic/claude-sonnet-4.6"
		}
		model = normalizeModelForOpenRouter(model)
		systemPrompt := `You are a senior technical architect, animation pipeline coordinator, and spec-driven development lead. Your job is to analyze the provided project generation prompt and produce a comprehensive end-to-end execution plan — with each task defined as a specification — that maps every deliverable to the available skills, Kilo agents, and commands from the animation-template-llm-training project.

## Available Kilo Commands (from AGENTS.md)
| Command | Purpose |
|---------|---------|
| /pipeline "<topic>" | Full pipeline: generate assets + render all videos |
| /render [scene] | Render-only: render Remotion MP4s from existing assets |
| /generate-assets "<topic>" | Generate SVG + narration + MP3 via Gemini |
| /serve | Start Flask dev server (port 5177) |

## Available Kilo Agents
| Agent | Specializes in |
|-------|---------------|
| remotion-dev | Remotion composition development, debugging, scene creation, React/TypeScript |
| infographic-builder | SVG infographic generation, narration scripts, Gemini API calls, asset pipeline |
| architect | General project setup, package.json, config files, git, CI/CD, README |

## Available Asset Generation Pipeline
1. Flask server (server.py on port 5177) — POST /api/generate/infographic → generated-assets/infographic.svg
2. Flask server — POST /api/generate/audio → generated-assets/narration.txt  
3. macOS say + ffmpeg → generated-assets/narration.mp3
4. Delivery Pilot server — POST /api/infographics/generate → generated PNG infographic (Gemini via Azure Key Vault)
5. Excalidraw export → SVG files for hand-drawn architecture sketches
6. npm run render:all → remotion/exports/*.mp4 (uses svgData, audioSrc, script as input props)

## Template Decision Patterns (from llm_thinking_log.md)
| Pattern ID | Description |
|------------|-------------|
| project-naming | Derive kebab-case project name from root cause analysis of the exam question |
| duck-typed-contracts | Subagent naive/resilient share identical interface, coordinator swaps them |
| heuristic-models | Use simulated/controlled workloads for reproducible CLI benchmarks |
| token-calc | Token consumption = words × 1.35, display cost savings in tables |
| 4-scene-timeline | Scenes: Title(0-4s) → NaiveFails(4-12s) → ResilientWins(12-22s) → Metrics(22-28s) |
| color-story | red(danger/naive) → violet(tech/solution) → green(success/metrics) → cyan(highlight) |
| img-fallback | Gemini API → Excalidraw → placeholder SVG |
| mp4-fallback | Remotion → Canvas captureStream() → ffmpeg mock → static PNG sequence |
| web-speech-api | window.speechSynthesis for browser narration, sync with GSAP timeline |
| vanilla-widgets | Zero-dependency interactive UI: calculator, sandbox terminal simulator |
| qa-audit | Audit for overclaims, separate simulation from production, markdown issue table |
| file-split | Single-responsibility files: domain.js, infrastructure.js, subagent-*.js, coordinator.js |
| gh-pages-verify | Poll GitHub Actions workflow until live URL returns HTTP 200 |
| emoji-rich | Every heading, metric, badge, file-tree entry carries ≥1 relevant emoji |
| seven-words | ≤7 words per on-screen label, symbols & graphics explain, not words |
| modal-lightbox | Click any image → fullscreen viewport-scale modal with object-fit: contain |
| modular-config | Per-agent configs over one monolithic config file |

## 🎯 Your Task — Spec-Driven Execution Plan
Analyze the project generation prompt and produce a complete, spec-driven plan. Every task MUST be defined as a specification with clear acceptance criteria. For each section:

### 1. 🤖 Agent Plan — Task-to-Agent Specification Map
Map every task to the appropriate Kilo agent with spec-level detail:
- remotion-dev tasks: scene composition specs, animation timelines, debugging checkpoints
- infographic-builder tasks: SVG/narration generation specs, API call contracts
- architect tasks: project scaffolding, config files, CI/CD, README, SEO assets
- Kilo commands: which multi-step workflows each command automates, with exact arguments

### 2. 🧠 Skill Map — Requirement-to-Pattern Traceability
For each project requirement, link it to at least one template decision pattern from the table above. Include the reasoning (why this pattern applies).

### 3. 📋 Execution Plan — Phased Spec Breakdown
List every file to create in dependency order, with FULL spec-level detail:
- File path, agent assignment, dependencies, complexity (small/medium/large)
- **spec**: a one-sentence acceptance criterion (what must this file deliver?)
- **priority**: critical/high/medium/low
- **tags**: [backend, frontend, config, asset, render, deploy, docs, seo]
- **estimatedLines**: approximate line count
- **templatePattern**: which decision pattern from the table above applies
- Phase grouping: Scaffold → Core Code → Config & CI → Assets & Narration → Remotion Scenes → Render → Deploy & Verify

### 4. 🖼️ Asset Pipeline — API-to-Output Trace
Every generated asset must specify: the exact API endpoint, the topic/prompt, the output file path, the responsible agent, and the fallback if the API fails.

### 5. ⚠️ Fallback Strategy — Per-Component Resilience
Define a specific fallback for every component that could fail. Do NOT use a single sentence — list each failure scenario individually:
- Gemini API unavailable? (Excalidraw fallback)
- ffmpeg missing? (Web Speech API browser-side)
- Remotion headless fails? (Canvas API captureStream or static PNG sequence)
- Flask server down? (direct file creation)
- npm install fails? (manual dependency instructions)
- GitHub Pages deploy timeout? (manual deploy steps)

### 6. 🎬 Compositions — Detailed Scene Specs
Design every Remotion scene as a specification. For each composition include:
- Learning objective (single concept the viewer must grasp)
- Visual elements (icons, colors, animations)
- Assets used (from assetPipeline above)
- Key animation timeline (frame ranges with actions)
- Narration trigger (exact frame and ≤7 word label)
- Agent responsible and the exact Kilo command to render it
- Transition to next scene
- acceptanceCriteria: what must be true for this scene to be "done"

Return your answer as a JSON object with this exact schema:
{
  "agentPlan": {
    "remotionDev": ["spec: create Scene1Title.tsx — title card with spring entrance animation, 120f duration, cyan-on-dark theme, acceptance: renders without React errors and produces exports/scene1-title.mp4", "spec: wire Audio sync in Scene2Naive — useAudioData() hook, play() on frame 60, acceptance: audio plays in sync with red flash overlay"],
    "infographicBuilder": ["spec: generate SVG infographic via POST /api/generate/infographic with topic='subagent retry architecture', acceptance: SVG renders cleanly in browser with correct viewBox", "spec: write narration script for all 4 scenes, ≤7 words per label, acceptance: script parses as valid JS array of {text, selector, durationHint} objects"],
    "kiloCommands": ["/pipeline \"subagent-resilience-demo\" — runs generate-assets → renders all 4 scenes → stitches FullVideo", "/generate-assets \"subagent-resilience-demo\" — generates SVG + narration + MP3 for all scenes", "/render scene1 — renders only Scene1Title to exports/scene1-title.mp4"]
  },
  "skillMap": [
    {
      "skill": "duck-typed-contracts",
      "appliedTo": "src/subagent-naive.js and src/subagent-resilient.js share identical async process(item) → Result interface",
      "decisionPattern": "duck-typed-contracts: Subagent naive/resilient share identical interface, coordinator swaps them",
      "reasoning": "Enables side-by-side benchmarking without changing coordinator code"
    },
    {
      "skill": "4-scene-timeline",
      "appliedTo": "GSAP timeline in index.html + Remotion FullVideo.tsx",
      "decisionPattern": "4-scene-timeline: Title(0-4s) → NaiveFails(4-12s) → ResilientWins(12-22s) → Metrics(22-28s)",
      "reasoning": "Chronological storytelling makes the comparison intuitive for exam candidates"
    }
  ],
  "executionPlan": {
    "phases": [
      {
        "phase": "1. Project Scaffold",
        "steps": [
          {
            "order": 1,
            "file": "package.json",
            "agent": "architect",
            "dependsOn": [],
            "complexity": "small",
            "priority": "critical",
            "tags": ["config"],
            "estimatedLines": 25,
            "templatePattern": "modular-config",
            "action": "Create ESM package.json (\"type\":\"module\") with scripts: start, naive, resilient, benchmark, render:all",
            "spec": "package.json must define \"type\":\"module\", all render scripts, and pass JSON.parse() validation"
          }
        ]
      }
    ],
    "fileCreationOrder": [
      {"order": 1, "file": "package.json", "agent": "architect", "spec": "Valid ESM package.json with all required scripts"},
      {"order": 2, "file": "src/domain.js", "agent": "architect", "spec": "Domain model exports entity classes and simulated test corpus"}
    ]
  },
  "compositions": [
    {
      "name": "Scene1Title",
      "description": "Title card introducing the architectural problem and the solution being demonstrated",
      "learningObjective": "Understand what problem this demo solves and why resilient architecture matters",
      "durationInFrames": 120,
      "visualElements": ["shield-lock-icon", "cyan-gradient-title", "subtitle-fade-in"],
      "assetsUsed": ["generated-assets/infographic.svg"],
      "keyAnimation": "f0-30: title springs in from scale(0.5), f30-60: subtitle fades up, f60-90: problem statement appears, f90-120: hold for reading",
      "transitionToNext": "fade-to-black (15 frames), then Scene2 starts from black",
      "narrationTrigger": "frame 10 — label: '🔒 Resilient Subagent Architecture'",
      "agent": "remotion-dev",
      "kiloCommand": "/render scene1",
      "acceptanceCriteria": "Rendered MP4 exists at exports/scene1-title.mp4, is non-empty, title text is legible at 1080p, no React console errors"
    }
  ],
  "assetPipeline": [
    {
      "step": 1,
      "endpoint": "POST /api/generate/infographic",
      "topic": "subagent retry architecture diagram showing coordinator → subagent pool → local recovery loop",
      "output": "generated-assets/infographic.svg",
      "agent": "infographic-builder",
      "fallback": "Draw equivalent diagram in Excalidraw, export as SVG to same path"
    },
    {
      "step": 2,
      "endpoint": "macOS say + ffmpeg",
      "topic": "convert narration.txt to MP3 using: say -v Samantha -f narration.txt -o narration.aiff && ffmpeg -i narration.aiff narration.mp3",
      "output": "generated-assets/narration.mp3",
      "agent": "infographic-builder",
      "fallback": "Use Web Speech API directly in browser (index.html already has this as primary narration)"
    }
  ],
  "fallbackStrategy": {
    "geminiApi": "If Gemini API unavailable: draw diagrams in Excalidraw, export SVG to generated-assets/. Write narration scripts manually as JS arrays.",
    "ffmpeg": "If ffmpeg missing: skip MP3 generation entirely. index.html already uses Web Speech API for browser narration — Remotion scenes can use silence or the Web Speech API output captured separately.",
    "remotion": "If Remotion headless render fails: use Canvas API captureStream() to record browser animations as WebM, or generate a static PNG sequence with frame numbers.",
    "flaskServer": "If Flask server unavailable: create generated-assets/ directory manually and write SVG/narration files directly using file I/O.",
    "npmInstall": "If npm install fails: provide a MANUAL_SETUP.md with exact dependency version list and manual download instructions.",
    "ghPages": "If GitHub Pages deploy times out: provide a MANUAL_DEPLOY.md with step-by-step gh-pages branch creation and GitHub UI settings instructions."
  },
  "totalDurationEstimate": "~600 frames at 30fps (20 seconds total for full video)",
  "showNotTell": "Use ≤7 word labels with icons instead of sentences (e.g., '🔒 Sandbox isolates context' not 'The subagent runs in an isolated sandbox environment that prevents...'). Every visual element IS the explanation — no voiceover describing what the viewer already sees.",
  "colorStory": "cyan(trust/tech #00f0ff) → violet(solution #a855f7) → green(success #10b981) → red(danger/naive #ef4444)",
  "trainingPace": "3 quick intro scenes, 2 deep-dive with pause points, 1 review card, total ~20 minutes hands-on time"
}

Rules — follow these strictly:
- Map EVERY task to one of: remotion-dev, infographic-builder, or architect.
- The fileCreationOrder must be sorted by dependency (package.json before src/ before remotion/ before exports/).
- Every composition MUST have a kiloCommand that can render it.
- Every asset in assetPipeline MUST specify its endpoint, topic, output path, agent, AND fallback.
- Be specific about template decision patterns — use exact names from the table above.
- Every step in executionPlan.phases MUST include: order, file, agent, dependsOn, complexity, priority, tags, estimatedLines, templatePattern, action, AND spec.
- The fallback strategy MUST list each failure scenario individually, not as one combined sentence.
- Do NOT skip any section. Every field in the schema is mandatory.`
		payload := map[string]any{
			"model": model,
			"messages": []map[string]string{
				{"role": "system", "content": systemPrompt},
				{"role": "user", "content": fmt.Sprintf("Here is the project generation prompt. Extract Remotion animation instructions from it:\n\n%s", req.PromptText)},
			},
			"temperature": 0.3,
			"max_tokens":  8192,
		}
		body, _ := json.Marshal(payload)
		ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
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
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadGateway)
			json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("OpenRouter HTTP %d", hresp.StatusCode), "detail": string(b)})
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
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadGateway)
			detail := map[string]any{"error": "invalid OpenRouter response"}
			if err != nil {
				detail["parseError"] = err.Error()
			} else {
				detail["parseError"] = "no choices returned"
			}
			if len(b) > 0 {
				end := len(b)
				if end > 500 { end = 500 }
				detail["rawBody"] = string(b[:end])
			}
			json.NewEncoder(w).Encode(detail)
			return
		}

		content := strings.TrimSpace(orResp.Choices[0].Message.Content)
		content = strings.TrimPrefix(content, "```json")
		content = strings.TrimPrefix(content, "```")
		content = strings.TrimSuffix(content, "```")
		content = strings.TrimSpace(content)

		var parsed map[string]any
		if err := json.Unmarshal([]byte(content), &parsed); err != nil {
			json.NewEncoder(w).Encode(map[string]any{
				"raw":   content,
				"type":  "text",
				"model": model,
			})
			return
		}
		parsed["type"] = "structured"
		parsed["model"] = model
		json.NewEncoder(w).Encode(parsed)
	}
}
