package pages

import (
	"fzwds/pkg/enums"
	"fzwds/pkg/tui/figurethisout"
	"fzwds/pkg/tui/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type AboutPage struct {
	name		PageName
	uiContext 	*figurethisout.UIContext
	helpKeys	[]figurethisout.HelpKeymap
}

func NewAboutPage(uiContext *figurethisout.UIContext) Page {
	return &AboutPage {
		name: About,
		uiContext: uiContext,
		helpKeys: []figurethisout.HelpKeymap {
			{Key: "q", Value: "quit"},
		},
	}
}

func (p AboutPage) Switch() tea.Cmd {
	return nil
}

func (p AboutPage) GetPageName() PageName {
	return p.name
}

func (p AboutPage) Update(msg tea.Msg) (Page, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "m":
			return p, SwitchPageCmd(Title)
		case "s":
			return p, SwitchPageCmd(Stats)
		case "q":
			return p, tea.Quit
		}
	}

	return p, nil
}

func (p AboutPage) View() string {
	base := styles.TextBody.Render
	accent := styles.TextAccent.Render
	yellow_bold := styles.TextYellow.Bold(true).Render

	style := lipgloss.NewStyle().
		Width(p.uiContext.ContentWidth).
		PaddingTop(1)

	return style.Render(lipgloss.JoinVertical(
		lipgloss.Left,
		// TODO: hyperlink styling after v2 upgrade
		accent("Fuzzwords") + base(" is a word game inspired by ") + accent("BombParty: https://jklm.fun/"),
		"",
		base("In BombParty, players respond to a prompt (a sequence of letters) by typing " +
			"a word containing those letters in ") + accent("consecutive order") + base("."),
		"",
		yellow_bold(" - Example: ") +
		styles.HighlightPromptAnswer("RWO", "OVERWORK", enums.PromptModeClassic) +
		base(" solves the prompt ") + accent("RWO") + base(", but ") +
		styles.HighlightPromptAnswer("RWO", "REWROTE", enums.PromptModeClassic) +
		base(" does not"),
		"",
		base("Fuzzwords allows for ") + accent("\"fuzzy\" matching") +
		base(", meaning the letters of the prompt must still be used in the same order " +
			"as in the prompt, but they do not need to be consecutive."),
		"",
		yellow_bold(" - Example: ") +
		styles.HighlightPromptAnswer("RWO", "OVERWORK", enums.PromptModeFuzzy) +
		base(" and ") +
		styles.HighlightPromptAnswer("RWO", "REWROTE", enums.PromptModeFuzzy) +
		base(" both solve ") + accent("RWO") + base(", but ") +
		accent("WARRIOR") + base(" does not"),

		// TODO: rules on extra lives, game modes (endless/max lives), etc
		// TODO: scrollbar and/or rule pages
	))
}
