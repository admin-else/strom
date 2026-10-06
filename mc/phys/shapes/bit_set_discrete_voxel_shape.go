package shapes

import (
	"math"
	"math/bits"

	"github.com/admin-else/strom/mc/core"
)

// BitSetDiscreteVoxelShape mirrors net.minecraft.world.phys.shapes.BitSetDiscreteVoxelShape.
// Java's BitSet is a []uint64 bit set here.
type BitSetDiscreteVoxelShape struct {
	storage          bitSet
	sizeX            int
	sizeY            int
	sizeZ            int
	xMin, yMin, zMin int
	xMax, yMax, zMax int
}

// NewBitSetDiscreteVoxelShape mirrors BitSetDiscreteVoxelShape(int, int, int).
func NewBitSetDiscreteVoxelShape(xSize int, ySize int, zSize int) (ret *BitSetDiscreteVoxelShape) {
	return &BitSetDiscreteVoxelShape{
		storage: newBitSet(xSize * ySize * zSize),
		sizeX:   xSize,
		sizeY:   ySize,
		sizeZ:   zSize,
		xMin:    xSize,
		yMin:    ySize,
		zMin:    zSize,
	}
}

// WithFilledBounds mirrors BitSetDiscreteVoxelShape.withFilledBounds(...).
func WithFilledBounds(xSize int, ySize int, zSize int, xMin int, yMin int, zMin int, xMax int, yMax int, zMax int) (ret *BitSetDiscreteVoxelShape) {
	shape := NewBitSetDiscreteVoxelShape(xSize, ySize, zSize)
	shape.xMin, shape.yMin, shape.zMin = xMin, yMin, zMin
	shape.xMax, shape.yMax, shape.zMax = xMax, yMax, zMax
	for x := xMin; x < xMax; x++ {
		for y := yMin; y < yMax; y++ {
			for z := zMin; z < zMax; z++ {
				shape.fillUpdateBounds(x, y, z, false)
			}
		}
	}
	return shape
}

// NewBitSetDiscreteVoxelShapeFrom mirrors BitSetDiscreteVoxelShape(DiscreteVoxelShape).
func NewBitSetDiscreteVoxelShapeFrom(voxelShape DiscreteVoxelShape) (ret *BitSetDiscreteVoxelShape) {
	ret = &BitSetDiscreteVoxelShape{
		sizeX: voxelShape.xSize(),
		sizeY: voxelShape.ySize(),
		sizeZ: voxelShape.zSize(),
	}
	if other, ok := voxelShape.(*BitSetDiscreteVoxelShape); ok {
		ret.storage = other.storage.clone()
	} else {
		ret.storage = newBitSet(ret.sizeX * ret.sizeY * ret.sizeZ)
		for x := 0; x < ret.sizeX; x++ {
			for y := 0; y < ret.sizeY; y++ {
				for z := 0; z < ret.sizeZ; z++ {
					if voxelShape.isFull(x, y, z) {
						ret.storage.set(ret.getIndex(x, y, z))
					}
				}
			}
		}
	}
	ret.xMin = voxelShape.firstFull(core.AxisX)
	ret.yMin = voxelShape.firstFull(core.AxisY)
	ret.zMin = voxelShape.firstFull(core.AxisZ)
	ret.xMax = voxelShape.lastFull(core.AxisX)
	ret.yMax = voxelShape.lastFull(core.AxisY)
	ret.zMax = voxelShape.lastFull(core.AxisZ)
	return
}

func (s *BitSetDiscreteVoxelShape) xSize() int { return s.sizeX }
func (s *BitSetDiscreteVoxelShape) ySize() int { return s.sizeY }
func (s *BitSetDiscreteVoxelShape) zSize() int { return s.sizeZ }

func (s *BitSetDiscreteVoxelShape) getIndex(x int, y int, z int) (ret int) {
	return (x*s.sizeY+y)*s.sizeZ + z
}

// IsFull mirrors BitSetDiscreteVoxelShape.isFull(int, int, int).
func (s *BitSetDiscreteVoxelShape) isFull(x int, y int, z int) (ret bool) {
	return s.storage.get(s.getIndex(x, y, z))
}

