package utils

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"os/exec"
	"runtime"
)

func Average(arr []int) float64 {
	if len(arr) == 0 {
		return 0
	}

	var total float64
	for i := range arr {
		total += float64(arr[i])
	}
	avg := total / float64(len(arr))
	return math.Round((avg * 10)) / 10
}

// Preserves order, which is necessary for binary search
func Remove[T any](list []T, i int) []T {
    return append(list[:i], list[i+1:]...)
}

func GetLongestStr(list []string) string {
	var longest string
	for _, str := range list {
		if len(str) > len(longest) {
			longest = str
		}
	}
	return longest
}

func ArrToMap(lines []string) map[string]bool {
	var word_map = make(map[string]bool)
	for _, word := range lines {
		word_map[word] = true
	}
	return word_map
}

func ReadLines(path string, min_len int) ([]string, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var lines []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
		line := scanner.Text()
		if len(line) >= min_len {
			lines = append(lines, line)
		}
    }

    return lines, scanner.Err()
} 

func ClearWindow() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func FormatTime(seconds int) string {
	if seconds < 3600 {
		minutes := seconds / 60
		sec := seconds % 60
		return fmt.Sprintf("%d:%02d", minutes, sec)
	} else {
		hours := seconds / 3600
		remainingMinutes := (seconds % 3600) / 60
		remainingSeconds := seconds % 60
		return fmt.Sprintf("%d:%02d:%02d", hours, remainingMinutes, remainingSeconds)
	}
}
