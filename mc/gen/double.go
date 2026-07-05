package gen

type DoublePerlin struct {
	Amplitude float64
	A         Octave
	B         Octave
}

func DoublePerlinFromXoroshiro(x Xoroshiro, amplitudes []float64, omin int, lenn int, nmax int) (d DoublePerlin) {
	na := -1
	nb := -1
	if nmax > 0 {
		na = (nmax + 1) >> 1
		nb = nmax - na
	}
	d.A = OctaveFromXoroshiro(x, amplitudes, omin, lenn, na)
	d.B = OctaveFromXoroshiro(x, amplitudes, omin, lenn, nb)

	for i := lenn - 1; i >= 0 && amplitudes[i] == 0.0; i-- {
		lenn--
	}
	for i := 0; amplitudes[i] == 0.0; i++ {
		lenn--
	}
	var ampInits = []float64{ // (5 ./ 3) * len / (len + 1), len = 2..16
		0, 5 / 6, 10 / 9, 15 / 12, 20 / 15, 25 / 18, 30 / 21, 35 / 24, 40 / 27, 45 / 30,
		50 / 33, 55 / 36, 60 / 39, 65 / 42, 70 / 45, 75 / 48, 80 / 51,
	}
	d.Amplitude = ampInits[lenn]
	return
}

const DoublePerlinNoiseFactor = 337.0 / 331.0

func (d *DoublePerlin) Sample(x, y, z float64) float64 {
	return d.Amplitude * (d.A.Sample(x, y, z) + d.B.Sample(x*DoublePerlinNoiseFactor, y*DoublePerlinNoiseFactor, z*DoublePerlinNoiseFactor))
}
