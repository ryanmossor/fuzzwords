package game

import (
	"fmt"
	"fzw/src/enums"
	"fzw/src/utils"
	"math"
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

func (s PlayerStats) AverageSolveLength() float64 {
	return utils.Average(s.SolveLengths)
}

func (s *PlayerStats) GenerateFinalStats() string {
	var stat_lines []string
	var sb strings.Builder

	stat_lines = append(stat_lines, fmt.Sprintf("Prompts solved: %d", s.PromptsSolved))
	stat_lines = append(stat_lines, fmt.Sprintf("Prompts failed: %d", s.PromptsFailed))
	stat_lines = append(stat_lines, fmt.Sprintf("Extra lives gained: %d", s.ExtraLivesGained))
	stat_lines = append(stat_lines, fmt.Sprintf("Fewest turns for extra life: %d", s.FewestExtraLifeSolves))
	stat_lines = append(stat_lines, fmt.Sprintf("Longest solve: %s (%d letters)", s.LongestSolve, len(s.LongestSolve)))
	stat_lines = append(stat_lines, fmt.Sprintf("Average solve length: %.1f letters", s.AverageSolveLength()))

	row_contents_len_max := len(utils.GetLongestStr(stat_lines)) + 2 // 2 padding chars for space before/after row contents
	if row_contents_len_max % 2 == 1 {
		row_contents_len_max++
	}

	BuildTableHeader(&sb, row_contents_len_max)

	for _, str := range stat_lines {
		PadTableRow(&sb, str, row_contents_len_max)
	}

	sb.WriteString("└")
	for i := 0; i < row_contents_len_max; i++ {
		sb.WriteString("─")
	}
	sb.WriteString("┘\n")

	return sb.String()
}

func BuildTableHeader(sb *strings.Builder, row_contents_len_max int) {
	sb.WriteString("┌")
	for i := 0; i < row_contents_len_max; i++ {
		sb.WriteString("─")
	}
	sb.WriteString("┐\n")


	sb.WriteString("│")
	header_padding := int(math.Floor(float64(row_contents_len_max - 10) / 2))
	for i := 0; i < header_padding; i++ {
		sb.WriteString(" ")
	}
	sb.WriteString("GAME STATS")
	for i := 0; i < header_padding; i++ {
		sb.WriteString(" ")
	}
	sb.WriteString("│\n")


	sb.WriteString("├")
	for i := 0; i < row_contents_len_max; i++ {
		sb.WriteString("─")
	}
	sb.WriteString("┤\n")
}

func PadTableRow(sb *strings.Builder, str string, row_contents_len_max int) {
	sb.WriteString(fmt.Sprintf("│ %s", str))
	for i := len(str); i <= row_contents_len_max - 2; i++ {
		sb.WriteString(" ")
	}
	sb.WriteString("│\n")
}
