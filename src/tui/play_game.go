package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m model) GameSwitch() (model, tea.Cmd) {
	m = m.SwitchPage(game_page)
	return m, nil
}

func (m model) GameUpdate(msg tea.Msg) (model, tea.Cmd) {
	return m, nil
}

func (m model) GameView() string {
	return lipgloss.JoinVertical(
		lipgloss.Center,
		m.theme.TextAccent().Bold(true).Render("Starting game..."),
	)
}
