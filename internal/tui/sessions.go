package tui

import (
	"strconv"

	"github.com/charmbracelet/lipgloss"
	"github.com/hrpofficial736/promtrace/internal/store"
	"github.com/hrpofficial736/promtrace/internal/util"
)

func RenderSessions(sessions []*store.Session) string {
	titleText := RenderText(Heading, util.GetPromtraceHeadingAsciiText()+"\nSessions: ", 0, 0)

	cols := []string{"SESSION ID", "CALLS", "STARTED AT", "AVG LATENCY", "TOKENS", "COST"}
	var rows [][]string

	for _, s := range sessions {
		row := []string{
			s.ID,
			strconv.Itoa(s.TotalCalls),
			s.StartedAt.Format("Jan 02 15:04:05"),
			strconv.Itoa(s.AvgLatency),
			strconv.Itoa(s.TotalTokens),
			strconv.Itoa(s.TotalCost),
		}

		rows = append(rows, row)
	}

	table := RenderTable(cols, rows)

	return lipgloss.NewStyle().Padding(0, 1).Render(lipgloss.JoinVertical(lipgloss.Left, titleText, table))
}
