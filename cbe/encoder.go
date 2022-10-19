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
	"github.com/kstenerud/go-concise-encoding/ce/events"
	"github.com/kstenerud/go-concise-encoding/conversions"
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
	opts   *options.CBEEncoderOptions
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
	if opts == nil {
		o := options.DefaultCBEEncoderOptions()
		opts = &o
	} else {
		opts.ApplyDefaults()
	}
	_this.opts = opts
	_this.writer.Init()
}

// Prepare the encoder for encoding. All events will be encoded to writer.
// PrepareToEncode MUST be called before using the encoder.
func (_this *Encoder) PrepareToEncode(writer io.Writer) {
	_this.writer.SetWriter(writer)
}

// ============================================================================

// DataEventReceiver

func (_this *Encoder) OnPadding() {
	_this.writer.WriteType(cbeTypePadding)
}

func (_this *Encoder) OnComment(bool, []byte) {
	// CBE doesn't have comments, so do nothing.
}

func (_this *Encoder) OnBeginDocument() {
	_this.writer.WriteSingleByte(CBESignatureByte)
}

func (_this *Encoder) OnVersion(version uint64) {
	_this.writer.WriteULEB(version)
}

func (_this *Encoder) OnNull() {
	_this.writer.WriteType(cbeTypeNull)
}

