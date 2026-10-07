package world

import (
	"testing"

	"github.com/admin-else/strom/mc/data"
	"github.com/admin-else/strom/mc/level"
)

func sectionWithBlock(t *testing.T, version string, blockName string) *level.Chunk {
	t.Helper()
	stateID, err := data.StateIdFromBlocKAndStateMap(version, blockName, nil)
	if err != nil {
		t.Fatalf("state id for %s: %v", blockName, err)
	}
	section := level.Section{}
	section.Blocks, err = level.MakeBlockFormat(version).FullWith(stateID)
	if err != nil {
		t.Fatalf("blocks: %v", err)
	}
	section.Biomes, _ = level.MakeBiomeFormat(version).FullWith(0)
	return &level.Chunk{Sections: []level.Section{section}}
}

// TestFluidAt checks water, lava and empty resolve from the block state.
func TestFluidAt(t *testing.T) {
	version := "26.4-snapshot-2"
	for _, tc := range []struct {
		block string
		kind  FluidKind
		id    string
	}{
		{"water", FluidWater, "minecraft:water"},
		{"lava", FluidLava, "minecraft:lava"},
		{"air", FluidEmpty, "minecraft:empty"},
	} {
		w := NewWorld(version, -64, 384)
		w.storeChunk(ChunkPos{0, 0}, sectionWithBlock(t, version, tc.block))
		got, err := w.FluidAt(1, -64+1, 1)
		if err != nil {
			t.Fatalf("%s: %v", tc.block, err)
		}
		if got.Kind != tc.kind {
			t.Errorf("%s: kind = %d (%s), want %d (%s)", tc.block, got.Kind, got.Kind, tc.kind, tc.kind)
		}
		if got.ID != tc.id {
			t.Errorf("%s: id = %q, want %q", tc.block, got.ID, tc.id)
		}
	}
}

// TestFluidKindString checks the exposed kind names.
func TestFluidKindString(t *testing.T) {
	if FluidEmpty.String() != "empty" || FluidWater.String() != "water" || FluidLava.String() != "lava" {
		t.Errorf("unexpected kind strings: %s %s %s", FluidEmpty, FluidWater, FluidLava)
	}
}
