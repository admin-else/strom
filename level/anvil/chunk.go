package anvil

type Chunk struct {
	DataVersion      int32
	Heightmaps       map[string][]int64
	InhabitedTime    int64
	XPos, YPos, ZPos int32
	Sections         []Section
}

type Section struct {
	Y           int8
	BlockStates struct {
		Data    []int64
		Palette []struct {
			Name       string
			Properties map[string]string `nbt:"properties,omitempty"`
		}
	}
	Biomes []struct {
		Data    []int64
		Palette []string
	}
}
