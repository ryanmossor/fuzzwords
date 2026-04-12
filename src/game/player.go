package game

import (
	"fzwds/src/utils"
	"slices"
	"strings"
	"time"
)

type Player struct {
	HealthCurrent 			int
	HealthDisplay			string
	LettersUsed				[]string
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

type TurnResult struct {
	ExtraLifeGranted	bool
}

func (g *GameState) HandleCorrectAnswer(answer string) {
	turn := g.CurrentTurn()
	turn.TotalTurnDuration = time.Since(turn.TurnStart)
	turn.Solved = true
	turn.Answer = answer
	turn.UniqueLetterCount = utils.CountUniqueLetters(answer)

	g.Player.Streak++
	turn.Streak = g.Player.Streak

	for _, c := range strings.ToUpper(answer) {
		ch := string(c)

		// TODO: consolidate LettersUsed/LettersRemaining, make []rune instead of []string?
		if strings.Contains(g.Alphabet, ch) && !slices.Contains(g.Player.LettersUsed, ch) {
			g.Player.LettersUsed = append(g.Player.LettersUsed, ch)
			turn.NewLettersUsed = append(turn.NewLettersUsed, c)
		}

		g.Player.LettersRemaining[c] = true
	}

	if len(g.Player.LettersUsed) >= len(g.Alphabet) {
		g.Player.LettersUsed = nil
		// TODO having letters remaining AND letters used seems redundant? consider consolidating into single map
		g.Player.LettersRemaining = utils.StringToCharMap(g.Alphabet)

		if g.Player.HealthCurrent < g.Settings.HealthMax {
			g.Player.HealthCurrent++
			turn.Health++
		}
		turn.ExtraLifeGained = true
	}

	slices.Sort(g.Player.LettersUsed)
}

func (g *GameState) HandleFailedTurn() {
	turn := g.CurrentTurn()

	g.Player.Streak = 0
	turn.Streak = 0

	turn.Strikes++
	if turn.Strikes == g.Settings.PromptStrikes {
		turn.TotalTurnDuration = time.Since(turn.TurnStart)
	}

	g.Player.HealthCurrent--
	turn.Health--
}

func (g *GameState) GrantExtraLife() {
	g.Player.LettersUsed = nil
	g.Player.LettersRemaining = utils.StringToCharMap(g.Alphabet)

	if g.Player.HealthCurrent < g.Settings.HealthMax {
		g.Player.HealthCurrent++
	}
}
