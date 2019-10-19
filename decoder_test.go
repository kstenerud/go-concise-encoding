package cbe

import (
	"testing"
)

func TestDecodeSmallInt(t *testing.T) {
	assertDecoded(t, ContainerTypeNone, []byte{0x00}, uint64(0))
	assertDecoded(t, ContainerTypeNone, []byte{0x01}, uint64(1))
	assertDecoded(t, ContainerTypeNone, []byte{100}, uint64(100))
	assertDecoded(t, ContainerTypeNone, []byte{0xff}, int64(-1))
	assertDecoded(t, ContainerTypeNone, []byte{0x9c}, int64(-100))
}

func TestDecodeInt8(t *testing.T) {
	assertDecoded(t, ContainerTypeNone, []byte{0x68, 0x6a}, uint64(106))
	assertDecoded(t, ContainerTypeNone, []byte{0x68, 0xff}, uint64(0xff))
	assertDecoded(t, ContainerTypeNone, []byte{0x69, 0xff}, int64(-0xff))
}

func TestDecodeInt16(t *testing.T) {
	assertDecoded(t, ContainerTypeNone, []byte{0x6a, 0x00, 0x01}, uint64(0x0100))
	assertDecoded(t, ContainerTypeNone, []byte{0x6a, 0xff, 0xff}, uint64(0xffff))
	assertDecoded(t, ContainerTypeNone, []byte{0x6b, 0x00, 0x01}, int64(-0x0100))
	assertDecoded(t, ContainerTypeNone, []byte{0x6b, 0xff, 0xff}, int64(-0xffff))
}

func TestDecodeInt32(t *testing.T) {
	assertDecoded(t, ContainerTypeNone, []byte{0x6c, 0x00, 0x00, 0x00, 0x01}, uint64(0x01000000))
	assertDecoded(t, ContainerTypeNone, []byte{0x6c, 0xff, 0xff, 0xff, 0xff}, uint64(0xffffffff))
	assertDecoded(t, ContainerTypeNone, []byte{0x6d, 0x00, 0x00, 0x00, 0x01}, int64(-0x01000000))
	assertDecoded(t, ContainerTypeNone, []byte{0x6d, 0xff, 0xff, 0xff, 0xff}, int64(-0xffffffff))
}

func TestDecodeInt64(t *testing.T) {
	assertDecoded(t, ContainerTypeNone, []byte{0x6e, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}, uint64(0x0100000000000000))
	assertDecoded(t, ContainerTypeNone, []byte{0x6e, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, uint64(0xffffffffffffffff))
	assertDecoded(t, ContainerTypeNone, []byte{0x6f, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}, int64(-0x0100000000000000))
	assertDecoded(t, ContainerTypeNone, []byte{0x6f, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80}, int64(-0x8000000000000000))
}

func TestDecodeStringSmall(t *testing.T) {
	for i := 0; i < 15; i++ {
		value := generateString(i)
		encoded := []byte{byte(0x80 + i)}
		encoded = append(encoded, []byte(value)...)
		assertDecoded(t, ContainerTypeNone, encoded, value)
	}
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

// func TestDecodePiecemeal2(t *testing.T) {
// 	value := []interface{}{0x100}
// 	encoded := []byte{0x77, 0x6a, 0x00, 0x01, 0x7b}
// 	assertDecodedPiecemeal(t, encoded, 1, 5, value)
// }

// func TestDecodePiecemeal3(t *testing.T) {
// 	value := []interface{}{1, 0x1234, 0x56789abc, uint64(0xfedcba9876543210)}
// 	encoded := []byte{0x77, 0x01, 0x6a, 0x34, 0x12, 0x6c, 0xbc, 0x9a, 0x78, 0x56, 0x6e, 0x10, 0x32, 0x54, 0x76, 0x98, 0xba, 0xdc, 0xfe, 0x7b}
// 	assertDecodedPiecemeal(t, encoded, 1, 20, value)
// }
