package tui

import (
	"fzwds/src/game"
	"fzwds/src/tui/animations"
	"fzwds/src/utils"
	"strings"

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

	m.footer_keymaps = []FooterKeymap{
		{key: "esc", value: "clear input"},
		{key: "ctrl+q", value: "quit"},
	}

	// Reset damage animation to ensure it doesn't keep playing from previous failed turn
	m.state.game_ui.player_damaged = false
	m.state.game_ui.input_restricted = false
	m.state.game_ui.prev_answer = ""
	m.state.game_ui.validation_msg = ""

	m.state.game = game.InitializeGame(m.game_settings)
	m.state.game.StartGame()
	m.state.game.NewTurn(true)

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

		if m.state.game.IsGameOver() {
			return m.GameOverSwitch(false, false)
		}

		turn_failure_msg := m.state.game.GetTurnFailureMessage()
		if turn_failure_msg == "" {
			m.anim_mgr.InitAnimations(animations.StrikeCounter)
			m.state.game.StartTurn(utils.RandomBetween(m.game_settings.TurnDurationMin, 30))
		} else {
			// Strike limit reached, show failure message and start new turn
			possible_solve := m.highlightPromptAnswer(
				m.state.game.CurrentTurn.Prompt,
				m.state.game.CurrentTurn.PossibleAnswer,
				m.state.game.Settings.PromptMode)
			updated_msg := strings.ReplaceAll(turn_failure_msg, "{solve}", possible_solve)
			turn_failure_msg = m.theme.TextRed().Render(updated_msg)

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
			return m.GameOverSwitch(false, true)

		case "enter":
			answer := strings.ToLower(strings.TrimSpace(m.text_input.Value()))
            m.text_input.Reset()

			var is_valid bool
			is_valid, m.state.game_ui.validation_msg = m.state.game.ValidateAnswer(answer)
			if !is_valid {
				m.state.game_ui.prev_answer = answer
				break
			} else {
				m.state.game_ui.prev_answer = ""
			}

			result := m.state.game.HandleCorrectAnswer(answer)
			if result.ExtraLifeGranted {
				m.state.game.GrantExtraLife()
				m.anim_mgr.InitAnimations(animations.ExtraLife)
			}

			// Reset damage animation to ensure it doesn't keep playing from previous failed turn
			m.anim_mgr.DeactivateAnimations(animations.ValidationMessage)
			m.state.game_ui.player_damaged = false

			// Check if win condition met (no more available words, max lives)
			if m.state.game.IsGameOver() {
				return m.GameOverSwitch(true, false)
			}

			m.state.game.NewTurn(false)

			return m, m.debounceInputCmd(300)
		}
	}

	var update_input_cmd tea.Cmd
	m.text_input, update_input_cmd = m.text_input.Update(msg)

	return m, update_input_cmd
}
