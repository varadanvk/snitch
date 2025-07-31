package notifications

import (
	"fmt"
	"log"
	"math/rand"
	"os/exec"
	"time"

	"snitch-tui/src/core"

	"github.com/gen2brain/beeep"
)

// Professional notification messages
var NotificationMessages = map[string][]string{
	"distracted": {
		"Distraction detected - consider refocusing on your current task",
		"Non-productive activity identified",
		"Activity may not align with your current goals",
		"Consider returning to your primary task",
	},
	"productive": {
		"Productive activity detected - great work!",
		"Good focus on productive tasks",
		"Productive session in progress",
		"Maintaining good work habits",
	},
	"reminder": {
		"Task reminder - how is your progress?",
		"Checking in on your current task",
		"Time to review your current activity",
		"Task status check",
	},
}

// Manager handles sending notifications
type Manager struct {
	lastNotification time.Time
	minInterval      time.Duration
}

// NewManager creates a new notification manager
func NewManager(minInterval time.Duration) *Manager {
	return &Manager{
		minInterval: minInterval,
	}
}

// SendActivityNotification sends a notification based on activity
func (nm *Manager) SendActivityNotification(activity core.Activity) error {
	// Rate limiting
	if time.Since(nm.lastNotification) < nm.minInterval {
		return nil // Skip notification due to rate limiting
	}

	var messages []string
	if activity.IsProductive {
		messages = NotificationMessages["productive"]
	} else {
		messages = NotificationMessages["distracted"]
	}

	// Pick a random professional message
	message := messages[rand.Intn(len(messages))]

	// Add context
	contextMessage := fmt.Sprintf("%s\n\nDetected: %s in %s",
		message, activity.Activity, activity.Application)

	// Send notification
	err := nm.sendNotification("Snitch Productivity Monitor", contextMessage)
	if err != nil {
		return fmt.Errorf("failed to send notification: %w", err)
	}

	nm.lastNotification = time.Now()
	return nil
}

// SendCustomNotification sends a custom notification
func (nm *Manager) SendCustomNotification(title, message string) error {
	return nm.sendNotification(title, message)
}

// sendNotification sends a notification using multiple methods
func (nm *Manager) sendNotification(title, message string) error {
	// Try beeep first (cross-platform)
	err := beeep.Notify(title, message, "")
	if err != nil {
		log.Printf("beeep notification failed: %v", err)
	}

	// Also try macOS native notification with sound
	cmd := exec.Command("osascript", "-e",
		fmt.Sprintf(`display notification "%s" with title "%s" sound name "Ping"`,
			message, title))

	if err := cmd.Run(); err != nil {
		log.Printf("macOS notification failed: %v", err)
		return err
	}

	return nil
}

// SetMinInterval updates the minimum interval between notifications
func (nm *Manager) SetMinInterval(interval time.Duration) {
	nm.minInterval = interval
}
