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
	container   reflect.Value
	nextBuilder ObjectBuilder
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

func (_this *sliceBuilder) SetParent(parent ObjectBuilder) {
	_this.parent = parent
}

func (_this *sliceBuilder) reset() {
	_this.container = reflect.MakeSlice(_this.dstType, 0, defaultSliceCap)
	_this.nextBuilder = _this.elemBuilder
}

func (_this *sliceBuilder) newElem() reflect.Value {
	return reflect.New(_this.dstType.Elem()).Elem()
}

func (_this *sliceBuilder) storeValue(value reflect.Value) {
	_this.container = reflect.Append(_this.container, value)
}

func (_this *sliceBuilder) BuildFromNil(ignored reflect.Value) {
	object := _this.newElem()
	_this.nextBuilder.BuildFromNil(object)
	_this.nextBuilder = _this.elemBuilder
	_this.storeValue(object)
}

func (_this *sliceBuilder) BuildFromBool(value bool, ignored reflect.Value) {
	object := _this.newElem()
	_this.nextBuilder.BuildFromBool(value, object)
	_this.nextBuilder = _this.elemBuilder
	_this.storeValue(object)
}

func (_this *sliceBuilder) BuildFromInt(value int64, ignored reflect.Value) {
	object := _this.newElem()
	_this.nextBuilder.BuildFromInt(value, object)
	_this.nextBuilder = _this.elemBuilder
	_this.storeValue(object)
}

func (_this *sliceBuilder) BuildFromUint(value uint64, ignored reflect.Value) {
	object := _this.newElem()
	_this.nextBuilder.BuildFromUint(value, object)
	_this.nextBuilder = _this.elemBuilder
	_this.storeValue(object)
}

func (_this *sliceBuilder) BuildFromBigInt(value *big.Int, ignored reflect.Value) {
	object := _this.newElem()
	_this.nextBuilder.BuildFromBigInt(value, object)
	_this.nextBuilder = _this.elemBuilder
	_this.storeValue(object)
}

func (_this *sliceBuilder) BuildFromFloat(value float64, ignored reflect.Value) {
	object := _this.newElem()
	_this.nextBuilder.BuildFromFloat(value, object)
	_this.nextBuilder = _this.elemBuilder
	_this.storeValue(object)
}

func (_this *sliceBuilder) BuildFromBigFloat(value *big.Float, ignored reflect.Value) {
	object := _this.newElem()
	_this.nextBuilder.BuildFromBigFloat(value, object)
	_this.nextBuilder = _this.elemBuilder
	_this.storeValue(object)
}

func (_this *sliceBuilder) BuildFromDecimalFloat(value compact_float.DFloat, ignored reflect.Value) {
	object := _this.newElem()
	_this.nextBuilder.BuildFromDecimalFloat(value, object)
	_this.nextBuilder = _this.elemBuilder
	_this.storeValue(object)
}

func (_this *sliceBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, ignored reflect.Value) {
	object := _this.newElem()
	_this.nextBuilder.BuildFromBigDecimalFloat(value, object)
	_this.nextBuilder = _this.elemBuilder
	_this.storeValue(object)
}

func (_this *sliceBuilder) BuildFromUUID(value []byte, ignored reflect.Value) {
	object := _this.newElem()
	_this.nextBuilder.BuildFromUUID(value, object)
	_this.nextBuilder = _this.elemBuilder
	_this.storeValue(object)
}

func (_this *sliceBuilder) BuildFromString(value string, ignored reflect.Value) {
	object := _this.newElem()
	_this.nextBuilder.BuildFromString(value, object)
	_this.nextBuilder = _this.elemBuilder
	_this.storeValue(object)
}

func (_this *sliceBuilder) BuildFromBytes(value []byte, ignored reflect.Value) {
	object := _this.newElem()
	_this.nextBuilder.BuildFromBytes(value, object)
	_this.nextBuilder = _this.elemBuilder
	_this.storeValue(object)
}

func (_this *sliceBuilder) BuildFromURI(value *url.URL, ignored reflect.Value) {
	object := _this.newElem()
	_this.nextBuilder.BuildFromURI(value, object)
	_this.nextBuilder = _this.elemBuilder
	_this.storeValue(object)
}

func (_this *sliceBuilder) BuildFromTime(value time.Time, ignored reflect.Value) {
	object := _this.newElem()
	_this.nextBuilder.BuildFromTime(value, object)
	_this.nextBuilder = _this.elemBuilder
	_this.storeValue(object)
}

func (_this *sliceBuilder) BuildFromCompactTime(value *compact_time.Time, ignored reflect.Value) {
	object := _this.newElem()
	_this.nextBuilder.BuildFromCompactTime(value, object)
	_this.nextBuilder = _this.elemBuilder
	_this.storeValue(object)
}

func (_this *sliceBuilder) BuildBeginList() {
	_this.nextBuilder.PrepareForListContents()
	_this.nextBuilder = _this.elemBuilder
}

func (_this *sliceBuilder) BuildBeginMap() {
	_this.nextBuilder.PrepareForMapContents()
	_this.nextBuilder = _this.elemBuilder
}

func (_this *sliceBuilder) BuildEndContainer() {
	object := _this.container
	_this.reset()
	_this.parent.NotifyChildContainerFinished(object)
}

// Marker: next obj is ID, obj after that also gets stored to registry
// Reference: next obj is ID, fetch and set obj from registry when available

// To store marker, need:
// - ID
// - Value
// - Registry

// To store reference, need:
// - ID
// - ptr to zero value
// - Registry

// In registry:
// NotifyMarker(id, value)
// - store marker
// - check for unmet references, set value
// NotifyReference(id, ptr-to-zero-value)
// - If marker exists, set value
// - If marker absent, store reference
// Rules about reference existing can be done in rules object

// Fetch ID by hijacking root current builder with marker ID builder.
// - marker ID builder has ptr to real builder, restores it after fetching ID
// -- need wrapper around real builder that waits for real builder to build, then NotifyMarker with result
// - where is ID kept during this?

// build from reference:
// - hijack to get ID
// - set using zero value
// - NotifyReference - need to pass in dst type?
// - registry will set later

func (_this *sliceBuilder) BuildFromMarker(id interface{}) {
	//_this.nextBuilder = newMarkerBuilder(...)
	// or maybe pass in builder?
	panic("TODO: sliceBuilder.Marker")
}

func (_this *sliceBuilder) BuildFromReference(id interface{}) {
	// Change this to "BuildFromObject" and just pass in finished object?
	// Need a global slow-setter that accepts anything for anything
	panic("TODO: sliceBuilder.Reference")
}

func (_this *sliceBuilder) PrepareForListContents() {
	_this.root.SetCurrentBuilder(_this)
}

func (_this *sliceBuilder) PrepareForMapContents() {
	builderPanicBadEvent(_this, _this.dstType, "PrepareForMapContents")
}

func (_this *sliceBuilder) NotifyChildContainerFinished(value reflect.Value) {
	_this.root.SetCurrentBuilder(_this)
	_this.storeValue(value)
}
