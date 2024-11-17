package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m model) GamePromptView() string {
	return lipgloss.JoinVertical(
		lipgloss.Center,
		"\n\n",
		m.theme.TextAccent().Bold(true).Render(strings.ToUpper(m.game_state.CurrentTurn.Prompt)),
	) 
}
