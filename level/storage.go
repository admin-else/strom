package level

import (
	"errors"
	"math"
	"slices"
)

type DynamicPalletLookup[T comparable] interface {
	GetIndex(T) uint32
	GetT(uint33 uint32) T
	Resize(storage *Storage[T]) error
}

type Storage[T comparable] struct {
	Palette                    []T
	CustomPallet               DynamicPalletLookup[T]
	Data                       []uint64
	PossibleBitsPerEntryValues []uint8 // These should be sorted
	BitsPerEntry               uint8
	Len                        int
	MaxAtCurrentBpe            int
	ElementsPerLong            int
	Mask                       uint64
}

var CannotGrowPalletErr = errors.New("cannot grow pallet")
var CannotGrowPalletNoPossibleBpeErr = errors.New("cannot grow pallet, no possible bpe provided")
var Bpe0ButNot1LenPaletteErr = errors.New("bpe 0 but not 1, len palette")
var InvalidPaletteIndexErr = errors.New("invalid palette index")

func NewStorage[T comparable](len int, possibleBpe []uint8) (s *Storage[T]) {
	s = new(Storage[T])
	s.Len = len
	s.PossibleBitsPerEntryValues = possibleBpe
	s.Resize(possibleBpe[0])
	return
}

func ImportStorage[T comparable](data []uint64, bpe uint8, palette []T, len int, possibleBpe []uint8) (s *Storage[T], err error) {
	s = new(Storage[T])
	s.Resize(bpe)
	s.Data = data
	s.Palette = palette
	s.Len = len
	s.PossibleBitsPerEntryValues = possibleBpe
	err = s.CheckData()
	return
}

func ImportStoragePaletteFunc[T comparable](data []uint64, bpe uint8, paletteF DynamicPalletLookup[T], len int, possibleBpe []uint8) (s *Storage[T], err error) {
	s = new(Storage[T])
	s.Resize(bpe)
	s.Data = data
	s.CustomPallet = paletteF
	s.Len = len
	s.PossibleBitsPerEntryValues = possibleBpe
	err = s.CheckData()
	return
}

func (s *Storage[T]) CheckData() (err error) {
	if s.BitsPerEntry == 0 && len(s.Palette) != 1 {
		err = Bpe0ButNot1LenPaletteErr
		return
	}
	for i := range s.Len {
		if s.GetLow(i) >= uint32(len(s.Palette)) {
			err = InvalidPaletteIndexErr
			return
		}
	}
	return
}

func (s *Storage[T]) Resize(bpe uint8) {
	if bpe == 0 {
		s.Data = nil
		s.BitsPerEntry = 0
		s.ElementsPerLong = 0
		s.Mask = 0
		s.MaxAtCurrentBpe = 1
		return
	}
	elementsPerLong := int(64 / int(bpe))
	numberOfLongs := int(math.Ceil(float64(s.Len) / float64(elementsPerLong)))
	newData := make([]uint64, numberOfLongs)
	if s.BitsPerEntry == 0 {
		s.Data = newData
		s.BitsPerEntry = bpe
		s.ElementsPerLong = elementsPerLong
		s.Mask = (uint64(1) << uint32(bpe)) - 1
		s.MaxAtCurrentBpe = 1 << s.BitsPerEntry
		return
	}
	for i := range s.Len {
		v := s.GetLow(i)
		newData[i/elementsPerLong] |= uint64(v) << ((i % elementsPerLong) * int(bpe))
	}
	s.Data = newData
	s.BitsPerEntry = bpe
	s.ElementsPerLong = elementsPerLong
	s.Mask = (uint64(1) << uint32(bpe)) - 1
	s.MaxAtCurrentBpe = 1 << s.BitsPerEntry
	return
}

func (s *Storage[T]) GrowPallet(v T) (index int, err error) {
	index = len(s.Palette)
	if index+1 > s.MaxAtCurrentBpe {
		if s.PossibleBitsPerEntryValues == nil {
			err = CannotGrowPalletNoPossibleBpeErr
			return
		}

		bpei := slices.Index(s.PossibleBitsPerEntryValues, s.BitsPerEntry)
		if bpei == len(s.PossibleBitsPerEntryValues)-1 {
			err = CannotGrowPalletErr
			return
		}
		newBpe := s.PossibleBitsPerEntryValues[bpei+1]
		s.Resize(newBpe)
	}
	s.Palette = append(s.Palette, v)
	return
}

func (s *Storage[T]) Set(i int, v T) (err error) {
	pi := slices.Index(s.Palette, v)
	if pi != -1 {
		if s.BitsPerEntry == 0 {
			return
		}
		s.SetLow(i, uint32(pi))
		return
	}
	pi, err = s.GrowPallet(v)
	if err != nil {
		return
	}
	s.SetLow(i, uint32(pi))
	return
}

func (s *Storage[T]) Get(i int) (v T) {
	if s.BitsPerEntry == 0 {
		return s.Palette[0]
	}
	return s.Palette[s.GetLow(i)]
}

func (s *Storage[T]) GetLow(i int) (j uint32) {
	j = uint32((s.Data[i/s.ElementsPerLong] >> ((i % s.ElementsPerLong) * int(s.BitsPerEntry))) & s.Mask)
	return
}

func (s *Storage[T]) SetLow(i int, v uint32) {
	s.Data[i/s.ElementsPerLong] &= ^(s.Mask << ((i % s.ElementsPerLong) * int(s.BitsPerEntry)))
	s.Data[i/s.ElementsPerLong] |= uint64(v) << ((i % s.ElementsPerLong) * int(s.BitsPerEntry))
}

func (s *Storage[T]) RealValues() []T {
	ret := make([]T, s.Len)
	for i := range s.Len {
		ret[i] = s.Get(i)
	}
	return ret
}

func (s *Storage[T]) SetMany(i, j int, v T) (err error) {
	if i == 0 && j == s.Len-1 && slices.Contains(s.PossibleBitsPerEntryValues, 0) {
		s.Resize(0)
		s.Palette = []T{v}
		return
	}
	pi := slices.Index(s.Palette, v)
	if pi != -1 {
		if s.BitsPerEntry == 0 {
			return
		}
		s.SetLow(i, uint32(pi))
		return
	}
	pi, err = s.GrowPallet(v)
	if err != nil {
		return
	}
	s.SetManyLow(i, j, uint32(pi))
	return
}

func (s *Storage[T]) SetManyLow(i, j int, v uint32) {
	n := j - i
	if n > s.ElementsPerLong*2 {
		for k := i; k < j; k++ {
			s.SetLow(k, v)
		}
		return
	}
	for i%s.ElementsPerLong == 0 {
		s.SetLow(i, v)
		i++
	}

	var vFilled uint64
	for k := range s.ElementsPerLong {
		vFilled |= uint64(v) << (k * int(s.BitsPerEntry))
	}

	t := i / s.ElementsPerLong
	for k := range (j - i) / s.ElementsPerLong {
		s.Data[t+k] = vFilled
	}

	for j%s.ElementsPerLong != 0 {
		s.SetLow(j, v)
		j--
	}
}
