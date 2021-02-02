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
	"time"

	"github.com/kstenerud/go-concise-encoding/buffer"
	"github.com/kstenerud/go-concise-encoding/conversions"
	"github.com/kstenerud/go-concise-encoding/events"
	"github.com/kstenerud/go-concise-encoding/internal/common"
	"github.com/kstenerud/go-concise-encoding/options"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
	"github.com/kstenerud/go-uleb128"
)

// Receives data events, constructing a CBE document from them.
//
// Note: This is a LOW LEVEL API. Error reporting is done via panics. Be sure
// to recover() at an appropriate location when calling this struct's methods
// directly (with the exception of constructors and initializers, which are not
// designed to panic).
type Encoder struct {
	buff buffer.StreamingWriteBuffer
	opts options.CBEEncoderOptions
}

// Create a new CBE encoder.
// If opts is nil, default options will be used.
func NewEncoder(opts *options.CBEEncoderOptions) *Encoder {
	_this := &Encoder{}
	_this.Init(opts)
	return _this
}

// Initialize this encoder.
// If opts is nil, default options will be used.
func (_this *Encoder) Init(opts *options.CBEEncoderOptions) {
	opts = opts.WithDefaultsApplied()
	_this.opts = *opts
	_this.buff.Init(_this.opts.BufferSize)
}

// Prepare the encoder for encoding. All events will be encoded to writer.
// PrepareToEncode MUST be called before using the encoder.
func (_this *Encoder) PrepareToEncode(writer io.Writer) {
	_this.buff.SetWriter(writer)
}

// ============================================================================

// DataEventReceiver

func (_this *Encoder) OnPadding(count int) {
	dst := _this.buff.RequireBytes(count)
	for i := 0; i < count; i++ {
		dst[i] = byte(cbeTypePadding)
	}
	_this.buff.UseBytes(count)
}

func (_this *Encoder) OnBeginDocument() {
}

func (_this *Encoder) OnVersion(version uint64) {
	_this.encodeByte(cbeDocumentHeader)
	_this.encodeULEB(version)
}

func (_this *Encoder) OnNA() {
	_this.encodeType(cbeTypeNA)
}

func (_this *Encoder) OnBool(value bool) {
	if value {
		_this.OnTrue()
	} else {
		_this.OnFalse()
	}
}

func (_this *Encoder) OnTrue() {
	_this.encodeType(cbeTypeTrue)
}

func (_this *Encoder) OnFalse() {
	_this.encodeType(cbeTypeFalse)
}

func (_this *Encoder) OnInt(value int64) {
	if value >= 0 {
		_this.OnPositiveInt(uint64(value))
	} else {
		_this.OnNegativeInt(uint64(-value))
	}
}

func (_this *Encoder) OnPositiveInt(value uint64) {
	switch {
	case fitsInSmallint(value):
		_this.encodeType(cbeTypeField(value))
	case fitsInUint8(value):
		_this.encodeTyped8Bits(cbeTypePosInt8, uint8(value))
	case fitsInUint16(value):
		_this.encodeTyped16Bits(cbeTypePosInt16, uint16(value))
	case fitsInUint32(value):
		_this.encodeTyped32Bits(cbeTypePosInt32, uint32(value))
	case fitsInUint48(value):
		_this.encodeTypedInt(cbeTypePosInt, value)
	default:
		_this.encodeTyped64Bits(cbeTypePosInt64, value)
	}
}

func (_this *Encoder) OnNegativeInt(value uint64) {
	switch {
	case fitsInSmallint(value):
		// Note: Must encode smallint using signed value
		_this.encodeType(cbeTypeField(-int64(value)))
	case fitsInUint8(value):
		_this.encodeTyped8Bits(cbeTypeNegInt8, uint8(value))
	case fitsInUint16(value):
		_this.encodeTyped16Bits(cbeTypeNegInt16, uint16(value))
	case fitsInUint32(value):
		_this.encodeTyped32Bits(cbeTypeNegInt32, uint32(value))
	case fitsInUint48(value):
		_this.encodeTypedInt(cbeTypeNegInt, value)
	default:
		_this.encodeTyped64Bits(cbeTypeNegInt64, value)
	}
}

func (_this *Encoder) OnBigInt(value *big.Int) {
	if value == nil {
		_this.encodeType(cbeTypeNA)
		return
	}

	if common.IsBigIntNegative(value) {
		if value.IsInt64() {
			_this.OnNegativeInt(uint64(-value.Int64()))
			return
		}
		// TODO: -0x7fffffffffffffff to -0xffffffffffffffff don't need big int encoding
		_this.encodeTypedBigInt(cbeTypeNegInt, value)
		return
	}

	if value.IsUint64() {
		_this.OnPositiveInt(value.Uint64())
		return
	}
	_this.encodeTypedBigInt(cbeTypePosInt, value)
}

