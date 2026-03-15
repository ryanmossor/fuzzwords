package animations

import (
	"fzwds/src/utils"
	"time"
)

type DamageShakeAnim struct {
	BaseAnim
	Frames		[]int
}

func (a *DamageShakeAnim) Update(now time.Time) {
	if !a.AdvanceFrame(now) {
		return
	}

	if a.Frame >= len(a.Frames) {
		a.Active = false
	}
}

func (a *DamageShakeAnim) Effect(text string) string {
	padding := a.Frames[a.Frame]
	if padding % 2 == 0 {
		return utils.RightPad(text, padding)
	}
	return utils.LeftPad(text, padding)
}
