package tui

import (
	"encoding/json"
	"fzwds/src/game"
	"fzwds/src/utils"
	"log/slog"
	"os"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SettingsState struct {
	selected		int
}

func (m model) SettingsSwitch() (model, tea.Cmd) {
	m = m.SwitchPage(settings_page)
	m.state.settings.selected = 0

	m.footer_keymaps = []footer_keymaps{
		{key: "↑/↓", value: "scroll"},
		{key: "←/→", value: "change"},
		{key: "ctrl+r", value: "defaults"},
        {key: "m", value: "main menu"},
		{key: "enter", value: "play"},
	}

	return m, nil
}

func (m model) SettingsUpdate(msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down", "tab":
			if m.state.settings.selected < len(m.settings_schema) - 1 {
				m.state.settings.selected++
			}
		case "k", "up", "shift+tab":
			if m.state.settings.selected > 0 {
				m.state.settings.selected--
			}
		case "+", "=", "right", "l":
			current_setting := m.settings_schema[m.state.settings.selected]
			current_val := m.game_settings_copy.GetSetting(current_setting.PropName)
			var new_val any
			if len(current_setting.ValidValues) > 0 {
				current_val_idx := -1
				for i, val := range current_setting.ValidValues {
					if utils.ValuesEqual(val.Value, current_val) {
						current_val_idx = i
						break
					}
				}
				var next_idx int
				if current_setting.Type == "int" {
					next_idx = min(current_val_idx + 1, len(current_setting.ValidValues) - 1)
				} else {
					// % len allows for circular indexing -- wraps around if current_val_idx + 1 < 0
					next_idx = (current_val_idx + 1) % len(current_setting.ValidValues)
				}
				new_val = current_setting.ValidValues[next_idx].Value
			} else if current_setting.Type == "int" {
				v := current_val.(int)
				new_val = v + 1
				if current_setting.Max != nil && new_val.(int) > *current_setting.Max {
					new_val = *current_setting.Max
				}
			}
			m.game_settings_copy.SetSetting(current_setting.PropName, new_val, m.settings_schema)
		case "-", "left", "h": 
			current_setting := m.settings_schema[m.state.settings.selected]
			current_val := m.game_settings_copy.GetSetting(current_setting.PropName)
			var new_val any
			if len(current_setting.ValidValues) > 0 {
				current_val_idx := -1
				for i, val := range current_setting.ValidValues {
					if utils.ValuesEqual(val.Value, current_val) {
						current_val_idx = i
						break
					}
				}

				var prev_idx int
				if current_setting.Type == "int" {
					prev_idx = max(current_val_idx - 1, 0)
				} else {
					// Adding arr len and calculating % allows for circular indexing -- wraps around if current_val_idx - 1 < 0
					prev_idx = (current_val_idx - 1 + len(current_setting.ValidValues)) % len(current_setting.ValidValues)
				}

				new_val = current_setting.ValidValues[prev_idx].Value
			} else if current_setting.Type == "int" {
				v := current_val.(int)
				new_val = v - 1
				if current_setting.Min != nil && new_val.(int) < *current_setting.Min {
					new_val = *current_setting.Min
				}
			}
			m.game_settings_copy.SetSetting(current_setting.PropName, new_val, m.settings_schema)
		case "ctrl+r":
			m.game_settings_copy = game.InitializeSettings()
		case "b":
			// TODO: beginner preset
		// case "m":
			// TODO: medium preset
		case "d":
			// TODO: difficult preset
		case "x":
			// TODO: expert preset
		case "enter":
			m.game_settings = &m.game_settings_copy

			// TODO: return m, cmd that updates game_settings
			// TODO: abstract this save logic to common func shared with root initalization of settings
			marshaled_settings, err := json.MarshalIndent(m.game_settings, "", "    ")
			if err != nil {
				slog.Error("Error marshaling validated settings JSON", "error", err)
			}

			if err := os.WriteFile(m.settings_path, marshaled_settings, 0644); err != nil {
                slog.Error("Error writing settings.json", "error", err)
			}

            return m.GameSwitch()
		case "m", "esc":
			return m.MainMenuSwitch()
		}
	}

	return m, nil
}

func (m model) SettingsView() string {
	base := m.theme.Base().Render
	accent := m.theme.TextAccent().Render 

	var lines []string

	for i, setting := range m.settings_schema {
		if setting.Disabled {
			continue
		}

		var default_val, sub_desc string
		current_val := m.game_settings_copy.GetSetting(setting.PropName)

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
					break
				}
			}
		}
		default_text := accent("  " + default_val + "    ")
		
		if m.state.settings.selected == i {
			setting_val_int, err := strconv.Atoi(default_val)
			if err != nil {
				default_text = accent("← " + default_val + " →  ")
			} else if setting_val_int == *setting.Max {
				default_text = accent("← " + default_val + "    ")
			} else if setting_val_int == *setting.Min {
				default_text = accent("  " + default_val + " →  ")
			} else {
				default_text = accent("← " + default_val + " →  ")
			}
		}

		display_name := accent(setting.DisplayName)
		row_1_space := m.width_content - lipgloss.Width(display_name) - lipgloss.Width(default_text) - 3

		var content string
		description := setting.Description

		// i == selected expands desc only for selected
		if i == m.state.settings.selected && setting.Description != "" {
			row_2_space := m.width_content - lipgloss.Width(description) - lipgloss.Width(sub_desc) - 5
			var row_2 string

			if setting.Description != "n/a" {
				row_2 = lipgloss.JoinHorizontal(
					lipgloss.Top,
					description,
					m.theme.Base().Width(row_2_space).Render(),
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
					m.theme.Base().Width(row_1_space).Render(),
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
					m.theme.Base().Width(row_1_space).Render(),
					default_text,
				),
			)
		}

		line := m.CreateBox(content, i == m.state.settings.selected)
		lines = append(lines, line)
	} 

	return m.theme.Base().Render(lipgloss.JoinVertical(
		lipgloss.Left,
		lines...,
	)) 
}
