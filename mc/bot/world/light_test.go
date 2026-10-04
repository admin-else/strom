package world

import (
	"testing"

	"github.com/admin-else/strom/mc/level"
)

func setBit(mask []byte, i int) []byte {
	if i/8 >= len(mask) {
		mask = append(mask, make([]byte, i/8-len(mask)+1)...)
	}
	mask[i/8] |= 1 << (uint(i) % 8)
	return mask
}

func filled(v uint8) []byte {
	b := make([]byte, level.LightNibbleBytes)
	for i := range b {
		b[i] = v<<4 | v
	}
	return b
}

// TestBuildSectionLight checks the ClientboundLightUpdatePacketData mapping:
// masks are bitsets over sectionsCount+2 light sections, arrays follow the set
// bits in increasing order, empty bits are all-zero layers and unset sky bits
// are open sky (15).
func TestBuildSectionLight(t *testing.T) {
	const sections = 3
	var skyMask, blockMask, emptySkyMask []byte
	skyMask = setBit(skyMask, 0)
	skyMask = setBit(skyMask, 2)
	emptySkyMask = setBit(emptySkyMask, 1)
	blockMask = setBit(blockMask, 2)

	skyLight := [][]byte{filled(7), filled(9)}
	blockLight := [][]byte{filled(3)}

	light := buildSectionLight(skyMask, blockMask, emptySkyMask, nil, skyLight, blockLight, sections)
	if len(light) != sections+2 {
		t.Fatalf("light sections = %d, want %d", len(light), sections+2)
	}
	if got := light[0].Sky.Get(0); got != 7 {
		t.Errorf("sky[0] = %d, want 7 (first array)", got)
	}
	if got := light[1].Sky.Get(0); got != 0 {
		t.Errorf("sky[1] = %d, want 0 (empty mask)", got)
	}
	if got := light[2].Sky.Get(0); got != 9 {
		t.Errorf("sky[2] = %d, want 9 (second array)", got)
	}
	if got := light[3].Sky.Get(0); got != 15 {
		t.Errorf("sky[3] = %d, want 15 (open sky)", got)
	}
	if got := light[4].Sky.Get(0); got != 15 {
		t.Errorf("sky[4] = %d, want 15 (open sky)", got)
	}
	if got := light[2].Block.Get(0); got != 3 {
		t.Errorf("block[2] = %d, want 3", got)
	}
	if got := light[0].Block.Get(0); got != 0 {
		t.Errorf("block[0] = %d, want 0", got)
	}
}

// TestUpdateLightMerges checks that a standalone light update only replaces the
// sections named by the masks and leaves the others intact.
func TestUpdateLightMerges(t *testing.T) {
	w := NewWorld("26.4-snapshot-2", -64, 48)
	chunk := &level.Chunk{
		Sections: make([]level.Section, 3),
		Light: []level.SectionLight{
			{Sky: level.DataLayer{DefaultValue: 15}},
			{Sky: level.DataLayer{DefaultValue: 15}},
			{Sky: level.DataLayer{DefaultValue: 15}},
			{Sky: level.DataLayer{DefaultValue: 15}},
			{Sky: level.DataLayer{DefaultValue: 15}},
		},
	}
	w.storeChunk(ChunkPos{0, 0}, chunk)

	blockMask := setBit(nil, 2)
	emptySkyMask := setBit(nil, 3)
	w.updateLight(ChunkPos{0, 0}, nil, blockMask, emptySkyMask, nil, nil, [][]byte{filled(4)})

	if got := chunk.Light[2].Block.Get(0); got != 4 {
		t.Errorf("block[2] = %d, want 4", got)
	}
	if got := chunk.Light[2].Sky.Get(0); got != 15 {
		t.Errorf("sky[2] = %d, want untouched 15", got)
	}
	if got := chunk.Light[3].Sky.Get(0); got != 0 {
		t.Errorf("sky[3] = %d, want 0 (empty mask)", got)
	}
	if got := chunk.Light[0].Sky.Get(0); got != 15 {
		t.Errorf("sky[0] = %d, want untouched 15", got)
	}
}

// TestChunkLightAtOffset verifies chunk section j reads light section j+1.
func TestChunkLightAtOffset(t *testing.T) {
	chunk := &level.Chunk{
		Sections: make([]level.Section, 3),
		Light:    make([]level.SectionLight, 5),
	}
	chunk.Light[0].Sky = level.DataLayer{DefaultValue: 1}
	chunk.Light[1].Sky = level.DataLayer{DefaultValue: 2}
	chunk.Light[2].Sky = level.DataLayer{DefaultValue: 3}
	chunk.Light[3].Sky = level.DataLayer{DefaultValue: 15}

	for section, want := range map[int]uint8{0: 2, 1: 3, 2: 15} {
		if _, sky := chunk.LightAt(section, 0, 0, 0); sky != want {
			t.Errorf("LightAt(section=%d) sky = %d, want %d", section, sky, want)
		}
	}
}
