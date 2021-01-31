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

package builder

import (
	"fmt"
	"math/big"
	"reflect"
	"time"

	"github.com/kstenerud/go-concise-encoding/events"
	"github.com/kstenerud/go-concise-encoding/internal/common"
	"github.com/kstenerud/go-concise-encoding/options"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
)

// RootBuilder adapts DataEventReceiver to ObjectBuilder, coordinating the
// build via sub-builders.
//
// Use GetBuiltObject() to fetch the final result.
//
// Note: This is a LOW LEVEL API. Error reporting is done via panics. Be sure
// to recover() at an appropriate location when calling this struct's methods
// directly (with the exception of constructors, initializers, and
// GetBuildObject(), which are not designed to panic).
type RootBuilder struct {
	dstType              reflect.Type
	currentBuilder       ObjectBuilder
	object               reflect.Value
	referenceFiller      ReferenceFiller
	chunkedData          []byte
	chunkedFunction      func([]byte)
	chunkRemainingLength uint64
	moreChunksFollow     bool
	context              Context
}

// -----------
// RootBuilder
// -----------

// Create a new root builder to build objects of dstType.
// If opts is nil, default options will be used.
func NewRootBuilder(session *Session, dstType reflect.Type, opts *options.BuilderOptions) *RootBuilder {
	_this := &RootBuilder{}
	_this.Init(session, dstType, opts)
	return _this
}

// Initialize this root builder to build objects of dstType.
// If opts is nil, default options will be used.
func (_this *RootBuilder) Init(session *Session, dstType reflect.Type, opts *options.BuilderOptions) {
	_this.dstType = dstType
	_this.context = context(opts,
		session.opts.CustomBinaryBuildFunction,
		session.opts.CustomTextBuildFunction,
		_this.referenceFiller.NotifyMarker,
		_this.referenceFiller.NotifyReference,
		session.GetBuilderGeneratorForType)
	_this.object = reflect.New(dstType).Elem()
	_this.chunkedData = make([]byte, 0, 128)

	generator := session.GetBuilderGeneratorForType(dstType)
	_this.context.StackBuilder(newTopLevelBuilder(_this, generator))
	_this.referenceFiller.Init()
}

func (_this *RootBuilder) String() string {
	return fmt.Sprintf("%v<%v>", reflect.TypeOf(_this), _this.context.CurrentBuilder)
}

// Get the object that was built after using this root builder as a DataEventReceiver.
func (_this *RootBuilder) GetBuiltObject() interface{} {
	// TODO: Verify this behavior
	if !_this.object.IsValid() {
		return nil
	}

	v := _this.object
	switch v.Kind() {
	case reflect.Struct, reflect.Array:
		return v.Addr().Interface()
	default:
		if _this.object.CanInterface() {
			return _this.object.Interface()
		}
		return nil
	}
}

// Callback from topLevelBuilder
func (_this *RootBuilder) NotifyChildContainerFinished(ctx *Context, value reflect.Value) {
	_this.object = value
}

// ---------------------------
// DataEventReceiver Callbacks
// ---------------------------

