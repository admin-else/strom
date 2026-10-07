package level

import "testing"

// TestPackLight checks the light coords match LightCoordsUtil.pack and its
// FULL_SKY/FULL_BRIGHT constants.
func TestPackLight(t *testing.T) {
	if got := PackLight(7, 13); got != 7<<4|13<<20 {
		t.Errorf("PackLight(7, 13) = %d, want %d", got, 7<<4|13<<20)
	}
	if got := PackLight(15, 15); got != 15728880 {
		t.Errorf("PackLight(15, 15) = %d, want FULL_BRIGHT (15728880)", got)
	}
	if got := PackLight(0, 15); got != 15728640 {
		t.Errorf("PackLight(0, 15) = %d, want FULL_SKY (15728640)", got)
	}
}
