package tui

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func (m model) StatsSwitch() (model, tea.Cmd) {
	m = m.SwitchPage(stats_page)

	m.footer_keymaps = []FooterKeymap{
		{key: "q", value: "quit"},
	}

	return m, nil
}

func (m model) StatsUpdate(msg tea.Msg) (model, tea.Cmd) {
    return m, nil
}

func (m model) StatsView() string {
	return m.theme.Base().Render(lipgloss.JoinVertical(
		lipgloss.Center,
        "🚧 Under construction 🚧",
	))
}
