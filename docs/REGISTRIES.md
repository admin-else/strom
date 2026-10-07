# Configuration registries (client `RegistryData`)

The server sends its **dynamic** registries during the configuration phase as
`ClientboundRegistryDataPacket` (`registry_data`, config id 7). strom now
decodes them while it acknowledges the configuration phase, so a consumer can
resolve wire ids to resource ids instead of hardcoding them.

## API

`mc/registry` holds the store (no proto dependency):

```go
type Registry struct{ /* name, id -> resource id, resource id -> id */ }
func (r *Registry) Name() string
func (r *Registry) Add(resourceID string) (id int)
func (r *Registry) Len() int
func (r *Registry) ID(resourceID string) (id int, ok bool)
func (r *Registry) ResourceID(id int) (resourceID string, ok bool)
func (r *Registry) Entries() []string // id order, a copy

type Store struct{ /* keyed by registry name */ }
func NewStore() *Store
func (s *Store) Set(r *Registry)
func (s *Store) Registry(name string) (r *Registry, ok bool)
func (s *Store) Names() []string // sorted
func (s *Store) RegistryID(registryName, resourceID string) (id int, ok bool)
func (s *Store) ResourceID(registryName string, id int) (resourceID string, ok bool)
```

The wire lists entries in id order, so `Registry` index == id. The optional
per-element NBT payload is decoded (to consume it) but not stored: this package
only carries id maps, not element data. Rendering-only element fields (biome
colours, special effects) stay in neon.

Capture at connect time:

```go
store := registry.NewStore()
conn, err := client.Login(addr, account,
    client.WithVersion("26.4-snapshot-2"),
    client.WithRegistries(store),
)
```

`client.IgnoreConfigStore(c, store)` is the underlying config handler; a nil
store keeps the old ignore-only behaviour (`client.IgnoreConfig`).

Query through the shared world:

```go
w := world.NewWorld(conn.Version, -64, 384)
w.SetRegistries(store)
id, ok := w.RegistryID("minecraft:worldgen/biome", "minecraft:plains")
name, ok := w.ResourceID("minecraft:worldgen/biome", id)
```

Probe / verification harness:

```bash
go run ./cmd/strom registry -addr 127.0.0.1:25565 -version 26.4-snapshot-2 -acc dev
go run ./cmd/strom registry -registry minecraft:worldgen/biome   # full id dump
go run ./cmd/strom registry -json                                # every registry
```

## What the 26.4 server sends

32 registries (exactly `RegistryDataLoader.SYNCHRONIZED_REGISTRIES`), including
`minecraft:worldgen/biome`. On the local 26.4-snapshot-2 server the biome
registry has **67** entries in wire order, matching neon's captured
`biome/registry_order.json`:

- `minecraft:dappled_forest` = **8**
- `minecraft:sulfur_caves` = **54**

`cmd/strom registry` asserts both ids and exits non-zero on mismatch.

## `minecraft:entity_type` is NOT a dynamic registry

`RegistryData` only carries `SYNCHRONIZED_REGISTRIES`. `entity_type` (and the
other built-in registries: block, item, fluid, particle_type, sound_event, ...)
lives in the version's **static** built-in registry, so it is never sent over
the wire. The 32 captured registries are exactly the synchronized set; there is
no `entity_type` packet.

Consequences for the neon follow-up:

- Do **not** source entity_type ids from `RegistryData`; they are fixed by the
  protocol version.
- Do **not** source them from minecraft-data's `entities.json` either: its `id`
  is not the wire id (26.4: `zombie=150`, `player=155`).
- The wire ids are the `BuiltInRegistries.ENTITY_TYPE` registration order, i.e.
  the declaration order in
  `reference/26.4-snapshot-2/.../world/entity/EntityTypes.java`:
  `cow=30`, `zombie=154`, `player=159`. This matches neon's current constants.
  Empirically confirmed against the live server with `/summon`
  (`Cow` -> type 30, `Zombie` -> type 154).

So neon can delete `registry_order.json` and query biomes through the store,
but should keep the entity_type constants (or regenerate them from the
`EntityTypes.java` order) — they cannot come from configuration registries.
