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
	"reflect"
	"time"

	"github.com/kstenerud/go-concise-encoding/buffer"
	"github.com/kstenerud/go-concise-encoding/conversions"
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
	buff           buffer.StreamingWriteBuffer
	opts           options.CBEEncoderOptions
	skipFirstMap   bool
	skipFirstList  bool
	containerDepth int
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
	_this.skipFirstList = _this.opts.ImpliedStructure == options.ImpliedStructureList
	_this.skipFirstMap = _this.opts.ImpliedStructure == options.ImpliedStructureMap
}

// Prepare the encoder for encoding. All events will be encoded to writer.
// PrepareToEncode MUST be called before using the encoder.
func (_this *Encoder) PrepareToEncode(writer io.Writer) {
	_this.buff.SetWriter(writer)
}

func (_this *Encoder) Reset() {
	_this.buff.Reset()
	_this.skipFirstList = _this.opts.ImpliedStructure == options.ImpliedStructureList
	_this.skipFirstMap = _this.opts.ImpliedStructure == options.ImpliedStructureMap
	_this.containerDepth = 0
}

// ============================================================================

// DataEventReceiver

func (_this *Encoder) OnPadding(count int) {
	dst := _this.buff.Allocate(count)
	for i := 0; i < count; i++ {
		dst[i] = byte(cbeTypePadding)
	}
}

func (_this *Encoder) OnBeginDocument() {
}

func (_this *Encoder) OnVersion(version uint64) {
	if _this.opts.ImpliedStructure != options.ImpliedStructureNone {
		return
	}
	_this.encodeULEB(version)
}

func (_this *Encoder) OnNil() {
	_this.encodeTypeOnly(cbeTypeNil)
}

func (_this *Encoder) OnBool(value bool) {
	if value {
		_this.OnTrue()
	} else {
		_this.OnFalse()
	}
}

func (_this *Encoder) OnTrue() {
	_this.encodeTypeOnly(cbeTypeTrue)
}

func (_this *Encoder) OnFalse() {
	_this.encodeTypeOnly(cbeTypeFalse)
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
		_this.encodeTypeOnly(cbeTypeField(value))
	case fitsInUint8(value):
		_this.encodeTyped8Bits(cbeTypePosInt8, uint8(value))
	case fitsInUint16(value):
		_this.encodeTyped16Bits(cbeTypePosInt16, uint16(value))
	case fitsInUint32(value):
		_this.encodeTyped32Bits(cbeTypePosInt32, uint32(value))
	case fitsInUint48(value):
		_this.encodeUint(cbeTypePosInt, value)
	default:
		_this.encodeTyped64Bits(cbeTypePosInt64, value)
	}
}

func (_this *Encoder) OnNegativeInt(value uint64) {
	switch {
	case fitsInSmallint(value):
		// Note: Must encode smallint using signed value
		_this.encodeTypeOnly(cbeTypeField(-int64(value)))
	case fitsInUint8(value):
		_this.encodeTyped8Bits(cbeTypeNegInt8, uint8(value))
	case fitsInUint16(value):
		_this.encodeTyped16Bits(cbeTypeNegInt16, uint16(value))
	case fitsInUint32(value):
		_this.encodeTyped32Bits(cbeTypeNegInt32, uint32(value))
	case fitsInUint48(value):
		_this.encodeUint(cbeTypeNegInt, value)
	default:
		_this.encodeTyped64Bits(cbeTypeNegInt64, value)
	}
}

