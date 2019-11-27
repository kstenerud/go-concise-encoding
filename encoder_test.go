package cbe

import (
	"net/url"
	"testing"
	"time"
)

func TestEncodePadding(t *testing.T) {
	assertEncoded(t, InlineContainerTypeNone, func(e *Encoder) error { return e.Padding(1) }, []byte{0x7f})
	assertEncoded(t, InlineContainerTypeNone, func(e *Encoder) error { return e.Padding(2) }, []byte{0x7f, 0x7f})
}

func TestEncodeNil(t *testing.T) {
	assertEncoded(t, InlineContainerTypeNone, func(e *Encoder) error { return e.Nil() }, []byte{0x7e})
}

func TestEncodeIntSmall(t *testing.T) {
	assertEncoded(t, InlineContainerTypeNone, func(e *Encoder) error { return e.Int(0) }, []byte{0})
	assertEncoded(t, InlineContainerTypeNone, func(e *Encoder) error { return e.Int(1) }, []byte{1})
	assertEncoded(t, InlineContainerTypeNone, func(e *Encoder) error { return e.Int(99) }, []byte{99})
	assertEncoded(t, InlineContainerTypeNone, func(e *Encoder) error { return e.Int(100) }, []byte{100})
	assertEncoded(t, InlineContainerTypeNone, func(e *Encoder) error { return e.Int(-1) }, []byte{255})
	assertEncoded(t, InlineContainerTypeNone, func(e *Encoder) error { return e.Int(-99) }, []byte{157})
	assertEncoded(t, InlineContainerTypeNone, func(e *Encoder) error { return e.Int(-100) }, []byte{156})
}

func TestEncodeInt8(t *testing.T) {
	assertEncoded(t, InlineContainerTypeNone, func(e *Encoder) error { return e.Int(101) }, []byte{0x68, 0x65})
	assertEncoded(t, InlineContainerTypeNone, func(e *Encoder) error { return e.Int(-101) }, []byte{0x69, 0x65})
	assertEncoded(t, InlineContainerTypeNone, func(e *Encoder) error { return e.Int(255) }, []byte{0x68, 0xff})
	assertEncoded(t, InlineContainerTypeNone, func(e *Encoder) error { return e.Int(-255) }, []byte{0x69, 0xff})
}

func TestEncodeInt16(t *testing.T) {
	assertEncoded(t, InlineContainerTypeNone, func(e *Encoder) error { return e.Int(0x100) }, []byte{0x6a, 0x00, 0x01})
	assertEncoded(t, InlineContainerTypeNone, func(e *Encoder) error { return e.Int(0x7fff) }, []byte{0x6a, 0xff, 0x7f})
	assertEncoded(t, InlineContainerTypeNone, func(e *Encoder) error { return e.Int(-0x7fff) }, []byte{0x6b, 0xff, 0x7f})
	assertEncoded(t, InlineContainerTypeNone, func(e *Encoder) error { return e.Int(0xffff) }, []byte{0x6a, 0xff, 0xff})
	assertEncoded(t, InlineContainerTypeNone, func(e *Encoder) error { return e.Int(-0xffff) }, []byte{0x6b, 0xff, 0xff})
}

func TestEncodeInt21(t *testing.T) {
	assertEncoded(t, InlineContainerTypeNone, func(e *Encoder) error { return e.Int(0x10000) }, []byte{0x66, 0x84, 0x80, 0x00})
}

func TestEncodeInt32(t *testing.T) {
	assertEncoded(t, InlineContainerTypeNone, func(e *Encoder) error { return e.Int(0x7fffffff) }, []byte{0x6c, 0xff, 0xff, 0xff, 0x7f})
	assertEncoded(t, InlineContainerTypeNone, func(e *Encoder) error { return e.Int(-0x7fffffff) }, []byte{0x6d, 0xff, 0xff, 0xff, 0x7f})
	assertEncoded(t, InlineContainerTypeNone, func(e *Encoder) error { return e.Int(0xffffffff) }, []byte{0x6c, 0xff, 0xff, 0xff, 0xff})
	assertEncoded(t, InlineContainerTypeNone, func(e *Encoder) error { return e.Int(-0xffffffff) }, []byte{0x6d, 0xff, 0xff, 0xff, 0xff})
}

func TestEncodeInt49(t *testing.T) {
	assertEncoded(t, InlineContainerTypeNone, func(e *Encoder) error { return e.Int(0x100000000) }, []byte{0x66, 0x90, 0x80, 0x80, 0x80, 0})
}

