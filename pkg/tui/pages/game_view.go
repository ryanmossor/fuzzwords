package pages
//
// import (
// 	"fzwds/pkg/game"
// 	"fzwds/pkg/tui/animations"
// 	"fzwds/pkg/tui/pages"
// 	"fzwds/pkg/tui/styles"
// 	"fzwds/pkg/utils"
// 	"log/slog"
// 	"reflect"
// 	"strings"
//
// 	"github.com/charmbracelet/bubbles/textinput"
// 	tea "github.com/charmbracelet/bubbletea"
// 	"github.com/charmbracelet/lipgloss"
// )
//
// type turnUIState struct {
// 	prompt		string
// 	strikes		int
// 	prevAnswer	string
// }
//
// type gameUIState struct {
// 	playerDamaged		bool
// 	inputRestricted		bool
// 	gameOver			bool
// 	gameQuit			bool
// 	gameWon				bool
// 	health				int
// 	gameMsg				string
// 	possibleFinalAnswer	string
// 	lettersUsed 		map[rune]bool
// 	turn				turnUIState
// 	stats				game.PlayerStats
// }
//
// func (m model) GameView() string {
// 	prompt := styles.TextAccent.
// 		Bold(true).
// 		Render(strings.ToUpper(m.state.game.turn.prompt))
//
// 	var game_msg string
// 	if m.state.game.gameMsg != "" {
// 		game_msg = m.renderValidationMsg()
// 	} else {
// 		game_msg = m.highlightPromptAnswer(
// 			m.state.game.turn.prompt,
// 			m.gameInput.Value(),
// 			m.game.Settings().PromptMode)
// 	}
//
// 	return lipgloss.JoinVertical(
// 		lipgloss.Center,
// 		"",
// 		"",
// 		"",
// 		"",
// 		"",
// 		"",
// 		"",
// 		"",
// 		"",
// 		prompt,
// 		m.GameStrikeCounterView(),
// 		"",
// 		"",
// 		"",
// 		"",
// 		"",
// 		"",
// 		game_msg,
// 		"",
// 		m.GetBlockInputView(),
// 	)
// }
//
// func (m model) GameSwitch() (model, tea.Cmd) {
// 	m = m.SwitchPage(pages.GamePage)
//
// 	m.footerKeymaps = []footerKeymap {
// 		{key: "esc", value: "clear input"},
// 		{key: "ctrl+q", value: "quit"},
// 	}
//
// 	var events []game.GameEvent
// 	m.game, events = game.NewGame(&m.settings.Game)
//
// 	m.state.game = gameUIState {
// 		lettersUsed: utils.StringToCharMap(m.game.Settings().Alphabet.Letters()),
// 	}
// 	m.state.gameOver = gameOverState {
// 		viewCache: make(map[string]string),
// 	}
// 	m.state.gameReview = gameReviewState {
// 		viewCache: make(map[int]*turnDisplay),
// 	}
//
// 	var cmds []tea.Cmd
// 	for _, e := range events {
// 		cmds = append(cmds, m.handleGameEvent(e)...)
// 	}
//
// 	m.state.game.health = m.game.Settings().HealthInitial
//
// 	m.gameInput = m.initBlockTextInput()
// 	cmds = append(cmds, textinput.Blink)
//
// 	return m, tea.Batch(cmds...)
// }
//
// func (m model) GameUpdate(msg tea.Msg) (model, tea.Cmd) {
// 	var cmds []tea.Cmd
//
// 	switch msg := msg.(type) {
// 	case TickMsg:
// 		events := m.game.AdvanceTime(msg.Time)
// 		if len(events) == 0 {
// 			return m, nil
// 		}
//
// 		for _, e := range events {
// 			cmds = append(cmds, m.handleGameEvent(e)...)
// 		}
//
//         if m.state.game.gameOver {
// 			var cmd tea.Cmd
// 			m, cmd = m.GameOverSwitch()
// 			cmds = append(cmds, cmd)
// 		}
//
//         return m, tea.Batch(cmds...)
//
// 	case tea.KeyMsg:
//         if m.state.game.inputRestricted {
//             return m, nil
//         }
//
//         key := msg.String()
// 		if key != "enter" {
// 			m.state.game.gameMsg = ""
// 		}
// 		m.animManager.DeactivateAnimations(animations.ValidationMessage)
//
// 		// TODO: skips -- sacrifice life for skip, earn extra skips from extra lifes if already full
// 		// TODO: mute key combo for alert sound in game?
// 		switch key {
// 		case "up":
// 			if m.state.game.turn.prevAnswer != "" {
// 				m.gameInput.SetValue(m.state.game.turn.prevAnswer)
// 			}
// 			return m, nil
//
// 		case "esc":
// 			m.gameInput.Reset()
// 			return m, nil
//
// 		case "ctrl+q":
// 			events := m.game.QuitGame()
// 			for _, e := range events {
// 				cmds = append(cmds, m.handleGameEvent(e)...)
// 			}
//
// 			m, cmd := m.GameOverSwitch()
// 			cmds = append(cmds, cmd)
//
// 			return m, tea.Batch(cmds...)
//
// 		case "enter":
// 			answer := strings.ToLower(strings.TrimSpace(m.gameInput.Value()))
//             m.gameInput.Reset()
//
// 			events := m.game.SubmitAnswer(answer)
// 			for _, e := range events {
// 				cmds = append(cmds, m.handleGameEvent(e)...)
// 			}
//
// 			cmds = append(cmds, m.debounceInputCmd(300))
//
// 			if m.state.game.gameOver {
// 				var cmd tea.Cmd
// 				m, cmd = m.GameOverSwitch()
// 				cmds = append(cmds, cmd)
// 			}
//
// 			return m, tea.Batch(cmds...)
// 		}
// 	}
//
// 	var update_input_cmd tea.Cmd
// 	m.gameInput, update_input_cmd = m.gameInput.Update(msg)
//
// 	return m, update_input_cmd
// }
//
// func (m *model) handleGameEvent(event game.GameEvent) []tea.Cmd {
// 	var cmds []tea.Cmd
//
// 	switch e := event.(type) {
// 	case game.NewTurnEvent:
// 		m.state.game.turn = turnUIState {
// 			prompt: e.Prompt,
// 			strikes: 0,
// 			prevAnswer: "",
// 		}
//
// 	case game.AnswerAcceptedEvent:
// 		msg := styles.TextGreen.Render("✓ " + strings.ToUpper(e.Answer) + "  ")
// 		m.state.game.gameMsg = msg
// 		m.animManager.DeactivateAnimations(animations.ValidationMessage)
// 		for _, c := range e.NewLettersUsed {
// 			m.state.game.lettersUsed[c] = true
// 		}
//
// 	case game.AnswerRejectedEvent:
// 		m.state.game.turn.prevAnswer = e.Answer
//
// 		var msg string
//
// 		switch e.Reason {
// 		case game.RejectionEmpty:
// 			msg = "No answer given"
// 		case game.RejectionInvalidWord:
// 			msg = "Invalid word: " + strings.ToUpper(e.Answer)
// 		case game.RejectionPromptMismatch:
// 			msg = strings.ToUpper(e.Answer) + " does not satisfy prompt"
// 		case game.RejectionAlreadyUsed:
// 			msg = "🔒 " + strings.ToUpper(e.Answer) + " already used"
// 		}
//
// 		m.state.game.gameMsg = styles.TextRed.Render(msg)
//
// 	case game.StrikeEvent:
// 		m.state.game.turn.strikes = e.StrikeCount
//
// 		if e.Strikeout {
// 			m.state.game.gameMsg = styles.TextRed.Render("Prompt " + strings.ToUpper(e.Prompt) + " failed")
// 			m.animManager.InitAnimations(animations.ValidationMessage)
// 			m.gameInput.Reset()
// 			cmds = append(cmds, m.debounceInputCmd(500))
// 		} else {
// 			m.animManager.InitAnimations(animations.StrikeCounter)
// 		}
//
// 	case game.PlayerDamagedEvent:
// 		m.state.game.playerDamaged = true
// 		m.state.game.health = e.Health
// 		cmds = append(cmds,
// 			m.togglePlayerDamagedCmd(),
// 			m.terminalBellCmd(false),
// 		)
//
// 	case game.ExtraLifeEvent:
// 		m.state.game.health = e.Health
// 		m.animManager.InitAnimations(animations.ExtraLife)
// 		m.state.game.lettersUsed = utils.StringToCharMap(m.game.Settings().Alphabet.Letters())
//
// 	case game.GameQuitEvent:
// 		m.state.game.gameQuit = true
//
// 	case game.GameOverEvent:
// 		cmds = append(cmds, m.debounceInputCmd(500))
// 		m.state.game.gameOver = true
// 		m.state.game.possibleFinalAnswer = e.PossibleAnswer
// 		m.state.game.stats = e.Stats
// 		m.state.gameReview.turns = e.Turns
//
// 	case game.GameWonEvent:
// 		m.state.game.playerDamaged = false
// 		m.state.game.gameOver = true
// 		m.state.game.gameWon = true
// 		m.state.game.gameMsg = ""
// 		m.animManager.InitAnimations(animations.GameOverWin)
// 		m.state.game.stats = e.Stats
// 		m.state.gameReview.turns = e.Turns
//
// 	default:
// 		slog.Warn("Game event not handled", "type", reflect.TypeOf(e).String(), "event", e)
// 	}
//
// 	return cmds
// }
