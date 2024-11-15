package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m model) AboutSwitch() (model, tea.Cmd) {
	m = m.SwitchPage(about_page)
	m.footerCmds = []footerCmd{
		{key: "esc", value: "main menu"},
		{key: "s", value: "settings"},
		{key: "q", value: "quit"},
	}
	return m, nil
}

func (m model) AboutUpdate(msg tea.Msg) (model, tea.Cmd) {
	return m, nil
}

func (m model) AboutView() string {
	// base := m.theme.Base().Width(m.widthContent).Render
	// bold := m.theme.TextAccent().Bold(true).Render
	// bold := m.theme.TextAccent().Bold(true).Render

	base := m.theme.Base().Render
	// accent := m.theme.TextAccent().Bold(true).Render 
	accent := m.theme.TextAccent().Render 

	// first_line_accent := m.theme.TextAccent().Render 
	// accent := m.theme.TextAccent().Width(m.widthContent).Render

	return lipgloss.JoinVertical(
		lipgloss.Left,
		accent("Fuzzwords") + base(" is a word game inspired by ") + accent("BombParty: https://jklm.fun/"),
		"",
		"BombParty challenges players with a prompt that must be fulfilled by a word containing the letters of the prompt in " + 
			accent("consecutive order") +
			base(" (e.g., the word OVE") + 
			accent("RWO") +
			base("RKED satisfies the prompt ") +
			accent("RWO") +
			base(")."),
		"",
		"Fuzzwords allows for " + 
			accent("\"fuzzy\" matching") +
			base(", meaning the letters of the prompt must appear in the answer in the given order, but are not required to be consecutive."),
		"",
		"For example, words like IN" + accent("V") + base("ES") + accent("TM") + base("ENT and ") +
			accent("V") + base("EN") + accent("T") + base("RILOQUIS") + accent("M") +
			base(" satisfy the prompt ") +
			accent("VTM") +
			base(". However, MOTIVE does not, as the letters are not in the correct order."),
		// TODO: rules on extra lives, game modes (endless/LMS), etc
	)
}
