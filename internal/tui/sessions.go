package tui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/hrpofficial736/promtrace/internal/store"
	"github.com/hrpofficial736/promtrace/internal/util"
	"strconv"
)

func RenderSessions(sessions []*store.Session) string {
	titleText := RenderText(Heading, "SESSIONS: ")

	cols := []string{"SESSION ID", "CALLS", "STARTED AT", "AVG LATENCY", "TOKENS", "COST"}
	var rows [][]string

	for _, s := range sessions {
		row := []string{
			s.ID,
			strconv.Itoa(s.TotalCalls),
			s.StartedAt.Format("Jan 02 15:04:05"),
			strconv.Itoa(s.AvgLatency) + "ms",
			strconv.Itoa(s.TotalTokens),
			util.FmtCost(s.TotalCost),
		}

		rows = append(rows, row)
	}

	table := RenderTable(cols, rows)

	return BoxStyle.Border(lipgloss.ASCIIBorder()).Render(lipgloss.JoinVertical(lipgloss.Left, titleText, table))
}
