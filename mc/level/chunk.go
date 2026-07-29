package level

import (
	"io"
	"log/slog"
)

type Chunk struct {
	Sections []Section
	Version  string
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
		err = SectionEncodePacket(w, s)
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
