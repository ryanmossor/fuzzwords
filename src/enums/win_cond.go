package enums

import (
	"encoding/json"
	"log/slog"
	"strings"
)

type WinCondition int
const (
	Endless WinCondition = iota
	MaxLives
	Debug
)

var (
	WinCondName = map[int]string{
		0: "endless",
		1: "max lives",
		2: "debug",
	}

	WinCondValue = map[string]int{
		"endless": 0,
		"max lives": 1,
		"debug": 2,
	}
)

func (w WinCondition) String() string {
	return WinCondName[int(w)]
}

func ParseWinCond(s string) WinCondition {
	s = strings.TrimSpace(strings.ToLower(s))
	value, ok := WinCondValue[s]
	if !ok {
		slog.Warn("Invalid win condition - defaulting to endless", "winCond", s)
		return WinCondition(0)
	}
	return WinCondition(value)
}

func (w WinCondition) MarshalJSON() ([]byte, error) {
	return json.Marshal(w.String())
}

func (w *WinCondition) UnmarshalJSON(data []byte) (err error) {
	var cond string
	if err := json.Unmarshal(data, &cond); err != nil {
		slog.Error("Error parsing win condition", "errMsg", err)
		return err
	}
	*w = ParseWinCond(cond)

	return nil
}
