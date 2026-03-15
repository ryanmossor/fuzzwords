package animations

import (
	"fzwds/src/utils"
	"time"
)

type DamageShakeAnim struct {
	BaseAnim
	Frames		[]int
}

func NewDamageShakeAnim(target EffectTarget, max_amplitude int) *DamageShakeAnim {
	return &DamageShakeAnim {
		BaseAnim: BaseAnim {
			FrameInterval:	time.Second / 30,
			PrevFrame:		time.Now(),
			Frame:			0,
			Loop:			false,
			Active:			false,
			Target:			target,
		},
		Frames: 			utils.FillDescending(max_amplitude, 0),
	}
}

func (a *DamageShakeAnim) Update(now time.Time) {
	if !a.AdvanceFrame(now) {
		return
	}

	if a.Frame >= len(a.Frames) {
		a.Active = false
	}
}

func (a *DamageShakeAnim) ApplyEffect(text string) string {
	padding := a.Frames[a.Frame]
	if padding % 2 == 0 {
		return utils.RightPad(text, padding)
	}
	return utils.LeftPad(text, padding)
}
