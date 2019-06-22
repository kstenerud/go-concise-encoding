package cbe

import (
	"testing"
	"time"
)

func TestEncodePadding(t *testing.T) {
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) { e.Padding(1) }, []byte{0x7f})
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) { e.Padding(2) }, []byte{0x7f, 0x7f})
}

func TestEncodeNil(t *testing.T) {
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) { e.Nil() }, []byte{0x6f})
}

func TestEncodeIntSmall(t *testing.T) {
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) { e.Int(0) }, []byte{0})
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) { e.Int(1) }, []byte{1})
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) { e.Int(104) }, []byte{104})
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) { e.Int(105) }, []byte{105})
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) { e.Int(-1) }, []byte{255})
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) { e.Int(-105) }, []byte{151})
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) { e.Int(-106) }, []byte{150})
}

func TestEncodeInt8(t *testing.T) {
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) { e.Int(107) }, []byte{0x6a, 0x6b})
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) { e.Int(-107) }, []byte{0x7a, 0x6b})
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) { e.Int(255) }, []byte{0x6a, 0xff})
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) { e.Int(-255) }, []byte{0x7a, 0xff})

}

func TestEncodeInt16(t *testing.T) {
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) { e.Int(0x100) }, []byte{0x6b, 0x00, 0x01})
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) { e.Int(0x7fff) }, []byte{0x6b, 0xff, 0x7f})
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) { e.Int(-0x7fff) }, []byte{0x7b, 0xff, 0x7f})
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) { e.Int(0xffff) }, []byte{0x6b, 0xff, 0xff})
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) { e.Int(-0xffff) }, []byte{0x7b, 0xff, 0xff})
}

func TestEncodeInt32(t *testing.T) {
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) { e.Int(0x10000) }, []byte{0x6c, 0x00, 0x00, 0x01, 0x00})
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) { e.Int(0x7fffffff) }, []byte{0x6c, 0xff, 0xff, 0xff, 0x7f})
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) { e.Int(-0x7fffffff) }, []byte{0x7c, 0xff, 0xff, 0xff, 0x7f})
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) { e.Int(0xffffffff) }, []byte{0x6c, 0xff, 0xff, 0xff, 0xff})
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) { e.Int(-0xffffffff) }, []byte{0x7c, 0xff, 0xff, 0xff, 0xff})
}

func TestEncodeInt64(t *testing.T) {
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) { e.Int(0x100000000) }, []byte{0x6d, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00})
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) { e.Int(0x7fffffffffffffff) }, []byte{0x6d, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x7f})
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) { e.Int(-0x7fffffffffffffff) }, []byte{0x7d, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x7f})
}

func TestEncodeFloat64(t *testing.T) {
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) { e.Float(1.0123) }, []byte{0x73, 0x51, 0xda, 0x1b, 0x7c, 0x61, 0x32, 0xf0, 0x3f})
}

func TestEncodeList(t *testing.T) {
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) { e.ListBegin() }, []byte{0x93})
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) { e.ListBegin(); e.ListEnd() }, []byte{0x93, 0x95})
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) { e.ListBegin(); e.Int(1); e.String("a"); e.ListEnd() }, []byte{0x93, 0x01, 0x81, 0x61, 0x95})
}

func TestEncodeMap(t *testing.T) {
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) { e.MapBegin() }, []byte{0x94})
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) { e.MapBegin(); e.ListEnd() }, []byte{0x94, 0x95})
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) {
		e.MapBegin()
		e.String("1")
		e.Uint(1)
		e.String("2")
		e.Uint(2)
		e.ListEnd()
	},
		[]byte{0x94, 0x81, 0x31, 0x01, 0x81, 0x32, 0x02, 0x95})
}

func TestEncodeInlineList(t *testing.T) {
	assertEncoded(t, ContainerTypeList, func(e *CbeEncoder) {
		e.Nil()
		e.Int(5)
		e.String("")
	}, []byte{0x6f, 0x05, 0x80})
}

func TestEncodeInlineMap(t *testing.T) {
	assertEncoded(t, ContainerTypeMap, func(e *CbeEncoder) {
		e.Int(1)
		e.String("")
	}, []byte{0x01, 0x80})
}

