package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

// unused for now
func (m model) GameHudUpdate(msg tea.Msg) (model, tea.Cmd) {
	return m, nil
}

func (m model) GameHudView() string {
	if !m.game_active {
		return ""
	}

	// bold := m.theme.TextAccent().Bold(true).Render
	accent := m.theme.TextAccent().Render
	base := m.theme.Base().Render

	var fields []string

	health := m.player.HealthDisplay
	strikes := fmt.Sprintf("Strikes: %d / %d", m.turn.Strikes, m.settings.PromptStrikesMax)
	// TODO: total time elapsed in rightmost column

	fields = []string{
		health,
		strikes,
		// timeElapsed,
	}

	header := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(m.renderer.NewStyle().Foreground(m.theme.Border())).
		Row(fields...).
		Width(m.width_container).
		StyleFunc(func(row, col int) lipgloss.Style {
			return m.theme.Base().
				Padding(0, 1).
				AlignHorizontal(lipgloss.Center)
		}).
		Render()

	letters_remaining := []string{}
	for _, c := range m.settings.Alphabet {
		letter := string(c)
		if m.player.LettersRemaining[letter] {
			letters_remaining = append(letters_remaining, base(letter))
		} else {
			letters_remaining = append(letters_remaining, accent(letter))
		}
	}

	return lipgloss.JoinVertical(
		lipgloss.Center,
		header,
		strings.Join(letters_remaining, " "))
}
