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

package cte

import (
	"math/big"
	"time"

	"github.com/kstenerud/go-concise-encoding/events"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
)

type listEncoder struct{}

var globalListEncoder listEncoder

func (_this *listEncoder) String() string { return "listEncoder" }

func (_this *listEncoder) prepareToWrite(ctx *EncoderContext) {
	ctx.ContainerHasObjects = true
	ctx.WriteCurrentPrefix()
}

func (_this *listEncoder) Begin(ctx *EncoderContext) {
	ctx.Stream.WriteListBegin()
	ctx.IncreaseIndent()
	ctx.SetStandardIndentPrefix()
	ctx.ContainerHasObjects = false
}

func (_this *listEncoder) End(ctx *EncoderContext) {
	ctx.DecreaseIndent()
	if ctx.ContainerHasObjects {
		ctx.WriteBasicIndent()
	}
	ctx.Stream.WriteListEnd()
	ctx.Unstack()
	ctx.CurrentEncoder.ChildContainerFinished(ctx)
}

func (_this *listEncoder) ChildContainerFinished(ctx *EncoderContext) {
	ctx.ContainerHasObjects = true
	ctx.SetStandardIndentPrefix()
}

func (_this *listEncoder) EncodeBool(ctx *EncoderContext, value bool) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteBool(value)
}
func (_this *listEncoder) EncodeTrue(ctx *EncoderContext) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteTrue()
}
func (_this *listEncoder) EncodeFalse(ctx *EncoderContext) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteFalse()
}
func (_this *listEncoder) EncodePositiveInt(ctx *EncoderContext, value uint64) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WritePositiveInt(value)
}
func (_this *listEncoder) EncodeNegativeInt(ctx *EncoderContext, value uint64) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteNegativeInt(value)
}
func (_this *listEncoder) EncodeInt(ctx *EncoderContext, value int64) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteInt(value)
}
func (_this *listEncoder) EncodeBigInt(ctx *EncoderContext, value *big.Int) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteBigInt(value)
}
func (_this *listEncoder) EncodeFloat(ctx *EncoderContext, value float64) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteFloat(value)
}
func (_this *listEncoder) EncodeBigFloat(ctx *EncoderContext, value *big.Float) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteBigFloat(value)
}
func (_this *listEncoder) EncodeDecimalFloat(ctx *EncoderContext, value compact_float.DFloat) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteDecimalFloat(value)
}
func (_this *listEncoder) EncodeBigDecimalFloat(ctx *EncoderContext, value *apd.Decimal) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteBigDecimalFloat(value)
}
func (_this *listEncoder) EncodeNan(ctx *EncoderContext, signaling bool) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteNan(signaling)
}
func (_this *listEncoder) EncodeTime(ctx *EncoderContext, value time.Time) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteTime(value)
}
func (_this *listEncoder) EncodeCompactTime(ctx *EncoderContext, value compact_time.Time) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteCompactTime(value)
}
func (_this *listEncoder) EncodeUUID(ctx *EncoderContext, value []byte) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteUUID(value)
}
func (_this *listEncoder) BeginList(ctx *EncoderContext) {
	_this.prepareToWrite(ctx)
	ctx.BeginStandardList()
}
func (_this *listEncoder) BeginMap(ctx *EncoderContext) {
	_this.prepareToWrite(ctx)
	ctx.BeginStandardMap()
}
func (_this *listEncoder) BeginMarkup(ctx *EncoderContext) {
	_this.prepareToWrite(ctx)
	ctx.BeginStandardMarkup()
}
func (_this *listEncoder) BeginMetadata(ctx *EncoderContext) {
	_this.prepareToWrite(ctx)
	ctx.BeginStandardMetadata()
}
func (_this *listEncoder) BeginComment(ctx *EncoderContext) {
	_this.prepareToWrite(ctx)
	ctx.BeginStandardComment()
}
func (_this *listEncoder) BeginMarker(ctx *EncoderContext) {
	_this.prepareToWrite(ctx)
	ctx.BeginStandardMarker()
}
func (_this *listEncoder) BeginReference(ctx *EncoderContext) {
	_this.prepareToWrite(ctx)
	ctx.BeginStandardReference()
}
func (_this *listEncoder) BeginConcatenate(ctx *EncoderContext) {
	_this.prepareToWrite(ctx)
	ctx.BeginStandardConcatenate()
}
func (_this *listEncoder) BeginConstant(ctx *EncoderContext, name []byte, explicitValue bool) {
	_this.prepareToWrite(ctx)
	ctx.BeginStandardConstant(name, explicitValue)
}
func (_this *listEncoder) BeginNA(ctx *EncoderContext) {
	_this.prepareToWrite(ctx)
	ctx.BeginStandardNA()
}
func (_this *listEncoder) EncodeArray(ctx *EncoderContext, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	_this.prepareToWrite(ctx)
	ctx.ArrayEngine.EncodeArray(arrayType, elementCount, data)
}
func (_this *listEncoder) EncodeStringlikeArray(ctx *EncoderContext, arrayType events.ArrayType, data string) {
	_this.prepareToWrite(ctx)
	ctx.ArrayEngine.EncodeStringlikeArray(arrayType, data)
}
func (_this *listEncoder) BeginArray(ctx *EncoderContext, arrayType events.ArrayType) {
	_this.prepareToWrite(ctx)
	ctx.BeginStandardArray(arrayType)
}

// =============================================================================

