package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sort"
	"strings"
	"syscall"
	"time"

	"arrogance/firebase"

	"firebase.google.com/go/v4/auth"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
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

	// Navigation and layout styles
	navStyle = lipgloss.NewStyle().
			Border(lipgloss.Border{
			Top:         "─",
			Bottom:      "─",
			Left:        "│",
			Right:       "│",
			TopLeft:     "┌",
			TopRight:    "┐",
			BottomLeft:  "└",
			BottomRight: "┘",
		}).
		BorderForeground(lipgloss.Color("69")).
		Padding(0, 1)

	activeTabStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("69")).
			Foreground(lipgloss.Color("69")).
			Padding(0, 2)

	inactiveTabStyle = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("240")).
				Foreground(lipgloss.Color("240")).
				Padding(0, 2)

	docStyle = lipgloss.NewStyle().Padding(1, 2, 1, 2)

	highlightColor = lipgloss.Color("69")
	// Status indicators
	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00")).
			MarginLeft(2)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			MarginLeft(2)

	loadingStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFF00")).
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

// UsersLoadedMsg is sent when users are loaded from Firebase
type usersLoadedMsg struct {
	users []*auth.UserRecord
}

// UsersErrorMsg is sent when there's an error loading users
type usersErrorMsg struct {
	err error
}

type routinesErrorMsg struct {
	err error
}

type routinesLoadedMsg struct {
	routines []map[string]interface{}
}

// TabChangeMsg is sent when the active tab changes
type tabChangeMsg struct {
	index int
}

func tick() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Model represents the application state
type Model struct {
	// Core app state
	title      string
	message    string
	firebase   *firebase.AppClient
	authSvc    *firebase.AuthService
	storeSvc   *firebase.FirestoreService
	loading    bool
	error      string
	spinnerIdx int

	// Navigation and layout
	width       int
	height      int
	activeTab   int
	tabs        []string
	currentView string

	// User components
	userTable   table.Model
	userList    []*auth.UserRecord
	userLoading bool
	userError   string

	// Routine components
	routineTable   table.Model
	routineList    []map[string]interface{}
	routineLoading bool
	routineError   string
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
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab", "right", "l":
			// Switch to next tab
			m.activeTab = (m.activeTab + 1) % len(m.tabs)
			m.currentView = m.getViewForActiveTab()

			switch m.currentView {
			case UsersView:
				m.userLoading = true
				return m, fetchUsers(m.authSvc)
			case RoutinesView:
				m.routineLoading = true
				return m, fetchRoutines(m.storeSvc)
			}

			return m, nil
		case "shift+tab", "left", "h":
			// Switch to previous tab
			m.activeTab = (m.activeTab - 1 + len(m.tabs)) % len(m.tabs)
			m.currentView = m.getViewForActiveTab()

			switch m.currentView {
			case UsersView:
				m.userLoading = true
				return m, fetchUsers(m.authSvc)
			case RoutinesView:
				m.routineLoading = true
				return m, fetchRoutines(m.storeSvc)
			}

			return m, nil
		}

		// If we're viewing the user table, pass the key to the table

	case tea.WindowSizeMsg:
		// Update the model with the new window size
		m.width = msg.Width
		m.height = msg.Height

		// Update table height based on window size
		m.userTable.SetHeight(m.height - 13) // Adjust height for header and footer

		// Keep the same view
		return m, nil

	case firebaseInitMsg:
		// Update model with Firebase client and services
		m.firebase = msg.client
		m.authSvc = msg.authSvc
		m.storeSvc = msg.storeSvc
		m.loading = false
		m.message = "Firebase initialized successfully!"
		m.currentView = m.getViewForActiveTab()

		// If we're on the Users tab, load users immediately
		if m.activeTab == UsersTab {
			m.userLoading = true
			return m, fetchUsers(m.authSvc)
		}

		return m, nil

	case firebaseErrorMsg:
		// Update model with Firebase initialization error
		m.loading = false
		m.error = fmt.Sprintf("Failed to initialize Firebase: %v", msg.err)
		m.currentView = ErrorView
		return m, nil

	case usersLoadedMsg:
		// Update model with loaded users
		m.userList = msg.users
		m.userLoading = false

		// Convert users to table rows
		rows := []table.Row{}

		sortedUserList := make([]*auth.UserRecord, len(m.userList))
		copy(sortedUserList, m.userList)

		sort.Slice(sortedUserList, func(i, j int) bool {
			return sortedUserList[i].UserMetadata.CreationTimestamp < sortedUserList[j].UserMetadata.CreationTimestamp
		})

		for _, user := range sortedUserList {
			created := time.Unix(user.UserMetadata.CreationTimestamp/1000, 0).Format("02 Jan 2006, 15:04")
			lastLogin := "-"
			if user.UserMetadata.LastLogInTimestamp > 0 {
				lastLogin = time.Unix(user.UserMetadata.LastLogInTimestamp/1000, 0).Format("02 Jan 2006, 15:04")
			}
			lastActivity := "-"
			if user.UserMetadata.LastRefreshTimestamp > 0 {
				lastActivity = time.Unix(user.UserMetadata.LastRefreshTimestamp/1000, 0).Format("02 Jan 2006, 15:04")
			}

			rows = append(rows, table.Row{
				user.UserInfo.UID,
				user.UserInfo.Email,
				user.UserInfo.DisplayName,
				created,
				lastLogin,
				lastActivity,
			})
		}

		// Update the table with the rows
		m.userTable.SetRows(rows)
		return m, nil

	case usersErrorMsg:
		// Update model with user loading error
		m.userLoading = false
		m.userError = fmt.Sprintf("Failed to load users: %v", msg.err)
		return m, nil

	case routinesLoadedMsg:
		// Update routine table
		m.routineLoading = false
		m.routineError = "Temp message" // TODO: Update with actual routines data
		return m, nil

	case routinesErrorMsg:
		// Update model with user loading error
		m.routineLoading = false
		m.routineError = fmt.Sprintf("Failed to load users: %v", msg.err)
		return m, nil

	case tabChangeMsg:
		// Update the active tab
		if msg.index >= 0 && msg.index < len(m.tabs) {
			m.activeTab = msg.index
			m.currentView = m.getViewForActiveTab()

			// If switching to Users tab and no users are loaded yet, load them
			if m.activeTab == UsersTab && len(m.userList) == 0 && !m.userLoading && m.authSvc != nil {
				m.userLoading = true
				return m, fetchUsers(m.authSvc)
			}
		}
		return m, nil

	case tickMsg:
		// Update spinner index for any view that needs animation
		m.spinnerIdx = (m.spinnerIdx + 1) % len(spinnerChars)

		// Continue ticking if we're in a loading state anywhere in the app
		if m.loading || m.userLoading {
			return m, tick()
		}
	}

	return m, tea.Batch(cmds...)
}

