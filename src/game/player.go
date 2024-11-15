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
	player.UpdateHealthDisplay()

	return player
}

// TODO: move to somewhere in tui? may need to change from hearts to colorized vertical bars for terminal compatibility
func (p *Player) UpdateHealthDisplay() {
	health_display := ""

	i := 0
	for i < p.HealthCurrent {
		// 🧡💛💚💙🩵💜🖤🤍🤎
		health_display += "🩵"
		i++
	}
	for i < p._gameSettings.HealthMax {
		health_display += "🤍"
		i++
	}

	p.HealthDisplay = health_display
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
			p.UpdateHealthDisplay()
		}
	}

	slices.Sort(p.LettersUsed)
	p.Stats.UpdateSolvedStats(answer)
}

func (p *Player) HandleFailedTurn() {
	p.HealthCurrent--
	p.UpdateHealthDisplay()

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
