package theme

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	Black      = lipgloss.Color("#000000")
	White      = lipgloss.Color("#ffffff")

	Background = lipgloss.AdaptiveColor{Dark: "#1E1E2E", Light: "#EFF1F5"}
	Border     = lipgloss.AdaptiveColor{Dark: "#585B70", Light: "#ACB0BE"} // Surface 2
	Body       = lipgloss.AdaptiveColor{Dark: "#A6ADC8", Light: "#6C6F85"} // Subtext 0
	Accent     = lipgloss.AdaptiveColor{Dark: "#FFFFFF", Light: "#11181C"}
	Dim        = lipgloss.AdaptiveColor{Dark: "#6C7086", Light: "#ACB0BE"} // Overlay 0
	InputBg    = lipgloss.AdaptiveColor{Dark: "#45475A", Light: "#CCD0DA"} // Surface 1
	Highlight  = lipgloss.AdaptiveColor{Dark: "#74C7EC", Light: "#209FB5"}

	Red        = lipgloss.AdaptiveColor{Dark: "#F38BA8", Light: "#D20F39"}
	Orange     = lipgloss.AdaptiveColor{Dark: "#FAB387", Light: "#FE640B"}
	Yellow     = lipgloss.AdaptiveColor{Dark: "#F9E2AF", Light: "#DF8E1D"}
	Green      = lipgloss.AdaptiveColor{Dark: "#A6E3A1", Light: "#40A02B"}
	Blue       = lipgloss.AdaptiveColor{Dark: "#89B4FA", Light: "#1E66F5"}
	Indigo     = lipgloss.AdaptiveColor{Dark: "#7287FD", Light: "#8839EF"} // swapped indigo/purple light values
	Purple     = lipgloss.AdaptiveColor{Dark: "#CBA6F7", Light: "#7287FD"}
)