func encodeMapSeparator(ctx *EncoderContext) {
	ctx.Stream.AddString(" = ")
}

type mapKeyEncoder struct{}

var globalMapKeyEncoder mapKeyEncoder

func (_this *mapKeyEncoder) String() string { return "mapKeyEncoder" }

func (_this *mapKeyEncoder) prepareToWrite(ctx *EncoderContext) {
	ctx.WriteCurrentPrefix()
	ctx.SetStandardMapValuePrefix()
	ctx.ChangeEncoder(&globalMapValueEncoder)
	ctx.ContainerHasObjects = true
}

func (_this *mapKeyEncoder) prepareForContainer(ctx *EncoderContext) {
	ctx.WriteCurrentPrefix()
}

func (_this *mapKeyEncoder) Begin(ctx *EncoderContext) {
	ctx.Stream.WriteMapBegin()
	ctx.IncreaseIndent()
	ctx.SetStandardMapKeyPrefix()
	ctx.ContainerHasObjects = false
}

func (_this *mapKeyEncoder) End(ctx *EncoderContext) {
	ctx.DecreaseIndent()
	if ctx.ContainerHasObjects {
		ctx.WriteBasicIndent()
	}
	ctx.Stream.WriteMapEnd()
	ctx.Unstack()
	ctx.CurrentEncoder.ChildContainerFinished(ctx)
}

func (_this *mapKeyEncoder) ChildContainerFinished(ctx *EncoderContext) {
	ctx.SetStandardMapValuePrefix()
	ctx.ChangeEncoder(&globalMapValueEncoder)
	ctx.ContainerHasObjects = true
}

func (_this *mapKeyEncoder) EncodeBool(ctx *EncoderContext, value bool) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteBool(value)
}
func (_this *mapKeyEncoder) EncodeTrue(ctx *EncoderContext) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteTrue()
}
func (_this *mapKeyEncoder) EncodeFalse(ctx *EncoderContext) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteFalse()
}
func (_this *mapKeyEncoder) EncodePositiveInt(ctx *EncoderContext, value uint64) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WritePositiveInt(value)
}
func (_this *mapKeyEncoder) EncodeNegativeInt(ctx *EncoderContext, value uint64) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteNegativeInt(value)
}
func (_this *mapKeyEncoder) EncodeInt(ctx *EncoderContext, value int64) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteInt(value)
}
func (_this *mapKeyEncoder) EncodeBigInt(ctx *EncoderContext, value *big.Int) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteBigInt(value)
}
func (_this *mapKeyEncoder) EncodeFloat(ctx *EncoderContext, value float64) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteFloat(value)
}
func (_this *mapKeyEncoder) EncodeBigFloat(ctx *EncoderContext, value *big.Float) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteBigFloat(value)
}
func (_this *mapKeyEncoder) EncodeDecimalFloat(ctx *EncoderContext, value compact_float.DFloat) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteDecimalFloat(value)
}
func (_this *mapKeyEncoder) EncodeBigDecimalFloat(ctx *EncoderContext, value *apd.Decimal) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteBigDecimalFloat(value)
}
func (_this *mapKeyEncoder) EncodeTime(ctx *EncoderContext, value time.Time) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteTime(value)
}
func (_this *mapKeyEncoder) EncodeCompactTime(ctx *EncoderContext, value compact_time.Time) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteCompactTime(value)
}
func (_this *mapKeyEncoder) EncodeUUID(ctx *EncoderContext, value []byte) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteUUID(value)
}
func (_this *mapKeyEncoder) BeginMetadata(ctx *EncoderContext) {
	_this.prepareForContainer(ctx)
	ctx.BeginStandardMetadata()
}
func (_this *mapKeyEncoder) BeginComment(ctx *EncoderContext) {
	_this.prepareForContainer(ctx)
	ctx.BeginStandardComment()
}
func (_this *mapKeyEncoder) BeginMarker(ctx *EncoderContext) {
	_this.prepareForContainer(ctx)
	ctx.BeginStandardMarker()
}
func (_this *mapKeyEncoder) BeginReference(ctx *EncoderContext) {
	_this.prepareForContainer(ctx)
	ctx.BeginStandardReference()
}
func (_this *mapKeyEncoder) BeginConcatenate(ctx *EncoderContext) {
	// TODO ?
	_this.prepareForContainer(ctx)
	ctx.BeginStandardConcatenate()
}
func (_this *mapKeyEncoder) BeginConstant(ctx *EncoderContext, name []byte, explicitValue bool) {
	_this.prepareForContainer(ctx)
	ctx.BeginStandardConstant(name, explicitValue)
}
func (_this *mapKeyEncoder) EncodeArray(ctx *EncoderContext, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	_this.prepareToWrite(ctx)
	ctx.ArrayEngine.EncodeArray(arrayType, elementCount, data)
}
func (_this *mapKeyEncoder) EncodeStringlikeArray(ctx *EncoderContext, arrayType events.ArrayType, data string) {
	_this.prepareToWrite(ctx)
	ctx.ArrayEngine.EncodeStringlikeArray(arrayType, data)
}
func (_this *mapKeyEncoder) BeginArray(ctx *EncoderContext, arrayType events.ArrayType) {
	_this.prepareForContainer(ctx)
	ctx.BeginStandardArray(arrayType)
}

// =============================================================================

type mapValueEncoder struct{}

var globalMapValueEncoder mapValueEncoder

func (_this *mapValueEncoder) String() string { return "mapValueEncoder" }

