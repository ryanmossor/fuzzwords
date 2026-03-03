package tui

import (
	"fzwds/src/constants"
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

type LogoTickMsg struct{}
type LogoCompleteMsg struct{}
type LogoRestartMsg struct{}
type LogoUnhideMsg struct{}
func (m *model) initMainMenuLogoAnimCmd() tea.Cmd {
	return tea.Tick(5 * time.Second, func(t time.Time) tea.Msg {
		return LogoInitMsg{}
	})
}

func (m *model) mainMenuLogoUpdateCmd() tea.Cmd {
	if m.state.title.logo_anim_idx == len(constants.GAME_TITLE) {
		return tea.Tick(1500 * time.Millisecond, func(t time.Time) tea.Msg {
			return LogoCompleteMsg{}
		})
	}

	return tea.Tick(250 * time.Millisecond, func(t time.Time) tea.Msg {
		return LogoTickMsg{}
	})
}

type TogglePlayerDamagedMsg struct{}
func (m *model) setPlayerDamagedStateCmd() tea.Cmd {
	m.state.game_ui.player_damaged = true
    return tea.Tick(time.Millisecond * time.Duration(250), func(t time.Time) tea.Msg {
		return TogglePlayerDamagedMsg{}
	})
}

type DamageShakeAnimationMsg struct{}
func (m *model) damageShakeAnimationCmd(count int) tea.Cmd {
	if !m.enable_animations {
		return nil
	}

	m.state.game_ui.damage_anim_padding = count * 2
	return tea.Tick(time.Second / time.Duration(m.anim_fps), func(t time.Time) tea.Msg {
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

type ExtraLifeAnimInitMsg struct{}
type ExtraLifeAnimTickMsg struct{}
type ExtraLifeAnimCompleteMsg struct{}
func (m *model) extraLifeAnimInitMsg() tea.Cmd {
	m.state.game_ui.extra_life_anim.active = true
	return tea.Cmd(func() tea.Msg {
		return ExtraLifeAnimTickMsg{}
	})
}

func (m *model) extraLifeAnimTickMsg() tea.Cmd {
	anim := m.state.game_ui.extra_life_anim
	if anim.cur_frame == anim.total_frames {
		return tea.Cmd(func() tea.Msg {
			return ExtraLifeAnimCompleteMsg{}
		})
	}

	return tea.Tick(time.Second / time.Duration(anim.fps), func(t time.Time) tea.Msg {
		return ExtraLifeAnimTickMsg{}
	})
}
