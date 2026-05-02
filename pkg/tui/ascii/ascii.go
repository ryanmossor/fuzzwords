package ascii

// Returns height of tallest ASCII glyph in the provided string.
func MaxHeight(str string, glyphSet map[rune][]string) int {
	maxHeight := 0
	for _, c := range str {
		if len(CoderMini[c]) > maxHeight {
			maxHeight = len(CoderMini[c])
		}
	}
	return maxHeight
}

var CoderMini = map[rune][]string {
	'f': {
		"  ▄▄",
		" ██ ",
		"▀██▀",
		" ██ ",
		" ██ ",
	},
	'u': {
		"     ",
		"     ",
		"██ ██",
		"██ ██",
		"▀██▀█",
	},
	'z': {
		"     ",
		"     ",
		"▀▀▀██",
		"  ▄█▀",
		"▄██▄▄",
	},
	'w': {
		"       ",
		"       ",
		"██   ██",
		"██ █ ██",
		" ██▀██ ",
	},
	'o': {
		"     ",
		"     ",
		"▄███▄",
		"██ ██",
		"▀███▀",
	},
	'r': {
		"     ",
		"     ",
		"████▄",
		"██ ▀▀",
		"██   ",
	},
	'd': {
		"   ▄▄",
		"   ██",
		"▄████",
		"██ ██",
		"▀████",
	},
	's': {
		"     ",
		"     ",
		"▄█▀▀▀",
		"▀███▄",
		"▄▄▄█▀",
	},
}
