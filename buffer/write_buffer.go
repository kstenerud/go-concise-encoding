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
	"fmt"
	"io"
)

const minWriterFreeSpace = 300

type StreamingWriteBuffer struct {
	writer io.Writer
	bytes  []byte
	pos    []byte
	used   int
}

func NewWriteBuffer(bufferSize int) *StreamingWriteBuffer {
	_this := &StreamingWriteBuffer{}
	_this.Init(bufferSize)
	return _this
}

func (_this *StreamingWriteBuffer) Init(bufferSize int) {
	if bufferSize < minWriterFreeSpace*2 {
		bufferSize = minWriterFreeSpace * 2
	}
	_this.bytes = make([]byte, bufferSize, bufferSize)
	_this.pos = _this.bytes
}

func (_this *StreamingWriteBuffer) SetWriter(writer io.Writer) {
	_this.writer = writer
}

func (_this *StreamingWriteBuffer) Reset() {
	_this.writer = nil
	_this.pos = _this.bytes
	_this.used = 0
}

func (_this *StreamingWriteBuffer) flushUnconditionally() {
	data := _this.bytes[:_this.used]
	_this.pos = _this.bytes
	_this.used = 0

	_, err := _this.writer.Write(data)
	if err != nil {
		panic(err)
	}
}

func (_this *StreamingWriteBuffer) Flush() {
	if _this.used != 0 {
		_this.flushUnconditionally()
	}
}

func (_this *StreamingWriteBuffer) FlushIfNeeded() {
	if len(_this.pos) < minWriterFreeSpace {
		_this.flushUnconditionally()
	}
}

func (_this *StreamingWriteBuffer) growTo(byteCount int) {
	_this.Flush()
	if len(_this.bytes) < byteCount {
		_this.bytes = make([]byte, byteCount, byteCount)
		_this.pos = _this.bytes
	}
}

// Allocates byteCount + minWriterFreeSpace
func (_this *StreamingWriteBuffer) RequireBytes(byteCount int) []byte {
	byteCount += minWriterFreeSpace
	if byteCount > len(_this.pos) {
		_this.growTo(byteCount)
	}
	return _this.pos
}

func (_this *StreamingWriteBuffer) UseBytes(usedByteCount int) {
	_this.used += usedByteCount
	_this.pos = _this.pos[usedByteCount:]
}

// Warning: AddByte doesn't check for the end of the buffer
func (_this *StreamingWriteBuffer) AddByte(b byte) {
	_this.RequireBytes(1)
	_this.pos[0] = b
	_this.UseBytes(1)
}

func (_this *StreamingWriteBuffer) AddString(str string) {
	length := len(str)
	buf := _this.RequireBytes(length)
	copy(buf, str)
	_this.UseBytes(length)
}

func (_this *StreamingWriteBuffer) AddBytes(bytes []byte) {
	length := len(bytes)
	buf := _this.RequireBytes(length)
	copy(buf, bytes)
	_this.UseBytes(length)
}

func (_this *StreamingWriteBuffer) AddFmt(format string, args ...interface{}) {
	_this.AddString(fmt.Sprintf(format, args...))
}

// Add a formatted string, but strip the specified number of characters from the
// beginning of the result before adding.
func (_this *StreamingWriteBuffer) AddFmtStripped(stripByteCount int, format string, args ...interface{}) {
	_this.AddString(fmt.Sprintf(format, args...)[stripByteCount:])
}
