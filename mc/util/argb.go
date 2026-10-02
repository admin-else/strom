package util

import "math"

// ARGB mirrors the integer part of net.minecraft.util.ARGB. The Vec3/Vector3f
// overloads are omitted because they depend on world.phys/joml types that are
// not part of this headless subset; add them when a consumer needs them.

const argbLinearChannelDepth = 1024

var srgbToLinear = makeSrgbToLinear()
var linearToSrgb = makeLinearToSrgb()

func computeSrgbToLinear(x float32) float32 {
	if x >= 0.04045 {
		return float32(math.Pow((float64(x)+0.055)/1.055, 2.4))
	}
	return x / 12.92
}

func computeLinearToSrgb(x float32) float32 {
	if x >= 0.0031308 {
		return float32(1.055*math.Pow(float64(x), 0.4166666666666667) - 0.055)
	}
	return 12.92 * x
}

func makeSrgbToLinear() (lookup [256]int16) {
	for i := range lookup {
		channel := float32(i) / 255.0
		lookup[i] = int16(round32(computeSrgbToLinear(channel) * 1023.0))
	}
	return lookup
}

func makeLinearToSrgb() (lookup [1024]byte) {
	for i := range lookup {
		channel := float32(i) / 1023.0
		lookup[i] = byte(round32(computeLinearToSrgb(channel) * 255.0))
	}
	return lookup
}

func round32(value float32) int {
	return int(math.Floor(float64(value) + 0.5))
}

// SRGBToLinearChannel mirrors ARGB.srgbToLinearChannel(int).
func SRGBToLinearChannel(srgb int) float32 {
	return float32(srgbToLinear[srgb]) / 1023.0
}

// LinearToSrgbChannel mirrors ARGB.linearToSrgbChannel(float).
func LinearToSrgbChannel(linear float32) int {
	return int(linearToSrgb[Floor32(linear*1023.0)]) & 0xFF
}

// Floor32 mirrors Mth.floor(float).
func Floor32(value float32) int {
	return int(math.Floor(float64(value)))
}

// ArgbAlpha mirrors ARGB.alpha(int).
func ArgbAlpha(color int) int {
	return color >> 24 & 0xFF
}

// ArgbRed mirrors ARGB.red(int).
func ArgbRed(color int) int {
	return color >> 16 & 0xFF
}

// ArgbGreen mirrors ARGB.green(int).
func ArgbGreen(color int) int {
	return color >> 8 & 0xFF
}

// ArgbBlue mirrors ARGB.blue(int).
func ArgbBlue(color int) int {
	return color & 0xFF
}

// ArgbColor4 mirrors ARGB.color(int, int, int, int).
func ArgbColor4(galpha int, red int, green int, blue int) int {
	return (galpha&0xFF)<<24 | (red&0xFF)<<16 | (green&0xFF)<<8 | blue&0xFF
}

// ArgbColor3 mirrors ARGB.color(int, int, int).
func ArgbColor3(red int, green int, blue int) int {
	return ArgbColor4(255, red, green, blue)
}

// ArgbMultiply mirrors ARGB.multiply(int, int).
func ArgbMultiply(lhs int, rhs int) int {
	if lhs == -1 {
		return rhs
	}
	if rhs == -1 {
		return lhs
	}
	return ArgbColor4(
		ArgbAlpha(lhs)*ArgbAlpha(rhs)/255,
		ArgbRed(lhs)*ArgbRed(rhs)/255,
		ArgbGreen(lhs)*ArgbGreen(rhs)/255,
		ArgbBlue(lhs)*ArgbBlue(rhs)/255,
	)
}

// ArgbOpaque mirrors ARGB.opaque(int).
func ArgbOpaque(color int) int {
	return color | -16777216
}

// ArgbTransparent mirrors ARGB.transparent(int).
func ArgbTransparent(color int) int {
	return color & 16777215
}

// ArgbAs8BitChannel mirrors ARGB.as8BitChannel(float).
func ArgbAs8BitChannel(value float32) int {
	return Floor32(value * 255.0)
}

// ArgbAlphaFloat mirrors ARGB.alphaFloat(int).
func ArgbAlphaFloat(color int) float32 {
	return float32(ArgbAlpha(color)) / 255.0
}

// ArgbRedFloat mirrors ARGB.redFloat(int).
func ArgbRedFloat(color int) float32 {
	return float32(ArgbRed(color)) / 255.0
}

// ArgbGreenFloat mirrors ARGB.greenFloat(int).
func ArgbGreenFloat(color int) float32 {
	return float32(ArgbGreen(color)) / 255.0
}

// ArgbBlueFloat mirrors ARGB.blueFloat(int).
func ArgbBlueFloat(color int) float32 {
	return float32(ArgbBlue(color)) / 255.0
}

// ToABGR mirrors ARGB.toABGR(int).
func ToABGR(color int) int {
	return color&-16711936 | (color&0xFF0000)>>16 | (color&0xFF)<<16
}

// FromABGR mirrors ARGB.fromABGR(int).
func FromABGR(color int) int {
	return ToABGR(color)
}
