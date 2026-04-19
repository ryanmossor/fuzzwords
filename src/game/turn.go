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
	turnNumber			int
	finalTurn			bool
	turnStart			time.Time
	totalTurnDuration	time.Duration

	sourceWord 			string
	prompt 	   			string
	answer				string
	guesses				int

	strikes	   			int
	strikeStart			time.Time
	strikeDuration		time.Duration

	solved				bool
	extraLifeGained		bool

	lettersUsed			map[rune]bool
	newLettersUsed		[]rune
	uniqueLetterCount	int
	streak				int
	health				int
}

func (t Turn) TurnNumber() int {
	return t.turnNumber
}

func (t Turn) FinalTurn() bool {
	return t.finalTurn
}

func (t Turn) TotalTurnDuration() time.Duration {
	return t.totalTurnDuration
}

// TODO: some way of preventing this if game is still active (but allow access in debug mode)?
func (t Turn) SourceWord() string {
	return t.sourceWord
}

func (t Turn) Prompt() string {
	return t.prompt
}

func (t Turn) Answer() string {
	return t.answer
}

func (t Turn) Guesses() int {
	return t.guesses
}

func (t Turn) Strikes() int {
	return t.strikes
}

func (t Turn) Solved() bool {
	return t.solved
}

func (t Turn) ExtraLifeGained() bool {
	return t.extraLifeGained
}

func (t Turn) LettersUsed() map[rune]bool {
	return t.lettersUsed
}

func (t Turn) NewLettersUsed() []rune {
	return t.newLettersUsed
}

func (t Turn) UniqueLetterCount() int {
	return t.uniqueLetterCount
}

func (t Turn) Streak() int {
	return t.streak
}

func (t Turn) Health() int {
	return t.health
}

type TurnTransition int
const (
	TransitionFirstTurn TurnTransition = iota
	TransitionSolved
	TransitionTimeout
)

// TODO: ensure next prompt is different from previous if previous prompt was failed
func (g *Game) newTurn(reason TurnTransition) Turn {
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

	switch reason {
	case TransitionFirstTurn:
		turn_duration = 30 * time.Second

	case TransitionSolved:
		remaining := max(0, g.TimeRemaining())
		turn_duration = max(remaining, time.Duration(g.Settings.TurnDurationMin) * time.Second)

	case TransitionTimeout:
		// Reset to random time between 15s (or turn min if larger) and 30s
		min_sec := max(g.Settings.TurnDurationMin, 15)
		max_sec := 30
		rand_sec := utils.RandomBetween(min_sec, max_sec)
		turn_duration = time.Duration(rand_sec) * time.Second
	}

	now := time.Now()
	turn := Turn {
		turnNumber: g.CurrentTurnNumber() + 1,
		turnStart: now,

		sourceWord: word,
		prompt: prompt,
		answer: "",
		guesses: 0,

		strikes: 0,
		strikeStart: now,
		strikeDuration: turn_duration,

		solved: false,
		extraLifeGained: false,

		lettersUsed: maps.Clone(g.player.lettersUsed),
		newLettersUsed: make([]rune, 0, 16),
		health: g.player.healthCurrent,
	}
	g.turns = append(g.turns, turn)

	return turn
}

func (g *Game) startStrikeTimer() {
	min_sec := max(15, g.Settings.TurnDurationMin)
	duration_sec := utils.RandomBetween(min_sec, 30)

	g.currentTurn().strikeStart = time.Now()
	g.currentTurn().strikeDuration = time.Duration(duration_sec) * time.Second
}

func (g Game) TimeRemaining() time.Duration {
	return g.currentTurn().strikeStart.
		Add(g.currentTurn().strikeDuration).
		Sub(time.Now())
}

type answerResult struct {
	accepted	bool
	// TODO make this an enum which UI then generates display message from?
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
		is_match = utils.IsFuzzyMatch(answer, g.currentTurn().prompt)
	}
	if g.Settings.PromptMode == enums.PromptModeClassic {
		is_match = strings.Contains(answer, g.currentTurn().prompt)
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
		"prompt", g.currentTurn().prompt,
		"sourceWord", g.currentTurn().sourceWord,
		"answer", answer,
		"accepted", result.accepted,
		"reason", result.reason,
		"promptMode", g.Settings.PromptMode.String())

	if result.accepted {
		word_idx, found := slices.BinarySearch(g.wordLists.available, answer)
		assert.Assert(found, "Validated answer not found in available word list",
			"startUnixTs", g.startUnixTs,
			"prompt", g.currentTurn().prompt,
			"answer", answer,
			"wordIdx", word_idx,
			"actualWordAtIdx", g.wordLists.available[word_idx],
			"remainingWords", len(g.wordLists.available),
			"alreadyUsed", g.wordLists.used[answer])

		g.wordLists.available = utils.Remove(g.wordLists.available, word_idx)
		g.wordLists.used[answer] = true
	}

	if incr_guess_count {
		g.currentTurn().guesses++
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

func (g *Game) handleCorrectAnswer(answer string) Turn {
	turn := g.currentTurn()
	turn.totalTurnDuration = time.Since(turn.turnStart)
	turn.solved = true
	turn.answer = answer
	turn.uniqueLetterCount = utils.CountUniqueLetters(answer)

	g.player.streak++
	turn.streak = g.player.streak

	for _, c := range strings.ToUpper(answer) {
		if used, is_in_alphabet := g.player.lettersUsed[c]; !used && is_in_alphabet {
			turn.newLettersUsed = append(turn.newLettersUsed, c)
			g.player.lettersUsed[c] = true
		}
	}

	all_used := true
	for _, c := range g.Settings.Alphabet.Letters() {
		if !g.player.lettersUsed[c] {
			all_used = false
			break
		}
	}

	if all_used {
		g.player.lettersUsed = utils.StringToCharMap(g.Settings.Alphabet.Letters())

		if g.player.healthCurrent < g.Settings.HealthMax {
			g.player.healthCurrent++
			turn.health++
		}

		turn.extraLifeGained = true
	}

	return *turn
}
