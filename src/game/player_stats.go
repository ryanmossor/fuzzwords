package game

import (
	"fzwds/src/enums"
	"fzwds/src/utils"
	"strings"
)

type PlayerStats struct {
	PromptsSolved 			int
	PromptsFailed			int // TODO: store list of failed?
	CurrentStreak			int
	LongestStreak			int
	ExtraLivesGained		int
	FewestExtraLifeSolves	int
	LongestSolve			string
	MostUniqueLetters		string
	LetterCounts			map[string]int
	SolveLengths			[]int
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

	s.CurrentStreak++
	if s.CurrentStreak > s.LongestStreak {
		s.LongestStreak = s.CurrentStreak
	}

	s.SolveLengths = append(s.SolveLengths, len(answer))

	if len(answer) > len(s.LongestSolve) {
		s.LongestSolve = answer
	}

	if utils.CountUniqueLetters(answer) > utils.CountUniqueLetters(s.MostUniqueLetters) {
		s.MostUniqueLetters = answer
	}

	for _, ch := range strings.ToUpper(answer) {
		s.LetterCounts[string(ch)] += 1
	}
}

func (s *PlayerStats) UpdateFailedStats() {
	s.PromptsFailed++
	s.CurrentStreak = 0
}

func (s PlayerStats) AverageSolveLength() float64 {
	return utils.Average(s.SolveLengths)
}
