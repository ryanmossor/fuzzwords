package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

func (m model) HeaderUpdate(msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// if m.hasMenu {
			switch msg.String() {
			case "s":
				return m.SettingsSwitch()
			case "a":
				return m.AboutSwitch()
			case "m":
				return m.MainMenuSwitch()
			// case "q":
				// return m, tea.Quit
				// fmt.Println("You pressed", msg.String())
			}
		// }
	}

	return m, nil
}

func (m model) HeaderView() string {
	// bold := m.theme.TextAccent().Bold(true).Render
	accent := m.theme.TextAccent().Render
	base := m.theme.Base().Render

	// back := base("← ") + bold("esc") + base(" back")
	menu := accent("[m]") + base("ain menu")
	about := accent("[a]") + base("bout")
	settings := accent("[s]") + base("ettings")

	// cart :=
	// 	accent("c") +
	// 		base(" cart") +
	// 		accent(fmt.Sprintf(" $%2v", total/100)) +
	// 		base(fmt.Sprintf(" [%d]", count))

	// don't think i need game, gameover pages for header? only play, config, about
	switch m.page {
	case splash_page:
		menu = accent("[m]ain menu")
	case about_page:
		about = accent("[a]bout")
	case settings_page:
		settings = accent("[s]ettings")
	}

	var tabs []string

	switch m.size {
	// case small:
	// 	tabs = []string{
	// 		menu,
	// 		about,
	// 		settings,
	// 	}
	// case medium:
	// 	tabs = []string{
	// 		menu,
	// 		about,
	// 		settings,
	// 	}
	default:
		tabs = []string{
			menu,
			about,
			settings,
		}
	}

	return table.New().
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
}
