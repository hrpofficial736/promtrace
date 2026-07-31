package tui

import (
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	ltable "github.com/charmbracelet/lipgloss/table"
	"github.com/hrpofficial736/promtrace/internal/logger"
	"github.com/hrpofficial736/promtrace/internal/store"
	"github.com/hrpofficial736/promtrace/internal/util"
)

type WatchModel struct {
	store  store.Store
	cursor int
	traces []*store.Trace
	limit  int
	info   string
}

func NewWatchModel(s store.Store, l int) WatchModel {
	return WatchModel{
		store: s,
		limit: l,
	}
}

type tickMsg time.Time
type tracesMsg []*store.Trace

func tickEvery() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m WatchModel) fetchTraces() tea.Cmd {
	return func() tea.Msg {
		traces, err := m.store.ListTraces(m.limit)
		if err != nil {
			logger.Log.Error("error", "error", err)
			return err
		}
		return tracesMsg(traces)
	}
}

func (m WatchModel) Init() tea.Cmd {
	return tea.Batch(m.fetchTraces(), tickEvery())
}

func (m WatchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down":
			if m.cursor < len(m.traces)-1 {
				m.cursor++
			}
		case "enter":
			if len(m.traces) > 0 {
				m.info = RenderTraceInfoContainer(m.traces[m.cursor]) +
					lipgloss.NewStyle().
						Padding(1, 1).
						Render(
							HintStyle.Render("press 'q' to quit and 'backspace' to go back to watch traces."),
						)
			}
			return m, nil
		case "backspace":
			m.info = ""
			return m, nil
		}
	case tickMsg:
		return m, tea.Batch(m.fetchTraces(), tickEvery())
	case tracesMsg:
		m.traces = msg
		if m.cursor >= len(m.traces) {
			m.cursor = max(0, len(m.traces)-1)
		}
	}
	return m, nil
}

func (m WatchModel) View() string {
	if m.info != "" {
		return m.info
	}
	const width = 80

	// ── Header ───────────────────────────────────────────────────────────────
	left := HeadingStyle.Render("WATCH YOUR TRACES HERE")
	right := HintStyle.Render("auto-refresh 1s")
	gap := strings.Repeat(" ", max(0, width-lipgloss.Width(left)-lipgloss.Width(right)))
	header := left + gap + right

	divider := lipgloss.NewStyle().Foreground(ColorBorder).Render(strings.Repeat("─", width))

	// ── Body ─────────────────────────────────────────────────────────────────
	var body string

	if len(m.traces) == 0 {
		body = lipgloss.NewStyle().
			Width(width).
			Align(lipgloss.Center).
			Padding(2, 0).
			Render(MutedStyle.Render("waiting for traces…"))
	} else {
		var rows [][]string
		for _, t := range m.traces {
			rows = append(rows, []string{
				t.Timestamp.Format("15:04:05"),
				t.ID,
				t.Model,
				strconv.Itoa(int(t.LatencyMs)) + "ms",
				strconv.Itoa(t.Tokens),
				util.FmtCost(t.Cost),
			})
		}

		cursor := m.cursor
		t := ltable.New().
			Border(lipgloss.RoundedBorder()).
			BorderStyle(lipgloss.NewStyle().Foreground(ColorBorder)).
			Headers("TIME", "ID", "MODEL", "LATENCY", "TOKENS", "COST").
			Rows(rows...).
			StyleFunc(func(row, col int) lipgloss.Style {
				switch {
				case row == ltable.HeaderRow:
					return TableHeaderStyle.
						Padding(0, 2).
						Align(lipgloss.Center)
				case row == cursor:
					return SelectedRowStyle.Padding(0, 2)
				case row%2 == 0:
					return TableRowStyle.Padding(0, 2)
				default:
					return TableRowAltStyle.Padding(0, 2)
				}
			})

		body = t.Render()
	}

	// ── Footer ────────────────────────────────────────────────────────────────
	footer := HintStyle.Render("↑/↓ navigate   q quit")

	container := lipgloss.NewStyle().
		Padding(1, 1).
		Render(
			lipgloss.JoinVertical(lipgloss.Left, header, divider, body, footer),
		)

	return container
}
