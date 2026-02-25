package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m model) AboutSwitch() (model, tea.Cmd) {
	m = m.SwitchPage(about_page)
	m.footer_cmds = []footerCmd{
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
	highlight := m.theme.TextHighlight().Render
	lavender_bold := m.theme.TextLavender().Bold(true).Render

	return lipgloss.JoinVertical(
		lipgloss.Left,
		accent("Fuzzwords") + base(" is a word game inspired by ") + accent("BombParty: https://jklm.fun/"),
		"",
		base("In BombParty, players respond to a prompt (a sequence of letters) by typing " +
			"a word containing those letters in ") + accent("consecutive order") + base("."),
		"",
		lavender_bold(" ▶ Example:") +
		accent(" OVE") + highlight("RWO") + accent("RKED") +
		base(" solves the prompt ") + accent("RWO"),
		"",
		base("Fuzzwords allows for ") +
		accent("\"fuzzy\" matching") +
		base(", meaning the letters of the prompt must still be used in the same order " +
			"as in the prompt, but they do not need to be consecutive."),
		"",
		lavender_bold(" ▶ Example:") +
		accent(" IN") + highlight("V") + accent("ES") + highlight("TM") + accent("ENT") +
		base(" and ") +
		highlight("V") + accent("EN") + highlight("T") + accent("RILOQUIS") + highlight("M") +
		base(" both solve the prompt ") + accent("VTM"),

		// TODO: rules on extra lives, game modes (endless/max lives), etc
		// TODO: scrollbar and/or rule pages
	)
}
