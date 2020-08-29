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
	"net/url"
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
	session              *Session
	chunkedData          []byte
	chunkedFunction      func([]byte)
	chunkRemainingLength uint64
	moreChunksFollow     bool
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
	opts = opts.WithDefaultsApplied()
	_this.session = session
	_this.dstType = dstType
	_this.object = reflect.New(dstType).Elem()

	builder := session.GetBuilderForType(dstType).NewInstance(_this, _this, opts)
	_this.currentBuilder = newTopLevelBuilder(_this, builder)
	_this.referenceFiller.Init()
}

func (_this *RootBuilder) String() string {
	return fmt.Sprintf("%v<%v>", reflect.TypeOf(_this), _this.currentBuilder)
}

func (_this *RootBuilder) NotifyMarker(id interface{}, value reflect.Value) {
	_this.referenceFiller.NotifyMarker(id, value)
}

func (_this *RootBuilder) NotifyReference(lookingForID interface{}, valueSetter func(value reflect.Value)) {
	_this.referenceFiller.NotifyReference(lookingForID, valueSetter)
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

func (_this *RootBuilder) SetCurrentBuilder(builder ObjectBuilder) {
	_this.currentBuilder = builder
}

// -------------
// ObjectBuilder
// -------------

func (_this *RootBuilder) InitTemplate(_ *Session) {
	PanicBadEvent(_this, "InitTemplate")
}
func (_this *RootBuilder) NewInstance(_ *RootBuilder, _ ObjectBuilder, _ *options.BuilderOptions) ObjectBuilder {
	PanicBadEvent(_this, "NewInstance")
	return nil
}
func (_this *RootBuilder) SetParent(_ ObjectBuilder) {
	PanicBadEvent(_this, "SetParent")
}
func (_this *RootBuilder) BuildFromNil(_ reflect.Value) {
	PanicBadEvent(_this, "Nil")
}
func (_this *RootBuilder) BuildFromBool(_ bool, _ reflect.Value) {
	PanicBadEvent(_this, "Bool")
}
func (_this *RootBuilder) BuildFromInt(_ int64, _ reflect.Value) {
	PanicBadEvent(_this, "Int")
}
func (_this *RootBuilder) BuildFromUint(_ uint64, _ reflect.Value) {
	PanicBadEvent(_this, "Uint")
}
func (_this *RootBuilder) BuildFromBigInt(_ *big.Int, _ reflect.Value) {
	PanicBadEvent(_this, "BigInt")
}
func (_this *RootBuilder) BuildFromFloat(_ float64, _ reflect.Value) {
	PanicBadEvent(_this, "Float")
}
func (_this *RootBuilder) BuildFromBigFloat(_ *big.Float, _ reflect.Value) {
	PanicBadEvent(_this, "BigFloat")
}
func (_this *RootBuilder) BuildFromDecimalFloat(_ compact_float.DFloat, _ reflect.Value) {
	PanicBadEvent(_this, "Float")
}
func (_this *RootBuilder) BuildFromBigDecimalFloat(_ *apd.Decimal, _ reflect.Value) {
	PanicBadEvent(_this, "DecimalFloat")
}
func (_this *RootBuilder) BuildFromUUID(_ []byte, _ reflect.Value) {
	PanicBadEvent(_this, "UUID")
}
func (_this *RootBuilder) BuildFromString(_ []byte, _ reflect.Value) {
	PanicBadEvent(_this, "String")
}
func (_this *RootBuilder) BuildFromVerbatimString(_ []byte, _ reflect.Value) {
	PanicBadEvent(_this, "VerbatimString")
}
func (_this *RootBuilder) BuildFromURI(_ *url.URL, _ reflect.Value) {
	PanicBadEvent(_this, "URI")
}
func (_this *RootBuilder) BuildFromCustomBinary(_ []byte, _ reflect.Value) {
	PanicBadEvent(_this, "CustomBinary")
}
func (_this *RootBuilder) BuildFromCustomText(_ []byte, _ reflect.Value) {
	PanicBadEvent(_this, "CustomText")
}
func (_this *RootBuilder) BuildFromArray(arrayType events.ArrayType, _ []byte, _ reflect.Value) {
	PanicBadEvent(_this, "TypedArray(%v)", arrayType)
}
func (_this *RootBuilder) BuildFromTime(_ time.Time, _ reflect.Value) {
	PanicBadEvent(_this, "Time")
}
func (_this *RootBuilder) BuildFromCompactTime(_ *compact_time.Time, _ reflect.Value) {
	PanicBadEvent(_this, "CompactTime")
}
func (_this *RootBuilder) BuildBeginList() {
	PanicBadEvent(_this, "List")
}
func (_this *RootBuilder) BuildBeginMap() {
	PanicBadEvent(_this, "Map")
}
func (_this *RootBuilder) BuildEndContainer() {
	PanicBadEvent(_this, "End")
}
func (_this *RootBuilder) BuildBeginMarker(_ interface{}) {
	PanicBadEvent(_this, "Marker")
}
func (_this *RootBuilder) BuildFromReference(_ interface{}) {
	PanicBadEvent(_this, "Reference")
}
func (_this *RootBuilder) PrepareForListContents() {
	PanicBadEvent(_this, "PrepareForListContents")
}
func (_this *RootBuilder) PrepareForMapContents() {
	PanicBadEvent(_this, "PrepareForMapContents")
}
func (_this *RootBuilder) NotifyChildContainerFinished(value reflect.Value) {
	_this.object = value
}

// -----------------------
// ObjectIteratorCallbacks
// -----------------------

func (_this *RootBuilder) OnBeginDocument()   {}
func (_this *RootBuilder) OnVersion(_ uint64) {}
func (_this *RootBuilder) OnPadding(_ int)    {}
func (_this *RootBuilder) OnNil() {
	_this.currentBuilder.BuildFromNil(_this.object)
}
func (_this *RootBuilder) OnBool(value bool) {
	_this.currentBuilder.BuildFromBool(value, _this.object)
}
func (_this *RootBuilder) OnTrue() {
	_this.currentBuilder.BuildFromBool(true, _this.object)
}
func (_this *RootBuilder) OnFalse() {
	_this.currentBuilder.BuildFromBool(false, _this.object)
}
func (_this *RootBuilder) OnPositiveInt(value uint64) {
	_this.currentBuilder.BuildFromUint(value, _this.object)
}
func (_this *RootBuilder) OnNegativeInt(value uint64) {
	if value <= 0x7fffffffffffffff {
		_this.currentBuilder.BuildFromInt(-(int64(value)), _this.object)
		return
	}
	bi := &big.Int{}
	bi.SetUint64(value)
	_this.currentBuilder.BuildFromBigInt(bi.Neg(bi), _this.object)
}
func (_this *RootBuilder) OnInt(value int64) {
	_this.currentBuilder.BuildFromInt(value, _this.object)
}
func (_this *RootBuilder) OnBigInt(value *big.Int) {
	_this.currentBuilder.BuildFromBigInt(value, _this.object)
}
func (_this *RootBuilder) OnFloat(value float64) {
	_this.currentBuilder.BuildFromFloat(value, _this.object)
}
func (_this *RootBuilder) OnBigFloat(value *big.Float) {
	_this.currentBuilder.BuildFromBigFloat(value, _this.object)
}
func (_this *RootBuilder) OnDecimalFloat(value compact_float.DFloat) {
	_this.currentBuilder.BuildFromDecimalFloat(value, _this.object)
}
func (_this *RootBuilder) OnBigDecimalFloat(value *apd.Decimal) {
	_this.currentBuilder.BuildFromBigDecimalFloat(value, _this.object)
}
func (_this *RootBuilder) OnNan(signaling bool) {
	nan := common.QuietNan
	if signaling {
		nan = common.SignalingNan
	}
	_this.currentBuilder.BuildFromFloat(nan, _this.object)
}
func (_this *RootBuilder) OnUUID(value []byte) {
	_this.currentBuilder.BuildFromUUID(value, _this.object)
}
func (_this *RootBuilder) OnTime(value time.Time) {
	_this.currentBuilder.BuildFromTime(value, _this.object)
}
func (_this *RootBuilder) OnCompactTime(value *compact_time.Time) {
	_this.currentBuilder.BuildFromCompactTime(value, _this.object)
}
func (_this *RootBuilder) OnArray(arrayType events.ArrayType, elementCount uint64, value []byte) {
	_this.currentBuilder.BuildFromArray(arrayType, value, _this.object)
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
	_this.currentBuilder.BuildBeginList()
}
func (_this *RootBuilder) OnMap() {
	_this.currentBuilder.BuildBeginMap()
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
	_this.currentBuilder.BuildEndContainer()
}
func (_this *RootBuilder) OnMarker() {
	originalBuilder := _this.currentBuilder
	_this.currentBuilder = newMarkerIDBuilder(func(id interface{}) {
		_this.currentBuilder = originalBuilder
		_this.currentBuilder.BuildBeginMarker(id)
	})
}
func (_this *RootBuilder) OnReference() {
	originalBuilder := _this.currentBuilder
	_this.currentBuilder = newMarkerIDBuilder(func(id interface{}) {
		_this.currentBuilder = originalBuilder
		_this.currentBuilder.BuildFromReference(id)
	})
}
func (_this *RootBuilder) OnEndDocument() {
	// Nothing to do
}
