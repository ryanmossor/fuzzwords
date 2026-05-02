package pages

import (
	"fzwds/pkg/enums"
	"fzwds/pkg/game"
	"fzwds/pkg/tui/commands"
	"fzwds/pkg/tui/figurethisout"
	"fzwds/pkg/tui/styles"
	"fzwds/pkg/utils"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type settingsMenuCategory string
const (
	preferences	 = "preferences"
	gameSettings = "game_settings"
)

type SettingsPage struct {
	name			PageName
	uiContext 		*figurethisout.UIContext
	helpKeys		[]figurethisout.HelpKeymap
	selected 		int
	category		settingsMenuCategory
	schemaList		[]game.SettingsSchemaItem
	title			string
	settingsCopy	game.Settings
}

// TODO: pass settings in separately rather than via UIContext?
// would allow for separating Prefs from Game settings
func NewGameSettingsPage(uiContext *figurethisout.UIContext) Page {
	return &SettingsPage {
		name: Settings,
		uiContext: uiContext,
		helpKeys: []figurethisout.HelpKeymap {
			{Key: "↑/↓", Value: "scroll"},
			{Key: "←/→", Value: "change"},
			{Key: "enter", Value: "play"},
			{Key: "ctrl+d", Value: "defaults"},
			{Key: "m", Value: "menu"},
		},
		category: gameSettings,
		selected: 0,
		schemaList: uiContext.Schema.Game,
		title: "Game Settings",
		settingsCopy: *uiContext.Settings,
	}
}

func NewPreferencesPage(uiContext *figurethisout.UIContext) Page {
	return &SettingsPage {
		name: Preferences,
		uiContext: uiContext,
		helpKeys: []figurethisout.HelpKeymap {
			{Key: "↑/↓", Value: "scroll"},
			{Key: "←/→", Value: "change"},
			{Key: "enter", Value: "save"},
			{Key: "ctrl+d", Value: "defaults"},
			{Key: "m", Value: "menu"},
		},
		// TODO: can category be removed too?
		category: preferences,
		selected: 0,
		schemaList: uiContext.Schema.Prefs,
		title: "General Preferences",
		settingsCopy: *uiContext.Settings,
	}
}

func (p SettingsPage) GetPageName() PageName {
	return p.name
}

func (p SettingsPage) Switch() tea.Cmd {
	// p.settingsCopy = *p.uiContext.Settings
	return nil
}

func (p *SettingsPage) Update(msg tea.Msg) (Page, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down", "tab":
			if p.selected == len(p.schemaList) - 1 {
				p.selected = 0
			} else {
				p.selected += 1
			}
			// if p.selected == 0 {
			// 	m.gotoTop = true
			// }

		case "k", "up", "shift+tab":
			if p.selected == 0 {
				p.selected = len(p.schemaList) - 1
			} else {
				p.selected -= 1
			}
			// if m.state.settings.selected == len(m.state.settings.schemaList) - 1 {
			// 	m.gotoBottom = true
			// }

		case "+", "=", "right", "l":
			setting := p.schemaList[p.selected]
			is_bell_being_enabled := setting.PropName == "BellEnabled" && !p.settingsCopy.Prefs.BellEnabled

			p.changeCurrentSetting(Next, p.schemaList)

			if is_bell_being_enabled {
				cmds = append(cmds, commands.TerminalBellCmd(p.uiContext.Settings.Prefs, true))
			}

		case "-", "left", "h":
			setting := p.schemaList[p.selected]
			is_bell_being_enabled := setting.PropName == "BellEnabled" && !p.settingsCopy.Prefs.BellEnabled

			p.changeCurrentSetting(Prev, p.schemaList)

			if is_bell_being_enabled {
				cmds = append(cmds, commands.TerminalBellCmd(p.uiContext.Settings.Prefs, true))
			}

		case "ctrl+d":
			switch p.category {
			case preferences:
				p.settingsCopy.Prefs = game.GetDefaultGeneralPreferences()
			case gameSettings:
				p.settingsCopy.Game = game.GetDefaultGameSettings()
			}

		case "enter":
			p.uiContext.Settings = &p.settingsCopy
			cmds = append(cmds, commands.SaveSettingsCmd(*p.uiContext.Settings, p.uiContext.SettingsPath))

			var cmd tea.Cmd
			switch p.category {
			case preferences:
				p.uiContext.AnimManager.SetAnimationStatus(p.uiContext.Settings.Prefs.AnimationsEnabled)
				cmd = SwitchPageCmd(Title)

			case gameSettings:
				if p.uiContext.Settings.Game.Dictionary == enums.Pokemon {
					cmd = SwitchPageCmd(PokemonGenMenu)
				} else {
					cmd = SwitchPageCmd(Game)
				}
			}
			cmds = append(cmds, cmd)

		case "m", "esc":
			p.settingsCopy = *p.uiContext.Settings
			cmds = append(cmds, SwitchPageCmd(Title))
		}
	}

	return p, tea.Batch(cmds...)
}

