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
	ctx.WriteIndent()
	ctx.Stream.WriteListBegin()
	ctx.IncreaseIndent()
}

func (_this *listEncoder) End(ctx *EncoderContext) {
	ctx.DecreaseIndent()
	ctx.WriteIndent()
	ctx.Stream.WriteListEnd()
}

func (_this *listEncoder) EncodeBool(ctx *EncoderContext, value bool) {
	ctx.WriteIndent()
	ctx.Stream.WriteBool(value)
}
func (_this *listEncoder) EncodeTrue(ctx *EncoderContext) {
	ctx.WriteIndent()
	ctx.Stream.WriteTrue()
}
func (_this *listEncoder) EncodeFalse(ctx *EncoderContext) {
	ctx.WriteIndent()
	ctx.Stream.WriteFalse()
}
func (_this *listEncoder) EncodePositiveInt(ctx *EncoderContext, value uint64) {
	ctx.WriteIndent()
	ctx.Stream.WritePositiveInt(value)
}
func (_this *listEncoder) EncodeNegativeInt(ctx *EncoderContext, value uint64) {
	ctx.WriteIndent()
	ctx.Stream.WriteNegativeInt(value)
}
func (_this *listEncoder) EncodeInt(ctx *EncoderContext, value int64) {
	ctx.WriteIndent()
	ctx.Stream.WriteInt(value)
}
func (_this *listEncoder) EncodeBigInt(ctx *EncoderContext, value *big.Int) {
	ctx.WriteIndent()
	ctx.Stream.WriteBigInt(value)
}
func (_this *listEncoder) EncodeFloat(ctx *EncoderContext, value float64) {
	ctx.WriteIndent()
	ctx.Stream.WriteFloat(value)
}
func (_this *listEncoder) EncodeBigFloat(ctx *EncoderContext, value *big.Float) {
	ctx.WriteIndent()
	ctx.Stream.WriteBigFloat(value)
}
func (_this *listEncoder) EncodeDecimalFloat(ctx *EncoderContext, value compact_float.DFloat) {
	ctx.WriteIndent()
	ctx.Stream.WriteDecimalFloat(value)
}
func (_this *listEncoder) EncodeBigDecimalFloat(ctx *EncoderContext, value *apd.Decimal) {
	ctx.WriteIndent()
	ctx.Stream.WriteBigDecimalFloat(value)
}
func (_this *listEncoder) EncodeNan(ctx *EncoderContext, signaling bool) {
	ctx.WriteIndent()
	ctx.Stream.WriteNan(signaling)
}
func (_this *listEncoder) EncodeTime(ctx *EncoderContext, value time.Time) {
	ctx.WriteIndent()
	ctx.Stream.WriteTime(value)
}
func (_this *listEncoder) EncodeCompactTime(ctx *EncoderContext, value compact_time.Time) {
	ctx.WriteIndent()
	ctx.Stream.WriteCompactTime(value)
}
func (_this *listEncoder) EncodeUUID(ctx *EncoderContext, value []byte) {
	ctx.WriteIndent()
	ctx.Stream.WriteUUID(value)
}
func (_this *listEncoder) BeginList(ctx *EncoderContext) {
	ctx.WriteIndent()
	ctx.BeginList()
}
func (_this *listEncoder) BeginMap(ctx *EncoderContext) {
	ctx.WriteIndent()
	ctx.BeginMap()
}
func (_this *listEncoder) BeginMarkup(ctx *EncoderContext) {
	ctx.WriteIndent()
	ctx.BeginMarkup()
}
func (_this *listEncoder) BeginMetadata(ctx *EncoderContext) {
	ctx.WriteIndent()
	ctx.BeginMetadata()
}
func (_this *listEncoder) BeginComment(ctx *EncoderContext) {
	ctx.WriteIndent()
	ctx.BeginComment()
}
func (_this *listEncoder) BeginMarker(ctx *EncoderContext) {
	ctx.WriteIndent()
	ctx.BeginMarker()
}
func (_this *listEncoder) BeginReference(ctx *EncoderContext) {
	ctx.WriteIndent()
	ctx.BeginReference()
}
func (_this *listEncoder) BeginConcatenate(ctx *EncoderContext) {
	ctx.WriteIndent()
	ctx.BeginConcatenate()
}
func (_this *listEncoder) BeginConstant(ctx *EncoderContext, name []byte, explicitValue bool) {
	ctx.WriteIndent()
	ctx.BeginConstant(name, explicitValue)
}
func (_this *listEncoder) BeginNA(ctx *EncoderContext) {
	ctx.WriteIndent()
	ctx.BeginNA()
}
func (_this *listEncoder) EncodeArray(ctx *EncoderContext, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	ctx.WriteIndent()
	ctx.Stream.WriteArray(arrayType, elementCount, data)
}
func (_this *listEncoder) EncodeStringlikeArray(ctx *EncoderContext, arrayType events.ArrayType, data string) {
	ctx.WriteIndent()
	ctx.Stream.WriteStringlikeArray(arrayType, data)
}
func (_this *listEncoder) BeginArray(ctx *EncoderContext, arrayType events.ArrayType) {
	ctx.WriteIndent()
	ctx.BeginArray(arrayType)
}

