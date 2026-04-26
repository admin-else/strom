package level_test

import (
	"testing"

	"github.com/admin-else/strom/level"
)

func heightmapFormat() level.RealStorageFormat {
	return level.RealStorageFormat{
		AvailableBpes: []uint8{9},
		BiggestDirect: true,
		Len:           256,
	}
}

func blockStorageFormat() level.RealStorageFormat {
	return level.RealStorageFormat{
		AvailableBpes: []uint8{0, 4, 5, 6, 7, 8, 15},
		BiggestDirect: true,
		Len:           4096,
	}
}

func TestRealStorageImportDataDirect(t *testing.T) {
	format := heightmapFormat()
	data := make([]uint64, 256)
	data[0] = 42

	s, err := format.ImportDataDirect(data)
	if err != nil {
		t.Fatalf("ImportDataDirect failed: %v", err)
	}

	got, err := s.Get(0)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if got != 42 {
		t.Errorf("Get(0) = %d, want 42", got)
	}
}

func TestRealStorageImportDataDirectMultipleBpe(t *testing.T) {
	format := level.RealStorageFormat{
		AvailableBpes: []uint8{4, 8},
		BiggestDirect: true,
		Len:           256,
	}

	_, err := format.ImportDataDirect([]uint64{1})
	if err != level.AvailableBpeMustBeLenOneErr {
		t.Errorf("expected AvailableBpeMustBeLenOneErr, got %v", err)
	}
}

func TestRealStorageImportUnsortedBpe(t *testing.T) {
	format := level.RealStorageFormat{
		AvailableBpes: []uint8{9, 4},
		BiggestDirect: true,
		Len:           256,
	}

	_, err := format.Import([]uint64{1}, 4, nil)
	if err != level.AvailableBpeMustBeSortedErr {
		t.Errorf("expected AvailableBpeMustBeSortedErr, got %v", err)
	}
}

func TestRealStorageImportBpeNotAvailable(t *testing.T) {
	format := heightmapFormat()

	_, err := format.Import([]uint64{1}, 4, nil)
	if err != level.BpeNotAvailableErr {
		t.Errorf("expected BpeNotAvailableErr, got %v", err)
	}
}

func TestRealStorageGetSet(t *testing.T) {
	format := blockStorageFormat()
	s, err := format.Import(nil, 4, []int32{0, 1, 2, 3})
	if err != nil {
		t.Fatalf("Import failed: %v", err)
	}

	s.Set(0, 1)
	s.Set(1, 2)
	s.Set(2, 3)

	for i, want := range []int32{1, 2, 3} {
		got, err := s.Get(int32(i))
		if err != nil {
			t.Fatalf("Get(%d) failed: %v", i, err)
		}
		if got != want {
			t.Errorf("Get(%d) = %d, want %d", i, got, want)
		}
	}
}

func TestRealStorageOutOfBounds(t *testing.T) {
	format := blockStorageFormat()
	s, err := format.Import(nil, 4, []int32{0, 1})
	if err != nil {
		t.Fatalf("Import failed: %v", err)
	}

	_, err = s.Get(4096)
	if err != level.OutOfBoundsErr {
		t.Errorf("expected OutOfBoundsErr, got %v", err)
	}

	_, err = s.Get(-1)
	if err != level.OutOfBoundsErr {
		t.Errorf("expected OutOfBoundsErr, got %v", err)
	}
}

func TestRealStorageBpeZero(t *testing.T) {
	format := level.RealStorageFormat{
		AvailableBpes: []uint8{0, 4, 8},
		BiggestDirect: true,
		Len:           64,
	}

	s, err := format.Import(nil, 0, []int32{42})
	if err != nil {
		t.Fatalf("Import failed: %v", err)
	}

	for i := int32(0); i < 64; i++ {
		got, err := s.Get(i)
		if err != nil {
			t.Fatalf("Get(%d) failed: %v", i, err)
		}
		if got != 42 {
			t.Errorf("Get(%d) = %d, want 42", i, got)
		}
	}
}

