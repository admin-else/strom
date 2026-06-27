#!/usr/bin/sh
set -xeuo pipefail

# Go get / install does not clone git submodules so we have to do it like this

rm -rf minecraft-data/
git clone --depth 1 git@github.com:PrismarineJS/minecraft-data.git
mv minecraft-data minecraft-data-untrimmed

mkdir minecraft-data
mkdir minecraft-data/pc

mv minecraft-data-untrimmed/data/dataPaths.json minecraft-data

copy-mc() {
  mv minecraft-data-untrimmed/data/pc/$1 minecraft-data/pc
}

copy-mc common

copy-mc 1.21.8
copy-mc 1.21.9
copy-mc 1.21.11

rm -rf minecraft-data-untrimmed