func (_this *RootBuilder) OnBeginDocument()   {}
func (_this *RootBuilder) OnVersion(_ uint64) {}
func (_this *RootBuilder) OnPadding(_ int)    {}
func (_this *RootBuilder) OnNA() {
	_this.context.CurrentBuilder.BuildFromNil(&_this.context, _this.object)
}
func (_this *RootBuilder) OnBool(value bool) {
	_this.context.CurrentBuilder.BuildFromBool(&_this.context, value, _this.object)
}
func (_this *RootBuilder) OnTrue() {
	_this.context.CurrentBuilder.BuildFromBool(&_this.context, true, _this.object)
}
func (_this *RootBuilder) OnFalse() {
	_this.context.CurrentBuilder.BuildFromBool(&_this.context, false, _this.object)
}
func (_this *RootBuilder) OnPositiveInt(value uint64) {
	_this.context.CurrentBuilder.BuildFromUint(&_this.context, value, _this.object)
}
func (_this *RootBuilder) OnNegativeInt(value uint64) {
	if value <= 0x7fffffffffffffff {
		_this.context.CurrentBuilder.BuildFromInt(&_this.context, -(int64(value)), _this.object)
		return
	}
	bi := &big.Int{}
	bi.SetUint64(value)
	_this.context.CurrentBuilder.BuildFromBigInt(&_this.context, bi.Neg(bi), _this.object)
}
func (_this *RootBuilder) OnInt(value int64) {
	_this.context.CurrentBuilder.BuildFromInt(&_this.context, value, _this.object)
}
func (_this *RootBuilder) OnBigInt(value *big.Int) {
	_this.context.CurrentBuilder.BuildFromBigInt(&_this.context, value, _this.object)
}
func (_this *RootBuilder) OnFloat(value float64) {
	_this.context.CurrentBuilder.BuildFromFloat(&_this.context, value, _this.object)
}
func (_this *RootBuilder) OnBigFloat(value *big.Float) {
	_this.context.CurrentBuilder.BuildFromBigFloat(&_this.context, value, _this.object)
}
func (_this *RootBuilder) OnDecimalFloat(value compact_float.DFloat) {
	_this.context.CurrentBuilder.BuildFromDecimalFloat(&_this.context, value, _this.object)
}
func (_this *RootBuilder) OnBigDecimalFloat(value *apd.Decimal) {
	_this.context.CurrentBuilder.BuildFromBigDecimalFloat(&_this.context, value, _this.object)
}
func (_this *RootBuilder) OnNan(signaling bool) {
	nan := common.QuietNan
	if signaling {
		nan = common.SignalingNan
	}
	_this.context.CurrentBuilder.BuildFromFloat(&_this.context, nan, _this.object)
}
func (_this *RootBuilder) OnUUID(value []byte) {
	_this.context.CurrentBuilder.BuildFromUUID(&_this.context, value, _this.object)
}
func (_this *RootBuilder) OnTime(value time.Time) {
	_this.context.CurrentBuilder.BuildFromTime(&_this.context, value, _this.object)
}
func (_this *RootBuilder) OnCompactTime(value compact_time.Time) {
	_this.context.CurrentBuilder.BuildFromCompactTime(&_this.context, value, _this.object)
}
func (_this *RootBuilder) OnArray(arrayType events.ArrayType, elementCount uint64, value []byte) {
	_this.context.CurrentBuilder.BuildFromArray(&_this.context, arrayType, value, _this.object)
}
func (_this *RootBuilder) OnArrayBegin(arrayType events.ArrayType) {
	_this.chunkedFunction = func(bytes []byte) {
		elementCount := common.ByteCountToElementCount(arrayType.ElementSize(), uint64(len(bytes)))
		_this.OnArray(arrayType, elementCount, bytes)
	}
	_this.chunkedData = _this.chunkedData[:0]
}
func (_this *RootBuilder) OnArrayChunk(length uint64, moreChunksFollow bool) {
	_this.chunkRemainingLength = length
	_this.moreChunksFollow = moreChunksFollow
	if !_this.moreChunksFollow && _this.chunkRemainingLength == 0 {
		_this.chunkedFunction(_this.chunkedData)
		_this.chunkedData = _this.chunkedData[:0]
	}
}
func (_this *RootBuilder) OnArrayData(data []byte) {
	_this.chunkedData = append(_this.chunkedData, data...)
	_this.chunkRemainingLength -= uint64(len(data))
	if !_this.moreChunksFollow && _this.chunkRemainingLength == 0 {
		_this.chunkedFunction(_this.chunkedData)
		_this.chunkedData = _this.chunkedData[:0]
	}
}
func (_this *RootBuilder) OnList() {
	_this.context.CurrentBuilder.BuildInitiateList(&_this.context)
}
func (_this *RootBuilder) OnMap() {
	_this.context.CurrentBuilder.BuildInitiateMap(&_this.context)
}
func (_this *RootBuilder) OnMarkup() {
	panic("TODO: RootBuilder.OnMarkup")
}
func (_this *RootBuilder) OnMetadata() {
	panic("TODO: RootBuilder.OnMetadata")
}
func (_this *RootBuilder) OnComment() {
	panic("TODO: RootBuilder.OnComment")
}
func (_this *RootBuilder) OnEnd() {
	_this.context.CurrentBuilder.BuildEndContainer(&_this.context)
}
func (_this *RootBuilder) OnMarker() {
	_this.context.StackBuilder(newMarkerIDBuilder())
}
func (_this *RootBuilder) OnReference() {
	_this.context.StackBuilder(newReferenceIDBuilder())
}
func (_this *RootBuilder) OnConcatenate() {
	_this.context.CurrentBuilder.BuildConcatenate(&_this.context)
}
func (_this *RootBuilder) OnConstant(name []byte, explicitValue bool) {
	if !explicitValue {
		panic(fmt.Errorf("Cannot build from constant %s without explicit value", string(name)))
	}
}
func (_this *RootBuilder) OnEndDocument() {
	// Nothing to do
}
