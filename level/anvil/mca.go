package anvil

import (
	"compress/gzip"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/admin-else/strom/nbt"
)

type CompressionStrategy uint8

const (
	ChunksWidthRegionFile = 32
	ChunksRegionFile      = ChunksWidthRegionFile * ChunksWidthRegionFile
)

const (
	_ CompressionStrategy = iota
	CompressionStrategyGZIP
	CompressionStrategyZLIB
	CompressionStrategyNONE
	CompressionStrategyLZ4
)

func (c CompressionStrategy) Wrap(r io.ReadCloser) (rNew io.ReadCloser, shouldClose bool, err error) {
	switch c {
	case CompressionStrategyGZIP:
		shouldClose = true
		rNew, err = gzip.NewReader(r)
		return
	case CompressionStrategyZLIB:
		shouldClose = true
		rNew, err = zlib.NewReader(r)
		return
	case CompressionStrategyNONE:
		rNew = r
		return
	case CompressionStrategyLZ4:
		err = fmt.Errorf("lz4 compression not supported")
		return
	default:
		err = fmt.Errorf("unknown compression strategy %d", c)
		return
	}
}

type ChunkHeader struct {
	Chunk, Time [ChunksRegionFile]int32
}

func ReadChunkHeader(r *os.File) (ch *ChunkHeader, err error) {
	var h ChunkHeader
	err = binary.Read(r, binary.BigEndian, &h)
	if err != nil {
		return
	}
	ch = &h
	return
}

func (c *ChunkHeader) ChunkAt(f *os.File, x, z int32) (n *nbt.Tag, err error) {
	off, size := chunkLocation(c.Chunk[chunkXZtoIndex(x, z)])
	_, err = f.Seek(int64(off), 0)
	if err != nil {
		return
	}
	_ = size // i dont really know why mc has this maybe to point out how much space is wasted?
	var realSize int32
	err = binary.Read(f, binary.BigEndian, &realSize)
	if err != nil {
		return
	}
	var compressionStrategy CompressionStrategy
	err = binary.Read(f, binary.BigEndian, &compressionStrategy)
	if err != nil {
		return
	}
	r, shouldClose, err := compressionStrategy.Wrap(f)
	if err != nil {
		return
	}
	if shouldClose {
		defer r.Close()
	}
	n = &nbt.Tag{}
	err = n.Decode(r)
	if err != nil {
		return
	}
	return
}

func chunkXZtoIndex(x, z int32) int32 {
	return (x % ChunksWidthRegionFile) + (z%ChunksWidthRegionFile)*ChunksWidthRegionFile
}

func chunkLocation(v int32) (offset, l int32) {
	return ((v >> 8) & 0xFFFFFF) * 4096, (v & 0xFF) * 4096
}
