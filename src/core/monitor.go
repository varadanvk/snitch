package core

import (
	"fmt"
	"image"
	"os/exec"
	"strings"

	"github.com/kbinani/screenshot"
)

// ScreenMonitor handles screen capture and window detection
type ScreenMonitor struct {
	lastCapture image.Image
}

// NewScreenMonitor creates a new screen monitor
func NewScreenMonitor() *ScreenMonitor {
	return &ScreenMonitor{}
}

// CaptureScreen captures the current screen
func (sm *ScreenMonitor) CaptureScreen() (image.Image, error) {
	bounds := screenshot.GetDisplayBounds(0)
	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		return nil, fmt.Errorf("failed to capture screen: %w", err)
	}
	
	sm.lastCapture = img
	return img, nil
}

// GetLastCapture returns the last captured screen
func (sm *ScreenMonitor) GetLastCapture() image.Image {
	return sm.lastCapture
}

// WindowInfo holds information about the active window
type WindowInfo struct {
	Application string
	Title       string
}

// GetActiveWindow returns information about the currently active window (macOS)
func (sm *ScreenMonitor) GetActiveWindow() (WindowInfo, error) {
	// Get active application
	cmd := exec.Command("osascript", "-e", 
		`tell application "System Events" to get name of first application process whose frontmost is true`)
	appOutput, err := cmd.Output()
	if err != nil {
		return WindowInfo{}, fmt.Errorf("failed to get active app: %w", err)
	}
	app := strings.TrimSpace(string(appOutput))

	// Get window title
	cmd = exec.Command("osascript", "-e", 
		fmt.Sprintf(`tell application "System Events" to get title of front window of application process "%s"`, app))
	titleOutput, err := cmd.Output()
	title := "Unknown"
	if err == nil {
		title = strings.TrimSpace(string(titleOutput))
	}

	return WindowInfo{
		Application: app,
		Title:       title,
	}, nil
}