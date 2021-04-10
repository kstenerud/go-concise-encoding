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

func (_this *listEncoder) completeObject(ctx *EncoderContext) {
	ctx.SetStandardIndentPrefix()
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
	ctx.CurrentEncoder.ChildContainerFinished(ctx, true)
}

func (_this *listEncoder) ChildContainerFinished(ctx *EncoderContext, isVisibleChild bool) {
	ctx.ContainerHasObjects = true
	_this.completeObject(ctx)
}

func (_this *listEncoder) BeginNA(ctx *EncoderContext) {
	_this.prepareToWrite(ctx)
	ctx.BeginNA()
}
func (_this *listEncoder) EncodeNil(ctx *EncoderContext) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteNil()
	_this.completeObject(ctx)
}
func (_this *listEncoder) EncodeBool(ctx *EncoderContext, value bool) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteBool(value)
	_this.completeObject(ctx)
}
func (_this *listEncoder) EncodeTrue(ctx *EncoderContext) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteTrue()
	_this.completeObject(ctx)
}
func (_this *listEncoder) EncodeFalse(ctx *EncoderContext) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteFalse()
	_this.completeObject(ctx)
}
func (_this *listEncoder) EncodePositiveInt(ctx *EncoderContext, value uint64) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WritePositiveInt(value)
	_this.completeObject(ctx)
}
func (_this *listEncoder) EncodeNegativeInt(ctx *EncoderContext, value uint64) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteNegativeInt(value)
	_this.completeObject(ctx)
}
func (_this *listEncoder) EncodeInt(ctx *EncoderContext, value int64) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteInt(value)
	_this.completeObject(ctx)
}
func (_this *listEncoder) EncodeBigInt(ctx *EncoderContext, value *big.Int) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteBigInt(value)
	_this.completeObject(ctx)
}
func (_this *listEncoder) EncodeFloat(ctx *EncoderContext, value float64) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteFloat(value)
	_this.completeObject(ctx)
}
func (_this *listEncoder) EncodeBigFloat(ctx *EncoderContext, value *big.Float) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteBigFloat(value)
	_this.completeObject(ctx)
}
func (_this *listEncoder) EncodeDecimalFloat(ctx *EncoderContext, value compact_float.DFloat) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteDecimalFloat(value)
	_this.completeObject(ctx)
}
func (_this *listEncoder) EncodeBigDecimalFloat(ctx *EncoderContext, value *apd.Decimal) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteBigDecimalFloat(value)
	_this.completeObject(ctx)
}
func (_this *listEncoder) EncodeNan(ctx *EncoderContext, signaling bool) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteNan(signaling)
	_this.completeObject(ctx)
}
func (_this *listEncoder) EncodeTime(ctx *EncoderContext, value time.Time) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteTime(value)
	_this.completeObject(ctx)
}
func (_this *listEncoder) EncodeCompactTime(ctx *EncoderContext, value compact_time.Time) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteCompactTime(value)
	_this.completeObject(ctx)
}
func (_this *listEncoder) EncodeUUID(ctx *EncoderContext, value []byte) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteUUID(value)
	_this.completeObject(ctx)
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
func (_this *listEncoder) BeginConstant(ctx *EncoderContext, name []byte, explicitValue bool) {
	_this.prepareToWrite(ctx)
	ctx.BeginStandardConstant(name, explicitValue)
}
func (_this *listEncoder) EncodeArray(ctx *EncoderContext, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	_this.prepareToWrite(ctx)
	ctx.WriteArray(arrayType, elementCount, data)
	_this.completeObject(ctx)
}
func (_this *listEncoder) EncodeStringlikeArray(ctx *EncoderContext, arrayType events.ArrayType, data string) {
	_this.prepareToWrite(ctx)
	ctx.WriteStringlikeArray(arrayType, data)
	_this.completeObject(ctx)
}
func (_this *listEncoder) BeginArray(ctx *EncoderContext, arrayType events.ArrayType) {
	_this.prepareToWrite(ctx)
	ctx.BeginStandardArray(arrayType)
}

// =============================================================================

func encodeMapSeparator(ctx *EncoderContext) {
	ctx.Stream.WriteString(" = ")
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
	ctx.CurrentEncoder.ChildContainerFinished(ctx, true)
}

