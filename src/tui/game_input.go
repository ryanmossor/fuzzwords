package tui

import (
	"fzwds/src/enums"
	"fzwds/src/game"
	"fzwds/src/utils"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type EnableInputMsg time.Time

func (m model) GameSwitch() (model, tea.Cmd) {
	m = m.SwitchPage(game_page)

	// TODO: move these to game state?
	m.game_active = true
	m.game_over = false

	m.game_state = game.InitializeGame()
	m.game_state.NewTurn()

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
		// clear validation msg while debouncing enter presses (yes this is kinda scuffed)
		if msg.String() != "enter" {
			m.game_state.CurrentTurn.ValidationMsg = ""
		}

		switch msg.String() {
		case "esc":
			m.text_input.Reset()
		case "enter":
			if m.state.game.restrict_input {
				return m, nil
			}

			m.state.game.restrict_input = true

			// TODO: trim answer & take only first word before any spaces/symbols
			m.game_state.CurrentTurn.Answer = strings.ToLower(m.text_input.Value())
			m.game_state.CurrentTurn.ValidateAnswer(&m.game_state.WordLists, m.game_state.Settings)

			if m.game_state.CurrentTurn.IsValid {
				m.game_state.HandleCorrectAnswer()

				if len(m.game_state.WordLists.Available) == 0 {
					m.game_active = false
					win_msg := m.theme.TextGreen().Bold(true).Render("YOU WIN!")
					return m.GameOverSwitch(win_msg)
				}

				m.game_state.NewTurn()
			}

			if (m.game_state.Settings.WinCondition == enums.MaxLives && m.game_state.Player.HealthCurrent == m.game_state.Settings.HealthMax) {
				// TODO: are both of these flags needed?
				m.game_active = false
				win_msg := m.theme.TextGreen().Bold(true).Render("YOU WIN!")
				return m.GameOverSwitch(win_msg)
			}

			if m.game_state.CurrentTurn.Strikes == m.game_state.Settings.PromptStrikesMax {
				m.game_state.HandleFailedTurn()

				if m.game_state.Player.HealthCurrent == 0 {
					// TODO: are both of these flags needed?
					m.game_active = false
					game_over_msg := m.theme.TextRed().Bold(true).Render("===== GAME OVER =====")
					return m.GameOverSwitch(game_over_msg)
				} else {
					m.game_state.NewTurn()
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

	var colorized_text string
	var border_color lipgloss.TerminalColor

	if m.game_state.CurrentTurn.ValidationMsg != "" {
		colorized_text, border_color = m.renderValidationMsg()
	} else {
		colorized_text, border_color = m.renderColorizedInput()
	}

	return lipgloss.JoinVertical(
		lipgloss.Center,
		colorized_text,
		"\n",
		lipgloss.NewStyle().
			BorderForeground(border_color).
			BorderStyle(lipgloss.RoundedBorder()).
			Width(50).
			Render(m.text_input.View()),
	) 
}

func (m model) setInputBorderColor(answer string) lipgloss.TerminalColor {
	if m.game_state.WordLists.FULL_MAP[strings.ToLower(answer)] {
		return m.theme.green
	}
	return m.theme.red
}

func (m model) renderColorizedInput() (string, lipgloss.TerminalColor) {
	border_color := m.theme.Border()
	accent := m.theme.TextAccent().Render
	blue := m.theme.TextBlue().Render

	prompt_upper := strings.ToUpper(m.game_state.CurrentTurn.Prompt)
	answer_upper := strings.ToUpper(m.text_input.Value())
	var sb strings.Builder
	 
	switch m.game_state.Settings.PromptMode {
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

		if m.game_state.Settings.HighlightInput && utils.IsFuzzyMatch(answer_upper, prompt_upper) {
			border_color = m.setInputBorderColor(answer_upper)
		}
	case enums.Classic:
		if !strings.Contains(answer_upper, prompt_upper) {
			sb.WriteString(accent(answer_upper))
			return sb.String(), border_color
		}
		
		sub_idx := strings.Index(answer_upper, prompt_upper)
		sb.WriteString(accent(answer_upper[0:sub_idx]))
		sb.WriteString(blue(answer_upper[sub_idx:sub_idx + len(prompt_upper)]))
		sb.WriteString(accent(answer_upper[sub_idx + len(prompt_upper):]))

		if m.game_state.Settings.HighlightInput {
			border_color = m.setInputBorderColor(answer_upper)
		}
	}

	return sb.String(), border_color
}

func (m *model) renderValidationMsg() (string, lipgloss.TerminalColor) {
	border_color := m.theme.Border()

	if strings.Contains(m.game_state.CurrentTurn.ValidationMsg, "Correct") {
		return m.theme.TextGreen().Bold(true).Render(m.game_state.CurrentTurn.ValidationMsg), border_color
	}

	if m.game_state.CurrentTurn.Strikes > 0 {
		m.text_input.PromptStyle = m.theme.TextRed()
		border_color = m.theme.red
	}

	return m.theme.TextRed().Render(m.game_state.CurrentTurn.ValidationMsg), border_color
}
