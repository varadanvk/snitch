package ui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

)

// updateMain handles main menu navigation and selection
func (m *Model) updateMain(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	oldCursor := m.cursor
	switch msg.String() {
	case "ctrl+c", "q":
		if m.core.IsMonitoring() {
			m.core.StopMonitoring()
		}
		return m, tea.Quit
	case "up", "k":
		for {
			if m.cursor > 0 {
				m.cursor--
			} else {
				break
			}
			// Check if current position is valid
			if !m.isMenuItemDisabled(m.cursor) {
				break
			}
		}
	case "down", "j":
		for {
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			} else {
				break
			}
			// Check if current position is valid
			if !m.isMenuItemDisabled(m.cursor) {
				break
			}
		}
	case "enter", " ":
		return m.handleMainSelection()
	}

	// Mark for redraw if cursor moved
	if oldCursor != m.cursor {
		m.needsRedraw = true
	}

	return m, nil
}

// handleMainSelection handles main menu item selection
func (m *Model) handleMainSelection() (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch m.cursor {
	case 0: // Start Monitoring
		if !m.core.IsMonitoring() {
			m.core.StartMonitoring()
			m.status = "[ACTIVE] Monitoring started - AI analysis active"
			// Move cursor to stop monitoring
			m.cursor = 1
			m.needsRedraw = true
			// Start the tick timer
			cmds = append(cmds, TickCmd())
		}
	case 1: // Stop Monitoring
		if m.core.IsMonitoring() {
			m.core.StopMonitoring()
			m.status = "[STOPPED] Monitoring stopped"
			// Move cursor to start monitoring
			m.cursor = 0
			m.needsRedraw = true
		}
	case 2: // View Activity Log
		m.currentView = "activity"
		m.status = "Activity Log - Press 'b' to go back"
		m.needsRedraw = true
	case 3: // Productivity Stats
		m.currentView = "stats"
		m.status = "Productivity Statistics - Press 'b' to go back"
		m.needsRedraw = true
	case 4: // Settings
		m.currentView = "settings"
		m.status = "Settings - Press 'b' to go back"
		m.needsRedraw = true
	case 5: // AI Setup
		m.currentView = "setup"
		m.status = "AI Setup - Press 'b' to go back"
		m.needsRedraw = true
	case 6: // Set Current Task
		m.currentView = "task"
		m.taskInput.Focus()
		m.taskInput.SetValue(m.core.GetCurrentTask()) // Pre-fill with current task
		m.taskMessage = ""
		m.status = "Set Current Task"
		m.needsRedraw = true
		cmds = append(cmds, textinput.Blink)
	case 7: // Quit
		if m.core.IsMonitoring() {
			m.core.StopMonitoring()
		}
		return m, tea.Quit
	}
	return m, tea.Batch(cmds...)
}

// updateActivity handles activity view navigation
func (m *Model) updateActivity(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg.String() {
	case "esc", "b":
		m.currentView = "main"
		m.status = "Snitch AI Productivity Monitor - Ready"
		m.needsRedraw = true
	case "ctrl+c", "q":
		return m, tea.Quit
	case "enter":
		// Show detail view for selected row
		selectedRow := m.activityTable.Cursor()
		if selectedRow < len(m.activities) {
			m.selectedActivity = &m.activities[selectedRow]
			m.currentView = "activity_detail"
			m.status = "Activity Details - Press 'b' to go back"
			m.needsRedraw = true
		}
	default:
		m.activityTable, cmd = m.activityTable.Update(msg)
	}
	return m, cmd
}

// updateActivityDetail handles activity detail view navigation
func (m *Model) updateActivityDetail(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "b":
		m.currentView = "activity"
		m.status = "Activity Log - Press 'b' to go back"
		m.selectedActivity = nil
	case "ctrl+c", "q":
		return m, tea.Quit
	}
	return m, nil
}

// handleTableClick handles mouse clicks on activity table
func (m *Model) handleTableClick(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	// Simple click handling - if clicked in table area, show details
	if msg.Y >= 3 && msg.Y <= 13 { // Approximate table area
		selectedRow := m.activityTable.Cursor()
		if selectedRow < len(m.activities) {
			m.selectedActivity = &m.activities[selectedRow]
			m.currentView = "activity_detail"
			m.status = "Activity Details - Press 'b' to go back"
		}
	}
	return m, nil
}

