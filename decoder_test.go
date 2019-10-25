package cbe

import (
	"testing"
	"time"
)

func TestDecodeSmallInt(t *testing.T) {
	assertDecoded(t, ContainerTypeNone, []byte{0x00}, int64(0))
	assertDecoded(t, ContainerTypeNone, []byte{0x01}, int64(1))
	assertDecoded(t, ContainerTypeNone, []byte{100}, int64(100))
	assertDecoded(t, ContainerTypeNone, []byte{0xff}, int64(-1))
	assertDecoded(t, ContainerTypeNone, []byte{0x9c}, int64(-100))
}

func TestDecodeInt8(t *testing.T) {
	assertDecoded(t, ContainerTypeNone, []byte{0x68, 0x6a}, int64(106))
	assertDecoded(t, ContainerTypeNone, []byte{0x68, 0xff}, int64(0xff))
	assertDecoded(t, ContainerTypeNone, []byte{0x69, 0xff}, int64(-0xff))
}

func TestDecodeInt21(t *testing.T) {
	assertDecoded(t, ContainerTypeNone, []byte{0x66, 0x84, 0x80, 0x00}, int64(0x10000))
	assertDecoded(t, ContainerTypeNone, []byte{0x67, 0x84, 0x80, 0x00}, int64(-0x10000))
}

func TestDecodeInt16(t *testing.T) {
	assertDecoded(t, ContainerTypeNone, []byte{0x6a, 0x00, 0x01}, int64(0x0100))
	assertDecoded(t, ContainerTypeNone, []byte{0x6a, 0xff, 0xff}, int64(0xffff))
	assertDecoded(t, ContainerTypeNone, []byte{0x6b, 0x00, 0x01}, int64(-0x0100))
	assertDecoded(t, ContainerTypeNone, []byte{0x6b, 0xff, 0xff}, int64(-0xffff))
}

func TestDecodeInt32(t *testing.T) {
	assertDecoded(t, ContainerTypeNone, []byte{0x6c, 0x00, 0x00, 0x00, 0x01}, int64(0x01000000))
	assertDecoded(t, ContainerTypeNone, []byte{0x6c, 0xff, 0xff, 0xff, 0xff}, int64(0xffffffff))
	assertDecoded(t, ContainerTypeNone, []byte{0x6d, 0x00, 0x00, 0x00, 0x01}, int64(-0x01000000))
	assertDecoded(t, ContainerTypeNone, []byte{0x6d, 0xff, 0xff, 0xff, 0xff}, int64(-0xffffffff))
}

func TestDecodeInt49(t *testing.T) {
	assertDecoded(t, ContainerTypeNone, []byte{0x66, 0x90, 0x80, 0x80, 0x80, 0x00}, int64(0x100000000))
	assertDecoded(t, ContainerTypeNone, []byte{0x67, 0x90, 0x80, 0x80, 0x80, 0x00}, int64(-0x100000000))
}