func (_this *mapValueEncoder) prepareToWrite(ctx *EncoderContext) {
	ctx.WriteCurrentPrefix()
	ctx.SetStandardMapKeyPrefix()
	ctx.ChangeEncoder(&globalMapKeyEncoder)
	ctx.ContainerHasObjects = true
}

func (_this *mapValueEncoder) prepareForContainer(ctx *EncoderContext) {
	ctx.WriteCurrentPrefix()
	ctx.ContainerHasObjects = true
}

func (_this *mapValueEncoder) ChildContainerFinished(ctx *EncoderContext) {
	ctx.SetStandardMapKeyPrefix()
	ctx.ChangeEncoder(&globalMapKeyEncoder)
	ctx.ContainerHasObjects = true
}

func (_this *mapValueEncoder) EncodeBool(ctx *EncoderContext, value bool) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteBool(value)
}
func (_this *mapValueEncoder) EncodeTrue(ctx *EncoderContext) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteTrue()
}
func (_this *mapValueEncoder) EncodeFalse(ctx *EncoderContext) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteFalse()
}
func (_this *mapValueEncoder) EncodePositiveInt(ctx *EncoderContext, value uint64) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WritePositiveInt(value)
}
func (_this *mapValueEncoder) EncodeNegativeInt(ctx *EncoderContext, value uint64) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteNegativeInt(value)
}
func (_this *mapValueEncoder) EncodeInt(ctx *EncoderContext, value int64) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteInt(value)
}
func (_this *mapValueEncoder) EncodeBigInt(ctx *EncoderContext, value *big.Int) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteBigInt(value)
}
func (_this *mapValueEncoder) EncodeFloat(ctx *EncoderContext, value float64) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteFloat(value)
}
func (_this *mapValueEncoder) EncodeBigFloat(ctx *EncoderContext, value *big.Float) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteBigFloat(value)
}
func (_this *mapValueEncoder) EncodeDecimalFloat(ctx *EncoderContext, value compact_float.DFloat) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteDecimalFloat(value)
}
func (_this *mapValueEncoder) EncodeBigDecimalFloat(ctx *EncoderContext, value *apd.Decimal) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteBigDecimalFloat(value)
}
func (_this *mapValueEncoder) EncodeNan(ctx *EncoderContext, signaling bool) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteNan(signaling)
}
func (_this *mapValueEncoder) EncodeTime(ctx *EncoderContext, value time.Time) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteTime(value)
}
func (_this *mapValueEncoder) EncodeCompactTime(ctx *EncoderContext, value compact_time.Time) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteCompactTime(value)
}
func (_this *mapValueEncoder) EncodeUUID(ctx *EncoderContext, value []byte) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteUUID(value)
}
func (_this *mapValueEncoder) BeginList(ctx *EncoderContext) {
	_this.prepareForContainer(ctx)
	ctx.BeginStandardList()
}
func (_this *mapValueEncoder) BeginMap(ctx *EncoderContext) {
	_this.prepareForContainer(ctx)
	ctx.BeginStandardMap()
}
func (_this *mapValueEncoder) BeginMarkup(ctx *EncoderContext) {
	_this.prepareForContainer(ctx)
	ctx.BeginStandardMarkup()
}
func (_this *mapValueEncoder) BeginMetadata(ctx *EncoderContext) {
	_this.prepareForContainer(ctx)
	ctx.BeginStandardMetadata()
}
func (_this *mapValueEncoder) BeginComment(ctx *EncoderContext) {
	_this.prepareForContainer(ctx)
	ctx.ChangeEncoder(&globalMapKeyEncoder)
	ctx.BeginStandardComment()
}
func (_this *mapValueEncoder) BeginMarker(ctx *EncoderContext) {
	_this.prepareForContainer(ctx)
	ctx.BeginStandardMarker()
}
func (_this *mapValueEncoder) BeginReference(ctx *EncoderContext) {
	_this.prepareForContainer(ctx)
	ctx.BeginStandardReference()
}
func (_this *mapValueEncoder) BeginConcatenate(ctx *EncoderContext) {
	_this.prepareForContainer(ctx)
	ctx.BeginStandardConcatenate()
}
func (_this *mapValueEncoder) BeginConstant(ctx *EncoderContext, name []byte, explicitValue bool) {
	_this.prepareForContainer(ctx)
	ctx.BeginStandardConstant(name, explicitValue)
}
func (_this *mapValueEncoder) BeginNA(ctx *EncoderContext) {
	_this.prepareForContainer(ctx)
	ctx.BeginStandardNA()
}
func (_this *mapValueEncoder) EncodeArray(ctx *EncoderContext, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	_this.prepareToWrite(ctx)
	ctx.ArrayEngine.EncodeArray(arrayType, elementCount, data)
}
func (_this *mapValueEncoder) EncodeStringlikeArray(ctx *EncoderContext, arrayType events.ArrayType, data string) {
	_this.prepareToWrite(ctx)
	ctx.ArrayEngine.EncodeStringlikeArray(arrayType, data)
}
func (_this *mapValueEncoder) BeginArray(ctx *EncoderContext, arrayType events.ArrayType) {
	_this.prepareForContainer(ctx)
	ctx.BeginStandardArray(arrayType)
}

// =============================================================================

type metadataKeyEncoder struct{}

var globalMetadataKeyEncoder metadataKeyEncoder

func (_this *metadataKeyEncoder) String() string { return "metadataKeyEncoder" }

