package data_test

import (
	"testing"

	"github.com/admin-else/strom/mc/data"
)

func TestLookupBiomeById(t *testing.T) {
	b, ok := data.LookupBiomeById("1.21.8", 10)
	if !ok {
		t.Fatal("Biome not found")
	}
	t.Logf("%v", b)
}

// TestBiomePrecipitationAndTemperature checks the snake_case has_precipitation
// JSON key is decoded (it does not match the Go field name by case folding).
func TestBiomePrecipitationAndTemperature(t *testing.T) {
	b, ok := data.LookupBiomeByName("26.4-snapshot-2", "plains")
	if !ok {
		t.Fatal("plains not found")
	}
	if !b.HasPrecipitation {
		t.Error("plains has_precipitation = false, want true")
	}
	if b.Temperature != 0.8 {
		t.Errorf("plains temperature = %v, want 0.8", b.Temperature)
	}
}
