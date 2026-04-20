package tui

import (
	"fmt"
	"fzwds/src/tui/styles"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m model) FooterView() string {
	var footer_text_right string
	if m.game.GameActive() || m.page == game_review_page || m.page == game_over_page {
		footer_text_right = fmt.Sprintf("%s ─ %s ─ %s",
			m.game.Settings().Dictionary.String(),
			m.game.Settings().PromptMode.String(),
			m.game.Settings().WinCondition.String())
	}

	pad := 2
	max_footer_width := max(0, m.width_container - lipgloss.Width(footer_text_right) - pad)
	footer_line := strings.Repeat("─", max_footer_width) + footer_text_right + strings.Repeat("─", pad)

	if m.state.game.playerDamaged {
		footer_line = styles.TextRed.Render(footer_line)
	} else {
		footer_line = styles.TextDim.Render(footer_line)
	}

	keymaps := []string{}
	for _, k := range m.footer_keymaps {
		keymaps = append(keymaps,
			fmt.Sprintf("%s %s",
				styles.TextBlue.Bold(true).Render(k.key),
				styles.TextBody.Render(k.value)),
		)
	}

	return lipgloss.JoinVertical(
		lipgloss.Center,
		// TODO move footer msg, inline text, keymaps(?) to config struct per page that is
		// retrieved in root View() and passed to FooterView()
		m.state.footer.footer_msg,
		footer_line,
		strings.Join(keymaps, styles.TextBody.Render(" • ")),
		"",
	)
}
