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

type naCatEncoder struct{}

var globalNACatEncoder naCatEncoder

func (_this *naCatEncoder) String() string { return "naCatEncoder" }

func (_this *naCatEncoder) Begin(ctx *EncoderContext) {
	ctx.Stream.AddString("@na")
}

func (_this *naCatEncoder) prepareToWrite(ctx *EncoderContext) {
	ctx.Stream.WriteSeparator()
	ctx.Unstack()
}

func (_this *naCatEncoder) EncodeBool(ctx *EncoderContext, value bool) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteBool(value)
	ctx.CurrentEncoder.ChildContainerFinished(ctx, true)
}
func (_this *naCatEncoder) EncodeTrue(ctx *EncoderContext) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteTrue()
	ctx.CurrentEncoder.ChildContainerFinished(ctx, true)
}
func (_this *naCatEncoder) EncodeFalse(ctx *EncoderContext) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteFalse()
	ctx.CurrentEncoder.ChildContainerFinished(ctx, true)
}
func (_this *naCatEncoder) EncodePositiveInt(ctx *EncoderContext, value uint64) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WritePositiveInt(value)
	ctx.CurrentEncoder.ChildContainerFinished(ctx, true)
}
func (_this *naCatEncoder) EncodeNegativeInt(ctx *EncoderContext, value uint64) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteNegativeInt(value)
	ctx.CurrentEncoder.ChildContainerFinished(ctx, true)
}
func (_this *naCatEncoder) EncodeInt(ctx *EncoderContext, value int64) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteInt(value)
	ctx.CurrentEncoder.ChildContainerFinished(ctx, true)
}
func (_this *naCatEncoder) EncodeBigInt(ctx *EncoderContext, value *big.Int) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteBigInt(value)
	ctx.CurrentEncoder.ChildContainerFinished(ctx, true)
}
func (_this *naCatEncoder) EncodeFloat(ctx *EncoderContext, value float64) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteFloat(value)
	ctx.CurrentEncoder.ChildContainerFinished(ctx, true)
}
func (_this *naCatEncoder) EncodeBigFloat(ctx *EncoderContext, value *big.Float) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteBigFloat(value)
	ctx.CurrentEncoder.ChildContainerFinished(ctx, true)
}
func (_this *naCatEncoder) EncodeDecimalFloat(ctx *EncoderContext, value compact_float.DFloat) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteDecimalFloat(value)
	ctx.CurrentEncoder.ChildContainerFinished(ctx, true)
}
func (_this *naCatEncoder) EncodeBigDecimalFloat(ctx *EncoderContext, value *apd.Decimal) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteBigDecimalFloat(value)
	ctx.CurrentEncoder.ChildContainerFinished(ctx, true)
}
func (_this *naCatEncoder) EncodeNan(ctx *EncoderContext, signaling bool) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteNan(signaling)
	ctx.CurrentEncoder.ChildContainerFinished(ctx, true)
}
func (_this *naCatEncoder) EncodeTime(ctx *EncoderContext, value time.Time) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteTime(value)
	ctx.CurrentEncoder.ChildContainerFinished(ctx, true)
}
func (_this *naCatEncoder) EncodeCompactTime(ctx *EncoderContext, value compact_time.Time) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteCompactTime(value)
	ctx.CurrentEncoder.ChildContainerFinished(ctx, true)
}
func (_this *naCatEncoder) EncodeUUID(ctx *EncoderContext, value []byte) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteUUID(value)
	ctx.CurrentEncoder.ChildContainerFinished(ctx, true)
}
func (_this *naCatEncoder) BeginList(ctx *EncoderContext) {
	_this.prepareToWrite(ctx)
	ctx.BeginStandardList()
}
func (_this *naCatEncoder) BeginMap(ctx *EncoderContext) {
	_this.prepareToWrite(ctx)
	ctx.BeginStandardMap()
}
func (_this *naCatEncoder) BeginMarkup(ctx *EncoderContext) {
	_this.prepareToWrite(ctx)
	ctx.BeginStandardMarkup()
}
func (_this *naCatEncoder) BeginMarker(ctx *EncoderContext) {
	_this.prepareToWrite(ctx)
	ctx.BeginStandardMarker()
}
func (_this *naCatEncoder) BeginReference(ctx *EncoderContext) {
	_this.prepareToWrite(ctx)
	ctx.BeginStandardReference()
}
func (_this *naCatEncoder) BeginConstant(ctx *EncoderContext, name []byte, explicitValue bool) {
	_this.prepareToWrite(ctx)
	ctx.BeginStandardConstant(name, explicitValue)
}
func (_this *naCatEncoder) BeginNA(ctx *EncoderContext) {
	// Only unstack
	ctx.Unstack()
	ctx.CurrentEncoder.ChildContainerFinished(ctx, true)
}
func (_this *naCatEncoder) EncodeArray(ctx *EncoderContext, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	_this.prepareToWrite(ctx)
	ctx.ArrayEngine.EncodeArray(arrayType, elementCount, data)
	ctx.CurrentEncoder.ChildContainerFinished(ctx, true)
}
func (_this *naCatEncoder) EncodeStringlikeArray(ctx *EncoderContext, arrayType events.ArrayType, data string) {
	_this.prepareToWrite(ctx)
	ctx.ArrayEngine.EncodeStringlikeArray(arrayType, data)
	ctx.CurrentEncoder.ChildContainerFinished(ctx, true)
}
func (_this *naCatEncoder) BeginArray(ctx *EncoderContext, arrayType events.ArrayType) {
	_this.prepareToWrite(ctx)
	ctx.BeginStandardArray(arrayType)
}

