#!/usr/bin/env python3
"""Rebuild the 26.4-snapshot-2 particle protocol from the vanilla jar.

The PrismarineJS minecraft-data generator does not support `26.4-snapshot-2`
yet, so the protocol schema is derived from `26.2` by
`patch_protocol_26_4_snapshot_2.py`. That derivation leaves the particle data
stale: it carries the particle registry ids and the `particleData`/`Particle`
option layouts of an older release, and the `world_particles` packet keeps the
old field order. The result is that a live `world_particles` packet fails to
decode (`index not found` on the particle mapper).

Do **not** patch this from the wiki. Derive it from the jar instead:

  * particle registry ids come from the vanilla data generator report
    `registries.json` (`minecraft:particle_type`), which mirrors
    `net/minecraft/core/particles/ParticleTypes.java` registration order;
  * particle option codecs are transcribed from the decompiled
    `net/minecraft/core/particles/*.java` `streamCodec` methods;
  * the `world_particles` field order comes from the decompiled
    `net/minecraft/network/protocol/game/ClientboundLevelParticlesPacket.java`.

The vanilla report is produced with:

    java -DbundlerMainClass=net.minecraft.data.Main \\
      -jar ~/.gradle/caches/fabric-loom/26.4-snapshot-2/minecraft-server.jar \\
      --reports out
    # report: out/generated/reports/registries.json

Usage:
    patch_particles_26_4_snapshot_2.py <registries.json> <protocol.json>

The protocol file is patched in place.
"""
import json
import sys

REGISTRIES, PROTOCOL = sys.argv[1:3]


def field(name, typ):
    return {"name": name, "type": typ}


# Particle options whose type is not a SimpleParticleType. Each entry is the
# order of fields written by that particle's `streamCodec` in the decompiled
# net/minecraft/core/particles/*.java. SimpleParticleType particles write no
# data and fall through to the switch default (void).
PARTICLE_OPTIONS = {
    # BlockParticleOption.streamCodec -> Block.BLOCK_STATE_REGISTRY_STREAM_CODEC
    "block": "varint",
    "block_crumble": "varint",
    "block_marker": "varint",
    "dust_pillar": "varint",
    "falling_dust": "varint",
    # DustParticleOptions.STREAM_CODEC = INT color, FLOAT scale
    "dust": ["container", [field("color", "i32"), field("scale", "f32")]],
    # DustColorTransitionOptions.STREAM_CODEC = INT fromColor, INT toColor, FLOAT scale
    "dust_color_transition": [
        "container",
        [
            field("fromColor", "i32"),
            field("toColor", "i32"),
            field("scale", "f32"),
        ],
    ],
    # PowerParticleOption.streamCodec = FLOAT power
    "dragon_breath": "f32",
    # SpellParticleOption.streamCodec = INT color, FLOAT power
    "effect": ["container", [field("color", "i32"), field("power", "f32")]],
    "instant_effect": ["container", [field("color", "i32"), field("power", "f32")]],
    # ColorParticleOption.streamCodec = INT color
    "entity_effect": "i32",
    "tinted_leaves": "i32",
    "flash": "i32",
    # ItemParticleOption.streamCodec -> ItemStackTemplate.STREAM_CODEC
    "item": "ItemStackTemplate",
    # SculkChargeParticleOptions.STREAM_CODEC = FLOAT roll
    "sculk_charge": "f32",
    # ShriekParticleOption.STREAM_CODEC = VAR_INT delay
    "shriek": "varint",
    # VibrationParticleOption.STREAM_CODEC = PositionSource, VAR_INT arrivalInTicks
    "vibration": [
        "container",
        [
            {
                "name": "positionType",
                "type": [
                    "mapper",
                    {"type": "varint", "mappings": {"0": "block", "1": "entity"}},
                ],
            },
            {
                "name": "position",
                "type": [
                    "switch",
                    {
                        "compareTo": "positionType",
                        "fields": {
                            "block": "position",
                            "entity": [
                                "container",
                                [
                                    field("entityId", "varint"),
                                    field("entityEyeHeight", "f32"),
                                ],
                            ],
                        },
                    },
                ],
            },
            field("ticks", "varint"),
        ],
    ],
    # TrailParticleOption.STREAM_CODEC = Vec3, INT color, VAR_INT duration
    "trail": [
        "container",
        [
            field("target", "vec3f64"),
            field("color", "i32"),
            field("duration", "varint"),
        ],
    ],
    # GeyserParticleOptions.streamCodec = INT waterBlocks
    "geyser": ["container", [field("waterBlocks", "i32")]],
    "geyser_plume": ["container", [field("waterBlocks", "i32")]],
    # GeyserBaseParticleOptions.streamCodec = INT waterBlocks, FLOAT burstImpulseBase
    "geyser_base": [
        "container",
        [field("waterBlocks", "i32"), field("burstImpulseBase", "f32")],
    ],
    "geyser_poof": [
        "container",
        [field("waterBlocks", "i32"), field("burstImpulseBase", "f32")],
    ],
}

