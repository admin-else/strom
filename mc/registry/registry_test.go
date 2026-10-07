package registry

import (
	"slices"
	"testing"
)

func TestRegistryIDOrderIsWireOrder(t *testing.T) {
	r := NewRegistry("minecraft:worldgen/biome")
	for id, name := range []string{"minecraft:badlands", "minecraft:dappled_forest", "minecraft:sulfur_caves"} {
		got := r.Add(name)
		if got != id {
			t.Fatalf("Add(%q) = %d, want %d", name, got, id)
		}
	}

	if got, ok := r.ID("minecraft:dappled_forest"); !ok || got != 1 {
		t.Fatalf("ID(dappled_forest) = %d, %v; want 1, true", got, ok)
	}
	if got, ok := r.ResourceID(2); !ok || got != "minecraft:sulfur_caves" {
		t.Fatalf("ResourceID(2) = %q, %v; want sulfur_caves, true", got, ok)
	}
	if _, ok := r.ResourceID(3); ok {
		t.Fatalf("ResourceID(3) reported ok for an out-of-range id")
	}
	if _, ok := r.ID("minecraft:missing"); ok {
		t.Fatalf("ID(missing) reported ok")
	}
	if r.Len() != 3 {
		t.Fatalf("Len() = %d, want 3", r.Len())
	}
}

func TestRegistryAddDuplicateKeepsOriginalID(t *testing.T) {
	r := NewRegistry("minecraft:entity_type")
	if got := r.Add("minecraft:cow"); got != 0 {
		t.Fatalf("first Add = %d, want 0", got)
	}
	r.Add("minecraft:zombie")
	if got := r.Add("minecraft:cow"); got != 0 {
		t.Fatalf("duplicate Add = %d, want 0", got)
	}
	if r.Len() != 2 {
		t.Fatalf("Len() = %d, want 2", r.Len())
	}
}

func TestRegistryEntriesIsCopy(t *testing.T) {
	r := NewRegistry("minecraft:worldgen/biome")
	r.Add("minecraft:plains")
	entries := r.Entries()
	entries[0] = "mutated"
	if got, _ := r.ResourceID(0); got != "minecraft:plains" {
		t.Fatalf("mutating Entries() changed the registry: %q", got)
	}
}

func TestStoreLookups(t *testing.T) {
	store := NewStore()
	biomes := NewRegistry("minecraft:worldgen/biome")
	biomes.Add("minecraft:plains")
	biomes.Add("minecraft:dappled_forest")
	store.Set(biomes)
	store.Set(NewRegistry("minecraft:entity_type"))

	if got, ok := store.RegistryID("minecraft:worldgen/biome", "minecraft:dappled_forest"); !ok || got != 1 {
		t.Fatalf("RegistryID = %d, %v; want 1, true", got, ok)
	}
	if got, ok := store.ResourceID("minecraft:worldgen/biome", 0); !ok || got != "minecraft:plains" {
		t.Fatalf("ResourceID = %q, %v; want plains, true", got, ok)
	}
	if _, ok := store.RegistryID("minecraft:missing", "minecraft:plains"); ok {
		t.Fatalf("RegistryID on missing registry reported ok")
	}
	if got, ok := store.Registry("minecraft:entity_type"); !ok || got.Len() != 0 {
		t.Fatalf("Registry(entity_type) = %v, %v; want empty registry, true", got, ok)
	}

	want := []string{"minecraft:entity_type", "minecraft:worldgen/biome"}
	if got := store.Names(); !slices.Equal(got, want) {
		t.Fatalf("Names() = %v, want %v", got, want)
	}
}
