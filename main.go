package main

import (
	"flag"
	"fmt"
	"fzwds/src/tui"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func main() {
	cache_dir, err := os.UserCacheDir()
	if err != nil {
        // panic("Cache dir not found" + err.Error())
		cache_dir = os.TempDir()
	}
	path := filepath.Join(cache_dir, "fuzzwords")
	os.MkdirAll(path, os.ModePerm)

	log_file, err := os.OpenFile(filepath.Join(path, "log.json"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        panic("Failed to open log file: " + err.Error())
    }
	defer log_file.Close()

	debug := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	var log_level slog.Leveler
	if *debug {
		log_level = slog.LevelDebug
	} else {
		log_level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: log_level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key != slog.TimeKey {
				return a
			}
			t := a.Value.Time()
			a.Value = slog.StringValue(t.Format(time.DateTime)) // format as YYYY-MM-DD HH:mm:ss
			return a
		},
	}

    fileHandler := slog.NewJSONHandler(log_file, opts)
    slog.SetDefault(slog.New(fileHandler))

	menu := tui.NewModel(lipgloss.DefaultRenderer(), *debug)
	prog := tea.NewProgram(
        menu,
        tea.WithAltScreen(),
        tea.WithMouseCellMotion(), // enable mouse support for scroll wheel usage
    )
	_, err = prog.Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		slog.Error("Error running program", "errMsg", err)
		os.Exit(1)
	}
}
 
