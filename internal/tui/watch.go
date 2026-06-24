package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	ltable "github.com/charmbracelet/lipgloss/table"
	"github.com/hrpofficial736/promtrace/internal/logger"
	"github.com/hrpofficial736/promtrace/internal/store"
	"github.com/hrpofficial736/promtrace/internal/util"
	"strconv"
	"time"
)

type WatchModel struct {
	store  store.Store
	cursor int
	traces []*store.Trace
}

func NewWatchModel(s store.Store) WatchModel {
	m := WatchModel{
		store: s,
	}

	return m
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
		traces, err := m.store.ListTraces(20)
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

	// title section

	title := RenderText(Heading, util.GetPromtraceHeadingAsciiText()+"\n", 0, 1) + RenderText(Hint, "\n 🌏 Live Traces ", 0, 1)

	var rows [][]string

	for _, t := range m.traces {
		rows = append(rows, []string{
			t.Timestamp.Format("Jan 02 15:04:05"),
			t.Model,
			strconv.Itoa(int(t.LatencyMs)),
			strconv.Itoa(t.Tokens),
			strconv.Itoa(t.Cost),
		})
	}

	t := ltable.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(ColorTertiary)).
		Headers("TIME", "MODEL", "LATENCY", "TOKENS", "COST").
		Rows(rows...).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == -1 {
				return lipgloss.NewStyle().
					Foreground(ColorSelected).
					Align(lipgloss.Center, lipgloss.Center).
					Bold(true).
					Padding(0, 2)
			}

			if row == m.cursor {
				return lipgloss.NewStyle().
					Foreground(ColorPrimary).
					Background(ColorSelected).
					Align(lipgloss.Center, lipgloss.Center).
					Bold(true).
					Padding(0, 0)
			}
			return lipgloss.NewStyle().
				Foreground(ColorSecondary).
				Align(lipgloss.Center, lipgloss.Center).
				Padding(0, 0)
		})

	help := RenderText(Hint, "\n press ↑/↓ to navigate and q to quit\n", 0, 0)

	return lipgloss.JoinVertical(lipgloss.Left, title, t.Render(), help)

}
