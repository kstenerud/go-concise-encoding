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

// Imposes the structural rules that enforce a well-formed concise encoding
// document.
package rules

import (
	"math"
	"math/big"
	"time"

	"github.com/kstenerud/go-concise-encoding/events"
	"github.com/kstenerud/go-concise-encoding/internal/common"
	"github.com/kstenerud/go-concise-encoding/options"
	"github.com/kstenerud/go-concise-encoding/version"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
)

// Rules is a DataEventsReceiver passthrough object that constrains the order
// and contents of events to ensure that they form a valid and complete Concise
// Encoding document.
//
// Put this right after your event generator in the event receiver chain to
// enforce correctly formed documents.
//
// Note: This is a LOW LEVEL API. Error reporting is done via panics. Be sure
// to recover() at an appropriate location when calling this struct's methods
// directly (with the exception of constructors and initializers, which are not
// designed to panic).
type Rules struct {
	context  Context
	receiver events.DataEventReceiver
}

// Create a new rules set.
// If opts = nil, defaults are used.
func NewRules(nextReceiver events.DataEventReceiver, opts *options.RuleOptions) *Rules {
	_this := &Rules{}
	_this.Init(nextReceiver, opts)
	return _this
}

// Initialize a rules set.
// If opts = nil, defaults are used.
func (_this *Rules) Init(nextReceiver events.DataEventReceiver, opts *options.RuleOptions) {
	opts = opts.WithDefaultsApplied()
	_this.receiver = nextReceiver
	_this.context.Init(version.ConciseEncodingVersion, opts)
}

// Reset the rules set back to its initial state.
func (_this *Rules) Reset() {
	_this.context.Reset()
}

func (_this *Rules) SetNextReceiver(nextReceiver events.DataEventReceiver) {
	_this.receiver = nextReceiver
}

func (_this *Rules) OnBeginDocument() {
	_this.context.CurrentEntry.Rule.OnBeginDocument(&_this.context)
	_this.receiver.OnBeginDocument()
}

func (_this *Rules) OnVersion(version uint64) {
	_this.context.CurrentEntry.Rule.OnVersion(&_this.context, version)
	_this.receiver.OnVersion(version)
}

func (_this *Rules) OnPadding(count int) {
	_this.context.CurrentEntry.Rule.OnPadding(&_this.context)
	_this.receiver.OnPadding(count)
}

func (_this *Rules) OnNA() {
	_this.context.NotifyNewObject()
	_this.context.CurrentEntry.Rule.OnNonKeyableObject(&_this.context)
	_this.receiver.OnNA()
}

func (_this *Rules) OnBool(value bool) {
	_this.context.NotifyNewObject()
	_this.context.CurrentEntry.Rule.OnKeyableObject(&_this.context)
	_this.receiver.OnBool(value)
}

func (_this *Rules) OnTrue() {
	_this.context.NotifyNewObject()
	_this.context.CurrentEntry.Rule.OnKeyableObject(&_this.context)
	_this.receiver.OnTrue()
}

func (_this *Rules) OnFalse() {
	_this.context.NotifyNewObject()
	_this.context.CurrentEntry.Rule.OnKeyableObject(&_this.context)
	_this.receiver.OnFalse()
}

func (_this *Rules) OnPositiveInt(value uint64) {
	_this.context.NotifyNewObject()
	_this.context.CurrentEntry.Rule.OnPositiveInt(&_this.context, value)
	_this.receiver.OnPositiveInt(value)
}

func (_this *Rules) OnNegativeInt(value uint64) {
	_this.context.NotifyNewObject()
	_this.context.CurrentEntry.Rule.OnKeyableObject(&_this.context)
	_this.receiver.OnNegativeInt(value)
}

func (_this *Rules) OnInt(value int64) {
	_this.context.NotifyNewObject()
	_this.context.CurrentEntry.Rule.OnInt(&_this.context, value)
	_this.receiver.OnInt(value)
}

func (_this *Rules) OnBigInt(value *big.Int) {
	if value == nil {
		_this.OnNA()
		return
	}

	_this.context.NotifyNewObject()
	_this.context.CurrentEntry.Rule.OnBigInt(&_this.context, value)
	_this.receiver.OnBigInt(value)
}

