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

	"github.com/cockroachdb/apd/v2"
	compact_float "github.com/kstenerud/go-compact-float"
	compact_time "github.com/kstenerud/go-compact-time"
	"github.com/kstenerud/go-concise-encoding/conversions"
	"github.com/kstenerud/go-concise-encoding/events"
	"github.com/kstenerud/go-concise-encoding/internal/common"
	"github.com/kstenerud/go-concise-encoding/options"
)

// Receives data events, constructing a CBE document from them.
//
// Note: This is a LOW LEVEL API. Error reporting is done via panics. Be sure
// to recover() at an appropriate location when calling this struct's methods
// directly (with the exception of constructors and initializers, which are not
// designed to panic).
type Encoder struct {
	writer Writer
	opts   options.CBEEncoderOptions
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
	_this.writer.Init()
}

// Prepare the encoder for encoding. All events will be encoded to writer.
// PrepareToEncode MUST be called before using the encoder.
func (_this *Encoder) PrepareToEncode(writer io.Writer) {
	_this.writer.SetWriter(writer)
}

// ============================================================================

// DataEventReceiver

func (_this *Encoder) OnPadding(count int) {
	if count == 1 {
		_this.writer.WriteType(cbeTypePadding)
		return
	}

	_this.writer.ExpandBuffer(count)
	for i := 0; i < count; i++ {
		_this.writer.Buffer[i] = byte(cbeTypePadding)
	}
	_this.writer.FlushBuffer(count)
}

func (_this *Encoder) OnComment(bool, []byte) {
	// Ignored
}

func (_this *Encoder) OnBeginDocument() {
}

func (_this *Encoder) OnVersion(version uint64) {
	_this.writer.WriteByte(cbeDocumentHeader)
	_this.writer.WriteULEB(version)
}

func (_this *Encoder) OnNull() {
	_this.writer.WriteType(cbeTypeNull)
}

func (_this *Encoder) OnBool(value bool) {
	if value {
		_this.OnTrue()
	} else {
		_this.OnFalse()
	}
}

func (_this *Encoder) OnTrue() {
	_this.writer.WriteType(cbeTypeTrue)
}

func (_this *Encoder) OnFalse() {
	_this.writer.WriteType(cbeTypeFalse)
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
		_this.writer.WriteType(cbeTypeField(value))
	case fitsInUint8(value):
		_this.writer.WriteTyped8Bits(cbeTypePosInt8, uint8(value))
	case fitsInUint16(value):
		_this.writer.WriteTyped16Bits(cbeTypePosInt16, uint16(value))
	case fitsInUint32(value):
		_this.writer.WriteTyped32Bits(cbeTypePosInt32, uint32(value))
	case fitsInUint48(value):
		_this.writer.WriteTypedInt(cbeTypePosInt, value)
	default:
		_this.writer.WriteTyped64Bits(cbeTypePosInt64, value)
	}
}

func (_this *Encoder) OnNegativeInt(value uint64) {
	switch {
	case value == 0:
		_this.writer.WriteTyped8Bits(cbeTypeNegInt8, uint8(value))
	case fitsInSmallint(value):
		// Note: Must encode smallint using signed value
		_this.writer.WriteType(cbeTypeField(-int64(value)))
	case fitsInUint8(value):
		_this.writer.WriteTyped8Bits(cbeTypeNegInt8, uint8(value))
	case fitsInUint16(value):
		_this.writer.WriteTyped16Bits(cbeTypeNegInt16, uint16(value))
	case fitsInUint32(value):
		_this.writer.WriteTyped32Bits(cbeTypeNegInt32, uint32(value))
	case fitsInUint48(value):
		_this.writer.WriteTypedInt(cbeTypeNegInt, value)
	default:
		_this.writer.WriteTyped64Bits(cbeTypeNegInt64, value)
	}
}

func (_this *Encoder) OnBigInt(value *big.Int) {
	if value == nil {
		_this.writer.WriteType(cbeTypeNull)
		return
	}

	if common.IsBigIntNegative(value) {
		if value.IsInt64() {
			_this.OnNegativeInt(uint64(-value.Int64()))
			return
		}
		// TODO: -0x7fffffffffffffff to -0xffffffffffffffff don't need big int encoding
		_this.writer.WriteTypedBigInt(cbeTypeNegInt, value)
		return
	}

	if value.IsUint64() {
		_this.OnPositiveInt(value.Uint64())
		return
	}
	_this.writer.WriteTypedBigInt(cbeTypePosInt, value)
}

