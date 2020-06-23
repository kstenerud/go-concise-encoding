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
	"time"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
	"github.com/kstenerud/go-uleb128"
)

type CBEEncoderOptions struct {
}

// Receives data events, constructing a CBE document from them.
type CBEEncoder struct {
	buff    buffer
	options CBEEncoderOptions
}

func NewCBEEncoder(options *CBEEncoderOptions) *CBEEncoder {
	_this := &CBEEncoder{}
	_this.Init(options)
	return _this
}

func (_this *CBEEncoder) Init(options *CBEEncoderOptions) {
	if options != nil {
		_this.options = *options
	}
}

// Get the document that resulted from the data events this encoder received.
func (_this *CBEEncoder) GetBuiltDocument() []byte {
	return _this.buff.bytes
}

func (_this *CBEEncoder) OnPadding(count int) {
	dst := _this.buff.Allocate(count)
	for i := 0; i < count; i++ {
		dst[i] = byte(cbeTypePadding)
	}
}

func (_this *CBEEncoder) OnVersion(version uint64) {
	_this.encodeULEB(version)
}

func (_this *CBEEncoder) OnNil() {
	_this.encodeTypeOnly(cbeTypeNil)
}

func (_this *CBEEncoder) OnBool(value bool) {
	if value {
		_this.OnTrue()
	} else {
		_this.OnFalse()
	}
}

func (_this *CBEEncoder) OnTrue() {
	_this.encodeTypeOnly(cbeTypeTrue)
}

func (_this *CBEEncoder) OnFalse() {
	_this.encodeTypeOnly(cbeTypeFalse)
}

func (_this *CBEEncoder) OnInt(value int64) {
	if value >= 0 {
		_this.OnPositiveInt(uint64(value))
	} else {
		_this.OnNegativeInt(uint64(-value))
	}
}

func (_this *CBEEncoder) OnPositiveInt(value uint64) {
	switch {
	case fitsInSmallint(value):
		_this.encodeTypeOnly(cbeTypeField(value))
	case fitsInUint8(value):
		_this.encodeTyped8Bits(cbeTypePosInt8, uint8(value))
	case fitsInUint16(value):
		_this.encodeTyped16Bits(cbeTypePosInt16, uint16(value))
	case fitsInUint21(value):
		_this.encodeUint(cbeTypePosInt, value)
	case fitsInUint32(value):
		_this.encodeTyped32Bits(cbeTypePosInt32, uint32(value))
	case fitsInUint49(value):
		_this.encodeUint(cbeTypePosInt, value)
	default:
		_this.encodeTyped64Bits(cbeTypePosInt64, value)
	}
}

func (_this *CBEEncoder) OnNegativeInt(value uint64) {
	switch {
	case fitsInSmallint(value):
		// Note: Must encode smallint using signed value
		_this.encodeTypeOnly(cbeTypeField(-int64(value)))
	case fitsInUint8(value):
		_this.encodeTyped8Bits(cbeTypeNegInt8, uint8(value))
	case fitsInUint16(value):
		_this.encodeTyped16Bits(cbeTypeNegInt16, uint16(value))
	case fitsInUint21(value):
		_this.encodeUint(cbeTypeNegInt, value)
	case fitsInUint32(value):
		_this.encodeTyped32Bits(cbeTypeNegInt32, uint32(value))
	case fitsInUint49(value):
		_this.encodeUint(cbeTypeNegInt, value)
	default:
		_this.encodeTyped64Bits(cbeTypeNegInt64, value)
	}
}

func (_this *CBEEncoder) OnBigInt(value *big.Int) {
	if isBigIntNegative(value) {
		_this.encodeTypedBigInt(cbeTypeNegInt, value)
	} else {
		_this.encodeTypedBigInt(cbeTypePosInt, value)
	}
}

