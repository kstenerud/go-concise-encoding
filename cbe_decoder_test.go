package cbe

import (
	"testing"

	"github.com/kstenerud/go-cbe/rules"
	"github.com/kstenerud/go-equivalence"
)

func TestCBEDecodeSmallInt(t *testing.T) {
	assertDecodedEncoded(t, []byte{cbeCodecVersion, 0x00}, Tokens(0))
	assertDecodedEncoded(t, []byte{cbeCodecVersion, 0x01}, Tokens(1))
	assertDecodedEncoded(t, []byte{cbeCodecVersion, 100}, Tokens(100))
	assertDecodedEncoded(t, []byte{cbeCodecVersion, 0xff}, Tokens(neg(1)))
	assertDecodedEncoded(t, []byte{cbeCodecVersion, 0x9c}, Tokens(neg(100)))
}

func TestCBEDecodeInt8(t *testing.T) {
	assertDecodedEncoded(t, []byte{cbeCodecVersion, 0x68, 0x65}, Tokens(101))
	assertDecodedEncoded(t, []byte{cbeCodecVersion, 0x68, 0xff}, Tokens(0xff))
	assertDecodedEncoded(t, []byte{cbeCodecVersion, 0x69, 0xff}, Tokens(neg(0xff)))
}

func TestCBEDecodeInt21(t *testing.T) {
	assertDecodedEncoded(t, []byte{cbeCodecVersion, 0x66, 0x84, 0x80, 0x00}, Tokens(0x10000))
	assertDecodedEncoded(t, []byte{cbeCodecVersion, 0x67, 0x84, 0x80, 0x00}, Tokens(neg(0x10000)))
}

func TestCBEDecodeInt16(t *testing.T) {
	assertDecodedEncoded(t, []byte{cbeCodecVersion, 0x6a, 0x00, 0x01}, Tokens(0x0100))
	assertDecodedEncoded(t, []byte{cbeCodecVersion, 0x6a, 0xff, 0xff}, Tokens(0xffff))
	assertDecodedEncoded(t, []byte{cbeCodecVersion, 0x6b, 0x00, 0x01}, Tokens(neg(0x0100)))
	assertDecodedEncoded(t, []byte{cbeCodecVersion, 0x6b, 0xff, 0xff}, Tokens(neg(0xffff)))
}

func TestCBEDecodeInt32(t *testing.T) {
	assertDecodedEncoded(t, []byte{cbeCodecVersion, 0x6c, 0x00, 0x00, 0x00, 0x01}, Tokens(0x01000000))
	assertDecodedEncoded(t, []byte{cbeCodecVersion, 0x6c, 0xff, 0xff, 0xff, 0xff}, Tokens(0xffffffff))
	assertDecodedEncoded(t, []byte{cbeCodecVersion, 0x6d, 0x00, 0x00, 0x00, 0x01}, Tokens(neg(0x01000000)))
	assertDecodedEncoded(t, []byte{cbeCodecVersion, 0x6d, 0xff, 0xff, 0xff, 0xff}, Tokens(neg(0xffffffff)))
}

func TestCBEDecodeInt49(t *testing.T) {
	assertDecodedEncoded(t, []byte{cbeCodecVersion, 0x66, 0x90, 0x80, 0x80, 0x80, 0x00}, Tokens(0x100000000))
	assertDecodedEncoded(t, []byte{cbeCodecVersion, 0x67, 0x90, 0x80, 0x80, 0x80, 0x00}, Tokens(neg(0x100000000)))
}

