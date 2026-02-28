package utils

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
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

func FilterWordList(words []string, min_len int) []string {
	var filtered []string
	for _, word := range words {
		if len(word) > min_len {
			filtered = append(filtered, word)
		}
	}
	return filtered
}

func FormatTime(seconds int) string {
	if seconds < 3600 {
		m := seconds / 60
		s := seconds % 60
		return fmt.Sprintf("%dm%02ds", m, s)
	} else {
		h := seconds / 3600
		m := (seconds % 3600) / 60
		s := seconds % 60
		return fmt.Sprintf("%dh%02dm%02ds", h, m, s)
	}
}

func CreateFuzzyPrompt(word string, prompt_len int) string {
	if len(word) <= prompt_len {
		return word
	}

	var prompt string
	rand_min := 0

	for i := prompt_len; i > 0; i-- {
		rand_max := len(word) - i
		rand_idx := rand.Intn(rand_max - rand_min + 1) + rand_min

		if i == prompt_len && rand_idx == rand_max {
			return prompt + word[rand_idx:]
		}

		prompt += string(word[rand_idx])
		rand_min = rand_idx + 1
	}

	return prompt
}

func IsFuzzyMatch(answer string, prompt string) bool {
    sub_idx := 0
	for i := range answer {
        if answer[i] == prompt[sub_idx] {
            sub_idx++
            if sub_idx == len(prompt) {
                return true
            }
        }
    }

    return false
}

func LeftPad(s string, n int) string {
	if n <= 0 {
		return s
	}
	return strings.Repeat(" ", n) + s
}

func RightPad(s string, n int) string {
	if n <= 0 {
		return s
	}
	return s + strings.Repeat(" ", n)
}

func ParseInt(n any) (int, bool) {
	switch v := n.(type) {
	case int:
		return v, true
	case int8:
		return int(v), true
	case int16:
		return int(v), true
	case int32:
		return int(v), true
	case int64:
		return int(v), true
	case float32:
		if float64(v) != math.Trunc(float64(v)) {
			return 0, false
		}
		return int(v), true
	case float64:
		if v != math.Trunc(v) {
			return 0, false
		}
		return int(v), true
	case string:
		vInt, err := strconv.Atoi(v)
		if err != nil {
			return 0, false
		}
		return vInt, true
	default:
		return 0, false
	}
}

func ValuesEqual(a, b any) bool {
	if int_a, ok_a := ParseInt(a); ok_a {
		if int_b, ok_b := ParseInt(b); ok_b {
			return int_a == int_b
		}
	}

	float_a, ok_a := a.(float64)
	float_b, ok_b := b.(float64)
	if ok_a && ok_b {
		return float_a == float_b
	}

	return a == b
}
