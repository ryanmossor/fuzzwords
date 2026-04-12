package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m model) FooterView() string {
	var footer_text_right string
	if m.state.game.GameActive || m.page == game_review_page || m.page == game_over_page {
		footer_text_right = fmt.Sprintf("%s / %s / %s",
			m.state.game.Settings.Dictionary.String(),
			m.state.game.Settings.PromptMode.String(),
			m.state.game.Settings.WinCondition.String())
	}

	pad := 2
	max_footer_width := max(0, m.width_container - len(footer_text_right) - pad)
	footer_line := strings.Repeat("─", max_footer_width) + footer_text_right + strings.Repeat("─", pad)

	if m.state.game_ui.player_damaged {
		footer_line = m.theme.TextRed().Render(footer_line)
	} else {
		footer_line = m.theme.TextDim().Render(footer_line)
	}

	table := m.theme.Base().
		Width(m.width_container).
		PaddingBottom(1).
		Align(lipgloss.Center)

	keymaps := []string{}
	for _, k := range m.footer_keymaps {
		keymaps = append(keymaps,
			fmt.Sprintf(" %s %s  ",
				m.theme.TextAccent().Bold(true).Render(k.key),
				m.theme.Base().Render(k.value)),
		)
	}

	return lipgloss.JoinVertical(
		lipgloss.Center,
		// TODO move footer msg, inline text, keymaps(?) to config struct per page that is
		// retrieved in root View() and passed to FooterView()
		m.state.footer.footer_msg,
		footer_line,
		table.Render(lipgloss.JoinHorizontal(lipgloss.Center, keymaps...)),
	)
}
