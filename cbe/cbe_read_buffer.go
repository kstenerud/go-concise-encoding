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

package cbe

import (
	"fmt"
	"io"
	"math"
	"math/big"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
	"github.com/kstenerud/go-concise-encoding/buffer"
	"github.com/kstenerud/go-uleb128"
)

const (
	defaultReadBufferSize = 2048
	defaultMinFreeBytes   = 48
)

type CBEReadBuffer struct {
	buffer   buffer.StreamingReadBuffer
	position int
}

// Create a new CBE read buffer. The buffer will be empty until
// RefillIfNecessary is called.
//
// The "low water" amount where RefillIfNecessary() will actually refill is
// readBufferSize/50, with a minimum of 10.
//
// Note: if readBufferSize is less than 64, it will use the default (2048).
func NewCBEReadBuffer(reader io.Reader, readBufferSize int) *CBEReadBuffer {
	_this := &CBEReadBuffer{}
	_this.Init(reader, readBufferSize)
	return _this
}

// Init the read buffer. The buffer will be empty until RefillIfNecessary is called.
//
// The "low water" amount where RefillIfNecessary() will actually refill is
// readBufferSize/50, with a minimum of 10.
//
// Note: if readBufferSize is less than 64, it will use the default (2048).
func (_this *CBEReadBuffer) Init(reader io.Reader, readBufferSize int) {
	minFreeBytes := defaultMinFreeBytes
	if readBufferSize < 64 {
		readBufferSize = defaultReadBufferSize
	} else {
		minFreeBytes = readBufferSize / 50
		if minFreeBytes < 10 {
			minFreeBytes = 10
		}
	}
	_this.buffer.Init(reader, readBufferSize, minFreeBytes)
}

// Refill the buffer from the reader if we've hit the "low water" of unread
// bytes.
func (_this *CBEReadBuffer) RefillIfNecessary() {
	_this.position += _this.buffer.RefillIfNecessary(_this.position)
}

func (_this *CBEReadBuffer) HasUnreadData() bool {
	return _this.position < len(_this.buffer.Buffer)
}

func (_this *CBEReadBuffer) DecodeVersion() uint64 {
	asUint, asBig := _this.DecodeUint()
	if asBig != nil {
		_this.errorf("Version too big")
	}
	return asUint
}

func (_this *CBEReadBuffer) DecodeType() cbeTypeField {
	return cbeTypeField(_this.DecodeUint8())
}

func (_this *CBEReadBuffer) DecodeUint8() uint8 {
	value := uint8(_this.byteAtPositionOffset(0))
	_this.markBytesRead(1)
	return value
}

func (_this *CBEReadBuffer) DecodeUint16() uint16 {
	value := uint16(_this.byteAtPositionOffset(0)) |
		uint16(_this.byteAtPositionOffset(1))<<8
	_this.markBytesRead(2)
	return value
}

func (_this *CBEReadBuffer) DecodeUint32() uint32 {
	value := uint32(_this.byteAtPositionOffset(0)) |
		uint32(_this.byteAtPositionOffset(1))<<8 |
		uint32(_this.byteAtPositionOffset(2))<<16 |
		uint32(_this.byteAtPositionOffset(3))<<24
	_this.markBytesRead(4)
	return value
}

func (_this *CBEReadBuffer) DecodeUint64() uint64 {
	value := uint64(_this.byteAtPositionOffset(0)) |
		uint64(_this.byteAtPositionOffset(1))<<8 |
		uint64(_this.byteAtPositionOffset(2))<<16 |
		uint64(_this.byteAtPositionOffset(3))<<24 |
		uint64(_this.byteAtPositionOffset(4))<<32 |
		uint64(_this.byteAtPositionOffset(5))<<40 |
		uint64(_this.byteAtPositionOffset(6))<<48 |
		uint64(_this.byteAtPositionOffset(7))<<56
	_this.markBytesRead(8)
	return value
}

func (_this *CBEReadBuffer) DecodeUint() (asUint uint64, asBig *big.Int) {
	asUint, asBig, bytesDecoded, isComplete := uleb128.Decode(0, 0, _this.allUnreadBytes())
	if !isComplete {
		_this.unexpectedEOD()
	}
	_this.markBytesRead(bytesDecoded)
	return
}

func (_this *CBEReadBuffer) DecodeFloat32() float32 {
	return math.Float32frombits(_this.DecodeUint32())
}

func (_this *CBEReadBuffer) DecodeFloat64() float64 {
	return math.Float64frombits(_this.DecodeUint64())
}

