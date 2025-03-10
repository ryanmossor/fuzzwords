package game

import (
	"fzwds/src/enums"
	"fzwds/src/utils"
	"strings"
)

type PlayerStats struct {
	PromptsSolved 			int
	PromptsFailed			int // TODO: store list of failed?
	ExtraLivesGained		int
	FewestExtraLifeSolves	int
	LongestSolve			string
	LetterCounts			map[string]int
	SolveLengths			[]int
	// TODO: most unique letters in a solve
	ElapsedSeconds			int
}

func InitializePlayerStats() PlayerStats {
	letter_counts := make(map[string]int)
	for _, c := range enums.Alphabets[enums.FullAlphabet] {
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

	for _, ch := range strings.ToUpper(answer) {
		s.LetterCounts[string(ch)] += 1
	}
}

func (s *PlayerStats) UpdateFailedStats() {
	s.PromptsFailed++
}

func (s PlayerStats) AverageSolveLength() float64 {
	return utils.Average(s.SolveLengths)
}