// getViewForActiveTab returns the view type for the current active tab
func (m Model) getViewForActiveTab() string {
	if m.loading {
		return LoadingView
	}
	if m.error != "" {
		return ErrorView
	}

	switch m.activeTab {
	case UsersTab:
		return UsersView
	case RoutinesTab:
		return RoutinesView
	default:
		return HomeView
	}
}

// View renders the UI
func (m Model) View() string {
	// Special case for tests that don't set width/height
	if m.width == 0 || m.height == 0 {
		// Simple view for testing
		var sb strings.Builder
		sb.WriteString(m.title)
		sb.WriteString("\n\n")
		sb.WriteString(m.message)
		sb.WriteString("\n\n")
		sb.WriteString("Press 'q' or Ctrl+C to quit.")
		return sb.String()
	}

	var content string
	switch m.currentView {
	case LoadingView:
		content = m.loadingView()
	case ErrorView:
		content = m.errorView()
	case UsersView:
		content = m.usersView()
	case RoutinesView:
		content = m.routinesView()
	default:
		content = m.homeView()
	}

	// Make sure the content fits within the terminal dimensions
	return lipgloss.NewStyle().
		MaxHeight(m.height).
		MaxWidth(m.width).
		Render(content)
}

// renderTabs renders the navigation tabs
func (m Model) renderTabs() string {
	var renderedTabs []string

	for i, tab := range m.tabs {
		var style lipgloss.Style
		if i == m.activeTab {
			style = activeTabStyle
		} else {
			style = inactiveTabStyle
		}
		renderedTabs = append(renderedTabs, style.Render(tab))
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
}

// loadingView shows a loading screen
func (m Model) loadingView() string {
	spinner := spinnerChars[m.spinnerIdx]
	loading := loadingStyle.Render(spinner + " Loading Firebase... Please wait.")

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		titleStyle.Render(m.title)+"\n\n"+loading,
	)
}

// errorView shows an error screen
func (m Model) errorView() string {
	errorMsg := errorStyle.Render("Error: " + m.error)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		titleStyle.Render(m.title)+"\n\n"+errorMsg+"\n\nPress 'q' to quit.",
	)
}

