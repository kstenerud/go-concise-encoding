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
	compact_float "github.com/kstenerud/go-compact-float"
	compact_time "github.com/kstenerud/go-compact-time"
	"github.com/kstenerud/go-concise-encoding/internal/common"
	"github.com/kstenerud/go-uleb128"
)

type Reader struct {
	reader io.Reader
	buffer []byte
}

func NewReader() *Reader {
	_this := &Reader{}
	_this.Init()
	return _this
}

func (_this *Reader) Init() {
	_this.buffer = make([]byte, decoderStartBufferSize)
}

func (_this *Reader) SetReader(reader io.Reader) {
	_this.reader = reader
}

func (_this *Reader) ReadUint8() uint8 {
	if _, err := _this.reader.Read(_this.buffer[:1]); err != nil {
		_this.unexpectedError(err)
	}
	return _this.buffer[0]
}

func (_this *Reader) ReadUint16() uint16 {
	_this.readIntoBuffer(2)
	return uint16(_this.buffer[0]) | uint16(_this.buffer[1])<<8
}

func (_this *Reader) ReadUint32() uint32 {
	_this.readIntoBuffer(4)
	return uint32(_this.buffer[0]) | uint32(_this.buffer[1])<<8 | uint32(_this.buffer[2])<<16 | uint32(_this.buffer[3])<<24
}

func (_this *Reader) ReadUint64() uint64 {
	_this.readIntoBuffer(8)
	return uint64(_this.buffer[0]) | uint64(_this.buffer[1])<<8 | uint64(_this.buffer[2])<<16 | uint64(_this.buffer[3])<<24 |
		uint64(_this.buffer[4])<<32 | uint64(_this.buffer[5])<<40 | uint64(_this.buffer[6])<<48 | uint64(_this.buffer[7])<<56
}

func (_this *Reader) ReadVersion() uint64 {
	return _this.readSmallULEB128("version", math.MaxUint64)
}

func (_this *Reader) ReadTypeOrEOF() cbeTypeField {
	if _, err := _this.reader.Read(_this.buffer[:1]); err != nil {
		if err == io.EOF {
			return cbeTypeEOF
		}
		_this.unexpectedError(err)
	}

	return cbeTypeField(_this.buffer[0])
}

func (_this *Reader) ReadType() cbeTypeField {
	return cbeTypeField(_this.ReadUint8())
}

func (_this *Reader) ReadUint() (asUint uint64, asBig *big.Int) {
	byteCount := _this.readSmallULEB128("uint length field", maxBigIntBitCount/8)
	bytes := _this.ReadBytes(int(byteCount))
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
	words := make([]big.Word, wordCount)
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

func (_this *Reader) ReadFloat16() float32 {
	return common.Float32FromFloat16Bits(_this.ReadUint16())
}

func (_this *Reader) ReadFloat32() float32 {
	return math.Float32frombits(_this.ReadUint32())
}

func (_this *Reader) ReadFloat64() float64 {
	return math.Float64frombits(_this.ReadUint64())
}

func (_this *Reader) ReadDecimalFloat() (compact_float.DFloat, *apd.Decimal) {
	value, bigValue, _, err := compact_float.DecodeWithByteBuffer(_this.reader, _this.buffer)
	if err != nil {
		_this.unexpectedError(err)
	}

	return value, bigValue
}

func (_this *Reader) ReadDate() compact_time.Time {
	value, _, err := compact_time.DecodeDateWithBuffer(_this.reader, _this.buffer)
	if err != nil {
		_this.unexpectedError(err)
	}

	return value
}

func (_this *Reader) ReadTime() compact_time.Time {
	value, _, err := compact_time.DecodeTimeWithBuffer(_this.reader, _this.buffer)
	if err != nil {
		_this.unexpectedError(err)
	}

	return value
}

func (_this *Reader) ReadTimestamp() compact_time.Time {
	value, _, err := compact_time.DecodeTimestampWithBuffer(_this.reader, _this.buffer)
	if err != nil {
		_this.unexpectedError(err)
	}

	return value
}

func (_this *Reader) ReadArrayChunkHeader() (elementCount uint64, moreChunksFollow bool) {
	header := _this.readSmallULEB128("array chunk header", math.MaxUint64)
	return header >> 1, header&1 == 1
}

func (_this *Reader) ReadBytes(byteCount int) []byte {
	_this.readIntoBuffer(byteCount)
	return _this.buffer[:byteCount]
}

func (_this *Reader) ReadIdentifier() []byte {
	length := int(_this.readSmallULEB128("identifier length", 100000))
	if length == 0 {
		panic(fmt.Errorf("identifier cannot be empty"))
	}
	return _this.ReadBytes(length)
}

// ============================================================================

// Internal

// TODO: Check max big.int bit count
const maxBigIntBitCount = 8192

const decoderStartBufferSize = 127

func (_this *Reader) readSmallULEB128(name string, maxValue uint64) uint64 {
	asUint, asBig, _, err := uleb128.DecodeWithByteBuffer(_this.reader, _this.buffer)
	if err != nil {
		_this.unexpectedError(err)
	}

	if asBig != nil {
		_this.errorf("%v: %v is too big (max allowed value = %v)", asBig, name, maxValue)
	}
	if asUint > maxValue {
		_this.errorf("%v: %v is too big (max allowed value = %v)", asBig, name, maxValue)
	}
	return asUint
}

func (_this *Reader) readIntoBuffer(count int) {
	_this.expandBufferTo(count)
	dst := _this.buffer[:count]
	for len(dst) > 0 {
		if bytesRead, err := _this.reader.Read(dst); err != nil {
			_this.unexpectedError(err)
		} else {
			dst = dst[bytesRead:]
		}
	}
}

func (_this *Reader) expandBufferTo(minSize int) {
	if len(_this.buffer) < minSize {
		_this.buffer = make([]byte, minSize*2)
	}
}

func (_this *Reader) unexpectedError(err error) {
	// TODO: Diagnostics
	panic(err)
}

func (_this *Reader) errorf(format string, args ...interface{}) {
	// TODO: Diagnostics
	panic(fmt.Errorf(format, args...))
}
