package level_test

import (
	"fmt"
	"testing"

	"github.com/admin-else/strom/level"
	"github.com/admin-else/strom/util"
)

var possibleBpe = []uint8{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}

func FuzzStorage(f *testing.F) {
	f.Add([]byte{4, 0, 1})
	f.Add([]byte{4, 0, 1, 0, 2, 0, 1})
	f.Add([]byte{8, 0, 10, 1, 20, 2, 30, 3, 40, 4, 50})
	f.Add([]byte{4, 0, 5, 1, 10, 2, 15, 3, 20, 0, 25})
	f.Add([]byte{2, 0, 1, 1, 2, 0, 3, 1, 4})
	f.Add([]byte{8, 0, 42, 1, 42, 2, 42, 3, 42})

	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) < 1 {
			t.Skip()
		}

		length := int(data[0])%64 + 1
		s := level.NewStorage[int32](length, possibleBpe)

		ref := make([]int32, length)
		refSet := make([]bool, length)

		for i := 1; i+1 < len(data); i += 2 {
			idx := int(data[i]) % length
			val := int32(data[i+1])

			err := s.Set(idx, val)
			if err == level.CannotGrowPalletErr {
				continue
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			ref[idx] = val
			refSet[idx] = true
		}

		for i := range length {
			if !refSet[i] {
				continue
			}
			got := s.Get(i)
			if got != ref[i] {
				t.Errorf("Get(%d) = %d, want %d", i, got, ref[i])
			}
		}
	})
}

func TestSetOverwrite(t *testing.T) {
	s := level.NewStorage[int32](8, []uint8{2, 4})

	s.Set(0, 1)
	s.Set(0, 2)
	s.Set(0, 1)

	got := s.Get(0)
	if got != 1 {
		t.Errorf("Get(0) = %d, want 1 (|= does not clear old bits)", got)
	}
}

func TestResizePreservesData(t *testing.T) {
	s := level.NewStorage[int32](8, []uint8{2, 4})

	s.Set(0, 10)
	s.Set(1, 20)
	s.Set(2, 30)
	s.Set(3, 40)
	s.Set(4, 50)

	for i, want := range []int32{10, 20, 30, 40} {
		got := s.Get(i)
		if got != want {
			t.Errorf("Get(%d) = %d, want %d (Resize repacks with old BPE)", i, got, want)
		}
	}
}

func TestBpeZeroSingleValue(t *testing.T) {
	s := level.NewStorage[int32](64, possibleBpe)
	s.Palette = append(s.Palette, 42)

	for i := range 64 {
		if s.Get(i) != 42 {
			t.Errorf("Get(%d) = %d, want 42 (BPE 0 should always return first palette entry)", i, s.Get(i))
		}
	}

	s.Set(30, 42)
	s.Set(0, 42)

	for i := range 64 {
		if s.Get(i) != 42 {
			t.Errorf("Get(%d) = %d, want 42 after Set same value", i, s.Get(i))
		}
	}
}

func TestBpeZeroGrowPallet(t *testing.T) {
	s := level.NewStorage[int32](64, possibleBpe)
	s.Palette = append(s.Palette, 10)

	s.Set(0, 10)
	s.Set(1, 20)

	if s.Get(0) != 10 {
		t.Errorf("Get(0) = %d, want 10 after grow from BPE 0", s.Get(0))
	}
	if s.Get(1) != 20 {
		t.Errorf("Get(1) = %d, want 20 after grow from BPE 0", s.Get(1))
	}
}

func TestResizeElementsPerLong(t *testing.T) {
	s := level.NewStorage[int32](8, []uint8{2, 4})

	s.Set(0, 10)
	s.Set(1, 20)
	s.Set(2, 30)
	s.Set(3, 40)
	s.Set(4, 50)

	if s.ElementsPerLong != 16 {
		t.Errorf("ElementsPerLong = %d, want 16 (Resize stores old ElementsPerLong)", s.ElementsPerLong)
	}
}

func TestImportData(t *testing.T) {
	data := []uint64{
		2924043581155788451,
		3014291839833165474,
		2924008396783439527,
		3014291770978157730,
		2905993930227668138,
		2978157076574191776,
		2888296810347582631,
		2941952219986869920,
		3122589749976320164,
		2833795322612367018,
		3086525630071259805,
		2833830575707083949,
		3014115504847927964,
		2851845594302469291,
		2905958400778123933,
		3014291839833162913,
		2815745670742620830,
		2942057979530264736,
		2815745670876841123,
		2887944277684467868,
		2815745740402934945,
		2887944277549201564,
		3032094718851891360,
		2887908817896552616,
		2887944277684469920,
		2869788867884569248,
		2869894626083422368,
		2834183176950988960,
		2869894625949204127,
		3032376331669946015,
		2869894694802636456,
		2996136153003671199,
		2869929879174992040,
		2869894625948941983,
		2887944277684735139,
		2869894625948941983,
		22322953887,
	}
	n := int32(384 + 1)
	bpe := util.BpeByNum(float64(n)) // overworld world height
	s, err := level.ImportStorage[int32](data, bpe, util.SliceRange(0, n), level.ChunkColumns, nil)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(s.RealValues())
}
