#!/bin/sh

set -eu

GENERATED_FILE=../../internal/common/unicode-generated.go

if [ $# != 1 ]; then
	echo "Error: Missing UCD XML file"
	echo "generate.sh generates $GENERATED_FILE"
	echo
	echo "Usage: $0 /path/to/ucd.all.flat.xml"
	echo "Extract it from https://www.unicode.org/Public/UCD/latest/ucdxml/ucd.all.flat.zip"
	exit 1
fi

SRCPATH="$1"

go build && ./unicode-gen "$SRCPATH" >"$GENERATED_FILE"
