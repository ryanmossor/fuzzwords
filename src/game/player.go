package game

import (
	"fzwds/src/utils"
)

type Player struct {
	HealthCurrent   	int
	LettersUsed     	[]rune
	LettersRemaining	map[rune]bool
	streak          	int
	Stats           	PlayerStats
}

func (g *GameState) newPlayer() Player {
	player := Player{
		HealthCurrent:    g.Settings.HealthInitial,
		LettersRemaining: utils.StringToCharMap(g.Alphabet),
		LettersUsed:      make([]rune, 0, len(g.Alphabet)),
		streak:           0,
		Stats:            PlayerStats{},
	}
	return player
}
