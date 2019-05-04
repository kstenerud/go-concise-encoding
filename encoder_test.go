package cbe

import (
	"bytes"
	"testing"
)

// TODO:
// - Array: too many or too few bytes
// - Array: Tried to start different data type before completion
// - Array: Concenience functions
// - Comment: invalid characters
// - String: 0-15, more bytes
// - String: invalid characters
// - Container: Unbalanced containers
// - Container: Unterminated container
// - Container: Max depth exceeded
// - Map: Bad key
// - Map: Missing value
// - Readme examples
// - Spec examples?

func assertEncoded(t *testing.T, encoder *Encoder, expected []byte) {
	actual := encoder.Encoded()
	if !bytes.Equal(actual, expected) {
		t.Errorf("Expected %v, actual %v", expected, actual)
	}
}

func TestPadding(t *testing.T) {
	assertEncoded(t, NewEncoder(9).Padding(1), []byte{0x7f})
	assertEncoded(t, NewEncoder(9).Padding(2), []byte{0x7f, 0x7f})
}

func TestNil(t *testing.T) {
	assertEncoded(t, NewEncoder(9).Nil(), []byte{0x6f})
}

func TestIntSmall(t *testing.T) {
	assertEncoded(t, NewEncoder(9).Int(0), []byte{0})
	assertEncoded(t, NewEncoder(9).Int(1), []byte{1})
	assertEncoded(t, NewEncoder(9).Int(104), []byte{104})
	assertEncoded(t, NewEncoder(9).Int(105), []byte{105})
	assertEncoded(t, NewEncoder(9).Int(-1), []byte{255})
	assertEncoded(t, NewEncoder(9).Int(-105), []byte{151})
	assertEncoded(t, NewEncoder(9).Int(-106), []byte{150})
}

func TestInt8(t *testing.T) {
	assertEncoded(t, NewEncoder(9).Int(107), []byte{0x6a, 0x6b})
	assertEncoded(t, NewEncoder(9).Int(-107), []byte{0x7a, 0x6b})
	assertEncoded(t, NewEncoder(9).Int(255), []byte{0x6a, 0xff})
	assertEncoded(t, NewEncoder(9).Int(-255), []byte{0x7a, 0xff})

}

func TestInt16(t *testing.T) {
	assertEncoded(t, NewEncoder(9).Int(0x100), []byte{0x6b, 0x00, 0x01})
	assertEncoded(t, NewEncoder(9).Int(0x7fff), []byte{0x6b, 0xff, 0x7f})
	assertEncoded(t, NewEncoder(9).Int(-0x7fff), []byte{0x7b, 0xff, 0x7f})
	assertEncoded(t, NewEncoder(9).Int(0xffff), []byte{0x6b, 0xff, 0xff})
	assertEncoded(t, NewEncoder(9).Int(-0xffff), []byte{0x7b, 0xff, 0xff})
}

func TestInt32(t *testing.T) {
	assertEncoded(t, NewEncoder(9).Int(0x10000), []byte{0x6c, 0x00, 0x00, 0x01, 0x00})
	assertEncoded(t, NewEncoder(9).Int(0x7fffffff), []byte{0x6c, 0xff, 0xff, 0xff, 0x7f})
	assertEncoded(t, NewEncoder(9).Int(-0x7fffffff), []byte{0x7c, 0xff, 0xff, 0xff, 0x7f})
	assertEncoded(t, NewEncoder(9).Int(0xffffffff), []byte{0x6c, 0xff, 0xff, 0xff, 0xff})
	assertEncoded(t, NewEncoder(9).Int(-0xffffffff), []byte{0x7c, 0xff, 0xff, 0xff, 0xff})
}

func TestInt64(t *testing.T) {
	assertEncoded(t, NewEncoder(9).Int(0x100000000), []byte{0x6d, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00})
	assertEncoded(t, NewEncoder(9).Int(0x7fffffffffffffff), []byte{0x6d, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x7f})
	assertEncoded(t, NewEncoder(9).Int(-0x7fffffffffffffff), []byte{0x7d, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x7f})
}

func TestFloat64(t *testing.T) {
	assertEncoded(t, NewEncoder(9).Float(1.0123), []byte{0x73, 0x51, 0xda, 0x1b, 0x7c, 0x61, 0x32, 0xf0, 0x3f})
}

func TestList(t *testing.T) {
	assertEncoded(t, NewEncoder(9).ListBegin(), []byte{0x93})
	assertEncoded(t, NewEncoder(9).ListBegin().ListEnd(), []byte{0x93, 0x95})
	assertEncoded(t, NewEncoder(9).ListBegin().Int(1).String("a").ListEnd(), []byte{0x93, 0x01, 0x81, 0x61, 0x95})
}

func TestMap(t *testing.T) {
	assertEncoded(t, NewEncoder(9).MapBegin(), []byte{0x94})
	assertEncoded(t, NewEncoder(9).MapBegin().ListEnd(), []byte{0x94, 0x95})
	assertEncoded(t, NewEncoder(9).MapBegin().String("1").Uint(1).String("2").Uint(2).ListEnd(),
		[]byte{0x94, 0x81, 0x31, 0x01, 0x81, 0x32, 0x02, 0x95})
}

func TestBinary(t *testing.T) {
	assertEncoded(t, NewEncoder(9).Binary([]byte{}), []byte{0x91, 0x00})
	assertEncoded(t, NewEncoder(9).Binary([]byte{1}), []byte{0x91, 0x04, 0x01})
	assertEncoded(t, NewEncoder(9).Binary([]byte{1, 2}), []byte{0x91, 0x08, 0x01, 0x02})
	assertEncoded(t, NewEncoder(9).Binary([]byte{1, 2, 3}), []byte{0x91, 0x0c, 0x01, 0x02, 0x03})
	// TODO: Longer than 64 bytes
}

func TestString(t *testing.T) {
	assertEncoded(t, NewEncoder(9).String(""), []byte{0x80})
	assertEncoded(t, NewEncoder(9).String("0"), []byte{0x81, 0x30})
	assertEncoded(t, NewEncoder(9).String("01"), []byte{0x82, 0x30, 0x31})
	assertEncoded(t, NewEncoder(9).String("0123456789012345"), []byte{
		0x90, 0x40, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37,
		0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35,
	})
	// TODO: Longer than 64 bytes
}
