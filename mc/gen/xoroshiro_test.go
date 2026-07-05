package gen

import (
	"testing"
)

func TestXoroshiroFromString(t *testing.T) {
	x := XoroshiroFromString("minecraft:offset")
	if x.L != 0x080518cf6af25384 || x.H != 0x3f3dfb40a54febd5 {
		t.Errorf("bad L: %x, H: %x", x.L, x.H)
	}
}

func TestSeedChain(t *testing.T) {
	var worldSeed int64 = 70388202460762
	x := XoroshiroFromSeed(worldSeed)
	x.XorString("minecraft:entities/wither_skeleton")
	x.Mix()
	i := 0
	for {
		_ = x.NextBound(1)
		_ = x.NextBound(2)
		if x.NextFloat() < 0.025 {
			i++
		} else {
			break
		}
	}
	if i != 7 {
		t.Errorf("expected 7, got %d", i)
	}
}

func TestBiome(t *testing.T) {
	var worldSeed int64 = 70388202460762
	x := XoroshiroFromSeed(worldSeed)
	x.Mix()
	b := BiomeNoiseFromXoroshiro(x, false)
	_, _, _, _, _, _, combinedData, _ := b.Sample(-1, -16, 0)

	expected := [6]int64{3238, 3521, 2364, 1396, 0, 2202}
	if combinedData != expected {
		t.Errorf("combinedData = %v, want %v", combinedData, expected)
	}
}
