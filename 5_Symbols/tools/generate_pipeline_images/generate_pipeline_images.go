package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type Stage struct {
	ID    int
	Title string
	Desc  string
	Phase string
}

func main() {
	geminiKey := os.Getenv("GEMINI_API_KEY")
	if geminiKey == "" {
		log.Fatal("GEMINI_API_KEY not found in environment. Run sync_secrets first.")
	}

	stages := []Stage{
		{1, "Problem", "Identify and articulate the core problem the course will solve.", "Pre-Production"},
		{2, "Solution", "Define the course as the solution. Articulate what learners will be able to do.", "Pre-Production"},
		{3, "Research", "Deep-dive into the subject matter while self-learning the content.", "Pre-Production"},
		{4, "Outline", "Structure the course into modules and sections.", "Pre-Production"},
		{5, "Script", "Write the full narration script for each video.", "Pre-Production"},
		{6, "IVQ", "Design interactive quiz questions embedded at key checkpoints.", "Pre-Production"},
		{7, "Shot List", "Translate the script into a detailed shot-by-shot plan.", "Production"},
		{8, "Footage Mapping", "Map the raw recorded footage to each shot in the shot list.", "Production"},
		{9, "Edit List", "Create the editing decision list specifying cuts, transitions, etc.", "Post-Production"},
		{10, "Playlist", "Upload final videos to YouTube and organize them into a playlist.", "Publishing"},
		{11, "Thumbnail", "Design compelling video thumbnails that drive clicks.", "Publishing"},
	}

	for _, stage := range stages {
		fmt.Printf("Processing Stage %02d: %s...\n", stage.ID, stage.Title)
		
		prompt, err := generatePrompt(geminiKey, stage)
		if err != nil {
			log.Printf("Error generating prompt for stage %d: %v", stage.ID, err)
			continue
		}
		
		fmt.Printf("Generated Prompt: %s\n", prompt)
		
		savePrompt(stage.ID, stage.Title, prompt)
	}

	fmt.Println("Pipeline image generation trigger completed.")
}

func generatePrompt(apiKey string, stage Stage) (string, error) {
	prompt := fmt.Sprintf(`You are an expert AI Image Prompt Engineer.
Create a high-fidelity, detailed image generation prompt for Stage %02d of a course production pipeline.
Stage Title: %s
Description: %s
Phase: %s

Style Requirements:
- Theme: "Cloud Architect Certification"
- Visual Style: Minimalist, dark corporate, glassmorphism, tech-focused, professional.
- Color Palette: Deep blues, neon accents (Violet for Pre-Prod, Cyan for Prod, Green for Post-Prod, Amber for Publishing).
- Aspect Ratio: 16:9.
- Content: An abstract or symbolic representation of the stage. High-tech, futuristic aesthetic.

Return ONLY the refined prompt text, no commentary.`, stage.ID, stage.Title, stage.Desc, stage.Phase)

	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent?key=" + apiKey
	reqBody := map[string]any{
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

	b, _ := json.Marshal(reqBody)
	resp, err := http.Post(url, "application/json", bytes.NewReader(b))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("gemini error %d: %s", resp.StatusCode, body)
	}

	var gResp struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&gResp); err != nil {
		return "", err
	}

	if len(gResp.Candidates) > 0 && len(gResp.Candidates[0].Content.Parts) > 0 {
		return strings.TrimSpace(gResp.Candidates[0].Content.Parts[0].Text), nil
	}

	return "", fmt.Errorf("no prompt generated")
}

func savePrompt(id int, title, prompt string) {
	filename := fmt.Sprintf("3_Simulation/generated/pipeline/%02d_%s_prompt.txt", id, strings.ReplaceAll(strings.ToLower(title), " ", "_"))
	_ = os.WriteFile(filename, []byte(prompt), 0644)
}
