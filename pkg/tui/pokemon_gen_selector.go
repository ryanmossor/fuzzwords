package tui

import (
	"fmt"
	"fzwds/pkg/dictionary"
	"fzwds/pkg/game"
	"fzwds/pkg/tui/styles"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type pokemonMenuState struct {
	selected 	int
	genList		[]int
	genState	map[int]bool
}

func (m model) PokemonGenSelectorSwitch() (model, tea.Cmd) {
	m = m.SwitchPage(pokemonGenMenuPage)
	m.state.pokemonMenu.selected = 1

	m.footerKeymaps = []footerKeymap {
		{key: "↑/↓", value: "scroll"},
		{key: "←/→", value: "change"},
		{key: "esc", value: "back"},
		{key: "enter", value: "play"},
	}

	return m, nil
}

func (m model) PokemonGenSelectorUpdate(msg tea.Msg) (model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		m.state.footer.footerMsg = ""

		switch msg.String() {
		case "j", "down", "tab":
			if m.state.pokemonMenu.selected + 1 > len(dictionary.PokemonDictionary) {
				m.state.pokemonMenu.selected = 1
			} else {
				m.state.pokemonMenu.selected += 1
			}

			if m.state.pokemonMenu.selected == 1 {
				m.gotoTop = true
			}

		case "k", "up", "shift+tab":
			if m.state.pokemonMenu.selected - 1 <= 0 {
				m.state.pokemonMenu.selected = len(dictionary.PokemonDictionary)
			} else {
				m.state.pokemonMenu.selected -= 1
			}

			if m.state.pokemonMenu.selected == len(dictionary.PokemonDictionary) - 1 {
				m.gotoBottom = true
			}

		case "+", "=", "right", "l", "-", "left", "h":
			idx := m.state.pokemonMenu.selected
			m.state.pokemonMenu.genState[idx] = !m.state.pokemonMenu.genState[idx]

		case "enter":
			selected_gens := make([]int, 0, len(dictionary.PokemonDictionary))
			for gen, enabled := range m.state.pokemonMenu.genState {
				if enabled {
					selected_gens = append(selected_gens, gen)
				}
			}

			if len(selected_gens) == 0 {
				m.state.footer.footerMsg = styles.
					TextRed.
					Render("You must select at least one generation")
				return m, nil
			}

			m.state.pokemonMenu.genList = selected_gens
			m.settings.Game.PokemonGens = selected_gens

			cmds = append(cmds, m.saveSettingsCmd(*m.settings, m.settingsPath))

			var cmd tea.Cmd
			m, cmd = m.GameSwitch()
			cmds = append(cmds, cmd)

		case "esc":
			return m.SettingsSwitch(gameSettings)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m model) PokemonGenSelectorView() string {
	base := styles.TextBody.Render
	dim := styles.TextDim.Render
	accent := styles.TextAccent.Bold(true).Render

	var lines []string
	for i := range len(dictionary.PokemonDictionary) {
		gen := i + 1

		cur_val := m.state.pokemonMenu.genState[gen]
		var enabled_text string
		if cur_val == true {
			enabled_text = "on"
		} else {
			enabled_text = "off"
		}

		display_name := fmt.Sprintf(" Gen %d", gen)
		row_text := dim("  " + enabled_text + "  ")
		is_selected := m.state.pokemonMenu.selected == gen

		if is_selected {
			row_text = accent("◀ " + enabled_text + " ▶")
		}

		// TODO: better way of calculating width (eg max 50% of width container?)
		row_space := m.contentWidth - lipgloss.Width(display_name) - lipgloss.Width(row_text) - 3 - 26
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
		line := styles.CreatePokemonMenuItem(content, is_selected, apply_bottom_border, m.contentWidth - 26)
		lines = append(lines, line)
	}

	return lipgloss.JoinVertical(lipgloss.Center, lines...)
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
