package game

import (
	"fzwds/src/assert"
	"fzwds/src/dictionary"
	"fzwds/src/enums"
	"fzwds/src/utils"
	"log/slog"
	"time"
)

type GameState struct {
	GameActive			bool
	GameWon				bool
	EarlyQuit			bool
	startUnixTs			int64
	gameStart			time.Time
	gameEnd				time.Time
	Alphabet			string
	// Indexes of failed turns
	failedTurns			[]int

	Settings			GameSettings
	wordLists			wordLists
	Player				Player
	// TODO: consider making this a map[int]*Turn? key is turn number/idx
	// Would make accessing failed turns by idx easier
	turns				[]Turn
}

func NewGame(settings *GameSettings) GameState {
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

    word_lists := wordLists {
		fullDict: full_map,
		available: available,
        used: make(map[string]bool),
    }
	alphabet := enums.Alphabets[settings.Alphabet]

	g := GameState {
		GameActive: 		true,
		GameWon:			false,
		gameStart: 			time.Now(),
		startUnixTs:		time.Now().UnixMilli(),

		Alphabet: 			alphabet,
		failedTurns:		[]int{},

		Settings: 			*settings,
		wordLists: 			word_lists,
		turns:				make([]Turn, 0, 300), // 300 should cover most games before realloc needed
	}
	g.Player = g.newPlayer()
	g.newTurn(true)

	slog.Info("Initialized game",
		"startUnixTs", g.startUnixTs,
		"alphabet", g.Alphabet,
		"settings", g.Settings)

	return g
}

func (g GameState) determineWon() bool {
	all_words_used := len(g.wordLists.available) == 0
	max_lives_win := g.Settings.WinCondition == enums.WinConditionMaxLives &&
					 g.Player.HealthCurrent == g.Settings.HealthMax

	return all_words_used || max_lives_win
}

func (g *GameState) EndGame(early_quit bool) {
	if !g.GameActive {
		return
	}

	g.gameEnd = time.Now()
	g.GameActive = false
	g.EarlyQuit = early_quit
	g.GameWon = g.determineWon()

	turn := g.CurrentTurn()
	turn.TotalTurnDuration = time.Since(turn.turnStart)
	turn.FinalTurn = true

	g.Player.Stats = g.CalculateGameStats()
}

func (g *GameState) endGameIfOver() bool {
	over := g.Player.HealthCurrent == 0 || g.determineWon()
	if over {
		g.EndGame(false)
		return true
	}
	return false
}

func (g GameState) TurnCount() int {
	return len(g.turns)
}

func (g GameState) CurrentTurnNumber() int {
	return len(g.turns)
}

func (g GameState) CurrentTurn() *Turn {
	assert.Assert(len(g.turns) > 0, "Attempted to access current turn before game initialized")
	return &g.turns[len(g.turns) - 1]
}

func (g GameState) PreviousTurn() (*Turn, bool) {
	if len(g.turns) <= 1 {
		return nil, false
	}
	return &g.turns[len(g.turns) - 2], true
}

// TODO: this takes turn idx, but i'm referencing turns by number (1-based) in many places.
// Consider pros/cons of idx vs turn number
func (g GameState) GetTurn(idx int) *Turn {
	clamped_idx := utils.Clamp(idx, 0, len(g.turns) - 1)
	return &g.turns[clamped_idx]
}

func (g GameState) NextFailedTurnIdx(turn_idx_cur int) int {
	for i := turn_idx_cur; i < len(g.turns); i++ {
		if (g.turns[i].Strikes > 0 || !g.turns[i].Solved) && i > turn_idx_cur {
			return i
		}
	}
	return turn_idx_cur
}

func (g GameState) PrevFailedTurnIdx(turn_idx_cur int) int {
	for i := turn_idx_cur; i >= 0; i-- {
		if (g.turns[i].Strikes > 0 || !g.turns[i].Solved) && i < turn_idx_cur {
			return i
		}
	}
	return turn_idx_cur
}
