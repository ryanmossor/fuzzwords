package tui

import (
	"fmt"
	"fzwds/src/enums"
	"fzwds/src/game"
	"fzwds/src/tui/animations"
	"math/rand"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m model) GameView() string {
	prompt := m.theme.TextAccent().
		Bold(true).
		Render(strings.ToUpper(m.state.game.CurrentTurn.Prompt))

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
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		prompt,
		m.GameStrikeCounterView(),
		"",
		"",
		"",
		"",
		"",
		"",
		colorized_input,
		"",
		m.GetBlockInputView(),
	)
}

func (m model) GameSwitch() (model, tea.Cmd) {
	m = m.SwitchPage(game_page)

	m.state.game_ui.game_active = true
	m.state.game_ui.validation_msg = ""

	// Reset damage animation to ensure it doesn't keep playing from previous failed turn
	m.state.game_ui.player_damaged = false

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

	extra_life_anim := animations.NewRainbowScrollAnim(animations.ExtraLife, 30, false, m.theme.GetRainbowColors())
	m.animation_manager.Register(extra_life_anim)

	validation_msg_dmg_anim := animations.NewDamageShakeAnim(animations.ValidationMessage, 8)
	m.animation_manager.Register(validation_msg_dmg_anim)

	strike_dmg_anim := animations.NewDamageShakeAnim(animations.StrikeCounter, 6)
	m.animation_manager.Register(strike_dmg_anim)

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
			m.terminalBellCmd(false),
		)

		turn_duration_min := max(m.game_settings.TurnDurationMin, 10)
		turn_duration_max := 30
		turn_time := rand.Intn(turn_duration_max - turn_duration_min + 1) + turn_duration_min 
		m.state.game_ui.timer = time.Duration(turn_time) * time.Second

		if m.state.game.Player.HealthCurrent == 0 {
			return m.GameOverSwitch(false, false)
		} else if m.state.game.CurrentTurn.Strikes == m.state.game.Settings.PromptStrikes {
			m.state.game_ui.validation_msg = m.theme.TextRed().Render(
				fmt.Sprintf(
					"Prompt %s failed. Possible solve: ",
					strings.ToUpper(m.state.game.CurrentTurn.Prompt)))
			m.state.game_ui.validation_msg += m.highlightPromptAnswer(
				m.state.game.CurrentTurn.Prompt,
				m.state.game.CurrentTurn.SourceWord,
				m.state.game.Settings.PromptMode)

			m.animation_manager.InitAnimations(animations.ValidationMessage)

			m.text_input.Reset()
			cmds = append(cmds, m.debounceInputCmd(500))

			m.state.game.NewTurn()
		} else if m.state.game.CurrentTurn.Strikes < m.state.game.Settings.PromptStrikes {
			m.state.game_ui.validation_msg = ""
			m.animation_manager.InitAnimations(animations.StrikeCounter)
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
		m.animation_manager.DeactivateAnimations(animations.ValidationMessage)

		switch key {
		case "esc":
			m.text_input.Reset()
		case "ctrl+q":
			return m.GameOverSwitch(false, true)
		case "enter":
			// TODO pass answer to ValidateAnswer
			m.state.game.CurrentTurn.Answer = strings.ToLower(strings.TrimSpace(m.text_input.Value()))
            m.text_input.Reset()
			m.state.game_ui.validation_msg = m.state.game.ValidateAnswer()

			if !m.state.game.CurrentTurn.IsValid {
				break
			}

			m.state.game.HandleCorrectAnswer()
			if len(m.state.game.Player.LettersUsed) >= len(m.state.game.Alphabet) {
				m.state.game.GrantExtraLife()
				m.animation_manager.InitAnimations(animations.ExtraLife)
			}

			// Reset damage animation to ensure it doesn't keep playing from previous failed turn
			m.animation_manager.DeactivateAnimations(animations.ValidationMessage)
			m.state.game_ui.player_damaged = false

			if len(m.state.game.WordLists.Available) == 0 {
				return m.GameOverSwitch(true, false)
			} else if (
				m.state.game.Settings.WinCondition == enums.MaxLives &&
				m.state.game.Player.HealthCurrent == m.state.game.Settings.HealthMax) {
				return m.GameOverSwitch(true, false)
			}

			m.state.game.NewTurn()

			if m.state.game_ui.timer < time.Duration(m.game_settings.TurnDurationMin) * time.Second {
				m.state.game_ui.timer = time.Duration(m.game_settings.TurnDurationMin) * time.Second
			}

			cmds = append(cmds, m.debounceInputCmd(300))
			return m, tea.Batch(cmds...)
		}
	}

	var update_input_cmd tea.Cmd
	m.text_input, update_input_cmd = m.text_input.Update(msg)
	cmds = append(cmds, update_input_cmd)

	return m, tea.Batch(cmds...)
}
