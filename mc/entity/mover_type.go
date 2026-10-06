package entity

// MoverType mirrors net.minecraft.world.entity.MoverType.
type MoverType int

const (
	MoverTypeSELF MoverType = iota
	MoverTypePLAYER
	MoverTypePISTON
	MoverTypeSHULKER_BOX
	MoverTypeSHULKER
)

// IsServerAndClientSimulated mirrors MoverType.isServerAndClientSimulated().
func (m MoverType) IsServerAndClientSimulated() (ret bool) { return m != MoverTypeSELF }
