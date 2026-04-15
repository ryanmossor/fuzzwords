package game

import (
	"fmt"
	"fzwds/src/assert"
	"fzwds/src/enums"
	"fzwds/src/utils"
	"log/slog"
	"maps"
	"math/rand"
	"slices"
	"strings"
	"time"
)

type Turn struct {
	TurnNumber			int
	FinalTurn			bool
	turnStart			time.Time
	TotalTurnDuration	time.Duration

	// TODO: make private and getter returns nothing if game is active
	SourceWord 			string
	Prompt 	   			string
	Answer				string
	Guesses				int

	Strikes	   			int
	strikeStart			time.Time
	strikeDuration		time.Duration

	Solved				bool
	ExtraLifeGained		bool

	LettersRemaining	map[rune]bool
	NewLettersUsed		[]rune
	UniqueLetterCount	int
	Streak				int
	Health				int
}

// TODO: ensure next prompt is different from previous if previous prompt was failed
func (g *Game) newTurn(first_turn bool) {
	word_idx := rand.Intn(len(g.wordLists.available))
	word := g.wordLists.available[word_idx]

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
		prompt = createFuzzyPrompt(word, prompt_len, g.Settings.Dictionary)
	case enums.PromptModeClassic:
		// TODO: classic prompts can contain hyphens/symbols in pokemon names bc it's just a substring
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
		// Game start: default to 30s
		turn_duration = 30 * time.Second

	} else if g.TimeRemaining() <= 0 {
		// Timer expiration: reset to random time between 15s (or turn min if larger) and 30s
		turn_duration_min := max(g.Settings.TurnDurationMin, 15)
		turn_duration_max := 30
		rand_sec := utils.RandomBetween(turn_duration_min, turn_duration_max)
		turn_duration = time.Duration(rand_sec) * time.Second

	} else if g.TimeRemaining().Seconds() < float64(g.Settings.TurnDurationMin) {
		// Correct answer: reset timer to TurnDurationMin if timer is < TurnDurationMin
		turn_duration = time.Duration(g.Settings.TurnDurationMin) * time.Second

	} else {
		// Correct answer: do nothing if timer > TurnDurationMin
		turn_duration = g.TimeRemaining()
	}

	now := time.Now()
	g.turns = append(g.turns, Turn {
		TurnNumber: g.CurrentTurnNumber() + 1,
		turnStart: now,

		SourceWord: word,
		Prompt: prompt,
		Answer: "",
		Guesses: 0,

		Strikes: 0,
		strikeStart: now,
		strikeDuration: turn_duration,

		Solved: false,
		ExtraLifeGained: false,

		LettersRemaining: maps.Clone(g.Player.LettersRemaining),
		NewLettersUsed: make([]rune, 0, 16),
		Health: g.Player.HealthCurrent,
	})
	g.timerId++
}

func (g *Game) startStrikeTimer() {
	turn_duration_min := max(15, g.Settings.TurnDurationMin)
	duration_sec := utils.RandomBetween(turn_duration_min, 30)

	g.CurrentTurn().strikeStart = time.Now()
	g.CurrentTurn().strikeDuration = time.Duration(duration_sec) * time.Second
	g.timerId++
}

func (g Game) TimeRemaining() time.Duration {
	return g.CurrentTurn().strikeStart.
		Add(g.CurrentTurn().strikeDuration).
		Sub(time.Now())
}

type answerResult struct {
	accepted	bool
	reason		string
}

