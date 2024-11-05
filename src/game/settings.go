package game

import "fzw/src/enums"

type Settings struct {
	Alphabet				string
	HealthInitial			int
	HealthMax				int
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
		Alphabet: enums.DebugAlphabet,
		HealthInitial: 2,
		HealthMax: 3,
		PromptLenMax: 3,
		PromptLenMin: 2,
		PromptMode: enums.Fuzzy,
		PromptStrikesMax: 3,
		TurnDurationMin: 10,
		WinCondition: enums.Endless,
	}
}
