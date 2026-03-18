package game

import (
	fzwds "fzwds/src"
	"fzwds/src/enums"
	"fzwds/src/utils"
	"time"
)

type GameState struct {
	Alphabet			string
	GameActive			bool
	GameStart			time.Time
	GameEnd				time.Time
	Settings			Settings
	WordLists			WordLists
	Player				Player
	PreviousTurn		Turn
	CurrentTurn			Turn
	// TODO: cache next turn?
}

func InitializeGame(settings *Settings) GameState {
    word_lists := WordLists {
		FULL_MAP: fzwds.EnglishDictionaryMap,
		Available: utils.FilterWordList(fzwds.EnglishDictionary, settings.PromptLenMin),
        Used: make(map[string]bool),
    }
	alphabet := enums.Alphabets[settings.Alphabet]

	return GameState {
		Alphabet: alphabet,
		Settings: *settings,
		WordLists: word_lists,
		Player: InitializePlayer(settings, alphabet),
	}
}

func (g *GameState) StartGame() {
	g.GameStart = time.Now()
	g.GameActive = true
}

func (g *GameState) EndGame(won bool) {
	g.GameEnd = time.Now()
	g.GameActive = false
	if !won {
		g.Player.HealthCurrent = 0
	}
}

func (g *GameState) IsGameOver() bool {
	player_dead := g.Player.HealthCurrent == 0
	all_words_used := len(g.WordLists.Available) == 0
	max_lives_win := g.Settings.WinCondition == enums.MaxLives && g.Player.HealthCurrent == g.Settings.HealthMax

	if player_dead || all_words_used || max_lives_win {
		return true
	}

	return false
}
