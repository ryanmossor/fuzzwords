package tui

import (
	"github.com/charmbracelet/lipgloss"
)

type theme struct {
	renderer *lipgloss.Renderer

	border     lipgloss.TerminalColor
	background lipgloss.TerminalColor
	highlight  lipgloss.TerminalColor
	error      lipgloss.TerminalColor
	body       lipgloss.TerminalColor
	accent     lipgloss.TerminalColor

	base lipgloss.Style
	// form *huh.Theme
}

func BasicTheme(renderer *lipgloss.Renderer, highlight *string) theme {
	base := theme{
		renderer: renderer,
	}

	base.background = lipgloss.AdaptiveColor{Dark: "#000000", Light: "#FBFCFD"}
	base.border = lipgloss.AdaptiveColor{Dark: "#3A3F42", Light: "#D7DBDF"}
	base.body = lipgloss.AdaptiveColor{Dark: "#889096", Light: "#889096"}
	base.accent = lipgloss.AdaptiveColor{Dark: "#FFFFFF", Light: "#11181C"}
	if highlight != nil {
		base.highlight = lipgloss.Color(*highlight)
	} else {
		base.highlight = lipgloss.Color("#FF5C00")
	}
	base.error = lipgloss.Color("203")

	base.base = renderer.NewStyle().Foreground(base.body)
	// base.form = HuhTheme(base)

	return base
}

func (b theme) Body() lipgloss.TerminalColor {
	return b.body
}

func (b theme) Highlight() lipgloss.TerminalColor {
	return b.highlight
}

func (b theme) Background() lipgloss.TerminalColor {
	return b.background
}

func (b theme) Accent() lipgloss.TerminalColor {
	return b.accent
}

func (b theme) Base() lipgloss.Style {
	return b.base.Copy()
}

func (b theme) TextBody() lipgloss.Style {
	return b.Base().Foreground(b.body)
}

func (b theme) TextAccent() lipgloss.Style {
	return b.Base().Foreground(b.accent)
}

func (b theme) TextHighlight() lipgloss.Style {
	return b.Base().Foreground(b.highlight)
}

func (b theme) TextError() lipgloss.Style {
	return b.Base().Foreground(b.error)
}

// func (b theme) Form() *huh.Theme {
// 	return b.form
// }

func (b theme) Border() lipgloss.TerminalColor {
	return b.border
}
