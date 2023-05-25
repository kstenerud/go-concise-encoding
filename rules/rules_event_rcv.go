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
	"fmt"
	"math"
	"math/big"

	"github.com/cockroachdb/apd/v2"
	compact_float "github.com/kstenerud/go-compact-float"
	compact_time "github.com/kstenerud/go-compact-time"
	"github.com/kstenerud/go-concise-encoding/ce/events"
	"github.com/kstenerud/go-concise-encoding/configuration"
	"github.com/kstenerud/go-concise-encoding/internal/common"
	"github.com/kstenerud/go-concise-encoding/nullevent"
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
// If config = nil, defaults are used.
func NewRules(nextReceiver events.DataEventReceiver, config *configuration.RuleConfiguration) *RulesEventReceiver {
	_this := &RulesEventReceiver{}
	_this.Init(nextReceiver, config)
	return _this
}

var nullReceiver = &nullevent.NullEventReceiver{}

// Initialize a rules set.
// If config = nil, defaults are used.
func (_this *RulesEventReceiver) Init(nextReceiver events.DataEventReceiver, config *configuration.RuleConfiguration) {
	if config == nil {
		defaultConfig := configuration.DefaultRuleConfiguration()
		config = &defaultConfig
	} else {
		config.ApplyDefaults()
	}

	_this.receiver = nextReceiver
	if _this.receiver == nil {
		_this.receiver = nullReceiver
	}
	_this.context.Init(config)
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

func (_this *RulesEventReceiver) OnPadding() {
	_this.context.CurrentEntry.Rule.OnPadding(&_this.context)
	_this.receiver.OnPadding()
}

func (_this *RulesEventReceiver) OnComment(isMultiline bool, contents []byte) {
	// TODO: Validate comment contents
	_this.context.CurrentEntry.Rule.OnComment(&_this.context)
	_this.receiver.OnComment(isMultiline, contents)
}

func (_this *RulesEventReceiver) OnNull() {
	_this.context.NotifyNewObject(true)
	_this.context.CurrentEntry.Rule.OnNull(&_this.context)
	_this.receiver.OnNull()
}

func (_this *RulesEventReceiver) OnBoolean(value bool) {
	_this.context.NotifyNewObject(true)
	_this.context.CurrentEntry.Rule.OnKeyableObject(&_this.context, DataTypeBool)
	_this.receiver.OnBoolean(value)
}

func (_this *RulesEventReceiver) OnTrue() {
	_this.context.NotifyNewObject(true)
	_this.context.CurrentEntry.Rule.OnKeyableObject(&_this.context, DataTypeBool)
	_this.receiver.OnTrue()
}

func (_this *RulesEventReceiver) OnFalse() {
	_this.context.NotifyNewObject(true)
	_this.context.CurrentEntry.Rule.OnKeyableObject(&_this.context, DataTypeBool)
	_this.receiver.OnFalse()
}

func (_this *RulesEventReceiver) OnPositiveInt(value uint64) {
	_this.context.NotifyNewObject(true)
	_this.context.CurrentEntry.Rule.OnKeyableObject(&_this.context, DataTypeInt)
	_this.receiver.OnPositiveInt(value)
}

func (_this *RulesEventReceiver) OnNegativeInt(value uint64) {
	_this.context.NotifyNewObject(true)
	_this.context.CurrentEntry.Rule.OnKeyableObject(&_this.context, DataTypeInt)
	_this.receiver.OnNegativeInt(value)
}

func (_this *RulesEventReceiver) OnInt(value int64) {
	_this.context.NotifyNewObject(true)
	_this.context.CurrentEntry.Rule.OnKeyableObject(&_this.context, DataTypeInt)
	_this.receiver.OnInt(value)
}

func (_this *RulesEventReceiver) OnBigInt(value *big.Int) {
	if value == nil {
		_this.OnNull()
		return
	}

	_this.context.NotifyNewObject(true)
	_this.context.CurrentEntry.Rule.OnKeyableObject(&_this.context, DataTypeInt)
	_this.receiver.OnBigInt(value)
}

func (_this *RulesEventReceiver) OnFloat(value float64) {
	if math.IsNaN(value) {
		_this.OnNan(!common.HasQuietNanBitSet64(value))
		return
	}

	_this.context.NotifyNewObject(true)
	_this.context.CurrentEntry.Rule.OnKeyableObject(&_this.context, DataTypeFloat)
	_this.receiver.OnFloat(value)
}

func (_this *RulesEventReceiver) OnBigFloat(value *big.Float) {
	if value == nil {
		_this.OnNull()
		return
	}

	_this.context.NotifyNewObject(true)
	_this.context.CurrentEntry.Rule.OnKeyableObject(&_this.context, DataTypeFloat)
	_this.receiver.OnBigFloat(value)
}

func (_this *RulesEventReceiver) OnDecimalFloat(value compact_float.DFloat) {
	if value.IsNan() {
		_this.OnNan(value.IsSignalingNan())
		return
	}

	_this.context.NotifyNewObject(true)
	_this.context.CurrentEntry.Rule.OnKeyableObject(&_this.context, DataTypeFloat)
	_this.receiver.OnDecimalFloat(value)
}

func (_this *RulesEventReceiver) OnBigDecimalFloat(value *apd.Decimal) {
	if value == nil {
		_this.OnNull()
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

	_this.context.NotifyNewObject(true)
	_this.context.CurrentEntry.Rule.OnKeyableObject(&_this.context, DataTypeFloat)
	_this.receiver.OnBigDecimalFloat(value)
}

func (_this *RulesEventReceiver) OnNan(signaling bool) {
	_this.context.NotifyNewObject(true)
	_this.context.CurrentEntry.Rule.OnNonKeyableObject(&_this.context, DataTypeNan)
	_this.receiver.OnNan(signaling)
}

func (_this *RulesEventReceiver) OnUID(value []byte) {
	_this.context.NotifyNewObject(true)
	_this.context.CurrentEntry.Rule.OnKeyableObject(&_this.context, DataTypeUID)
	_this.receiver.OnUID(value)
}

func (_this *RulesEventReceiver) OnTime(value compact_time.Time) {
	_this.context.NotifyNewObject(true)
	_this.context.CurrentEntry.Rule.OnKeyableObject(&_this.context, DataTypeTime)
	_this.receiver.OnTime(value)
}

func (_this *RulesEventReceiver) validateArrayAPICall(arrayType events.ArrayType) {
	switch arrayType {
	case events.ArrayTypeCustomBinary, events.ArrayTypeCustomText:
		panic(fmt.Errorf("BUG: %v is not allowed in the array API. Use the custom type API instead", arrayType))
	case events.ArrayTypeMedia:
		panic(fmt.Errorf("BUG: %v is not allowed in the array API. Use the media API instead", arrayType))
	default:
	}
}

func (_this *RulesEventReceiver) validateCustomTypeAPICall(arrayType events.ArrayType) {
	switch arrayType {
	case events.ArrayTypeCustomBinary, events.ArrayTypeCustomText:
	default:
		panic(fmt.Errorf("BUG: %v is not allowed in the custom type API. Use the array API instead", arrayType))
	}
}

func (_this *RulesEventReceiver) OnArray(arrayType events.ArrayType, elementCount uint64, value []byte) {
	_this.validateArrayAPICall(arrayType)
	_this.context.NotifyNewObject(true)
	_this.context.CurrentEntry.Rule.OnArray(&_this.context, arrayType, elementCount, value)
	_this.receiver.OnArray(arrayType, elementCount, value)
}

func (_this *RulesEventReceiver) OnStringlikeArray(arrayType events.ArrayType, value string) {
	_this.validateArrayAPICall(arrayType)
	_this.context.NotifyNewObject(true)
	_this.context.CurrentEntry.Rule.OnStringlikeArray(&_this.context, arrayType, value)
	_this.receiver.OnStringlikeArray(arrayType, value)
}

func (_this *RulesEventReceiver) OnMedia(mediaType string, value []byte) {
	if len(mediaType) > 0xffffffff {
		panic(fmt.Errorf("media type is too long (%v bytes)", len(mediaType)))
	}
	_this.context.NotifyNewObject(true)
	_this.context.CurrentEntry.Rule.OnArray(&_this.context, events.ArrayTypeMedia, uint64(len(value)), value)
	_this.receiver.OnMedia(mediaType, value)
}

func (_this *RulesEventReceiver) OnCustomBinary(customType uint64, value []byte) {
	_this.context.NotifyNewObject(true)
	_this.context.CurrentEntry.Rule.OnArray(&_this.context, events.ArrayTypeCustomBinary, uint64(len(value)), value)
	_this.receiver.OnCustomBinary(customType, value)
}

func (_this *RulesEventReceiver) OnCustomText(customType uint64, value string) {
	_this.context.NotifyNewObject(true)
	_this.context.CurrentEntry.Rule.OnStringlikeArray(&_this.context, events.ArrayTypeCustomText, value)
	_this.receiver.OnCustomText(customType, value)
}

func (_this *RulesEventReceiver) OnArrayBegin(arrayType events.ArrayType) {
	_this.validateArrayAPICall(arrayType)
	_this.context.NotifyNewObject(true)
	_this.context.CurrentEntry.Rule.OnArrayBegin(&_this.context, arrayType)
	_this.receiver.OnArrayBegin(arrayType)
}

func (_this *RulesEventReceiver) OnMediaBegin(mediaType string) {
	_this.context.NotifyNewObject(true)
	_this.context.CurrentEntry.Rule.OnArrayBegin(&_this.context, events.ArrayTypeMedia)
	_this.receiver.OnMediaBegin(mediaType)
}

func (_this *RulesEventReceiver) OnCustomBegin(arrayType events.ArrayType, customType uint64) {
	_this.validateCustomTypeAPICall(arrayType)
	_this.context.NotifyNewObject(true)
	_this.context.CurrentEntry.Rule.OnArrayBegin(&_this.context, arrayType)
	_this.receiver.OnCustomBegin(arrayType, customType)
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
	_this.context.NotifyNewObject(true)
	_this.context.CurrentEntry.Rule.OnList(&_this.context)
	_this.receiver.OnList()
}

func (_this *RulesEventReceiver) OnMap() {
	_this.context.NotifyNewObject(true)
	_this.context.CurrentEntry.Rule.OnMap(&_this.context)
	_this.receiver.OnMap()
}

func (_this *RulesEventReceiver) OnEndContainer() {
	_this.context.CurrentEntry.Rule.OnEnd(&_this.context)
	_this.receiver.OnEndContainer()
}

func (_this *RulesEventReceiver) OnRecordType(identifier []byte) {
	_this.context.NotifyNewObject(false)
	_this.context.ValidateIdentifier(identifier)
	_this.context.CurrentEntry.Rule.OnRecordType(&_this.context, identifier)
	_this.receiver.OnRecordType(identifier)
}

func (_this *RulesEventReceiver) OnRecord(identifier []byte) {
	_this.context.NotifyNewObject(true)
	_this.context.ValidateIdentifier(identifier)
	_this.context.CurrentEntry.Rule.OnRecord(&_this.context, identifier)
	_this.receiver.OnRecord(identifier)
}

func (_this *RulesEventReceiver) OnNode() {
	_this.context.NotifyNewObject(true)
	_this.context.CurrentEntry.Rule.OnNode(&_this.context)
	_this.receiver.OnNode()
}

func (_this *RulesEventReceiver) OnEdge() {
	_this.context.NotifyNewObject(true)
	_this.context.CurrentEntry.Rule.OnEdge(&_this.context)
	_this.receiver.OnEdge()
}

func (_this *RulesEventReceiver) OnMarker(identifier []byte) {
	_this.context.ValidateIdentifier(identifier)
	_this.context.NotifyNewObject(true)
	_this.context.CurrentEntry.Rule.OnMarker(&_this.context, identifier)
	_this.receiver.OnMarker(identifier)
}

func (_this *RulesEventReceiver) OnReferenceLocal(identifier []byte) {
	_this.context.ValidateIdentifier(identifier)
	_this.context.NotifyNewObject(true)
	_this.context.CurrentEntry.Rule.OnReferenceLocal(&_this.context, identifier)
	_this.receiver.OnReferenceLocal(identifier)
}

func (_this *RulesEventReceiver) OnEndDocument() {
	_this.context.CurrentEntry.Rule.OnEndDocument(&_this.context)
	_this.receiver.OnEndDocument()
}
