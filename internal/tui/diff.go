package tui

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/hrpofficial736/promtrace/internal/diff"
	"github.com/hrpofficial736/promtrace/internal/util"
)

func RenderDiff(diff *diff.DiffResponseModel) string {

	titleText := RenderText(Heading, "TRACES DIFF\n")

	areSysPromptsIdentical := diff.SystemPrompt.Identical

	var sysPromptText string

	if areSysPromptsIdentical {
		sysPromptText = "identical"
	} else {
		sysPromptText = fmt.Sprintf(
			"%s\n%s\n",
			RemovedStyle.Render("\n- "+diff.SystemPrompt.Old),
			AddedStyle.Render("+ "+diff.SystemPrompt.New),
		)
	}

	areUserPromptsIdentical := diff.UserPrompt.Identical

	var userPromptText string

	if areUserPromptsIdentical {
		userPromptText = "identical"
	} else {
		userPromptText = fmt.Sprintf(
			"%s\n%s\n",
			RemovedStyle.Render("\n- "+diff.UserPrompt.Old),
			AddedStyle.Render("+ "+diff.UserPrompt.New),
		)
	}

	respSign := "+"

	if diff.ResponseLength.Delta < 0 {
		respSign = ""
	}

	responseLengthText := fmt.Sprintf(
		"response length: %d chars -> %d chars (%s%d)",
		diff.ResponseLength.A,
		diff.ResponseLength.B,
		respSign,
		diff.ResponseLength.Delta,
	)

	costSign := "+"

	if diff.Cost.Delta < 0 {
		costSign = ""
	}

	costText := fmt.Sprintf(
		"cost: %s -> %s (%s%d%%)",
		util.FmtCost(diff.Cost.A),
		util.FmtCost(diff.Cost.B),
		costSign,
		diff.Cost.Delta,
	)

	latencySign := "+"

	if diff.Latency.Delta < 0 {
		latencySign = ""
	}

	latencyText := fmt.Sprintf(
		"latency: %dms -> %dms (%s%d%%)",
		diff.Latency.A,
		diff.Latency.B,
		latencySign,
		diff.Latency.Delta,
	)

	container := lipgloss.JoinVertical(
		lipgloss.Left,
		titleText,
		"system prompt: "+sysPromptText,
		"user prompt: "+userPromptText,
		responseLengthText,
		costText,
		latencyText,
	)
	return BoxStyle.Border(lipgloss.ASCIIBorder()).Padding(0, 1).Render(container)
}
