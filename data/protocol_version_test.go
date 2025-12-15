package data_test

import (
	"testing"

	"github.com/admin-else/strom/data"
)

func TestLookUpProtocolVersion(t *testing.T) {
	version, err := data.LookUpProtocolVersionByVersion(772)
	if err != nil {
		t.Fatalf("Expected to find version 772, got error: %v", err)
	}
	if version.MinecraftVersion != "1.21.8" {
		t.Errorf("Expected version 1.21.8, got %s", version.MinecraftVersion)
	}
	t.Logf("Found version %v", version)
}
