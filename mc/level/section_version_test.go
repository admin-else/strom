package level

import (
	"testing"
)

func TestSectionFormatVersionBoundaries(t *testing.T) {
	if !hasSectionFluidCount("26.2") {
		t.Fatalf("26.2 sections must carry the fluid count short")
	}
	if !hasSectionFluidCount("26.4-snapshot-2") {
		t.Fatalf("26.4 sections must carry the fluid count short")
	}
	if hasSectionFluidCount("1.21.11") {
		t.Fatalf("1.21.11 sections must not carry the fluid count short")
	}
	if got := sectionBiomeEntries("26.2"); got != BiomesPerChunkSection {
		t.Fatalf("26.2 biome entries = %d, want %d", got, BiomesPerChunkSection)
	}
	if got := sectionBiomeEntries("1.21.11"); got != BiomesPerChunkSection {
		t.Fatalf("1.21.11 biome entries = %d, want %d", got, BiomesPerChunkSection)
	}
	if got := sectionBiomeEntries("26.4-snapshot-2"); got != BlocksPerChunkSection {
		t.Fatalf("26.4 biome entries = %d, want %d", got, BlocksPerChunkSection)
	}
}

func TestMakeBpeRangeCoversAllLocalAndDirect(t *testing.T) {
	bpes := makeBpeRange(4, 15)
	for want := uint8(0); want <= 15; want++ {
		if want > 0 && want < 4 {
			continue
		}
		if !containsBpe(bpes, want) {
			t.Fatalf("bits-per-entry %d missing from %v", want, bpes)
		}
	}
	if bpes[len(bpes)-1] != 15 {
		t.Fatalf("direct bits-per-entry must be last, got %v", bpes)
	}
	biomes := makeBpeRange(0, 7)
	if biomes[len(biomes)-1] != 7 {
		t.Fatalf("biome direct bits-per-entry must be last, got %v", biomes)
	}
	if !containsBpe(biomes, 4) || !containsBpe(biomes, 5) {
		t.Fatalf("biome local palettes 4 and 5 must be accepted, got %v", biomes)
	}
}

func containsBpe(bpes []uint8, value uint8) (ret bool) {
	for _, bpe := range bpes {
		if bpe == value {
			return true
		}
	}
	return false
}