func (_this *Encoder) OnFloat(value float64) {
	if math.IsInf(value, 0) {
		sign := 1
		if value < 0 {
			sign = -1
		}
		_this.encodeInfinity(sign)
		return
	}

	if math.IsNaN(value) {
		_this.encodeNaN(common.IsSignalingNan(value))
		return
	}

	if value == 0 {
		sign := 1
		if math.Float64bits(value) == 0x8000000000000000 {
			sign = -1
		}
		_this.encodeZero(sign)
		return
	}

	asfloat32 := float32(value)
	asFloat16 := math.Float32frombits(math.Float32bits(asfloat32) & 0xffff0000)

	if float64(asFloat16) == value {
		_this.encodeFloat16(asFloat16)
		return
	}

	if float64(asfloat32) == value {
		_this.encodeFloat32(asfloat32)
		return
	}

	_this.encodeFloat64(value)
}

func (_this *Encoder) OnBigFloat(value *big.Float) {
	if value == nil {
		_this.encodeType(cbeTypeNA)
		return
	}

	v, err := conversions.BigFloatToPBigDecimalFloat(value)
	if err != nil {
		_this.errorf("could not convert %v to apd.Decimal", value)
	}
	_this.OnBigDecimalFloat(v)
}

func (_this *Encoder) OnDecimalFloat(value compact_float.DFloat) {
	_this.encodeDecimalFloat(value)
}

func (_this *Encoder) OnBigDecimalFloat(value *apd.Decimal) {
	if value == nil {
		_this.encodeType(cbeTypeNA)
		return
	}

	_this.encodeBigDecimalFloat(value)
}

func (_this *Encoder) OnNan(signaling bool) {
	_this.encodeNaN(signaling)
}

func (_this *Encoder) OnUUID(value []byte) {
	const uuidSize = 16
	const dataOffset = 1
	dst := _this.buff.RequireBytes(uuidSize + dataOffset)
	dst[0] = byte(cbeTypeUUID)
	dst = dst[dataOffset:]
	copy(dst, value)
	_this.buff.UseBytes(uuidSize + dataOffset)
}

func (_this *Encoder) OnTime(value time.Time) {
	const dataOffset = 1
	dst := _this.buff.RequireBytes(compact_time.MaxEncodeLength + dataOffset)
	dst[0] = byte(cbeTypeTimestamp)
	dst = dst[dataOffset:]
	byteCount, _ := compact_time.EncodeGoTimestamp(value, dst)
	_this.buff.UseBytes(byteCount + dataOffset)
}

func (_this *Encoder) OnCompactTime(value compact_time.Time) {
	if value.IsZeroValue() {
		_this.encodeType(cbeTypeNA)
		_this.encodeType(cbeTypeNA)
		return
	}

	var timeType cbeTypeField
	switch value.TimeType {
	case compact_time.TypeDate:
		timeType = cbeTypeDate
	case compact_time.TypeTime:
		timeType = cbeTypeTime
	case compact_time.TypeTimestamp:
		timeType = cbeTypeTimestamp
	}
	const dataOffset = 1
	dst := _this.buff.RequireBytes(compact_time.MaxEncodeLength + dataOffset)
	dst[0] = byte(timeType)
	dst = dst[dataOffset:]
	byteCount, _ := value.Encode(dst)
	_this.buff.UseBytes(byteCount + dataOffset)
}

func (_this *Encoder) OnArray(arrayType events.ArrayType, elementCount uint64, value []byte) {
	if arrayType == events.ArrayTypeString && elementCount <= maxSmallStringLength {
		const dataOffset = 1
		dst := _this.buff.RequireBytes(int(elementCount) + dataOffset)
		dst[0] = byte(cbeTypeString0 + cbeTypeField(elementCount))
		dst = dst[dataOffset:]
		copy(dst, value)
		_this.buff.UseBytes(int(elementCount) + dataOffset)
		return
	}

	_this.encodeArrayHeader(arrayType)
	_this.encodeArrayChunkHeader(elementCount, 0)
	_this.encodeArrayData(value)
}

func (_this *Encoder) OnArrayBegin(arrayType events.ArrayType) {
	_this.encodeArrayHeader(arrayType)
}

func (_this *Encoder) OnArrayChunk(elementCount uint64, moreChunksFollow bool) {
	continuationBit := uint64(0)
	if moreChunksFollow {
		continuationBit = 1
	}
	_this.encodeArrayChunkHeader(elementCount, continuationBit)
}

func (_this *Encoder) OnArrayData(data []byte) {
	_this.encodeArrayData(data)
}

