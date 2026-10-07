#!/usr/bin/env python3
"""Build a 26.4-snapshot-2 protocol.json from the vendored 26.2 protocol.json.

The wiki (PrismarineJS minecraft-data) does not yet publish a 26.4 snapshot, so
the protocol schema is derived from the 26.2 data plus the vanilla 26.4 data
generator reports.

Reproduce the inputs:

  1. Vanilla packet reports (packet name <-> protocol id) for 26.2 and 26.4:
       java -DbundlerMainClass=net.minecraft.data.Main \\
         -jar ~/.gradle/caches/fabric-loom/<ver>/minecraft-server.jar --reports out
     The report lands in out/generated/reports/packets.json. Use the Fabric Loom
     server jars for both versions.

  2. Run this script:
       python3 mc/data/patch_protocol_26_4_snapshot_2.py \\
         mc/data/minecraft-data/pc/26.2/protocol.json \\
         <26.2>/packets.json <26.4>/packets.json \\
         mc/data/minecraft-data/pc/26.4-snapshot-2/protocol.json

Derivation rules:

  * packet ids come from the vanilla 26.4 report;
  * wiki packet names are matched to vanilla 26.4 names through the vanilla 26.2
    report ids, so only packets genuinely renamed between 26.2 and 26.4 need an
    explicit entry (VANILLA_RENAMES);
  * packet set changes (NEW_PACKETS) and field changes (FIELD_FIXES) come from a
    decompiler diff of net/minecraft/network/protocol/** between the two Loom
    client jars (see the 26.4 diff notes).

Usage: patch_protocol_26_4_snapshot_2.py <26.2 protocol.json> <26.2 packets.json> <26.4 packets.json> <out.json>

This only produces the wiki-derived skeleton. The particle registry ids,
particle option codecs and the `world_particles` field order are stale here;
run `patch_particles_26_4_snapshot_2.py <registries.json> <out.json>` afterwards
to rebuild them from the vanilla jar. See docs/DATA-GENERATION.md.
"""
import json
import sys

SRC, REP262, REP264, OUT = sys.argv[1:5]

# wiki state/direction key -> (report state key, report direction key)
KEYMAP = {
    ("handshaking", "toServer"): ("handshake", "serverbound"),
    ("handshaking", "toClient"): ("handshake", "clientbound"),
    ("status", "toServer"): ("status", "serverbound"),
    ("status", "toClient"): ("status", "clientbound"),
    ("login", "toServer"): ("login", "serverbound"),
    ("login", "toClient"): ("login", "clientbound"),
    ("configuration", "toServer"): ("configuration", "serverbound"),
    ("configuration", "toClient"): ("configuration", "clientbound"),
    ("play", "toServer"): ("play", "serverbound"),
    ("play", "toClient"): ("play", "clientbound"),
}

# vanilla 26.2 name -> vanilla 26.4 name, only for packets that were renamed.
# Everything else keeps its vanilla 26.2 name (and the wiki name is recovered from
# the wiki id, so no per-packet wiki alias table is needed).
VANILLA_RENAMES = {
    ("play", "serverbound"): {"swing": "punch"},
}

# New 26.4 packets (vanilla name), no wiki 26.2 counterpart, ordered fields.
NEW_PACKETS = {
    ("configuration", "toClient"): {
        "post_effects": [("postEffects", ["array", {"countType": "varint", "type": "string"}])],
    },
    ("play", "toClient"): {
        "add_transient_block": [("location", "position"), ("blockState", "varint")],
        "swing_animation": [
            ("entityId", "varint"),
            ("hand", ["mapper", {"type": "varint", "mappings": {"0": "main_hand", "1": "off_hand"}}]),
            ("animation", "varint"),
        ],
        "post_effects": [("postEffects", ["array", {"countType": "varint", "type": "string"}])],
    },
}

# Field-list replacements for existing wiki packet definitions (decompiler diff).
FIELD_FIXES = {
    ("handshaking", "toServer", "set_protocol"): [
        ("protocolVersion", "varint"),
        ("serverHost", "string"),
        ("serverPort", "u16"),
        ("nextState", "varint"),
    ],
    # ClientboundLoginPacket gained onlineMode before enforcesSecureChat.
    ("play", "toClient", "login"): [
        ("entityId", "i32"),
        ("isHardcore", "bool"),
        ("worldNames", ["array", {"countType": "varint", "type": "string"}]),
        ("maxPlayers", "varint"),
        ("viewDistance", "varint"),
        ("simulationDistance", "varint"),
        ("reducedDebugInfo", "bool"),
        ("enableRespawnScreen", "bool"),
        ("doLimitedCrafting", "bool"),
        ("worldState", "SpawnInfo"),
        ("onlineMode", "bool"),
        ("enforcesSecureChat", "bool"),
    ],
    # ServerboundAcceptTeleportationPacket now carries the accepted position.
    ("play", "toServer", "teleport_confirm"): [
        ("teleportId", "varint"),
        ("x", "f64"),
        ("y", "f64"),
        ("z", "f64"),
        ("yRot", "f32"),
        ("xRot", "f32"),
    ],
}

