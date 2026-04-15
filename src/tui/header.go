package tui

import (
	"fzwds/src/tui/styles"
	"fzwds/src/tui/theme"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

func (m model) HeaderUpdate(msg tea.Msg) (model, tea.Cmd) {
	// TODO: has_header flag
	if m.game.GameActive ||
	m.page == game_over_page ||
	m.page == game_review_page ||
	m.page == settings_page ||
	m.page == pokemon_gen_selector {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "s":
			return m.StatsSwitch()
		case "a":
			return m.AboutSwitch()
		case "m":
			return m.MainMenuSwitch()
		case "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) HeaderView() string {
	if m.page == game_page || m.page == game_over_page {
		return m.GameHudView()
	}
	if m.page == game_review_page {
		return m.GameReviewHudView()
	}

	bold := styles.TextAccent.Bold(true).Render
	accent := styles.TextAccent.Render
	body := styles.TextBody.Render

	// TODO: entire header could probably be top line, custom text from model state, bottom line
	if m.page == settings_page || m.page == pokemon_gen_selector {
		return lipgloss.JoinVertical(
			lipgloss.Center,
			styles.TextBorder.Render(strings.Repeat("─", m.width_container)),
			styles.TextBlue.Bold(true).Render(m.state.settings.title),
			styles.TextBorder.Render(strings.Repeat("─", m.width_container)))
	}

	menu := accent("m") + body(" main menu")
	about := accent("a") + body(" about")
	stats := accent("s") + body(" stats")

	switch m.page {
	case splash_page:
		menu = bold("m main menu")
	case about_page:
		about = bold("a about")
	case stats_page:
		stats = bold("s stats")
	}

	tabs := []string{
		menu,
		about,
		stats,
	}

	return table.New().
		Border(lipgloss.NormalBorder()).
		BorderLeft(false).
		BorderRight(false).
		BorderColumn(false).
		BorderStyle(lipgloss.NewStyle().Foreground(theme.Border)).
		Row(tabs...).
		Width(m.width_container).
		StyleFunc(func(row, col int) lipgloss.Style {
			return lipgloss.NewStyle().Padding(0, 1).AlignHorizontal(lipgloss.Center)
		}).
		Render()
}
