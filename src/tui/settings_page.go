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
	selected int
}

func (m model) SettingsSwitch() (model, tea.Cmd) {
	m = m.SwitchPage(settings_page)
	m.state.settings.selected = 0

	m.footer_keymaps = []footer_keymaps{
		{key: "↑/↓", value: "scroll"},
		{key: "←/→", value: "change"},
		{key: "ctrl+r", value: "reset defaults"},
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
			m.state.settings.selected = (m.state.settings.selected + 1 + len(m.settings_schema)) % len(m.settings_schema)
			if m.state.settings.selected == 0 {
				m.goto_top = true
			}
		case "k", "up", "shift+tab":
			m.state.settings.selected = (m.state.settings.selected - 1 + len(m.settings_schema)) % len(m.settings_schema)
			if m.state.settings.selected == len(m.settings_schema) - 1 {
				m.goto_bottom = true
			}
		case "+", "=", "right", "l":
			m.changeCurrentSetting(Next)
		case "-", "left", "h":
			m.changeCurrentSetting(Prev)
		case "ctrl+r":
			m.game_settings_copy = game.GetDefaultSettings()
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
	dim := m.theme.TextExtraDim().Render
	accent := m.theme.TextAccent().Bold(true).Render

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
		default_text := dim("  " + default_val + "    ")
		var display_name string

		if m.state.settings.selected == i {
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

		// Don't apply border to final setting box
		apply_bottom_border := i != len(m.settings_schema) - 1
		line := m.CreateSettingsMenuItem(content, i == m.state.settings.selected, apply_bottom_border)
		lines = append(lines, line)
	}

	return m.theme.Base().Render(lipgloss.JoinVertical(
		lipgloss.Center,
		lines...,
	))
}

type Direction int
const (
	Next Direction = 1
	Prev Direction = -1
)

func (m *model) changeCurrentSetting(dir Direction) {
	if m.state.settings.selected < 0 || m.state.settings.selected >= len(m.settings_schema) {
		return
	}

	dir_int := int(dir)

	setting := m.settings_schema[m.state.settings.selected]
	cur_val := m.game_settings_copy.GetSetting(setting.PropName)
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

	m.game_settings_copy.SetSetting(setting.PropName, new_val, m.settings_schema)
}
