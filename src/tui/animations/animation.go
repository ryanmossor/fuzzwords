package animations

import (
	"time"
)

type Animation interface {
	Init()
	Update(time time.Time)
	IsActive() bool
	ApplyEffect(string) string

	activate()
	deactivate()
	target() EffectTarget
}

type BaseAnim struct {
	FrameInterval	time.Duration
	PrevFrame     	time.Time
	Frame			int

	Loop 			bool // move to derived types?
	Active 			bool

	Target			EffectTarget
}

// If time since PrevFrame is >= FrameInterval, increment Frame and update time of PrevFrame to now.
// Returns true if frame advanced, false otherwise.
func (a *BaseAnim) AdvanceFrame(now time.Time) bool {
	if !a.Active {
		return false
	}

	if now.Sub(a.PrevFrame) >= a.FrameInterval {
		a.Frame++
		a.PrevFrame = now
		return true
	}

	return false
}

func (a *BaseAnim) Init() {
	a.Active = true
	a.Frame = 0
	a.PrevFrame = time.Time{}
}

func (a *BaseAnim) IsActive() bool {
	return a.Active
}

func (a *BaseAnim) activate() {
	a.Active = true
}

func (a *BaseAnim) deactivate() {
	a.Active = false
}

func (a *BaseAnim) target() EffectTarget {
	return a.Target
}
