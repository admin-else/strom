package anvil_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"testing"

	"github.com/admin-else/strom/level"
	"github.com/admin-else/strom/level/anvil"
	"github.com/admin-else/strom/nbt"
	"github.com/admin-else/strom/util"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestLevel(t *testing.T) {
	f, err := os.Open("./testdata/testworld/level.dat")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	l, err := anvil.ReadLevelData(f)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(l)
}

func TestRegionFile(t *testing.T) {
	f, err := os.Open("./testdata/testworld/region/r.0.0.mca")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	ch, err := anvil.ReadChunkHeader(f)
	if err != nil {
		t.Fatal(err)
	}

	n, err := ch.ChunkAt(f, 0, 0)
	if err != nil {
		return
	}

	fmt.Println(string(util.MustT(json.MarshalIndent(n, "", "  "))))
	fmt.Println("Chunk loaded successfully")
}

func loadPacketToChunk(b []byte) (chunk *level.Chunk, err error) {
	buffer := bytes.NewBuffer(b)
	chunk, err = level.ReadChunkFromChunkPacketData(buffer, "1.21.11", 384)
	return
}

func TestConverToPacket(t *testing.T) {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	f, err := os.Open("./testdata/r.0.0.mca")
	if err != nil {
		t.Errorf("Failed to open file: %v", err)
		return
	}
	defer f.Close()
	mca, err := anvil.ReadChunkHeader(f)
	if err != nil {
		t.Errorf("Failed to read chunk header: %v", err)
		return
	}
	var n *nbt.Tag
	n, err = mca.ChunkAt(f, 0, 0)
	if err != nil {
		t.Errorf("Failed to read chunk: %v", err)
		return
	}
	var anvilChunk anvil.Chunk
	err = nbt.Format.Decode(n.Value, &anvilChunk)
	if err != nil {
		t.Errorf("Failed to decode chunk: %v", err)
		return
	}
	var anvilPacketableChunk *level.Chunk
	anvilPacketableChunk, err = anvilChunk.ToStorage("1.21.11")
	if err != nil {
		t.Errorf("Failed to convert chunk: %v", err)
		return
	}
	buffer := bytes.NewBuffer(nil)
	err = anvilPacketableChunk.WriteChunkData(buffer)
	if err != nil {
		t.Errorf("Failed to write chunk data: %v", err)
		return
	}
	a := chunk00.ChunkData.Val
	b := buffer.Bytes()
	if !bytes.Equal(a, b) {
		t.Errorf("Chunk data is different: %v", cmp.Diff(a, b))
	}

	aChunk, err := loadPacketToChunk(a)
	if err != nil {
		t.Errorf("Failed to load chunk from packet: %v", err)
		return
	}

	bChunk, err := loadPacketToChunk(b)
	if err != nil {
		t.Errorf("Failed to load chunk from file: %v", err)
		return
	}
	var storageComparer = cmp.Comparer(func(a, b *level.Storage) bool {
		if a == nil && b == nil {
			return true
		}
		if a == nil || b == nil {
			return false
		}
		return a.Equals(b)
	})

	if !aChunk.Equals(anvilPacketableChunk) {
		t.Errorf("Chunk is different: %v", cmp.Diff(aChunk, anvilPacketableChunk, storageComparer, cmpopts.IgnoreUnexported(level.Chunk{}, level.Section{}, level.Storage{})))
	}
	if !bChunk.Equals(aChunk) {
		t.Errorf("Chunk is different: %v", cmp.Diff(bChunk, aChunk, storageComparer, cmpopts.IgnoreUnexported(level.Chunk{}, level.Section{}, level.Storage{})))
	}
}
