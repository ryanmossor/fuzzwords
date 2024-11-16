package game

import (
	"slices"
	"strings"
)

type Player struct {
	HealthCurrent 			int
	HealthDisplay			string
	LettersUsed				[]string
	LettersRemaining 		map[string]bool
	TurnsSinceLastExtraLife int
	Stats					PlayerStats
	_gameSettings			*Settings
}

func InitializePlayer(cfg *Settings) Player {
	player := Player{
		HealthCurrent: cfg.HealthInitial,
		LettersRemaining: alphabetToMap(cfg.Alphabet),
		Stats: InitializePlayerStats(),
		_gameSettings: cfg,
	}

	return player
}

func (p *Player) HandleCorrectAnswer(answer string) {
	p.TurnsSinceLastExtraLife++

	for _, c := range strings.ToUpper(answer) {
		ch := string(c)

		if strings.Contains(p._gameSettings.Alphabet, ch) && !slices.Contains(p.LettersUsed, ch) {
			p.LettersUsed = append(p.LettersUsed, ch)
		}

		p.LettersRemaining[ch] = true
	}

	if len(p.LettersUsed) >= len(p._gameSettings.Alphabet) {
		p.LettersUsed = nil
		p.LettersRemaining = alphabetToMap(p._gameSettings.Alphabet)

		p.Stats.ExtraLivesGained++
		if p.Stats.FewestExtraLifeSolves == 0 || p.TurnsSinceLastExtraLife < p.Stats.FewestExtraLifeSolves {
			p.Stats.FewestExtraLifeSolves = p.TurnsSinceLastExtraLife
		}
		p.TurnsSinceLastExtraLife = 0

		if p.HealthCurrent < p._gameSettings.HealthMax {
			p.HealthCurrent++
		}
	}

	slices.Sort(p.LettersUsed)
	p.Stats.UpdateSolvedStats(answer)
}

func (p *Player) HandleFailedTurn() {
	p.HealthCurrent--
	p.TurnsSinceLastExtraLife++
	p.Stats.UpdateFailedStats()
}

func alphabetToMap(alphabet string) map[string]bool {
	letters_remaining := make(map[string]bool)
	for _, c := range alphabet {
		letters_remaining[string(c)] = false
	}
	return letters_remaining
}
