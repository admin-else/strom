// Package registry stores the dynamic registries a server sends during the
// configuration phase. The entry order on the wire is the numeric id order
// (index == id), so a Registry is both an id -> resource-id and a
// resource-id -> id lookup.
package registry

import (
	"slices"
	"sync"
)

// Registry is a single dynamic registry, e.g. minecraft:worldgen/biome or
// minecraft:entity_type. Entries are stored in wire (id) order.
type Registry struct {
	name    string
	entries []string
	byName  map[string]int
}

// NewRegistry creates an empty registry with the given resource name.
func NewRegistry(name string) (r *Registry) {
	return &Registry{name: name, byName: make(map[string]int)}
}

// Name returns the registry's resource name, e.g. minecraft:entity_type.
func (r *Registry) Name() string { return r.name }

// Add appends a resource id at the next numeric id and returns that id.
// A duplicate resource id overwrites the earlier name lookup but keeps its
// original position, matching a registry that lists each element once.
func (r *Registry) Add(resourceID string) (id int) {
	if existing, ok := r.byName[resourceID]; ok {
		return existing
	}
	id = len(r.entries)
	r.entries = append(r.entries, resourceID)
	r.byName[resourceID] = id
	return id
}

// Len returns the number of entries.
func (r *Registry) Len() int { return len(r.entries) }

// ID returns the numeric id for a resource id.
func (r *Registry) ID(resourceID string) (id int, ok bool) {
	id, ok = r.byName[resourceID]
	return
}

// ResourceID returns the resource id for a numeric id.
func (r *Registry) ResourceID(id int) (resourceID string, ok bool) {
	if id < 0 || id >= len(r.entries) {
		return
	}
	return r.entries[id], true
}

// Entries returns the resource ids in id order. The returned slice is a copy.
func (r *Registry) Entries() (ret []string) {
	return slices.Clone(r.entries)
}

// Store holds the registries captured from a server, keyed by registry name.
type Store struct {
	mu         sync.RWMutex
	registries map[string]*Registry
}

// NewStore creates an empty registry store.
func NewStore() (s *Store) {
	return &Store{registries: make(map[string]*Registry)}
}

// Set replaces the registry stored under r.Name().
func (s *Store) Set(r *Registry) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.registries[r.Name()] = r
}

// Registry returns the registry stored under name.
func (s *Store) Registry(name string) (r *Registry, ok bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok = s.registries[name]
	return
}

// Names returns the stored registry names in sorted order.
func (s *Store) Names() (ret []string) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	ret = make([]string, 0, len(s.registries))
	for name := range s.registries {
		ret = append(ret, name)
	}
	slices.Sort(ret)
	return
}

// RegistryID returns the numeric id of resourceID inside registryName.
func (s *Store) RegistryID(registryName, resourceID string) (id int, ok bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, found := s.registries[registryName]
	if !found {
		return
	}
	return r.ID(resourceID)
}

// ResourceID returns the resource id of numeric id inside registryName.
func (s *Store) ResourceID(registryName string, id int) (resourceID string, ok bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, found := s.registries[registryName]
	if !found {
		return
	}
	return r.ResourceID(id)
}
