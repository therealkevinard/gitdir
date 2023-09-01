package clone

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/therealkevinard/gitdir/ui/styles"
)

type model struct {
	statusText string
	spinner    spinner.Model
	quitting   bool
}

func (m *model) Init() tea.Cmd { return m.spinner.Tick }

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case statusTextUpdate:
		m.statusText = string(msg)
		if string(msg) == "finished" {
			m.quitting = true
			return m, tea.Quit
		}

		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	default:
		return m, nil
	}
}

func (m *model) View() string {
	s := ""
	if m.quitting {
		return styles.QuitTextStyle.Render("okay, bye!")
	}
	if m.statusText != "" {
		s = m.spinner.View() + "  " + m.statusText
	}

	return styles.OKTextStyle.Render(s)
}
