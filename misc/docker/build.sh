#!/bin/sh

mkdir -p target/static
cp -r ../static/* target/static
go build -o target/hashfile 'github.com/PieterD/warp/cmd/hashfile'
GOOS=linux go build -o target/contentserver 'github.com/PieterD/warp/cmd/contentserver'
mkdir -p target/static/binary
GOOS=js GOARCH=wasm go build -o target/static/binary/raw 'github.com/PieterD/warp/cmd/gltest'
HASH=`./target/hashfile target/static/binary/raw`
mv target/static/binary/raw target/static/binary/$HASH
sed -i -e "s/ZZZ-ZZZ.ZZZ!ZZZ.ZZZ-ZZZ/$HASH/" target/static/index.html