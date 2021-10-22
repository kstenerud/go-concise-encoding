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
}

// Prepare the encoder for encoding. All events will be encoded to writer.
// PrepareToEncode MUST be called before using the encoder.
func (_this *EncoderEventReceiver) PrepareToEncode(writer io.Writer) {
	_this.context.Stream.SetWriter(writer)
}

func (_this *EncoderEventReceiver) OnBeginDocument() {
	_this.context.Begin()
}

func (_this *EncoderEventReceiver) OnVersion(version uint64) {
	_this.context.Stream.WriteVersion(version)
	_this.context.WriteIndent()
}

func (_this *EncoderEventReceiver) OnPadding(count int) {
	// Nothing to do
}

func (_this *EncoderEventReceiver) OnComment(isMultiline bool, contents []byte) {
	_this.context.BeforeComment()
	_this.context.Stream.WriteCommentBegin(isMultiline)
	_this.context.Stream.WriteBytes(contents)
	_this.context.Stream.WriteCommentEnd(isMultiline)
	_this.context.AfterComment(isMultiline)
}

func (_this *EncoderEventReceiver) OnNil() {
	_this.context.BeforeValue()
	_this.context.Stream.WriteNil()
	_this.context.AfterValue()
}

func (_this *EncoderEventReceiver) OnNA() {
	_this.context.BeforeValue()
	_this.context.Stream.WriteNA()
	_this.context.Stack(concatDecorator)
}

func (_this *EncoderEventReceiver) OnBool(value bool) {
	_this.context.BeforeValue()
	_this.context.Stream.WriteBool(value)
	_this.context.AfterValue()
}

func (_this *EncoderEventReceiver) OnTrue() {
	_this.context.BeforeValue()
	_this.context.Stream.WriteTrue()
	_this.context.AfterValue()
}

func (_this *EncoderEventReceiver) OnFalse() {
	_this.context.BeforeValue()
	_this.context.Stream.WriteFalse()
	_this.context.AfterValue()
}

func (_this *EncoderEventReceiver) OnPositiveInt(value uint64) {
	_this.context.BeforeValue()
	_this.context.Stream.WritePositiveInt(value)
	_this.context.AfterValue()
}

func (_this *EncoderEventReceiver) OnNegativeInt(value uint64) {
	_this.context.BeforeValue()
	_this.context.Stream.WriteNegativeInt(value)
	_this.context.AfterValue()
}

func (_this *EncoderEventReceiver) OnInt(value int64) {
	_this.context.BeforeValue()
	_this.context.Stream.WriteInt(value)
	_this.context.AfterValue()
}

func (_this *EncoderEventReceiver) OnBigInt(value *big.Int) {
	_this.context.BeforeValue()
	_this.context.Stream.WriteBigInt(value)
	_this.context.AfterValue()
}

func (_this *EncoderEventReceiver) OnFloat(value float64) {
	_this.context.BeforeValue()
	_this.context.Stream.WriteFloat(value)
	_this.context.AfterValue()
}

func (_this *EncoderEventReceiver) OnBigFloat(value *big.Float) {
	_this.context.BeforeValue()
	_this.context.Stream.WriteBigFloat(value)
	_this.context.AfterValue()
}

func (_this *EncoderEventReceiver) OnDecimalFloat(value compact_float.DFloat) {
	_this.context.BeforeValue()
	_this.context.Stream.WriteDecimalFloat(value)
	_this.context.AfterValue()
}

func (_this *EncoderEventReceiver) OnBigDecimalFloat(value *apd.Decimal) {
	_this.context.BeforeValue()
	_this.context.Stream.WriteBigDecimalFloat(value)
	_this.context.AfterValue()
}

func (_this *EncoderEventReceiver) OnNan(signaling bool) {
	_this.context.BeforeValue()
	_this.context.Stream.WriteNan(signaling)
	_this.context.AfterValue()
}

