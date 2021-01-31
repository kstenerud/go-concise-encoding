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

	"github.com/kstenerud/go-concise-encoding/internal/common"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
	"github.com/kstenerud/go-concise-encoding/buffer"
	"github.com/kstenerud/go-uleb128"
)

type ReadBuffer struct {
	buffer   buffer.StreamingReadBuffer
	position int
}

// Create a new CBE read buffer. The buffer will be empty until RefillIfNecessary() is
// called.
//
// readBufferSize determines the initial size of the buffer, and
// loWaterByteCount determines when RefillIfNecessary() refills the buffer from
// the reader.
func NewReadBuffer(reader io.Reader, readBufferSize int, loWaterByteCount int) *ReadBuffer {
	_this := &ReadBuffer{}
	_this.Init(reader, readBufferSize, loWaterByteCount)
	return _this
}

// Init the read buffer. The buffer will be empty until RefillIfNecessary() is
// called.
//
// readBufferSize determines the initial size of the buffer, and
// loWaterByteCount determines when RefillIfNecessary() refills the buffer from
// the reader.
func (_this *ReadBuffer) Init(reader io.Reader, readBufferSize int, loWaterByteCount int) {
	_this.buffer.Init(reader, readBufferSize, loWaterByteCount)
	_this.position = 0
}

func (_this *ReadBuffer) Reset() {
	_this.buffer.Reset()
	_this.position = 0
}

// Refill the buffer from the reader if we've hit the "low water" of unread
// bytes.
func (_this *ReadBuffer) RefillIfNecessary() {
	_this.position += _this.buffer.RefillIfNecessary(_this.position, _this.position)
}

func (_this *ReadBuffer) HasUnreadData() bool {
	return _this.position < len(_this.buffer.Buffer)
}

func (_this *ReadBuffer) DecodeUint8() uint8 {
	_this.buffer.RequireBytes(_this.position, 1)
	value := _this.byteAtPositionOffset(0)
	_this.markBytesRead(1)
	return value
}

func (_this *ReadBuffer) DecodeUint16() uint16 {
	_this.buffer.RequireBytes(_this.position, 2)
	value := uint16(_this.byteAtPositionOffset(0)) |
		uint16(_this.byteAtPositionOffset(1))<<8
	_this.markBytesRead(2)
	return value
}

func (_this *ReadBuffer) DecodeUint32() uint32 {
	_this.buffer.RequireBytes(_this.position, 4)
	value := uint32(_this.byteAtPositionOffset(0)) |
		uint32(_this.byteAtPositionOffset(1))<<8 |
		uint32(_this.byteAtPositionOffset(2))<<16 |
		uint32(_this.byteAtPositionOffset(3))<<24
	_this.markBytesRead(4)
	return value
}