func (_this *metadataKeyEncoder) prepareToWrite(ctx *EncoderContext) {
	ctx.WriteCurrentPrefix()
	ctx.SetStandardMapValuePrefix()
	ctx.ChangeEncoder(&globalMetadataValueEncoder)
	ctx.ContainerHasObjects = true
}

func (_this *metadataKeyEncoder) prepareForContainer(ctx *EncoderContext) {
	ctx.WriteCurrentPrefix()
	ctx.ContainerHasObjects = true
}

func (_this *metadataKeyEncoder) Begin(ctx *EncoderContext) {
	ctx.Stream.WriteMetadataBegin()
	ctx.IncreaseIndent()
	ctx.SetStandardMapKeyPrefix()
	ctx.ContainerHasObjects = false
}

func (_this *metadataKeyEncoder) End(ctx *EncoderContext) {
	ctx.DecreaseIndent()
	if ctx.ContainerHasObjects {
		ctx.WriteBasicIndent()
	}
	ctx.Stream.WriteMetadataEnd()
	ctx.ChangeEncoder(&globalPostInvisibleEncoder)
}

func (_this *metadataKeyEncoder) ChildContainerFinished(ctx *EncoderContext) {
	ctx.SetStandardMapValuePrefix()
	ctx.ChangeEncoder(&globalMetadataValueEncoder)
	ctx.ContainerHasObjects = true
}

func (_this *metadataKeyEncoder) EncodeBool(ctx *EncoderContext, value bool) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteBool(value)
}
func (_this *metadataKeyEncoder) EncodeTrue(ctx *EncoderContext) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteTrue()
}
func (_this *metadataKeyEncoder) EncodeFalse(ctx *EncoderContext) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteFalse()
}
func (_this *metadataKeyEncoder) EncodePositiveInt(ctx *EncoderContext, value uint64) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WritePositiveInt(value)
}
func (_this *metadataKeyEncoder) EncodeNegativeInt(ctx *EncoderContext, value uint64) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteNegativeInt(value)
}
func (_this *metadataKeyEncoder) EncodeInt(ctx *EncoderContext, value int64) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteInt(value)
}
func (_this *metadataKeyEncoder) EncodeBigInt(ctx *EncoderContext, value *big.Int) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteBigInt(value)
}
func (_this *metadataKeyEncoder) EncodeFloat(ctx *EncoderContext, value float64) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteFloat(value)
}
func (_this *metadataKeyEncoder) EncodeBigFloat(ctx *EncoderContext, value *big.Float) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteBigFloat(value)
}
func (_this *metadataKeyEncoder) EncodeDecimalFloat(ctx *EncoderContext, value compact_float.DFloat) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteDecimalFloat(value)
}
func (_this *metadataKeyEncoder) EncodeBigDecimalFloat(ctx *EncoderContext, value *apd.Decimal) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteBigDecimalFloat(value)
}
func (_this *metadataKeyEncoder) EncodeTime(ctx *EncoderContext, value time.Time) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteTime(value)
}
func (_this *metadataKeyEncoder) EncodeCompactTime(ctx *EncoderContext, value compact_time.Time) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteCompactTime(value)
}
func (_this *metadataKeyEncoder) EncodeUUID(ctx *EncoderContext, value []byte) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteUUID(value)
}
func (_this *metadataKeyEncoder) BeginMetadata(ctx *EncoderContext) {
	_this.prepareForContainer(ctx)
	ctx.BeginStandardMetadata()
}
func (_this *metadataKeyEncoder) BeginComment(ctx *EncoderContext) {
	_this.prepareForContainer(ctx)
	ctx.BeginStandardComment()
}
func (_this *metadataKeyEncoder) BeginMarker(ctx *EncoderContext) {
	_this.prepareForContainer(ctx)
	ctx.BeginStandardMarker()
}
func (_this *metadataKeyEncoder) BeginReference(ctx *EncoderContext) {
	_this.prepareForContainer(ctx)
	ctx.BeginStandardReference()
}
func (_this *metadataKeyEncoder) BeginConcatenate(ctx *EncoderContext) {
	// TODO ?
	_this.prepareForContainer(ctx)
	ctx.BeginStandardConcatenate()
}
func (_this *metadataKeyEncoder) BeginConstant(ctx *EncoderContext, name []byte, explicitValue bool) {
	_this.prepareForContainer(ctx)
	ctx.BeginStandardConstant(name, explicitValue)
}
func (_this *metadataKeyEncoder) EncodeArray(ctx *EncoderContext, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	_this.prepareToWrite(ctx)
	ctx.ArrayEngine.EncodeArray(arrayType, elementCount, data)
}
func (_this *metadataKeyEncoder) EncodeStringlikeArray(ctx *EncoderContext, arrayType events.ArrayType, data string) {
	_this.prepareToWrite(ctx)
	ctx.ArrayEngine.EncodeStringlikeArray(arrayType, data)
}
func (_this *metadataKeyEncoder) BeginArray(ctx *EncoderContext, arrayType events.ArrayType) {
	_this.prepareForContainer(ctx)
	ctx.BeginStandardArray(arrayType)
}

// =============================================================================

type metadataValueEncoder struct{}

var globalMetadataValueEncoder metadataValueEncoder

func (_this *metadataValueEncoder) String() string { return "metadataValueEncoder" }

func (_this *metadataValueEncoder) prepareToWrite(ctx *EncoderContext) {
	ctx.WriteCurrentPrefix()
	ctx.SetStandardMapKeyPrefix()
	ctx.ChangeEncoder(&globalMetadataKeyEncoder)
	ctx.ContainerHasObjects = true
}

