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

type Encoder struct {
	context EncoderContext
	buffer  *EncodeBuffer
}

// Create a new encoder.
// If opts = nil, defaults are used.
func NewEncoder(opts *options.CTEEncoderOptions) *Encoder {
	_this := &Encoder{}
	_this.Init(opts)
	return _this
}

// Initialize an encoder.
// If opts = nil, defaults are used.
func (_this *Encoder) Init(opts *options.CTEEncoderOptions) {
	opts = opts.WithDefaultsApplied()
	_this.context.Init(opts)
}

// Reset the encoder back to its initial state.
func (_this *Encoder) Reset() {
	_this.context.Reset()
}

// Prepare the encoder for encoding. All events will be encoded to writer.
// PrepareToEncode MUST be called before using the encoder.
func (_this *Encoder) PrepareToEncode(writer io.Writer) {
	_this.context.Stream.SetWriter(writer)
}

func (_this *Encoder) OnBeginDocument() {
	// TODO: Reset?
}

func (_this *Encoder) OnVersion(version uint64) {
	_this.context.Stream.WriteVersion(version)
}

func (_this *Encoder) OnPadding(count int) {
	// Nothing to do
}

func (_this *Encoder) OnNA() {
	_this.context.ApplyPrefix()
	_this.context.Stream.WriteNA()
	_this.context.BeginNA()
}

// TODO: Need callback for child end

func (_this *Encoder) OnBool(value bool) {
	_this.context.ApplyPrefix()
	_this.context.Stream.WriteBool(value)
	_this.context.NotifyObject()
}

func (_this *Encoder) OnTrue() {
	_this.context.ApplyPrefix()
	_this.context.Stream.WriteTrue()
	_this.context.NotifyObject()
}

func (_this *Encoder) OnFalse() {
	_this.context.ApplyPrefix()
	_this.context.Stream.WriteFalse()
	_this.context.NotifyObject()
}

func (_this *Encoder) OnPositiveInt(value uint64) {
	_this.context.ApplyPrefix()
	_this.context.Stream.WritePositiveInt(value)
	_this.context.NotifyObject()
}

func (_this *Encoder) OnNegativeInt(value uint64) {
	_this.context.ApplyPrefix()
	_this.context.Stream.WriteNegativeInt(value)
	_this.context.NotifyObject()
}

func (_this *Encoder) OnInt(value int64) {
	_this.context.ApplyPrefix()
	_this.context.Stream.WriteInt(value)
	_this.context.NotifyObject()
}

func (_this *Encoder) OnBigInt(value *big.Int) {
	if value == nil {
		_this.OnNA()
		return
	}

	_this.context.ApplyPrefix()
	_this.context.Stream.WriteBigInt(value)
	_this.context.NotifyObject()
}

func (_this *Encoder) OnFloat(value float64) {
	if math.IsNaN(value) {
		_this.OnNan(common.IsSignalingNan(value))
		return
	}

	_this.context.ApplyPrefix()
	_this.context.Stream.WriteFloat(value)
	_this.context.NotifyObject()
}

func (_this *Encoder) OnBigFloat(value *big.Float) {
	if value == nil {
		_this.OnNA()
		return
	}

	_this.context.ApplyPrefix()
	_this.context.Stream.WriteBigFloat(value)
	_this.context.NotifyObject()
}

func (_this *Encoder) OnDecimalFloat(value compact_float.DFloat) {
	if value.IsNan() {
		_this.OnNan(value.IsSignalingNan())
		return
	}

	_this.context.ApplyPrefix()
	_this.context.Stream.WriteDecimalFloat(value)
	_this.context.NotifyObject()
}

func (_this *Encoder) OnBigDecimalFloat(value *apd.Decimal) {
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

	_this.context.ApplyPrefix()
	_this.context.Stream.WriteBigDecimalFloat(value)
	_this.context.NotifyObject()
}

func (_this *Encoder) OnNan(signaling bool) {
	_this.context.ApplyPrefix()
	_this.context.Stream.WriteNan(signaling)
	_this.context.NotifyObject()
}

func (_this *Encoder) OnUUID(value []byte) {
	_this.context.ApplyPrefix()
	_this.context.Stream.WriteUUID(value)
	_this.context.NotifyObject()
}

func (_this *Encoder) OnTime(value time.Time) {
	_this.context.ApplyPrefix()
	_this.context.Stream.WriteTime(value)
	_this.context.NotifyObject()
}

func (_this *Encoder) OnCompactTime(value compact_time.Time) {
	_this.context.ApplyPrefix()
	_this.context.Stream.WriteCompactTime(value)
	_this.context.NotifyObject()
}

func (_this *Encoder) OnArray(arrayType events.ArrayType, elementCount uint64, value []byte) {
	_this.context.ApplyPrefix()
	switch arrayType {
	case events.ArrayTypeString:
		// TODO: unquoted-safe
		// TODO: escapes
		_this.context.Stream.AddByte('"')
		_this.context.Stream.AddBytes(value)
		_this.context.Stream.AddByte('"')
	}
	_this.context.NotifyObject()
	_this.context.Stream.Flush()
}

func (_this *Encoder) OnArrayBegin(arrayType events.ArrayType) {
	// TODO: Stack array prefixer
}

func (_this *Encoder) OnArrayChunk(length uint64, moreChunksFollow bool) {
	// TODO: Chunking
}

func (_this *Encoder) OnArrayData(data []byte) {
	// _this.context.CurrentPrefixer.RenderArrayPortion(&_this.context, value)
	// TODO: Detect end of array, and unstack
}

func (_this *Encoder) OnConcatenate() {
	_this.context.BeginConcatenate()
}

func (_this *Encoder) OnList() {
	_this.context.BeginList()
}

func (_this *Encoder) OnMap() {
	_this.context.BeginMap()
}

func (_this *Encoder) OnMarkup() {
	_this.context.BeginMarkup()
}

func (_this *Encoder) OnMetadata() {
	_this.context.BeginMetadata()
}

func (_this *Encoder) OnComment() {
	_this.context.BeginComment()
}

func (_this *Encoder) OnEnd() {
	_this.context.EndContainer()
	_this.context.NotifyObject()
}

func (_this *Encoder) OnMarker() {
	_this.context.BeginMarker()
}

func (_this *Encoder) OnReference() {
	_this.context.BeginReference()
}

func (_this *Encoder) OnConstant(name []byte, explicitValue bool) {
	// _this.context.CurrentEncoder.OnConstant(name, explicitValue)
}

func (_this *Encoder) OnEndDocument() {
	// TODO: Do nothing?
}
