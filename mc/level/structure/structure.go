package structure

import (
	"github.com/admin-else/strom/mc/nbt"
	"github.com/go-viper/mapstructure/v2"
)

func Load(fileName string) (s *Structure, err error) {
	n, err := nbt.ReadUnstructuredFilePath(fileName)
	if err != nil {
		return
	}
	s = &Structure{}
	err = mapstructure.Decode(n.Value, s)
	return
}

type Structure struct {
	DataVersion int32
	Blocks      []struct {
		Pos   [3]int32
		State int
	}
	Palette []struct {
		Name       string
		Properties map[string]string
	}
	Size [3]int32
}

func (s *Structure) GetBlock(x, y, z int32) (name string) {
	for _, block := range s.Blocks {
		if block.Pos[0] == x && block.Pos[1] == y && block.Pos[2] == z {
			palletIndex := block.State

			if palletIndex >= len(s.Palette) {
				return ""
			}

			return s.Palette[palletIndex].Name
		}
	}
	return ""
}
