package animations

import (
	"log/slog"
	"strings"
	"time"
)

type AnimationManager struct {
	// TODO: may want to consider making this map[string][]Animation to register multiple
	// animations on a single key instead of needing a unique key (target) for each anim
	Animations map[string]Animation
}

func InitAnimManager() AnimationManager {
	return AnimationManager {
		Animations: make(map[string]Animation),
	}
}

func (m *AnimationManager) Get(key string) (Animation, bool) {
	anim, ok := m.Animations[key]
	return anim, ok
}

func (m *AnimationManager) Register(key string, val Animation) {
	m.Animations[key] = val
}

func (m *AnimationManager) InitAnimations(target_prefix EffectTarget) {
	for key, anim := range m.Animations {
		if strings.HasPrefix(key, string(target_prefix)) {
			slog.Debug("Initializing animation for target",
				"targetPrefix", target_prefix,
				"anim", anim)
			anim.Init()
		}
	}
}

func (m *AnimationManager) Update(now time.Time) {
	for _, a := range m.Animations {
		a.Update(now)
	}
}

type EffectTarget string
func (m *AnimationManager) EffectsFor(target string) []TextEffect {
    effects := []TextEffect{}
	slog.Debug("EffectsFor", "target", target)
    for key, a := range m.Animations {
        if strings.HasPrefix(key, target) && a.IsActive() {
            effects = append(effects, a.Effect())
        }
    }
    return effects
}

type TextEffect func(string) string
func ApplyTextEffects(text string, effects ...TextEffect) string {
	slog.Debug("Applying text effects", "text", text, "effectsLen", len(effects))
	out := text
	for _, e := range effects {
		out = e(out)
	}
	return out
}

const (
	ExtraLife EffectTarget = "extra_life"
)
