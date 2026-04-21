package litematic

import (
	"errors"
	"io"
	"math"
	"os"
	"strings"
	"unsafe"

	"github.com/admin-else/strom/data"
	"github.com/admin-else/strom/level"
	"github.com/admin-else/strom/nbt"
	"github.com/go-viper/mapstructure/v2"
)

type Position struct {
	X, Y, Z int32
}

type MetaData struct {
	Author        string
	Description   string
	Name          string
	TimeCreated   int64
	TimeModified  int64
	TotalBlocks   int32
	TotalVolume   int32
	RegionCount   int32
	EnclosingSize Position
}

type BlockState struct {
	Name       string
	Properties map[string]string
}

type Region struct {
	Position          Position
	Size              Position
	BlockStatePalette []BlockState
	BlockStates       []int64
}

func (r *Region) GetBlockAt(x, y, z int32) (name string, properties map[string]string, err error) {

	// this is probably fine
	uintSlice := unsafe.Slice((*uint64)(unsafe.Pointer(unsafe.SliceData(r.BlockStates))), len(r.BlockStates))
	if r.Size.X < 0 {
		r.Size.X *= -1
	}
	if r.Size.Y < 0 {
		r.Size.Y *= -1
	}
	if r.Size.Z < 0 {
		r.Size.Z *= -1
	}

	if x >= r.Size.X || y >= r.Size.Y || z >= r.Size.Z || x < 0 || y < 0 || z < 0 {
		err = OutOfBoundsErr
		return
	}
	i := x + z*r.Size.X + y*r.Size.X*r.Size.Z
	bpe := int32(math.Ceil(math.Log2(float64(len(r.BlockStatePalette)))))
	i = level.GetEntryFromLongs(uintSlice, i, bpe)
	if i >= int32(len(r.BlockStatePalette)) {
		err = BlockStateNotInPalletErr
		return
	}
	block := r.BlockStatePalette[i]
	name = block.Name
	properties = block.Properties
	return
}

type File struct {
	Version              int32
	SubVersion           int32
	MinecraftDataVersion int32
	Metadata             MetaData
	Regions              map[string]Region
}

type Structure struct {
	File
	BlockStatePalletRealized map[string][]int32
}

func LoadFromPath(path string) (s *Structure, err error) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()
	s, err = LoadFromFile(f)
	return
}

func LoadFromFile(r io.Reader) (s *Structure, err error) {
	n, err := nbt.ReadUnstructuredFile(r)
	if err != nil {
		return
	}
	return Load(n.Value)
}

func Load(v any) (s *Structure, err error) {
	s = &Structure{}
	err = mapstructure.Decode(v, &s.File)
	if err != nil {
		return
	}
	ret, err := data.LookUpVersionByDataVersion(s.MinecraftDataVersion)
	if err != nil {
		return
	}
	s.BlockStatePalletRealized = make(map[string][]int32)
	for k, r := range s.Regions {
		s.BlockStatePalletRealized[k] = make([]int32, len(r.BlockStatePalette))
		for i, b := range r.BlockStatePalette {
			var stateId int32
			stateId, err = data.StateIdFromBlocKAndStateMap(ret.MinecraftVersion, strings.TrimPrefix(b.Name, "minecraft:"), b.Properties)
			if err != nil {
				return
			}
			s.BlockStatePalletRealized[k][i] = stateId
		}
	}
	return
}

var BlockStateNotInPalletErr = errors.New("block state not in pallet")
var OutOfBoundsErr = errors.New("out of bounds")
var RegionNotInFileErr = errors.New("region not in file")

func (s *Structure) GetBlockAt(name string, x, y, z int32) (blockState int32, err error) {
	r, ok := s.Regions[name]
	if !ok {
		err = RegionNotInFileErr
		return
	}
	uintSlice := make([]uint64, len(r.BlockStates))
	for i, v := range r.BlockStates {
		uintSlice[i] = uint64(v)
	}
	if r.Size.X < 0 {
		r.Size.X *= -1
	}
	if r.Size.Y < 0 {
		r.Size.Y *= -1
	}
	if r.Size.Z < 0 {
		r.Size.Z *= -1
	}

	if x >= r.Size.X || y >= r.Size.Y || z >= r.Size.Z {
		err = OutOfBoundsErr
		return
	}
	i := x + z*r.Size.X + y*r.Size.X*r.Size.Z
	bpe := int32(math.Ceil(math.Log2(float64(len(r.BlockStatePalette)))))
	i = level.GetEntryFromLongs(uintSlice, i, bpe)
	mapSlice := s.BlockStatePalletRealized[name]
	if i >= int32(len(mapSlice)) {
		err = BlockStateNotInPalletErr
		return
	}
	blockState = mapSlice[i]
	return
}
