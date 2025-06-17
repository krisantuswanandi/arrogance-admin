package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// CheckServiceAccount verifies that the service account file exists and is valid
func CheckServiceAccount() (string, error) {
	// Check for environment variable with service account path
	serviceAccountPath := os.Getenv("FIREBASE_SERVICE_ACCOUNT")
	if serviceAccountPath == "" {
		// Default to looking for the service account in the current directory or config directory
		homeDir, err := os.UserHomeDir()
		if err == nil {
			// Try in the config directory first
			configPath := filepath.Join(homeDir, ".config", "arrogance", "service-account.json")
			if _, err := os.Stat(configPath); err == nil {
				serviceAccountPath = configPath
				log.Printf("Found service account at: %s", configPath)
			} else {
				log.Printf("No service account found at: %s (%v)", configPath, err)
			}
		}

		// If still not found, try in the current directory
		if serviceAccountPath == "" {
			localPath := "service-account.json"
			if _, err := os.Stat(localPath); err == nil {
				serviceAccountPath = localPath
				log.Printf("Found service account at: %s", localPath)
			} else {
				log.Printf("No service account found at: %s (%v)", localPath, err)
			}
		}
	}

	if serviceAccountPath == "" {
		return "", fmt.Errorf("no service account file found")
	}

	// Verify the service account file is readable and valid JSON
	content, err := os.ReadFile(serviceAccountPath)
	if err != nil {
		return "", fmt.Errorf("could not read service account file: %w", err)
	}

	// Check if it's valid JSON
	var js map[string]interface{}
	if err := json.Unmarshal(content, &js); err != nil {
		return "", fmt.Errorf("service account file is not valid JSON: %w", err)
	}

	// Check for required fields
	requiredFields := []string{"type", "project_id", "private_key_id", "private_key", "client_email"}
	for _, field := range requiredFields {
		if _, ok := js[field]; !ok {
			return "", fmt.Errorf("service account file is missing required field: %s", field)
		}
	}

	// Ensure it's a service account
	if js["type"] != "service_account" {
		return "", fmt.Errorf("file is not a service account (type: %s)", js["type"])
	}

	log.Printf("Service account file validated successfully. Project ID: %s", js["project_id"])
	return serviceAccountPath, nil
}