func TestRealStorageBpeZeroPaletteLen(t *testing.T) {
	format := level.RealStorageFormat{
		AvailableBpes: []uint8{0, 4, 8},
		BiggestDirect: true,
		Len:           64,
	}

	_, err := format.Import(nil, 0, []int32{1, 2})
	if err != level.Bpe0PalletMustBeLenOneErr {
		t.Errorf("expected Bpe0PalletMustBeLenOneErr, got %v", err)
	}
}

func TestRealStorageResize(t *testing.T) {
	format := level.RealStorageFormat{
		AvailableBpes: []uint8{4, 8},
		BiggestDirect: false,
		Len:           16,
	}

	s, err := format.Import(nil, 4, []int32{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15})
	if err != nil {
		t.Fatalf("Import failed: %v", err)
	}

	for i := int32(0); i < 16; i++ {
		s.Set(i, i)
	}

	for i := int32(0); i < 16; i++ {
		got, err := s.Get(i)
		if err != nil {
			t.Fatalf("Get(%d) failed: %v", i, err)
		}
		if got != i {
			t.Errorf("Get(%d) = %d, want %d", i, got, i)
		}
	}
}

func TestRealStorageResizePreservesData(t *testing.T) {
	format := level.RealStorageFormat{
		AvailableBpes: []uint8{4, 5, 6, 8},
		BiggestDirect: false,
		Len:           8,
	}

	s, err := format.Import(nil, 4, []int32{0, 1})
	if err != nil {
		t.Fatalf("Import failed: %v", err)
	}

	s.Set(0, 10)
	s.Set(1, 20)
	s.Set(2, 30)
	s.Set(3, 40)

	for i, want := range []int32{10, 20, 30, 40} {
		got, err := s.Get(int32(i))
		if err != nil {
			t.Fatalf("Get(%d) failed: %v", i, err)
		}
		if got != want {
			t.Errorf("Get(%d) = %d, want %d", i, got, want)
		}
	}
}

func TestRealStoragePaletteGrowth(t *testing.T) {
	format := level.RealStorageFormat{
		AvailableBpes: []uint8{4, 5, 6, 8},
		BiggestDirect: false,
		Len:           16,
	}

	s, err := format.Import(nil, 4, []int32{0})
	if err != nil {
		t.Fatalf("Import failed: %v", err)
	}

	s.Set(0, 10)
	s.Set(1, 20)

	got, err := s.Get(0)
	if err != nil {
		t.Fatalf("Get(0) failed: %v", err)
	}
	if got != 10 {
		t.Errorf("Get(0) = %d, want 10", got)
	}

	got, err = s.Get(1)
	if err != nil {
		t.Fatalf("Get(1) failed: %v", err)
	}
	if got != 20 {
		t.Errorf("Get(1) = %d, want 20", got)
	}
}

func TestRealStorageNoAvailableBpe(t *testing.T) {
	format := level.RealStorageFormat{
		AvailableBpes: []uint8{4},
		BiggestDirect: false,
		Len:           16,
	}

	s, err := format.Import(nil, 4, []int32{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15})
	if err != nil {
		t.Fatalf("Import failed: %v", err)
	}

	err = s.Set(0, 16)
	if err != level.NoAvailableBpeErr {
		t.Errorf("expected NoAvailableBpeErr, got %v", err)
	}
}

func TestRealStorageBiggestDirect(t *testing.T) {
	format := level.RealStorageFormat{
		AvailableBpes: []uint8{4, 8},
		BiggestDirect: true,
		Len:           64,
	}

	s, err := format.Import(nil, 8, nil)
	if err != nil {
		t.Fatalf("Import failed: %v", err)
	}

	for i := int32(0); i < 64; i++ {
		s.Set(i, i)
	}

	for i := int32(0); i < 64; i++ {
		got, err := s.Get(i)
		if err != nil {
			t.Fatalf("Get(%d) failed: %v", i, err)
		}
		if got != i {
			t.Errorf("Get(%d) = %d, want %d", i, got, i)
		}
	}
}

