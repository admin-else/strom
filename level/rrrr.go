package level

import (
	"errors"
	"math"
	"slices"
)

type RealStorageFormat struct {
	AvailableBpes []uint8
	BiggestDirect bool
	Len           int32
}

var BpeNotAvailableErr = errors.New("bpe not available")
var OutOfBoundsErr = errors.New("out of bounds")
var AvailableBpeMustBeLenOneErr = errors.New("available bpe must be 1")
var NoAvailableBpeErr = errors.New("no available bpe")
var Bpe0PalletMustBeLenOneErr = errors.New("bpe 0 pallet must be 1")
var RealInvalidPaletteIndexErr = errors.New("invalid palette index")
var AvailableBpeMustBeSortedErr = errors.New("available bpe must be sorted")
var DirectModeCantResizeErr = errors.New("direct mode cannot resize")

func (r RealStorageFormat) Import(data []uint64, bpe uint8, palette []int32) (s *RealStorage, err error) {
	if !slices.IsSorted(r.AvailableBpes) {
		err = AvailableBpeMustBeSortedErr
		return
	}

	if !slices.Contains(r.AvailableBpes, bpe) {
		err = BpeNotAvailableErr
		return
	}
	s = new(RealStorage)
	s.format = r
	s.resize(bpe)
	if bpe == 0 && len(palette) != 1 {
		err = Bpe0PalletMustBeLenOneErr
		return
	}
	if data != nil {
		s.data = data
	}
	s.palette = palette
	return
}

func (r RealStorageFormat) ImportDataDirect(data []uint64) (s *RealStorage, err error) {
	if len(r.AvailableBpes) != 1 {
		err = AvailableBpeMustBeLenOneErr
		return
	}
	return r.Import(data, r.AvailableBpes[0], nil)
}

type RealStorage struct {
	format  RealStorageFormat
	data    []uint64
	bpe     uint8
	palette []int32

	elementsPerLong, maxAtCurrentBpe int32

	mask uint64
}

func (r *RealStorage) doWeResizeDirect(bpe uint8) bool {
	if !r.format.BiggestDirect || r.format.AvailableBpes == nil {
		return false
	}
	return slices.Index(r.format.AvailableBpes, bpe) == len(r.format.AvailableBpes)-1
}

// only resizes upwards
func (r *RealStorage) resize(bpe uint8) {
	if bpe == 0 {
		r.data = nil
		r.bpe = 0
		r.elementsPerLong = 0
		r.mask = 0
		r.maxAtCurrentBpe = 1
		return
	}
	numberOfElements := int32(64 / int(bpe))
	numberOfLongs := int(math.Ceil(float64(r.format.Len) / float64(numberOfElements)))
	newData := make([]uint64, numberOfLongs)
	mask := (uint64(1) << uint32(bpe)) - 1
	maxAtCurrentBpe := int32(1 << bpe)
	if r.bpe == 0 {
		r.bpe = bpe
		r.elementsPerLong = numberOfElements
		r.data = newData
		r.mask = mask
		r.maxAtCurrentBpe = maxAtCurrentBpe
		return
	}
	resizeDirect := r.doWeResizeDirect(bpe)
	for i := range r.format.Len {
		var v int32
		if resizeDirect {
			v, _ = r.Get(i)
		} else {
			v = int32((r.data[i/r.elementsPerLong] >> ((i % r.elementsPerLong) * int32(r.bpe))) & r.mask)
		}
		r.data[i/r.elementsPerLong] |= uint64(v) << ((i % r.elementsPerLong) * int32(bpe))
	}
	if resizeDirect {
		r.palette = nil
	}
	r.elementsPerLong = numberOfElements
	r.data = newData
	r.mask = mask
	r.maxAtCurrentBpe = maxAtCurrentBpe

	return
}

func (r *RealStorage) Get(i int32) (s int32, err error) {
	if i < 0 || i >= r.format.Len {
		err = OutOfBoundsErr
		return
	}
	if r.bpe == 0 { // single mode
		if len(r.palette) != 1 {
			err = Bpe0PalletMustBeLenOneErr
			return
		}
		s = r.palette[0]
		return
	}

	v := int32((r.data[i/r.elementsPerLong] >> ((i % r.elementsPerLong) * int32(r.bpe))) & r.mask)
	if r.palette == nil { // direct
		s = v
		return
	}
	if v < 0 || v >= int32(len(r.palette)) {
		err = RealInvalidPaletteIndexErr
		return
	}

	s = r.palette[v]
	return
}

func (r *RealStorage) set(i int32, v int32) {
	shift := (i % r.elementsPerLong) * int32(r.bpe)
	r.data[i/r.elementsPerLong] = (r.data[i/r.elementsPerLong] & ^(r.mask << shift)) | (uint64(v) << shift)
}

func (r *RealStorage) Set(i int32, v int32) (err error) {
	if i < 0 || i >= r.format.Len {
		err = OutOfBoundsErr
		return
	}
	if r.palette != nil {
		iv := slices.Index(r.palette, v)
		if iv == -1 {
			iv, err = r.appendPalette(v)
		}
		v = int32(iv)
	}
	r.set(i, v)
	return
}

func (r *RealStorage) appendPalette(v int32) (i int, err error) {
	i = len(r.palette)
	if r.palette == nil {
		err = DirectModeCantResizeErr
		return
	}
	if i < int(r.maxAtCurrentBpe) {
		r.palette = append(r.palette, v)
		return
	}
	if r.format.AvailableBpes == nil {
		err = NoAvailableBpeErr
		return
	}

	bpei := slices.Index(r.format.AvailableBpes, r.bpe)
	if bpei == len(r.format.AvailableBpes)-1 {
		err = NoAvailableBpeErr
		return
	}
	newBpe := r.format.AvailableBpes[bpei+1]
	r.resize(newBpe)
	if r.palette == nil {
		i = int(v) // direct mode
		return
	}

	r.palette = append(r.palette, v)
	i = len(r.palette) - 1
	return
}
