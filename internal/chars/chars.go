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
	"fmt"
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
	if str[0] == '-' {
		panic(fmt.Errorf("Identifier may not start with '-'"))
	}

	for i := 0; i < len(str); {
		byteCount := CalculateRuneByteCount(str[i])
		switch byteCount {
		case 0:
			panic(fmt.Errorf("Identifier contains invalid UTF-8 sequence"))
		case 1:
			if !IsRuneValidIdentifier(rune(str[i])) {
				return false
			}
			i++
		case 2:
			if i+1 >= len(str) {
				panic(fmt.Errorf("Identifier contains incomplete UTF-8 sequence"))
			}
			r := (rune(str[i]&0x1f) << 6) | rune(str[i+1]&0x3f)
			if !IsRuneValidIdentifier(r) {
				return false
			}
			i += 2
		case 3:
			if i+2 >= len(str) {
				panic(fmt.Errorf("Identifier contains incomplete UTF-8 sequence"))
			}
			r := (rune(str[i]&0x0f) << 12) | (rune(str[i+1]&0x3f) << 6) | rune(str[i+1]&0x3f)
			if !IsRuneValidIdentifier(r) {
				return false
			}
			i += 3
		case 4:
			if i+3 >= len(str) {
				panic(fmt.Errorf("Identifier contains incomplete UTF-8 sequence"))
			}
			r := (rune(str[i]&0x07) << 18) | (rune(str[i+1]&0x3f) << 12) | (rune(str[i+2]&0x3f) << 6) | rune(str[i+3]&0x3f)
			if !IsRuneValidIdentifier(r) {
				return false
			}
			i += 4
		}
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

func getBitArrayValue(array []byte, index int) bool {
	bits := array[index>>3]
	return bits&(1<<(index&7)) != 0
}
