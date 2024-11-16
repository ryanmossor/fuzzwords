package tui

import (
	"github.com/charmbracelet/lipgloss"
)

type theme struct {
	renderer *lipgloss.Renderer

	border     	lipgloss.TerminalColor
	background 	lipgloss.TerminalColor
	highlight  	lipgloss.TerminalColor
	red      	lipgloss.TerminalColor
	body       	lipgloss.TerminalColor
	accent     	lipgloss.TerminalColor
	yellow     	lipgloss.TerminalColor
	green      	lipgloss.TerminalColor
	blue		lipgloss.TerminalColor
	orange		lipgloss.TerminalColor
	dim			lipgloss.TerminalColor
	extra_dim	lipgloss.TerminalColor

	base lipgloss.Style
	// form *huh.Theme
}

func BasicTheme(renderer *lipgloss.Renderer) theme {
	base := theme{
		renderer: renderer,
	}

	// TODO: look into ANSI colors for increased compatibility
	base.background = lipgloss.AdaptiveColor{Dark: "#1E1E2E", Light: "#EFF1F5"}
	base.border = lipgloss.AdaptiveColor{Dark: "#585B70", Light: "#ACB0BE"} // Surface 2
	// base.body = lipgloss.AdaptiveColor{Dark: "#CDD6F4", Light: "#4C4F69"} // Text
	base.body = lipgloss.AdaptiveColor{Dark: "#A6ADC8", Light: "#6C6F85"} // Subtext 0
	base.accent = lipgloss.AdaptiveColor{Dark: "#FFFFFF", Light: "#11181C"}
	// base.accent = lipgloss.AdaptiveColor{Dark: "#CDD6F4", Light: "#4C4F69"} // Text
	base.yellow = lipgloss.AdaptiveColor{Dark: "#F9E2AF", Light: "#DF8E1D"}
	base.green = lipgloss.AdaptiveColor{Dark: "#A6E3A1", Light: "#40A02B"}
	base.blue = lipgloss.AdaptiveColor{Dark: "#74C7EC", Light: "#209FB5"}
	base.orange = lipgloss.AdaptiveColor{Dark: "#FAB387", Light: "#FE640B"}
	base.dim = lipgloss.AdaptiveColor{Dark: "#878787", Light: "#ACB0BE"} // not part of catppuccin palette
	// base.extra_dim = lipgloss.AdaptiveColor{Dark: "#585B70", Light: "#ACB0BE"} // Surface 2
	base.extra_dim = lipgloss.AdaptiveColor{Dark: "#6C7086", Light: "#ACB0BE"} // Overlay 0

	base.highlight = lipgloss.Color("#74C7EC")

	// base.error = lipgloss.Color("203")
	base.red = lipgloss.AdaptiveColor{Dark: "#F38BA8", Light: "#D20F39"}

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

func (b theme) TextYellow() lipgloss.Style {
	return b.Base().Foreground(b.yellow)
}

func (b theme) TextGreen() lipgloss.Style {
	return b.Base().Foreground(b.green)
}

func (b theme) TextOrange() lipgloss.Style {
	return b.Base().Foreground(b.orange)
}

func (b theme) TextBlue() lipgloss.Style {
	return b.Base().Foreground(b.blue)
}

func (b theme) TextRed() lipgloss.Style {
	return b.Base().Foreground(b.red)
}

func (b theme) TextDim() lipgloss.Style {
	return b.Base().Foreground(b.dim)
}

func (b theme) TextExtraDim() lipgloss.Style {
	return b.Base().Foreground(b.extra_dim)
}

// func (b theme) Form() *huh.Theme {
// 	return b.form
// }

func (b theme) Border() lipgloss.TerminalColor {
	return b.dim
}
