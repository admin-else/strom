package gen

import (
	"testing"

	"github.com/admin-else/strom/mc/data"
)

func TestBiomeNoise_SampleBiome(t *testing.T) {
	var worldSeed int64 = 70388202460762
	x := XoroshiroFromSeed(worldSeed)
	x.Mix()
	b := BiomeNoiseFromXoroshiro(x, false)
	mnd, err := data.LoadMultiNoise("1.21.11")
	if err != nil {
		t.Fatalf("LoadMultiNoise: %v", err)
	}

	biome, err := b.SampleBiome(-1, -16, 0, mnd, data.MultiNoiseOverworld)
	if err != nil {
		t.Fatalf("SampleBiome: %v", err)
	}

	expected := "minecraft:bamboo_jungle"
	if biome != expected {
		t.Errorf("SampleBiome(-1, -16, 0) = %s, want %s", biome, expected)
	}
}
