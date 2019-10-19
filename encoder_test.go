package cbe

import (
	"testing"
	"time"
)

func TestEncodePadding(t *testing.T) {
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) error { return e.Padding(1) }, []byte{0x7f})
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) error { return e.Padding(2) }, []byte{0x7f, 0x7f})
}

func TestEncodeNil(t *testing.T) {
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) error { return e.Nil() }, []byte{0x7e})
}

func TestEncodeIntSmall(t *testing.T) {
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) error { return e.Int(0) }, []byte{0})
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) error { return e.Int(1) }, []byte{1})
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) error { return e.Int(99) }, []byte{99})
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) error { return e.Int(100) }, []byte{100})
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) error { return e.Int(-1) }, []byte{255})
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) error { return e.Int(-99) }, []byte{157})
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) error { return e.Int(-100) }, []byte{156})
}

func TestEncodeInt8(t *testing.T) {
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) error { return e.Int(101) }, []byte{0x68, 0x65})
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) error { return e.Int(-101) }, []byte{0x69, 0x65})
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) error { return e.Int(255) }, []byte{0x68, 0xff})
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) error { return e.Int(-255) }, []byte{0x69, 0xff})
}

func TestEncodeInt16(t *testing.T) {
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) error { return e.Int(0x100) }, []byte{0x6a, 0x00, 0x01})
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) error { return e.Int(0x7fff) }, []byte{0x6a, 0xff, 0x7f})
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) error { return e.Int(-0x7fff) }, []byte{0x6b, 0xff, 0x7f})
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) error { return e.Int(0xffff) }, []byte{0x6a, 0xff, 0xff})
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) error { return e.Int(-0xffff) }, []byte{0x6b, 0xff, 0xff})
}

func TestEncodeInt32(t *testing.T) {
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) error { return e.Int(0x10000) }, []byte{0x6c, 0x00, 0x00, 0x01, 0x00})
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) error { return e.Int(0x7fffffff) }, []byte{0x6c, 0xff, 0xff, 0xff, 0x7f})
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) error { return e.Int(-0x7fffffff) }, []byte{0x6d, 0xff, 0xff, 0xff, 0x7f})
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) error { return e.Int(0xffffffff) }, []byte{0x6c, 0xff, 0xff, 0xff, 0xff})
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) error { return e.Int(-0xffffffff) }, []byte{0x6d, 0xff, 0xff, 0xff, 0xff})
}

func TestEncodeInt64(t *testing.T) {
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) error { return e.Int(0x100000000) }, []byte{0x6e, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00})
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) error { return e.Int(0x7fffffffffffffff) }, []byte{0x6e, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x7f})
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) error { return e.Int(-0x7fffffffffffffff) }, []byte{0x6f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x7f})
}

func TestEncodeFloat64(t *testing.T) {
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) error { return e.Float(1.0123) }, []byte{0x71, 0x51, 0xda, 0x1b, 0x7c, 0x61, 0x32, 0xf0, 0x3f})
}

func TestEncodeList(t *testing.T) {
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) error { return e.ListBegin() }, []byte{0x77})
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) error { return ShortCircuit(e.ListBegin(), e.ContainerEnd()) }, []byte{0x77, 0x7b})
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) error {
		return ShortCircuit(e.ListBegin(), e.Int(1), e.String("a"), e.ContainerEnd())
	}, []byte{0x77, 0x01, 0x81, 0x61, 0x7b})
}

func TestEncodeMap(t *testing.T) {
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) error { return e.MapBegin() }, []byte{0x78})
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) error { return ShortCircuit(e.MapBegin(), e.ContainerEnd()) }, []byte{0x78, 0x7b})
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) error {
		return ShortCircuit(
			e.MapBegin(),
			e.String("1"),
			e.Uint(1),
			e.String("2"),
			e.Uint(2),
			e.ContainerEnd())
	},
		[]byte{0x78, 0x81, 0x31, 0x01, 0x81, 0x32, 0x02, 0x7b})
}

func TestEncodeInlineList(t *testing.T) {
	assertEncoded(t, ContainerTypeList, func(e *CbeEncoder) error {
		return ShortCircuit(
			e.Nil(),
			e.Int(5),
			e.String(""))
	}, []byte{0x7e, 0x05, 0x80})
}

func TestEncodeInlineMap(t *testing.T) {
	assertEncoded(t, ContainerTypeUnorderedMap, func(e *CbeEncoder) error {
		return ShortCircuit(
			e.Int(1),
			e.String(""))
	}, []byte{0x01, 0x80})
}

