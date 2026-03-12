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

func (m model) GameOverSwitch(win, early_quit bool) (model, tea.Cmd) {
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
			"Possible solve for final prompt %s: ",
			strings.ToUpper(m.state.game.CurrentTurn.Prompt)))

		m.state.game_ui.validation_msg += m.highlightPromptAnswer(
			m.state.game.CurrentTurn.Prompt,
			m.state.game.CurrentTurn.SourceWord,
			m.state.game.Settings.PromptMode)

        m.state.game_ui.game_over_msg = red.Bold(true).Render("☠️ GAME OVER ☠️")
    }

	m = m.SwitchPage(game_over_page)

	m.footer_keymaps = []footer_keymaps{
		{key: "m", value: "main menu"},
        {key: "s", value: "change settings"},
		{key: "enter", value: "new game"},
		{key: "q", value: "quit"},
	}

	// Briefly prevent key presses on game over screen
	cmds := []tea.Cmd{ m.debounceInputCmd(500) }
	if !early_quit && !win {
		cmds = append(cmds, m.terminalBellCmd(false))
	}

	return m, tea.Batch(cmds...)
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
	stats := m.state.game.Player.Stats

	var longest_solve, longest_count string
	if stats.LongestSolve == "" {
		longest_solve = "-"
	} else {
		longest_solve = fmt.Sprintf("%s", stats.LongestSolve)
		longest_count = fmt.Sprintf("(%d)", len(stats.LongestSolve))
	}

	var most_unique_solve, most_unique_count string
	if stats.MostUniqueLetters == "" {
		most_unique_solve = "-"
	} else {
		most_unique_solve = fmt.Sprintf("%s", stats.MostUniqueLetters)
		most_unique_count = fmt.Sprintf("(%d)", utils.CountUniqueLetters(stats.MostUniqueLetters))
	}

	fastest_extra_life := fmt.Sprintf("%d turns", stats.FewestExtraLifeSolves)
	if stats.FewestExtraLifeSolves == 0 {
		fastest_extra_life = "-"
	}

	solves_per_min := "0"
    if stats.PromptsSolved > 0 {
		spm := float64(stats.PromptsSolved) / (float64(stats.ElapsedSeconds) / 60.0)
		solves_per_min = fmt.Sprintf("%.1f", spm)
	}

	rows := [][]string{
		{"Time survived", utils.FormatTime(stats.ElapsedSeconds)},
		{"Prompts solved", strconv.Itoa(stats.PromptsSolved)},
		{"Prompts failed", strconv.Itoa(stats.PromptsFailed)},
		{"Solves per minute", solves_per_min},
		{"Average solve length", fmt.Sprintf("%.1f letters", stats.AverageSolveLength())},
		{"Longest word used", longest_solve, longest_count},
		{"Most unique letters", most_unique_solve, most_unique_count},
		{"Extra lives gained", strconv.Itoa(stats.ExtraLivesGained)},
		{"Fastest extra life", fastest_extra_life},
	}
		
	stats_table := table.New().
		Border(lipgloss.HiddenBorder()).
		BorderColumn(false).
		BorderStyle(m.renderer.NewStyle().Foreground(m.theme.Border())).
		Rows(rows...).
		StyleFunc(func(row, col int) lipgloss.Style {
			var style lipgloss.Style

			if row % 2 == 0 {
				style = m.theme.TextAccent()
			} else {
				style = m.theme.Base()
			}

			if col == 0 && stats.PromptsSolved > 0 {
				// Pad 1st col to offset extra width of 3rd col (counts for longest/most unique words)
				// 3rd col only populated if at least 1 prompt was solved
				style = style.PaddingLeft(len(longest_count) + 1)
			} else if col == 1 {
				style = style.
					Align(lipgloss.Right).
					PaddingRight(1).
					PaddingLeft(5)
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
