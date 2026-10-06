// Package control drives a headless player with the ported Entity physics and
// input handling. It resolves block collision shapes from minecraft-data and
// feeds them to the entity package, so a bot can walk, fall, jump and step up
// against the world it receives from the server.
package control

import (
	"github.com/admin-else/strom/mc/bot/world"
	"github.com/admin-else/strom/mc/data"
	"github.com/admin-else/strom/mc/entity"
	"github.com/admin-else/strom/mc/phys"
	"github.com/admin-else/strom/mc/phys/shapes"
	"github.com/admin-else/strom/mc/util"
)

// Level implements entity.CollisionGetter on top of a bot world plus
// minecraft-data collision shapes. Entity collisions are not modelled yet.
type Level struct {
	World   *world.World
	Version string
}

// NewLevel wraps a shared world for the given protocol version.
func NewLevel(w *world.World, version string) (ret *Level) {
	return &Level{World: w, Version: version}
}

// GetEntityCollisions mirrors Level.getEntityCollisions for the collision slice.
func (l *Level) GetEntityCollisions(source *entity.Entity, box phys.AABB) (ret []*shapes.VoxelShape) {
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
