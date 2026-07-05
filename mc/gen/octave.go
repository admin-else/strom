package gen

type Octave []PerlinNoise

var octaveStrHashes = []Xoroshiro{
	XoroshiroFromString("octave_-16"),
	XoroshiroFromString("octave_-15"),
	XoroshiroFromString("octave_-14"),
	XoroshiroFromString("octave_-13"),
	XoroshiroFromString("octave_-12"),
	XoroshiroFromString("octave_-11"),
	XoroshiroFromString("octave_-10"),
	XoroshiroFromString("octave_-9"),
	XoroshiroFromString("octave_-8"),
	XoroshiroFromString("octave_-7"),
	XoroshiroFromString("octave_-6"),
	XoroshiroFromString("octave_-5"),
	XoroshiroFromString("octave_-4"),
	XoroshiroFromString("octave_-3"),
	XoroshiroFromString("octave_-2"),
	XoroshiroFromString("octave_-1"),
	XoroshiroFromString("octave_0"),
}

var octaveLacunaInits = []float64{1, .5, .25, 1. / 8, 1. / 16, 1. / 32, 1. / 64, 1. / 128, 1. / 256, 1. / 512, 1. / 1024,
	1. / 2048, 1. / 4096, 1. / 8192, 1. / 16384, 1. / 32768, 1. / 65536}

var octavePersistInits = []float64{0, 1, 2. / 3, 4. / 7, 8. / 15, 16. / 31, 32. / 63, 64. / 127, 128. / 255, 256. / 511, 512. / 1023,
	1024. / 2047, 2048. / 4095, 4096. / 8191, 8192. / 16383, 16384. / 32767, 32768. / 65535}

func OctaveFromXoroshiro(x *Xoroshiro, amplitudes []float64, omin int, lenn int, nmax int) (o Octave) {
	if (-omin < 0 || -omin >= len(octaveLacunaInits)) || (lenn < 0 || lenn >= len(octavePersistInits)) {
		panic("octave initialization out of range")
	}
	octaveX := x.Split()
	i := 0
	n := 0
	lacuna := octaveLacunaInits[-omin]
	persist := octavePersistInits[lenn]

	for i < lenn && n != nmax {
		if amplitudes[i] == 0 {
			i++
			lacuna *= 2
			persist *= 0.5
			continue
		}
		copyX := octaveX
		copyX.Xor(octaveStrHashes[16+omin+i])
		p := PerlinNoiseFromXoroshiro(copyX)
		p.amplitude = amplitudes[i] * persist
		p.lacunarity = lacuna
		o = append(o, p)
		i++
		n++
		lacuna *= 2
		persist *= 0.5
	}
	return
}

func (o *Octave) Sample(x, y, z float64) (v float64) {
	for _, n := range *o {
		v += n.Sample(x*n.lacunarity, y*n.lacunarity, z*n.lacunarity, 0, 0) * n.amplitude
	}
	return
}
