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
	"math"
	"math/big"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
	"github.com/kstenerud/go-uleb128"
)

type cbeDecodeBuffer struct {
	data     []byte
	position int
}

func (_this *cbeDecodeBuffer) Init(data []byte) {
	_this.data = data
	_this.position = 0
}

func (_this *cbeDecodeBuffer) HasUnreadData() bool {
	return _this.position < len(_this.data)
}

func (_this *cbeDecodeBuffer) DecodeUint8() uint8 {
	value := uint8(_this.data[_this.position])
	_this.position++
	return value
}

func (_this *cbeDecodeBuffer) DecodeUint16() uint16 {
	value := uint16(_this.data[_this.position]) |
		uint16(_this.data[_this.position+1])<<8
	_this.position += 2
	return value
}

func (_this *cbeDecodeBuffer) DecodeUint32() uint32 {
	value := uint32(_this.data[_this.position]) |
		uint32(_this.data[_this.position+1])<<8 |
		uint32(_this.data[_this.position+2])<<16 |
		uint32(_this.data[_this.position+3])<<24
	_this.position += 4
	return value
}

func (_this *cbeDecodeBuffer) DecodeUint64() uint64 {
	value := uint64(_this.data[_this.position]) |
		uint64(_this.data[_this.position+1])<<8 |
		uint64(_this.data[_this.position+2])<<16 |
		uint64(_this.data[_this.position+3])<<24 |
		uint64(_this.data[_this.position+4])<<32 |
		uint64(_this.data[_this.position+5])<<40 |
		uint64(_this.data[_this.position+6])<<48 |
		uint64(_this.data[_this.position+7])<<56
	_this.position += 8
	return value
}

func (_this *cbeDecodeBuffer) DecodeUint() (asUint uint64, asBig *big.Int) {
	asUint, asBig, bytesDecoded, isComplete := uleb128.Decode(0, 0, _this.data[_this.position:])
	if !isComplete {
		_this.unexpectedEOD()
	}
	_this.position += bytesDecoded
	return
}

func (_this *cbeDecodeBuffer) DecodeVersion() uint64 {
	asUint, asBig := _this.DecodeUint()
	if asBig != nil {
		_this.errorf("Version too big")
	}
	return asUint
}

func (_this *cbeDecodeBuffer) DecodeType() cbeTypeField {
	return cbeTypeField(_this.DecodeUint8())
}

func (_this *cbeDecodeBuffer) DecodeFloat32() float32 {
	return math.Float32frombits(_this.DecodeUint32())
}

func (_this *cbeDecodeBuffer) DecodeFloat64() float64 {
	return math.Float64frombits(_this.DecodeUint64())
}

func (_this *cbeDecodeBuffer) DecodeDecimalFloat() (compact_float.DFloat, *apd.Decimal) {
	value, bigValue, bytesDecoded, err := compact_float.Decode(_this.data[_this.position:])
	if err != nil {
		if err == compact_float.ErrorIncomplete {
			_this.unexpectedEOD()
		}
		_this.unexpectedError(err)
	}
	_this.position += bytesDecoded
	return value, bigValue
}

func (_this *cbeDecodeBuffer) DecodeDate() *compact_time.Time {
	value, bytesDecoded, err := compact_time.DecodeDate(_this.data[_this.position:])
	if err != nil {
		if err == compact_time.ErrorIncomplete {
			_this.unexpectedEOD()
		}
		_this.unexpectedError(err)
	}
	if err := value.Validate(); err != nil {
		_this.unexpectedError(err)
	}
	_this.position += bytesDecoded
	return value
}

func (_this *cbeDecodeBuffer) DecodeTime() *compact_time.Time {
	value, bytesDecoded, isComplete := compact_time.DecodeTime(_this.data[_this.position:])
	if err := value.Validate(); err != nil {
		_this.unexpectedError(err)
	}
	if !isComplete {
		_this.unexpectedEOD()
	}
	_this.position += bytesDecoded
	return value
}

func (_this *cbeDecodeBuffer) DecodeTimestamp() *compact_time.Time {
	value, bytesDecoded, err := compact_time.DecodeTimestamp(_this.data[_this.position:])
	if err != nil {
		if err == compact_time.ErrorIncomplete {
			_this.unexpectedEOD()
		}
		_this.unexpectedError(err)
	}
	if err := value.Validate(); err != nil {
		_this.unexpectedError(err)
	}
	_this.position += bytesDecoded
	return value
}

func (_this *cbeDecodeBuffer) DecodeBytes(byteCount int) []byte {
	value := _this.data[_this.position : _this.position+byteCount]
	_this.position += byteCount
	return value
}

func (_this *cbeDecodeBuffer) DecodeChunkHeader() (length uint64, isFinalChunk bool) {
	asUint, asBig := _this.DecodeUint()
	if asBig != nil {
		_this.errorf("Chunk header length field is too big")
	}
	return asUint >> 1, asUint&1 == 1
}

func (_this *cbeDecodeBuffer) unexpectedEOD() {
	_this.errorf("Unexpected end of document")
}

func (_this *cbeDecodeBuffer) unexpectedError(err error) {
	// TODO: Diagnostics
	panic(err)
}

func (_this *cbeDecodeBuffer) errorf(format string, args ...interface{}) {
	// TODO: Diagnostics
	panic(fmt.Errorf(format, args...))
}
