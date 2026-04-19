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
	gameActive			bool
	gameWon				bool
	startUnixTs			int64
	gameStart			time.Time
	gameEnd				time.Time
	// Indexes of failed turns
	failedTurns			[]int

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
		gameActive: 	true,
		gameWon:		false,
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
	g.newTurn(TransitionFirstTurn)

	var events []GameEvent
	events = append(events, NewTurnEvent{ Prompt: g.currentTurn().prompt })

	slog.Info("Initialized game",
		"startUnixTs", g.startUnixTs,
		"alphabet", g.Settings.Alphabet.Letters(),
		"settings", g.Settings)

	return g, events
}

func (g *Game) SubmitAnswer(answer string) []GameEvent {
	var events []GameEvent

	if g.TimeRemaining() <= 0 {
		events = g.handleTimeout()
		return events
	}

	result := g.validateAnswer(answer)
	if !result.accepted {
		events = append(events, AnswerRejectedEvent {
			Answer: answer,
			Reason: result.reason,
		})
		return events
	}

	life_gained := g.handleCorrectAnswer(answer)
	if life_gained {
		events = append(events, ExtraLifeEvent{ Health: uint(g.Player.healthCurrent) })
	}

	events = append(events, AnswerAcceptedEvent {
		Answer: answer,
		NewLettersUsed: g.currentTurn().newLettersUsed,
	})

	if g.determineWon() {
		g.endGame()
		events = append(events, GameWonEvent{ Stats: g.Player.stats })
		return events
	}

	g.newTurn(TransitionSolved)

	events = append(events, NewTurnEvent{ Prompt: g.currentTurn().prompt })

	return events
}

func (g *Game) AdvanceTime(now time.Time) []GameEvent {
	if !g.gameActive {
		return nil
	}

	turn_end := g.currentTurn().strikeStart.Add(g.currentTurn().strikeDuration)
	if now.After(turn_end) {
		events := g.handleTimeout()
		return events
	}

	return nil
}

// Handle timer expiry. Will increment strike counter, advance to
// next turn, or end the game depending on current game state.
func (g *Game) handleTimeout() []GameEvent {
	turn := g.currentTurn()
	var events []GameEvent

	g.Player.streak = 0
	turn.streak = 0

	g.Player.healthCurrent--
	turn.health--
	events = append(events, PlayerDamagedEvent{ Health: uint(g.Player.healthCurrent) })

	turn.strikes++
	strike_evt := StrikeEvent{}

	if g.Player.healthCurrent <= 0 {
		g.endGame()
		events = append(events, strike_evt, GameOverEvent{
			PossibleAnswer: turn.sourceWord,
			Stats: g.Player.stats,
		})
		return events
	}

	if turn.strikes == g.Settings.PromptStrikes {
		turn.totalTurnDuration = time.Since(turn.turnStart)

		strike_evt.Strikeout = true
		strike_evt.Message = fmt.Sprintf("Prompt %s failed", strings.ToUpper(turn.prompt))

		g.newTurn(TransitionTimeout)
		events = append(events, NewTurnEvent{ Prompt: g.currentTurn().prompt })
	} else {
		g.startStrikeTimer()
	}
	strike_evt.StrikeCount = g.currentTurn().strikes
	events = append(events, strike_evt)

	return events
}

func (g Game) determineWon() bool {
	all_words_used := len(g.wordLists.available) == 0
	max_lives_win := g.Settings.WinCondition == enums.WinConditionMaxLives &&
					 g.Player.healthCurrent == g.Settings.HealthMax

	return all_words_used || max_lives_win
}

func (g *Game) QuitGame() []GameEvent {
	if !g.gameActive {
		return nil
	}
	g.endGame()

	return []GameEvent {
		GameOverEvent{
			PossibleAnswer: g.currentTurn().sourceWord,
			Stats: g.Player.stats,
		},
		GameQuitEvent{},
	}
}

func (g *Game) endGame() {
	if !g.gameActive {
		return
	}

	g.gameEnd = time.Now()
	g.gameActive = false
	g.gameWon = g.determineWon()

	turn := g.currentTurn()
	turn.totalTurnDuration = time.Since(turn.turnStart)
	turn.finalTurn = true

	g.Player.stats = g.calculateGameStats()
}

func (g Game) GameActive() bool {
	return g.gameActive
}

func (g Game) GameWon() bool {
	return g.gameWon
}

func (g Game) TurnCount() int {
	return len(g.turns)
}

func (g Game) CurrentTurnNumber() int {
	return len(g.turns)
}

func (g Game) currentTurn() *Turn {
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
		if (g.turns[i].strikes > 0 || !g.turns[i].solved) && i > turn_idx_cur {
			return i
		}
	}
	return turn_idx_cur
}

func (g Game) PrevFailedTurnIdx(turn_idx_cur int) int {
	for i := turn_idx_cur; i >= 0; i-- {
		if (g.turns[i].strikes > 0 || !g.turns[i].solved) && i < turn_idx_cur {
			return i
		}
	}
	return turn_idx_cur
}
