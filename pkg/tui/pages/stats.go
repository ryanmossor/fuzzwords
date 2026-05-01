package pages

import (
	"fzwds/pkg/tui/figurethisout"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type StatsPage struct {
	name		PageName
	uiContext 	*figurethisout.UIContext
	helpKeys	[]figurethisout.HelpKeymap
}

func NewStatsPage(uiContext *figurethisout.UIContext) Page {
	return &StatsPage {
		name: 		Stats,
		uiContext: 	uiContext,
		helpKeys: 	[]figurethisout.HelpKeymap {
			{Key: "q", Value: "quit"},
		},
	}

}

func (p StatsPage) Switch() tea.Cmd {
	return nil
}

func (p StatsPage) GetPageName() PageName {
	return p.name
}

func (p StatsPage) Update(msg tea.Msg) (Page, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "m":
			return p, SwitchPageCmd(Title)
		case "a":
			return p, SwitchPageCmd(About)
		case "q":
			return p, tea.Quit
		}
	}
    return p, nil
}

func (p StatsPage) View() string {
	return lipgloss.NewStyle().
		Height(p.uiContext.ContentHeight).
		AlignVertical(lipgloss.Center).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Center,
				"🚧 Under construction 🚧",
	))
}
