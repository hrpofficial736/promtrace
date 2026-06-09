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
	Cost         int
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
	TotalCost   int
	AvgLatency  int
	TotalTokens int
}

type trendData struct {
	Date       time.Time
	Calls      int
	Tokens     int
	Cost       int
	AvgLatency int
}

type Stats struct {
	TotalCalls  int
	TotalTokens int
	TotalCost   int
	AvgLatency  int
	Trend       []*trendData
}
