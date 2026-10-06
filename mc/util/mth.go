package util

import "math"

// Mth is the subset of net.minecraft.util.Mth used by the renderer and the
// physics layer. Mth is headless-reusable (see AGENTS.md) and lives in strom.
// Only the methods that the currently translated units need are ported; more
// are added as consumers appear.

// Mth.PI and friends.
const (
	MthPI       float32 = float32(math.Pi)
	MthHalfPI   float32 = float32(math.Pi / 2)
	MthTwoPI    float32 = float32(math.Pi * 2)
	MthDegToRad float32 = float32(math.Pi / 180.0)
	MthRadToDeg float32 = float32(180.0) / float32(math.Pi)
	MthEpsilon  float32 = 1.0e-5
)

// SIN_SCALE mirrors Mth.SIN_SCALE (10430.378350470453).
const sinScale = 10430.378350470453

// sinTable mirrors Mth.SIN: a 65536-entry quantized sine table.
var sinTable = buildSinTable()

func buildSinTable() (table [65536]float32) {
	for i := range table {
		table[i] = float32(math.Sin(float64(i) / sinScale))
	}
	table[0] = 0.0
	table[16384] = 1.0
	table[32768] = 0.0
	table[49152] = -1.0
	return
}

// Sin mirrors Mth.sin(double).
func Sin(i float64) (ret float32) {
	if i >= 0.0 {
		return sinTable[int((int64(i*sinScale+0.5))&65535)]
	}
	return -sinTable[int((int64(-i*sinScale+0.5))&65535)]
}

// Cos mirrors Mth.cos(double).
func Cos(i float64) (ret float32) {
	if i >= 0.0 {
		return sinTable[int((int64(i*sinScale+16384.0+0.5))&65535)]
	}
	return sinTable[int((int64(-i*sinScale+16384.0+0.5))&65535)]
}

// Sqrt mirrors Mth.sqrt(float).
func Sqrt(x float32) (ret float32) { return float32(math.Sqrt(float64(x))) }

// FloorFloat mirrors Mth.floor(float).
func FloorFloat(v float32) (ret int) { return toInt32(math.Floor(float64(v))) }

// FloorDouble mirrors Mth.floor(double).
func FloorDouble(v float64) (ret int) { return toInt32(math.Floor(v)) }

// LFloor mirrors Mth.lfloor(double).
func LFloor(v float64) (ret int64) { return toInt64(math.Floor(v)) }

// AbsFloat mirrors Mth.abs(float).
func AbsFloat(v float32) (ret float32) { return float32(math.Abs(float64(v))) }

// AbsInt mirrors Mth.abs(int).
func AbsInt(v int) (ret int) {
	if v < 0 {
		return -v
	}
	return v
}

// CeilFloat mirrors Mth.ceil(float).
func CeilFloat(v float32) (ret int) { return toInt32(math.Ceil(float64(v))) }

// CeilDouble mirrors Mth.ceil(double).
func CeilDouble(v float64) (ret int) { return toInt32(math.Ceil(v)) }

// toInt32 mirrors a Java `(int)` narrowing cast of a double: NaN becomes 0 and
// out-of-range values saturate to the 32-bit bounds (JLS 5.1.3).
func toInt32(v float64) (ret int) {
	if math.IsNaN(v) {
		return 0
	}
	if v >= math.MaxInt32 {
		return math.MaxInt32
	}
	if v <= math.MinInt32 {
		return math.MinInt32
	}
	return int(v)
}

// toInt64 mirrors a Java `(long)` narrowing cast of a double.
func toInt64(v float64) (ret int64) {
	if math.IsNaN(v) {
		return 0
	}
	if v >= math.MaxInt64 {
		return math.MaxInt64
	}
	if v <= math.MinInt64 {
		return math.MinInt64
	}
	return int64(v)
}

// BinarySearch mirrors Mth.binarySearch(int, int, IntPredicate).
func BinarySearch(from int, to int, condition func(int) bool) (ret int) {
	length := to - from
	for length > 0 {
		half := length / 2
		middle := from + half
		if condition(middle) {
			length = half
		} else {
			from = middle + 1
			length -= half + 1
		}
	}
	return from
}

// ClampedLerpFloat mirrors Mth.clampedLerp(float, float, float).
func ClampedLerpFloat(factor float32, minValue float32, maxValue float32) (ret float32) {
	if factor < 0.0 {
		return minValue
	}
	if factor > 1.0 {
		return maxValue
	}
	return LerpFloat(factor, minValue, maxValue)
}

// ClampedLerpDouble mirrors Mth.clampedLerp(double, double, double).
func ClampedLerpDouble(factor float64, minValue float64, maxValue float64) (ret float64) {
	if factor < 0.0 {
		return minValue
	}
	if factor > 1.0 {
		return maxValue
	}
	return LerpDouble(factor, minValue, maxValue)
}

// LerpFloat mirrors Mth.lerp(float, float, float).
func LerpFloat(alpha float32, p0 float32, p1 float32) (ret float32) {
	return p0 + alpha*(p1-p0)
}

