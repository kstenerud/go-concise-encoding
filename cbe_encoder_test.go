package cbe

import (
	"net/url"
	"testing"

	"github.com/kstenerud/go-cbe/rules"

	"github.com/kstenerud/go-compact-time"
)

func TestEncodePadding(t *testing.T) {
	assertEncoded(t, func(e *CBEEncoder) error { return e.Padding(1) }, []byte{cbeCodecVersion, 0x7f})
	assertEncoded(t, func(e *CBEEncoder) error { return e.Padding(2) }, []byte{cbeCodecVersion, 0x7f, 0x7f})
}

func TestEncodeNil(t *testing.T) {
	assertEncoded(t, func(e *CBEEncoder) error { return e.Nil() }, []byte{cbeCodecVersion, 0x7e})
}

func TestEncodeIntSmall(t *testing.T) {
	assertEncoded(t, func(e *CBEEncoder) error { return e.Int(0) }, []byte{cbeCodecVersion, 0})
	assertEncoded(t, func(e *CBEEncoder) error { return e.Int(1) }, []byte{cbeCodecVersion, 1})
	assertEncoded(t, func(e *CBEEncoder) error { return e.Int(99) }, []byte{cbeCodecVersion, 99})
	assertEncoded(t, func(e *CBEEncoder) error { return e.Int(100) }, []byte{cbeCodecVersion, 100})
	assertEncoded(t, func(e *CBEEncoder) error { return e.Int(-1) }, []byte{cbeCodecVersion, 255})
	assertEncoded(t, func(e *CBEEncoder) error { return e.Int(-99) }, []byte{cbeCodecVersion, 157})
	assertEncoded(t, func(e *CBEEncoder) error { return e.Int(-100) }, []byte{cbeCodecVersion, 156})
}

func TestEncodeInt8(t *testing.T) {
	assertEncoded(t, func(e *CBEEncoder) error { return e.Int(101) }, []byte{cbeCodecVersion, 0x68, 0x65})
	assertEncoded(t, func(e *CBEEncoder) error { return e.Int(-101) }, []byte{cbeCodecVersion, 0x69, 0x65})
	assertEncoded(t, func(e *CBEEncoder) error { return e.Int(255) }, []byte{cbeCodecVersion, 0x68, 0xff})
	assertEncoded(t, func(e *CBEEncoder) error { return e.Int(-255) }, []byte{cbeCodecVersion, 0x69, 0xff})
}

func TestEncodeInt16(t *testing.T) {
	assertEncoded(t, func(e *CBEEncoder) error { return e.Int(0x100) }, []byte{cbeCodecVersion, 0x6a, 0x00, 0x01})
	assertEncoded(t, func(e *CBEEncoder) error { return e.Int(0x7fff) }, []byte{cbeCodecVersion, 0x6a, 0xff, 0x7f})
	assertEncoded(t, func(e *CBEEncoder) error { return e.Int(-0x7fff) }, []byte{cbeCodecVersion, 0x6b, 0xff, 0x7f})
	assertEncoded(t, func(e *CBEEncoder) error { return e.Int(0xffff) }, []byte{cbeCodecVersion, 0x6a, 0xff, 0xff})
	assertEncoded(t, func(e *CBEEncoder) error { return e.Int(-0xffff) }, []byte{cbeCodecVersion, 0x6b, 0xff, 0xff})
}

func TestEncodeInt21(t *testing.T) {
	assertEncoded(t, func(e *CBEEncoder) error { return e.Int(0x10000) }, []byte{cbeCodecVersion, 0x66, 0x84, 0x80, 0x00})
}

func TestEncodeInt32(t *testing.T) {
	assertEncoded(t, func(e *CBEEncoder) error { return e.Int(0x7fffffff) }, []byte{cbeCodecVersion, 0x6c, 0xff, 0xff, 0xff, 0x7f})
	assertEncoded(t, func(e *CBEEncoder) error { return e.Int(-0x7fffffff) }, []byte{cbeCodecVersion, 0x6d, 0xff, 0xff, 0xff, 0x7f})
	assertEncoded(t, func(e *CBEEncoder) error { return e.Int(0xffffffff) }, []byte{cbeCodecVersion, 0x6c, 0xff, 0xff, 0xff, 0xff})
	assertEncoded(t, func(e *CBEEncoder) error { return e.Int(-0xffffffff) }, []byte{cbeCodecVersion, 0x6d, 0xff, 0xff, 0xff, 0xff})
}

