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
	"io"
	"math/big"
	"time"

	"github.com/kstenerud/go-concise-encoding/events"
	"github.com/kstenerud/go-concise-encoding/options"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
)

type Encoder interface {
	Begin(ctx *EncoderContext)
	End(ctx *EncoderContext)
	ChildContainerFinished(ctx *EncoderContext, isVisibleChild bool)
	EncodeNil(ctx *EncoderContext)
	BeginNA(ctx *EncoderContext)
	EncodeBool(ctx *EncoderContext, value bool)
	EncodeTrue(ctx *EncoderContext)
	EncodeFalse(ctx *EncoderContext)
	EncodePositiveInt(ctx *EncoderContext, value uint64)
	EncodeNegativeInt(ctx *EncoderContext, value uint64)
	EncodeInt(ctx *EncoderContext, value int64)
	EncodeBigInt(ctx *EncoderContext, value *big.Int)
	EncodeFloat(ctx *EncoderContext, value float64)
	EncodeBigFloat(ctx *EncoderContext, value *big.Float)
	EncodeDecimalFloat(ctx *EncoderContext, value compact_float.DFloat)
	EncodeBigDecimalFloat(ctx *EncoderContext, value *apd.Decimal)
	EncodeNan(ctx *EncoderContext, signaling bool)
	EncodeTime(ctx *EncoderContext, value time.Time)
	EncodeCompactTime(ctx *EncoderContext, value compact_time.Time)
	EncodeUUID(ctx *EncoderContext, value []byte)
	EncodeIdentifier(ctx *EncoderContext, value []byte)
	BeginList(ctx *EncoderContext)
	BeginMap(ctx *EncoderContext)
	BeginMarkup(ctx *EncoderContext)
	BeginComment(ctx *EncoderContext)
	BeginMarker(ctx *EncoderContext)
	BeginReference(ctx *EncoderContext)
	BeginConstant(ctx *EncoderContext, name []byte, explicitValue bool)
	EncodeArray(ctx *EncoderContext, arrayType events.ArrayType, elementCount uint64, data []uint8)
	EncodeStringlikeArray(ctx *EncoderContext, arrayType events.ArrayType, data string)
	BeginArray(ctx *EncoderContext, arrayType events.ArrayType)
	BeginArrayChunk(ctx *EncoderContext, elementCount uint64, moreChunksFollow bool)
	EncodeArrayData(ctx *EncoderContext, data []byte)
}

type EncoderEventReceiver struct {
	context EncoderContext
}

// Create a new encoder.
// If opts = nil, defaults are used.
func NewEncoder(opts *options.CTEEncoderOptions) *EncoderEventReceiver {
	_this := &EncoderEventReceiver{}
	_this.Init(opts)
	return _this
}

// Initialize an encoder.
// If opts = nil, defaults are used.
func (_this *EncoderEventReceiver) Init(opts *options.CTEEncoderOptions) {
	opts = opts.WithDefaultsApplied()
	_this.context.Init(opts)
}

// Reset the encoder back to its initial state.
func (_this *EncoderEventReceiver) Reset() {
	_this.context.Reset()
}

// Prepare the encoder for encoding. All events will be encoded to writer.
// PrepareToEncode MUST be called before using the encoder.
func (_this *EncoderEventReceiver) PrepareToEncode(writer io.Writer) {
	_this.context.Stream.SetWriter(writer)
}

func (_this *EncoderEventReceiver) OnBeginDocument() {
	_this.context.Reset()
}

func (_this *EncoderEventReceiver) OnVersion(version uint64) {
	_this.context.Stream.WriteVersion(version)
}

func (_this *EncoderEventReceiver) OnPadding(count int) {
	// Nothing to do
}

func (_this *EncoderEventReceiver) OnNil() {
	_this.context.CurrentEncoder.EncodeNil(&_this.context)
}

func (_this *EncoderEventReceiver) OnNA() {
	_this.context.CurrentEncoder.BeginNA(&_this.context)
}

func (_this *EncoderEventReceiver) OnBool(value bool) {
	_this.context.CurrentEncoder.EncodeBool(&_this.context, value)
}

func (_this *EncoderEventReceiver) OnTrue() {
	_this.context.CurrentEncoder.EncodeTrue(&_this.context)
}

func (_this *EncoderEventReceiver) OnFalse() {
	_this.context.CurrentEncoder.EncodeFalse(&_this.context)
}

