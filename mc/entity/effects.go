package entity

// MobEffect mirrors the net.minecraft.world.effect.MobEffects ids the movement
// code reads.
type MobEffect string

// Effect ids mirror the MobEffects registry keys.
const (
	MobEffectJUMP_BOOST     MobEffect = "jump_boost"
	MobEffectSLOW_FALLING   MobEffect = "slow_falling"
	MobEffectLEVITATION     MobEffect = "levitation"
	MobEffectDOLPHINS_GRACE MobEffect = "dolphins_grace"
)

// EffectInstance mirrors net.minecraft.world.effect.MobEffectInstance for the
// movement slice. Duration -1 means infinite.
type EffectInstance struct {
	Duration  int
	Amplifier int
}

// HasEffect mirrors LivingEntity.hasEffect(Holder<MobEffect>).
func (e *Entity) HasEffect(effect MobEffect) (ret bool) {
	_, ok := e.effects[effect]
	return ok
}

// GetEffectAmplifier mirrors LivingEntity.getEffect(...).getAmplifier().
func (e *Entity) GetEffectAmplifier(effect MobEffect) (amplifier int, ok bool) {
	instance, ok := e.effects[effect]
	if !ok {
		return 0, false
	}
	return instance.Amplifier, true
}

// GetEffect mirrors LivingEntity.getEffect(Holder<MobEffect>).
func (e *Entity) GetEffect(effect MobEffect) (instance EffectInstance, ok bool) {
	instance, ok = e.effects[effect]
	return
}

// AddEffect mirrors LivingEntity.addEffect(MobEffectInstance).
func (e *Entity) AddEffect(effect MobEffect, amplifier int, duration int) {
	if e.effects == nil {
		e.effects = make(map[MobEffect]EffectInstance)
	}
	e.effects[effect] = EffectInstance{Duration: duration, Amplifier: amplifier}
}

// RemoveEffect mirrors LivingEntity.removeEffect.
func (e *Entity) RemoveEffect(effect MobEffect) { delete(e.effects, effect) }

// Effects returns a copy of the active effects.
func (e *Entity) Effects() (ret map[MobEffect]EffectInstance) {
	ret = make(map[MobEffect]EffectInstance, len(e.effects))
	for effect, instance := range e.effects {
		ret[effect] = instance
	}
	return
}

// TickEffects mirrors LivingEntity.tickEffects for the movement slice: finite
// durations count down and the effect is removed at zero.
func (e *Entity) TickEffects() {
	for effect, instance := range e.effects {
		if instance.Duration < 0 {
			continue
		}
		instance.Duration--
		if instance.Duration <= 0 {
			delete(e.effects, effect)
			continue
		}
		e.effects[effect] = instance
	}
}
