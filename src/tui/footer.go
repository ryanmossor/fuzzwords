package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m model) FooterView() string {
	bold := m.theme.TextAccent().Bold(true).Render
	base := m.theme.Base().Render
	dim := m.theme.TextDim().Render
	red := m.theme.TextRed().Render

	var footer_text string
	if m.state.game_ui.game_active {
		footer_text = fmt.Sprintf("%s/%s",
			m.state.game.Settings.PromptMode.String(),
			m.state.game.Settings.WinCondition.String())
	}

	right_pad := 3
	max_footer_width := max(0, m.width_container - len(footer_text) - right_pad)
	footer_line := strings.Repeat("─", max_footer_width) + footer_text + strings.Repeat("─", right_pad)

	if m.state.game_ui.player_damaged {
		footer_line = red(footer_line)
	} else {
		footer_line = dim(footer_line)
	}

	table := m.theme.Base().
		Width(m.width_container).
		PaddingBottom(1).
		Align(lipgloss.Center)

	keymaps := []string{}
	for _, k := range m.footer_keymaps {
		keymaps = append(keymaps, bold(" " + k.key + " ") + base(k.value + "  "))
	}

	return lipgloss.JoinVertical(
		lipgloss.Center,
		footer_line,
		table.Render(lipgloss.JoinHorizontal(lipgloss.Center, keymaps...)),
	)
} 
