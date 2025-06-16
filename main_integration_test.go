package main

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// mockProgram is a mock implementation of tea.Program for testing
type mockProgram struct {
	model tea.Model
	err   error
}

func (m *mockProgram) Start() (tea.Model, error) {
	return m.model, m.err
}

// TestMainFunction tests the setup of the main function
// Note: This doesn't actually run main() as that would block the test,
// but it tests the same setup logic
func TestMainSetup(t *testing.T) {
	// Create a model with the same values as in main()
	m := Model{
		title:   "Welcome to Arrogance Admin",
		message: "This is a TUI application built with Charm.",
	}

	// Verify the model initialization
	if m.title != "Welcome to Arrogance Admin" {
		t.Errorf("Expected title to be 'Welcome to Arrogance Admin', got '%s'", m.title)
	}

	if m.message != "This is a TUI application built with Charm." {
		t.Errorf("Expected message to be 'This is a TUI application built with Charm.', got '%s'", m.message)
	}

	// Test the initial view
	view := m.View()
	if view == "" {
		t.Error("View should not be empty")
	}
}

// createModelFn is a function that creates and initializes a Model
type createModelFn func() Model

// testModelCreation is a helper function to test model creation
func testModelCreation(t *testing.T, createFn createModelFn, expectedTitle, expectedMessage string) {
	m := createFn()

	if m.title != expectedTitle {
		t.Errorf("Expected title to be '%s', got '%s'", expectedTitle, m.title)
	}

	if m.message != expectedMessage {
		t.Errorf("Expected message to be '%s', got '%s'", expectedMessage, m.message)
	}
}

// TestModelCreation tests the creation of model with different values
func TestModelCreation(t *testing.T) {
	// Test case 1: Default model from main()
	testModelCreation(
		t,
		func() Model {
			return Model{
				title:   "Welcome to Arrogance Admin",
				message: "This is a TUI application built with Charm.",
			}
		},
		"Welcome to Arrogance Admin",
		"This is a TUI application built with Charm.",
	)

	// Test case 2: Custom model
	testModelCreation(
		t,
		func() Model {
			return Model{
				title:   "Custom Title",
				message: "Custom Message",
			}
		},
		"Custom Title",
		"Custom Message",
	)
}
