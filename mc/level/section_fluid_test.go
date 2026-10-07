package level_test

import (
	"bytes"
	"testing"

	"github.com/admin-else/strom/mc/level"
)

func TestSectionFluidCountRoundTrip(t *testing.T) {
	for _, version := range []string{"26.2", "26.4-snapshot-2"} {
		section := level.Section{BlockCount: 3, FluidCount: 7}
		section.Blocks, _ = level.MakeBlockFormat(version).FullWith(0)
		section.Biomes, _ = level.MakeBiomeFormat(version).FullWith(0)

		var buffer bytes.Buffer
		if err := level.SectionEncodePacketVersion(&buffer, section, version); err != nil {
			t.Fatalf("%s encode: %v", version, err)
		}
		got, err := level.SectionDecodePacket(&buffer, version)
		if err != nil {
			t.Fatalf("%s decode: %v", version, err)
		}
		if got.FluidCount != 7 {
			t.Fatalf("%s fluid count = %d, want 7", version, got.FluidCount)
		}
	}
}
