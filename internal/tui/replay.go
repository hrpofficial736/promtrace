package tui

import (
	"fmt"
	"strconv"

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

	titleText := RenderText(Heading, util.GetPromtraceHeadingAsciiText()+"\nReplay Comparison: ", 0, 0)

	modelsText := fmt.Sprintf(`
	%s:  %sms, $%s, %s chars response
	%s: %sms, $%s, %s chars response
`,
		r1.Model, strconv.Itoa(r1.Latency), strconv.Itoa(r1.Cost), strconv.Itoa(r1.ResponseLength),
		r2.Model, strconv.Itoa(r2.Latency), strconv.Itoa(r2.Cost), strconv.Itoa(r2.ResponseLength),
	)

	return lipgloss.NewStyle().Padding(0, 1).Render(lipgloss.JoinVertical(lipgloss.Left, titleText, modelsText))

}
