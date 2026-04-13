package game

import (
	"fzwds/src/utils"
)

type Player struct {
	HealthCurrent 			int
	LettersUsed				[]rune
	LettersRemaining 		map[rune]bool
	Streak					int
	Stats					PlayerStats
}

func InitializePlayer(cfg *GameSettings, alphabet string) Player {
	player := Player{
		HealthCurrent: cfg.HealthInitial,
		LettersRemaining: utils.StringToCharMap(alphabet),
	}

	return player
}
