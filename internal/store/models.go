package store

import (
	"time"
)

type Trace struct {
	ID           string
	SessionID    string
	Timestamp    time.Time
	Host         string
	Method       string
	Path         string
	Model        string
	Tokens       int
	Cost         float64
	SystemPrompt string
	UserPrompt   string
	RequestBody  string
	Response     string
	StatusCode   int
	LatencyMs    int64
	CreatedAt    time.Time
}

type Session struct {
	ID          string
	TotalCalls  int
	StartedAt   time.Time
	TotalCost   float64
	AvgLatency  int
	TotalTokens int
}

type trendData struct {
	Date       time.Time
	Calls      int
	Tokens     int
	Cost       float64
	AvgLatency float64
}

type Stats struct {
	TotalCalls  int
	TotalTokens int
	TotalCost   float64
	AvgLatency  float64
	Trend       []*trendData
}
