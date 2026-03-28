package game

import (
	"fzwds/src/dictionary"
	"fzwds/src/enums"
	"fzwds/src/utils"
	"log/slog"
	"time"
)

type GameState struct {
	Alphabet			string
	GameActive			bool
	GameStart			time.Time
	GameEnd				time.Time
	Settings			GameSettings
	WordLists			WordLists
	Player				Player
	PreviousTurn		Turn
	CurrentTurn			Turn
	StartUnixTs			int64
	// TODO: cache next turn?
}

func InitializeGame(settings *GameSettings) GameState {
	var full_map map[string]bool
	var available []string

	switch settings.Dictionary {
	case enums.English:
		full_map = dictionary.EnglishDictionaryMap
		available = utils.FilterWordList(dictionary.EnglishDictionary, settings.PromptLenMin)
	case enums.Pokemon:
		available = dictionary.GetSelectedPokemonGenList(settings.PokemonGens...)
		full_map = utils.ArrToMap(available)
	}

    word_lists := WordLists {
		FULL_MAP: full_map,
		Available: available,
        Used: make(map[string]bool),
    }
	alphabet := enums.Alphabets[settings.Alphabet]

	g := GameState {
		Alphabet: alphabet,
		GameActive: true,
		GameStart: time.Now(),
		Settings: *settings,
		WordLists: word_lists,
		Player: InitializePlayer(settings, alphabet),
		StartUnixTs: time.Now().UnixMilli(),
	}
	g.NewTurn(true)

	slog.Info("Initialized game",
		"startUnixTs", g.StartUnixTs,
		"alphabet", g.Alphabet,
		"settings", g.Settings)

	return g
}

func (g *GameState) EndGame(won bool) {
	g.GameEnd = time.Now()
	g.GameActive = false
	g.Player.Stats.TimeSurvived = int(g.GameEnd.Sub(g.GameStart).Seconds())
	if !won {
		g.Player.HealthCurrent = 0
	}
}

func (g *GameState) IsGameOver() bool {
	player_dead := g.Player.HealthCurrent == 0
	all_words_used := len(g.WordLists.Available) == 0
	max_lives_win := g.Settings.WinCondition == enums.WinConditionMaxLives &&
					 g.Player.HealthCurrent == g.Settings.HealthMax

	if player_dead || all_words_used || max_lives_win {
		return true
	}

	return false
}
