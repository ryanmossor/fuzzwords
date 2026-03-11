package tui

import (
	"fmt"
	"fzwds/src/game"
	"fzwds/src/tui/animations"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

func (m model) GameHudUpdate(msg tea.Msg) (model, tea.Cmd) {
	return m, nil
}

func (m model) GameHudView() string {
	if !m.state.game_ui.game_active {
		return ""
	}

	dim := m.theme.TextDim().Render
	yellow := m.theme.TextYellow().Bold(true).Render
	red := m.theme.TextRed().Render

	health := m.RenderHealthDisplay()

    var timer_display string
	if m.state.game_ui.timer >= 10 * time.Second {
        timer_display = fmt.Sprintf("%.0fs", m.state.game_ui.timer.Seconds())
	} else {
        timer_display = fmt.Sprintf("%.1fs", m.state.game_ui.timer.Seconds())
    }

    if m.state.game_ui.timer < 5 * time.Second {
        timer_display = red(timer_display)
    }

	var fields []string
	if m.state.game_ui.player_damaged {
		fields = []string{
			red("Health: " + health),
			"⏳ " + timer_display,
		}
	} else {
		fields = []string{
			"Health: " + health,
			"⏳ " + timer_display,
		}
	}

	var border_style lipgloss.Style
	if m.state.game_ui.player_damaged {
		border_style = m.renderer.NewStyle().Foreground(m.theme.red)
	} else {
		border_style = m.renderer.NewStyle().Foreground(m.theme.Border())
	}

	header := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(border_style).
		BorderColumn(false).
		Row(fields...).
		Width(m.width_container).
		StyleFunc(func(row, col int) lipgloss.Style {
			if col == 0 {
				return m.theme.Base().Align(lipgloss.Left).PaddingLeft(5)
			}
			return m.theme.Base().Align(lipgloss.Right).PaddingRight(5)
		}).
		Render()

	if effects := m.animation_manager.EffectsFor(string(animations.ExtraLife)); len(effects) > 0 {
		animated_alphabet := animations.ApplyTextEffects(
			strings.Join(strings.Split(m.state.game.Alphabet, ""), " "),
			effects...)

		return lipgloss.JoinVertical(
			lipgloss.Center,
			m.DebugView(),
			header,
			animated_alphabet)
	}

	letters_remaining := []string{}
	for _, c := range m.state.game.Alphabet {
		letter := string(c)
		if m.state.game.Player.LettersRemaining[letter] {
			letters_remaining = append(letters_remaining, dim(letter))
		} else if m.state.game_ui.player_damaged {
			letters_remaining = append(letters_remaining, red(letter))
		} else {
			letters_remaining = append(letters_remaining, yellow(letter))
		}
	}

	return lipgloss.JoinVertical(
		lipgloss.Center,
		m.DebugView(),
		header,
		strings.Join(letters_remaining, " "))
}

func (m model) RenderHealthDisplay() string {
	health_icons := strings.Split(m.game_settings.HealthDisplay, ";")
	if len(health_icons) != 2 {
		health_icons = strings.Split(game.GetDefaultSettings().HealthDisplay, ";")
	}
	health_icon_full := health_icons[0]
	health_icon_empty := health_icons[1]

	var full_style, bracket_style lipgloss.Style
	if m.state.game_ui.player_damaged {
		full_style = m.theme.TextRed()
		bracket_style = m.theme.TextRed()
	} else {
		full_style = m.theme.TextGreen()
		bracket_style = m.theme.Base()
	}

	var sb strings.Builder
	if strings.HasPrefix(health_icon_full, "#") {
		sb.WriteString(bracket_style.Render("["))
	}

	i := 0
	for i < m.state.game.Player.HealthCurrent {
		if m.state.game_ui.player_damaged {
			sb.WriteString(full_style.Render(health_icon_full))
		} else {
			sb.WriteString(full_style.Render(health_icon_full))
		}
		i++
	}

	for i < m.state.game.Settings.HealthMax {
		sb.WriteString(m.theme.Base().Render(health_icon_empty))
		i++
	}

	if strings.HasPrefix(health_icon_full, "#") {
		sb.WriteString(bracket_style.Render("]"))
	}
	return strings.TrimSpace(sb.String())
}
