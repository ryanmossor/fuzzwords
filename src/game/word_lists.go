package game

import "strings"

type wordLists struct {
	available  []string
	used	   map[string]bool
	fullDict   map[string]bool
}

func (g *GameState) WordInDictionary(answer string) bool {
	return g.wordLists.fullDict[strings.ToLower(answer)]
}
