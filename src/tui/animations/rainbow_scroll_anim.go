package animations

import (
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

type RainbowScrollAnim struct {
	BaseAnim
	Offset 			int
	TotalFrames		int
	Colors			[]lipgloss.Style
}

func NewRainbowScrollAnim(
	target EffectTarget,
	total_frames int,
	loop bool,
	colors []lipgloss.Style,
) *RainbowScrollAnim {
	return &RainbowScrollAnim {
		BaseAnim: BaseAnim {
			FrameInterval:	time.Second / 20,
			PrevFrame:		time.Now(),
			Frame:			0,
			Loop:			loop,
			Active:			false,
			Target:			target,
		},
		Offset: 			0,
		TotalFrames: 		total_frames,
		Colors: 			colors,
	}
}

func (a *RainbowScrollAnim) Update(now time.Time) {
	if !a.AdvanceFrame(now) {
		return
	}

	a.Offset = (a.Offset - 1 + len(a.Colors)) % len(a.Colors)
	if !a.Loop && a.Frame >= a.TotalFrames {
		a.Active = false
	}
}

func (a *RainbowScrollAnim) Init() {
	a.BaseAnim.Init()
	a.Offset = 0
}

func (a *RainbowScrollAnim) ApplyEffect(text string) string {
	var out strings.Builder

	i := 0
	for _, c := range text {
		if string(c) == " " {
			out.WriteString(string(c))
			continue
		}
		style := a.Colors[(i + a.Offset) % len(a.Colors)]
		out.WriteString(style.Bold(true).Render(string(c)))
		i++
	}

	return out.String()
}
