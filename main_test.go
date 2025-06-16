package main

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func init() {
	// Enable test mode for all tests
	SetTestMode(true)
}

func TestModelInit(t *testing.T) {
	defer SetTestMode(false)

	m := Model{
		title:      "Test Title",
		message:    "Test Message",
		loading:    false,
		spinnerIdx: 0,
	}

	cmd := m.Init()
	if cmd != nil {
		t.Error("Expected Init to return nil command in test mode, got non-nil")
	}
}

func TestModelUpdate(t *testing.T) {
	m := Model{
		title:      "Test Title",
		message:    "Test Message",
		loading:    false,
		spinnerIdx: 0,
	}

	// Test case 1: Quit on 'q' key press
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	updatedModel, cmd := m.Update(keyMsg)

	// Check the updated model
	updatedM, ok := updatedModel.(Model)
	if !ok {
		t.Error("Failed to cast updated model to Model type")
		return
	}

	// Model should remain unchanged
	if updatedM.title != m.title || updatedM.message != m.message {
		t.Errorf("Model was unexpectedly changed. Got title=%s, message=%s",
			updatedM.title, updatedM.message)
	}

	// Command should be tea.Quit
	if cmd == nil {
		t.Error("Expected quit command, got nil")
	}

	// Test case 2: Quit on 'ctrl+c'
	keyMsg = tea.KeyMsg{Type: tea.KeyCtrlC}
	updatedModel, cmd = m.Update(keyMsg)

	// Check the updated model
	updatedM, ok = updatedModel.(Model)
	if !ok {
		t.Error("Failed to cast updated model to Model type")
		return
	}

	// Model should remain unchanged
	if updatedM.title != m.title || updatedM.message != m.message {
		t.Errorf("Model was unexpectedly changed. Got title=%s, message=%s",
			updatedM.title, updatedM.message)
	}

	// Command should be tea.Quit
	if cmd == nil {
		t.Error("Expected quit command, got nil")
	}

	// Test case 3: Other key should not quit
	keyMsg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
	updatedModel, cmd = m.Update(keyMsg)

	// Check that command is nil (no quit)
	if cmd != nil {
		t.Error("Expected nil command for non-quit key, got non-nil")
	}
}

func TestModelView(t *testing.T) {
	// Create a simple model with loading=false to test basic view
	m := Model{
		title:      "Test Title",
		message:    "Test Message",
		spinnerIdx: 0,
		loading:    false,
	}

	view := m.View()

	// Check that view contains the title
	if !strings.Contains(view, "Test Title") {
		t.Error("View does not contain the title")
	}

	// Check that view contains the message
	if !strings.Contains(view, "Test Message") {
		t.Error("View does not contain the message")
	}

	// Check that view contains the quit instructions
	if !strings.Contains(view, "Press 'q' or Ctrl+C to quit.") {
		t.Error("View does not contain quit instructions")
	}
}

// TestNewModel tests the creation of a model with expected default values
func TestNewModel(t *testing.T) {
	// Create a model the same way as in main()
	m := Model{
		title:      "Welcome to Arrogance Admin",
		message:    "This is a TUI application built with Charm.",
		spinnerIdx: 0,
	}

	// Verify the model has the expected values
	if m.title != "Welcome to Arrogance Admin" {
		t.Errorf("Expected title to be 'Welcome to Arrogance Admin', got '%s'", m.title)
	}

	if m.message != "This is a TUI application built with Charm." {
		t.Errorf("Expected message to be 'This is a TUI application built with Charm.', got '%s'", m.message)
	}
}
