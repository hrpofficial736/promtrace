package tui

import "github.com/charmbracelet/lipgloss"

// Colors
var (
	ColorPrimary   = lipgloss.Color("#F8FAFC") // near-white for primary text
	ColorSecondary = lipgloss.Color("#CBD5E1") // slate-300 for secondary text
	ColorMuted     = lipgloss.Color("#64748B") // slate-500 for muted / hints
	ColorAccent    = lipgloss.Color("#818CF8") // indigo-400 for accents
	ColorAccentDim = lipgloss.Color("#4F46E5") // indigo-600 for dimmed accent
	ColorSuccess   = lipgloss.Color("#22C55E") // green-500
	ColorError     = lipgloss.Color("#EF4444") // red-500
	ColorWarning   = lipgloss.Color("#F59E0B") // amber-500
	ColorBorder    = lipgloss.Color("#334155") // slate-700 for borders
	ColorHighlight = lipgloss.Color("#1E293B") // slate-800 for row highlights

	// Legacy aliases used by watch.go — kept for backward compatibility
	ColorTertiary   = ColorAccent
	ColorSelected   = ColorMuted
	ColorHeaderText = ColorMuted
	ColorPrettyText = ColorAccent
)

// Text styles
var (
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(ColorSecondary)

	HeadingStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorAccent)

	// BoldHeadingStyle is a legacy alias kept for backward compatibility.
	BoldHeadingStyle = HeadingStyle

	LabelStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorMuted).
			MarginBottom(0)

	BodyStyle = lipgloss.NewStyle().
			Foreground(ColorSecondary)

	MutedStyle = lipgloss.NewStyle().
			Foreground(ColorMuted)

	HintStyle = lipgloss.NewStyle().
			Foreground(ColorMuted)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(ColorSuccess)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(ColorError)

	AddedStyle = lipgloss.NewStyle().
			Foreground(ColorSuccess)

	RemovedStyle = lipgloss.NewStyle().
			Foreground(ColorError)
)

// Layout / container styles
var (
	// BoxStyle is a plain padding wrapper.
	BoxStyle = lipgloss.NewStyle().Padding(1, 2)

	// PanelStyle is a bordered panel with subtle styling.
	PanelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorBorder).
			Padding(1, 2)

	// AccentBorderPanelStyle is a panel with an accent-colored border.
	AccentBorderPanelStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(ColorAccent).
				Padding(1, 2)

	// TableHeaderStyle styles the header row of a table.
	TableHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(ColorAccent).
				Align(lipgloss.Center).
				Padding(0, 2)

	// TableRowStyle styles ordinary table rows.
	TableRowStyle = lipgloss.NewStyle().
			Foreground(ColorSecondary).
			Align(lipgloss.Center).
			Padding(0, 2)

	// TableRowAltStyle styles alternate table rows.
	TableRowAltStyle = lipgloss.NewStyle().
				Foreground(ColorPrimary).
				Background(ColorHighlight).
				Align(lipgloss.Center).
				Padding(0, 2)

	// SelectedRowStyle highlights the currently selected row.
	SelectedRowStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(ColorPrimary).
				Background(ColorAccentDim).
				Align(lipgloss.Center).
				Padding(0, 2)
)
