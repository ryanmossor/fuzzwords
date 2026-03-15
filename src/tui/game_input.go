package tui

import (
	"fzwds/src/enums"
	"fzwds/src/tui/animations"
	"fzwds/src/utils"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

func (m model) wordInDictionary(answer string) bool {
	return m.state.game.WordLists.FULL_MAP[strings.ToLower(answer)]
}

// Highlight prompt letters in current answer
func (m model) highlightPromptAnswer(prompt, answer string, prompt_mode enums.PromptMode) string {
	accent := m.theme.TextAccent().Render
	highlight := m.theme.TextHighlight().Render

	prompt_upper := strings.ToUpper(prompt)
	answer_upper := strings.ToUpper(answer)
	var sb strings.Builder
	 
	switch prompt_mode {
	case enums.Fuzzy:
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
	case enums.Classic:
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
	prompt_upper := strings.ToUpper(m.state.game.CurrentTurn.Prompt)
	answer_upper := strings.ToUpper(m.text_input.Value())

	is_match := false
	switch m.state.game.Settings.PromptMode {
	case enums.Fuzzy:
		is_match = utils.IsFuzzyMatch(answer_upper, prompt_upper)
	case enums.Classic:
		is_match = strings.Contains(answer_upper, prompt_upper)
	}

	if m.state.game.Settings.HighlightInput {
		valid_word := m.wordInDictionary(answer_upper)
		if is_match && valid_word {
			return m.theme.green
		} else if is_match && !valid_word {
			return m.theme.red
		}
	} 

	if m.state.game_ui.player_damaged {
		return m.theme.red
	}

	return default_color
}

func (m *model) renderValidationMsg() string {
	// TODO: check if m.state.game.CurrentTurn.IsValid?
	if strings.HasPrefix(m.state.game_ui.validation_msg, "✓") {
		return m.theme.TextGreen().Render(utils.RightPad(m.state.game_ui.validation_msg, 2))
	}

	var msg string
	if !m.state.game_ui.game_active {
		msg = m.state.game_ui.validation_msg
	} else {
		msg, _ = m.animation_manager.ApplyAnimations(
			string(animations.ValidationMessage),
			m.state.game_ui.validation_msg)

		// Input box will shake if msg and input width are not both even/odd
		if len(utils.StripANSICodes(msg)) % 2 != m.text_input.CharLimit % 2 {
			msg = utils.RightPad(msg, 1)
		}
	}

	return m.theme.TextRed().Render(msg)
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
	border_color := m.getInputAccentColor(m.theme.Border())
	input := m.TextInputRoundedBorderStyle(border_color).Render(m.text_input.View())

	return lipgloss.JoinHorizontal(lipgloss.Center, input)
}

// Initialize text input for use with block style borders
func (m model) initBlockTextInput() textinput.Model {
	text_input := textinput.New()
	text_input.Prompt = "  "
	text_input.CharLimit = 40
	text_input.Width = text_input.CharLimit - len(text_input.Prompt) - 1

	input_bg_style := lipgloss.NewStyle().Background(m.theme.input_bg)

	text_input.TextStyle = input_bg_style
	text_input.PromptStyle = input_bg_style.
		Foreground(m.theme.body).
		BorderLeftForeground(m.theme.blue)
	text_input.Cursor.TextStyle = input_bg_style
	text_input.PlaceholderStyle = input_bg_style.Foreground(m.theme.dim)

	text_input.Focus()

	return text_input
}

// Get text input with block border styling applied
func (m model) GetBlockInputView() string {
	border_color := m.getInputAccentColor(m.text_input.PromptStyle.GetBorderLeftForeground())
	return m.TextInputBlockBorderStyle(border_color).
		Render(m.text_input.View())
}
