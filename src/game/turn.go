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
	Msg			string
}

func (g *GameState) NewTurn() {
	word_idx := rand.Intn(len(g.WordLists.Available))
	word := g.WordLists.Available[word_idx]
	prompt_str := ""

	switch g.Settings.PromptMode {
	case enums.Fuzzy:
		min_idx := 0
		loop_len := min(len(word), g.Settings.PromptLenMax)
		for i := loop_len; i > 0; i-- {
			substr := word[min_idx:]
			rand_max := len(substr) - i
			rand_idx := 0
			if rand_max > 0 {
				rand_idx = rand.Intn(rand_max)
			}
			min_idx += rand_idx + 1
			c := substr[rand_idx]
			prompt_str += string(c)
		}
	case enums.Classic:
		if len(word) <= g.Settings.PromptLenMax {
			prompt_str = word
		} else {
			rand_max := len(word) - g.Settings.PromptLenMax
			rand_idx := rand.Intn(rand_max)
			prompt_str = word[rand_idx:g.Settings.PromptLenMax + rand_idx]
		}
	}

	slog.Debug("New turn", 
		"prompt", prompt_str,
		"sourceWord", word,
		"promptMode", g.Settings.PromptMode.String())

	next_turn := Turn{ 
		SourceWord: word,
		Prompt: prompt_str,
		IsValid: true,
		Strikes: 0,
		Msg: g.CurrentTurn.Msg,
	}

	g.PreviousTurn = g.CurrentTurn
	g.CurrentTurn = next_turn
}

// TODO: maybe return validation msg rather than storing on turn object? store in tui state instead
func (t *Turn) ValidateAnswer(word_lists *WordLists, cfg Settings) {
	slog.Debug("Validating answer", 
		"promptStr", t.Prompt,
		"answer", t.Answer,
		"sourceWord", t.SourceWord,
		"promptMode", cfg.PromptMode.String())

	is_valid := true
	msg := "✓ Correct!"

	if len(t.Answer) == 0 {
		is_valid = false
		msg = "No answer given"
	}

	answer_upper := strings.ToUpper(t.Answer)

	if is_valid && !word_lists.FULL_MAP[t.Answer] {
		is_valid = false
		msg = fmt.Sprintf("Invalid word: %s", answer_upper)
	}

	if is_valid && ((cfg.PromptMode == enums.Fuzzy && !utils.IsFuzzyMatch(t.Answer, t.Prompt)) ||
		(cfg.PromptMode == enums.Classic && !strings.Contains(t.Answer, t.Prompt))) {
			is_valid = false
			msg = fmt.Sprintf("%s does not satisfy prompt", answer_upper)
		}
	
	if is_valid && word_lists.Used[t.Answer] {
		is_valid = false
		msg = fmt.Sprintf("🔒 %s already used", answer_upper)
	}

	if !is_valid {
		t.Strikes++
	}

	if !is_valid && t.Strikes == cfg.PromptStrikesMax {
		msg = fmt.Sprintf("Prompt %s failed. Possible answer: %s", strings.ToUpper(t.Prompt), strings.ToUpper(t.SourceWord))
	}

	t.IsValid = is_valid
	t.Msg = msg

	if is_valid {
		word_idx, _ := slices.BinarySearch(word_lists.Available, t.Answer)
		word_lists.Available = utils.Remove(word_lists.Available, word_idx)
		word_lists.Used[t.Answer] = true
	}
}
