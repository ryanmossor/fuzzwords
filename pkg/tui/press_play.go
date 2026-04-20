package tui

import (
	"fzwds/pkg/tui/styles"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type PressPlayState struct {
	visible bool
}

var (
	hidden  = styles.TextBody.Render("Press       to play")
	visible = styles.TextBody.Render("Press ") +
			  styles.TextAccent.Bold(true).Render("ENTER") +
			  styles.TextBody.Render(" to play")
)

type PressPlayTickMsg struct {}
func (m model) pressPlayFlashCmd() tea.Cmd {
	if !m.app_settings.Prefs.AnimationsEnabled {
		return nil
	}
	return tea.Every(850 * time.Millisecond, func(t time.Time) tea.Msg {
		return PressPlayTickMsg{}
	})
}

func (m model) PressPlayInit() tea.Cmd {
	return m.pressPlayFlashCmd()
}

func (m model) PressPlayView() string {
	if !m.state.pressPlay.visible && m.app_settings.Prefs.AnimationsEnabled {
		return hidden
	}
	return visible
}
