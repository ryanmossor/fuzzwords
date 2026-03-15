package tui

import (
	"fmt"
	"fzwds/src/tui/animations"
)

func (m model) GameStrikeCounterView() string {
	if m.state.game.CurrentTurn.Strikes == 0 {
		return ""
	}

	strike_counter := "Strikes: " + m.theme.TextRed().Render(fmt.Sprintf("%d/%d",
		m.state.game.CurrentTurn.Strikes,
		m.state.game.Settings.PromptStrikes))
	strike_counter, _ = m.animation_manager.ApplyAnimations(string(animations.StrikeCounter), strike_counter)

	return strike_counter
}
