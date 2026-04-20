package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m model) StatsSwitch() (model, tea.Cmd) {
	m = m.SwitchPage(statsPage)

	m.footerKeymaps = []footerKeymap{
		{key: "q", value: "quit"},
	}

	return m, nil
}

func (m model) StatsUpdate(msg tea.Msg) (model, tea.Cmd) {
    return m, nil
}

func (m model) StatsView() string {
	return lipgloss.JoinVertical(
		lipgloss.Center,
        "🚧 Under construction 🚧",
	)
}
