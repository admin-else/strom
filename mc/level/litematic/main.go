package litematic

import (
	"errors"
	"io"
	"os"
	"strings"
	"unsafe"

	data2 "git.anygate.cloud/anygatecloud/strom/mc/data"
	"git.anygate.cloud/anygatecloud/strom/mc/level"
	"git.anygate.cloud/anygatecloud/strom/mc/nbt"
	"git.anygate.cloud/anygatecloud/strom/mc/util"
	"github.com/go-viper/mapstructure/v2"
)

var OutOfBoundsErr = errors.New("out of bounds")

const mcNamespaced = "minecraft:"

type Position struct {
	X, Y, Z int32
}

func (p Position) FixNegatives() Position {
	if p.X < 0 {
		p.X *= -1
	}
	if p.Y < 0 {
		p.Y *= -1
	}
	if p.Z < 0 {
		p.Z *= -1
	}
	return p
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
	*level.Storage
}

func (s *Structure) GetBlockAt(region string, x, y, z int32) (stateId int32, err error) {
	r := s.Regions[region]
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
	stateId, err = r.Get(i)
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
	Version string
}

// LoadFromPath loads a litematic file from disk.
func LoadFromPath(path string) (s *Structure, err error) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()
	s, err = LoadFromFile(f)
	return
}

// LoadFromFile reads and parses a litematic file from a reader.
func LoadFromFile(r io.Reader) (s *Structure, err error) {
	n, err := nbt.ReadUnstructuredFile(r)
	if err != nil {
		return
	}
	return Load(n.Value)
}

// Load decodes raw NBT data into a litematic Structure.
func Load(v any) (s *Structure, err error) {
	s = &Structure{}
	err = mapstructure.Decode(v, &s.File)
	if err != nil {
		return
	}
	ret, err := data2.LookUpVersionByDataVersion(s.MinecraftDataVersion)
	if err != nil {
		return
	}
	s.Version = ret.MinecraftVersion
	for k, region := range s.Regions {
		palette := make([]int32, len(region.BlockStatePalette))
		for i, b := range region.BlockStatePalette {
			var paletteEntry int32
			paletteEntry, err = data2.StateIdFromBlocKAndStateMap(s.Version, strings.TrimPrefix(b.Name, mcNamespaced), b.Properties)
			if err != nil {
				return
			}
			palette[i] = paletteEntry
		}

		sizeFixed := region.Size.FixNegatives()
		l := sizeFixed.X * sizeFixed.Y * sizeFixed.Z
		bpe := util.BpeByNum(float64(len(region.BlockStatePalette)))
		uintSlice := unsafe.Slice((*uint64)(unsafe.Pointer(unsafe.SliceData(region.BlockStates))), len(region.BlockStates))
		region.Storage, err = level.StorageFormat{
			AvailableBpes: []uint8{bpe},
			BiggestDirect: false,
			Len:           l,
		}.Import(uintSlice, bpe, palette)
		if err != nil {
			return
		}
		s.Regions[k] = region
	}

	return
}
