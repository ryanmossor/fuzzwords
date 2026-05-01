package pages

// import (
// 	"fzwds/pkg/enums"
// 	"fzwds/pkg/tui/animations"
// 	"fzwds/pkg/tui/styles"
// 	"fzwds/pkg/tui/theme"
// 	"fzwds/pkg/utils"
// 	"strings"
//
// 	"github.com/charmbracelet/bubbles/cursor"
// 	"github.com/charmbracelet/bubbles/textinput"
// 	"github.com/charmbracelet/lipgloss"
// )

// // Get accent color for input box based on HighlightInput setting, damage state, etc.
// // Style applied to border if rounded-style input box, or left accent bar if block-style input box.
// func (m model) getInputAccentColor(default_color lipgloss.TerminalColor) lipgloss.TerminalColor {
// 	if m.state.game.playerDamaged {
// 		return theme.Red
// 	}
//
// 	prompt_upper := strings.ToUpper(m.state.game.turn.prompt)
// 	answer_upper := strings.ToUpper(m.gameInput.Value())
//
// 	is_match := false
// 	switch m.game.Settings().PromptMode {
// 	case enums.PromptModeFuzzy:
// 		is_match = utils.IsFuzzyMatch(answer_upper, prompt_upper)
// 	case enums.PromptModeClassic:
// 		is_match = strings.Contains(answer_upper, prompt_upper)
// 	}
//
// 	if m.game.Settings().HighlightInput {
// 		valid_word := m.game.WordInDictionary(answer_upper)
// 		if is_match && valid_word {
// 			return theme.Green
// 		} else if is_match && !valid_word {
// 			return theme.Red
// 		}
// 	}
//
// 	return default_color
// }
//
// func (m *model) renderValidationMsg() string {
// 	msg, changed := m.animManager.ApplyAnimations(string(animations.ValidationMessage), m.state.game.gameMsg)
// 	if !changed {
// 		return msg
// 	}
//
// 	// Prevent input box from shaking by ensuring msg and input width are both even/odd
// 	raw_str := utils.StripANSICodes(msg)
// 	if len(raw_str) % 2 != m.gameInput.CharLimit % 2 {
// 		if raw_str[0] == ' ' {
// 			msg = utils.LeftPad(msg, 1)
// 		} else {
// 			msg = utils.RightPad(msg, 1)
// 		}
// 	}
//
// 	return msg
// }
//
// // Initialize text input for use with rounded borders
// func (m model) initRoundedTextInput() textinput.Model {
// 	text_input := textinput.New()
// 	text_input.Prompt = " "
// 	text_input.CharLimit = 40
// 	text_input.Width = 40
//
// 	if !m.settings.Prefs.AnimationsEnabled {
// 		text_input.Cursor.SetMode(cursor.CursorStatic)
// 	}
//
// 	text_input.Focus()
//
// 	return text_input
// }
//
// // Get text input with rounded border styling applied
// func (m model) GetRoundedInputView() string {
// 	border_color := m.getInputAccentColor(theme.Border)
// 	input := styles.
// 		TextInputRoundedBorderStyle(border_color, m.gameInput.CharLimit).
// 		Render(m.gameInput.View())
//
// 	return lipgloss.JoinHorizontal(lipgloss.Center, input)
// }
//
// // Initialize text input for use with block style borders
// func (m model) initBlockTextInput() textinput.Model {
// 	text_input := textinput.New()
// 	text_input.Prompt = "  "
// 	text_input.CharLimit = 40
// 	text_input.Width = text_input.CharLimit - len(text_input.Prompt) - 1
//
// 	if !m.settings.Prefs.AnimationsEnabled {
// 		text_input.Cursor.SetMode(cursor.CursorStatic)
// 	}
//
// 	input_bg_style := lipgloss.NewStyle().Background(theme.InputBg)
//
// 	text_input.TextStyle = input_bg_style
// 	text_input.PromptStyle = input_bg_style.
// 		Foreground(theme.Body).
// 		BorderLeftForeground(theme.Highlight)
// 	text_input.Cursor.TextStyle = input_bg_style
// 	text_input.PlaceholderStyle = input_bg_style.Foreground(theme.Dim)
//
// 	text_input.Focus()
//
// 	return text_input
// }
//
// // Get text input with block border styling applied
// func (m model) GetBlockInputView() string {
// 	border_color := m.getInputAccentColor(m.gameInput.PromptStyle.GetBorderLeftForeground())
// 	return styles.TextInputBlockBorderStyle(border_color, m.gameInput.CharLimit).
// 		Render(m.gameInput.View())
// }
