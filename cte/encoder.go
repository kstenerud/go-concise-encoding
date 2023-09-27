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
	"fmt"
	"io"
	"math/big"

	"github.com/cockroachdb/apd/v2"
	compact_float "github.com/kstenerud/go-compact-float"
	compact_time "github.com/kstenerud/go-compact-time"
	"github.com/kstenerud/go-concise-encoding/ce/events"
	"github.com/kstenerud/go-concise-encoding/configuration"
)

type EncoderEventReceiver struct {
	context EncoderContext
}

// Create a new encoder.
// If config = nil, defaults are used.
func NewEncoder(config *configuration.Configuration) *EncoderEventReceiver {
	_this := &EncoderEventReceiver{}
	_this.Init(config)
	return _this
}

// Initialize an encoder.
// If config = nil, defaults are used.
func (_this *EncoderEventReceiver) Init(config *configuration.Configuration) {
	_this.context.Init(config)
}

// Prepare the encoder for encoding. All events will be encoded to writer.
// PrepareToEncode MUST be called before using the encoder.
func (_this *EncoderEventReceiver) PrepareToEncode(writer io.Writer) {
	_this.context.Stream.SetWriter(writer)
}

func (_this *EncoderEventReceiver) OnBeginDocument() {
	_this.context.Begin()
	_this.context.Stream.WriteByteNotLF('c')
}

func (_this *EncoderEventReceiver) OnVersion(version uint64) {
	_this.context.Stream.WriteVersion(version)
	_this.context.WriteNewlineAndOriginAndIndent()
}

func (_this *EncoderEventReceiver) OnPadding() {
	// CTE doesn't have padding, so do nothing.
}

func (_this *EncoderEventReceiver) OnComment(isMultiline bool, contents []byte) {
	_this.context.BeforeComment()
	_this.context.Stream.WriteCommentBegin(isMultiline)
	if isMultiline {
		_this.context.Stream.WriteBytesPossibleLF(contents)
	} else {
		_this.context.Stream.WriteBytesNotLF(contents)
	}
	_this.context.Stream.WriteCommentEnd(isMultiline)
	_this.context.AfterComment(isMultiline)
}

func (_this *EncoderEventReceiver) OnNull() {
	_this.context.BeforeValue()
	_this.context.Stream.WriteNull()
	_this.context.AfterValue()
}

func (_this *EncoderEventReceiver) OnBoolean(value bool) {
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
	_this.context.Stream.WritePositiveInt(value, 10)
	_this.context.AfterValue()
}

func (_this *EncoderEventReceiver) OnNegativeInt(value uint64) {
	_this.context.BeforeValue()
	_this.context.Stream.WriteNegativeInt(value, 10)
	_this.context.AfterValue()
}

func (_this *EncoderEventReceiver) OnInt(value int64) {
	_this.context.BeforeValue()
	_this.context.Stream.WriteInt(value, 10)
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

func (_this *EncoderEventReceiver) OnTime(value compact_time.Time) {
	_this.context.BeforeValue()
	_this.context.Stream.WriteTime(value)
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

func (_this *EncoderEventReceiver) OnMedia(mediaType string, data []byte) {
	_this.context.BeforeValue()
	_this.context.EncodeMedia(mediaType, data)
	_this.context.AfterValue()
}

func (_this *EncoderEventReceiver) OnCustomBinary(customType uint64, data []byte) {
	_this.context.BeforeValue()
	_this.context.EncodeCustomBinary(customType, data)
	_this.context.AfterValue()
}

func (_this *EncoderEventReceiver) OnCustomText(customType uint64, data string) {
	_this.context.BeforeValue()
	_this.context.EncodeCustomText(customType, data)
	_this.context.AfterValue()
}

func (_this *EncoderEventReceiver) OnArrayBegin(arrayType events.ArrayType) {
	_this.context.BeforeValue()
	_this.context.BeginArray(arrayType, func() { _this.context.AfterValue() })
}

func (_this *EncoderEventReceiver) OnMediaBegin(mediaType string) {
	_this.context.BeforeValue()
	_this.context.BeginMedia(mediaType, func() { _this.context.AfterValue() })
}

func (_this *EncoderEventReceiver) OnCustomBegin(arrayType events.ArrayType, customType uint64) {
	_this.context.BeforeValue()
	switch arrayType {
	case events.ArrayTypeCustomBinary:
		_this.context.BeginCustomBinary(customType, func() { _this.context.AfterValue() })
	case events.ArrayTypeCustomText:
		_this.context.BeginCustomText(customType, func() { _this.context.AfterValue() })
	default:
		panic(fmt.Errorf("BUG: EncoderEventReceiver.OnCustomBegin: Cannot handle type %v", arrayType))
	}
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

func (_this *EncoderEventReceiver) OnRecordType(id []byte) {
	_this.context.BeforeValue()
	_this.context.BeginContainer()
	_this.context.Stream.WriteRecordTypeBegin(id)
	_this.context.Indent()
	_this.context.Stack(recordTypeDecorator)
}

func (_this *EncoderEventReceiver) OnRecord(id []byte) {
	_this.context.BeforeValue()
	_this.context.BeginContainer()
	_this.context.Stream.WriteRecordBegin(id)
	_this.context.Indent()
	_this.context.Stack(recordDecorator)
}

func (_this *EncoderEventReceiver) OnEdge() {
	_this.context.BeforeValue()
	_this.context.BeginContainer()
	_this.context.Stream.WriteEdgeBegin()
	_this.context.Indent()
	_this.context.Stack(edgeDecorator)
}

func (_this *EncoderEventReceiver) OnNode() {
	_this.context.BeforeValue()
	_this.context.Stream.WriteNodeBegin()
	_this.context.Indent()
	_this.context.Stack(nodeValueDecorator)
}

func (_this *EncoderEventReceiver) OnEndContainer() {
	_this.context.EndContainer()
}

func (_this *EncoderEventReceiver) OnMarker(id []byte) {
	_this.context.BeforeValue()
	_this.context.Stream.WriteMarkerBegin(id)
	_this.context.Stack(concatDecorator)
}

func (_this *EncoderEventReceiver) OnReferenceLocal(id []byte) {
	_this.context.BeforeValue()
	_this.context.Stream.WriteLocalReference(id)
	_this.context.AfterValue()
}

func (_this *EncoderEventReceiver) OnRemoteReference() {
	_this.context.BeforeValue()
	_this.context.Stream.WriteRemoteReferenceBegin()
	_this.context.Stack(concatDecorator)
}

func (_this *EncoderEventReceiver) OnEndDocument() {
	// Nothing to do
}

func (_this *EncoderEventReceiver) OnError() {}