// updateActivityTable updates the activity table with recent data
func (m *Model) updateActivityTable() {
	recent := m.core.GetRecentActivities(20)
	m.activities = recent // Store for detail view
	rows := []table.Row{}

	for i, activity := range recent {
		status := "[DISTRACTED]"
		if activity.IsProductive {
			status = "[PRODUCTIVE]"
		}

		score := fmt.Sprintf("%.1f", activity.ProductivityScore*100)

		// Truncate long activity descriptions
		activityDesc := activity.Activity
		if len(activityDesc) > 25 {
			activityDesc = activityDesc[:22] + "..."
		}

		rows = append(rows, table.Row{
			fmt.Sprintf("%d", i), // Hidden ID for click handling
			activity.Timestamp.Format("15:04:05"),
			status,
			activityDesc,
			activity.Application,
			score + "%",
		})
	}

	m.activityTable.SetRows(rows)
}

// updateSettings handles settings view navigation and editing
func (m *Model) updateSettings(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if m.settingsEditing {
		switch msg.String() {
		case "esc":
			m.settingsEditing = false
			m.settingsInput.Blur()
			m.settingsMessage = ""
		case "enter":
			return m.handleSettingsSave()
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	} else {
		switch msg.String() {
		case "esc", "b":
			m.currentView = "main"
			m.status = "Snitch AI Productivity Monitor - Ready"
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.settingsCursor > 0 {
				m.settingsCursor--
			}
		case "down", "j":
			if m.settingsCursor < len(m.settingsItems)-1 {
				m.settingsCursor++
			}
		case "enter", " ":
			return m.handleSettingsEdit()
		}
	}
	return m, cmd
}

// handleSettingsEdit starts editing a setting
func (m *Model) handleSettingsEdit() (tea.Model, tea.Cmd) {
	cfg := m.core.GetConfig()
	var currentValue string

	switch m.settingsCursor {
	case 0: // Monitoring Interval
		currentValue = fmt.Sprintf("%d", cfg.MonitoringInterval)
	case 1: // Notification Interval
		currentValue = fmt.Sprintf("%d", cfg.NotificationInterval)
	case 2: // Save Screenshots
		currentValue = fmt.Sprintf("%t", cfg.SaveScreenshots)
	case 3: // Snitch Mode
		currentValue = fmt.Sprintf("%t", cfg.SnitchMode)
	case 4: // Productive Apps
		currentValue = strings.Join(cfg.ProductiveApps, ", ")
	case 5: // Distracting Apps
		currentValue = strings.Join(cfg.DistractingApps, ", ")
	}

	m.settingsEditing = true
	m.settingsInput.SetValue(currentValue)
	m.settingsInput.Focus()
	m.settingsMessage = ""

	return m, textinput.Blink
}

// handleSettingsSave saves the edited setting
func (m *Model) handleSettingsSave() (tea.Model, tea.Cmd) {
	cfg := m.core.GetConfig()
	newValue := strings.TrimSpace(m.settingsInput.Value())

	if newValue == "" {
		m.settingsMessage = "[ERROR] Value cannot be empty"
		return m, nil
	}

	switch m.settingsCursor {
	case 0: // Monitoring Interval
		if val, err := strconv.Atoi(newValue); err == nil && val > 0 {
			cfg.MonitoringInterval = val
			m.settingsMessage = "[SUCCESS] Monitoring interval updated"
		} else {
			m.settingsMessage = "[ERROR] Invalid number (must be > 0)"
			return m, nil
		}
	case 1: // Notification Interval
		if val, err := strconv.Atoi(newValue); err == nil && val > 0 {
			cfg.NotificationInterval = val
			m.settingsMessage = "[SUCCESS] Notification interval updated"
		} else {
			m.settingsMessage = "[ERROR] Invalid number (must be > 0)"
			return m, nil
		}
	case 2: // Save Screenshots
		if val, err := strconv.ParseBool(newValue); err == nil {
			cfg.SaveScreenshots = val
			m.settingsMessage = "[SUCCESS] Save screenshots updated"
		} else {
			m.settingsMessage = "[ERROR] Invalid boolean (true/false)"
			return m, nil
		}
	case 3: // Snitch Mode
		if val, err := strconv.ParseBool(newValue); err == nil {
			cfg.SnitchMode = val
			m.settingsMessage = "[SUCCESS] Snitch mode updated"
		} else {
			m.settingsMessage = "[ERROR] Invalid boolean (true/false)"
			return m, nil
		}
	case 4: // Productive Apps
		apps := []string{}
		for _, app := range strings.Split(newValue, ",") {
			app = strings.TrimSpace(app)
			if app != "" {
				apps = append(apps, app)
			}
		}
		cfg.ProductiveApps = apps
		m.settingsMessage = "[SUCCESS] Productive apps updated"
	case 5: // Distracting Apps
		apps := []string{}
		for _, app := range strings.Split(newValue, ",") {
			app = strings.TrimSpace(app)
			if app != "" {
				apps = append(apps, app)
			}
		}
		cfg.DistractingApps = apps
		m.settingsMessage = "[SUCCESS] Distracting apps updated"
	}

	// Save configuration - This needs to be implemented by the concrete core type
	// For now, we'll assume the core handles config saving internally

	m.settingsEditing = false
	m.settingsInput.Blur()

	return m, nil
}

// updateSetup handles AI setup navigation
func (m *Model) updateSetup(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch m.setupStep {
	case 0: // Choose backend
		switch msg.String() {
		case "esc", "b":
			m.currentView = "main"
			m.status = "Snitch AI Productivity Monitor - Ready"
			m.setupStep = 0
		case "ctrl+c", "q":
			return m, tea.Quit
		case "1":
			// Choose Groq
			m.setupStep = 1
			m.textInput.Focus()
			m.textInput.SetValue("")
			m.setupMessage = ""
			m.status = "Enter your Groq API key"
			cmd = textinput.Blink
		case "2":
			// Choose Ollama
			cfg := m.core.GetConfig()
			cfg.AIBackend = "ollama"
			// Save configuration - handled by core
			m.setupMessage = "[SUCCESS] Ollama backend selected! Make sure Ollama is running with 'ollama pull llava'"
			m.setupStep = 2
		}
	case 1: // Enter Groq API key
		switch msg.String() {
		case "esc":
			m.setupStep = 0
			m.textInput.Blur()
			m.setupMessage = ""
			m.status = "AI Setup - Choose your backend"
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			apiKey := strings.TrimSpace(m.textInput.Value())
			if len(apiKey) > 10 { // Basic validation
				// Save the API key
				cfg := m.core.GetConfig()
				cfg.GroqAPIKey = apiKey
				cfg.AIBackend = "groq"
				// Save configuration - handled by core

				// Recreate analyzer - would be handled by core implementation

				m.setupMessage = "[SUCCESS] Groq API key saved! AI analysis is now enabled."
				m.setupStep = 2
				m.textInput.Blur()
			} else {
				m.setupMessage = "[ERROR] Please enter a valid API key (should be longer than 10 characters)"
			}
		}
	case 2: // Confirmation
		switch msg.String() {
		case "esc", "b", "enter":
			m.currentView = "main"
			m.status = "Snitch AI Productivity Monitor - Ready"
			m.setupStep = 0
			m.setupMessage = ""
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, cmd
}

// updateTask handles task input
func (m *Model) updateTask(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg.String() {
	case "esc", "b":
		m.currentView = "main"
		m.status = "Snitch AI Productivity Monitor - Ready"
		m.taskInput.Blur()
		m.taskMessage = ""
	case "ctrl+c", "q":
		return m, tea.Quit
	case "enter":
		task := strings.TrimSpace(m.taskInput.Value())
		if len(task) > 0 {
			// Set the current task
			m.core.SetCurrentTask(task)
			m.taskMessage = "[SUCCESS] Current task set successfully!"

			// Auto-return to main after 2 seconds or on next key press
			m.currentView = "main"
			m.status = fmt.Sprintf("[TASK] Current task: %s", task)
			m.taskInput.Blur()
		} else {
			m.taskMessage = "[ERROR] Please enter a task description"
		}
	case "tab":
		// Focus the input if not already focused
		if !m.taskInput.Focused() {
			m.taskInput.Focus()
			cmd = textinput.Blink
		}
	}

	return m, cmd
}

// updateStats handles stats view navigation
func (m *Model) updateStats(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "b":
		m.currentView = "main"
		m.status = "Snitch AI Productivity Monitor - Ready"
	case "ctrl+c", "q":
		return m, tea.Quit
	}
	return m, nil
}