func (_this *mapKeyEncoder) ChildContainerFinished(ctx *EncoderContext, isVisibleChild bool) {
	if isVisibleChild {
		ctx.SetStandardMapValuePrefix()
		ctx.ChangeEncoder(&globalMapValueEncoder)
	}
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
func (_this *mapKeyEncoder) BeginConstant(ctx *EncoderContext, name []byte, explicitValue bool) {
	_this.prepareForContainer(ctx)
	ctx.BeginStandardConstant(name, explicitValue)
}
func (_this *mapKeyEncoder) EncodeArray(ctx *EncoderContext, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	_this.prepareToWrite(ctx)
	ctx.WriteArray(arrayType, elementCount, data)
}
func (_this *mapKeyEncoder) EncodeStringlikeArray(ctx *EncoderContext, arrayType events.ArrayType, data string) {
	_this.prepareToWrite(ctx)
	ctx.WriteStringlikeArray(arrayType, data)
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

func (_this *mapValueEncoder) ChildContainerFinished(ctx *EncoderContext, isVisibleChild bool) {
	if isVisibleChild {
		ctx.SetStandardMapKeyPrefix()
		ctx.ChangeEncoder(&globalMapKeyEncoder)
	}
	ctx.ContainerHasObjects = true
}

func (_this *mapValueEncoder) BeginNA(ctx *EncoderContext) {
	_this.prepareForContainer(ctx)
	ctx.BeginNA()
}
func (_this *mapValueEncoder) EncodeNil(ctx *EncoderContext) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteNil()
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
func (_this *mapValueEncoder) BeginComment(ctx *EncoderContext) {
	_this.prepareForContainer(ctx)
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
func (_this *mapValueEncoder) BeginConstant(ctx *EncoderContext, name []byte, explicitValue bool) {
	_this.prepareForContainer(ctx)
	ctx.BeginStandardConstant(name, explicitValue)
}
func (_this *mapValueEncoder) EncodeArray(ctx *EncoderContext, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	_this.prepareToWrite(ctx)
	ctx.WriteArray(arrayType, elementCount, data)
}
func (_this *mapValueEncoder) EncodeStringlikeArray(ctx *EncoderContext, arrayType events.ArrayType, data string) {
	_this.prepareToWrite(ctx)
	ctx.WriteStringlikeArray(arrayType, data)
}
func (_this *mapValueEncoder) BeginArray(ctx *EncoderContext, arrayType events.ArrayType) {
	_this.prepareForContainer(ctx)
	ctx.BeginStandardArray(arrayType)
}

// =============================================================================

type markupNameEncoder struct{}

var globalMarkupNameEncoder markupNameEncoder

func (_this *markupNameEncoder) String() string { return "markupNameEncoder" }

func (_this *markupNameEncoder) Begin(ctx *EncoderContext) {
	ctx.Stream.WriteMarkupBegin()
	ctx.IncreaseIndent()
	ctx.ClearPrefix()
	ctx.ContainerHasObjects = false
}

func (_this *markupNameEncoder) EncodeIdentifier(ctx *EncoderContext, data []uint8) {
	ctx.WriteCurrentPrefix()
	ctx.SetMarkupAttributeKeyPrefix()
	ctx.SwitchToMarkupAttributes()
	ctx.WriteIdentifier(data)
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

func (_this *markupKeyEncoder) ChildContainerFinished(ctx *EncoderContext, isVisibleChild bool) {
	if isVisibleChild {
		ctx.SetMarkupAttributeValuePrefix()
		ctx.ChangeEncoder(&globalMarkupValueEncoder)
	}
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
func (_this *markupKeyEncoder) BeginConstant(ctx *EncoderContext, name []byte, explicitValue bool) {
	_this.prepareForContainer(ctx)
	ctx.BeginStandardConstant(name, explicitValue)
}
func (_this *markupKeyEncoder) EncodeArray(ctx *EncoderContext, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	_this.prepareToWrite(ctx)
	ctx.WriteArray(arrayType, elementCount, data)
}
func (_this *markupKeyEncoder) EncodeStringlikeArray(ctx *EncoderContext, arrayType events.ArrayType, data string) {
	_this.prepareToWrite(ctx)
	ctx.WriteStringlikeArray(arrayType, data)
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

func (_this *markupValueEncoder) ChildContainerFinished(ctx *EncoderContext, isVisibleChild bool) {
	if isVisibleChild {
		ctx.SetMarkupAttributeKeyPrefix()
		ctx.ChangeEncoder(&globalMarkupKeyEncoder)
	}
}

func (_this *markupValueEncoder) BeginNA(ctx *EncoderContext) {
	_this.prepareForContainer(ctx)
	ctx.BeginNA()
}
func (_this *markupValueEncoder) EncodeNil(ctx *EncoderContext) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteNil()
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
func (_this *markupValueEncoder) BeginComment(ctx *EncoderContext) {
	_this.prepareForContainer(ctx)
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
func (_this *markupValueEncoder) BeginConstant(ctx *EncoderContext, name []byte, explicitValue bool) {
	_this.prepareForContainer(ctx)
	ctx.BeginStandardConstant(name, explicitValue)
}
func (_this *markupValueEncoder) EncodeArray(ctx *EncoderContext, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	_this.prepareToWrite(ctx)
	ctx.WriteArray(arrayType, elementCount, data)
}
func (_this *markupValueEncoder) EncodeStringlikeArray(ctx *EncoderContext, arrayType events.ArrayType, data string) {
	_this.prepareToWrite(ctx)
	ctx.WriteStringlikeArray(arrayType, data)
}
func (_this *markupValueEncoder) BeginArray(ctx *EncoderContext, arrayType events.ArrayType) {
	_this.prepareForContainer(ctx)
	ctx.BeginStandardArray(arrayType)
}

// =============================================================================

type markupContentsEncoder struct{}

var globalMarkupContentsEncoder markupContentsEncoder

func (_this *markupContentsEncoder) String() string { return "markupContentsEncoder" }

func (_this *markupContentsEncoder) beginObject(ctx *EncoderContext) {
	if !ctx.ContainerHasObjects {
		ctx.Stream.WriteByte(',')
		ctx.SetStandardIndentPrefix()
	}
	ctx.WriteCurrentPrefix()
	ctx.ClearPrefix()
	ctx.ContainerHasObjects = true
}

func (_this *markupContentsEncoder) beginContainer(ctx *EncoderContext) {
	if !ctx.ContainerHasObjects {
		ctx.Stream.WriteByte(',')
		ctx.SetStandardIndentPrefix()
	}
	ctx.WriteCurrentPrefix()
}

func (_this *markupContentsEncoder) Begin(ctx *EncoderContext) {
	ctx.ContainerHasObjects = false
	ctx.ClearPrefix()
}

func (_this *markupContentsEncoder) End(ctx *EncoderContext) {
	ctx.DecreaseIndent()
	if ctx.ContainerHasObjects {
		ctx.WriteBasicIndent()
	}
	ctx.Stream.WriteMarkupEnd()
	ctx.Unstack()
	ctx.CurrentEncoder.ChildContainerFinished(ctx, true)
}

func (_this *markupContentsEncoder) ChildContainerFinished(ctx *EncoderContext, isVisibleChild bool) {
	ctx.SetStandardIndentPrefix()
	ctx.ContainerHasObjects = true
}

func (_this *markupContentsEncoder) BeginMarkup(ctx *EncoderContext) {
	_this.beginContainer(ctx)
	ctx.BeginStandardMarkup()
}
func (_this *markupContentsEncoder) BeginComment(ctx *EncoderContext) {
	_this.beginContainer(ctx)
	ctx.BeginStandardComment()
}
func (_this *markupContentsEncoder) EncodeArray(ctx *EncoderContext, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	_this.beginObject(ctx)
	ctx.WriteMarkupContentStringData(data)
}
func (_this *markupContentsEncoder) EncodeStringlikeArray(ctx *EncoderContext, arrayType events.ArrayType, data string) {
	_this.beginObject(ctx)
	ctx.WriteMarkupContentString(data)
}
func (_this *markupContentsEncoder) BeginArray(ctx *EncoderContext, arrayType events.ArrayType) {
	_this.beginObject(ctx)
	// TODO: begin markup content string array instead of standard array
	ctx.BeginStandardArray(arrayType)
}

// =============================================================================

type commentEncoder struct{}

var globalCommentEncoder commentEncoder

func (_this *commentEncoder) String() string { return "commentEncoder" }

func (_this *commentEncoder) beginItem(ctx *EncoderContext) {
	ctx.ContainerHasObjects = true
	ctx.Stream.WriteByte(' ')
}

func (_this *commentEncoder) Begin(ctx *EncoderContext) {
	ctx.ContainerHasObjects = false
	ctx.Stream.WriteString("/*")
}

func (_this *commentEncoder) End(ctx *EncoderContext) {
	if ctx.ContainerHasObjects {
		ctx.Stream.WriteByte(' ')
	}
	ctx.Stream.WriteString("*/")
	ctx.Unstack()
	ctx.SetStandardIndentPrefix()
	ctx.CurrentEncoder.ChildContainerFinished(ctx, false)
}

func (_this *commentEncoder) ChildContainerFinished(ctx *EncoderContext, isVisibleChild bool) {
	ctx.ContainerHasObjects = true
}

func (_this *commentEncoder) BeginComment(ctx *EncoderContext) {
	_this.beginItem(ctx)
	ctx.BeginStandardComment()
}
func (_this *commentEncoder) EncodeArray(ctx *EncoderContext, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	_this.beginItem(ctx)
	ctx.WriteCommentStringData(data)
}
func (_this *commentEncoder) EncodeStringlikeArray(ctx *EncoderContext, arrayType events.ArrayType, data string) {
	_this.beginItem(ctx)
	ctx.WriteCommentString(data)
}
func (_this *commentEncoder) BeginArray(ctx *EncoderContext, arrayType events.ArrayType) {
	_this.beginItem(ctx)
	// TODO: begin comment content string array instead of standard array
	ctx.BeginStandardArray(arrayType)
}
