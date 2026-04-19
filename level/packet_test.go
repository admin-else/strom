package level_test

import (
	"bytes"
	"slices"
	"testing"

	"github.com/admin-else/strom/data"
	"github.com/admin-else/strom/level"
	"github.com/admin-else/strom/util"
)

var RedMushroomBlock = util.MustOk(data.LookupBlockByName("1.21.11", "red_mushroom_block"))
var TheVoid = util.MustOk(data.LookupBiomeByName("1.21.11", "the_void"))

func TestPackSections(t *testing.T) {
	blocks := util.MakeSingleValuedSlice(RedMushroomBlock.DefaultState, level.BlocksPerChunkSection)
	biomes := util.MakeSingleValuedSlice(TheVoid.Id, level.BiomesPerChunkSection)

	b := bytes.NewBuffer(nil)
	err := level.PackSection(blocks, biomes, b)
	if err != nil {
		t.Fatal(err)
		return
	}
}

func TestDataLongsTest(t *testing.T) {
	var bpes = []int32{1, 2, 3, 4, 5, 10, 31}
	var testValues = [][]int32{{1}, {1, 2}, {1, 2, 3}, {1, 2, 3, 4}, {1, 2, 3, 4, 5}, {999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999, 999}}

	for _, bpe := range bpes {
		maxValue := int32((1 << bpe) - 1)
		t.Logf("BPE: %v, Maxvalue: %v", bpe, maxValue)
		for _, values := range testValues {
			t.Logf("Data: %v", values)

			skip := false
			for _, v := range values {
				if v >= maxValue {
					skip = true
					break
				}
			}
			if skip {
				t.Log("Skipping because of overflow")
				continue
			}

			longs := level.DataToLongs(values, bpe)

			t.Logf("Longs: %v", longs)

			dataNew := level.LongsToData(longs, int32(len(values)), bpe)

			if !slices.Equal(values, dataNew) {
				t.Error("Data is not equal")
			}
			t.Log("Data conversion successful")
		}
	}
}
