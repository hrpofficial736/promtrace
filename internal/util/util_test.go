package util

import (
	"testing"
	"time"
)

func TestGenerateID(t *testing.T) {
	id := GenerateID()
	if len(id) != 8 {
		t.Errorf("expected 8 chars, got %d", len(id))
	}

	id2 := GenerateID()
	if id == id2 {
		t.Error("expected unique IDs")
	}
}

func TestParseDuration(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected time.Duration
		wantErr  bool
	}{
		{"7 days", "7d", 7 * 24 * time.Hour, false},
		{"1 day", "1d", 24 * time.Hour, false},
		{"24 hours", "24h", 24 * time.Hour, false},
		{"30 minutes", "30m", 30 * time.Minute, false},
		{"invalid", "abc", 0, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ParseDuration(tc.input)
			if tc.wantErr && err == nil {
				t.Error("expected error")
			}

			if !tc.wantErr && got != tc.expected {
				t.Errorf("got %v, expected %v", got, tc.expected)
			}
		})
	}
}
