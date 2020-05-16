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

package concise_encoding

import "fmt"

type Utf8Validator struct {
	bytesRemaining int
	accumulator    int
}

func (this *Utf8Validator) Reset() {
	this.bytesRemaining = 0
	this.accumulator = 0
}

func (this *Utf8Validator) AddByte(byteValue int) {
	const continuationMask = 0xc0
	const continuationMatch = 0x80
	if this.bytesRemaining > 0 {
		if byteValue&continuationMask != continuationMatch {
			panic(fmt.Errorf("UTF-8 encoding: expected continuation bit (0x80) in byte [0x%02x]", byteValue))
		}
		this.bytesRemaining--
		this.accumulator = (this.accumulator << 6) | (byteValue & ^continuationMask)
		return
	}

	const initiator1ByteMask = 0x80
	const initiator1ByteMatch = 0x80
	if byteValue&initiator1ByteMask != initiator1ByteMatch {
		this.bytesRemaining = 0
		this.accumulator = byteValue
		if byteValue == 0 {
			panic(fmt.Errorf("UTF-8 encoding: NUL byte is not allowed"))
		}
		return
	}

	const initiator2ByteMask = 0xe0
	const initiator2ByteMatch = 0xc0
	const firstByte2ByteMask = 0x1f
	if (byteValue & initiator2ByteMask) == initiator2ByteMatch {
		this.bytesRemaining = 1
		this.accumulator = byteValue & firstByte2ByteMask
		return
	}

	const initiator3ByteMask = 0xf0
	const initiator3ByteMatch = 0xe0
	const firstByte3ByteMask = 0x0f
	if (byteValue & initiator3ByteMask) == initiator3ByteMatch {
		this.bytesRemaining = 2
		this.accumulator = byteValue & firstByte3ByteMask
		return
	}

	const initiator4ByteMask = 0xf8
	const initiator4ByteMatch = 0xf0
	const firstByte4ByteMask = 0x07
	if (byteValue & initiator4ByteMask) == initiator4ByteMatch {
		this.bytesRemaining = 3
		this.accumulator = byteValue & firstByte4ByteMask
		return
	}

	panic(fmt.Errorf("UTF-8 encoding: Invalid byte [0x%02x]", byteValue))
}

func (this *Utf8Validator) IsCompleteCharacter() bool {
	return this.bytesRemaining == 0
}

func (this *Utf8Validator) Character() int {
	return this.accumulator
}
