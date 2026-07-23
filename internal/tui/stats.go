package tui

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/lipgloss"
	"github.com/hrpofficial736/promtrace/internal/store"
	"github.com/hrpofficial736/promtrace/internal/util"
)

// statCard renders a single summary card with a label and a prominent value.
func statCard(label, value string) string {
	top := LabelStyle.Render(label)
	bottom := TitleStyle.Bold(true).Render(value)
	inner := lipgloss.JoinVertical(lipgloss.Left, top, bottom)
	return PanelStyle.Render(inner)
}

// RenderStats renders a stats overview with summary cards and a daily trend table.
func RenderStats(stats *store.Stats) string {
	// ── Summary cards ────────────────────────────────────────────────────────
	cards := lipgloss.JoinHorizontal(
		lipgloss.Top,
		statCard("TOTAL CALLS", strconv.Itoa(stats.TotalCalls)),
		statCard("TOTAL TOKENS", strconv.Itoa(stats.TotalTokens)),
		statCard("TOTAL COST", util.FmtCost(stats.TotalCost)),
		statCard("AVG LATENCY", fmt.Sprintf("%.2fms", stats.AvgLatency)),
	)

	// ── Daily trend table ─────────────────────────────────────────────────────
	trendHeading := HeadingStyle.Render("DAILY TREND")

	cols := []string{"DATE", "CALLS", "TOKENS", "COST", "AVG LATENCY"}
	var rows [][]string
	for _, t := range stats.Trend {
		rows = append(rows, []string{
			t.Date.Format("02 Jan 2006"),
			strconv.Itoa(t.Calls),
			strconv.Itoa(t.Tokens),
			util.FmtCost(t.Cost),
			fmt.Sprintf("%.2fms", t.AvgLatency),
		})
	}
	trendTable := RenderTable(cols, rows)

	// ── Compose ───────────────────────────────────────────────────────────────
	body := lipgloss.JoinVertical(
		lipgloss.Left,
		cards,
		"",
		trendHeading,
		"",
		trendTable,
	)

	return lipgloss.NewStyle().Padding(1, 2).Render(body)
}
