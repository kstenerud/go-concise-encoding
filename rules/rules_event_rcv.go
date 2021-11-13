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

package rules

import (
	"math"
	"math/big"
	"time"

	"github.com/cockroachdb/apd/v2"
	compact_float "github.com/kstenerud/go-compact-float"
	compact_time "github.com/kstenerud/go-compact-time"
	"github.com/kstenerud/go-concise-encoding/events"
	"github.com/kstenerud/go-concise-encoding/internal/common"
	"github.com/kstenerud/go-concise-encoding/options"
	"github.com/kstenerud/go-concise-encoding/version"
)

// RulesEventReceiver is a DataEventsReceiver passthrough object that constrains
// the order and contents of events to ensure that they form a valid and
// complete Concise Encoding document.
//
// Put this right after your event generator in the event receiver chain to
// enforce correctly formed documents.
//
// Note: This is a LOW LEVEL API. Error reporting is done via panics. Be sure
// to recover() at an appropriate location when calling this struct's methods
// directly (with the exception of constructors and initializers, which are not
// designed to panic).
type RulesEventReceiver struct {
	context  Context
	receiver events.DataEventReceiver
}

// Create a new rules set.
// If opts = nil, defaults are used.
func NewRules(nextReceiver events.DataEventReceiver, opts *options.RuleOptions) *RulesEventReceiver {
	_this := &RulesEventReceiver{}
	_this.Init(nextReceiver, opts)
	return _this
}

var nullReceiver = &events.NullEventReceiver{}

// Initialize a rules set.
// If opts = nil, defaults are used.
func (_this *RulesEventReceiver) Init(nextReceiver events.DataEventReceiver, opts *options.RuleOptions) {
	opts = opts.WithDefaultsApplied()
	_this.receiver = nextReceiver
	if _this.receiver == nil {
		_this.receiver = nullReceiver
	}
	_this.context.Init(version.ConciseEncodingVersion, opts)
}

// Reset the rules set back to its initial state.
func (_this *RulesEventReceiver) Reset() {
	_this.context.Reset()
}

func (_this *RulesEventReceiver) SetNextReceiver(nextReceiver events.DataEventReceiver) {
	_this.receiver = nextReceiver
}

func (_this *RulesEventReceiver) OnBeginDocument() {
	_this.context.CurrentEntry.Rule.OnBeginDocument(&_this.context)
	_this.receiver.OnBeginDocument()
}

func (_this *RulesEventReceiver) OnVersion(version uint64) {
	_this.context.CurrentEntry.Rule.OnVersion(&_this.context, version)
	_this.receiver.OnVersion(version)
}

func (_this *RulesEventReceiver) OnPadding(count int) {
	_this.context.CurrentEntry.Rule.OnPadding(&_this.context)
	_this.receiver.OnPadding(count)
}

func (_this *RulesEventReceiver) OnComment(isMultiline bool, contents []byte) {
	// TODO: Validate comment contents
	_this.context.CurrentEntry.Rule.OnComment(&_this.context)
	_this.receiver.OnComment(isMultiline, contents)
}

func (_this *RulesEventReceiver) OnNil() {
	_this.context.NotifyNewObject()
	_this.context.CurrentEntry.Rule.OnNil(&_this.context)
	_this.receiver.OnNil()
}

func (_this *RulesEventReceiver) OnBool(value bool) {
	_this.context.NotifyNewObject()
	_this.context.CurrentEntry.Rule.OnKeyableObject(&_this.context, DataTypeBool)
	_this.receiver.OnBool(value)
}

func (_this *RulesEventReceiver) OnTrue() {
	_this.context.NotifyNewObject()
	_this.context.CurrentEntry.Rule.OnKeyableObject(&_this.context, DataTypeBool)
	_this.receiver.OnTrue()
}

func (_this *RulesEventReceiver) OnFalse() {
	_this.context.NotifyNewObject()
	_this.context.CurrentEntry.Rule.OnKeyableObject(&_this.context, DataTypeBool)
	_this.receiver.OnFalse()
}

func (_this *RulesEventReceiver) OnPositiveInt(value uint64) {
	_this.context.NotifyNewObject()
	_this.context.CurrentEntry.Rule.OnKeyableObject(&_this.context, DataTypeInt)
	_this.receiver.OnPositiveInt(value)
}

func (_this *RulesEventReceiver) OnNegativeInt(value uint64) {
	_this.context.NotifyNewObject()
	_this.context.CurrentEntry.Rule.OnKeyableObject(&_this.context, DataTypeInt)
	_this.receiver.OnNegativeInt(value)
}

func (_this *RulesEventReceiver) OnInt(value int64) {
	_this.context.NotifyNewObject()
	_this.context.CurrentEntry.Rule.OnKeyableObject(&_this.context, DataTypeInt)
	_this.receiver.OnInt(value)
}

func (_this *RulesEventReceiver) OnBigInt(value *big.Int) {
	if value == nil {
		_this.OnNil()
		return
	}

	_this.context.NotifyNewObject()
	_this.context.CurrentEntry.Rule.OnKeyableObject(&_this.context, DataTypeInt)
	_this.receiver.OnBigInt(value)
}

func (_this *RulesEventReceiver) OnFloat(value float64) {
	if math.IsNaN(value) {
		_this.OnNan(common.IsSignalingNan(value))
		return
	}

	_this.context.NotifyNewObject()
	_this.context.CurrentEntry.Rule.OnKeyableObject(&_this.context, DataTypeFloat)
	_this.receiver.OnFloat(value)
}

func (_this *RulesEventReceiver) OnBigFloat(value *big.Float) {
	if value == nil {
		_this.OnNil()
		return
	}

	_this.context.NotifyNewObject()
	_this.context.CurrentEntry.Rule.OnKeyableObject(&_this.context, DataTypeFloat)
	_this.receiver.OnBigFloat(value)
}

