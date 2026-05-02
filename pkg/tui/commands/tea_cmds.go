package commands

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
func TickCmd(fps int) tea.Cmd {
	return tea.Tick(time.Second / time.Duration(fps), func(t time.Time) tea.Msg {
		return TickMsg{t}
	})
}

type EnableInputMsg time.Time
func DebounceInputCmd(duration_ms int) tea.Cmd {
    // m.state.game.inputRestricted = true
    return tea.Tick(time.Millisecond * time.Duration(duration_ms), func(t time.Time) tea.Msg {
		return EnableInputMsg(t)
	})
}

type TogglePlayerDamagedMsg struct{}
func TogglePlayerDamagedCmd() tea.Cmd {
    return tea.Tick(time.Millisecond * time.Duration(500), func(t time.Time) tea.Msg {
		return TogglePlayerDamagedMsg{}
	})
}

func TerminalBellCmd(prefs game.GeneralPreferences, force bool) tea.Cmd {
	if force || prefs.BellEnabled {
		return func() tea.Msg {
			fmt.Fprint(os.Stdout, "\a") // send BEL character
			return nil
		}
	}
	return nil
}

func SaveSettingsCmd(settings game.Settings, path string) tea.Cmd {
	return func() tea.Msg {
		settings.WriteSettings(path)
		return nil
	}
}
