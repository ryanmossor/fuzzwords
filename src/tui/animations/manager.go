package animations

import (
	"log/slog"
	"strings"
	"time"
)

type AnimationManager struct {
	animations map[EffectTarget]Animation
}

func InitAnimManager() AnimationManager {
	return AnimationManager {
		animations: make(map[EffectTarget]Animation),
	}
}

func (m *AnimationManager) Get(key EffectTarget) (Animation, bool) {
	anim, ok := m.animations[key]
	return anim, ok
}

func (m *AnimationManager) Register(anim Animation) {
	m.animations[anim.target()] = anim
	slog.Debug("Registered animation", "target", anim.target(), "animations", m.animations)
}

func (m *AnimationManager) InitAnimations(target_prefix EffectTarget) {
	for key, anim := range m.animations {
		if strings.HasPrefix(string(key), string(target_prefix)) {
			slog.Debug("Initializing animation for target",
				"targetPrefix", target_prefix,
				"anim", anim)
			anim.Init()
		}
	}
}

func (m *AnimationManager) DeactivateAnimations(target_prefix EffectTarget) {
	for key, anim := range m.animations {
		if strings.HasPrefix(string(key), string(target_prefix)) && anim.IsActive() {
			slog.Debug("Deactivating animations", "key", key)
			anim.deactivate()
		}
	}
}

func (m *AnimationManager) Update(now time.Time) {
	for _, a := range m.animations {
		a.Update(now)
	}
}

// Apply all active animations for target to provided input text.
// First return value is output string with all active animations applied.
// Second return value is bool indicating whether input string was changed.
func (m *AnimationManager) ApplyAnimations(target, text string) (string, bool) {
	out := text
	changed := false
    for key, a := range m.animations {
        if strings.HasPrefix(string(key), target) && a.IsActive() {
			slog.Debug("Applying text effect", "target", target, "text", text)
			out = a.ApplyEffect(out)
			changed = true
        }
    }
    return out, changed
}

type EffectTarget string
const (
	ExtraLife 			EffectTarget = "extra_life"
	StrikeCounter 		EffectTarget = "strike_counter"
	TitleLogo 			EffectTarget = "title_logo"
	ValidationMessage 	EffectTarget = "validation_message"
	GameOverWin		 	EffectTarget = "game_over_win"
)
