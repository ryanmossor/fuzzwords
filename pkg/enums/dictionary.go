package enums

import (
	"encoding/json"
	"log/slog"
	"strings"
)

type Dictionary int
const (
	English Dictionary = iota
	Pokemon
)

var (
	DictionaryName = map[int]string{
		0: "english",
		1: "pokemon",
	}

	DictionaryValue = map[string]int{
		"english": 0,
		"pokemon": 1,
	}
)

func (d Dictionary) String() string {
	return DictionaryName[int(d)]
}

func ParseDictionary(s string) Dictionary {
	s = strings.TrimSpace(strings.ToLower(s))
	value, ok := DictionaryValue[s]
	if !ok {
		slog.Warn("Invalid dictionary - defaulting to english", "dictionary", s)
		return Dictionary(1)
	}
	return Dictionary(value)
}

func (d Dictionary) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *Dictionary) UnmarshalJSON(data []byte) (err error) {
	var dictionary string
	if err := json.Unmarshal(data, &dictionary); err != nil {
		slog.Error("Error parsing dictionary", "errMsg", err)
		return err
	}
	*d = ParseDictionary(dictionary)

	return nil
}
