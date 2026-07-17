package tui

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/hrpofficial736/promtrace/internal/store"
	"github.com/hrpofficial736/promtrace/internal/util"
	"strconv"
)

func RenderTraceInfoContainer(t *store.Trace) string {

	heading := RenderText(Heading, "TRACE INFORMATION")

	basicInfo := RenderText(Hint, fmt.Sprintf(`
ID: %s
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
		util.FmtCost(t.Cost),
	))

	systemPrompt := RenderText(Heading, "System Prompt") + "\n" + t.SystemPrompt + "\n"

	userPrompt := RenderText(Heading, "User Prompt") + "\n" + t.UserPrompt + "\n"

	responseText := RenderText(Heading, "Response") + "\n" + t.Response

	container := lipgloss.JoinVertical(lipgloss.Left, heading, basicInfo, systemPrompt, userPrompt, responseText)

	return BoxStyle.Width(100).Border(lipgloss.ASCIIBorder()).Padding(0, 1).Render(container)
}
