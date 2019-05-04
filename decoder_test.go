package cbe

import (
	"fmt"
	"testing"
	"time"
)

type Nil int

type testCallbacks struct {
	data interface{}
}

func (callbacks *testCallbacks) OnNil() error {
	callbacks.data = new(Nil)
	return nil
}

func (callbacks *testCallbacks) OnBool(value bool) error {
	callbacks.data = value
	return nil
}

func (callbacks *testCallbacks) OnInt(value int) error {
	callbacks.data = value
	return nil
}

func (callbacks *testCallbacks) OnInt64(value int64) error {
	callbacks.data = value
	return nil
}

func (callbacks *testCallbacks) OnUint(value uint) error {
	callbacks.data = value
	return nil
}

func (callbacks *testCallbacks) OnUint64(value uint64) error {
	callbacks.data = value
	return nil
}

func (callbacks *testCallbacks) OnFloat32(value float32) error {
	callbacks.data = value
	return nil
}

func (callbacks *testCallbacks) OnFloat64(value float64) error {
	callbacks.data = value
	return nil
}

func (callbacks *testCallbacks) OnTime(value time.Time) error {
	callbacks.data = value
	return nil
}

func (callbacks *testCallbacks) OnListBegin() error {
	// singleData.listBegan = true
	return nil
}

func (callbacks *testCallbacks) OnListEnd() error {
	// singleData.listEnded = true
	return nil
}

func (callbacks *testCallbacks) OnMapBegin() error {
	// singleData.mapBegan = true
	return nil
}

func (callbacks *testCallbacks) OnMapEnd() error {
	// singleData.mapEnded = true
	return nil
}

func (callbacks *testCallbacks) OnStringBegin(byteCount uint64) error {
	// singleData.stringLength = byteCount
	return nil
}

func (callbacks *testCallbacks) OnStringData(bytes []byte) error {
	// singleData.stringData = string(bytes)
	return nil
}

func (callbacks *testCallbacks) OnCommentBegin(byteCount uint64) error {
	// singleData.commentLength = byteCount
	return nil
}

func (callbacks *testCallbacks) OnCommentData(bytes []byte) error {
	// singleData.commentData = string(bytes)
	return nil
}

func (callbacks *testCallbacks) OnBinaryBegin(byteCount uint64) error {
	// singleData.binaryLength = byteCount
	return nil
}

func (callbacks *testCallbacks) OnBinaryData(bytes []byte) error {
	// singleData.binaryData = bytes
	return nil
}

func assertDecoded(t *testing.T, encoded []byte, expected interface{}) {
	callbacks := new(testCallbacks)
	decoder := NewDecoder(9, callbacks)
	err := decoder.Feed(encoded)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	actual := callbacks.data
	if actual != expected {
		t.Errorf("Expected %v, actual %v", expected, actual)
	}
}

func TestDecodeSmallInt(t *testing.T) {
	assertDecoded(t, []byte{0x00}, uint(0))
	assertDecoded(t, []byte{0x01}, uint(1))
	assertDecoded(t, []byte{105}, uint(105))
	assertDecoded(t, []byte{0xff}, int(-1))
	assertDecoded(t, []byte{0x96}, int(-106))
}

func TestDecodeInt8(t *testing.T) {
	assertDecoded(t, []byte{0x6a, 0x6a}, uint(106))
	assertDecoded(t, []byte{0x6a, 0xff}, uint(0xff))
	assertDecoded(t, []byte{0x7a, 0xff}, int(-0xff))
}

func TestDecodeInt16(t *testing.T) {
	assertDecoded(t, []byte{0x6b, 0x00, 0x01}, uint(0x0100))
	assertDecoded(t, []byte{0x6b, 0xff, 0xff}, uint(0xffff))
	assertDecoded(t, []byte{0x7b, 0x00, 0x01}, int(-0x0100))
	assertDecoded(t, []byte{0x7b, 0xff, 0xff}, int(-0xffff))
}

func TestDecodeInt32(t *testing.T) {
	assertDecoded(t, []byte{0x6c, 0x00, 0x00, 0x00, 0x01}, uint(0x01000000))
	assertDecoded(t, []byte{0x6c, 0xff, 0xff, 0xff, 0xff}, uint(0xffffffff))
	assertDecoded(t, []byte{0x7c, 0x00, 0x00, 0x00, 0x01}, int(-0x01000000))
	assertDecoded(t, []byte{0x7c, 0xff, 0xff, 0xff, 0xff}, int64(-0xffffffff))
}

func TestDecodeInt64(t *testing.T) {
	assertDecoded(t, []byte{0x6d, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}, uint64(0x0100000000000000))
	assertDecoded(t, []byte{0x6d, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, uint64(0xffffffffffffffff))
	assertDecoded(t, []byte{0x7d, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}, int64(-0x0100000000000000))
	assertDecoded(t, []byte{0x7d, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80}, int64(-0x8000000000000000))
}
