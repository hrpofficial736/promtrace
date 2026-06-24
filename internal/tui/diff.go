package tui

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/lipgloss"
	"github.com/hrpofficial736/promtrace/internal/diff"
	"github.com/hrpofficial736/promtrace/internal/util"
)

func RenderDiff(diff *diff.DiffResponseModel) string {

	titleText := RenderText(Heading, util.GetPromtraceHeadingAsciiText()+"\nTraces Diff: ", 0, 0)

	sysPromptStatus := diff.SystemPrompt.Identical

	sysPromptText := RenderText(Heading, "system prompt: ", 0, 0)

	sysPromptVal := ""

	if sysPromptStatus {
		sysPromptVal = "identical"
	}

	sysPromptText += sysPromptVal

	if !sysPromptStatus {
		sysPromptText += fmt.Sprintf(`
		Old:
		%s
		New: 
		%s
		`, diff.SystemPrompt.Old, diff.SystemPrompt.New)
	}

	sysPromptText = RenderText(Hint, sysPromptText, 0, 0)

	userPromptStatus := diff.UserPrompt.Identical

	userPromptText := RenderText(Heading, "user prompt: ", 0, 0)

	userPromptVal := ""

	if userPromptStatus {
		userPromptVal = "identical"
	}

	userPromptText += userPromptVal

	if !userPromptStatus {
		userPromptText += fmt.Sprintf(`
		Old: 
		%s
		New:
		%s
		`, diff.UserPrompt.Old, diff.UserPrompt.New)
	}

	userPromptText = RenderText(Hint, userPromptText, 0, 0)

	responseSign := "+"
	costSign := "+"
	latencySign := "+"

	if diff.ResponseLength.Delta < 0 {
		responseSign = "-"
	}

	if diff.Cost.Delta < 0 {
		costSign = "-"
	}

	if diff.Latency.Delta < 0 {
		latencySign = "-"
	}

	numbersText := RenderText(Hint, fmt.Sprintf(`
response length: %s chars -> %s chars (%s%s)
cost: $%s -> $%s (%s%.2f%%)
latency: %sms -> %sms (%s%.2f%%)`,
		strconv.Itoa(diff.ResponseLength.A), strconv.Itoa(diff.ResponseLength.B), responseSign, strconv.Itoa(diff.ResponseLength.Delta),
		strconv.Itoa(diff.Cost.A), strconv.Itoa(diff.Cost.B), costSign, diff.Cost.PctChange,
		strconv.Itoa(diff.Latency.A), strconv.Itoa(diff.Latency.B), latencySign, diff.Latency.PctChange,
	), 0, 0)

	container := lipgloss.JoinVertical(lipgloss.Left, titleText, sysPromptText, userPromptText, numbersText)
	content := lipgloss.NewStyle().Padding(0, 1).Render(container)

	return content
}
