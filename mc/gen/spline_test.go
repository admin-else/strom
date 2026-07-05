package gen

import (
	"math"
	"testing"
)

func TestSpline(t *testing.T) {
	v := BiomeDepthSpline.Sample([]float32{0.236413047, 0.139604524, -0.339296669, 0.220234469})
	var e float32 = 0.025035858154296875 - 0.015
	if math.Abs(float64(e-v)) > 1e-8 { // float impercision
		t.Errorf("bad value: %f != %f", v, e)
	}
}
