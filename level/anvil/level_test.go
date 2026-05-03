package anvil_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/admin-else/strom/level"
	"github.com/admin-else/strom/level/anvil"
	"github.com/admin-else/strom/nbt"
	"github.com/admin-else/strom/util"
	"github.com/google/go-cmp/cmp"
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

func TestConverToPacket(t *testing.T) {
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
	var chunk anvil.Chunk
	err = nbt.Format.Decode(n.Value, &chunk)
	if err != nil {
		t.Errorf("Failed to decode chunk: %v", err)
		return
	}
	var packetChunk *level.Chunk
	packetChunk, err = chunk.ToStorage("1.21.11")
	if err != nil {
		t.Errorf("Failed to convert chunk: %v", err)
		return
	}
	buffer := bytes.NewBuffer(nil)
	err = packetChunk.WriteChunkData(buffer)
	if err != nil {
		t.Errorf("Failed to write chunk data: %v", err)
		return
	}
	a := chunk00.ChunkData.Val
	b := buffer.Bytes()
	if !bytes.Equal(a, b) {
		t.Errorf("Chunk data is different: %v", cmp.Diff(a, b))
	}

}
