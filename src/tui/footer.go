package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m model) FooterView() string {
	bold := m.theme.TextAccent().Bold(true).Render
	base := m.theme.Base().Render
	text := m.theme.TextBody().Render
	dim := m.theme.TextDim().Render
	red := m.theme.TextRed().Render

	var footer, game_mode string
	if m.state.game_ui.game_active {
		game_mode = fmt.Sprintf("%s mode", m.state.game.Settings.PromptMode.String())
	}
	footer += game_mode

	max_footer_width := max(0, m.width_container - len(game_mode) - 3)
	if m.state.game_ui.player_damaged {
		footer = red(strings.Repeat("─", max_footer_width) + footer + strings.Repeat("─", 3))
	} else {
		footer = dim(strings.Repeat("─", max_footer_width)) + text(footer) + dim(strings.Repeat("─", 3))
	}

	table := m.theme.Base().
		Width(m.width_container - 5).
		PaddingBottom(1).
		Align(lipgloss.Center)

	keymaps := []string{}
	for _, k := range m.footer_keymaps {
		keymaps = append(keymaps, bold(" " + k.key + " ") + base(k.value + "  "))
	}

	return lipgloss.JoinVertical(
		lipgloss.Center,
		footer,
		table.Render(
			lipgloss.JoinHorizontal(
				lipgloss.Center,
				keymaps...),
		))
} 
