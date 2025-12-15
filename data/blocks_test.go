package data_test

import (
	"fmt"
	"testing"

	"github.com/admin-else/queser/data"
)

func TestLookupBlockByStateId(t *testing.T) {
	b, ok := data.LookupBlockByStateId("1.21.8", 0)
	if !ok {
		t.Fatal("Block not found")
	}
	fmt.Println(b)
}
