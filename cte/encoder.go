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
	"math"
	"math/big"
	"time"

	"github.com/kstenerud/go-concise-encoding/internal/common"

	"github.com/kstenerud/go-concise-encoding/events"
	"github.com/kstenerud/go-concise-encoding/options"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
)

type Encoder interface {
	Begin(ctx *EncoderContext)
	End(ctx *EncoderContext)
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
	BeginList(ctx *EncoderContext)
	BeginMap(ctx *EncoderContext)
	BeginMarkup(ctx *EncoderContext)
	BeginMetadata(ctx *EncoderContext)
	BeginComment(ctx *EncoderContext)
	BeginMarker(ctx *EncoderContext)
	BeginReference(ctx *EncoderContext)
	BeginConcatenate(ctx *EncoderContext)
	BeginConstant(ctx *EncoderContext, name []byte, explicitValue bool)
	BeginNA(ctx *EncoderContext)
	EncodeArray(ctx *EncoderContext, arrayType events.ArrayType, elementCount uint64, data []uint8)
	EncodeStringlikeArray(ctx *EncoderContext, arrayType events.ArrayType, data string)
	BeginArray(ctx *EncoderContext, arrayType events.ArrayType)
	BeginArrayChunk(ctx *EncoderContext, length uint64, moreChunksFollow bool)
	EncodeArrayData(ctx *EncoderContext, data []byte)
}

type RootEncoder struct {
	context EncoderContext
	buffer  *EncodeBuffer
}

// Create a new encoder.
// If opts = nil, defaults are used.
func NewEncoder(opts *options.CTEEncoderOptions) *RootEncoder {
	_this := &RootEncoder{}
	_this.Init(opts)
	return _this
}

// Initialize an encoder.
// If opts = nil, defaults are used.
func (_this *RootEncoder) Init(opts *options.CTEEncoderOptions) {
	opts = opts.WithDefaultsApplied()
	_this.context.Init(opts)
}

// Reset the encoder back to its initial state.
func (_this *RootEncoder) Reset() {
	_this.context.Reset()
}

// Prepare the encoder for encoding. All events will be encoded to writer.
// PrepareToEncode MUST be called before using the encoder.
func (_this *RootEncoder) PrepareToEncode(writer io.Writer) {
	_this.context.Stream.SetWriter(writer)
}

func (_this *RootEncoder) OnBeginDocument() {
	// TODO: Reset?
}

func (_this *RootEncoder) OnVersion(version uint64) {
	_this.context.Stream.WriteVersion(version)
}

func (_this *RootEncoder) OnPadding(count int) {
	// Nothing to do
}

func (_this *RootEncoder) OnNA() {
	_this.context.CurrentEncoder.BeginNA(&_this.context)
}

func (_this *RootEncoder) OnBool(value bool) {
	_this.context.CurrentEncoder.EncodeBool(&_this.context, value)
}

func (_this *RootEncoder) OnTrue() {
	_this.context.CurrentEncoder.EncodeTrue(&_this.context)
}

func (_this *RootEncoder) OnFalse() {
	_this.context.CurrentEncoder.EncodeFalse(&_this.context)
}

func (_this *RootEncoder) OnPositiveInt(value uint64) {
	_this.context.CurrentEncoder.EncodePositiveInt(&_this.context, value)
}

func (_this *RootEncoder) OnNegativeInt(value uint64) {
	_this.context.CurrentEncoder.EncodeNegativeInt(&_this.context, value)
}

func (_this *RootEncoder) OnInt(value int64) {
	_this.context.CurrentEncoder.EncodeInt(&_this.context, value)
}

func (_this *RootEncoder) OnBigInt(value *big.Int) {
	if value == nil {
		_this.OnNA()
		return
	}

	_this.context.CurrentEncoder.EncodeBigInt(&_this.context, value)
}

func (_this *RootEncoder) OnFloat(value float64) {
	if math.IsNaN(value) {
		_this.OnNan(common.IsSignalingNan(value))
		return
	}

	_this.context.CurrentEncoder.EncodeFloat(&_this.context, value)
}

func (_this *RootEncoder) OnBigFloat(value *big.Float) {
	if value == nil {
		_this.OnNA()
		return
	}

	_this.context.CurrentEncoder.EncodeBigFloat(&_this.context, value)
}

func (_this *RootEncoder) OnDecimalFloat(value compact_float.DFloat) {
	if value.IsNan() {
		_this.OnNan(value.IsSignalingNan())
		return
	}

	_this.context.CurrentEncoder.EncodeDecimalFloat(&_this.context, value)
}

func (_this *RootEncoder) OnBigDecimalFloat(value *apd.Decimal) {
	if value == nil {
		_this.OnNA()
		return
	}
	if value.Form == apd.NaNSignaling {
		_this.OnNan(true)
		return
	}
	if value.Form == apd.NaN {
		_this.OnNan(false)
		return
	}

	_this.context.CurrentEncoder.EncodeBigDecimalFloat(&_this.context, value)
}

func (_this *RootEncoder) OnNan(signaling bool) {
	_this.context.CurrentEncoder.EncodeNan(&_this.context, signaling)
}

func (_this *RootEncoder) OnUUID(value []byte) {
	_this.context.CurrentEncoder.EncodeUUID(&_this.context, value)
}

func (_this *RootEncoder) OnTime(value time.Time) {
	_this.context.CurrentEncoder.EncodeTime(&_this.context, value)
}

func (_this *RootEncoder) OnCompactTime(value compact_time.Time) {
	_this.context.CurrentEncoder.EncodeCompactTime(&_this.context, value)
}

func (_this *RootEncoder) OnArray(arrayType events.ArrayType, elementCount uint64, value []byte) {
	_this.context.CurrentEncoder.EncodeArray(&_this.context, arrayType, elementCount, value)
}

func (_this *RootEncoder) OnStringlikeArray(arrayType events.ArrayType, value string) {
	_this.context.CurrentEncoder.EncodeStringlikeArray(&_this.context, arrayType, value)
}

func (_this *RootEncoder) OnArrayBegin(arrayType events.ArrayType) {
	// TODO: Stack array prefixer
}

func (_this *RootEncoder) OnArrayChunk(length uint64, moreChunksFollow bool) {
	// TODO: Chunking
}

func (_this *RootEncoder) OnArrayData(data []byte) {
	// _this.context.CurrentPrefixer.RenderArrayPortion(&_this.context, value)
	// TODO: Detect end of array, and unstack
}

func (_this *RootEncoder) OnConcatenate() {
	_this.context.BeginConcatenate()
}

func (_this *RootEncoder) OnList() {
	_this.context.BeginList()
}

func (_this *RootEncoder) OnMap() {
	_this.context.BeginMap()
}

func (_this *RootEncoder) OnMarkup() {
	_this.context.BeginMarkup()
}

func (_this *RootEncoder) OnMetadata() {
	_this.context.BeginMetadata()
}

func (_this *RootEncoder) OnComment() {
	_this.context.BeginComment()
}

func (_this *RootEncoder) OnEnd() {
	_this.context.EndContainer()
}

func (_this *RootEncoder) OnMarker() {
	_this.context.BeginMarker()
}

func (_this *RootEncoder) OnReference() {
	_this.context.BeginReference()
}

func (_this *RootEncoder) OnConstant(name []byte, explicitValue bool) {
	// _this.context.CurrentEncoder.OnConstant(name, explicitValue)
}

func (_this *RootEncoder) OnEndDocument() {
	_this.context.Stream.Flush()
}
