package tui

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/lipgloss"
	"github.com/hrpofficial736/promtrace/internal/store"
	"github.com/hrpofficial736/promtrace/internal/util"
)

func RenderStats(stats *store.Stats) string {
	titleText := RenderText(Heading, util.GetPromtraceHeadingAsciiText()+"\nStats: ")

	summary := RenderText(Hint, fmt.Sprintf(`
Total Calls: %s
Total Tokens: %s
Total Cost: %s
Avg Latency: %sms
	`, strconv.Itoa(stats.TotalCalls), strconv.Itoa(stats.TotalTokens), util.FmtCost(stats.TotalCost), strconv.Itoa(stats.AvgLatency)))

	cols := []string{"DATE", "CALLS", "TOKENS", "COST", "AVG LATENCY"}

	var rows [][]string

	for _, t := range stats.Trend {
		row := []string{t.Date.Format("2006-01-02"), strconv.Itoa(t.Calls), strconv.Itoa(t.Tokens), strconv.Itoa(t.Cost), strconv.Itoa(t.AvgLatency) + "ms"}

		rows = append(rows, row)
	}

	trend := RenderTable(cols, rows)

	container := lipgloss.NewStyle().Padding(0, 1).Render(lipgloss.JoinVertical(lipgloss.Left, titleText, summary, trend))

	return container
}
