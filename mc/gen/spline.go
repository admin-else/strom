package gen

type Spline interface {
	Sample(t []float32) float32
}

type FixedSpline struct {
	Val float32
}

func (s FixedSpline) Sample(t []float32) float32 {
	return s.Val
}

type SubSpline struct {
	Loc, Der float32
	Val      Spline
}
type SplineA struct {
	TargetIndex int
	Points      []SubSpline
}

func ternery[T any](a bool, b T, c T) T {
	if a {
		return b
	}
	return c
}

func (s SplineA) Sample(in []float32) float32 {
	f := in[s.TargetIndex]

	i := 0
	for _, p := range s.Points {
		if f <= p.Loc {
			break
		}
		i++
	}
	if i == 0 || i == len(s.Points) {
		if i != 0 {
			i--
		}
		v := s.Points[i].Val.Sample(in)
		return v + s.Points[i].Der*(f-s.Points[i].Loc)
	}
	s1 := s.Points[i-1]
	s2 := s.Points[i]
	g := s.Points[i-1].Loc
	h := s.Points[i].Loc
	k := (f - g) / (h - g)
	l := s.Points[i-1].Der
	m := s.Points[i].Der
	n := s1.Val.Sample(in)
	o := s2.Val.Sample(in)
	p := l*(h-g) - (o - n)
	q := -m*(h-g) + (o - n)
	r := lerp32(k, n, o) + k*(1.0-k)*lerp32(k, p, q)
	return r
}

const (
	splineIndexContinentalNess = iota
	splineIndexErosion
	splineIndexWeirdness
)

func getOffsetValue(weirdness float32, continentalness float32) float32 {
	f0 := 1.0 - (1.0-continentalness)*0.5
	f1 := 0.5 * (1.0 - continentalness)
	f2 := (weirdness + 1.17) * 0.46082947
	off := f2*f0 - f1
	if weirdness < -0.7 {
		if off > -0.2222 {
			return off
		}
		return -0.2222
	}
	if off > 0 {
		return off
	}
	return 0

}

func createSpline38219(f float32, bl bool) Spline {
	s := SplineA{
		TargetIndex: splineIndexWeirdness,
		Points:      nil,
	}
	i := getOffsetValue(-1.0, f)
	k := getOffsetValue(1.0, f)
	l := 1.0 - (1.0-f)*0.5
	u := 0.5 * (1.0 - f)
	l = u/(0.46082947*l) - 1.17

	if -0.65 < l && l < 1.0 {
		u = getOffsetValue(-0.65, f)
		p := getOffsetValue(-0.75, f)
		q := (p - i) * 4.0
		r := getOffsetValue(l, f)
		z := (k - r) / (1.0 - l)

		s.Points = append(s.Points, SubSpline{-1.0, q, FixedSpline{i}})
		s.Points = append(s.Points, SubSpline{-0.75, 0, FixedSpline{p}})
		s.Points = append(s.Points, SubSpline{-0.65, 0, FixedSpline{u}})
		s.Points = append(s.Points, SubSpline{l - 0.01, 0, FixedSpline{r}})
		s.Points = append(s.Points, SubSpline{l, z, FixedSpline{r}})
		s.Points = append(s.Points, SubSpline{1.0, z, FixedSpline{k}})
	} else {
		u = (k - i) * 0.5
		if bl {
			s.Points = append(s.Points, SubSpline{-1.0, 0, FixedSpline{ternery(i > 0.2, i, 0.2)}})
			s.Points = append(s.Points, SubSpline{0.0, u, FixedSpline{lerp32(0.5, i, k)}})
		} else {
			s.Points = append(s.Points, SubSpline{-1.0, u, FixedSpline{i}})
		}
		s.Points = append(s.Points, SubSpline{1.0, u, FixedSpline{k}})
	}
	return s
}

