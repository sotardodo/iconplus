package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type vaultLoginResponse struct {
	Auth struct {
		ClientToken string `json:"client_token"`
	} `json:"auth"`
}

type vaultSecretResponse struct {
	Data struct {
		Data map[string]string `json:"data"`
	} `json:"data"`
}

func applyVaultSecrets() {
	addr := os.Getenv("VAULT_ADDR")
	if addr == "" {
		return
	}

	token, err := resolveVaultToken(addr)
	if err != nil {
		log.Printf("Vault auth skipped: %v", err)
		return
	}

	secrets, err := readVaultSecret(addr, token, "golang")
	if err != nil {
		log.Printf("Vault secret read skipped: %v", err)
		return
	}

	if password := secrets["db_password"]; password != "" {
		os.Setenv("DB_PASSWORD", password)
		log.Println("DB password loaded from Vault")
	}
}

func resolveVaultToken(addr string) (string, error) {
	if token := os.Getenv("VAULT_TOKEN"); token != "" {
		return token, nil
	}

	jwtBytes, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/token")
	if err != nil {
		return "", fmt.Errorf("kubernetes service account token unavailable")
	}

	payload := map[string]string{
		"role": getEnv("VAULT_ROLE", "golang-app"),
		"jwt":  string(jwtBytes),
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Post(addr+"/v1/auth/kubernetes/login", "application/json", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		responseBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("vault login failed: %s", string(responseBody))
	}

	var login vaultLoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&login); err != nil {
		return "", err
	}
	if login.Auth.ClientToken == "" {
		return "", fmt.Errorf("vault login returned empty token")
	}

	return login.Auth.ClientToken, nil
}

func readVaultSecret(addr, token, path string) (map[string]string, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest(http.MethodGet, addr+"/v1/secret/data/"+path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Vault-Token", token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		responseBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("vault secret read failed: %s", string(responseBody))
	}

	var secret vaultSecretResponse
	if err := json.NewDecoder(resp.Body).Decode(&secret); err != nil {
		return nil, err
	}

	return secret.Data.Data, nil
}
