package world

import (
	"testing"

	"github.com/admin-else/strom/mc/level"
)

// TestWorldHeightAtMinYOffset checks that HeightAt converts the raw packed
// heightmap value (relative to the world's minimum Y) back to an absolute Y.
func TestWorldHeightAtMinYOffset(t *testing.T) {
	data := make([]int64, 37)
	index := 3*level.ChunkWidth + 2 // local z=3, x=2
	data[index/7] = int64(100) << ((index % 7) * 9)

	w := NewWorld("26.4-snapshot-2", -64, 384)
	w.storeChunk(ChunkPos{0, 0}, &level.Chunk{
		Sections:   make([]level.Section, 24),
		Heightmaps: []level.Heightmap{{Type: "world_surface", Data: data}},
	})

	h, err := w.HeightAt("world_surface", 2, 3)
	if err != nil {
		t.Fatal(err)
	}
	if h != 100-64 {
		t.Errorf("HeightAt = %d, want %d", h, 100-64)
	}
}
