package control

import (
	botentity "github.com/admin-else/strom/mc/bot/entity"
	"github.com/admin-else/strom/mc/entity"
	"github.com/admin-else/strom/mc/phys"
	"github.com/admin-else/strom/mc/phys/shapes"
)

// EntityProvider supplies the collision shapes of the tracked remote entities so
// the local player collides with them. The local player itself must be excluded
// by the provider.
type EntityProvider interface {
	EntityCollisionShapes(source *entity.Entity, box phys.AABB) (ret []*shapes.VoxelShape)
}

// entityTracker adapts a bot entity module into an EntityProvider, approximating
// every tracked entity with a player-sized box.
type entityTracker struct {
	module *botentity.Module
}

// NewEntityTracker adapts a bot entity module into an EntityProvider. The module
// only tracks server-spawned entities, so the local player is never included.
func NewEntityTracker(module *botentity.Module) (ret EntityProvider) {
	return &entityTracker{module: module}
}

// EntityCollisionShapes mirrors Level.getEntityCollisions.
func (t *entityTracker) EntityCollisionShapes(source *entity.Entity, box phys.AABB) (ret []*shapes.VoxelShape) {
	for _, tracked := range t.module.Entities() {
		shape := shapes.Box(
			tracked.Position.X-0.3, tracked.Position.Y, tracked.Position.Z-0.3,
			tracked.Position.X+0.3, tracked.Position.Y+1.8, tracked.Position.Z+0.3,
		)
		bounds, ok := shape.Bounds()
		if ok && bounds.Intersects(box) {
			ret = append(ret, shape)
		}
	}
	return
}
