#!/usr/bin/sh
set -xeuo pipefail

# Go get / install does not clone git submodules so we have to do it like this

cd "$(dirname "$0")"

rm -rf minecraft-data/
git clone --depth 1 git@github.com:admin-else/minecraft-data.git
mv minecraft-data minecraft-data-untrimmed

mkdir minecraft-data
mkdir minecraft-data/pc

mv minecraft-data-untrimmed/data/dataPaths.json minecraft-data

copy-mc() {
  mv minecraft-data-untrimmed/data/pc/$1 minecraft-data/pc
}

copy-mc common

# Keep in sync with the version list in mc/proto_generator/main.go
VERSIONS="1.8 1.9 1.9.2 1.9.4 1.10 1.10.1 1.10.2 1.11 1.11.2 1.12 1.12.1 1.12.2 1.13 1.13.1 1.13.2 1.14 1.14.1 1.14.3 1.14.4 1.15 1.15.1 1.15.2 1.16 1.16.1 1.16.2 1.16.3 1.16.4 1.16.5 1.17 1.17.1 1.18 1.18.1 1.18.2 1.19 1.19.2 1.19.3 1.19.4 1.20 1.20.1 1.20.2 1.20.3 1.20.4 1.20.5 1.20.6 1.21 1.21.1 1.21.3 1.21.4 1.21.5 1.21.6 1.21.8 1.21.9 1.21.11 26.1 26.2"
for v in $VERSIONS; do
  copy-mc "$v"
done

rm -rf minecraft-data-untrimmed
