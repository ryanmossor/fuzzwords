package tui

import (
	"fmt"
	"fzwds/src/enums"
	"fzwds/src/game"
	"fzwds/src/utils"
	"math/rand"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var win_msg string = "===== YOU WIN! ====="
var game_over_msg string = "===== GAME OVER ====="

func (m model) GameSwitch() (model, tea.Cmd) {
	m = m.SwitchPage(game_page)

	m.state.game_ui.game_active = true
	m.state.game_ui.validation_msg = ""
	m.state.game_ui.player_damaged = false

	m.state.game = game.InitializeGame(m.game_settings)
	m.state.game.NewTurn()

	m.state.game_ui.start_time = time.Now()
    m.state.game_ui.timer = (30 + 1) * time.Second

	m.footer_keymaps = []footer_keymaps{
		{key: "esc", value: "clear input"},
		{key: "ctrl+q", value: "quit"},
	}

	m.state.game_ui.input_restricted = false
	m.text_input.Reset()

	cmds := []tea.Cmd{textinput.Blink, m.setTurnTickerCmd()}

	return m, tea.Batch(cmds...)
}

func (m model) GameUpdate(msg tea.Msg) (model, tea.Cmd) {
	red_bold := m.theme.TextRed().Bold(true).Render
	red := m.theme.TextRed().Render
	green := m.theme.TextGreen().Bold(true).Render

	switch msg := msg.(type) {
    case TurnTimerTickMsg:
        cmds := []tea.Cmd{}

		if m.state.game_ui.timer <= 0 {
            m.state.game.HandleFailedTurn()
			cmds = append(cmds, m.setPlayerDamagedStateCmd(), m.damageShakeAnimationCmd(8))

            turn_duration_min := max(m.game_settings.TurnDurationMin, 10)
            turn_duration_max := 30
            turn_time := rand.Intn(turn_duration_max - turn_duration_min + 1) + turn_duration_min 
            m.state.game_ui.timer = time.Duration(turn_time) * time.Second

            if m.state.game.Player.HealthCurrent == 0 {
                return m.GameOverSwitch(red_bold(game_over_msg), false)
			} else if m.state.game.CurrentTurn.Strikes == m.state.game.Settings.PromptStrikesMax {
				m.state.game_ui.validation_msg = red(
					fmt.Sprintf(
						"Prompt %s failed. Possible solve: ",
						strings.ToUpper(m.state.game.CurrentTurn.Prompt)))
				m.state.game_ui.validation_msg += m.colorizeInput(m.state.game.CurrentTurn.SourceWord)

				m.text_input.Reset()
				cmds = append(cmds, m.debounceInputCmd(500))

                m.state.game.NewTurn()
            } else if m.state.game.CurrentTurn.Strikes < m.state.game.Settings.PromptStrikesMax {
                m.state.game_ui.validation_msg = ""
			}
		}

        cmds = append(cmds, m.setTurnTickerCmd())
        return m, tea.Batch(cmds...)
	case tea.KeyMsg:
        if m.state.game_ui.input_restricted {
            return m, nil
        }

		var cmds []tea.Cmd

        key := msg.String()
		if key != "enter" {
			m.state.game_ui.validation_msg = ""
		}

		switch key {
		case "esc":
			m.text_input.Reset()
		case "ctrl+q":
			return m.GameOverSwitch(red_bold(game_over_msg), false)
		case "enter":
			m.state.game.CurrentTurn.Answer = strings.ToLower(strings.TrimSpace(m.text_input.Value()))
            m.text_input.Reset()
			m.state.game_ui.validation_msg = m.state.game.ValidateAnswer()

			if m.state.game.CurrentTurn.IsValid {
				m.state.game.HandleCorrectAnswer()
				if len(m.state.game.Player.LettersUsed) >= len(m.state.game.Alphabet) {
					m.state.game.GrantExtraLife()
					cmds = append(cmds, m.extraLifeAnimInitMsg())
				}

				// Reset damage animation to ensure it doesn't keep playing from previous failed turn
				m.state.game_ui.player_damaged = false
				m.state.game_ui.damage_anim_padding = 0

				// TODO: move win condition check to game_over?
				if len(m.state.game.WordLists.Available) == 0 {
					return m.GameOverSwitch(green(win_msg), true)
				} else if (m.state.game.Settings.WinCondition == enums.Debug && m.state.game.Player.Stats.PromptsSolved == 10) {
                    return m.GameOverSwitch(green("stop stalling and do some work"), true)
                } else if (m.state.game.Settings.WinCondition == enums.MaxLives && m.state.game.Player.HealthCurrent == m.state.game.Settings.HealthMax) {
                    return m.GameOverSwitch(green(win_msg), true)
                }

				m.state.game.NewTurn()

				if m.state.game_ui.timer < time.Duration(m.game_settings.TurnDurationMin) * time.Second {
					m.state.game_ui.timer = time.Duration(m.game_settings.TurnDurationMin) * time.Second
				}

				cmds = append(cmds, m.debounceInputCmd(300))
				return m, tea.Batch(cmds...)
            }
		}
	case DamageShakeAnimationMsg:
		if m.state.game_ui.damage_anim_padding > 0 {
			m.state.game_ui.damage_anim_padding -= 2
			return m, tea.Tick(time.Second / time.Duration(m.anim_fps), func(t time.Time) tea.Msg {
				return DamageShakeAnimationMsg{}
			})
		}
	}

	var cmd tea.Cmd
	m.text_input, cmd = m.text_input.Update(msg)

	return m, cmd
}

func (m model) GameInputView() string {
	if !m.state.game_ui.game_active {
		return ""
	}

	var colorized_input string
	var border_color lipgloss.TerminalColor

	if m.state.game_ui.validation_msg != "" {
		colorized_input = m.renderValidationMsg()
		border_color = m.getInputBorderColor()

		if m.state.game_ui.player_damaged {
			m.text_input.PromptStyle = m.theme.TextRed()
		} else {
			m.text_input.Reset()
		}
	} else {
		colorized_input = m.colorizeInput(m.text_input.Value())
		border_color = m.getInputBorderColor()
	}

	return lipgloss.JoinVertical(
		lipgloss.Center,
		colorized_input,
		"\n",
		lipgloss.NewStyle().
			BorderForeground(border_color).
			BorderStyle(lipgloss.RoundedBorder()).
			Width(50).
			Render(m.text_input.View()),
	) 
}

func (m model) wordInDictionary(answer string) bool {
	return m.state.game.WordLists.FULL_MAP[strings.ToLower(answer)]
}

func (m model) colorizeInput(answer string) string {
	accent := m.theme.TextAccent().Render
	highlight := m.theme.TextHighlight().Render

	prompt_upper := strings.ToUpper(m.state.game.CurrentTurn.Prompt)
	answer_upper := strings.ToUpper(answer)
	var sb strings.Builder
	 
	switch m.state.game.Settings.PromptMode {
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

func (m model) getInputBorderColor() lipgloss.TerminalColor {
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

	return m.theme.Border()
}

func (m *model) renderValidationMsg() string {
	if strings.HasPrefix(m.state.game_ui.validation_msg, "✓") {
		return m.theme.TextGreen().Render(strings.TrimSpace(m.state.game_ui.validation_msg))
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
