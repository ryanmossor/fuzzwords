package enums

import (
	"encoding/json"
	"log/slog"
	"strings"
)

type PromptMode uint8
const (
	Fuzzy PromptMode = iota
	Classic
)

var (
	PromptModeName = map[uint8]string{
		0: "fuzzy",
		1: "classic",
	}

	PromptModeValue = map[string]uint8{
		"fuzzy": 0,
		"classic": 1,
	}
)

func (m PromptMode) String() string {
	return PromptModeName[uint8(m)]
}

func ParsePromptMode(s string) PromptMode {
	s = strings.TrimSpace(strings.ToLower(s))
	value, ok := PromptModeValue[s]
	if !ok {
		slog.Error("Invalid prompt mode - defaulting to fuzzy", "promptMode", s)
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

