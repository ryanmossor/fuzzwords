package tui

import (
	"fmt"
	"fzwds/src/utils"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m model) GamePromptView() string {
	if m.state.game.CurrentTurn.Strikes == 0 {
		return lipgloss.JoinVertical(
			lipgloss.Center,
			"\n\n",
			m.theme.TextAccent().Bold(true).Render(strings.ToUpper(m.state.game.CurrentTurn.Prompt)),
		)
	}

	strike_label := "Strikes: "
	plain_strikes := fmt.Sprintf("%d/%d", m.state.game.CurrentTurn.Strikes, m.state.game.Settings.PromptStrikes)
	colored_strikes := m.theme.TextRed().Render(plain_strikes)
	strike_counter := strike_label + colored_strikes
	strike_counter, padding_spaces := m.applyDamageShakeAnimation(strike_counter)

	// Length of "Strikes: x/y" plus padding; excludes terminal color codes
	strike_counter_visible_len := len(strike_label) + len(plain_strikes) + padding_spaces

	prompt := m.state.game.CurrentTurn.Prompt
	if (len(prompt) % 2) == (strike_counter_visible_len % 2) {
		// Add padding to prompt if necessary to prevent prompt from shaking
		prompt = utils.RightPad(m.state.game.CurrentTurn.Prompt, 1)
	}

	return lipgloss.JoinVertical(
		lipgloss.Center,
		"\n\n\n\n",
		m.theme.TextAccent().Bold(true).Render(strings.ToUpper(prompt)),
		strike_counter,
	)
}
