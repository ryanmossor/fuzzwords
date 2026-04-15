package tui

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type TurnTimerExpiredMsg struct {
	timerId		uint
	duration 	time.Duration
}
func (m model) turnTimerExpiredCmd(timer_id uint, duration time.Duration) tea.Cmd {
    return tea.Tick(duration, func(t time.Time) tea.Msg {
		return TurnTimerExpiredMsg{ timerId: timer_id, duration: duration }
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
