package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

type wordLists struct {
	FULL_MAP   map[string]bool
	available  []string
	used	   map[string]bool
}

type config struct {
	minPromptLen		int
	maxPromptLen		int
	minTurnDuration		int
	maxPromptStrikes	int
	startingHealth		int
	maxHealth			int
	// TODO: add cfg for hints after each strike?
	// hintsEnabled		bool
	// charsPerHint		int
}

type player struct {
	curHealth  int
	cfg		   config
}

type turn struct {
	sourceWord string
	prompt 	   string
	answer     string
	strikes	   int
}

type result struct {
	isValid	   bool
	reason	   string
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
	wordList, err := readLines("./wordlist.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	words := wordLists{FULL_MAP: arrToMap(wordList), available: wordList, used: make(map[string]bool) }

	cfg := config{
		minPromptLen: 2,
		maxPromptLen: 3,
		minTurnDuration: 10,
		maxPromptStrikes: 3,
		startingHealth: 2,
		maxHealth: 5,
	}

	p := player{curHealth: cfg.startingHealth, cfg: cfg}

	reader := bufio.NewReader(os.Stdin)
	clear()

	for len(wordList) > 0 {
		var t turn = generatePrompt(wordList, p)

		for t.strikes < p.cfg.maxPromptStrikes {
			fmt.Println()
			fmt.Fprintf(os.Stdin, "[ Health: %d / %d ]\n", p.curHealth, p.cfg.maxHealth)
			fmt.Fprintf(os.Stdin, "[ Strikes: %d / %d ]\n", t.strikes, p.cfg.maxPromptStrikes)
			fmt.Println("Prompt:", strings.ToUpper(t.prompt))
			fmt.Print("Answer: ")

			answer, _ := reader.ReadString('\n')
			t.answer = strings.ToLower(strings.TrimSpace(answer))

			result := validateAnswer(&t, &words)
			if result.isValid {
				fmt.Println("Correct!")
				time.Sleep(750 * time.Millisecond)
				break
			} else {
				t.strikes += 1
				fmt.Println(result.reason)
			}
		}

		if t.strikes == p.cfg.maxPromptStrikes {
			fmt.Println("Prompt failed. Possible answer:", t.sourceWord)
			p.curHealth -= 1

			if p.curHealth == 0 {
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

func generatePrompt(wordList []string, p player) turn {
	wordIdx := rand.Intn(len(wordList))
	word := wordList[wordIdx]

	promptStr := ""
	minIdx := 0

	loopLen := min(len(word), p.cfg.maxPromptLen)
	for i := loopLen; i > 0; i-- {
		substr := string(word[minIdx:])
		randomMax := len(substr) - i
		randomIdx := 0
		if randomMax > 0 {
			randomIdx = rand.Intn(randomMax)
		}
		minIdx += randomIdx + 1
		c := substr[randomIdx]
		promptStr += string(c)
	}

	return turn{sourceWord: word, prompt: promptStr, strikes: 0}
}

func validateAnswer(t *turn, wordLists *wordLists) result {
	if wordLists.used[t.answer] {
		return result{isValid: false, reason: "Word has already been used. Try again."}
	} else if !wordLists.FULL_MAP[t.answer] {
		return result{isValid: false, reason: "Invalid word. Try again."}
	}

	subIdx := 0
	for i := 0; i < len(t.prompt); i++ {
		substr := t.answer[subIdx:]
		currentPromptChar := t.prompt[i]

		if !strings.Contains(substr, string(currentPromptChar)) {
			return result{isValid: false, reason: "Word does not satisfy the prompt. Try again."}
		}

		subIdx += strings.Index(substr, string(currentPromptChar)) + 1
	}

	wordLists.available = remove(wordLists.available, binarySearch(wordLists.available, t.answer))
	wordLists.used[t.answer] = true

	return result{isValid: true}
}

// func processLetters(t turn, p player) {
// 	TODO: keep track of letters used; increase health when all used
// }

func binarySearch(list []string, target string) int {
	low, high := 0, len(list) - 1
	for low <= high {
		mid := (low + high) / 2
		if list[mid] == target {
			return mid
		} else if list[mid] < target {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}

	return -1
}

// Preserves order, which is necessary for binary search
func remove(list []string, i int) []string {
    return append(list[:i], list[i+1:]...)
}

func arrToMap(lines []string) map[string]bool {
	var wordMap = make(map[string]bool)
	for _, word := range lines {
		wordMap[word] = true
	}
	return wordMap
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
