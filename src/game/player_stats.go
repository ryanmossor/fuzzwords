package game

import (
	"fzwds/src/utils"
)

type PlayerStats struct {
	PromptsSolved 			int
	PromptsFailed			int
	CurrentStreak			int
	LongestStreak			int
	ExtraLivesGained		int
	FewestExtraLifeSolves	int
	LongestSolve			string
	MostUniqueLetters		string
	SolveLengths			[]int
	AverageSolveLength		float64
	TimeSurvived			int
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
}

func (s *PlayerStats) UpdateFailedStats() {
	s.PromptsFailed++
	s.CurrentStreak = 0
}
