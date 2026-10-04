package data_test

import (
	"testing"

	"github.com/admin-else/strom/mc/data"
)

// TestItems26_4AgainstServerReport checks item ids against the vanilla
// registries.json report (minecraft:item registry).
func TestItems26_4AgainstServerReport(t *testing.T) {
	const version = "26.4-snapshot-2"
	cases := map[string]int32{
		"air":           0,
		"stone":         1,
		"sulfur":        26,
		"cinnabar":      40,
		"diamond_block": 129,
		"oak_stairs":    515,
		"diamond":       1012,
	}
	for name, id := range cases {
		item, ok := data.LookupItemByName(version, name)
		if !ok {
			t.Errorf("item %s not found", name)
			continue
		}
		if item.Id != id {
			t.Errorf("item %s id = %d, want %d", name, item.Id, id)
		}
		byId, ok := data.LookupItemById(version, id)
		if !ok || byId.Name != name {
			t.Errorf("LookupItemById(%d) = %v, want %s", id, byId, name)
		}
	}
}
