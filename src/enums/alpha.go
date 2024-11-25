package enums

import (
	"encoding/json"
	"log/slog"
	"strings"
)

type Alphabet uint8 
const (
	EasyAlphabet Alphabet = iota
	MediumAlphabet
	FullAlphabet
	DebugAlphabet
)

var (
	AlphabetName = map[uint8]string{
		0: "easy",
		1: "medium",
		2: "full",
		3: "debug",
	}

	AlphabetValue = map[string]uint8{
		"easy": 0,
		"medium": 1,
		"full": 2,
		"debug": 3,
	}

	Alphabets = map[Alphabet]string{
		EasyAlphabet: "ABCDEFGHILMNOPRSTUWY", // J, K, Q, V, X, Z removed
		MediumAlphabet: "ABCDEFGHIJKLMNOPQRSTUVWY", // X and Z removed
		FullAlphabet: "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		DebugAlphabet: "ABC",
	}
)

func (a Alphabet) String() string {
	return AlphabetName[uint8(a)]
}

func ParseAlphabet(s string) Alphabet {
	s = strings.TrimSpace(strings.ToLower(s))
	value, ok := AlphabetValue[s]
	if !ok {
		slog.Error("Invalid alphabet - defaulting to easy", "alphabet", s)
		return Alphabet(0)
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
