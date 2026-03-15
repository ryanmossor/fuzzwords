package animations

import (
	"log/slog"
	"strings"
	"time"
)

type AnimationManager struct {
	animations map[string]Animation
}

func InitAnimManager() AnimationManager {
	return AnimationManager {
		animations: make(map[string]Animation),
	}
}

func (m *AnimationManager) Get(key string) (Animation, bool) {
	anim, ok := m.animations[key]
	return anim, ok
}

func (m *AnimationManager) Register(key string, val Animation) {
	m.animations[key] = val
}

func (m *AnimationManager) InitAnimations(target_prefix EffectTarget) {
	for key, anim := range m.animations {
		if strings.HasPrefix(key, string(target_prefix)) {
			slog.Debug("Initializing animation for target",
				"targetPrefix", target_prefix,
				"anim", anim)
			anim.Init()
		}
	}
}

func (m *AnimationManager) Update(now time.Time) {
	for _, a := range m.animations {
		a.Update(now)
	}
}

type EffectTarget string
func (m *AnimationManager) ApplyAnimations(target, text string) (string, bool) {
	out := text
	changed := false
    for key, a := range m.animations {
        if strings.HasPrefix(key, target) && a.IsActive() {
			slog.Debug("Applying text effect", "target", target, "text", text)
			out = a.Effect(out)
			changed = true
        }
    }
    return out, changed
}

const (
	ExtraLife EffectTarget = "extra_life"
	TitleLogo EffectTarget = "title_logo"
)
