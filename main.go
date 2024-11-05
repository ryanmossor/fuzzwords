package main

import (
	"bufio"
	"fmt"
	"fzw/src/enums"
	"fzw/src/game"
	"fzw/src/utils"
	"os"
	"strings"
	"time"
)

func main() {
	cfg := game.InitializeSettings()
	player := game.InitializePlayer(cfg)

	word_list, err := utils.ReadLines("./wordlist.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	words := game.WordLists{
		FULL_MAP: utils.ArrToMap(word_list),
		Available: word_list,
		Used: make(map[string]bool),
	}

	reader := bufio.NewReader(os.Stdin)
	utils.ClearWindow()

	for len(words.Available) > 0 {
		turn := game.NewTurn(words.Available, cfg)

		fmt.Fprintf(os.Stdin, "[ Health: %s ]\n", player.HealthDisplay)
		fmt.Println("Unused letters:", player.LettersRemaining)
		fmt.Println()

		for turn.Strikes < cfg.PromptStrikesMax {
			fmt.Fprintf(os.Stdin, "[ Strikes: %d / %d ]\n", turn.Strikes, cfg.PromptStrikesMax)
			fmt.Println("Prompt:", strings.ToUpper(turn.Prompt))
			fmt.Print("Answer: ")

			answer, _ := reader.ReadString('\n')
			turn.Answer = strings.ToLower(strings.TrimSpace(answer))

			turn.ValidateAnswer(&words, cfg)
			fmt.Println(turn.Msg)

			if turn.IsValid {
				player.HandleCorrectAnswer(turn.Answer, &player, cfg)
				time.Sleep(750 * time.Millisecond)
				break
			} else {
				turn.Strikes += 1
				fmt.Println()
			}
		}

		if cfg.WinCondition == enums.MaxLives && player.HealthCurrent == cfg.HealthMax {
			fmt.Println("Max lives achieved -- you win!")
			os.Exit(0)
		}

		if turn.Strikes == cfg.PromptStrikesMax {
			fmt.Println("Prompt failed. Possible answer:", turn.SourceWord)
			player.DecrementHealth()

			if player.HealthCurrent == 0 {
				fmt.Println()
				fmt.Println("===== GAME OVER =====")
				fmt.Println()
				os.Exit(0)
			} else {
				time.Sleep(2 * time.Second)
			}
		}

		utils.ClearWindow()
	}

	fmt.Println("Congratulations, you used every word in the dictionary.")
	os.Exit(0)
}
