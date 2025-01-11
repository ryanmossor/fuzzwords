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

	// TODO: move these to game state?
	m.game_active = true
	m.game_over = false
	m.state.game.validation_msg = ""

	m.game_state = game.InitializeGame(m.game_settings)
	m.game_state.NewTurn()

	m.game_start_time = time.Now()

    m.game_timer.remaining_time = 30 * time.Second

	m.footer_cmds = []footerCmd{
		{key: "esc", value: "clear input"},
		{key: "ctrl+q", value: "quit"},
	}

	m.state.game.restrict_input = false
	m.text_input.Reset()

	cmds := []tea.Cmd{textinput.Blink, setTurnTickerCmd()}

	return m, tea.Batch(cmds...)
}

func (m model) GameUpdate(msg tea.Msg) (model, tea.Cmd) {
	red := m.theme.TextRed().Bold(true).Render
	green := m.theme.TextGreen().Bold(true).Render

	switch msg := msg.(type) {
    case TurnTimerTickMsg:
		if m.game_timer.remaining_time <= 0 {
            m.game_state.HandleFailedTurn()

            turn_duration_min := max(m.game_settings.TurnDurationMin, 10)
            turn_duration_max := 30
            turn_time := rand.Intn(turn_duration_max - turn_duration_min + 1) + turn_duration_min 
            m.game_timer.remaining_time = time.Duration(turn_time) * time.Second

            if m.game_state.Player.HealthCurrent == 0 {
                return m.GameOverSwitch(red(game_over_msg))
            } else {
                m.state.game.validation_msg = fmt.Sprintf(
                    "Prompt %s failed. Possible answer: %s",
                    strings.ToUpper(m.game_state.CurrentTurn.Prompt),
                    strings.ToUpper(m.game_state.CurrentTurn.SourceWord))

                m.game_state.NewTurn()
            }

            m.text_input.Reset()
		}

		m.game_timer.remaining_time -= time.Millisecond * 100

		return m, setTurnTickerCmd()
	case tea.KeyMsg:
		// clear validation msg while debouncing enter presses (yes this is kinda scuffed)
		if msg.String() != "enter" {
			m.state.game.validation_msg = ""
		}

		switch msg.String() {
		case "esc":
			m.text_input.Reset()
		case "ctrl+q":
			return m.GameOverSwitch(red(game_over_msg))
		case "enter":
			if m.state.game.restrict_input {
				return m, nil
			}

			m.state.game.restrict_input = true

			// TODO: trim answer & take only first word before any spaces/symbols
			m.game_state.CurrentTurn.Answer = strings.ToLower(m.text_input.Value())
			m.state.game.validation_msg = m.game_state.ValidateAnswer()

			if m.game_state.CurrentTurn.IsValid {
				m.game_state.HandleCorrectAnswer()

				if len(m.game_state.WordLists.Available) == 0 {
					return m.GameOverSwitch(green(win_msg))
				}

				m.game_state.NewTurn()

                if m.game_timer.remaining_time < time.Duration(m.game_settings.TurnDurationMin) * time.Second {
                    m.game_timer.remaining_time = time.Duration(m.game_settings.TurnDurationMin) * time.Second
                }
			}

			if (m.game_state.Settings.WinCondition == enums.Debug && m.game_state.Player.Stats.PromptsSolved == 10) {
				return m.GameOverSwitch(green("stop stalling and do some work"))
			}

			if (m.game_state.Settings.WinCondition == enums.MaxLives && m.game_state.Player.HealthCurrent == m.game_state.Settings.HealthMax) {
				return m.GameOverSwitch(green(win_msg))
			}

			m.text_input.Reset()

            return m, debounceInputCmd(300)
		}
	}

	var cmd tea.Cmd
	m.text_input, cmd = m.text_input.Update(msg)

	return m, cmd
}

func (m model) GameInputView() string {
	if !m.game_active {
		return ""
	}

	var colorized_text string
	var border_color lipgloss.TerminalColor

	if m.state.game.validation_msg != "" {
		colorized_text, border_color = m.renderValidationMsg()
	} else {
		colorized_text, border_color = m.renderColorizedInput()
	}

	return lipgloss.JoinVertical(
		lipgloss.Center,
		colorized_text,
		"\n",
		lipgloss.NewStyle().
			BorderForeground(border_color).
			BorderStyle(lipgloss.RoundedBorder()).
			Width(50).
			Render(m.text_input.View()),
	) 
}

func (m model) setInputBorderColor(answer string) lipgloss.TerminalColor {
	if m.game_state.WordLists.FULL_MAP[strings.ToLower(answer)] {
		return m.theme.green
	}
	return m.theme.red
}

func (m model) renderColorizedInput() (string, lipgloss.TerminalColor) {
	border_color := m.theme.Border()
	accent := m.theme.TextAccent().Render
	blue := m.theme.TextBlue().Render

	prompt_upper := strings.ToUpper(m.game_state.CurrentTurn.Prompt)
	answer_upper := strings.ToUpper(m.text_input.Value())
	var sb strings.Builder
	 
	switch m.game_state.Settings.PromptMode {
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

		if m.game_state.Settings.HighlightInput && utils.IsFuzzyMatch(answer_upper, prompt_upper) {
			border_color = m.setInputBorderColor(answer_upper)
		}
	case enums.Classic:
		if !strings.Contains(answer_upper, prompt_upper) {
			sb.WriteString(accent(answer_upper))
			return sb.String(), border_color
		}
		
		sub_idx := strings.Index(answer_upper, prompt_upper)
		sb.WriteString(accent(answer_upper[0:sub_idx]))
		sb.WriteString(blue(answer_upper[sub_idx:sub_idx + len(prompt_upper)]))
		sb.WriteString(accent(answer_upper[sub_idx + len(prompt_upper):]))

		if m.game_state.Settings.HighlightInput {
			border_color = m.setInputBorderColor(answer_upper)
		}
	}

	return sb.String(), border_color
}

func (m *model) renderValidationMsg() (string, lipgloss.TerminalColor) {
	border_color := m.theme.Border()

	if strings.Contains(m.state.game.validation_msg, "Correct") {
		return m.theme.TextGreen().Bold(true).Render(m.state.game.validation_msg), border_color
	}

    m.text_input.PromptStyle = m.theme.TextRed()
    border_color = m.theme.red

	return m.theme.TextRed().Render(m.state.game.validation_msg), border_color
}
