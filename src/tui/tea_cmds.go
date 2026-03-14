package tui

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

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

type DamageShakeAnimationMsg struct{}
func (m *model) damageShakeAnimationCmd(count int) tea.Cmd {
	if !m.enable_animations {
		return nil
	}

	m.state.game_ui.damage_anim_padding = count * 2
	return tea.Tick(time.Second / time.Duration(m.FPS), func(t time.Time) tea.Msg {
		if m.state.game_ui.damage_anim_padding > 0 {
			return DamageShakeAnimationMsg{}
		}
		return nil
	})
}

type TurnTimerTickMsg struct{}
func (m *model) setTurnTickerCmd() tea.Cmd {
	if m.state.game_ui.timer > time.Second * 10 {
		m.state.game_ui.timer -= time.Second
		return tea.Tick(time.Second, func(t time.Time) tea.Msg {
			return TurnTimerTickMsg{}
		})
	}

	m.state.game_ui.timer -= time.Millisecond * 100
	return tea.Tick(time.Millisecond * 100, func(t time.Time) tea.Msg {
		return TurnTimerTickMsg{}
	})
}

func (m model) terminalBellCmd(force bool) tea.Cmd {
	if force || m.game_settings.BellEnabled {
		return func() tea.Msg {
			// Send BEL character
			fmt.Fprint(os.Stdout, "\a")
			return nil
		}
	}
	return nil
}
