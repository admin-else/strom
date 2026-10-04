#!/usr/bin/env python3
"""Regenerate the 26.4-snapshot-2 items.json from the vanilla registries report.

The item protocol id is the numeric id used in the slot wire format, so a stale
items.json silently mislabels every inventory slot. The report's
`minecraft:item` registry is authoritative for the name <-> id mapping.

Usage:
    gen_items_26_4_snapshot_2.py <registries.json> <base items.json> <out items.json>

The base file supplies displayName/stackSize for items that already existed.
"""
import json
import re
import sys

REGISTRIES, BASE, OUT = sys.argv[1:4]


def display_name(name):
    return name.replace("_", " ").title()


def main():
    registries = json.load(open(REGISTRIES))
    base = {i["name"]: i for i in json.load(open(BASE))}
    entries = registries["minecraft:item"]["entries"]

    out = []
    for full_name, entry in entries.items():
        name = full_name.split(":", 1)[1]
        old = base.get(name, {})
        out.append({
            "id": entry["protocol_id"],
            "name": name,
            "displayName": old.get("displayName", display_name(name)),
            "stackSize": old.get("stackSize", 64),
        })
    out.sort(key=lambda i: i["id"])
    json.dump(out, open(OUT, "w"), indent=2)
    print("wrote", OUT, "items:", len(out))


if __name__ == "__main__":
    main()
