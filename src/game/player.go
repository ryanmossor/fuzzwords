package game

import (
	"fzw/src/utils"
	"slices"
	"strings"
)

type Player struct {
	HealthCurrent 		int
	HealthMax			int
	HealthDisplay		string
	LettersUsed			[]string
	LettersRemaining 	[]string
}

func InitializePlayer(cfg Settings) Player {
	player := Player{
		HealthCurrent: cfg.HealthInitial,
		HealthMax: cfg.HealthMax,
		LettersUsed: nil,
		LettersRemaining: strings.Split(cfg.Alphabet, ""),
	}
	player.UpdateHealthDisplay()
	return player
}

func (p *Player) IncrementHealth(cfg Settings) {
	p.LettersUsed = nil
	p.LettersRemaining = strings.Split(cfg.Alphabet, "")

	if p.HealthCurrent < p.HealthMax {
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
	for i < p.HealthMax {
		health_display += "🤍"
		i++
	}

	p.HealthDisplay = health_display
}

func (p *Player) HandleCorrectAnswer(answer string, player *Player, cfg Settings) {
	for i := 0; i < len(answer); i++ {
		c := strings.ToUpper(string(answer[i]))

		if strings.Contains(cfg.Alphabet, c) && !slices.Contains(player.LettersUsed, c) {
			player.LettersUsed = append(player.LettersUsed, c)
		}

		if slices.Contains(player.LettersRemaining, c) {
			player.LettersRemaining = utils.Remove(player.LettersRemaining, slices.Index(player.LettersRemaining, c))
		}
	}

	if len(player.LettersUsed) >= len(cfg.Alphabet) {
		player.IncrementHealth(cfg)
	}

	slices.Sort(player.LettersUsed)
}
