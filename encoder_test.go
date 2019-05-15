package cbe

import (
	"bytes"
	"testing"
)

// TODO:
// - Array: Tried to start different data type before completion
// - Comment: invalid characters
// - String: invalid characters
// - Container: Unbalanced containers
// - Container: Unterminated container
// - Container: Max depth exceeded
// - Map: Bad key
// - Map: Missing value
// - Readme examples
// - Spec examples?

func testPanics(function func()) (didPanic bool) {
	defer func() {
		if r := recover(); r != nil {
			didPanic = true
		}
	}()
	didPanic = false
	function()
	return didPanic
}

func assertPanics(t *testing.T, function func()) {
	if !testPanics(function) {
		t.Errorf("Should have panicked but didn't")
	}
}

func assertSuccess(t *testing.T, err error) {
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func assertFailure(t *testing.T, err error) {
	if err == nil {
		t.Errorf("Unexpected success")
	}
}

func assertEncoded(t *testing.T, function func(*CbeEncoder), expected []byte) {
	encoder := NewCbeEncoder(100)
	function(encoder)
	actual := encoder.Encoded()
	if !bytes.Equal(actual, expected) {
		t.Errorf("Expected %v, actual %v", expected, actual)
	}
}

func TestPadding(t *testing.T) {
	assertEncoded(t, func(e *CbeEncoder) { e.Padding(1) }, []byte{0x7f})
	assertEncoded(t, func(e *CbeEncoder) { e.Padding(2) }, []byte{0x7f, 0x7f})
}

func TestNil(t *testing.T) {
	assertEncoded(t, func(e *CbeEncoder) { e.Nil() }, []byte{0x6f})
}

func TestIntSmall(t *testing.T) {
	assertEncoded(t, func(e *CbeEncoder) { e.Int(0) }, []byte{0})
	assertEncoded(t, func(e *CbeEncoder) { e.Int(1) }, []byte{1})
	assertEncoded(t, func(e *CbeEncoder) { e.Int(104) }, []byte{104})
	assertEncoded(t, func(e *CbeEncoder) { e.Int(105) }, []byte{105})
	assertEncoded(t, func(e *CbeEncoder) { e.Int(-1) }, []byte{255})
	assertEncoded(t, func(e *CbeEncoder) { e.Int(-105) }, []byte{151})
	assertEncoded(t, func(e *CbeEncoder) { e.Int(-106) }, []byte{150})
}

func TestInt8(t *testing.T) {
	assertEncoded(t, func(e *CbeEncoder) { e.Int(107) }, []byte{0x6a, 0x6b})
	assertEncoded(t, func(e *CbeEncoder) { e.Int(-107) }, []byte{0x7a, 0x6b})
	assertEncoded(t, func(e *CbeEncoder) { e.Int(255) }, []byte{0x6a, 0xff})
	assertEncoded(t, func(e *CbeEncoder) { e.Int(-255) }, []byte{0x7a, 0xff})

}

func TestInt16(t *testing.T) {
	assertEncoded(t, func(e *CbeEncoder) { e.Int(0x100) }, []byte{0x6b, 0x00, 0x01})
	assertEncoded(t, func(e *CbeEncoder) { e.Int(0x7fff) }, []byte{0x6b, 0xff, 0x7f})
	assertEncoded(t, func(e *CbeEncoder) { e.Int(-0x7fff) }, []byte{0x7b, 0xff, 0x7f})
	assertEncoded(t, func(e *CbeEncoder) { e.Int(0xffff) }, []byte{0x6b, 0xff, 0xff})
	assertEncoded(t, func(e *CbeEncoder) { e.Int(-0xffff) }, []byte{0x7b, 0xff, 0xff})
}

func TestInt32(t *testing.T) {
	assertEncoded(t, func(e *CbeEncoder) { e.Int(0x10000) }, []byte{0x6c, 0x00, 0x00, 0x01, 0x00})
	assertEncoded(t, func(e *CbeEncoder) { e.Int(0x7fffffff) }, []byte{0x6c, 0xff, 0xff, 0xff, 0x7f})
	assertEncoded(t, func(e *CbeEncoder) { e.Int(-0x7fffffff) }, []byte{0x7c, 0xff, 0xff, 0xff, 0x7f})
	assertEncoded(t, func(e *CbeEncoder) { e.Int(0xffffffff) }, []byte{0x6c, 0xff, 0xff, 0xff, 0xff})
	assertEncoded(t, func(e *CbeEncoder) { e.Int(-0xffffffff) }, []byte{0x7c, 0xff, 0xff, 0xff, 0xff})
}

func TestInt64(t *testing.T) {
	assertEncoded(t, func(e *CbeEncoder) { e.Int(0x100000000) }, []byte{0x6d, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00})
	assertEncoded(t, func(e *CbeEncoder) { e.Int(0x7fffffffffffffff) }, []byte{0x6d, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x7f})
	assertEncoded(t, func(e *CbeEncoder) { e.Int(-0x7fffffffffffffff) }, []byte{0x7d, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x7f})
}

func TestFloat64(t *testing.T) {
	assertEncoded(t, func(e *CbeEncoder) { e.Float(1.0123) }, []byte{0x73, 0x51, 0xda, 0x1b, 0x7c, 0x61, 0x32, 0xf0, 0x3f})
}

func TestList(t *testing.T) {
	assertEncoded(t, func(e *CbeEncoder) { e.ListBegin() }, []byte{0x93})
	assertEncoded(t, func(e *CbeEncoder) { e.ListBegin(); e.ListEnd() }, []byte{0x93, 0x95})
	assertEncoded(t, func(e *CbeEncoder) { e.ListBegin(); e.Int(1); e.String("a"); e.ListEnd() }, []byte{0x93, 0x01, 0x81, 0x61, 0x95})
}

func TestMap(t *testing.T) {
	assertEncoded(t, func(e *CbeEncoder) { e.MapBegin() }, []byte{0x94})
	assertEncoded(t, func(e *CbeEncoder) { e.MapBegin(); e.ListEnd() }, []byte{0x94, 0x95})
	assertEncoded(t, func(e *CbeEncoder) {
		e.MapBegin()
		e.String("1")
		e.Uint(1)
		e.String("2")
		e.Uint(2)
		e.ListEnd()
	},
		[]byte{0x94, 0x81, 0x31, 0x01, 0x81, 0x32, 0x02, 0x95})
}

func TestBinary(t *testing.T) {
	assertEncoded(t, func(e *CbeEncoder) { e.Bytes([]byte{}) }, []byte{0x91, 0x00})
	assertEncoded(t, func(e *CbeEncoder) { e.Bytes([]byte{1}) }, []byte{0x91, 0x04, 0x01})
	assertEncoded(t, func(e *CbeEncoder) { e.Bytes([]byte{1, 2}) }, []byte{0x91, 0x08, 0x01, 0x02})
	assertEncoded(t, func(e *CbeEncoder) { e.Bytes([]byte{1, 2, 3}) }, []byte{0x91, 0x0c, 0x01, 0x02, 0x03})
}

func TestString(t *testing.T) {
	assertEncoded(t, func(e *CbeEncoder) { e.String("") }, []byte{0x80})
	assertEncoded(t, func(e *CbeEncoder) { e.String("0") }, []byte{0x81, 0x30})
	assertEncoded(t, func(e *CbeEncoder) { e.String("01") }, []byte{0x82, 0x30, 0x31})
	assertEncoded(t, func(e *CbeEncoder) { e.String("0123456789012345") }, []byte{
		0x90, 0x40, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37,
		0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35,
	})
}

func TestBinaryTooLong(t *testing.T) {
	encoder := NewCbeEncoder(100)
	assertSuccess(t, encoder.BinaryBegin(10))
	assertFailure(t, encoder.BinaryData(make([]byte, 11)))
}

func TestBinaryTooShort(t *testing.T) {
	encoder := NewCbeEncoder(100)
	assertSuccess(t, encoder.BinaryBegin(10))
	assertSuccess(t, encoder.BinaryData(make([]byte, 9)))
	assertPanics(t, func() {
		encoder.Encoded()
	})
}

func TestStringTooLong(t *testing.T) {
	encoder := NewCbeEncoder(100)
	assertSuccess(t, encoder.StringBegin(6))
	assertFailure(t, encoder.StringData([]byte("abcdefg")))
}

func TestStringTooShort(t *testing.T) {
	encoder := NewCbeEncoder(100)
	assertSuccess(t, encoder.StringBegin(8))
	assertSuccess(t, encoder.StringData([]byte("abcdefg")))
	assertPanics(t, func() {
		encoder.Encoded()
	})
}

func TestCommentTooLong(t *testing.T) {
	encoder := NewCbeEncoder(100)
	assertSuccess(t, encoder.CommentBegin(6))
	assertFailure(t, encoder.CommentData([]byte("abcdefg")))
}

func TestCommentTooShort(t *testing.T) {
	encoder := NewCbeEncoder(100)
	assertSuccess(t, encoder.CommentBegin(8))
	assertSuccess(t, encoder.CommentData([]byte("abcdefg")))
	assertPanics(t, func() {
		encoder.Encoded()
	})
}
