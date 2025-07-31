package ui

import (
	"fmt"
	"strings"
	"time"

)

// viewMain renders the main menu view
func (m *Model) viewMain() string {
	// Styles
	headerStyle := GetHeaderStyle()
	selectedStyle := GetSelectedStyle()
	normalStyle := GetNormalStyle()
	disabledStyle := GetDisabledStyle()
	statusStyle := GetStatusStyle()

	// Build the view
	s := headerStyle.Render("SNITCH AI PRODUCTIVITY MONITOR") + "\n\n"

	// Menu
	for i, choice := range m.choices {
		cursor := " "
		style := normalStyle

		// Check if this menu item is disabled
		disabled := m.isMenuItemDisabled(i)

		if m.cursor == i && !disabled {
			cursor = ">"
			style = selectedStyle
		} else if disabled {
			style = disabledStyle
		}

		s += style.Render(fmt.Sprintf("%s %s", cursor, choice)) + "\n"
	}

	// Status
	s += statusStyle.Render(fmt.Sprintf("\n%s", m.status))

	// Monitoring indicator with spinner
	if m.core.IsMonitoring() {
		sessionDuration := time.Since(m.sessionStart).Round(time.Second)
		s += statusStyle.Render(fmt.Sprintf("\n%s Monitoring active - Session: %s",
			m.spinner.View(), sessionDuration))

		// Show current task if set
		currentTask := m.core.GetCurrentTask()
		if currentTask != "" {
			s += statusStyle.Render(fmt.Sprintf("\n[TASK] %s", currentTask))
		}

		// Show quick stats
		if stats := m.core.GetStats(); stats.TotalTime > 0 {
			productivityPercent := stats.ProductivityRate * 100
			s += statusStyle.Render(fmt.Sprintf("\n[STATS] Productivity: %.1f%% | Activities: %d",
				productivityPercent, len(m.activities)))
		}
	}

	s += statusStyle.Render("\n\nPress q to quit, ↑/↓ to navigate, enter to select")

	return s
}

// viewActivity renders the activity log view
func (m *Model) viewActivity() string {
	headerStyle := GetHeaderStyle()

	s := headerStyle.Render("RECENT ACTIVITY LOG") + "\n\n"
	s += m.activityTable.View() + "\n"
	s += "\nPress 'enter' to view details, 'b' to go back, 'q' to quit, ↑/↓ to navigate"
	return s
}

// viewActivityDetail renders the activity detail view
func (m *Model) viewActivityDetail() string {
	if m.selectedActivity == nil {
		return "No activity selected"
	}

	headerStyle := GetHeaderStyle()
	labelStyle := GetLabelStyle()
	valueStyle := GetValueStyle()

	activity := m.selectedActivity

	s := headerStyle.Render("ACTIVITY DETAILS") + "\n\n"

	s += labelStyle.Render("Timestamp:") + "\n"
	s += valueStyle.Render(activity.Timestamp.Format("Monday, January 2, 2006 at 3:04:05 PM")) + "\n\n"

	s += labelStyle.Render("Application:") + "\n"
	s += valueStyle.Render(activity.Application) + "\n\n"

	s += labelStyle.Render("Window Title:") + "\n"
	s += valueStyle.Render(activity.WindowTitle) + "\n\n"

	s += labelStyle.Render("Activity Description:") + "\n"
	s += valueStyle.Render(activity.Activity) + "\n\n"

	s += labelStyle.Render("Status:") + "\n"
	statusText := "[DISTRACTING]"
	if activity.IsProductive {
		statusText = "[PRODUCTIVE]"
	}
	statusStyle := GetProductivityStatusStyle(activity.IsProductive)
	s += statusStyle.Render(statusText) + "\n\n"

	s += labelStyle.Render("Productivity Score:") + "\n"
	scoreText := fmt.Sprintf("%.1f%% ", activity.ProductivityScore*100)

	// Add visual indicator
	if activity.ProductivityScore >= 0.8 {
		scoreText += "[EXCELLENT]"
	} else if activity.ProductivityScore >= 0.6 {
		scoreText += "[GOOD]"
	} else if activity.ProductivityScore >= 0.4 {
		scoreText += "[FAIR]"
	} else {
		scoreText += "[POOR]"
	}

	scoreStyle := GetProductivityScoreStyle(activity.ProductivityScore)
	s += scoreStyle.Render(scoreText) + "\n\n"

	s += labelStyle.Render("Category:") + "\n"
	s += valueStyle.Render(strings.Title(activity.Category)) + "\n\n"

	s += labelStyle.Render("Duration:") + "\n"
	s += valueStyle.Render(fmt.Sprintf("%d seconds", activity.Duration)) + "\n\n"

	// Add some analysis context
	s += labelStyle.Render("Analysis Context:") + "\n"
	if activity.IsProductive {
		s += valueStyle.Render("[PRODUCTIVE] This activity aligns with productive work patterns") + "\n"
	} else {
		s += valueStyle.Render("[WARNING] This activity may be distracting from your main goals") + "\n"
	}

	s += "\n\nPress 'b' to go back, 'q' to quit"
	return s
}

