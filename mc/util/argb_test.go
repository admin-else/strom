package util

import "testing"

func TestToABGR(t *testing.T) {
	// ARGB 0x11223344 -> ABGR 0x11443322
	if got := ToABGR(0x11223344); got != 0x11443322 {
		t.Errorf("ToABGR = %#x, want %#x", got, 0x11443322)
	}
	if got := FromABGR(0x11443322); got != 0x11223344 {
		t.Errorf("FromABGR = %#x, want %#x", got, 0x11223344)
	}
}

func TestArgbChannels(t *testing.T) {
	color := ArgbColor4(0x11, 0x22, 0x33, 0x44)
	if ArgbAlpha(color) != 0x11 || ArgbRed(color) != 0x22 || ArgbGreen(color) != 0x33 || ArgbBlue(color) != 0x44 {
		t.Errorf("channels = %d %d %d %d", ArgbAlpha(color), ArgbRed(color), ArgbGreen(color), ArgbBlue(color))
	}
}
