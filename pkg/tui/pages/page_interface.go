package pages

import tea "github.com/charmbracelet/bubbletea"

type PageName string
const (
	About 			PageName = "about"
	GameOver 		PageName = "game over"
	Game 			PageName = "game"
	GameReview 		PageName = "review"
	PokemonGenMenu 	PageName = "pokemon gens"
	Preferences		PageName = "preferences"
	Settings 		PageName = "settings"
    Stats 			PageName = "stats"
	Title 			PageName = "title screen"
)

type Page interface {
	GetPageName() PageName
	Switch() tea.Cmd
	Update(msg tea.Msg) (Page, tea.Cmd)
	View() string
}

type SwitchPageMsg struct {
	PageName	PageName
}
func SwitchPageCmd(pageName PageName) tea.Cmd {
	return func() tea.Msg {
		return SwitchPageMsg{ PageName: pageName }
	}
}
