package tui

import (
	"fmt"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

func (m model) GameOverSwitch(game_over_msg string) (model, tea.Cmd) {
	m.game_over_msg = game_over_msg
	m = m.SwitchPage(game_over_page)

	m.footer_cmds = []footerCmd{
		{key: "m", value: "main menu"},
		{key: "enter", value: "new game"},
		{key: "q", value: "quit"},
	}

	return m, nil
}

func (m model) GameOverUpdate(msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "m":
			return m.MainMenuSwitch()
		case "enter":
			return m.GameSwitch()
		case "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) GameOverView() string {
	longest_solve := m.player.Stats.LongestSolve
	if longest_solve == "" {
		longest_solve = "-"
	}

	fastest_extra_life := fmt.Sprintf("%d turns", m.player.Stats.FewestExtraLifeSolves)
	if m.player.Stats.FewestExtraLifeSolves == 0 {
		fastest_extra_life = "-"
	}

	stats := [][]string{
		{"Prompts solved", strconv.Itoa(m.player.Stats.PromptsSolved)},
		{"Prompts failed", strconv.Itoa(m.player.Stats.PromptsFailed)},
		{"Average solve length", fmt.Sprintf("%.1f letters", m.player.Stats.AverageSolveLength())},
		{"Longest word used", fmt.Sprintf("%s (%d)", longest_solve, len(m.player.Stats.LongestSolve))},
		{"Extra lives gained", strconv.Itoa(m.player.Stats.ExtraLivesGained)},
		{"Fastest extra life", fastest_extra_life},
	}
		
	stats_table := table.New().
		Border(lipgloss.HiddenBorder()).
		BorderColumn(false).
		BorderStyle(m.renderer.NewStyle().Foreground(m.theme.Border())).
		Rows(stats...).
		StyleFunc(func(row, col int) lipgloss.Style {
			var style lipgloss.Style

			if row % 2 == 0 {
				style = m.theme.TextAccent().Padding(0, 1)
			} else {
				style = m.theme.Base().Padding(0, 1)
			}

			if col == 1 {
				style = style.Align(lipgloss.Right).Padding(0, 1, 0, 2)
			}
			return style
		}).
		Render()

	return lipgloss.JoinVertical(
		lipgloss.Center,
		m.game_over_msg,
		"",
		stats_table,
	)
}
