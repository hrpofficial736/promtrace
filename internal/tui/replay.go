package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/hrpofficial736/promtrace/internal/util"
)

type ReplayComparisonResponseStruct struct {
	Model          string
	ResponseLength int
	Latency        int
	Cost           float64
}

func RenderReplay(r1, r2 *ReplayComparisonResponseStruct) string {
	// ── Title ─────────────────────────────────────────────────────────────────
	title := HeadingStyle.Render("Replay Comparison")

	// ── Comparison table ──────────────────────────────────────────────────────
	cols := []string{"", "MODEL", "LATENCY", "COST", "RESPONSE LENGTH"}
	rows := [][]string{
		{
			RenderBadge("original", ColorMuted),
			r1.Model,
			fmt.Sprintf("%dms", r1.Latency),
			util.FmtCost(r1.Cost),
			fmt.Sprintf("%d chars", r1.ResponseLength),
		},
		{
			RenderBadge("replayed", ColorAccent),
			r2.Model,
			fmt.Sprintf("%dms", r2.Latency),
			util.FmtCost(r2.Cost),
			fmt.Sprintf("%d chars", r2.ResponseLength),
		},
	}
	table := RenderTable(cols, rows)

	// ── Delta summary ─────────────────────────────────────────────────────────
	latencyDelta := r2.Latency - r1.Latency
	costDelta := r2.Cost - r1.Cost
	lengthDelta := r2.ResponseLength - r1.ResponseLength

	var latencyStr, costStr string

	if latencyDelta < 0 {
		latencyStr = SuccessStyle.Render(fmt.Sprintf("%+dms", latencyDelta)) +
			MutedStyle.Render("  faster")
	} else if latencyDelta > 0 {
		latencyStr = ErrorStyle.Render(fmt.Sprintf("%+dms", latencyDelta)) +
			MutedStyle.Render("  slower")
	} else {
		latencyStr = MutedStyle.Render("no change")
	}

	if costDelta < 0 {
		costStr = SuccessStyle.Render(fmt.Sprintf("%+.6f", costDelta)) +
			MutedStyle.Render("  cheaper")
	} else if costDelta > 0 {
		costStr = ErrorStyle.Render(fmt.Sprintf("%+.6f", costDelta)) +
			MutedStyle.Render("  more expensive")
	} else {
		costStr = MutedStyle.Render("no change")
	}

	lengthStr := BodyStyle.Render(fmt.Sprintf("%+d chars", lengthDelta))

	deltas := lipgloss.JoinVertical(lipgloss.Left,
		RenderKeyValue("latency delta", latencyStr),
		RenderKeyValue("cost delta", costStr),
		RenderKeyValue("length delta", lengthStr),
	)

	deltaSection := RenderSection("Delta Summary", deltas)

	// ── Compose ───────────────────────────────────────────────────────────────
	content := lipgloss.JoinVertical(lipgloss.Left, title, "", table, "", deltaSection)
	return lipgloss.NewStyle().Padding(1, 2).Render(content)
}
