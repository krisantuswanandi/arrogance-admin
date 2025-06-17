package main

import (
	"fmt"
	"os"
)

func main() {
	// Enable verbose logging
	os.Setenv("FIREBASE_DEBUG", "true")

	// Check if the service account file exists and is valid
	serviceAccountPath, err := CheckServiceAccount()
	if err != nil {
		fmt.Printf("Service account validation error: %v\n", err)
		os.Exit(1)
	}

	// Set the service account path in environment
	fmt.Printf("Using validated service account at: %s\n", serviceAccountPath)
	os.Setenv("FIREBASE_SERVICE_ACCOUNT", serviceAccountPath)

	// Also set GOOGLE_APPLICATION_CREDENTIALS which is used by the Firebase Admin SDK
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", serviceAccountPath)

	// Call the real main function
	realMain()
}
