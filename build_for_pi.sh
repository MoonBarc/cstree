#!/bin/bash

USERNAME=lawson
PI_HOST=cstree.local
PI_PROJECT_LOCATION="cstree"
TARBALL=pi.tar.xz

echo "compiling..."
# zig = best
CC="zig cc -target aarch64-linux-gnu" CGO_ENABLED=1 GOARCH=arm64 GOOS=linux go build .

echo "uploading..."
tar czf $TARBALL --exclude=$TARBALL .
scp $TARBALL $USERNAME@$PI_HOST:$PI_PROJECT_LOCATION
ssh $USERNAME@$PI_HOST "cd $PI_PROJECT_LOCATION; rm -rf ./build; mkdir build; cd build; tar xf ../$TARBALL"

echo "restarting server..."
ssh $USERNAME@$PI_HOST "sudo systemctl restart cstree"

rm ./treed
rm ./pi.tar.xz