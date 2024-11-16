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

	accent := m.theme.TextAccent().Render
	blue := m.theme.TextBlue().Render

	prompt_upper := strings.ToUpper(m.turn.Prompt)
	answer_upper := strings.ToUpper(m.text_input.Value())
	var sb strings.Builder
	 
	switch m.settings.PromptMode {
	case enums.Fuzzy:
		prompt_idx := 0
		for _, c := range answer_upper {
			curr_char := string(c)

			if prompt_idx < len(prompt_upper) && curr_char == string(prompt_upper[prompt_idx]) {
				sb.WriteString(blue(curr_char))
				prompt_idx++
			} else {
				sb.WriteString(accent(curr_char))
			}
		}
	case enums.Classic:
		if !strings.Contains(answer_upper, prompt_upper) {
			sb.WriteString(accent(answer_upper))
			break
		}
		
		sub_idx := strings.Index(answer_upper, prompt_upper)
		sb.WriteString(accent(answer_upper[0:sub_idx]))
		sb.WriteString(blue(answer_upper[sub_idx:sub_idx + len(prompt_upper)]))
		sb.WriteString(accent(answer_upper[sub_idx + len(prompt_upper):]))
	}

	// TODO: show possible answer after striking out
	var turn_msg string
	if !m.turn.IsValid && m.turn.Strikes < m.settings.PromptStrikesMax {
		turn_msg = m.theme.TextRed().Render(m.turn.Msg)
	} else if !m.turn.IsValid && m.turn.Strikes == m.settings.PromptStrikesMax {
		turn_msg = fmt.Sprintf("Prompt failed. Possible answer: %s", m.turn.SourceWord)
	}
	// 	turn_msg = m.theme.TextHighlight().Render(m.turn.Msg)

	return lipgloss.JoinVertical(
		lipgloss.Center,
		sb.String(),
		"",
		turn_msg,
		m.InputField.Render(m.text_input.View()),
	) 
}
