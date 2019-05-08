package cbe

import (
	"bytes"
	"testing"
)

func assertMarshaled(t *testing.T, value interface{}, expected []byte) {
	encoder := NewEncoder(100)
	marshaler := NewMarshaler(encoder)
	marshaler.Marshal(value)
	actual := encoder.Encoded()
	if !bytes.Equal(actual, expected) {
		t.Errorf("Expected %v, actual %v", expected, actual)
	}
}

func TestNil(t *testing.T) {
	assertMarshaled(t, nil, []byte{0x6f})
}

func Test0(t *testing.T) {
	assertMarshaled(t, 0, []byte{0x00})
}

func TestN1(t *testing.T) {
	assertMarshaled(t, -1, []byte{0xff})
}

func Test200(t *testing.T) {
	assertMarshaled(t, 200, []byte{0x6a, 200})
}

func Test2000(t *testing.T) {
	assertMarshaled(t, 2000, []byte{0x6b, 0xd0, 0x07})
}

func Test1_5(t *testing.T) {
	assertMarshaled(t, 1.5, []byte{0x72, 0x00, 0x00, 0xc0, 0x3f})
}

func Test1_0123(t *testing.T) {
	assertMarshaled(t, 1.0123, []byte{0x73, 0x51, 0xda, 0x1b, 0x7c, 0x61, 0x32, 0xf0, 0x3f})
}
