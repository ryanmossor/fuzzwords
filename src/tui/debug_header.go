package tui

import (
	"fmt"
	"runtime"
	// "strconv"

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
		// "VH " + strconv.Itoa(m.viewport_height),
		// "VW " + strconv.Itoa(m.viewport_width),
		// "CW " + strconv.Itoa(m.width_container),
		// "coloredStrikeLen " + m.debug_map["coloredStrikeLen"],
		// "visibleLen " + m.debug_map["visibleLen"],
		// "strikeLen " + m.debug_map["strikeLen"],

		"viewSize " + m.debug_map["viewSize"] + " B",
		// "runeCount " + m.debug_map["runeCount"],

		// fmt.Sprintf("heightContainer %d", m.height_container),
		// fmt.Sprintf("heightContent %d", m.height_content),
		// fmt.Sprintf("sel %d", m.state.game_review.selected_turn),
		// fmt.Sprintf("visible %s", m.debug_map["visibleRows"]),
		// fmt.Sprintf("visStart %s", m.debug_map["visStart"]),

		m.debug_map["tableTime"],
		// fmt.Sprintf("reviewCacheLen %d", len(m.state.game_review.view_cache)),
		// fmt.Sprintf("rowWidStr %s", m.debug_map["rowWidStr"]),
		fmt.Sprintf("Turn %d", m.state.game.CurrentTurnNumber()),

		// "selected " + m.debug_map["selected"],
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
	// stats = append(stats, fmt.Sprintf("Current alloc: %v MiB", memStats.HeapAlloc / 1024 / 1024))
	// Total heap space reserved (used and unused)
	stats = append(stats, fmt.Sprintf("Heap reserved: %v MiB", memStats.HeapSys / 1024 / 1024))
	// Cumulative memory requested by program
	stats = append(stats, fmt.Sprintf("Total mem: %v MiB", memStats.Sys / 1024 / 1024))

	stats = append(stats, m.debug_map["fps"])

	return stats
}
