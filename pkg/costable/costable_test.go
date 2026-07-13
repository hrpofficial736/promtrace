package costable

import "testing"

func TestCalculateCost(t *testing.T) {
	cost := CalculateCost("gpt-4o", 1000, 500)

	if cost == 0 {
		t.Errorf("expected non zero cost for gpt-4o")
	}
}

func TestCalculateCostUnknownModel(t *testing.T) {
	cost := CalculateCost("nonexistent-model", 1000, 500)

	if cost != 0 {
		t.Errorf("expected 0 for unknown model, got %f", cost)
	}
}
