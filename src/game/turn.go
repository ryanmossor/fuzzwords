package game

import (
	"fmt"
	"fzwds/src/assert"
	"fzwds/src/enums"
	"fzwds/src/utils"
	"log/slog"
	"math/rand"
	"slices"
	"strings"
	"time"
)

type Turn struct {
	SourceWord 		string
	PossibleAnswer	string
	Prompt 	   		string
	Strikes	   		int
	TurnStart		time.Time
	TurnDuration	time.Duration
}

// TODO: ensure next prompt is different from previous if previous prompt was failed
func (g *GameState) NewTurn(first_turn bool) {
	word_idx := rand.Intn(len(g.WordLists.Available))
	word := g.WordLists.Available[word_idx]

	assert.Assert(word != "", "Random word must not be empty", "word", word, "wordIdx", word_idx)

	var prompt string
	prompt_len := utils.RandomBetween(g.Settings.PromptLenMin, g.Settings.PromptLenMax)

	assert.Assert(prompt_len >= g.Settings.PromptLenMin, "Prompt len must be >= PromptLenMin",
		"randPromptLen", prompt_len,
		"promptLenMin", g.Settings.PromptLenMin)
	assert.Assert(prompt_len <= g.Settings.PromptLenMax, "Prompt len must be <= PromptLenMax",
		"randPromptLen", prompt_len,
		"promptLenMax", g.Settings.PromptLenMax)

	switch g.Settings.PromptMode {
	case enums.PromptModeFuzzy:
		prompt = utils.CreateFuzzyPrompt(word, prompt_len)
	case enums.PromptModeClassic:
		if len(word) <= g.Settings.PromptLenMax {
			prompt = word
		} else {
			rand_max := len(word) - prompt_len
			rand_idx := rand.Intn(rand_max)
			prompt = word[rand_idx:prompt_len + rand_idx]
		}
	}

	assert.Assert(prompt != "", "Prompt must not be empty",
		"word", word,
		"wordIdx", word_idx,
		"prompt", prompt)

	var turn_duration time.Duration
	if first_turn {
		turn_duration = 30 * time.Second
	} else if g.TimeRemaining() <= 0 {
		turn_duration_min := max(g.Settings.TurnDurationMin, 10)
		turn_duration_max := 30
		rand_sec := utils.RandomBetween(turn_duration_min, turn_duration_max)
		turn_duration = time.Duration(rand_sec) * time.Second
	} else if g.TimeRemaining().Seconds() < float64(g.Settings.TurnDurationMin) {
		turn_duration = time.Duration(g.Settings.TurnDurationMin) * time.Second
	} else {
		turn_duration = g.TimeRemaining()
	}

	next_turn := Turn {
		SourceWord: word,
		PossibleAnswer: g.getPossibleAnswer(prompt, word),
		Prompt: prompt,
		Strikes: 0,
		TurnStart: time.Now(),
		TurnDuration: turn_duration,
	}

	g.PreviousTurn = g.CurrentTurn
	g.CurrentTurn = next_turn
}

func (g *GameState) StartTurn(duration_sec int) {
	g.CurrentTurn.TurnStart = time.Now()
	g.CurrentTurn.TurnDuration = time.Duration(duration_sec) * time.Second
}

func (g *GameState) TimeRemaining() time.Duration {
	return g.CurrentTurn.TurnStart.
		Add(g.CurrentTurn.TurnDuration).
		Sub(time.Now())
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

	fuzzy_match := g.Settings.PromptMode == enums.PromptModeFuzzy && utils.IsFuzzyMatch(answer, g.CurrentTurn.Prompt)
	classic_match := g.Settings.PromptMode == enums.PromptModeClassic && strings.Contains(answer, g.CurrentTurn.Prompt)
	if is_valid && !(fuzzy_match || classic_match) {
		is_valid = false
		msg = fmt.Sprintf("%s does not satisfy prompt", answer_upper)
	}

	if is_valid && g.WordLists.Used[answer] {
		is_valid = false
		msg = fmt.Sprintf("🔒 %s already used", answer_upper)
	}

	slog.Debug("Answer validated",
		"startUnixTs", g.StartUnixTs,
		"prompt", g.CurrentTurn.Prompt,
		"sourceWord", g.CurrentTurn.SourceWord,
		"possibleAnswer", g.CurrentTurn.PossibleAnswer,
		"answer", answer,
		"isValid", is_valid,
		"validationMsg", msg,
		"promptMode", g.Settings.PromptMode.String())

	if is_valid {
		word_idx, found := slices.BinarySearch(g.WordLists.Available, answer)
		assert.Assert(found, "Validated answer not found in available word list",
			"startUnixTs", g.StartUnixTs,
			"prompt", g.CurrentTurn.Prompt,
			"answer", answer,
			"wordIdx", word_idx,
			"actualWordAtIdx", g.WordLists.Available[word_idx],
			"remainingWords", len(g.WordLists.Available),
			"alreadyUsed", g.WordLists.Used[answer])

		g.WordLists.Available = utils.Remove(g.WordLists.Available, word_idx)
		g.WordLists.Used[answer] = true
	}

	return is_valid, msg
}

func (g *GameState) getPossibleAnswer(prompt, source_word string) string {
	is_valid_word := g.WordLists.FULL_MAP[prompt]
	has_been_used := g.WordLists.Used[prompt]
	if is_valid_word && !has_been_used {
		return prompt
	}

	possible_answer := source_word
	for _, word := range g.WordLists.Available {
		if len(word) < len(prompt) || len(word) > len(possible_answer) {
			continue
		}

		is_match := false
		if g.Settings.PromptMode == enums.PromptModeFuzzy {
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

func (g *GameState) GetTurnFailureMessage() string {
	if g.CurrentTurn.Strikes == g.Settings.PromptStrikes {
		return fmt.Sprintf(
			"Prompt %s failed. Possible solve: {solve}",
			strings.ToUpper(g.CurrentTurn.Prompt))
	}

	return ""
}
