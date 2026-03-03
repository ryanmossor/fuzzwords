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
	if !m.state.game_ui.game_active {
		return ""
	}

	dim := m.theme.TextExtraDim().Render
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

	game_mode := fmt.Sprintf("Mode: %s", m.state.game.Settings.PromptMode.String())

	var fields []string
	if m.state.game_ui.player_damaged {
		fields = []string{
			red(health),
			red(game_mode),
			"Time: " + timer_display,
		}
	} else {
		fields = []string{
			health,
			game_mode,
			"Time: " + timer_display,
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
		Row(fields...).
		Width(m.width_container).
		StyleFunc(func(row, col int) lipgloss.Style {
			return m.theme.Base().AlignHorizontal(lipgloss.Center)
		}).
		Render()

	letters_remaining := []string{}
	var colored_alphabet string
	for _, c := range m.state.game.Alphabet {
		letter := string(c)
		// TODO: make this better
		if m.state.game_ui.extra_life_anim.active {
			colored_alphabet = m.ApplyExtraLifeFlashAnim()
			break
		} else if m.state.game.Player.LettersRemaining[letter] {
			letters_remaining = append(letters_remaining, dim(letter))
		} else if m.state.game_ui.player_damaged {
			letters_remaining = append(letters_remaining, red(letter))
		} else {
			letters_remaining = append(letters_remaining, yellow(letter))
		}
	}

	if colored_alphabet != "" {
		return lipgloss.JoinVertical(
			lipgloss.Center,
			m.DebugView(),
			header,
			colored_alphabet)
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

	for i < m.state.game.Player.HealthCurrent {
		if m.state.game_ui.player_damaged {
			health_display.WriteString(red("█"))
		} else {
			health_display.WriteString(green("█"))
		}

		if i < m.state.game.Settings.HealthMax - 1 {
			health_display.WriteString(" ")
		}
		i++
	}

	for i < m.state.game.Settings.HealthMax {
		health_display.WriteString(base("▒"))
		if i < m.state.game.Settings.HealthMax - 1 {
			health_display.WriteString(" ")
		}
		i++
	}

	return health_display.String()
}

func (m model) ApplyExtraLifeFlashAnim() string {
	var colors []lipgloss.Style = m.theme.GetRainbowColors()
	colored_letters := []string{}

	for i, c := range m.state.game.Alphabet {
		color_idx := (i - m.state.game_ui.extra_life_anim.offset) % len(colors)
		if color_idx < 0 {
			color_idx += len(colors)
		}
		color := colors[color_idx]
		colored_letters = append(colored_letters, color.Bold(true).Render(string(c)))
	}

	return strings.Join(colored_letters, " ")
}
