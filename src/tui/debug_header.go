package tui

import (
	"fmt"
	"runtime"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

func (p page) String() string {
	switch p {
	case splash_page:
		return "Splash"
	case about_page:
		return "About"
	case settings_page:
		return "Settings"
	case game_page:
		return "Game"
	case game_over_page:
		return "Game Over"
	default:
		return "Unknown page"
	}
}

func (s size) String() string {
	switch s {
	case undersized:
		return "Undersized"
	case small:
		return "Small"
	case medium:
		return "Medium"
	case large:
		return "Large"
	default:
		return "Unknown size"
	}
}

// TODO: make debug view a panel that appears left/right of main view rather than a finnicky header
func (m model) DebugView() string {
	if !m.debug {
		return ""
	}

	tabs := []string{
		"viewSize " + m.debug_map["viewSize"] + " B",
		fmt.Sprintf("input %v", !m.state.game.inputRestricted),
		"keyPress " + m.debug_map["keyPress"],

		// fmt.Sprintf("heightContainer %d", m.height_container),
		// fmt.Sprintf("heightContent %d", m.height_content),

		fmt.Sprintf("Turn %d", m.game.CurrentTurnNumber()),

		// string(m.page),
		// m.size.String(),
	}

	return table.New().
		Border(lipgloss.HiddenBorder()).
		BorderBottom(false).
		Row(tabs...).
		Row(m.memStatsView()...).
		Width(m.width_container).
		Render()
}

func (m model) memStatsView() []string {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	var stats []string

	// Total memory allocated and in use
	stats = append(stats, fmt.Sprintf("curAlloc: %v MiB", memStats.HeapAlloc / 1024 / 1024))
	// Total heap space reserved (used and unused)
	stats = append(stats, fmt.Sprintf("heapResv: %v MiB", memStats.HeapSys / 1024 / 1024))
	// Cumulative memory requested by program
	stats = append(stats, fmt.Sprintf("memTotal: %v MiB", memStats.Sys / 1024 / 1024))

	stats = append(stats, m.debug_map["renderTime"])

	return stats
}
