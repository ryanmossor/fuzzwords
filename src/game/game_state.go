package game

import (
	fzwds "fzwds/src"
	"fzwds/src/enums"
	"fzwds/src/utils"
)

type GameState struct {
	Alphabet			string
	Settings			Settings
	WordLists			WordLists
	Player				Player
	PreviousTurn		Turn
	CurrentTurn			Turn
	// TODO: cache next turn?
}

func InitializeGame(settings *Settings) GameState {
	word_list := fzwds.EnglishDictionary
    word_lists := WordLists{
        FULL_MAP: utils.ArrToMap(word_list),
        Available: utils.FilterWordList(word_list, settings.PromptLenMin),
        Used: make(map[string]bool),
    }
	alphabet := enums.Alphabets[settings.Alphabet]

	return GameState{
		Alphabet: alphabet,
		Settings: *settings,
		WordLists: word_lists,
		Player: InitializePlayer(settings, alphabet),
	}
}
