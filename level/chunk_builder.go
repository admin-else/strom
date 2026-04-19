package level

import (
	"bytes"
	"errors"

	"github.com/admin-else/strom/data"
	"github.com/admin-else/strom/util"
)

type ChunkBuilder struct {
	Biomes, Blocks [][]int32
	Bottom, Top    int32
	Version        string
}

func NewChunkBuilder(bottom, top int32, version string) *ChunkBuilder {
	return &ChunkBuilder{Bottom: bottom, Top: top, Version: version}
}

func NewOverworldChunkBuilder(version string) *ChunkBuilder {
	return NewChunkBuilder(-64, 320, version)
}

func (cb *ChunkBuilder) ColumnHeight() int32 {
	absBottom := cb.Bottom
	if absBottom < 0 {
		absBottom *= -1
	}
	return absBottom + cb.Top
}

func (cb *ChunkBuilder) FillAllById(block, biome int32) {
	numSections := cb.ColumnHeight() / ChunkWidth
	blockSection := util.MakeSingleValuedSlice(block, BlocksPerChunkSection)
	biomeSection := util.MakeSingleValuedSlice(biome, BiomesPerChunkSection)
	cb.Blocks = make([][]int32, numSections)
	cb.Biomes = make([][]int32, numSections)
	for i := range cb.Blocks {
		cb.Blocks[i] = blockSection
		cb.Biomes[i] = biomeSection
	}
	return
}

func (cb *ChunkBuilder) FillAll(biome, block string) (err error) {
	blockData, ok := data.LookupBlockByName(cb.Version, block)
	if !ok {
		err = data.UnknownBlockNameErr
		return
	}
	biomeData, ok := data.LookupBiomeByName(cb.Version, biome)
	if !ok {
		err = data.UnknownBiomeNameErr
		return
	}
	cb.FillAllById(blockData.DefaultState, biomeData.Id)
	return
}

var InvalidChunkPositionErr = errors.New("invalid chunk position")

func (cb *ChunkBuilder) SetBlockId(x, y, z, id int32) (err error) {
	if x < 0 || x >= ChunkWidth || z < 0 || z >= ChunkWidth || y < cb.Bottom || y >= cb.Top {
		err = InvalidChunkPositionErr
		return
	}

	y += cb.Bottom * -1
	section := util.FloorDiv(y, ChunkWidth)
	i := (y%ChunkWidth)*ChunkWidth*ChunkWidth + z*ChunkWidth + x
	cb.Blocks[section][i] = id
	return
}

func (cb *ChunkBuilder) SetBlock(x, y, z int32, block string) (err error) {
	blockData, ok := data.LookupBlockByName(cb.Version, block)
	if !ok {
		err = data.UnknownBlockNameErr
		return
	}
	return cb.SetBlockId(x, y, z, blockData.DefaultState)
}

func (cb *ChunkBuilder) Build() (data []byte, err error) {
	b := bytes.NewBuffer(nil)
	err = PackSections(cb.Blocks, cb.Biomes, cb.Version, b)
	if err != nil {
		return
	}
	data = b.Bytes()
	return
}
