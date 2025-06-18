package firebase

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

// AppClient holds Firebase client instances
type AppClient struct {
	App       *firebase.App
	Auth      *auth.Client
	Firestore *firestore.Client
}

var (
	// Client is a global Firebase client
	Client *AppClient
)

// InitFirebase initializes Firebase services
func InitFirebase() (*AppClient, error) {
	ctx := context.Background()

	// Get service account path from environment
	serviceAccountPath := os.Getenv("FIREBASE_SERVICE_ACCOUNT")
	if serviceAccountPath == "" {
		// Look for service account in common locations
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

	var app *firebase.App
	var err error

	if serviceAccountPath != "" {
		// Initialize with service account file if available
		opt := option.WithCredentialsFile(serviceAccountPath)
		app, err = firebase.NewApp(ctx, nil, opt)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize Firebase with service account: %w", err)
		}
	} else {
		// Initialize with default credentials (useful for development or when running on GCP)
		app, err = firebase.NewApp(ctx, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize Firebase with default credentials: %w", err)
		}
	}

	// Initialize Auth client
	authClient, err := app.Auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Auth client: %w", err)
	}

	// Initialize Firestore client
	firestoreClient, err := app.Firestore(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Firestore client: %w", err)
	}

	client := &AppClient{
		App:       app,
		Auth:      authClient,
		Firestore: firestoreClient,
	}

	// Set the global client
	Client = client

	return client, nil
}

// CloseFirebase closes Firebase connections
func CloseFirebase() error {
	if Client != nil && Client.Firestore != nil {
		return Client.Firestore.Close()
	}
	return nil
}
