package ui

import (
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"

	"github.com/varadanvk/snitch/src/config"
	"github.com/varadanvk/snitch/src/core"
)

// SnitchCore interface defines the core functionality needed by the UI
type SnitchCore interface {
	IsMonitoring() bool
	StartMonitoring()
	StopMonitoring()
	GetStats() core.ProductivityStats
	GetRecentActivities(count int) []core.Activity
	SetCurrentTask(task string)
	GetCurrentTask() string
	GetConfig() *config.Config
}

// Model represents the main TUI model with all state and components
type Model struct {
	core        SnitchCore
	choices     []string
	cursor      int
	status      string
	lastUpdate  time.Time
	currentView string // "main", "activity", "settings", "stats", "activity_detail", "setup", "task"

	// Bubble components
	activityTable table.Model
	spinner       spinner.Model
	progress      progress.Model
	textInput     textinput.Model

	// Activity detail view - now handled by ActivityComponent

	// Setup state
	setupStep    int // 0: choose backend, 1: enter groq key, 2: confirm
	setupMessage string

	// Task state
	taskInput   textinput.Model
	taskMessage string

	// Settings state
	settingsCursor  int
	settingsEditing bool
	settingsInput   textinput.Model
	settingsMessage string
	settingsItems   []string

	// Rendering state
	needsRedraw    bool
	lastRenderTime time.Time
	renderCount    int

	// Session tracking for monitoring
	sessionStart time.Time

	// UI Components (to be integrated)
	// These will hold component instances once fully integrated
	mainMenuComponent interface{}
	activityComponent interface{}
	settingsComponent interface{}
	statsComponent    interface{}
	setupComponent    interface{}
	taskComponent     interface{}

	// Activity state (temporary - will move to component)
	selectedActivity *core.Activity
	activities       []core.Activity
}

// tickMsg represents a tick message for updating the UI
type tickMsg time.Time

// NewModel creates a new Model instance with initialized components
func NewModel(core SnitchCore) *Model {
	// Initialize table for activity log
	columns := []table.Column{
		{Title: "ID", Width: 0}, // Hidden column for click handling
		{Title: "Time", Width: 10},
		{Title: "Status", Width: 12},
		{Title: "Activity", Width: 30},
		{Title: "Application", Width: 18},
		{Title: "Score", Width: 8},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	// Apply table styles (will be moved to styles.go)
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(GetTableBorderStyle()).
		BorderForeground(GetTableBorderColor()).
		BorderBottom(true).
		Bold(true).
		Foreground(GetTableHeaderColor())
	s.Selected = s.Selected.
		Foreground(GetTableSelectedForeground()).
		Background(GetTableSelectedBackground()).
		Bold(false)
	t.SetStyles(s)

	// Initialize spinner
	sp := spinner.New()
	sp.Spinner = spinner.Line
	sp.Style = GetSpinnerStyle()

	// Initialize progress bar
	prog := progress.New(progress.WithDefaultGradient())

	// Initialize text input for setup
	ti := textinput.New()
	ti.Placeholder = "Enter your Groq API key..."
	ti.CharLimit = 100
	ti.Width = 50

	// Initialize text input for task
	taskInput := textinput.New()
	taskInput.Placeholder = "Enter your current task (e.g., 'Working on user authentication feature')..."
	taskInput.CharLimit = 200
	taskInput.Width = 60

	// Initialize text input for settings
	settingsInput := textinput.New()
	settingsInput.Placeholder = "Enter new value..."
	settingsInput.CharLimit = 100
	settingsInput.Width = 40

	m := &Model{
		core: core,
		choices: []string{
			"Start Monitoring",
			"Stop Monitoring",
			"View Activity Log",
			"Productivity Stats",
			"Settings",
			"AI Setup",
			"Set Current Task",
			"Quit",
		},
		status:        "Snitch AI Productivity Monitor - Ready",
		currentView:   "main",
		activityTable: t,
		spinner:       sp,
		progress:      prog,
		textInput:     ti,
		taskInput:     taskInput,
		setupStep:     0,
		cursor:        0,
		settingsInput: settingsInput,
		settingsItems: []string{
			"Monitoring Interval",
			"Notification Interval",
			"Save Screenshots",
			"Snitch Mode",
			"Productive Apps",
			"Distracting Apps",
		},
		needsRedraw:    true,
		lastRenderTime: time.Now(),
		sessionStart:   time.Now(),
	}

	// Ensure cursor starts on a valid position
	m.cursor = m.findValidCursor(0, 1)

	return m
}

// Helper function to check if a menu item is disabled
func (m *Model) isMenuItemDisabled(index int) bool {
	if index == 0 && m.core.IsMonitoring() { // Start Monitoring
		return true
	}
	if index == 1 && !m.core.IsMonitoring() { // Stop Monitoring
		return true
	}
	return false
}

// Helper function to find a valid cursor position
func (m *Model) findValidCursor(start, direction int) int {
	cursor := start
	for i := 0; i < len(m.choices); i++ {
		if !m.isMenuItemDisabled(cursor) {
			return cursor
		}
		cursor += direction
		if cursor >= len(m.choices) {
			cursor = 0
		} else if cursor < 0 {
			cursor = len(m.choices) - 1
		}
	}
	return start // fallback
}

// GetCurrentView returns the current view name
func (m *Model) GetCurrentView() string {
	return m.currentView
}

// SetCurrentView sets the current view
func (m *Model) SetCurrentView(view string) {
	m.currentView = view
	m.needsRedraw = true
}

// GetStatus returns the current status message
func (m *Model) GetStatus() string {
	return m.status
}

// SetStatus sets the status message
func (m *Model) SetStatus(status string) {
	m.status = status
}

// GetSessionStart returns the session start time
func (m *Model) GetSessionStart() time.Time {
	return m.sessionStart
}

// SetSessionStart sets the session start time
func (m *Model) SetSessionStart(t time.Time) {
	m.sessionStart = t
}