// =============================================================================

type constantEncoder struct{}

var globalConstantEncoder constantEncoder

func (_this *constantEncoder) String() string { return "constantEncoder" }

func (_this *constantEncoder) Begin(ctx *EncoderContext) {
}

func (_this *constantEncoder) prepareToWrite(ctx *EncoderContext) {
	ctx.Stream.WriteSeparator()
	ctx.Unstack()
}

func (_this *constantEncoder) EncodeNA(ctx *EncoderContext) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteNA()
	ctx.CurrentEncoder.ChildContainerFinished(ctx, true)
}
func (_this *constantEncoder) BeginNACat(ctx *EncoderContext) {
	_this.prepareToWrite(ctx)
	ctx.BeginNACat()
}
func (_this *constantEncoder) EncodeBool(ctx *EncoderContext, value bool) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteBool(value)
	ctx.CurrentEncoder.ChildContainerFinished(ctx, true)
}
func (_this *constantEncoder) EncodeTrue(ctx *EncoderContext) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteTrue()
	ctx.CurrentEncoder.ChildContainerFinished(ctx, true)
}
func (_this *constantEncoder) EncodeFalse(ctx *EncoderContext) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteFalse()
	ctx.CurrentEncoder.ChildContainerFinished(ctx, true)
}
func (_this *constantEncoder) EncodePositiveInt(ctx *EncoderContext, value uint64) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WritePositiveInt(value)
	ctx.CurrentEncoder.ChildContainerFinished(ctx, true)
}
func (_this *constantEncoder) EncodeNegativeInt(ctx *EncoderContext, value uint64) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteNegativeInt(value)
	ctx.CurrentEncoder.ChildContainerFinished(ctx, true)
}
func (_this *constantEncoder) EncodeInt(ctx *EncoderContext, value int64) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteInt(value)
	ctx.CurrentEncoder.ChildContainerFinished(ctx, true)
}
func (_this *constantEncoder) EncodeBigInt(ctx *EncoderContext, value *big.Int) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteBigInt(value)
	ctx.CurrentEncoder.ChildContainerFinished(ctx, true)
}
func (_this *constantEncoder) EncodeFloat(ctx *EncoderContext, value float64) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteFloat(value)
	ctx.CurrentEncoder.ChildContainerFinished(ctx, true)
}
func (_this *constantEncoder) EncodeBigFloat(ctx *EncoderContext, value *big.Float) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteBigFloat(value)
	ctx.CurrentEncoder.ChildContainerFinished(ctx, true)
}
func (_this *constantEncoder) EncodeDecimalFloat(ctx *EncoderContext, value compact_float.DFloat) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteDecimalFloat(value)
	ctx.CurrentEncoder.ChildContainerFinished(ctx, true)
}
func (_this *constantEncoder) EncodeBigDecimalFloat(ctx *EncoderContext, value *apd.Decimal) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteBigDecimalFloat(value)
	ctx.CurrentEncoder.ChildContainerFinished(ctx, true)
}
func (_this *constantEncoder) EncodeNan(ctx *EncoderContext, signaling bool) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteNan(signaling)
	ctx.CurrentEncoder.ChildContainerFinished(ctx, true)
}
func (_this *constantEncoder) EncodeTime(ctx *EncoderContext, value time.Time) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteTime(value)
	ctx.CurrentEncoder.ChildContainerFinished(ctx, true)
}
func (_this *constantEncoder) EncodeCompactTime(ctx *EncoderContext, value compact_time.Time) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteCompactTime(value)
	ctx.CurrentEncoder.ChildContainerFinished(ctx, true)
}
func (_this *constantEncoder) EncodeUUID(ctx *EncoderContext, value []byte) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteUUID(value)
	ctx.CurrentEncoder.ChildContainerFinished(ctx, true)
}
func (_this *constantEncoder) BeginList(ctx *EncoderContext) {
	_this.prepareToWrite(ctx)
	ctx.BeginStandardList()
}
func (_this *constantEncoder) BeginMap(ctx *EncoderContext) {
	_this.prepareToWrite(ctx)
	ctx.BeginStandardMap()
}
func (_this *constantEncoder) BeginMarkup(ctx *EncoderContext) {
	_this.prepareToWrite(ctx)
	ctx.BeginStandardMarkup()
}
func (_this *constantEncoder) BeginMarker(ctx *EncoderContext) {
	_this.prepareToWrite(ctx)
	ctx.BeginStandardMarker()
}
func (_this *constantEncoder) BeginReference(ctx *EncoderContext) {
	_this.prepareToWrite(ctx)
	ctx.BeginStandardReference()
}
func (_this *constantEncoder) BeginConstant(ctx *EncoderContext, name []byte, explicitValue bool) {
	_this.prepareToWrite(ctx)
	ctx.BeginStandardConstant(name, explicitValue)
}
func (_this *constantEncoder) EncodeArray(ctx *EncoderContext, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	_this.prepareToWrite(ctx)
	ctx.ArrayEngine.EncodeArray(arrayType, elementCount, data)
	ctx.CurrentEncoder.ChildContainerFinished(ctx, true)
}
func (_this *constantEncoder) EncodeStringlikeArray(ctx *EncoderContext, arrayType events.ArrayType, data string) {
	_this.prepareToWrite(ctx)
	ctx.ArrayEngine.EncodeStringlikeArray(arrayType, data)
	ctx.CurrentEncoder.ChildContainerFinished(ctx, true)
}
func (_this *constantEncoder) BeginArray(ctx *EncoderContext, arrayType events.ArrayType) {
	_this.prepareToWrite(ctx)
	ctx.BeginStandardArray(arrayType)
}

