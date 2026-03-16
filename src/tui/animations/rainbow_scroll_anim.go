package animations

import (
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

type rainbowScrollAnim struct {
	baseAnim
	offset 			int
	totalFrames		int
	colors			[]lipgloss.Style
}

func NewRainbowScrollAnim(
	target effectTarget,
	total_frames int,
	loop bool,
	colors []lipgloss.Style,
) *rainbowScrollAnim {
	return &rainbowScrollAnim {
		baseAnim: baseAnim {
			frameInterval:	time.Second / 20,
			prevFrame:		time.Now(),
			frame:			0,
			loop:			loop,
			active:			false,
			target:			target,
		},
		offset: 			0,
		totalFrames: 		total_frames,
		colors: 			colors,
	}
}

func (a *rainbowScrollAnim) update(now time.Time) {
	if !a.advanceFrame(now) {
		return
	}

	a.offset = (a.offset - 1 + len(a.colors)) % len(a.colors)
	if !a.loop && a.frame >= a.totalFrames {
		a.active = false
	}
}

func (a *rainbowScrollAnim) init() {
	a.baseAnim.init()
	a.offset = 0
}

func (a *rainbowScrollAnim) applyAnimation(text string) string {
	var out strings.Builder

	i := 0
	for _, c := range text {
		if string(c) == " " {
			out.WriteString(string(c))
			continue
		}
		style := a.colors[(i + a.offset) % len(a.colors)]
		out.WriteString(style.Bold(true).Render(string(c)))
		i++
	}

	return out.String()
}
