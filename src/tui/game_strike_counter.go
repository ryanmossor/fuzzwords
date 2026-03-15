package tui

import "fmt"

func (m model) GameStrikeCounterView() string {
	if m.state.game.CurrentTurn.Strikes == 0 {
		return ""
	}

	strike_count := "Strikes: " + m.theme.TextRed().Render(fmt.Sprintf("%d/%d",
		m.state.game.CurrentTurn.Strikes,
		m.state.game.Settings.PromptStrikes))
	strike_count, _ = m.applyDamageShakeAnimation(strike_count)

	return strike_count
}
