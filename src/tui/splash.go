package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
) 

func (m model) MainMenuSwitch() (model, tea.Cmd) {
	m = m.SwitchPage(splash_page)
	return m, nil
}

func (m model) MainMenuUpdate(msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return m.GameSwitch()
		case "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) MainMenuView() string {
	// base := m.theme.Base().Width(m.widthContent).Render
	base := m.theme.Base().Render
	accent := m.theme.TextAccent().Render 

	return lipgloss.JoinVertical(
		lipgloss.Center,

		// accent("   ,d8888b                        d8b         "),
		// accent("   88P'                           88P         "),
		// accent("d888888P                         d88          "),
		// accent("  ?88'd88888P ?88   d8P  d8P d888888   .d888b,"),
		// accent("  88P    d8P' d88  d8P' d8P'd8P' ?88   ?8b,   "),
		// accent(" d88   d8P'   ?8b ,88b ,88' 88b  ,88b    `?8b "),
		// accent("d88'  d88888P'`?888P'888P'  `?88P'`88b`?888P' "),

		// accent("   ___                         __            "),
		// accent(" /'___\\                       /\\ \\           "),
		// accent("/\\ \\__/  ____    __  __  __   \\_\\ \\    ____  "),
		// accent("\\ \\ ,__\\/\\_ ,`\\ /\\ \\/\\ \\/\\ \\  /'_` \\  /',__\\ "),
		// accent(" \\ \\ \\_/\\/_/  /_\\ \\ \\_/ \\_/ \\/\\ \\L\\ \\/\\__, `\\"),
		// accent("  \\ \\_\\   /\\____\\\\ \\___x___/'\\ \\___,_\\/\\____/"),
		// accent("   \\/_/   \\/____/ \\/__//__/   \\/__,_ /\\/___/ "),

		// "",

		// accent(" ___ __       __   __  "),
		// accent("|__   / |  | |  \\ /__` "),
		// accent("|    /_ |/\\| |__/ .__/ "),

		// "",
											
		accent("███████╗███████╗██╗    ██╗██████╗ ███████╗"),
		accent("██╔════╝╚══███╔╝██║    ██║██╔══██╗██╔════╝"),
		accent("█████╗    ███╔╝ ██║ █╗ ██║██║  ██║███████╗"),
		accent("██╔══╝   ███╔╝  ██║███╗██║██║  ██║╚════██║"),
		accent("██║     ███████╗╚███╔███╔╝██████╔╝███████║"),
		accent("╚═╝     ╚══════╝ ╚══╝╚══╝ ╚═════╝ ╚══════╝"),

		"",

		"\n\nPress " + accent("ENTER") + base(" to play"),
	)
}
