package data

import (
	"math"
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

func TestMultiNoiseLookup_BVHMatchesBruteForce(t *testing.T) {
	data, err := LoadMultiNoise("1.21.11")
	if err != nil {
		t.Fatalf("LoadMultiNoise: %v", err)
	}

	brute := func(preset string, q point6) (string, int64) {
		best := result{dist: int64(math.MaxInt64)}
		for _, p := range data.Presets[preset] {
			d := q.dist(p)
			if d < best.dist {
				best.dist = d
				best.biome = p.Biome
			}
		}
		return best.biome, best.dist
	}

	presets := []string{MultiNoiseOverworld, MultiNoiseNether}
	for _, preset := range presets {
		for i := 0; i < 1000; i++ {
			q := point6{
				t: int64(i%20000 - 10000),
				h: int64((i*7)%20000 - 10000),
				c: int64((i*13)%20000 - 10000),
				e: int64((i*17)%20000 - 10000),
				d: int64((i*3)%20000 - 10000),
				w: int64((i*11)%20000 - 10000),
			}
			wantBiome, wantDist := brute(preset, q)
			gotBiome, err := data.Lookup(preset, q.t, q.h, q.c, q.e, q.d, q.w)
			if err != nil {
				t.Fatalf("Lookup(%s): %v", preset, err)
			}
			if gotBiome != wantBiome {
				t.Errorf("preset=%s query=%+v got=%s want=%s dist=%d", preset, q, gotBiome, wantBiome, wantDist)
			}
		}
	}
}
