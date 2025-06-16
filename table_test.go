package main

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// MockKeyMsg creates a tea.KeyMsg for the given key string
func mockKeyMsg(key string) tea.Msg {
	if key == "ctrl+c" {
		return tea.KeyMsg{Type: tea.KeyCtrlC}
	}
	
	if len(key) == 1 {
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(key)}
	}
	
	return nil
}

// TestKeyHandling is a table-driven test for key handling
func TestKeyHandling(t *testing.T) {
	tests := []struct {
		name        string
		key         string
		expectQuit  bool
	}{
		{"Quit on q", "q", true},
		{"Quit on ctrl+c", "ctrl+c", true},
		{"Don't quit on other keys", "a", false},
		{"Don't quit on other keys", "z", false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Model{
				title:   "Test",
				message: "Test",
			}
			
			_, cmd := m.Update(mockKeyMsg(tt.key))
			
			if tt.expectQuit && cmd == nil {
				t.Error("Expected quit command, got nil")
			}
			
			if !tt.expectQuit && cmd != nil {
				t.Error("Expected nil command, got non-nil")
			}
		})
	}
}

// TestStyleRendering tests that styles are applied correctly
func TestStyleRendering(t *testing.T) {
	// Reset styles to simplify testing and avoid terminal color codes
	origTitleStyle := titleStyle
	origNormalStyle := normalStyle
	
	defer func() {
		// Restore original styles after test
		titleStyle = origTitleStyle
		normalStyle = origNormalStyle
	}()
	
	// Use simplified styles for testing
	titleStyle = lipgloss.NewStyle()
	normalStyle = lipgloss.NewStyle()
	
	m := Model{
		title:   "Title",
		message: "Message",
	}
	
	view := m.View()
	
	// Check the structure of the view
	lines := strings.Split(view, "\n")
	
	if len(lines) < 5 {
		t.Fatalf("Expected at least 5 lines in view, got %d", len(lines))
	}
	
	if !strings.Contains(lines[0], "Title") {
		t.Errorf("First line should contain title, got: %s", lines[0])
	}
	
	if !strings.Contains(lines[2], "Message") {
		t.Errorf("Third line should contain message, got: %s", lines[2])
	}
	
	if !strings.Contains(lines[4], "Press 'q' or Ctrl+C to quit") {
		t.Errorf("Fifth line should contain quit instructions, got: %s", lines[4])
	}
}
