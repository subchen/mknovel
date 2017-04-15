#!/bin/sh

set -e
ROOT=$(cd $(dirname $0); pwd)
DIST=$ROOT/bin
NAME=mknovel
VERSION=$(cat VERSION.txt)

mkdir -p $DIST

ldflags="-s -w"
ldflags="$ldflags -X 'main.BuildVersion=$(git rev-list HEAD --count)'"
ldflags="$ldflags -X 'main.BuildGitCommit=$(git describe --abbrev=0 --always)'"
ldflags="$ldflags -X 'main.BuildDate=$(date -u -R)'"
ldflags="$ldflags -X 'main.VERSION=$VERSION'"

# build and zip
echo "building for linux"
cd $ROOT && GOOS=linux GOARCH=amd64 go build -ldflags "$ldflags" -o $DIST/$NAME-linux-$VERSION

echo "building for darwin"
cd $ROOT && GOOS=darwin GOARCH=amd64 go build -ldflags "$ldflags" -o $DIST/$NAME-darwin-$VERSION

echo "building for windows"
cd $ROOT && GOOS=windows GOARCH=amd64 go build -ldflags "$ldflags" -o $DIST/$NAME-$VERSION.exe
