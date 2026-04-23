package anvil_test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/admin-else/strom/level/anvil"
	"github.com/admin-else/strom/util"
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
