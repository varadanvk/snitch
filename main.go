package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/varadanvk/snitch/src/config"
	"github.com/varadanvk/snitch/src/core"
	"github.com/varadanvk/snitch/src/ml"
	"github.com/varadanvk/snitch/src/notifications"
	"github.com/varadanvk/snitch/src/ui"
)

const Version = "v1.0.0"

// SnitchCore manages the core monitoring functionality
type SnitchCore struct {
	configManager       *config.Manager
	screenMonitor       *core.ScreenMonitor
	activityHistory     *core.ActivityHistory
	analyzer            *ml.Analyzer
	notificationManager *notifications.Manager

	currentTask    string
	isMonitoring   bool
	monitoringStop chan bool
	sessionStart   time.Time
	mu             sync.RWMutex
}

// NewSnitchCore creates a new Snitch core instance
func NewSnitchCore() *SnitchCore {
	configManager := config.NewManager()
	cfg := configManager.Get()

	// Determine AI backend
	var backend ml.AIBackendType
	switch cfg.AIBackend {
	case "groq":
		backend = ml.BackendGroq
	case "ollama":
		backend = ml.BackendOllama
	default:
		backend = ml.BackendGroq // Default to Groq
	}

	return &SnitchCore{
		configManager:       configManager,
		screenMonitor:       core.NewScreenMonitor(),
		activityHistory:     core.NewActivityHistory(),
		analyzer:            ml.NewAnalyzer(backend, cfg.OllamaURL, cfg.OllamaModel, cfg.GroqAPIKey),
		notificationManager: notifications.NewManager(time.Duration(cfg.NotificationInterval) * time.Second),
		monitoringStop:      make(chan bool),
		sessionStart:        time.Now(),
	}
}

// StartMonitoring begins the monitoring process
func (sc *SnitchCore) StartMonitoring() {
	sc.mu.Lock()
	if sc.isMonitoring {
		sc.mu.Unlock()
		return
	}
	sc.isMonitoring = true
	sc.sessionStart = time.Now()
	sc.mu.Unlock()

	go sc.monitoringLoop()
}

// StopMonitoring stops the monitoring process
func (sc *SnitchCore) StopMonitoring() {
	sc.mu.Lock()
	if !sc.isMonitoring {
		sc.mu.Unlock()
		return
	}
	sc.isMonitoring = false
	sc.mu.Unlock()

	sc.monitoringStop <- true
}

// monitoringLoop is the main monitoring loop
func (sc *SnitchCore) monitoringLoop() {
	cfg := sc.configManager.Get()
	ticker := time.NewTicker(time.Duration(cfg.MonitoringInterval) * time.Second)
	defer ticker.Stop()

	log.Println("Monitoring loop started")

	for {
		select {
		case <-sc.monitoringStop:
			log.Println("Monitoring loop stopped")
			return
		case <-ticker.C:
			// Capture screen
			img, err := sc.screenMonitor.CaptureScreen()
			if err != nil {
				log.Printf("Error capturing screen: %v", err)
				continue
			}

			// Get window info
			windowInfo, err := sc.screenMonitor.GetActiveWindow()
			if err != nil {
				log.Printf("Error getting window info: %v", err)
				continue
			}

			// Analyze activity
			activity, err := sc.analyzer.AnalyzeActivity(img, windowInfo, cfg.MonitoringInterval, sc.GetCurrentTask())
			if err != nil {
				log.Printf("Error analyzing activity: %v", err)
				continue
			}

			// Add to history
			sc.activityHistory.Add(activity)

			// Send notification if appropriate
			if !activity.IsProductive {
				sc.notificationManager.SendActivityNotification(activity)
			}
		}
	}
}

// IsMonitoring returns whether monitoring is active
func (sc *SnitchCore) IsMonitoring() bool {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	return sc.isMonitoring
}

// GetStats returns productivity statistics
func (sc *SnitchCore) GetStats() core.ProductivityStats {
	return sc.activityHistory.CalculateStats()
}

// GetRecentActivities returns recent activities
func (sc *SnitchCore) GetRecentActivities(count int) []core.Activity {
	return sc.activityHistory.GetRecent(count)
}

// SetCurrentTask sets the current task
func (sc *SnitchCore) SetCurrentTask(task string) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.currentTask = task
	log.Printf("Current task set to: %s", task)
}

// GetCurrentTask returns the current task
func (sc *SnitchCore) GetCurrentTask() string {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	return sc.currentTask
}

// GetConfig returns the current configuration
func (sc *SnitchCore) GetConfig() *config.Config {
	return sc.configManager.Get()
}

// SaveConfig saves the current configuration to disk
func (sc *SnitchCore) SaveConfig() error {
	return sc.configManager.Save()
}

func main() {
	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	// Create the core application instance
	core := NewSnitchCore()

	// Create the UI model with the core
	model := ui.NewModel(core)

	// Create and run the TUI program
	program := ui.NewProgram(model)
	if err := ui.RunProgram(program); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
