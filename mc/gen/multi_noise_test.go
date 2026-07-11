package gen

import (
	"testing"
)

func TestLoadMultiNoise(t *testing.T) {
	for _, version := range []string{"1.21.8", "1.21.9", "1.21.11"} {
		data, err := LoadMultiNoise(version)
		if err != nil {
			t.Fatalf("LoadMultiNoise(%q): %v", version, err)
		}
		if len(data.Presets[MultiNoiseOverworld]) == 0 {
			t.Errorf("%s: overworld preset empty", version)
		}
		if len(data.Presets[MultiNoiseNether]) == 0 {
			t.Errorf("%s: nether preset empty", version)
		}
	}
}

func TestMultiNoiseLookup_NetherOrigin(t *testing.T) {
	data, err := LoadMultiNoise("1.21.11")
	if err != nil {
		t.Fatalf("LoadMultiNoise: %v", err)
	}
	biome, err := data.Lookup(MultiNoiseNether, 0, 0, 0, 0, 0, 0)
	if err != nil {
		t.Fatalf("Lookup: %v", err)
	}
	if biome != "minecraft:nether_wastes" {
		t.Errorf("nether origin = %s, want minecraft:nether_wastes", biome)
	}
}

func TestBiomeNoise_SampleBiome(t *testing.T) {
	var worldSeed int64 = 70388202460762
	x := XoroshiroFromSeed(worldSeed)
	x.Mix()
	b := BiomeNoiseFromXoroshiro(x, false)
	data, err := LoadMultiNoise("1.21.11")
	if err != nil {
		t.Fatalf("LoadMultiNoise: %v", err)
	}

	biome, err := b.SampleBiome(-1, -16, 0, data, MultiNoiseOverworld)
	if err != nil {
		t.Fatalf("SampleBiome: %v", err)
	}

	expected := "minecraft:bamboo_jungle"
	if biome != expected {
		t.Errorf("SampleBiome(-1, -16, 0) = %s, want %s", biome, expected)
	}
}
