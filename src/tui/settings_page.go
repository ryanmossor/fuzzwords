package tui

import (
	"encoding/json"
	"fzwds/src/enums"
	"fzwds/src/game"
	"log/slog"
	"os"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type settingsState struct {
	selected		int
}

func (m model) SettingsSwitch() (model, tea.Cmd) {
	m = m.SwitchPage(settings_page)
	m.state.settings.selected = 0

	m.footer_cmds = []footerCmd{
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
			if m.state.settings.selected < len(m.settings_menu_json) - 1 {
				m.state.settings.selected++
			}
		case "k", "up", "shift+tab":
			if m.state.settings.selected > 0 {
				m.state.settings.selected--
			}
		case "+", "=", "right", "l":
			m.changeSetting(m.state.settings.selected, 1)
		case "-", "left", "h": 
			m.changeSetting(m.state.settings.selected, -1)
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
	// lines = append(lines, base("\nChange your ") + accent("game settings") + base(" here.\n"))

	for i, setting := range m.settings_menu_json {
		if setting.Disabled {
			continue
		}

		name := accent(setting.Name)

		default_val := strconv.Itoa(m.getIntSetting(setting.PropName))
		var sub_desc string
		if setting.ValidValues != nil {
			default_val, sub_desc = m.getStringSetting(setting.PropName)
		}
		default_text := base("  ") + accent(default_val) + base("    ")
		
		if m.state.settings.selected == i {
			default_text = accent("← " + default_val + " →  ")
		}
		row_1_space := m.width_content - lipgloss.Width(name) - lipgloss.Width(default_text) - 3

		var content string
		description := setting.Description

		// i == selected expands desc only for selected
		if i == m.state.settings.selected && setting.Description != "" {
			row_2_space := m.width_content - lipgloss.Width(description) - lipgloss.Width(sub_desc) - 5
			var row_2 string

			// if strings.TrimSpace(setting.Description) != "" {
			if setting.Description != "n/a" {
			// if sub_desc != "" {
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
					name,
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
					name,
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

func (m model) getIntSetting(propName string) int {
	switch propName {
	case "HealthInitial":
		return m.game_settings_copy.HealthInitial
	case "HealthMax":
		return m.game_settings_copy.HealthMax
	case "PromptLenMin":
		return m.game_settings_copy.PromptLenMin
	case "PromptLenMax":
		return m.game_settings_copy.PromptLenMax
	case "PromptStrikesMax":
		return m.game_settings_copy.PromptStrikesMax
	case "TurnDurationMin":
		return m.game_settings_copy.TurnDurationMin
	}

	return 0
}

func (m model) getStringSetting(propName string) (string, string) {
	var val string

	switch propName {
	case "Alphabet":
		val = m.game_settings_copy.Alphabet.String()
	case "PromptMode":
		val = m.game_settings_copy.PromptMode.String()
	case "WinCondition":
		val = m.game_settings_copy.WinCondition.String()
	}

	return val, m.getSubDescription(propName, val)
}

func (m *model) changeSetting(selected int, count int) {
	switch m.settings_menu_json[selected].PropName {
	case "Alphabet":
		alphabet_idx := int(m.game_settings_copy.Alphabet) + count
		m.game_settings_copy.SetAlphabet(alphabet_idx)
	case "HealthInitial":
		m.game_settings_copy.SetHealthInitial(m.game_settings_copy.HealthInitial + count)
	case "HealthMax":
		m.game_settings_copy.SetHealthMax(m.game_settings_copy.HealthMax + count)
	case "PromptLenMin":
		m.game_settings_copy.SetPromptLenMin(m.game_settings_copy.PromptLenMin + count)
	case "PromptLenMax":
		m.game_settings_copy.SetPromptLenMax(m.game_settings_copy.PromptLenMax + count)
	case "PromptMode":
		if m.game_settings_copy.PromptMode == enums.Fuzzy {
			m.game_settings_copy.SetPromptMode(enums.Classic.String())
		} else {
			m.game_settings_copy.SetPromptMode(enums.Fuzzy.String())
		}
	// case "PromptStrikesMax":
	// 	return m.game_settings.PromptStrikesMax
	// case "TurnDurationMin":
	// 	return m.game_settings.TurnDurationMin
	case "WinCondition":
		if m.game_settings_copy.WinCondition == enums.Endless {
			m.game_settings_copy.SetWinCondition(enums.MaxLives.String())
		} else {
			m.game_settings_copy.SetWinCondition(enums.Endless.String())
		}
	}
}

func (m model) getSubDescription(propName string, val string) string {
	var setting game.Config
	for _, s := range m.settings_menu_json {
		if s.PropName == propName {
			setting = s
			break
		}
	}

	for _, v := range setting.ValidValues {
		if v.Name == val {
			return v.Description
		}
	}

	return ""
}