func (_this *Encoder) OnList() {
	_this.encodeType(cbeTypeList)
}

func (_this *Encoder) OnMap() {
	_this.encodeType(cbeTypeMap)
}

func (_this *Encoder) OnMarkup() {
	_this.encodeType(cbeTypeMarkup)
}

func (_this *Encoder) OnMetadata() {
	_this.encodeType(cbeTypeMetadata)
}

func (_this *Encoder) OnComment() {
	_this.encodeType(cbeTypeComment)
}

func (_this *Encoder) OnEnd() {
	_this.encodeType(cbeTypeEndContainer)
}

func (_this *Encoder) OnMarker() {
	_this.encodeType(cbeTypeMarker)
}

func (_this *Encoder) OnReference() {
	_this.encodeType(cbeTypeReference)
}

func (_this *Encoder) OnConcatenate() {
	_this.encodeType(cbeTypeConcatenate)
}

func (_this *Encoder) OnConstant(name []byte, explicitValue bool) {
	if !explicitValue {
		_this.errorf("Cannot encode constant %s without explicit value", string(name))
	}
}

func (_this *Encoder) OnEndDocument() {
	_this.buff.Flush()
	_this.reset()
}

// ============================================================================

// Internal

func (_this *Encoder) reset() {
	_this.buff.Reset()
}

const (
	maxSmallStringLength = 15

	bitMask48 = uint64((1 << 48) - 1)
)

func fitsInSmallint(value uint64) bool {
	return value <= uint64(cbeSmallIntMax)
}

func fitsInUint8(value uint64) bool {
	return value <= math.MaxUint8
}

func fitsInUint16(value uint64) bool {
	return value <= math.MaxUint16
}

func fitsInUint32(value uint64) bool {
	return value <= math.MaxUint32
}

func fitsInUint48(value uint64) bool {
	return value == (value & bitMask48)
}

func (_this *Encoder) encodeULEB(value uint64) {
	dst := _this.buff.RequireBytes(uleb128.EncodedSizeUint64(value))
	byteCount, _ := uleb128.EncodeUint64(value, dst)
	_this.buff.UseBytes(byteCount)
}

func (_this *Encoder) encodeByte(value byte) {
	_this.buff.AddByte(value)
}

func (_this *Encoder) encodeType(value cbeTypeField) {
	_this.encodeByte(byte(value))
}

func (_this *Encoder) encodeTyped8Bits(typeValue cbeTypeField, value byte) {
	dst := _this.buff.RequireBytes(2)
	dst[0] = byte(typeValue)
	dst[1] = value
	_this.buff.UseBytes(2)
}

func (_this *Encoder) encodeTyped16Bits(typeValue cbeTypeField, value uint16) {
	dst := _this.buff.RequireBytes(3)
	dst[0] = byte(typeValue)
	dst[1] = byte(value)
	dst[2] = byte(value >> 8)
	_this.buff.UseBytes(3)
}

func (_this *Encoder) encodeTyped32Bits(typeValue cbeTypeField, value uint32) {
	dst := _this.buff.RequireBytes(5)
	dst[0] = byte(typeValue)
	dst[1] = byte(value)
	dst[2] = byte(value >> 8)
	dst[3] = byte(value >> 16)
	dst[4] = byte(value >> 24)
	_this.buff.UseBytes(5)
}

func (_this *Encoder) encodeTyped64Bits(typeValue cbeTypeField, value uint64) {
	dst := _this.buff.RequireBytes(9)
	dst[0] = byte(typeValue)
	dst[1] = byte(value)
	dst[2] = byte(value >> 8)
	dst[3] = byte(value >> 16)
	dst[4] = byte(value >> 24)
	dst[5] = byte(value >> 32)
	dst[6] = byte(value >> 40)
	dst[7] = byte(value >> 48)
	dst[8] = byte(value >> 56)
	_this.buff.UseBytes(9)
}

func (_this *Encoder) encodeTypedInt(cbeType cbeTypeField, value uint64) {
	byteCount := 0
	for accum := value; accum > 0; byteCount++ {
		accum >>= 8
	}
	const dataOffset = 2
	dst := _this.buff.RequireBytes(byteCount + dataOffset)
	dst[0] = byte(cbeType)
	dst[1] = byte(byteCount)
	dst = dst[dataOffset:]

	for i := 0; value > 0; i++ {
		dst[i] = byte(value)
		value >>= 8
	}
	_this.buff.UseBytes(byteCount + dataOffset)
}

