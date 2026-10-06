package phys

import "github.com/admin-else/strom/mc/core"

// AxisStepOrder mirrors Direction.axisStepOrder(Vec3). It lives in phys rather
// than core because core cannot import phys (phys imports core for Direction).
func AxisStepOrder(movement Vec3) (ret []core.DirectionAxis) {
	if absF(movement.X) < absF(movement.Z) {
		return []core.DirectionAxis{core.AxisY, core.AxisZ, core.AxisX}
	}
	return []core.DirectionAxis{core.AxisY, core.AxisX, core.AxisZ}
}

func absF(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}