func (_this *EncoderEventReceiver) OnPositiveInt(value uint64) {
	_this.context.CurrentEncoder.EncodePositiveInt(&_this.context, value)
}

func (_this *EncoderEventReceiver) OnNegativeInt(value uint64) {
	_this.context.CurrentEncoder.EncodeNegativeInt(&_this.context, value)
}

func (_this *EncoderEventReceiver) OnInt(value int64) {
	_this.context.CurrentEncoder.EncodeInt(&_this.context, value)
}

func (_this *EncoderEventReceiver) OnBigInt(value *big.Int) {
	_this.context.CurrentEncoder.EncodeBigInt(&_this.context, value)
}

func (_this *EncoderEventReceiver) OnFloat(value float64) {
	_this.context.CurrentEncoder.EncodeFloat(&_this.context, value)
}

func (_this *EncoderEventReceiver) OnBigFloat(value *big.Float) {
	_this.context.CurrentEncoder.EncodeBigFloat(&_this.context, value)
}

func (_this *EncoderEventReceiver) OnDecimalFloat(value compact_float.DFloat) {
	_this.context.CurrentEncoder.EncodeDecimalFloat(&_this.context, value)
}

func (_this *EncoderEventReceiver) OnBigDecimalFloat(value *apd.Decimal) {
	_this.context.CurrentEncoder.EncodeBigDecimalFloat(&_this.context, value)
}

func (_this *EncoderEventReceiver) OnNan(signaling bool) {
	_this.context.CurrentEncoder.EncodeNan(&_this.context, signaling)
}

func (_this *EncoderEventReceiver) OnUUID(value []byte) {
	_this.context.CurrentEncoder.EncodeUUID(&_this.context, value)
}

func (_this *EncoderEventReceiver) OnIdentifier(value []byte) {
	_this.context.CurrentEncoder.EncodeIdentifier(&_this.context, value)
}

func (_this *EncoderEventReceiver) OnTime(value time.Time) {
	_this.context.CurrentEncoder.EncodeTime(&_this.context, value)
}

func (_this *EncoderEventReceiver) OnCompactTime(value compact_time.Time) {
	_this.context.CurrentEncoder.EncodeCompactTime(&_this.context, value)
}

func (_this *EncoderEventReceiver) OnArray(arrayType events.ArrayType, elementCount uint64, value []byte) {
	_this.context.CurrentEncoder.EncodeArray(&_this.context, arrayType, elementCount, value)
}

func (_this *EncoderEventReceiver) OnStringlikeArray(arrayType events.ArrayType, value string) {
	_this.context.CurrentEncoder.EncodeStringlikeArray(&_this.context, arrayType, value)
}

func (_this *EncoderEventReceiver) OnArrayBegin(arrayType events.ArrayType) {
	_this.context.CurrentEncoder.BeginArray(&_this.context, arrayType)
}

func (_this *EncoderEventReceiver) OnArrayChunk(elementCount uint64, moreChunksFollow bool) {
	_this.context.CurrentEncoder.BeginArrayChunk(&_this.context, elementCount, moreChunksFollow)
}

func (_this *EncoderEventReceiver) OnArrayData(data []byte) {
	_this.context.CurrentEncoder.EncodeArrayData(&_this.context, data)
}

func (_this *EncoderEventReceiver) OnList() {
	_this.context.CurrentEncoder.BeginList(&_this.context)
}

func (_this *EncoderEventReceiver) OnMap() {
	_this.context.CurrentEncoder.BeginMap(&_this.context)
}

func (_this *EncoderEventReceiver) OnMarkup() {
	_this.context.CurrentEncoder.BeginMarkup(&_this.context)
}

func (_this *EncoderEventReceiver) OnComment() {
	_this.context.CurrentEncoder.BeginComment(&_this.context)
}

func (_this *EncoderEventReceiver) OnEnd() {
	_this.context.CurrentEncoder.End(&_this.context)
}

func (_this *EncoderEventReceiver) OnMarker() {
	_this.context.CurrentEncoder.BeginMarker(&_this.context)
}

func (_this *EncoderEventReceiver) OnReference() {
	_this.context.CurrentEncoder.BeginReference(&_this.context)
}

func (_this *EncoderEventReceiver) OnConstant(name []byte, explicitValue bool) {
	_this.context.CurrentEncoder.BeginConstant(&_this.context, name, explicitValue)
}

func (_this *EncoderEventReceiver) OnEndDocument() {
	// Nothing to do
}
