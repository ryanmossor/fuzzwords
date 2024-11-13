package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
) 

func (m model) MainMenuSwitch() (model, tea.Cmd) {
	m = m.SwitchPage(splash_page)
	return m, nil
}

func (m model) MainMenuUpdate(msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return m.GameSwitch()
		case "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) MainMenuView() string {
	// base := m.theme.Base().Width(m.widthContent).Render
	base := m.theme.Base().Render
	accent := m.theme.TextAccent().Render 

	var title []string
	switch m.size {
	case large:
		title = append(title, accent("███████╗███████╗██╗    ██╗██████╗ ███████╗"))
		title = append(title, accent("██╔════╝╚══███╔╝██║    ██║██╔══██╗██╔════╝"))
		title = append(title, accent("█████╗    ███╔╝ ██║ █╗ ██║██║  ██║███████╗"))
		title = append(title, accent("██╔══╝   ███╔╝  ██║███╗██║██║  ██║╚════██║"))
		title = append(title, accent("██║     ███████╗╚███╔███╔╝██████╔╝███████║"))
		title = append(title, accent("╚═╝     ╚══════╝ ╚══╝╚══╝ ╚═════╝ ╚══════╝"))
	default:
		title = append(title, accent(" ___ __       __   __  "))
		title = append(title, accent("|__   / |  | |  \\ /__` "))
		title = append(title, accent("|    /_ |/\\| |__/ .__/ "))
	}

	title = append(title, "\n\nPress " + accent("ENTER") + base(" to play"))

	return lipgloss.JoinVertical(
		lipgloss.Center,
		title...
	)
}
