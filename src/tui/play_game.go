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
	m.game_start_time = time.Now()

	m.footerCmds = []footerCmd{
		{key: "esc", value: "clear input"},
	}

	return m, textinput.Blink
}

func (m model) GameUpdate(msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.text_input.Reset()
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

	// var turn_msg string
	// if !m.turn.IsValid && m.turn.Strikes < m.settings.PromptStrikesMax {
	// 	turn_msg = m.theme.TextError().Render(m.turn.Msg)
	// } else if !m.turn.IsValid && m.turn.Strikes == m.settings.PromptStrikesMax {
	// 	turn_msg = fmt.Sprintf("Prompt failed. Possible answer: %s", m.turn.SourceWord)
	// }
	// 	turn_msg = m.theme.TextHighlight().Render(m.turn.Msg)

	// Show WIP answer with prompt letters highlighted
	// var colorized_answer []string
	 
	// switch m.settings.PromptMode {
	// case enums.Fuzzy:
	// 	prompt_caps := strings.ToUpper(m.turn.Prompt)
	// 	prompt_idx := 0
	// 	for _, c := range strings.ToUpper(m.text_input.Value()) {
	// 		curr_char := string(c)

	// 		if prompt_idx < len(m.turn.Prompt) && curr_char == string(prompt_caps[prompt_idx]) {
	// 			colorized_answer = append(colorized_answer, m.theme.TextHighlight().Render(curr_char))
	// 			prompt_idx++
	// 		} else {
	// 			colorized_answer = append(colorized_answer, m.theme.TextAccent().Render(curr_char))
	// 		}
	// 	}
	// case enums.Classic:
	// 	// TODO
	// }

	prompt_ascii := []string{"", "", "", "", "", ""}

	for _, c := range strings.Split(m.turn.Prompt, "") {
		for i, l := range utils.MedFont[strings.ToUpper(c)] {
			prompt_ascii[i] += m.theme.TextAccent().Render(l)
			prompt_ascii[i] += " "
		}
	}
	// fmt.Println(output)

	// for _, l := range output {
	// 	fmt.Println(l)
	// }

	return lipgloss.JoinVertical(
		lipgloss.Center,
		// debug_info,
		// "",
		// m.theme.TextAccent().Render("Prompt: " + strings.ToUpper(m.turn.Prompt)),
		// m.theme.TextAccent().Render(output...),
		prompt_ascii...
		// "",
		// strings.Join(colorized_answer, ""),
		// "",
		// turn_msg,
		// m.InputField.Render(m.text_input.View()),
	) 
}
