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

func InitializePlayer(cfg *Settings, alphabet string) Player {
	player := Player{
		HealthCurrent: cfg.HealthInitial,
		LettersRemaining: alphabetToMap(alphabet),
		Stats: InitializePlayerStats(),
	}

	return player
}

func (g *GameState) HandleCorrectAnswer(answer string) {
	g.Player.TurnsSinceLastExtraLife++

	for _, c := range strings.ToUpper(answer) {
		ch := string(c)

		if strings.Contains(g.Alphabet, ch) && !slices.Contains(g.Player.LettersUsed, ch) {
			g.Player.LettersUsed = append(g.Player.LettersUsed, ch)
		}

		g.Player.LettersRemaining[ch] = true
	}

	slices.Sort(g.Player.LettersUsed)
	g.Player.Stats.UpdateSolvedStats(answer)
}

func (g *GameState) ShouldGrantExtraLife() bool {
	if len(g.Player.LettersUsed) < len(g.Alphabet) {
		return false
	}

	g.Player.LettersUsed = nil
	g.Player.LettersRemaining = alphabetToMap(g.Alphabet)

	g.Player.Stats.ExtraLivesGained++
	if g.Player.Stats.FewestExtraLifeSolves == 0 || g.Player.TurnsSinceLastExtraLife < g.Player.Stats.FewestExtraLifeSolves {
		g.Player.Stats.FewestExtraLifeSolves = g.Player.TurnsSinceLastExtraLife
	}
	g.Player.TurnsSinceLastExtraLife = 0

	if g.Player.HealthCurrent < g.Settings.HealthMax {
		g.Player.HealthCurrent++
	}

	return true
}

func (g *GameState) HandleFailedTurn() {
	g.CurrentTurn.Strikes++
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
