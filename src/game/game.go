package game

import (
	"fmt"
	"fzwds/src/assert"
	"fzwds/src/dictionary"
	"fzwds/src/enums"
	"fzwds/src/utils"
	"log/slog"
	"strings"
	"time"
)

type Game struct {
	GameActive			bool
	GameWon				bool
	Quit				bool
	startUnixTs			int64
	gameStart			time.Time
	gameEnd				time.Time
	// Indexes of failed turns
	failedTurns			[]int
	// Incrementing timer id counter; allows caller to skip HandleTurnTimeout on stale timeouts
	// TODO: replace with something better (eg, game controls timeouts, caller checks for events?)
	timerId				uint

	Settings			GameSettings
	wordLists			wordLists
	Player				Player
	// TODO: consider making this a map[int]*Turn? key is turn number/idx
	// Would make accessing failed turns by idx easier
	turns				[]Turn
}

func NewGame(settings *GameSettings) (Game, []GameEvent) {
	var full_map map[string]bool
	var available []string

	switch settings.Dictionary {
	case enums.English:
		full_map = dictionary.EnglishDictionaryMap
		available = utils.FilterWordList(dictionary.EnglishDictionary, settings.PromptLenMin)
	case enums.Pokemon:
		available = dictionary.GetSelectedPokemonGenList(settings.PokemonGens...)
		full_map = utils.ArrToMap(available)
	}

	g := Game {
		GameActive: 	true,
		GameWon:		false,
		gameStart: 		time.Now(),
		startUnixTs:	time.Now().UnixMilli(),
		failedTurns:	[]int{},
		Settings: 		*settings,
		wordLists: 		wordLists {
			fullDict: 		full_map,
			available: 		available,
			used: 			make(map[string]bool),
		},
		Player:			newPlayer(*settings),
		turns:			make([]Turn, 0, 300), // 300 should cover most games before realloc needed
	}
	g.newTurn(true)

	var events []GameEvent
	events = append(events, TimerTickEvent {
		TimerId: g.timerId,
		Duration: g.CurrentTurn().strikeDuration,
	})

	slog.Info("Initialized game",
		"startUnixTs", g.startUnixTs,
		"alphabet", g.Settings.Alphabet.Letters(),
		"settings", g.Settings,
		"events", events)

	return g, events
}

func (g Game) TimerId() uint {
	return g.timerId
}

func (g *Game) SubmitAnswer(answer string) []GameEvent {
	var events []GameEvent

	answer_res := g.validateAnswer(answer)
	if !answer_res.accepted {
		events = append(events, AnswerRejectedEvent {
			Answer: answer,
			Reason: answer_res.reason,
		})
		return events
	}
	events = append(events, AnswerAcceptedEvent{ Answer: answer })

	life_gained := g.handleCorrectAnswer(answer)
	if life_gained {
		events = append(events, ExtraLifeEvent{})
	}

	if g.determineWon() {
		g.endGame()
		events = append(events, GameWonEvent{})
		return events
	}

	g.newTurn(false)

	events = append(events, TimerTickEvent {
		TimerId: g.timerId,
		Duration: g.CurrentTurn().strikeDuration,
	})

	return events
}

// Handle timer expiry. Will increment strike counter, advance to
// next turn, or end the game depending on current game state.
// TODO: make this internal w/ AdvanceTime(time.Time) exposed to caller
// Instead of caller drive turn timeouts, game engine manages them and caller polls for updates?
func (g *Game) HandleTurnTimeout() []GameEvent {
	turn := g.CurrentTurn()
	var events []GameEvent

	g.Player.streak = 0
	turn.Streak = 0

	g.Player.HealthCurrent--
	turn.Health--
	events = append(events, PlayerDamagedEvent{})

	turn.Strikes++
	strike_evt := StrikeEvent{}

	if g.Player.HealthCurrent == 0 {
		g.endGame()
		events = append(events, strike_evt, GameOverEvent{})
		return events
	}

	if turn.Strikes == g.Settings.PromptStrikes {
		turn.TotalTurnDuration = time.Since(turn.turnStart)

		strike_evt.Strikeout = true
		strike_evt.Message = fmt.Sprintf("Prompt %s failed", strings.ToUpper(turn.Prompt))

		g.newTurn(false)
	} else {
		g.startStrikeTimer()
	}
	events = append(events, strike_evt)

	events = append(events, TimerTickEvent {
		TimerId: g.timerId,
		Duration: g.CurrentTurn().strikeDuration,
	})

	return events
}

func (g Game) determineWon() bool {
	all_words_used := len(g.wordLists.available) == 0
	max_lives_win := g.Settings.WinCondition == enums.WinConditionMaxLives &&
					 g.Player.HealthCurrent == g.Settings.HealthMax

	return all_words_used || max_lives_win
}

func (g *Game) QuitGame() []GameEvent {
	if !g.GameActive {
		return nil
	}
	g.Quit = true
	g.endGame()
	return []GameEvent{ GameQuitEvent{} }
}

func (g *Game) endGame() {
	if !g.GameActive {
		return
	}

	g.gameEnd = time.Now()
	g.GameActive = false
	g.GameWon = !g.Quit && g.determineWon()

	turn := g.CurrentTurn()
	turn.TotalTurnDuration = time.Since(turn.turnStart)
	turn.FinalTurn = true

	g.Player.Stats = g.CalculateGameStats()
}

func (g Game) TurnCount() int {
	return len(g.turns)
}

func (g Game) CurrentTurnNumber() int {
	return len(g.turns)
}

func (g Game) CurrentTurn() *Turn {
	assert.Assert(len(g.turns) > 0, "Attempted to access current turn before game initialized")
	return &g.turns[len(g.turns) - 1]
}

func (g Game) PreviousTurn() (*Turn, bool) {
	if len(g.turns) <= 1 {
		return nil, false
	}
	return &g.turns[len(g.turns) - 2], true
}

// TODO: this takes turn idx, but i'm referencing turns by number (1-based) in many places.
// Consider pros/cons of idx vs turn number
func (g Game) GetTurn(idx int) *Turn {
	clamped_idx := utils.Clamp(idx, 0, len(g.turns) - 1)
	return &g.turns[clamped_idx]
}

func (g Game) NextFailedTurnIdx(turn_idx_cur int) int {
	for i := turn_idx_cur; i < len(g.turns); i++ {
		if (g.turns[i].Strikes > 0 || !g.turns[i].Solved) && i > turn_idx_cur {
			return i
		}
	}
	return turn_idx_cur
}

func (g Game) PrevFailedTurnIdx(turn_idx_cur int) int {
	for i := turn_idx_cur; i >= 0; i-- {
		if (g.turns[i].Strikes > 0 || !g.turns[i].Solved) && i < turn_idx_cur {
			return i
		}
	}
	return turn_idx_cur
}
