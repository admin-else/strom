package level

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"

	"github.com/admin-else/strom/proto_base"
)

const (
	BlocksPerChunkSection = 4096
	BiomesPerChunkSection = 64
)

func UnpackArrayPalette(r io.Reader, bitsPerEntry uint8, numberOfEntries int) (data []int32, err error) {
	var palletLen proto_base.VarInt
	palletLen, err = palletLen.Decode(r)
	if err != nil {
		return
	}
	var pallet []int32 // remember we can't use make because this is user-controlled data
	for range palletLen {
		var entry proto_base.VarInt
		entry, err = proto_base.VarInt(0).Decode(r)
		if err != nil {
			return
		}
		pallet = append(pallet, int32(entry))
	}
	data, err = UnpackLongData(r, bitsPerEntry, numberOfEntries)
	for i, b := range data {
		data[i] = pallet[b] // maybe check bounds?
	}

	return
}

func UnpackLongData(r io.Reader, bitsPerEntry uint8, numberOfEntries int) (data []int32, err error) {
	entriesPerLong := int(64 / bitsPerEntry)
	numberOfLongs := int(math.Ceil(float64(numberOfEntries) / float64(entriesPerLong)))
	var dataL []uint64
	for range numberOfLongs {
		var entry uint64
		err = binary.Read(r, binary.BigEndian, &entry)
		if err != nil {
			return
		}
		dataL = append(dataL, entry)
	}

	data = LongsToData(dataL, int32(numberOfEntries), int32(bitsPerEntry))
	return
}

func UnpackSingleValuePalette(r io.Reader, numberOfEntries int) (data []int32, err error) {
	entry, err := proto_base.VarInt(0).Decode(r)
	if err != nil {
		return
	}
	data = make([]int32, numberOfEntries)
	for i := range data {
		data[i] = int32(entry)
	}
	return
}

func LongsToData(data []uint64, n, bpe int32) (ret []int32) {
	mask := (uint64(1) << uint32(bpe)) - 1
	elementsPerLong := 64 / bpe

	ret = make([]int32, n)
	for i := range n {
		ret[i] = int32((data[i/elementsPerLong] >> ((i % elementsPerLong) * bpe)) & mask)
	}
	return
}

func UnpackBlockData(r io.Reader) (blocks [BlocksPerChunkSection]int32, err error) {
	var bitsPerEntry uint8
	err = binary.Read(r, binary.BigEndian, &bitsPerEntry)
	if err != nil {
		return
	}
	var blocksSlice []int32
	fmt.Println("blocks bpe", bitsPerEntry)
	switch bitsPerEntry {
	case 0:
		blocksSlice, err = UnpackSingleValuePalette(r, BlocksPerChunkSection)
	case 1, 2, 3, 4:
		blocksSlice, err = UnpackArrayPalette(r, bitsPerEntry, BlocksPerChunkSection)
	case 5, 6, 7, 8:
		err = errors.New("unimplemented")
	default:
		blocksSlice, err = UnpackLongData(r, bitsPerEntry, BlocksPerChunkSection)
	}
	if err != nil {
		return
	}
	blocks = [BlocksPerChunkSection]int32(blocksSlice)
	return
}

func UnpackBiomeData(r io.Reader) (biomes [BiomesPerChunkSection]int32, err error) {
	// net.minecraft.world.chunk.PaletteProvider#forBiomes
	var bitsPerEntry uint8
	err = binary.Read(r, binary.BigEndian, &bitsPerEntry)
	if err != nil {
		return
	}
	var biomesSlice []int32
	fmt.Println("biome bpe", bitsPerEntry)
	switch bitsPerEntry {
	case 0:
		biomesSlice, err = UnpackSingleValuePalette(r, BiomesPerChunkSection)
	case 1, 2, 3:
		biomesSlice, err = UnpackArrayPalette(r, bitsPerEntry, BiomesPerChunkSection)
	default:
		biomesSlice, err = UnpackLongData(r, bitsPerEntry, BiomesPerChunkSection)
	}
	if err != nil {
		return
	}
	biomes = [BiomesPerChunkSection]int32(biomesSlice)
	return
}

type ChunkSection struct {
	BlockData [BlocksPerChunkSection]int32
	BiomeData [BiomesPerChunkSection]int32
}

func UnpackSection(r io.Reader) (s ChunkSection, err error) {
	var blockCount int16
	err = binary.Read(r, binary.BigEndian, &blockCount)
	if err != nil {
		return
	}
	s.BlockData, err = UnpackBlockData(r)
	if err != nil {
		return
	}
	s.BiomeData, err = UnpackBiomeData(r)
	if err != nil {
		return
	}

	return
}
