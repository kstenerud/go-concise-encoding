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

type naEncoder struct{}

var globalNAEncoder naEncoder

func (_this *naEncoder) EncodeBool(ctx *EncoderContext, value bool) {
	ctx.Stream.WriteSeparator()
	ctx.unstack()
	ctx.Stream.WriteBool(value)
}
func (_this *naEncoder) EncodeTrue(ctx *EncoderContext) {
	ctx.Stream.WriteSeparator()
	ctx.unstack()
	ctx.Stream.WriteTrue()
}
func (_this *naEncoder) EncodeFalse(ctx *EncoderContext) {
	ctx.Stream.WriteSeparator()
	ctx.unstack()
	ctx.Stream.WriteFalse()
}
func (_this *naEncoder) EncodePositiveInt(ctx *EncoderContext, value uint64) {
	ctx.Stream.WriteSeparator()
	ctx.unstack()
	ctx.Stream.WritePositiveInt(value)
}
func (_this *naEncoder) EncodeNegativeInt(ctx *EncoderContext, value uint64) {
	ctx.Stream.WriteSeparator()
	ctx.unstack()
	ctx.Stream.WriteNegativeInt(value)
}
func (_this *naEncoder) EncodeInt(ctx *EncoderContext, value int64) {
	ctx.Stream.WriteSeparator()
	ctx.unstack()
	ctx.Stream.WriteInt(value)
}
func (_this *naEncoder) EncodeBigInt(ctx *EncoderContext, value *big.Int) {
	ctx.Stream.WriteSeparator()
	ctx.unstack()
	ctx.Stream.WriteBigInt(value)
}
func (_this *naEncoder) EncodeFloat(ctx *EncoderContext, value float64) {
	ctx.Stream.WriteSeparator()
	ctx.unstack()
	ctx.Stream.WriteFloat(value)
}
func (_this *naEncoder) EncodeBigFloat(ctx *EncoderContext, value *big.Float) {
	ctx.Stream.WriteSeparator()
	ctx.unstack()
	ctx.Stream.WriteBigFloat(value)
}
func (_this *naEncoder) EncodeDecimalFloat(ctx *EncoderContext, value compact_float.DFloat) {
	ctx.Stream.WriteSeparator()
	ctx.unstack()
	ctx.Stream.WriteDecimalFloat(value)
}
func (_this *naEncoder) EncodeBigDecimalFloat(ctx *EncoderContext, value *apd.Decimal) {
	ctx.Stream.WriteSeparator()
	ctx.unstack()
	ctx.Stream.WriteBigDecimalFloat(value)
}
func (_this *naEncoder) EncodeNan(ctx *EncoderContext, signaling bool) {
	ctx.Stream.WriteSeparator()
	ctx.unstack()
	ctx.Stream.WriteNan(signaling)
}
func (_this *naEncoder) EncodeTime(ctx *EncoderContext, value time.Time) {
	ctx.Stream.WriteSeparator()
	ctx.unstack()
	ctx.Stream.WriteTime(value)
}
func (_this *naEncoder) EncodeCompactTime(ctx *EncoderContext, value compact_time.Time) {
	ctx.Stream.WriteSeparator()
	ctx.unstack()
	ctx.Stream.WriteCompactTime(value)
}
func (_this *naEncoder) EncodeUUID(ctx *EncoderContext, value []byte) {
	ctx.Stream.WriteSeparator()
	ctx.unstack()
	ctx.Stream.WriteUUID(value)
}
func (_this *naEncoder) BeginList(ctx *EncoderContext) {
	ctx.Stream.WriteSeparator()
	ctx.unstack()
	ctx.BeginList()
}
func (_this *naEncoder) BeginMap(ctx *EncoderContext) {
	ctx.Stream.WriteSeparator()
	ctx.unstack()
	ctx.BeginMap()
}
func (_this *naEncoder) BeginMarkup(ctx *EncoderContext) {
	ctx.Stream.WriteSeparator()
	ctx.unstack()
	ctx.BeginMarkup()
}
func (_this *naEncoder) BeginMarker(ctx *EncoderContext) {
	ctx.Stream.WriteSeparator()
	ctx.unstack()
	ctx.BeginMarker()
}
func (_this *naEncoder) BeginReference(ctx *EncoderContext) {
	ctx.Stream.WriteSeparator()
	ctx.unstack()
	ctx.BeginReference()
}
func (_this *naEncoder) BeginConcatenate(ctx *EncoderContext) {
	ctx.Stream.WriteSeparator()
	ctx.unstack()
	ctx.BeginConcatenate()
}
func (_this *naEncoder) BeginConstant(ctx *EncoderContext, name []byte, explicitValue bool) {
	ctx.Stream.WriteSeparator()
	ctx.unstack()
	ctx.BeginConstant(name, explicitValue)
}
func (_this *naEncoder) BeginNA(ctx *EncoderContext) {
	// Only unstack
	ctx.Stream.WriteSeparator() // TODO: Remove me
	ctx.unstack()
	ctx.Stream.WriteNA() // TODO: Remove me
}
func (_this *naEncoder) EncodeArray(ctx *EncoderContext, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	ctx.Stream.WriteSeparator()
	ctx.unstack()
	ctx.Stream.WriteArray(arrayType, elementCount, data)
}
func (_this *naEncoder) EncodeStringlikeArray(ctx *EncoderContext, arrayType events.ArrayType, data string) {
	ctx.Stream.WriteSeparator()
	ctx.unstack()
	ctx.Stream.WriteStringlikeArray(arrayType, data)
}
func (_this *naEncoder) BeginArray(ctx *EncoderContext, arrayType events.ArrayType) {
	ctx.Stream.WriteSeparator()
	ctx.unstack()
	ctx.BeginArray(arrayType)
}
