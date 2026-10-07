package world

import "github.com/admin-else/strom/mc/data"

// FluidKind is the fluid at a position, matching the vanilla FluidTags the
// client checks for the eye-fluid fog.
type FluidKind uint8

const (
	FluidEmpty FluidKind = iota
	FluidWater
	FluidLava
)

// String returns the lowercase kind name.
func (k FluidKind) String() (ret string) {
	switch k {
	case FluidWater:
		return "water"
	case FluidLava:
		return "lava"
	default:
		return "empty"
	}
}

// FluidState is the fluid at a world position: the kind plus the vanilla fluid
// resource id.
type FluidState struct {
	Kind FluidKind
	ID   string
}

// FluidAt returns the fluid at the given world coordinates so the client can
// pick the eye-fluid fog. The block-state decode is reused; water, lava and
// bubble columns (and waterlogged blocks) are recognised by block name.
func (w *World) FluidAt(x, y, z int32) (ret FluidState, err error) {
	stateId, err := w.GetBlock(x, y, z)
	if err != nil {
		return
	}
	block, properties, err := data.FromBlockState(w.version, stateId)
	if err != nil {
		return
	}
	switch block.Name {
	case "water", "bubble_column":
		ret = FluidState{Kind: FluidWater, ID: "minecraft:water"}
	case "lava":
		ret = FluidState{Kind: FluidLava, ID: "minecraft:lava"}
	default:
		if properties["waterlogged"] == "true" {
			ret = FluidState{Kind: FluidWater, ID: "minecraft:water"}
		} else {
			ret = FluidState{Kind: FluidEmpty, ID: "minecraft:empty"}
		}
	}
	return
}
