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

func (m model) DebugView() string {
	if !m.debug {
		return ""
	}

	// tabs := []string{
	// 	"VH " + strconv.Itoa(m.viewport_height),
	// 	"VW " + strconv.Itoa(m.viewport_width),
	// 	"CW " + strconv.Itoa(m.width_container),
	// 	m.size.String(),
	// }

	return table.New().
		Border(lipgloss.HiddenBorder()).
		BorderStyle(m.renderer.NewStyle().Foreground(m.theme.Border())).
		Row(memStatsView()...).
		// Row(tabs...).
		Width(m.width_container).
		StyleFunc(func(row, col int) lipgloss.Style {
			return m.theme.Base().
				Padding(0, 1).
				AlignHorizontal(lipgloss.Center)
		}).
		Render()
}

func memStatsView() []string {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	var stats []string

	// Print total memory allocated and still in use (in bytes)
	// stats = append(stats, fmt.Sprintf("Total Alloc %v MiB", memStats.TotalAlloc/1024/1024))
	stats = append(stats, fmt.Sprintf("Sys %v MiB", memStats.Sys/1024/1024))
	stats = append(stats, fmt.Sprintf("Heap Alloc %v MiB", memStats.HeapAlloc/1024/1024))
	stats = append(stats, fmt.Sprintf("Heap Sys %v MiB", memStats.HeapSys/1024/1024))

	return stats
}
