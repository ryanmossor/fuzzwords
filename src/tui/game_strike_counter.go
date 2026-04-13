package tui

import (
	"fmt"
	"fzwds/src/tui/animations"
	"fzwds/src/tui/styles"
)

func (m model) GameStrikeCounterView() string {
	if m.state.game.CurrentTurn().Strikes == 0 {
		return ""
	}

	strike_counter := styles.TextBody.Render("Strikes: ")
	strike_counter += styles.TextRed.Render(fmt.Sprintf("%d/%d",
		m.state.game.CurrentTurn().Strikes,
		m.state.game.Settings.PromptStrikes))

	strike_counter, _ = m.anim_mgr.ApplyAnimations(
		string(animations.StrikeCounter),
		strike_counter)

	return strike_counter
}