// =============================================================================

type postInvisibleEncoder struct{}

var globalPostInvisibleEncoder postInvisibleEncoder

func (_this *postInvisibleEncoder) String() string { return "postInvisibleEncoder" }

func (_this *postInvisibleEncoder) removeSelf(ctx *EncoderContext) Encoder {
	ctx.Unstack()
	ctx.ClearPrefix()
	return ctx.CurrentEncoder
}

func (_this *postInvisibleEncoder) EncodeNA(ctx *EncoderContext) {
	_this.removeSelf(ctx).EncodeNA(ctx)
}
func (_this *postInvisibleEncoder) BeginNACat(ctx *EncoderContext) {
	_this.removeSelf(ctx).BeginNACat(ctx)
}
func (_this *postInvisibleEncoder) EncodeBool(ctx *EncoderContext, value bool) {
	_this.removeSelf(ctx).EncodeBool(ctx, value)
}
func (_this *postInvisibleEncoder) EncodeTrue(ctx *EncoderContext) {
	_this.removeSelf(ctx).EncodeTrue(ctx)
}
func (_this *postInvisibleEncoder) EncodeFalse(ctx *EncoderContext) {
	_this.removeSelf(ctx).EncodeFalse(ctx)
}
func (_this *postInvisibleEncoder) EncodePositiveInt(ctx *EncoderContext, value uint64) {
	_this.removeSelf(ctx).EncodePositiveInt(ctx, value)
}
func (_this *postInvisibleEncoder) EncodeNegativeInt(ctx *EncoderContext, value uint64) {
	_this.removeSelf(ctx).EncodeNegativeInt(ctx, value)
}
func (_this *postInvisibleEncoder) EncodeInt(ctx *EncoderContext, value int64) {
	_this.removeSelf(ctx).EncodeInt(ctx, value)
}
func (_this *postInvisibleEncoder) EncodeBigInt(ctx *EncoderContext, value *big.Int) {
	_this.removeSelf(ctx).EncodeBigInt(ctx, value)
}
func (_this *postInvisibleEncoder) EncodeFloat(ctx *EncoderContext, value float64) {
	_this.removeSelf(ctx).EncodeFloat(ctx, value)
}
func (_this *postInvisibleEncoder) EncodeBigFloat(ctx *EncoderContext, value *big.Float) {
	_this.removeSelf(ctx).EncodeBigFloat(ctx, value)
}
func (_this *postInvisibleEncoder) EncodeDecimalFloat(ctx *EncoderContext, value compact_float.DFloat) {
	_this.removeSelf(ctx).EncodeDecimalFloat(ctx, value)
}
func (_this *postInvisibleEncoder) EncodeBigDecimalFloat(ctx *EncoderContext, value *apd.Decimal) {
	_this.removeSelf(ctx).EncodeBigDecimalFloat(ctx, value)
}
func (_this *postInvisibleEncoder) EncodeNan(ctx *EncoderContext, signaling bool) {
	_this.removeSelf(ctx).EncodeNan(ctx, signaling)
}
func (_this *postInvisibleEncoder) EncodeTime(ctx *EncoderContext, value time.Time) {
	_this.removeSelf(ctx).EncodeTime(ctx, value)
}
func (_this *postInvisibleEncoder) EncodeCompactTime(ctx *EncoderContext, value compact_time.Time) {
	_this.removeSelf(ctx).EncodeCompactTime(ctx, value)
}
func (_this *postInvisibleEncoder) EncodeUUID(ctx *EncoderContext, value []byte) {
	_this.removeSelf(ctx).EncodeUUID(ctx, value)
}
func (_this *postInvisibleEncoder) BeginList(ctx *EncoderContext) {
	_this.removeSelf(ctx).BeginList(ctx)
}
func (_this *postInvisibleEncoder) BeginMap(ctx *EncoderContext) {
	_this.removeSelf(ctx).BeginMap(ctx)
}
func (_this *postInvisibleEncoder) BeginMarkup(ctx *EncoderContext) {
	_this.removeSelf(ctx).BeginMarkup(ctx)
}
func (_this *postInvisibleEncoder) BeginMetadata(ctx *EncoderContext) {
	_this.removeSelf(ctx).BeginMetadata(ctx)
}
func (_this *postInvisibleEncoder) BeginComment(ctx *EncoderContext) {
	_this.removeSelf(ctx).BeginComment(ctx)
}
func (_this *postInvisibleEncoder) BeginMarker(ctx *EncoderContext) {
	_this.removeSelf(ctx).BeginMarker(ctx)
}
func (_this *postInvisibleEncoder) BeginReference(ctx *EncoderContext) {
	_this.removeSelf(ctx).BeginReference(ctx)
}
func (_this *postInvisibleEncoder) BeginConstant(ctx *EncoderContext, name []byte, explicitValue bool) {
	_this.removeSelf(ctx).BeginConstant(ctx, name, explicitValue)
}
func (_this *postInvisibleEncoder) EncodeArray(ctx *EncoderContext, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	_this.removeSelf(ctx).EncodeArray(ctx, arrayType, elementCount, data)
}
func (_this *postInvisibleEncoder) EncodeStringlikeArray(ctx *EncoderContext, arrayType events.ArrayType, data string) {
	_this.removeSelf(ctx).EncodeStringlikeArray(ctx, arrayType, data)
}
func (_this *postInvisibleEncoder) BeginArray(ctx *EncoderContext, arrayType events.ArrayType) {
	_this.removeSelf(ctx).BeginArray(ctx, arrayType)
}

