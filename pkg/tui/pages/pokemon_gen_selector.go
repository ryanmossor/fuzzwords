package pages

import (
	"fmt"
	"fzwds/pkg/dictionary"
	"fzwds/pkg/game"
	"fzwds/pkg/tui/commands"
	"fzwds/pkg/tui/figurethisout"
	"fzwds/pkg/tui/styles"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type PokemonMenuPage struct {
	name		PageName
	uiContext 	*figurethisout.UIContext
	helpKeys	[]figurethisout.HelpKeymap

	selected 	int
	genList		[]int
	genState	map[int]bool
}

func NewPokemonMenuPage(uiContext *figurethisout.UIContext) Page {
	return &PokemonMenuPage {
		name: About,
		uiContext: uiContext,
		helpKeys: []figurethisout.HelpKeymap {
			{Key: "↑/↓", Value: "scroll"},
			{Key: "←/→", Value: "change"},
			{Key: "esc", Value: "back"},
			{Key: "enter", Value: "play"},
		},
		selected: 1,
		genList: []int{},
		genState: initSelectedPokemonGens(uiContext.Settings),
	}
}

func (p PokemonMenuPage) GetPageName() PageName {
	return p.name
}

func (p PokemonMenuPage) Switch() tea.Cmd {
	return nil
}

func (p PokemonMenuPage) Update(msg tea.Msg) (Page, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// TODO re enable footer msg after i decide where it should live
		// m.state.footer.footerMsg = ""

		switch msg.String() {
		case "j", "down", "tab":
			if p.selected + 1 > len(dictionary.PokemonDictionary) {
				p.selected = 1
			} else {
				p.selected += 1
			}

			// if m.state.pokemonMenu.selected == 1 {
			// 	m.gotoTop = true
			// }

		case "k", "up", "shift+tab":
			if p.selected - 1 <= 0 {
				p.selected = len(dictionary.PokemonDictionary)
			} else {
				p.selected -= 1
			}

			// if p.selected == len(dictionary.PokemonDictionary) - 1 {
			// 	m.gotoBottom = true
			// }

		case "+", "=", "right", "l", "-", "left", "h":
			idx := p.selected
			p.genState[idx] = !p.genState[idx]

		case "enter":
			selected_gens := make([]int, 0, len(dictionary.PokemonDictionary))
			for gen, enabled := range p.genState {
				if enabled {
					selected_gens = append(selected_gens, gen)
				}
			}

			if len(selected_gens) == 0 {
				// m.state.footer.footerMsg = styles.
				// 	TextRed.
				// 	Render("You must select at least one generation")
				return p, nil
			}

			p.genList = selected_gens
			p.uiContext.Settings.Game.PokemonGens = selected_gens

			cmds = append(cmds,
				commands.SaveSettingsCmd(*p.uiContext.Settings, p.uiContext.SettingsPath),
				SwitchPageCmd(Game),
			)

		case "esc":
			return p, SwitchPageCmd(Settings)
		}
	}

	return p, tea.Batch(cmds...)
}

func (p PokemonMenuPage) View() string {
	base := styles.TextBody.Render
	dim := styles.TextDim.Render
	accent := styles.TextAccent.Bold(true).Render

	var lines []string
	for i := range len(dictionary.PokemonDictionary) {
		gen := i + 1

		cur_val := p.genState[gen]
		var enabled_text string
		if cur_val == true {
			enabled_text = "on"
		} else {
			enabled_text = "off"
		}

		display_name := fmt.Sprintf(" Gen %d", gen)
		row_text := dim("  " + enabled_text + "  ")
		is_selected := p.selected == gen

		if is_selected {
			row_text = accent("◀ " + enabled_text + " ▶")
		}

		// TODO: better way of calculating width (eg max 50% of width container?)
		row_space := p.uiContext.ContentWidth - lipgloss.Width(display_name) - lipgloss.Width(row_text) - 3 - 26
		row_space = max(0, row_space)

		gen_len := len(dictionary.PokemonDictionary[gen])
		desc := fmt.Sprintf(" %s - %s",
			dictionary.PokemonDictionary[gen][0],
			dictionary.PokemonDictionary[gen][gen_len - 1])

		var content string

		if is_selected {
			content = lipgloss.JoinVertical(
				lipgloss.Left,
				lipgloss.JoinHorizontal(
					lipgloss.Center,
					accent(display_name),
					strings.Repeat(" ", row_space),
					row_text),
				base(desc))
		} else {
			content = lipgloss.JoinHorizontal(
				lipgloss.Center,
				dim(display_name),
				strings.Repeat(" ", row_space),
				row_text)
		}

		apply_bottom_border := gen != len(dictionary.PokemonDictionary)
		// TODO: better width calculation
		line := styles.CreatePokemonMenuItem(content, is_selected, apply_bottom_border, p.uiContext.ContentHeight - 26)
		lines = append(lines, line)
	}

	return lipgloss.NewStyle().AlignVertical(lipgloss.Center).Height(p.uiContext.ContentHeight).Render(
		lipgloss.JoinVertical(lipgloss.Center, lines...),
	)
}

func initSelectedPokemonGens(settings *game.Settings) map[int]bool {
	gen_map := make(map[int]bool)
	for _, gen := range settings.Game.PokemonGens {
		if _, ok := dictionary.PokemonDictionary[gen]; ok {
			gen_map[gen] = true
		}
	}
	return gen_map
}
