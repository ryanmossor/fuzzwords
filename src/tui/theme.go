package tui

import (
	"image/color"
	"os"

	"charm.land/lipgloss/v2"
)

type theme struct {
	// renderer *lipgloss.Renderer

	border     	color.Color
	background 	color.Color
	highlight  	color.Color
	body       	color.Color
	accent     	color.Color
	dim			color.Color
	lavender	color.Color
	input_bg	color.Color

	red      	color.Color
	orange		color.Color
	yellow     	color.Color
	green      	color.Color
	blue		color.Color
	indigo		color.Color
	purple		color.Color

	base lipgloss.Style
}

// func BasicTheme(renderer *lipgloss.Renderer) theme {
func BasicTheme() theme {
	// base := theme{
	// 	renderer: renderer,
	// }

	hasDark := lipgloss.HasDarkBackground(os.Stdin, os.Stdout)
	lightDark := lipgloss.LightDark(hasDark)
	// color := lightDark(lipgloss.Color("#0000ff"), lipgloss.Color("#000099"))

	// TODO: look into ANSI colors for increased compatibility
	return theme {
		background: lightDark(lipgloss.Color("#EFF1F5"), lipgloss.Color("#1E1E2E")),
		border: lightDark(lipgloss.Color("#ACB0BE"), lipgloss.Color("#585B70")), // Surface 2
		body: lightDark(lipgloss.Color("#6C6F85"), lipgloss.Color("#A6ADC8")), // Subtext 0
		accent: lightDark(lipgloss.Color("#11181C"), lipgloss.Color("#FFFFFF")),
		dim: lightDark(lipgloss.Color("#ACB0BE"), lipgloss.Color("#6C7086")), // Overlay 0
		lavender: lightDark(lipgloss.Color("#7287FD"), lipgloss.Color("#B4BEFE")),
		input_bg: lightDark(lipgloss.Color("#CCD0DA"), lipgloss.Color("#45475A")), // Surface 1
		// input_bg: lightDark(lipgloss.Color("#CCD0DA"), lipgloss.Color("#313244")), // Surface 0

		highlight: lightDark(lipgloss.Color("#209FB5"), lipgloss.Color("#74C7EC")),

		red: lightDark(lipgloss.Color("#D20F39"), lipgloss.Color("#F38BA8")),
		orange: lightDark(lipgloss.Color("#FE640B"), lipgloss.Color("#FAB387")),
		yellow: lightDark(lipgloss.Color("#DF8E1D"), lipgloss.Color("#F9E2AF")),
		green: lightDark(lipgloss.Color("#40A02B"), lipgloss.Color("#A6E3A1")),
		blue: lightDark(lipgloss.Color("#209FB5"), lipgloss.Color("#74C7EC")),
		indigo: lightDark(lipgloss.Color("#8839EF"), lipgloss.Color("#7287FD")), // swapped indigo/purple light values
		purple: lightDark(lipgloss.Color("#7287FD"), lipgloss.Color("#CBA6F7")),
	}

	// base.base = renderer.NewStyle().Foreground(base.body)

	// return base
}

func (b theme) Body() color.Color {
	return b.body
}

func (b theme) Highlight() color.Color {
	return b.highlight
}

func (b theme) Background() color.Color {
	return b.background
}

func (b theme) Accent() color.Color {
	return b.accent
}

func (b theme) Base() lipgloss.Style {
	return b.base
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

func (b theme) TextRed() lipgloss.Style {
	return b.Base().Foreground(b.red)
}

func (b theme) TextOrange() lipgloss.Style {
	return b.Base().Foreground(b.orange)
}

func (b theme) TextYellow() lipgloss.Style {
	return b.Base().Foreground(b.yellow)
}

func (b theme) TextGreen() lipgloss.Style {
	return b.Base().Foreground(b.green)
}

func (b theme) TextBlue() lipgloss.Style {
	return b.Base().Foreground(b.blue)
}

func (b theme) TextIndigo() lipgloss.Style {
	return b.Base().Foreground(b.indigo)
}

func (b theme) TextPurple() lipgloss.Style {
	return b.Base().Foreground(b.purple)
}

func (b theme) TextLavender() lipgloss.Style {
	return b.Base().Foreground(b.lavender)
}

func (b theme) TextDim() lipgloss.Style {
	return b.Base().Foreground(b.dim)
}

func (b theme) Border() color.Color {
	return b.border
}

func (b theme) GetRainbowColors() []lipgloss.Style {
	return []lipgloss.Style{
		b.TextRed(),
		b.TextOrange(),
		b.TextYellow(),
		b.TextGreen(),
		b.TextBlue(),
		b.TextIndigo(),
		b.TextPurple(),
	}
}
