package entity

// Abilities mirrors net.minecraft.world.entity.player.Abilities.
type Abilities struct {
	Invulnerable bool
	Flying       bool
	Mayfly       bool
	Instabuild   bool
	FlyingSpeed  float32
	WalkingSpeed float32
}

// DefaultAbilities mirrors Abilities.createDefault (survival player).
func DefaultAbilities() (ret Abilities) {
	return Abilities{
		Invulnerable: false,
		Flying:       false,
		Mayfly:       false,
		Instabuild:   false,
		FlyingSpeed:  0.05,
		WalkingSpeed: 0.1,
	}
}
