package tokenizer

import "testing"

func TestEstimateTokens(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"empty string", "", 0},
		{"short text", "hello", 1},
		{"16 chars", "hello world 1234", 4},
		{"single char", "a", 0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := EstimateTokens(tc.input)
			if got != tc.expected {
				t.Errorf("got %d, want %d", got, tc.expected)
			}
		})
	}
}
