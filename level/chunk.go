package level

import (
	"encoding/binary"
	"io"
	"math"

	"github.com/admin-else/strom/data"
	"github.com/admin-else/strom/util"
)

type BlockState struct {
	Id         string
	Properties map[string]string
}

const (
	ChunkWidth            = 16
	ChunkColumns          = ChunkWidth * ChunkWidth
	BlocksPerChunkSection = ChunkWidth * ChunkWidth * ChunkWidth
	ChunkBiomesWidth      = 4
	BiomesPerChunkSection = ChunkBiomesWidth * ChunkBiomesWidth * ChunkBiomesWidth
)

type Section struct {
	Blocks *Storage[*BlockState]
	Biomes *Storage[string]
}

func PossibleBpesBlocks(version string) (ret []uint8) {
	maxBpe := util.BpeByNum(float64(len(data.BlocksForVersion(version))))
	ret = []uint8{0, 4, 5, 6, 7, 8, maxBpe}
	return
}

func PossibleBpesBiomes(version string) (ret []uint8) {
	maxBpe := uint8(math.Ceil(math.Log2(float64(len(data.BiomesForVersion(version))))))
	ret = []uint8{0, 1, 2, 3, maxBpe}
	return
}

func SectionFromPacket(r io.Reader, version string) (s *Section, err error) {
	s = &Section{}
	s.Blocks = NewStorage[*BlockState](BlocksPerChunkSection, PossibleBpesBlocks(version))
	s.Biomes = NewStorage[string](BiomesPerChunkSection, PossibleBpesBiomes(version))

	var bitsPerEntryBlocks uint8
	err = binary.Read(r, binary.BigEndian, &bitsPerEntryBlocks)
	if err != nil {
		return
	}

	return
}
