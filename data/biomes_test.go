package data_test

import (
	"testing"

	"github.com/admin-else/strom/data"
)

func TestLookupBiomeById(t *testing.T) {
	b, ok := data.LookupBiomeById("1.21.8", 10)
	if !ok {
		t.Fatal("Biome not found")
	}
	t.Logf("%v", b)
}