func (_this *Encoder) encodeTypedBigInt(cbeType cbeTypeField, value *big.Int) {
	if value == nil {
		_this.encodeType(cbeTypeNA)
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
	_this.encodeType(cbeType)
	_this.encodeULEB(uint64(byteCount))
	dst := _this.buff.RequireBytes(byteCount)
	iDst := 0
	for _, word := range words {
		for iPart := 0; iPart < bytesPerWord; iPart++ {
			dst[iDst] = byte(word)
			iDst++
			word >>= 8
			if iDst >= byteCount {
				break
			}
		}
	}
	_this.buff.UseBytes(byteCount)
}

func (_this *Encoder) encodeFloat16(value float32) {
	_this.encodeTyped16Bits(cbeTypeFloat16, uint16(math.Float32bits(value)>>16))
}

func (_this *Encoder) encodeFloat32(value float32) {
	_this.encodeTyped32Bits(cbeTypeFloat32, math.Float32bits(value))
}

func (_this *Encoder) encodeFloat64(value float64) {
	_this.encodeTyped64Bits(cbeTypeFloat64, math.Float64bits(value))
}

func (_this *Encoder) encodeDecimalFloat(value compact_float.DFloat) {
	const dataOffset = 1
	dst := _this.buff.RequireBytes(compact_float.MaxEncodeLength() + dataOffset)
	dst[0] = byte(cbeTypeDecimal)
	dst = dst[dataOffset:]
	byteCount, _ := compact_float.Encode(value, dst)
	_this.buff.UseBytes(byteCount + dataOffset)
}

func (_this *Encoder) encodeBigDecimalFloat(value *apd.Decimal) {
	const dataOffset = 1
	dst := _this.buff.RequireBytes(compact_float.MaxEncodeLengthBig(value) + dataOffset)
	dst[0] = byte(cbeTypeDecimal)
	dst = dst[dataOffset:]
	byteCount, _ := compact_float.EncodeBig(value, dst)
	_this.buff.UseBytes(byteCount + dataOffset)
}

func (_this *Encoder) encodeZero(sign int) {
	const maxEncodedLength = 2
	const dataOffset = 1
	dst := _this.buff.RequireBytes(maxEncodedLength + dataOffset)
	dst[0] = byte(cbeTypeDecimal)
	dst = dst[dataOffset:]
	byteCount := 0
	if sign < 0 {
		byteCount, _ = compact_float.EncodeNegativeZero(dst)
	} else {
		byteCount, _ = compact_float.EncodeZero(dst)
	}
	_this.buff.UseBytes(byteCount + dataOffset)
}

func (_this *Encoder) encodeInfinity(sign int) {
	const maxEncodedLength = 2
	const dataOffset = 1
	dst := _this.buff.RequireBytes(maxEncodedLength + dataOffset)
	dst[0] = byte(cbeTypeDecimal)
	dst = dst[dataOffset:]
	byteCount := 0
	if sign < 0 {
		byteCount, _ = compact_float.EncodeNegativeInfinity(dst)
	} else {
		byteCount, _ = compact_float.EncodeInfinity(dst)
	}
	_this.buff.UseBytes(byteCount + dataOffset)
}

func (_this *Encoder) encodeNaN(signaling bool) {
	const maxEncodedLength = 2
	const dataOffset = 1
	dst := _this.buff.RequireBytes(maxEncodedLength + dataOffset)
	dst[0] = byte(cbeTypeDecimal)
	dst = dst[dataOffset:]
	byteCount := 0
	if signaling {
		byteCount, _ = compact_float.EncodeSignalingNan(dst)
	} else {
		byteCount, _ = compact_float.EncodeQuietNan(dst)
	}
	_this.buff.UseBytes(byteCount + dataOffset)
}

func (_this *Encoder) encodeArrayHeader(arrayType events.ArrayType) {
	if isTypedArray[arrayType] {
		dst := _this.buff.RequireBytes(2)
		dst[0] = byte(cbeTypeArray)
		dst[1] = byte(arrayTypeToCBEType[arrayType])
		_this.buff.UseBytes(2)
	} else {
		dst := _this.buff.RequireBytes(1)
		dst[0] = byte(arrayTypeToCBEType[arrayType])
		_this.buff.UseBytes(1)
	}
}

func (_this *Encoder) encodeArrayChunkHeader(elementCount uint64, moreChunksFollow uint64) {
	_this.encodeULEB((elementCount << 1) | moreChunksFollow)
}

func (_this *Encoder) encodeArrayData(data []byte) {
	dst := _this.buff.RequireBytes(len(data))
	copy(dst, data)
	_this.buff.UseBytes(len(data))
}

func (_this *Encoder) unexpectedError(err error, encoding interface{}) {
	_this.errorf("unexpected error [%v] while encoding %v", err, encoding)
}

func (_this *Encoder) errorf(format string, args ...interface{}) {
	panic(fmt.Errorf(format, args...))
}
