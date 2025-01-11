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

func (m model) GameOverSwitch(msg string, win bool) (model, tea.Cmd) {
	// TODO: are both of these flags needed?
	m.game_active = false
	m.game_over = true

	m.game_state.Player.Stats.ElapsedSeconds = int(time.Since(m.game_start_time).Seconds())

	m.game_over_msg = msg
	m = m.SwitchPage(game_over_page)

    if win {
        m.state.game.validation_msg = ""
    } else {
        m.state.game.validation_msg = fmt.Sprintf(
            "Possible answer for final prompt %s: %s",
            strings.ToUpper(m.game_state.CurrentTurn.Prompt),
            strings.ToUpper(m.game_state.CurrentTurn.SourceWord))
    }

	m.footer_cmds = []footerCmd{
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
		if m.state.game.restrict_input {
			return m, nil
		}

		switch msg.String() {
		case "m":
			m.game_over = false
			return m.MainMenuSwitch()
		case "s":
			m.game_over = false
			return m.SettingsSwitch()
		case "enter":
			m.game_over = false
			return m.GameSwitch()
		case "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) GameOverView() string {
	longest_solve := fmt.Sprintf("%s (%d)",
		m.game_state.Player.Stats.LongestSolve,
		len(m.game_state.Player.Stats.LongestSolve))

	if m.game_state.Player.Stats.LongestSolve == "" {
		longest_solve = "-"
	}

	fastest_extra_life := fmt.Sprintf("%d turns", m.game_state.Player.Stats.FewestExtraLifeSolves)
	if m.game_state.Player.Stats.FewestExtraLifeSolves == 0 {
		fastest_extra_life = "-"
	}

	var solves_per_min float64 = 0
    if m.game_state.Player.Stats.PromptsSolved > 0 {
        solves_per_min = float64(m.game_state.Player.Stats.PromptsSolved) / (float64(m.game_state.Player.Stats.ElapsedSeconds) / 60.0)
    }

	stats := [][]string{
		{"Time survived", utils.FormatTime(m.game_state.Player.Stats.ElapsedSeconds)},
		{"Prompts solved", strconv.Itoa(m.game_state.Player.Stats.PromptsSolved)},
		{"Prompts failed", strconv.Itoa(m.game_state.Player.Stats.PromptsFailed)},
		{"Solves per minute", fmt.Sprintf("%.1f", solves_per_min)},
		{"Average solve length", fmt.Sprintf("%.1f letters", m.game_state.Player.Stats.AverageSolveLength())},
		{"Longest word used", longest_solve},
		{"Extra lives gained", strconv.Itoa(m.game_state.Player.Stats.ExtraLivesGained)},
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

    validation_msg, _ := m.renderValidationMsg()

	return lipgloss.JoinVertical(
		lipgloss.Center,
		m.game_over_msg,
		"",
		stats_table,
		"",
        validation_msg,
	)
}