func createFlatOffsetSpline(f, g, h, i, j, k float32) Spline {
	l := 0.5 * (g - f)
	if l < k {
		l = k
	}
	m := 5.0 * (h - g)
	n := m
	if l < m {
		n = l
	}

	return SplineA{
		TargetIndex: splineIndexWeirdness,
		Points: []SubSpline{
			{-1, l, FixedSpline{f}},
			{-0.4, n, FixedSpline{g}},
			{0, m, FixedSpline{h}},
			{0.4, 2.0 * (i - h), FixedSpline{i}},
			{1, 0.7 * (j - i), FixedSpline{j}},
		},
	}
}

func landSpline(f, g, h, i, j, k float32, bl bool) Spline {
	s1 := createSpline38219(lerp32(i, 0.6, 1.5), bl)
	s2 := createSpline38219(lerp32(i, 0.6, 1.0), bl)
	s3 := createSpline38219(i, bl)
	ih := 0.5 * i
	s4 := createFlatOffsetSpline(f-0.15, ih, ih, ih, i*0.6, 0.5)
	s5 := createFlatOffsetSpline(f, j*i, g*i, ih, i*0.6, 0.5)
	s6 := createFlatOffsetSpline(f, j, j, g, h, 0.5)
	s7 := createFlatOffsetSpline(f, j, j, g, h, 0.5)
	s8 := SplineA{
		TargetIndex: splineIndexWeirdness,
		Points: []SubSpline{
			{-1, 0, FixedSpline{f}},
			{-0.4, 0, s6},
			{0, 0, FixedSpline{h + 0.07}},
		},
	}
	s9 := createFlatOffsetSpline(-0.02, k, k, g, h, 0)
	s := SplineA{
		TargetIndex: splineIndexErosion,
		Points: []SubSpline{
			{-0.85, 0, s1},
			{-0.7, 0, s2},
			{-0.4, 0, s3},
			{-0.35, 0, s4},
			{-0.1, 0, s5},
			{0.2, 0, s6},
		},
	}
	if bl {
		s.Points = append(s.Points, SubSpline{0.4, 0, s7})
		s.Points = append(s.Points, SubSpline{0.45, 0, s8})
		s.Points = append(s.Points, SubSpline{0.55, 0, s8})
		s.Points = append(s.Points, SubSpline{0.58, 0, s7})
	}
	s.Points = append(s.Points, SubSpline{0.7, 0, s9})
	return s
}

func CreateBiomeDepthSpline() Spline {
	sp1 := landSpline(-0.15, 0.00, 0.0, 0.1, 0.00, -0.03, false)
	sp2 := landSpline(-0.10, 0.03, 0.1, 0.1, 0.01, -0.03, false)
	sp3 := landSpline(-0.10, 0.03, 0.1, 0.7, 0.01, -0.03, true)
	sp4 := landSpline(-0.05, 0.03, 0.1, 1.0, 0.01, 0.01, true)

	return SplineA{
		TargetIndex: splineIndexContinentalNess,
		Points: []SubSpline{
			{Loc: -1.10, Val: FixedSpline{Val: 0.044}, Der: 0.0},
			{Loc: -1.02, Val: FixedSpline{Val: -0.2222}, Der: 0.0},
			{Loc: -0.51, Val: FixedSpline{Val: -0.2222}, Der: 0.0},
			{Loc: -0.44, Val: FixedSpline{Val: -0.12}, Der: 0.0},
			{Loc: -0.18, Val: FixedSpline{Val: -0.12}, Der: 0.0},
			{Loc: -0.16, Val: sp1, Der: 0.0},
			{Loc: -0.15, Val: sp1, Der: 0.0},
			{Loc: -0.10, Val: sp2, Der: 0.0},
			{Loc: 0.25, Val: sp3, Der: 0.0},
			{Loc: 1.00, Val: sp4, Der: 0.0},
		},
	}
}

var BiomeDepthSpline = CreateBiomeDepthSpline()
