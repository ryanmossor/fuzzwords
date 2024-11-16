package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type pressPlayState struct {
	visible bool
}

type PressPlayTickMsg struct {}

func (m model) PressPlayInit() tea.Cmd {
	return tea.Every(time.Millisecond * 700, func(t time.Time) tea.Msg {
		return PressPlayTickMsg{}
	})
}

func (m model) PressPlayUpdate(msg tea.Msg) (model, tea.Cmd) {
	switch msg.(type) {
	case PressPlayTickMsg:
		m.state.press_play.visible = !m.state.press_play.visible
		return m, tea.Every(time.Millisecond * 700, func(t time.Time) tea.Msg {
			return PressPlayTickMsg{}
		})
	}
	return m, nil
}

func (m model) PressPlayView() string {
	if !m.state.press_play.visible {
		return ""
	}

	accent := m.theme.TextAccent().Render
	green := m.theme.TextBlue().Bold(true).Render 
	return accent("Press ") + green("ENTER") + accent(" to play")
}