func (_this *Encoder) OnBoolean(value bool) {
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
		value = value.Neg(value)
		if value.IsUint64() {
			_this.OnNegativeInt(uint64(value.Uint64()))
			return
		}
		value = value.Neg(value)
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
		_this.writer.WriteNaN(!common.HasQuietNanBitSet64(value))
		return
	}

	if value == 0 {
		sign := 1
		if common.IsNegativeFloat(value) {
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
	if value.IsNegativeZero() {
		_this.writer.WriteZero(-1)
	} else if value.IsZero() {
		_this.writer.WriteZero(1)
	} else {
		_this.writer.WriteDecimalFloat(value)
	}
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

var ctimeToCBEType = [...]cbeTypeField{
	compact_time.TimeTypeDate:      cbeTypeDate,
	compact_time.TimeTypeTime:      cbeTypeTime,
	compact_time.TimeTypeTimestamp: cbeTypeTimestamp,
}

func (_this *Encoder) OnTime(value compact_time.Time) {
	if value.IsZeroValue() {
		_this.writer.WriteType(cbeTypeNull)
		return
	}

	_this.writer.ExpandBufferTo(value.EncodedSize() + 1)
	_this.writer.Buffer[0] = byte(ctimeToCBEType[value.Type])
	count := value.EncodeToBytes(_this.writer.Buffer[1:])
	_this.writer.FlushBufferFirstBytes(count + 1)
}

const maxSmallArrayLength = 15

type arrayTypeInfo struct {
	shortArrayType       cbeTypeField
	hasSmallArraySupport bool
	isPlane7f            bool
}

var arrayInfo = [events.NumArrayTypes]arrayTypeInfo{
	events.ArrayTypeString: {
		shortArrayType:       cbeTypeString0,
		hasSmallArraySupport: true,
		isPlane7f:            false,
	},
	events.ArrayTypeUint16: {
		shortArrayType:       cbeTypeShortArrayUint16,
		hasSmallArraySupport: true,
		isPlane7f:            true,
	},
	events.ArrayTypeUint32: {
		shortArrayType:       cbeTypeShortArrayUint32,
		hasSmallArraySupport: true,
		isPlane7f:            true,
	},
	events.ArrayTypeUint64: {
		shortArrayType:       cbeTypeShortArrayUint64,
		hasSmallArraySupport: true,
		isPlane7f:            true,
	},
	events.ArrayTypeInt8: {
		shortArrayType:       cbeTypeShortArrayInt8,
		hasSmallArraySupport: true,
		isPlane7f:            true,
	},
	events.ArrayTypeInt16: {
		shortArrayType:       cbeTypeShortArrayInt16,
		hasSmallArraySupport: true,
		isPlane7f:            true,
	},
	events.ArrayTypeInt32: {
		shortArrayType:       cbeTypeShortArrayInt32,
		hasSmallArraySupport: true,
		isPlane7f:            true,
	},
	events.ArrayTypeInt64: {
		shortArrayType:       cbeTypeShortArrayInt64,
		hasSmallArraySupport: true,
		isPlane7f:            true,
	},
	events.ArrayTypeFloat16: {
		shortArrayType:       cbeTypeShortArrayFloat16,
		hasSmallArraySupport: true,
		isPlane7f:            true,
	},
	events.ArrayTypeFloat32: {
		shortArrayType:       cbeTypeShortArrayFloat32,
		hasSmallArraySupport: true,
		isPlane7f:            true,
	},
	events.ArrayTypeFloat64: {
		shortArrayType:       cbeTypeShortArrayFloat64,
		hasSmallArraySupport: true,
		isPlane7f:            true,
	},
	events.ArrayTypeUID: {
		shortArrayType:       cbeTypeShortArrayUID,
		hasSmallArraySupport: true,
		isPlane7f:            true,
	},
}

func (_this *Encoder) OnArray(arrayType events.ArrayType, elementCount uint64, value []byte) {
	if elementCount <= maxSmallArrayLength {
		info := arrayInfo[arrayType]
		if info.hasSmallArraySupport {
			if info.isPlane7f {
				_this.writer.WriteType(cbeTypePlane7f)
			}
			_this.writer.WriteSingleByte(byte(info.shortArrayType) | byte(elementCount))
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

func (_this *Encoder) OnMedia(mediaType string, value []byte) {
	_this.OnMediaBegin(mediaType)
	_this.writer.WriteArrayChunkHeader(uint64(len(value)), 0)
	_this.writer.WriteBytes(value)
}

func (_this *Encoder) OnCustomBinary(customType uint64, value []byte) {
	_this.writer.WriteType(cbeTypeCustomType)
	_this.writer.WriteULEB(uint64(customType))
	_this.writer.WriteArrayChunkHeader(uint64(len(value)), 0)
	_this.writer.WriteBytes(value)
}

func (_this *Encoder) OnCustomText(customType uint64, value string) {
	panic(fmt.Errorf("CBE encoder cannot encode custom text"))
}

func (_this *Encoder) OnArrayBegin(arrayType events.ArrayType) {
	_this.writer.WriteArrayHeader(arrayType)
}

func (_this *Encoder) OnMediaBegin(mediaType string) {
	_this.writer.WriteType(cbeTypePlane7f)
	_this.writer.WriteType(cbeTypeMedia)
	_this.writer.WriteULEB(uint64(len(mediaType)))
	_this.writer.WriteString(mediaType)
}

func (_this *Encoder) OnCustomBegin(arrayType events.ArrayType, customType uint64) {
	_this.writer.WriteType(cbeTypeCustomType)
	_this.writer.WriteULEB(uint64(customType))
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

func (_this *Encoder) OnStructTemplate(id []byte) {
	_this.writer.WriteType(cbeTypePlane7f)
	_this.writer.WriteType(cbeTypeStructTemplate)
	_this.writer.WriteIdentifier(id)
}

func (_this *Encoder) OnStructInstance(id []byte) {
	_this.writer.WriteType(cbeTypeStructInstance)
	_this.writer.WriteIdentifier(id)
}

func (_this *Encoder) OnNode() {
	_this.writer.WriteType(cbeTypeNode)
}

func (_this *Encoder) OnEdge() {
	_this.writer.WriteType(cbeTypeEdge)
}

func (_this *Encoder) OnEndContainer() {
	_this.writer.WriteType(cbeTypeEndContainer)
}

func (_this *Encoder) OnMarker(id []byte) {
	_this.writer.WriteType(cbeTypePlane7f)
	_this.writer.WriteType(cbeTypeMarker)
	_this.writer.WriteIdentifier(id)
}

func (_this *Encoder) OnReferenceLocal(id []byte) {
	_this.writer.WriteType(cbeTypeLocalReference)
	_this.writer.WriteIdentifier(id)
}

func (_this *Encoder) OnRemoteReference() {
	_this.writer.WriteType(cbeTypePlane7f)
	_this.writer.WriteType(cbeTypeRemoteReference)
}

func (_this *Encoder) OnEndDocument() {
	// Nothing to do
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

func (_this *Encoder) errorf(format string, args ...interface{}) {
	panic(fmt.Errorf(format, args...))
}