// homeView shows the home screen
func (m Model) homeView() string {
	// Layout
	doc := strings.Builder{}

	// Render navigation bar
	nav := m.renderTabs()
	navBar := navStyle.Width(m.width - 4).Render(nav)
	doc.WriteString(navBar)
	doc.WriteString("\n")

	// Content
	var contentText string
	if m.loading {
		// Show loading status
		spinner := spinnerChars[m.spinnerIdx]
		contentText = fmt.Sprintf("Welcome to Arrogance Admin!\n\n%s Firebase is connecting... Please wait.", spinner)
	} else if m.error != "" {
		// Show error status
		contentText = fmt.Sprintf("Welcome to Arrogance Admin!\n\nError: %s\n\nPlease check your Firebase configuration.", m.error)
	} else if m.firebase != nil {
		// Show connected status with service details
		contentText = "Welcome to Arrogance Admin!\n\nFirebase is initialized and ready to use.\n\n"

		if m.authSvc != nil {
			contentText += "✓ Auth service ready\n"
		} else {
			contentText += "✗ Auth service not available\n"
		}

		if m.storeSvc != nil {
			contentText += "✓ Firestore service ready\n"
		} else {
			contentText += "✗ Firestore service not available\n"
		}

		contentText += "\nUse the tabs above to navigate."
	} else {
		// Firebase client is null but not loading - error state
		contentText = "Welcome to Arrogance Admin!\n\nFirebase initialization failed.\n\nPlease check your configuration and restart the application."
	}

	content := lipgloss.NewStyle().
		Width(m.width-4).
		Height(m.height-10).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(1, 2).
		Render(contentText)

	doc.WriteString(content)

	// Footer
	footer := lipgloss.NewStyle().
		Width(m.width-4).
		Align(lipgloss.Left).
		Padding(0, 2).
		Render("Press 'q' to quit, tab/arrow keys to navigate")

	doc.WriteString("\n" + footer)

	// Full view
	return docStyle.Render(doc.String())
}

// usersView shows the users screen
func (m Model) usersView() string {
	// Layout
	doc := strings.Builder{}

	// Render navigation bar
	nav := m.renderTabs()
	navBar := navStyle.Width(m.width - 4).Render(nav)
	doc.WriteString(navBar)
	doc.WriteString("\n")

	// Content
	var content string
	if m.userLoading {
		// Show loading spinner
		spinner := spinnerChars[m.spinnerIdx]
		content = lipgloss.NewStyle().
			Width(m.width-8).
			Height(m.height-10).
			Padding(2, 2).
			Render(loadingStyle.Render(spinner + " Loading users..."))
	} else if m.userError != "" {
		// Show error message
		content = lipgloss.NewStyle().
			Width(m.width-8).
			Height(m.height-10).
			Padding(2, 2).
			Render(errorStyle.Render("Error loading users: " + m.userError))
	} else if len(m.userList) == 0 {
		// Show empty state
		content = lipgloss.NewStyle().
			Width(m.width-8).
			Height(m.height-10).
			Padding(2, 2).
			Render("No users found in Firebase Authentication.\n\nTo add users, use the Firebase Console or Authentication SDK.")
	} else {
		// Adjust table dimensions based on terminal size
		contentWidth := m.width - 8

		// Calculate column widths
		totalWidth := 0
		for _, col := range m.userTable.Columns() {
			totalWidth += col.Width
		}

		// Adjust column widths if necessary
		if totalWidth > contentWidth {
			ratio := float64(contentWidth) / float64(totalWidth)
			columns := m.userTable.Columns()
			for i := range columns {
				columns[i].Width = int(float64(columns[i].Width) * ratio)
			}
			m.userTable.SetColumns(columns)
		}

		tableView := m.userTable.View()
		usersCount := fmt.Sprintf("\nTotal users: %d", len(m.userList))

		content = lipgloss.NewStyle().
			Width(m.width-8).
			Padding(1, 2).
			Render(tableView + usersCount)
	}

	contentBox := lipgloss.NewStyle().
		Width(m.width - 4).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Render(content)

	doc.WriteString(contentBox)

	// Footer
	footerText := "Press 'q' to quit, tab/arrow keys to navigate"
	if !m.userLoading && m.userError == "" && len(m.userList) > 0 {
		footerText += ", up/down to select users"
	}

	footer := lipgloss.NewStyle().
		Width(m.width-4).
		Align(lipgloss.Left).
		Padding(0, 2).
		Render(footerText)

	doc.WriteString("\n" + footer)

	// Full view
	return docStyle.Render(doc.String())
}

