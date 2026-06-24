package tui

import (
	"github.com/charmbracelet/lipgloss"
	ltable "github.com/charmbracelet/lipgloss/table"
)

func RenderStatus(success bool, text string) string {
	var out string
	if success {
		out = "✅ " + text
	} else {
		out = "❌ " + text
	}
	return AddedStyle.Render(out)
}

type TextType int

const (
	Heading = iota
	Hint
)

func RenderText(textType TextType, text string, x, y int) string {
	switch textType {
	case Heading:
		return TitleStyle.Padding(x, y).Render(text)
	case Hint:
		return HintStyle.Padding(x, y).Render(text)
	default:
		return HintStyle.Padding(x, y).Render(text)
	}
}

func RenderTable(cols []string, rows [][]string) string {

	t := ltable.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(ColorTertiary)).
		Headers(cols...).
		Rows(rows...)

	return t.Render()
}