func (_this *CBEReadBuffer) DecodeDecimalFloat() (compact_float.DFloat, *apd.Decimal) {
	value, bigValue, bytesDecoded, err := compact_float.Decode(_this.allUnreadBytes())
	if err != nil {
		if err == compact_float.ErrorIncomplete {
			_this.buffer.RequestAndRetry(_this.position, bytesDecoded*2, func(positionOffset int) {
				_this.position += positionOffset
				value, bigValue, bytesDecoded, err = compact_float.Decode(_this.allUnreadBytes())
			})
			if err == compact_float.ErrorIncomplete {
				_this.unexpectedEOD()
			}
		}
		_this.unexpectedError(err)
	}
	_this.markBytesRead(bytesDecoded)
	return value, bigValue
}

func (_this *CBEReadBuffer) DecodeDate() *compact_time.Time {
	value, bytesDecoded, err := compact_time.DecodeDate(_this.allUnreadBytes())
	if err != nil {
		if err == compact_time.ErrorIncomplete {
			_this.buffer.RequestAndRetry(_this.position, bytesDecoded*2, func(positionOffset int) {
				_this.position += positionOffset
				value, bytesDecoded, err = compact_time.DecodeDate(_this.allUnreadBytes())
			})
			if err == compact_time.ErrorIncomplete {
				_this.unexpectedEOD()
			}
		}
		_this.unexpectedError(err)
	}
	if err := value.Validate(); err != nil {
		_this.unexpectedError(err)
	}
	_this.markBytesRead(bytesDecoded)
	return value
}

func (_this *CBEReadBuffer) DecodeTime() *compact_time.Time {
	value, bytesDecoded, isComplete := compact_time.DecodeTime(_this.allUnreadBytes())
	if err := value.Validate(); err != nil {
		_this.unexpectedError(err)
	}
	if !isComplete {
		_this.buffer.RequestAndRetry(_this.position, bytesDecoded*2, func(positionOffset int) {
			_this.position += positionOffset
			value, bytesDecoded, isComplete = compact_time.DecodeTime(_this.allUnreadBytes())
			if err := value.Validate(); err != nil {
				_this.unexpectedError(err)
			}
		})
		if !isComplete {
			_this.unexpectedEOD()
		}
	}
	_this.markBytesRead(bytesDecoded)
	return value
}

func (_this *CBEReadBuffer) DecodeTimestamp() *compact_time.Time {
	value, bytesDecoded, err := compact_time.DecodeTimestamp(_this.allUnreadBytes())
	if err != nil {
		if err == compact_time.ErrorIncomplete {
			_this.buffer.RequestAndRetry(_this.position, bytesDecoded*2, func(positionOffset int) {
				_this.position += positionOffset
				value, bytesDecoded, err = compact_time.DecodeTimestamp(_this.allUnreadBytes())
			})
			if err == compact_time.ErrorIncomplete {
				_this.unexpectedEOD()
			}
		}
		_this.unexpectedError(err)
	}
	if err := value.Validate(); err != nil {
		_this.unexpectedError(err)
	}
	_this.markBytesRead(bytesDecoded)
	return value
}

func (_this *CBEReadBuffer) DecodeChunkHeader() (length uint64, isFinalChunk bool) {
	asUint, asBig := _this.DecodeUint()
	if asBig != nil {
		_this.errorf("Chunk header length field is too big")
	}
	return asUint >> 1, asUint&1 == 1
}

func (_this *CBEReadBuffer) DecodeBytes(byteCount int) []byte {
	_this.buffer.RequireBytes(_this.position, byteCount)
	value := _this.buffer.Buffer[_this.position : _this.position+byteCount]
	_this.markBytesRead(byteCount)
	return value
}

func (_this *CBEReadBuffer) byteAtPositionOffset(offset int) byte {
	return _this.buffer.Buffer[_this.position+offset]
}

func (_this *CBEReadBuffer) markBytesRead(byteCount int) {
	_this.position += byteCount
}

func (_this *CBEReadBuffer) allUnreadBytes() []byte {
	return _this.buffer.Buffer[_this.position:]
}

func (_this *CBEReadBuffer) unexpectedEOD() {
	_this.errorf("Unexpected end of document")
}

func (_this *CBEReadBuffer) unexpectedError(err error) {
	// TODO: Diagnostics
	panic(err)
}

func (_this *CBEReadBuffer) errorf(format string, args ...interface{}) {
	// TODO: Diagnostics
	panic(fmt.Errorf(format, args...))
}
