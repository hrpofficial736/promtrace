package diff

import (
	"github.com/hrpofficial736/promtrace/internal/store"
)

func diffStrings(s1, s2 string) {
	if s1 == s2 {
		// identical
	} else {
		// - removed
		// + added
	}
}

func diffNumbers(a, b int) {
	if a == b {
		// identical
	} else {
		pct := float64(b-a) / float64(a) * 100
		// sign := "+"
		if pct < 0 {
			// sign = ""
		}

		// print numbers with % change
	}
}

func DiffTraces(t1, t2 *store.Trace) error {
	// system prompt diff
	diffStrings(t1.SystemPrompt, t2.SystemPrompt)

	// user prompt diff
	diffStrings(t1.UserPrompt, t2.UserPrompt)

	// response length
	diffNumbers(len(t1.Response), len(t2.Response))

	// cost
	diffNumbers(t1.Cost, t2.Cost)

	// latency
	diffNumbers(int(t1.LatencyMs), int(t2.LatencyMs))

	return nil
}
