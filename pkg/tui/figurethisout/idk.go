package figurethisout

import (
	"fzwds/pkg/game"
	"fzwds/pkg/tui/animations"
)

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

	// TODO: separate settings struct for path, schema etc?
	// pass settings only to game-related pages which need it, prefs only to other pages
	// maybe include prefs on uiContext?
	SettingsPath		string
	Settings			*game.Settings
	Schema				game.SettingsSchema

	// footer msg?
}
