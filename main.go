package main

import (
	"fmt"
	"fzwds/src/tui"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	tea "github.com/charmbracelet/bubbletea"
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

	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
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

	menu := tui.NewModel()
	prog := tea.NewProgram(menu, tea.WithAltScreen())
	_, err = prog.Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		slog.Error("Error running program", "errMsg", err)
		os.Exit(1)
	}
}
 