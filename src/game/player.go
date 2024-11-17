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
}

func InitializePlayer(cfg *Settings) Player {
	player := Player{
		HealthCurrent: cfg.HealthInitial,
		LettersRemaining: alphabetToMap(cfg.Alphabet),
		Stats: InitializePlayerStats(),
	}

	return player
}

func (g *GameState) HandleCorrectAnswer() {
	g.Player.TurnsSinceLastExtraLife++

	for _, c := range strings.ToUpper(g.CurrentTurn.Answer) {
		ch := string(c)

		if strings.Contains(g.Settings.Alphabet, ch) && !slices.Contains(g.Player.LettersUsed, ch) {
			g.Player.LettersUsed = append(g.Player.LettersUsed, ch)
		}

		g.Player.LettersRemaining[ch] = true
	}

	if len(g.Player.LettersUsed) >= len(g.Settings.Alphabet) {
		g.Player.LettersUsed = nil
		g.Player.LettersRemaining = alphabetToMap(g.Settings.Alphabet)

		g.Player.Stats.ExtraLivesGained++
		if g.Player.Stats.FewestExtraLifeSolves == 0 || g.Player.TurnsSinceLastExtraLife < g.Player.Stats.FewestExtraLifeSolves {
			g.Player.Stats.FewestExtraLifeSolves = g.Player.TurnsSinceLastExtraLife
		}
		g.Player.TurnsSinceLastExtraLife = 0

		if g.Player.HealthCurrent < g.Settings.HealthMax {
			g.Player.HealthCurrent++
		}
	}

	slices.Sort(g.Player.LettersUsed)
	g.Player.Stats.UpdateSolvedStats(g.CurrentTurn.Answer)
}

func (g *GameState) HandleFailedTurn() {
	g.Player.HealthCurrent--
	g.Player.TurnsSinceLastExtraLife++
	g.Player.Stats.UpdateFailedStats()
}

func alphabetToMap(alphabet string) map[string]bool {
	letters_remaining := make(map[string]bool)
	for _, c := range alphabet {
		letters_remaining[string(c)] = false
	}
	return letters_remaining
}
