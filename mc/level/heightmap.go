package level

// HeightmapTypes maps the network heightmap kind to its wire type id.
var HeightmapTypes = map[string]int32{
	"world_surface_wg":       0,
	"world_surface":          1,
	"ocean_floor_wg":         2,
	"ocean_floor":            3,
	"motion_blocking":        4,
	"motion_blocking_no_leaves": 5,
}

// Heightmap stores the raw packed height values for one heightmap kind.
type Heightmap struct {
	Type string
	Data []int64
}

// heightAt unpacks the 9-bit height for column (x, z) from a packed heightmap.
func heightAt(data []int64, x, z int) (h int32, ok bool) {
	index := z*ChunkWidth + x
	bitsPer := 9
	valuesPerLong := 64 / bitsPer
	total := len(data) * valuesPerLong
	if index < 0 || index >= total {
		return
	}
	longIndex := index / valuesPerLong
	offset := uint(index%valuesPerLong) * uint(bitsPer)
	if longIndex >= len(data) {
		return
	}
	h = int32(data[longIndex] >> offset & 0x1FF)
	ok = true
	return
}
