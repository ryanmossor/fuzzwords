package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type settingsState struct {
	selected		int
}

func (m model) SettingsSwitch() (model, tea.Cmd) {
	m = m.SwitchPage(settings_page)
	m.footer_cmds = []footerCmd{
		{key: "↑/↓", value: "scroll"},
		{key: "←/→", value: "change"},
		{key: "ctrl+r", value: "restore defaults"},
		{key: "enter", value: "save"},
	}
	return m, nil
}

func (m model) SettingsUpdate(msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down", "tab":
			// TODO: next setting
		case "k", "up", "shift+tab":
			// TODO: prev setting
		case "+", "=", "right", "l":
			// TODO: increase setting if applicable
			// if m.IsCartEmpty() {
			// 	return m, nil
			// }
			// productVariantID := m.VisibleCartItems()[m.state.cart.selected].ProductVariantID
			// return m.UpdateCart(productVariantID, 1)
		case "-", "left", "h": 
			// TODO: decrease setting if applicable
		case "ctrl+r":
			// TODO: reset settings to defaults
		case "b":
			// TODO: beginner preset
		case "m":
			// TODO: medium preset
		case "d":
			// TODO: difficult preset
		case "x":
			// TODO: expert preset
		case "enter":
			// TODO: save settings
		case "esc":
			return m.MainMenuSwitch()
		}
	}

	return m, nil
}

func (m model) SettingsView() string {
	// base := m.theme.Base().Width(m.widthContent).Render
	base := m.theme.Base().Render
	dim := m.theme.TextDim().Render
	accent := m.theme.TextAccent().Render 

	var lines []string
	lines = append(lines, base("Change your ") + accent("game settings") + base(" here.\n"))
	lines = append(lines, accent("Alphabet: "))
	lines = append(lines, dim("Easy / ") + accent("Medium") + dim(" / Full"))

	return lipgloss.JoinVertical(
		lipgloss.Left,
		lines...,
		// base("Change your game settings here."),
	)
}
