package game

import (
	"fzwds/src/utils"
	"log/slog"
	"time"
)

type PlayerStats struct {
	timePlayed            time.Duration
	promptsSolved         int
	promptsFailed         int
	solvesPerMinute       float64
	averageSolveLength    float64
	longestStreak         int
	extraLivesGained      int
	fewestExtraLifeSolves int
	mostUniqueCount       int
	mostUniqueWord        string
	longestSolve          string
}

func (s PlayerStats) TimePlayed() time.Duration {
	return s.timePlayed
}

func (s PlayerStats) PromptsSolved() int {
	return s.promptsSolved
}

func (s PlayerStats) SolvesPerMinute() float64 {
	return s.solvesPerMinute
}

func (s PlayerStats) AverageSolveLength() float64 {
	return s.averageSolveLength
}

func (s PlayerStats) LongestStreak() int {
	return s.longestStreak
}

func (s PlayerStats) ExtraLivesGained() int {
	return s.extraLivesGained
}

func (s PlayerStats) FewestExtraLifeSolves() int {
	return s.fewestExtraLifeSolves
}

func (s PlayerStats) MostUniqueCount() int {
	return s.mostUniqueCount
}

func (s PlayerStats) MostUniqueWord() string {
	return s.mostUniqueWord
}

func (s PlayerStats) LongestSolve() string {
	return s.longestSolve
}

type Player struct {
	healthCurrent    int
	streak           int
	lettersUsed		 map[rune]bool
	stats            PlayerStats
}

func newPlayer(settings GameSettings) Player {
	player := Player {
		healthCurrent:  settings.HealthInitial,
		streak:         0,
		lettersUsed:	utils.StringToCharMap(settings.Alphabet.Letters()),
		stats:          PlayerStats{},
	}
	return player
}

func (g *Game) calculateGameStats() PlayerStats {
	start := time.Now()

	stats := PlayerStats{}
	stats.timePlayed = g.gameEnd.Sub(g.gameStart)

	solve_lengths := make([]int, 0, len(g.turns))
	solve_len_idx := 0

	turns_since_last_extra_life := 0
	longest_streak := 0

	for i, turn := range g.turns {
		turns_since_last_extra_life++

		if turn.solved {
			stats.promptsSolved++

			if turn.streak > longest_streak {
				longest_streak = turn.streak
			}

			solve_lengths = append(solve_lengths, len(turn.answer))
			solve_len_idx++

			if len(turn.answer) > len(stats.longestSolve) {
				stats.longestSolve = turn.answer
			}

			if turn.uniqueLetterCount > stats.mostUniqueCount {
				stats.mostUniqueWord = turn.answer
				stats.mostUniqueCount = turn.uniqueLetterCount
			}

			if turn.extraLifeGained {
				stats.extraLivesGained++
				if stats.fewestExtraLifeSolves == 0 || turns_since_last_extra_life < stats.fewestExtraLifeSolves {
					stats.fewestExtraLifeSolves = turns_since_last_extra_life
				}
				turns_since_last_extra_life = 0
			}
		} else {
			stats.promptsFailed++
			g.failedTurns = append(g.failedTurns, i)
		}
	}

	stats.averageSolveLength = utils.Average(solve_lengths)
	stats.solvesPerMinute = float64(stats.promptsSolved) / (float64(stats.timePlayed.Seconds()) / 60.0)
	stats.longestStreak = longest_streak

	elapsed := time.Since(start)

	slog.Debug("Calculated stats for game",
		"startUnixTx", g.startUnixTs,
		"turns", len(g.turns),
		"gameDuration", utils.FormatTime(stats.timePlayed),
		"calcTimeMs", elapsed.Milliseconds(),
	)

	return stats
}
