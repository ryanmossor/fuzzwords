package tui

import (
	"fmt"
	"strconv"
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

	base := m.theme.Base().Render
	accent := m.theme.TextAccent().Render
	yellow := m.theme.TextYellow().Bold(true).Render
	error := m.theme.TextError().Render

	var fields []string

	health := m.player.HealthDisplay

	var strikes string
	if m.turn.Strikes > 0 {
		strikes = accent("Strikes: ") + error(strconv.Itoa(m.turn.Strikes)) + accent(" / " + strconv.Itoa(m.settings.PromptStrikesMax))
	} else {
		strikes = accent(fmt.Sprintf("Strikes: %d / %d", m.turn.Strikes, m.settings.PromptStrikesMax))
	}

	// elapsed_sec := int(time.Since(m.game_start_time).Seconds())
	// elapsed_formatted := "⏱  " + utils.FormatTime(elapsed_sec)

	game_mode := fmt.Sprintf("Mode: %s", m.settings.PromptMode.String())

	fields = []string{
		health,
		strikes,
		// elapsed_formatted,
		game_mode,
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
	for _, c := range m.settings.Alphabet {
		letter := string(c)
		if m.player.LettersRemaining[letter] {
			letters_remaining = append(letters_remaining, base(letter))
		} else {
			letters_remaining = append(letters_remaining, yellow(letter))
		}
	}

	return lipgloss.JoinVertical(
		lipgloss.Center,
		header,
		strings.Join(letters_remaining, " "))
}
