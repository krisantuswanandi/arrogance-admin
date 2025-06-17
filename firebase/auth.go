package firebase

import (
	"context"
	"errors"
	"log"
	"os"

	"firebase.google.com/go/v4/auth"
)

// AuthService provides authentication-related functionality
type AuthService struct {
	client *auth.Client
}

// NewAuthService creates a new AuthService
func NewAuthService(client *auth.Client) *AuthService {
	return &AuthService{
		client: client,
	}
}

// VerifyIDToken verifies a Firebase ID token
func (s *AuthService) VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error) {
	if s.client == nil {
		return nil, errors.New("auth client not initialized")
	}
	return s.client.VerifyIDToken(ctx, idToken)
}

// GetUser gets a user by their UID
func (s *AuthService) GetUser(ctx context.Context, uid string) (*auth.UserRecord, error) {
	if s.client == nil {
		return nil, errors.New("auth client not initialized")
	}
	return s.client.GetUser(ctx, uid)
}

// ListUsers lists users with pagination
func (s *AuthService) ListUsers(ctx context.Context, maxResults uint32, pageToken string) (*auth.UserIterator, error) {
	if s.client == nil {
		log.Println("Auth client is nil")
		return nil, errors.New("auth client not initialized")
	}

	log.Println("Auth client is available, creating user iterator")

	// Print the GOOGLE_APPLICATION_CREDENTIALS environment variable
	if creds := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"); creds != "" {
		log.Printf("Using GOOGLE_APPLICATION_CREDENTIALS: %s", creds)
	}

	// Using the Users method to get an iterator
	iter := s.client.Users(ctx, pageToken)
	log.Println("User iterator created successfully")

	return iter, nil
}

// CreateUser creates a new user
func (s *AuthService) CreateUser(ctx context.Context, params *auth.UserToCreate) (string, error) {
	if s.client == nil {
		return "", errors.New("auth client not initialized")
	}
	user, err := s.client.CreateUser(ctx, params)
	if err != nil {
		return "", err
	}
	return user.UID, nil
}

// UpdateUser updates a user
func (s *AuthService) UpdateUser(ctx context.Context, uid string, params *auth.UserToUpdate) error {
	if s.client == nil {
		return errors.New("auth client not initialized")
	}
	_, err := s.client.UpdateUser(ctx, uid, params)
	return err
}

// DeleteUser deletes a user
func (s *AuthService) DeleteUser(ctx context.Context, uid string) error {
	if s.client == nil {
		return errors.New("auth client not initialized")
	}
	return s.client.DeleteUser(ctx, uid)
}
