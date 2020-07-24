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
	dstType         reflect.Type
	currentBuilder  ObjectBuilder
	object          reflect.Value
	referenceFiller ReferenceFiller
	session         *Session
}

// -----------
// RootBuilder
// -----------

// Create a new root builder to build objects of dstType. If options is nil,
// default options will be used.
func NewRootBuilder(session *Session, dstType reflect.Type, options *options.BuilderOptions) *RootBuilder {
	_this := &RootBuilder{}
	_this.Init(session, dstType, options)
	return _this
}

// Initialize this root builder to build objects of dstType. If options is nil,
// default options will be used.
func (_this *RootBuilder) Init(session *Session, dstType reflect.Type, options *options.BuilderOptions) {
	_this.session = session
	_this.dstType = dstType
	_this.object = reflect.New(dstType).Elem()

	builder := session.GetBuilderForType(dstType).CloneFromTemplate(_this, _this, options.WithDefaultsApplied())
	_this.currentBuilder = newTopLevelBuilder(_this, builder)
	_this.referenceFiller.Init()
}

func (_this *RootBuilder) String() string {
	return fmt.Sprintf("%v<%v>", reflect.TypeOf(_this), _this.currentBuilder)
}

func (_this *RootBuilder) GetMarkerRegistry() *ReferenceFiller {
	return &_this.referenceFiller
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

func (_this *RootBuilder) PostCacheInitBuilder(_ *Session) {
	BuilderPanicBadEvent(_this, "PostCacheInitBuilder")
}
func (_this *RootBuilder) CloneFromTemplate(_ *RootBuilder, _ ObjectBuilder, _ *options.BuilderOptions) ObjectBuilder {
	BuilderPanicBadEvent(_this, "CloneFromTemplate")
	return nil
}
func (_this *RootBuilder) SetParent(parent ObjectBuilder) {
	BuilderPanicBadEvent(_this, "SetParent")
}
func (_this *RootBuilder) BuildFromNil(_ reflect.Value) {
	BuilderPanicBadEvent(_this, "Nil")
}
func (_this *RootBuilder) BuildFromBool(_ bool, _ reflect.Value) {
	BuilderPanicBadEvent(_this, "Bool")
}
func (_this *RootBuilder) BuildFromInt(_ int64, _ reflect.Value) {
	BuilderPanicBadEvent(_this, "Int")
}
func (_this *RootBuilder) BuildFromUint(_ uint64, _ reflect.Value) {
	BuilderPanicBadEvent(_this, "Uint")
}
func (_this *RootBuilder) BuildFromBigInt(_ *big.Int, _ reflect.Value) {
	BuilderPanicBadEvent(_this, "BigInt")
}
func (_this *RootBuilder) BuildFromFloat(_ float64, _ reflect.Value) {
	BuilderPanicBadEvent(_this, "Float")
}
func (_this *RootBuilder) BuildFromBigFloat(_ *big.Float, _ reflect.Value) {
	BuilderPanicBadEvent(_this, "BigFloat")
}
func (_this *RootBuilder) BuildFromDecimalFloat(_ compact_float.DFloat, _ reflect.Value) {
	BuilderPanicBadEvent(_this, "Float")
}
func (_this *RootBuilder) BuildFromBigDecimalFloat(_ *apd.Decimal, _ reflect.Value) {
	BuilderPanicBadEvent(_this, "DecimalFloat")
}
func (_this *RootBuilder) BuildFromUUID(_ []byte, _ reflect.Value) {
	BuilderPanicBadEvent(_this, "UUID")
}
func (_this *RootBuilder) BuildFromString(_ string, _ reflect.Value) {
	BuilderPanicBadEvent(_this, "String")
}
func (_this *RootBuilder) BuildFromVerbatimString(_ string, _ reflect.Value) {
	BuilderPanicBadEvent(_this, "VerbatimString")
}
func (_this *RootBuilder) BuildFromBytes(_ []byte, _ reflect.Value) {
	BuilderPanicBadEvent(_this, "Bytes")
}
func (_this *RootBuilder) BuildFromCustomBinary(_ []byte, _ reflect.Value) {
	BuilderPanicBadEvent(_this, "CustomBinary")
}
func (_this *RootBuilder) BuildFromCustomText(_ string, _ reflect.Value) {
	BuilderPanicBadEvent(_this, "CustomText")
}
func (_this *RootBuilder) BuildFromURI(_ *url.URL, _ reflect.Value) {
	BuilderPanicBadEvent(_this, "URI")
}
func (_this *RootBuilder) BuildFromTime(_ time.Time, _ reflect.Value) {
	BuilderPanicBadEvent(_this, "Time")
}
func (_this *RootBuilder) BuildFromCompactTime(_ *compact_time.Time, _ reflect.Value) {
	BuilderPanicBadEvent(_this, "CompactTime")
}
func (_this *RootBuilder) BuildBeginList() {
	BuilderPanicBadEvent(_this, "List")
}
func (_this *RootBuilder) BuildBeginMap() {
	BuilderPanicBadEvent(_this, "Map")
}
func (_this *RootBuilder) BuildEndContainer() {
	BuilderPanicBadEvent(_this, "End")
}
func (_this *RootBuilder) BuildBeginMarker(_ interface{}) {
	BuilderPanicBadEvent(_this, "Marker")
}
func (_this *RootBuilder) BuildFromReference(_ interface{}) {
	BuilderPanicBadEvent(_this, "Reference")
}
func (_this *RootBuilder) PrepareForListContents() {
	BuilderPanicBadEvent(_this, "PrepareForListContents")
}
func (_this *RootBuilder) PrepareForMapContents() {
	BuilderPanicBadEvent(_this, "PrepareForMapContents")
}
func (_this *RootBuilder) NotifyChildContainerFinished(value reflect.Value) {
	_this.object = value
}

// -----------------------
// ObjectIteratorCallbacks
// -----------------------

func (_this *RootBuilder) OnVersion(version uint64) {}
func (_this *RootBuilder) OnPadding(count int)      {}
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
func (_this *RootBuilder) OnBytes(value []byte) {
	_this.currentBuilder.BuildFromBytes(value, _this.object)
}
func (_this *RootBuilder) OnString(value string) {
	_this.currentBuilder.BuildFromString(value, _this.object)
}
func (_this *RootBuilder) OnVerbatimString(value string) {
	_this.currentBuilder.BuildFromVerbatimString(value, _this.object)
}
func (_this *RootBuilder) OnURI(value string) {
	u, err := url.Parse(value)
	if err != nil {
		panic(err)
	}
	_this.currentBuilder.BuildFromURI(u, _this.object)
}
func (_this *RootBuilder) OnCustomBinary(value []byte) {
	_this.currentBuilder.BuildFromCustomBinary(value, _this.object)
}
func (_this *RootBuilder) OnCustomText(value string) {
	_this.currentBuilder.BuildFromCustomText(value, _this.object)
}
func (_this *RootBuilder) OnBytesBegin() {
	panic("TODO: RootBuilder.OnBytesBegin")
}
func (_this *RootBuilder) OnStringBegin() {
	panic("TODO: RootBuilder.OnStringBegin")
}
func (_this *RootBuilder) OnVerbatimStringBegin() {
	panic("TODO: RootBuilder.OnVerbatimStringBegin")
}
func (_this *RootBuilder) OnURIBegin() {
	panic("TODO: RootBuilder.OnURIBegin")
}
func (_this *RootBuilder) OnCustomBinaryBegin() {
	panic("TODO: RootBuilder.OnCustomBinaryBegin")
}
func (_this *RootBuilder) OnCustomTextBegin() {
	panic("TODO: RootBuilder.OnCustomTextBegin")
}
func (_this *RootBuilder) OnArrayChunk(length uint64, isFinalChunk bool) {
	panic("TODO: RootBuilder.OnArrayChunk")
}
func (_this *RootBuilder) OnArrayData(data []byte) {
	panic("TODO: RootBuilder.OnArrayData")
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