# Helper-type replacements not tied to a packet name (decompiler diff).
# CommonPlayerSpawnInfo dropped the seed long and encodes GameType as a varint.
TYPE_FIXES = {
    ("play", "toClient", "SpawnInfo"): [
        ("dimension", "varint"),
        ("name", "string"),
        (
            "gamemode",
            ["mapper", {"type": "varint", "mappings": {"0": "survival", "1": "creative", "2": "adventure", "3": "spectator"}}],
        ),
        ("previousGamemode", "varint"),
        ("isDebug", "bool"),
        ("isFlat", "bool"),
        ("death", ["option", "GlobalPos"]),
        ("portalCooldown", "varint"),
        ("seaLevel", "varint"),
    ],
}

# shared top-level defs for wiki packet names with no per-state definition
SHARED_DEF = {
    ("login", "toServer"): {"cookie_response": "packet_common_cookie_response"},
    ("login", "toClient"): {"cookie_request": "packet_common_cookie_request"},
    ("configuration", "toServer"): {
        "settings": "packet_common_settings",
        "cookie_response": "packet_common_cookie_response",
        "select_known_packs": "packet_common_select_known_packs",
        "custom_click_action": "packet_common_custom_click_action",
    },
    ("configuration", "toClient"): {"store_cookie": "packet_common_store_cookie"},
}


def wiki_to_container(fields):
    return ["container", [{"name": fn, "type": ft} for fn, ft in fields]]


def report_by_id(report, rstate, rdir):
    return {v["protocol_id"]: k.split(":")[-1] for k, v in report.get(rstate, {}).get(rdir, {}).items()}


def main():
    proto = json.load(open(SRC))
    rep262 = json.load(open(REP262))
    rep264 = json.load(open(REP264))

    for (wst, wdir), (rst, rdir) in KEYMAP.items():
        packet = proto[wst][wdir]["types"]["packet"]
        name_field = [f for f in packet[1] if f["name"] == "name"][0]
        switch_field = [f for f in packet[1] if f["name"] == "params"][0]["type"][1]["fields"]
        wiki_ids = {int(k, 0): v for k, v in name_field["type"][1]["mappings"].items()}
        vanilla262 = report_by_id(rep262, rst, rdir)
        vanilla264 = report_by_id(rep264, rst, rdir)
        renames = VANILLA_RENAMES.get((rst, rdir), {})

        # wiki 26.2 id -> vanilla 26.4 name
        id_to_vanilla264 = {}
        for wid, wname in wiki_ids.items():
            v262 = vanilla262.get(wid)
            if v262 is None:
                continue  # packet not present in 26.2 report (legacy ping etc.)
            id_to_vanilla264[wid] = renames.get(v262, v262)

        # wiki id -> wiki name, for ids that survived
        preserved = {wid: wiki_ids[wid] for wid in id_to_vanilla264}

        # map each 26.4 vanilla packet to a wiki name
        new_mappings = {}
        new_switch = {}
        used = set()
        for pid in sorted(vanilla264):
            vname = vanilla264[pid]
            wname = None
            for wid, v264 in id_to_vanilla264.items():
                if v264 == vname and wid not in used:
                    wname = preserved[wid]
                    used.add(wid)
                    break
            if wname is None:
                if vname in NEW_PACKETS.get((wst, wdir), {}):
                    fields = NEW_PACKETS[(wst, wdir)][vname]
                    proto[wst][wdir]["types"]["packet_" + vname] = wiki_to_container(fields)
                    wname = vname
                else:
                    raise SystemExit(f"unmatched 26.4 packet {rst}/{rdir}/{vname} (id {pid})")
            new_mappings[str(pid)] = wname
            if "packet_" + wname in proto[wst][wdir]["types"]:
                new_switch[wname] = "packet_" + wname
            elif wname in SHARED_DEF.get((wst, wdir), {}):
                new_switch[wname] = SHARED_DEF[(wst, wdir)][wname]
            elif "packet_common_" + wname in proto["types"]:
                new_switch[wname] = "packet_common_" + wname
            else:
                raise SystemExit(f"no wiki packet def for {wst}/{wdir}/{wname}")

        name_field["type"][1]["mappings"] = new_mappings
        switch_field.clear()
        for pid in sorted(int(k) for k in new_mappings):
            wname = new_mappings[str(pid)]
            switch_field[wname] = new_switch[wname]

    # login_success stays in the login state; the play "login" packet is separate
    # and already exists in the wiki 26.2 schema, so nothing to move.

    for (wst, wdir, wname), fields in FIELD_FIXES.items():
        proto[wst][wdir]["types"]["packet_" + wname] = wiki_to_container(fields)

    for (wst, wdir, tname), fields in TYPE_FIXES.items():
        proto[wst][wdir]["types"][tname] = wiki_to_container(fields)

    json.dump(proto, open(OUT, "w"), indent=2)
    print("wrote", OUT)


if __name__ == "__main__":
    main()