func (_this *metadataValueEncoder) prepareForContainer(ctx *EncoderContext) {
	ctx.WriteCurrentPrefix()
	ctx.ContainerHasObjects = true
}

func (_this *metadataValueEncoder) ChildContainerFinished(ctx *EncoderContext) {
	ctx.SetStandardMapKeyPrefix()
	ctx.ChangeEncoder(&globalMetadataKeyEncoder)
	ctx.ContainerHasObjects = true
}

func (_this *metadataValueEncoder) EncodeBool(ctx *EncoderContext, value bool) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteBool(value)
}
func (_this *metadataValueEncoder) EncodeTrue(ctx *EncoderContext) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteTrue()
}
func (_this *metadataValueEncoder) EncodeFalse(ctx *EncoderContext) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteFalse()
}
func (_this *metadataValueEncoder) EncodePositiveInt(ctx *EncoderContext, value uint64) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WritePositiveInt(value)
}
func (_this *metadataValueEncoder) EncodeNegativeInt(ctx *EncoderContext, value uint64) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteNegativeInt(value)
}
func (_this *metadataValueEncoder) EncodeInt(ctx *EncoderContext, value int64) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteInt(value)
}
func (_this *metadataValueEncoder) EncodeBigInt(ctx *EncoderContext, value *big.Int) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteBigInt(value)
}
func (_this *metadataValueEncoder) EncodeFloat(ctx *EncoderContext, value float64) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteFloat(value)
}
func (_this *metadataValueEncoder) EncodeBigFloat(ctx *EncoderContext, value *big.Float) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteBigFloat(value)
}
func (_this *metadataValueEncoder) EncodeDecimalFloat(ctx *EncoderContext, value compact_float.DFloat) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteDecimalFloat(value)
}
func (_this *metadataValueEncoder) EncodeBigDecimalFloat(ctx *EncoderContext, value *apd.Decimal) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteBigDecimalFloat(value)
}
func (_this *metadataValueEncoder) EncodeNan(ctx *EncoderContext, signaling bool) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteNan(signaling)
}
func (_this *metadataValueEncoder) EncodeTime(ctx *EncoderContext, value time.Time) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteTime(value)
}
func (_this *metadataValueEncoder) EncodeCompactTime(ctx *EncoderContext, value compact_time.Time) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteCompactTime(value)
}
func (_this *metadataValueEncoder) EncodeUUID(ctx *EncoderContext, value []byte) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteUUID(value)
}
func (_this *metadataValueEncoder) BeginList(ctx *EncoderContext) {
	_this.prepareForContainer(ctx)
	ctx.BeginStandardList()
}
func (_this *metadataValueEncoder) BeginMap(ctx *EncoderContext) {
	_this.prepareForContainer(ctx)
	ctx.BeginStandardMap()
}
func (_this *metadataValueEncoder) BeginMarkup(ctx *EncoderContext) {
	_this.prepareForContainer(ctx)
	ctx.BeginStandardMarkup()
}
func (_this *metadataValueEncoder) BeginMetadata(ctx *EncoderContext) {
	_this.prepareForContainer(ctx)
	ctx.BeginStandardMetadata()
}
func (_this *metadataValueEncoder) BeginComment(ctx *EncoderContext) {
	_this.prepareForContainer(ctx)
	ctx.ChangeEncoder(&globalMetadataKeyEncoder)
	ctx.BeginStandardComment()
}
func (_this *metadataValueEncoder) BeginMarker(ctx *EncoderContext) {
	_this.prepareForContainer(ctx)
	ctx.BeginStandardMarker()
}
func (_this *metadataValueEncoder) BeginReference(ctx *EncoderContext) {
	_this.prepareForContainer(ctx)
	ctx.BeginStandardReference()
}
func (_this *metadataValueEncoder) BeginConcatenate(ctx *EncoderContext) {
	_this.prepareForContainer(ctx)
	ctx.BeginStandardConcatenate()
}
func (_this *metadataValueEncoder) BeginConstant(ctx *EncoderContext, name []byte, explicitValue bool) {
	_this.prepareForContainer(ctx)
	ctx.BeginStandardConstant(name, explicitValue)
}
func (_this *metadataValueEncoder) BeginNA(ctx *EncoderContext) {
	_this.prepareForContainer(ctx)
	ctx.BeginStandardNA()
}
func (_this *metadataValueEncoder) EncodeArray(ctx *EncoderContext, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	_this.prepareToWrite(ctx)
	ctx.ArrayEngine.EncodeArray(arrayType, elementCount, data)
}
func (_this *metadataValueEncoder) EncodeStringlikeArray(ctx *EncoderContext, arrayType events.ArrayType, data string) {
	_this.prepareToWrite(ctx)
	ctx.ArrayEngine.EncodeStringlikeArray(arrayType, data)
}
func (_this *metadataValueEncoder) BeginArray(ctx *EncoderContext, arrayType events.ArrayType) {
	_this.prepareForContainer(ctx)
	ctx.BeginStandardArray(arrayType)
}

// =============================================================================

type markupNameEncoder struct{}

var globalMarkupNameEncoder markupNameEncoder

func (_this *markupNameEncoder) String() string { return "markupNameEncoder" }

