package level

import "io"

type Chunk struct {
	Sections []Section
	version  string
}

func ReadChunkFromChunkPacketData(r io.Reader, version string, worldHeight int) (err error) {
	c := &Chunk{}
	c.version = version
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
