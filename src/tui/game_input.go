package tui

import (
	"fmt"
	"fzw/src/enums"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m model) GameInputUpdate(msg tea.Msg) (model, tea.Cmd) {
	return m, nil
}

func (m model) GameInputView() string {
	if !m.game_active {
		return ""
	}

	var colorized_answer []string
	 
	switch m.settings.PromptMode {
	case enums.Fuzzy:
		prompt_caps := strings.ToUpper(m.turn.Prompt)
		prompt_idx := 0
		for _, c := range strings.ToUpper(m.text_input.Value()) {
			curr_char := string(c)

			if prompt_idx < len(m.turn.Prompt) && curr_char == string(prompt_caps[prompt_idx]) {
				colorized_answer = append(colorized_answer, m.theme.TextHighlight().Render(curr_char))
				prompt_idx++
			} else {
				colorized_answer = append(colorized_answer, m.theme.TextAccent().Render(curr_char))
			}
		}
	case enums.Classic:
		// TODO
	}

	// TODO: show possible answer after striking out
	var turn_msg string
	if !m.turn.IsValid && m.turn.Strikes < m.settings.PromptStrikesMax {
		turn_msg = m.theme.TextError().Render(m.turn.Msg)
	} else if !m.turn.IsValid && m.turn.Strikes == m.settings.PromptStrikesMax {
		turn_msg = fmt.Sprintf("Prompt failed. Possible answer: %s", m.turn.SourceWord)
	}
	// 	turn_msg = m.theme.TextHighlight().Render(m.turn.Msg)

	return lipgloss.JoinVertical(
		lipgloss.Center,
		strings.Join(colorized_answer, ""),
		"",
		turn_msg,
		m.InputField.Render(m.text_input.View()),
	) 
}
