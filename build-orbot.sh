#!/bin/sh

rm -f ../OrbotLib.aar
rm -f ../OrbotLib-sources.jar

if [ -d IPtProxy ]; then
  cd IPtProxy
  git clean -f
  if [ -d lyrebird ]; then
    cd lyrebird
    git reset --hard
    cd ..
  fi
  if [ -d snowflake ]; then
    cd snowflake
    git reset --hard
    cd ..
  fi
  cd ..
  git submodule update --init --recursive
  cd IPtProxy
else
  git submodule update --init --recursive
  cd IPtProxy || exit 1
fi


printf '\n\n--- Apply patches to lyrebird and snowflake...\n'
patch --directory=lyrebird --strip=1 < lyrebird.patch
patch --directory=snowflake --strip=1 < snowflake.patch

cd ..
cp OrbotTun.go/* IPtProxy/IPtProxy.go/


TEMPDIR="$TMPDIR/IPtProxy"
printf '\n\n--- Prepare build environment at %s...\n' "$TEMPDIR"
cd IPtProxy
CURRENT=$PWD
rm -rf "$TEMPDIR"
mkdir -p "$TEMPDIR"
cp -a . "$TEMPDIR/"

printf '\n\n--- Compile %s...\n' "$OUTPUT"
export PATH=~/go/bin:$PATH
cd "$TEMPDIR/IPtProxy.go" || exit 1
gomobile init 
gomobile bind -o "OrbotLib.aar" -target=android -androidapi=19 -v -trimpath
cp -v OrbotLib* "$CURRENT/../.."
printf '\n\nDone\n'