func TestEncodeInt64(t *testing.T) {
	assertEncoded(t, InlineContainerTypeNone, func(e *Encoder) error { return e.Int(0x7fffffffffffffff) }, []byte{0x6e, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x7f})
	assertEncoded(t, InlineContainerTypeNone, func(e *Encoder) error { return e.Int(-0x7fffffffffffffff) }, []byte{0x6f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x7f})
}

func TestEncodeFloat32(t *testing.T) {
	assertEncoded(t, InlineContainerTypeNone, func(e *Encoder) error { return e.Float(1.5) }, []byte{0x70, 0x00, 0x00, 0xc0, 0x3f})
}

func TestEncodeFloat64(t *testing.T) {
	assertEncoded(t, InlineContainerTypeNone, func(e *Encoder) error { return e.Float(1.0123) }, []byte{0x71, 0x51, 0xda, 0x1b, 0x7c, 0x61, 0x32, 0xf0, 0x3f})
}

func TestEncodeFloatRounded(t *testing.T) {
	assertEncoded(t, InlineContainerTypeNone, func(e *Encoder) error { return e.FloatRounded(1.0123, 5) }, []byte{0x65, 0x12, 0xcf, 0x0b})
}

func newDate(year int, month int, day int) time.Time {
	location := time.UTC
	hour := 0
	minute := 0
	second := 0
	nanosecond := 0
	return time.Date(year, time.Month(month), day, hour, minute, second, nanosecond, location)
}

func newTime(hour int, minute int, second int, nanosecond int, timezone string) time.Time {
	year := 0
	month := 1
	day := 1
	location := time.UTC
	if len(timezone) > 0 {
		var err error
		location, err = time.LoadLocation(timezone)
		if err != nil {
			panic(err)
		}
	}
	return time.Date(year, time.Month(month), day, hour, minute, second, nanosecond, location)
}

func newTimestamp(year int, month int, day int, hour int, minute int, second int, nanosecond int, timezone string) time.Time {
	location := time.UTC
	if len(timezone) > 0 {
		var err error
		location, err = time.LoadLocation(timezone)
		if err != nil {
			panic(err)
		}
	}
	return time.Date(year, time.Month(month), day, hour, minute, second, nanosecond, location)
}

func TestEncodeDate(t *testing.T) {
	assertEncoded(t, InlineContainerTypeNone, func(e *Encoder) error { return e.Date(newDate(2000, 1, 1)) }, []byte{0x99, 0x21, 0x00, 0x00})
}

func TestEncodeTime(t *testing.T) {
	assertEncoded(t, InlineContainerTypeNone, func(e *Encoder) error { return e.Time(newTime(10, 10, 10, 0, "Asia/Tokyo")) },
		[]byte{0x9a, 0x50, 0x8a, 0x02, 0x0e, 'S', '/', 'T', 'o', 'k', 'y', 'o'})
}

func TestEncodeTimestamp(t *testing.T) {
	assertEncoded(t, InlineContainerTypeNone, func(e *Encoder) error {
		return e.Timestamp(newTimestamp(2020, 8, 30, 15, 33, 14, 19577323, "Asia/Singapore"))
	},
		[]byte{0x9b, 0x3b, 0xe1, 0xf3, 0xb8, 0x9e, 0xab, 0x12, 0x00, 0x50, 0x16, 'S', '/', 'S', 'i', 'n', 'g', 'a', 'p', 'o', 'r', 'e'})
}

func TestEncodeList(t *testing.T) {
	assertEncoded(t, InlineContainerTypeNone, func(e *Encoder) error { return e.ListBegin() }, []byte{0x75})
	assertEncoded(t, InlineContainerTypeNone, func(e *Encoder) error { return ShortCircuit(e.ListBegin(), e.ContainerEnd()) }, []byte{0x75, 0x79})
	assertEncoded(t, InlineContainerTypeNone, func(e *Encoder) error {
		return ShortCircuit(e.ListBegin(), e.Int(1), e.String("a"), e.ContainerEnd())
	}, []byte{0x75, 0x01, 0x81, 0x61, 0x79})
}

func TestEncodeMap(t *testing.T) {
	assertEncoded(t, InlineContainerTypeNone, func(e *Encoder) error { return e.MapBegin() }, []byte{0x76})
	assertEncoded(t, InlineContainerTypeNone, func(e *Encoder) error { return ShortCircuit(e.MapBegin(), e.ContainerEnd()) }, []byte{0x76, 0x79})
	assertEncoded(t, InlineContainerTypeNone, func(e *Encoder) error {
		return ShortCircuit(
			e.MapBegin(),
			e.String("1"),
			e.PositiveInt(1),
			e.String("2"),
			e.PositiveInt(2),
			e.ContainerEnd())
	},
		[]byte{0x76, 0x81, 0x31, 0x01, 0x81, 0x32, 0x02, 0x79})
}

