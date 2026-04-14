package tui

import (
	"fmt"
	"fzwds/src/assert"
	"fzwds/src/game"
	"fzwds/src/tui/animations"
	"fzwds/src/tui/styles"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

func (m model) GameHudView() string {
	return lipgloss.JoinVertical(
		lipgloss.Center,
		m.renderTopBar(),
		m.renderRemainingLetters())
}

func (m model) renderHealthDisplay(health_current int) string {
	assert.Assert(health_current >= 0, "Health cannot be less than 0", "health", health_current)

	// TODO: perform this check once on game startup rather than per redraw?
	health_icons := strings.Split(m.app_settings.Game.HealthDisplay, ";")
	if len(health_icons) != 2 {
		health_icons = strings.Split(game.GetDefaultSettings().Game.HealthDisplay, ";")
	}
	health_icon_full := health_icons[0]
	health_icon_empty := health_icons[1]

	var full_style, bracket_style lipgloss.Style
	if m.state.game_ui.player_damaged {
		full_style = styles.TextRed
		bracket_style = styles.TextRed
	} else {
		full_style = styles.TextHighlight
		bracket_style = styles.TextBody
	}

	var sb strings.Builder
	if strings.HasPrefix(health_icon_full, "#") {
		sb.WriteString(bracket_style.Render("["))
	}

	health_max := m.state.game.Settings.HealthMax
	sb.WriteString(full_style.Render(strings.Repeat(health_icon_full, health_current)))
	sb.WriteString(styles.TextBody.Render(strings.Repeat(health_icon_empty, health_max - health_current)))

	if strings.HasPrefix(health_icon_full, "#") {
		sb.WriteString(bracket_style.Render("]"))
	}

	return strings.TrimSpace(sb.String())
}

func (m model) renderTopBar() string {
	red := styles.TextRed.Render

    var timer_display string
	if m.state.game_ui.player_damaged || !m.state.game.GameActive {
		timer_display = "⌛️ 0.0s"
	} else if m.state.game.TimeRemaining().Seconds() <= 9.9 {
		timer_display = fmt.Sprintf("⏳ %.1fs", m.state.game.TimeRemaining().Seconds())
	} else {
		timer_display = fmt.Sprintf("⏳  %.0fs", m.state.game.TimeRemaining().Seconds())
    }

    if m.state.game.GameActive && (m.state.game.TimeRemaining().Seconds() < 5 || m.state.game_ui.player_damaged) {
		// TODO: pulsing yellow/orange/red anim when below 5s; red 0.0 on damaged
        timer_display = red(timer_display)
    }

	var text_style, border_style lipgloss.Style
	if m.state.game_ui.player_damaged {
		text_style = styles.TextRed
		border_style = styles.TextRed
	} else {
		text_style = styles.TextBody
		border_style = styles.TextBorder
	}

	row_items := []string {
		m.renderHealthDisplay(m.state.game.Player.HealthCurrent),
		text_style.Render(timer_display),
	}

	header := table.New().
		Border(lipgloss.NormalBorder()).
		BorderLeft(false).
		BorderRight(false).
		BorderStyle(border_style).
		BorderColumn(false).
		Row(row_items...).
		Width(m.width_container).
		StyleFunc(func(row, col int) lipgloss.Style {
			if col == 0 {
				return lipgloss.NewStyle().Align(lipgloss.Left).PaddingLeft(8)
			}
			return lipgloss.NewStyle().Align(lipgloss.Right).PaddingRight(8)
		}).
		Render()

	return header
}

func (m model) renderRemainingLetters() string {
	if !m.state.game.GameActive {
		return ""
	}

	letters, changed := m.anim_mgr.ApplyAnimations(
		string(animations.ExtraLife),
		strings.Join(strings.Split(m.state.game.Alphabet, ""), " "))
	if changed {
		return letters
	}

	var out strings.Builder
	for i, c := range m.state.game.Alphabet {
		if m.state.game.Player.LettersRemaining[c] {
			out.WriteString(styles.TextDim.Render(string(c)))
		} else if m.state.game_ui.player_damaged {
			out.WriteString(styles.TextRed.Bold(true).Render(string(c)))
		} else {
			out.WriteString(styles.TextYellow.Bold(true).Render(string(c)))
		}

		if i < len(m.state.game.Alphabet) - 1 {
			out.WriteRune(' ')
		}
	}

	return out.String()
}
