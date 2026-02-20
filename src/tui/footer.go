package tui

import "github.com/charmbracelet/lipgloss"

func (m model) FooterView() string {
	bold := m.theme.TextAccent().Bold(true).Render
	base := m.theme.Base().Render

	var border_style lipgloss.TerminalColor
	if m.state.game.damaged {
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

	commands := []string{}
	for _, cmd := range m.footer_cmds {
		commands = append(commands, bold(" " + cmd.key + " ") + base(cmd.value + "  "))
	}

	return lipgloss.JoinVertical(
		lipgloss.Center,
		table.Render(
			lipgloss.JoinHorizontal(
				lipgloss.Center,
				commands...,
			),
		))
} 
