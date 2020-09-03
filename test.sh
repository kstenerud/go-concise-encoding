#!/bin/bash

realpath() {
    [[ $1 = /* ]] && echo "$1" || echo "$PWD/${1#./}"
}

cd "$(dirname "$(realpath "$0")")"

echo
echo "Test normally"
echo "-------------"
go test ./...

echo
echo "Test purego"
echo "-----------"
go test -tags purego ./...

echo
echo "Test 32-bit"
echo "-----------"
GOARCH=386 go test ./...
