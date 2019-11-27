package cbe

import (
	"net/url"
	"testing"
	"time"
)

func TestSizeNil(t *testing.T) {
	assertMarshaledSize(t, InlineContainerTypeNone, nil, 1)
}

func TestSizeBool(t *testing.T) {
	assertMarshaledSize(t, InlineContainerTypeNone, true, 1)
	assertMarshaledSize(t, InlineContainerTypeNone, false, 1)
}

func TestSizeInt(t *testing.T) {
	assertMarshaledSize(t, InlineContainerTypeNone, 0, 1)
	assertMarshaledSize(t, InlineContainerTypeNone, 1, 1)
	assertMarshaledSize(t, InlineContainerTypeNone, -1, 1)
	assertMarshaledSize(t, InlineContainerTypeNone, 0xff, 2)
	assertMarshaledSize(t, InlineContainerTypeNone, -0xff, 2)
	assertMarshaledSize(t, InlineContainerTypeNone, 0xffff, 3)
	assertMarshaledSize(t, InlineContainerTypeNone, -0xffff, 3)
	assertMarshaledSize(t, InlineContainerTypeNone, 0xffffffff, 5)
	assertMarshaledSize(t, InlineContainerTypeNone, -0xffffffff, 5)
	assertMarshaledSize(t, InlineContainerTypeNone, 0xffffffffff, 7)
	assertMarshaledSize(t, InlineContainerTypeNone, -0xffffffffff, 7)
	assertMarshaledSize(t, InlineContainerTypeNone, 0xffffffffffff, 8)
	assertMarshaledSize(t, InlineContainerTypeNone, -0xffffffffffff, 8)
	assertMarshaledSize(t, InlineContainerTypeNone, 0xfffffffffffff, 9)
	assertMarshaledSize(t, InlineContainerTypeNone, -0xfffffffffffff, 9)
}

func TestSizeFloat(t *testing.T) {
	assertMarshaledSize(t, InlineContainerTypeNone, 0.0, 5)
	assertMarshaledSize(t, InlineContainerTypeNone, -0.0, 5)
	assertMarshaledSize(t, InlineContainerTypeNone, 0.1, 9)
	assertMarshaledSize(t, InlineContainerTypeNone, -0.1, 9)
	assertMarshaledSize(t, InlineContainerTypeNone, 0.25, 5)
	assertMarshaledSize(t, InlineContainerTypeNone, -0.25, 5)
}

func TestSizeTime(t *testing.T) {
	location, err := time.LoadLocation("America/Vancouver")
	if err != nil {
		t.Fatal(err)
	}
	date := time.Date(2055, time.Month(2), 14, 8, 30, 0, 55, location)
	assertMarshaledSize(t, InlineContainerTypeNone, date, 22)
}

func TestSizeString(t *testing.T) {
	assertMarshaledSize(t, InlineContainerTypeNone, "", 1)
	assertMarshaledSize(t, InlineContainerTypeNone, "1", 2)
	assertMarshaledSize(t, InlineContainerTypeNone, "12", 3)
	assertMarshaledSize(t, InlineContainerTypeNone, "123456789abcdef", 16)
	assertMarshaledSize(t, InlineContainerTypeNone, "123456789abcdefg", 18)
	assertMarshaledSize(t, InlineContainerTypeNone, generateString(2000), 2003)
}

func TestSizeURI(t *testing.T) {
	testURL, err := url.Parse("http://example.com")
	if err != nil {
		t.Fatal(err)
	}
	assertMarshaledSize(t, InlineContainerTypeNone, testURL, 20)
}

func TestSizeBytes(t *testing.T) {
	assertMarshaledSize(t, InlineContainerTypeNone, make([]byte, 0), 2)
	assertMarshaledSize(t, InlineContainerTypeNone, make([]byte, 1), 3)
	assertMarshaledSize(t, InlineContainerTypeNone, make([]byte, 2000), 2003)
}

func TestSizeList(t *testing.T) {
	assertMarshaledSize(t, InlineContainerTypeNone, make([]int, 0), 2)
	assertMarshaledSize(t, InlineContainerTypeNone, []interface{}{0, 1000}, 6)
}

func TestSizeMap(t *testing.T) {
	assertMarshaledSize(t, InlineContainerTypeNone, map[interface{}]interface{}{0: "1"}, 5)
}