func (_this *markupNameEncoder) prepareToWrite(ctx *EncoderContext) {
	ctx.WriteCurrentPrefix()
	ctx.SetMarkupAttributeKeyPrefix()
	ctx.SwitchToMarkupAttributes()
}

func (_this *markupNameEncoder) prepareForContainer(ctx *EncoderContext) {
	ctx.WriteCurrentPrefix()
	ctx.SetMarkupAttributeKeyPrefix()
}

func (_this *markupNameEncoder) Begin(ctx *EncoderContext) {
	ctx.Stream.WriteMarkupBegin()
	ctx.IncreaseIndent()
	ctx.ClearPrefix()
	ctx.ContainerHasObjects = false
}

func (_this *markupNameEncoder) ChildContainerFinished(ctx *EncoderContext) {
	ctx.SwitchToMarkupAttributes()
}

func (_this *markupNameEncoder) EncodeBool(ctx *EncoderContext, value bool) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteBool(value)
}
func (_this *markupNameEncoder) EncodeTrue(ctx *EncoderContext) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteTrue()
}
func (_this *markupNameEncoder) EncodeFalse(ctx *EncoderContext) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteFalse()
}
func (_this *markupNameEncoder) EncodePositiveInt(ctx *EncoderContext, value uint64) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WritePositiveInt(value)
}
func (_this *markupNameEncoder) EncodeNegativeInt(ctx *EncoderContext, value uint64) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteNegativeInt(value)
}
func (_this *markupNameEncoder) EncodeInt(ctx *EncoderContext, value int64) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteInt(value)
}
func (_this *markupNameEncoder) EncodeBigInt(ctx *EncoderContext, value *big.Int) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteBigInt(value)
}
func (_this *markupNameEncoder) EncodeFloat(ctx *EncoderContext, value float64) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteFloat(value)
}
func (_this *markupNameEncoder) EncodeBigFloat(ctx *EncoderContext, value *big.Float) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteBigFloat(value)
}
func (_this *markupNameEncoder) EncodeDecimalFloat(ctx *EncoderContext, value compact_float.DFloat) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteDecimalFloat(value)
}
func (_this *markupNameEncoder) EncodeBigDecimalFloat(ctx *EncoderContext, value *apd.Decimal) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteBigDecimalFloat(value)
}
func (_this *markupNameEncoder) EncodeTime(ctx *EncoderContext, value time.Time) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteTime(value)
}
func (_this *markupNameEncoder) EncodeCompactTime(ctx *EncoderContext, value compact_time.Time) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteCompactTime(value)
}
func (_this *markupNameEncoder) EncodeUUID(ctx *EncoderContext, value []byte) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteUUID(value)
}
func (_this *markupNameEncoder) BeginMarker(ctx *EncoderContext) {
	_this.prepareForContainer(ctx)
	ctx.BeginStandardMarker()
}
func (_this *markupNameEncoder) BeginReference(ctx *EncoderContext) {
	_this.prepareForContainer(ctx)
	ctx.BeginStandardReference()
}
func (_this *markupNameEncoder) BeginConcatenate(ctx *EncoderContext) {
	// TODO ?
	_this.prepareForContainer(ctx)
	ctx.BeginStandardConcatenate()
}
func (_this *markupNameEncoder) BeginConstant(ctx *EncoderContext, name []byte, explicitValue bool) {
	_this.prepareForContainer(ctx)
	ctx.BeginStandardConstant(name, explicitValue)
}
func (_this *markupNameEncoder) EncodeArray(ctx *EncoderContext, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	_this.prepareToWrite(ctx)
	ctx.ArrayEngine.EncodeArray(arrayType, elementCount, data)
}
func (_this *markupNameEncoder) EncodeStringlikeArray(ctx *EncoderContext, arrayType events.ArrayType, data string) {
	_this.prepareToWrite(ctx)
	ctx.ArrayEngine.EncodeStringlikeArray(arrayType, data)
}
func (_this *markupNameEncoder) BeginArray(ctx *EncoderContext, arrayType events.ArrayType) {
	_this.prepareForContainer(ctx)
	ctx.BeginStandardArray(arrayType)
}

// =============================================================================

type markupKeyEncoder struct{}

var globalMarkupKeyEncoder markupKeyEncoder

func (_this *markupKeyEncoder) String() string { return "markupKeyEncoder" }

func (_this *markupKeyEncoder) prepareToWrite(ctx *EncoderContext) {
	ctx.WriteCurrentPrefix()
	ctx.SetMarkupAttributeValuePrefix()
	ctx.ChangeEncoder(&globalMarkupValueEncoder)
}

func (_this *markupKeyEncoder) prepareForContainer(ctx *EncoderContext) {
	ctx.WriteCurrentPrefix()
}

func (_this *markupKeyEncoder) End(ctx *EncoderContext) {
	if ctx.ContainerHasObjects {
		ctx.Stream.WriteMarkupContentsBegin()
	}
	ctx.SwitchToMarkupContents()
}

func (_this *markupKeyEncoder) ChildContainerFinished(ctx *EncoderContext) {
	ctx.SetMarkupAttributeValuePrefix()
	ctx.ChangeEncoder(&globalMarkupValueEncoder)
}

