package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config holds all application configuration
type Config struct {
	MonitoringInterval    int                   `json:"monitoring_interval"`
	NotificationInterval  int                   `json:"notification_interval"`
	Sensitivity           string                `json:"sensitivity"`
	FocusedHours          map[string]int        `json:"focused_hours"`
	Theme                 string                `json:"theme"`
	SaveScreenshots       bool                  `json:"save_screenshots"`
	ProductiveApps        []string              `json:"productive_apps"`
	DistractingApps       []string              `json:"distracting_apps"`
	SnitchMode            bool                  `json:"snitch_mode"`
	AccountabilityBuddies []AccountabilityBuddy `json:"accountability_buddies"`
	
	// AI Configuration
	AIBackend   string `json:"ai_backend"`    // "ollama" or "groq"
	OllamaURL   string `json:"ollama_url"`
	OllamaModel string `json:"ollama_model"`
	GroqAPIKey  string `json:"groq_api_key"`
}

type AccountabilityBuddy struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Enabled     bool   `json:"enabled"`
}

// Manager handles configuration loading/saving
type Manager struct {
	configPath string
	config     *Config
}

// NewManager creates a new configuration manager
func NewManager() *Manager {
	homeDir, _ := os.UserHomeDir()
	configDir := filepath.Join(homeDir, ".snitch")
	configPath := filepath.Join(configDir, "config.json")

	// Create config directory
	os.MkdirAll(configDir, 0755)

	manager := &Manager{
		configPath: configPath,
		config:     defaultConfig(),
	}

	// Load existing config
	manager.Load()

	return manager
}

// defaultConfig returns default configuration
func defaultConfig() *Config {
	return &Config{
		MonitoringInterval:    3,
		NotificationInterval:  15,
		Sensitivity:           "medium",
		FocusedHours:          map[string]int{"start": 9, "end": 17},
		Theme:                 "system",
		SaveScreenshots:       false,
		ProductiveApps:        []string{"Code", "Terminal", "Xcode", "IntelliJ", "Sublime Text", "Vim"},
		DistractingApps:       []string{"Safari", "Chrome", "YouTube", "Twitter", "Instagram", "TikTok"},
		SnitchMode:            false,
		AccountabilityBuddies: []AccountabilityBuddy{},
		
		// AI defaults
		AIBackend:   "groq",
		OllamaURL:   "http://localhost:11434",
		OllamaModel: "llava",
		GroqAPIKey:  "", // User needs to set this
	}
}

// Get returns the current configuration
func (m *Manager) Get() *Config {
	return m.config
}

// Update modifies configuration values
func (m *Manager) Update(updates map[string]interface{}) error {
	// This is a simplified update - in practice you'd want type-safe updates
	data, err := json.Marshal(updates)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, m.config)
}

// Load reads configuration from disk
func (m *Manager) Load() error {
	data, err := os.ReadFile(m.configPath)
	if err != nil {
		// File doesn't exist, use defaults
		return nil
	}

	return json.Unmarshal(data, m.config)
}

// Save writes configuration to disk
func (m *Manager) Save() error {
	data, err := json.MarshalIndent(m.config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(m.configPath, data, 0644)
}
