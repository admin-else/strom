package gen

import "math"

type PerlinNoise struct {
	d                                                        [257]byte
	h2                                                       byte
	xoffset, yoffset, zoffset, amplitude, lacunarity, d2, t2 float64
}

func indexedLerp(idx uint8, a, b, c float64) float64 {
	switch idx & 0xf {
	case 0:
		return a + b
	case 1:
		return -a + b
	case 2:
		return a - b
	case 3:
		return -a - b
	case 4:
		return a + c
	case 5:
		return -a + c
	case 6:
		return a - c
	case 7:
		return -a - c
	case 8:
		return b + c
	case 9:
		return -b + c
	case 10:
		return b - c
	case 11:
		return -b - c
	case 12:
		return a + b
	case 13:
		return -b + c
	case 14:
		return -a + b
	case 15:
		return -b - c
	}
	panic("unreachable")
}

func lerp(part, from, to float64) float64 {
	return from + part*(to-from)
}

func (n PerlinNoise) Sample(x, y, z float64, yamp, ymin float64) (v float64) {
	var h1, h2, h3 uint8
	var t1, t2, t3 float64

	if y == 0 {
		y = n.d2
		h2 = n.h2
		t2 = n.t2
	} else {
		y += n.yoffset
		yFloored := math.Floor(y)
		y -= yFloored
		h2 = uint8(int(yFloored))
		t2 = y * y * y * (y*(y*6.0-15.0) + 10.0)
	}

	x += n.xoffset
	z += n.zoffset

	xFloored := math.Floor(x)
	zFloored := math.Floor(z)
	x -= xFloored
	z -= zFloored

	h1 = uint8(int(xFloored))
	h3 = uint8(int(zFloored))

	t1 = x * x * x * (x*(x*6.0-15.0) + 10.0)
	t3 = z * z * z * (z*(z*6.0-15.0) + 10.0)

	if yamp != 0 {
		y -= math.Floor(math.Min(y, ymin)/yamp) * yamp
	}

	a1 := n.d[h1] + h2
	b1 := n.d[h1+1] + h2

	a2 := n.d[a1] + h3
	b2 := n.d[b1] + h3
	a3 := n.d[a1+1] + h3
	b3 := n.d[b1+1] + h3

	l1 := indexedLerp(n.d[a2], x, y, z)
	l2 := indexedLerp(n.d[b2], x-1, y, z)
	l3 := indexedLerp(n.d[a3], x, y-1, z)
	l4 := indexedLerp(n.d[b3], x-1, y-1, z)
	l5 := indexedLerp(n.d[a2+1], x, y, z-1)
	l6 := indexedLerp(n.d[b2+1], x-1, y, z-1)
	l7 := indexedLerp(n.d[a3+1], x, y-1, z-1)
	l8 := indexedLerp(n.d[b3+1], x-1, y-1, z-1)

	l1 = lerp(t1, l1, l2)
	l3 = lerp(t1, l3, l4)
	l5 = lerp(t1, l5, l6)
	l7 = lerp(t1, l7, l8)

	l1 = lerp(t2, l1, l3)
	l5 = lerp(t2, l5, l7)

	return lerp(t3, l1, l5)
}

func PerlinNoiseFromXoroshiro(x Xoroshiro) (p PerlinNoise) {
	p.xoffset = x.NextFloat64() * 256.0
	p.yoffset = x.NextFloat64() * 256.0
	p.zoffset = x.NextFloat64() * 256.0
	p.amplitude = 1
	p.lacunarity = 1
	for i := range 256 {
		p.d[i] = uint8(i)
	}
	for i := range int32(256) {
		j := x.NextBound(uint32(256-i)) + i
		n := p.d[i]
		p.d[i] = p.d[j]
		p.d[j] = n
	}
	p.d[256] = p.d[0]
	i2 := math.Floor(p.yoffset)
	d2 := p.yoffset - i2
	p.h2 = uint8(int(i2))
	p.d2 = d2
	p.t2 = d2 * d2 * d2 * (d2*(d2*6.0-15.0) + 10.0)
	return
}
