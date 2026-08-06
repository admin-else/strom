package anvil

import (
	"strings"
	"unsafe"

	data2 "git.anygate.cloud/anygatecloud/strom/mc/data"
	level2 "git.anygate.cloud/anygatecloud/strom/mc/level"
	"git.anygate.cloud/anygatecloud/strom/mc/util"
)

type Chunk struct {
	DataVersion      int32
	Heightmaps       map[string][]int64
	InhabitedTime    int64
	XPos, YPos, ZPos int32
	Sections         []Section
}

const McPrefix = "minecraft:"

func (c Chunk) ToStorage(version string) (toChunk *level2.Chunk, err error) {
	f := level2.MakeBlockFormat(version)
	toChunk = new(level2.Chunk)
	toChunk.Version = version

	toChunk.Sections = make([]level2.Section, len(c.Sections))
	for i, s := range c.Sections {
		pallete := make([]int32, len(s.BlockStates.Palette))
		for j, p := range s.BlockStates.Palette {
			pallete[j], err = data2.StateIdFromBlocKAndStateMap(version, strings.TrimPrefix(p.Name, McPrefix), p.Properties)
			if err != nil {
				return
			}
		}
		// this is probably fine

		if s.BlockStates.Data != nil {
			uintSlice := unsafe.Slice((*uint64)(unsafe.Pointer(unsafe.SliceData(s.BlockStates.Data))), len(s.BlockStates.Data))
			bpe := util.BpeByNum(float64(len(s.BlockStates.Palette)))
			toChunk.Sections[i].Blocks, err = f.Import(uintSlice, bpe, pallete)
			if err != nil {
				return
			}
			toChunk.Sections[i].BlockCount = int16(level2.BlocksPerChunkSection - level2.CountStorage(toChunk.Sections[i].Blocks, 0))
		} else {
			var fillerId = int32(0)
			if len(pallete) == 1 {
				fillerId = pallete[0]
				if fillerId != 0 {
					toChunk.Sections[i].BlockCount = level2.BlocksPerChunkSection
				}
			}
			toChunk.Sections[i].Blocks, err = f.FullWith(fillerId)
			if err != nil {
				return
			}
		}
	}
	f = level2.MakeBiomeFormat(version)
	for i, s := range c.Sections {
		pallete := make([]int32, len(s.Biomes.Palette))
		for j, p := range s.Biomes.Palette {
			var b *data2.Biome
			b, ok := data2.LookupBiomeByName(version, strings.TrimPrefix(p, McPrefix))
			if !ok {
				err = data2.UnknownBiomeNameErr
				return
			}
			pallete[j] = b.Id
		}

		if s.Biomes.Data != nil {

			// this is probably fine
			uintSlice := unsafe.Slice((*uint64)(unsafe.Pointer(unsafe.SliceData(s.Biomes.Data))), len(s.Biomes.Data))
			bpe := util.BpeByNum(float64(len(pallete)))
			toChunk.Sections[i].Biomes, err = f.Import(uintSlice, bpe, pallete)
			if err != nil {
				return
			}
		} else {
			var fillerId = int32(0)
			if len(pallete) == 1 {
				fillerId = pallete[0]
			}
			toChunk.Sections[i].Biomes, err = f.FullWith(fillerId)
			if err != nil {
				return
			}
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
			Properties map[string]string `nbt:"Properties,omitempty"`
		}
	} `nbt:"block_states,omitempty"`
	Biomes struct {
		Data    []int64 `nbt:"data,omitempty"`
		Palette []string
	} `nbt:"biomes,omitempty"`
}
