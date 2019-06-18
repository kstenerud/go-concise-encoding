package cbe

import (
	"testing"
	"time"
)

func TestSizeNil(t *testing.T) {
	assertMarshaledSize(t, ContainerTypeNone, nil, 1)
}

func TestSizeBool(t *testing.T) {
	assertMarshaledSize(t, ContainerTypeNone, true, 1)
	assertMarshaledSize(t, ContainerTypeNone, false, 1)
}

func TestSizeInt(t *testing.T) {
	assertMarshaledSize(t, ContainerTypeNone, 0, 1)
	assertMarshaledSize(t, ContainerTypeNone, 1, 1)
	assertMarshaledSize(t, ContainerTypeNone, -1, 1)
	assertMarshaledSize(t, ContainerTypeNone, 0xff, 2)
	assertMarshaledSize(t, ContainerTypeNone, -0xff, 2)
	assertMarshaledSize(t, ContainerTypeNone, 0xffff, 3)
	assertMarshaledSize(t, ContainerTypeNone, -0xffff, 3)
	assertMarshaledSize(t, ContainerTypeNone, 0xffffffff, 5)
	assertMarshaledSize(t, ContainerTypeNone, -0xffffffff, 5)
	assertMarshaledSize(t, ContainerTypeNone, 0xffffffffff, 9)
	assertMarshaledSize(t, ContainerTypeNone, -0xffffffffff, 9)
}

func TestSizeFloat(t *testing.T) {
	assertMarshaledSize(t, ContainerTypeNone, 0.0, 5)
	assertMarshaledSize(t, ContainerTypeNone, -0.0, 5)
	assertMarshaledSize(t, ContainerTypeNone, 0.1, 9)
	assertMarshaledSize(t, ContainerTypeNone, -0.1, 9)
	assertMarshaledSize(t, ContainerTypeNone, 0.25, 5)
	assertMarshaledSize(t, ContainerTypeNone, -0.25, 5)
}

func TestSizeTime(t *testing.T) {
	assertMarshaledSize(t, ContainerTypeNone, time.Now(), 9)
}

func TestSizeString(t *testing.T) {
	assertMarshaledSize(t, ContainerTypeNone, "", 1)
	assertMarshaledSize(t, ContainerTypeNone, "1", 2)
	assertMarshaledSize(t, ContainerTypeNone, "12", 3)
	assertMarshaledSize(t, ContainerTypeNone, "123456789abcdef", 16)
	assertMarshaledSize(t, ContainerTypeNone, "123456789abcdefg", 18)
	assertMarshaledSize(t, ContainerTypeNone, generateString(2000), 2003)
}

func TestSizeBytes(t *testing.T) {
	assertMarshaledSize(t, ContainerTypeNone, make([]byte, 0), 2)
	assertMarshaledSize(t, ContainerTypeNone, make([]byte, 1), 3)
	assertMarshaledSize(t, ContainerTypeNone, make([]byte, 2000), 2003)
}

func TestSizeList(t *testing.T) {
	assertMarshaledSize(t, ContainerTypeNone, make([]int, 0), 2)
	assertMarshaledSize(t, ContainerTypeNone, []interface{}{0, 1000}, 6)
}

func TestSizeMap(t *testing.T) {
	assertMarshaledSize(t, ContainerTypeNone, map[interface{}]interface{}{0: "1"}, 5)
}