func TestCBEDecodeInt64(t *testing.T) {
	assertDecodedEncoded(t, []byte{cbeCodecVersion, 0x6e, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}, Tokens(0x0100000000000000))
	assertDecodedEncoded(t, []byte{cbeCodecVersion, 0x6e, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, Tokens(uint64(0xffffffffffffffff)))
	assertDecodedEncoded(t, []byte{cbeCodecVersion, 0x6f, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}, Tokens(neg(0x0100000000000000)))
	assertDecodedEncoded(t, []byte{cbeCodecVersion, 0x6f, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80}, Tokens(neg(0x8000000000000000)))
}

func TestCBEDecodeInt(t *testing.T) {
	assertDecodedEncoded(t, []byte{cbeCodecVersion, 0x66, 0xff, 0xff, 0x7f}, Tokens(0x1fffff))
	assertDecodedEncoded(t, []byte{cbeCodecVersion, 0x67, 0xff, 0xff, 0x7f}, Tokens(neg(0x1fffff)))
}

func TestCBEDecodeFloat(t *testing.T) {
	assertDecodedEncoded(t, []byte{cbeCodecVersion, 0x70, 0x00, 0x00, 0x00, 0x00}, Tokens(0.0))
	assertDecodedEncoded(t, []byte{cbeCodecVersion, 0x70, 0x22, 0x24, 0x6c, 0xc9}, Tokens(-967234.125))
	assertDecodedEncoded(t, []byte{cbeCodecVersion, 0x71, 0x51, 0xda, 0x1b, 0x7c, 0x61, 0x32, 0xf0, 0x3f}, Tokens(1.0123))
	assertDecoded(t, []byte{cbeCodecVersion, 0x65, 0x06, 0x01}, Tokens(0.1))
	assertDecoded(t, []byte{cbeCodecVersion, 0x65, 0x0a, 0x13}, Tokens(0.19))
	assertDecoded(t, []byte{cbeCodecVersion, 0x65, 0x82, 0x74, 0xdc, 0xe9, 0x87, 0x22}, Tokens(1.94659234e101))
}

func TestCBEDecodeBool(t *testing.T) {
	assertDecodedEncoded(t, []byte{cbeCodecVersion, 0x7c}, Tokens(false))
	assertDecodedEncoded(t, []byte{cbeCodecVersion, 0x7d}, Tokens(true))
}

func TestCBEDecodeTime(t *testing.T) {
	assertDecodedEncoded(t, []byte{cbeCodecVersion, 0x99, 0x2f, 0x00, 0x1e}, Tokens(Date(2015, 1, 15)))
	assertDecodedEncoded(t, []byte{cbeCodecVersion, 0x9a, 0xbb, 0xce, 0x4a, 0x06}, Tokens(Time(23, 14, 43, 100000000, "")))
	assertDecodedEncoded(t, []byte{cbeCodecVersion, 0x9a, 0xba, 0xce, 0x4a, 0x06, 0x10, 'E', '/', 'B', 'e', 'r', 'l', 'i', 'n'}, Tokens(Time(23, 14, 43, 100000000, "Europe/Berlin")))
	assertDecodedEncoded(t, []byte{cbeCodecVersion, 0x9b, 0x03, 0xa6, 0x5d, 0x1b, 0x00, 0x00, 0x00, 0x04, 0x33}, Tokens(TS(1955, 11, 11, 22, 38, 0, 1, "")))
	assertDecodedEncoded(t, []byte{cbeCodecVersion, 0x9b, 0x40, 0x56, 0xd0, 0x0a, 0x3a, 0x1a, 'M', '/', 'L', 'o', 's', '_', 'A', 'n', 'g', 'e', 'l', 'e', 's'}, Tokens(TS(1985, 10, 26, 1, 22, 16, 0, "America/Los_Angeles")))
}

func TestCBEDecodeNil(t *testing.T) {
	assertDecodedEncoded(t, []byte{cbeCodecVersion, 0x7e}, Tokens(Nil()))
}

func TestCBEDecodePadding(t *testing.T) {
	assertDecoded(t, []byte{cbeCodecVersion, 0x7f, 0x00}, Tokens(0))
	assertDecoded(t, []byte{cbeCodecVersion, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x00}, Tokens(0))
}

func TestCBEDecodeStringSmall(t *testing.T) {
	encoded := []byte{cbeCodecVersion, byte(0x80)}
	assertDecoded(t, encoded, Tokens(str(), chunk(0, true)))

	for i := 1; i < 15; i++ {
		value := genString(i)
		encoded = []byte{cbeCodecVersion, byte(0x80 + i)}
		encoded = append(encoded, []byte(value)...)
		assertDecoded(t, encoded, Tokens(str(), chunk(uint64(i), true), strBytes(value)))
	}
}

func TestCBEDecodeList(t *testing.T) {
	assertDecodedEncoded(t, []byte{cbeCodecVersion, 0x7a, 0x7b}, Tokens(list(), end()))
	assertDecodedEncoded(t, []byte{cbeCodecVersion, 0x7a, 0x01, 0x7b}, Tokens(list(), 1, end()))
	assertDecoded(t, []byte{cbeCodecVersion, 0x7a, 0x81, 0x31, 0x01, 0x7b}, Tokens(list(), str(), chunk(1, true), strBytes("1"), 1, end()))
}

func TestCBEDecodeMap(t *testing.T) {
	assertDecodedEncoded(t, []byte{cbeCodecVersion, 0x79, 0x7b}, Tokens(Map(), end()))
	assertDecoded(t, []byte{cbeCodecVersion, 0x79, 0x81, 0x31, 0x01, 0x7b}, Tokens(Map(), str(), chunk(1, true), strBytes("1"), 1, end()))
	assertDecoded(t, []byte{cbeCodecVersion, 0x79, 0x81, 0x31, 0x01, 0x81, 0x32, 0x02, 0x7b},
		Tokens(Map(), str(), chunk(1, true), strBytes("1"), 1, str(), chunk(1, true), strBytes("2"), 2, end()))
}

func TestCBEDecodeMarkup(t *testing.T) {
	assertDecodedEncoded(t, []byte{cbeCodecVersion, 0x78, 1, 0x7b, 0x7b}, Tokens(markup(), 1, end(), end()))
	assertDecoded(t, []byte{cbeCodecVersion, 0x78, 1, 2, 3, 0x7b, 0x81, 'a', 0x7b}, Tokens(markup(), 1, 2, 3, end(), str(), chunk(1, true), strBytes("a"), end()))
}

func TestCBEDecodeMetadata(t *testing.T) {
	assertDecodedEncoded(t, []byte{cbeCodecVersion, 0x77, 0x7b, 1}, Tokens(meta(), end(), 1))
	assertDecoded(t, []byte{cbeCodecVersion, 0x77, 0x81, 0x31, 0x01, 0x7b, 1}, Tokens(meta(), str(), chunk(1, true), strBytes("1"), 1, end(), 1))
	assertDecoded(t, []byte{cbeCodecVersion, 0x77, 0x81, 0x31, 0x01, 0x81, 0x32, 0x02, 0x7b, 1},
		Tokens(meta(), str(), chunk(1, true), strBytes("1"), 1, str(), chunk(1, true), strBytes("2"), 2, end(), 1))
}

func TestCBEDecodeComment(t *testing.T) {
	assertDecodedWithAutoEnd(t, []byte{cbeCodecVersion, 0x76, 0x7b}, Tokens(cmt(), end()))
	assertDecodedWithAutoEnd(t, []byte{cbeCodecVersion, 0x76, 0x81, 'a', 0x7b}, Tokens(cmt(), str(), chunk(1, true), strBytes("a"), end()))
	assertDecodedWithAutoEnd(t, []byte{cbeCodecVersion, 0x76, 0x81, 0x31, 0x76, 0x7b, 0x7b}, Tokens(cmt(), str(), chunk(1, true), strBytes("1"), cmt(), end(), end()))
}

func TestCBEDecodeString0(t *testing.T) {
	encoded := []byte{cbeCodecVersion, 0x90, 0x00}
	assertDecodedEncoded(t, encoded, Tokens(str(), chunk(0, true)))
}

func TestCBEDecodeString1(t *testing.T) {
	value := genString(1)
	encoded := []byte{cbeCodecVersion, 0x90, 0x02}
	encoded = append(encoded, []byte(value)...)
	assertDecodedEncoded(t, encoded, Tokens(str(), chunk(1, true), strBytes(value)))
}

func TestCBEDecodeString16(t *testing.T) {
	value := genString(16)
	encoded := []byte{cbeCodecVersion, 0x90, 0x20}
	encoded = append(encoded, []byte(value)...)
	assertDecodedEncoded(t, encoded, Tokens(str(), chunk(16, true), strBytes(value)))
}

func TestCBEDecodeString64(t *testing.T) {
	value := genString(64)
	encoded := []byte{cbeCodecVersion, 0x90, 0x81, 0x00}
	encoded = append(encoded, []byte(value)...)
	assertDecodedEncoded(t, encoded, Tokens(str(), chunk(64, true), strBytes(value)))
}

func TestCBEDecodeString5000(t *testing.T) {
	value := genString(5000)
	encoded := []byte{cbeCodecVersion, 0x90, 0xce, 0x10}
	encoded = append(encoded, []byte(value)...)
	assertDecodedEncoded(t, encoded, Tokens(str(), chunk(5000, true), strBytes(value)))
}

func TestCBEDecodeStringChunked(t *testing.T) {
	encoded := []byte{cbeCodecVersion, 0x90, 0x11}
	encoded = append(encoded, []byte(genString(8))...)
	encoded = append(encoded, 0x05)
	encoded = append(encoded, []byte(genString(2))...)
	encoded = append(encoded, 0x04)
	encoded = append(encoded, []byte(genString(2))...)
	assertDecoded(t, encoded, Tokens(str(), chunk(8, false), strBytes(genString(8)), chunk(2, false), strBytes(genString(2)), chunk(2, true), strBytes(genString(2))))
}

func TestCBEDecodeURI(t *testing.T) {
	value := "http://example.com"
	encoded := []byte{cbeCodecVersion, 0x92, 0x24}
	encoded = append(encoded, []byte(value)...)
	assertDecodedEncoded(t, encoded, Tokens(uri(), chunk(18, true), strBytes(value)))
}

func TestCBEDecodeURIChunked(t *testing.T) {
	encoded := []byte{cbeCodecVersion, 0x92, 0x11}
	encoded = append(encoded, []byte("http://e")...)
	encoded = append(encoded, 0x0d)
	encoded = append(encoded, []byte("xample")...)
	encoded = append(encoded, 0x08)
	encoded = append(encoded, []byte(".com")...)
	assertDecoded(t, encoded, Tokens(uri(), chunk(8, false), strBytes("http://e"), chunk(6, false), strBytes("xample"), chunk(4, true), strBytes(".com")))
}

func TestCBEDecodeBytes0(t *testing.T) {
	encoded := []byte{cbeCodecVersion, 0x91, 0x00}
	assertDecodedEncoded(t, encoded, Tokens(bin(), chunk(0, true)))
}

func TestCBEDecodeBytes1(t *testing.T) {
	value := genBytes(1)
	encoded := []byte{cbeCodecVersion, 0x91, 0x02}
	encoded = append(encoded, value...)
	assertDecodedEncoded(t, encoded, Tokens(bin(), chunk(1, true), value))
}

func TestCBEDecodeBytes16(t *testing.T) {
	value := genBytes(16)
	encoded := []byte{cbeCodecVersion, 0x91, 0x20}
	encoded = append(encoded, value...)
	assertDecodedEncoded(t, encoded, Tokens(bin(), chunk(16, true), value))
}

func TestCBEDecodeBytes64(t *testing.T) {
	value := genBytes(64)
	encoded := []byte{cbeCodecVersion, 0x91, 0x81, 0x00}
	encoded = append(encoded, []byte(value)...)
	assertDecodedEncoded(t, encoded, Tokens(bin(), chunk(64, true), value))
}

func TestCBEDecodeBytes5000(t *testing.T) {
	value := genBytes(5000)
	encoded := []byte{cbeCodecVersion, 0x91, 0xce, 0x10}
	encoded = append(encoded, []byte(value)...)
	assertDecodedEncoded(t, encoded, Tokens(bin(), chunk(5000, true), value))
}

func TestCBEDecodeBytesChunked(t *testing.T) {
	encoded := []byte{cbeCodecVersion, 0x91, 0x11}
	encoded = append(encoded, genBytes(8)...)
	encoded = append(encoded, 0x05)
	encoded = append(encoded, genBytes(2)...)
	encoded = append(encoded, 0x04)
	encoded = append(encoded, genBytes(2)...)
	assertDecoded(t, encoded, Tokens(bin(), chunk(8, false), genBytes(8), chunk(2, false), genBytes(2), chunk(2, true), genBytes(2)))
}

func TestCBEDecodeMarkerReference(t *testing.T) {
	assertDecodedEncoded(t, []byte{cbeCodecVersion, 0x7a, 0x97, 1, 0x22, 0x7b}, Tokens(list(), marker(), 1, 0x22, end()))
	assertDecodedEncoded(t, []byte{cbeCodecVersion, 0x7a, 0x97, 1, 0x22, 0x98, 1, 0x7b}, Tokens(list(), marker(), 1, 0x22, ref(), 1, end()))
}

func TestCBEDecodeStringInvalid(t *testing.T) {
	assertSuccess(t, tryDecode([]byte{cbeCodecVersion, 0x83, 0x40, 0x41, 0x41}))
	assertFailure(t, tryDecode([]byte{cbeCodecVersion, 0x83, 0x40, 0x81, 0x41}))
}

func TestCBEDecodeBytesTooShort(t *testing.T) {
	assertSuccess(t, tryDecode([]byte{cbeCodecVersion, 0x91, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}))
	assertFailure(t, tryDecode([]byte{cbeCodecVersion, 0x91, 0x10, 0x00}))
}

func TestCBEDecodeStringTooShort(t *testing.T) {
	assertSuccess(t, tryDecode([]byte{cbeCodecVersion, 0x82, 0x40, 0x40}))
	assertFailure(t, tryDecode([]byte{cbeCodecVersion, 0x82, 0x40}))
}

func TestCBEDecodeMapNilKey(t *testing.T) {
	assertSuccess(t, tryDecode([]byte{cbeCodecVersion, 0x79, 0x00, 0x00, 0x7b}))
	assertFailure(t, tryDecode([]byte{cbeCodecVersion, 0x79, 0x7f, 0x00, 0x7b}))
}

func TestCBEDecodeMapListKey(t *testing.T) {
	assertSuccess(t, tryDecode([]byte{cbeCodecVersion, 0x79, 0x00, 0x00, 0x7b}))
	assertFailure(t, tryDecode([]byte{cbeCodecVersion, 0x79, 0x7a, 0x7b, 0x00, 0x7b}))
}

func TestCBEDecodeMapMapKey(t *testing.T) {
	assertSuccess(t, tryDecode([]byte{cbeCodecVersion, 0x79, 0x00, 0x00, 0x7b}))
	assertFailure(t, tryDecode([]byte{cbeCodecVersion, 0x79, 0x79, 0x7b, 0x00, 0x7b}))
}

func TestCBEDecodeUnbalancedContainers(t *testing.T) {
	assertSuccess(t, tryDecode([]byte{cbeCodecVersion, 0x7a, 0x7b}))
	assertFailure(t, tryDecode([]byte{cbeCodecVersion, 0x7a}))
}

func TestCBEDecodeCloseNoContainer(t *testing.T) {
	assertFailure(t, tryDecode([]byte{cbeCodecVersion, 0x7b}))
}

func TestCBEDecodeMapMissingValue(t *testing.T) {
	assertSuccess(t, tryDecode([]byte{cbeCodecVersion, 0x79, 0x00, 0x00, 0x7b}))
	assertFailure(t, tryDecode([]byte{cbeCodecVersion, 0x79, 0x00, 0x7b}))
}

func TestCBEDecodeInlineList(t *testing.T) {
	encoded := []byte{0x00}
	expected := Tokens(list(), 0, end())
	actual, err := decodeDocumentCommon(InlineContainerTypeList, true, rules.DefaultLimits(), encoded)
	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}
	if !equivalence.IsEquivalent(actual, expected) {
		t.Errorf("Expected %v, actual %v", expected, actual)
	}
}

