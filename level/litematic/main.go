package main

import (
	"github.com/admin-else/strom/data"
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
	for k, r := range s.Regions {
		s.BlockStatePalletRealized[k] = make([]int32, len(r.BlockStatePalette))
		for i, b := range r.BlockStatePalette {
			var stateId int32
			stateId, err = data.StateIdFromBlocKAndStateMap(ret.MinecraftVersion, b.Name, b.Properties)
			if err != nil {
				return
			}
			s.BlockStatePalletRealized[k][i] = stateId
		}
	}
	return
}
