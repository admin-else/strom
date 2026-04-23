package level

import (
	"encoding/binary"
	"errors"
	"io"
	"math"

	"github.com/admin-else/strom/data"
	"github.com/admin-else/strom/proto_base"
	"github.com/admin-else/strom/util"
)

const (
	ChunkWidth            = 16
	ChunkColumns          = ChunkWidth * ChunkWidth
	BlocksPerChunkSection = ChunkWidth * ChunkWidth * ChunkWidth
	BiomesPerChunkSection = 4 * 4 * 4
)

var InvalidPalletIndexErr = errors.New("invalid palette index")

func UnpackArrayPalette(r io.Reader, bitsPerEntry uint8, numberOfEntries int) (data []int32, err error) {
	palletLen, err := proto_base.DecodeVarInt(r)
	if err != nil {
		return
	}
	var pallet []int32 // remember we can't use make because this is user-controlled data
	for range palletLen {
		var entry int32
		entry, err = proto_base.DecodeVarInt(r)
		if err != nil {
			return
		}
		pallet = append(pallet, entry)
	}
	data, err = UnpackLongData(r, bitsPerEntry, numberOfEntries)

	for i, b := range data {
		if b < 0 || b >= palletLen {
			err = InvalidPalletIndexErr
			return
		}
		data[i] = pallet[b]
	}

	return
}

func PackArrayPallet(unique map[int32]struct{}, bpe uint8, data []int32, w io.Writer) (err error) {
	pallet := util.SetToSliceOrdered(unique) // this may improve performance in combination with zlib i dont know tho
	err = proto_base.EncodeVarInt(w, int32(len(pallet)))
	if err != nil {
		return
	}
	var palletMap = make(map[int32]int32)
	for i, entry := range pallet {
		palletMap[entry] = int32(i)
		err = proto_base.EncodeVarInt(w, entry)
		if err != nil {
			return
		}
	}
	palletedData := make([]int32, len(data))
	for i, entry := range data {
		palletedData[i] = palletMap[entry]
	}
	return PackLongData(palletedData, bpe, w)
}

var BadBitsPerEntryErr = errors.New("bad bits per entry")

func UnpackLongData(r io.Reader, bitsPerEntry uint8, numberOfEntries int) (data []int32, err error) {
	if bitsPerEntry < 1 || bitsPerEntry > 31 {
		err = BadBitsPerEntryErr
		return
	}

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

	data = LongsToData(dataL, int32(numberOfEntries), bitsPerEntry)
	return
}

func PackLongData(data []int32, bpe uint8, w io.Writer) (err error) {
	dataL, err := DataToLongs(data, bpe)
	if err != nil {
		return
	}
	for _, entry := range dataL {
		err = binary.Write(w, binary.BigEndian, entry)
	}
	return
}

func UnpackSingleValuePalette(r io.Reader, numberOfEntries int) (data []int32, err error) {
	entry, err := proto_base.DecodeVarInt(r)
	if err != nil {
		return
	}
	data = make([]int32, numberOfEntries)
	for i := range data {
		data[i] = entry
	}
	return
}

func LongsToData(data []uint64, n int32, bpe uint8) (ret []int32) {
	mask := (uint64(1) << uint32(bpe)) - 1
	elementsPerLong := int32(64 / bpe)

	ret = make([]int32, n)
	for i := range n {
		ret[i] = int32((data[i/elementsPerLong] >> ((i % elementsPerLong) * int32(bpe))) & mask)
	}
	return
}

var BpeTooSmallErr = errors.New("bits per entry too small")

func DataToLongs(data []int32, bpe uint8) (ret []uint64, err error) {
	maxValue := uint64((1 << bpe) - 1)

	elementsPerLong := 64 / int32(bpe)
	retLen := int32(math.Ceil(float64(len(data)) / float64(elementsPerLong)))
	ret = make([]uint64, retLen)
	for i := range retLen {
		for j := range elementsPerLong {
			dataI := i*elementsPerLong + j
			if dataI >= int32(len(data)) {
				break
			}
			toInsert := uint64(data[dataI])
			if toInsert > maxValue {
				err = BpeTooSmallErr
				return
			}
			ret[i] |= toInsert << (uint32(j) * uint32(bpe))
		}
	}
	return
}

func GetEntryFromLongs(data []uint64, i, bpe int32) (ret int32) {
	mask := (uint64(1) << uint32(bpe)) - 1
	elementsPerLong := 64 / bpe
	ret = int32((data[i/elementsPerLong] >> ((i % elementsPerLong) * bpe)) & mask)
	return
}

// UnpackBlockData minecraft does this at net.minecraft.world.chunk.PaletteProvider
func UnpackBlockData(r io.Reader) (blocks []int32, err error) {
	var bitsPerEntry uint8
	err = binary.Read(r, binary.BigEndian, &bitsPerEntry)
	if err != nil {
		return
	}
	switch bitsPerEntry {
	case 0:
		blocks, err = UnpackSingleValuePalette(r, BlocksPerChunkSection)
	case 1, 2, 3, 4:
		blocks, err = UnpackArrayPalette(r, 4, BlocksPerChunkSection)
	case 5, 6, 7, 8:
		blocks, err = UnpackArrayPalette(r, bitsPerEntry, BlocksPerChunkSection)
	default:
		blocks, err = UnpackLongData(r, bitsPerEntry, BlocksPerChunkSection)
	}
	if err != nil {
		return
	}
	return
}

