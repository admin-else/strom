package gen

type NoiseParam int

const (
	NoiseParamTemperature NoiseParam = iota
	NoiseParamHumidity
	NoiseParamContinentalNess
	NoiseParamErosion
	NoiseParamShift
	NoiseParamWeirdness
)

type NoiseParamMapV struct {
	A []float64
	S string
	O int
}

var NoiseParamMap = map[NoiseParam]NoiseParamMapV{
	NoiseParamShift:           {A: []float64{1, 1, 1, 0}, S: "minecraft:offset", O: -3},
	NoiseParamTemperature:     {A: []float64{1.5, 0, 1, 0, 0, 0}, S: "minecraft:temperature", O: -10},
	NoiseParamHumidity:        {A: []float64{1, 1, 0, 0, 0, 0}, S: "minecraft:vegetation", O: -8},
	NoiseParamContinentalNess: {A: []float64{1, 1, 2, 2, 2, 1, 1, 1, 1}, S: "minecraft:continentalness", O: -9},
	NoiseParamErosion:         {A: []float64{1, 1, 0, 1, 1}, S: "minecraft:erosion", O: -9},
	NoiseParamWeirdness:       {A: []float64{1, 2, 1, 0, 0, 0}, S: "minecraft:ridge", O: -7},
}

func initClimateSeed(n NoiseParam, x Xoroshiro, large bool, nmax int) DoublePerlin {
	v := NoiseParamMap[n]
	s := v.S
	omin := v.O
	if large && n != NoiseParamShift && n != NoiseParamWeirdness {
		s += "_large"
		omin -= 2
	}
	x.XorString(s)
	return DoublePerlinFromXoroshiro(&x, v.A, omin, len(v.A), nmax)
}

type BiomeNoise struct {
	Temperature, Humidity, ContinentalNess, Erosion, Shift, Weirdness DoublePerlin
}

func BiomeNoiseFromXoroshiro(x Xoroshiro, large bool) (v BiomeNoise) {
	x = x.Split()

	v.Temperature = initClimateSeed(NoiseParamTemperature, x, large, -1)
	v.Humidity = initClimateSeed(NoiseParamHumidity, x, large, -1)
	v.ContinentalNess = initClimateSeed(NoiseParamContinentalNess, x, large, -1)
	v.Erosion = initClimateSeed(NoiseParamErosion, x, large, -1)
	v.Shift = initClimateSeed(NoiseParamShift, x, large, -1)
	v.Weirdness = initClimateSeed(NoiseParamWeirdness, x, large, -1)
	return
}

func (b BiomeNoise) Sample(x, y, z int) (t, h, c, e, d, w float64, combinedData [6]int64, id int) {
	px := float64(x)
	pz := float64(z)
	px += b.Shift.Sample(px, 0, pz) * 4
	pz += b.Shift.Sample(pz, float64(x), 0) * 4

	c = b.ContinentalNess.Sample(px, 0, pz)
	e = b.Erosion.Sample(px, 0, pz)
	w = b.Weirdness.Sample(px, 0, pz)
	t = b.Temperature.Sample(px, 0, pz)
	h = b.Humidity.Sample(px, 0, pz)

	combinedData = [6]int64{
		int64(float32(t) * 10000.0),
		int64(float32(h) * 10000.0),
		int64(float32(c) * 10000.0),
		int64(float32(e) * 10000.0),
		int64(float32(d) * 10000.0),
		int64(float32(w) * 10000.0),
	}
	return
}
