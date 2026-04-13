package game

import "strings"

type WordLists struct {
	FULL_MAP   map[string]bool
	Available  []string
	Used	   map[string]bool
}

func (g *GameState) WordInDictionary(answer string) bool {
	return g.wordLists.FULL_MAP[strings.ToLower(answer)]
}
