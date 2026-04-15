package tui

import (
	"encoding/json"
	"fzwds/src/enums"
	"fzwds/src/game"
	"fzwds/src/tui/styles"
	"fzwds/src/utils"
	"log/slog"
	"os"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SettingsState struct {
	selected 		int
	category		SettingsMenuCategory
	schemaList		[]game.SettingsSchemaItem
	lastSel			map[SettingsMenuCategory]int
	title			string
}

type SettingsMenuCategory string
const (
	preferences	 = "preferences"
	game_settings = "game_settings"
)

func (m model) SettingsSwitch(category SettingsMenuCategory) (model, tea.Cmd) {
	m = m.SwitchPage(settings_page)
	m.state.settings.category = category
	m.state.settings.selected = m.state.settings.lastSel[category]

	switch category {
	case preferences:
		m.state.settings.schemaList = m.app_settings_schema.Prefs
		m.state.settings.title = "General Preferences"

		m.footer_keymaps = []FooterKeymap {
			{key: "↑/↓", value: "scroll"},
			{key: "←/→", value: "change"},
			{key: "enter", value: "save"},
			{key: "ctrl+d", value: "defaults"},
			{key: "m", value: "menu"},
		}
	case game_settings:
		m.state.settings.schemaList = m.app_settings_schema.Game
		m.state.settings.title = "Game Settings"

		m.footer_keymaps = []FooterKeymap {
			{key: "↑/↓", value: "scroll"},
			{key: "←/→", value: "change"},
			{key: "enter", value: "play"},
			{key: "ctrl+d", value: "defaults"},
			{key: "m", value: "menu"},
		}
	}

	// TODO: presets in header? beginner/med/hard/expert
	// preset text highlighted if currently selected settings match a preset
	// would need custom equals function

	return m, nil
}

func (m model) SettingsUpdate(msg tea.Msg) (model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down", "tab":
			m.state.settings.selected = (m.state.settings.selected + 1 + len(m.state.settings.schemaList)) % len(m.state.settings.schemaList)
			m.state.settings.lastSel[m.state.settings.category] = m.state.settings.selected
			if m.state.settings.selected == 0 {
				m.goto_top = true
			}

		case "k", "up", "shift+tab":
			m.state.settings.selected = (m.state.settings.selected - 1 + len(m.state.settings.schemaList)) % len(m.state.settings.schemaList)
			m.state.settings.lastSel[m.state.settings.category] = m.state.settings.selected
			if m.state.settings.selected == len(m.state.settings.schemaList) - 1 {
				m.goto_bottom = true
			}

		case "+", "=", "right", "l":
			setting := m.state.settings.schemaList[m.state.settings.selected]
			is_bell_being_enabled := setting.PropName == "BellEnabled" && !m.app_settings_copy.Prefs.BellEnabled

			m.changeCurrentSetting(Next, m.state.settings.schemaList)

			if is_bell_being_enabled {
				cmds = append(cmds, m.terminalBellCmd(true))
			}

		case "-", "left", "h":
			setting := m.state.settings.schemaList[m.state.settings.selected]
			is_bell_being_enabled := setting.PropName == "BellEnabled" && !m.app_settings_copy.Prefs.BellEnabled

			m.changeCurrentSetting(Prev, m.state.settings.schemaList)

			if is_bell_being_enabled {
				cmds = append(cmds, m.terminalBellCmd(true))
			}

		case "ctrl+d":
			switch m.state.settings.category {
			case preferences:
				m.app_settings_copy.Prefs = game.GetDefaultGeneralPreferences()
			case game_settings:
				m.app_settings_copy.Game = game.GetDefaultGameSettings()
			}

		case "enter":
			m.app_settings = &m.app_settings_copy

			// TODO: return m, cmd that updates game_settings
			// TODO: abstract this save logic to common func shared with root initalization of settings
			marshaled_settings, err := json.MarshalIndent(m.app_settings, "", "    ")
			if err != nil {
				slog.Error("Error marshaling validated settings JSON", "error", err)
			}

			if err := os.WriteFile(m.app_settings_path, marshaled_settings, 0644); err != nil {
				slog.Error("Error writing settings.json", "error", err)
			}

			switch m.state.settings.category {
			case preferences:
				m.anim_mgr.SetAnimationStatus(m.app_settings.Prefs.AnimationsEnabled)
				return m.MainMenuSwitch()
			case game_settings:
				if m.app_settings.Game.Dictionary == enums.Pokemon {
					return m.PokemonGenSelectorSwitch()
				} else {
					return m.GameSwitch()
				}
			}

		case "m", "esc":
			m.app_settings_copy = *m.app_settings
			return m.MainMenuSwitch()
		}
	}

	return m, tea.Batch(cmds...)
}

func (m model) SettingsView() string {
	base := styles.TextBody.Render
	dim := styles.TextDim.Render
	accent := styles.TextAccent.Bold(true).Render

	var lines []string
	for i, setting := range m.state.settings.schemaList {
		if setting.Disabled {
			continue
		}

		var default_val, sub_desc string
		current_val := m.app_settings_copy.GetSetting(setting.PropName)

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
		is_selected := m.state.settings.selected == i

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

		row_1_space := m.width_content - lipgloss.Width(display_name) - lipgloss.Width(default_text) - 3

		var content string
		description := setting.Description

		// Show description for selected item only
		if is_selected && setting.Description != "" {
			row_2_space := m.width_content - lipgloss.Width(description) - lipgloss.Width(sub_desc) - 5
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
		apply_bottom_border := i != len(m.state.settings.schemaList) - 1
		line := styles.CreateSettingsMenuItem(content, is_selected, apply_bottom_border, m.width_content - 2)
		lines = append(lines, line)
	}

	return lipgloss.JoinVertical(lipgloss.Center, lines...)
}

type Direction int
const (
	Next Direction = 1
	Prev Direction = -1
)

func (m *model) changeCurrentSetting(dir Direction, schema []game.SettingsSchemaItem) {
	if m.state.settings.selected < 0 || m.state.settings.selected >= len(schema) {
		return
	}

	dir_int := int(dir)

	setting := schema[m.state.settings.selected]
	cur_val := m.app_settings_copy.GetSetting(setting.PropName)
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

	m.app_settings_copy.SetSetting(setting.PropName, new_val, m.app_settings_schema)
}
