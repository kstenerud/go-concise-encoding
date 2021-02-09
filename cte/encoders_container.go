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

func (_this *listEncoder) Begin(ctx *EncoderContext) {
	ctx.WriteCurrentPrefix()
	ctx.Stream.WriteListBegin()
	ctx.IncreaseIndent()
	ctx.SetStandardIndentPrefix()
	ctx.ContainerHasObjects = false
}

func (_this *listEncoder) End(ctx *EncoderContext) {
	ctx.DecreaseIndent()
	if ctx.ContainerHasObjects {
		ctx.WriteIndent()
	}
	ctx.Stream.WriteListEnd()
	ctx.ContainerHasObjects = true
}

func (_this *listEncoder) ChildContainerFinished(ctx *EncoderContext) {
	ctx.SetStandardIndentPrefix()
}

func (_this *listEncoder) EncodeBool(ctx *EncoderContext, value bool) {
	ctx.WriteCurrentPrefix()
	ctx.Stream.WriteBool(value)
	ctx.ContainerHasObjects = true
}
func (_this *listEncoder) EncodeTrue(ctx *EncoderContext) {
	ctx.WriteCurrentPrefix()
	ctx.Stream.WriteTrue()
	ctx.ContainerHasObjects = true
}
func (_this *listEncoder) EncodeFalse(ctx *EncoderContext) {
	ctx.WriteCurrentPrefix()
	ctx.Stream.WriteFalse()
	ctx.ContainerHasObjects = true
}
func (_this *listEncoder) EncodePositiveInt(ctx *EncoderContext, value uint64) {
	ctx.WriteCurrentPrefix()
	ctx.Stream.WritePositiveInt(value)
	ctx.ContainerHasObjects = true
}
func (_this *listEncoder) EncodeNegativeInt(ctx *EncoderContext, value uint64) {
	ctx.WriteCurrentPrefix()
	ctx.Stream.WriteNegativeInt(value)
	ctx.ContainerHasObjects = true
}
func (_this *listEncoder) EncodeInt(ctx *EncoderContext, value int64) {
	ctx.WriteCurrentPrefix()
	ctx.Stream.WriteInt(value)
	ctx.ContainerHasObjects = true
}
func (_this *listEncoder) EncodeBigInt(ctx *EncoderContext, value *big.Int) {
	ctx.WriteCurrentPrefix()
	ctx.Stream.WriteBigInt(value)
	ctx.ContainerHasObjects = true
}
func (_this *listEncoder) EncodeFloat(ctx *EncoderContext, value float64) {
	ctx.WriteCurrentPrefix()
	ctx.Stream.WriteFloat(value)
	ctx.ContainerHasObjects = true
}
func (_this *listEncoder) EncodeBigFloat(ctx *EncoderContext, value *big.Float) {
	ctx.WriteCurrentPrefix()
	ctx.Stream.WriteBigFloat(value)
	ctx.ContainerHasObjects = true
}
func (_this *listEncoder) EncodeDecimalFloat(ctx *EncoderContext, value compact_float.DFloat) {
	ctx.WriteCurrentPrefix()
	ctx.Stream.WriteDecimalFloat(value)
	ctx.ContainerHasObjects = true
}
func (_this *listEncoder) EncodeBigDecimalFloat(ctx *EncoderContext, value *apd.Decimal) {
	ctx.WriteCurrentPrefix()
	ctx.Stream.WriteBigDecimalFloat(value)
	ctx.ContainerHasObjects = true
}
func (_this *listEncoder) EncodeNan(ctx *EncoderContext, signaling bool) {
	ctx.WriteCurrentPrefix()
	ctx.Stream.WriteNan(signaling)
	ctx.ContainerHasObjects = true
}
func (_this *listEncoder) EncodeTime(ctx *EncoderContext, value time.Time) {
	ctx.WriteCurrentPrefix()
	ctx.Stream.WriteTime(value)
	ctx.ContainerHasObjects = true
}
func (_this *listEncoder) EncodeCompactTime(ctx *EncoderContext, value compact_time.Time) {
	ctx.WriteCurrentPrefix()
	ctx.Stream.WriteCompactTime(value)
	ctx.ContainerHasObjects = true
}
func (_this *listEncoder) EncodeUUID(ctx *EncoderContext, value []byte) {
	ctx.WriteCurrentPrefix()
	ctx.Stream.WriteUUID(value)
	ctx.ContainerHasObjects = true
}
func (_this *listEncoder) BeginList(ctx *EncoderContext) {
	ctx.WriteCurrentPrefix()
	ctx.BeginList()
	ctx.ContainerHasObjects = true
}
func (_this *listEncoder) BeginMap(ctx *EncoderContext) {
	ctx.WriteCurrentPrefix()
	ctx.BeginMap()
	ctx.ContainerHasObjects = true
}
func (_this *listEncoder) BeginMarkup(ctx *EncoderContext) {
	ctx.WriteCurrentPrefix()
	ctx.BeginMarkup()
	ctx.ContainerHasObjects = true
}
func (_this *listEncoder) BeginMetadata(ctx *EncoderContext) {
	ctx.WriteCurrentPrefix()
	ctx.BeginMetadata()
	ctx.ContainerHasObjects = true
}
func (_this *listEncoder) BeginComment(ctx *EncoderContext) {
	ctx.WriteCurrentPrefix()
	ctx.BeginComment()
	ctx.ContainerHasObjects = true
}
func (_this *listEncoder) BeginMarker(ctx *EncoderContext) {
	ctx.WriteCurrentPrefix()
	ctx.BeginMarker()
	ctx.ContainerHasObjects = true
}
func (_this *listEncoder) BeginReference(ctx *EncoderContext) {
	ctx.WriteCurrentPrefix()
	ctx.BeginReference()
	ctx.ContainerHasObjects = true
}
func (_this *listEncoder) BeginConcatenate(ctx *EncoderContext) {
	ctx.WriteCurrentPrefix()
	ctx.BeginConcatenate()
	ctx.ContainerHasObjects = true
}
func (_this *listEncoder) BeginConstant(ctx *EncoderContext, name []byte, explicitValue bool) {
	ctx.WriteCurrentPrefix()
	ctx.BeginConstant(name, explicitValue)
	ctx.ContainerHasObjects = true
}
func (_this *listEncoder) BeginNA(ctx *EncoderContext) {
	ctx.WriteCurrentPrefix()
	ctx.BeginNA()
	ctx.ContainerHasObjects = true
}
func (_this *listEncoder) EncodeArray(ctx *EncoderContext, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	ctx.WriteCurrentPrefix()
	ctx.Stream.WriteArray(arrayType, elementCount, data)
	ctx.ContainerHasObjects = true
}
func (_this *listEncoder) EncodeStringlikeArray(ctx *EncoderContext, arrayType events.ArrayType, data string) {
	ctx.WriteCurrentPrefix()
	ctx.Stream.WriteStringlikeArray(arrayType, data)
	ctx.ContainerHasObjects = true
}
func (_this *listEncoder) BeginArray(ctx *EncoderContext, arrayType events.ArrayType) {
	ctx.WriteCurrentPrefix()
	ctx.BeginArray(arrayType)
	ctx.ContainerHasObjects = true
}

