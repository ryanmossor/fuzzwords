package game

import (
	"fzw/src/enums"
	"fzw/src/utils"
	"strings"
)

type PlayerStats struct {
	PromptsSolved 			int
	PromptsFailed			int // TODO: store list of failed?
	TimeSurvived			int // TODO: TimeStarted/TimeDied unix timestamps on either player or stats struct; format as 1h23m45s or 1:23:45
	ExtraLivesGained		int
	FewestExtraLifeSolves	int
	LongestSolve			string
	LetterCounts			map[string]int
	SolveLengths			[]int
	// TODO: most unique letters in a solve
}

func InitializePlayerStats() PlayerStats {
	letter_counts := make(map[string]int)
	for _, c := range enums.FullAlphabet {
		letter_counts[string(c)] = 0
	}

	return PlayerStats{ LetterCounts: letter_counts }
}

func (s *PlayerStats) UpdateSolvedStats(answer string) {
	s.PromptsSolved++
	s.SolveLengths = append(s.SolveLengths, len(answer))

	if len(answer) > len(s.LongestSolve) {
		s.LongestSolve = answer
	}

	for _, c := range strings.ToUpper(answer) {
		s.LetterCounts[string(c)] += 1
	}
}

func (s *PlayerStats) UpdateFailedStats() {
	s.PromptsFailed++
}

func (s *PlayerStats) AverageSolveLength() float64 {
	return utils.Average(s.SolveLengths)
}

// TODO: func to generate stats table w/ str builder at end of game?
