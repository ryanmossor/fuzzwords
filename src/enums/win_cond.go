package enums

type WinCondition int
const (
	Endless WinCondition = iota
	MaxLives
	Debug
)

func (w WinCondition) String() string {
	switch w {
	case Endless:
		return "endless"
	case MaxLives:
		return "max lives"
	case Debug:
		return "debug"
	default:
		return "unknown game mode"
	}
}
