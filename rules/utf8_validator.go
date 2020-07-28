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

package rules

import "fmt"

// UTF8 Validator takes a stream of bytes and ensures that they form a complete
// UTF-8 character.
type UTF8Validator struct {
	bytesRemaining int
	accumulator    rune
}

func (_this *UTF8Validator) Reset() {
	_this.bytesRemaining = 0
	_this.accumulator = 0
}

// Add a byte to the UTF-8 character that is being built. When the character is
// complete. IsCompleteCharacter() will return true.
//
// This method panics if the UTF-8 sequence is invalid.
func (_this *UTF8Validator) AddByte(byteValue byte) {
	const continuationMask = 0xc0
	const continuationMatch = 0x80
	if _this.bytesRemaining > 0 {
		if byteValue&continuationMask != continuationMatch {
			panic(fmt.Errorf("UTF-8 encoding: expected continuation bit (0x80) in byte [0x%02x]", byteValue))
		}
		_this.bytesRemaining--
		_this.accumulator = (_this.accumulator << 6) | (rune(byteValue) & ^continuationMask)
		return
	}

	const initiator1ByteMask = 0x80
	const initiator1ByteMatch = 0x80
	if byteValue&initiator1ByteMask != initiator1ByteMatch {
		_this.bytesRemaining = 0
		_this.accumulator = rune(byteValue)
		if byteValue == 0 {
			panic(fmt.Errorf("UTF-8 encoding: NUL byte is not allowed"))
		}
		return
	}

	const initiator2ByteMask = 0xe0
	const initiator2ByteMatch = 0xc0
	const firstByte2ByteMask = 0x1f
	if (byteValue & initiator2ByteMask) == initiator2ByteMatch {
		_this.bytesRemaining = 1
		_this.accumulator = rune(byteValue) & firstByte2ByteMask
		return
	}

	const initiator3ByteMask = 0xf0
	const initiator3ByteMatch = 0xe0
	const firstByte3ByteMask = 0x0f
	if (byteValue & initiator3ByteMask) == initiator3ByteMatch {
		_this.bytesRemaining = 2
		_this.accumulator = rune(byteValue) & firstByte3ByteMask
		return
	}

	const initiator4ByteMask = 0xf8
	const initiator4ByteMatch = 0xf0
	const firstByte4ByteMask = 0x07
	if (byteValue & initiator4ByteMask) == initiator4ByteMatch {
		_this.bytesRemaining = 3
		_this.accumulator = rune(byteValue) & firstByte4ByteMask
		return
	}

	panic(fmt.Errorf("UTF-8 encoding: Invalid byte [0x%02x]", byteValue))
}

// Returns true if the last added byte completed the UTF-8 character.
func (_this *UTF8Validator) IsCompleteCharacter() bool {
	return _this.bytesRemaining == 0
}

// Get the fully built UTF-8 character. Don't call this until
// IsCompleteCharacter() returns true.
func (_this *UTF8Validator) GetCharacter() rune {
	return _this.accumulator
}
