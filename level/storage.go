package level

import (
	"errors"
	"math"
	"slices"
)

type Storage[T comparable] struct {
	Palette                    []T
	Data                       []uint64
	PossibleBitsPerEntryValues []uint8 // These should be sorted
	BitsPerEntry               uint8
	Len                        int
	MaxAtCurrentBpe            int
	ElementsPerLong            int
	Mask                       uint64
}

var CannotGrowPalletErr = errors.New("cannot grow pallet")

func NewStorage[T comparable](len int, possibleBpe []uint8) (s *Storage[T]) {
	s = new(Storage[T])
	s.Len = len
	s.PossibleBitsPerEntryValues = possibleBpe
	s.Resize(possibleBpe[0])
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
	if s.BitsPerEntry == 0 {
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