func (_this *CBEEncoder) OnFloat(value float64) {
	if math.IsInf(value, 0) {
		sign := 1
		if value < 0 {
			sign = -1
		}
		_this.encodeInfinity(sign)
		return
	}

	if math.IsNaN(value) {
		_this.encodeNaN(isSignalingNan(value))
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
	if float64(asfloat32) == value {
		_this.encodeFloat32(asfloat32)
		return
	}

	_this.encodeFloat64(value)
}

func (_this *CBEEncoder) OnBigFloat(value *big.Float) {
	v, _, err := apd.NewFromString(bigFloatToString(value))
	if err != nil {
		panic(fmt.Errorf("Could not convert %v to apd.Decimal", value))
	}
	_this.OnBigDecimalFloat(v)
}

func (_this *CBEEncoder) OnDecimalFloat(value compact_float.DFloat) {
	_this.encodeDecimalFloat(value)
}

func (_this *CBEEncoder) OnBigDecimalFloat(value *apd.Decimal) {
	_this.encodeBigDecimalFloat(value)
}

func (_this *CBEEncoder) OnNan(signaling bool) {
	_this.encodeNaN(signaling)
}

func (_this *CBEEncoder) OnUUID(value []byte) {
	dst := _this.buff.Allocate(17)
	dst[0] = byte(cbeTypeUUID)
	dst = dst[1:]
	copy(dst, value[:])
}

func (_this *CBEEncoder) OnComplex(value complex128) {
	panic("TODO: CBEEncoder.OnComplex")
}

func (_this *CBEEncoder) OnTime(value time.Time) {
	_this.OnCompactTime(compact_time.AsCompactTime(value))
}

func (_this *CBEEncoder) OnCompactTime(value *compact_time.Time) {
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

func (_this *CBEEncoder) OnBytes(value []byte) {
	_this.encodeTypedBytes(cbeTypeBytes, value)
}

func (_this *CBEEncoder) OnURI(value string) {
	_this.encodeTypedBytes(cbeTypeURI, []byte(value))
}

func (_this *CBEEncoder) OnString(value string) {
	bytes := []byte(value)
	stringLength := len(bytes)

	if stringLength > maxSmallStringLength {
		_this.encodeTypedBytes(cbeTypeString, bytes)
		return
	}

	dst := _this.buff.Allocate(stringLength + 1)
	dst[0] = byte(cbeTypeString0 + cbeTypeField(stringLength))
	dst = dst[1:]
	copy(dst, bytes)
}

func (_this *CBEEncoder) OnCustom(value []byte) {
	_this.encodeTypedBytes(cbeTypeCustom, value)
}

func (_this *CBEEncoder) OnBytesBegin() {
	_this.encodeTypeOnly(cbeTypeBytes)
}

func (_this *CBEEncoder) OnStringBegin() {
	_this.encodeTypeOnly(cbeTypeString)
}

func (_this *CBEEncoder) OnURIBegin() {
	_this.encodeTypeOnly(cbeTypeURI)
}

func (_this *CBEEncoder) OnCustomBegin() {
	_this.encodeTypeOnly(cbeTypeCustom)
}

func (_this *CBEEncoder) OnArrayChunk(length uint64, isFinalChunk bool) {
	continuationBit := uint64(0)
	if isFinalChunk {
		continuationBit = 1
	}
	_this.encodeULEB((uint64(length) << 1) | continuationBit)
}

func (_this *CBEEncoder) OnArrayData(data []byte) {
	dst := _this.buff.Allocate(len(data))
	copy(dst, data)
}

func (_this *CBEEncoder) OnList() {
	_this.encodeTypeOnly(cbeTypeList)
}

func (_this *CBEEncoder) OnMap() {
	_this.encodeTypeOnly(cbeTypeMap)
}

func (_this *CBEEncoder) OnMarkup() {
	_this.encodeTypeOnly(cbeTypeMarkup)
}

func (_this *CBEEncoder) OnMetadata() {
	_this.encodeTypeOnly(cbeTypeMetadata)
}

func (_this *CBEEncoder) OnComment() {
	_this.encodeTypeOnly(cbeTypeComment)
}

func (_this *CBEEncoder) OnEnd() {
	_this.encodeTypeOnly(cbeTypeEndContainer)
}

func (_this *CBEEncoder) OnMarker() {
	_this.encodeTypeOnly(cbeTypeMarker)
}

func (_this *CBEEncoder) OnReference() {
	_this.encodeTypeOnly(cbeTypeReference)
}

func (_this *CBEEncoder) OnEndDocument() {
}

// ============================================================================

const (
	minBufferCap         = 64
	maxSmallStringLength = 15

	bitMask21 = uint64((1 << 21) - 1)
	bitMask49 = uint64((1 << 49) - 1)
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

func fitsInUint21(value uint64) bool {
	return value == (value & bitMask21)
}

func fitsInUint32(value uint64) bool {
	return value <= math.MaxUint32
}

func fitsInUint49(value uint64) bool {
	return value == (value & bitMask49)
}

func (_this *CBEEncoder) encodeVersion(version uint64) {
	_this.encodeULEB(version)
}

func (_this *CBEEncoder) encodeTypeOnly(value cbeTypeField) {
	_this.buff.Allocate(1)[0] = byte(value)
}

func (_this *CBEEncoder) encodeULEB(value uint64) {
	dst := _this.buff.Allocate(uleb128.EncodedSizeUint64(value))
	byteCount, _ := uleb128.EncodeUint64(value, dst)
	_this.buff.CorrectAllocation(byteCount)
}

func (_this *CBEEncoder) encodeTypedULEB(cbeType cbeTypeField, value uint64) {
	dst := _this.buff.Allocate(uleb128.EncodedSizeUint64(value) + 1)
	dst[0] = byte(cbeType)
	dst = dst[1:]
	byteCount, _ := uleb128.EncodeUint64(value, dst)
	_this.buff.CorrectAllocation(byteCount + 1)
}

func (_this *CBEEncoder) encodeTypedBigInt(cbeType cbeTypeField, value *big.Int) {
	if value == nil {
		_this.encodeTypeOnly(cbeTypeNil)
		return
	}
	dst := _this.buff.Allocate(uleb128.EncodedSize(value) + 1)
	dst[0] = byte(cbeType)
	dst = dst[1:]
	byteCount, _ := uleb128.Encode(value, dst)
	_this.buff.CorrectAllocation(byteCount + 1)
}

func (_this *CBEEncoder) encodeTypedBytes(cbeType cbeTypeField, bytes []byte) {
	bytesLength := len(bytes)
	lengthField := uint64(bytesLength<<1) | 1
	dst := _this.buff.Allocate(uleb128.EncodedSizeUint64(lengthField) + bytesLength + 1)
	dst[0] = byte(cbeType)
	dst = dst[1:]
	byteCount, _ := uleb128.EncodeUint64(lengthField, dst)
	dst = dst[byteCount:]
	copy(dst, bytes)
	_this.buff.CorrectAllocation(byteCount + bytesLength + 1)
}

func (_this *CBEEncoder) encodeTyped8Bits(typeValue cbeTypeField, value byte) {
	dst := _this.buff.Allocate(2)
	dst[0] = byte(typeValue)
	dst[1] = value
}

func (_this *CBEEncoder) encodeTyped16Bits(typeValue cbeTypeField, value uint16) {
	dst := _this.buff.Allocate(3)
	dst[0] = byte(typeValue)
	dst[1] = byte(value)
	dst[2] = byte(value >> 8)
}

func (_this *CBEEncoder) encodeTyped32Bits(typeValue cbeTypeField, value uint32) {
	dst := _this.buff.Allocate(5)
	dst[0] = byte(typeValue)
	dst[1] = byte(value)
	dst[2] = byte(value >> 8)
	dst[3] = byte(value >> 16)
	dst[4] = byte(value >> 24)
}

func (_this *CBEEncoder) encodeTyped64Bits(typeValue cbeTypeField, value uint64) {
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

func (_this *CBEEncoder) encodeUint(typeValue cbeTypeField, value uint64) {
	_this.encodeTypedULEB(typeValue, value)
}

func (_this *CBEEncoder) encodeFloat32(value float32) {
	_this.encodeTyped32Bits(cbeTypeFloat32, math.Float32bits(value))
}

func (_this *CBEEncoder) encodeFloat64(value float64) {
	_this.encodeTyped64Bits(cbeTypeFloat64, math.Float64bits(value))
}

func (_this *CBEEncoder) encodeDecimalFloat(value compact_float.DFloat) {
	dst := _this.buff.Allocate(compact_float.MaxEncodeLength() + 1)
	dst[0] = byte(cbeTypeDecimal)
	dst = dst[1:]
	byteCount, _ := compact_float.Encode(value, dst)
	_this.buff.CorrectAllocation(byteCount + 1)
}

func (_this *CBEEncoder) encodeBigDecimalFloat(value *apd.Decimal) {
	dst := _this.buff.Allocate(compact_float.MaxEncodeLengthBig(value) + 1)
	dst[0] = byte(cbeTypeDecimal)
	dst = dst[1:]
	byteCount, _ := compact_float.EncodeBig(value, dst)
	_this.buff.CorrectAllocation(byteCount + 1)
}

func (_this *CBEEncoder) encodeZero(sign int) {
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

func (_this *CBEEncoder) encodeInfinity(sign int) {
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

func (_this *CBEEncoder) encodeNaN(signaling bool) {
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
