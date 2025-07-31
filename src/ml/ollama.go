package ml

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"net/http"
	"strings"
	"time"

	"snitch-tui/src/core"
)

// OllamaRequest represents a request to Ollama API
type OllamaRequest struct {
	Model  string   `json:"model"`
	Prompt string   `json:"prompt"`
	Images []string `json:"images,omitempty"`
	Stream bool     `json:"stream"`
}

// OllamaResponse represents a response from Ollama API
type OllamaResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

// OllamaAnalyzer uses Ollama API for real image analysis
type OllamaAnalyzer struct {
	ollamaURL string
	model     string
	client    *http.Client
}

// NewOllamaAnalyzer creates a new Ollama-based analyzer
func NewOllamaAnalyzer(ollamaURL, model string) *OllamaAnalyzer {
	if ollamaURL == "" {
		ollamaURL = "http://localhost:11434"
	}
	if model == "" {
		model = "llava"
	}

	return &OllamaAnalyzer{
		ollamaURL: ollamaURL,
		model:     model,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// AnalyzeScreenshot analyzes a screenshot using Ollama's vision model
func (oa *OllamaAnalyzer) AnalyzeScreenshot(img image.Image, windowInfo core.WindowInfo, monitoringInterval int) (core.Activity, error) {
	// Convert image to base64
	imageB64, err := oa.imageToBase64(img)
	if err != nil {
		return core.Activity{}, fmt.Errorf("failed to encode image: %w", err)
	}

	// Create prompt for productivity analysis
	prompt := fmt.Sprintf(`Analyze this screenshot and determine what activity the user is doing. 

Current application: %s
Window title: %s

Please respond with a JSON object containing:
{
  "activity": "brief description of what the user is doing",
  "is_productive": true/false,
  "productivity_score": 0.0-1.0,
  "category": "work/break/distraction",
  "confidence": 0.0-1.0
}

Focus on identifying:
- Code editing, development work, documentation
- Communication (email, messaging, meetings)
- Research, reading technical content
- Social media, entertainment, gaming
- Shopping, news browsing

Be concise and accurate.`, windowInfo.Application, windowInfo.Title)

	// Make request to Ollama
	response, err := oa.queryOllama(prompt, imageB64)
	if err != nil {
		// Fallback to simple heuristic if Ollama fails
		return oa.fallbackAnalysis(windowInfo, monitoringInterval), nil
	}

	// Parse response
	activity, err := oa.parseOllamaResponse(response, windowInfo, monitoringInterval)
	if err != nil {
		// Fallback if parsing fails
		return oa.fallbackAnalysis(windowInfo, monitoringInterval), nil
	}

	return activity, nil
}

// imageToBase64 converts an image to base64 string
func (oa *OllamaAnalyzer) imageToBase64(img image.Image) (string, error) {
	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

// queryOllama sends a request to Ollama API
func (oa *OllamaAnalyzer) queryOllama(prompt, imageB64 string) (string, error) {
	request := OllamaRequest{
		Model:  oa.model,
		Prompt: prompt,
		Images: []string{imageB64},
		Stream: false,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", err
	}

	resp, err := oa.client.Post(
		oa.ollamaURL+"/api/generate",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return "", fmt.Errorf("failed to query Ollama: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Ollama API returned status %d", resp.StatusCode)
	}

	var ollamaResp OllamaResponse
	err = json.NewDecoder(resp.Body).Decode(&ollamaResp)
	if err != nil {
		return "", err
	}

	return ollamaResp.Response, nil
}

// OllamaAnalysisResult represents the parsed result from Ollama
type OllamaAnalysisResult struct {
	Activity          string  `json:"activity"`
	IsProductive      bool    `json:"is_productive"`
	ProductivityScore float64 `json:"productivity_score"`
	Category          string  `json:"category"`
	Confidence        float64 `json:"confidence"`
}

// parseOllamaResponse parses Ollama's response into an Activity
func (oa *OllamaAnalyzer) parseOllamaResponse(response string, windowInfo core.WindowInfo, monitoringInterval int) (core.Activity, error) {
	// Try to extract JSON from the response
	jsonStart := strings.Index(response, "{")
	jsonEnd := strings.LastIndex(response, "}") + 1

	if jsonStart == -1 || jsonEnd <= jsonStart {
		return core.Activity{}, fmt.Errorf("no JSON found in response")
	}

	jsonStr := response[jsonStart:jsonEnd]

	var result OllamaAnalysisResult
	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		return core.Activity{}, fmt.Errorf("failed to parse JSON: %w", err)
	}

	activityType := "distracting"
	if result.IsProductive {
		activityType = "productive"
	}

	return core.Activity{
		Timestamp:         time.Now(),
		Type:              activityType,
		Activity:          result.Activity,
		Application:       windowInfo.Application,
		WindowTitle:       windowInfo.Title,
		IsProductive:      result.IsProductive,
		Duration:          monitoringInterval,
		ProductivityScore: result.ProductivityScore,
		Category:          result.Category,
	}, nil
}

// fallbackAnalysis provides simple heuristic-based analysis when Ollama fails
func (oa *OllamaAnalyzer) fallbackAnalysis(windowInfo core.WindowInfo, monitoringInterval int) core.Activity {
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

// IsOllamaAvailable checks if Ollama is running and the model is available
func (oa *OllamaAnalyzer) IsOllamaAvailable() bool {
	resp, err := oa.client.Get(oa.ollamaURL + "/api/tags")
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}
