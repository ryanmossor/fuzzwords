package tui

import "github.com/charmbracelet/lipgloss"

func (m model) FooterView() string {
	bold := m.theme.TextAccent().Bold(true).Render
	base := m.theme.Base().Render

	var border_style lipgloss.TerminalColor
	if m.state.game_ui.player_damaged {
		border_style = m.theme.red
	} else {
		border_style = m.theme.Border()
	}

	table := m.theme.Base().
		Width(m.width_container).
		BorderTop(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(border_style).
		PaddingBottom(1).
		Align(lipgloss.Center)

	keymaps := []string{}
	for _, k := range m.footer_keymaps {
		keymaps = append(keymaps, bold(" " + k.key + " ") + base(k.value + "  "))
	}

	return lipgloss.JoinVertical(
		lipgloss.Center,
		table.Render(
			lipgloss.JoinHorizontal(
				lipgloss.Center,
				keymaps...,
			),
		))
} 
