package util

// Mth is the subset of net.minecraft.util.Mth used by the renderer. Mth is
// headless-reusable (see AGENTS.md) and lives in strom. Only the methods that
// the currently translated units need are ported; more are added as consumers
// appear.

var multiplyDeBruijnBitPosition = [32]int{
	0, 1, 28, 2, 29, 14, 24, 3, 30, 22, 20, 15, 25, 17, 4, 8,
	31, 27, 13, 23, 21, 19, 16, 7, 26, 12, 18, 6, 11, 5, 10, 9,
}

// IsMultipleOf mirrors Mth.isMultipleOf(int, int).
func IsMultipleOf(dividend int, divisor int) bool {
	return dividend%divisor == 0
}

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
