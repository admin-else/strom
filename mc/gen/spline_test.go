package gen

import "testing"

func TestSpline(t *testing.T) {
	v := BiomeDepthSpline.Sample([]float32{0.236413047, 0.220234469, 0.139604524, -0.339296669})
	if v != 0.025035858154296875 {
		t.Errorf("bad value: %f", v)
	}
}