func (g *Game) validateAnswer(answer string) answerResult {
	result := answerResult{ accepted: true }
	incr_guess_count := true

	if len(answer) == 0 {
		incr_guess_count = false
		result.accepted = false
		result.reason = "No answer given"
	}

	if result.accepted && !g.wordLists.fullDict[answer] {
		result.accepted = false
		result.reason = fmt.Sprintf("Invalid word: %s", strings.ToUpper(answer))
	}

	is_match := false
	if g.Settings.PromptMode == enums.PromptModeFuzzy {
		is_match = utils.IsFuzzyMatch(answer, g.CurrentTurn().Prompt)
	}
	if g.Settings.PromptMode == enums.PromptModeClassic {
		is_match = strings.Contains(answer, g.CurrentTurn().Prompt)
	}

	if result.accepted && !is_match {
		result.accepted = false
		result.reason = fmt.Sprintf("%s does not satisfy prompt", strings.ToUpper(answer))
	}

	if result.accepted && g.wordLists.used[answer] {
		result.accepted = false
		result.reason = fmt.Sprintf("🔒 %s already used", strings.ToUpper(answer))
	}

	slog.Debug("Answer validated",
		"startUnixTs", g.startUnixTs,
		"prompt", g.CurrentTurn().Prompt,
		"sourceWord", g.CurrentTurn().SourceWord,
		"answer", answer,
		"accepted", result.accepted,
		"reason", result.reason,
		"promptMode", g.Settings.PromptMode.String())

	if result.accepted {
		word_idx, found := slices.BinarySearch(g.wordLists.available, answer)
		assert.Assert(found, "Validated answer not found in available word list",
			"startUnixTs", g.startUnixTs,
			"prompt", g.CurrentTurn().Prompt,
			"answer", answer,
			"wordIdx", word_idx,
			"actualWordAtIdx", g.wordLists.available[word_idx],
			"remainingWords", len(g.wordLists.available),
			"alreadyUsed", g.wordLists.used[answer])

		g.wordLists.available = utils.Remove(g.wordLists.available, word_idx)
		g.wordLists.used[answer] = true
	}

	if incr_guess_count {
		g.CurrentTurn().Guesses++
	}

	return result
}

func createFuzzyPrompt(word string, prompt_len int, dict enums.Dictionary) string {
	stripped_word := word
	if dict == enums.Pokemon {
		stripped_word = utils.StripNumbersAndSymbols(word)
	}

	if len(stripped_word) <= prompt_len {
		return stripped_word
	}

	var prompt string
	rand_min := 0

	for i := prompt_len; i > 0; i-- {
		rand_max := len(stripped_word) - i
		rand_idx := utils.RandomBetween(rand_min, rand_max)

		if i == prompt_len && rand_idx == rand_max {
			return prompt + stripped_word[rand_idx:]
		}

		prompt += string(stripped_word[rand_idx])
		rand_min = rand_idx + 1
	}

	return prompt
}

func (g *Game) handleCorrectAnswer(answer string) bool {
	turn := g.CurrentTurn()
	turn.TotalTurnDuration = time.Since(turn.turnStart)
	turn.Solved = true
	turn.Answer = answer
	turn.UniqueLetterCount = utils.CountUniqueLetters(answer)

	g.Player.streak++
	turn.Streak = g.Player.streak

	for _, c := range strings.ToUpper(answer) {
		// TODO: consolidate LettersUsed/LettersRemaining, make []rune instead of []string?
		if !slices.Contains(g.Player.LettersUsed, c) && strings.ContainsRune(g.Settings.Alphabet.Letters(), c) {
			g.Player.LettersUsed = append(g.Player.LettersUsed, c)
			turn.NewLettersUsed = append(turn.NewLettersUsed, c)
		}

		g.Player.LettersRemaining[c] = true
	}

	if len(g.Player.LettersUsed) >= len(g.Settings.Alphabet.Letters()) {
		g.Player.LettersUsed = make([]rune, 0, len(g.Settings.Alphabet.Letters()))
		// TODO having letters remaining AND letters used seems redundant? consider consolidating into single map
		g.Player.LettersRemaining = utils.StringToCharMap(g.Settings.Alphabet.Letters())

		if g.Player.HealthCurrent < g.Settings.HealthMax {
			g.Player.HealthCurrent++
			turn.Health++
		}

		turn.ExtraLifeGained = true
	}

	return turn.ExtraLifeGained
}
