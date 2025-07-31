package core

import (
	"sync"
	"time"
)

// Activity represents a single tracked activity
type Activity struct {
	Timestamp         time.Time `json:"timestamp"`
	Type              string    `json:"type"`              // "productive" or "distracting"
	Activity          string    `json:"activity"`          // description
	Application       string    `json:"application"`       // detected application
	WindowTitle       string    `json:"window_title"`      // window title if available
	IsProductive      bool      `json:"is_productive"`
	Duration          int       `json:"duration"`
	ProductivityScore float64   `json:"productivity_score"` // 0-1 score
	Category          string    `json:"category"`          // "work", "break", "distraction"
}

// ActivityHistory manages the history of activities
type ActivityHistory struct {
	activities []Activity
	mu         sync.RWMutex
}

// NewActivityHistory creates a new activity history
func NewActivityHistory() *ActivityHistory {
	return &ActivityHistory{
		activities: make([]Activity, 0),
	}
}

// Add adds a new activity to the history
func (ah *ActivityHistory) Add(activity Activity) {
	ah.mu.Lock()
	defer ah.mu.Unlock()
	
	ah.activities = append(ah.activities, activity)
	
	// Keep only last 1000 activities
	if len(ah.activities) > 1000 {
		ah.activities = ah.activities[len(ah.activities)-1000:]
	}
}

// GetRecent returns the most recent activities
func (ah *ActivityHistory) GetRecent(count int) []Activity {
	ah.mu.RLock()
	defer ah.mu.RUnlock()
	
	if len(ah.activities) < count {
		return ah.activities
	}
	return ah.activities[len(ah.activities)-count:]
}

// GetAll returns all activities
func (ah *ActivityHistory) GetAll() []Activity {
	ah.mu.RLock()
	defer ah.mu.RUnlock()
	
	// Return a copy to prevent external modification
	result := make([]Activity, len(ah.activities))
	copy(result, ah.activities)
	return result
}

// Count returns the total number of activities
func (ah *ActivityHistory) Count() int {
	ah.mu.RLock()
	defer ah.mu.RUnlock()
	
	return len(ah.activities)
}

// ProductivityStats holds productivity statistics
type ProductivityStats struct {
	TotalTime        time.Duration
	ProductiveTime   time.Duration
	DistractingTime  time.Duration
	ProductivityRate float64
	TopActivities    map[string]int
	TopApps          map[string]int
}

// CalculateStats computes productivity statistics from activities
func (ah *ActivityHistory) CalculateStats() ProductivityStats {
	ah.mu.RLock()
	defer ah.mu.RUnlock()
	
	stats := ProductivityStats{
		TopActivities: make(map[string]int),
		TopApps:       make(map[string]int),
	}
	
	for _, activity := range ah.activities {
		duration := time.Duration(activity.Duration) * time.Second
		stats.TotalTime += duration
		
		if activity.IsProductive {
			stats.ProductiveTime += duration
		} else {
			stats.DistractingTime += duration
		}
		
		stats.TopActivities[activity.Activity]++
		stats.TopApps[activity.Application]++
	}
	
	if stats.TotalTime > 0 {
		stats.ProductivityRate = float64(stats.ProductiveTime) / float64(stats.TotalTime)
	}
	
	return stats
}