// viewStats renders the productivity statistics view
func (m *Model) viewStats() string {
	headerStyle := GetHeaderStyle()
	statStyle := GetStatStyle()

	s := headerStyle.Render("PRODUCTIVITY STATISTICS") + "\n\n"

	stats := m.core.GetStats()

	// Session info
	sessionDuration := time.Since(m.sessionStart).Round(time.Second)
	s += statStyle.Render(fmt.Sprintf("Session Duration: %s", sessionDuration)) + "\n"
	s += statStyle.Render(fmt.Sprintf("Total Tracked: %s", stats.TotalTime.Round(time.Second))) + "\n"
	s += statStyle.Render(fmt.Sprintf("Productive Time: %s", stats.ProductiveTime.Round(time.Second))) + "\n"
	s += statStyle.Render(fmt.Sprintf("Distracting Time: %s", stats.DistractingTime.Round(time.Second))) + "\n"

	// Productivity rate with progress bar
	s += "\n" + statStyle.Render("Productivity Rate:") + "\n"
	s += m.progress.ViewAs(stats.ProductivityRate) + fmt.Sprintf(" %.1f%%", stats.ProductivityRate*100) + "\n"

	// Top activities
	if len(stats.TopActivities) > 0 {
		s += "\n" + statStyle.Render("Top Activities:") + "\n"
		count := 0
		for activity, freq := range stats.TopActivities {
			if count >= 5 {
				break
			}
			s += statStyle.Render(fmt.Sprintf("  • %s (%d times)", activity, freq)) + "\n"
			count++
		}
	}

	// Top apps
	if len(stats.TopApps) > 0 {
		s += "\n" + statStyle.Render("Top Applications:") + "\n"
		count := 0
		for app, freq := range stats.TopApps {
			if count >= 5 {
				break
			}
			s += statStyle.Render(fmt.Sprintf("  • %s (%d times)", app, freq)) + "\n"
			count++
		}
	}

	s += "\n\nPress 'b' to go back, 'q' to quit"
	return s
}

// viewTask renders the task input view
func (m *Model) viewTask() string {
	headerStyle := GetHeaderStyle()
	labelStyle := GetLabelStyle() 
	valueStyle := GetValueStyle()
	successStyle := GetSuccessStyle()
	warningStyle := GetWarningStyle()

	s := headerStyle.Render("SET CURRENT TASK") + "\n\n"

	s += labelStyle.Render("What are you working on right now?") + "\n\n"

	s += m.taskInput.View() + "\n\n"

	if m.taskMessage != "" {
		if strings.Contains(m.taskMessage, "[SUCCESS]") {
			s += successStyle.Render(m.taskMessage) + "\n\n"
		} else {
			s += warningStyle.Render(m.taskMessage) + "\n\n"
		}
	}

	currentTask := m.core.GetCurrentTask()
	if currentTask != "" {
		s += labelStyle.Render("Current Task:") + "\n"
		s += valueStyle.Render("Task: " + currentTask) + "\n\n"
	}

	s += valueStyle.Render("Examples:") + "\n"
	s += valueStyle.Render("• Working on user authentication feature") + "\n"
	s += valueStyle.Render("• Debugging payment processing issue") + "\n"
	s += valueStyle.Render("• Writing documentation for API endpoints") + "\n"
	s += valueStyle.Render("• Reviewing pull requests") + "\n"
	s += valueStyle.Render("• Planning sprint for next week") + "\n\n"

	s += "Press Enter to save, 'b' to go back, 'q' to quit"
	return s
}

