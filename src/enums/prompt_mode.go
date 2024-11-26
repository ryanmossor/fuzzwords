package enums

import (
	"encoding/json"
	"log/slog"
	"strings"
)

type PromptMode int
const (
	Fuzzy PromptMode = iota
	Classic
)

var (
	PromptModeName = map[int]string{
		0: "fuzzy",
		1: "classic",
	}

	PromptModeValue = map[string]int{
		"fuzzy": 0,
		"classic": 1,
	}
)

func (m PromptMode) String() string {
	return PromptModeName[int(m)]
}

func ParsePromptMode(s string) PromptMode {
	s = strings.TrimSpace(strings.ToLower(s))
	value, ok := PromptModeValue[s]
	if !ok {
		slog.Warn("Invalid prompt mode - defaulting to fuzzy", "promptMode", s)
		return PromptMode(0)
	}
	return PromptMode(value)
}

func (m PromptMode) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.String())
}

func (m *PromptMode) UnmarshalJSON(data []byte) (err error) {
	var mode string
	if err := json.Unmarshal(data, &mode); err != nil {
		slog.Error("Error parsing prompt mode", "errMsg", err)
		return err
	}
	*m = ParsePromptMode(mode)

	return nil
}