func (_this *Encoder) OnFloat(value float64) {
	if math.IsInf(value, 0) {
		sign := 1
		if value < 0 {
			sign = -1
		}
		_this.writer.WriteInfinity(sign)
		return
	}

	if math.IsNaN(value) {
		_this.writer.WriteNaN(common.IsSignalingNan(value))
		return
	}

	if value == 0 {
		sign := 1
		if math.Float64bits(value) == 0x8000000000000000 {
			sign = -1
		}
		_this.writer.WriteZero(sign)
		return
	}

	asfloat32 := float32(value)
	asFloat16 := math.Float32frombits(math.Float32bits(asfloat32) & 0xffff0000)

	if float64(asFloat16) == value {
		_this.writer.WriteFloat16(asFloat16)
		return
	}

	if float64(asfloat32) == value {
		_this.writer.WriteFloat32(asfloat32)
		return
	}

	_this.writer.WriteFloat64(value)
}

func (_this *Encoder) OnBigFloat(value *big.Float) {
	if value == nil {
		_this.writer.WriteType(cbeTypeNull)
		return
	}

	if f64, acc := value.Float64(); acc == big.Exact {
		_this.OnFloat(f64)
		return
	}

	// TODO: Big float rounding needs a configuration policy
	v, err := conversions.BigFloatToPBigDecimalFloat(value)
	if err != nil {
		_this.errorf("could not convert %v to apd.Decimal", value)
	}
	_this.OnBigDecimalFloat(v)
}

func (_this *Encoder) OnDecimalFloat(value compact_float.DFloat) {
	_this.writer.WriteDecimalFloat(value)
}

func (_this *Encoder) OnBigDecimalFloat(value *apd.Decimal) {
	if value == nil {
		_this.writer.WriteType(cbeTypeNull)
		return
	}

	_this.writer.WriteBigDecimalFloat(value)
}

func (_this *Encoder) OnNan(signaling bool) {
	_this.writer.WriteNaN(signaling)
}

func (_this *Encoder) OnUID(value []byte) {
	_this.writer.WriteType(cbeTypeUID)
	_this.writer.WriteBytes(value)
}

func (_this *Encoder) OnTime(value time.Time) {
	_this.writer.ExpandBuffer(compact_time.EncodedSizeGoTime(value) + 1)
	_this.writer.Buffer[0] = byte(cbeTypeTimestamp)
	count := compact_time.EncodeGoTimestampToBytes(value, _this.writer.Buffer[1:])
	_this.writer.FlushBuffer(count + 1)
}

var ctimeToCBEType = [...]cbeTypeField{
	compact_time.TimeTypeDate:      cbeTypeDate,
	compact_time.TimeTypeTime:      cbeTypeTime,
	compact_time.TimeTypeTimestamp: cbeTypeTimestamp,
}

func (_this *Encoder) OnCompactTime(value compact_time.Time) {
	if value.IsZeroValue() {
		_this.writer.WriteType(cbeTypeNull)
		return
	}

	_this.writer.ExpandBuffer(value.EncodedSize() + 1)
	_this.writer.Buffer[0] = byte(ctimeToCBEType[value.Type])
	count := value.EncodeToBytes(_this.writer.Buffer[1:])
	_this.writer.FlushBuffer(count + 1)
}

const maxSmallArrayLength = 15

type arrayTypeInfo struct {
	shortArrayType       cbeTypeField
	hasSmallArraySupport bool
	isPlane2             bool
}

var arrayInfo = [events.NumArrayTypes]arrayTypeInfo{
	events.ArrayTypeString: arrayTypeInfo{
		shortArrayType:       cbeTypeString0,
		hasSmallArraySupport: true,
		isPlane2:             false,
	},
	events.ArrayTypeUint16: arrayTypeInfo{
		shortArrayType:       cbeTypeShortArrayUint16,
		hasSmallArraySupport: true,
		isPlane2:             true,
	},
	events.ArrayTypeUint32: arrayTypeInfo{
		shortArrayType:       cbeTypeShortArrayUint32,
		hasSmallArraySupport: true,
		isPlane2:             true,
	},
	events.ArrayTypeUint64: arrayTypeInfo{
		shortArrayType:       cbeTypeShortArrayUint64,
		hasSmallArraySupport: true,
		isPlane2:             true,
	},
	events.ArrayTypeInt8: arrayTypeInfo{
		shortArrayType:       cbeTypeShortArrayInt8,
		hasSmallArraySupport: true,
		isPlane2:             true,
	},
	events.ArrayTypeInt16: arrayTypeInfo{
		shortArrayType:       cbeTypeShortArrayInt16,
		hasSmallArraySupport: true,
		isPlane2:             true,
	},
	events.ArrayTypeInt32: arrayTypeInfo{
		shortArrayType:       cbeTypeShortArrayInt32,
		hasSmallArraySupport: true,
		isPlane2:             true,
	},
	events.ArrayTypeInt64: arrayTypeInfo{
		shortArrayType:       cbeTypeShortArrayInt64,
		hasSmallArraySupport: true,
		isPlane2:             true,
	},
	events.ArrayTypeFloat16: arrayTypeInfo{
		shortArrayType:       cbeTypeShortArrayFloat16,
		hasSmallArraySupport: true,
		isPlane2:             true,
	},
	events.ArrayTypeFloat32: arrayTypeInfo{
		shortArrayType:       cbeTypeShortArrayFloat32,
		hasSmallArraySupport: true,
		isPlane2:             true,
	},
	events.ArrayTypeFloat64: arrayTypeInfo{
		shortArrayType:       cbeTypeShortArrayFloat64,
		hasSmallArraySupport: true,
		isPlane2:             true,
	},
	events.ArrayTypeUID: arrayTypeInfo{
		shortArrayType:       cbeTypeShortArrayUID,
		hasSmallArraySupport: true,
		isPlane2:             true,
	},
}

