#!/bin/bash

set -e 

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
echo "working dir $DIR"
mkdir -p $DIR/dist

os=$(go env GOOS)
arch=$(go env GOARCH)
version=$(cat $DIR/version.go | grep "const VERSION" | awk '{print $NF}' | sed 's/"//g')
goversion=$(go version | awk '{print $3}')

echo "... running tests"
./test.sh || exit 1

for os in linux ; do
    echo "... building v$version for $os/$arch"
    BUILD=$(mktemp -d ${TMPDIR:-/tmp}/oauth2_proxy.XXXXXX)
    TARGET="oauth2_proxy-$version.$os-$arch.$goversion"
    GOOS=$os GOARCH=$arch CGO_ENABLED=0 go build -o $BUILD/$TARGET/oauth2_proxy || exit 1
    mv $BUILD/$TARGET/oauth2_proxy $DIR/dist/
done
