package tui

import (
	"fzwds/src/constants"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var letters = map[byte][]string {
	'f': {
		"  ▄▄",
		" ██ ",
		"▀██▀",
		" ██ ",
		" ██ ",
	},
	'u': {
		"     ",
		"     ",
		"██ ██",
		"██ ██",
		"▀██▀█",
	},
	'z': {
		"     ",
		"     ",
		"▀▀▀██",
		"  ▄█▀",
		"▄██▄▄",
	},
	'w': {
		"       ",
		"       ",
		"██   ██",
		"██ █ ██",
		" ██▀██ ",
	},
	'o': {
		"     ",
		"     ",
		"▄███▄",
		"██ ██",
		"▀███▀",
	},
	'r': {
		"     ",
		"     ",
		"████▄",
		"██ ▀▀",
		"██   ",
	},
	'd': {
		"   ▄▄",
		"   ██",
		"▄████",
		"██ ██",
		"▀████",
	},
	's': {
		"     ",
		"     ",
		"▄█▀▀▀",
		"▀███▄",
		"▄▄▄█▀",
	},
}

func (m model) MainMenuInit() tea.Cmd {
	return tea.Tick(3500 * time.Millisecond, func(t time.Time) tea.Msg {
		return LogoInitMsg{}
	})
}

func (m model) MainMenuSwitch() (model, tea.Cmd) {
	m = m.SwitchPage(splash_page)
	m.footer_cmds = []footerCmd{
		{key: "q", value: "quit"},
	}

	return m, nil
}

func (m model) MainMenuUpdate(msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return m.SettingsSwitch()
		case "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) MainMenuView() string {
	base := m.theme.base.Render
	highlight := m.theme.TextBlue().Render
	yellow := m.theme.TextYellow().Render

	HEADER_LEN := 2
	logo := make([]string, HEADER_LEN + len(letters['f']))
	logo[0] = "\n"
	logo[1] = "\n"

	switch m.size {
	case large:
		if !m.enable_animations {
			logo[2] = yellow("  ▄▄                                          ▄▄      ")
			logo[3] = yellow(" ██                                           ██      ")
			logo[4] = yellow("▀██▀ ██ ██ ▀▀▀██ ▀▀▀██ ██   ██ ▄███▄ ████▄ ▄████ ▄█▀▀▀")
			logo[5] = yellow(" ██  ██ ██   ▄█▀   ▄█▀ ██ █ ██ ██ ██ ██ ▀▀ ██ ██ ▀███▄")
			logo[6] = yellow(" ██  ▀██▀█ ▄██▄▄ ▄██▄▄  ██▀██  ▀███▀ ██    ▀████ ▄▄▄█▀")
		} else if !m.state.title.init {
			logo[2] = yellow("  ▄▄                  ▄▄      ")
			logo[3] = yellow(" ██                   ██      ")
			logo[4] = yellow("▀██▀ ▀▀▀██ ██   ██ ▄████ ▄█▀▀▀")
			logo[5] = yellow(" ██    ▄█▀ ██ █ ██ ██ ██ ▀███▄")
			logo[6] = yellow(" ██  ▄██▄▄  ██▀██  ▀████ ▄▄▄█▀")
		} else {
			prompt_idx := 0
			for i := range m.state.title.logo_idx {
				current_title_char := constants.GAME_TITLE[i]

				style := base
				for j := prompt_idx; j < len(constants.TITLE_PROMPT); j++ {
					c := constants.TITLE_PROMPT[j]
					is_prompt_letter := c == current_title_char
					if is_prompt_letter {
						style = highlight
						prompt_idx++
						break
					}
				}

				letter_arr := letters[current_title_char]
				for k, line := range letter_arr {
					idx := HEADER_LEN + k
					logo[idx] = logo[idx] + style(line) + " "
				}
			}
		}
	default:
		logo = append(logo, base(" ___ __       __   __  "))
		logo = append(logo, base("|__   / |  | |  \\ /__` "))
		logo = append(logo, base("|    /_ |/\\| |__/ .__/ "))
	}

	logo = append(logo, "\n\n\n")
	logo = append(logo, m.PressPlayView())

	return lipgloss.JoinVertical(
		lipgloss.Center,
		logo...
	)
}
