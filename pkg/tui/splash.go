package tui

import (
	"fzwds/pkg/constants"
	"fzwds/pkg/tui/animations"
	"fzwds/pkg/tui/styles"

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
	m.anim_mgr.InitAnimations(animations.TitleLogo)
	return nil
}

func (m model) MainMenuSwitch() (model, tea.Cmd) {
	// Don't switch if already here; prevents title anim reload
	if m.page == splash_page {
		return m, nil
	}

	m = m.SwitchPage(splash_page)
	m.footer_keymaps = []FooterKeymap{
		{key: "ctrl+p", value: "preferences"},
		{key: "q", value: "quit"},
	}
	m.anim_mgr.InitAnimations(animations.TitleLogo)

	return m, m.pressPlayFlashCmd()
}

func (m model) MainMenuUpdate(msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return m.SettingsSwitch(game_settings)
		case "ctrl+p":
			return m.SettingsSwitch(preferences)
		}

	case PressPlayTickMsg:
		m.state.pressPlay.visible = !m.state.pressPlay.visible
		return m, m.pressPlayFlashCmd()
	}

	return m, nil
}

func (m model) MainMenuView() string {
	yellow := styles.TextYellow

	// Initialize []string of size equal to height of each "glyph".
	// This maintains consistent vertical spacing on title screen even when no glyphs are displayed.
	logo := make([]string, len(letters['f']))

	switch m.size {
	case large:
		a, _ := m.anim_mgr.Get(animations.TitleLogo)
		anim, ok := a.(*animations.TitleScreenLogoAnim)
		if !m.settings.Prefs.AnimationsEnabled || !ok {
			// Display yellow logo if animation state could not be retrieved
			for _, ch := range constants.FULL_GAME_TITLE {
				logo = drawGlyph(byte(ch), logo, yellow)
			}

			logo = append([]string{"\n", "\n"}, logo...) // prepend top padding
			logo = append(logo, "\n\n\n") // append bottom padding
			logo = append(logo, m.PressPlayView())

			return lipgloss.JoinVertical(lipgloss.Center, logo...)
		}

		switch anim.Phase {
		case animations.AbbreviatedTitlePhase:
			for _, ch := range constants.ABBR_GAME_TITLE {
				logo = drawGlyph(byte(ch), logo, yellow)
			}

		case animations.TypingFullTitlePhase, animations.FullTitlePausePhase:
			base := styles.TextBody
			highlight := styles.TextHighlight

			prompt_idx := 0
			for i := range anim.TypedLetters {
				current_title_char := constants.FULL_GAME_TITLE[i]

				style := base
				for j := prompt_idx; j < len(constants.ABBR_GAME_TITLE); j++ {
					ch := constants.ABBR_GAME_TITLE[j]

					is_prompt_letter := ch == current_title_char
					if is_prompt_letter {
						style = highlight
						prompt_idx++
						break
					}
				}

				logo = drawGlyph(byte(current_title_char), logo, style)
			}

		case animations.FullTitleRainbowScrollPhase:
			for i, ch := range constants.FULL_GAME_TITLE {
				style_idx := (anim.ColorIdx + i + len(anim.Colors)) % len(anim.Colors)
				style := anim.Colors[style_idx]
				logo = drawGlyph(byte(ch), logo, style)
			}

		case animations.TitleResetPhase:
			// Do nothing; logo hidden before anim restarts
		}

	default:
		for _, c := range constants.ABBR_GAME_TITLE {
			logo = drawGlyph(byte(c), logo, yellow)
		}
	}

	logo = append([]string{"\n\n\n"}, logo...) // prepend top padding
	logo = append(logo, "\n\n\n\n") // append bottom padding
	logo = append(logo, m.PressPlayView())

	return lipgloss.JoinVertical(
		lipgloss.Center,
		logo...
	)
}

func drawGlyph(char byte, logo []string, style lipgloss.Style) []string {
	char_glyph := letters[char]
	for i, line := range char_glyph {
		logo[i] = logo[i] + style.Render(line) + " "
	}
	return logo
}
