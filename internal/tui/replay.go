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
	Cost           int
}

func RenderReplay(r1, r2 *ReplayComparisonResponseStruct) string {

	titleText := RenderText(Heading, util.GetPromtraceHeadingAsciiText()+"\nReplay Comparison: ")

	cols := []string{"", "MODEL", "LATENCY", "COST", "RESPONSE"}

	rows := [][]string{
		{"original", r1.Model, fmt.Sprintf("%dms", r1.Latency), util.FmtCost(r1.Cost), fmt.Sprintf("%d chars", r1.ResponseLength)},
		{"replayed", r2.Model, fmt.Sprintf("%dms", r2.Latency), util.FmtCost(r2.Cost), fmt.Sprintf("%d chars", r2.ResponseLength)},
	}

	table := RenderTable(cols, rows)

	return lipgloss.NewStyle().Padding(0, 1).Render(lipgloss.JoinVertical(lipgloss.Left, titleText, table))

}