func (_this *EncoderEventReceiver) OnUID(value []byte) {
	_this.context.BeforeValue()
	_this.context.Stream.WriteUID(value)
	_this.context.AfterValue()
}

func (_this *EncoderEventReceiver) OnTime(value time.Time) {
	_this.context.BeforeValue()
	_this.context.Stream.WriteTime(value)
	_this.context.AfterValue()
}

func (_this *EncoderEventReceiver) OnCompactTime(value compact_time.Time) {
	_this.context.BeforeValue()
	_this.context.Stream.WriteCompactTime(value)
	_this.context.AfterValue()
}

func (_this *EncoderEventReceiver) OnArray(arrayType events.ArrayType, elementCount uint64, value []byte) {
	_this.context.BeforeValue()
	_this.context.EncodeArray(arrayType, elementCount, value)
	_this.context.AfterValue()
}

func (_this *EncoderEventReceiver) OnStringlikeArray(arrayType events.ArrayType, value string) {
	_this.context.BeforeValue()
	_this.context.EncodeStringlikeArray(arrayType, value)
	_this.context.AfterValue()
}

func (_this *EncoderEventReceiver) OnArrayBegin(arrayType events.ArrayType) {
	_this.context.BeforeValue()
	_this.context.BeginArray(arrayType, func() { _this.context.AfterValue() })
}

func (_this *EncoderEventReceiver) OnArrayChunk(elementCount uint64, moreChunksFollow bool) {
	_this.context.BeginArrayChunk(elementCount, moreChunksFollow)
}

func (_this *EncoderEventReceiver) OnArrayData(data []byte) {
	_this.context.EncodeArrayData(data)
}

func (_this *EncoderEventReceiver) OnList() {
	_this.context.BeforeValue()
	_this.context.BeginContainer()
	_this.context.Stream.WriteListBegin()
	_this.context.Indent()
	_this.context.Stack(listDecorator)
}

func (_this *EncoderEventReceiver) OnMap() {
	_this.context.BeforeValue()
	_this.context.BeginContainer()
	_this.context.Stream.WriteMapBegin()
	_this.context.Indent()
	_this.context.Stack(mapKeyDecorator)
}

func (_this *EncoderEventReceiver) OnMarkup(id []byte) {
	_this.context.BeforeValue()
	_this.context.BeginContainer()
	_this.context.Stream.WriteMarkupBegin(id)
	_this.context.Indent()
	_this.context.Stack(markupKeyDecorator)
}

func (_this *EncoderEventReceiver) OnEdge() {
	_this.context.BeforeValue()
	_this.context.Stream.WriteEdgeBegin()
	_this.context.Stack(edgeSourceDecorator)
}

func (_this *EncoderEventReceiver) OnNode() {
	_this.context.BeforeValue()
	_this.context.Stream.WriteNodeBegin()
	_this.context.Stack(nodeValueDecorator)
}

func (_this *EncoderEventReceiver) OnEnd() {
	_this.context.EndContainer()
}

func (_this *EncoderEventReceiver) OnMarker(id []byte) {
	_this.context.BeforeValue()
	_this.context.Stream.WriteMarkerBegin(id)
	_this.context.Stack(concatDecorator)
}

func (_this *EncoderEventReceiver) OnReference(id []byte) {
	_this.context.BeforeValue()
	_this.context.Stream.WriteReference(id)
	_this.context.AfterValue()
}

func (_this *EncoderEventReceiver) OnRIDReference() {
	_this.context.BeforeValue()
	_this.context.Stream.WriteReferenceBegin()
	_this.context.Stack(concatDecorator)
}

func (_this *EncoderEventReceiver) OnConstant(name []byte) {
	_this.context.BeforeValue()
	_this.context.Stream.WriteConstant(name)
	_this.context.AfterValue()
}

func (_this *EncoderEventReceiver) OnEndDocument() {
	// Nothing to do
}
