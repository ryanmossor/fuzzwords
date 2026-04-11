package game

import (
	"fzwds/src/dictionary"
	"fzwds/src/enums"
	"fzwds/src/utils"
	"log/slog"
	"time"
)

type GameState struct {
	Alphabet			string
	GameActive			bool
	GameStart			time.Time
	GameEnd				time.Time
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
	// TODO: make this a helper method which returns 0 if len(turns) == 0 else CurrentTurn().TurnNumber
	CurrentTurnNumber	int
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
		CurrentTurnNumber: 	0,
	}
	g.NewTurn(true)

	slog.Info("Initialized game",
		"startUnixTs", g.StartUnixTs,
		"alphabet", g.Alphabet,
		"settings", g.Settings)

	return g
}

func (g *GameState) EndGame(won bool) {
	if !g.GameActive {
		return
	}

	g.GameEnd = time.Now()
	g.GameActive = false
	g.GameWon = won

	turn := g.CurrentTurn()
	turn.TotalTurnDuration = time.Since(turn.TurnStart)
	turn.FinalTurn = true

	if !won {
		g.Player.HealthCurrent = 0
	}
	g.Player.Stats = g.CalculateGameStats()
}

func (g *GameState) IsGameOver() bool {
	player_dead := g.Player.HealthCurrent == 0
	all_words_used := len(g.WordLists.Available) == 0
	max_lives_win := g.Settings.WinCondition == enums.WinConditionMaxLives &&
					 g.Player.HealthCurrent == g.Settings.HealthMax

	if player_dead || all_words_used || max_lives_win {
		return true
	}

	return false
}

func (g *GameState) CalculateGameStats() PlayerStats {
	start := time.Now()

	stats := PlayerStats{}
	stats.TimeSurvived = int(g.GameEnd.Sub(g.GameStart).Seconds())

	solve_lengths := make([]int, 0, len(g.turns))
	solve_len_idx := 0

	turns_since_last_extra_life := 0
	longest_streak := 0

	for i, turn := range g.turns {
		turns_since_last_extra_life++

		if turn.Solved {
			stats.PromptsSolved++

			if turn.Streak > longest_streak {
				longest_streak = turn.Streak
			}

			solve_lengths = append(solve_lengths, len(turn.Answer))
			solve_len_idx++

			if len(turn.Answer) > len(stats.LongestSolve) {
				stats.LongestSolve = turn.Answer
			}

			if utils.CountUniqueLetters(turn.Answer) > utils.CountUniqueLetters(stats.MostUniqueLetters) {
				stats.MostUniqueLetters = turn.Answer
			}

			if turn.ExtraLifeGained {
				stats.ExtraLivesGained++
				if stats.FewestExtraLifeSolves == 0 || turns_since_last_extra_life < stats.FewestExtraLifeSolves {
					stats.FewestExtraLifeSolves = turns_since_last_extra_life
				}
				turns_since_last_extra_life = 0
			}
		} else {
			stats.PromptsFailed++
			g.FailedTurns = append(g.FailedTurns, i)
		}
	}

	stats.AverageSolveLength = utils.Average(solve_lengths)
	stats.LongestStreak = longest_streak

	elapsed := time.Since(start)

	slog.Debug("Calculated stats for game",
		"startUnixTx", g.StartUnixTs,
		"turns", len(g.turns),
		"gameDuration", utils.FormatTime(int(g.GameEnd.Sub(g.GameStart).Seconds())),
		"calcTimeMs", elapsed.Milliseconds(),
	)

	return stats
}