func TestCBEDecodeInlineMap(t *testing.T) {
	encoded := []byte{0x00, 0x01}
	expected := Tokens(Map(), 0, 1, end())
	actual, err := decodeDocumentCommon(InlineContainerTypeMap, true, rules.DefaultLimits(), encoded)
	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}
	if !equivalence.IsEquivalent(actual, expected) {
		t.Errorf("Expected %v, actual %v", expected, actual)
	}
}

func TestCBEDecodePiecemeal(t *testing.T) {
	encoded := []byte{cbeCodecVersion, 0x7a, 0x01, 0x7b}
	assertDecodedPiecemeal(t, encoded, 1, 3, Tokens(list(), 1, end()))
}

func TestCBEDecodePiecemeal2(t *testing.T) {
	encoded := []byte{cbeCodecVersion, 0x7a, 0x6a, 0x00, 0x01, 0x7b}
	assertDecodedPiecemeal(t, encoded, 1, 5, Tokens(list(), 0x100, end()))
}

func TestCBEDecodePiecemeal3(t *testing.T) {
	encoded := []byte{cbeCodecVersion, 0x7a, 0x01, 0x6a, 0x34, 0x12, 0x6c, 0xbc, 0x9a, 0x78, 0x56, 0x6e, 0x10, 0x32, 0x54, 0x76, 0x98, 0xba, 0xdc, 0xfe, 0x7b}
	assertDecodedPiecemeal(t, encoded, 1, 20, Tokens(list(), 1, 0x1234, 0x56789abc, uint64(0xfedcba9876543210), end()))
}
