package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	ltable "github.com/charmbracelet/lipgloss/table"
)

// TextType selects a pre-defined text style in RenderText.
type TextType int

const (
	Title   TextType = iota
	Heading TextType = iota
	Label   TextType = iota
	Body    TextType = iota
	Muted   TextType = iota
	Hint    TextType = iota
)

// RenderText renders text using the named style.
func RenderText(textType TextType, text string) string {
	switch textType {
	case Title:
		return TitleStyle.Render(text)
	case Heading:
		return HeadingStyle.Render(text)
	case Label:
		return LabelStyle.Render(text)
	case Muted:
		return MutedStyle.Render(text)
	case Hint:
		return HintStyle.Render(text)
	default:
		return BodyStyle.Render(text)
	}
}

// RenderTable renders a bordered, styled table with alternating row colors.
func RenderTable(cols []string, rows [][]string) string {
	t := ltable.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(ColorBorder)).
		Headers(cols...).
		Rows(rows...).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == ltable.HeaderRow {
				return TableHeaderStyle
			}
			if row%2 == 0 {
				return TableRowAltStyle
			}
			return TableRowStyle
		})

	return t.Render()
}

// RenderKeyValue renders a single label: value line.
func RenderKeyValue(label, value string) string {
	return LabelStyle.Render(label+":") + " " + BodyStyle.Render(value)
}

// RenderSection wraps content in an accent-bordered panel with a heading.
func RenderSection(title, content string) string {
	heading := HeadingStyle.Render(title)
	inner := lipgloss.JoinVertical(lipgloss.Left, heading, "", content)
	return AccentBorderPanelStyle.Render(inner)
}

// RenderDivider returns a horizontal rule of the given width using muted characters.
func RenderDivider(width int) string {
	return MutedStyle.Render(strings.Repeat("─", width))
}

// RenderBadge renders a small inline badge with a given background color.
func RenderBadge(text string, color lipgloss.Color) string {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorPrimary).
		Background(color).
		Padding(0, 1).
		Render(text)
}

// RenderStatus renders a success or failure indicator line.
func RenderStatus(success bool, text string) string {
	if success {
		return SuccessStyle.Render("✓ " + text)
	}
	return ErrorStyle.Render("✗ " + text)
}

// RenderCommandDescription renders a command name and description block.
func RenderCommandDescription(cmdName, description string) string {
	header := HeadingStyle.Render(fmt.Sprintf("'%s'", cmdName)) + BodyStyle.Render(" command")
	divider := RenderDivider(40)
	return lipgloss.JoinVertical(lipgloss.Left, header, divider, BodyStyle.Render(description))
}

// BoxWrapper wraps content in the default box style.
func BoxWrapper(content string) string {
	return BoxStyle.Render(content)
}