func TestEncodeBytes(t *testing.T) {
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) error { return e.Bytes([]byte{}) }, []byte{0x91, 0x00})
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) error { return e.Bytes([]byte{1}) }, []byte{0x91, 0x01, 0x01})
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) error { return e.Bytes([]byte{1, 2}) }, []byte{0x91, 0x02, 0x01, 0x02})
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) error { return e.Bytes([]byte{1, 2, 3}) }, []byte{0x91, 0x03, 0x01, 0x02, 0x03})
}

func TestEncodeBytesLong(t *testing.T) {
	bytesLength := 500
	bytes := generateBytes(bytesLength)
	encoded := append([]byte{0x91, 0x83, 0x74}, bytes...)

	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) error { return e.Bytes(bytes) }, encoded)
}

func TestEncodeString(t *testing.T) {
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) error { return e.String("") }, []byte{0x80})
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) error { return e.String("0") }, []byte{0x81, 0x30})
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) error { return e.String("01") }, []byte{0x82, 0x30, 0x31})
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) error { return e.String("0123456789012345") }, []byte{
		0x90, 0x10, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37,
		0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35,
	})
}

func TestEncodeStringInvalid(t *testing.T) {
	encoder := NewCbeEncoder(ContainerTypeNone, nil, 100)
	assertFailure(t, encoder.String(string([]byte{0x40, 0x81, 0x42, 0x43, 0x44})))
}

func TestEncodeCommentInvalid(t *testing.T) {
	encoder := NewCbeEncoder(ContainerTypeNone, nil, 100)
	assertFailure(t, encoder.Comment(string([]byte{0x40, 0x81, 0x42, 0x43, 0x44})))
	assertFailure(t, encoder.Comment("A comment\nwith a newline"))
}

func TestEncodeBytesTooLong(t *testing.T) {
	encoder := NewCbeEncoder(ContainerTypeNone, nil, 100)
	assertSuccess(t, encoder.BytesBegin(10))
	_, err := encoder.ArrayData(make([]byte, 11))
	assertFailure(t, err)
}

func TestEncodeBytesTooShort(t *testing.T) {
	encoder := NewCbeEncoder(ContainerTypeNone, nil, 100)
	assertSuccess(t, encoder.BytesBegin(10))
	_, err := encoder.ArrayData(make([]byte, 9))
	assertSuccess(t, err)
	assertFailure(t, encoder.End())
}

func TestEncodeStringTooLong(t *testing.T) {
	encoder := NewCbeEncoder(ContainerTypeNone, nil, 100)
	assertSuccess(t, encoder.StringBegin(6))
	_, err := encoder.ArrayData([]byte("abcdefg"))
	assertFailure(t, err)
}

func TestEncodeStringTooShort(t *testing.T) {
	encoder := NewCbeEncoder(ContainerTypeNone, nil, 100)
	assertSuccess(t, encoder.StringBegin(8))
	_, err := encoder.ArrayData([]byte("abcdefg"))
	assertSuccess(t, err)
	assertFailure(t, encoder.End())
}

func TestEncodeCommentTooLong(t *testing.T) {
	encoder := NewCbeEncoder(ContainerTypeNone, nil, 100)
	assertSuccess(t, encoder.CommentBegin(6))
	_, err := encoder.ArrayData([]byte("abcdefg"))
	assertFailure(t, err)
}

func TestEncodeCommentTooShort(t *testing.T) {
	encoder := NewCbeEncoder(ContainerTypeNone, nil, 100)
	assertSuccess(t, encoder.CommentBegin(8))
	_, err := encoder.ArrayData([]byte("abcdefg"))
	assertSuccess(t, err)
	assertFailure(t, encoder.End())
}

func TestEncodeChangeArrayType(t *testing.T) {
	encoder := NewCbeEncoder(ContainerTypeNone, nil, 100)
	assertSuccess(t, encoder.BytesBegin(10))
	_, err := encoder.ArrayData(make([]byte, 5))
	assertSuccess(t, err)
	assertFailure(t, encoder.StringBegin(10))
}

func TestEncodeUnbalancedContainers(t *testing.T) {
	encoder := NewCbeEncoder(ContainerTypeNone, nil, 100)
	assertSuccess(t, encoder.ListBegin())
	assertFailure(t, encoder.End())
}

func TestEncodeCloseListTooManyTimes(t *testing.T) {
	encoder := NewCbeEncoder(ContainerTypeNone, nil, 100)
	assertSuccess(t, encoder.ListBegin())
	assertSuccess(t, encoder.ContainerEnd())
	assertFailure(t, encoder.ContainerEnd())
}

func TestEncodeCloseMapTooManyTimes(t *testing.T) {
	encoder := NewCbeEncoder(ContainerTypeNone, nil, 100)
	assertSuccess(t, encoder.MapBegin())
	assertSuccess(t, encoder.ContainerEnd())
	assertFailure(t, encoder.ContainerEnd())
}

