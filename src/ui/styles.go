package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// Color constants
const (
	ColorPrimary      = "#2563EB"
	ColorSuccess      = "#10B981"
	ColorWarning      = "#F59E0B"
	ColorError        = "#EF4444"
	ColorSecondary    = "#6B7280"
	ColorMuted        = "#9CA3AF"
	ColorDark         = "#374151"
	ColorWhite        = "#FFFFFF"
	ColorBackground   = "#1D4ED8"
	ColorTableBorder  = "#6B7280"
)

// Header styles
func GetHeaderStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(ColorPrimary)).
		MarginBottom(1)
}

// Menu styles
func GetSelectedMenuStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorWhite)).
		Background(lipgloss.Color(ColorBackground)).
		Padding(0, 1)
}

func GetNormalMenuStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorSecondary)).
		Padding(0, 1)
}

func GetDisabledMenuStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorMuted)).
		Padding(0, 1)
}

// Status styles
func GetStatusStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorSuccess)).
		MarginTop(1)
}

func GetSuccessStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorSuccess)).
		MarginLeft(2)
}

func GetWarningStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorWarning)).
		MarginLeft(2)
}

func GetErrorStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorError)).
		MarginLeft(2)
}

// Table styles
func GetTableBorderStyle() lipgloss.Border {
	return lipgloss.NormalBorder()
}

func GetTableBorderColor() lipgloss.Color {
	return lipgloss.Color(ColorTableBorder)
}

func GetTableHeaderColor() lipgloss.Color {
	return lipgloss.Color(ColorDark)
}

func GetTableSelectedForeground() lipgloss.Color {
	return lipgloss.Color(ColorWhite)
}

func GetTableSelectedBackground() lipgloss.Color {
	return lipgloss.Color(ColorBackground)
}

// Spinner style
func GetSpinnerStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(ColorPrimary))
}

// Label and value styles
func GetLabelStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(ColorSecondary))
}

func GetValueStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorDark)).
		MarginLeft(2)
}

func GetBoldLabelStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(ColorDark))
}

// Activity detail styles
func GetActivityStatusStyle(isProductive bool) lipgloss.Style {
	color := ColorError
	if isProductive {
		color = ColorSuccess
	}
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(color))
}

func GetProductivityScoreStyle(score float64) lipgloss.Style {
	var color string
	if score >= 0.7 {
		color = ColorSuccess
	} else if score >= 0.4 {
		color = ColorWarning
	} else {
		color = ColorError
	}
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(color))
}

// Setup styles
func GetSetupOptionStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorBackground)).
		Bold(true).
		MarginLeft(2)
}

// Settings styles
func GetSettingsSelectedStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorWhite)).
		Background(lipgloss.Color(ColorBackground)).
		Padding(0, 1)
}

func GetSettingsNormalStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorSecondary)).
		Padding(0, 1)
}

// Statistics styles
func GetStatStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorSecondary)).
		MarginLeft(2)
}

// Task styles
func GetTaskSuccessStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorSuccess)).
		MarginLeft(2)
}

func GetTaskWarningStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorWarning)).
		MarginLeft(2)
}

// Utility functions for consistent color application
func ApplyProductiveColor() lipgloss.Color {
	return lipgloss.Color(ColorSuccess)
}

func ApplyDistractingColor() lipgloss.Color {
	return lipgloss.Color(ColorError)
}

func ApplyPrimaryColor() lipgloss.Color {
	return lipgloss.Color(ColorPrimary)
}

func ApplySecondaryColor() lipgloss.Color {
	return lipgloss.Color(ColorSecondary)
}

func ApplyWarningColor() lipgloss.Color {
	return lipgloss.Color(ColorWarning)
}

// Additional style functions for views
func GetSelectedStyle() lipgloss.Style {
	return GetSelectedMenuStyle()
}

func GetNormalStyle() lipgloss.Style {
	return GetNormalMenuStyle()
}

func GetDisabledStyle() lipgloss.Style {
	return GetDisabledMenuStyle()
}

func GetProductivityStatusStyle(isProductive bool) lipgloss.Style {
	return GetActivityStatusStyle(isProductive)
}

func GetOptionStyle() lipgloss.Style {
	return GetSetupOptionStyle()
}