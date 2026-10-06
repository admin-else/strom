package phys

// HitResultType mirrors HitResult.Type.
type HitResultType int

const (
	HitResultMISS HitResultType = iota
	HitResultBLOCK
	HitResultENTITY
)

// HitResult mirrors net.minecraft.world.phys.HitResult. Java's distanceTo(Entity)
// is deferred with the entity unit.
type HitResult struct {
	Location Vec3
}

// NewHitResult mirrors HitResult(Vec3).
func NewHitResult(location Vec3) (ret HitResult) { return HitResult{Location: location} }

// GetLocation mirrors HitResult.getLocation().
func (h HitResult) GetLocation() (ret Vec3) { return h.Location }

// DistanceToXYZ mirrors HitResult.distanceTo(Entity) for an explicit position.
func (h HitResult) DistanceToXYZ(x float64, y float64, z float64) (ret float64) {
	xd := h.Location.X - x
	yd := h.Location.Y - y
	zd := h.Location.Z - z
	return xd*xd + yd*yd + zd*zd
}
