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
	"io"
	"math"
	"math/big"
	"unicode/utf8"

	"github.com/cockroachdb/apd/v2"
	compact_float "github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-concise-encoding/events"
	"github.com/kstenerud/go-concise-encoding/internal/common"
	"github.com/kstenerud/go-uleb128"
)

type Writer struct {
	writer       io.Writer
	stringWriter io.StringWriter
	adapter      StringWriterAdapter
	Buffer       []byte
}

func NewWriter() *Writer {
	_this := &Writer{}
	_this.Init()
	return _this
}

func (_this *Writer) Init() {
	_this.Buffer = make([]byte, writerStartBufferSize)
	_this.adapter.Init(_this)
}

// Set the actual io.Writer (or io.StringWriter) that will write the data.
func (_this *Writer) SetWriter(writer io.Writer) {
	_this.writer = writer
	if sw, ok := _this.writer.(io.StringWriter); ok {
		_this.stringWriter = sw
	} else {
		_this.stringWriter = &_this.adapter
	}
}

// Make sure the internal buffer is at least this big.
func (_this *Writer) ExpandBufferTo(minSize int) {
	if len(_this.Buffer) < minSize {
		_this.Buffer = make([]byte, minSize*2)
	}
}

// Flush count bytes from the start of the buffer
func (_this *Writer) FlushBufferFirstBytes(count int) {
	_this.writeBytes(_this.Buffer[:count])
}

// Flush the buffer from the start to the end position
func (_this *Writer) FlushBufferPortion(start, end int) {
	_this.writeBytes(_this.Buffer[start:end])
}

func (_this *Writer) WriteSingleByte(b byte) {
	_this.Buffer[0] = b
	_this.FlushBufferFirstBytes(1)
}

func (_this *Writer) WriteRune(r rune) {
	n := utf8.EncodeRune(_this.Buffer, r)
	_this.FlushBufferFirstBytes(n)
}

func (_this *Writer) WriteBytes(b []byte) {
	_this.writeBytes(b)
}

func (_this *Writer) WriteString(str string) {
	if _, err := _this.stringWriter.WriteString(str); err != nil {
		panic(err)
	}
}

func (_this *Writer) writeBytes(b []byte) {
	if _, err := _this.writer.Write(b); err != nil {
		panic(err)
	}
}

func (_this *Writer) WriteULEB(value uint64) {
	byteCount := uleb128.EncodeUint64ToBytes(value, _this.Buffer)
	_this.FlushBufferFirstBytes(byteCount)
}

func (_this *Writer) WriteType(t cbeTypeField) {
	_this.WriteSingleByte(byte(t))
}

func (_this *Writer) WriteTyped8Bits(typeValue cbeTypeField, value byte) {
	_this.Buffer[0] = byte(typeValue)
	_this.Buffer[1] = byte(value)
	_this.FlushBufferFirstBytes(2)
}

func (_this *Writer) WriteTyped16Bits(typeValue cbeTypeField, value uint16) {
	_this.Buffer[0] = byte(typeValue)
	_this.Buffer[1] = byte(value)
	_this.Buffer[2] = byte(value >> 8)
	_this.FlushBufferFirstBytes(3)
}

func (_this *Writer) WriteTyped32Bits(typeValue cbeTypeField, value uint32) {
	_this.Buffer[0] = byte(typeValue)
	_this.Buffer[1] = byte(value)
	_this.Buffer[2] = byte(value >> 8)
	_this.Buffer[3] = byte(value >> 16)
	_this.Buffer[4] = byte(value >> 24)
	_this.FlushBufferFirstBytes(5)
}

func (_this *Writer) WriteTyped64Bits(typeValue cbeTypeField, value uint64) {
	_this.Buffer[0] = byte(typeValue)
	_this.Buffer[1] = byte(value)
	_this.Buffer[2] = byte(value >> 8)
	_this.Buffer[3] = byte(value >> 16)
	_this.Buffer[4] = byte(value >> 24)
	_this.Buffer[5] = byte(value >> 32)
	_this.Buffer[6] = byte(value >> 40)
	_this.Buffer[7] = byte(value >> 48)
	_this.Buffer[8] = byte(value >> 56)
	_this.FlushBufferFirstBytes(9)
}

func (_this *Writer) WriteTypedInt(cbeType cbeTypeField, value uint64) {
	_this.Buffer[0] = byte(cbeType)

	byteCount := 0
	for accum := value; accum > 0; byteCount++ {
		_this.Buffer[byteCount+2] = byte(accum)
		accum >>= 8
	}
	_this.Buffer[1] = byte(byteCount)
	_this.FlushBufferFirstBytes(byteCount + 2)
}