func (_this *ReadBuffer) DecodeUint64() uint64 {
	_this.buffer.RequireBytes(_this.position, 8)
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

func (_this *ReadBuffer) DecodeVersion() uint64 {
	return _this.DecodeSmallULEB128("version", 0xffffffffffffffff)
}

func (_this *ReadBuffer) DecodeType() cbeTypeField {
	return cbeTypeField(_this.DecodeUint8())
}

func (_this *ReadBuffer) DecodeULEB128() (asUint uint64, asBig *big.Int) {
	asUint, asBig, bytesDecoded, isComplete := uleb128.Decode(0, 0, _this.allUnreadBytes())
	if isComplete {
		goto complete
	}
	_this.buffer.RequestAndRetry(_this.position, bytesDecoded*2, func(positionOffset int) {
		_this.position += positionOffset
		asUint, asBig, bytesDecoded, isComplete = uleb128.Decode(0, 0, _this.allUnreadBytes())
	})
	if !isComplete {
		_this.unexpectedEOD()
	}

complete:
	_this.markBytesRead(bytesDecoded)
	return
}

func (_this *ReadBuffer) DecodeSmallULEB128(name string, maxValue uint64) uint64 {
	asUint, asBig, bytesDecoded, isComplete := uleb128.Decode(0, 0, _this.allUnreadBytes())
	if isComplete {
		goto complete
	}
	_this.buffer.RequestAndRetry(_this.position, bytesDecoded*2, func(positionOffset int) {
		_this.position += positionOffset
		asUint, asBig, bytesDecoded, isComplete = uleb128.Decode(0, 0, _this.allUnreadBytes())
	})
	if !isComplete {
		_this.unexpectedEOD()
	}

complete:
	if asBig != nil {
		_this.errorf("%v: %v is too big", asBig, name)
	}
	if asUint > maxValue {
		_this.errorf("%v: %v is too big (max allowed value = %v)", asBig, name, maxValue)
	}
	_this.markBytesRead(bytesDecoded)
	return asUint
}

// TODO: Check max big.int bit count
const maxBigIntBitCount = 8192

func (_this *ReadBuffer) DecodeUint() (asUint uint64, asBig *big.Int) {
	byteCount := _this.DecodeSmallULEB128("uint length field", maxBigIntBitCount/8)
	bytes := _this.DecodeBytes(int(byteCount))
	if byteCount <= 8 {
		for i := 0; i < len(bytes); i++ {
			asUint |= uint64(bytes[i]) << (i * 8)
		}
		return
	}
	bytesPerWord := common.BytesPerInt
	wordCount := len(bytes) / bytesPerWord
	if len(bytes)%bytesPerWord != 0 {
		wordCount++
	}
	words := make([]big.Word, wordCount, wordCount)
	iWord := 0
	for iByte := 0; iByte < len(bytes); {
		word := big.Word(0)
		for iPart := 0; iPart < bytesPerWord && iByte < len(bytes); iPart++ {
			word |= big.Word(bytes[iByte]) << (iPart * 8)
			iByte++
		}
		words[iWord] = word
		iWord++
		word = 0
	}
	asBig = big.NewInt(0)
	asBig.SetBits(words)
	return
}

func (_this *ReadBuffer) DecodeFloat16() float32 {
	return math.Float32frombits(uint32(_this.DecodeUint16()) << 16)
}

func (_this *ReadBuffer) DecodeFloat32() float32 {
	return math.Float32frombits(_this.DecodeUint32())
}

func (_this *ReadBuffer) DecodeFloat64() float64 {
	return math.Float64frombits(_this.DecodeUint64())
}

func (_this *ReadBuffer) DecodeDecimalFloat() (compact_float.DFloat, *apd.Decimal) {
	value, bigValue, bytesDecoded, err := compact_float.Decode(_this.allUnreadBytes())
	if err == nil {
		goto complete
	}
	if err != compact_float.ErrorIncomplete {
		_this.unexpectedError(err)
	}
	_this.buffer.RequestAndRetry(_this.position, bytesDecoded*2, func(positionOffset int) {
		_this.position += positionOffset
		value, bigValue, bytesDecoded, err = compact_float.Decode(_this.allUnreadBytes())
	})
	if err == compact_float.ErrorIncomplete {
		_this.unexpectedEOD()
	}

complete:
	_this.markBytesRead(bytesDecoded)
	return value, bigValue
}

func (_this *ReadBuffer) DecodeDate() compact_time.Time {
	value, bytesDecoded, err := compact_time.DecodeDate(_this.allUnreadBytes())
	if err == compact_time.ErrorIncomplete {
		_this.buffer.RequestAndRetry(_this.position, bytesDecoded*2, func(positionOffset int) {
			_this.position += positionOffset
			value, bytesDecoded, err = compact_time.DecodeDate(_this.allUnreadBytes())
		})
		if err == compact_time.ErrorIncomplete {
			_this.unexpectedEOD()
		}
	}

	if err != nil {
		_this.unexpectedError(err)
	}
	_this.markBytesRead(bytesDecoded)
	return value
}

func (_this *ReadBuffer) DecodeTime() compact_time.Time {
	value, bytesDecoded, err := compact_time.DecodeTime(_this.allUnreadBytes())
	if err == compact_time.ErrorIncomplete {
		_this.buffer.RequestAndRetry(_this.position, bytesDecoded*2, func(positionOffset int) {
			_this.position += positionOffset
			value, bytesDecoded, err = compact_time.DecodeTime(_this.allUnreadBytes())
		})
		if err == compact_time.ErrorIncomplete {
			_this.unexpectedEOD()
		}
	}

	if err != nil {
		_this.unexpectedError(err)
	}
	_this.markBytesRead(bytesDecoded)
	return value
}

func (_this *ReadBuffer) DecodeTimestamp() compact_time.Time {
	value, bytesDecoded, err := compact_time.DecodeTimestamp(_this.allUnreadBytes())
	if err == compact_time.ErrorIncomplete {
		_this.buffer.RequestAndRetry(_this.position, bytesDecoded*2, func(positionOffset int) {
			_this.position += positionOffset
			value, bytesDecoded, err = compact_time.DecodeTimestamp(_this.allUnreadBytes())
		})
		if err == compact_time.ErrorIncomplete {
			_this.unexpectedEOD()
		}
	}

	if err != nil {
		_this.unexpectedError(err)
	}
	_this.markBytesRead(bytesDecoded)
	return value
}

func (_this *ReadBuffer) DecodeArrayChunkHeader() (length uint64, moreChunksFollow bool) {
	asUint, asBig := _this.DecodeULEB128()
	if asBig != nil {
		_this.errorf("Chunk header length field is too big")
	}
	return asUint >> 1, asUint&1 == 1
}

func (_this *ReadBuffer) DecodeBytes(byteCount int) []byte {
	_this.buffer.RequireBytes(_this.position, byteCount)
	value := _this.buffer.Buffer[_this.position : _this.position+byteCount]
	_this.markBytesRead(byteCount)
	return value
}

// ============================================================================

// Internal

func (_this *ReadBuffer) byteAtPositionOffset(offset int) byte {
	return _this.buffer.Buffer[_this.position+offset]
}

func (_this *ReadBuffer) markBytesRead(byteCount int) {
	_this.position += byteCount
}

func (_this *ReadBuffer) allUnreadBytes() []byte {
	return _this.buffer.Buffer[_this.position:]
}

func (_this *ReadBuffer) unexpectedEOD() {
	_this.errorf("Unexpected end of document")
}

func (_this *ReadBuffer) unexpectedError(err error) {
	// TODO: Diagnostics
	panic(err)
}

func (_this *ReadBuffer) errorf(format string, args ...interface{}) {
	// TODO: Diagnostics
	panic(fmt.Errorf(format, args...))
}
