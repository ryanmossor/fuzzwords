package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"slices"
	"strings"
	"time"
)

type Alphabet string 
const (
	DebugAlphabet 	= "ABC"
	EasyAlphabet 	= "ABCDEFGHILMNOPRSTUVW" // J, K, Q, X, Y, Z removed
	MediumAlphabet 	= "ABCDEFGHIJKLMNOPQRSTUVWY" // X and Z removed
	FullAlphabet 	= "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

type PromptMode int
const (
	Fuzzy PromptMode = iota
	Classic
)

type WinCondition int
const (
	Endless WinCondition = iota
	MaxLives
)

type WordLists struct {
	FULL_MAP   map[string]bool
	available  []string
	used	   map[string]bool
}

type GameSettings struct {
	alphabet				string
	health_initial			int
	health_max				int
	prompt_len_min			int
	prompt_len_max			int
	prompt_mode				PromptMode
	prompt_strikes_max		int
	turn_duration_min		int
	win_condition			WinCondition
	// TODO: add cfg for hints after each strike?
	// hints_enabled			bool
	// hint_chars_per_turn		int
}

type Player struct {
	health_current 		int
	health_max			int
	health_display		string
	letters_used		[]string
	letters_remaining 	[]string
}

type Turn struct {
	source_word string
	prompt 	   	string
	answer     	string
	strikes	   	int
}

type Result struct {
	is_valid	bool
	msg	   	   	string
}

func clear() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func main() {
	word_list, err := readLines("./wordlist.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	words := WordLists{
		FULL_MAP: arrToMap(word_list),
		available: word_list,
		used: make(map[string]bool),
	}

	cfg := GameSettings{
		alphabet: DebugAlphabet,
		health_initial: 2,
		health_max: 3,
		prompt_len_max: 3,
		prompt_len_min: 2,
		prompt_mode: Fuzzy,
		prompt_strikes_max: 3,
		turn_duration_min: 10,
		win_condition: Endless,
	}

	player := InitializePlayer(cfg)

	reader := bufio.NewReader(os.Stdin)
	clear()

	for len(words.available) > 0 {
		turn := generatePrompt(words.available, cfg)

		fmt.Fprintf(os.Stdin, "[ Health: %s ]\n", player.health_display)
		fmt.Println("Unused letters:", player.letters_remaining)
		fmt.Println()

		for turn.strikes < cfg.prompt_strikes_max {
			fmt.Fprintf(os.Stdin, "[ Strikes: %d / %d ]\n", turn.strikes, cfg.prompt_strikes_max)
			fmt.Println("Prompt:", strings.ToUpper(turn.prompt))
			fmt.Print("Answer: ")

			answer, _ := reader.ReadString('\n')
			turn.answer = strings.ToLower(strings.TrimSpace(answer))

			result := validateAnswer(&turn, &words, cfg)
			if result.is_valid {
				fmt.Println("Correct!")
				processLetters(turn, &player, cfg)
				time.Sleep(750 * time.Millisecond)
				break
			} else {
				turn.strikes += 1
				fmt.Println(result.msg)
				fmt.Println()
			}
		}

		if cfg.win_condition == MaxLives && player.health_current == cfg.health_max {
			fmt.Println("Max lives achieved -- you win!")
			os.Exit(0)
		}

		if turn.strikes == cfg.prompt_strikes_max {
			fmt.Println("Prompt failed. Possible answer:", turn.source_word)
			player.DecrementHealth()

			if player.health_current == 0 {
				fmt.Println()
				fmt.Println("===== GAME OVER =====")
				fmt.Println()
				os.Exit(0)
			} else {
				time.Sleep(3 * time.Second)
			}
		}

		clear()
	}

	fmt.Println("Congratulations, you used every word in the dictionary.")
	os.Exit(0)
}

func generatePrompt(word_list []string, cfg GameSettings) Turn {
	word_idx := rand.Intn(len(word_list))
	word := word_list[word_idx]
	prompt_str := ""

	switch cfg.prompt_mode {
	case Fuzzy:
		min_idx := 0
		loop_len := min(len(word), cfg.prompt_len_max)
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
	case Classic:
		if len(word) <= cfg.prompt_len_max {
			prompt_str = word
		} else {
			randomMax := len(word) - cfg.prompt_len_max
			randomIdx := rand.Intn(randomMax)
			prompt_str = word[randomIdx:cfg.prompt_len_max + randomIdx]
		}
	}

	return Turn{ source_word: word, prompt: prompt_str, strikes: 0 }
}

func validateAnswer(turn *Turn, word_lists *WordLists, cfg GameSettings) Result {
	if word_lists.used[turn.answer] {
		return Result{ is_valid: false, msg: "Word has already been used. Try again." }
	} else if !word_lists.FULL_MAP[turn.answer] {
		return Result{ is_valid: false, msg: "Invalid word. Try again." }
	}

	switch cfg.prompt_mode {
	case Fuzzy:
		sub_idx := 0
		for i := 0; i < len(turn.prompt); i++ {
			substr := turn.answer[sub_idx:]
			current_prompt_char := turn.prompt[i]

			if !strings.Contains(substr, string(current_prompt_char)) {
				return Result{ is_valid: false, msg: "Word does not satisfy the prompt. Try again." }
			}

			sub_idx += strings.Index(substr, string(current_prompt_char)) + 1
		}
	case Classic:
		if !strings.Contains(turn.answer, string(turn.prompt)) {
			return Result{ is_valid: false, msg: "Word does not satisfy the prompt. Try again." }
		}
	}

	word_idx, _ := slices.BinarySearch(word_lists.available, turn.answer)
	word_lists.available = remove(word_lists.available, word_idx)
	word_lists.used[turn.answer] = true

	return Result{ is_valid: true }
}

func processLetters(turn Turn, player *Player, cfg GameSettings) {
	for i := 0; i < len(turn.answer); i++ {
		c := strings.ToUpper(string(turn.answer[i]))

		if strings.Contains(cfg.alphabet, c) && !slices.Contains(player.letters_used, c) {
			player.letters_used = append(player.letters_used, c)
		}

		if slices.Contains(player.letters_remaining, c) {
			player.letters_remaining = remove(player.letters_remaining, slices.Index(player.letters_remaining, c))
		}
	}

	if len(player.letters_used) >= len(cfg.alphabet) {
		player.IncrementHealth(cfg)
	}

	slices.Sort(player.letters_used)
}

// Preserves order, which is necessary for binary search
func remove(list []string, i int) []string {
    return append(list[:i], list[i+1:]...)
}

func arrToMap(lines []string) map[string]bool {
	var word_map = make(map[string]bool)
	for _, word := range lines {
		word_map[word] = true
	}
	return word_map
}

func readLines(path string) ([]string, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var lines []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }

    return lines, scanner.Err()
} 

func (p *Player) IncrementHealth(cfg GameSettings) {
	p.letters_used = nil
	p.letters_remaining = strings.Split(cfg.alphabet, "")

	if p.health_current < p.health_max {
		p.health_current++
		p.UpdateHealthDisplay()
	}
}

func (p *Player) DecrementHealth() {
	p.health_current--
	p.UpdateHealthDisplay()
}

func (p *Player) UpdateHealthDisplay() {
	health_display := ""

	i := 0
	for i < p.health_current {
		// 🧡💛💚💙🩵💜🖤🤍🤎
		health_display += "🩵"
		i++
	}
	for i < p.health_max {
		health_display += "🤍"
		i++
	}

	p.health_display = health_display
}

func InitializePlayer(cfg GameSettings) Player {
	player := Player{
		health_current: cfg.health_initial,
		health_max: cfg.health_max,
		letters_used: nil,
		letters_remaining: strings.Split(cfg.alphabet, ""),
	}
	player.UpdateHealthDisplay()
	return player
}
