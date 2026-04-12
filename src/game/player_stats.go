package game

import (
	"fzwds/src/utils"
	"log/slog"
	"time"
)

type PlayerStats struct {
	PromptsSolved 			int
	PromptsFailed			int
	LongestStreak			int
	ExtraLivesGained		int
	FewestExtraLifeSolves	int
	LongestSolve			string
	MostUniqueWord			string
	MostUniqueCount			int
	AverageSolveLength		float64
	TimePlayed				int
}

func (g *GameState) CalculateGameStats() PlayerStats {
	start := time.Now()

	stats := PlayerStats{}
	stats.TimePlayed = int(g.GameEnd.Sub(g.GameStart).Seconds())

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

			if turn.UniqueLetterCount > stats.MostUniqueCount {
				stats.MostUniqueWord = turn.Answer
				stats.MostUniqueCount = turn.UniqueLetterCount
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