func (_this *markupKeyEncoder) EncodeBool(ctx *EncoderContext, value bool) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteBool(value)
}
func (_this *markupKeyEncoder) EncodeTrue(ctx *EncoderContext) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteTrue()
}
func (_this *markupKeyEncoder) EncodeFalse(ctx *EncoderContext) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteFalse()
}
func (_this *markupKeyEncoder) EncodePositiveInt(ctx *EncoderContext, value uint64) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WritePositiveInt(value)
}
func (_this *markupKeyEncoder) EncodeNegativeInt(ctx *EncoderContext, value uint64) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteNegativeInt(value)
}
func (_this *markupKeyEncoder) EncodeInt(ctx *EncoderContext, value int64) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteInt(value)
}
func (_this *markupKeyEncoder) EncodeBigInt(ctx *EncoderContext, value *big.Int) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteBigInt(value)
}
func (_this *markupKeyEncoder) EncodeFloat(ctx *EncoderContext, value float64) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteFloat(value)
}
func (_this *markupKeyEncoder) EncodeBigFloat(ctx *EncoderContext, value *big.Float) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteBigFloat(value)
}
func (_this *markupKeyEncoder) EncodeDecimalFloat(ctx *EncoderContext, value compact_float.DFloat) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteDecimalFloat(value)
}
func (_this *markupKeyEncoder) EncodeBigDecimalFloat(ctx *EncoderContext, value *apd.Decimal) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteBigDecimalFloat(value)
}
func (_this *markupKeyEncoder) EncodeTime(ctx *EncoderContext, value time.Time) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteTime(value)
}
func (_this *markupKeyEncoder) EncodeCompactTime(ctx *EncoderContext, value compact_time.Time) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteCompactTime(value)
}
func (_this *markupKeyEncoder) EncodeUUID(ctx *EncoderContext, value []byte) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteUUID(value)
}
func (_this *markupKeyEncoder) BeginMetadata(ctx *EncoderContext) {
	_this.prepareForContainer(ctx)
	ctx.BeginStandardMetadata()
}
func (_this *markupKeyEncoder) BeginComment(ctx *EncoderContext) {
	_this.prepareForContainer(ctx)
	ctx.BeginStandardComment()
}
func (_this *markupKeyEncoder) BeginMarker(ctx *EncoderContext) {
	_this.prepareForContainer(ctx)
	ctx.BeginStandardMarker()
}
func (_this *markupKeyEncoder) BeginReference(ctx *EncoderContext) {
	_this.prepareForContainer(ctx)
	ctx.BeginStandardReference()
}
func (_this *markupKeyEncoder) BeginConcatenate(ctx *EncoderContext) {
	// TODO ?
	_this.prepareForContainer(ctx)
	ctx.BeginStandardConcatenate()
}
func (_this *markupKeyEncoder) BeginConstant(ctx *EncoderContext, name []byte, explicitValue bool) {
	_this.prepareForContainer(ctx)
	ctx.BeginStandardConstant(name, explicitValue)
}
func (_this *markupKeyEncoder) EncodeArray(ctx *EncoderContext, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	_this.prepareToWrite(ctx)
	ctx.ArrayEngine.EncodeArray(arrayType, elementCount, data)
}
func (_this *markupKeyEncoder) EncodeStringlikeArray(ctx *EncoderContext, arrayType events.ArrayType, data string) {
	_this.prepareToWrite(ctx)
	ctx.ArrayEngine.EncodeStringlikeArray(arrayType, data)
}
func (_this *markupKeyEncoder) BeginArray(ctx *EncoderContext, arrayType events.ArrayType) {
	_this.prepareForContainer(ctx)
	ctx.BeginStandardArray(arrayType)
}

// =============================================================================

type markupValueEncoder struct{}

var globalMarkupValueEncoder markupValueEncoder

func (_this *markupValueEncoder) String() string { return "markupValueEncoder" }

func (_this *markupValueEncoder) prepareToWrite(ctx *EncoderContext) {
	ctx.WriteCurrentPrefix()
	ctx.SetMarkupAttributeKeyPrefix()
	ctx.ChangeEncoder(&globalMarkupKeyEncoder)
}

func (_this *markupValueEncoder) prepareForContainer(ctx *EncoderContext) {
	ctx.WriteCurrentPrefix()
}

func (_this *markupValueEncoder) ChildContainerFinished(ctx *EncoderContext) {
	ctx.SetMarkupAttributeKeyPrefix()
	ctx.ChangeEncoder(&globalMarkupKeyEncoder)
}

