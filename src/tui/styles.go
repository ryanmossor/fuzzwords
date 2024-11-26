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
		style = base.BorderForeground(m.theme.Accent()).BorderStyle(lipgloss.DoubleBorder())
	} else {
		style = base.BorderForeground(m.theme.Border()).BorderStyle(lipgloss.NormalBorder())
	}

	return style.PaddingLeft(1).Render(padded)
}