func TestEncodeCloseNoContainer(t *testing.T) {
	encoder := NewCbeEncoder(ContainerTypeNone, nil, 100)
	assertFailure(t, encoder.ContainerEnd())
}

func TestEncodeMapMissingValue(t *testing.T) {
	encoder := NewCbeEncoder(ContainerTypeNone, nil, 100)
	assertSuccess(t, encoder.MapBegin())
	assertSuccess(t, encoder.Int(1))
	assertFailure(t, encoder.ContainerEnd())
}

func TestEncodeMapNilKey(t *testing.T) {
	encoder := NewCbeEncoder(ContainerTypeNone, nil, 100)
	assertSuccess(t, encoder.MapBegin())
	assertFailure(t, encoder.Nil())
}

func TestEncodeMapListKey(t *testing.T) {
	encoder := NewCbeEncoder(ContainerTypeNone, nil, 100)
	assertSuccess(t, encoder.MapBegin())
	assertFailure(t, encoder.ListBegin())
}

func TestEncodeMapMapKey(t *testing.T) {
	encoder := NewCbeEncoder(ContainerTypeNone, nil, 100)
	assertSuccess(t, encoder.MapBegin())
	assertFailure(t, encoder.MapBegin())
}

func TestEncodeMapBytesKey(t *testing.T) {
	encoder := NewCbeEncoder(ContainerTypeNone, nil, 100)
	assertSuccess(t, encoder.MapBegin())
	assertSuccess(t, encoder.Bytes([]byte{1, 2, 3}))
	assertSuccess(t, encoder.Bytes([]byte{4, 5, 6}))
	assertSuccess(t, encoder.ContainerEnd())
}

func TestEncodeMapWithComments(t *testing.T) {
	encoder := NewCbeEncoder(ContainerTypeNone, nil, 100)
	assertSuccess(t, encoder.Comment("A comment before the map"))
	assertSuccess(t, encoder.MapBegin())
	assertSuccess(t, encoder.Comment("A comment"))
	assertSuccess(t, encoder.String("a key"))
	assertSuccess(t, encoder.Comment("Another comment"))
	assertSuccess(t, encoder.Bool(true))
	assertSuccess(t, encoder.Comment("Yet another comment"))
	assertSuccess(t, encoder.ContainerEnd())
	assertSuccess(t, encoder.Comment("A comment after the map"))
}

func TestEncodeListWithComments(t *testing.T) {
	encoder := NewCbeEncoder(ContainerTypeNone, nil, 100)
	assertSuccess(t, encoder.Comment("A comment before the list"))
	assertSuccess(t, encoder.ListBegin())
	assertSuccess(t, encoder.Comment("A comment"))
	assertSuccess(t, encoder.String("a string"))
	assertSuccess(t, encoder.Comment("Another comment"))
	assertSuccess(t, encoder.Bool(true))
	assertSuccess(t, encoder.Comment("Yet another comment"))
	assertSuccess(t, encoder.ContainerEnd())
	assertSuccess(t, encoder.Comment("A comment after the list"))
}

func TestEncodeContainerLimitExceeded(t *testing.T) {
	encoder := NewCbeEncoder(ContainerTypeNone, nil, 4)
	assertSuccess(t, encoder.ListBegin())
	assertSuccess(t, encoder.ListBegin())
	assertSuccess(t, encoder.MapBegin())
	assertSuccess(t, encoder.Bool(true))
	assertSuccess(t, encoder.ListBegin())
	assertFailure(t, encoder.MapBegin())
}

func TestEncodeToExternalBuffer(t *testing.T) {
	assertEncodesToExternalBuffer(t, ContainerTypeNone, nil, 1)
	assertEncodesToExternalBuffer(t, ContainerTypeNone, true, 1)
	assertEncodesToExternalBuffer(t, ContainerTypeNone, false, 1)
	assertEncodesToExternalBuffer(t, ContainerTypeNone, 0, 1)
	assertEncodesToExternalBuffer(t, ContainerTypeNone, -150, 2)
	assertEncodesToExternalBuffer(t, ContainerTypeNone, 1.1, 9)
	assertEncodesToExternalBuffer(t, ContainerTypeNone, 1.5, 5)
	assertEncodesToExternalBuffer(t, ContainerTypeNone, "test", 5)
	assertEncodesToExternalBuffer(t, ContainerTypeNone, "a longer string to test with", 30)
	assertEncodesToExternalBuffer(t, ContainerTypeNone, []byte{0x01}, 3)
	assertEncodesToExternalBuffer(t, ContainerTypeNone, time.Now(), 12)
	assertEncodesToExternalBuffer(t, ContainerTypeNone, map[interface{}]interface{}{"test": 1}, 8)
	assertEncodesToExternalBuffer(t, ContainerTypeNone, []interface{}{"test", 1}, 8)
}