func TestEncodeInt49(t *testing.T) {
	assertEncoded(t, func(e *CBEEncoder) error { return e.Int(0x100000000) }, []byte{cbeCodecVersion, 0x66, 0x90, 0x80, 0x80, 0x80, 0})
}

func TestEncodeInt64(t *testing.T) {
	assertEncoded(t, func(e *CBEEncoder) error { return e.Int(0x7fffffffffffffff) }, []byte{cbeCodecVersion, 0x6e, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x7f})
	assertEncoded(t, func(e *CBEEncoder) error { return e.Int(-0x7fffffffffffffff) }, []byte{cbeCodecVersion, 0x6f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x7f})
}

func TestEncodeFloat32(t *testing.T) {
	assertEncoded(t, func(e *CBEEncoder) error { return e.Float(1.5) }, []byte{cbeCodecVersion, 0x70, 0x00, 0x00, 0xc0, 0x3f})
}

func TestEncodeFloat64(t *testing.T) {
	assertEncoded(t, func(e *CBEEncoder) error { return e.Float(1.0123) }, []byte{cbeCodecVersion, 0x71, 0x51, 0xda, 0x1b, 0x7c, 0x61, 0x32, 0xf0, 0x3f})
}

func TestEncodeFloatRounded(t *testing.T) {
	assertEncoded(t, func(e *CBEEncoder) error { return e.FloatRounded(1.0123, 5) }, []byte{cbeCodecVersion, 0x65, 0x12, 0xcf, 0x0b})
}

func newDate(year int, month int, day int) *compact_time.Time {
	return compact_time.NewDate(year, month, day)
}

func newTime(hour int, minute int, second int, nanosecond int, areaLocation string) *compact_time.Time {
	return compact_time.NewTime(hour, minute, second, nanosecond, areaLocation)
}

func newTimestamp(year int, month int, day int, hour int, minute int, second int, nanosecond int, areaLocation string) *compact_time.Time {
	return compact_time.NewTimestamp(year, month, day, hour, minute, second, nanosecond, areaLocation)
}

// TODO: Go time
func TestEncodeDate(t *testing.T) {
	assertEncoded(t, func(e *CBEEncoder) error { return e.CompactTime(newDate(2000, 1, 1)) }, []byte{cbeCodecVersion, 0x99, 0x21, 0x00, 0x00})
}

func TestEncodeTime(t *testing.T) {
	assertEncoded(t, func(e *CBEEncoder) error { return e.CompactTime(newTime(10, 10, 10, 0, "Asia/Tokyo")) },
		[]byte{cbeCodecVersion, 0x9a, 0x50, 0x8a, 0x02, 0x0e, 'S', '/', 'T', 'o', 'k', 'y', 'o'})
}

func TestEncodeTimestamp(t *testing.T) {
	assertEncoded(t, func(e *CBEEncoder) error {
		return e.CompactTime(newTimestamp(2020, 8, 30, 15, 33, 14, 19577323, "Asia/Singapore"))
	},
		[]byte{cbeCodecVersion, 0x9b, 0x3b, 0xe1, 0xf3, 0xb8, 0x9e, 0xab, 0x12, 0x00, 0x50, 0x16, 'S', '/', 'S', 'i', 'n', 'g', 'a', 'p', 'o', 'r', 'e'})
}

func TestEncodeList(t *testing.T) {
	assertEncoded(t, func(e *CBEEncoder) error { return e.BeginList() }, []byte{cbeCodecVersion, 0x7a})
	assertEncoded(t, func(e *CBEEncoder) error { return ShortCircuit(e.BeginList(), e.EndContainer()) }, []byte{cbeCodecVersion, 0x7a, 0x7b})
	assertEncoded(t, func(e *CBEEncoder) error {
		return ShortCircuit(e.BeginList(), e.Int(1), e.String("a"), e.EndContainer())
	}, []byte{cbeCodecVersion, 0x7a, 0x01, 0x81, 0x61, 0x7b})
}

