package anvil_test

import (
	"os"
	"testing"

	"github.com/admin-else/strom/level/anvil"
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
