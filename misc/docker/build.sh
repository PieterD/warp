#!/bin/sh -x

rm -rf target
mkdir -p target/static
go build -o target/hashfile 'github.com/PieterD/warp/cmd/hashfile'
go build -o target/genindex 'github.com/PieterD/warp/cmd/genindex'
GOOS=linux go build -o target/contentserver 'github.com/PieterD/warp/cmd/contentserver'
cp ../../app/gltest/static/* target/static
GOOS=js GOARCH=wasm go build -o target/static/_binary 'github.com/PieterD/warp/app/gltest'
cp "$GOROOT/misc/wasm/wasm_exec.js" 'target/static/_wasm_exec.js'
./target/genindex > target/static/index.html
