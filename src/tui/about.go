package tui

import (
	"fzwds/src/enums"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func (m model) AboutSwitch() (model, tea.Cmd) {
	m = m.SwitchPage(about_page)
	m.footer_keymaps = []FooterKeymap{
		{key: "q", value: "quit"},
	}
	return m, nil
}

func (m model) AboutUpdate(msg tea.Msg) (model, tea.Cmd) {
	return m, nil
}

func (m model) AboutView() string {
	base := m.theme.Base().Render
	accent := m.theme.TextAccent().Render
	yellow_bold := m.theme.TextYellow().Bold(true).Render

	return lipgloss.JoinVertical(
		lipgloss.Left,
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
