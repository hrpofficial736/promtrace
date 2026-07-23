package diff

import (
	"github.com/hrpofficial736/promtrace/internal/store"
)

type StringDiffModel struct {
	Identical bool
	Old       string
	New       string
}

type NumericDiffModel[T any] struct {
	A         T
	B         T
	Delta     T
	PctChange float64
}

type DiffResponseModel struct {
	SystemPrompt   *StringDiffModel
	UserPrompt     *StringDiffModel
	ResponseLength *NumericDiffModel[int]
	Cost           *NumericDiffModel[float64]
	Latency        *NumericDiffModel[int]
}

func diffStrings(s1, s2 string) *StringDiffModel {
	response := &StringDiffModel{}
	if s1 == s2 {
		response.Identical = true
	} else {
		response.Old = s1
		response.New = s2
	}

	return response
}

func diffNumbers[T int | float64](a, b T) *NumericDiffModel[T] {

	response := &NumericDiffModel[T]{A: a, B: b, Delta: b - a}

	if a != 0 {
		response.PctChange = (float64(b-a) / float64(a) * 100)
	}

	return response

}

func DiffTraces(t1, t2 *store.Trace) (*DiffResponseModel, error) {
	diffResponse := &DiffResponseModel{}
	// system prompt diff
	diffResponse.SystemPrompt = diffStrings(t1.SystemPrompt, t2.SystemPrompt)

	// user prompt diff
	diffResponse.UserPrompt = diffStrings(t1.UserPrompt, t2.UserPrompt)

	// response length
	diffResponse.ResponseLength = diffNumbers(len(t1.Response), len(t2.Response))

	// cost
	diffResponse.Cost = diffNumbers(t1.Cost, t2.Cost)

	// latency
	diffResponse.Latency = diffNumbers(int(t1.LatencyMs), int(t2.LatencyMs))

	return diffResponse, nil
}
