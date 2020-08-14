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

// Streaming & buffering mechanisms to support document creation/consumption.
package buffer

import (
	"fmt"
	"io"
)

var UnexpectedEOD = fmt.Errorf("unexpected end of document")

// A read buffer that can dynamically resize and stream data in from a reader.
type StreamingReadBuffer struct {
	Buffer       []byte
	reader       io.Reader
	minFreeBytes int
	isEOF        bool
}

// Create a new buffer. The buffer will be empty until a refill, request, or
// require method is called.
func NewStreamingReadBuffer(reader io.Reader, bufferSize int, minFreeBytes int) *StreamingReadBuffer {
	_this := &StreamingReadBuffer{}
	_this.Init(reader, bufferSize, minFreeBytes)
	return _this
}

// Initialize the buffer. Note that the buffer will be empty until a refill,
// request, or require method is called.
func (_this *StreamingReadBuffer) Init(reader io.Reader, bufferSize int, minFreeBytes int) {
	if len(_this.Buffer) < bufferSize {
		_this.Buffer = make([]byte, 0, bufferSize)
	}
	_this.reader = reader
	_this.minFreeBytes = minFreeBytes
}

// Remove all references that won't be needed for re-use so that the GC can reclaim them.
// Note: This won't free the internal document buffer.
func (_this *StreamingReadBuffer) Reset() {
	_this.reader = nil
	_this.Buffer = _this.Buffer[:0]
}

func (_this *StreamingReadBuffer) ByteAtOffset(offset int) byte {
	return _this.Buffer[offset]
}

func (_this *StreamingReadBuffer) HasByteAtOffset(offset int) bool {
	return offset < len(_this.Buffer)
}

// Refill the internal buffer from the reader only if there are less than
// the minimum count of bytes free and the reader hasn't reached EOF.
//
// The existing "unread" data (data after position) may be moved to the start
// of the buffer, in which case positionOffset will be set to the offset from
// old position to new.
func (_this *StreamingReadBuffer) RefillIfNecessary(startOffset, position int) (positionOffset int) {
	if !_this.isEOF && _this.unreadByteCount(position) < _this.minFreeBytes {
		return _this.Refill(startOffset)
	}
	return 0
}

// Refill the internal buffer from the reader if it hasn't reached EOF.
//
// The existing "unread" data (data after position) may be moved to the start
// of the buffer, in which case positionOffset will be set to the offset from
// old position to new.
func (_this *StreamingReadBuffer) Refill(position int) (positionOffset int) {
	if _this.isEOF {
		return
	}

	_this.moveUnreadBytesToStart(position)
	_this.readFromReader(len(_this.Buffer))
	positionOffset = -position
	return
}

// Request that the specified number of unread bytes be made available (position
// marks the threshold of the "unread" data).
//
// If necessary, The buffer may read more bytes from the reader, and may resize.
//
// The number of bytes read may not be as many as requested.
//
// The existing "unread" data (data after position) may be moved to the start
// of the buffer, in which case positionOffset will be set to the offset from
// old position to new.
func (_this *StreamingReadBuffer) RequestBytes(position int, byteCount int) (positionOffset int) {
	if _this.isEOF {
		return
	}

	unreadCount := _this.unreadByteCount(position)
	if byteCount <= unreadCount {
		return
	}

	if byteCount > cap(_this.Buffer) {
		// Make sure to at least double the buffer size to avoid massive
		// copying due to a loop that increments by a tiny amount every time.
		newSize := byteCount + cap(_this.Buffer)
		newSlice := make([]byte, newSize, newSize)
		copy(newSlice, _this.Buffer[position:])
		_this.Buffer = newSlice[:unreadCount]
	}

	return _this.Refill(position)
}

// Require that the specified number of unread bytes be made available (position
// marks the threshold of the "unread" data).
//
// If necessary, The buffer may read more bytes from the reader, and may resize.
//
// The existing "unread" data (data after position) may be moved to the start
// of the buffer, in which case positionOffset will be set to the offset from
// old position to new.
//
// If the required number of bytes cannot be read, this method panics with
// UnexpectedEOD.
func (_this *StreamingReadBuffer) RequireBytes(position int, byteCount int) (positionOffset int) {
	positionOffset = _this.RequestBytes(position, byteCount)
	if byteCount+position+positionOffset > len(_this.Buffer) {
		panic(UnexpectedEOD)
	}
	return
}

// Request more bytes and retry an operation.
func (_this *StreamingReadBuffer) RequestAndRetry(position int, requestedByteCount int, operation func(positionOffset int)) {
	_this.RequestBytes(position, requestedByteCount)
	operation(-position)
}

// Require more bytes and retry an operation.
func (_this *StreamingReadBuffer) RequireAndRetry(position int, requiredByteCount int, operation func(positionOffset int)) {
	_this.RequireBytes(position, requiredByteCount)
	operation(-position)
}

func (_this *StreamingReadBuffer) IsEOF() bool {
	return _this.isEOF
}

// ============================================================================

// Internal

func (_this *StreamingReadBuffer) unreadByteCount(position int) int {
	return len(_this.Buffer) - position
}

func (_this *StreamingReadBuffer) moveUnreadBytesToStart(position int) {
	if position > 0 {
		unreadCount := _this.unreadByteCount(position)
		copy(_this.Buffer, _this.Buffer[position:])
		_this.Buffer = _this.Buffer[:unreadCount]
	}
	return
}

func (_this *StreamingReadBuffer) increaseSizeTo(position int, byteCount int) {
	oldData := _this.Buffer[position:len(_this.Buffer)]
	newData := make([]byte, len(oldData), byteCount)
	copy(newData, oldData)
}

func (_this *StreamingReadBuffer) readFromReader(startPosition int) {
	_this.Buffer = _this.Buffer[:cap(_this.Buffer)]
	bytesRead, err := _this.reader.Read(_this.Buffer[startPosition:])
	if err != nil {
		if err != io.EOF {
			panic(err)
		}
		_this.isEOF = true
	}
	_this.Buffer = _this.Buffer[:startPosition+bytesRead]
}