// =============================================================================

func encodeMapSeparator(ctx *EncoderContext) {
	ctx.Stream.AddString(" = ")
}

type mapKeyEncoder struct{}

var globalMapKeyEncoder mapKeyEncoder

func (_this *mapKeyEncoder) prepareToWrite(ctx *EncoderContext) {
	ctx.WriteCurrentPrefix()
	ctx.SetStandardMapValuePrefix()
	ctx.ChangeEncoder(&globalMapValueEncoder)
	ctx.ContainerHasObjects = true
}

func (_this *mapKeyEncoder) Begin(ctx *EncoderContext) {
	ctx.WriteCurrentPrefix()
	ctx.Stream.WriteMapBegin()
	ctx.IncreaseIndent()
	ctx.SetStandardMapKeyPrefix()
	ctx.ContainerHasObjects = false
}

func (_this *mapKeyEncoder) End(ctx *EncoderContext) {
	ctx.DecreaseIndent()
	if ctx.ContainerHasObjects {
		ctx.WriteIndent()
	}
	ctx.Stream.WriteMapEnd()
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
func (_this *mapKeyEncoder) EncodeNan(ctx *EncoderContext, signaling bool) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteNan(signaling)
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
	// TODO
	// _this.prepareToWrite(ctx)
	ctx.BeginMetadata()
}
func (_this *mapKeyEncoder) BeginComment(ctx *EncoderContext) {
	// TODO
	// _this.prepareToWrite(ctx)
	ctx.BeginComment()
}
func (_this *mapKeyEncoder) BeginMarker(ctx *EncoderContext) {
	_this.prepareToWrite(ctx)
	ctx.BeginMarker()
}
func (_this *mapKeyEncoder) BeginReference(ctx *EncoderContext) {
	_this.prepareToWrite(ctx)
	ctx.BeginReference()
}
func (_this *mapKeyEncoder) BeginConcatenate(ctx *EncoderContext) {
	// TODO ?
	_this.prepareToWrite(ctx)
	ctx.BeginConcatenate()
}
func (_this *mapKeyEncoder) BeginConstant(ctx *EncoderContext, name []byte, explicitValue bool) {
	_this.prepareToWrite(ctx)
	ctx.BeginConstant(name, explicitValue)
}
func (_this *mapKeyEncoder) BeginNA(ctx *EncoderContext) {
	_this.prepareToWrite(ctx)
	ctx.BeginNA()
}
func (_this *mapKeyEncoder) EncodeArray(ctx *EncoderContext, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteArray(arrayType, elementCount, data)
}
func (_this *mapKeyEncoder) EncodeStringlikeArray(ctx *EncoderContext, arrayType events.ArrayType, data string) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteStringlikeArray(arrayType, data)
}
func (_this *mapKeyEncoder) BeginArray(ctx *EncoderContext, arrayType events.ArrayType) {
	_this.prepareToWrite(ctx)
	ctx.BeginArray(arrayType)
}

