#!/bin/bash

set -eu

realpath() {
    [[ $1 = /* ]] && echo "$1" || echo "$PWD/${1#./}"
}

cd "$(dirname "$(realpath "$0")")"

echo
echo "Test normally"
echo "-------------"
go test -coverprofile=coverage.out ./...

echo
echo "Benchmark"
echo "---------"

go test -run Benchmark* -bench Benchmark* -benchmem -memprofile memprofile.out -cpuprofile cpuprofile.out

echo
echo "Test purego"
echo "-----------"
go test -tags purego ./...

echo
echo "Test 32-bit"
echo "-----------"
GOARCH=386 go test ./...

echo
echo "To see coverage: go tool cover -html=coverage.out"
echo "To see profile:  go tool pprof cpuprofile.out"
