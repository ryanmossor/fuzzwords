package figurethisout

import "fzwds/pkg/tui/animations"

type HelpKeymap struct {
	Key		string
	Value	string
}

type Size int
const (
	Undersized Size = iota
	Small
	Medium
	Large
)

type UIContext struct {
	DebugMap			map[string]string
	Size				Size

	ContainerWidth		int
	ContainerHeight		int

	ContentWidth		int
	ContentHeight		int

	ViewportWidth		int
	ViewportHeight		int

	FPS					int
	AnimManager			animations.AnimationManager
	InputRestricted		bool

	// footer msg?
}
