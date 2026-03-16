package animations

import (
	"log/slog"
	"strings"
	"time"
)

type AnimationManager struct {
	animations map[effectTarget]animation
}

func NewAnimationManager() AnimationManager {
	return AnimationManager {
		animations: make(map[effectTarget]animation),
	}
}

func (m *AnimationManager) Get(key effectTarget) (animation, bool) {
	anim, ok := m.animations[key]
	return anim, ok
}

func (m *AnimationManager) Register(anims ...animation) {
	for _, a := range anims {
		m.animations[a.getTarget()] = a
		slog.Debug("Registered animation", "target", a.getTarget())
	}
}

func (m *AnimationManager) InitAnimations(target_prefix effectTarget) {
	for key, a := range m.animations {
		if strings.HasPrefix(string(key), string(target_prefix)) {
			slog.Debug("Initializing animation for target", "targetPrefix", target_prefix, "anim", a)
			a.init()
		}
	}
}

func (m *AnimationManager) DeactivateAnimations(target_prefix effectTarget) {
	for key, a := range m.animations {
		if strings.HasPrefix(string(key), string(target_prefix)) && a.isActive() {
			slog.Debug("Deactivating animations", "key", key)
			a.deactivate()
		}
	}
}

func (m *AnimationManager) Update(now time.Time) {
	for _, a := range m.animations {
		a.update(now)
	}
}

// Apply all active animations for target to provided input text.
// First return value is output string with all active animations applied.
// Second return value is bool indicating whether input string was changed.
func (m *AnimationManager) ApplyAnimations(target, text string, animations_enabled bool) (string, bool) {
	if !animations_enabled {
		return text, false
	}

	out := text
	changed := false
    for key, a := range m.animations {
        if strings.HasPrefix(string(key), target) && a.isActive() {
			out = a.applyAnimation(out)
			changed = true
        }
    }
    return out, changed
}

type effectTarget string
const (
	ExtraLife 			effectTarget = "extra_life"
	GameOverWin		 	effectTarget = "game_over_win"
	StrikeCounter 		effectTarget = "strike_counter"
	TitleLogo 			effectTarget = "title_logo"
	ValidationMessage 	effectTarget = "validation_message"
)
