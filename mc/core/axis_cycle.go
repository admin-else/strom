package core

import "github.com/admin-else/strom/mc/util"

// AxisCycle mirrors net.minecraft.core.AxisCycle. Java models it as an enum with
// per-constant overrides; Go uses a value type with a switch.
type AxisCycle int

const (
	AxisCycleNONE AxisCycle = iota
	AxisCycleFORWARD
	AxisCycleBACKWARD
)

// AxisCycleVALUES mirrors AxisCycle.VALUES.
var AxisCycleVALUES = []AxisCycle{AxisCycleNONE, AxisCycleFORWARD, AxisCycleBACKWARD}

// CycleInt mirrors AxisCycle.cycle(int, int, int, Direction.Axis).
func (c AxisCycle) CycleInt(x int, y int, z int, axis DirectionAxis) (ret int) {
	switch c {
	case AxisCycleNONE:
		return axis.ChooseInt(x, y, z)
	case AxisCycleFORWARD:
		return axis.ChooseInt(z, x, y)
	default:
		return axis.ChooseInt(y, z, x)
	}
}

// CycleDouble mirrors AxisCycle.cycle(double, double, double, Direction.Axis).
func (c AxisCycle) CycleDouble(x float64, y float64, z float64, axis DirectionAxis) (ret float64) {
	switch c {
	case AxisCycleNONE:
		return axis.ChooseDouble(x, y, z)
	case AxisCycleFORWARD:
		return axis.ChooseDouble(z, x, y)
	default:
		return axis.ChooseDouble(y, z, x)
	}
}

// CycleAxis mirrors AxisCycle.cycle(Direction.Axis).
func (c AxisCycle) CycleAxis(axis DirectionAxis) (ret DirectionAxis) {
	switch c {
	case AxisCycleNONE:
		return axis
	case AxisCycleFORWARD:
		return DirectionAxisVALUES[floorModInt(int(axis)+1, 3)]
	default:
		return DirectionAxisVALUES[floorModInt(int(axis)-1, 3)]
	}
}

// Inverse mirrors AxisCycle.inverse().
func (c AxisCycle) Inverse() (ret AxisCycle) {
	switch c {
	case AxisCycleNONE:
		return AxisCycleNONE
	case AxisCycleFORWARD:
		return AxisCycleBACKWARD
	default:
		return AxisCycleFORWARD
	}
}

// AxisCycleBetween mirrors AxisCycle.between(Direction.Axis, Direction.Axis).
func AxisCycleBetween(from DirectionAxis, to DirectionAxis) (ret AxisCycle) {
	return AxisCycleVALUES[floorModInt(int(to)-int(from), 3)]
}

func floorModInt(a int, b int) (ret int) { return util.PositiveModuloInt(a, b) }
