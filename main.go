package main

import (
	"fmt"
	"fzw/src/tui"
	"log/slog"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// func main() {
//     cfg := game.InitializeSettings()
//     player := game.InitializePlayer(&cfg)

//     word_list, err := utils.ReadLines("./wordlist.txt")
//     if err != nil {
//         fmt.Fprintf(os.Stderr, "error: %v\n", err)
//         os.Exit(1)
//     }

//     words := game.WordLists{
//         FULL_MAP: utils.ArrToMap(word_list),
//         Available: word_list,
//         Used: make(map[string]bool),
//     }

//     reader := bufio.NewReader(os.Stdin)
//     input_chan := make(chan string)
//     stop_input := make(chan bool)
//     var ignore_input bool

//     // Goroutine to continuously read input
//     go func() {
//         for {
//             select {
//             case <-stop_input: // Stop reading input when receiving from stop_input
//                 return
//             default:
//                 input, _ := reader.ReadString('\n')
//                 if !ignore_input { // Only send to input_chan if ignore_input is false
//                     input_chan <- strings.TrimSpace(input)
//                 }
//             }
//         }
//     }()

//     utils.ClearWindow()

// 	for len(words.Available) > 0 {
// 		turn := game.NewTurn(words.Available, cfg)

// 		fmt.Fprintf(os.Stdin, "[ Health: %s ]\n", player.HealthDisplay)
// 		fmt.Println("Unused letters:", player.LettersRemaining)
// 		fmt.Println()

// 		turn_loop:
// 		for turn.Strikes < cfg.PromptStrikesMax {
// 			fmt.Fprintf(os.Stdin, "[ Strikes: %d / %d ]\n", turn.Strikes, cfg.PromptStrikesMax)

// 			fmt.Println("Prompt:", strings.ToUpper(turn.Prompt))
// 			ignore_input = false
// 			fmt.Print("Answer: ")
// 			answer := <-input_chan
// 			ignore_input = true // ignore further Enter presses until next iteration of turn loop

// 			turn.Answer = strings.ToLower(strings.TrimSpace(answer))
// 			turn.ValidateAnswer(&words, cfg)
// 			fmt.Println(turn.Msg)

// 			if turn.IsValid {
// 				player.HandleCorrectAnswer(turn.Answer)
// 				time.Sleep(750 * time.Millisecond)
// 				break turn_loop
// 			} else {
// 				turn.Strikes += 1
// 				fmt.Println()
// 			}

// 			if cfg.WinCondition == enums.MaxLives && player.HealthCurrent == cfg.HealthMax {
// 				fmt.Println("Max lives achieved -- you win!")
// 				close(stop_input) // Stop input goroutine on game over
// 				os.Exit(0)
// 			}

// 			if turn.Strikes == cfg.PromptStrikesMax {
// 				fmt.Println("Prompt failed. Possible answer:", turn.SourceWord)
// 				player.HandleFailedTurn()

// 				if player.HealthCurrent == 0 {
// 					fmt.Println()
// 					fmt.Println("===== GAME OVER =====")
// 					fmt.Println()
// 					fmt.Println(player.Stats.GenerateFinalStats())
// 					close(stop_input)
// 					os.Exit(0)
// 				} else {
// 					time.Sleep(2 * time.Second)
// 				}
// 			}
// 		}

// 		utils.ClearWindow()
// 	}

// 	fmt.Println("Congratulations, you used every word in the dictionary.")
// 	close(stop_input)
// 	os.Exit(0)
// }

// TUI
func main() {
	cache_dir, err := os.UserCacheDir()
	if err != nil {
        // panic("Cache dir not found" + err.Error())
		cache_dir = os.TempDir()
	}
	path := filepath.Join(cache_dir, "fuzzwords")
	os.MkdirAll(path, os.ModePerm)

	log_file, err := os.OpenFile(filepath.Join(path, "log.json"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        panic("Failed to open log file: " + err.Error())
    }
	defer log_file.Close()

    fileHandler := slog.NewJSONHandler(log_file, &slog.HandlerOptions{ Level: slog.LevelDebug })
    slog.SetDefault(slog.New(fileHandler))

	renderer := lipgloss.DefaultRenderer()
	menu := tui.NewModel(renderer)
	prog := tea.NewProgram(menu, tea.WithAltScreen())
	_, err = prog.Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		slog.Error("Error running program", "errMsg", err)
		os.Exit(1)
	}
}
 