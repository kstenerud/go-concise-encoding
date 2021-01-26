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
	"time"

	"github.com/kstenerud/go-concise-encoding/events"
	"github.com/kstenerud/go-concise-encoding/options"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
)

// Receives data events, constructing a CTE document from them.
//
// Note: This is a LOW LEVEL API. Error reporting is done via panics. Be sure
// to recover() at an appropriate location when calling Encoder's methods
// directly (with the exception of constructors and initializers, which are not
// designed to panic).
type OldEncoder struct {
	stream      EncodeBuffer
	engine      encoderEngine
	arrayEngine arrayEncoderEngine
	opts        options.CTEEncoderOptions
}

// Create a new CTE encoder, which will receive data events and write a document
// to writer. If opts is nil, default options will be used.
func NewOldEncoder(opts *options.CTEEncoderOptions) *OldEncoder {
	_this := &OldEncoder{}
	_this.Init(opts)
	return _this
}

// Initialize this encoder, which will receive data events and write a document
// to writer. If opts is nil, default options will be used.
func (_this *OldEncoder) Init(opts *options.CTEEncoderOptions) {
	opts = opts.WithDefaultsApplied()
	_this.opts = *opts
	_this.stream.Init(_this.opts.BufferSize)
	_this.engine.Init(&_this.stream, _this.opts.Indent)
	_this.arrayEngine.Init(&_this.engine, opts)
}

// Prepare the encoder for encoding. All events will be encoded to writer.
// PrepareToEncode MUST be called before using the encoder.
func (_this *OldEncoder) PrepareToEncode(writer io.Writer) {
	_this.stream.SetWriter(writer)
}

func (_this *OldEncoder) Reset() {
	_this.stream.Reset()
	_this.engine.Reset()
	_this.arrayEngine.Reset()
}

// ============================================================================

// DataEventReceiver

func (_this *OldEncoder) OnBeginDocument() {
	// Nothing to do
}

func (_this *OldEncoder) OnPadding(_ int) {
	// Nothing to do
}

func (_this *OldEncoder) OnVersion(version uint64) {
	_this.engine.AddVersion(version)
}

func (_this *OldEncoder) OnNA() {
	_this.engine.BeginObject()
	_this.stream.WriteNA()
	_this.engine.CompleteObject()
}

func (_this *OldEncoder) OnBool(value bool) {
	if value {
		_this.OnTrue()
	} else {
		_this.OnFalse()
	}
}

func (_this *OldEncoder) OnTrue() {
	_this.engine.BeginObject()
	_this.stream.WriteTrue()
	_this.engine.CompleteObject()
}

func (_this *OldEncoder) OnFalse() {
	_this.engine.BeginObject()
	_this.stream.WriteFalse()
	_this.engine.CompleteObject()
}

func (_this *OldEncoder) OnInt(value int64) {
	if value >= 0 {
		_this.OnPositiveInt(uint64(value))
	} else {
		_this.OnNegativeInt(uint64(-value))
	}
}

func (_this *OldEncoder) OnBigInt(value *big.Int) {
	_this.engine.BeginObject()
	_this.stream.WriteBigInt(value)
	_this.engine.CompleteObject()
}

func (_this *OldEncoder) OnPositiveInt(value uint64) {
	switch _this.engine.Awaiting {
	case awaitingMarkerID:
		_this.engine.CompleteMarker(value)
	case awaitingReferenceID:
		_this.engine.CompleteReference(value)
	default:
		_this.engine.BeginObject()
		_this.stream.WritePositiveInt(value)
		_this.engine.CompleteObject()
	}
}

func (_this *OldEncoder) OnNegativeInt(value uint64) {
	_this.engine.BeginObject()
	_this.stream.WriteNegativeInt(value)
	_this.engine.CompleteObject()
}

func (_this *OldEncoder) onInfinity(isPositive bool) {
	_this.engine.BeginObject()
	if isPositive {
		_this.stream.WritePosInfinity()
	} else {
		_this.stream.WriteNegInfinity()
	}
	_this.engine.CompleteObject()
}

