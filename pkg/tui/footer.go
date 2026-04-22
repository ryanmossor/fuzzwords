package tui

import (
	"fmt"
	"fzwds/pkg/tui/pages"
	"fzwds/pkg/tui/styles"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type footerState struct {
	footerMsg		string
}

func (m model) FooterView() string {
	var footer_text_right string
	if m.game.GameActive() || m.page == pages.GameReviewPage || m.page == pages.GameOverPage {
		footer_text_right = fmt.Sprintf("%s ─ %s ─ %s",
			m.game.Settings().Dictionary.String(),
			m.game.Settings().PromptMode.String(),
			m.game.Settings().WinCondition.String())
	}

	pad := 2
	max_footer_width := max(0, m.containerWidth - lipgloss.Width(footer_text_right) - pad)
	footer_line := strings.Repeat("─", max_footer_width) + footer_text_right + strings.Repeat("─", pad)

	if m.state.game.playerDamaged {
		footer_line = styles.TextRed.Render(footer_line)
	} else {
		footer_line = styles.TextDim.Render(footer_line)
	}

	keymaps := []string{}
	for _, k := range m.footerKeymaps {
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
		m.state.footer.footerMsg,
		footer_line,
		strings.Join(keymaps, styles.TextBody.Render(" • ")),
		"",
	)
}
