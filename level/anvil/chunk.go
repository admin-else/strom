package anvil

import (
	"strings"
	"unsafe"

	"github.com/admin-else/strom/data"
	"github.com/admin-else/strom/level"
	"github.com/admin-else/strom/util"
)

type Chunk struct {
	DataVersion      int32
	Heightmaps       map[string][]int64
	InhabitedTime    int64
	XPos, YPos, ZPos int32
	Sections         []Section
}

const McPrefix = "minecraft:"

func (c Chunk) ToStorage(version string) (toChunk *level.Chunk, err error) {
	f := level.MakeBlockFormat(version)
	toChunk = new(level.Chunk)

	toChunk.Sections = make([]level.Section, len(c.Sections))
	for i, s := range c.Sections {
		pallete := make([]int32, len(s.BlockStates.Palette))
		for j, p := range s.BlockStates.Palette {
			pallete[j], err = data.StateIdFromBlocKAndStateMap(version, strings.TrimPrefix(p.Name, McPrefix), p.Properties)
			if err != nil {
				return
			}
		}
		// this is probably fine
		uintSlice := unsafe.Slice((*uint64)(unsafe.Pointer(unsafe.SliceData(s.BlockStates.Data))), len(s.BlockStates.Data))

		toChunk.Sections[i].Blocks, err = f.Import(uintSlice, util.BpeByNum(float64(len(s.BlockStates.Palette))), pallete)
		if err != nil {
			return
		}
		toChunk.Sections[i].BlockCount = int16(level.CountStorage(toChunk.Sections[i].Blocks, 0))
	}

	for i, s := range c.Sections {
		pallete := make([]int32, len(s.Biomes.Palette))
		for j, p := range s.Biomes.Palette {
			var b *data.Biome
			b, ok := data.LookupBiomeByName(version, strings.TrimPrefix(p, McPrefix))
			if !ok {
				err = data.UnknownBiomeNameErr
				return
			}
			pallete[j] = b.Id
		}
		// this is probably fine
		uintSlice := unsafe.Slice((*uint64)(unsafe.Pointer(unsafe.SliceData(s.Biomes.Data))), len(s.Biomes.Data))

		toChunk.Sections[i].Biomes, err = f.Import(uintSlice, util.BpeByNum(float64(len(pallete))), pallete)
		if err != nil {
			return
		}
	}
	return
}

type Section struct {
	Y           int8
	BlockStates struct {
		Data    []int64 `nbt:"data,omitempty"`
		Palette []struct {
			Name       string
			Properties map[string]string `nbt:"properties,omitempty"`
		}
	}
	Biomes struct {
		Data    []int64 `nbt:"data,omitempty"`
		Palette []string
	}
}
