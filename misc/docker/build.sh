#!/bin/sh -e

readonly USAGE="usage: $0 <appname> [buildargs]"

function die { echo "$*" 1>&2 ; exit 1; }
cd "$(dirname "$0")" || die "failed to find base dir"

APPNAME=$1
shift || die "Missing argument: app name"
INDEXARGS=$1
shift || INDEXARGS=""
shift && die "excess arguments. $USAGE"

echo "Regenerating target"
rm -rf target
mkdir -p target/static

echo "Building tools"
go build -o target/hashfile 'github.com/PieterD/warp/cmd/hashfile'
go build -o target/genindex 'github.com/PieterD/warp/cmd/genindex'
GOOS=linux go build -o target/contentserver 'github.com/PieterD/warp/cmd/contentserver'

echo "Copying ${APPNAME} statics"
cp ../../app/${APPNAME}/static/* target/static

echo "Building ${APPNAME} WASM binary"
GOOS=js GOARCH=wasm go build -o target/static/_binary 'github.com/PieterD/warp/app/'${APPNAME}

echo "Copying ${GOROOT}/misc/wasm/wasm_exec.js"
cp "$GOROOT/misc/wasm/wasm_exec.js" 'target/static/_wasm_exec.js'

echo "Generating index with '${INDEXARGS}'"
./target/genindex ${INDEXARGS} > target/static/index.html

echo "Application ${APPNAME} is ready!"