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
	SourceWord 		string
	PossibleAnswer	string
	Prompt 	   		string
	Strikes	   		int
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
			rand_max := len(word) - prompt_len
			rand_idx := rand.Intn(rand_max)
			prompt = word[rand_idx:prompt_len + rand_idx]
		}
	}

	next_turn := Turn{ 
		SourceWord: word,
		PossibleAnswer: g.getShortPossibleAnswer(prompt),
		Prompt: prompt,
		Strikes: 0,
	}

	g.PreviousTurn = g.CurrentTurn
	g.CurrentTurn = next_turn
}

func (g *GameState) ValidateAnswer(answer string) (bool, string) {
	is_valid := true
	answer_upper := strings.ToUpper(answer)
	msg := fmt.Sprintf("✓ %s", answer_upper)

	if len(answer) == 0 {
		is_valid = false
		msg = "No answer given"
	}

	if is_valid && !g.WordLists.FULL_MAP[answer] {
		is_valid = false
		msg = fmt.Sprintf("Invalid word: %s", answer_upper)
	}

	fuzzy_match := g.Settings.PromptMode == enums.Fuzzy && utils.IsFuzzyMatch(answer, g.CurrentTurn.Prompt)
	classic_match := g.Settings.PromptMode == enums.Classic && strings.Contains(answer, g.CurrentTurn.Prompt)
	if is_valid && !(fuzzy_match || classic_match) {
		is_valid = false
		msg = fmt.Sprintf("%s does not satisfy prompt", answer_upper)
	}
	
	if is_valid && g.WordLists.Used[answer] {
		is_valid = false
		msg = fmt.Sprintf("🔒 %s already used", answer_upper)
	}

	slog.Info("Answer validated",
		"prompt", g.CurrentTurn.Prompt,
		"sourceWord", g.CurrentTurn.SourceWord,
		"possibleAnswer", g.CurrentTurn.PossibleAnswer,
		"answer", answer,
		"isValid", is_valid,
		"validationMsg", msg,
		"promptMode", g.Settings.PromptMode.String())

	if is_valid {
		word_idx, _ := slices.BinarySearch(g.WordLists.Available, answer)
		g.WordLists.Available = utils.Remove(g.WordLists.Available, word_idx)
		g.WordLists.Used[answer] = true
	}

	return is_valid, msg
}

func (g *GameState) getShortPossibleAnswer(prompt string) string {
	is_valid_word := g.WordLists.FULL_MAP[prompt]
	has_been_used := g.WordLists.Used[prompt]
	if is_valid_word && !has_been_used {
		return prompt
	}

	possible_answer := strings.Repeat("a", 50)
	for _, word := range g.WordLists.Available {
		if len(word) < len(prompt) || len(word) > len(possible_answer) {
			continue
		}

		is_match := false
		if g.Settings.PromptMode == enums.Fuzzy {
			is_match = utils.IsFuzzyMatch(word, prompt)
		} else {
			is_match = strings.Contains(word, prompt)
		}

		if !is_match {
			continue
		}

		if len(word) < len(possible_answer) {
			possible_answer = word
		}

		// Accept a word up to 2 chars longer than length of prompt
		if len(possible_answer) <= len(prompt) + 2 {
			break
		}
	}

	return possible_answer
}
