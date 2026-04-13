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
	TurnStart			time.Time
	TotalTurnDuration	time.Duration

	SourceWord 			string
	Prompt 	   			string
	Answer				string
	Guesses				int

	Strikes	   			int
	StrikeStart			time.Time
	StrikeDuration		time.Duration

	Solved				bool
	ExtraLifeGained		bool

	LettersRemaining	map[rune]bool
	NewLettersUsed		[]rune
	UniqueLetterCount	int
	Streak				int
	Health				int
	// may be able to get rid of validation_msg on ui state? maybe store on game state instead?
	// - don't need to colorize in UI if not showing possible answer anymore, so maybe
	// just use PrevTurn().Solved w/ GameState.ValidationMsg instead? entirely red/green depending on if solved
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
		// Timer expiration: reset to random time between 10s (or turn min if larger) and 30s
		turn_duration_min := max(g.Settings.TurnDurationMin, 10)
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
		TurnStart: now,

		SourceWord: word,
		Prompt: prompt,
		Answer: "",
		Guesses: 0,

		Strikes: 0,
		StrikeStart: now,
		StrikeDuration: turn_duration,

		Solved: false,
		ExtraLifeGained: false,

		LettersRemaining: maps.Clone(g.Player.LettersRemaining),
		NewLettersUsed: make([]rune, 0, 16),
		Health: g.Player.HealthCurrent,
	})
}

func (g *GameState) StartStrikeTimer() {
	turn_duration_min := max(15, g.Settings.TurnDurationMin)
	duration_sec := utils.RandomBetween(turn_duration_min, 30)

	g.CurrentTurn().StrikeStart = time.Now()
	g.CurrentTurn().StrikeDuration = time.Duration(duration_sec) * time.Second
}

func (g GameState) TimeRemaining() time.Duration {
	return g.CurrentTurn().StrikeStart.
		Add(g.CurrentTurn().StrikeDuration).
		Sub(time.Now())
}

type AnswerResult struct {
	IsValid				bool
	ExtraLifeGained		bool
	GameOver			bool
	Msg					string
}

func (g *GameState) SubmitAnswer(answer string) AnswerResult {
	is_valid, msg := g.validateAnswer(answer)
	if !is_valid {
		return AnswerResult{ IsValid: false, Msg: msg }
	}

	g.handleCorrectAnswer(answer)
	result := AnswerResult{
		IsValid: true,
		ExtraLifeGained: g.CurrentTurn().ExtraLifeGained,
		Msg: msg,
	}

	if g.EndGameIfOver() {
		result.GameOver = true
		return result
	}

	g.NewTurn(false)
	return result
}

func (g *GameState) validateAnswer(answer string) (bool, string) {
	is_valid := true
	incr_guess_count := true
	answer_upper := strings.ToUpper(answer)
	msg := fmt.Sprintf("✓ %s", answer_upper)

	if len(answer) == 0 {
		is_valid = false
		incr_guess_count = false
		msg = "No answer given"
	}

	if is_valid && !g.WordLists.FULL_MAP[answer] {
		is_valid = false
		msg = fmt.Sprintf("Invalid word: %s", answer_upper)
	}

	is_match := false
	if g.Settings.PromptMode == enums.PromptModeFuzzy {
		is_match = utils.IsFuzzyMatch(answer, g.CurrentTurn().Prompt)
	}
	if g.Settings.PromptMode == enums.PromptModeClassic {
		is_match = strings.Contains(answer, g.CurrentTurn().Prompt)
	}

	if is_valid && !is_match {
		is_valid = false
		msg = fmt.Sprintf("%s does not satisfy prompt", answer_upper)
	}

	if is_valid && g.WordLists.Used[answer] {
		is_valid = false
		msg = fmt.Sprintf("🔒 %s already used", answer_upper)
	}

	slog.Debug("Answer validated",
		"startUnixTs", g.StartUnixTs,
		"prompt", g.CurrentTurn().Prompt,
		"sourceWord", g.CurrentTurn().SourceWord,
		"answer", answer,
		"isValid", is_valid,
		"validationMsg", msg,
		"promptMode", g.Settings.PromptMode.String())

	if is_valid {
		word_idx, found := slices.BinarySearch(g.WordLists.Available, answer)
		assert.Assert(found, "Validated answer not found in available word list",
			"startUnixTs", g.StartUnixTs,
			"prompt", g.CurrentTurn().Prompt,
			"answer", answer,
			"wordIdx", word_idx,
			"actualWordAtIdx", g.WordLists.Available[word_idx],
			"remainingWords", len(g.WordLists.Available),
			"alreadyUsed", g.WordLists.Used[answer])

		g.WordLists.Available = utils.Remove(g.WordLists.Available, word_idx)
		g.WordLists.Used[answer] = true
	}

	if incr_guess_count {
		g.CurrentTurn().Guesses++
	}

	return is_valid, msg
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

func (g GameState) CurrentTurn() *Turn {
	assert.Assert(len(g.turns) > 0, "Attempted to access current turn before game initialized")
	return &g.turns[len(g.turns) - 1]
}

func (g GameState) PreviousTurn() (*Turn, bool) {
	if len(g.turns) <= 1 {
		return nil, false
	}
	return &g.turns[len(g.turns) - 2], true
}

// TODO: this takes turn idx, but i'm referencing turns by number (1-based) in many places.
// Consider pros/cons of idx vs turn number
func (g GameState) GetTurn(idx int) *Turn {
	clamped_idx := utils.Clamp(idx, 0, len(g.turns) - 1)
	return &g.turns[clamped_idx]
}

func (g GameState) NextFailedTurnIdx(turn_idx_cur int) int {
	for i := turn_idx_cur; i < len(g.turns); i++ {
		if (g.turns[i].Strikes > 0 || !g.turns[i].Solved) && i > turn_idx_cur {
			return i
		}
	}
	return turn_idx_cur
}

func (g GameState) PrevFailedTurnIdx(turn_idx_cur int) int {
	for i := turn_idx_cur; i >= 0; i-- {
		if (g.turns[i].Strikes > 0 || !g.turns[i].Solved) && i < turn_idx_cur {
			return i
		}
	}
	return turn_idx_cur
}