func TestEncodeBytes(t *testing.T) {
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) { e.Bytes([]byte{}) }, []byte{0x91, 0x00})
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) { e.Bytes([]byte{1}) }, []byte{0x91, 0x01, 0x01})
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) { e.Bytes([]byte{1, 2}) }, []byte{0x91, 0x02, 0x01, 0x02})
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) { e.Bytes([]byte{1, 2, 3}) }, []byte{0x91, 0x03, 0x01, 0x02, 0x03})
}

func TestEncodeBytesLong(t *testing.T) {
	bytesLength := 500
	bytes := generateBytes(bytesLength)
	encoded := append([]byte{0x91, 0x83, 0x74}, bytes...)

	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) { e.Bytes(bytes) }, encoded)
}

func TestEncodeString(t *testing.T) {
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) { e.String("") }, []byte{0x80})
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) { e.String("0") }, []byte{0x81, 0x30})
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) { e.String("01") }, []byte{0x82, 0x30, 0x31})
	assertEncoded(t, ContainerTypeNone, func(e *CbeEncoder) { e.String("0123456789012345") }, []byte{
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
	assertFailure(t, encoder.BytesData(make([]byte, 11)))
}

func TestEncodeBytesTooShort(t *testing.T) {
	encoder := NewCbeEncoder(ContainerTypeNone, nil, 100)
	assertSuccess(t, encoder.BytesBegin(10))
	assertSuccess(t, encoder.BytesData(make([]byte, 9)))
	assertFailure(t, encoder.End())
}

func TestEncodeStringTooLong(t *testing.T) {
	encoder := NewCbeEncoder(ContainerTypeNone, nil, 100)
	assertSuccess(t, encoder.StringBegin(6))
	assertFailure(t, encoder.StringData([]byte("abcdefg")))
}

func TestEncodeStringTooShort(t *testing.T) {
	encoder := NewCbeEncoder(ContainerTypeNone, nil, 100)
	assertSuccess(t, encoder.StringBegin(8))
	assertSuccess(t, encoder.StringData([]byte("abcdefg")))
	assertFailure(t, encoder.End())
}

func TestEncodeCommentTooLong(t *testing.T) {
	encoder := NewCbeEncoder(ContainerTypeNone, nil, 100)
	assertSuccess(t, encoder.CommentBegin(6))
	assertFailure(t, encoder.CommentData([]byte("abcdefg")))
}

func TestEncodeCommentTooShort(t *testing.T) {
	encoder := NewCbeEncoder(ContainerTypeNone, nil, 100)
	assertSuccess(t, encoder.CommentBegin(8))
	assertSuccess(t, encoder.CommentData([]byte("abcdefg")))
	assertFailure(t, encoder.End())
}

func TestEncodeChangeArrayType(t *testing.T) {
	encoder := NewCbeEncoder(ContainerTypeNone, nil, 100)
	assertSuccess(t, encoder.BytesBegin(10))
	assertSuccess(t, encoder.BytesData(make([]byte, 5)))
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
	assertSuccess(t, encoder.ListEnd())
	assertFailure(t, encoder.ListEnd())
}

func TestEncodeCloseMapTooManyTimes(t *testing.T) {
	encoder := NewCbeEncoder(ContainerTypeNone, nil, 100)
	assertSuccess(t, encoder.MapBegin())
	assertSuccess(t, encoder.MapEnd())
	assertFailure(t, encoder.MapEnd())
}

func TestEncodeCloseNoContainer(t *testing.T) {
	encoder := NewCbeEncoder(ContainerTypeNone, nil, 100)
	assertFailure(t, encoder.ListEnd())
}

func TestEncodeMapMissingValue(t *testing.T) {
	encoder := NewCbeEncoder(ContainerTypeNone, nil, 100)
	assertSuccess(t, encoder.MapBegin())
	assertSuccess(t, encoder.Int(1))
	assertFailure(t, encoder.MapEnd())
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
	assertSuccess(t, encoder.MapEnd())
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
	assertSuccess(t, encoder.MapEnd())
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
	assertSuccess(t, encoder.ListEnd())
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
	assertEncodesToExternalBuffer(t, ContainerTypeNone, time.Now(), 9)
	assertEncodesToExternalBuffer(t, ContainerTypeNone, map[interface{}]interface{}{"test": 1}, 8)
	assertEncodesToExternalBuffer(t, ContainerTypeNone, []interface{}{"test", 1}, 8)
}
