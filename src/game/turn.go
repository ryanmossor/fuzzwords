package game

import (
	"fzw/src/enums"
	"fzw/src/utils"
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

func NewTurn(word_list []string, cfg Settings) Turn {
	word_idx := rand.Intn(len(word_list))
	word := word_list[word_idx]
	prompt_str := ""

	switch cfg.PromptMode {
	case enums.Fuzzy:
		min_idx := 0
		loop_len := min(len(word), cfg.PromptLenMax)
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
		if len(word) <= cfg.PromptLenMax {
			prompt_str = word
		} else {
			rand_max := len(word) - cfg.PromptLenMax
			rand_idx := rand.Intn(rand_max)
			prompt_str = word[rand_idx:cfg.PromptLenMax + rand_idx]
		}
	}

	return Turn{ 
		SourceWord: word,
		Prompt: prompt_str,
		Strikes: 0,
	}
}

func (t *Turn) ValidateAnswer(word_lists *WordLists, cfg Settings) {
	if !word_lists.FULL_MAP[t.Answer] {
		t.IsValid = false
		t.Msg = "Invalid word. Try again." 
		return
	}

	switch cfg.PromptMode {
	case enums.Fuzzy:
		sub_idx := 0
		for i := 0; i < len(t.Prompt); i++ {
			substr := t.Answer[sub_idx:]
			current_prompt_char := t.Prompt[i]

			if !strings.Contains(substr, string(current_prompt_char)) {
				t.IsValid = false
				t.Msg = "Word does not satisfy the prompt. Try again." 
				return
			}

			sub_idx += strings.Index(substr, string(current_prompt_char)) + 1
		}
	case enums.Classic:
		if !strings.Contains(t.Answer, string(t.Prompt)) {
			t.IsValid = false
			t.Msg = "Word does not satisfy the prompt. Try again." 
			return
		}
	}
	
	if word_lists.Used[t.Answer] {
		t.IsValid = false
		t.Msg = "Word has already been used. Try again."
		return
	}

	word_idx, _ := slices.BinarySearch(word_lists.Available, t.Answer)
	word_lists.Available = utils.Remove(word_lists.Available, word_idx)
	word_lists.Used[t.Answer] = true

	t.IsValid = true
	t.Msg = "Correct!"
}
