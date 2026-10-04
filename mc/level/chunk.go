package level

import (
	"io"
	"log/slog"
)

type Chunk struct {
	Sections      []Section
	Light         []SectionLight
	Heightmaps    []Heightmap
	BlockEntities []BlockEntity
	Version       string
}

// BlockEntity is the raw block entity record carried alongside chunk data.
type BlockEntity struct {
	LocalX  uint8
	LocalZ  uint8
	Y       int16
	Type    int32
	NbtData any
}

// LightAt returns the block and sky light nibble values at the given section-local
// coordinates. section is a chunk section index; bx, by, bz are in [0, 16).
//
// The light array covers sectionsCount+2 light sections: vanilla's
// LevelLightEngine spans minSectionY-1 (one below) through maxSectionY+1 (one
// above), so chunk section j maps to light index j+1.
func (c *Chunk) LightAt(section int, bx, by, bz int) (block, sky uint8) {
	lightIndex := section + 1
	if lightIndex < 0 || lightIndex >= len(c.Light) {
		return
	}
	index := (by*ChunkWidth+bz)*ChunkWidth + bx
	block = c.Light[lightIndex].Block.Get(index)
	sky = c.Light[lightIndex].Sky.Get(index)
	return
}

// HeightAt returns the height value for the given kind at column (x, z), where
// x and z are chunk-local in [0, 16). ok is false when the kind is absent.
func (c *Chunk) HeightAt(kind string, x, z int) (h int32, ok bool) {
	for _, hm := range c.Heightmaps {
		if hm.Type == kind {
			return heightAt(hm.Data, x, z)
		}
	}
	return
}

// ReadChunkFromChunkPacketData decodes chunk packet data from a reader into a Chunk.
func ReadChunkFromChunkPacketData(r io.Reader, version string, worldHeight int) (c *Chunk, err error) {
	c = &Chunk{}
	c.Version = version
	nSections := worldHeight / ChunkWidth

	c.Sections = make([]Section, nSections)
	for i := range c.Sections {
		c.Sections[i], err = SectionDecodePacket(r, version)
		if err != nil {
			return
		}
	}
	return
}

func (c *Chunk) WriteChunkData(w io.Writer) (err error) {
	for _, s := range c.Sections {
		err = SectionEncodePacketVersion(w, s, c.Version)
		if err != nil {
			return
		}
	}
	return
}

func (c *Chunk) Equals(other *Chunk) bool {
	if c.Version != other.Version {
		slog.Debug("Chunk versions don't match: ", "version 1", c.Version, "version2", other.Version)
		return false
	}
	if len(c.Sections) != len(other.Sections) {
		slog.Debug("Chunk section len does not match: ", "chunk1 len", len(c.Sections), "chunk 2 len", len(other.Sections))
		return false
	}
	for i, s := range c.Sections {
		if !SectionEquals(s, other.Sections[i]) {
			return false
		}
	}
	return true
}
