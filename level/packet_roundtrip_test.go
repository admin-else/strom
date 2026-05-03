package level_test

import (
	"bytes"
	"testing"

	"github.com/admin-else/strom/level"
)

func FuzzReadChunkFromChunkPacketData(f *testing.F) {
	f.Fuzz(func(t *testing.T, raw []byte) {
		chunk, err := level.ReadChunkFromChunkPacketData(bytes.NewReader(raw), "1.21.11", 384)
		if err != nil {
			t.Skipf("decode failed: %v", err)
			return
		}

		b := bytes.NewBuffer(nil)
		err = chunk.WriteChunkData(b)
		if err != nil {
			t.Errorf("Failed to write chunk data: %v", err)
			return
		}
		chunk2, err := level.ReadChunkFromChunkPacketData(bytes.NewReader(b.Bytes()), "1.21.11", 384)
		if err != nil {
			t.Errorf("Failed to decode chunk data: %v", err)
			return
		}
		if !chunk.Equals(chunk2) {
			t.Errorf("Chunks are not equal")
			return
		}
	})
}