func TestDecodeInt64(t *testing.T) {
	assertDecoded(t, ContainerTypeNone, []byte{0x6e, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}, int64(0x0100000000000000))
	assertDecoded(t, ContainerTypeNone, []byte{0x6e, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, uint64(0xffffffffffffffff))
	assertDecoded(t, ContainerTypeNone, []byte{0x6f, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}, int64(-0x0100000000000000))
	assertDecoded(t, ContainerTypeNone, []byte{0x6f, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80}, int64(-0x8000000000000000))
}

func TestDecodeInt(t *testing.T) {
	assertDecoded(t, ContainerTypeNone, []byte{0x66, 0xff, 0xff, 0x7f}, int64(0x1fffff))
	assertDecoded(t, ContainerTypeNone, []byte{0x67, 0xff, 0xff, 0x7f}, int64(-0x1fffff))
}

func TestDecodeFloat(t *testing.T) {
	assertDecoded(t, ContainerTypeNone, []byte{0x70, 0x00, 0x00, 0x00, 0x00}, float32(0))
	assertDecoded(t, ContainerTypeNone, []byte{0x70, 0x22, 0x24, 0x6c, 0xc9}, float32(-967234.125))
	assertDecoded(t, ContainerTypeNone, []byte{0x71, 0x51, 0xda, 0x1b, 0x7c, 0x61, 0x32, 0xf0, 0x3f}, float64(1.0123))
	assertDecoded(t, ContainerTypeNone, []byte{0x65, 0x06, 0x01}, float64(0.1))
	assertDecoded(t, ContainerTypeNone, []byte{0x65, 0x0a, 0x13}, float64(0.19))
	assertDecoded(t, ContainerTypeNone, []byte{0x65, 0x82, 0x74, 0xdc, 0xe9, 0x87, 0x22}, float64(1.94659234e101))
}

func TestDecodeBool(t *testing.T) {
	assertDecoded(t, ContainerTypeNone, []byte{0x7c}, false)
	assertDecoded(t, ContainerTypeNone, []byte{0x7d}, true)
}

func asDate(year int, month int, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

func asTime(hour int, minute int, second int, nanosecond int, tz string) time.Time {
	location, err := time.LoadLocation(tz)
	if err != nil {
		panic(err)
	}
	return time.Date(0, 1, 1, hour, minute, second, nanosecond, location)
}

func asTimestamp(year int, month int, day int, hour int, minute int, second int, nanosecond int, tz string) time.Time {
	location, err := time.LoadLocation(tz)
	if err != nil {
		panic(err)
	}
	return time.Date(year, time.Month(month), day, hour, minute, second, nanosecond, location)
}

func TestDecodeTime(t *testing.T) {
	assertDecoded(t, ContainerTypeNone, []byte{0x99, 0x2f, 0x00, 0x1e}, asDate(2015, 1, 15))
	assertDecoded(t, ContainerTypeNone, []byte{0x9a, 0xbb, 0xce, 0x8a, 0x3e}, asTime(23, 14, 43, 1000000000, "Etc/UTC"))
	assertDecoded(t, ContainerTypeNone, []byte{0x9a, 0xba, 0xce, 0x8a, 0x3e, 0x10, 'E', '/', 'B', 'e', 'r', 'l', 'i', 'n'}, asTime(23, 14, 43, 1000000000, "Europe/Berlin"))
	assertDecoded(t, ContainerTypeNone, []byte{0x9b, 0x03, 0xa6, 0x5d, 0x1b, 0x00, 0x00, 0x00, 0x04, 0x33}, asTimestamp(1955, 11, 11, 22, 38, 0, 1, "Etc/UTC"))
	assertDecoded(t, ContainerTypeNone, []byte{0x9b, 0x40, 0x56, 0xd0, 0x0a, 0x3a, 0x1a, 'M', '/', 'L', 'o', 's', '_', 'A', 'n', 'g', 'e', 'l', 'e', 's'}, asTimestamp(1985, 10, 26, 1, 22, 16, 0, "America/Los_Angeles"))
}

func TestDecodeNil(t *testing.T) {
	assertDecoded(t, ContainerTypeNone, []byte{0x7e}, NilValue)
}

func TestDecodePadding(t *testing.T) {
	assertDecoded(t, ContainerTypeNone, []byte{0x7f}, nil)
	assertDecoded(t, ContainerTypeNone, []byte{0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f}, nil)
}

func TestDecodeStringSmall(t *testing.T) {
	for i := 0; i < 15; i++ {
		value := generateString(i)
		encoded := []byte{byte(0x80 + i)}
		encoded = append(encoded, []byte(value)...)
		assertDecoded(t, ContainerTypeNone, encoded, value)
	}
}

func asList(values ...interface{}) []interface{} {
	return values
}

func asMap(values ...interface{}) map[interface{}]interface{} {
	result := make(map[interface{}]interface{})
	var key interface{}
	for i, v := range values {
		if i&1 == 0 {
			key = v
		} else {
			result[key] = v
		}
	}
	return result
}

func TestDecodeList(t *testing.T) {
	assertDecoded(t, ContainerTypeNone, []byte{0x77, 0x7b}, asList())
	assertDecoded(t, ContainerTypeNone, []byte{0x77, 0x01, 0x7b}, asList(1))
	assertDecoded(t, ContainerTypeNone, []byte{0x77, 0x81, 0x31, 0x01, 0x7b}, asList("1", 1))
}

func TestDecodeMap(t *testing.T) {
	assertDecoded(t, ContainerTypeNone, []byte{0x78, 0x7b}, asMap())
	assertDecoded(t, ContainerTypeNone, []byte{0x78, 0x81, 0x31, 0x01, 0x7b}, asMap("1", 1))
	assertDecoded(t, ContainerTypeNone, []byte{0x78, 0x81, 0x31, 0x01, 0x81, 0x32, 0x02, 0x7b}, asMap("1", 1, "2", 2))

	assertDecoded(t, ContainerTypeNone, []byte{0x79, 0x81, 0x31, 0x01, 0x7b}, asMap("1", 1))
	assertDecoded(t, ContainerTypeNone, []byte{0x7a, 0x81, 0x31, 0x01, 0x7b}, asMap("1", 1))
}

func TestDecodeString0(t *testing.T) {
	value := generateString(0)
	encoded := []byte{0x90, 0x00}
	encoded = append(encoded, []byte(value)...)
	assertDecoded(t, ContainerTypeNone, encoded, value)
}

func TestDecodeString1(t *testing.T) {
	value := generateString(1)
	encoded := []byte{0x90, 0x01}
	encoded = append(encoded, []byte(value)...)
	assertDecoded(t, ContainerTypeNone, encoded, value)
}

func TestDecodeString16(t *testing.T) {
	value := generateString(16)
	encoded := []byte{0x90, 0x10}
	encoded = append(encoded, []byte(value)...)
	assertDecoded(t, ContainerTypeNone, encoded, value)
}

func TestDecodeString128(t *testing.T) {
	value := generateString(128)
	encoded := []byte{0x90, 0x81, 0x00}
	encoded = append(encoded, []byte(value)...)
	assertDecoded(t, ContainerTypeNone, encoded, value)
}

func TestDecodeString10000(t *testing.T) {
	value := generateString(10000)
	encoded := []byte{0x90, 0xce, 0x10}
	encoded = append(encoded, []byte(value)...)
	assertDecoded(t, ContainerTypeNone, encoded, value)
}

func TestDecodeComment0(t *testing.T) {
	value := generateString(0)
	encoded := []byte{0x93, 0x00}
	encoded = append(encoded, []byte(value)...)
	assertDecoded(t, ContainerTypeNone, encoded, value)
}

func TestDecodeComment1(t *testing.T) {
	value := generateString(1)
	encoded := []byte{0x93, 0x01}
	encoded = append(encoded, []byte(value)...)
	assertDecoded(t, ContainerTypeNone, encoded, value)
}

func TestDecodeComment16(t *testing.T) {
	value := generateString(16)
	encoded := []byte{0x93, 0x10}
	encoded = append(encoded, []byte(value)...)
	assertDecoded(t, ContainerTypeNone, encoded, value)
}

func TestDecodeComment128(t *testing.T) {
	value := generateString(128)
	encoded := []byte{0x93, 0x81, 0x00}
	encoded = append(encoded, []byte(value)...)
	assertDecoded(t, ContainerTypeNone, encoded, value)
}

func TestDecodeComment10000(t *testing.T) {
	value := generateString(10000)
	encoded := []byte{0x93, 0xce, 0x10}
	encoded = append(encoded, []byte(value)...)
	assertDecoded(t, ContainerTypeNone, encoded, value)
}

func TestDecodeBytes0(t *testing.T) {
	value := generateString(0)
	encoded := []byte{0x90, 0x00}
	encoded = append(encoded, []byte(value)...)
	assertDecoded(t, ContainerTypeNone, encoded, value)
}

func TestDecodeBytes1(t *testing.T) {
	value := generateString(1)
	encoded := []byte{0x91, 0x01}
	encoded = append(encoded, []byte(value)...)
	assertDecoded(t, ContainerTypeNone, encoded, value)
}

func TestDecodeBytes16(t *testing.T) {
	value := generateString(16)
	encoded := []byte{0x91, 0x10}
	encoded = append(encoded, []byte(value)...)
	assertDecoded(t, ContainerTypeNone, encoded, value)
}

func TestDecodeBytes128(t *testing.T) {
	value := generateString(128)
	encoded := []byte{0x91, 0x81, 0x00}
	encoded = append(encoded, []byte(value)...)
	assertDecoded(t, ContainerTypeNone, encoded, value)
}

func TestDecodeBytes10000(t *testing.T) {
	value := generateString(10000)
	encoded := []byte{0x91, 0xce, 0x10}
	encoded = append(encoded, []byte(value)...)
	assertDecoded(t, ContainerTypeNone, encoded, value)
}

func TestDecodeStringInvalid(t *testing.T) {
	assertFailure(t, tryDecode(9, []byte{0x83, 0x40, 0x81, 0x41}))
}

func TestDecodeCommentInvalid(t *testing.T) {
	assertFailure(t, tryDecode(9, []byte{0x93, 0x03, 0x40, 0x01, 0x41}))
}

func TestDecodeBytesTooShort(t *testing.T) {
	assertFailure(t, tryDecode(9, []byte{0x91, 0x08, 0x00}))
}

func TestDecodeBytesTooLong(t *testing.T) {
	assertFailure(t, tryDecode(9, []byte{0x91, 0x01, 0x00, 0x00}))
}

func TestDecodeStringTooShort(t *testing.T) {
	assertFailure(t, tryDecode(9, []byte{0x82, 0x00}))
}

func TestDecodeStringTooLong(t *testing.T) {
	assertFailure(t, tryDecode(9, []byte{0x82, 0x40, 0x41, 0x42}))
}

func TestDecodeCommentTooShort(t *testing.T) {
	assertFailure(t, tryDecode(9, []byte{0x93, 0x08, 0x00}))
}

func TestDecodeCommentTooLong(t *testing.T) {
	assertFailure(t, tryDecode(9, []byte{0x93, 0x01, 0x40, 0x41}))
}

func TestDecodeMapNilKey(t *testing.T) {
	assertFailure(t, tryDecode(9, []byte{0x78, 0x6f, 0x00, 0x7b}))
}

func TestDecodeMapListKey(t *testing.T) {
	assertFailure(t, tryDecode(9, []byte{0x78, 0x93, 0x95, 0x00, 0x7b}))
}

func TestDecodeMapMapKey(t *testing.T) {
	assertFailure(t, tryDecode(9, []byte{0x78, 0x94, 0x95, 0x00, 0x7b}))
}

func TestDecodeUnbalancedContainers(t *testing.T) {
	assertFailure(t, tryDecode(9, []byte{0x77}))
}

func TestDecodeListClosedTooManyTimes(t *testing.T) {
	assertFailure(t, tryDecode(9, []byte{0x77, 0x7b, 0x7b}))
}

func TestDecodeMapClosedTooManyTimes(t *testing.T) {
	assertFailure(t, tryDecode(9, []byte{0x78, 0x7b, 0x7b}))
}

func TestDecodeCloseNoContainer(t *testing.T) {
	assertFailure(t, tryDecode(9, []byte{0x7b}))
}

func TestDecodeMapMissingValue(t *testing.T) {
	assertFailure(t, tryDecode(9, []byte{0x78, 0x00, 0x7b}))
}

func TestDecodeContainerLimitExceeded(t *testing.T) {
	assertFailure(t, tryDecode(4, []byte{0x77, 0x77, 0x78, 0x00, 0x77, 0x77, 0x7b, 0x7b, 0x7b, 0x7b, 0x7b}))
}

func TestDecodeInlineList(t *testing.T) {
	assertDecoded(t, ContainerTypeList, []byte{0x00, 0x01}, []interface{}{0, 1})
}

func TestDecodeInlineMap(t *testing.T) {
	assertDecoded(t, ContainerTypeUnorderedMap, []byte{0x00, 0x01}, map[interface{}]interface{}{0: 1})
}

func TestDecodePiecemeal(t *testing.T) {
	value := []interface{}{1}
	encoded := []byte{0x77, 0x01, 0x7b}
	assertDecodedPiecemeal(t, encoded, 1, 3, value)
}

func TestDecodePiecemeal2(t *testing.T) {
	value := []interface{}{0x100}
	encoded := []byte{0x77, 0x6a, 0x00, 0x01, 0x7b}
	assertDecodedPiecemeal(t, encoded, 1, 5, value)
}

func TestDecodePiecemeal3(t *testing.T) {
	value := []interface{}{1, 0x1234, 0x56789abc, uint64(0xfedcba9876543210)}
	encoded := []byte{0x77, 0x01, 0x6a, 0x34, 0x12, 0x6c, 0xbc, 0x9a, 0x78, 0x56, 0x6e, 0x10, 0x32, 0x54, 0x76, 0x98, 0xba, 0xdc, 0xfe, 0x7b}
	assertDecodedPiecemeal(t, encoded, 1, 20, value)
}
