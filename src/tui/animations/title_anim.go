package animations

import (
	"fzwds/src/constants"
	"time"

	"github.com/charmbracelet/lipgloss"
)

type TitleScreenLogoAnim struct {
	baseAnim

	// Indicates currently active phase of animation
	Phase 			TitleScreenLogoPhase

	// Timestamp at which the current phase began
	phaseStart		time.Time

	// Index specifying how many glyphs of the full title logo have been "typed"
	TypedLetters	int

	// List of colors used to create rainbow scroll effect on full title logo
	Colors			[]lipgloss.Style

	// Current frame's starting index to use with []Colors for producing rainbow scroll effect
	ColorIdx		int
}

func NewTitleScreenLogoAnim(colors []lipgloss.Style) *TitleScreenLogoAnim {
	return &TitleScreenLogoAnim {
		baseAnim: baseAnim {
			frameInterval:	time.Second / 30,
			prevFrame:		time.Now(),
			frame:			0,
			loop:			true,
			active:			true,
			target:			TitleLogo,
		},
		Phase:				0,
		phaseStart:			time.Now(),
		TypedLetters:		0,
		ColorIdx: 			0,
		Colors: 			colors,
	}
}

func (a *TitleScreenLogoAnim) init() {
	a.baseAnim.init()
	a.frameInterval = time.Second * 5
	a.ColorIdx = 0
	a.Phase = 0
	a.phaseStart = time.Now()
	a.TypedLetters = 0
}

type TitleScreenLogoPhase int
const (
	AbbreviatedTitlePhase TitleScreenLogoPhase = iota
	TypingFullTitlePhase
	FullTitlePausePhase
	FullTitleRainbowScrollPhase
	TitleResetPhase
)

func (a *TitleScreenLogoAnim) update(now time.Time) {
	if !a.advanceFrame(now) {
		return
	}

	switch a.Phase {
	case AbbreviatedTitlePhase:
		// Wait 5 seconds on abbreviated logo
		if now.After(a.phaseStart.Add(5 * time.Second)) {
			a.nextPhase(now, time.Millisecond * 250)
		}
	case TypingFullTitlePhase:
		// "Typing" effect; display additional char of full title every 250ms by incrementing TypedLetters
		if a.TypedLetters >= len(constants.FULL_GAME_TITLE) {
			a.nextPhase(now, time.Millisecond * 1500)
		} else {
			a.TypedLetters++
		}
	case FullTitlePausePhase:
		// Wait 1.5s on fully typed logo
		if now.After(a.phaseStart.Add(time.Millisecond * 1500)) {
			a.nextPhase(now, time.Second / 12)
		}
	case FullTitleRainbowScrollPhase:
		// Apply rainbow scroll effect to full logo for 10s
		a.ColorIdx = (a.ColorIdx - 1 + len(a.Colors)) % len(a.Colors)
		if now.After(a.phaseStart.Add(10 * time.Second)) {
			a.nextPhase(now, time.Millisecond * 750)
		}
	case TitleResetPhase:
		// Reset to first phase to repeat animation
		if now.After(a.phaseStart.Add(a.frameInterval)) {
			a.init()
		}
	}
}

func (a *TitleScreenLogoAnim) nextPhase(now time.Time, frame_interval time.Duration) {
	a.frameInterval = frame_interval
	a.frame = 0
	a.Phase++
	a.phaseStart = now
	a.prevFrame = now
}

func (a *TitleScreenLogoAnim) applyAnimation(text string) string {
	// Because title screen anim is more complicated, effects/coloring are delegated to title screen view
	return text
}
