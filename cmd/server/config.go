package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type config struct {
	supabaseURL       string
	supabaseAnon      string
	axiomToken        string
	axiomDataset      string
	axiomAPIURL       string
	axiomQueryURL     string
	port              string
	azureAccountName  string
	azureAccountKey   string
	azureKeyVaultName string
	azureTenantID     string
	azureClientID     string
	azureClientSecret string
	googleClientID    string
}

func loadDotEnv(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		val = strings.Trim(val, "'\"")
		if os.Getenv(key) == "" {
			os.Setenv(key, val)
		}
	}
}

func loadConfig() config {
	loadDotEnv(".env")
	cfg := config{
		supabaseURL:       mustEnv("SUPABASE_URL"),
		supabaseAnon:      mustEnv("SUPABASE_ANON_KEY"),
		axiomToken:        os.Getenv("AXIOM_TOKEN"),
		axiomDataset:      mustEnv("AXIOM_DATASET"),
		axiomAPIURL:       envOr("AXIOM_API_URL", "https://api.axiom.co"),
		axiomQueryURL:     envOr("AXIOM_QUERY_URL", "https://api.axiom.co"),
		port:              envOr("PORT", "8080"),
		azureKeyVaultName: os.Getenv("AZURE_KEYVAULT_NAME"),
		azureTenantID:     os.Getenv("AZURE_TENANT_ID"),
		azureClientID:     os.Getenv("AZURE_CLIENT_ID"),
		azureClientSecret: os.Getenv("AZURE_CLIENT_SECRET"),
	}
	if connStr := os.Getenv("AZURE_STORAGE_CONN_STR"); connStr != "" {
		cfg.azureAccountName, cfg.azureAccountKey = parseStorageConnStr(connStr)
	}
	cfg.googleClientID = firstNonEmpty(
		cfg.getSecret("claude-architect-GOOGLE-CLIENT-ID"),
		cfg.getSecret("google-oauth-client-id"),
		cfg.getSecret("GOOGLE_CLIENT_ID"),
	)
	return cfg
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}

func getSecretFromKeyVault(vaultName, tenantID, clientID, clientSecret, secretName string) (string, error) {
	secretName = strings.ReplaceAll(strings.ToLower(secretName), "_", "-")

	tokenURL := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", tenantID)
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("scope", "https://vault.azure.net/.default")

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("oauth token error (%d): %s", resp.StatusCode, b)
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", err
	}

	secretURL := fmt.Sprintf("https://%s.vault.azure.net/secrets/%s?api-version=7.4", vaultName, secretName)
	reqSecret, err := http.NewRequest("GET", secretURL, nil)
	if err != nil {
		return "", err
	}
	reqSecret.Header.Set("Authorization", "Bearer "+tokenResp.AccessToken)

	respSecret, err := http.DefaultClient.Do(reqSecret)
	if err != nil {
		return "", err
	}
	defer respSecret.Body.Close()

	if respSecret.StatusCode >= 400 {
		b, _ := io.ReadAll(respSecret.Body)
		return "", fmt.Errorf("keyvault get secret error (%d): %s", respSecret.StatusCode, b)
	}

	var secretResp struct {
		Value string `json:"value"`
	}
	if err := json.NewDecoder(respSecret.Body).Decode(&secretResp); err != nil {
		return "", err
	}

	return secretResp.Value, nil
}

func (c config) getSecret(secretName string) string {
	if c.azureKeyVaultName != "" && c.azureTenantID != "" && c.azureClientID != "" && c.azureClientSecret != "" {
		val, err := getSecretFromKeyVault(c.azureKeyVaultName, c.azureTenantID, c.azureClientID, c.azureClientSecret, secretName)
		if err == nil && val != "" {
			log.Printf("Successfully loaded secret '%s' from Key Vault '%s'", secretName, c.azureKeyVaultName)
			return val
		}
		log.Printf("Keyvault getSecret failed for %s, falling back to env: %v", secretName, err)
	}
	return os.Getenv(secretName)
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("required env var %s is not set", key)
	}
	return v
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func maskSecret(v string) string {
	if v == "" {
		return ""
	}
	n := len(v)
	if n <= 8 {
		return fmt.Sprintf("•••• (%d chars)", n)
	}
	return fmt.Sprintf("%s••••%s (%d chars)", v[:2], v[n-2:], n)
}

func parseStorageConnStr(connStr string) (accountName, accountKey string) {
	for _, part := range strings.Split(connStr, ";") {
		k, v, ok := strings.Cut(part, "=")
		if !ok {
			continue
		}
		switch k {
		case "AccountName":
			accountName = v
		case "AccountKey":
			accountKey = v
		}
	}
	return
}
