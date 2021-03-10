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

type topLevelEncoder struct{}

var globalTopLevelEncoder topLevelEncoder

func (_this *topLevelEncoder) String() string { return "topLevelEncoder" }

func (_this *topLevelEncoder) ChildContainerFinished(ctx *EncoderContext, isVisibleChild bool) {
	// Nothing to do
}

func (_this *topLevelEncoder) EncodeNA(ctx *EncoderContext) {
	ctx.WriteCurrentPrefix()
	ctx.Stream.WriteNA()
}
func (_this *topLevelEncoder) EncodeBool(ctx *EncoderContext, value bool) {
	ctx.WriteCurrentPrefix()
	ctx.Stream.WriteBool(value)
}
func (_this *topLevelEncoder) EncodeTrue(ctx *EncoderContext) {
	ctx.WriteCurrentPrefix()
	ctx.Stream.WriteTrue()
}
func (_this *topLevelEncoder) EncodeFalse(ctx *EncoderContext) {
	ctx.WriteCurrentPrefix()
	ctx.Stream.WriteFalse()
}
func (_this *topLevelEncoder) EncodePositiveInt(ctx *EncoderContext, value uint64) {
	ctx.WriteCurrentPrefix()
	ctx.Stream.WritePositiveInt(value)
}
func (_this *topLevelEncoder) EncodeNegativeInt(ctx *EncoderContext, value uint64) {
	ctx.WriteCurrentPrefix()
	ctx.Stream.WriteNegativeInt(value)
}
func (_this *topLevelEncoder) EncodeInt(ctx *EncoderContext, value int64) {
	ctx.WriteCurrentPrefix()
	ctx.Stream.WriteInt(value)
}
func (_this *topLevelEncoder) EncodeBigInt(ctx *EncoderContext, value *big.Int) {
	ctx.WriteCurrentPrefix()
	ctx.Stream.WriteBigInt(value)
}
func (_this *topLevelEncoder) EncodeFloat(ctx *EncoderContext, value float64) {
	ctx.WriteCurrentPrefix()
	ctx.Stream.WriteFloat(value)
}
func (_this *topLevelEncoder) EncodeBigFloat(ctx *EncoderContext, value *big.Float) {
	ctx.WriteCurrentPrefix()
	ctx.Stream.WriteBigFloat(value)
}
func (_this *topLevelEncoder) EncodeDecimalFloat(ctx *EncoderContext, value compact_float.DFloat) {
	ctx.WriteCurrentPrefix()
	ctx.Stream.WriteDecimalFloat(value)
}
func (_this *topLevelEncoder) EncodeBigDecimalFloat(ctx *EncoderContext, value *apd.Decimal) {
	ctx.WriteCurrentPrefix()
	ctx.Stream.WriteBigDecimalFloat(value)
}
func (_this *topLevelEncoder) EncodeNan(ctx *EncoderContext, signaling bool) {
	ctx.WriteCurrentPrefix()
	ctx.Stream.WriteNan(signaling)
}
func (_this *topLevelEncoder) EncodeTime(ctx *EncoderContext, value time.Time) {
	ctx.WriteCurrentPrefix()
	ctx.Stream.WriteTime(value)
}
func (_this *topLevelEncoder) EncodeCompactTime(ctx *EncoderContext, value compact_time.Time) {
	ctx.WriteCurrentPrefix()
	ctx.Stream.WriteCompactTime(value)
}
func (_this *topLevelEncoder) EncodeUUID(ctx *EncoderContext, value []byte) {
	ctx.WriteCurrentPrefix()
	ctx.Stream.WriteUUID(value)
}
func (_this *topLevelEncoder) BeginList(ctx *EncoderContext) {
	ctx.WriteCurrentPrefix()
	ctx.BeginStandardList()
}
func (_this *topLevelEncoder) BeginMap(ctx *EncoderContext) {
	ctx.WriteCurrentPrefix()
	ctx.BeginStandardMap()
}
func (_this *topLevelEncoder) BeginMarkup(ctx *EncoderContext) {
	ctx.WriteCurrentPrefix()
	ctx.BeginStandardMarkup()
}
func (_this *topLevelEncoder) BeginMetadata(ctx *EncoderContext) {
	ctx.WriteCurrentPrefix()
	ctx.BeginStandardMetadata()
}
func (_this *topLevelEncoder) BeginComment(ctx *EncoderContext) {
	ctx.WriteCurrentPrefix()
	ctx.BeginStandardComment()
}
func (_this *topLevelEncoder) BeginMarker(ctx *EncoderContext) {
	ctx.WriteCurrentPrefix()
	ctx.BeginStandardMarker()
}
func (_this *topLevelEncoder) BeginReference(ctx *EncoderContext) {
	ctx.WriteCurrentPrefix()
	ctx.BeginStandardReference()
}
func (_this *topLevelEncoder) BeginConcatenate(ctx *EncoderContext) {
	ctx.WriteCurrentPrefix()
	ctx.BeginStandardConcatenate()
}
func (_this *topLevelEncoder) BeginConstant(ctx *EncoderContext, name []byte, explicitValue bool) {
	ctx.WriteCurrentPrefix()
	ctx.BeginStandardConstant(name, explicitValue)
}
func (_this *topLevelEncoder) EncodeArray(ctx *EncoderContext, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	ctx.WriteCurrentPrefix()
	ctx.ArrayEngine.EncodeArray(arrayType, elementCount, data)
}
func (_this *topLevelEncoder) EncodeStringlikeArray(ctx *EncoderContext, arrayType events.ArrayType, data string) {
	ctx.WriteCurrentPrefix()
	ctx.ArrayEngine.EncodeStringlikeArray(arrayType, data)
}
func (_this *topLevelEncoder) BeginArray(ctx *EncoderContext, arrayType events.ArrayType) {
	ctx.WriteCurrentPrefix()
	ctx.BeginStandardArray(arrayType)
}
