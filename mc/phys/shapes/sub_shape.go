package shapes

import (
	"github.com/admin-else/strom/mc/core"
	"github.com/admin-else/strom/mc/util"
)

// SubShape mirrors net.minecraft.world.phys.shapes.SubShape.
type SubShape struct {
	parent                 DiscreteVoxelShape
	startX, startY, startZ int
	endX, endY, endZ       int
	sizeX, sizeY, sizeZ    int
}

func newSubShape(parent DiscreteVoxelShape, startX int, startY int, startZ int, endX int, endY int, endZ int) (ret *SubShape) {
	return &SubShape{
		parent: parent,
		startX: startX, startY: startY, startZ: startZ,
		endX: endX, endY: endY, endZ: endZ,
		sizeX: endX - startX, sizeY: endY - startY, sizeZ: endZ - startZ,
	}
}

func (s *SubShape) xSize() int { return s.sizeX }
func (s *SubShape) ySize() int { return s.sizeY }
func (s *SubShape) zSize() int { return s.sizeZ }

// IsFull mirrors SubShape.isFull(int, int, int).
func (s *SubShape) isFull(x int, y int, z int) (ret bool) {
	return s.parent.isFull(s.startX+x, s.startY+y, s.startZ+z)
}

// Fill mirrors SubShape.fill(int, int, int).
func (s *SubShape) fill(x int, y int, z int) {
	s.parent.fill(s.startX+x, s.startY+y, s.startZ+z)
}

// FirstFull mirrors SubShape.firstFull(Direction.Axis).
func (s *SubShape) firstFull(axis core.DirectionAxis) (ret int) {
	return s.clampToShape(axis, s.parent.firstFull(axis))
}

// LastFull mirrors SubShape.lastFull(Direction.Axis).
func (s *SubShape) lastFull(axis core.DirectionAxis) (ret int) {
	return s.clampToShape(axis, s.parent.lastFull(axis))
}

func (s *SubShape) clampToShape(axis core.DirectionAxis, parentResult int) (ret int) {
	start := axis.ChooseInt(s.startX, s.startY, s.startZ)
	end := axis.ChooseInt(s.endX, s.endY, s.endZ)
	return util.ClampInt(parentResult, start, end) - start
}