func (s *BitSetDiscreteVoxelShape) fillUpdateBounds(x int, y int, z int, updateBounds bool) {
	s.storage.set(s.getIndex(x, y, z))
	if updateBounds {
		s.xMin = minInt(s.xMin, x)
		s.yMin = minInt(s.yMin, y)
		s.zMin = minInt(s.zMin, z)
		s.xMax = maxInt(s.xMax, x+1)
		s.yMax = maxInt(s.yMax, y+1)
		s.zMax = maxInt(s.zMax, z+1)
	}
}

// Fill mirrors BitSetDiscreteVoxelShape.fill(int, int, int).
func (s *BitSetDiscreteVoxelShape) fill(x int, y int, z int) {
	s.fillUpdateBounds(x, y, z, true)
}

// IsEmpty mirrors BitSetDiscreteVoxelShape.isEmpty().
func (s *BitSetDiscreteVoxelShape) isEmpty() (ret bool) { return s.storage.isEmpty() }

// FirstFull mirrors BitSetDiscreteVoxelShape.firstFull(Direction.Axis).
func (s *BitSetDiscreteVoxelShape) firstFull(axis core.DirectionAxis) (ret int) {
	return axis.ChooseInt(s.xMin, s.yMin, s.zMin)
}

// LastFull mirrors BitSetDiscreteVoxelShape.lastFull(Direction.Axis).
func (s *BitSetDiscreteVoxelShape) lastFull(axis core.DirectionAxis) (ret int) {
	return axis.ChooseInt(s.xMax, s.yMax, s.zMax)
}

// BitSetJoin mirrors BitSetDiscreteVoxelShape.join(...).
func BitSetJoin(first DiscreteVoxelShape, second DiscreteVoxelShape, xMerger IndexMerger, yMerger IndexMerger, zMerger IndexMerger, op BooleanOp) (ret *BitSetDiscreteVoxelShape) {
	shape := NewBitSetDiscreteVoxelShape(xMerger.Size()-1, yMerger.Size()-1, zMerger.Size()-1)
	boundMin := [3]int{math.MaxInt32, math.MaxInt32, math.MaxInt32}
	boundMax := [3]int{math.MinInt32, math.MinInt32, math.MinInt32}
	xMerger.ForMergedIndexes(func(x1 int, x2 int, xr int) bool {
		updatedSlice := false
		yMerger.ForMergedIndexes(func(y1 int, y2 int, yr int) bool {
			updatedColumn := false
			zMerger.ForMergedIndexes(func(z1 int, z2 int, zr int) bool {
				if op(DvsIsFullWide(first, x1, y1, z1), DvsIsFullWide(second, x2, y2, z2)) {
					shape.storage.set(shape.getIndex(xr, yr, zr))
					boundMin[2] = minInt(boundMin[2], zr)
					boundMax[2] = maxInt(boundMax[2], zr)
					updatedColumn = true
				}
				return true
			})
			if updatedColumn {
				boundMin[1] = minInt(boundMin[1], yr)
				boundMax[1] = maxInt(boundMax[1], yr)
				updatedSlice = true
			}
			return true
		})
		if updatedSlice {
			boundMin[0] = minInt(boundMin[0], xr)
			boundMax[0] = maxInt(boundMax[0], xr)
		}
		return true
	})
	shape.xMin, shape.yMin, shape.zMin = boundMin[0], boundMin[1], boundMin[2]
	shape.xMax, shape.yMax, shape.zMax = boundMax[0]+1, boundMax[1]+1, boundMax[2]+1
	return shape
}

