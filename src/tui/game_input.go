package tui

import (
	"fmt"
	"fzwds/src/enums"
	"fzwds/src/game"
	"fzwds/src/tui/animations"
	"fzwds/src/utils"
	"math/rand"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m model) GameSwitch() (model, tea.Cmd) {
	m = m.SwitchPage(game_page)

	m.state.game_ui.game_active = true
	m.state.game_ui.validation_msg = ""

	// Reset damage animation to ensure it doesn't keep playing from previous failed turn
	m.state.game_ui.player_damaged = false
	m.state.game_ui.damage_anim_padding = 0

	m.state.game = game.InitializeGame(m.game_settings)
	m.state.game.NewTurn()

	m.state.game_ui.start_time = time.Now()
    m.state.game_ui.timer = (30 + 1) * time.Second

	m.footer_keymaps = []footer_keymaps{
		{key: "esc", value: "clear input"},
		{key: "ctrl+q", value: "quit"},
	}

	m.text_input = m.initBlockTextInput()
	m.state.game_ui.input_restricted = false

	extra_life_anim := &animations.RainbowScrollAnim {
		BaseAnim: animations.BaseAnim {
			FrameInterval:	time.Second / 30,
			PrevFrame:		time.Now(),
			Frame:			0,
			Loop:			true,
			Active:			false,
			Target:			animations.ExtraLife,
		},
		Offset: 			0,
		TotalFrames: 		30,
		Colors: 			m.theme.GetRainbowColors(),
	}
	m.animation_manager.Register("extra_life", extra_life_anim)

	return m, tea.Batch(
		textinput.Blink,
		m.setTurnTickerCmd(),
	)
}

func (m model) GameUpdate(msg tea.Msg) (model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
    case TurnTimerTickMsg:
		if m.state.game_ui.timer > 0 {
			return m, m.setTurnTickerCmd()
		}

		m.state.game.HandleFailedTurn()
		cmds = append(cmds,
			m.setPlayerDamagedStateCmd(),
			m.damageShakeAnimationCmd(8),
			m.terminalBellCmd(false),
		)

		turn_duration_min := max(m.game_settings.TurnDurationMin, 10)
		turn_duration_max := 30
		turn_time := rand.Intn(turn_duration_max - turn_duration_min + 1) + turn_duration_min 
		m.state.game_ui.timer = time.Duration(turn_time) * time.Second

		if m.state.game.Player.HealthCurrent == 0 {
			return m.GameOverSwitch(false)
		} else if m.state.game.CurrentTurn.Strikes == m.state.game.Settings.PromptStrikes {
			m.state.game_ui.validation_msg = m.theme.TextRed().Render(
				// TODO: turn handler already sets msg to Possible solve:...
				// If msg contains "Possible solve", split on space and colorize final word
				fmt.Sprintf(
					"Prompt %s failed. Possible solve: ",
					strings.ToUpper(m.state.game.CurrentTurn.Prompt)))
			m.state.game_ui.validation_msg += m.highlightPromptAnswer(
				m.state.game.CurrentTurn.Prompt,
				m.state.game.CurrentTurn.SourceWord,
				m.state.game.Settings.PromptMode)

			m.text_input.Reset()
			cmds = append(cmds, m.debounceInputCmd(500))

			m.state.game.NewTurn()
		} else if m.state.game.CurrentTurn.Strikes < m.state.game.Settings.PromptStrikes {
			m.state.game_ui.validation_msg = ""
		}

        cmds = append(cmds, m.setTurnTickerCmd())
        return m, tea.Batch(cmds...)
	case tea.KeyMsg:
        if m.state.game_ui.input_restricted {
            return m, nil
        }

        key := msg.String()
		if key != "enter" {
			m.state.game_ui.validation_msg = ""
		}

		switch key {
		case "esc":
			m.text_input.Reset()
		case "ctrl+q":
			return m.GameOverSwitch(false)
		case "enter":
			m.state.game.CurrentTurn.Answer = strings.ToLower(strings.TrimSpace(m.text_input.Value()))
            m.text_input.Reset()
			m.state.game_ui.validation_msg = m.state.game.ValidateAnswer()

			if !m.state.game.CurrentTurn.IsValid {
				break
			}

			m.state.game.HandleCorrectAnswer()
			if len(m.state.game.Player.LettersUsed) >= len(m.state.game.Alphabet) {
				m.state.game.GrantExtraLife()
				m.animation_manager.InitAnimations("extra_life")
			}

			// Reset damage animation to ensure it doesn't keep playing from previous failed turn
			m.state.game_ui.player_damaged = false
			m.state.game_ui.damage_anim_padding = 0

			// TODO: move win condition check to game_over?
			if len(m.state.game.WordLists.Available) == 0 {
				return m.GameOverSwitch(true)
			} else if (
				m.state.game.Settings.WinCondition == enums.MaxLives &&
				m.state.game.Player.HealthCurrent == m.state.game.Settings.HealthMax) {
				return m.GameOverSwitch(true)
			}

			m.state.game.NewTurn()

			if m.state.game_ui.timer < time.Duration(m.game_settings.TurnDurationMin) * time.Second {
				m.state.game_ui.timer = time.Duration(m.game_settings.TurnDurationMin) * time.Second
			}

			cmds = append(cmds, m.debounceInputCmd(300))
			return m, tea.Batch(cmds...)
		}
	case DamageShakeAnimationMsg:
		if m.state.game_ui.damage_anim_padding > 0 {
			m.state.game_ui.damage_anim_padding -= 2
			return m, tea.Tick(time.Second / time.Duration(m.FPS), func(t time.Time) tea.Msg {
				return DamageShakeAnimationMsg{}
			})
		}
	}

	var update_input_cmd tea.Cmd
	m.text_input, update_input_cmd = m.text_input.Update(msg)
	cmds = append(cmds, update_input_cmd)

	return m, tea.Batch(cmds...)
}

func (m model) GameInputView() string {
	if !m.state.game_ui.game_active {
		return ""
	}

	var colorized_input string
	if m.state.game_ui.validation_msg != "" {
		colorized_input = m.renderValidationMsg()
	} else {
		colorized_input = m.highlightPromptAnswer(
			m.state.game.CurrentTurn.Prompt,
			m.text_input.Value(),
			m.state.game.Settings.PromptMode)
	}

	return lipgloss.JoinVertical(
		lipgloss.Center,
		colorized_input,
		"",
		m.getStyledBlockTextInput(),
		"",
	) 
}

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
		msg, _ = m.applyDamageShakeAnimation(m.state.game_ui.validation_msg)
		// Add padding to msg if necessary to prevent input box from shaking
		if len(utils.StripANSICodes(msg)) % 2 == 1 {
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
func (m model) getStyledRoundedTextInput() string {
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
func (m model) getStyledBlockTextInput() string {
	border_color := m.getInputAccentColor(m.text_input.PromptStyle.GetBorderLeftForeground())
	return m.TextInputBlockBorderStyle(border_color).
		Render(m.text_input.View())
}
