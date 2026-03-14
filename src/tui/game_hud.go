package tui

import (
	"fmt"
	"fzwds/src/game"
	"fzwds/src/tui/animations"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

func (m model) GameHudView() string {
	return lipgloss.JoinVertical(
		lipgloss.Center,
		m.DebugView(),
		m.renderTopBar(),
		m.renderRemainingLetters())
}

func (m model) renderHealthDisplay() string {
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

	health_cur := m.state.game.Player.HealthCurrent
	health_max := m.state.game.Settings.HealthMax
	sb.WriteString(full_style.Render(strings.Repeat(health_icon_full, health_cur)))
	sb.WriteString(full_style.Render(strings.Repeat(health_icon_empty, health_max - health_cur)))

	if strings.HasPrefix(health_icon_full, "#") {
		sb.WriteString(bracket_style.Render("]"))
	}

	return strings.TrimSpace(sb.String())
}

func (m model) renderTopBar() string {
	red := m.theme.TextRed().Render

    var timer_display string
	if m.state.game_ui.player_damaged || !m.state.game_ui.game_active {
		timer_display = "⌛️ 0.0s"
	} else if m.state.game_ui.timer >= 10 * time.Second {
        timer_display = fmt.Sprintf("⏳  %.0fs", m.state.game_ui.timer.Seconds())
	} else {
        timer_display = fmt.Sprintf("⏳ %.1fs", m.state.game_ui.timer.Seconds())
    }

    if m.state.game_ui.game_active && (m.state.game_ui.timer < (5 * time.Second) || m.state.game_ui.player_damaged) {
		// TODO: pulsing yellow/orange/red anim when below 5s; red 0.0 on damaged
        timer_display = red(timer_display)
    }

	var text_style, border_style lipgloss.Style
	if m.state.game_ui.player_damaged {
		text_style = m.theme.TextRed()
		border_style = m.theme.Base().Foreground(m.theme.red)
	} else {
		text_style = m.theme.TextBody()
		border_style = m.theme.Base().Foreground(m.theme.Border())
	}

	row_items := []string{
		m.renderHealthDisplay(),
		text_style.Render(timer_display),
	}

	header := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(border_style).
		BorderColumn(false).
		Row(row_items...).
		Width(m.width_container).
		StyleFunc(func(row, col int) lipgloss.Style {
			if col == 0 {
				return m.theme.Base().Align(lipgloss.Left).PaddingLeft(7)
			}
			return m.theme.Base().Align(lipgloss.Right).PaddingRight(7)
		}).
		Render()

	return header
}

func (m model) renderRemainingLetters() string {
	if !m.state.game_ui.game_active {
		return ""
	}

	letters, changed := m.animation_manager.ApplyAnimations(
		string(animations.ExtraLife),
		strings.Join(strings.Split(m.state.game.Alphabet, ""), " "))
	if changed {
		return letters
	}

	out := []string{}
	for _, c := range m.state.game.Alphabet {
		letter := string(c)

		if m.state.game.Player.LettersRemaining[letter] {
			out = append(out, m.theme.TextDim().Render(letter))
		} else if m.state.game_ui.player_damaged {
			out = append(out, m.theme.TextRed().Bold(true).Render(letter))
		} else {
			out = append(out, m.theme.TextYellow().Bold(true).Render(letter))
		}
	}

	return strings.Join(out, " ")
}