// =============================================================================

type referenceEncoder struct{}

var globalReferenceEncoder referenceEncoder

func (_this *referenceEncoder) String() string { return "referenceEncoder" }

func (_this *referenceEncoder) complete(ctx *EncoderContext) {
	ctx.Unstack()
	ctx.CurrentEncoder.ChildContainerFinished(ctx, true)
}

func (_this *referenceEncoder) Begin(ctx *EncoderContext) {
	ctx.Stream.AddByte('$')
}

func (_this *referenceEncoder) ChildContainerFinished(ctx *EncoderContext, isVisibleChild bool) {
	_this.complete(ctx)
}

func (_this *referenceEncoder) EncodePositiveInt(ctx *EncoderContext, value uint64) {
	ctx.Stream.WritePositiveInt(value)
	_this.complete(ctx)
}
func (_this *referenceEncoder) EncodeInt(ctx *EncoderContext, value int64) {
	ctx.Stream.WriteInt(value)
	_this.complete(ctx)
}
func (_this *referenceEncoder) EncodeBigInt(ctx *EncoderContext, value *big.Int) {
	ctx.Stream.WriteBigInt(value)
	_this.complete(ctx)
}
func (_this *referenceEncoder) BeginConstant(ctx *EncoderContext, name []byte, explicitValue bool) {
	panic("TODO: referenceEncoder.BeginConstant")
}
func (_this *referenceEncoder) EncodeArray(ctx *EncoderContext, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	ctx.ArrayEngine.EncodeArray(arrayType, elementCount, data)
	_this.complete(ctx)
}
func (_this *referenceEncoder) EncodeStringlikeArray(ctx *EncoderContext, arrayType events.ArrayType, data string) {
	ctx.ArrayEngine.EncodeStringlikeArray(arrayType, data)
	_this.complete(ctx)
}
func (_this *referenceEncoder) BeginArray(ctx *EncoderContext, arrayType events.ArrayType) {
	ctx.BeginStandardArray(arrayType)
}

