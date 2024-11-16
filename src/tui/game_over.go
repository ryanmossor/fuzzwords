package tui

import (
	"fmt"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

func (m model) GameOverSwitch() (model, tea.Cmd) {
	m = m.SwitchPage(game_over_page)
	m.footer_cmds = []footerCmd{
		{key: "enter", value: "play again"},
		{key: "q", value: "quit"},
	}

	return m, nil
}

func (m model) GameOverUpdate(msg tea.Msg) (model, tea.Cmd) {
	return m, nil
}

func (m model) GameOverView() string {
	longest_solve := m.player.Stats.LongestSolve
	if longest_solve == "" {
		longest_solve = "-"
	}
	stats := [][]string{
		{"Prompts solved", strconv.Itoa(m.player.Stats.PromptsSolved)},
		{"Prompts failed", strconv.Itoa(m.player.Stats.PromptsFailed)},
		{"Extra lives gained", strconv.Itoa(m.player.Stats.ExtraLivesGained)},
		{"Fewest turns for extra life", strconv.Itoa(m.player.Stats.FewestExtraLifeSolves)},
		{"Longest solve", fmt.Sprintf("%s (%d)", longest_solve, len(m.player.Stats.LongestSolve))},
		{"Average solve length", fmt.Sprintf("%.1f", m.player.Stats.AverageSolveLength())},
	}
		
	stats_table := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(m.renderer.NewStyle().Foreground(m.theme.Border())).
		Rows(stats...).
		StyleFunc(func(row, col int) lipgloss.Style {
			return m.theme.Base().Padding(0, 1)
		}).
		Render()

	return lipgloss.JoinVertical(
		lipgloss.Center,
		m.theme.TextRed().Bold(true).Render("GAME OVER\n"),
		stats_table,
	)
}
