#!/usr/bin/sh
set -xeuo pipefail

# Copy minecraft-data from the local AnyGate checkout instead of cloning upstream.
# This is used when upstream does not yet have the target version or when we want
# to test local changes before pushing.

LOCAL_MINECRAFT_DATA="${LOCAL_MINECRAFT_DATA:-$HOME/src/anygate/minecraft-data}"

cd "$(dirname "$0")"

rm -rf minecraft-data/
mkdir -p minecraft-data/pc

cp "$LOCAL_MINECRAFT_DATA/data/dataPaths.json" minecraft-data/

copy-mc() {
  cp -r "$LOCAL_MINECRAFT_DATA/data/pc/$1" minecraft-data/pc
}

copy-mc common

# Keep in sync with the version list in mc/proto_generator/main.go
VERSIONS="1.8 1.9 1.9.2 1.9.4 1.10 1.10.1 1.10.2 1.11 1.11.2 1.12 1.12.1 1.12.2 1.13 1.13.1 1.13.2 1.14 1.14.1 1.14.3 1.14.4 1.15 1.15.1 1.15.2 1.16 1.16.1 1.16.2 1.16.3 1.16.4 1.16.5 1.17 1.17.1 1.18 1.18.1 1.18.2 1.19 1.19.2 1.19.3 1.19.4 1.20 1.20.1 1.20.2 1.20.3 1.20.4 1.20.5 1.20.6 1.21 1.21.1 1.21.3 1.21.4 1.21.5 1.21.6 1.21.8 1.21.9 1.21.11 26.1 26.2"
for v in $VERSIONS; do
  copy-mc "$v"
done
