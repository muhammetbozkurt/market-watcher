package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type ServiceAccount struct {
	Type                    string `json:"type"`
	ProjectID               string `json:"project_id"`
	PrivateKeyID            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientID                string `json:"client_id"`
	AuthURI                 string `json:"auth_uri"`
	TokenURI                string `json:"token_uri"`
	AuthProviderX509CertURL string `json:"auth_provider_x509_cert_url"`
	ClientX509CertURL       string `json:"client_x509_cert_url"`
	UniverseDomain          string `json:"universe_domain"`
}

type Config struct {
	BigQueryProjectID      string
	BigQueryServiceAccount ServiceAccount
	GeminiAPIKey           string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	bigqueryServiceAccount := os.Getenv("BIGQUERY_SERVICE_ACCOUNT")
	if bigqueryServiceAccount == "" {
		return nil, fmt.Errorf("BIGQUERY_SERVICE_ACCOUNT is not set")
	}

	bigqueryProjectID := os.Getenv("BIGQUERY_PROJECT_ID")
	if bigqueryProjectID == "" {
		return nil, fmt.Errorf("BIGQUERY_PROJECT_ID is not set")
	}

	geminiAPIKey := os.Getenv("GEMINI_API_KEY")
	if geminiAPIKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY is not set")
	}

	var ServiceAccount ServiceAccount
	err = json.Unmarshal([]byte(bigqueryServiceAccount), &ServiceAccount)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal bigquery service account: %w", err)
	}

	return &Config{
		BigQueryProjectID:      bigqueryProjectID,
		BigQueryServiceAccount: ServiceAccount,
		GeminiAPIKey:           geminiAPIKey,
	}, nil
}
