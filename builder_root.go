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

package concise_encoding

import (
	"fmt"
	"math/big"
	"net/url"
	"reflect"
	"time"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
)

// RootBuilder adapts DataEventReceiver to ObjectBuilder, coordinating the build.
// Use GetBuiltObject() to fetch the final result.
type RootBuilder struct {
	dstType        reflect.Type
	currentBuilder ObjectBuilder
	object         reflect.Value
}

// -----------
// RootBuilder
// -----------

func NewRootBuilder(dstType reflect.Type, options *BuilderOptions) *RootBuilder {
	_this := &RootBuilder{}
	_this.Init(dstType, options)
	return _this
}

func (_this *RootBuilder) Init(dstType reflect.Type, options *BuilderOptions) {
	_this.dstType = dstType
	_this.object = reflect.New(dstType).Elem()

	builder := getBuilderForType(dstType).CloneFromTemplate(_this, _this, applyDefaultBuilderOptions(options))
	if builder.IsContainerOnly() {
		builder = newTopLevelContainerBuilder(_this, builder)
	}
	_this.currentBuilder = builder
}

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

func (_this *RootBuilder) setCurrentBuilder(builder ObjectBuilder) {
	_this.currentBuilder = builder
}

// -------------
// ObjectBuilder
// -------------

func (_this *RootBuilder) IsContainerOnly() bool {
	panic(fmt.Errorf("BUG: IsContainerOnly should never be called on RootBuilder"))
}

func (_this *RootBuilder) PostCacheInitBuilder() {
	panic(fmt.Errorf("BUG: PostCacheInitBuilder should never be called on RootBuilder"))
}
func (_this *RootBuilder) CloneFromTemplate(_ *RootBuilder, _ ObjectBuilder, _ *BuilderOptions) ObjectBuilder {
	panic(fmt.Errorf("BUG: CloneFromTemplate should never be called on RootBuilder"))
}
func (_this *RootBuilder) BuildFromNil(_ reflect.Value) {
	panic(fmt.Errorf("BUG: BuildFromNil should never be called on RootBuilder"))
}
func (_this *RootBuilder) BuildFromBool(_ bool, _ reflect.Value) {
	panic(fmt.Errorf("BUG: BuildFromBool should never be called on RootBuilder"))
}
func (_this *RootBuilder) BuildFromInt(_ int64, _ reflect.Value) {
	panic(fmt.Errorf("BUG: BuildFromInt should never be called on RootBuilder"))
}
func (_this *RootBuilder) BuildFromUint(_ uint64, _ reflect.Value) {
	panic(fmt.Errorf("BUG: BuildFromUint should never be called on RootBuilder"))
}
func (_this *RootBuilder) BuildFromBigInt(_ *big.Int, _ reflect.Value) {
	panic(fmt.Errorf("BUG: BuildFromBigInt should never be called on RootBuilder"))
}
func (_this *RootBuilder) BuildFromFloat(_ float64, _ reflect.Value) {
	panic(fmt.Errorf("BUG: BuildFromFloat should never be called on RootBuilder"))
}
func (_this *RootBuilder) BuildFromBigFloat(_ *big.Float, _ reflect.Value) {
	panic(fmt.Errorf("BUG: BuildFromBigFloat should never be called on RootBuilder"))
}
func (_this *RootBuilder) BuildFromDecimalFloat(_ compact_float.DFloat, _ reflect.Value) {
	panic(fmt.Errorf("BUG: BuildFromDecimalFloat should never be called on RootBuilder"))
}
func (_this *RootBuilder) BuildFromBigDecimalFloat(_ *apd.Decimal, _ reflect.Value) {
	panic(fmt.Errorf("BUG: BuildFromBigDecimalFloat should never be called on RootBuilder"))
}
func (_this *RootBuilder) BuildFromUUID(_ []byte, _ reflect.Value) {
	panic(fmt.Errorf("BUG: BuildFromUUID should never be called on RootBuilder"))
}
func (_this *RootBuilder) BuildFromString(_ string, _ reflect.Value) {
	panic(fmt.Errorf("BUG: BuildFromString should never be called on RootBuilder"))
}
func (_this *RootBuilder) BuildFromBytes(_ []byte, _ reflect.Value) {
	panic(fmt.Errorf("BUG: BuildFromBytes should never be called on RootBuilder"))
}
func (_this *RootBuilder) BuildFromURI(_ *url.URL, _ reflect.Value) {
	panic(fmt.Errorf("BUG: BuildFromURI should never be called on RootBuilder"))
}
func (_this *RootBuilder) BuildFromTime(_ time.Time, _ reflect.Value) {
	panic(fmt.Errorf("BUG: BuildFromTime should never be called on RootBuilder"))
}
func (_this *RootBuilder) BuildFromCompactTime(_ *compact_time.Time, _ reflect.Value) {
	panic(fmt.Errorf("BUG: BuildFromCompactTime should never be called on RootBuilder"))
}
func (_this *RootBuilder) BuildBeginList() {
	panic(fmt.Errorf("BUG: BuildBeginList should never be called on RootBuilder"))
}
func (_this *RootBuilder) BuildBeginMap() {
	panic(fmt.Errorf("BUG: BuildBeginMap should never be called on RootBuilder"))
}
func (_this *RootBuilder) BuildEndContainer() {
	panic(fmt.Errorf("BUG: BuildEndContainer should never be called on RootBuilder"))
}
func (_this *RootBuilder) BuildFromMarker(_ interface{}) {
	panic(fmt.Errorf("BUG: BuildFromMarker should never be called on RootBuilder"))
}
func (_this *RootBuilder) BuildFromReference(_ interface{}) {
	panic(fmt.Errorf("BUG: BuildFromReference should never be called on RootBuilder"))
}
func (_this *RootBuilder) PrepareForListContents() {
	panic(fmt.Errorf("BUG: PrepareForListContents should never be called on RootBuilder"))
}
func (_this *RootBuilder) PrepareForMapContents() {
	panic(fmt.Errorf("BUG: PrepareForMapContents should never be called on RootBuilder"))
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
	nan := quietNan
	if signaling {
		nan = signalingNan
	}
	_this.currentBuilder.BuildFromFloat(nan, _this.object)
}
func (_this *RootBuilder) OnComplex(value complex128) {
	panic("TODO: RootBuilder.OnComplex")
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
func (_this *RootBuilder) OnURI(value string) {
	u, err := url.Parse(value)
	if err != nil {
		panic(err)
	}
	_this.currentBuilder.BuildFromURI(u, _this.object)
}
func (_this *RootBuilder) OnCustom(value []byte) {
	panic("TODO: RootBuilder.OnCustom")
}
func (_this *RootBuilder) OnBytesBegin() {
	panic("TODO: RootBuilder.OnBytesBegin")
}
func (_this *RootBuilder) OnStringBegin() {
	panic("TODO: RootBuilder.OnStringBegin")
}
func (_this *RootBuilder) OnURIBegin() {
	panic("TODO: RootBuilder.OnURIBegin")
}
func (_this *RootBuilder) OnCustomBegin() {
	panic("TODO: RootBuilder.OnCustomBegin")
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
	panic("TODO: RootBuilder.OnMarker")
}
func (_this *RootBuilder) OnReference() {
	panic("TODO: RootBuilder.OnReference")
}
func (_this *RootBuilder) OnEndDocument() {
	// Nothing to do
}
