package components

import (
	"snitch-tui/src/core"

	"github.com/charmbracelet/bubbles/table"
)

// ActivityComponent handles activity log display
type ActivityComponent struct {
	table            table.Model
	activities       []core.Activity
	selectedActivity *core.Activity
}

// NewActivityComponent creates a new activity component
func NewActivityComponent() *ActivityComponent {
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

	return &ActivityComponent{
		table:      t,
		activities: make([]core.Activity, 0),
	}
}

// GetTable returns the activity table
func (ac *ActivityComponent) GetTable() table.Model {
	return ac.table
}

// GetActivities returns the activities list
func (ac *ActivityComponent) GetActivities() []core.Activity {
	return ac.activities
}

// GetSelectedActivity returns the currently selected activity
func (ac *ActivityComponent) GetSelectedActivity() *core.Activity {
	return ac.selectedActivity
}

// SetSelectedActivity sets the currently selected activity
func (ac *ActivityComponent) SetSelectedActivity(activity *core.Activity) {
	ac.selectedActivity = activity
}

// SetActivities updates the activities list
func (ac *ActivityComponent) SetActivities(activities []core.Activity) {
	ac.activities = activities
}

// UpdateTable updates the table with new activity data
func (ac *ActivityComponent) UpdateTable(activities []core.Activity) {
	ac.activities = activities
	// Table update logic would go here - this is a placeholder
	// In a full implementation, this would populate the table rows
}

// UpdateData updates the component data (alias for UpdateTable for compatibility)
func (ac *ActivityComponent) UpdateData(activities []core.Activity) {
	ac.UpdateTable(activities)
}

// SetCurrentView sets the current view mode for the activity component
func (ac *ActivityComponent) SetCurrentView(view string) {
	// Placeholder for view state management
	// This would typically manage different activity views like "list", "detail", etc.
}