func (m Model) routinesView() string {
	// Layout
	doc := strings.Builder{}

	// Render navigation bar
	nav := m.renderTabs()
	navBar := navStyle.Width(m.width - 4).Render(nav)
	doc.WriteString(navBar)
	doc.WriteString("\n")

	// Content
	var content string
	if m.routineLoading {
		// Show loading spinner
		spinner := spinnerChars[m.spinnerIdx]
		content = lipgloss.NewStyle().
			Width(m.width-8).
			Height(m.height-10).
			Padding(2, 2).
			Render(loadingStyle.Render(spinner + " Loading routines..."))
	} else if m.routineError != "" {
		// Show error message
		content = lipgloss.NewStyle().
			Width(m.width-8).
			Height(m.height-10).
			Padding(2, 2).
			Render(errorStyle.Render("Error loading routines: " + m.routineError))
	} else if len(m.routineList) == 0 {
		// Show empty state
		content = lipgloss.NewStyle().
			Width(m.width-8).
			Height(m.height-10).
			Padding(2, 2).
			Render("No routines found in Firestore.")
	} else {
		// Adjust table dimensions based on terminal size
		contentWidth := m.width - 8

		// Calculate column widths
		totalWidth := 0
		for _, col := range m.userTable.Columns() {
			totalWidth += col.Width
		}

		// Adjust column widths if necessary
		if totalWidth > contentWidth {
			ratio := float64(contentWidth) / float64(totalWidth)
			columns := m.userTable.Columns()
			for i := range columns {
				columns[i].Width = int(float64(columns[i].Width) * ratio)
			}
			m.userTable.SetColumns(columns)
		}

		tableView := m.userTable.View()
		usersCount := fmt.Sprintf("\nTotal users: %d", len(m.userList))

		content = lipgloss.NewStyle().
			Width(m.width-8).
			Padding(1, 2).
			Render(tableView + usersCount)
	}

	contentBox := lipgloss.NewStyle().
		Width(m.width - 4).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Render(content)

	doc.WriteString(contentBox)

	// Footer
	footerText := "Press 'q' to quit, tab/arrow keys to navigate"
	if !m.userLoading && m.userError == "" && len(m.userList) > 0 {
		footerText += ", up/down to select routines"
	}

	footer := lipgloss.NewStyle().
		Width(m.width-4).
		Align(lipgloss.Left).
		Padding(0, 2).
		Render(footerText)

	doc.WriteString("\n" + footer)

	// Full view
	return docStyle.Render(doc.String())
}

// Application constants
const (
	// Tab indices
	HomeTab     = 0
	UsersTab    = 1
	RoutinesTab = 2

	// View types for content
	LoadingView  = "loading"
	ErrorView    = "error"
	HomeView     = "home"
	UsersView    = "users"
	RoutinesView = "routines"
)

// Helper functions

// getTermSize returns the terminal size
func getTermSize() (int, int, error) {
	// Get terminal size using golang.org/x/term
	fd := int(os.Stdout.Fd())
	width, height, err := term.GetSize(fd)
	if err != nil {
		return 80, 24, err // Default fallback size
	}

	return width, height, nil
}