// =============================================================================

type mapKeyEncoder struct{}

var globalMapKeyEncoder mapKeyEncoder

func (_this *mapKeyEncoder) Begin(ctx *EncoderContext) {
	ctx.WriteIndent()
	ctx.Stream.WriteMapBegin()
	ctx.IncreaseIndent()
}

func (_this *mapKeyEncoder) End(ctx *EncoderContext) {
	ctx.DecreaseIndent()
	ctx.WriteIndent()
	ctx.Stream.WriteMapEnd()
}

func (_this *mapKeyEncoder) EncodeBool(ctx *EncoderContext, value bool) {
	ctx.WriteIndent()
	ctx.ChangeEncoder(&globalMapValueEncoder)
	ctx.Stream.WriteBool(value)
}
func (_this *mapKeyEncoder) EncodeTrue(ctx *EncoderContext) {
	ctx.WriteIndent()
	ctx.ChangeEncoder(&globalMapValueEncoder)
	ctx.Stream.WriteTrue()
}
func (_this *mapKeyEncoder) EncodeFalse(ctx *EncoderContext) {
	ctx.WriteIndent()
	ctx.ChangeEncoder(&globalMapValueEncoder)
	ctx.Stream.WriteFalse()
}
func (_this *mapKeyEncoder) EncodePositiveInt(ctx *EncoderContext, value uint64) {
	ctx.WriteIndent()
	ctx.ChangeEncoder(&globalMapValueEncoder)
	ctx.Stream.WritePositiveInt(value)
}
func (_this *mapKeyEncoder) EncodeNegativeInt(ctx *EncoderContext, value uint64) {
	ctx.WriteIndent()
	ctx.ChangeEncoder(&globalMapValueEncoder)
	ctx.Stream.WriteNegativeInt(value)
}
func (_this *mapKeyEncoder) EncodeInt(ctx *EncoderContext, value int64) {
	ctx.WriteIndent()
	ctx.ChangeEncoder(&globalMapValueEncoder)
	ctx.Stream.WriteInt(value)
}
func (_this *mapKeyEncoder) EncodeBigInt(ctx *EncoderContext, value *big.Int) {
	ctx.WriteIndent()
	ctx.ChangeEncoder(&globalMapValueEncoder)
	ctx.Stream.WriteBigInt(value)
}
func (_this *mapKeyEncoder) EncodeFloat(ctx *EncoderContext, value float64) {
	ctx.WriteIndent()
	ctx.ChangeEncoder(&globalMapValueEncoder)
	ctx.Stream.WriteFloat(value)
}
func (_this *mapKeyEncoder) EncodeBigFloat(ctx *EncoderContext, value *big.Float) {
	ctx.WriteIndent()
	ctx.ChangeEncoder(&globalMapValueEncoder)
	ctx.Stream.WriteBigFloat(value)
}
func (_this *mapKeyEncoder) EncodeDecimalFloat(ctx *EncoderContext, value compact_float.DFloat) {
	ctx.WriteIndent()
	ctx.ChangeEncoder(&globalMapValueEncoder)
	ctx.Stream.WriteDecimalFloat(value)
}
func (_this *mapKeyEncoder) EncodeBigDecimalFloat(ctx *EncoderContext, value *apd.Decimal) {
	ctx.WriteIndent()
	ctx.ChangeEncoder(&globalMapValueEncoder)
	ctx.Stream.WriteBigDecimalFloat(value)
}
func (_this *mapKeyEncoder) EncodeNan(ctx *EncoderContext, signaling bool) {
	ctx.WriteIndent()
	ctx.ChangeEncoder(&globalMapValueEncoder)
	ctx.Stream.WriteNan(signaling)
}
func (_this *mapKeyEncoder) EncodeTime(ctx *EncoderContext, value time.Time) {
	ctx.WriteIndent()
	ctx.ChangeEncoder(&globalMapValueEncoder)
	ctx.Stream.WriteTime(value)
}
func (_this *mapKeyEncoder) EncodeCompactTime(ctx *EncoderContext, value compact_time.Time) {
	ctx.WriteIndent()
	ctx.ChangeEncoder(&globalMapValueEncoder)
	ctx.Stream.WriteCompactTime(value)
}
func (_this *mapKeyEncoder) EncodeUUID(ctx *EncoderContext, value []byte) {
	ctx.WriteIndent()
	ctx.ChangeEncoder(&globalMapValueEncoder)
	ctx.Stream.WriteUUID(value)
}
func (_this *mapKeyEncoder) BeginMetadata(ctx *EncoderContext) {
	ctx.WriteIndent()
	ctx.ChangeEncoder(&globalMapValueEncoder)
	ctx.BeginMetadata()
}
func (_this *mapKeyEncoder) BeginComment(ctx *EncoderContext) {
	ctx.WriteIndent()
	ctx.ChangeEncoder(&globalMapValueEncoder)
	ctx.BeginComment()
}
func (_this *mapKeyEncoder) BeginMarker(ctx *EncoderContext) {
	ctx.WriteIndent()
	ctx.ChangeEncoder(&globalMapValueEncoder)
	ctx.BeginMarker()
}
func (_this *mapKeyEncoder) BeginReference(ctx *EncoderContext) {
	ctx.WriteIndent()
	ctx.ChangeEncoder(&globalMapValueEncoder)
	ctx.BeginReference()
}
func (_this *mapKeyEncoder) BeginConcatenate(ctx *EncoderContext) {
	ctx.WriteIndent()
	ctx.ChangeEncoder(&globalMapValueEncoder)
	ctx.BeginConcatenate()
}
func (_this *mapKeyEncoder) BeginConstant(ctx *EncoderContext, name []byte, explicitValue bool) {
	ctx.WriteIndent()
	ctx.ChangeEncoder(&globalMapValueEncoder)
	ctx.BeginConstant(name, explicitValue)
}
func (_this *mapKeyEncoder) BeginNA(ctx *EncoderContext) {
	ctx.WriteIndent()
	ctx.ChangeEncoder(&globalMapValueEncoder)
	ctx.BeginNA()
}
func (_this *mapKeyEncoder) EncodeArray(ctx *EncoderContext, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	ctx.WriteIndent()
	ctx.ChangeEncoder(&globalMapValueEncoder)
	ctx.Stream.WriteArray(arrayType, elementCount, data)
}
func (_this *mapKeyEncoder) EncodeStringlikeArray(ctx *EncoderContext, arrayType events.ArrayType, data string) {
	ctx.WriteIndent()
	ctx.ChangeEncoder(&globalMapValueEncoder)
	ctx.Stream.WriteStringlikeArray(arrayType, data)
}
func (_this *mapKeyEncoder) BeginArray(ctx *EncoderContext, arrayType events.ArrayType) {
	ctx.WriteIndent()
	ctx.ChangeEncoder(&globalMapValueEncoder)
	ctx.BeginArray(arrayType)
}