func (_this *Rules) OnFloat(value float64) {
	if math.IsNaN(value) {
		_this.OnNan(common.IsSignalingNan(value))
		return
	}

	_this.context.NotifyNewObject()
	_this.context.CurrentEntry.Rule.OnFloat(&_this.context, value)
	_this.receiver.OnFloat(value)
}

func (_this *Rules) OnBigFloat(value *big.Float) {
	if value == nil {
		_this.OnNA()
		return
	}

	_this.context.NotifyNewObject()
	_this.context.CurrentEntry.Rule.OnBigFloat(&_this.context, value)
	_this.receiver.OnBigFloat(value)
}

func (_this *Rules) OnDecimalFloat(value compact_float.DFloat) {
	if value.IsNan() {
		_this.OnNan(value.IsSignalingNan())
		return
	}

	_this.context.NotifyNewObject()
	_this.context.CurrentEntry.Rule.OnDecimalFloat(&_this.context, value)
	_this.receiver.OnDecimalFloat(value)
}

func (_this *Rules) OnBigDecimalFloat(value *apd.Decimal) {
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

	_this.context.NotifyNewObject()
	_this.context.CurrentEntry.Rule.OnBigDecimalFloat(&_this.context, value)
	_this.receiver.OnBigDecimalFloat(value)
}

func (_this *Rules) OnNan(signaling bool) {
	_this.context.NotifyNewObject()
	_this.context.CurrentEntry.Rule.OnNonKeyableObject(&_this.context)
	_this.receiver.OnNan(signaling)
}

func (_this *Rules) OnUUID(value []byte) {
	_this.context.NotifyNewObject()
	_this.context.CurrentEntry.Rule.OnKeyableObject(&_this.context)
	_this.receiver.OnUUID(value)
}

func (_this *Rules) OnTime(value time.Time) {
	_this.context.NotifyNewObject()
	_this.context.CurrentEntry.Rule.OnKeyableObject(&_this.context)
	_this.receiver.OnTime(value)
}

func (_this *Rules) OnCompactTime(value compact_time.Time) {
	_this.context.NotifyNewObject()
	_this.context.CurrentEntry.Rule.OnKeyableObject(&_this.context)
	_this.receiver.OnCompactTime(value)
}

func (_this *Rules) OnArray(arrayType events.ArrayType, elementCount uint64, value []byte) {
	_this.context.NotifyNewObject()
	_this.context.CurrentEntry.Rule.OnArray(&_this.context, arrayType, elementCount, value)
	_this.receiver.OnArray(arrayType, elementCount, value)
}

func (_this *Rules) OnArrayBegin(arrayType events.ArrayType) {
	_this.context.NotifyNewObject()
	_this.context.CurrentEntry.Rule.OnArrayBegin(&_this.context, arrayType)
	_this.receiver.OnArrayBegin(arrayType)
}

func (_this *Rules) OnArrayChunk(length uint64, moreChunksFollow bool) {
	_this.context.CurrentEntry.Rule.OnArrayChunk(&_this.context, length, moreChunksFollow)
	_this.receiver.OnArrayChunk(length, moreChunksFollow)
}

func (_this *Rules) OnArrayData(data []byte) {
	_this.context.CurrentEntry.Rule.OnArrayData(&_this.context, data)
	_this.receiver.OnArrayData(data)
}

func (_this *Rules) OnConcatenate() {
	panic("TODO: Rules.OnConcatenate")
}

func (_this *Rules) OnList() {
	_this.context.NotifyNewObject()
	_this.context.CurrentEntry.Rule.OnList(&_this.context)
	_this.receiver.OnList()
}

func (_this *Rules) OnMap() {
	_this.context.NotifyNewObject()
	_this.context.CurrentEntry.Rule.OnMap(&_this.context)
	_this.receiver.OnMap()
}

func (_this *Rules) OnMarkup() {
	_this.context.NotifyNewObject()
	_this.context.CurrentEntry.Rule.OnMarkup(&_this.context)
	_this.receiver.OnMarkup()
}

func (_this *Rules) OnMetadata() {
	_this.context.NotifyNewObject()
	_this.context.CurrentEntry.Rule.OnMetadata(&_this.context)
	_this.receiver.OnMetadata()
}

func (_this *Rules) OnComment() {
	_this.context.NotifyNewObject()
	_this.context.CurrentEntry.Rule.OnComment(&_this.context)
	_this.receiver.OnComment()
}

