package world

// MotionBlockingHeightAt returns the MOTION_BLOCKING height for world column
// (x, z), i.e. the Y of the first air block above the terrain surface. It
// mirrors ClientLevel.getHeight(Heightmap.Types.MOTION_BLOCKING, x, z), which
// the weather column renderer uses to stop rain/snow at the surface.
func (w *World) MotionBlockingHeightAt(x, z int32) (h int32, err error) {
	return w.HeightAt("motion_blocking", x, z)
}