func bitSetForAllBoxes(voxelShape DiscreteVoxelShape, consumer IntLineConsumer, mergeNeighbors bool) {
	shape := NewBitSetDiscreteVoxelShapeFrom(voxelShape)
	for y := 0; y < shape.sizeY; y++ {
		for x := 0; x < shape.sizeX; x++ {
			lastStartZ := -1
			for z := 0; z <= shape.sizeZ; z++ {
				if DvsIsFullWide(shape, x, y, z) {
					if mergeNeighbors {
						if lastStartZ == -1 {
							lastStartZ = z
						}
					} else {
						consumer(x, y, z, x+1, y+1, z+1)
					}
				} else if lastStartZ != -1 {
					endX := x
					endY := y
					shape.clearZStrip(lastStartZ, z, x, y)
					for shape.isZStripFull(lastStartZ, z, endX+1, y) {
						shape.clearZStrip(lastStartZ, z, endX+1, y)
						endX++
					}
					for shape.isXZRectangleFull(x, endX+1, lastStartZ, z, endY+1) {
						for cx := x; cx <= endX; cx++ {
							shape.clearZStrip(lastStartZ, z, cx, endY+1)
						}
						endY++
					}
					consumer(x, y, lastStartZ, endX+1, endY+1, z)
					lastStartZ = -1
				}
			}
		}
	}
}

func (s *BitSetDiscreteVoxelShape) isZStripFull(startZ int, endZ int, x int, y int) (ret bool) {
	if x < s.sizeX && y < s.sizeY {
		return s.storage.nextClearBit(s.getIndex(x, y, startZ)) >= s.getIndex(x, y, endZ)
	}
	return false
}

func (s *BitSetDiscreteVoxelShape) isXZRectangleFull(startX int, endX int, startZ int, endZ int, y int) (ret bool) {
	for x := startX; x < endX; x++ {
		if !s.isZStripFull(startZ, endZ, x, y) {
			return false
		}
	}
	return true
}

func (s *BitSetDiscreteVoxelShape) clearZStrip(startZ int, endZ int, x int, y int) {
	s.storage.clearRange(s.getIndex(x, y, startZ), s.getIndex(x, y, endZ))
}

// IsInterior mirrors BitSetDiscreteVoxelShape.isInterior(int, int, int).
func (s *BitSetDiscreteVoxelShape) IsInterior(x int, y int, z int) (ret bool) {
	isInterior := x > 0 && x < s.sizeX-1 && y > 0 && y < s.sizeY-1 && z > 0 && z < s.sizeZ-1
	return isInterior &&
		s.isFull(x, y, z) &&
		s.isFull(x-1, y, z) &&
		s.isFull(x+1, y, z) &&
		s.isFull(x, y-1, z) &&
		s.isFull(x, y+1, z) &&
		s.isFull(x, y, z-1) &&
		s.isFull(x, y, z+1)
}

func minInt(a int, b int) (ret int) {
	if a < b {
		return a
	}
	return b
}

func maxInt(a int, b int) (ret int) {
	if a > b {
		return a
	}
	return b
}

// bitSet is a minimal stand-in for java.util.BitSet with the operations the
// voxel shape uses.
type bitSet struct {
	words []uint64
}

func newBitSet(n int) (ret bitSet) {
	if n < 0 {
		n = 0
	}
	return bitSet{words: make([]uint64, (n+63)/64)}
}

func (b *bitSet) get(i int) (ret bool) {
	if i < 0 {
		return false
	}
	w := i >> 6
	if w >= len(b.words) {
		return false
	}
	return b.words[w]>>(uint(i)&63)&1 == 1
}

func (b *bitSet) set(i int) {
	w := i >> 6
	for w >= len(b.words) {
		b.words = append(b.words, 0)
	}
	b.words[w] |= uint64(1) << (uint(i) & 63)
}

func (b *bitSet) clearRange(from int, to int) {
	for i := from; i < to; i++ {
		w := i >> 6
		if w < len(b.words) {
			b.words[w] &^= uint64(1) << (uint(i) & 63)
		}
	}
}

func (b bitSet) isEmpty() (ret bool) {
	for _, w := range b.words {
		if w != 0 {
			return false
		}
	}
	return true
}

func (b bitSet) clone() (ret bitSet) {
	c := make([]uint64, len(b.words))
	copy(c, b.words)
	return bitSet{words: c}
}

func (b bitSet) nextClearBit(from int) (ret int) {
	if from < 0 {
		from = 0
	}
	i := from
	for {
		w := i >> 6
		if w >= len(b.words) {
			return i
		}
		word := ^b.words[w]
		word &^= (uint64(1) << (uint(i) & 63)) - 1
		if word != 0 {
			return w*64 + bits.TrailingZeros64(word)
		}
		i = (w + 1) * 64
	}
}