// =============================================================================

type markerIDEncoder struct{}

var globalMarkerIDEncoder markerIDEncoder

func (_this *markerIDEncoder) String() string { return "markerIDEncoder" }

func (_this *markerIDEncoder) complete(ctx *EncoderContext) {
	ctx.Unstack()
	ctx.Stream.AddByte(':')
	ctx.ClearPrefix()
}

func (_this *markerIDEncoder) Begin(ctx *EncoderContext) {
	ctx.Stream.AddByte('&')
}

func (_this *markerIDEncoder) ChildContainerFinished(ctx *EncoderContext, isVisibleChild bool) {
	_this.complete(ctx)
}

func (_this *markerIDEncoder) EncodePositiveInt(ctx *EncoderContext, value uint64) {
	ctx.Stream.WritePositiveInt(value)
	_this.complete(ctx)
}
func (_this *markerIDEncoder) EncodeInt(ctx *EncoderContext, value int64) {
	ctx.Stream.WriteInt(value)
	_this.complete(ctx)
}
func (_this *markerIDEncoder) EncodeBigInt(ctx *EncoderContext, value *big.Int) {
	ctx.Stream.WriteBigInt(value)
	_this.complete(ctx)
}
func (_this *markerIDEncoder) BeginConstant(ctx *EncoderContext, name []byte, explicitValue bool) {
	panic("TODO: markerIDEncoder.BeginConstant")
}
func (_this *markerIDEncoder) EncodeArray(ctx *EncoderContext, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	ctx.ArrayEngine.EncodeArray(arrayType, elementCount, data)
	_this.complete(ctx)
}
func (_this *markerIDEncoder) EncodeStringlikeArray(ctx *EncoderContext, arrayType events.ArrayType, data string) {
	ctx.ArrayEngine.EncodeStringlikeArray(arrayType, data)
	_this.complete(ctx)
}
func (_this *markerIDEncoder) BeginArray(ctx *EncoderContext, arrayType events.ArrayType) {
	ctx.BeginStandardArray(arrayType)
}

// =============================================================================

