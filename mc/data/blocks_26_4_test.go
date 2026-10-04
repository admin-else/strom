package data_test

import (
	"maps"
	"testing"

	"github.com/admin-else/strom/mc/data"
)

// TestBlocks26_4AgainstServerReport locks the 26.4-snapshot-2 block state ids and
// property mapping to values taken from the vanilla data generator report
// (generated/reports/blocks.json).
func TestBlocks26_4AgainstServerReport(t *testing.T) {
	const version = "26.4-snapshot-2"
	cases := []struct {
		stateId int32
		name    string
		props   map[string]string
	}{
		{6865, "diamond_block", nil},
		{88, "bedrock", nil},
		{89, "water", map[string]string{"level": "0"}},
		{104, "water", map[string]string{"level": "15"}},
		{8533, "birch_pressure_plate", map[string]string{"powered": "true"}},
		{8534, "birch_pressure_plate", map[string]string{"powered": "false"}},
		{5463, "oak_stairs", map[string]string{"facing": "north", "half": "top", "shape": "straight", "waterlogged": "true"}},
		{5474, "oak_stairs", map[string]string{"facing": "north", "half": "bottom", "shape": "straight", "waterlogged": "false"}},
		{8707, "repeater", map[string]string{"delay": "1", "facing": "north", "locked": "false", "powered": "false"}},
		{8767, "repeater", map[string]string{"delay": "4", "facing": "east", "locked": "false", "powered": "false"}},
		{8596, "snow", map[string]string{"layers": "8"}},
		{18631, "sea_pickle", map[string]string{"pickles": "4", "waterlogged": "false"}},
	}
	for _, c := range cases {
		block, props, err := data.FromBlockState(version, c.stateId)
		if err != nil {
			t.Errorf("FromBlockState(%d): %v", c.stateId, err)
			continue
		}
		if block.Name != c.name {
			t.Errorf("state %d = %s, want %s", c.stateId, block.Name, c.name)
		}
		if !maps.Equal(props, c.props) {
			t.Errorf("state %d %s properties = %v, want %v", c.stateId, c.name, props, c.props)
		}
		roundTrip, err := block.IdFromStateMap(props)
		if err != nil {
			t.Errorf("IdFromStateMap(%d): %v", c.stateId, err)
			continue
		}
		if roundTrip != c.stateId {
			t.Errorf("state %d round-trips to %d", c.stateId, roundTrip)
		}
	}
}
