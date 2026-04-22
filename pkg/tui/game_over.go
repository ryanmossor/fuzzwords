package tui

import (
	"fmt"
	"fzwds/pkg/tui/animations"
	"fzwds/pkg/tui/pages"
	"fzwds/pkg/tui/styles"
	"fzwds/pkg/tui/theme"
	"fzwds/pkg/utils"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

type gameOverState struct {
	viewCache		map[string]string
}

func (m model) GameOverSwitch() (model, tea.Cmd) {
	m = m.SwitchPage(pages.GameOverPage)

	m.footerKeymaps = []footerKeymap {
		{key: "r", value: "review"},
		{key: "enter", value: "new game"},
		{key: "s", value: "settings"},
		{key: "m", value: "main menu"},
		{key: "q", value: "quit"},
	}

	if m.state.game.gameWon {
		m.animManager.InitAnimations(animations.GameOverWin)
	}

	return m, nil
}

func (m model) GameOverUpdate(msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// TODO can this check be in main model? or maybe just keep here since i want to refactor pages/components
		if m.state.game.inputRestricted {
			return m, nil
		}

		switch msg.String() {
		case "m":
			m.animManager.DeactivateAnimations(animations.GameOverWin)
			return m.TitleScreenSwitch()
		case "s":
			m.animManager.DeactivateAnimations(animations.GameOverWin)
			return m.SettingsSwitch(gameSettings)
		case "r":
			m.animManager.DeactivateAnimations(animations.GameOverWin)
			return m.GameReviewSwitch()
		case "enter":
			m.animManager.DeactivateAnimations(animations.GameOverWin)
			return m.GameSwitch()
		case "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m *model) GameOverView() string {
	if m.state.gameOver.viewCache["fullView"] != "" {
		return m.state.gameOver.viewCache["fullView"]
	}

	view := lipgloss.JoinVertical(
		lipgloss.Center,
		m.renderGameOverTitleMsg(),
		"",
		m.renderGameOverStatTable(),
	)

	if !m.state.game.gameWon || !m.settings.Prefs.AnimationsEnabled {
		m.state.gameOver.viewCache["fullView"] = view
	}

	return view
}

func (m *model) renderGameOverTitleMsg() string {
	if cached := m.state.gameOver.viewCache["title"]; cached != "" {
		return cached
	}

	var title string

	if m.state.game.gameWon {
		var changed bool
		title, changed = m.animManager.ApplyAnimations(
			string(animations.GameOverWin),
			"===== YOU WIN! =====")

		if !changed {
			title = styles.TextYellow.Bold(true).Render(title)
		}

		// Don't cache animated title
		if m.settings.Prefs.AnimationsEnabled {
			return title
		}
	} else {
		title = styles.TextRed.Bold(true).Render("☠️ GAME OVER ☠️")
	}

	m.state.gameOver.viewCache["title"] = title

	return title
}

func (m *model) renderGameOverStatTable() string {
	if m.state.gameOver.viewCache["stats"] != "" {
		return m.state.gameOver.viewCache["stats"]
	}

	stats := m.state.game.stats

	var longest_streak string
	if stats.LongestStreak() == 0 {
		longest_streak = "-"
	} else {
		longest_streak = fmt.Sprintf("%d words", stats.LongestStreak())
	}

	var longest_solve, longest_count string
	if stats.LongestSolve() == "" {
		longest_solve = "-"
	} else {
		longest_solve = stats.LongestSolve()
		longest_count = fmt.Sprintf("(%d)", len(stats.LongestSolve()))
	}

	var most_unique_solve, most_unique_count string
	if stats.MostUniqueWord() == "" {
		most_unique_solve = "-"
	} else {
		most_unique_solve = stats.MostUniqueWord()
		most_unique_count = fmt.Sprintf("(%d)", stats.MostUniqueCount())
	}

	fastest_extra_life := fmt.Sprintf("%d turns", stats.FewestExtraLifeSolves())
	if stats.FewestExtraLifeSolves() == 0 {
		fastest_extra_life = "-"
	}

	solves_per_min := "0"
    if stats.PromptsSolved() > 0 {
		solves_per_min = fmt.Sprintf("%.1f", stats.SolvesPerMinute())
	}

	rows := [][]string {
		{"Time played", utils.FormatTime(stats.TimePlayed())},
		{"Prompts solved", fmt.Sprintf("%d / %d", stats.PromptsSolved(), len(m.state.gameReview.turns))},
		{"Solves per minute", solves_per_min},
		{"Longest streak", longest_streak},
		{"Average solve length", fmt.Sprintf("%.1f letters", stats.AverageSolveLength())},
		{"Longest word used", longest_solve, longest_count},
		{"Most unique letters", most_unique_solve, most_unique_count},
		{"Extra lives gained", strconv.Itoa(stats.ExtraLivesGained())},
		{"Fastest extra life", fastest_extra_life},
	}

	stats_table := table.New().
		Border(lipgloss.HiddenBorder()).
		BorderColumn(false).
		BorderStyle(lipgloss.NewStyle().Foreground(theme.Border)).
		Rows(rows...).
		StyleFunc(func(row, col int) lipgloss.Style {
			var style lipgloss.Style

			if row % 2 == 0 {
				style = styles.TextAccent
			} else {
				style = styles.TextBody
			}

			if col == 0 && stats.PromptsSolved() > 0 {
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

	lines := []string{ stats_table }

	if !m.state.game.gameWon {
		msg := styles.TextRed.Render(fmt.Sprintf(
			"Possible solve for final prompt %s: ",
			strings.ToUpper(m.state.game.turn.prompt)))
		msg += m.highlightPromptAnswer(
			m.state.game.turn.prompt,
			m.state.game.possibleFinalAnswer,
			m.game.Settings().PromptMode)

		lines = append(lines, "", msg)
	}

	view := lipgloss.JoinVertical(lipgloss.Center, lines...)
	m.state.gameOver.viewCache["stats"] = view

	return view
}
