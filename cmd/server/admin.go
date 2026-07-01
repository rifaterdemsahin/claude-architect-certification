package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

func adminBackupSupabaseHandler(cfg config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !isAdminRequest(r) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		dbURL := cfg.getSecret("SUPABASE_DB_URL")
		if dbURL == "" {
			dbURL = os.Getenv("SUPABASE_DB_URL")
		}

		var buf bytes.Buffer
		var dumped bool
		if dbURL != "" && !strings.Contains(dbURL, "[password]") {
			cmd := exec.Command("pg_dump", dbURL, "--no-owner", "--no-acl", "--schema=public")
			cmd.Stdout = &buf
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err == nil && buf.Len() > 0 {
				dumped = true
			} else {
				log.Printf("pg_dump failed or returned empty: %v", err)
			}
		}

		if !dumped {
			log.Println("Falling back to snapshotting local SQL files...")
			buf.Reset()
			buf.WriteString("-- Fallback snapshot from local schema.sql and seed.sql\n\n")
			if schema, err := os.ReadFile("5_Symbols/supabase/schema.sql"); err == nil {
				buf.Write(schema)
			}
			if seed, err := os.ReadFile("5_Symbols/supabase/seed.sql"); err == nil {
				buf.WriteString("\n\n-- SEED DATA --\n\n")
				buf.Write(seed)
			}
		}

		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", "attachment; filename=\"supabase_backup.sql\"")
		w.Write(buf.Bytes())
	}
}

func adminBackupAzureHandler(cfg config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !isAdminRequest(r) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		secrets := map[string]string{
			"SUPABASE_DB_URL":                    cfg.getSecret("SUPABASE_DB_URL"),
			"OPENROUTER_API_KEY":                 cfg.getSecret("OPENROUTER_API_KEY"),
			"GEMINI_API_KEY":                     cfg.getSecret("GEMINI_API_KEY"),
			"claude-architect-GOOGLE-CLIENT-ID":  cfg.getSecret("claude-architect-GOOGLE-CLIENT-ID"),
			"GOOGLE-IMAGEN-API-KEY":              cfg.getSecret("GOOGLE-IMAGEN-API-KEY"),
			"AXIOM_TOKEN":                        cfg.getSecret("AXIOM_TOKEN"),
			"AXIOM_ORG_ID":                       cfg.getSecret("AXIOM_ORG_ID"),
			"AdminPassword":                      cfg.getSecret("AdminPassword"),
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Disposition", "attachment; filename=\"azure_secrets_backup.json\"")
		json.NewEncoder(w).Encode(secrets)
	}
}

func adminStatusHandler(cfg config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-store")
		json.NewEncoder(w).Encode(map[string]any{"admin": isAdminRequest(r)})
	}
}

func adminLoginHandler(cfg config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req struct {
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		adminPass := cfg.getSecret("ClaudeCertificateSiteAdminPassword")
		if adminPass == "" {
			adminPass = cfg.getSecret("AdminPassword")
		}
		if adminPass == "" {
			adminPass = "admin"
		}

		if req.Password != adminPass {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"success":false,"error":"Invalid password"}`))
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "admin_logged_in",
			Value:    "true",
			Path:     "/",
			MaxAge:   3600 * 24,
			HttpOnly: false,
			SameSite: http.SameSiteLaxMode,
		})

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"success":true}`))
	}
}

func adminLogoutHandler(cfg config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "admin_logged_in",
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			Expires:  time.Unix(0, 0),
			HttpOnly: false,
			SameSite: http.SameSiteLaxMode,
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"success":true}`))
	}
}

func adminGDriveCredentialsHandler(cfg config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !isAdminRequest(r) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error":"Unauthorized"}`))
			return
		}

		clientID := cfg.getSecret("claude-architect-GOOGLE-CLIENT-ID")
		if clientID == "" {
			clientID = cfg.getSecret("google-oauth-client-id")
		}
		if clientID == "" {
			clientID = cfg.getSecret("GOOGLE_CLIENT_ID")
		}

		apiKey := cfg.getSecret("GOOGLE-IMAGEN-API-KEY")
		if apiKey == "" {
			apiKey = cfg.getSecret("GOOGLE-SEARCH-API-KEY")
		}
		if apiKey == "" {
			apiKey = cfg.getSecret("GOOGLE_API_KEY")
		}

		saveCredentialsToDotEnv(clientID, apiKey)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"clientId": clientID,
			"apiKey":   apiKey,
		})
	}
}

func saveCredentialsToDotEnv(clientId, apiKey string) {
	if clientId == "" || apiKey == "" {
		return
	}
	content, err := os.ReadFile(".env")
	if err != nil {
		content = []byte{}
	}
	lines := strings.Split(string(content), "\n")
	hasClientID := false
	hasAPIKey := false
	for i, line := range lines {
		if strings.HasPrefix(line, "GOOGLE_CLIENT_ID=") {
			lines[i] = "GOOGLE_CLIENT_ID=" + clientId
			hasClientID = true
		}
		if strings.HasPrefix(line, "GOOGLE_API_KEY=") {
			lines[i] = "GOOGLE_API_KEY=" + apiKey
			hasAPIKey = true
		}
	}
	if !hasClientID {
		lines = append(lines, "GOOGLE_CLIENT_ID="+clientId)
	}
	if !hasAPIKey {
		lines = append(lines, "GOOGLE_API_KEY="+apiKey)
	}
	os.WriteFile(".env", []byte(strings.Join(lines, "\n")), 0644)
}
