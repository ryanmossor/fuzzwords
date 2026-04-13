package tui

import (
	"fzwds/src/game"
	"fzwds/src/tui/animations"
	"fzwds/src/tui/styles"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m model) GameView() string {
	prompt := styles.TextAccent.
		Bold(true).
		Render(strings.ToUpper(m.state.game.CurrentTurn().Prompt))

	var colorized_input string
	if m.state.game_ui.validation_msg != "" {
		colorized_input = m.renderValidationMsg()
	} else {
		colorized_input = m.highlightPromptAnswer(
			m.state.game.CurrentTurn().Prompt,
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

	m.footer_keymaps = []FooterKeymap {
		{key: "esc", value: "clear input"},
		{key: "ctrl+q", value: "quit"},
	}

	// Reset damage animation to ensure it doesn't keep playing from previous failed turn
	m.state.game_ui.player_damaged = false
	m.state.game_ui.input_restricted = false
	m.state.game_ui.prev_answer = ""
	m.state.game_ui.validation_msg = ""
	m.state.game_ui.game_over_seen = false

	m.state.game_review.selected_turn = 0
	m.state.game_review.visible_row_start = 0
	m.state.game_review.view_cache = make(map[int]*TurnDisplay, 0)
	m.state.game = game.InitializeGame(&m.app_settings.Game)

	m.text_input = m.initBlockTextInput()
	return m, textinput.Blink
}

func (m model) GameUpdate(msg tea.Msg) (model, tea.Cmd) {

	switch msg := msg.(type) {
	case TurnTimerExpiredMsg:
		var cmds []tea.Cmd
		m.state.game.HandleFailedTurn()
		cmds = append(cmds,
			m.setPlayerDamagedStateCmd(),
			m.terminalBellCmd(false),
		)

        if m.state.game.EndGameIfOver() {
			return m.GameOverSwitch()
		}

		turn_failure_msg := m.state.game.GetTurnFailureMessage()
		if turn_failure_msg == "" {
			m.anim_mgr.InitAnimations(animations.StrikeCounter)
			m.state.game.StartTurn()
		} else {
			m.anim_mgr.InitAnimations(animations.ValidationMessage)

			m.text_input.Reset()
			cmds = append(cmds, m.debounceInputCmd(500))

			m.state.game.NewTurn(false)
		}
		m.state.game_ui.validation_msg = turn_failure_msg

        return m, tea.Batch(cmds...)

	case tea.KeyMsg:
        if m.state.game_ui.input_restricted {
            return m, nil
        }

        key := msg.String()
		if key != "enter" {
			m.state.game_ui.validation_msg = ""
		}
		m.anim_mgr.DeactivateAnimations(animations.ValidationMessage)

		// TODO: skips -- sacrifice life for skip, earn extra skips from extra lifes if already full
		// TODO: mute key combo for alert sound in game?
		switch key {
		case "up":
			if m.state.game_ui.prev_answer != "" {
				m.text_input.SetValue(m.state.game_ui.prev_answer)
			}
			return m, nil

		case "esc":
			m.text_input.Reset()
			return m, nil

		case "ctrl+q":
			m.state.game.EndGame(false, true)
			return m.GameOverSwitch()

		case "enter":
			answer := strings.ToLower(strings.TrimSpace(m.text_input.Value()))
            m.text_input.Reset()

			result := m.state.game.SubmitAnswer(answer)
			m.state.game_ui.validation_msg = result.ValidationMsg
			if !result.IsValid {
				m.state.game_ui.prev_answer = answer
				break
			}
			m.state.game_ui.prev_answer = ""

			if result.ExtraLifeGained {
				m.anim_mgr.InitAnimations(animations.ExtraLife)
			}

			// Reset damage animation to ensure it doesn't keep playing from previous failed turn
			m.anim_mgr.DeactivateAnimations(animations.ValidationMessage)
			m.state.game_ui.player_damaged = false

			if m.state.game.EndGameIfOver() {
				return m.GameOverSwitch()
			}

			m.state.game.NewTurn(false)

			return m, m.debounceInputCmd(300)
		}
	}

	var update_input_cmd tea.Cmd
	m.text_input, update_input_cmd = m.text_input.Update(msg)

	return m, update_input_cmd
}
