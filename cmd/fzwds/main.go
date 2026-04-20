package main

import (
	"flag"
	"fmt"
	"fzwds/pkg/tui"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
        // panic("Cache dir not found" + err.Error())
		cacheDir = os.TempDir()
	}
	path := filepath.Join(cacheDir, "fuzzwords")
	os.MkdirAll(path, os.ModePerm)

	logFile, err := os.OpenFile(filepath.Join(path, "log.json"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        panic("Failed to open log file: " + err.Error())
    }
	defer logFile.Close()

	debug := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	var logLevel slog.Leveler
	if *debug {
		logLevel = slog.LevelDebug
	} else {
		logLevel = slog.LevelInfo
	}

	opts := &slog.HandlerOptions {
		Level: logLevel,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key != slog.TimeKey {
				return a
			}
			t := a.Value.Time()
			a.Value = slog.StringValue(t.Format(time.DateTime)) // format as YYYY-MM-DD HH:mm:ss
			return a
		},
	}

    fileHandler := slog.NewJSONHandler(logFile, opts)
    slog.SetDefault(slog.New(fileHandler))

	schema := tui.LoadSchema()
	settings, path := tui.LoadSettings(schema)

	model := tui.NewModel(*debug, settings, schema, path)
	app := tea.NewProgram(model,
        tea.WithAltScreen(),
        tea.WithMouseCellMotion(),
    )
	_, err = app.Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		slog.Error("Error running program", "errMsg", err)
		os.Exit(1)
	}
}
