package tui

import (
	"fmt"
	"fzw/src/game"
	"fzw/src/utils"
	"os"
	"runtime"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func memStatsView() string {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	var sb strings.Builder

	// Print total memory allocated and still in use (in bytes)
	sb.WriteString(fmt.Sprintf("Total Alloc = %v MiB", memStats.TotalAlloc/1024/1024))
	sb.WriteString(" | ")
	sb.WriteString(fmt.Sprintf("Sys = %v MiB\n", memStats.Sys/1024/1024))
	sb.WriteString(fmt.Sprintf("Heap Alloc = %v MiB", memStats.HeapAlloc/1024/1024))
	sb.WriteString(" | ")
	sb.WriteString(fmt.Sprintf("Heap Sys = %v MiB", memStats.HeapSys/1024/1024))

	return sb.String()
}

func (m model) GameSwitch() (model, tea.Cmd) {
	m = m.SwitchPage(game_page)
	m.game_active = true

	// TODO: initialize word lists in background on program load
    word_list, err := utils.ReadLines("./wordlist.txt")
    if err != nil {
        fmt.Fprintf(os.Stderr, "error: %v\n", err)
        os.Exit(1)
    }

    m.word_lists = game.WordLists{
        FULL_MAP: utils.ArrToMap(word_list),
        Available: word_list,
        Used: make(map[string]bool),
    }
	
	m.turn = game.NewTurn(m.word_lists.Available, m.settings)

	return m, textinput.Blink
}

func (m model) GameUpdate(msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.turn.Answer = strings.ToLower(m.text_input.Value())
			m.turn.ValidateAnswer(&m.word_lists, m.settings)

			// may need to move out of switch/case?
			if m.turn.IsValid {
				m.player.HandleCorrectAnswer(m.turn.Answer)
				m.turn = game.NewTurn(m.word_lists.Available, m.settings)
				m.text_input.Reset()
			} else {
				m.turn.Strikes++
			}
		}
	}

	var cmd tea.Cmd
	m.text_input, cmd = m.text_input.Update(msg)

	return m, cmd
}

func (m model) GameView() string {
	debug_info := ""
	if m.debug {
		debug_info = memStatsView()
	}

	return lipgloss.JoinVertical(
		lipgloss.Center,
		debug_info,
		m.player.HealthDisplay,
		"",
		"Prompt: " + strings.ToUpper(m.turn.Prompt),
		"",
		strings.Join(m.player.LettersRemaining, " "),
		"",
		m.InputField.Render(m.text_input.View()),
	) 
}
