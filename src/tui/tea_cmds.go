package tui

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type TurnTimerExpiredMsg struct{}
func (m model) turnTimerExpiredCmd() tea.Cmd {
	return func() tea.Msg {
		return TurnTimerExpiredMsg{}
	}
}

type EnableInputMsg time.Time
func (m *model) debounceInputCmd(duration_ms int) tea.Cmd {
    m.state.game_ui.input_restricted = true

    return tea.Tick(time.Millisecond * time.Duration(duration_ms), func(t time.Time) tea.Msg {
		return EnableInputMsg(t)
	})
}

type TogglePlayerDamagedMsg struct{}
func (m *model) setPlayerDamagedStateCmd() tea.Cmd {
	m.state.game_ui.player_damaged = true
    return tea.Tick(time.Millisecond * time.Duration(400), func(t time.Time) tea.Msg {
		return TogglePlayerDamagedMsg{}
	})
}

func (m model) terminalBellCmd(force bool) tea.Cmd {
	if force || m.app_settings.Prefs.BellEnabled {
		return func() tea.Msg {
			// Send BEL character
			fmt.Fprint(os.Stdout, "\a")
			return nil
		}
	}
	return nil
}