type postRIDCatEncoder struct{}

var globalPostRIDCatEncoder postRIDCatEncoder

func (_this *postRIDCatEncoder) String() string { return "postRIDCatEncoder" }

func (_this *postRIDCatEncoder) ChildContainerFinished(ctx *EncoderContext, isVisibleChild bool) {
	ctx.Unstack()
}

func (_this *postRIDCatEncoder) EncodePositiveInt(ctx *EncoderContext, value uint64) {
	ctx.Stream.WriteConcat()
	ctx.Stream.WritePositiveInt(value)
	ctx.Unstack()
}
func (_this *postRIDCatEncoder) EncodeInt(ctx *EncoderContext, value int64) {
	ctx.Stream.WriteConcat()
	ctx.Stream.WriteInt(value)
	ctx.Unstack()
}
func (_this *postRIDCatEncoder) EncodeBigInt(ctx *EncoderContext, value *big.Int) {
	ctx.Stream.WriteConcat()
	ctx.Stream.WriteBigInt(value)
	ctx.Unstack()
}
func (_this *postRIDCatEncoder) EncodeArray(ctx *EncoderContext, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	ctx.Stream.WriteConcat()
	ctx.ArrayEngine.EncodeArray(arrayType, elementCount, data)
	ctx.Unstack()
}
func (_this *postRIDCatEncoder) EncodeStringlikeArray(ctx *EncoderContext, arrayType events.ArrayType, data string) {
	ctx.Stream.WriteConcat()
	ctx.ArrayEngine.EncodeStringlikeArray(arrayType, data)
	ctx.Unstack()
}
func (_this *postRIDCatEncoder) BeginArray(ctx *EncoderContext, arrayType events.ArrayType) {
	ctx.Stream.WriteConcat()
	ctx.BeginStandardArray(arrayType)
}

// =============================================================================

type postStreamRIDCatEncoder struct{}

var globalPostStreamRIDCatEncoder postStreamRIDCatEncoder

func (_this *postStreamRIDCatEncoder) String() string { return "postStreamRIDCatEncoder" }

func (_this *postStreamRIDCatEncoder) ChildContainerFinished(ctx *EncoderContext, isVisibleChild bool) {
	ctx.Unstack()
	ctx.CurrentEncoder.ChildContainerFinished(ctx, true)
}

func (_this *postStreamRIDCatEncoder) EncodePositiveInt(ctx *EncoderContext, value uint64) {
	ctx.Stream.WriteConcat()
	ctx.Stream.WritePositiveInt(value)
	ctx.Unstack()
	ctx.CurrentEncoder.ChildContainerFinished(ctx, true)
}
func (_this *postStreamRIDCatEncoder) EncodeInt(ctx *EncoderContext, value int64) {
	ctx.Stream.WriteConcat()
	ctx.Stream.WriteInt(value)
	ctx.Unstack()
	ctx.CurrentEncoder.ChildContainerFinished(ctx, true)
}
func (_this *postStreamRIDCatEncoder) EncodeBigInt(ctx *EncoderContext, value *big.Int) {
	ctx.Stream.WriteConcat()
	ctx.Stream.WriteBigInt(value)
	ctx.Unstack()
	ctx.CurrentEncoder.ChildContainerFinished(ctx, true)
}
func (_this *postStreamRIDCatEncoder) EncodeArray(ctx *EncoderContext, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	ctx.Stream.WriteConcat()
	ctx.ArrayEngine.EncodeArray(arrayType, elementCount, data)
	ctx.Unstack()
	ctx.CurrentEncoder.ChildContainerFinished(ctx, true)
}
func (_this *postStreamRIDCatEncoder) EncodeStringlikeArray(ctx *EncoderContext, arrayType events.ArrayType, data string) {
	ctx.Stream.WriteConcat()
	ctx.ArrayEngine.EncodeStringlikeArray(arrayType, data)
	ctx.Unstack()
	ctx.CurrentEncoder.ChildContainerFinished(ctx, true)
}
func (_this *postStreamRIDCatEncoder) BeginArray(ctx *EncoderContext, arrayType events.ArrayType) {
	ctx.Stream.WriteConcat()
	ctx.BeginStandardArray(arrayType)
}
