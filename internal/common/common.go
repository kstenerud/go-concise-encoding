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

package common

import (
	"fmt"
	"math"
	"math/big"
	"reflect"
	"unicode"
	"unicode/utf8"
)

// Numeric

func Float64GetExponent(v float64) int {
	return int((math.Float64bits(v)>>52)&0x7ff) - 1023
}

// Determines how many numeric digits can be stored per X bits

var bitsToHexDigitsTable = []int{0, 1, 1, 1, 1, 2, 2, 2, 2, 3, 3, 3, 3, 4, 4, 4, 4}
var bitsToDecimalDigitsTable = []int{0, 1, 1, 1, 1, 2, 2, 2, 3, 3}
var decimalDigitsToBitsTable = []int{0, 4, 7}

func BitsToDecimalDigits(bitCount int) int {
	return (bitCount/10)*3 + bitsToDecimalDigitsTable[bitCount%10]
}

func DecimalDigitsToBits(digitCount int) int {
	triadCount := digitCount / 3
	remainder := digitCount % 3
	return triadCount*10 + decimalDigitsToBitsTable[remainder]
}

func BitsToHexDigits(bitCount int) int {
	return (bitCount/16)*4 + bitsToHexDigitsTable[bitCount&15]
}

func HexDigitsToBits(digitCount int) int {
	return digitCount * 4
}

// Architecture

func Is64BitArch() bool {
	return oneIf64Bit == 1
}

// Reflect

func IsFieldExported(name string) bool {
	ch, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(ch)
}

func IsSignalingNan(value float64) bool {
	return math.Float64bits(value)&QuietNanBit == 0
}

func IsBigIntNegative(value *big.Int) bool {
	return value.Cmp(BigInt0) < 0
}

func IsPointer(v reflect.Value) bool {
	return kindProperties[v.Kind()]&KindPropertyPointer != 0
}

func IsLengthable(v reflect.Value) bool {
	return kindProperties[v.Kind()]&KindPropertyLengthable != 0
}

func IsNullable(v reflect.Value) bool {
	return kindProperties[v.Kind()]&KindPropertyNullable != 0
}

func IsNil(v reflect.Value) bool {
	if !v.IsValid() {
		return true
	}

	return IsNullable(v) && v.IsNil()
}

func NameOf(x interface{}) string {
	return fmt.Sprintf("%v", reflect.TypeOf(x))
}

// Utility

func CloneBytes(bytes []byte) []byte {
	bytesCopy := make([]byte, len(bytes))
	copy(bytesCopy, bytes)
	return bytesCopy
}

var requiresLowercaseAdjust [256]bool

func init() {
	for i := 'A'; i <= 'Z'; i++ {
		requiresLowercaseAdjust[i] = true
	}
}

// Convert ASCII characters A-Z to a-z, ignoring locale.
func ASCIIBytesToLower(bytes []byte) (didChange bool) {
	const lowercaseAdjust = byte('a' - 'A')

	for i, b := range bytes {
		if requiresLowercaseAdjust[b] {
			bytes[i] += lowercaseAdjust
			didChange = true
		}
	}
	return
}

func ASCIIToLower(s string) string {
	asBytes := []byte(s)
	if ASCIIBytesToLower(asBytes) {
		return string(asBytes)
	}
	return s
}

func ByteCountToElementCount(elementBitWidth int, byteCount uint64) uint64 {
	return (byteCount * 8) / uint64(elementBitWidth)
}

func ElementCountToByteCount(elementBitWidth int, elementCount uint64) uint64 {
	byteCount := (elementCount * uint64(elementBitWidth)) / 8
	if elementBitWidth == 1 && elementCount&7 != 0 {
		byteCount++
	}
	return byteCount
}
