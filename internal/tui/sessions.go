package tui

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/lipgloss"
	"github.com/hrpofficial736/promtrace/internal/store"
	"github.com/hrpofficial736/promtrace/internal/util"
)

// RenderSessions renders a styled sessions list with a heading, count subtitle,
// and a formatted table.
func RenderSessions(sessions []*store.Session) string {
	// ── Header ────────────────────────────────────────────────────────────────
	heading := HeadingStyle.Render("SESSIONS")
	subtitle := MutedStyle.Render(fmt.Sprintf("%d session(s) found", len(sessions)))

	// ── Table ─────────────────────────────────────────────────────────────────
	cols := []string{"SESSION ID", "CALLS", "STARTED AT", "AVG LATENCY", "TOKENS", "COST"}
	var rows [][]string
	for _, s := range sessions {
		rows = append(rows, []string{
			s.ID,
			strconv.Itoa(s.TotalCalls),
			s.StartedAt.Format("Jan 02 15:04:05"),
			strconv.Itoa(s.AvgLatency) + "ms",
			strconv.Itoa(s.TotalTokens),
			util.FmtCost(s.TotalCost),
		})
	}
	table := RenderTable(cols, rows)

	// ── Compose ───────────────────────────────────────────────────────────────
	inner := lipgloss.JoinVertical(
		lipgloss.Left,
		heading,
		subtitle,
		"",
		table,
	)

	return lipgloss.NewStyle().Padding(1, 2).Render(inner)
}
