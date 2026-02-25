package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type pressPlayState struct {
	visible bool
}

type PressPlayTickMsg struct {}
type LogoInitMsg struct{}

func pressPlayFlashCmd() tea.Cmd {
	return tea.Every(700 * time.Millisecond, func(t time.Time) tea.Msg {
		return PressPlayTickMsg{}
	})
}

func (m model) PressPlayInit() tea.Cmd {
	return pressPlayFlashCmd()
}

func (m model) PressPlayUpdate(msg tea.Msg) (model, tea.Cmd) {
	switch msg.(type) {
	case PressPlayTickMsg:
		m.state.press_play.visible = !m.state.press_play.visible
		return m, pressPlayFlashCmd()
	}
	return m, nil
}

func (m model) PressPlayView() string {
	if !m.state.press_play.visible {
		return ""
	}

	base := m.theme.Base().Render
	accent := m.theme.TextAccent().Bold(true).Render 
	return base("Press ") + accent("ENTER") + base(" to play")
}
