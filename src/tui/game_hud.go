package tui

import (
	"fmt"
	"strconv"
	"strings"

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

	base := m.theme.Base().Render
	dim := m.theme.TextExtraDim().Render
	yellow := m.theme.TextYellow().Bold(true).Render
	red := m.theme.TextRed().Render

	health := m.RenderHealthDisplay()

	var strikes string
	if m.game_state.CurrentTurn.Strikes > 0 {
		strikes = base("Strikes: ") + red(strconv.Itoa(m.game_state.CurrentTurn.Strikes)) + base(" / " + strconv.Itoa(m.game_state.Settings.PromptStrikesMax))
	} else {
		strikes = fmt.Sprintf("Strikes: %d / %d", m.game_state.CurrentTurn.Strikes, m.game_state.Settings.PromptStrikesMax)
	}

	game_mode := fmt.Sprintf("Mode: %s", m.game_state.Settings.PromptMode.String())

	fields := []string{
		health,
		game_mode,
		strikes,
	}

	header := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(m.renderer.NewStyle().Foreground(m.theme.Border())).
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

	var health_display strings.Builder
	i := 0

	for i < m.game_state.Player.HealthCurrent {
		health_display.WriteString(green("█"))
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
