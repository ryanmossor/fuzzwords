package game

import (
	"fzwds/src/utils"
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
		LettersRemaining: utils.StringToMap(alphabet),
		Stats: InitializePlayerStats(),
	}

	return player
}

type TurnResult struct {
	ExtraLifeGranted	bool
}

func (g *GameState) HandleCorrectAnswer(answer string) TurnResult {
	result := TurnResult { ExtraLifeGranted: false }
	g.Player.TurnsSinceLastExtraLife++

	for _, c := range strings.ToUpper(answer) {
		ch := string(c)

		if strings.Contains(g.Alphabet, ch) && !slices.Contains(g.Player.LettersUsed, ch) {
			g.Player.LettersUsed = append(g.Player.LettersUsed, ch)
		}

		g.Player.LettersRemaining[ch] = true
	}

	if len(g.Player.LettersUsed) >= len(g.Alphabet) {
		g.GrantExtraLife()
		result.ExtraLifeGranted = true
	}

	slices.Sort(g.Player.LettersUsed)
	g.Player.Stats.UpdateSolvedStats(answer)

	return result
}

func (g *GameState) GrantExtraLife() {
	g.Player.LettersUsed = nil
	g.Player.LettersRemaining = utils.StringToMap(g.Alphabet)

	g.Player.Stats.ExtraLivesGained++
	if g.Player.Stats.FewestExtraLifeSolves == 0 || g.Player.TurnsSinceLastExtraLife < g.Player.Stats.FewestExtraLifeSolves {
		g.Player.Stats.FewestExtraLifeSolves = g.Player.TurnsSinceLastExtraLife
	}
	g.Player.TurnsSinceLastExtraLife = 0

	if g.Player.HealthCurrent < g.Settings.HealthMax {
		g.Player.HealthCurrent++
	}
}

func (g *GameState) HandleFailedTurn() {
	g.CurrentTurn.Strikes++
	g.Player.HealthCurrent--
	g.Player.TurnsSinceLastExtraLife++
	g.Player.Stats.UpdateFailedStats()
}
