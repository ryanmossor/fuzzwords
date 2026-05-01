package pages

import (
	"fzwds/pkg/constants"
	"fzwds/pkg/game"
	"fzwds/pkg/tui/animations"
	"fzwds/pkg/tui/commands"
	"fzwds/pkg/tui/figurethisout"
	"fzwds/pkg/tui/styles"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	hidden  = styles.TextBody.Render("Press       to play")
	visible = styles.TextBody.Render("Press ") +
			  styles.TextAccent.Bold(true).Render("ENTER") +
			  styles.TextBody.Render(" to play")
)

type TitlePage struct {
	name 				PageName
	pressPlayVisible	bool
	pressPlayId 		int // used to ignore stale press play commands
	uiContext 			*figurethisout.UIContext
	settings			*game.Settings
	helpKeys			[]figurethisout.HelpKeymap
}

func NewTitlePage(uiContext *figurethisout.UIContext, settings *game.Settings) Page {
	return &TitlePage {
		name: 				Title,
		pressPlayVisible: 	true,
		pressPlayId: 		0,
		uiContext: 			uiContext,
		settings: 			settings,
		helpKeys: 			[]figurethisout.HelpKeymap {
			{Key: "ctrl+p", Value: "preferences"},
			{Key: "q", Value: "quit"},
		},
	}
}

func (p TitlePage) GetPageName() PageName {
	return p.name
}

type PressPlayTickMsg struct {
	Id		int
}
func (p TitlePage) pressPlayFlashCmd() tea.Cmd {
	if !p.settings.Prefs.AnimationsEnabled {
		return nil
	}
	return tea.Tick(850 * time.Millisecond, func(t time.Time) tea.Msg {
		return PressPlayTickMsg{ Id: p.pressPlayId }
	})
}

func (p *TitlePage) Switch() tea.Cmd {
	p.uiContext.AnimManager.InitAnimations(animations.TitleLogo)
	return tea.Batch(
		commands.TickCmd(p.uiContext.FPS),
		p.pressPlayFlashCmd(),
	)
}

func (p *TitlePage) Update(msg tea.Msg) (Page, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "s":
			return p, SwitchPageCmd(Stats)
		case "a":
			return p, SwitchPageCmd(About)
		case "enter":
			// return m.SettingsSwitch(gameSettings)
			return p, SwitchPageCmd(Settings)
		case "ctrl+p":
			// return m.SettingsSwitch(preferences)
			// TODO: preferences, settings pages each constructed w/ NewSettingsPage
			return p, SwitchPageCmd(Settings)
		case "q":
			return p, tea.Quit
		}

	case PressPlayTickMsg:
		// Drop stale tick messages from previous visits
		if msg.Id != p.pressPlayId {
			return p, nil
		}
		p.pressPlayVisible = !p.pressPlayVisible
		p.pressPlayId += 1
		return p, p.pressPlayFlashCmd()
	}

	return p, nil
}

func (p TitlePage) View() string {
	yellow := styles.TextYellow

	// Initialize []string of size equal to height of each "glyph".
	// This maintains consistent vertical spacing on title screen even when no glyphs are displayed.
	logo := make([]string, len(letters['f']))

	switch p.uiContext.Size {
	case figurethisout.Large:
		a, _ := p.uiContext.AnimManager.Get(animations.TitleLogo)
		anim, ok := a.(*animations.TitleScreenLogoAnim)
		if !p.settings.Prefs.AnimationsEnabled || !ok {
			// Display yellow logo if animation state could not be retrieved
			for _, ch := range constants.FULL_GAME_TITLE {
				logo = drawGlyph(byte(ch), logo, yellow)
			}

			logo = append([]string{"\n", "\n"}, logo...) // prepend top padding
			logo = append(logo, "\n\n\n") // append bottom padding
			logo = append(logo, p.pressPlayText())

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
	logo = append(logo, p.pressPlayText())

	return lipgloss.JoinVertical(
		lipgloss.Center,
		logo...
	)
}

func (p TitlePage) pressPlayText() string {
	if p.settings.Prefs.AnimationsEnabled && !p.pressPlayVisible {
		return hidden
	}
	return visible
}

func drawGlyph(char byte, logo []string, style lipgloss.Style) []string {
	char_glyph := letters[char]
	for i, line := range char_glyph {
		logo[i] = logo[i] + style.Render(line) + " "
	}
	return logo
}

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