func (_this *Rules) OnEnd() {
	_this.context.CurrentEntry.Rule.OnEnd(&_this.context)
	_this.receiver.OnEnd()
}

func (_this *Rules) OnMarker() {
	_this.context.CurrentEntry.Rule.OnMarker(&_this.context)
	_this.receiver.OnMarker()
}

func (_this *Rules) OnReference() {
	_this.context.NotifyNewObject()
	_this.context.CurrentEntry.Rule.OnReference(&_this.context)
	_this.receiver.OnReference()
}

func (_this *Rules) OnConstant(name []byte, explicitValue bool) {
	// TODO: Probably can handle this all without calling current rule
	_this.receiver.OnConstant(name, explicitValue)
}

func (_this *Rules) OnEndDocument() {
	_this.context.CurrentEntry.Rule.OnEndDocument(&_this.context)
	_this.receiver.OnEndDocument()
}

const maxMarkerIDRuneCount = 50
const maxMarkerIDByteCount = 4 * maxMarkerIDRuneCount // max 4 bytes per rune

type EventRule interface {
	OnBeginDocument(ctx *Context)
	OnEndDocument(ctx *Context)
	OnChildContainerEnded(ctx *Context, cType DataType)
	OnVersion(ctx *Context, version uint64)
	OnPadding(ctx *Context)
	OnKeyableObject(ctx *Context)
	OnNonKeyableObject(ctx *Context)
	OnInt(ctx *Context, value int64)
	OnPositiveInt(ctx *Context, value uint64)
	OnBigInt(ctx *Context, value *big.Int)
	OnFloat(ctx *Context, value float64)
	OnBigFloat(ctx *Context, value *big.Float)
	OnDecimalFloat(ctx *Context, value compact_float.DFloat)
	OnBigDecimalFloat(ctx *Context, value *apd.Decimal)
	OnList(ctx *Context)
	OnMap(ctx *Context)
	OnMarkup(ctx *Context)
	OnMetadata(ctx *Context)
	OnComment(ctx *Context)
	OnEnd(ctx *Context)
	OnMarker(ctx *Context)
	OnReference(ctx *Context)
	OnArray(ctx *Context, arrayType events.ArrayType, elementCount uint64, data []uint8)
	OnArrayBegin(ctx *Context, arrayType events.ArrayType)
	OnArrayChunk(ctx *Context, length uint64, moreChunksFollow bool)
	OnArrayData(ctx *Context, data []byte)
}

type DataType uint

const (
	DataTypeInvalid = 1 << iota
	DataTypeKeyable
	DataTypeAnyType

	AllowKeyable = DataTypeKeyable
	AllowAnyType = AllowKeyable | DataTypeAnyType
)

const keyableTypes = (1 << events.ArrayTypeString) | (1 << events.ArrayTypeResourceID)

func isKeyableType(arrayType events.ArrayType) bool {
	return ((1 << arrayType) & keyableTypes) != 0
}

var (
	beginDocumentRule       BeginDocumentRule
	endDocumentRule         EndDocumentRule
	terminalRule            TerminalRule
	versionRule             VersionRule
	topLevelRule            TopLevelRule
	listRule                ListRule
	mapKeyRule              MapKeyRule
	mapValueRule            MapValueRule
	markupNameRule          MarkupNameRule
	markupKeyRule           MarkupKeyRule
	markupValueRule         MarkupValueRule
	markupContentsRule      MarkupContentsRule
	commentRule             CommentRule
	metadataKeyRule         MetaKeyRule
	metadataValueRule       MetaValueRule
	metadataCompleteRule    MetaCompletionRule
	arrayRule               ArrayRule
	arrayChunkRule          ArrayChunkRule
	stringRule              StringRule
	stringChunkRule         StringChunkRule
	markerIDKeyableRule     MarkerIDKeyableRule
	markerIDAnyTypeRule     MarkerIDAnyTypeRule
	markedObjectKeyableRule MarkedObjectKeyableRule
	markedObjectAnyTypeRule MarkedObjectAnyTypeRule
	referenceKeyableRule    ReferenceKeyableRule
	referenceAnyTypeRule    ReferenceAnyTypeRule
	tlReferenceRIDRule      TLReferenceRIDRule
	stringBuilderRule       StringBuilderRule
	stringBuilderChunkRule  StringBuilderChunkRule
)
