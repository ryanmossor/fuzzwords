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

type EnableInputMsg time.Time

func (m model) GameSwitch() (model, tea.Cmd) {
	m = m.SwitchPage(game_page)
	m.game_active = true
	m.game_over = false
	m.player = game.InitializePlayer(&m.settings)

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

	m.footer_cmds = []footerCmd{
		{key: "esc", value: "clear input"},
	}

	m.state.game.restrict_input = false
	m.text_input.Reset()

	return m, textinput.Blink
}

func (m model) GameUpdate(msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		m.text_input.Prompt = " > "
		m.text_input.PromptStyle = m.default_prompt_style

		switch msg.String() {
		case "esc":
			m.text_input.Reset()
		case "enter":
			if m.state.game.restrict_input {
				return m, nil
			}

			m.state.game.restrict_input = true

			// TODO: trim answer & take only first word before any spaces/symbols
			m.turn.Answer = strings.ToLower(m.text_input.Value())
			m.turn.ValidateAnswer(&m.word_lists, m.settings)

			// may need to move out of switch/case?
			if m.turn.IsValid {
				m.player.HandleCorrectAnswer(m.turn.Answer)

				if len(m.word_lists.Available) == 0 {
					m.game_active = false
					win_msg := m.theme.TextGreen().Bold(true).Render("YOU WIN!")
					return m.GameOverSwitch(win_msg)
				}

				m.text_input.Prompt = " ✓ "
				m.text_input.PromptStyle = m.theme.TextGreen().Bold(true)
				m.turn = game.NewTurn(m.word_lists.Available, m.settings)

				// time.Sleep(750 * time.Millisecond)
				// m.text_input.Reset()
			} else {
				m.turn.Strikes++
				m.text_input.Prompt = " ✗ "
				m.text_input.PromptStyle = m.theme.TextRed().Bold(true)
			}

			if (m.settings.WinCondition == enums.MaxLives && m.player.HealthCurrent == m.settings.HealthMax) {
				// TODO: are both of these flags needed?
				m.game_active = false
				win_msg := m.theme.TextGreen().Bold(true).Render("YOU WIN!")
				return m.GameOverSwitch(win_msg)
			}

			if m.turn.Strikes == m.settings.PromptStrikesMax {
				m.player.HandleFailedTurn()

				if m.player.HealthCurrent == 0 {
					// TODO: are both of these flags needed?
					m.game_active = false
					game_over_msg := m.theme.TextRed().Bold(true).Render("=== GAME OVER ===")
					return m.GameOverSwitch(game_over_msg)
				} else {
					m.turn = game.NewTurn(m.word_lists.Available, m.settings)
					// m.text_input.Reset()
					// TODO: debounce while sleeping -- bug causing increase of strikes if spamming enter
					// time.Sleep(2 * time.Second)
				}
			}

			m.text_input.Reset()

			// Debounce addtional enter presses
			return m, tea.Tick(time.Millisecond * 300, func(t time.Time) tea.Msg {
				return EnableInputMsg(t)
			})
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

	border_color := m.theme.Border()
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

		if m.settings.HighlightInput && utils.IsFuzzyMatch(answer_upper, prompt_upper) {
			border_color = m.setInputBorderColor(answer_upper)
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

		if m.settings.HighlightInput {
			border_color = m.setInputBorderColor(answer_upper)
		}
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
			BorderForeground(border_color).
			BorderStyle(lipgloss.RoundedBorder()).
			Width(50).
			Render(m.text_input.View()),
	) 
}

func (m model) setInputBorderColor(answer string) lipgloss.TerminalColor {
	if m.word_lists.FULL_MAP[strings.ToLower(answer)] {
		return m.theme.green
	}
	return m.theme.red
}
