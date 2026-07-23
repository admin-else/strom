#!/usr/bin/sh
set -xeuo pipefail

# Copy minecraft-data from the local AnyGate checkout instead of cloning upstream.
# This is used when upstream does not yet have the target version or when we want
# to test local changes before pushing.

LOCAL_MINECRAFT_DATA="${LOCAL_MINECRAFT_DATA:-$HOME/src/anygate/minecraft-data}"

rm -rf minecraft-data/
mkdir -p minecraft-data/pc

cp "$LOCAL_MINECRAFT_DATA/data/dataPaths.json" minecraft-data/

copy-mc() {
  cp -r "$LOCAL_MINECRAFT_DATA/data/pc/$1" minecraft-data/pc
}

copy-mc common

copy-mc 1.21.8
copy-mc 1.21.9
copy-mc 1.21.11
copy-mc 26.1
copy-mc 26.2
