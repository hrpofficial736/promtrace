package store

import (
	"testing"
)

func TestSaveAndGetTrace(t *testing.T) {
	st, err := NewSQLiteStore(":memory:")
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	defer func() {
		err := st.Close()
		if err != nil {
			t.Fatalf("failed to close store: %v", err)
		}
	}()

	trace := &Trace{
		ID:        "test-001",
		SessionID: "session-abc",
		Host:      "api.openai.com",
		Method:    "POST",
		Model:     "gpt-4o",
		Tokens:    100,
		Cost:      24,
	}

	err = st.SaveTrace(trace)
	if err != nil {
		t.Fatalf("SaveTrace failed: %v", err)
	}

	got, err := st.GetTrace("test-001")
	if err != nil {
		t.Fatalf("GetTrace failed: %v", err)
	}

	if got.ID != "test-001" {
		t.Errorf("ID: got %s, want test-001", got.ID)
	}

	if got.Model != "gpt-4o" {
		t.Errorf("Model: got %s, want gpt-4o", got.Model)
	}

	if got.Tokens != 100 {
		t.Errorf("Tokens: got %d, want 100", got.Tokens)
	}

}
