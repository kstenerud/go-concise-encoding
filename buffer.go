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

type buffer struct {
	bytes             []byte
	lastAllocatedSize int
}

func (this *buffer) Bytes() []byte {
	return this.bytes
}

func (this *buffer) Allocate(byteCount int) []byte {
	length := len(this.bytes)
	if cap(this.bytes)-length < byteCount {
		this.grow(byteCount)
	} else {
		this.bytes = this.bytes[:length+byteCount]
	}
	this.lastAllocatedSize = byteCount
	return this.bytes[length:]
}

func (this *buffer) CorrectAllocation(usedByteCount int) {
	unused := this.lastAllocatedSize - usedByteCount
	this.bytes = this.bytes[:len(this.bytes)-unused]
}

// ============================================================================

func (this *buffer) grow(byteCount int) {
	length := len(this.bytes)
	growAmount := cap(this.bytes)
	if byteCount > growAmount {
		if byteCount > minBufferCap {
			growAmount = byteCount
		} else {
			growAmount = minBufferCap
		}
	}
	newCap := cap(this.bytes) + growAmount
	newBytes := make([]byte, length+byteCount, newCap)
	oldBytes := this.bytes
	copy(newBytes, oldBytes)
	this.bytes = newBytes
}
