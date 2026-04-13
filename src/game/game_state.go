package game

import (
	"fmt"
	"fzwds/src/dictionary"
	"fzwds/src/enums"
	"fzwds/src/utils"
	"log/slog"
	"slices"
	"strings"
	"time"
)

type GameState struct {
	Alphabet			string
	GameActive			bool
	EarlyQuit			bool
	GameStart			time.Time
	GameStop			time.Time
	Settings			GameSettings
	WordLists			WordLists
	Player				Player
	// TODO: consider making this a map[int]*Turn? key is turn number/idx
	// Would make accessing failed turns by idx easier
	turns				[]Turn
	// Indexes of failed turns
	FailedTurns			[]int
	StartUnixTs			int64
	GameWon				bool
}

func InitializeGame(settings *GameSettings) GameState {
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

    word_lists := WordLists {
		FULL_MAP: full_map,
		Available: available,
        Used: make(map[string]bool),
    }
	alphabet := enums.Alphabets[settings.Alphabet]

	g := GameState {
		StartUnixTs:		time.Now().UnixMilli(),

		Settings: 			*settings,
		Alphabet: 			alphabet,
		WordLists: 			word_lists,

		Player: 			InitializePlayer(settings, alphabet),
		GameActive: 		true,
		GameWon:			false,
		GameStart: 			time.Now(),

		// Prealloc 500 turns; should cover most games before slice needs to expand
		turns:				make([]Turn, 0, 500),
		FailedTurns:		[]int{},
	}
	g.NewTurn(true)

	slog.Info("Initialized game",
		"startUnixTs", g.StartUnixTs,
		"alphabet", g.Alphabet,
		"settings", g.Settings)

	return g
}

func (g *GameState) EndGame(won, early_quit bool) {
	if !g.GameActive {
		return
	}

	g.GameStop = time.Now()
	g.GameActive = false
	g.EarlyQuit = early_quit
	g.GameWon = won

	turn := g.CurrentTurn()
	turn.TotalTurnDuration = time.Since(turn.TurnStart)
	turn.FinalTurn = true

	if !won {
		g.Player.HealthCurrent = 0
	}
	g.Player.Stats = g.CalculateGameStats()
}

func (g *GameState) EndGameIfOver() bool {
	all_words_used := len(g.WordLists.Available) == 0
	max_lives_win := g.Settings.WinCondition == enums.WinConditionMaxLives &&
					 g.Player.HealthCurrent == g.Settings.HealthMax

	won := all_words_used || max_lives_win
	player_dead := g.Player.HealthCurrent == 0

	over := player_dead || won
	if !over {
		return false
	}

	g.EndGame(won, false)
	return true
}

func (g GameState) TurnCount() int {
	return len(g.turns)
}

func (g GameState) CurrentTurnNumber() int {
	return len(g.turns)
}

func (g *GameState) handleCorrectAnswer(answer string) {
	turn := g.CurrentTurn()
	turn.TotalTurnDuration = time.Since(turn.TurnStart)
	turn.Solved = true
	turn.Answer = answer
	turn.UniqueLetterCount = utils.CountUniqueLetters(answer)

	g.Player.Streak++
	turn.Streak = g.Player.Streak

	for _, c := range strings.ToUpper(answer) {
		ch := string(c)

		// TODO: consolidate LettersUsed/LettersRemaining, make []rune instead of []string?
		if strings.Contains(g.Alphabet, ch) && !slices.Contains(g.Player.LettersUsed, ch) {
			g.Player.LettersUsed = append(g.Player.LettersUsed, ch)
			turn.NewLettersUsed = append(turn.NewLettersUsed, c)
		}

		g.Player.LettersRemaining[c] = true
	}

	if len(g.Player.LettersUsed) >= len(g.Alphabet) {
		g.Player.LettersUsed = nil
		// TODO having letters remaining AND letters used seems redundant? consider consolidating into single map
		g.Player.LettersRemaining = utils.StringToCharMap(g.Alphabet)

		if g.Player.HealthCurrent < g.Settings.HealthMax {
			g.Player.HealthCurrent++
			turn.Health++
		}
		turn.ExtraLifeGained = true
	}

	slices.Sort(g.Player.LettersUsed)
}

type StrikeResult struct {
	Strikeout 		bool
	GameOver		bool
	Msg				string
}

// Handle timer expiry. Will increment strike counter, advance to
// next turn, or end the game depending on current game state.
func (g *GameState) HandleTurnTimeout() StrikeResult {
	turn := g.CurrentTurn()
	result := StrikeResult{}

	g.Player.Streak = 0
	turn.Streak = 0

	g.Player.HealthCurrent--
	turn.Health--

	turn.Strikes++

	if g.EndGameIfOver() {
		result.GameOver = true
		return result
	}

	if turn.Strikes == g.Settings.PromptStrikes {
		turn.TotalTurnDuration = time.Since(turn.TurnStart)
		result.Msg = fmt.Sprintf("Prompt %s failed", strings.ToUpper(turn.Prompt))
		result.Strikeout = true
		g.NewTurn(false)
	} else {
		g.StartStrikeTimer()
	}

	return result
}