func (_this *Encoder) OnArray(arrayType events.ArrayType, elementCount uint64, value []byte) {
	if elementCount <= maxSmallArrayLength {
		info := arrayInfo[arrayType]
		if info.hasSmallArraySupport {
			if info.isPlane2 {
				_this.writer.WriteByte(cbeTypePlane2)
			}
			_this.writer.WriteByte(byte(info.shortArrayType) | byte(elementCount))
			_this.writer.WriteBytes(value)
			return
		}
	}

	_this.writer.WriteArrayHeader(arrayType)
	_this.writer.WriteArrayChunkHeader(elementCount, 0)
	_this.writer.WriteBytes(value)
}

func (_this *Encoder) OnStringlikeArray(arrayType events.ArrayType, value string) {
	elementCount := uint64(len(value))

	if arrayType == events.ArrayTypeString && elementCount <= maxSmallStringLength {
		_this.writer.WriteType(cbeTypeString0 + cbeTypeField(elementCount))
		_this.writer.WriteString(value)
	} else {
		_this.writer.WriteArrayHeader(arrayType)
		_this.writer.WriteArrayChunkHeader(elementCount, 0)
		_this.writer.WriteString(value)
	}
}

func (_this *Encoder) OnArrayBegin(arrayType events.ArrayType) {
	_this.writer.WriteArrayHeader(arrayType)
}

func (_this *Encoder) OnArrayChunk(elementCount uint64, moreChunksFollow bool) {
	continuationBit := uint64(0)
	if moreChunksFollow {
		continuationBit = 1
	}
	_this.writer.WriteArrayChunkHeader(elementCount, continuationBit)
}

func (_this *Encoder) OnArrayData(data []byte) {
	_this.writer.WriteBytes(data)
}

func (_this *Encoder) OnList() {
	_this.writer.WriteType(cbeTypeList)
}

func (_this *Encoder) OnMap() {
	_this.writer.WriteType(cbeTypeMap)
}

func (_this *Encoder) OnMarkup(id []byte) {
	_this.writer.WriteType(cbeTypeMarkup)
	_this.writer.WriteIdentifier(id)
}

func (_this *Encoder) OnNode() {
	_this.writer.WriteType(cbeTypeNode)
}

func (_this *Encoder) OnEdge() {
	_this.writer.WriteType(cbeTypeEdge)
}

func (_this *Encoder) OnEnd() {
	_this.writer.WriteType(cbeTypeEndContainer)
}

func (_this *Encoder) OnMarker(id []byte) {
	_this.writer.WriteType(cbeTypeMarker)
	_this.writer.WriteIdentifier(id)
}

func (_this *Encoder) OnReference(id []byte) {
	_this.writer.WriteType(cbeTypeReference)
	_this.writer.WriteIdentifier(id)
}

func (_this *Encoder) OnRemoteReference() {
	_this.writer.WriteType(cbeTypePlane2)
	_this.writer.WriteType(cbeTypeRemoteReference)
}

func (_this *Encoder) OnConstant(name []byte) {
	_this.errorf("CBE Cannot encode constant (%s)", string(name))
}

func (_this *Encoder) OnEndDocument() {
	// Nothing to do
}

// ============================================================================

// Internal

const (
	encoderStartBufferSize = 32

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

func (_this *Encoder) unexpectedError(err error, encoding interface{}) {
	_this.errorf("unexpected error [%v] while encoding %v", err, encoding)
}

func (_this *Encoder) errorf(format string, args ...interface{}) {
	panic(fmt.Errorf(format, args...))
}
