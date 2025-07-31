package ml

import (
	"fmt"
	"image"
	"net/http"
	"snitch-tui/src/core"
	"time"
)

// AIBackendType defines which AI backend to use
type AIBackendType int

const (
	BackendOllama AIBackendType = iota
	BackendGroq
)

// AIAnalyzer is a unified interface for AI-powered activity analysis
type AIAnalyzer struct {
	ollamaAnalyzer *OllamaAnalyzer
	groqAnalyzer   *GroqAnalyzer
	backend        AIBackendType
}

// NewAIAnalyzer creates a new AIAnalyzer with the specified backend
func NewAIAnalyzer(backend AIBackendType, ollamaURL, ollamaModel, groqAPIKey string) *AIAnalyzer {
	var ollamaAnalyzer *OllamaAnalyzer
	var groqAnalyzer *GroqAnalyzer

	if backend == BackendOllama {
		ollamaAnalyzer = NewOllamaAnalyzer(ollamaURL, ollamaModel)
	} else if backend == BackendGroq {
		groqAnalyzer = NewGroqAnalyzer(groqAPIKey)
	}

	return &AIAnalyzer{
		ollamaAnalyzer: ollamaAnalyzer,
		groqAnalyzer:   groqAnalyzer,
		backend:        backend,
	}
}

// AnalyzeActivity analyzes a screenshot and window info to determine activity using the selected backend
func (a *AIAnalyzer) AnalyzeActivity(img image.Image, windowInfo core.WindowInfo, monitoringInterval int, currentTask string) (core.Activity, error) {
	switch a.backend {
	case BackendOllama:
		if a.ollamaAnalyzer == nil {
			return core.Activity{}, fmt.Errorf("Ollama analyzer not initialized")
		}
		return a.ollamaAnalyzer.AnalyzeScreenshot(img, windowInfo, monitoringInterval)
	case BackendGroq:
		if a.groqAnalyzer == nil {
			return core.Activity{}, fmt.Errorf("Groq analyzer not initialized")
		}
		return a.groqAnalyzer.AnalyzeScreenshot(img, windowInfo, monitoringInterval, currentTask)
	default:
		return core.Activity{}, fmt.Errorf("unknown AI backend")
	}
}

// SetBackend switches the backend at runtime
func (a *AIAnalyzer) SetBackend(backend AIBackendType) {
	a.backend = backend
}

// Helper to provide a default HTTP client (30s timeout)
func defaultHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 30 * time.Second,
	}
}
