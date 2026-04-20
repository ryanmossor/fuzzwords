package tui

import (
	"fmt"
	"fzwds/pkg/game"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type TickMsg struct {
	Time	time.Time
}
// Global tick timer
func (m model) tickCmd() tea.Cmd {
	return tea.Tick(time.Second / time.Duration(m.fps), func(t time.Time) tea.Msg {
		return TickMsg{t}
	})
}

type EnableInputMsg time.Time
func (m *model) debounceInputCmd(duration_ms int) tea.Cmd {
    m.state.game.inputRestricted = true
    return tea.Tick(time.Millisecond * time.Duration(duration_ms), func(t time.Time) tea.Msg {
		return EnableInputMsg(t)
	})
}

type TogglePlayerDamagedMsg struct{}
func (m *model) togglePlayerDamagedCmd() tea.Cmd {
    return tea.Tick(time.Millisecond * time.Duration(500), func(t time.Time) tea.Msg {
		return TogglePlayerDamagedMsg{}
	})
}

func (m model) terminalBellCmd(force bool) tea.Cmd {
	if force || m.settings.Prefs.BellEnabled {
		return func() tea.Msg {
			// Send BEL character
			fmt.Fprint(os.Stdout, "\a")
			return nil
		}
	}
	return nil
}

func(m model) saveSettingsCmd(settings game.Settings, path string) tea.Cmd {
	return func() tea.Msg {
		writeSettings(settings, path)
		return nil
	}
}
