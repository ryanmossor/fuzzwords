package tui

import (
	"fmt"
	"fzw/src/enums"
	"fzw/src/game"
	"fzw/src/utils"
	"os"
	"runtime"
	"strings"
	"time"

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
    word_list, err := utils.ReadLines("./wordlist.txt", m.settings.PromptLenMin)
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
	m.prompt_display = m.theme.TextAccent().Render(m.turn.Prompt)

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
				m.prompt_display = m.turn.Prompt
				// time.Sleep(750 * time.Millisecond)
				// m.text_input.Reset()
			} else {
				m.turn.Strikes++
			}

			if m.settings.WinCondition == enums.MaxLives && m.player.HealthCurrent == m.settings.HealthMax {
				// TODO: replace with switch to game over/stats view
				fmt.Println("Max lives achieved -- you win!")
				os.Exit(0)
			}

			if m.turn.Strikes == m.settings.PromptStrikesMax {
				m.player.HandleFailedTurn()

				if m.player.HealthCurrent == 0 {
					// fmt.Println()
					// fmt.Println("===== GAME OVER =====")
					// fmt.Println()
					m.player.Stats.GenerateFinalStats()
					
					// TODO: replace with switch to game over/stats view
					os.Exit(0)
				} else {
					m.turn = game.NewTurn(m.word_lists.Available, m.settings)
					m.prompt_display = m.turn.Prompt
					// m.text_input.Reset()
					// TODO: debounce while sleeping -- bug causing increase of strikes if spamming enter
					time.Sleep(2 * time.Second)
				}
			}

			m.text_input.Reset()
		}
	}

	// TODO: tolower answer & prompt
	switch m.settings.PromptMode {
	case enums.Fuzzy:
		sub_idx := 0
		colorized := strings.Split(m.turn.Prompt, "")
		for i := range len(m.turn.Prompt) {
			substr := m.text_input.Value()[sub_idx:]
			current_prompt_char := string(m.turn.Prompt[i])

			if !strings.Contains(substr, current_prompt_char) {
				break
			}

			colorized[i] = m.theme.TextHighlight().Render(current_prompt_char)
			sub_idx += strings.Index(substr, current_prompt_char) + 1
		}
		m.prompt_display = strings.Join(colorized, "")
	// case enums.Classic:
	// 	if !strings.Contains(t.Answer, string(t.Prompt)) {
	// 		t.IsValid = false
	// 		t.Msg = "Word does not satisfy the prompt. Try again." 
	// 		return
	// 	}
	}

	var cmd tea.Cmd
	m.text_input, cmd = m.text_input.Update(msg)

	return m, cmd
}

func (m model) GameView() string {
	// debug_info := ""
	// if m.debug {
	// 	// debug_info = memStatsView()
	// 	debug_info = fmt.Sprintf("answer: %s | strikes: %d | isValid: %t | msg: %s", m.turn.Answer, m.turn.Strikes, m.turn.IsValid, m.turn.Msg)
	// }

	// TODO: show possible answer after striking out

	var turn_msg string
	if !m.turn.IsValid && m.turn.Strikes < m.settings.PromptStrikesMax {
		turn_msg = m.theme.TextError().Render(m.turn.Msg)
	} else if !m.turn.IsValid && m.turn.Strikes == m.settings.PromptStrikesMax {
		turn_msg = fmt.Sprintf("Prompt failed. Possible answer: %s", m.turn.SourceWord)
	}
	// 	turn_msg = m.theme.TextHighlight().Render(m.turn.Msg)

	return lipgloss.JoinVertical(
		lipgloss.Center,
		// debug_info,
		// "Prompt: " + strings.ToUpper(m.turn.Prompt),
		// "Prompt: " + strings.ToUpper(m.prompt_display),
		m.theme.TextAccent().Render("Prompt: " + m.prompt_display),
		"",
		turn_msg,
		m.InputField.Render(m.text_input.View()),
	) 
}
