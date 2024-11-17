package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

func (m model) HeaderUpdate(msg tea.Msg) (model, tea.Cmd) {
	// TODO: has_header flag
	if m.game_active || m.game_over {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "s":
			return m.SettingsSwitch()
		case "a":
			return m.AboutSwitch()
		case "m":
			return m.MainMenuSwitch()
		case "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) HeaderView() string {
	if m.page == game_page {
		return ""
	}
	if m.page == game_over_page {
		// No content, but used for top margin
		return "\n\n\n"
	}

	bold := m.theme.TextAccent().Bold(true).Render
	accent := m.theme.TextAccent().Render
	base := m.theme.Base().Render

	menu := accent("[m]") + base("ain menu")
	about := accent("[a]") + base("bout")
	settings := accent("[s]") + base("ettings")

	switch m.page {
	case splash_page:
		menu = bold("[m]ain menu")
	case about_page:
		about = bold("[a]bout")
	case settings_page:
		settings = bold("[s]ettings")
	}

	tabs := []string{
		menu,
		about,
		settings,
	}
		
	header := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(m.renderer.NewStyle().Foreground(m.theme.Border())).
		Row(tabs...).
		Width(m.width_container).
		StyleFunc(func(row, col int) lipgloss.Style {
			return m.theme.Base().
				Padding(0, 1).
				AlignHorizontal(lipgloss.Center)
		}).
		Render()

	return lipgloss.JoinVertical(
		lipgloss.Center,
		m.DebugView(),
		header,
	)
}