func (_this *Encoder) OnBigInt(value *big.Int) {
	if common.IsBigIntNegative(value) {
		if value.IsInt64() {
			_this.OnNegativeInt(uint64(-value.Int64()))
			return
		}
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
	v, _, err := apd.NewFromString(conversions.BigFloatToString(value))
	if err != nil {
		panic(fmt.Errorf("could not convert %v to apd.Decimal", value))
	}
	_this.OnBigDecimalFloat(v)
}

func (_this *Encoder) OnDecimalFloat(value compact_float.DFloat) {
	_this.encodeDecimalFloat(value)
}

func (_this *Encoder) OnBigDecimalFloat(value *apd.Decimal) {
	_this.encodeBigDecimalFloat(value)
}

func (_this *Encoder) OnNan(signaling bool) {
	_this.encodeNaN(signaling)
}

func (_this *Encoder) OnUUID(value []byte) {
	dst := _this.buff.Allocate(17)
	dst[0] = byte(cbeTypeUUID)
	dst = dst[1:]
	copy(dst, value[:])
}

func (_this *Encoder) OnTime(value time.Time) {
	_this.OnCompactTime(compact_time.AsCompactTime(value))
}

func (_this *Encoder) OnCompactTime(value *compact_time.Time) {
	var timeType cbeTypeField
	switch value.TimeIs {
	case compact_time.TypeDate:
		timeType = cbeTypeDate
	case compact_time.TypeTime:
		timeType = cbeTypeTime
	case compact_time.TypeTimestamp:
		timeType = cbeTypeTimestamp
	}
	dst := _this.buff.Allocate(compact_time.MaxEncodeLength + 1)
	dst[0] = byte(timeType)
	dst = dst[1:]
	byteCount, _ := compact_time.Encode(value, dst)
	_this.buff.CorrectAllocation(byteCount + 1)
}

func (_this *Encoder) OnURI(value []byte) {
	_this.encodeTypedByteArray(cbeTypeURI, value)
}

func (_this *Encoder) OnString(value []byte) {
	stringLength := len(value)

	if stringLength > maxSmallStringLength {
		_this.encodeTypedByteArray(cbeTypeString, value)
		return
	}

	dst := _this.buff.Allocate(stringLength + 1)
	dst[0] = byte(cbeTypeString0 + cbeTypeField(stringLength))
	dst = dst[1:]
	copy(dst, value)
}

func (_this *Encoder) OnVerbatimString(value []byte) {
	_this.encodeTypedByteArray(cbeTypeVerbatimString, value)
}

func (_this *Encoder) OnCustomBinary(value []byte) {
	_this.encodeTypedByteArray(cbeTypeCustomBinary, value)
}

func (_this *Encoder) OnCustomText(value []byte) {
	_this.encodeTypedByteArray(cbeTypeCustomText, value)
}

func (_this *Encoder) OnTypedArray(elemType reflect.Type, value []byte) {
	// TODO: Typed array support
	_this.encodeArrayUint8(value)
}

func (_this *Encoder) OnStringBegin() {
	_this.encodeTypeOnly(cbeTypeString)
}

func (_this *Encoder) OnVerbatimStringBegin() {
	_this.encodeTypeOnly(cbeTypeVerbatimString)
}

func (_this *Encoder) OnURIBegin() {
	_this.encodeTypeOnly(cbeTypeURI)
}

func (_this *Encoder) OnCustomBinaryBegin() {
	_this.encodeTypeOnly(cbeTypeCustomBinary)
}

func (_this *Encoder) OnCustomTextBegin() {
	_this.encodeTypeOnly(cbeTypeCustomText)
}

func (_this *Encoder) OnTypedArrayBegin(elemType reflect.Type) {
	// TODO: Typed array support
	_this.encodeTypeOnly(cbeTypeArray)
	_this.encodeTypeOnly(cbeTypePosInt8)
}

func (_this *Encoder) OnArrayChunk(length uint64, moreChunksFollow bool) {
	continuationBit := uint64(0)
	if moreChunksFollow {
		continuationBit = 1
	}
	_this.encodeULEB((length << 1) | continuationBit)
}

func (_this *Encoder) OnArrayData(data []byte) {
	dst := _this.buff.Allocate(len(data))
	copy(dst, data)
}

func (_this *Encoder) OnList() {
	if _this.skipFirstList {
		_this.skipFirstList = false
		return
	}

	_this.encodeTypeOnly(cbeTypeList)
	_this.containerDepth++
}

func (_this *Encoder) OnMap() {
	if _this.skipFirstMap {
		_this.skipFirstMap = false
		return
	}

	_this.encodeTypeOnly(cbeTypeMap)
	_this.containerDepth++
}

func (_this *Encoder) OnMarkup() {
	_this.encodeTypeOnly(cbeTypeMarkup)
	_this.containerDepth += 2
}

func (_this *Encoder) OnMetadata() {
	_this.encodeTypeOnly(cbeTypeMetadata)
	_this.containerDepth++
}

func (_this *Encoder) OnComment() {
	_this.encodeTypeOnly(cbeTypeComment)
	_this.containerDepth++
}

func (_this *Encoder) OnEnd() {
	if _this.containerDepth <= 0 {
		return
	}
	_this.containerDepth--
	_this.encodeTypeOnly(cbeTypeEndContainer)
}

func (_this *Encoder) OnMarker() {
	_this.encodeTypeOnly(cbeTypeMarker)
}

func (_this *Encoder) OnReference() {
	_this.encodeTypeOnly(cbeTypeReference)
}

func (_this *Encoder) OnEndDocument() {
	_this.buff.Flush()
}

// ============================================================================

// Internal

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

func (_this *Encoder) encodeVersion(version uint64) {
	_this.encodeULEB(version)
}

func (_this *Encoder) encodeTypeOnly(value cbeTypeField) {
	_this.buff.Allocate(1)[0] = byte(value)
}

func (_this *Encoder) encodeULEB(value uint64) {
	dst := _this.buff.Allocate(uleb128.EncodedSizeUint64(value))
	byteCount, _ := uleb128.EncodeUint64(value, dst)
	_this.buff.CorrectAllocation(byteCount)
}

func (_this *Encoder) encodeTypedInt(cbeType cbeTypeField, value uint64) {
	byteCount := 0
	for accum := value; accum > 0; byteCount++ {
		accum >>= 8
	}
	dst := _this.buff.Allocate(byteCount + 2)
	dst[0] = byte(cbeType)
	dst[1] = byte(byteCount)
	dst = dst[2:]

	for i := 0; value > 0; i++ {
		dst[i] = byte(value)
		value >>= 8
	}
}

func (_this *Encoder) encodeTypedBigInt(cbeType cbeTypeField, value *big.Int) {
	if value == nil {
		_this.encodeTypeOnly(cbeTypeNil)
		return
	}
	words := value.Bits()
	lastWordByteCount := 0
	lastWord := words[len(words)-1]
	for lastWord != 0 {
		lastWordByteCount++
		lastWord >>= 8
	}
	bytesPerWord := common.BytesPerInt()
	byteCount := (len(words)-1)*bytesPerWord + lastWordByteCount
	_this.encodeTypeOnly(cbeType)
	_this.encodeULEB(uint64(byteCount))
	dst := _this.buff.Allocate(byteCount)
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
}

func (_this *Encoder) beginArrayChunk(byteCount int, dataOffset int, moreChunksFollow uint64) (bytes []byte, dataBegin int) {
	lengthField := uint64(byteCount<<1) | moreChunksFollow
	lengthFieldLength := uleb128.EncodedSizeUint64(lengthField)
	bytes = _this.buff.Allocate(lengthFieldLength + byteCount + dataOffset)
	encodedCount, _ := uleb128.EncodeUint64(lengthField, bytes[dataOffset:])
	dataBegin = dataOffset + encodedCount
	return
}

func (_this *Encoder) beginTypedArray(elemType cbeTypeField, firstChunkLength int, moreChunksFollow uint64) (dst []byte) {
	dst, dataBegin := _this.beginArrayChunk(firstChunkLength, 2, moreChunksFollow)
	dst[0] = byte(cbeTypeArray)
	dst[1] = byte(elemType)
	return dst[dataBegin:]
}

func (_this *Encoder) encodeArrayUint8(data []uint8) {
	dst := _this.beginTypedArray(cbeTypePosInt8, len(data), 0)
	copy(dst, data)
}

func (_this *Encoder) encodeTypedByteArray(cbeType cbeTypeField, data []byte) {
	dst, dataBegin := _this.beginArrayChunk(len(data), 1, 0)
	dst[0] = byte(cbeType)
	dst = dst[dataBegin:]
	copy(dst, data)
}

func (_this *Encoder) encodeTyped8Bits(typeValue cbeTypeField, value byte) {
	dst := _this.buff.Allocate(2)
	dst[0] = byte(typeValue)
	dst[1] = value
}

func (_this *Encoder) encodeTyped16Bits(typeValue cbeTypeField, value uint16) {
	dst := _this.buff.Allocate(3)
	dst[0] = byte(typeValue)
	dst[1] = byte(value)
	dst[2] = byte(value >> 8)
}

func (_this *Encoder) encodeTyped32Bits(typeValue cbeTypeField, value uint32) {
	dst := _this.buff.Allocate(5)
	dst[0] = byte(typeValue)
	dst[1] = byte(value)
	dst[2] = byte(value >> 8)
	dst[3] = byte(value >> 16)
	dst[4] = byte(value >> 24)
}

func (_this *Encoder) encodeTyped64Bits(typeValue cbeTypeField, value uint64) {
	dst := _this.buff.Allocate(9)
	dst[0] = byte(typeValue)
	dst[1] = byte(value)
	dst[2] = byte(value >> 8)
	dst[3] = byte(value >> 16)
	dst[4] = byte(value >> 24)
	dst[5] = byte(value >> 32)
	dst[6] = byte(value >> 40)
	dst[7] = byte(value >> 48)
	dst[8] = byte(value >> 56)
}

func (_this *Encoder) encodeUint(typeValue cbeTypeField, value uint64) {
	_this.encodeTypedInt(typeValue, value)
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
	dst := _this.buff.Allocate(compact_float.MaxEncodeLength() + 1)
	dst[0] = byte(cbeTypeDecimal)
	dst = dst[1:]
	byteCount, _ := compact_float.Encode(value, dst)
	_this.buff.CorrectAllocation(byteCount + 1)
}

func (_this *Encoder) encodeBigDecimalFloat(value *apd.Decimal) {
	dst := _this.buff.Allocate(compact_float.MaxEncodeLengthBig(value) + 1)
	dst[0] = byte(cbeTypeDecimal)
	dst = dst[1:]
	byteCount, _ := compact_float.EncodeBig(value, dst)
	_this.buff.CorrectAllocation(byteCount + 1)
}

func (_this *Encoder) encodeZero(sign int) {
	maxEncodedLength := 2
	dst := _this.buff.Allocate(maxEncodedLength + 1)
	dst[0] = byte(cbeTypeDecimal)
	dst = dst[1:]
	byteCount := 0
	if sign < 0 {
		byteCount, _ = compact_float.EncodeNegativeZero(dst)
	} else {
		byteCount, _ = compact_float.EncodeZero(dst)
	}
	_this.buff.CorrectAllocation(byteCount + 1)
}

func (_this *Encoder) encodeInfinity(sign int) {
	maxEncodedLength := 2
	dst := _this.buff.Allocate(maxEncodedLength + 1)
	dst[0] = byte(cbeTypeDecimal)
	dst = dst[1:]
	byteCount := 0
	if sign < 0 {
		byteCount, _ = compact_float.EncodeNegativeInfinity(dst)
	} else {
		byteCount, _ = compact_float.EncodeInfinity(dst)
	}
	_this.buff.CorrectAllocation(byteCount + 1)
}

func (_this *Encoder) encodeNaN(signaling bool) {
	maxEncodedLength := 2
	dst := _this.buff.Allocate(maxEncodedLength + 1)
	dst[0] = byte(cbeTypeDecimal)
	dst = dst[1:]
	byteCount := 0
	if signaling {
		byteCount, _ = compact_float.EncodeSignalingNan(dst)
	} else {
		byteCount, _ = compact_float.EncodeQuietNan(dst)
	}
	_this.buff.CorrectAllocation(byteCount + 1)
}
