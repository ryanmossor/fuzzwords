package tui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m model) GamePromptView() string {
	if m.game_state.CurrentTurn.Strikes > 0 {
		strike_count := "\nStrikes: " + m.theme.TextRed().Render(
			fmt.Sprintf("%s/%s",
				strconv.Itoa(m.game_state.CurrentTurn.Strikes),
				strconv.Itoa(m.game_state.Settings.PromptStrikesMax)))

		return lipgloss.JoinVertical(
			lipgloss.Center,
			"\n\n\n\n",
			m.theme.TextAccent().Bold(true).Render(strings.ToUpper(m.game_state.CurrentTurn.Prompt)),
			strike_count,
		)
	}

	return lipgloss.JoinVertical(
		lipgloss.Center,
		"\n\n",
		m.theme.TextAccent().Bold(true).Render(strings.ToUpper(m.game_state.CurrentTurn.Prompt)),
	) 
}
