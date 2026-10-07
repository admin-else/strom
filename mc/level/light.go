package level

const (
	LightNibbleBytes = 2048
)

// DataLayer holds a packed 4-bit light nibble array for one light section.
// When Data is nil the layer is homogeneously filled with DefaultValue, which
// mirrors vanilla's DataLayer(defaultValue) and avoids allocating an array for
// the common open-sky (15) and dark (0) sections.
type DataLayer struct {
	Data         []byte
	DefaultValue uint8
}

// Get returns the 4-bit value at index within the layer.
func (d *DataLayer) Get(index int) (v uint8) {
	if d.Data == nil {
		return d.DefaultValue
	}
	return nibbleAt(d.Data, index)
}

// GetXYZ returns the 4-bit value at the given section-local coordinates.
func (d *DataLayer) GetXYZ(x, y, z int) (v uint8) {
	return d.Get((y*ChunkWidth+z)*ChunkWidth + x)
}

// IsEmpty reports whether the layer is an all-zero homogeneous layer, matching
// vanilla DataLayer.isEmpty().
func (d *DataLayer) IsEmpty() bool {
	return d.Data == nil && d.DefaultValue == 0
}

// SectionLight holds the block and sky light layers for one light section.
type SectionLight struct {
	Sky   DataLayer
	Block DataLayer
}

// nibbleAt reads the 4-bit value at index within a packed nibble array.
func nibbleAt(data []byte, index int) (v uint8) {
	if index < 0 || index/2 >= len(data) {
		return
	}
	b := data[index/2]
	if index%2 == 0 {
		return b & 0xF
	}
	return b >> 4 & 0xF
}

// PackLight mirrors LightCoordsUtil.pack: the block nibble in bits 4-7 and the
// sky nibble in bits 20-23. World.LightAt supplies the two nibbles.
func PackLight(block, sky uint8) (ret uint32) {
	return uint32(block)<<4 | uint32(sky)<<20
}
