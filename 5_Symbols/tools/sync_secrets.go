package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {
	vaultName := os.Getenv("AZURE_KEYVAULT_NAME")
	tenantID := os.Getenv("AZURE_TENANT_ID")
	clientID := os.Getenv("AZURE_CLIENT_ID")
	clientSecret := os.Getenv("AZURE_CLIENT_SECRET")

	if vaultName == "" || tenantID == "" || clientID == "" || clientSecret == "" {
		log.Fatal("Missing Azure Key Vault environment variables")
	}

	secretNames := []string{"GEMINI_API_KEY", "SUPABASE_URL", "SUPABASE_ANON_KEY", "AXIOM_TOKEN", "AZURE_STORAGE_CONN_STR"}
	
	envContent := ""
	if _, err := os.Stat(".env"); err == nil {
		content, _ := os.ReadFile(".env")
		envContent = string(content)
	}

	for _, name := range secretNames {
		fmt.Printf("Fetching %s...\n", name)
		val, err := getSecret(vaultName, tenantID, clientID, clientSecret, name)
		if err != nil {
			log.Printf("Failed to fetch %s: %v", name, err)
			continue
		}
		
		line := fmt.Sprintf("%s=%s", name, val)
		if strings.Contains(envContent, name+"=") {
			// Update existing
			lines := strings.Split(envContent, "\n")
			for i, l := range lines {
				if strings.HasPrefix(l, name+"=") {
					lines[i] = line
					break
				}
			}
			envContent = strings.Join(lines, "\n")
		} else {
			// Append new
			if envContent != "" && !strings.HasSuffix(envContent, "\n") {
				envContent += "\n"
			}
			envContent += line + "\n"
		}
	}

	err := os.WriteFile(".env", []byte(envContent), 0600)
	if err != nil {
		log.Fatalf("Failed to write .env: %v", err)
	}
	fmt.Println("Successfully updated .env with secrets from Key Vault")
}

func getSecret(vaultName, tenantID, clientID, clientSecret, secretName string) (string, error) {
	// Key Vault secret names only allow alphanumeric and hyphens
	kvSecretName := strings.ReplaceAll(strings.ToLower(secretName), "_", "-")

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

	secretURL := fmt.Sprintf("https://%s.vault.azure.net/secrets/%s?api-version=7.4", vaultName, kvSecretName)
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
