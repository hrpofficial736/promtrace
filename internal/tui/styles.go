package tui

import "github.com/charmbracelet/lipgloss"

// colors
var (
	ColorPrimary   = lipgloss.Color("#FFFFFF")
	ColorSecondary = lipgloss.Color("#6B7280")
	ColorSuccess   = lipgloss.Color("#22C55E")
	ColorError     = lipgloss.Color("#EF4444")
)

// styles

var (
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary).
			PaddingLeft(1)

	AddedStyle = lipgloss.NewStyle().
			Foreground(ColorSuccess)

	RemovedStyle = lipgloss.NewStyle().
			Foreground(ColorError)

	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorSecondary).
			Padding(0, 1)
)
