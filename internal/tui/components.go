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
	Title = iota
	Heading
	Hint
	Body
)

func RenderText(textType TextType, text string) string {
	switch textType {
	case Title:
		return TitleStyle.Render(text)
	case Heading:
		return BoldHeadingStyle.Render(text)
	case Hint:
		return HintStyle.Render(text)
	default:
		return BodyStyle.Render(text)
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

func RenderCommandDescription(cmdName, description string) string {
	return RenderText(Heading, cmdName) + "\n" + RenderText(Body, description)
}

func BoxWrapper(content string) string {
	return BoxStyle.Render(content)
}