func (_this *RulesEventReceiver) OnDecimalFloat(value compact_float.DFloat) {
	if value.IsNan() {
		_this.OnNan(value.IsSignalingNan())
		return
	}

	_this.context.NotifyNewObject()
	_this.context.CurrentEntry.Rule.OnKeyableObject(&_this.context, DataTypeFloat)
	_this.receiver.OnDecimalFloat(value)
}

func (_this *RulesEventReceiver) OnBigDecimalFloat(value *apd.Decimal) {
	if value == nil {
		_this.OnNil()
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

	_this.context.NotifyNewObject()
	_this.context.CurrentEntry.Rule.OnKeyableObject(&_this.context, DataTypeFloat)
	_this.receiver.OnBigDecimalFloat(value)
}

func (_this *RulesEventReceiver) OnNan(signaling bool) {
	_this.context.NotifyNewObject()
	_this.context.CurrentEntry.Rule.OnNonKeyableObject(&_this.context, DataTypeNan)
	_this.receiver.OnNan(signaling)
}

func (_this *RulesEventReceiver) OnUID(value []byte) {
	_this.context.NotifyNewObject()
	_this.context.CurrentEntry.Rule.OnKeyableObject(&_this.context, DataTypeUID)
	_this.receiver.OnUID(value)
}

func (_this *RulesEventReceiver) OnTime(value time.Time) {
	_this.context.NotifyNewObject()
	_this.context.CurrentEntry.Rule.OnKeyableObject(&_this.context, DataTypeTime)
	_this.receiver.OnTime(value)
}

func (_this *RulesEventReceiver) OnCompactTime(value compact_time.Time) {
	_this.context.NotifyNewObject()
	_this.context.CurrentEntry.Rule.OnKeyableObject(&_this.context, DataTypeTime)
	_this.receiver.OnCompactTime(value)
}

func (_this *RulesEventReceiver) OnArray(arrayType events.ArrayType, elementCount uint64, value []byte) {
	_this.context.NotifyNewObject()
	_this.context.CurrentEntry.Rule.OnArray(&_this.context, arrayType, elementCount, value)
	_this.receiver.OnArray(arrayType, elementCount, value)
}

func (_this *RulesEventReceiver) OnStringlikeArray(arrayType events.ArrayType, value string) {
	_this.context.NotifyNewObject()
	_this.context.CurrentEntry.Rule.OnStringlikeArray(&_this.context, arrayType, value)
	_this.receiver.OnStringlikeArray(arrayType, value)
}

func (_this *RulesEventReceiver) OnArrayBegin(arrayType events.ArrayType) {
	_this.context.NotifyNewObject()
	_this.context.CurrentEntry.Rule.OnArrayBegin(&_this.context, arrayType)
	_this.receiver.OnArrayBegin(arrayType)
}

func (_this *RulesEventReceiver) OnArrayChunk(length uint64, moreChunksFollow bool) {
	_this.context.CurrentEntry.Rule.OnArrayChunk(&_this.context, length, moreChunksFollow)
	_this.receiver.OnArrayChunk(length, moreChunksFollow)
}

func (_this *RulesEventReceiver) OnArrayData(data []byte) {
	_this.context.CurrentEntry.Rule.OnArrayData(&_this.context, data)
	_this.receiver.OnArrayData(data)
}

func (_this *RulesEventReceiver) OnList() {
	_this.context.NotifyNewObject()
	_this.context.CurrentEntry.Rule.OnList(&_this.context)
	_this.receiver.OnList()
}

func (_this *RulesEventReceiver) OnMap() {
	_this.context.NotifyNewObject()
	_this.context.CurrentEntry.Rule.OnMap(&_this.context)
	_this.receiver.OnMap()
}

func (_this *RulesEventReceiver) OnMarkup(identifier []byte) {
	_this.context.ValidateIdentifier(identifier)
	_this.context.NotifyNewObject()
	_this.context.CurrentEntry.Rule.OnMarkup(&_this.context, identifier)
	_this.receiver.OnMarkup(identifier)
}

func (_this *RulesEventReceiver) OnEnd() {
	_this.context.CurrentEntry.Rule.OnEnd(&_this.context)
	_this.receiver.OnEnd()
}

func (_this *RulesEventReceiver) OnNode() {
	_this.context.CurrentEntry.Rule.OnNode(&_this.context)
	_this.receiver.OnNode()
}

func (_this *RulesEventReceiver) OnEdge() {
	_this.context.CurrentEntry.Rule.OnEdge(&_this.context)
	_this.receiver.OnEdge()
}

func (_this *RulesEventReceiver) OnMarker(identifier []byte) {
	_this.context.ValidateMarkerID(identifier)
	_this.context.CurrentEntry.Rule.OnMarker(&_this.context, identifier)
	_this.receiver.OnMarker(identifier)
}

func (_this *RulesEventReceiver) OnReference(identifier []byte) {
	_this.context.ValidateMarkerID(identifier)
	_this.context.NotifyNewObject()
	_this.context.CurrentEntry.Rule.OnReference(&_this.context, identifier)
	_this.receiver.OnReference(identifier)
}

func (_this *RulesEventReceiver) OnConstant(identifier []byte) {
	_this.context.ValidateIdentifier(identifier)
	_this.context.CurrentEntry.Rule.OnConstant(&_this.context, identifier)
	_this.receiver.OnConstant(identifier)
}

func (_this *RulesEventReceiver) OnEndDocument() {
	_this.context.CurrentEntry.Rule.OnEndDocument(&_this.context)
	_this.receiver.OnEndDocument()
}
