#!/bin/sh

TARGET=android
OUTPUT=IPtProxy.aar

git submodule update --init --recursive

cd IPtProxy
git clean -f
cd obfs4 || exit 1
git reset --hard
cd ../snowflake || exit 1
git reset --hard
cd ..

printf '\n\n--- Apply patches to Obfs4proxy and Snowflake...\n'
patch --directory=obfs4 --strip=1 < obfs4.patch
patch --directory=snowflake --strip=1 < snowflake.patch

cd ..
cp OrbotTun.go/* IPtProxy/IPtProxy.go/
cd IPtProxy/IPtProxy.go
printf '\n\n--- Compile %s...\n' "$OUTPUT"
export PATH=~/go/bin:$PATH
gomobile init 
gomobile bind -o ../../OrbotLib.aar -v
printf '\n\nDone\n'