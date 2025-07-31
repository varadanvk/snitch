package ml

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"net/http"
	"snitch-tui/src/core"
	"strings"
	"time"
)

// GroqRequest follows OpenAI API format
type GroqRequest struct {
	Model          string          `json:"model"`
	Messages       []Message       `json:"messages"`
	Stream         bool            `json:"stream"`
	ResponseFormat *ResponseFormat `json:"response_format,omitempty"`
}

type ResponseFormat struct {
	Type string `json:"type"`
}

type Message struct {
	Role    string    `json:"role"`
	Content []Content `json:"content"`
}

type Content struct {
	Type     string    `json:"type"`
	Text     string    `json:"text,omitempty"`
	ImageURL *ImageURL `json:"image_url,omitempty"`
}

type ImageURL struct {
	URL string `json:"url"`
}

// GroqResponse follows OpenAI API format
type GroqResponse struct {
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Message struct {
		Content string `json:"content"`
	} `json:"message"`
}

type GroqAnalyzer struct {
	groqURL string
	apiKey  string
	client  *http.Client
}

type GroqAnalysisResult struct {
	Activity          string  `json:"activity"`
	IsProductive      bool    `json:"is_productive"`
	ProductivityScore float64 `json:"productivity_score"`
	Category          string  `json:"category"`
	Confidence        float64 `json:"confidence"`
	TaskAlignment     float64 `json:"task_alignment"`
}

// NewGroqAnalyzer creates a new Groq analyzer
func NewGroqAnalyzer(apiKey string) *GroqAnalyzer {
	return &GroqAnalyzer{
		groqURL: "https://api.groq.com/openai/v1/chat/completions",
		apiKey:  apiKey,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (ga *GroqAnalyzer) queryGroq(prompt, imageB64 string) (string, error) {
	// Create the request with vision model and JSON mode
	request := GroqRequest{
		Model: "meta-llama/llama-4-scout-17b-16e-instruct", // Correct Groq vision model
		Messages: []Message{
			{
				Role: "user",
				Content: []Content{
					{
						Type: "text",
						Text: prompt,
					},
					{
						Type: "image_url",
						ImageURL: &ImageURL{
							URL: "data:image/png;base64," + imageB64,
						},
					},
				},
			},
		},
		Stream: false,
		ResponseFormat: &ResponseFormat{
			Type: "json_object",
		},
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", ga.groqURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+ga.apiKey)

	// Make the request
	resp, err := ga.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to query Groq: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Read the error response for debugging
		var errorBody bytes.Buffer
		errorBody.ReadFrom(resp.Body)
		return "", fmt.Errorf("Groq API returned status %d: %s", resp.StatusCode, errorBody.String())
	}

	var groqResp GroqResponse
	err = json.NewDecoder(resp.Body).Decode(&groqResp)
	if err != nil {
		return "", fmt.Errorf("failed to decode Groq response: %w", err)
	}

	if len(groqResp.Choices) == 0 {
		return "", fmt.Errorf("no choices in Groq response")
	}

	return groqResp.Choices[0].Message.Content, nil
}

func (ga *GroqAnalyzer) imageToBase64(img image.Image) (string, error) {
	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func (ga *GroqAnalyzer) AnalyzeScreenshot(img image.Image, windowInfo core.WindowInfo, monitoringInterval int, currentTask string) (core.Activity, error) {
	imageB64, err := ga.imageToBase64(img)
	if err != nil {
		return core.Activity{}, fmt.Errorf("failed to encode image: %w", err)
	}

	taskContext := ""
	if currentTask != "" {
		taskContext = fmt.Sprintf("\nCurrent task: %s", currentTask)
	}

	prompt := fmt.Sprintf(`Analyze this screenshot and determine what activity the user is doing. 

Current application: %s
Window title: %s%s

Please respond with ONLY a JSON object containing:
{
  "activity": "brief description of what the user is doing",
  "is_productive": true/false,
  "productivity_score": 0.0-1.0,
  "category": "work/break/distraction",
  "confidence": 0.0-1.0,
  "task_alignment": 0.0-1.0
}

Focus on identifying:
- Code editing, development work, documentation
- Communication (email, messaging, meetings)
- Research, reading technical content
- Social media, entertainment, gaming
- Shopping, news browsing

%sConsider how well the current activity aligns with the stated task when setting task_alignment and productivity scores.

Be concise and accurate. Return ONLY the JSON, no other text.`,
		windowInfo.Application,
		windowInfo.Title,
		taskContext,
		func() string {
			if currentTask != "" {
				return fmt.Sprintf("IMPORTANT: The user should be working on: %s. ", currentTask)
			}
			return ""
		}())

	response, err := ga.queryGroq(prompt, imageB64)
	if err != nil {
		return core.Activity{}, fmt.Errorf("failed to query Groq: %w", err)
	}

	// Parse the JSON response
	var result GroqAnalysisResult

	// Try to extract JSON from the response
	jsonStart := strings.Index(response, "{")
	jsonEnd := strings.LastIndex(response, "}") + 1

	if jsonStart == -1 || jsonEnd <= jsonStart {
		return core.Activity{}, fmt.Errorf("no JSON found in response: %s", response)
	}

	jsonStr := response[jsonStart:jsonEnd]

	err = json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		return core.Activity{}, fmt.Errorf("failed to unmarshal Groq response: %w, response: %s", err, jsonStr)
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
