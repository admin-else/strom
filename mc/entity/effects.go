package entity

// MobEffect mirrors the net.minecraft.world.effect.MobEffects ids the movement
// code reads. Only the amplifier is tracked.
type MobEffect string

// Effect ids mirror the MobEffects registry keys.
const (
	MobEffectJUMP_BOOST     MobEffect = "jump_boost"
	MobEffectSLOW_FALLING   MobEffect = "slow_falling"
	MobEffectLEVITATION     MobEffect = "levitation"
	MobEffectDOLPHINS_GRACE MobEffect = "dolphins_grace"
)

// HasEffect mirrors LivingEntity.hasEffect(Holder<MobEffect>).
func (e *Entity) HasEffect(effect MobEffect) (ret bool) {
	_, ok := e.effects[effect]
	return ok
}

// GetEffectAmplifier mirrors LivingEntity.getEffect(...).getAmplifier(). ok is
// false when the effect is absent.
func (e *Entity) GetEffectAmplifier(effect MobEffect) (amplifier int, ok bool) {
	amplifier, ok = e.effects[effect]
	return
}

// AddEffect mirrors LivingEntity.addEffect for the movement slice.
func (e *Entity) AddEffect(effect MobEffect, amplifier int) {
	if e.effects == nil {
		e.effects = make(map[MobEffect]int)
	}
	e.effects[effect] = amplifier
}

// RemoveEffect mirrors LivingEntity.removeEffect.
func (e *Entity) RemoveEffect(effect MobEffect) { delete(e.effects, effect) }
