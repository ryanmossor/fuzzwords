package tui

import (
	"fzwds/src/enums"
	"fzwds/src/tui/animations"
	"fzwds/src/tui/styles"
	"fzwds/src/tui/theme"
	"fzwds/src/utils"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

// Highlight prompt letters in current answer
func (m model) highlightPromptAnswer(prompt, answer string, prompt_mode enums.PromptMode) string {
	accent := styles.TextAccent.Render
	highlight := styles.TextHighlight.Render

	prompt_upper := strings.ToUpper(prompt)
	answer_upper := strings.ToUpper(answer)
	var sb strings.Builder

	switch prompt_mode {
	case enums.PromptModeFuzzy:
		prompt_idx := 0
		for _, c := range answer_upper {
			curr_char := string(c)

			if prompt_idx < len(prompt_upper) && curr_char == string(prompt_upper[prompt_idx]) {
				sb.WriteString(highlight(curr_char))
				prompt_idx++
			} else {
				sb.WriteString(accent(curr_char))
			}
		}
	case enums.PromptModeClassic:
		if !strings.Contains(answer_upper, prompt_upper) {
			sb.WriteString(accent(answer_upper))
			return sb.String()
		}

		sub_idx := strings.Index(answer_upper, prompt_upper)
		sb.WriteString(accent(answer_upper[0:sub_idx]))
		sb.WriteString(highlight(answer_upper[sub_idx:sub_idx + len(prompt_upper)]))
		sb.WriteString(accent(answer_upper[sub_idx + len(prompt_upper):]))
	}

	return sb.String()
}

// Get accent color for input box based on HighlightInput setting, damage state, etc.
// Style applied to border if rounded-style input box, or left accent bar if block-style input box.
func (m model) getInputAccentColor(default_color lipgloss.TerminalColor) lipgloss.TerminalColor {
	if m.state.game_ui.player_damaged {
		return theme.Red
	}

	prompt_upper := strings.ToUpper(m.game.CurrentTurn().Prompt)
	answer_upper := strings.ToUpper(m.text_input.Value())

	is_match := false
	switch m.game.Settings.PromptMode {
	case enums.PromptModeFuzzy:
		is_match = utils.IsFuzzyMatch(answer_upper, prompt_upper)
	case enums.PromptModeClassic:
		is_match = strings.Contains(answer_upper, prompt_upper)
	}

	if m.game.Settings.HighlightInput {
		valid_word := m.game.WordInDictionary(answer_upper)
		if is_match && valid_word {
			return theme.Green
		} else if is_match && !valid_word {
			return theme.Red
		}
	}

	return default_color
}

func (m *model) renderValidationMsg() string {
	if strings.HasPrefix(m.state.game_ui.validation_msg, "✓") {
		return styles.TextGreen.Render(utils.RightPad(m.state.game_ui.validation_msg, 2))
	}

	var msg string
	msg, _ = m.anim_mgr.ApplyAnimations(
		string(animations.ValidationMessage),
		m.state.game_ui.validation_msg)

	// Prevent input box from shaking by ensuring msg and input width are both even/odd
	raw_str := utils.StripANSICodes(msg)
	if len(raw_str) % 2 != m.text_input.CharLimit % 2 {
		if raw_str[0] == ' ' {
			msg = utils.LeftPad(msg, 1)
		} else {
			msg = utils.RightPad(msg, 1)
		}
	}

	return styles.TextRed.Render(msg)
}

// Initialize text input for use with rounded borders
func (m model) initRoundedTextInput() textinput.Model {
	text_input := textinput.New()
	text_input.Prompt = " "
	text_input.CharLimit = 40
	text_input.Width = 40

	text_input.Focus()

	return text_input
}

// Get text input with rounded border styling applied
func (m model) GetRoundedInputView() string {
	border_color := m.getInputAccentColor(theme.Border)
	input := styles.
		TextInputRoundedBorderStyle(border_color, m.text_input.CharLimit).
		Render(m.text_input.View())

	return lipgloss.JoinHorizontal(lipgloss.Center, input)
}

// Initialize text input for use with block style borders
func (m model) initBlockTextInput() textinput.Model {
	text_input := textinput.New()
	text_input.Prompt = "  "
	text_input.CharLimit = 40
	text_input.Width = text_input.CharLimit - len(text_input.Prompt) - 1

	input_bg_style := lipgloss.NewStyle().Background(theme.InputBg)

	text_input.TextStyle = input_bg_style
	text_input.PromptStyle = input_bg_style.
		Foreground(theme.Body).
		BorderLeftForeground(theme.Highlight)
	text_input.Cursor.TextStyle = input_bg_style
	text_input.PlaceholderStyle = input_bg_style.Foreground(theme.Dim)

	text_input.Focus()

	return text_input
}

// Get text input with block border styling applied
func (m model) GetBlockInputView() string {
	border_color := m.getInputAccentColor(m.text_input.PromptStyle.GetBorderLeftForeground())
	return styles.TextInputBlockBorderStyle(border_color, m.text_input.CharLimit).
		Render(m.text_input.View())
}
