package world

import (
	"testing"

	"github.com/admin-else/strom/mc/level"
	"github.com/admin-else/strom/mc/registry"
)

// TestBiomeAt checks the biome palette id is resolved through the captured
// registry and joined with the minecraft-data precipitation and temperature.
func TestBiomeAt(t *testing.T) {
	store := registry.NewStore()
	biomes := registry.NewRegistry(BiomeRegistryName)
	biomes.Add("minecraft:badlands")
	biomes.Add("minecraft:desert")
	biomes.Add("minecraft:plains")
	store.Set(biomes)

	section := level.Section{}
	section.Biomes, _ = level.MakeBiomeFormat("26.4-snapshot-2").FullWith(2)

	w := NewWorld("26.4-snapshot-2", -64, 384)
	w.SetRegistries(store)
	w.storeChunk(ChunkPos{0, 0}, &level.Chunk{Sections: []level.Section{section}})

	got, err := w.BiomeAt(4, -64+8, 5)
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != "minecraft:plains" {
		t.Errorf("biome id = %q, want minecraft:plains", got.ID)
	}
	if !got.HasPrecipitation {
		t.Errorf("plains must have precipitation")
	}
	if got.Temperature != 0.8 {
		t.Errorf("plains temperature = %v, want 0.8", got.Temperature)
	}
}

// TestBiomeAtRequiresRegistry checks the error when no registry was captured.
func TestBiomeAtRequiresRegistry(t *testing.T) {
	section := level.Section{}
	section.Biomes, _ = level.MakeBiomeFormat("26.4-snapshot-2").FullWith(0)
	w := NewWorld("26.4-snapshot-2", -64, 384)
	w.storeChunk(ChunkPos{0, 0}, &level.Chunk{Sections: []level.Section{section}})
	if _, err := w.BiomeAt(0, -64+8, 0); err != BiomeRegistryAbsentErr {
		t.Errorf("err = %v, want BiomeRegistryAbsentErr", err)
	}
}
