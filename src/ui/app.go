package ui

import (
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// TickCmd creates a tick command for regular UI updates
func TickCmd() tea.Cmd {
	return tea.Tick(time.Second*2, func(t time.Time) tea.Msg { // Reduced frequency
		return tickMsg(t)
	})
}

// Init initializes the TUI model and returns initial commands
func (m *Model) Init() tea.Cmd {
	// Only start ticker if monitoring is already active
	var cmds []tea.Cmd
	if m.core.IsMonitoring() {
		cmds = append(cmds, TickCmd())
	}
	cmds = append(cmds, m.spinner.Tick, textinput.Blink)
	return tea.Batch(cmds...)
}

// Update handles all TUI updates and message routing
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	// Handle text input first if we're in setup or task mode
	if m.currentView == "setup" && m.setupStep == 1 {
		m.textInput, cmd = m.textInput.Update(msg)
		cmds = append(cmds, cmd)
	}
	if m.currentView == "task" {
		m.taskInput, cmd = m.taskInput.Update(msg)
		cmds = append(cmds, cmd)
	}
	if m.currentView == "settings" && m.settingsEditing {
		m.settingsInput, cmd = m.settingsInput.Update(msg)
		cmds = append(cmds, cmd)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.currentView {
		case "main":
			return m.updateMain(msg)
		case "activity":
			return m.updateActivity(msg)
		case "activity_detail":
			return m.updateActivityDetail(msg)
		case "settings":
			return m.updateSettings(msg)
		case "setup":
			return m.updateSetup(msg)
		case "task":
			return m.updateTask(msg)
		case "stats":
			return m.updateStats(msg)
		}
	case tea.MouseMsg:
		if m.currentView == "activity" && msg.Type == tea.MouseLeft {
			return m.handleTableClick(msg)
		}
	case tickMsg:
		m.lastUpdate = time.Time(msg)
		
		// Only update if monitoring is active or we're in activity view
		if m.core.IsMonitoring() || m.currentView == "activity" {
			if m.currentView == "activity" {
				m.updateActivityTable()
			}
			m.needsRedraw = true
		}
		
		// Only continue ticking if monitoring is active
		if m.core.IsMonitoring() {
			cmds = append(cmds, TickCmd())
		}
	case tea.Msg: // Handle spinner updates
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}
	
	return m, tea.Batch(cmds...)
}

// View renders the current view based on the model state
func (m *Model) View() string {
	switch m.currentView {
	case "main":
		return m.viewMain()
	case "activity":
		return m.viewActivity()
	case "activity_detail":
		return m.viewActivityDetail()
	case "settings":
		return m.viewSettings()
	case "setup":
		return m.viewSetup()
	case "task":
		return m.viewTask()
	case "stats":
		return m.viewStats()
	}
	return m.viewMain()
}

// NewProgram creates a new tea.Program with the given model
func NewProgram(model *Model) *tea.Program {
	return tea.NewProgram(
		model,
		tea.WithAltScreen(),       // Use alternate screen buffer
		tea.WithMouseCellMotion(), // Enable mouse support
	)
}

// RunProgram runs the TUI program and returns any error
func RunProgram(program *tea.Program) error {
	_, err := program.Run()
	return err
}