package tui

import "github.com/charmbracelet/lipgloss"

// colors
var (
	ColorPrimary    = lipgloss.Color("#FFFFFF")
	ColorSecondary  = lipgloss.Color("#D1D5DB")
	ColorTertiary   = lipgloss.Color("#5E4987")
	ColorSuccess    = lipgloss.Color("#22C55E")
	ColorError      = lipgloss.Color("#EF4444")
	ColorHeaderText = lipgloss.Color("#A78BFA")
	ColorSelected   = lipgloss.Color("#7C3AED")
)

// styles

var (
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorSelected)

	HintStyle = lipgloss.NewStyle().
			Bold(false).
			Foreground(ColorHeaderText)

	AddedStyle = lipgloss.NewStyle().
			Foreground(ColorSuccess)

	RemovedStyle = lipgloss.NewStyle().
			Foreground(ColorError)

	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(ColorTertiary).
			Padding(0, 1)
)
