package client

import (
	"testing"

	"github.com/admin-else/strom/mc/nbt"
	"github.com/admin-else/strom/mc/proto_generated/v1_21_8"
	"github.com/admin-else/strom/mc/registry"
)

func TestConfigIgnorerCapturesRegistryData(t *testing.T) {
	store := registry.NewStore()
	ci := &ConfigIgnorer{registries: store}

	packet := &v1_21_8.ConfigurationToClientPacketRegistryData{Id: "minecraft:worldgen/biome"}
	packet.Entries = append(packet.Entries, struct {
		Key   string
		Value *nbt.Anon
	}{Key: "minecraft:plains"})
	err := ci.OnRegistryData(packet)
	if err != nil {
		t.Fatalf("OnRegistryData: %v", err)
	}

	if got, ok := store.RegistryID("minecraft:worldgen/biome", "minecraft:plains"); !ok || got != 0 {
		t.Fatalf("RegistryID = %d, %v; want 0, true", got, ok)
	}
}

func TestConfigIgnorerWithoutStoreIgnoresRegistryData(t *testing.T) {
	ci := &ConfigIgnorer{}
	err := ci.OnRegistryData(&v1_21_8.ConfigurationToClientPacketRegistryData{Id: "minecraft:worldgen/biome"})
	if err != nil {
		t.Fatalf("OnRegistryData with nil store: %v", err)
	}
}