func (_this *markupValueEncoder) EncodeBool(ctx *EncoderContext, value bool) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteBool(value)
}
func (_this *markupValueEncoder) EncodeTrue(ctx *EncoderContext) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteTrue()
}
func (_this *markupValueEncoder) EncodeFalse(ctx *EncoderContext) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteFalse()
}
func (_this *markupValueEncoder) EncodePositiveInt(ctx *EncoderContext, value uint64) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WritePositiveInt(value)
}
func (_this *markupValueEncoder) EncodeNegativeInt(ctx *EncoderContext, value uint64) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteNegativeInt(value)
}
func (_this *markupValueEncoder) EncodeInt(ctx *EncoderContext, value int64) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteInt(value)
}
func (_this *markupValueEncoder) EncodeBigInt(ctx *EncoderContext, value *big.Int) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteBigInt(value)
}
func (_this *markupValueEncoder) EncodeFloat(ctx *EncoderContext, value float64) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteFloat(value)
}
func (_this *markupValueEncoder) EncodeBigFloat(ctx *EncoderContext, value *big.Float) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteBigFloat(value)
}
func (_this *markupValueEncoder) EncodeDecimalFloat(ctx *EncoderContext, value compact_float.DFloat) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteDecimalFloat(value)
}
func (_this *markupValueEncoder) EncodeBigDecimalFloat(ctx *EncoderContext, value *apd.Decimal) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteBigDecimalFloat(value)
}
func (_this *markupValueEncoder) EncodeNan(ctx *EncoderContext, signaling bool) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteNan(signaling)
}
func (_this *markupValueEncoder) EncodeTime(ctx *EncoderContext, value time.Time) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteTime(value)
}
func (_this *markupValueEncoder) EncodeCompactTime(ctx *EncoderContext, value compact_time.Time) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteCompactTime(value)
}
func (_this *markupValueEncoder) EncodeUUID(ctx *EncoderContext, value []byte) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteUUID(value)
}
func (_this *markupValueEncoder) BeginList(ctx *EncoderContext) {
	_this.prepareForContainer(ctx)
	ctx.BeginStandardList()
}
func (_this *markupValueEncoder) BeginMap(ctx *EncoderContext) {
	_this.prepareForContainer(ctx)
	ctx.BeginStandardMap()
}
func (_this *markupValueEncoder) BeginMarkup(ctx *EncoderContext) {
	_this.prepareForContainer(ctx)
	ctx.BeginStandardMarkup()
}
func (_this *markupValueEncoder) BeginMetadata(ctx *EncoderContext) {
	_this.prepareForContainer(ctx)
	ctx.BeginStandardMetadata()
}
func (_this *markupValueEncoder) BeginComment(ctx *EncoderContext) {
	_this.prepareForContainer(ctx)
	ctx.ChangeEncoder(&globalMarkupKeyEncoder)
	ctx.BeginStandardComment()
}
func (_this *markupValueEncoder) BeginMarker(ctx *EncoderContext) {
	_this.prepareForContainer(ctx)
	ctx.BeginStandardMarker()
}
func (_this *markupValueEncoder) BeginReference(ctx *EncoderContext) {
	_this.prepareForContainer(ctx)
	ctx.BeginStandardReference()
}
func (_this *markupValueEncoder) BeginConcatenate(ctx *EncoderContext) {
	_this.prepareForContainer(ctx)
	ctx.BeginStandardConcatenate()
}
func (_this *markupValueEncoder) BeginConstant(ctx *EncoderContext, name []byte, explicitValue bool) {
	_this.prepareForContainer(ctx)
	ctx.BeginStandardConstant(name, explicitValue)
}
func (_this *markupValueEncoder) BeginNA(ctx *EncoderContext) {
	_this.prepareForContainer(ctx)
	ctx.BeginStandardNA()
}
func (_this *markupValueEncoder) EncodeArray(ctx *EncoderContext, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	_this.prepareToWrite(ctx)
	ctx.ArrayEngine.EncodeArray(arrayType, elementCount, data)
}
func (_this *markupValueEncoder) EncodeStringlikeArray(ctx *EncoderContext, arrayType events.ArrayType, data string) {
	_this.prepareToWrite(ctx)
	ctx.ArrayEngine.EncodeStringlikeArray(arrayType, data)
}
func (_this *markupValueEncoder) BeginArray(ctx *EncoderContext, arrayType events.ArrayType) {
	_this.prepareForContainer(ctx)
	ctx.BeginStandardArray(arrayType)
}

// =============================================================================

type markupContentsEncoder struct{}

var globalMarkupContentsEncoder markupContentsEncoder

func (_this *markupContentsEncoder) String() string { return "markupContentsEncoder" }

func (_this *markupContentsEncoder) prepareToWrite(ctx *EncoderContext) {
	if !ctx.ContainerHasObjects {
		ctx.Stream.AddByte(',')
		ctx.Stream.AddBytes(ctx.indenter.Get())
	}
	ctx.ContainerHasObjects = true
}

func (_this *markupContentsEncoder) Begin(ctx *EncoderContext) {
	ctx.ContainerHasObjects = false
}

func (_this *markupContentsEncoder) End(ctx *EncoderContext) {
	ctx.DecreaseIndent()
	if ctx.ContainerHasObjects {
		ctx.WriteBasicIndent()
	}
	ctx.Stream.WriteMarkupEnd()
	ctx.Unstack()
	ctx.CurrentEncoder.ChildContainerFinished(ctx)
}

func (_this *markupContentsEncoder) ChildContainerFinished(ctx *EncoderContext) {
	ctx.ContainerHasObjects = true
	ctx.SetStandardIndentPrefix()
}

func (_this *markupContentsEncoder) BeginMarkup(ctx *EncoderContext) {
	_this.prepareToWrite(ctx)
	ctx.BeginStandardMarkup()
}
func (_this *markupContentsEncoder) BeginComment(ctx *EncoderContext) {
	_this.prepareToWrite(ctx)
	ctx.BeginStandardComment()
}
func (_this *markupContentsEncoder) EncodeArray(ctx *EncoderContext, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	_this.prepareToWrite(ctx)
	ctx.ArrayEngine.EncodeMarkupContentStringData(data)
}
func (_this *markupContentsEncoder) EncodeStringlikeArray(ctx *EncoderContext, arrayType events.ArrayType, data string) {
	_this.prepareToWrite(ctx)
	ctx.ArrayEngine.EncodeMarkupContentString(data)
}
func (_this *markupContentsEncoder) BeginArray(ctx *EncoderContext, arrayType events.ArrayType) {
	_this.prepareToWrite(ctx)
	// TODO: begin content string array
	ctx.BeginStandardArray(arrayType)
}
