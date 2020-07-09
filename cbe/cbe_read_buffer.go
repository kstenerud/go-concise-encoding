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
	"math"
	"math/big"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
	"github.com/kstenerud/go-concise-encoding/buffer"
	"github.com/kstenerud/go-uleb128"
)

type CBEReadBuffer struct {
	buffer.ReadBuffer
}

func (_this *CBEReadBuffer) DecodeVersion() uint64 {
	asUint, asBig := _this.DecodeUint()
	if asBig != nil {
		_this.ReadBuffer.Errorf("Version too big")
	}
	return asUint
}

func (_this *CBEReadBuffer) DecodeType() cbeTypeField {
	return cbeTypeField(_this.DecodeUint8())
}

func (_this *CBEReadBuffer) DecodeUint8() uint8 {
	value := uint8(_this.ReadBuffer.ByteAtPositionOffset(0))
	_this.ReadBuffer.MarkBytesRead(1)
	return value
}

func (_this *CBEReadBuffer) DecodeUint16() uint16 {
	value := uint16(_this.ReadBuffer.ByteAtPositionOffset(0)) |
		uint16(_this.ReadBuffer.ByteAtPositionOffset(1))<<8
	_this.ReadBuffer.MarkBytesRead(2)
	return value
}

func (_this *CBEReadBuffer) DecodeUint32() uint32 {
	value := uint32(_this.ReadBuffer.ByteAtPositionOffset(0)) |
		uint32(_this.ReadBuffer.ByteAtPositionOffset(1))<<8 |
		uint32(_this.ReadBuffer.ByteAtPositionOffset(2))<<16 |
		uint32(_this.ReadBuffer.ByteAtPositionOffset(3))<<24
	_this.ReadBuffer.MarkBytesRead(4)
	return value
}

func (_this *CBEReadBuffer) DecodeUint64() uint64 {
	value := uint64(_this.ReadBuffer.ByteAtPositionOffset(0)) |
		uint64(_this.ReadBuffer.ByteAtPositionOffset(1))<<8 |
		uint64(_this.ReadBuffer.ByteAtPositionOffset(2))<<16 |
		uint64(_this.ReadBuffer.ByteAtPositionOffset(3))<<24 |
		uint64(_this.ReadBuffer.ByteAtPositionOffset(4))<<32 |
		uint64(_this.ReadBuffer.ByteAtPositionOffset(5))<<40 |
		uint64(_this.ReadBuffer.ByteAtPositionOffset(6))<<48 |
		uint64(_this.ReadBuffer.ByteAtPositionOffset(7))<<56
	_this.ReadBuffer.MarkBytesRead(8)
	return value
}

func (_this *CBEReadBuffer) DecodeUint() (asUint uint64, asBig *big.Int) {
	asUint, asBig, bytesDecoded, isComplete := uleb128.Decode(0, 0, _this.ReadBuffer.AllUnreadBytes())
	if !isComplete {
		_this.ReadBuffer.UnexpectedEOD()
	}
	_this.ReadBuffer.MarkBytesRead(bytesDecoded)
	return
}

func (_this *CBEReadBuffer) DecodeFloat32() float32 {
	return math.Float32frombits(_this.DecodeUint32())
}

func (_this *CBEReadBuffer) DecodeFloat64() float64 {
	return math.Float64frombits(_this.DecodeUint64())
}

func (_this *CBEReadBuffer) DecodeDecimalFloat() (compact_float.DFloat, *apd.Decimal) {
	value, bigValue, bytesDecoded, err := compact_float.Decode(_this.ReadBuffer.AllUnreadBytes())
	if err != nil {
		if err == compact_float.ErrorIncomplete {
			_this.ReadBuffer.RefillAndRetry(func() {
				value, bigValue, bytesDecoded, err = compact_float.Decode(_this.ReadBuffer.AllUnreadBytes())
			})
			if err == compact_float.ErrorIncomplete {
				_this.ReadBuffer.UnexpectedEOD()
			}
		}
		_this.ReadBuffer.UnexpectedError(err)
	}
	_this.ReadBuffer.MarkBytesRead(bytesDecoded)
	return value, bigValue
}

func (_this *CBEReadBuffer) DecodeDate() *compact_time.Time {
	value, bytesDecoded, err := compact_time.DecodeDate(_this.ReadBuffer.AllUnreadBytes())
	if err != nil {
		if err == compact_time.ErrorIncomplete {
			_this.ReadBuffer.RefillAndRetry(func() {
				value, bytesDecoded, err = compact_time.DecodeDate(_this.ReadBuffer.AllUnreadBytes())
			})
			if err == compact_time.ErrorIncomplete {
				_this.ReadBuffer.UnexpectedEOD()
			}
		}
		_this.ReadBuffer.UnexpectedError(err)
	}
	if err := value.Validate(); err != nil {
		_this.ReadBuffer.UnexpectedError(err)
	}
	_this.ReadBuffer.MarkBytesRead(bytesDecoded)
	return value
}

func (_this *CBEReadBuffer) DecodeTime() *compact_time.Time {
	value, bytesDecoded, isComplete := compact_time.DecodeTime(_this.ReadBuffer.AllUnreadBytes())
	if err := value.Validate(); err != nil {
		_this.ReadBuffer.UnexpectedError(err)
	}
	if !isComplete {
		_this.ReadBuffer.RefillAndRetry(func() {
			value, bytesDecoded, isComplete = compact_time.DecodeTime(_this.ReadBuffer.AllUnreadBytes())
			if err := value.Validate(); err != nil {
				_this.ReadBuffer.UnexpectedError(err)
			}
		})
		if !isComplete {
			_this.ReadBuffer.UnexpectedEOD()
		}
	}
	_this.ReadBuffer.MarkBytesRead(bytesDecoded)
	return value
}

func (_this *CBEReadBuffer) DecodeTimestamp() *compact_time.Time {
	value, bytesDecoded, err := compact_time.DecodeTimestamp(_this.ReadBuffer.AllUnreadBytes())
	if err != nil {
		if err == compact_time.ErrorIncomplete {
			_this.ReadBuffer.RefillAndRetry(func() {
				value, bytesDecoded, err = compact_time.DecodeTimestamp(_this.ReadBuffer.AllUnreadBytes())
			})
			if err == compact_time.ErrorIncomplete {
				_this.ReadBuffer.UnexpectedEOD()
			}
		}
		_this.ReadBuffer.UnexpectedError(err)
	}
	if err := value.Validate(); err != nil {
		_this.ReadBuffer.UnexpectedError(err)
	}
	_this.ReadBuffer.MarkBytesRead(bytesDecoded)
	return value
}

func (_this *CBEReadBuffer) DecodeChunkHeader() (length uint64, isFinalChunk bool) {
	asUint, asBig := _this.DecodeUint()
	if asBig != nil {
		_this.ReadBuffer.Errorf("Chunk header length field is too big")
	}
	return asUint >> 1, asUint&1 == 1
}

func (_this *CBEReadBuffer) DecodeBytes(byteCount int) []byte {
	_this.ReadBuffer.RequireBytes(byteCount)
	value := _this.ReadBuffer.NextUnreadBytes(byteCount)
	_this.ReadBuffer.MarkBytesRead(byteCount)
	return value
}
