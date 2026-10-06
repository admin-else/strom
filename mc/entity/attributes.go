package entity

// Attribute mirrors net.minecraft.world.entity.ai.attributes.Attributes: the
// subset of attributes the movement code reads.
type Attribute string

// Attribute ids mirror the Attributes registry keys.
const (
	AttributeGRAVITY                   Attribute = "gravity"
	AttributeJUMP_STRENGTH             Attribute = "jump_strength"
	AttributeMOVEMENT_SPEED            Attribute = "movement_speed"
	AttributeFRICTION_MODIFIER         Attribute = "friction_modifier"
	AttributeAIR_DRAG_MODIFIER         Attribute = "air_drag_modifier"
	AttributeWATER_MOVEMENT_EFFICIENCY Attribute = "water_movement_efficiency"
	AttributeMOVEMENT_EFFICIENCY       Attribute = "movement_efficiency"
	AttributeSNEAKING_SPEED            Attribute = "sneaking_speed"
	AttributeFLYING_SPEED              Attribute = "flying_speed"
)

// DefaultAttributes mirrors the vanilla default attribute values a player has.
func DefaultAttributes() (ret map[Attribute]float64) {
	return map[Attribute]float64{
		AttributeGRAVITY:                   0.08,
		AttributeJUMP_STRENGTH:             0.42,
		AttributeMOVEMENT_SPEED:            0.1,
		AttributeFRICTION_MODIFIER:         0.0,
		AttributeAIR_DRAG_MODIFIER:         0.0,
		AttributeWATER_MOVEMENT_EFFICIENCY: 0.0,
		AttributeMOVEMENT_EFFICIENCY:       0.0,
		AttributeSNEAKING_SPEED:            0.3,
		AttributeFLYING_SPEED:              0.05,
	}
}

// AttributeValue mirrors LivingEntity.getAttributeValue(Attribute).
func (e *Entity) AttributeValue(attribute Attribute) (ret float64) {
	if e.attributes == nil {
		return 0.0
	}
	return e.attributes[attribute]
}

// SetAttributeValue mirrors AttributeInstance.setBaseValue for the movement slice.
func (e *Entity) SetAttributeValue(attribute Attribute, value float64) {
	if e.attributes == nil {
		e.attributes = DefaultAttributes()
	}
	e.attributes[attribute] = value
}