func TestEncodeMap(t *testing.T) {
	assertEncoded(t, func(e *CBEEncoder) error { return e.BeginMap() }, []byte{cbeCodecVersion, 0x79})
	assertEncoded(t, func(e *CBEEncoder) error { return ShortCircuit(e.BeginMap(), e.EndContainer()) }, []byte{cbeCodecVersion, 0x79, 0x7b})
	assertEncoded(t, func(e *CBEEncoder) error {
		return ShortCircuit(
			e.BeginMap(),
			e.String("1"),
			e.PositiveInt(1),
			e.String("2"),
			e.PositiveInt(2),
			e.EndContainer())
	},
		[]byte{cbeCodecVersion, 0x79, 0x81, 0x31, 0x01, 0x81, 0x32, 0x02, 0x7b})
}

func TestEncodeMetadata(t *testing.T) {
	assertEncoded(t, func(e *CBEEncoder) error { return e.BeginMetadata() }, []byte{cbeCodecVersion, 0x77})
	assertEncoded(t, func(e *CBEEncoder) error { return ShortCircuit(e.BeginMetadata(), e.EndContainer()) }, []byte{cbeCodecVersion, 0x77, 0x7b})
	assertEncoded(t, func(e *CBEEncoder) error {
		return ShortCircuit(
			e.BeginMetadata(),
			e.String("1"),
			e.PositiveInt(1),
			e.String("2"),
			e.PositiveInt(2),
			e.EndContainer())
	},
		[]byte{cbeCodecVersion, 0x77, 0x81, 0x31, 0x01, 0x81, 0x32, 0x02, 0x7b})
}

func TestEncodeComment(t *testing.T) {
	assertEncoded(t, func(e *CBEEncoder) error { return e.BeginComment() }, []byte{cbeCodecVersion, 0x76})
	assertEncoded(t, func(e *CBEEncoder) error { return ShortCircuit(e.BeginComment(), e.EndContainer()) }, []byte{cbeCodecVersion, 0x76, 0x7b})
	assertEncoded(t, func(e *CBEEncoder) error {
		return ShortCircuit(e.BeginComment(), e.String("a\n"), e.EndContainer())
	}, []byte{cbeCodecVersion, 0x76, 0x82, 0x61, 0x0a, 0x7b})
}

func TestEncodeCommentInvalid(t *testing.T) {
	encoder := NewCBEEncoder(InlineContainerTypeNone, nil, rules.DefaultLimits())
	assertSuccess(t, encoder.BeginComment())
	assertFailure(t, encoder.String(string([]byte{0x40, 0x81, 0x42, 0x43, 0x44})))
}

// TODO
// func TestEncodeInlineList(t *testing.T) {
// 	assertEncoded(t, InlineContainerTypeList, func(e *CBEEncoder) error {
// 		return ShortCircuit(
// 			e.Nil(),
// 			e.Int(5),
// 			e.String(""))
// 	}, []byte{codecVersion, 0x7e, 0x05, 0x80})
// }

// TODO
// func TestEncodeInlineMap(t *testing.T) {
// 	assertEncoded(t, InlineContainerTypeMap, func(e *CBEEncoder) error {
// 		return ShortCircuit(
// 			e.Int(1),
// 			e.String(""))
// 	}, []byte{codecVersion, 0x01, 0x80})
// }

func TestEncodeBytes(t *testing.T) {
	assertEncoded(t, func(e *CBEEncoder) error { return e.Bytes([]byte{}) }, []byte{cbeCodecVersion, 0x91, 0x00})
	assertEncoded(t, func(e *CBEEncoder) error { return e.Bytes([]byte{1}) }, []byte{cbeCodecVersion, 0x91, 0x02, 0x01})
	assertEncoded(t, func(e *CBEEncoder) error { return e.Bytes([]byte{1, 2}) }, []byte{cbeCodecVersion, 0x91, 0x04, 0x01, 0x02})
	assertEncoded(t, func(e *CBEEncoder) error { return e.Bytes([]byte{1, 2, 3}) }, []byte{cbeCodecVersion, 0x91, 0x06, 0x01, 0x02, 0x03})
}

