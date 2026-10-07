package world

import (
	"errors"
	"strings"

	"github.com/admin-else/strom/mc/data"
	"github.com/admin-else/strom/mc/level"
)

// BiomeRegistryName is the dynamic registry the server sends the biome palette
// in during configuration.
const BiomeRegistryName = "minecraft:worldgen/biome"

var (
	BiomeRegistryAbsentErr = errors.New("biome registry not attached")
	BiomeNotInRegistryErr  = errors.New("biome id not in the captured registry")
	BiomeDataAbsentErr     = errors.New("biome not in the minecraft-data definitions")
)

// BiomeInfo is the raw biome state at a position: the resource id from the
// server registry plus the precipitation and temperature the client uses to
// decide rain vs snow and to tint the sky.
type BiomeInfo struct {
	NumericID        int32
	ID               string
	HasPrecipitation bool
	Temperature      float64
}

// BiomeAt returns the biome at the given world coordinates. The numeric id is
// read from the section's biome palette and resolved through the captured biome
// registry; precipitation and temperature come from minecraft-data.
func (w *World) BiomeAt(x, y, z int32) (ret BiomeInfo, err error) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	chunkX := floorDiv(x, level.ChunkWidth)
	chunkZ := floorDiv(z, level.ChunkWidth)
	chunk, ok := w.chunks[ChunkPos{chunkX, chunkZ}]
	if !ok {
		err = ChunkNotLoadedErr
		return
	}
	sectionIndex := (int(y) - w.minY) / level.ChunkWidth
	if sectionIndex < 0 || sectionIndex >= len(chunk.Sections) {
		err = OutOfBoundsErr
		return
	}
	lx, ly, lz := blockToLocal(x, y, z)
	numericID, getErr := chunk.Sections[sectionIndex].Biomes.Get(level.BiomeIndex(w.version, lx, ly, lz))
	if getErr != nil {
		err = getErr
		return
	}
	ret.NumericID = numericID
	if w.registries == nil {
		err = BiomeRegistryAbsentErr
		return
	}
	resourceID, found := w.registries.ResourceID(BiomeRegistryName, int(numericID))
	if !found {
		err = BiomeNotInRegistryErr
		return
	}
	ret.ID = resourceID
	definition, found := data.LookupBiomeByName(w.version, strings.TrimPrefix(resourceID, "minecraft:"))
	if !found {
		err = BiomeDataAbsentErr
		return
	}
	ret.HasPrecipitation = definition.HasPrecipitation
	ret.Temperature = definition.Temperature
	return
}