// viewSetup renders the AI setup view
func (m *Model) viewSetup() string {
	headerStyle := GetHeaderStyle()
	labelStyle := GetLabelStyle()
	valueStyle := GetValueStyle()
	warningStyle := GetWarningStyle()
	successStyle := GetSuccessStyle()
	optionStyle := GetOptionStyle()

	cfg := m.core.GetConfig()

	switch m.setupStep {
	case 0: // Choose backend
		s := headerStyle.Render("AI SETUP - CHOOSE YOUR BACKEND") + "\n\n"

		s += labelStyle.Render("Current Configuration:") + "\n"
		backendStatus := cfg.AIBackend
		if cfg.AIBackend == "groq" {
			if cfg.GroqAPIKey != "" {
				backendStatus += " [CONFIGURED]"
			} else {
				backendStatus += " [MISSING KEY]"
			}
		} else if cfg.AIBackend == "ollama" {
			backendStatus += " (local)"
		}
		s += valueStyle.Render("Backend: " + strings.Title(backendStatus)) + "\n\n"

		s += labelStyle.Render("Choose an AI backend:") + "\n\n"

		s += optionStyle.Render("1. Groq (Recommended)") + "\n"
		s += valueStyle.Render("   • Fast cloud-based AI") + "\n"
		s += valueStyle.Render("   • Requires API key (free tier available)") + "\n"
		s += valueStyle.Render("   • Best performance") + "\n\n"

		s += optionStyle.Render("2. Ollama (Local)") + "\n"
		s += valueStyle.Render("   • Runs locally on your machine") + "\n"
		s += valueStyle.Render("   • Complete privacy") + "\n"
		s += valueStyle.Render("   • Requires Ollama installation") + "\n\n"

		s += "Press 1 for Groq, 2 for Ollama, or 'b' to go back"
		return s

	case 1: // Enter API key
		s := headerStyle.Render("GROQ API KEY SETUP") + "\n\n"

		s += labelStyle.Render("Enter your Groq API key:") + "\n\n"

		s += m.textInput.View() + "\n\n"

		if m.setupMessage != "" {
			if strings.Contains(m.setupMessage, "[ERROR]") {
				s += warningStyle.Render(m.setupMessage) + "\n\n"
			} else {
				s += successStyle.Render(m.setupMessage) + "\n\n"
			}
		}

		s += valueStyle.Render("Get your free API key at: https://console.groq.com/") + "\n"
		s += valueStyle.Render("Press Enter to save, Esc to go back") + "\n"

		return s

	case 2: // Confirmation
		s := headerStyle.Render("SETUP COMPLETE!") + "\n\n"

		if m.setupMessage != "" {
			s += successStyle.Render(m.setupMessage) + "\n\n"
		}

		s += labelStyle.Render("Current Configuration:") + "\n"
		s += valueStyle.Render("Backend: " + strings.Title(cfg.AIBackend)) + "\n"
		if cfg.AIBackend == "groq" && cfg.GroqAPIKey != "" {
			s += valueStyle.Render("API Key: " + cfg.GroqAPIKey[:8] + "..." + cfg.GroqAPIKey[len(cfg.GroqAPIKey)-4:]) + "\n"
		}
		s += valueStyle.Render("Status: Ready for AI-powered analysis!") + "\n\n"

		s += "Press Enter or 'b' to return to main menu"
		return s

	default:
		return "Unknown setup step"
	}
}

// viewSettings renders the settings view
func (m *Model) viewSettings() string {
	headerStyle := GetHeaderStyle()
	normalStyle := GetNormalStyle()
	selectedStyle := GetSelectedStyle()
	labelStyle := GetLabelStyle()
	valueStyle := GetValueStyle()
	successStyle := GetSuccessStyle()
	errorStyle := GetErrorStyle()

	s := headerStyle.Render("SETTINGS") + "\n\n"

	cfg := m.core.GetConfig()

	if m.settingsEditing {
		// Show editing interface
		s += labelStyle.Render(fmt.Sprintf("Editing: %s", m.settingsItems[m.settingsCursor])) + "\n\n"
		s += m.settingsInput.View() + "\n\n"

		if m.settingsMessage != "" {
			if strings.Contains(m.settingsMessage, "[SUCCESS]") {
				s += successStyle.Render(m.settingsMessage) + "\n\n"
			} else {
				s += errorStyle.Render(m.settingsMessage) + "\n\n"
			}
		}

		s += valueStyle.Render("Press Enter to save, Esc to cancel") + "\n"
	} else {
		// Show settings menu
		for i, setting := range m.settingsItems {
			cursor := " "
			style := normalStyle

			if m.settingsCursor == i {
				cursor = ">"
				style = selectedStyle
			}

			var currentValue string
			switch i {
			case 0:
				currentValue = fmt.Sprintf("%d seconds", cfg.MonitoringInterval)
			case 1:
				currentValue = fmt.Sprintf("%d seconds", cfg.NotificationInterval)
			case 2:
				currentValue = fmt.Sprintf("%t", cfg.SaveScreenshots)
			case 3:
				currentValue = fmt.Sprintf("%t", cfg.SnitchMode)
			case 4:
				currentValue = fmt.Sprintf("%d apps", len(cfg.ProductiveApps))
			case 5:
				currentValue = fmt.Sprintf("%d apps", len(cfg.DistractingApps))
			}

			s += style.Render(fmt.Sprintf("%s %s: %s", cursor, setting, currentValue)) + "\n"
		}

		if m.settingsMessage != "" {
			if strings.Contains(m.settingsMessage, "[SUCCESS]") {
				s += "\n" + successStyle.Render(m.settingsMessage) + "\n"
			} else {
				s += "\n" + errorStyle.Render(m.settingsMessage) + "\n"
			}
		}

		// Show detailed view of selected setting
		s += "\n" + labelStyle.Render("Current Configuration:") + "\n"
		switch m.settingsCursor {
		case 4: // Productive Apps
			s += valueStyle.Render("Productive Apps:") + "\n"
			for _, app := range cfg.ProductiveApps {
				s += valueStyle.Render(fmt.Sprintf("  + %s", app)) + "\n"
			}
		case 5: // Distracting Apps
			s += valueStyle.Render("Distracting Apps:") + "\n"
			for _, app := range cfg.DistractingApps {
				s += valueStyle.Render(fmt.Sprintf("  - %s", app)) + "\n"
			}
		}

		s += "\n\nPress Enter to edit, ↑/↓ to navigate, 'b' to go back, 'q' to quit"
	}

	return s
}