# ItemStackTemplate.STREAM_CODEC = Item(varint), VAR_INT count, DataComponentPatch.
# DataComponentPatch = VAR_INT positive count, VAR_INT negative count,
# positive entries (SlotComponent), negative entries (SlotComponentType only).
ITEM_STACK_TEMPLATE = [
    "container",
    [
        field("itemId", "varint"),
        field("count", "varint"),
        field("addedComponentCount", "varint"),
        field("removedComponentCount", "varint"),
        field(
            "components",
            ["array", {"count": "addedComponentCount", "type": "SlotComponent"}],
        ),
        field(
            "removeComponents",
            [
                "array",
                {
                    "count": "removedComponentCount",
                    "type": ["container", [field("type", "SlotComponentType")]],
                },
            ],
        ),
    ],
]

# ClientboundLevelParticlesPacket.STREAM_CODEC field order.
WORLD_PARTICLES = [
    "container",
    [
        field("particle", "Particle"),
        field("overrideLimiter", "bool"),
        field("alwaysShow", "bool"),
        field("x", "f64"),
        field("y", "f64"),
        field("z", "f64"),
        field("xDist", "f32"),
        field("yDist", "f32"),
        field("zDist", "f32"),
        field("xMaxSpeed", "f32"),
        field("yMaxSpeed", "f32"),
        field("zMaxSpeed", "f32"),
        field("count", "varint"),
        field("randomizationType", "varint"),
    ],
]


def particle_mappings(registries):
    entries = registries["minecraft:particle_type"]["entries"]
    by_id = {}
    for full_name, entry in entries.items():
        name = full_name.split(":", 1)[1]
        by_id[entry["protocol_id"]] = name
    return {str(i): by_id[i] for i in sorted(by_id)}


def rebuild_particles(proto, registries):
    mappings = particle_mappings(registries)
    missing = [
        name
        for name in PARTICLE_OPTIONS
        if name not in mappings.values()
    ]
    if missing:
        raise SystemExit(f"particle options not present in the registry: {missing}")

    proto["types"]["ItemStackTemplate"] = ITEM_STACK_TEMPLATE
    proto["types"]["Particle"] = [
        "container",
        [
            field("type", ["mapper", {"type": "varint", "mappings": mappings}]),
            field(
                "data",
                [
                    "switch",
                    {
                        "compareTo": "type",
                        "fields": PARTICLE_OPTIONS,
                        "default": "void",
                    },
                ],
            ),
        ],
    ]
    proto["play"]["toClient"]["types"]["packet_world_particles"] = WORLD_PARTICLES


def main():
    registries = json.load(open(REGISTRIES))
    proto = json.load(open(PROTOCOL))
    rebuild_particles(proto, registries)
    json.dump(proto, open(PROTOCOL, "w"), indent=2)
    print("patched particles in", PROTOCOL)


if __name__ == "__main__":
    main()
