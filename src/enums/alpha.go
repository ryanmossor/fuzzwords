package enums

import (
	"encoding/json"
	"log/slog"
	"strings"
)

type Alphabet int 
const (
	DebugAlphabet Alphabet = iota
	EasyAlphabet
	MediumAlphabet
	FullAlphabet
)

var (
	AlphabetName = map[int]string{
		0: "debug",
		1: "easy",
		2: "medium",
		3: "full",
	}

	AlphabetValue = map[string]int{
		"debug": 0,
		"easy": 1,
		"medium": 2,
		"full": 3,
	}

	Alphabets = map[Alphabet]string{
		DebugAlphabet: "ABC",
		EasyAlphabet: "ABCDEFGHILMNOPRSTUWY", // J, K, Q, V, X, Z removed
		MediumAlphabet: "ABCDEFGHIJKLMNOPQRSTUVWY", // X and Z removed
		FullAlphabet: "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
	}
)

func (a Alphabet) String() string {
	return AlphabetName[int(a)]
}

func ParseAlphabet(s string) Alphabet {
	s = strings.TrimSpace(strings.ToLower(s))
	value, ok := AlphabetValue[s]
	if !ok {
		slog.Warn("Invalid alphabet - defaulting to easy", "alphabet", s)
		return Alphabet(1)
	}
	return Alphabet(value)
}

func (a Alphabet) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.String())
}

func (a *Alphabet) UnmarshalJSON(data []byte) (err error) {
	var alphabet string
	if err := json.Unmarshal(data, &alphabet); err != nil {
		slog.Error("Error parsing alphabet", "errMsg", err)
		return err
	}
	*a = ParseAlphabet(alphabet)

	return nil
}
