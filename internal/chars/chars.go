// Copyright 2019 Karl Stenerud
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE SOFTWARE.

package chars

import (
	"unicode/utf8"
)

const EOFMarker = 0x100

func CalculateRuneByteCount(startByte byte) int {
	return int(runeByteCounts[startByte>>3])
}

// Returns the index of the start of the last UTF-8 rune in data, and whether
// it's complete or not.
//
// If data is empty, returns (0, true)
// If there are no rune starts, returns (0, false)
// If there is a rune start, returns (index-of-rune-start, is-a-complete-rune)
func IndexOfLastRuneStart(data []byte) (index int, isCompleteRune bool) {
	dataLength := len(data)
	if dataLength == 0 {
		return 0, true
	}

	for index = dataLength - 1; index >= 0; index-- {
		runeByteCount := CalculateRuneByteCount(data[index])
		if runeByteCount > 0 {
			isCompleteRune = index+runeByteCount == dataLength
			return
		}
	}
	return
}

func (_this Properties) HasProperty(property Properties) bool {
	return _this&property != 0
}

func IsRuneSafeFor(r rune, flags SafetyFlags) bool {
	if getBitArrayValue(stringlikeSafe[:], int(r)) {
		return true
	}
	unsafety := stringlikeUnsafe[r]
	if unsafety == 0 {
		// unsafety 0 actually means "all". when the corresponding stringlikeSafe
		// entry indicates "unsafe". This keeps the stringlikeUnsafe map small.
		unsafety = SafetyAll
	}
	return unsafety&flags == 0
}

func IsIdentifierSafe(str []byte) bool {
	if len(str) == 0 {
		return false
	}
	r, size := utf8.DecodeRune(str)
	if !IsRuneValidIdentifierFirstChar(r) {
		return false
	}
	str = str[size:]
	for len(str) > 0 {
		r, size := utf8.DecodeRune(str)
		if !IsRuneValidIdentifier(r) {
			return false
		}
		str = str[size:]
	}
	return true
}

func IsMarkerIDSafe(str []byte) bool {
	if len(str) == 0 {
		return false
	}
	r, size := utf8.DecodeRune(str)
	if !IsRuneValidIdentifierFirstChar(r) {
		return false
	}
	str = str[size:]
	for len(str) > 0 {
		r, size := utf8.DecodeRune(str)
		if !IsRuneValidMarkerID(r) {
			return false
		}
		str = str[size:]
	}
	return true
}

func ByteHasProperty(b byte, property Properties) bool {
	return properties[b].HasProperty(property)
}

type Byte uint8

func (_this Byte) HasProperty(property Properties) bool {
	return properties[_this].HasProperty(property)
}

type ByteWithEOF uint16

func (_this ByteWithEOF) HasProperty(property Properties) bool {
	return properties[_this].HasProperty(property)
}

func IsRuneValidIdentifier(r rune) bool {
	return getBitArrayValue(identifierSafe[:], int(r))
}

func IsRuneValidIdentifierFirstChar(r rune) bool {
	return getBitArrayValue(identifierFirstSafe[:], int(r))
}

func IsRuneValidMarkerID(r rune) bool {
	return r != ':' && getBitArrayValue(identifierSafe[:], int(r))
}

func getBitArrayValue(array []byte, index int) bool {
	bits := array[index>>3]
	return bits&(1<<(index&7)) != 0
}

var HexCharValues = [256]byte{}

const InvalidHexChar = byte(0x80)

func init() {
	for i := 0; i < len(HexCharValues); i++ {
		HexCharValues[i] = InvalidHexChar
	}
	for i := '0'; i <= '9'; i++ {
		HexCharValues[i] = byte(i - '0')
	}
	for i := 'A'; i <= 'F'; i++ {
		HexCharValues[i] = byte(i - 'A' + 10)
	}
	for i := 'a'; i <= 'f'; i++ {
		HexCharValues[i] = byte(i - 'a' + 10)
	}
}

var HexChars = "0123456789abcdef"
