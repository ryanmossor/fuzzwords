package tui

import (
	"fmt"
	"fzwds/pkg/tui/animations"
	"fzwds/pkg/tui/styles"
)

func (m model) GameStrikeCounterView() string {
	if m.state.game.turn.strikes == 0 {
		return ""
	}

	strike_counter := styles.TextBody.Render("Strikes: ")
	strike_counter += styles.TextRed.Render(fmt.Sprintf("%d/%d",
		m.state.game.turn.strikes,
		m.game.Settings().PromptStrikes))

	strike_counter, _ = m.animManager.ApplyAnimations(
		string(animations.StrikeCounter),
		strike_counter)

	return strike_counter
}