func TestEncodeMetadata(t *testing.T) {
	assertEncoded(t, InlineContainerTypeNone, func(e *Encoder) error { return e.MetadataBegin() }, []byte{0x77})
	assertEncoded(t, InlineContainerTypeNone, func(e *Encoder) error { return ShortCircuit(e.MetadataBegin(), e.ContainerEnd()) }, []byte{0x77, 0x79})
	assertEncoded(t, InlineContainerTypeNone, func(e *Encoder) error {
		return ShortCircuit(
			e.MetadataBegin(),
			e.String("1"),
			e.PositiveInt(1),
			e.String("2"),
			e.PositiveInt(2),
			e.ContainerEnd())
	},
		[]byte{0x77, 0x81, 0x31, 0x01, 0x81, 0x32, 0x02, 0x79})
}

func TestEncodeInlineList(t *testing.T) {
	assertEncoded(t, InlineContainerTypeList, func(e *Encoder) error {
		return ShortCircuit(
			e.Nil(),
			e.Int(5),
			e.String(""))
	}, []byte{0x7e, 0x05, 0x80})
}

func TestEncodeInlineMap(t *testing.T) {
	assertEncoded(t, InlineContainerTypeMap, func(e *Encoder) error {
		return ShortCircuit(
			e.Int(1),
			e.String(""))
	}, []byte{0x01, 0x80})
}

func TestEncodeBytes(t *testing.T) {
	assertEncoded(t, InlineContainerTypeNone, func(e *Encoder) error { return e.Bytes([]byte{}) }, []byte{0x91, 0x00})
	assertEncoded(t, InlineContainerTypeNone, func(e *Encoder) error { return e.Bytes([]byte{1}) }, []byte{0x91, 0x01, 0x01})
	assertEncoded(t, InlineContainerTypeNone, func(e *Encoder) error { return e.Bytes([]byte{1, 2}) }, []byte{0x91, 0x02, 0x01, 0x02})
	assertEncoded(t, InlineContainerTypeNone, func(e *Encoder) error { return e.Bytes([]byte{1, 2, 3}) }, []byte{0x91, 0x03, 0x01, 0x02, 0x03})
}

func TestEncodeBytesLong(t *testing.T) {
	bytesLength := 500
	bytes := generateBytes(bytesLength)
	encoded := append([]byte{0x91, 0x83, 0x74}, bytes...)

	assertEncoded(t, InlineContainerTypeNone, func(e *Encoder) error { return e.Bytes(bytes) }, encoded)
}

func TestEncodeString(t *testing.T) {
	assertEncoded(t, InlineContainerTypeNone, func(e *Encoder) error { return e.String("") }, []byte{0x80})
	assertEncoded(t, InlineContainerTypeNone, func(e *Encoder) error { return e.String("0") }, []byte{0x81, 0x30})
	assertEncoded(t, InlineContainerTypeNone, func(e *Encoder) error { return e.String("01") }, []byte{0x82, 0x30, 0x31})
	assertEncoded(t, InlineContainerTypeNone, func(e *Encoder) error { return e.String("0123456789012345") }, []byte{
		0x90, 0x10, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37,
		0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35,
	})
}

func TestEncodeComment(t *testing.T) {
	assertEncoded(t, InlineContainerTypeNone, func(e *Encoder) error { return e.Comment("") }, []byte{0x93, 0x00})
	assertEncoded(t, InlineContainerTypeNone, func(e *Encoder) error { return e.Comment("0") }, []byte{0x93, 0x01, 0x30})
	assertEncoded(t, InlineContainerTypeNone, func(e *Encoder) error { return e.Comment("01") }, []byte{0x93, 0x02, 0x30, 0x31})
	assertEncoded(t, InlineContainerTypeNone, func(e *Encoder) error { return e.Comment("0123456789012345") }, []byte{
		0x93, 0x10, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37,
		0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35,
	})
}

func newURL(urlStr string) *url.URL {
	result, err := url.Parse(urlStr)
	if err != nil {
		panic(err)
	}
	return result
}

func TestEncodeURI(t *testing.T) {
	assertEncoded(t, InlineContainerTypeNone, func(e *Encoder) error { return e.URI(newURL("http://test.org")) }, []byte{
		0x92, 0x0f, 0x68, 0x74, 0x74, 0x70, 0x3a, 0x2f, 0x2f, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x6f, 0x72, 0x67,
	})
}

