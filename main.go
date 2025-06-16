package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	// Styles
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FF5F87")).
			Background(lipgloss.Color("#383838")).
			PaddingLeft(4).
			PaddingRight(4).
			MarginTop(1).
			MarginBottom(1)

	normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			MarginLeft(2)
)

// Model represents the application state
type Model struct {
	title   string
	message string
}

// Initialize the application
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles events and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

// View renders the UI
func (m Model) View() string {
	s := titleStyle.Render(m.title) + "\n\n"
	s += normalStyle.Render(m.message) + "\n\n"
	s += normalStyle.Render("Press 'q' or Ctrl+C to quit.") + "\n"
	return s
}

func main() {
	m := Model{
		title:   "Welcome to Arrogance Admin",
		message: "This is a TUI application built with Charm.",
	}

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
