package enums

type PromptMode int
const (
	Fuzzy PromptMode = iota
	Classic
)

func (m PromptMode) String() string {
	switch m {
	case Fuzzy:
		return "Fuzzy"
	case Classic:
		return "Classic"
	default:
		return "Unknown game mode"
	}
}
