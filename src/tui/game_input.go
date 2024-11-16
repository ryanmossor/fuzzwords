package tui

import (
	"fmt"
	"fzw/src/enums"
	"fzw/src/game"
	"fzw/src/utils"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

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
			// TODO: trim answer & take only first word before any spaces/symbols
			m.turn.Answer = strings.ToLower(m.text_input.Value())
			m.turn.ValidateAnswer(&m.word_lists, m.settings)

			// may need to move out of switch/case?
			if m.turn.IsValid {
				m.player.HandleCorrectAnswer(m.turn.Answer)
				m.turn = game.NewTurn(m.word_lists.Available, m.settings)
				// time.Sleep(750 * time.Millisecond)
				// m.text_input.Reset()
			} else {
				m.turn.Strikes++
			}

			if m.settings.WinCondition == enums.MaxLives && m.player.HealthCurrent == m.settings.HealthMax {
				// TODO: replace with switch to game over/stats view
				fmt.Println("Max lives achieved -- you win!")
				return m, tea.Quit
			}

			if m.turn.Strikes == m.settings.PromptStrikesMax {
				m.player.HandleFailedTurn()

				if m.player.HealthCurrent == 0 {
					// fmt.Println()
					// fmt.Println("===== GAME OVER =====")
					// fmt.Println()
					m.player.Stats.GenerateFinalStats()
					
					// TODO: replace with switch to game over/stats view
					return m, tea.Quit
				} else {
					m.turn = game.NewTurn(m.word_lists.Available, m.settings)
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

func (m model) GameInputView() string {
	if !m.game_active {
		return ""
	}

	accent := m.theme.TextAccent().Render
	blue := m.theme.TextBlue().Render

	prompt_upper := strings.ToUpper(m.turn.Prompt)
	answer_upper := strings.ToUpper(m.text_input.Value())
	var sb strings.Builder
	 
	switch m.settings.PromptMode {
	case enums.Fuzzy:
		prompt_idx := 0
		for _, c := range answer_upper {
			curr_char := string(c)

			if prompt_idx < len(prompt_upper) && curr_char == string(prompt_upper[prompt_idx]) {
				sb.WriteString(blue(curr_char))
				prompt_idx++
			} else {
				sb.WriteString(accent(curr_char))
			}
		}
	case enums.Classic:
		if !strings.Contains(answer_upper, prompt_upper) {
			sb.WriteString(accent(answer_upper))
			break
		}
		
		sub_idx := strings.Index(answer_upper, prompt_upper)
		sb.WriteString(accent(answer_upper[0:sub_idx]))
		sb.WriteString(blue(answer_upper[sub_idx:sub_idx + len(prompt_upper)]))
		sb.WriteString(accent(answer_upper[sub_idx + len(prompt_upper):]))
	}

	// TODO: show possible answer after striking out
	var turn_msg string
	if !m.turn.IsValid && m.turn.Strikes < m.settings.PromptStrikesMax {
		turn_msg = m.theme.TextRed().Render(m.turn.Msg)
	} else if !m.turn.IsValid && m.turn.Strikes == m.settings.PromptStrikesMax {
		turn_msg = fmt.Sprintf("Prompt failed. Possible answer: %s", m.turn.SourceWord)
	}
	// 	turn_msg = m.theme.TextHighlight().Render(m.turn.Msg)

	return lipgloss.JoinVertical(
		lipgloss.Center,
		sb.String(),
		"",
		turn_msg,
		lipgloss.NewStyle().
			BorderForeground(m.theme.border).
			BorderStyle(lipgloss.RoundedBorder()).
			Width(50).
			Render(m.text_input.View()),
	) 
}
