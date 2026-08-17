package anvil_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"git.anygate.cloud/anygatecloud/strom/mc/level"
	anvil2 "git.anygate.cloud/anygatecloud/strom/mc/level/anvil"
	"git.anygate.cloud/anygatecloud/strom/mc/util"
)

func TestLevel(t *testing.T) {
	f, err := os.Open("../../../.devres/anvil_testdata/testworld/level.dat")
	if err != nil {
		t.Skipf("testdata not available: %v", err)
	}
	defer f.Close()

	l, err := anvil2.ReadLevelData(f)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(l)
}

func TestRegionFile(t *testing.T) {
	f, err := os.Open("../../../.devres/anvil_testdata/testworld/region/r.0.0.mca")
	if err != nil {
		t.Skipf("testdata not available: %v", err)
	}
	defer f.Close()
	ch, err := anvil2.ReadChunkHeader(f)
	if err != nil {
		t.Fatal(err)
	}

	n, err := ch.ChunkAt(f, 0, 0)
	if err != nil {
		return
	}

	fmt.Println(string(util.MustT(json.MarshalIndent(n, "", "  ")))[100:])
	fmt.Println("Chunk loaded successfully")
}

func loadPacketToChunk(b []byte) (chunk *level.Chunk, err error) {
	buffer := bytes.NewBuffer(b)
	chunk, err = level.ReadChunkFromChunkPacketData(buffer, "1.21.11", 384)
	return
}
