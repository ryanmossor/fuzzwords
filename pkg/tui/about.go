package tui

import (
	"fzwds/pkg/enums"
	"fzwds/pkg/tui/pages"
	"fzwds/pkg/tui/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m model) AboutSwitch() (model, tea.Cmd) {
	m = m.SwitchPage(pages.AboutPage)
	m.footerKeymaps = []footerKeymap{
		{key: "q", value: "quit"},
	}
	return m, nil
}

func (m model) AboutUpdate(msg tea.Msg) (model, tea.Cmd) {
	return m, nil
}

func (m model) AboutView() string {
	base := styles.TextBody.Render
	accent := styles.TextAccent.Render
	yellow_bold := styles.TextYellow.Bold(true).Render

	return lipgloss.JoinVertical(
		lipgloss.Left,
		// TODO: hyperlink styling after v2 upgrade
		accent("Fuzzwords") + base(" is a word game inspired by ") + accent("BombParty: https://jklm.fun/"),
		"",
		base("In BombParty, players respond to a prompt (a sequence of letters) by typing " +
			"a word containing those letters in ") + accent("consecutive order") + base("."),
		"",
		yellow_bold(" - Example: ") +
		m.highlightPromptAnswer("RWO", "OVERWORK", enums.PromptModeClassic) +
		base(" solves the prompt ") + accent("RWO") + base(", but ") +
		m.highlightPromptAnswer("RWO", "REWROTE", enums.PromptModeClassic) +
		base(" does not"),
		"",
		base("Fuzzwords allows for ") + accent("\"fuzzy\" matching") +
		base(", meaning the letters of the prompt must still be used in the same order " +
			"as in the prompt, but they do not need to be consecutive."),
		"",
		yellow_bold(" - Example: ") +
		m.highlightPromptAnswer("RWO", "OVERWORK", enums.PromptModeFuzzy) +
		base(" and ") +
		m.highlightPromptAnswer("RWO", "REWROTE", enums.PromptModeFuzzy) +
		base(" both solve ") + accent("RWO") + base(", but ") +
		accent("WARRIOR") + base(" does not"),

		// TODO: rules on extra lives, game modes (endless/max lives), etc
		// TODO: scrollbar and/or rule pages
	)
}
