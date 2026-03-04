package tui

import (
	"fmt"
	"fzwds/src/utils"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

func (m model) GameOverSwitch(win bool) (model, tea.Cmd) {
	m.state.game_ui.game_active = false
	m.state.game_ui.player_damaged = false
	m.state.game.Player.Stats.ElapsedSeconds = int(time.Since(m.state.game_ui.start_time).Seconds())

	red := m.theme.TextRed()
	green := m.theme.TextGreen().Bold(true)

    if win {
        m.state.game_ui.validation_msg = ""
        m.state.game_ui.game_over_msg = green.Render("===== YOU WIN! =====")
    } else {
		m.state.game_ui.validation_msg = red.Render(fmt.Sprintf(
			"Possible answer for final prompt %s: ",
			strings.ToUpper(m.state.game.CurrentTurn.Prompt)))
		m.state.game_ui.validation_msg += m.colorizeInput(m.state.game.CurrentTurn.SourceWord)
        m.state.game_ui.game_over_msg = red.Bold(true).Render("===== GAME OVER =====")
    }

	m = m.SwitchPage(game_over_page)

	m.footer_keymaps = []footer_keymaps{
		{key: "m", value: "main menu"},
        {key: "s", value: "change settings"},
		{key: "enter", value: "new game"},
		{key: "q", value: "quit"},
	}

	// Briefly prevent key presses on game over screen
	return m, m.debounceInputCmd(500)
}

func (m model) GameOverUpdate(msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.state.game_ui.input_restricted {
			return m, nil
		}

		switch msg.String() {
		case "m":
			return m.MainMenuSwitch()
		case "s":
			return m.SettingsSwitch()
		case "enter":
			return m.GameSwitch()
		case "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) GameOverView() string {
	longest_solve := fmt.Sprintf("%s (%d)",
		m.state.game.Player.Stats.LongestSolve,
		len(m.state.game.Player.Stats.LongestSolve))

	if m.state.game.Player.Stats.LongestSolve == "" {
		longest_solve = "-"
	}

	fastest_extra_life := fmt.Sprintf("%d turns", m.state.game.Player.Stats.FewestExtraLifeSolves)
	if m.state.game.Player.Stats.FewestExtraLifeSolves == 0 {
		fastest_extra_life = "-"
	}

	var solves_per_min float64 = 0
    if m.state.game.Player.Stats.PromptsSolved > 0 {
        solves_per_min = float64(m.state.game.Player.Stats.PromptsSolved) / (float64(m.state.game.Player.Stats.ElapsedSeconds) / 60.0)
    }

	stats := [][]string{
		{"Time survived", utils.FormatTime(m.state.game.Player.Stats.ElapsedSeconds)},
		{"Prompts solved", strconv.Itoa(m.state.game.Player.Stats.PromptsSolved)},
		{"Prompts failed", strconv.Itoa(m.state.game.Player.Stats.PromptsFailed)},
		{"Solves per minute", fmt.Sprintf("%.1f", solves_per_min)},
		{"Average solve length", fmt.Sprintf("%.1f letters", m.state.game.Player.Stats.AverageSolveLength())},
		{"Longest word used", longest_solve},
		{
			"Solve w/ most unique letters",
			fmt.Sprintf("%s (%d)",
				m.state.game.Player.Stats.MostUniqueLetters,
				utils.CountUniqueLetters(m.state.game.Player.Stats.MostUniqueLetters)),
		},
		{"Extra lives gained", strconv.Itoa(m.state.game.Player.Stats.ExtraLivesGained)},
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

    validation_msg := m.renderValidationMsg()

	return lipgloss.JoinVertical(
		lipgloss.Center,
		m.state.game_ui.game_over_msg,
		"",
		stats_table,
		"",
        validation_msg,
	)
}
