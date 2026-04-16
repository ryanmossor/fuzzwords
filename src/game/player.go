package game

import (
	"fzwds/src/utils"
	"log/slog"
	"time"
)

type PlayerStats struct {
	TimePlayed            time.Duration
	PromptsSolved         int
	PromptsFailed         int
	SolvesPerMinute       float64
	AverageSolveLength    float64
	LongestStreak         int
	ExtraLivesGained      int
	FewestExtraLifeSolves int
	MostUniqueCount       int
	MostUniqueWord        string
	LongestSolve          string
}

type Player struct {
	HealthCurrent    int
	LettersUsed      []rune
	LettersRemaining map[rune]bool
	streak           int
	Stats            PlayerStats
}

func newPlayer(settings GameSettings) Player {
	player := Player{
		HealthCurrent:    settings.HealthInitial,
		LettersRemaining: utils.StringToCharMap(settings.Alphabet.Letters()),
		LettersUsed:      make([]rune, 0, len(settings.Alphabet.Letters())),
		streak:           0,
		Stats:            PlayerStats{},
	}
	return player
}

func (g *Game) CalculateGameStats() PlayerStats {
	start := time.Now()

	stats := PlayerStats{}
	stats.TimePlayed = g.gameEnd.Sub(g.gameStart)

	solve_lengths := make([]int, 0, len(g.turns))
	solve_len_idx := 0

	turns_since_last_extra_life := 0
	longest_streak := 0

	for i, turn := range g.turns {
		turns_since_last_extra_life++

		if turn.solved {
			stats.PromptsSolved++

			if turn.streak > longest_streak {
				longest_streak = turn.streak
			}

			solve_lengths = append(solve_lengths, len(turn.answer))
			solve_len_idx++

			if len(turn.answer) > len(stats.LongestSolve) {
				stats.LongestSolve = turn.answer
			}

			if turn.uniqueLetterCount > stats.MostUniqueCount {
				stats.MostUniqueWord = turn.answer
				stats.MostUniqueCount = turn.uniqueLetterCount
			}

			if turn.extraLifeGained {
				stats.ExtraLivesGained++
				if stats.FewestExtraLifeSolves == 0 || turns_since_last_extra_life < stats.FewestExtraLifeSolves {
					stats.FewestExtraLifeSolves = turns_since_last_extra_life
				}
				turns_since_last_extra_life = 0
			}
		} else {
			stats.PromptsFailed++
			g.failedTurns = append(g.failedTurns, i)
		}
	}

	stats.AverageSolveLength = utils.Average(solve_lengths)
	stats.SolvesPerMinute = float64(stats.PromptsSolved) / (float64(stats.TimePlayed.Seconds()) / 60.0)
	stats.LongestStreak = longest_streak

	elapsed := time.Since(start)

	slog.Debug("Calculated stats for game",
		"startUnixTx", g.startUnixTs,
		"turns", len(g.turns),
		"gameDuration", utils.FormatTime(stats.TimePlayed),
		"calcTimeMs", elapsed.Milliseconds(),
	)

	return stats
}
