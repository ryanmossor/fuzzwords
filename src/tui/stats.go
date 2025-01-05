package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m model) StatsSwitch() (model, tea.Cmd) {
	m = m.SwitchPage(stats_page)

	m.footer_cmds = []footerCmd{
		{key: "↑/↓", value: "scroll"},
		{key: "←/→", value: "change"},
		{key: "ctrl+r", value: "restore defaults"},
		{key: "enter", value: "start"},
	}

	return m, nil
}

func (m model) StatsUpdate(msg tea.Msg) (model, tea.Cmd) {
    return m, nil
}

func (m model) StatsView() string {
	return m.theme.Base().Render(lipgloss.JoinVertical(
		lipgloss.Left,
		// lines...,
        "TODO",
	)) 
}
