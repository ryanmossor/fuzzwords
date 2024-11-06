package game

import (
	"fzw/src/utils"
	"slices"
	"strings"
)

type Player struct {
	HealthCurrent 		int
	HealthDisplay		string
	LettersUsed			[]string
	LettersRemaining 	[]string
	gameSettings		*Settings
}

func InitializePlayer(cfg *Settings) Player {
	player := Player{
		HealthCurrent: cfg.HealthInitial,
		LettersUsed: nil,
		LettersRemaining: strings.Split(cfg.Alphabet, ""),
		gameSettings: cfg,
	}
	player.UpdateHealthDisplay()

	return player
}

func (p *Player) IncrementHealth() {
	p.LettersUsed = nil
	p.LettersRemaining = strings.Split(p.gameSettings.Alphabet, "")

	if p.HealthCurrent < p.gameSettings.HealthMax {
		p.HealthCurrent++
		p.UpdateHealthDisplay()
	}
}

func (p *Player) DecrementHealth() {
	p.HealthCurrent--
	p.UpdateHealthDisplay()
}

func (p *Player) UpdateHealthDisplay() {
	health_display := ""

	i := 0
	for i < p.HealthCurrent {
		// 🧡💛💚💙🩵💜🖤🤍🤎
		health_display += "🩵"
		i++
	}
	for i < p.gameSettings.HealthMax {
		health_display += "🤍"
		i++
	}

	p.HealthDisplay = health_display
}

func (p *Player) HandleCorrectAnswer(answer string) {
	for i := 0; i < len(answer); i++ {
		c := strings.ToUpper(string(answer[i]))

		if strings.Contains(p.gameSettings.Alphabet, c) && !slices.Contains(p.LettersUsed, c) {
			p.LettersUsed = append(p.LettersUsed, c)
		}

		if slices.Contains(p.LettersRemaining, c) {
			p.LettersRemaining = utils.Remove(p.LettersRemaining, slices.Index(p.LettersRemaining, c))
		}
	}

	if len(p.LettersUsed) >= len(p.gameSettings.Alphabet) {
		p.IncrementHealth()
	}

	slices.Sort(p.LettersUsed)
}
