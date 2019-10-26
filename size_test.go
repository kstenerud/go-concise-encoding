package cbe

import (
	"net/url"
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
	assertMarshaledSize(t, ContainerTypeNone, 0xffffffffff, 7)
	assertMarshaledSize(t, ContainerTypeNone, -0xffffffffff, 7)
	assertMarshaledSize(t, ContainerTypeNone, 0xffffffffffff, 8)
	assertMarshaledSize(t, ContainerTypeNone, -0xffffffffffff, 8)
	assertMarshaledSize(t, ContainerTypeNone, 0xfffffffffffff, 9)
	assertMarshaledSize(t, ContainerTypeNone, -0xfffffffffffff, 9)
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
	location, err := time.LoadLocation("America/Vancouver")
	if err != nil {
		t.Fatal(err)
	}
	date := time.Date(2055, time.Month(2), 14, 8, 30, 0, 55, location)
	assertMarshaledSize(t, ContainerTypeNone, date, 22)
}

func TestSizeString(t *testing.T) {
	assertMarshaledSize(t, ContainerTypeNone, "", 1)
	assertMarshaledSize(t, ContainerTypeNone, "1", 2)
	assertMarshaledSize(t, ContainerTypeNone, "12", 3)
	assertMarshaledSize(t, ContainerTypeNone, "123456789abcdef", 16)
	assertMarshaledSize(t, ContainerTypeNone, "123456789abcdefg", 18)
	assertMarshaledSize(t, ContainerTypeNone, generateString(2000), 2003)
}

func TestSizeURI(t *testing.T) {
	testURL, err := url.Parse("http://example.com")
	if err != nil {
		t.Fatal(err)
	}
	assertMarshaledSize(t, ContainerTypeNone, testURL, 20)
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
