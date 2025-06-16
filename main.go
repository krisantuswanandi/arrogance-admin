package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"arrogance/firebase"

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

	// Spinner characters
	spinnerChars = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
)

// Custom message types
type firebaseInitMsg struct {
	client   *firebase.AppClient
	authSvc  *firebase.AuthService
	storeSvc *firebase.FirestoreService
}

type firebaseErrorMsg struct {
	err error
}

// TickMsg is a message that's sent when the timer ticks
type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Model represents the application state
type Model struct {
	title      string
	message    string
	firebase   *firebase.AppClient
	authSvc    *firebase.AuthService
	storeSvc   *firebase.FirestoreService
	loading    bool
	error      string
	spinnerIdx int
}

// Initialize the application
func (m Model) Init() tea.Cmd {
	// Skip Firebase initialization in test mode
	if IsTestMode() {
		return nil
	}

	return tea.Batch(
		initFirebase(),
		tick(),
	)
}

// initFirebase is a command that initializes Firebase
func initFirebase() tea.Cmd {
	return func() tea.Msg {
		// Initialize Firebase
		firebaseClient, err := firebase.InitFirebase()
		if err != nil {
			return firebaseErrorMsg{err: err}
		}

		// Create services
		authSvc := firebase.NewAuthService(firebaseClient.Auth)
		storeSvc := firebase.NewFirestoreService(firebaseClient.Firestore)

		return firebaseInitMsg{
			client:   firebaseClient,
			authSvc:  authSvc,
			storeSvc: storeSvc,
		}
	}
}

// Update handles events and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	case firebaseInitMsg:
		// Update model with Firebase client and services
		m.firebase = msg.client
		m.authSvc = msg.authSvc
		m.storeSvc = msg.storeSvc
		m.loading = false
		m.message = "This is a TUI application built with Charm. Firebase initialized successfully!"
		return m, nil
	case firebaseErrorMsg:
		// Update model with Firebase initialization error
		m.loading = false
		m.error = fmt.Sprintf("Failed to initialize Firebase: %v", msg.err)
		return m, nil
	case tickMsg:
		// Continue ticking to refresh UI while loading
		if m.loading {
			// Update spinner index
			m.spinnerIdx = (m.spinnerIdx + 1) % len(spinnerChars)
			return m, tick()
		}
	}

	return m, nil
}

// View renders the UI
func (m Model) View() string {
	s := titleStyle.Render(m.title) + "\n\n"

	if m.loading {
		loadingStyle := normalStyle.Foreground(lipgloss.Color("#FFFF00"))
		spinner := spinnerChars[m.spinnerIdx]
		s += loadingStyle.Render(spinner+" Loading... Please wait while Firebase initializes.") + "\n\n"

		// Add spinner animation
		spinnerFrame := spinnerChars[m.spinnerIdx%len(spinnerChars)]
		s += loadingStyle.Render(fmt.Sprintf(" %s  ", spinnerFrame)) + "\n\n"
		m.spinnerIdx += 1
	} else if m.error != "" {
		errorStyle := normalStyle.Foreground(lipgloss.Color("#FF0000"))
		s += errorStyle.Render("Error: "+m.error) + "\n\n"
	} else {
		s += normalStyle.Render(m.message) + "\n\n" // Show Firebase connection status
		if m.firebase != nil {
			statusStyle := normalStyle.Foreground(lipgloss.Color("#00FF00"))
			s += statusStyle.Render("✓ Firebase connected") + "\n"

			if m.authSvc != nil {
				s += statusStyle.Render("✓ Auth service ready") + "\n"
			}

			if m.storeSvc != nil {
				s += statusStyle.Render("✓ Firestore service ready") + "\n"
			}

			s += "\n"
		}
	}

	s += normalStyle.Render("Press 'q' or Ctrl+C to quit.") + "\n"
	return s
}

func main() {
	// Setup signal handling for graceful shutdown
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	// Initialize model with loading state
	m := Model{
		title:      "Welcome to Arrogance Admin",
		message:    "Initializing Firebase...",
		loading:    true,
		spinnerIdx: 0,
	}

	// Start the application
	p := tea.NewProgram(m)

	// Handle cleanup on exit
	go func() {
		<-signalCh
		fmt.Println("Shutting down...")
		_ = firebase.CloseFirebase()
		os.Exit(0)
	}()

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		_ = firebase.CloseFirebase()
		os.Exit(1)
	}

	// Clean up Firebase on exit
	_ = firebase.CloseFirebase()
}
