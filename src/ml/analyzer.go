package ml

import (
	"image"
	"log"
	"strings"
	"time"

	"snitch-tui/src/core"
)

// Analyzer handles activity analysis using the AIAnalyzer
type Analyzer struct {
	aiAnalyzer *AIAnalyzer
	useAI      bool
}

// NewAnalyzer creates a new activity analyzer with backend selection
func NewAnalyzer(backend AIBackendType, ollamaURL, ollamaModel, groqAPIKey string) *Analyzer {
	aiAnalyzer := NewAIAnalyzer(backend, ollamaURL, ollamaModel, groqAPIKey)

	// Check if AI is available based on backend
	useAI := false
	switch backend {
	case BackendOllama:
		if aiAnalyzer.ollamaAnalyzer != nil {
			useAI = aiAnalyzer.ollamaAnalyzer.IsOllamaAvailable()
		}
		if useAI {
			log.Println("Ollama is available - using AI-powered analysis")
		} else {
			log.Println("Ollama not available - using fallback heuristic analysis")
		}
	case BackendGroq:
		useAI = groqAPIKey != ""
		if useAI {
			log.Println("Groq API key configured - using AI-powered analysis")
		} else {
			log.Println("Groq API key not configured - using fallback heuristic analysis")
		}
	}

	return &Analyzer{
		aiAnalyzer: aiAnalyzer,
		useAI:      useAI,
	}
}

// AnalyzeActivity analyzes a screenshot and window info to determine activity
func (a *Analyzer) AnalyzeActivity(img image.Image, windowInfo core.WindowInfo, monitoringInterval int, currentTask string) (core.Activity, error) {
	if a.useAI {
		activity, err := a.aiAnalyzer.AnalyzeActivity(img, windowInfo, monitoringInterval, currentTask)
		if err != nil {
			log.Printf("AI analysis failed, using fallback: %v", err)
			return a.fallbackAnalysis(windowInfo, monitoringInterval), nil
		}
		return activity, nil
	}

	return a.fallbackAnalysis(windowInfo, monitoringInterval), nil
}

// SwitchBackend changes the AI backend at runtime
func (a *Analyzer) SwitchBackend(backend AIBackendType) {
	if a.aiAnalyzer != nil {
		a.aiAnalyzer.SetBackend(backend)
		log.Printf("Switched to %s backend", getBackendName(backend))
	}
}

// RefreshAIStatus checks if AI is available and updates the flag
func (a *Analyzer) RefreshAIStatus() {
	if a.aiAnalyzer == nil {
		a.useAI = false
		return
	}

	switch a.aiAnalyzer.backend {
	case BackendOllama:
		if a.aiAnalyzer.ollamaAnalyzer != nil {
			a.useAI = a.aiAnalyzer.ollamaAnalyzer.IsOllamaAvailable()
		}
	case BackendGroq:
		a.useAI = a.aiAnalyzer.groqAnalyzer != nil
	}

	if a.useAI {
		log.Printf("AI is now available using %s backend", getBackendName(a.aiAnalyzer.backend))
	} else {
		log.Println("AI is not available - using fallback analysis")
	}
}

// getBackendName returns a human-readable name for the backend
func getBackendName(backend AIBackendType) string {
	switch backend {
	case BackendOllama:
		return "Ollama"
	case BackendGroq:
		return "Groq"
	default:
		return "Unknown"
	}
}

// fallbackAnalysis provides simple heuristic-based analysis when AI fails
func (a *Analyzer) fallbackAnalysis(windowInfo core.WindowInfo, monitoringInterval int) core.Activity {
	appLower := strings.ToLower(windowInfo.Application)
	titleLower := strings.ToLower(windowInfo.Title)

	// Simple heuristics based on app names and window titles
	isProductive := false
	activity := "unknown activity"
	category := "unknown"
	score := 0.5

	// Productive indicators
	if strings.Contains(appLower, "code") || strings.Contains(appLower, "xcode") ||
		strings.Contains(appLower, "terminal") || strings.Contains(appLower, "vim") ||
		strings.Contains(titleLower, "code") || strings.Contains(titleLower, "programming") {
		isProductive = true
		activity = "coding/development"
		category = "work"
		score = 0.8
	} else if strings.Contains(appLower, "mail") || strings.Contains(titleLower, "email") {
		isProductive = true
		activity = "email communication"
		category = "work"
		score = 0.7
	} else if strings.Contains(appLower, "slack") || strings.Contains(appLower, "teams") ||
		strings.Contains(appLower, "zoom") || strings.Contains(appLower, "meet") {
		isProductive = true
		activity = "team communication"
		category = "work"
		score = 0.75
	} else if strings.Contains(appLower, "safari") || strings.Contains(appLower, "chrome") ||
		strings.Contains(appLower, "firefox") {
		// Browser - depends on content
		if strings.Contains(titleLower, "github") || strings.Contains(titleLower, "stackoverflow") ||
			strings.Contains(titleLower, "documentation") || strings.Contains(titleLower, "docs") {
			isProductive = true
			activity = "research/documentation"
			category = "work"
			score = 0.7
		} else if strings.Contains(titleLower, "youtube") || strings.Contains(titleLower, "netflix") ||
			strings.Contains(titleLower, "twitter") || strings.Contains(titleLower, "facebook") {
			isProductive = false
			activity = "entertainment/social media"
			category = "distraction"
			score = 0.2
		} else {
			activity = "web browsing"
			category = "break"
			score = 0.4
		}
	}

	activityType := "distracting"
	if isProductive {
		activityType = "productive"
	}

	return core.Activity{
		Timestamp:         time.Now(),
		Type:              activityType,
		Activity:          activity,
		Application:       windowInfo.Application,
		WindowTitle:       windowInfo.Title,
		IsProductive:      isProductive,
		Duration:          monitoringInterval,
		ProductivityScore: score,
		Category:          category,
	}
}
