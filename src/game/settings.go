package game

import "fzw/src/enums"

// TODO: individual setting struct w/ name, default value, optional help text?
type Settings struct {
	Alphabet				string
	HealthInitial			int
	HealthMax				int
	HighlightInput			bool
	PromptLenMin			int
	PromptLenMax			int
	PromptMode				enums.PromptMode
	PromptStrikesMax		int
	TurnDurationMin			int
	WinCondition			enums.WinCondition
	// TODO: add cfg for hints after each strike?
	// hints_enabled			bool
	// hint_chars_per_turn		int
}

func InitializeSettings() Settings {
	return Settings{
		Alphabet: enums.FullAlphabet,
		HealthInitial: 2,
		HealthMax: 3,
		HighlightInput: true,
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
