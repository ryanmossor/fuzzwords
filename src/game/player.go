package game

import (
	"fzwds/src/utils"
)

type Player struct {
	HealthCurrent   	int
	LettersUsed     	[]rune
	LettersRemaining	map[rune]bool
	Streak          	int
	Stats           	PlayerStats
}

func (g *GameState) InitializePlayer() Player {
	player := Player{
		HealthCurrent:    g.Settings.HealthInitial,
		LettersRemaining: utils.StringToCharMap(g.Alphabet),
		LettersUsed:      make([]rune, 0, len(g.Alphabet)),
		Streak:           0,
		Stats:            PlayerStats{},
	}
	return player
}