// =============================================================================

type mapValueEncoder struct{}

var globalMapValueEncoder mapValueEncoder

func (_this *mapValueEncoder) prepareToWrite(ctx *EncoderContext) {
	ctx.WriteCurrentPrefix()
	ctx.SetStandardMapKeyPrefix()
	ctx.ChangeEncoder(&globalMapKeyEncoder)
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
	ctx.BeginList()
}
func (_this *mapValueEncoder) BeginMap(ctx *EncoderContext) {
	ctx.BeginMap()
}
func (_this *mapValueEncoder) BeginMarkup(ctx *EncoderContext) {
	ctx.BeginMarkup()
}
func (_this *mapValueEncoder) BeginMetadata(ctx *EncoderContext) {
	ctx.BeginMetadata()
}
func (_this *mapValueEncoder) BeginComment(ctx *EncoderContext) {
	ctx.ChangeEncoder(&globalMapKeyEncoder)
	ctx.BeginComment()
}
func (_this *mapValueEncoder) BeginMarker(ctx *EncoderContext) {
	_this.prepareToWrite(ctx)
	ctx.BeginMarker()
}
func (_this *mapValueEncoder) BeginReference(ctx *EncoderContext) {
	_this.prepareToWrite(ctx)
	ctx.BeginReference()
}
func (_this *mapValueEncoder) BeginConcatenate(ctx *EncoderContext) {
	_this.prepareToWrite(ctx)
	ctx.BeginConcatenate()
}
func (_this *mapValueEncoder) BeginConstant(ctx *EncoderContext, name []byte, explicitValue bool) {
	_this.prepareToWrite(ctx)
	ctx.BeginConstant(name, explicitValue)
}
func (_this *mapValueEncoder) BeginNA(ctx *EncoderContext) {
	_this.prepareToWrite(ctx)
	ctx.BeginNA()
}
func (_this *mapValueEncoder) EncodeArray(ctx *EncoderContext, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteArray(arrayType, elementCount, data)
}
func (_this *mapValueEncoder) EncodeStringlikeArray(ctx *EncoderContext, arrayType events.ArrayType, data string) {
	_this.prepareToWrite(ctx)
	ctx.Stream.WriteStringlikeArray(arrayType, data)
}
func (_this *mapValueEncoder) BeginArray(ctx *EncoderContext, arrayType events.ArrayType) {
	_this.prepareToWrite(ctx)
	ctx.BeginArray(arrayType)
}
