package data_test

import (
	"fmt"
	"testing"

	"github.com/admin-else/strom/data"
)

func TestLookupBlockByStateId(t *testing.T) {
	b, ok := data.LookupBlockByStateId("1.21.8", 0)
	if !ok {
		t.Fatal("Block not found")
	}
	fmt.Println(b)
}

func TestBlockStates(t *testing.T) {
	var stateId int32 = 8301
	b, stateMap, err := data.FromBlockState("1.21.8", stateId)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Logf("Block state map: %v, block: %v", stateMap, b.Name)
	stateId2, err := b.IdFromStateMap(stateMap)
	if err != nil {
		t.Fatal(err)
		return
	}
	if stateId2 != stateId {
		t.Fatal("state id mismatch")
	}
}