func TestRealStorageDirectModeNoPalette(t *testing.T) {
	format := level.RealStorageFormat{
		AvailableBpes: []uint8{9},
		BiggestDirect: true,
		Len:           256,
	}

	data := make([]uint64, 36)
	data[0] = 0x1FFFF

	s, err := format.ImportDataDirect(data)
	if err != nil {
		t.Fatalf("ImportDataDirect failed: %v", err)
	}

	got, err := s.Get(0)
	if err != nil {
		t.Fatalf("Get(0) failed: %v", err)
	}
	if got != 0x1FFFF {
		t.Errorf("Get(0) = %d, want %d", got, 0x1FFFF)
	}
}

func TestRealStorageInvalidPaletteIndex(t *testing.T) {
	format := level.RealStorageFormat{
		AvailableBpes: []uint8{4, 8},
		BiggestDirect: false,
		Len:           16,
	}

	data := []uint64{0xFFFF}
	s, err := format.Import(data, 4, []int32{0, 1})
	if err != nil {
		t.Fatalf("Import failed: %v", err)
	}

	_, err = s.Get(0)
	if err != level.RealInvalidPaletteIndexErr {
		t.Errorf("expected RealInvalidPaletteIndexErr, got %v", err)
	}
}

func TestRealStorageOverwrite(t *testing.T) {
	format := blockStorageFormat()
	s, err := format.Import(nil, 4, []int32{0, 1, 2})
	if err != nil {
		t.Fatalf("Import failed: %v", err)
	}

	s.Set(0, 1)
	s.Set(0, 2)
	s.Set(0, 1)

	got, err := s.Get(0)
	if err != nil {
		t.Fatalf("Get(0) failed: %v", err)
	}
	if got != 1 {
		t.Errorf("Get(0) = %d, want 1", got)
	}
}

func TestRealStorageSetOutOfBounds(t *testing.T) {
	format := blockStorageFormat()
	s, err := format.Import(nil, 4, []int32{0, 1})
	if err != nil {
		t.Fatalf("Import failed: %v", err)
	}

	err = s.Set(4096, 1)
	if err != level.OutOfBoundsErr {
		t.Errorf("expected OutOfBoundsErr, got %v", err)
	}

	err = s.Set(-1, 1)
	if err != level.OutOfBoundsErr {
		t.Errorf("expected OutOfBoundsErr, got %v", err)
	}
}

func TestRealStorageGetAfterResize(t *testing.T) {
	format := level.RealStorageFormat{
		AvailableBpes: []uint8{4, 5, 6, 8},
		BiggestDirect: false,
		Len:           32,
	}

	s, err := format.Import(nil, 4, []int32{0})
	if err != nil {
		t.Fatalf("Import failed: %v", err)
	}

	for i := int32(0); i < 16; i++ {
		s.Set(i, i)
	}

	for i := int32(0); i < 16; i++ {
		got, err := s.Get(i)
		if err != nil {
			t.Fatalf("Get(%d) failed: %v", i, err)
		}
		if got != i {
			t.Errorf("Get(%d) = %d, want %d", i, got, i)
		}
	}
}

func TestRealStorageResizeDirect(t *testing.T) {
	format := level.RealStorageFormat{
		AvailableBpes: []uint8{4, 8},
		BiggestDirect: true,
		Len:           64,
	}

	s, err := format.Import(nil, 4, []int32{})
	if err != nil {
		t.Fatalf("Import failed: %v", err)
	}

	for i := int32(0); i < 64; i++ {
		err = s.Set(i, i)
		if err != nil {
			t.Fatalf("Set(%d) failed: %v", i, err)
		}
	}

	for i := int32(0); i < 64; i++ {
		got, err := s.Get(i)
		if err != nil {
			t.Fatalf("Get(%d) failed: %v", i, err)
		}
		if got != i {
			t.Errorf("Get(%d) = %d, want %d", i, got, i)
		}
	}
}
