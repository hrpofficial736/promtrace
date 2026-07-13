package tui

import "github.com/charmbracelet/lipgloss"

// colors
var (
	ColorPrimary    = lipgloss.Color("#FFFFFF")
	ColorSecondary  = lipgloss.Color("#D1D5DB")
	ColorTertiary   = lipgloss.Color("#5E4987")
	ColorSuccess    = lipgloss.Color("#22C55E")
	ColorError      = lipgloss.Color("#EF4444")
	ColorHeaderText = lipgloss.Color("#999E9E")
	ColorPrettyText = lipgloss.Color("#E5226D")
	ColorSelected   = lipgloss.Color("#999E9E")
)

// styles

var (
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorHeaderText)

	BoldHeadingStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(ColorHeaderText)

	HintStyle = lipgloss.NewStyle().
			Bold(false).
			Foreground(ColorHeaderText)

	BodyStyle = lipgloss.NewStyle().
			Bold(false).
			Foreground(ColorSecondary)

	AddedStyle = lipgloss.NewStyle().
			Foreground(ColorSuccess)

	RemovedStyle = lipgloss.NewStyle().
			Foreground(ColorError)

	BoxStyle = lipgloss.NewStyle().Padding(1, 1)
)
