package level

import (
	"encoding/binary"
	"errors"
	"io"
	"math"
	"slices"
)

type StorageFormat struct {
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

func (r StorageFormat) Import(data []uint64, bpe uint8, palette []int32) (s *Storage, err error) {
	if !slices.IsSorted(r.AvailableBpes) {
		err = AvailableBpeMustBeSortedErr
		return
	}

	s = new(Storage)
	s.format = r
	s.resize(bpe)

	if !slices.Contains(r.AvailableBpes, bpe) {
		var newBpe uint8
		for _, newBpe = range r.AvailableBpes {
			if bpe < newBpe {
				continue
			}
		}
		var tmps *Storage
		tmps, err = StorageFormat{
			AvailableBpes: []uint8{bpe},
			BiggestDirect: false,
			Len:           r.Len,
		}.ImportDataDirect(data)
		if err != nil {
			return
		}
		tmps.resize(newBpe)
		data = tmps.data
		bpe = newBpe
	}
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

func (r StorageFormat) ImportDataDirect(data []uint64) (s *Storage, err error) {
	if len(r.AvailableBpes) != 1 {
		err = AvailableBpeMustBeLenOneErr
		return
	}
	return r.Import(data, r.AvailableBpes[0], nil)
}

func (r StorageFormat) ImportFromReader(reader io.Reader, bpe uint8, palette []int32) (s *Storage, err error) {
	if !slices.Contains(r.AvailableBpes, bpe) {
		err = BpeNotAvailableErr
		return
	}
	var longs []uint64
	if bpe == 0 {
		return r.Import(longs, bpe, palette)
	}
	elementsPerLong := int32(64 / int(bpe))
	numberOfLongs := int(math.Ceil(float64(r.Len) / float64(elementsPerLong)))
	for range numberOfLongs {
		var long uint64
		err = binary.Read(reader, binary.BigEndian, &long) // surely this works?
		if err != nil {
			return
		}
		longs = append(longs, long)
	}
	return r.Import(longs, bpe, palette)
}

func (r StorageFormat) FullWith(v int32) (s *Storage, err error) {
	return r.Import(nil, 0, []int32{v})
}

type Storage struct {
	format  StorageFormat
	data    []uint64
	bpe     uint8
	palette []int32

	elementsPerLong, maxAtCurrentBpe int32

	mask uint64
}

func (r *Storage) doWeResizeDirect(bpe uint8) bool {
	if !r.format.BiggestDirect || r.format.AvailableBpes == nil {
		return false
	}
	return slices.Index(r.format.AvailableBpes, bpe) == len(r.format.AvailableBpes)-1
}

// only resizes upwards
func (r *Storage) resize(bpe uint8) {
	if bpe == 0 {
		r.data = nil
		r.bpe = 0
		r.elementsPerLong = 0
		r.mask = 0
		r.maxAtCurrentBpe = 1
		return
	}
	elementsPerLong := int32(64 / int(bpe))
	numberOfLongs := int(math.Ceil(float64(r.format.Len) / float64(elementsPerLong)))
	newData := make([]uint64, numberOfLongs)
	mask := (uint64(1) << uint32(bpe)) - 1
	maxAtCurrentBpe := int32(1 << bpe)
	if r.bpe == 0 {
		r.bpe = bpe
		r.elementsPerLong = elementsPerLong
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
		newData[i/int32(elementsPerLong)] |= uint64(v) << ((i % int32(elementsPerLong)) * int32(bpe))
	}
	if resizeDirect {
		r.palette = nil
	}
	r.bpe = bpe
	r.elementsPerLong = elementsPerLong
	r.data = newData
	r.mask = mask
	r.maxAtCurrentBpe = maxAtCurrentBpe

	return
}

func (r *Storage) Get(i int32) (s int32, err error) {
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

func (r *Storage) set(i int32, v int32) {
	shift := (i % r.elementsPerLong) * int32(r.bpe)
	r.data[i/r.elementsPerLong] = (r.data[i/r.elementsPerLong] & ^(r.mask << shift)) | (uint64(v) << shift)
}

func (r *Storage) Set(i int32, v int32) (err error) {
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

func (r *Storage) appendPalette(v int32) (i int, err error) {
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

func (r *Storage) Export() (data []uint64, bpe uint8, palette []int32, format StorageFormat) {
	return r.data, r.bpe, r.palette, r.format
}

func CompareStorage(a, b *Storage) bool {
	if a.format.Len != b.format.Len {
		return false
	}
	for i := range a.format.Len {
		av, _ := a.Get(i)
		bv, _ := b.Get(i)
		if av == bv {
			return false
		}
	}
	return true
}

func CountStorage(s *Storage, v int32) (a int) {
	for i := range s.format.Len {
		if w, _ := s.Get(i); w == v {
			a++
		}
	}
	return
}
