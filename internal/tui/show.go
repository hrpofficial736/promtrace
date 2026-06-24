package tui

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/lipgloss"
	"github.com/hrpofficial736/promtrace/internal/store"
	"github.com/hrpofficial736/promtrace/internal/util"
)

func RenderTraceInfoContainer(t *store.Trace) string {
	titleText := RenderText(Heading, util.GetPromtraceHeadingAsciiText()+"\nTrace Information: "+t.ID+"\n", 0, 1)

	basicInfo := RenderText(Hint, fmt.Sprintf(`
 Call: %s
 Model: %s
 Time: %s
 Latency: %sms
 Tokens: %s
 Cost: %s
	`,
		t.ID,
		t.Model,
		t.Timestamp.Format("Jan 02 15:04:02"),
		strconv.Itoa(int(t.LatencyMs)),
		strconv.Itoa(t.Tokens),
		strconv.Itoa(t.Cost),
	), 0, 0)

	systemPrompt := BoxStyle.Render(
		RenderText(Hint, "System Prompt", 0, 0) + "\n" + t.SystemPrompt,
	)

	userPrompt := BoxStyle.Render(
		RenderText(Hint, "User Prompt", 0, 0) + "\n" + t.UserPrompt,
	)

	responseText := BoxStyle.Render(
		RenderText(Hint, "Response", 0, 0) + "\n" + t.Response,
	)

	tb := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(ColorTertiary).
		Bold(true).
		Padding(0, 1)

	container := lipgloss.JoinVertical(lipgloss.Left, basicInfo, systemPrompt, userPrompt, responseText)

	window := lipgloss.JoinVertical(lipgloss.Left, titleText, tb.Render(container))

	return lipgloss.NewStyle().Render(window)
}