func (_this *Writer) WriteTypedBigInt(cbeType cbeTypeField, value *big.Int) {
	if value == nil {
		_this.WriteType(cbeTypeNull)
		return
	}
	words := value.Bits()
	lastWordByteCount := 0
	lastWord := words[len(words)-1]
	for lastWord != 0 {
		lastWordByteCount++
		lastWord >>= 8
	}
	bytesPerWord := common.BytesPerInt
	byteCount := (len(words)-1)*bytesPerWord + lastWordByteCount
	_this.WriteType(cbeType)
	_this.WriteULEB(uint64(byteCount))
	buff := make([]byte, byteCount)
	iBuff := 0
	for _, word := range words {
		for iPart := 0; iPart < bytesPerWord; iPart++ {
			buff[iBuff] = byte(word)
			iBuff++
			word >>= 8
			if iBuff >= byteCount {
				break
			}
		}
	}
	_this.WriteBytes(buff)
}

func (_this *Writer) WriteFloat16(value float32) {
	_this.WriteTyped16Bits(cbeTypeFloat16, uint16(math.Float32bits(value)>>16))
}

func (_this *Writer) WriteFloat32(value float32) {
	_this.WriteTyped32Bits(cbeTypeFloat32, math.Float32bits(value))
}

func (_this *Writer) WriteFloat64(value float64) {
	_this.WriteTyped64Bits(cbeTypeFloat64, math.Float64bits(value))
}

func (_this *Writer) WriteDecimalFloat(value compact_float.DFloat) {
	_this.Buffer[0] = byte(cbeTypeDecimal)
	count := compact_float.EncodeToBytes(value, _this.Buffer[1:])
	_this.FlushBufferFirstBytes(count + 1)
}

func (_this *Writer) WriteBigDecimalFloat(value *apd.Decimal) {
	_this.ExpandBufferTo(compact_float.MaxEncodeLengthBig(value) + 1)

	_this.Buffer[0] = byte(cbeTypeDecimal)
	count := compact_float.EncodeBigToBytes(value, _this.Buffer[1:])
	_this.FlushBufferFirstBytes(count + 1)
}

func (_this *Writer) WriteZero(sign int) {
	if sign < 0 {
		_this.WriteTyped8Bits(cbeTypeNegInt8, 0)
	} else {
		_this.WriteType(cbeTypeField(0))
	}
}

func (_this *Writer) WriteInfinity(sign int) {
	_this.Buffer[0] = byte(cbeTypeDecimal)

	var count int
	if sign < 0 {
		count = compact_float.EncodeNegativeInfinity(_this.Buffer[1:])
	} else {
		count = compact_float.EncodeInfinity(_this.Buffer[1:])
	}
	_this.FlushBufferFirstBytes(count + 1)
}

func (_this *Writer) WriteNaN(signaling bool) {
	_this.Buffer[0] = byte(cbeTypeDecimal)

	var count int
	if signaling {
		count = compact_float.EncodeSignalingNan(_this.Buffer[1:])
	} else {
		count = compact_float.EncodeQuietNan(_this.Buffer[1:])
	}
	_this.FlushBufferFirstBytes(count + 1)
}

func (_this *Writer) WriteIdentifier(b []byte) {
	_this.WriteSingleByte(byte(len(b)))
	_this.WriteBytes(b)
}

func (_this *Writer) WriteArrayHeader(arrayType events.ArrayType) {
	byteCount := _this.WriteArrayHeaderToBytes(arrayType, _this.Buffer)
	_this.FlushBufferFirstBytes(byteCount)
}

func (_this *Writer) WriteArrayChunkHeader(elementCount uint64, moreChunksFollow uint64) {
	_this.WriteULEB((elementCount << 1) | moreChunksFollow)
}

func (_this *Writer) WriteArrayAndChunkHeader(arrayType events.ArrayType, elementCount uint64, moreChunksFollow uint64) {
	byteCount := _this.WriteArrayHeaderToBytes(arrayType, _this.Buffer)
	byteCount += _this.WriteArrayChunkHeaderToBytes(elementCount, moreChunksFollow, _this.Buffer[byteCount:])
	_this.FlushBufferFirstBytes(byteCount)
}

func (_this *Writer) WriteArrayHeaderToBytes(arrayType events.ArrayType, buffer []byte) int {
	if isPlane7fArray[arrayType] {
		_this.Buffer[0] = byte(cbeTypePlane7f)
		_this.Buffer[1] = byte(arrayTypeToCBEType[arrayType])
		return 2
	} else {
		_this.Buffer[0] = byte(arrayTypeToCBEType[arrayType])
		return 1
	}
}

func (_this *Writer) WriteArrayChunkHeaderToBytes(elementCount uint64, moreChunksFollow uint64, buffer []byte) int {
	value := (elementCount << 1) | moreChunksFollow
	return uleb128.EncodeUint64ToBytes(value, _this.Buffer)
}

// ============================================================================

const writerStartBufferSize = 32

type StringWriterAdapter struct {
	writer *Writer
}

func (_this *StringWriterAdapter) Init(writer *Writer) {
	_this.writer = writer
}

func (_this *StringWriterAdapter) WriteString(str string) (n int, err error) {
	n = len(str)
	_this.writer.ExpandBufferTo(n)
	copy(_this.writer.Buffer, str)
	_this.writer.FlushBufferFirstBytes(n)
	return
}