func PackBlockData(blocks []int32, version string, w io.Writer) (err error) {
	unique := util.GetUniqueSlice(blocks)
	bpe := util.BpeByNum(float64(len(unique)))

	if bpe > 8 {
		bpe = util.BpeByNum(float64(len(data.BlocksForVersion(version))))
	}

	err = binary.Write(w, binary.BigEndian, bpe)
	if err != nil {
		return
	}

	switch bpe {
	case 0:
		err = proto_base.EncodeVarInt(w, blocks[0])
	case 1, 2, 3, 4:
		err = PackArrayPallet(unique, 4, blocks, w) // MOJANK
	case 5, 6, 7, 8:
		err = PackArrayPallet(unique, bpe, blocks, w)
	default:
		// matches net.minecraft.world.chunk.PaletteProvider
		err = PackLongData(blocks, bpe, w)
	}

	// This should also maybe get a maxblockid check but my fuzzer does not trigger a fail here so its probably fine
	// also that would eat alot of performace

	return
}

var BadBiomeIdErr = errors.New("bad biome id")

func UnpackBiomeData(r io.Reader, version string) (biomes []int32, err error) {
	// net.minecraft.world.chunk.PaletteProvider#forBiomes
	var bitsPerEntry uint8
	err = binary.Read(r, binary.BigEndian, &bitsPerEntry)
	if err != nil {
		return
	}
	switch bitsPerEntry {
	case 0:
		biomes, err = UnpackSingleValuePalette(r, BiomesPerChunkSection)
	case 1, 2, 3:
		biomes, err = UnpackArrayPalette(r, bitsPerEntry, BiomesPerChunkSection)
	default:
		biomes, err = UnpackLongData(r, bitsPerEntry, BiomesPerChunkSection)
	}

	maxBiome := int32(len(data.BiomesForVersion(version)))
	for _, b := range biomes {
		if b < 0 || b >= maxBiome {
			err = BadBiomeIdErr
			return
		}
	}

	if err != nil {
		return
	}
	return
}

func PackBiomeData(biomes []int32, version string, w io.Writer) (err error) {
	unique := util.GetUniqueSlice(biomes)
	bpe := uint8(math.Ceil(math.Log2(float64(len(unique)))))

	if bpe > 3 {
		bpe = uint8(math.Ceil(math.Log2(float64(len(data.BiomesForVersion(version))))))
	}

	err = binary.Write(w, binary.BigEndian, bpe)
	if err != nil {
		return
	}

	switch bpe {
	case 0:
		err = proto_base.EncodeVarInt(w, biomes[0])
	case 1, 2, 3:
		err = PackArrayPallet(unique, bpe, biomes, w)
	default:
		// matches net.minecraft.world.chunk.PaletteProvider
		err = PackLongData(biomes, bpe, w)
	}
	return
}

func UnpackSection(r io.Reader, version string) (blocks, biomes []int32, err error) {
	var blockCount int16
	err = binary.Read(r, binary.BigEndian, &blockCount)
	if err != nil {
		return
	}
	blocks, err = UnpackBlockData(r)
	if err != nil {
		return
	}
	biomes, err = UnpackBiomeData(r, version)
	return
}

func UnpackNSections(r io.Reader, n int) (blocks, biomes [][]int32, version string, err error) {
	for range n {
		var blocksSection, biomesSections []int32
		blocksSection, biomesSections, err = UnpackSection(r, version)
		if err != nil {
			return
		}
		blocks = append(blocks, blocksSection)
		biomes = append(biomes, biomesSections)
	}
	return
}

var LenMustMatchErr = errors.New("length must match")
var BiomeMust64LongErr = errors.New("biome section must be 64 long")
var BlockMust4096LongErr = errors.New("block section must be 4096 long")

func PackSection(blocks, biomes []int32, version string, w io.Writer) (err error) {
	if len(blocks) != BlocksPerChunkSection {
		return BlockMust4096LongErr
	}
	if len(biomes) != BiomesPerChunkSection {
		return BiomeMust64LongErr
	}
	numBlocks := BlocksPerChunkSection - int32(util.CountSlice(blocks, 0))

	err = binary.Write(w, binary.BigEndian, int16(numBlocks))
	if err != nil {
		return
	}
	err = PackBlockData(blocks, version, w)
	if err != nil {
		return
	}
	err = PackBiomeData(biomes, version, w)
	return
}

func PackSections(blocks, biomes [][]int32, version string, w io.Writer) (err error) {
	if len(blocks) != len(biomes) {
		return LenMustMatchErr
	}
	for i, blocksSection := range blocks {
		err = PackSection(blocksSection, biomes[i], version, w)
		if err != nil {
			return
		}
	}
	return
}
