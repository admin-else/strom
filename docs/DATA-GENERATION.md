# Minecraft data generation (strom)

`strom` vendors PrismarineJS **minecraft-data** under
`mc/data/minecraft-data/pc/<version>/`. The generator in `mc/proto_generator/`
turns `protocol.json` into the typed wire packages under
`mc/proto_generated/`. This document records where the data comes from and how
to regenerate it without guessing.

## 1. Rule: the jar is the source of truth, not the wiki

Prismarine's wiki (`minecraft-data`) is community-maintained and is frequently
stale for fresh snapshots. Do **not** derive protocol fields, registry ids, or
particle types from it when the vanilla jar disagrees. Use, in order:

1. **Prismarine's official generator** when it supports the version:
   <https://github.com/PrismarineJS/minecraft-data-generator>. It is a Fabric
   mod that runs the integrated server and emits complete, authoritative files
   (`blocks`, `items`, `biomes`, `blockCollisionShapes`, `entities`,
   `particles`, `tints`, ...). Supported versions are listed in the repo's
   `versions.json`. This is the preferred path.
2. **The vanilla jar's own data generator reports** for anything the generator
   cannot emit yet:
   ```bash
   java -DbundlerMainClass=net.minecraft.data.Main \
     -jar ~/.gradle/caches/fabric-loom/<version>/minecraft-server.jar \
     --reports out
   # -> out/generated/reports/{registries,blocks,packets,commands}.json
   ```
   `registries.json` is authoritative for every built-in registry's
   `name -> protocol_id` mapping (blocks, items, particle types, ...), and
   `packets.json` maps packet names to protocol ids.
3. **The decompiled client source** for stream codec layouts (packet field
   order, particle option fields), diffed against the previous version. The
   neon repo carries the decompile under `reference/<version>/client-src/`;
   the upstream method is a decompiler diff of
   `net/minecraft/network/protocol/**` and
   `net/minecraft/core/particles/**`.

Never hand-transcribe ids or field orders from an older release's wiki entry.

## 2. Version support

As of this writing the Prismarine generator lists versions up to `26.1`.
`26.4-snapshot-2` is **not** supported yet, so strom builds it with stopgap
scripts driven by the vanilla reports plus a decompiler diff (section 3). When
the generator gains `26.4-snapshot-2`, regenerate the whole
`pc/26.4-snapshot-2/` directory with it and delete the stopgap scripts.

## 3. Pipeline for `26.4-snapshot-2`

Run from `strom/`, with `REPORTS` pointing at the extracted
`out/generated/reports`:

```bash
REPORTS=/home/adman/mc-test/server-26.4/generated/reports

# skeleton: packet ids from the vanilla reports, packet set/fields from a
# decompiler diff of the 26.2 and 26.4 jars
python3 mc/data/patch_protocol_26_4_snapshot_2.py \
  mc/data/minecraft-data/pc/26.2/protocol.json \
  <26.2-packets.json> <26.4-packets.json> \
  mc/data/minecraft-data/pc/26.4-snapshot-2/protocol.json

# jar-derived particle registry + option codecs + world_particles layout
python3 mc/data/patch_particles_26_4_snapshot_2.py \
  $REPORTS/registries.json \
  mc/data/minecraft-data/pc/26.4-snapshot-2/protocol.json

# blocks / items from the same reports
python3 mc/data/gen_blocks_26_4_snapshot_2.py $REPORTS/blocks.json <base> <out>
python3 mc/data/gen_items_26_4_snapshot_2.py  $REPORTS/registries.json <base> <out>

# rebuild the typed wire packages
go run ./mc/proto_generator
go test ./mc/data/ ./mc/proto/ ./mc/proto_generator/
```

`patch_particles_26_4_snapshot_2.py` is idempotent and only rewrites the
`Particle` type, the `ItemStackTemplate` type and the `packet_world_particles`
definition.

## 4. Case study: the 26.4 particle protocol

**Symptom.** A live `world_particles` (`ClientboundLevelParticlesPacket`)
decoded as `proto.UnCodablePacket` with `index not found`, while
`entity_metadata` decoded fine.

**Root cause.** The wiki-derived `26.4` protocol carried an older release's
particle data:

- the `Particle` mapper ids were from an earlier version (`flame=32`,
  `cloud=4`), but the 26.4 registry has `flame=39`, `cloud=11` and seven new
  particles (`sulfur_bubbles`, `noxious_gas`, `noxious_gas_cloud`, `geyser`,
  `geyser_base`, `geyser_poof`, `geyser_plume`);
- the `data` switch had no `default`, so simple particles (a
  `SimpleParticleType` writes no data) could not be skipped;
- option layouts had changed: `DustParticleOptions` is now
  `INT color, FLOAT scale` (was four floats),
  `DustColorTransitionOptions` is `INT fromColor, INT toColor, FLOAT scale`
  (was seven floats), `TrailParticleOption` added a `VAR_INT duration`;
- `packet_world_particles` had the old field order
  (`longDistance, alwaysShow, x, y, z, offset..., amount, particle`) instead of
  the real `particle, overrideLimiter, alwaysShow, x, y, z, xDist..., count,
  randomizationType`.

**Fix.** `patch_particles_26_4_snapshot_2.py` regenerates all of it from the
jar: registry ids from `registries.json`, option codecs from
`net/minecraft/core/particles/*.java`, and the packet layout from
`ClientboundLevelParticlesPacket.java`. After regenerating, a live 26.4
`world_particles` (flame and dust) decodes cleanly.

## 5. Other gotchas

- **`dataPaths.json` aliasing**: a version entry pointing at another version's
  directory makes local edits no-ops. Always point at `pc/<version>`.
- **Integer block-state properties** are absolute values with a non-zero
  minimum for some blocks (repeater `delay` 1..4, snow `layers` 1..8).
  Parsing the mixed-radix digit must add `MinValue`.
- **Property order** (mixed-radix significance) is not stored in the vanilla
  block report; derive it from the state-id sequence.
- The biome registry is dynamic and network-sent, so it needs the generator.

## 6. References

- minecraft-data: <https://github.com/PrismarineJS/minecraft-data>
- add-data-new-version: <https://github.com/PrismarineJS/minecraft-data/blob/master/doc/add-data-new-version.md>
- minecraft-data-generator: <https://github.com/PrismarineJS/minecraft-data-generator>
- minecraft-data-generator-server: <https://github.com/PrismarineJS/minecraft-data-generator-server>
