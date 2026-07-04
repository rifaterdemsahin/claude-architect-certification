package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
)

func main() {
	cfg := loadConfig()

	navConfigRaw, err := os.ReadFile("navigation_config.json")
	if err != nil {
		log.Fatalf("cannot read navigation_config.json: %v", err)
	}
	navConfigJS := template.JS(navConfigRaw)

	funcs := template.FuncMap{
		"inc": func(i int) int { return i + 1 },
	}
	tmpl := template.Must(
		template.New("").Funcs(funcs).ParseFiles(
			"5_Symbols/templates/index.html",
			"5_Symbols/templates/axiom_errors.html",
		),
	)

	fs := http.FileServer(http.Dir("."))
	mux := http.NewServeMux()
	mux.Handle("/shared/", observe(cfg, fs))
	mux.Handle("/navigation_config.json", observe(cfg, fs))
	mux.Handle("/api/config", observe(cfg, configHandler(cfg)))
	mux.Handle("/api/env-status", observe(cfg, envStatusHandler(cfg)))
	mux.Handle("/api/nav/favs", observe(cfg, navFavsHandler(cfg)))
	mux.Handle("/api/errors", observe(cfg, clientErrorsHandler(cfg)))
	mux.Handle("/api/research/upload", observe(cfg, researchUploadHandler(cfg)))
	mux.Handle("/api/research/files", observe(cfg, researchFilesHandler(cfg)))
	mux.Handle("/api/research/file", observe(cfg, researchFileHandler(cfg)))
	mux.Handle("/api/research/rename", observe(cfg, researchRenameHandler(cfg)))
	mux.Handle("/api/explanations/generate", observe(cfg, generateExplanationHandler(cfg)))
	mux.Handle("/api/explanations", observe(cfg, explanationsHandler(cfg)))
	mux.Handle("/api/images/generate", observe(cfg, imageGenerateHandler(cfg)))
	mux.Handle("/api/images/enhance-prompt", observe(cfg, imageEnhancePromptHandler(cfg)))
	mux.Handle("/api/images/save", observe(cfg, imageSaveHandler(cfg)))
	mux.Handle("/api/images/backfill-thumbnails", observe(cfg, imageThumbnailBackfillHandler(cfg)))
	mux.Handle("/api/search", observe(cfg, searchHandler(cfg)))
	mux.Handle("/api/admin/reindex-search", observe(cfg, adminReindexSearchHandler(cfg)))
	mux.Handle("/api/images/test-gemini", observe(cfg, imageTestGeminiHandler(cfg)))
	mux.Handle("/api/infographics/generate", observe(cfg, infographicGenerateHandler(cfg)))
	mux.Handle("/api/infographics/save", observe(cfg, infographicSaveHandler(cfg)))
	mux.Handle("/api/analogies/generate-prompt", observe(cfg, analogyPromptHandler(cfg)))
	mux.Handle("/api/analogies/generate", observe(cfg, analogyGenerateHandler(cfg)))
	mux.Handle("/api/analogies/save", observe(cfg, analogySaveHandler(cfg)))
	mux.Handle("/api/scripts/openrouter", observe(cfg, scriptGenerateHandler(cfg)))
	mux.Handle("/api/lowerthirds/generate", observe(cfg, lowerThirdGenerateHandler(cfg)))
	mux.Handle("/api/lowerthirds/openrouter", observe(cfg, openRouterGenerateHandler(cfg)))
	mux.Handle("/api/drawings/openrouter", observe(cfg, drawingGenerateHandler(cfg)))
	mux.Handle("/api/drawings/save", observe(cfg, drawingSaveHandler(cfg)))
	mux.Handle("/api/slides/openrouter", observe(cfg, slideGenerateHandler(cfg)))
	mux.Handle("/api/animations/generate-prompt", observe(cfg, animationPromptHandler(cfg)))
	mux.Handle("/api/animations/openrouter-prepare", observe(cfg, animationOpenRouterPrepareHandler(cfg)))
	mux.Handle("/api/animations/runpod/run", observe(cfg, animationRunpodRunHandler(cfg)))
	mux.Handle("/api/animations/runpod/status", observe(cfg, animationRunpodStatusHandler(cfg)))
	mux.Handle("/api/animations/runpod/generate-code", observe(cfg, animationRunpodGenerateCodeHandler(cfg)))
	mux.Handle("/api/ai/sanity-check", observe(cfg, sanityCheckHandler(cfg)))
	mux.Handle("/api/ai/fix-grammar", observe(cfg, fixGrammarHandler(cfg)))
	mux.Handle("/api/admin/login", observe(cfg, adminLoginHandler(cfg)))
	mux.Handle("/api/admin/logout", observe(cfg, adminLogoutHandler(cfg)))
	mux.Handle("/api/admin/status", observe(cfg, adminStatusHandler(cfg)))
	mux.Handle("/api/admin/backup/supabase", observe(cfg, adminBackupSupabaseHandler(cfg)))
	mux.Handle("/api/admin/backup/azure", observe(cfg, adminBackupAzureHandler(cfg)))
	mux.Handle("/api/axiom/logs", observe(cfg, axiomLogsHandler(cfg)))
	mux.Handle("/api/admin/gdrive-credentials", observe(cfg, adminGDriveCredentialsHandler(cfg)))
	mux.Handle("/admin/errors", observe(cfg, axiomErrorsHandler(tmpl, cfg, navConfigJS)))
	mux.Handle("/", observe(cfg, homeHandler(tmpl, cfg, navConfigJS)))

	addr := ":" + cfg.port
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
