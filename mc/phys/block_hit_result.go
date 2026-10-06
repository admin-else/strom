package phys

import "github.com/admin-else/strom/mc/core"

// BlockHitResult mirrors net.minecraft.world.phys.BlockHitResult. The network
// StreamCodec is deferred.
type BlockHitResult struct {
	HitResult
	direction      core.Direction
	blockPos       core.BlockPos
	miss           bool
	inside         bool
	worldBorderHit bool
}

// NewBlockHitResult mirrors BlockHitResult(Vec3, Direction, BlockPos, boolean).
func NewBlockHitResult(location Vec3, direction core.Direction, pos core.BlockPos, inside bool) (ret BlockHitResult) {
	return newBlockHitResult(false, location, direction, pos, inside, false)
}

// NewBlockHitResultFull mirrors BlockHitResult(Vec3, Direction, BlockPos, boolean, boolean).
func NewBlockHitResultFull(location Vec3, direction core.Direction, pos core.BlockPos, inside bool, worldBorderHit bool) (ret BlockHitResult) {
	return newBlockHitResult(false, location, direction, pos, inside, worldBorderHit)
}

// BlockHitResultMiss mirrors BlockHitResult.miss(Vec3, Direction, BlockPos).
func BlockHitResultMiss(location Vec3, direction core.Direction, pos core.BlockPos) (ret BlockHitResult) {
	return newBlockHitResult(true, location, direction, pos, false, false)
}

func newBlockHitResult(miss bool, location Vec3, direction core.Direction, blockPos core.BlockPos, inside bool, worldBorderHit bool) (ret BlockHitResult) {
	return BlockHitResult{
		HitResult:      NewHitResult(location),
		miss:           miss,
		direction:      direction,
		blockPos:       blockPos,
		inside:         inside,
		worldBorderHit: worldBorderHit,
	}
}

// WithDirection mirrors BlockHitResult.withDirection(Direction).
func (b BlockHitResult) WithDirection(direction core.Direction) (ret BlockHitResult) {
	return newBlockHitResult(b.miss, b.Location, direction, b.blockPos, b.inside, b.worldBorderHit)
}

// WithPosition mirrors BlockHitResult.withPosition(BlockPos).
func (b BlockHitResult) WithPosition(pos core.BlockPos) (ret BlockHitResult) {
	return newBlockHitResult(b.miss, b.Location, b.direction, pos, b.inside, b.worldBorderHit)
}

// HitBorder mirrors BlockHitResult.hitBorder().
func (b BlockHitResult) HitBorder() (ret BlockHitResult) {
	return newBlockHitResult(b.miss, b.Location, b.direction, b.blockPos, b.inside, true)
}

// GetBlockPos mirrors BlockHitResult.getBlockPos().
func (b BlockHitResult) GetBlockPos() (ret core.BlockPos) { return b.blockPos }

// GetDirection mirrors BlockHitResult.getDirection().
func (b BlockHitResult) GetDirection() (ret core.Direction) { return b.direction }

// GetType mirrors BlockHitResult.getType().
func (b BlockHitResult) GetType() (ret HitResultType) {
	if b.miss {
		return HitResultMISS
	}
	return HitResultBLOCK
}

// IsInside mirrors BlockHitResult.isInside().
func (b BlockHitResult) IsInside() (ret bool) { return b.inside }

// IsWorldBorderHit mirrors BlockHitResult.isWorldBorderHit().
func (b BlockHitResult) IsWorldBorderHit() (ret bool) { return b.worldBorderHit }