func TestEncodeStringInvalid(t *testing.T) {
	encoder := NewCbeEncoder(InlineContainerTypeNone, nil, 100)
	assertFailure(t, encoder.String(string([]byte{0x40, 0x81, 0x42, 0x43, 0x44})))
}

func TestEncodeCommentInvalid(t *testing.T) {
	encoder := NewCbeEncoder(InlineContainerTypeNone, nil, 100)
	assertFailure(t, encoder.Comment(string([]byte{0x40, 0x81, 0x42, 0x43, 0x44})))
	assertFailure(t, encoder.Comment("A comment\nwith a newline"))
}

func TestEncodeBytesTooLong(t *testing.T) {
	encoder := NewCbeEncoder(InlineContainerTypeNone, nil, 100)
	assertSuccess(t, encoder.BytesBegin(10))
	_, err := encoder.ArrayData(make([]byte, 11))
	assertFailure(t, err)
}

func TestEncodeBytesTooShort(t *testing.T) {
	encoder := NewCbeEncoder(InlineContainerTypeNone, nil, 100)
	assertSuccess(t, encoder.BytesBegin(10))
	_, err := encoder.ArrayData(make([]byte, 9))
	assertSuccess(t, err)
	assertFailure(t, encoder.End())
}

func TestEncodeStringTooLong(t *testing.T) {
	encoder := NewCbeEncoder(InlineContainerTypeNone, nil, 100)
	assertSuccess(t, encoder.StringBegin(6))
	_, err := encoder.ArrayData([]byte("abcdefg"))
	assertFailure(t, err)
}

func TestEncodeStringTooShort(t *testing.T) {
	encoder := NewCbeEncoder(InlineContainerTypeNone, nil, 100)
	assertSuccess(t, encoder.StringBegin(8))
	_, err := encoder.ArrayData([]byte("abcdefg"))
	assertSuccess(t, err)
	assertFailure(t, encoder.End())
}

func TestEncodeCommentTooLong(t *testing.T) {
	encoder := NewCbeEncoder(InlineContainerTypeNone, nil, 100)
	assertSuccess(t, encoder.CommentBegin(6))
	_, err := encoder.ArrayData([]byte("abcdefg"))
	assertFailure(t, err)
}

func TestEncodeURITooShort(t *testing.T) {
	encoder := NewCbeEncoder(InlineContainerTypeNone, nil, 100)
	assertSuccess(t, encoder.URIBegin(8))
	_, err := encoder.ArrayData([]byte("abcdefg"))
	assertSuccess(t, err)
	assertFailure(t, encoder.End())
}

func TestEncodeURITooLong(t *testing.T) {
	encoder := NewCbeEncoder(InlineContainerTypeNone, nil, 100)
	assertSuccess(t, encoder.URIBegin(6))
	_, err := encoder.ArrayData([]byte("abcdefg"))
	assertFailure(t, err)
}

func TestEncodeCommentTooShort(t *testing.T) {
	encoder := NewCbeEncoder(InlineContainerTypeNone, nil, 100)
	assertSuccess(t, encoder.CommentBegin(8))
	_, err := encoder.ArrayData([]byte("abcdefg"))
	assertSuccess(t, err)
	assertFailure(t, encoder.End())
}

func TestEncodeChangeArrayType(t *testing.T) {
	encoder := NewCbeEncoder(InlineContainerTypeNone, nil, 100)
	assertSuccess(t, encoder.BytesBegin(10))
	_, err := encoder.ArrayData(make([]byte, 5))
	assertSuccess(t, err)
	assertFailure(t, encoder.StringBegin(10))
}

func TestEncodeUnbalancedContainers(t *testing.T) {
	encoder := NewCbeEncoder(InlineContainerTypeNone, nil, 100)
	assertSuccess(t, encoder.ListBegin())
	assertFailure(t, encoder.End())
}

func TestEncodeCloseListTooManyTimes(t *testing.T) {
	encoder := NewCbeEncoder(InlineContainerTypeNone, nil, 100)
	assertSuccess(t, encoder.ListBegin())
	assertSuccess(t, encoder.ContainerEnd())
	assertFailure(t, encoder.ContainerEnd())
}

func TestEncodeCloseMapTooManyTimes(t *testing.T) {
	encoder := NewCbeEncoder(InlineContainerTypeNone, nil, 100)
	assertSuccess(t, encoder.MapBegin())
	assertSuccess(t, encoder.ContainerEnd())
	assertFailure(t, encoder.ContainerEnd())
}