// =============================================================================

type mapValueEncoder struct{}

var globalMapValueEncoder mapValueEncoder

func (_this *mapValueEncoder) EncodeBool(ctx *EncoderContext, value bool) {
	ctx.WriteIndent()
	ctx.ChangeEncoder(&globalMapKeyEncoder)
	ctx.Stream.WriteBool(value)
}
func (_this *mapValueEncoder) EncodeTrue(ctx *EncoderContext) {
	ctx.WriteIndent()
	ctx.ChangeEncoder(&globalMapKeyEncoder)
	ctx.Stream.WriteTrue()
}
func (_this *mapValueEncoder) EncodeFalse(ctx *EncoderContext) {
	ctx.WriteIndent()
	ctx.ChangeEncoder(&globalMapKeyEncoder)
	ctx.Stream.WriteFalse()
}
func (_this *mapValueEncoder) EncodePositiveInt(ctx *EncoderContext, value uint64) {
	ctx.WriteIndent()
	ctx.ChangeEncoder(&globalMapKeyEncoder)
	ctx.Stream.WritePositiveInt(value)
}
func (_this *mapValueEncoder) EncodeNegativeInt(ctx *EncoderContext, value uint64) {
	ctx.WriteIndent()
	ctx.ChangeEncoder(&globalMapKeyEncoder)
	ctx.Stream.WriteNegativeInt(value)
}
func (_this *mapValueEncoder) EncodeInt(ctx *EncoderContext, value int64) {
	ctx.WriteIndent()
	ctx.ChangeEncoder(&globalMapKeyEncoder)
	ctx.Stream.WriteInt(value)
}
func (_this *mapValueEncoder) EncodeBigInt(ctx *EncoderContext, value *big.Int) {
	ctx.WriteIndent()
	ctx.ChangeEncoder(&globalMapKeyEncoder)
	ctx.Stream.WriteBigInt(value)
}
func (_this *mapValueEncoder) EncodeFloat(ctx *EncoderContext, value float64) {
	ctx.WriteIndent()
	ctx.ChangeEncoder(&globalMapKeyEncoder)
	ctx.Stream.WriteFloat(value)
}
func (_this *mapValueEncoder) EncodeBigFloat(ctx *EncoderContext, value *big.Float) {
	ctx.WriteIndent()
	ctx.ChangeEncoder(&globalMapKeyEncoder)
	ctx.Stream.WriteBigFloat(value)
}
func (_this *mapValueEncoder) EncodeDecimalFloat(ctx *EncoderContext, value compact_float.DFloat) {
	ctx.WriteIndent()
	ctx.ChangeEncoder(&globalMapKeyEncoder)
	ctx.Stream.WriteDecimalFloat(value)
}
func (_this *mapValueEncoder) EncodeBigDecimalFloat(ctx *EncoderContext, value *apd.Decimal) {
	ctx.WriteIndent()
	ctx.ChangeEncoder(&globalMapKeyEncoder)
	ctx.Stream.WriteBigDecimalFloat(value)
}
func (_this *mapValueEncoder) EncodeNan(ctx *EncoderContext, signaling bool) {
	ctx.WriteIndent()
	ctx.ChangeEncoder(&globalMapKeyEncoder)
	ctx.Stream.WriteNan(signaling)
}
func (_this *mapValueEncoder) EncodeTime(ctx *EncoderContext, value time.Time) {
	ctx.WriteIndent()
	ctx.ChangeEncoder(&globalMapKeyEncoder)
	ctx.Stream.WriteTime(value)
}
func (_this *mapValueEncoder) EncodeCompactTime(ctx *EncoderContext, value compact_time.Time) {
	ctx.WriteIndent()
	ctx.ChangeEncoder(&globalMapKeyEncoder)
	ctx.Stream.WriteCompactTime(value)
}
func (_this *mapValueEncoder) EncodeUUID(ctx *EncoderContext, value []byte) {
	ctx.WriteIndent()
	ctx.ChangeEncoder(&globalMapKeyEncoder)
	ctx.Stream.WriteUUID(value)
}
func (_this *mapValueEncoder) BeginList(ctx *EncoderContext) {
	ctx.WriteIndent()
	ctx.ChangeEncoder(&globalMapKeyEncoder)
	ctx.BeginList()
}
func (_this *mapValueEncoder) BeginMap(ctx *EncoderContext) {
	ctx.WriteIndent()
	ctx.ChangeEncoder(&globalMapKeyEncoder)
	ctx.BeginMap()
}
func (_this *mapValueEncoder) BeginMarkup(ctx *EncoderContext) {
	ctx.WriteIndent()
	ctx.ChangeEncoder(&globalMapKeyEncoder)
	ctx.BeginMarkup()
}
func (_this *mapValueEncoder) BeginMetadata(ctx *EncoderContext) {
	ctx.WriteIndent()
	ctx.ChangeEncoder(&globalMapKeyEncoder)
	ctx.BeginMetadata()
}
func (_this *mapValueEncoder) BeginComment(ctx *EncoderContext) {
	ctx.WriteIndent()
	ctx.ChangeEncoder(&globalMapKeyEncoder)
	ctx.BeginComment()
}
func (_this *mapValueEncoder) BeginMarker(ctx *EncoderContext) {
	ctx.WriteIndent()
	ctx.ChangeEncoder(&globalMapKeyEncoder)
	ctx.BeginMarker()
}
func (_this *mapValueEncoder) BeginReference(ctx *EncoderContext) {
	ctx.WriteIndent()
	ctx.ChangeEncoder(&globalMapKeyEncoder)
	ctx.BeginReference()
}
func (_this *mapValueEncoder) BeginConcatenate(ctx *EncoderContext) {
	ctx.WriteIndent()
	ctx.ChangeEncoder(&globalMapKeyEncoder)
	ctx.BeginConcatenate()
}
func (_this *mapValueEncoder) BeginConstant(ctx *EncoderContext, name []byte, explicitValue bool) {
	ctx.WriteIndent()
	ctx.ChangeEncoder(&globalMapKeyEncoder)
	ctx.BeginConstant(name, explicitValue)
}
func (_this *mapValueEncoder) BeginNA(ctx *EncoderContext) {
	ctx.WriteIndent()
	ctx.ChangeEncoder(&globalMapKeyEncoder)
	ctx.BeginNA()
}
func (_this *mapValueEncoder) EncodeArray(ctx *EncoderContext, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	ctx.WriteIndent()
	ctx.ChangeEncoder(&globalMapKeyEncoder)
	ctx.Stream.WriteArray(arrayType, elementCount, data)
}
func (_this *mapValueEncoder) EncodeStringlikeArray(ctx *EncoderContext, arrayType events.ArrayType, data string) {
	ctx.WriteIndent()
	ctx.ChangeEncoder(&globalMapKeyEncoder)
	ctx.Stream.WriteStringlikeArray(arrayType, data)
}
func (_this *mapValueEncoder) BeginArray(ctx *EncoderContext, arrayType events.ArrayType) {
	ctx.WriteIndent()
	ctx.ChangeEncoder(&globalMapKeyEncoder)
	ctx.BeginArray(arrayType)
}
