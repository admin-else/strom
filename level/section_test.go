package level_test

import (
	"bytes"
	"testing"

	"github.com/admin-else/strom/level"
)

func FuzzSectionDeEnCodePacket(f *testing.F) {
	f.Fuzz(func(t *testing.T, raw []byte) {
		r := bytes.NewReader(raw)
		s1, err := level.SectionDecodePacket(r, TestVersion)
		if err != nil {
			t.Skipf("decode failed: %v", err)
			return
		}

		buf := bytes.NewBuffer(nil)
		err = level.SectionEncodePacket(buf, s1)
		if err != nil {
			t.Fatal(err)
			return
		}

		s2, err := level.SectionDecodePacket(buf, TestVersion)
		if err != nil {
			t.Fatal(err)
			return
		}

		if s1.BlockCount != s2.BlockCount {
			t.Error("BlockCount mismatch")
		}
		if level.CompareStorage(s1.Blocks, s2.Blocks) {
			t.Error("Blocks mismatch")
		}
		if level.CompareStorage(s1.Biomes, s2.Biomes) {
			t.Error("Biomes mismatch")
		}

	})
}
