package tui

import (
	"github.com/charmbracelet/lipgloss"
)

func (m model) CreateBox(content string, selected bool) string {
	total_width := m.width_content - 2
	padded := lipgloss.PlaceHorizontal(total_width, lipgloss.Left, content)
	base := m.theme.Base().Border(lipgloss.NormalBorder()).Width(total_width)

	var style lipgloss.Style
	if selected {
		style = base.BorderForeground(m.theme.Accent()).BorderStyle(lipgloss.RoundedBorder())
	} else {
		style = base.BorderForeground(m.theme.Border()).BorderStyle(lipgloss.RoundedBorder())
	}

	return style.PaddingLeft(1).Render(padded)
}

func (m model) CreateSettingsMenuItem(content string, is_selected, apply_bottom_border bool) string {
	total_width := m.width_content - 2
	padded := lipgloss.PlaceHorizontal(total_width, lipgloss.Left, content)
	base := m.theme.Base().
		BorderBottom(apply_bottom_border).
		BorderForeground(m.theme.Border()).
		BorderStyle(lipgloss.NormalBorder()).
		Width(total_width)

	return base.PaddingLeft(1).Render(padded)
}

func (m model) TextInputBlockBorderStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		BorderForeground(m.theme.input_bg).
		BorderStyle(lipgloss.InnerHalfBlockBorder()).
		Width(m.text_input.CharLimit)
}

func (m model) TextInputRoundedBorderStyle(border_color lipgloss.TerminalColor) lipgloss.Style {
	return lipgloss.NewStyle().
		BorderForeground(border_color).
		BorderStyle(lipgloss.RoundedBorder()).
		Width(m.text_input.CharLimit)
}
