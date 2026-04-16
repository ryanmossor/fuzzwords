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
	events = append(events, NewTurnEvent{ Prompt: g.currentTurn().prompt })

	slog.Info("Initialized game",
		"startUnixTs", g.startUnixTs,
		"alphabet", g.Settings.Alphabet.Letters(),
		"settings", g.Settings)

	return g, events
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
		events = append(events, ExtraLifeEvent{ Health: uint(g.Player.HealthCurrent) })
	}

	if g.determineWon() {
		g.endGame()
		events = append(events, GameWonEvent{})
		return events
	}

	g.newTurn(false)

	events = append(events, NewTurnEvent{ Prompt: g.currentTurn().prompt })

	return events
}

func (g *Game) AdvanceTime(now time.Time) []GameEvent {
	if !g.GameActive {
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

	g.Player.HealthCurrent--
	turn.health--
	events = append(events, PlayerDamagedEvent{ Health: uint(g.Player.HealthCurrent) })

	turn.strikes++
	strike_evt := StrikeEvent{}

	if g.Player.HealthCurrent <= 0 {
		g.endGame()
		events = append(events, strike_evt, GameOverEvent{ PossibleAnswer: turn.sourceWord })
		return events
	}

	if turn.strikes == g.Settings.PromptStrikes {
		turn.totalTurnDuration = time.Since(turn.turnStart)

		strike_evt.Strikeout = true
		strike_evt.Message = fmt.Sprintf("Prompt %s failed", strings.ToUpper(turn.prompt))

		g.newTurn(false)
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
					 g.Player.HealthCurrent == g.Settings.HealthMax

	return all_words_used || max_lives_win
}

func (g *Game) QuitGame() []GameEvent {
	if !g.GameActive {
		return nil
	}
	g.Quit = true
	g.endGame()
	return []GameEvent{ GameOverEvent{ g.currentTurn().sourceWord } }
}

func (g *Game) endGame() {
	if !g.GameActive {
		return
	}

	g.gameEnd = time.Now()
	g.GameActive = false
	g.GameWon = !g.Quit && g.determineWon()

	turn := g.currentTurn()
	turn.totalTurnDuration = time.Since(turn.turnStart)
	turn.finalTurn = true

	g.Player.Stats = g.CalculateGameStats()
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
