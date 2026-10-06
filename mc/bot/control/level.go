// Package control drives a headless player with the ported Entity physics and
// input handling. It resolves block collision shapes from minecraft-data and
// feeds them to the entity package, so a bot can walk, fall, jump and step up
// against the world it receives from the server.
package control

import (
	"github.com/admin-else/strom/mc/bot/world"
	"github.com/admin-else/strom/mc/core"
	"github.com/admin-else/strom/mc/data"
	"github.com/admin-else/strom/mc/entity"
	"github.com/admin-else/strom/mc/phys"
	"github.com/admin-else/strom/mc/phys/shapes"
	"github.com/admin-else/strom/mc/util"
)

// Level implements entity.CollisionGetter on top of a bot world plus
// minecraft-data collision shapes. Entity collisions are supplied by an
// optional EntityProvider.
type Level struct {
	World    *world.World
	Version  string
	Entities EntityProvider
}

// NewLevel wraps a shared world for the given protocol version.
func NewLevel(w *world.World, version string) (ret *Level) {
	return &Level{World: w, Version: version}
}

// GetEntityCollisions mirrors Level.getEntityCollisions for the collision slice.
func (l *Level) GetEntityCollisions(source *entity.Entity, box phys.AABB) (ret []*shapes.VoxelShape) {
	if l.Entities != nil {
		return l.Entities.EntityCollisionShapes(source, box)
	}
	return nil
}

// GetBlockCollisions mirrors Level.getBlockCollisions: every block state whose
// collision shape intersects box contributes its boxes as voxel shapes.
func (l *Level) GetBlockCollisions(source *entity.Entity, box phys.AABB) (ret []*shapes.VoxelShape) {
	minX := util.FloorDouble(box.MinX)
	maxX := util.FloorDouble(box.MaxX)
	minY := util.FloorDouble(box.MinY)
	maxY := util.FloorDouble(box.MaxY)
	minZ := util.FloorDouble(box.MinZ)
	maxZ := util.FloorDouble(box.MaxZ)

	for x := minX; x <= maxX; x++ {
		for y := minY; y <= maxY; y++ {
			for z := minZ; z <= maxZ; z++ {
				stateId, err := l.World.GetBlock(int32(x), int32(y), int32(z))
				if err != nil {
					continue
				}
				boxes, ok := data.CollisionShapeBoxesAt(l.Version, stateId)
				if !ok {
					continue
				}
				for _, b := range boxes {
					ret = append(ret, shapes.Box(b[0], b[1], b[2], b[3], b[4], b[5]).MoveXYZ(float64(x), float64(y), float64(z)))
				}
			}
		}
	}
	return
}

// blockName returns the block identifier at pos, or ok=false when the chunk is
// not loaded.
func (l *Level) blockName(pos core.BlockPos) (name string, ok bool) {
	stateId, err := l.World.GetBlock(int32(pos.X), int32(pos.Y), int32(pos.Z))
	if err != nil {
		return "", false
	}
	block, found := data.LookupBlockByStateId(l.Version, stateId)
	if !found {
		return "", false
	}
	return block.Name, true
}

// MovementContext implementation (net.minecraft.world.level.Level fluid/climb
// queries). Fluids and climbables are recognised by block name; block friction
// is approximated from the vanilla defaults.

// IsInWater mirrors Level.getFluidState(pos).is(FluidTags.WATER).
func (l *Level) IsInWater(pos core.BlockPos) (ret bool) {
	name, ok := l.blockName(pos)
	return ok && name == "water"
}

// IsInLava mirrors Level.getFluidState(pos).is(FluidTags.LAVA).
func (l *Level) IsInLava(pos core.BlockPos) (ret bool) {
	name, ok := l.blockName(pos)
	return ok && name == "lava"
}

// IsInShallowFluid approximates LivingEntity.isInShallowFluid by treating every
// fluid contact as shallow (the world store does not carry per-block fluid
// heights yet).
func (l *Level) IsInShallowFluid(pos core.BlockPos) (ret bool) { return true }

// OnClimbable approximates LivingEntity.onClimbable() by block name.
func (l *Level) OnClimbable(pos core.BlockPos) (ret bool) {
	name, ok := l.blockName(pos)
	if !ok {
		return false
	}
	switch name {
	case "ladder", "vine", "scaffolding", "twisting_vines", "weeping_vines",
		"cave_vines", "cave_vines_plant", "chain":
		return true
	}
	return false
}

// BlockFriction mirrors Block.getFriction() for the common blocks.
func (l *Level) BlockFriction(pos core.BlockPos) (ret float32) {
	name, _ := l.blockName(pos)
	switch name {
	case "ice", "packed_ice", "blue_ice", "frosted_ice":
		return 0.98
	case "slime_block":
		return 0.8
	}
	return 0.6
}

// NoCollision mirrors CollisionGetter.noCollision(AABB).
func (l *Level) NoCollision(source *entity.Entity, box phys.AABB) (ret bool) {
	return len(l.GetBlockCollisions(source, box)) == 0
}
