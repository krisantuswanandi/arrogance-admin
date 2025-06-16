package main

import (
	"arrogance/firebase"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// MockModel is a version of Model used for testing
type MockModel struct {
	Model
}

// Init overrides the real Init method for testing
func (m MockModel) Init() tea.Cmd {
	// Return nil instead of actual Firebase initialization for tests
	return nil
}

// NewMockModel creates a Model for testing without real Firebase initialization
func NewMockModel() MockModel {
	return MockModel{
		Model: Model{
			title:      "Test Title",
			message:    "Test Message",
			loading:    false,
			spinnerIdx: 0,
		},
	}
}

// CreateMockFirebaseClient creates a mock Firebase client for testing
func CreateMockFirebaseClient() *firebase.AppClient {
	return &firebase.AppClient{}
}

// CreateMockAuthService creates a mock Auth service for testing
func CreateMockAuthService() *firebase.AuthService {
	return &firebase.AuthService{}
}

// CreateMockFirestoreService creates a mock Firestore service for testing
func CreateMockFirestoreService() *firebase.FirestoreService {
	return &firebase.FirestoreService{}
}

// UpdateTestModel updates a test model with the necessary test values
func UpdateTestModel(t *testing.T, m *Model) {
	// Initialize model with test values
	m.spinnerIdx = 0
}
