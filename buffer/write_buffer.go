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

package buffer

import (
	"io"
)

const (
	minWriteBufferSize = 1024
)

type WriteBuffer struct {
	writer            io.Writer
	bytes             []byte
	lastAllocatedSize int
}

func NewWriteBuffer(writer io.Writer, bufferSize int) *WriteBuffer {
	_this := &WriteBuffer{}
	_this.Init(writer, bufferSize)
	return _this
}

func (_this *WriteBuffer) Init(writer io.Writer, bufferSize int) {
	if bufferSize < minWriteBufferSize {
		bufferSize = minWriteBufferSize
	}
	_this.bytes = make([]byte, 0, bufferSize)
	_this.writer = writer
}

func (_this *WriteBuffer) Bytes() []byte {
	return _this.bytes
}

func (_this *WriteBuffer) Allocate(byteCount int) []byte {
	length := len(_this.bytes)
	if cap(_this.bytes)-length < byteCount {
		_this.Flush()
		length = 0
		if cap(_this.bytes) < byteCount {
			_this.grow(byteCount)
		}
	}
	_this.bytes = _this.bytes[:length+byteCount]
	_this.lastAllocatedSize = byteCount
	return _this.bytes[length:]
}

func (_this *WriteBuffer) CorrectAllocation(usedByteCount int) {
	unused := _this.lastAllocatedSize - usedByteCount
	_this.bytes = _this.bytes[:len(_this.bytes)-unused]
}

func (_this *WriteBuffer) Flush() {
	_, err := _this.writer.Write(_this.bytes)
	if err != nil {
		panic(err)
	}
	_this.bytes = _this.bytes[:0]
}

// ============================================================================

func (_this *WriteBuffer) grow(byteCount int) {
	length := len(_this.bytes)
	growAmount := cap(_this.bytes)
	if byteCount > growAmount {
		growAmount = byteCount
	}
	newCap := cap(_this.bytes) + growAmount
	newBytes := make([]byte, length+byteCount, newCap)
	oldBytes := _this.bytes
	copy(newBytes, oldBytes)
	_this.bytes = newBytes
}
