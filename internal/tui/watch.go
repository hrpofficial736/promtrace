package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hrpofficial736/promtrace/internal/logger"
	"github.com/hrpofficial736/promtrace/internal/store"
	"strconv"
	"time"
)

type WatchModel struct {
	store store.Store
	table table.Model
}

func NewWatchModel(s store.Store) WatchModel {

	t := table.New(
		table.WithColumns([]table.Column{
			{Title: "time", Width: 20},
			{Title: "model", Width: 18},
			{Title: "latency", Width: 10},
			{Title: "tokens", Width: 8},
			{Title: "cost", Width: 8},
		}),
		table.WithFocused(true),
		table.WithHeight(15),
	)
	st := table.DefaultStyles()

	st.Header = st.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		Bold(true)

	t.SetStyles(st)

	m := WatchModel{
		store: s,
		table: t,
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
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tickMsg:
		return m, tea.Batch(m.fetchTraces(), tickEvery())
	case tracesMsg:
		var rows []table.Row

		for _, t := range msg {
			row := []string{t.Timestamp.String(), t.Model, fmt.Sprintf("%dms", t.LatencyMs), strconv.Itoa(t.Tokens), strconv.Itoa(t.Cost)}
			rows = append(rows, table.Row(row))
		}
		m.table.SetRows(rows)
	}

	var cmd tea.Cmd

	m.table, cmd = m.table.Update(msg)

	return m, cmd
}

func (m WatchModel) View() string {
	return m.table.View()
}