func (_this *OldEncoder) OnFloat(value float64) {
	_this.engine.BeginObject()
	_this.stream.WriteFloat(value)
	_this.engine.CompleteObject()
}

func (_this *OldEncoder) OnBigFloat(value *big.Float) {
	_this.engine.BeginObject()
	_this.stream.WriteBigFloat(value)
	_this.engine.CompleteObject()
}

func (_this *OldEncoder) OnDecimalFloat(value compact_float.DFloat) {
	_this.engine.BeginObject()
	_this.stream.WriteDecimalFloat(value)
	_this.engine.CompleteObject()
}

func (_this *OldEncoder) OnBigDecimalFloat(value *apd.Decimal) {
	_this.engine.BeginObject()
	_this.stream.WriteBigDecimalFloat(value)
	_this.engine.CompleteObject()
}

func (_this *OldEncoder) OnNan(signaling bool) {
	_this.engine.BeginObject()
	if signaling {
		_this.stream.WriteSignalingNan()
	} else {
		_this.stream.WriteQuietNan()
	}
	_this.engine.CompleteObject()
}

func (_this *OldEncoder) OnUUID(value []byte) {
	_this.engine.BeginObject()
	_this.stream.WriteUUID(value)
	_this.engine.CompleteObject()
}

func (_this *OldEncoder) OnTime(value time.Time) {
	_this.engine.BeginObject()
	_this.stream.WriteTime(value)
	_this.engine.CompleteObject()
}

func (_this *OldEncoder) OnCompactTime(value *compact_time.Time) {
	_this.engine.BeginObject()
	_this.stream.WriteCompactTime(value)
	_this.engine.CompleteObject()
}

func (_this *OldEncoder) OnArray(arrayType events.ArrayType, elementCount uint64, value []byte) {
	_this.OnArrayBegin(arrayType)
	_this.OnArrayChunk(elementCount, false)
	if elementCount > 0 {
		_this.OnArrayData(value)
	}
}

func (_this *OldEncoder) OnArrayBegin(arrayType events.ArrayType) {
	_this.arrayEngine.OnArrayBegin(arrayType)
}

func (_this *OldEncoder) OnArrayChunk(elementCount uint64, moreChunksFollow bool) {
	_this.arrayEngine.OnArrayChunk(elementCount, moreChunksFollow)
}

func (_this *OldEncoder) OnArrayData(data []byte) {
	_this.arrayEngine.OnArrayData(data)
}

func (_this *OldEncoder) OnList() {
	_this.engine.BeginContainer(awaitingListFirstItem, "[")
}

func (_this *OldEncoder) OnMap() {
	_this.engine.BeginContainer(awaitingMapFirstKey, "{")
}

func (_this *OldEncoder) OnMarkup() {
	_this.engine.BeginMarkup()
}

func (_this *OldEncoder) OnMetadata() {
	_this.engine.BeginPseudoContainer(awaitingMetaFirstKey, "(")
}

func (_this *OldEncoder) OnComment() {
	_this.engine.BeginComment()
}

func (_this *OldEncoder) OnEnd() {
	_this.engine.EndContainer()
}

func (_this *OldEncoder) OnMarker() {
	_this.engine.BeginMarker()
}

func (_this *OldEncoder) OnReference() {
	_this.engine.BeginReference()
}

func (_this *OldEncoder) OnConcatenate() {
	panic("TODO: CTE Encoder.OnConcatenate")
}

func (_this *OldEncoder) OnConstant(name []byte, explicitValue bool) {
	_this.engine.BeginObject()
	_this.stream.AddByte('#')
	_this.stream.AddNonemptyBytes(name)
	if explicitValue {
		_this.stream.AddByte(':')
	} else {
		_this.engine.CompleteObject()
	}
}

func (_this *OldEncoder) OnEndDocument() {
	_this.stream.Flush()
}

func (_this *OldEncoder) unexpectedError(err error, encoding interface{}) {
	_this.errorf("unexpected error [%v] while encoding %v", err, encoding)
}

func (_this *OldEncoder) errorf(format string, args ...interface{}) {
	panic(fmt.Errorf(format, args...))
}
