package tui

import (
	"io"
	"net/http"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hrpofficial736/promtrace/internal/replay"
	"github.com/hrpofficial736/promtrace/internal/store"
)

type SpinnerModel struct {
	spinner spinner.Model
	done    bool
	Err     error

	t         *store.Trace
	modelFlag string

	Res  *http.Response
	Body []byte
}

type DoneMsg struct {
	res  *http.Response
	body []byte
	err  error
}

func NewSpinnerModel(t *store.Trace, modelFlag string) SpinnerModel {
	s := spinner.New()
	s.Spinner = spinner.MiniDot
	s.Style = lipgloss.NewStyle().Foreground(ColorAccent)

	return SpinnerModel{
		spinner:   s,
		t:         t,
		modelFlag: modelFlag,
	}
}

func doReplay(t *store.Trace, modelFlag string) tea.Cmd {
	return func() tea.Msg {
		res, err := replay.ReplayRequest(t, modelFlag)
		if err != nil {
			return DoneMsg{
				err: err,
			}
		}

		body, _ := io.ReadAll(res.Body)
		err = res.Body.Close()
		if err != nil {
			return DoneMsg{
				err: err,
			}
		}
		return DoneMsg{res: res, body: body, err: err}
	}
}

func (m SpinnerModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, doReplay(m.t, m.modelFlag))
}

func (m SpinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case DoneMsg:
		m.done = true
		m.Res = msg.res
		m.Body = msg.body
		m.Err = msg.err
		return m, tea.Quit
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m SpinnerModel) View() string {
	if m.done {
		return ""
	}

	return m.spinner.View() + " replaying trace..."
}
