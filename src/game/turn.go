package game

import (
	"fmt"
	"fzwds/src/enums"
	"fzwds/src/utils"
	"log/slog"
	"math/rand"
	"slices"
	"strings"
)

type Turn struct {
	SourceWord 	string
	Prompt 	   	string
	Strikes	   	int

	Answer     	string
	IsValid		bool
}

// TODO: ensure next prompt is different from previous if previous prompt was failed
func (g *GameState) NewTurn() {
	word_idx := rand.Intn(len(g.WordLists.Available))
	word := g.WordLists.Available[word_idx]

	var prompt string
	prompt_len := rand.Intn(g.Settings.PromptLenMax - g.Settings.PromptLenMin + 1) + g.Settings.PromptLenMin 

	switch g.Settings.PromptMode {
	case enums.Fuzzy:
		prompt = utils.CreateFuzzyPrompt(word, prompt_len)
	case enums.Classic:
		if len(word) <= g.Settings.PromptLenMax {
			prompt = word
		} else {
			rand_max := len(word) - g.Settings.PromptLenMax
			rand_idx := rand.Intn(rand_max)
			prompt = word[rand_idx:g.Settings.PromptLenMax + rand_idx]
		}
	}

	next_turn := Turn{ 
		SourceWord: word,
		Prompt: prompt,
		Strikes: 0,
		IsValid: true,
	}

	g.PreviousTurn = g.CurrentTurn
	g.CurrentTurn = next_turn
}

func (g *GameState) ValidateAnswer() string {
	is_valid := true
	answer_upper := strings.ToUpper(g.CurrentTurn.Answer)
	msg := fmt.Sprintf("✓ %s", answer_upper)

	if len(g.CurrentTurn.Answer) == 0 {
		is_valid = false
		msg = "No answer given"
	}

	if is_valid && !g.WordLists.FULL_MAP[g.CurrentTurn.Answer] {
		is_valid = false
		msg = fmt.Sprintf("Invalid word: %s", answer_upper)
	}

	fuzzy_match := g.Settings.PromptMode == enums.Fuzzy && utils.IsFuzzyMatch(g.CurrentTurn.Answer, g.CurrentTurn.Prompt)
	classic_match := g.Settings.PromptMode == enums.Classic && strings.Contains(g.CurrentTurn.Answer, g.CurrentTurn.Prompt)
	if is_valid && !(fuzzy_match || classic_match) {
		is_valid = false
		msg = fmt.Sprintf("%s does not satisfy prompt", answer_upper)
	}
	
	if is_valid && g.WordLists.Used[g.CurrentTurn.Answer] {
		is_valid = false
		msg = fmt.Sprintf("🔒 %s already used", answer_upper)
	}

	slog.Debug("Answer validated", 
		"prompt", g.CurrentTurn.Prompt,
		"sourceWord", g.CurrentTurn.SourceWord,
		"answer", g.CurrentTurn.Answer,
		"isValid", is_valid,
		"validationMsg", msg,
		"promptMode", g.Settings.PromptMode.String())

	if !is_valid && g.CurrentTurn.Strikes == g.Settings.PromptStrikesMax {
		msg = fmt.Sprintf(
			"Prompt %s failed. Possible answer: %s",
			strings.ToUpper(g.CurrentTurn.Prompt),
			strings.ToUpper(g.CurrentTurn.SourceWord))
	}

	g.CurrentTurn.IsValid = is_valid

	if is_valid {
		word_idx, _ := slices.BinarySearch(g.WordLists.Available, g.CurrentTurn.Answer)
		g.WordLists.Available = utils.Remove(g.WordLists.Available, word_idx)
		g.WordLists.Used[g.CurrentTurn.Answer] = true
	}

	return msg
}
