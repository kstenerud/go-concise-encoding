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

func (this *cbeDecodeBuffer) Init(data []byte) {
	this.data = data
	this.position = 0
}

func (this *cbeDecodeBuffer) HasUnreadData() bool {
	return this.position < len(this.data)
}

func (this *cbeDecodeBuffer) DecodeUint8() uint8 {
	value := uint8(this.data[this.position])
	this.position++
	return value
}

func (this *cbeDecodeBuffer) DecodeUint16() uint16 {
	value := uint16(this.data[this.position]) |
		uint16(this.data[this.position+1])<<8
	this.position += 2
	return value
}

func (this *cbeDecodeBuffer) DecodeUint32() uint32 {
	value := uint32(this.data[this.position]) |
		uint32(this.data[this.position+1])<<8 |
		uint32(this.data[this.position+2])<<16 |
		uint32(this.data[this.position+3])<<24
	this.position += 4
	return value
}

func (this *cbeDecodeBuffer) DecodeUint64() uint64 {
	value := uint64(this.data[this.position]) |
		uint64(this.data[this.position+1])<<8 |
		uint64(this.data[this.position+2])<<16 |
		uint64(this.data[this.position+3])<<24 |
		uint64(this.data[this.position+4])<<32 |
		uint64(this.data[this.position+5])<<40 |
		uint64(this.data[this.position+6])<<48 |
		uint64(this.data[this.position+7])<<56
	this.position += 8
	return value
}

func (this *cbeDecodeBuffer) DecodeUint() (asUint uint64, asBig *big.Int) {
	asUint, asBig, bytesDecoded, isComplete := uleb128.Decode(0, 0, this.data[this.position:])
	if !isComplete {
		this.unexpectedEOD()
	}
	this.position += bytesDecoded
	return
}

func (this *cbeDecodeBuffer) DecodeVersion() uint64 {
	asUint, asBig := this.DecodeUint()
	if asBig != nil {
		this.errorf("Version too big")
	}
	return asUint
}

func (this *cbeDecodeBuffer) DecodeType() cbeTypeField {
	return cbeTypeField(this.DecodeUint8())
}

func (this *cbeDecodeBuffer) DecodeFloat32() float32 {
	return math.Float32frombits(this.DecodeUint32())
}

func (this *cbeDecodeBuffer) DecodeFloat64() float64 {
	return math.Float64frombits(this.DecodeUint64())
}

func (this *cbeDecodeBuffer) DecodeDecimalFloat() (compact_float.DFloat, *apd.Decimal) {
	value, bigValue, bytesDecoded, err := compact_float.Decode(this.data[this.position:])
	if err != nil {
		if err == compact_float.ErrorIncomplete {
			this.unexpectedEOD()
		}
		this.unexpectedError(err)
	}
	this.position += bytesDecoded
	return value, bigValue
}

func (this *cbeDecodeBuffer) DecodeDate() *compact_time.Time {
	value, bytesDecoded, err := compact_time.DecodeDate(this.data[this.position:])
	if err != nil {
		if err == compact_time.ErrorIncomplete {
			this.unexpectedEOD()
		}
		this.unexpectedError(err)
	}
	if err := value.Validate(); err != nil {
		this.unexpectedError(err)
	}
	this.position += bytesDecoded
	return value
}

func (this *cbeDecodeBuffer) DecodeTime() *compact_time.Time {
	value, bytesDecoded, isComplete := compact_time.DecodeTime(this.data[this.position:])
	if err := value.Validate(); err != nil {
		this.unexpectedError(err)
	}
	if !isComplete {
		this.unexpectedEOD()
	}
	this.position += bytesDecoded
	return value
}

func (this *cbeDecodeBuffer) DecodeTimestamp() *compact_time.Time {
	value, bytesDecoded, err := compact_time.DecodeTimestamp(this.data[this.position:])
	if err != nil {
		if err == compact_time.ErrorIncomplete {
			this.unexpectedEOD()
		}
		this.unexpectedError(err)
	}
	if err := value.Validate(); err != nil {
		this.unexpectedError(err)
	}
	this.position += bytesDecoded
	return value
}

func (this *cbeDecodeBuffer) DecodeBytes(byteCount int) []byte {
	value := this.data[this.position : this.position+byteCount]
	this.position += byteCount
	return value
}

func (this *cbeDecodeBuffer) DecodeChunkHeader() (length uint64, isFinalChunk bool) {
	asUint, asBig := this.DecodeUint()
	if asBig != nil {
		this.errorf("Chunk header length field is too big")
	}
	return asUint >> 1, asUint&1 == 1
}

func (this *cbeDecodeBuffer) unexpectedEOD() {
	this.errorf("Unexpected end of document")
}

func (this *cbeDecodeBuffer) unexpectedError(err error) {
	// TODO: Diagnostics
	panic(err)
}

func (this *cbeDecodeBuffer) errorf(format string, args ...interface{}) {
	// TODO: Diagnostics
	panic(fmt.Errorf(format, args...))
}
