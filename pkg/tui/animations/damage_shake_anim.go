package animations

import (
	"fzwds/pkg/utils"
	"time"
)

type damageShakeAnim struct {
	baseAnim
	frames		[]int
}

func NewDamageShakeAnim(target effectTarget, max_amplitude int) *damageShakeAnim {
	return &damageShakeAnim {
		baseAnim: baseAnim {
			// Higher than global tick rate, guarantees frame advance every tick to ensure smooth motion
			frameInterval:	time.Second / 60,
			prevFrame:		time.Now(),
			frame:			0,
			loop:			false,
			active:			false,
			target:			target,
		},
		frames: 			utils.FillDescending(max_amplitude, 0),
	}
}

func (a *damageShakeAnim) update(now time.Time) {
	if !a.advanceFrame(now) {
		return
	}

	if a.frame >= len(a.frames) {
		a.active = false
	}
}

func (a *damageShakeAnim) applyAnimation(text string) string {
	padding := a.frames[a.frame]
	if padding % 2 == 0 {
		return utils.RightPad(text, padding)
	}
	return utils.LeftPad(text, padding)
}
