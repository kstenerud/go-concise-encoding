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

const EndOfDocumentMarker = 0x100

func (_this CharProperty) HasProperty(property CharProperty) bool {
	return _this&property != 0
}

func StringRunesHaveProperty(str string, property CharProperty) bool {
	for _, ch := range str {
		if charProperties[ch].HasProperty(property) {
			return true
		}
	}
	return false
}

func StringBytesHaveProperty(str []byte, property CharProperty) bool {
	for _, ch := range str {
		if asciiProperties[ch].HasProperty(property) {
			return true
		}
	}
	return false
}

func GetRuneProperty(r rune) CharProperty {
	return charProperties[r]
}

func RuneHasProperty(r rune, property CharProperty) bool {
	return charProperties[r].HasProperty(property)
}

func ByteHasProperty(b byte, property CharProperty) bool {
	return asciiProperties[b].HasProperty(property)
}

type Byte uint8

func (_this Byte) HasProperty(property CharProperty) bool {
	return asciiProperties[_this].HasProperty(property)
}

type ByteWithEOF uint16

func (_this ByteWithEOF) HasProperty(property CharProperty) bool {
	return asciiProperties[_this].HasProperty(property)
}