func TestEncodeCloseMetadataTooManyTimes(t *testing.T) {
	encoder := NewCbeEncoder(InlineContainerTypeNone, nil, 100)
	assertSuccess(t, encoder.MetadataBegin())
	assertSuccess(t, encoder.ContainerEnd())
	assertFailure(t, encoder.ContainerEnd())
}

func TestEncodeCloseNoContainer(t *testing.T) {
	encoder := NewCbeEncoder(InlineContainerTypeNone, nil, 100)
	assertFailure(t, encoder.ContainerEnd())
}

func TestEncodeMapMissingValue(t *testing.T) {
	encoder := NewCbeEncoder(InlineContainerTypeNone, nil, 100)
	assertSuccess(t, encoder.MapBegin())
	assertSuccess(t, encoder.Int(1))
	assertFailure(t, encoder.ContainerEnd())
}

func TestEncodeMetadataMissingValue(t *testing.T) {
	encoder := NewCbeEncoder(InlineContainerTypeNone, nil, 100)
	assertSuccess(t, encoder.MetadataBegin())
	assertSuccess(t, encoder.Int(1))
	assertFailure(t, encoder.ContainerEnd())
}

func TestEncodeMapNilKey(t *testing.T) {
	encoder := NewCbeEncoder(InlineContainerTypeNone, nil, 100)
	assertSuccess(t, encoder.MapBegin())
	assertFailure(t, encoder.Nil())
}

func TestEncodeMapListKey(t *testing.T) {
	encoder := NewCbeEncoder(InlineContainerTypeNone, nil, 100)
	assertSuccess(t, encoder.MapBegin())
	assertFailure(t, encoder.ListBegin())
}

func TestEncodeMapMapKey(t *testing.T) {
	encoder := NewCbeEncoder(InlineContainerTypeNone, nil, 100)
	assertSuccess(t, encoder.MapBegin())
	assertFailure(t, encoder.MapBegin())
}

func TestEncodeMapBytesKey(t *testing.T) {
	encoder := NewCbeEncoder(InlineContainerTypeNone, nil, 100)
	assertSuccess(t, encoder.MapBegin())
	assertSuccess(t, encoder.Bytes([]byte{1, 2, 3}))
	assertSuccess(t, encoder.Bytes([]byte{4, 5, 6}))
	assertSuccess(t, encoder.ContainerEnd())
}

func TestEncodeMapWithComments(t *testing.T) {
	encoder := NewCbeEncoder(InlineContainerTypeNone, nil, 100)
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
	encoder := NewCbeEncoder(InlineContainerTypeNone, nil, 100)
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
	encoder := NewCbeEncoder(InlineContainerTypeNone, nil, 4)
	assertSuccess(t, encoder.ListBegin())
	assertSuccess(t, encoder.ListBegin())
	assertSuccess(t, encoder.MapBegin())
	assertSuccess(t, encoder.Bool(true))
	assertSuccess(t, encoder.ListBegin())
	assertFailure(t, encoder.MapBegin())
}

func TestEncodeToExternalBuffer(t *testing.T) {
	assertEncodesToExternalBuffer(t, InlineContainerTypeNone, nil, 1)
	assertEncodesToExternalBuffer(t, InlineContainerTypeNone, true, 1)
	assertEncodesToExternalBuffer(t, InlineContainerTypeNone, false, 1)
	assertEncodesToExternalBuffer(t, InlineContainerTypeNone, 0, 1)
	assertEncodesToExternalBuffer(t, InlineContainerTypeNone, -150, 2)
	assertEncodesToExternalBuffer(t, InlineContainerTypeNone, 1.1, 9)
	assertEncodesToExternalBuffer(t, InlineContainerTypeNone, 1.5, 5)
	assertEncodesToExternalBuffer(t, InlineContainerTypeNone, "test", 5)
	assertEncodesToExternalBuffer(t, InlineContainerTypeNone, "a longer string to test with", 30)
	assertEncodesToExternalBuffer(t, InlineContainerTypeNone, []byte{0x01}, 3)
	assertEncodesToExternalBuffer(t, InlineContainerTypeNone, time.Now(), 12)
	assertEncodesToExternalBuffer(t, InlineContainerTypeNone, map[interface{}]interface{}{"test": 1}, 8)
	assertEncodesToExternalBuffer(t, InlineContainerTypeNone, []interface{}{"test", 1}, 8)
}
