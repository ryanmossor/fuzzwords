package tui

import (
	// "fzwds/pkg/tui/pages"
	"fzwds/pkg/tui/pages"
	"fzwds/pkg/tui/styles"
	"fzwds/pkg/tui/theme"

	// "strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

func (m model) HeaderUpdate(msg tea.Msg) (model, tea.Cmd) {
	// TODO: has_header flag
	// if m.game.GameActive() ||
	// m.page == pages.GameOverPage ||
	// m.page == pages.GameReviewPage ||
	// m.page == pages.SettingsPage ||
	// m.page == pages.PokemonGenMenuPage {
	// 	return m, nil
	// }

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		// case "s":
		// 	return m.StatsSwitch()
		// case "a":
		// 	return m.AboutSwitch()
		// case "m":
		// 	return m.TitleScreenSwitch()
		case "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) HeaderView() string {
	// if m.page == pages.GamePage || m.page == pages.GameOverPage {
	// 	return m.GameHudView()
	// }
	// if m.page == pages.GameReviewPage {
	// 	return m.GameReviewHudView()
	// }

	bold := styles.TextAccent.Bold(true).Render
	accent := styles.TextAccent.Render
	body := styles.TextBody.Render

	// TODO: entire header could probably be top line, custom text from model state, bottom line
	// if m.page == pages.SettingsPage || m.page == pages.PokemonGenMenuPage {
	// 	return lipgloss.JoinVertical(
	// 		lipgloss.Center,
	// 		styles.TextBorder.Render(strings.Repeat("─", m.containerWidth)),
	// 		styles.TextBlue.Bold(true).Render(m.state.settings.title),
	// 		styles.TextBorder.Render(strings.Repeat("─", m.containerWidth)))
	// }

	menu := accent("m") + body(" main menu")
	about := accent("a") + body(" about")
	stats := accent("s") + body(" stats")

	switch m.currentPage.GetPageName() {
	case pages.Title:
		menu = bold("m main menu")
	case pages.About:
		about = bold("a about")
	case pages.Stats:
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
		Width(m.uiContext.ContainerWidth).
		StyleFunc(func(row, col int) lipgloss.Style {
			return lipgloss.NewStyle().Padding(0, 1).AlignHorizontal(lipgloss.Center)
		}).
		Render()
}
