#!/usr/bin/env sh

# Go get / install does not clone git submodules so we have to do it like this

rm -rf minecraft-data/
git clone --depth 1 git@github.com:PrismarineJS/minecraft-data.git
mv minecraft-data/data data
rm -rf minecraft-data
mv data minecraft-data
rm -rf minecraft-data/bedrock