func TestEncodeBytesLong(t *testing.T) {
	bytesLength := 500
	bytes := genBytes(bytesLength)
	encoded := append([]byte{cbeCodecVersion, 0x91, 0x87, 0x68}, bytes...)
	assertEncoded(t, func(e *CBEEncoder) error { return e.Bytes(bytes) }, encoded)
}

func TestEncodeString(t *testing.T) {
	assertEncoded(t, func(e *CBEEncoder) error { return e.String("") }, []byte{cbeCodecVersion, 0x80})
	assertEncoded(t, func(e *CBEEncoder) error { return e.String("0") }, []byte{cbeCodecVersion, 0x81, 0x30})
	assertEncoded(t, func(e *CBEEncoder) error { return e.String("01") }, []byte{cbeCodecVersion, 0x82, 0x30, 0x31})
	assertEncoded(t, func(e *CBEEncoder) error { return e.String("0123456789012345") }, []byte{
		cbeCodecVersion, 0x90, 0x20, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37,
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
	assertEncoded(t, func(e *CBEEncoder) error { return e.URI(newURL("http://test.org")) }, []byte{
		cbeCodecVersion, 0x92, 0x1e, 0x68, 0x74, 0x74, 0x70, 0x3a, 0x2f, 0x2f, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x6f, 0x72, 0x67,
	})
}

func TestEncodeStringInvalid(t *testing.T) {
	encoder := NewCBEEncoder(InlineContainerTypeNone, nil, rules.DefaultLimits())
	assertFailure(t, encoder.String(string([]byte{cbeCodecVersion, 0x40, 0x81, 0x42, 0x43, 0x44})))
}

func TestEncodeBytesTooLong(t *testing.T) {
	encoder := NewCBEEncoder(InlineContainerTypeNone, nil, rules.DefaultLimits())
	assertSuccess(t, encoder.BeginBytes())
	assertSuccess(t, encoder.beginChunk(10, true))
	_, err := encoder.ArrayData(make([]byte, 11))
	assertFailure(t, err)
}

func TestEncodeBytesTooShort(t *testing.T) {
	encoder := NewCBEEncoder(InlineContainerTypeNone, nil, rules.DefaultLimits())
	assertSuccess(t, encoder.BeginBytes())
	assertSuccess(t, encoder.BeginChunk(10, true))
	_, err := encoder.ArrayData(make([]byte, 9))
	assertSuccess(t, err)
	assertFailure(t, encoder.End())
}

func TestEncodeStringTooLong(t *testing.T) {
	encoder := NewCBEEncoder(InlineContainerTypeNone, nil, rules.DefaultLimits())
	assertSuccess(t, encoder.BeginString())
	assertSuccess(t, encoder.BeginChunk(6, true))
	_, err := encoder.ArrayData([]byte("abcdefg"))
	assertFailure(t, err)
}

func TestEncodeStringTooShort(t *testing.T) {
	encoder := NewCBEEncoder(InlineContainerTypeNone, nil, rules.DefaultLimits())
	assertSuccess(t, encoder.BeginString())
	assertSuccess(t, encoder.BeginChunk(8, true))
	_, err := encoder.ArrayData([]byte("abcdefg"))
	assertSuccess(t, err)
	assertFailure(t, encoder.End())
}

func TestEncodeURITooShort(t *testing.T) {
	encoder := NewCBEEncoder(InlineContainerTypeNone, nil, rules.DefaultLimits())
	assertSuccess(t, encoder.BeginURI())
	assertSuccess(t, encoder.BeginChunk(8, true))
	_, err := encoder.ArrayData([]byte("abcdefg"))
	assertSuccess(t, err)
	assertFailure(t, encoder.End())
}

func TestEncodeURITooLong(t *testing.T) {
	encoder := NewCBEEncoder(InlineContainerTypeNone, nil, rules.DefaultLimits())
	assertSuccess(t, encoder.BeginURI())
	assertSuccess(t, encoder.BeginChunk(6, true))
	_, err := encoder.ArrayData([]byte("abcdefg"))
	assertFailure(t, err)
}

func TestEncodeChangeArrayType(t *testing.T) {
	encoder := NewCBEEncoder(InlineContainerTypeNone, nil, rules.DefaultLimits())
	assertSuccess(t, encoder.BeginBytes())
	assertSuccess(t, encoder.BeginChunk(10, true))
	_, err := encoder.ArrayData(make([]byte, 5))
	assertSuccess(t, err)
	assertFailure(t, encoder.BeginString())
}

func TestEncodeUnbalancedContainers(t *testing.T) {
	encoder := NewCBEEncoder(InlineContainerTypeNone, nil, rules.DefaultLimits())
	assertSuccess(t, encoder.BeginList())
	assertFailure(t, encoder.End())
}

func TestEncodeCloseListTooManyTimes(t *testing.T) {
	encoder := NewCBEEncoder(InlineContainerTypeNone, nil, rules.DefaultLimits())
	assertSuccess(t, encoder.BeginList())
	assertSuccess(t, encoder.EndContainer())
	assertFailure(t, encoder.EndContainer())
}

func TestEncodeCloseMapTooManyTimes(t *testing.T) {
	encoder := NewCBEEncoder(InlineContainerTypeNone, nil, rules.DefaultLimits())
	assertSuccess(t, encoder.BeginMap())
	assertSuccess(t, encoder.EndContainer())
	assertFailure(t, encoder.EndContainer())
}

func TestEncodeCloseMetadataTooManyTimes(t *testing.T) {
	encoder := NewCBEEncoder(InlineContainerTypeNone, nil, rules.DefaultLimits())
	assertSuccess(t, encoder.BeginMetadata())
	assertSuccess(t, encoder.EndContainer())
	assertFailure(t, encoder.EndContainer())
}

func TestEncodeCloseNoContainer(t *testing.T) {
	encoder := NewCBEEncoder(InlineContainerTypeNone, nil, rules.DefaultLimits())
	assertFailure(t, encoder.EndContainer())
}

func TestEncodeMapMissingValue(t *testing.T) {
	encoder := NewCBEEncoder(InlineContainerTypeNone, nil, rules.DefaultLimits())
	assertSuccess(t, encoder.BeginMap())
	assertSuccess(t, encoder.Int(1))
	assertFailure(t, encoder.EndContainer())
}

func TestEncodeMetadataMissingValue(t *testing.T) {
	encoder := NewCBEEncoder(InlineContainerTypeNone, nil, rules.DefaultLimits())
	assertSuccess(t, encoder.BeginMetadata())
	assertSuccess(t, encoder.Int(1))
	assertFailure(t, encoder.EndContainer())
}

func TestEncodeMapNilKey(t *testing.T) {
	encoder := NewCBEEncoder(InlineContainerTypeNone, nil, rules.DefaultLimits())
	assertSuccess(t, encoder.BeginMap())
	assertFailure(t, encoder.Nil())
}

func TestEncodeMapListKey(t *testing.T) {
	encoder := NewCBEEncoder(InlineContainerTypeNone, nil, rules.DefaultLimits())
	assertSuccess(t, encoder.BeginMap())
	assertFailure(t, encoder.BeginList())
}

func TestEncodeMapMapKey(t *testing.T) {
	encoder := NewCBEEncoder(InlineContainerTypeNone, nil, rules.DefaultLimits())
	assertSuccess(t, encoder.BeginMap())
	assertFailure(t, encoder.BeginMap())
}

func TestEncodeMapBytesKey(t *testing.T) {
	encoder := NewCBEEncoder(InlineContainerTypeNone, nil, rules.DefaultLimits())
	assertSuccess(t, encoder.BeginMap())
	assertSuccess(t, encoder.Bytes([]byte{1, 2, 3}))
	assertSuccess(t, encoder.Bytes([]byte{4, 5, 6}))
	assertSuccess(t, encoder.EndContainer())
}

func TestEncodeMapWithComments(t *testing.T) {
	encoder := NewCBEEncoder(InlineContainerTypeNone, nil, rules.DefaultLimits())
	assertSuccess(t, encoder.BeginComment())
	assertSuccess(t, encoder.String("A comment before the map"))
	assertSuccess(t, encoder.EndContainer())
	assertSuccess(t, encoder.BeginMap())
	assertSuccess(t, encoder.BeginComment())
	assertSuccess(t, encoder.String("A comment"))
	assertSuccess(t, encoder.EndContainer())
	assertSuccess(t, encoder.String("a key"))
	assertSuccess(t, encoder.BeginComment())
	assertSuccess(t, encoder.String("Another comment"))
	assertSuccess(t, encoder.EndContainer())
	assertSuccess(t, encoder.Bool(true))
	assertSuccess(t, encoder.BeginComment())
	assertSuccess(t, encoder.String("Yet another comment"))
	assertSuccess(t, encoder.EndContainer())
	assertSuccess(t, encoder.EndContainer())
	assertSuccess(t, encoder.BeginComment())
	assertSuccess(t, encoder.String("A comment after the map"))
	assertSuccess(t, encoder.EndContainer())
}

func TestEncodeListWithComments(t *testing.T) {
	encoder := NewCBEEncoder(InlineContainerTypeNone, nil, rules.DefaultLimits())
	assertSuccess(t, encoder.BeginComment())
	assertSuccess(t, encoder.String("A comment before the list"))
	assertSuccess(t, encoder.EndContainer())
	assertSuccess(t, encoder.BeginList())
	assertSuccess(t, encoder.BeginComment())
	assertSuccess(t, encoder.String("A comment"))
	assertSuccess(t, encoder.EndContainer())
	assertSuccess(t, encoder.String("a string"))
	assertSuccess(t, encoder.BeginComment())
	assertSuccess(t, encoder.String("Another comment"))
	assertSuccess(t, encoder.EndContainer())
	assertSuccess(t, encoder.Bool(true))
	assertSuccess(t, encoder.BeginComment())
	assertSuccess(t, encoder.String("Yet another comment"))
	assertSuccess(t, encoder.EndContainer())
	assertSuccess(t, encoder.EndContainer())
	assertSuccess(t, encoder.BeginComment())
	assertSuccess(t, encoder.String("A comment after the list"))
	assertSuccess(t, encoder.EndContainer())
}

// TODO
// func TestEncodeContainerLimitExceeded(t *testing.T) {
// 	encoder := NewCBEEncoder(nil, 4)
// 	assertSuccess(t, encoder.BeginList())
// 	assertSuccess(t, encoder.BeginList())
// 	assertSuccess(t, encoder.BeginMap())
// 	assertSuccess(t, encoder.Bool(true))
// 	assertSuccess(t, encoder.BeginList())
// 	assertFailure(t, encoder.BeginMap())
// }

// TODO
// func TestEncodeToExternalBuffer(t *testing.T) {
// 	assertEncodesToExternalBuffer(t, nil, 2)
// 	assertEncodesToExternalBuffer(t, true, 2)
// 	assertEncodesToExternalBuffer(t, false, 2)
// 	assertEncodesToExternalBuffer(t, 0, 2)
// 	assertEncodesToExternalBuffer(t, -150, 3)
// 	assertEncodesToExternalBuffer(t, 1.1, 10)
// 	assertEncodesToExternalBuffer(t, 1.5, 6)
// 	assertEncodesToExternalBuffer(t, "test", 6)
// 	assertEncodesToExternalBuffer(t, "a longer string to test with", 31)
// 	assertEncodesToExternalBuffer(t, []byte{0x01}, 4)
// 	assertEncodesToExternalBuffer(t, time.Date(2000, time.Month(1), 1, 0, 0, 0, 0, time.UTC), 11)
// 	assertEncodesToExternalBuffer(t, map[interface{}]interface{}{"test": 1}, 9)
// 	assertEncodesToExternalBuffer(t, []interface{}{"test", 1}, 9)
// }
