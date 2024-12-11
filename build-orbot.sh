#!/bin/sh

rm -f ../OrbotLib.aar
rm -f ../OrbotLib-sources.jar

# should match Orbot's...
export MIN_ANDROID_SDK=23

if [ -d IPtProxy ]; then
  cd IPtProxy
  git clean -fdx
  git reset --hard
  cd ..
  git submodule update --init --recursive
else
  git submodule update --init --recursive
fi

cp OrbotTun.go/* IPtProxy/IPtProxy.go/

if test -d "$TMPDIR"; then
    :
elif test -d "$TMP"; then
    TMPDIR=$TMP
elif test -d /var/tmp; then
    TMPDIR=/var/tmp
else
    TMPDIR=/tmp
fi

TEMPDIR="$TMPDIR/IPtProxy"
printf '\n\n--- Prepare build environment at %s...\n' "$TEMPDIR"
cd IPtProxy/IPtProxy.go
go mod tidy
go get golang.org/x/mobile/bind
cd ..
CURRENT=$PWD
rm -rf "$TEMPDIR"
mkdir -p "$TEMPDIR"
cp -a . "$TEMPDIR/"

printf '\n\n--- Compile %s...\n' "$OUTPUT"
export PATH=~/go/bin:$PATH
cd "$TEMPDIR/IPtProxy.go" || exit 1
gomobile init 
gomobile bind -o "OrbotLib.aar" -ldflags="-w -s -checklinkname=0" -target=android -androidapi="$MIN_ANDROID_SDK" -v -trimpath
cp -v OrbotLib* "$CURRENT/../.."
printf '\n\nDone\n'
