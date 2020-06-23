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
	"math/big"
	"net/url"
	"reflect"
	"time"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
)

const defaultSliceCap = 4

type sliceBuilder struct {
	// Const data
	dstType reflect.Type

	// Cloned data (must be populated)
	elemBuilder ObjectBuilder

	// Clone inserted data
	root   *RootBuilder
	parent ObjectBuilder

	// Variable data (must be reset)
	container reflect.Value
}

func newSliceBuilder(dstType reflect.Type) ObjectBuilder {
	return &sliceBuilder{
		dstType: dstType,
	}
}

func (_this *sliceBuilder) IsContainerOnly() bool {
	return true
}

func (_this *sliceBuilder) PostCacheInitBuilder() {
	_this.elemBuilder = getBuilderForType(_this.dstType.Elem())
}

func (_this *sliceBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder, options *BuilderOptions) ObjectBuilder {
	that := &sliceBuilder{
		dstType: _this.dstType,
		parent:  parent,
		root:    root,
	}
	that.elemBuilder = _this.elemBuilder.CloneFromTemplate(root, that, options)
	that.reset()
	return that
}

func (_this *sliceBuilder) reset() {
	_this.container = reflect.MakeSlice(_this.dstType, 0, defaultSliceCap)
}

func (_this *sliceBuilder) newElem() reflect.Value {
	return reflect.New(_this.dstType.Elem()).Elem()
}

func (_this *sliceBuilder) storeValue(value reflect.Value) {
	_this.container = reflect.Append(_this.container, value)
}

func (_this *sliceBuilder) BuildFromNil(ignored reflect.Value) {
	object := _this.newElem()
	_this.elemBuilder.BuildFromNil(object)
	_this.storeValue(object)
}

func (_this *sliceBuilder) BuildFromBool(value bool, ignored reflect.Value) {
	object := _this.newElem()
	_this.elemBuilder.BuildFromBool(value, object)
	_this.storeValue(object)
}

func (_this *sliceBuilder) BuildFromInt(value int64, ignored reflect.Value) {
	object := _this.newElem()
	_this.elemBuilder.BuildFromInt(value, object)
	_this.storeValue(object)
}

func (_this *sliceBuilder) BuildFromUint(value uint64, ignored reflect.Value) {
	object := _this.newElem()
	_this.elemBuilder.BuildFromUint(value, object)
	_this.storeValue(object)
}

func (_this *sliceBuilder) BuildFromBigInt(value *big.Int, ignored reflect.Value) {
	object := _this.newElem()
	_this.elemBuilder.BuildFromBigInt(value, object)
	_this.storeValue(object)
}

func (_this *sliceBuilder) BuildFromFloat(value float64, ignored reflect.Value) {
	object := _this.newElem()
	_this.elemBuilder.BuildFromFloat(value, object)
	_this.storeValue(object)
}

func (_this *sliceBuilder) BuildFromBigFloat(value *big.Float, ignored reflect.Value) {
	object := _this.newElem()
	_this.elemBuilder.BuildFromBigFloat(value, object)
	_this.storeValue(object)
}

func (_this *sliceBuilder) BuildFromDecimalFloat(value compact_float.DFloat, ignored reflect.Value) {
	object := _this.newElem()
	_this.elemBuilder.BuildFromDecimalFloat(value, object)
	_this.storeValue(object)
}

func (_this *sliceBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, ignored reflect.Value) {
	object := _this.newElem()
	_this.elemBuilder.BuildFromBigDecimalFloat(value, object)
	_this.storeValue(object)
}

func (_this *sliceBuilder) BuildFromUUID(value []byte, ignored reflect.Value) {
	object := _this.newElem()
	_this.elemBuilder.BuildFromUUID(value, object)
	_this.storeValue(object)
}

func (_this *sliceBuilder) BuildFromString(value string, ignored reflect.Value) {
	object := _this.newElem()
	_this.elemBuilder.BuildFromString(value, object)
	_this.storeValue(object)
}

func (_this *sliceBuilder) BuildFromBytes(value []byte, ignored reflect.Value) {
	object := _this.newElem()
	_this.elemBuilder.BuildFromBytes(value, object)
	_this.storeValue(object)
}

func (_this *sliceBuilder) BuildFromURI(value *url.URL, ignored reflect.Value) {
	object := _this.newElem()
	_this.elemBuilder.BuildFromURI(value, object)
	_this.storeValue(object)
}

func (_this *sliceBuilder) BuildFromTime(value time.Time, ignored reflect.Value) {
	object := _this.newElem()
	_this.elemBuilder.BuildFromTime(value, object)
	_this.storeValue(object)
}

func (_this *sliceBuilder) BuildFromCompactTime(value *compact_time.Time, ignored reflect.Value) {
	object := _this.newElem()
	_this.elemBuilder.BuildFromCompactTime(value, object)
	_this.storeValue(object)
}

func (_this *sliceBuilder) BuildBeginList() {
	_this.elemBuilder.PrepareForListContents()
}

func (_this *sliceBuilder) BuildBeginMap() {
	_this.elemBuilder.PrepareForMapContents()
}

func (_this *sliceBuilder) BuildEndContainer() {
	object := _this.container
	_this.reset()
	_this.parent.NotifyChildContainerFinished(object)
}

func (_this *sliceBuilder) BuildFromMarker(id interface{}) {
	panic("TODO: sliceBuilder.Marker")
}

func (_this *sliceBuilder) BuildFromReference(id interface{}) {
	panic("TODO: sliceBuilder.Reference")
}

func (_this *sliceBuilder) PrepareForListContents() {
	_this.root.setCurrentBuilder(_this)
}

func (_this *sliceBuilder) PrepareForMapContents() {
	builderPanicBadEvent(_this, _this.dstType, "PrepareForMapContents")
}

func (_this *sliceBuilder) NotifyChildContainerFinished(value reflect.Value) {
	_this.root.setCurrentBuilder(_this)
	_this.storeValue(value)
}
