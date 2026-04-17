package tui

import (
	"fmt"
	"fzwds/src/game"
	"fzwds/src/tui/animations"
	"fzwds/src/tui/styles"
	"log/slog"
	"reflect"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m model) GameView() string {
	prompt := styles.TextAccent.
		Bold(true).
		Render(strings.ToUpper(m.state.game.turn.prompt))

	var game_msg string
	if m.state.game.gameMsg != "" {
		game_msg = m.renderValidationMsg()
	} else {
		game_msg = m.highlightPromptAnswer(
			m.state.game.turn.prompt,
			m.text_input.Value(),
			m.game.Settings.PromptMode)
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
		game_msg,
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

	m.state.game = GameUIState{}
	m.state.gameOver = GameOverState{ viewCache: make(map[string]string) }
	m.state.gameReview = GameReviewState{ viewCache: make(map[int]*TurnDisplay) }

	var events []game.GameEvent
	m.game, events = game.NewGame(&m.app_settings.Game)

	var cmds []tea.Cmd
	for _, e := range events {
		cmds = append(cmds, m.handleGameEvent(e)...)
	}

	m.state.game.health = uint(m.game.Settings.HealthInitial)

	m.text_input = m.initBlockTextInput()
	cmds = append(cmds, textinput.Blink)

	return m, tea.Batch(cmds...)
}

func (m model) GameUpdate(msg tea.Msg) (model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case TickMsg:
		events := m.game.AdvanceTime(msg.Time)
		if len(events) == 0 {
			return m, nil
		}

		for _, e := range events {
			cmds = append(cmds, m.handleGameEvent(e)...)
		}

        if m.state.game.gameOver {
			var cmd tea.Cmd
			m, cmd = m.GameOverSwitch()
			cmds = append(cmds, cmd)
		}

        return m, tea.Batch(cmds...)

	case tea.KeyMsg:
        if m.state.game.inputRestricted {
            return m, nil
        }

        key := msg.String()
		if key != "enter" {
			m.state.game.gameMsg = ""
		}
		m.anim_mgr.DeactivateAnimations(animations.ValidationMessage)

		// TODO: skips -- sacrifice life for skip, earn extra skips from extra lifes if already full
		// TODO: mute key combo for alert sound in game?
		switch key {
		case "up":
			if m.state.game.turn.prevAnswer != "" {
				m.text_input.SetValue(m.state.game.turn.prevAnswer)
			}
			return m, nil

		case "esc":
			m.text_input.Reset()
			return m, nil

		case "ctrl+q":
			events := m.game.QuitGame()
			for _, e := range events {
				cmds = append(cmds, m.handleGameEvent(e)...)
			}

			m, cmd := m.GameOverSwitch()
			cmds = append(cmds, cmd)

			return m, tea.Batch(cmds...)

		case "enter":
			answer := strings.ToLower(strings.TrimSpace(m.text_input.Value()))
            m.text_input.Reset()

			events := m.game.SubmitAnswer(answer)
			for _, e := range events {
				cmds = append(cmds, m.handleGameEvent(e)...)
			}

			cmds = append(cmds, m.debounceInputCmd(300))

			if m.state.game.gameOver {
				var cmd tea.Cmd
				m, cmd = m.GameOverSwitch()
				cmds = append(cmds, cmd)
			}

			return m, tea.Batch(cmds...)
		}
	}

	var update_input_cmd tea.Cmd
	m.text_input, update_input_cmd = m.text_input.Update(msg)

	return m, update_input_cmd
}

func (m *model) handleGameEvent(e game.GameEvent) []tea.Cmd {
	var cmds []tea.Cmd

	switch e := e.(type) {
	case game.NewTurnEvent:
		m.state.game.turn = TurnUIState {
			prompt: e.Prompt,
			strikes: 0,
			prevAnswer: "",
		}

	case game.AnswerAcceptedEvent:
		msg := fmt.Sprintf("✓ %s  ", strings.ToUpper(e.Answer))
		m.state.game.gameMsg = msg
		m.anim_mgr.DeactivateAnimations(animations.ValidationMessage)

	case game.AnswerRejectedEvent:
		m.state.game.turn.prevAnswer = e.Answer
		m.state.game.gameMsg = e.Reason

	case game.StrikeEvent:
		m.state.game.gameMsg = e.Message
		m.state.game.turn.strikes = e.StrikeCount

		if e.Strikeout {
			m.anim_mgr.InitAnimations(animations.ValidationMessage)
			m.text_input.Reset()
			cmds = append(cmds, m.debounceInputCmd(500))
		} else {
			m.anim_mgr.InitAnimations(animations.StrikeCounter)
		}

	case game.PlayerDamagedEvent:
		m.state.game.playerDamaged = true
		m.state.game.health = e.Health
		cmds = append(cmds,
			m.togglePlayerDamagedCmd(),
			m.terminalBellCmd(false),
		)

	case game.ExtraLifeEvent:
		m.state.game.health = e.Health
		m.anim_mgr.InitAnimations(animations.ExtraLife)

	case game.GameQuitEvent:
		m.state.game.gameQuit = true

	case game.GameOverEvent:
		cmds = append(cmds, m.debounceInputCmd(500))
		m.state.game.gameOver = true
		m.state.game.possibleFinalAnswer = e.PossibleAnswer
		m.state.game.stats = e.Stats

	case game.GameWonEvent:
		m.state.game.playerDamaged = false
		m.state.game.gameOver = true
		m.state.game.gameMsg = ""
		m.anim_mgr.InitAnimations(animations.GameOverWin)
		m.state.game.stats = e.Stats

	default:
		slog.Warn("Game event not handled", "type", reflect.TypeOf(e).String(), "event", e)
	}

	return cmds
}
