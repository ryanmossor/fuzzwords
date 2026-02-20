package tui

import (
	"fmt"
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
	if !m.game_active {
		return ""
	}

	// base := m.theme.Base().Render
	dim := m.theme.TextExtraDim().Render
	yellow := m.theme.TextYellow().Bold(true).Render
	red := m.theme.TextRed().Render

	health := m.RenderHealthDisplay()

    var timer_display string
	if m.game_timer.remaining_time >= 10 * time.Second {
        timer_display = fmt.Sprintf("%.0fs", m.game_timer.remaining_time.Seconds())
	} else {
        timer_display = fmt.Sprintf("%.1fs", m.game_timer.remaining_time.Seconds())
    }

    if m.game_timer.remaining_time < 5 * time.Second {
        timer_display = red(timer_display)
    }

	game_mode := fmt.Sprintf("Mode: %s", m.game_state.Settings.PromptMode.String())

	fields := []string{
		health,
		game_mode,
		"Time: " + timer_display,
	}

	var border_style lipgloss.Style
	if m.state.game.damaged {
		border_style = m.renderer.NewStyle().Foreground(m.theme.red)
	} else {
		border_style = m.renderer.NewStyle().Foreground(m.theme.Border())
	}

	header := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(border_style).
		Row(fields...).
		Width(m.width_container).
		StyleFunc(func(row, col int) lipgloss.Style {
			return m.theme.Base().AlignHorizontal(lipgloss.Center)
		}).
		Render()

	letters_remaining := []string{}
	for _, c := range m.game_state.Alphabet {
		letter := string(c)
		if m.game_state.Player.LettersRemaining[letter] {
			letters_remaining = append(letters_remaining, dim(letter))
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
	base := m.theme.Base().Render
	green := m.theme.TextGreen().Render
	red := m.theme.TextRed().Render

	var health_display strings.Builder
	i := 0

	for i < m.game_state.Player.HealthCurrent {
		if m.state.game.damaged {
			health_display.WriteString(red("█"))
		} else {
			health_display.WriteString(green("█"))
		}

		if i < m.game_state.Settings.HealthMax - 1 {
			health_display.WriteString(" ")
		}
		i++
	}

	for i < m.game_state.Settings.HealthMax {
		health_display.WriteString(base("▒"))
		if i < m.game_state.Settings.HealthMax - 1 {
			health_display.WriteString(" ")
		}
		i++
	}

	return health_display.String()
}
