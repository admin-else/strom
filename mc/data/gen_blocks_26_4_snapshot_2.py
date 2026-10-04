#!/usr/bin/env python3
"""Regenerate the 26.4-snapshot-2 blocks.json from the vanilla data generator report.

PrismarineJS minecraft-data does not publish 26.4 snapshots yet, and the
vendored 26.4-snapshot-2/blocks.json was a copy of 26.2 with stale global state
ids. This script rebuilds it from the authoritative vanilla report so that
LookupBlockByStateId/FromBlockState agree with the server.

Reproduce the input report with the vanilla 26.4 server jar:

    java -DbundlerMainClass=net.minecraft.data.Main \\
      -jar server.jar --reports out

The block report lands in out/generated/reports/blocks.json.

Usage:
    gen_blocks_26_4_snapshot_2.py <report blocks.json> <base blocks.json> <out blocks.json>

`base blocks.json` supplies non-state metadata (displayName, hardness, ...) for
blocks that already existed; state ids and property definitions always come from
the report.
"""
import itertools
import json
import re
import sys

REPORT, BASE, OUT = sys.argv[1:4]


def find_property_order(props, states, min_id):
    """Return property names most-significant-first for the report's id encoding.

    The report encodes each state as a mixed-radix number over the block's
    properties, with the least significant property varying fastest. The order
    is not recorded, so try permutations until the decomposition reproduces the
    report's id -> properties mapping for every state.
    """
    names = list(props.keys())

    def matches(order):
        k = len(order)
        weights = [1] * k
        for i in range(k - 2, -1, -1):
            weights[i] = weights[i + 1] * len(props[order[i + 1]])
        for state in states:
            digits = state["id"] - min_id
            values = state.get("properties", {})
            for i, name in enumerate(order):
                index = (digits // weights[i]) % len(props[name])
                if props[name][index] != values.get(name):
                    return False
        return True

    for order in itertools.permutations(names):
        if matches(order):
            return list(order)
    raise SystemExit("no property order reproduces the report states")


def build_states(props, order):
    states = []
    for name in order:
        values = props[name]
        if values == ["true", "false"]:
            states.append({"name": name, "type": "bool", "num_values": 2})
            continue
        is_int = all(re.fullmatch(r"-?\d+", v) for v in values)
        states.append({
            "name": name,
            "type": "int" if is_int else "enum",
            "num_values": len(values),
            "values": values,
        })
    return states


def display_name(name):
    return name.replace("_", " ").title()


def main():
    report = json.load(open(REPORT))
    base = json.load(open(BASE))
    by_name = {b["name"]: b for b in base}
    next_id = max(b["id"] for b in base) + 1

    out = []
    for full_name, entry in report.items():
        name = full_name.split(":", 1)[1]
        props = entry.get("properties", {})
        states = sorted(entry["states"], key=lambda s: s["id"])
        min_id = states[0]["id"]
        max_id = states[-1]["id"]
        default_id = next((s["id"] for s in states if s.get("default")), min_id)
        order = find_property_order(props, states, min_id) if props else []

        block = dict(by_name.get(name, {}))
        if "id" not in block:
            block["id"] = next_id
            next_id += 1
        block["name"] = name
        block.setdefault("displayName", display_name(name))
        block.setdefault("hardness", 0.0)
        block.setdefault("resistance", 0.0)
        block.setdefault("stackSize", 64)
        block.setdefault("diggable", True)
        block.setdefault("material", "default")
        block.setdefault("transparent", False)
        block.setdefault("emitLight", 0)
        block.setdefault("filterLight", 15)
        block.setdefault("drops", [])
        block.setdefault("boundingBox", "block")
        block["defaultState"] = default_id
        block["minStateId"] = min_id
        block["maxStateId"] = max_id
        block["states"] = build_states(props, order)
        out.append(block)

    json.dump(out, open(OUT, "w"), indent=2)
    print("wrote", OUT, "blocks:", len(out))


if __name__ == "__main__":
    main()
