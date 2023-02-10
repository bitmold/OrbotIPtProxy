#!/bin/sh

rm -f OrbotLib.aar
rm -f OrbotLib-sources.jar

if [ -d IPtProxy ]; then
  cd IPtProxy
  git clean -f
fi

if [ -d obfs4 ]; then
  cd obfs4 || exit 1
  git reset --hard
  cd ..
fi

if [ -d snowflake ]; then
  cd snowflake || exit 1
  git reset --hard
  cd ..
fi

git submodule update --init --recursive

printf '\n\n--- Apply patches to Obfs4proxy and Snowflake...\n'
patch --directory=obfs4 --strip=1 < obfs4.patch
patch --directory=snowflake --strip=1 < snowflake.patch

cd ..
cp OrbotTun.go/* IPtProxy/IPtProxy.go/ -v


cd IPtProxy/IPtProxy.go
printf '\n\n--- Compile %s...\n' "$OUTPUT"
export PATH=~/go/bin:$PATH
gomobile init 
gomobile bind -o ../../OrbotLib.aar -v
printf '\n\nDone\n'
