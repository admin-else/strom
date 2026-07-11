package data

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
