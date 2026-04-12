package tui

import (
	"fmt"
	"fzwds/src/tui/animations"
	"fzwds/src/utils"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

func (m model) GameOverSwitch() (model, tea.Cmd) {
	// Briefly prevent key presses on game over screen
	cmds := []tea.Cmd{ m.debounceInputCmd(500) }

	m = m.SwitchPage(game_over_page)

	m.footer_keymaps = []FooterKeymap {
		{key: m.theme.Base().Foreground(m.theme.black).Background(m.theme.blue).Bold(true).Render(" r review "), value: ""},
		{key: "enter", value: "new game"},
		{key: "m", value: "main menu"},
        {key: "s", value: "change settings"},
		{key: "q", value: "quit"},
	}

	if m.state.game.GameWon {
		m.anim_mgr.InitAnimations(animations.GameOverWin)
	}

	m.state.game_ui.player_damaged = false

	if m.state.game.GameWon {
		m.state.game_ui.validation_msg = ""
		m.state.game_ui.game_over_msg = "===== YOU WIN! ====="
		m.anim_mgr.InitAnimations(animations.GameOverWin)
	} else {
		red := m.theme.TextRed()
		m.state.game_ui.validation_msg = red.Render(fmt.Sprintf(
			"Possible solve for final prompt %s: ",
			strings.ToUpper(m.state.game.CurrentTurn().Prompt)))
		m.state.game_ui.validation_msg += m.highlightPromptAnswer(
			m.state.game.CurrentTurn().Prompt,
			m.state.game.CurrentTurn().SourceWord,
			m.state.game.Settings.PromptMode)

		m.state.game_ui.game_over_msg = red.Bold(true).Render("☠️ GAME OVER ☠️")
	}

	if !m.state.game.EarlyQuit && !m.state.game.GameWon {
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
			m.anim_mgr.DeactivateAnimations(animations.GameOverWin)
			return m.MainMenuSwitch()
		case "s":
			m.anim_mgr.DeactivateAnimations(animations.GameOverWin)
			return m.SettingsSwitch(game_settings)
		case "r":
			return m.GameReviewSwitch()
		case "enter":
			m.anim_mgr.DeactivateAnimations(animations.GameOverWin)
			return m.GameSwitch()
		case "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) GameOverView() string {
	stats := m.state.game.Player.Stats

	var longest_streak string
	if stats.LongestStreak == 0 {
		longest_streak = "-"
	} else {
		longest_streak = fmt.Sprintf("%d words", stats.LongestStreak)
	}

	var longest_solve, longest_count string
	if stats.LongestSolve == "" {
		longest_solve = "-"
	} else {
		longest_solve = fmt.Sprintf("%s", stats.LongestSolve)
		longest_count = fmt.Sprintf("(%d)", len(stats.LongestSolve))
	}

	var most_unique_solve, most_unique_count string
	if stats.MostUniqueWord == "" {
		most_unique_solve = "-"
	} else {
		most_unique_solve = fmt.Sprintf("%s", stats.MostUniqueWord)
		most_unique_count = fmt.Sprintf("(%d)", stats.MostUniqueCount)
	}

	fastest_extra_life := fmt.Sprintf("%d turns", stats.FewestExtraLifeSolves)
	if stats.FewestExtraLifeSolves == 0 {
		fastest_extra_life = "-"
	}

	solves_per_min := "0"
    if stats.PromptsSolved > 0 {
		spm := float64(stats.PromptsSolved) / (float64(stats.TimePlayed) / 60.0)
		solves_per_min = fmt.Sprintf("%.1f", spm)
	}

	rows := [][]string {
		{"Time played", utils.FormatTime(stats.TimePlayed)},
		{"Prompts solved", strconv.Itoa(stats.PromptsSolved)},
		{"Prompts failed", strconv.Itoa(stats.PromptsFailed)},
		{"Solves per minute", solves_per_min},
		{"Longest streak", longest_streak},
		{"Average solve length", fmt.Sprintf("%.1f letters", stats.AverageSolveLength)},
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

	game_over_msg, changed := m.anim_mgr.ApplyAnimations(
		string(animations.GameOverWin),
		m.state.game_ui.game_over_msg)
	if !changed {
		game_over_msg = m.theme.TextYellow().Bold(true).Render(game_over_msg)
	}

	return lipgloss.JoinVertical(
		lipgloss.Center,
		game_over_msg,
		"",
		stats_table,
		"",
        m.state.game_ui.validation_msg,
	)
}
