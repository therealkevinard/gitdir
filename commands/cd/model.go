package cd

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/therealkevinard/gitdir/ui/styles"
)

type stringListItem string

func (i stringListItem) FilterValue() string { return string(i) }

type model struct {
	list     list.Model
	choice   string
	quitting bool
	command  *Command
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(stringListItem)
			if ok {
				m.command.selection = string(i)
			}

			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

func (m *model) View() string {
	if m.choice != "" {
		return styles.QuitTextStyle.Render("")
	}
	if m.quitting {
		return styles.QuitTextStyle.Render("")
	}

	return "\n" + m.list.View()
}