func (p SettingsPage) View() string {
	base := styles.TextBody.Render
	dim := styles.TextDim.Render
	accent := styles.TextAccent.Bold(true).Render

	var lines []string
	for i, setting := range p.schemaList {
		if setting.Disabled {
			continue
		}

		var default_val, sub_desc string
		current_val := p.settingsCopy.GetSetting(setting.PropName)

		switch setting.Type {
		case "int":
			default_val = strconv.Itoa(current_val.(int))
		case "enum", "string":
			default_val = current_val.(string)
		}

		if setting.ValidValues != nil {
			sub_desc = ""
			for _, val := range setting.ValidValues {
				if utils.ValuesEqual(val.Value, current_val) {
					sub_desc = val.Description
					if val.DisplayText != "" {
						default_val = val.DisplayText
					}
					break
				}
			}
		}
		default_text := dim("  " + default_val + "    ")
		var display_name string
		is_selected := p.selected == i

		if is_selected {
			display_name = accent(setting.DisplayName)
			setting_val_int, err := strconv.Atoi(default_val)
			if err != nil {
				default_text = accent("◀ " + default_val + " ▶  ")
			} else if setting_val_int == *setting.Max {
				default_text = accent("◀ " + default_val + "    ")
			} else if setting_val_int == *setting.Min {
				default_text = accent("  " + default_val + " ▶  ")
			} else {
				default_text = accent("◀ " + default_val + " ▶  ")
			}
		} else {
			display_name = dim(setting.DisplayName)
		}

		row_1_space := p.uiContext.ContentWidth - lipgloss.Width(display_name) - lipgloss.Width(default_text) - 3
		row_1_space = max(0, row_1_space)

		var content string
		description := setting.Description

		// Show description for selected item only
		if is_selected && setting.Description != "" {
			row_2_space := p.uiContext.ContentWidth - lipgloss.Width(description) - lipgloss.Width(sub_desc) - 5
			var row_2 string

			if setting.Description != "n/a" {
				row_2 = lipgloss.JoinHorizontal(
					lipgloss.Top,
					base(description),
					strings.Repeat(" ", row_2_space),
					base(sub_desc),
				)
			} else {
				row_2 = lipgloss.JoinHorizontal(
					lipgloss.Top,
					base(sub_desc),
				)
			}

			content = lipgloss.JoinVertical(
				lipgloss.Left,
				lipgloss.JoinHorizontal(
					lipgloss.Top,
					display_name,
					strings.Repeat(" ", row_1_space),
					default_text,
				),
				row_2,
			)
		} else {
			content = lipgloss.JoinVertical(
				lipgloss.Left,
				lipgloss.JoinHorizontal(
					lipgloss.Top,
					display_name,
					strings.Repeat(" ", row_1_space),
					default_text,
				),
			)
		}

		// Don't apply border to final setting box
		apply_bottom_border := i != len(p.schemaList) - 1
		line := styles.CreateSettingsMenuItem(content, is_selected, apply_bottom_border, p.uiContext.ContentWidth - 2)
		lines = append(lines, line)
	}

	return lipgloss.NewStyle().AlignVertical(lipgloss.Center).Height(p.uiContext.ContentHeight).Render(
		lipgloss.JoinVertical(lipgloss.Center, lines...),
	)
}

type Direction int
const (
	Next Direction = 1
	Prev Direction = -1
)

func (p *SettingsPage) changeCurrentSetting(dir Direction, schema []game.SettingsSchemaItem) {
	if p.selected < 0 || p.selected >= len(schema) {
		return
	}

	dir_int := int(dir)

	setting := schema[p.selected]
	cur_val := p.settingsCopy.GetSetting(setting.PropName)
	var new_val any

	if len(setting.ValidValues) > 0 {
		cur_idx := -1
		for i, opt := range setting.ValidValues {
			if utils.ValuesEqual(opt.Value, cur_val) {
				cur_idx = i
				break
			}
		}

		var next_idx int
		if setting.Type == "int" {
			// Linear (clamp)
			next_idx = max(cur_idx + dir_int, 0)
			if next_idx >= len(setting.ValidValues) {
				next_idx = len(setting.ValidValues) - 1
			}
		} else {
			// Circular wrap-around
			next_idx = (cur_idx + dir_int + len(setting.ValidValues)) % len(setting.ValidValues)
		}

		new_val = setting.ValidValues[next_idx].Value
	} else if setting.Type == "int" {
		new_val = cur_val.(int) + dir_int

		if setting.Min != nil && new_val.(int) < *setting.Min {
			new_val = *setting.Min
		}
		if setting.Max != nil && new_val.(int) > *setting.Max {
			new_val = *setting.Max
		}
	} else {
		return
	}

	p.settingsCopy.SetSetting(setting.PropName, new_val, p.uiContext.Schema)
}
