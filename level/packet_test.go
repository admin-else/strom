package level_test

import (
	"bytes"
	"math"
	"slices"
	"testing"

	"github.com/admin-else/strom/data"
	"github.com/admin-else/strom/level"
	"github.com/admin-else/strom/util"
)

var TestVersion = "1.21.11"
var RedMushroomBlock = util.MustOk(data.LookupBlockByName(TestVersion, "red_mushroom_block"))
var TheVoid = util.MustOk(data.LookupBiomeByName(TestVersion, "the_void"))

func TestPackSections(t *testing.T) {
	blocks := util.MakeSingleValuedSlice(RedMushroomBlock.DefaultState, level.BlocksPerChunkSection)
	biomes := util.MakeSingleValuedSlice(TheVoid.Id, level.BiomesPerChunkSection)

	b := bytes.NewBuffer(nil)
	err := level.PackSection(blocks, biomes, TestVersion, b)
	if err != nil {
		t.Fatal(err)
		return
	}
}

func FuzzPackSection(f *testing.F) {
	corpusBlocks := util.MakeSingleValuedSlice(int32(1), level.BlocksPerChunkSection)
	corpusBiomes := util.MakeSingleValuedSlice(int32(1), level.BiomesPerChunkSection)
	b := bytes.NewBuffer(nil)
	_ = level.PackSection(corpusBlocks, corpusBiomes, TestVersion, b)
	f.Add(b.Bytes())

	f.Fuzz(func(t *testing.T, raw []byte) {
		r := bytes.NewReader(raw)
		t.Logf("Raw: %v", raw)
		blocks, biomes, err := level.UnpackSection(r, TestVersion)
		if err != nil {
			t.Skipf("unpacking failed: %v", err)
			return
		}
		t.Logf("Blocks: %#v, Biomes: %#v", blocks, biomes)

		buf := bytes.NewBuffer(nil)
		err = level.PackSection(blocks, biomes, TestVersion, buf)
		if err != nil {
			t.Fatal(err)
			return
		}
		t.Logf("Packed: %v", buf.Bytes())
		newBlocks, newBimoes, err := level.UnpackSection(buf, TestVersion)
		if err != nil {
			t.Fatal(err)
			return
		}

		if !slices.Equal(blocks, newBlocks) {
			t.Error("Blocks are not equal")
		}
		if !slices.Equal(biomes, newBimoes) {
			t.Error("Biomes are not equal")
		}
	})
}

func TestDataLongsTest(t *testing.T) {
	var bpes = []uint8{1, 2, 3, 4, 5, 10, 31}
	var testValues = [][]int32{{1}, {1, 2}, {1, 2, 3}, {1, 2, 3, 4}, {1, 2, 3, 4, 5}, {999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999}}

	for _, bpe := range bpes {
		t.Logf("BPE: %v", bpe)
		for _, values := range testValues {
			t.Logf("Data: %v", values)

			longs, err := level.DataToLongs(values, bpe)
			if err != nil {
				t.Logf("Error: %v", err)
				continue
			}
			t.Logf("Longs: %v", longs)

			dataNew := level.LongsToData(longs, int32(len(values)), bpe)
			t.Logf("Data new: %v", dataNew)
			if !slices.Equal(values, dataNew) {
				t.Error("Data is not equal")
			} else {
				t.Log("Data conversion successful")
			}
		}
	}
}

func TestWeirdBiomesData(t *testing.T) {
	biomes := []int32{48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 67, 66, 65, 57, 56, 55, 50, 49}

	var b = bytes.NewBuffer(nil)

	unique := util.GetUniqueSlice(biomes)
	bpe := uint8(math.Ceil(math.Log2(float64(len(unique)))))
	err := level.PackLongData(biomes, bpe, b)
	_ = err
}
