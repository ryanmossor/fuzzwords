package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
)

func main() {
	wordList, err := readLines("./wordlist.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	WORD_MAP := arrToMap(wordList)
	prompt := generatePrompt(WORD_MAP, wordList)
	fmt.Println("prompt:", prompt)

	isValid := isValidAnswer("tio", "overintellectualizations", WORD_MAP)
	fmt.Println(isValid)
}

func generatePrompt(wordMap map[string]bool, wordList []string) string {
	idx := rand.Intn(len(wordMap))
	word := wordList[idx]
	maxPromptLength := 3

	prompt := ""
	minIdx := 0

	for i := min(len(word), maxPromptLength); i > 0; i-- {
		substr := string(word[minIdx:])
		randomMax := len(substr) - i
		randomIdx := rand.Intn(randomMax)
		minIdx += randomIdx + 1
		c := substr[randomIdx]
		prompt += string(c)
	}

	return prompt
}

func isValidAnswer(prompt string, answer string, wordMap map[string]bool) bool {
	answerLc := strings.ToLower(answer)
	
	if !wordMap[answerLc] {
		return false
	}

	subIdx := 0
	for i := 0; i < len(prompt); i++ {
		substr := answerLc[subIdx:]
		currentPromptChar := prompt[i]

		if !strings.Contains(substr, string(currentPromptChar)) {
			return false
		}

		subIdx += strings.Index(substr, string(currentPromptChar)) + 1
	}

	return true
}

func arrToMap(lines []string) map[string]bool {
	var wordMap = make( map[string]bool)
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
