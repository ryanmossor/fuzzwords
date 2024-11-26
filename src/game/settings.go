package game

import (
	"fzwds/src/enums"
	"slices"
	"strings"
)

type ValidValue struct {
	Name			string `json:"name"`
	Description		string `json:"description"`
}

type Config struct {
	Name			string 			`json:"name"`
	PropName		string 			`json:"propName"`
	Disabled		bool 			`json:"disabled"`
	Default			string 			`json:"default"`
	Description		string 			`json:"description,omitempty"`
	Min				int 			`json:"min,omitempty"`
	Max				int 			`json:"max,omitempty"`
	ValidValues		[]ValidValue 	`json:"validValues,omitempty"`
}

// TODO: individual setting struct w/ name, default value, optional help text?
type Settings struct {
	Alphabet				enums.Alphabet		`json:"alphabet"`
	HealthInitial			int					`json:"healthInital"`
	HealthMax				int					`json:"healthMax"`
	HighlightInput			bool				`json:"highlightInput"`
	PromptLenMin			int					`json:"promptLenMin"`
	PromptLenMax			int					`json:"promptLenMax"`
	PromptMode				enums.PromptMode	`json:"promptMode"`
	PromptStrikesMax		int					`json:"promptStrikesMax"`
	TurnDurationMin			int					`json:"turnDurationMin"`
	WinCondition			enums.WinCondition	`json:"winCondition"`
	// TODO: add cfg for hints after each strike?
	// hints_enabled			bool
	// hint_chars_per_turn		int
}

func InitializeSettings() Settings {
	return Settings{
		Alphabet: enums.FullAlphabet,
		HealthInitial: 2,
		HealthMax: 3,
		HighlightInput: false,
		PromptLenMax: 3,
		PromptLenMin: 2,
		PromptMode: enums.Fuzzy,
		PromptStrikesMax: 2,
		TurnDurationMin: 10,
		WinCondition: enums.Endless,
	}
}

func EasySettings() Settings {
	cfg := InitializeSettings()
	cfg.Alphabet = enums.EasyAlphabet
	cfg.TurnDurationMin = 15
	return cfg
}

func MediumSettings() Settings {
	cfg := InitializeSettings()
	cfg.Alphabet = enums.MediumAlphabet
	cfg.TurnDurationMin = 10
	return cfg
}

func DifficultSettings() Settings {
	cfg := InitializeSettings()
	cfg.Alphabet = enums.FullAlphabet
	cfg.TurnDurationMin = 5
	cfg.PromptLenMax = 4
	cfg.PromptStrikesMax = 2
	return cfg
}

func ExpertSettings() Settings {
	cfg := InitializeSettings()
	cfg.Alphabet = enums.FullAlphabet
	cfg.HealthInitial = 1
	cfg.HealthMax = 1
	cfg.TurnDurationMin = 5
	cfg.PromptLenMax = 5
	cfg.PromptStrikesMax = 1
	return cfg
}

func (s *Settings) SetAlphabet(alphabet_idx int) *Settings {
	if alphabet_idx < 1 {
		alphabet_idx = len(enums.AlphabetValue) - 1
	} else if alphabet_idx > len(enums.AlphabetValue) - 1 {
		alphabet_idx = 1
	}
	s.Alphabet = enums.Alphabet(alphabet_idx)

	return s
}

var HEALTH_MIN = 1
var HEALTH_MAX = 10

func (s *Settings) SetHealthInitial(health int) *Settings {
	// health := s.HealthInitial + count

	if health < HEALTH_MIN {
		s.HealthInitial = HEALTH_MIN
	} else if health > HEALTH_MAX {
		s.HealthInitial = HEALTH_MAX
	} else {
		s.HealthInitial = health
	}

	if s.HealthInitial > s.HealthMax {
		s.HealthMax = s.HealthInitial
	}

	return s
}

func (s *Settings) SetHealthMax(health int) *Settings {
	// health := s.HealthMax + count

	if health < HEALTH_MIN {
		s.HealthMax = HEALTH_MIN
	} else if health > HEALTH_MAX {
		s.HealthMax = HEALTH_MAX
	} else {
		s.HealthMax = health
	}

	if s.HealthMax < s.HealthInitial {
		s.HealthInitial = s.HealthMax
	}

	return s
}

var PROMPT_LEN_MIN = 2
var PROMPT_LEN_MAX = 5

func (s *Settings) SetPromptLenMin(len int) *Settings {
	// len := s.PromptLenMin + count

	if len < PROMPT_LEN_MIN {
		s.PromptLenMin = PROMPT_LEN_MIN
	} else if len > PROMPT_LEN_MAX {
		s.PromptLenMin = PROMPT_LEN_MAX 
	} else {
		s.PromptLenMin = len
	}

	if s.PromptLenMin > s.PromptLenMax {
		s.PromptLenMax = s.PromptLenMin
	}

	return s
}

func (s *Settings) SetPromptLenMax(len int) *Settings {
	// len := s.PromptLenMax + count

	if len < PROMPT_LEN_MIN {
		s.PromptLenMax = PROMPT_LEN_MIN
	} else if len > PROMPT_LEN_MAX {
		s.PromptLenMax = PROMPT_LEN_MAX 
	} else {
		s.PromptLenMax = len
	}

	if s.PromptLenMax < s.PromptLenMin {
		s.PromptLenMin = s.PromptLenMax
	}

	return s
}

func (s *Settings) SetPromptMode(mode string) *Settings {
	switch strings.ToLower(mode) {
	case "classic":
		s.PromptMode = enums.Classic
	default:
		s.PromptMode = enums.Fuzzy
	}

	return s
}

func (s *Settings) SetWinCondition(cond string) *Settings {
	switch strings.ToLower(cond) {
	case "debug":
		s.WinCondition = enums.Debug
	case "max lives":
		s.WinCondition = enums.MaxLives
	default:
		s.WinCondition = enums.Endless
	}

	return s
}

var STRIKES_MIN = 1
var STRIKES_MAX = 3

func (s *Settings) SetPromptStrikesMax(strikes int) *Settings {
	if strikes < STRIKES_MIN {
		s.PromptStrikesMax = STRIKES_MIN
	} else if strikes > STRIKES_MAX {
		s.PromptStrikesMax = STRIKES_MAX
	} else {
		s.PromptStrikesMax = strikes
	}
	
	return s
}

var TURN_DURATION_MIN = 5
var TURN_DURATION_MAX = 60
var TURN_DURATION_INTERVALS = []int{ 5, 10, 15, 20, 25, 30, 45, 60 }

func (s *Settings) SetTurnDurationMin(duration int) *Settings {
	if !slices.Contains(TURN_DURATION_INTERVALS, duration) {
		s.TurnDurationMin = 5
	} else {
		s.TurnDurationMin = duration
	}

	return s
}

func (s *Settings) ValidateSettings() *Settings {
	if s.Alphabet != enums.EasyAlphabet && s.Alphabet != enums.MediumAlphabet && s.Alphabet != enums.FullAlphabet && s.Alphabet != enums.DebugAlphabet {
		s.Alphabet = enums.EasyAlphabet
	}

	return s.
		SetHealthInitial(s.HealthInitial).
		SetHealthMax(s.HealthMax).
		SetPromptLenMin(s.PromptLenMin).
		SetPromptLenMax(s.PromptLenMax).
		SetPromptMode(s.PromptMode.String()).
		SetPromptStrikesMax(s.PromptStrikesMax).
		SetTurnDurationMin(s.TurnDurationMin).
		SetWinCondition(s.WinCondition.String())
}
