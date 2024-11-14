package tui

import (
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

func (p page) String() string {
	switch p {
	case splash_page:
		return "Splash"
	case about_page:
		return "About"
	case settings_page:
		return "Settings"
	case game_page:
		return "Game"
	case game_over_page:
		return "Game Over"
	default:
		return "Unknown page"
	}
}

func (s size) String() string {
	switch s {
	case undersized:
		return "Undersized"
	case small:
		return "Small"
	case medium:
		return "Medium"
	case large:
		return "Large"
	default:
		return "Unknown size"
	}
}

func (m model) DebugUpdate(msg tea.Msg) (model, tea.Cmd) {
	return m, nil
}

func (m model) DebugView() string {
	if !m.debug {
		return ""
	}

	vw_vh := "VH: " + strconv.Itoa(m.viewport_height) + " | VW: " + strconv.Itoa(m.viewport_width)

	tabs := []string{
		vw_vh,
		m.page.String(),
		m.size.String(),
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
