package tui

import (
	"fmt"
	"fzwds/pkg/assert"
	"fzwds/pkg/game"
	"fzwds/pkg/tui/animations"
	"fzwds/pkg/tui/styles"
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
	health_icons := strings.Split(m.settings.Prefs.HealthDisplay, ";")
	if len(health_icons) != 2 {
		health_icons = strings.Split(game.GetDefaultSettings().Prefs.HealthDisplay, ";")
	}
	health_icon_full := health_icons[0]
	health_icon_empty := health_icons[1]

	var full_style, bracket_style lipgloss.Style
	if m.state.game.playerDamaged {
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

	health_max := m.game.Settings().HealthMax
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
	if !m.game.GameActive() {
		timer_display = "⌛️  ─  "
	} else if m.state.game.playerDamaged {
		timer_display = "⌛️ 0.0s"
	} else if m.game.TimeRemaining().Seconds() <= 9.9 {
		timer_display = fmt.Sprintf("⏳ %.1fs", m.game.TimeRemaining().Seconds())
	} else {
		timer_display = fmt.Sprintf("⏳  %.0fs", m.game.TimeRemaining().Seconds())
    }

    if m.game.GameActive() && (m.game.TimeRemaining().Seconds() < 5 || m.state.game.playerDamaged) {
		// TODO: pulsing yellow/orange/red anim when below 5s; red 0.0 on damaged
        timer_display = red(timer_display)
    }

	var text_style, border_style lipgloss.Style
	if m.state.game.playerDamaged {
		text_style = styles.TextRed
		border_style = styles.TextRed
	} else {
		text_style = styles.TextBody
		border_style = styles.TextBorder
	}

	row_items := []string {
		m.renderHealthDisplay(int(m.state.game.health)),
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
	if !m.game.GameActive() {
		return ""
	}

	letters, changed := m.anim_mgr.ApplyAnimations(
		string(animations.ExtraLife),
		strings.Join(strings.Split(m.game.Settings().Alphabet.Letters(), ""), " "))
	if changed {
		return letters
	}

	var out strings.Builder
	for i, c := range m.game.Settings().Alphabet.Letters() {
		if m.state.game.lettersUsed[c] {
			out.WriteString(styles.TextDim.Render(string(c)))
		} else if m.state.game.playerDamaged {
			out.WriteString(styles.TextRed.Bold(true).Render(string(c)))
		} else {
			out.WriteString(styles.TextYellow.Bold(true).Render(string(c)))
		}

		if i < len(m.game.Settings().Alphabet.Letters()) - 1 {
			out.WriteRune(' ')
		}
	}

	return out.String()
}
