package cbe

// import (
// 	"testing"
// )

// func TestDecodeSmallInt(t *testing.T) {
// 	assertDecoded(t, []byte{0x00}, uint64(0))
// 	assertDecoded(t, []byte{0x01}, uint64(1))
// 	assertDecoded(t, []byte{105}, uint64(105))
// 	assertDecoded(t, []byte{0xff}, int64(-1))
// 	assertDecoded(t, []byte{0x96}, int64(-106))
// }

// func TestDecodeInt8(t *testing.T) {
// 	assertDecoded(t, []byte{0x6a, 0x6a}, uint64(106))
// 	assertDecoded(t, []byte{0x6a, 0xff}, uint64(0xff))
// 	assertDecoded(t, []byte{0x7a, 0xff}, int64(-0xff))
// }

// func TestDecodeInt16(t *testing.T) {
// 	assertDecoded(t, []byte{0x6b, 0x00, 0x01}, uint64(0x0100))
// 	assertDecoded(t, []byte{0x6b, 0xff, 0xff}, uint64(0xffff))
// 	assertDecoded(t, []byte{0x7b, 0x00, 0x01}, int64(-0x0100))
// 	assertDecoded(t, []byte{0x7b, 0xff, 0xff}, int64(-0xffff))
// }

// func TestDecodeInt32(t *testing.T) {
// 	assertDecoded(t, []byte{0x6c, 0x00, 0x00, 0x00, 0x01}, uint64(0x01000000))
// 	assertDecoded(t, []byte{0x6c, 0xff, 0xff, 0xff, 0xff}, uint64(0xffffffff))
// 	assertDecoded(t, []byte{0x7c, 0x00, 0x00, 0x00, 0x01}, int64(-0x01000000))
// 	assertDecoded(t, []byte{0x7c, 0xff, 0xff, 0xff, 0xff}, int64(-0xffffffff))
// }

// func TestDecodeInt64(t *testing.T) {
// 	assertDecoded(t, []byte{0x6d, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}, uint64(0x0100000000000000))
// 	assertDecoded(t, []byte{0x6d, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, uint64(0xffffffffffffffff))
// 	assertDecoded(t, []byte{0x7d, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}, int64(-0x0100000000000000))
// 	assertDecoded(t, []byte{0x7d, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80}, int64(-0x8000000000000000))
// }

// func TestDecodeStringSmall(t *testing.T) {
// 	for i := 0; i < 15; i++ {
// 		value := generateString(i)
// 		encoded := []byte{byte(0x80 + i)}
// 		encoded = append(encoded, []byte(value)...)
// 		assertDecoded(t, encoded, value)
// 	}
// }

// func TestDecodeString0(t *testing.T) {
// 	value := generateString(0)
// 	encoded := []byte{0x90, 0 << 2}
// 	encoded = append(encoded, []byte(value)...)
// 	assertDecoded(t, encoded, value)
// }

// func TestDecodeString1(t *testing.T) {
// 	value := generateString(1)
// 	encoded := []byte{0x90, 1 << 2}
// 	encoded = append(encoded, []byte(value)...)
// 	assertDecoded(t, encoded, value)
// }

// func TestDecodeString16(t *testing.T) {
// 	value := generateString(16)
// 	encoded := []byte{0x90, 16 << 2}
// 	encoded = append(encoded, []byte(value)...)
// 	assertDecoded(t, encoded, value)
// }

// func TestDecodeString64(t *testing.T) {
// 	value := generateString(64)
// 	encoded := []byte{0x90, 0x01, 0x01}
// 	encoded = append(encoded, []byte(value)...)
// 	assertDecoded(t, encoded, value)
// }

// func TestDecodeString16384(t *testing.T) {
// 	value := generateString(16384)
// 	encoded := []byte{0x90, 0x02, 0x00, 0x01, 0x00}
// 	encoded = append(encoded, []byte(value)...)
// 	assertDecoded(t, encoded, value)
// }

// func TestDecodeComment0(t *testing.T) {
// 	value := generateString(0)
// 	encoded := []byte{0x92, 0 << 2}
// 	encoded = append(encoded, []byte(value)...)
// 	assertDecoded(t, encoded, value)
// }

// func TestDecodeComment1(t *testing.T) {
// 	value := generateString(1)
// 	encoded := []byte{0x92, 1 << 2}
// 	encoded = append(encoded, []byte(value)...)
// 	assertDecoded(t, encoded, value)
// }

// func TestDecodeComment16(t *testing.T) {
// 	value := generateString(16)
// 	encoded := []byte{0x92, 16 << 2}
// 	encoded = append(encoded, []byte(value)...)
// 	assertDecoded(t, encoded, value)
// }

// func TestDecodeComment64(t *testing.T) {
// 	value := generateString(64)
// 	encoded := []byte{0x92, 0x01, 0x01}
// 	encoded = append(encoded, []byte(value)...)
// 	assertDecoded(t, encoded, value)
// }

// func TestDecodeComment16384(t *testing.T) {
// 	value := generateString(16384)
// 	encoded := []byte{0x92, 0x02, 0x00, 0x01, 0x00}
// 	encoded = append(encoded, []byte(value)...)
// 	assertDecoded(t, encoded, value)
// }

// func TestDecodeBinary0(t *testing.T) {
// 	value := generateBinary(0)
// 	encoded := []byte{0x91, 0 << 2}
// 	encoded = append(encoded, []byte(value)...)
// 	assertDecoded(t, encoded, value)
// }

// func TestDecodeBinary1(t *testing.T) {
// 	value := generateBinary(1)
// 	encoded := []byte{0x91, 1 << 2}
// 	encoded = append(encoded, []byte(value)...)
// 	assertDecoded(t, encoded, value)
// }

// func TestDecodeBinary16(t *testing.T) {
// 	value := generateBinary(16)
// 	encoded := []byte{0x91, 16 << 2}
// 	encoded = append(encoded, []byte(value)...)
// 	assertDecoded(t, encoded, value)
// }

// func TestDecodeBinary64(t *testing.T) {
// 	value := generateBinary(64)
// 	encoded := []byte{0x91, 0x01, 0x01}
// 	encoded = append(encoded, []byte(value)...)
// 	assertDecoded(t, encoded, value)
// }

// func TestDecodeBinary16384(t *testing.T) {
// 	value := generateBinary(16384)
// 	encoded := []byte{0x91, 0x02, 0x00, 0x01, 0x00}
// 	encoded = append(encoded, []byte(value)...)
// 	assertDecoded(t, encoded, value)
// }
