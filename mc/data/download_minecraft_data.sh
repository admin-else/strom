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

# Legacy versions
copy-mc 1.8
copy-mc 1.12.2
copy-mc 1.14.4

# Modern versions
copy-mc 1.21.8
copy-mc 1.21.9
copy-mc 1.21.11
copy-mc 26.1
copy-mc 26.2

rm -rf minecraft-data-untrimmed