// LerpDouble mirrors Mth.lerp(double, double, double).
func LerpDouble(alpha float64, p0 float64, p1 float64) (ret float64) {
	return p0 + alpha*(p1-p0)
}

// SquareFloat mirrors Mth.square(float).
func SquareFloat(x float32) (ret float32) { return x * x }

// SquareDouble mirrors Mth.square(double).
func SquareDouble(x float64) (ret float64) { return x * x }

// SquareInt mirrors Mth.square(int).
func SquareInt(x int) (ret int) { return x * x }

// SquareLong mirrors Mth.square(long).
func SquareLong(x int64) (ret int64) { return x * x }

// LengthSquared2Double mirrors Mth.lengthSquared(double, double).
func LengthSquared2Double(x float64, y float64) (ret float64) { return x*x + y*y }

// LengthSquared2Float mirrors Mth.lengthSquared(float, float).
func LengthSquared2Float(x float32, y float32) (ret float32) { return x*x + y*y }

// LengthSquared3Double mirrors Mth.lengthSquared(double, double, double).
func LengthSquared3Double(x float64, y float64, z float64) (ret float64) { return x*x + y*y + z*z }

// LengthSquared3Float mirrors Mth.lengthSquared(float, float, float).
func LengthSquared3Float(x float32, y float32, z float32) (ret float32) { return x*x + y*y + z*z }

var multiplyDeBruijnBitPosition = [32]int{
	0, 1, 28, 2, 29, 14, 24, 3, 30, 22, 20, 15, 25, 17, 4, 8,
	31, 27, 13, 23, 21, 19, 16, 7, 26, 12, 18, 6, 11, 5, 10, 9,
}

// IsMultipleOf mirrors Mth.isMultipleOf(int, int).
func IsMultipleOf(dividend int, divisor int) bool {
	return dividend%divisor == 0
}

// PositiveModuloInt mirrors Mth.positiveModulo(int, int) / Math.floorMod.
func PositiveModuloInt(input int, mod int) (ret int) { return ((input % mod) + mod) % mod }

// ClampInt mirrors Mth.clamp(int, int, int).
func ClampInt(value int, minValue int, maxValue int) int {
	if value < minValue {
		return minValue
	}
	if value > maxValue {
		return maxValue
	}
	return value
}

// ClampLong mirrors Mth.clamp(long, long, long).
func ClampLong(value int64, minValue int64, maxValue int64) int64 {
	if value < minValue {
		return minValue
	}
	if value > maxValue {
		return maxValue
	}
	return value
}

// ClampFloat mirrors Mth.clamp(float, float, float).
func ClampFloat(value float32, minValue float32, maxValue float32) float32 {
	if value < minValue {
		return minValue
	}
	if value > maxValue {
		return maxValue
	}
	return value
}

// ClampDouble mirrors Mth.clamp(double, double, double).
func ClampDouble(value float64, minValue float64, maxValue float64) float64 {
	if value < minValue {
		return minValue
	}
	if value > maxValue {
		return maxValue
	}
	return value
}

// IsPowerOfTwo mirrors Mth.isPowerOfTwo(int).
func IsPowerOfTwo(input int) bool {
	return input != 0 && (input&(input-1)) == 0
}

// IsPowerOfTwoLong mirrors Mth.isPowerOfTwo(long).
func IsPowerOfTwoLong(input int64) bool {
	return input != 0 && (input&(input-1)) == 0
}

// SmallestEncompassingPowerOfTwo mirrors Mth.smallestEncompassingPowerOfTwo(int).
func SmallestEncompassingPowerOfTwo(input int) int {
	result := input - 1
	result |= result >> 1
	result |= result >> 2
	result |= result >> 4
	result |= result >> 8
	result |= result >> 16
	return result + 1
}

// CeilLog2 mirrors Mth.ceillog2(int).
func CeilLog2(input int) int {
	if !IsPowerOfTwo(input) {
		input = SmallestEncompassingPowerOfTwo(input)
	}
	return multiplyDeBruijnBitPosition[(int64(input)*125613361>>27)&31]
}

// Log2 mirrors Mth.log2(int).
func Log2(input int) int {
	if IsPowerOfTwo(input) {
		return CeilLog2(input)
	}
	return CeilLog2(input) - 1
}

// PositiveCeilDiv mirrors Mth.positiveCeilDiv(int, int).
func PositiveCeilDiv(input int, divisor int) int {
	return -(int)(floorDiv(int64(-input), int64(divisor)))
}

// RoundToward mirrors Mth.roundToward(int, int).
func RoundToward(input int, multiple int) int {
	return PositiveCeilDiv(input, multiple) * multiple
}

// RoundTowardLong mirrors Mth.roundToward(long, long).
func RoundTowardLong(input int64, multiple int64) int64 {
	return PositiveCeilDivLong(input, multiple) * multiple
}

// PositiveCeilDivLong mirrors Mth.positiveCeilDiv(long, long).
func PositiveCeilDivLong(input int64, divisor int64) int64 {
	return -floorDiv(-input, divisor)
}

func floorDiv(dividend int64, divisor int64) int64 {
	quotient := dividend / divisor
	if dividend%divisor != 0 && (dividend^divisor) < 0 {
		quotient--
	}
	return quotient
}
