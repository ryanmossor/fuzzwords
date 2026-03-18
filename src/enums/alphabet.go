package enums

import (
	"encoding/json"
	"log/slog"
	"strings"
)

type Alphabet int 
const (
	AlphabetEasy Alphabet = iota
	AlphabetMedium
	AlphabetFull
)

var (
	AlphabetName = map[int]string{
		0: "easy",
		1: "medium",
		2: "full",
	}

	AlphabetValue = map[string]int{
		"easy": 0,
		"medium": 1,
		"full": 2,
	}

	Alphabets = map[Alphabet]string{
		AlphabetEasy: "ABCDEFGHILMNOPRSTUWY", // J, K, Q, V, X, Z removed
		AlphabetMedium: "ABCDEFGHIJKLMNOPQRSTUVWY", // X and Z removed
		AlphabetFull: "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
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
