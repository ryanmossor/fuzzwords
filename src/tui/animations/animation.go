package animations

import (
	"time"
)

type animation interface {
	init()

	activate()
	deactivate()
	isActive() bool

	update(time time.Time)
	applyAnimation(string) string

	getTarget() effectTarget
}

type baseAnim struct {
	frameInterval	time.Duration
	prevFrame     	time.Time
	frame			int

	loop 			bool // move to derived types?
	active 			bool

	target			effectTarget
}

// If time since PrevFrame is >= FrameInterval, increment Frame and update time of PrevFrame to now.
// Returns true if frame advanced, false otherwise.
func (a *baseAnim) advanceFrame(now time.Time) bool {
	if !a.active {
		return false
	}

	if now.Sub(a.prevFrame) >= a.frameInterval {
		a.frame++
		a.prevFrame = now
		return true
	}

	return false
}

func (a *baseAnim) init() {
	a.active = true
	a.frame = 0
	a.prevFrame = time.Time{}
}

func (a *baseAnim) isActive() bool {
	return a.active
}

func (a *baseAnim) activate() {
	a.active = true
}

func (a *baseAnim) deactivate() {
	a.active = false
}

func (a *baseAnim) getTarget() effectTarget {
	return a.target
}
