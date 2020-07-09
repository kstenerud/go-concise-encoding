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

const (
	minReadBufferSize = 2048
	minFreeBytes      = 48
)

// ReadBuffer is a base class for read buffers used by CBE and CTE decoders.
type ReadBuffer struct {
	reader   io.Reader
	data     []byte
	position int
	isEOF    bool
}

func NewReadBuffer(reader io.Reader, bufferSize int) *ReadBuffer {
	_this := &ReadBuffer{}
	_this.Init(reader, bufferSize)
	return _this
}

func (_this *ReadBuffer) Init(reader io.Reader, bufferSize int) {
	if bufferSize < minReadBufferSize {
		bufferSize = minReadBufferSize
	}
	_this.reader = reader
	_this.data = make([]byte, 0, bufferSize)
	_this.position = 0
}

// Refill the internal buffer from the reader only if there are less than
// the minimum bytes free.
func (_this *ReadBuffer) RefillIfNecessary() error {
	if _this.unreadByteCount() < minFreeBytes {
		return _this.Refill()
	}
	return nil
}

// Refill the internal buffer from the reader.
// If the reader has reached EOF, this method does nothing.
func (_this *ReadBuffer) Refill() error {
	if _this.isEOF {
		return nil
	}

	_this.moveUnreadBytesToStart()
	unreadCount := len(_this.data)
	newData := _this.data[:cap(_this.data)]
	bytesRead, err := _this.reader.Read(newData[unreadCount:])
	if err == io.EOF {
		_this.isEOF = true
		err = nil
	}
	_this.data = newData[:unreadCount+bytesRead]
	return err
}

func (_this *ReadBuffer) HasUnreadData() bool {
	return _this.position < len(_this.data)
}

// Return the byte at the offset from the current position.
func (_this *ReadBuffer) ByteAtPositionOffset(offset int) byte {
	return _this.data[_this.position+offset]
}

// Return a slice of all remaining unread bytes in the current buffer.
func (_this *ReadBuffer) AllUnreadBytes() []byte {
	return _this.data[_this.position:]
}

// Return the specified number of unread bytes.
func (_this *ReadBuffer) NextUnreadBytes(length int) []byte {
	return _this.data[_this.position : _this.position+length]
}

// Mark the specified number of bytes read, advancing our position in the buffer.
func (_this *ReadBuffer) MarkBytesRead(byteCount int) {
	_this.position += byteCount
}

// Notify that you require the specified number of bytes. If necessary, The
// buffer may read more bytes from the reader or possibly even increase the
// buffer size.
func (_this *ReadBuffer) RequireBytes(requiredByteCount int) {
	if requiredByteCount <= _this.unreadByteCount() {
		return
	}

	if requiredByteCount > cap(_this.data) {
		_this.increaseSizeTo(requiredByteCount)
	}
	if err := _this.Refill(); err != nil {
		_this.UnexpectedError(err)
	}
	if requiredByteCount > len(_this.data) {
		_this.UnexpectedEOD()
	}
}

// Refill the buffer and retry an operation.
func (_this *ReadBuffer) RefillAndRetry(operation func()) {
	if _this.position > 0 {
		if err := _this.Refill(); err != nil {
			_this.UnexpectedError(err)
		}
		operation()
	}
	_this.UnexpectedEOD()
}

func (_this *ReadBuffer) UnexpectedEOD() {
	_this.Errorf("Unexpected end of document")
}

func (_this *ReadBuffer) UnexpectedError(err error) {
	// TODO: Diagnostics
	panic(err)
}

func (_this *ReadBuffer) Errorf(format string, args ...interface{}) {
	// TODO: Diagnostics
	panic(fmt.Errorf(format, args...))
}

func (_this *ReadBuffer) increaseSizeTo(byteCount int) {
	oldData := _this.data[_this.position:len(_this.data)]
	newData := make([]byte, len(oldData), byteCount)
	copy(newData, oldData)
}

func (_this *ReadBuffer) unreadByteCount() int {
	return len(_this.data) - _this.position
}

func (_this *ReadBuffer) moveUnreadBytesToStart() {
	if _this.position > 0 {
		unreadCount := _this.unreadByteCount()
		copy(_this.data, _this.data[_this.position:])
		_this.data = _this.data[:unreadCount]
		_this.position = 0
	}
}