// initUserTable initializes the user table with appropriate columns
func initUserTable() table.Model {
	columns := []table.Column{
		{Title: "UID", Width: 25},
		{Title: "Email", Width: 30},
		{Title: "Display Name", Width: 20},
		{Title: "Created", Width: 20},
		{Title: "Last Sign In", Width: 20},
		{Title: "Last Activity", Width: 20},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithHeight(10),
		table.WithFocused(true),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(highlightColor).
		BorderBottom(true).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(highlightColor).
		Bold(true)

	t.SetStyles(s)

	return t
}

// fetchUsers fetches users from Firebase and returns a tea.Cmd
func fetchUsers(authSvc *firebase.AuthService) tea.Cmd {
	return func() tea.Msg {
		if authSvc == nil {
			return usersErrorMsg{err: errors.New("auth service not initialized")}
		}

		// In test mode, return test data
		if IsTestMode() {
			testUsers := []*auth.UserRecord{
				{
					UserInfo: &auth.UserInfo{
						UID:         "test-uid-1",
						Email:       "test1@example.com",
						DisplayName: "Test User 1",
					},
					UserMetadata: &auth.UserMetadata{
						CreationTimestamp:    time.Now().Add(-24*time.Hour).Unix() * 1000,
						LastLogInTimestamp:   time.Now().Add(-1*time.Hour).Unix() * 1000,
						LastRefreshTimestamp: time.Now().Add(-1*time.Hour).Unix() * 1000,
					},
				},
				{
					UserInfo: &auth.UserInfo{
						UID:         "test-uid-2",
						Email:       "test2@example.com",
						DisplayName: "Test User 2",
					},
					UserMetadata: &auth.UserMetadata{
						CreationTimestamp:    time.Now().Add(-48*time.Hour).Unix() * 1000,
						LastLogInTimestamp:   time.Now().Add(-2*time.Hour).Unix() * 1000,
						LastRefreshTimestamp: time.Now().Add(-2*time.Hour).Unix() * 1000,
					},
				},
			}
			return usersLoadedMsg{users: testUsers}
		}

		// Fetch users from Firebase Auth
		ctx := context.Background()

		// Try to fetch users from Firebase Auth
		iter, err := authSvc.ListUsers(ctx, 1000, "")
		if err != nil {
			// If we fail to get users, return fallback data for development
			return usersErrorMsg{err: fmt.Errorf("failed to fetch users: %w", err)}
		}

		var users []*auth.UserRecord

		// Process user records from the iterator
		for {
			user, err := iter.Next()
			if err != nil {
				// Check if we've reached the end of the iterator
				if err.Error() == "no more items in iterator" {
					break
				}

				// If there's another error, return the error
				return usersErrorMsg{err: fmt.Errorf("error iterating users: %w", err)}
			}

			// Convert ExportedUserRecord to UserRecord
			userRecord := &auth.UserRecord{
				UserInfo: &auth.UserInfo{
					UID:         user.UID,
					Email:       user.Email,
					DisplayName: user.DisplayName,
					PhoneNumber: user.PhoneNumber,
					PhotoURL:    user.PhotoURL,
				},
				UserMetadata: &auth.UserMetadata{
					CreationTimestamp:    user.UserMetadata.CreationTimestamp,
					LastLogInTimestamp:   user.UserMetadata.LastLogInTimestamp,
					LastRefreshTimestamp: user.UserMetadata.LastRefreshTimestamp,
				},
				Disabled: user.Disabled,
			}

			users = append(users, userRecord)
		}

		// If no users were found, add sample users for development
		if len(users) == 0 {
			// Add sample users to demonstrate UI functionality
			sampleUsers := []*auth.UserRecord{
				{
					UserInfo: &auth.UserInfo{
						UID:         "sample-user-1",
						Email:       "sample1@example.com",
						DisplayName: "Sample User 1",
					},
					UserMetadata: &auth.UserMetadata{
						CreationTimestamp:  time.Now().Add(-7*24*time.Hour).Unix() * 1000,
						LastLogInTimestamp: time.Now().Add(-3*time.Hour).Unix() * 1000,
					},
				},
				{
					UserInfo: &auth.UserInfo{
						UID:         "sample-user-2",
						Email:       "sample2@example.com",
						DisplayName: "Sample User 2",
					},
					UserMetadata: &auth.UserMetadata{
						CreationTimestamp:  time.Now().Add(-14*24*time.Hour).Unix() * 1000,
						LastLogInTimestamp: time.Now().Add(-1*time.Hour).Unix() * 1000,
					},
				},
			}
			return usersLoadedMsg{users: sampleUsers}
		}

		return usersLoadedMsg{users: users}
	}
}

func fetchRoutines(storeSvg *firebase.FirestoreService) tea.Cmd {
	return func() tea.Msg {
		if storeSvg == nil {
			return routinesErrorMsg{err: errors.New("firestore service not initialized")}
		}

		// Fetch users from Firebase Auth
		ctx := context.Background()

		// Try to fetch users from Firebase Auth
		routines, err := storeSvg.List(ctx, "routines")

		if err != nil {
			// If we fail to get users, return fallback data for development
			return routinesErrorMsg{err: fmt.Errorf("failed to fetch routines: %w", err)}
		}

		// Process user records from the iterator
		a := routines[0]
		log.Println(a)

		// If no users were found, add sample users for development
		return routinesLoadedMsg{routines: routines}
	}
}

func realMain() {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	// Setup signal handling for graceful shutdown
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	// Get initial terminal size
	width, height, _ := getTermSize()

	// Initialize model with loading state
	m := Model{
		title:       "Arrogance Admin",
		message:     "Initializing Firebase...",
		loading:     true,
		spinnerIdx:  0,
		width:       width,
		height:      height,
		activeTab:   HomeTab,
		tabs:        []string{"Home", "Users", "Routines"},
		currentView: LoadingView,
		userLoading: false,
	}

	// Initialize user table
	m.userTable = initUserTable()

	// Start the application
	p := tea.NewProgram(m, tea.WithAltScreen())

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
