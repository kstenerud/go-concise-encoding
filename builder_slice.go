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

func (this *sliceBuilder) IsContainerOnly() bool {
	return true
}

func (this *sliceBuilder) PostCacheInitBuilder() {
	this.elemBuilder = getBuilderForType(this.dstType.Elem())
}

func (this *sliceBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder, options *BuilderOptions) ObjectBuilder {
	that := &sliceBuilder{
		dstType: this.dstType,
		parent:  parent,
		root:    root,
	}
	that.elemBuilder = this.elemBuilder.CloneFromTemplate(root, that, options)
	that.reset()
	return that
}

func (this *sliceBuilder) reset() {
	this.container = reflect.MakeSlice(this.dstType, 0, defaultSliceCap)
}

func (this *sliceBuilder) newElem() reflect.Value {
	return reflect.New(this.dstType.Elem()).Elem()
}

func (this *sliceBuilder) storeValue(value reflect.Value) {
	this.container = reflect.Append(this.container, value)
}

func (this *sliceBuilder) BuildFromNil(ignored reflect.Value) {
	object := this.newElem()
	this.elemBuilder.BuildFromNil(object)
	this.storeValue(object)
}

func (this *sliceBuilder) BuildFromBool(value bool, ignored reflect.Value) {
	object := this.newElem()
	this.elemBuilder.BuildFromBool(value, object)
	this.storeValue(object)
}

func (this *sliceBuilder) BuildFromInt(value int64, ignored reflect.Value) {
	object := this.newElem()
	this.elemBuilder.BuildFromInt(value, object)
	this.storeValue(object)
}

func (this *sliceBuilder) BuildFromUint(value uint64, ignored reflect.Value) {
	object := this.newElem()
	this.elemBuilder.BuildFromUint(value, object)
	this.storeValue(object)
}

func (this *sliceBuilder) BuildFromBigInt(value *big.Int, ignored reflect.Value) {
	object := this.newElem()
	this.elemBuilder.BuildFromBigInt(value, object)
	this.storeValue(object)
}

func (this *sliceBuilder) BuildFromFloat(value float64, ignored reflect.Value) {
	object := this.newElem()
	this.elemBuilder.BuildFromFloat(value, object)
	this.storeValue(object)
}

func (this *sliceBuilder) BuildFromBigFloat(value *big.Float, ignored reflect.Value) {
	object := this.newElem()
	this.elemBuilder.BuildFromBigFloat(value, object)
	this.storeValue(object)
}

func (this *sliceBuilder) BuildFromDecimalFloat(value compact_float.DFloat, ignored reflect.Value) {
	object := this.newElem()
	this.elemBuilder.BuildFromDecimalFloat(value, object)
	this.storeValue(object)
}

func (this *sliceBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, ignored reflect.Value) {
	object := this.newElem()
	this.elemBuilder.BuildFromBigDecimalFloat(value, object)
	this.storeValue(object)
}

func (this *sliceBuilder) BuildFromUUID(value []byte, ignored reflect.Value) {
	object := this.newElem()
	this.elemBuilder.BuildFromUUID(value, object)
	this.storeValue(object)
}

func (this *sliceBuilder) BuildFromString(value string, ignored reflect.Value) {
	object := this.newElem()
	this.elemBuilder.BuildFromString(value, object)
	this.storeValue(object)
}

func (this *sliceBuilder) BuildFromBytes(value []byte, ignored reflect.Value) {
	object := this.newElem()
	this.elemBuilder.BuildFromBytes(value, object)
	this.storeValue(object)
}

func (this *sliceBuilder) BuildFromURI(value *url.URL, ignored reflect.Value) {
	object := this.newElem()
	this.elemBuilder.BuildFromURI(value, object)
	this.storeValue(object)
}

func (this *sliceBuilder) BuildFromTime(value time.Time, ignored reflect.Value) {
	object := this.newElem()
	this.elemBuilder.BuildFromTime(value, object)
	this.storeValue(object)
}

func (this *sliceBuilder) BuildFromCompactTime(value *compact_time.Time, ignored reflect.Value) {
	object := this.newElem()
	this.elemBuilder.BuildFromCompactTime(value, object)
	this.storeValue(object)
}

func (this *sliceBuilder) BuildBeginList() {
	this.elemBuilder.PrepareForListContents()
}

func (this *sliceBuilder) BuildBeginMap() {
	this.elemBuilder.PrepareForMapContents()
}

func (this *sliceBuilder) BuildEndContainer() {
	object := this.container
	this.reset()
	this.parent.NotifyChildContainerFinished(object)
}

func (this *sliceBuilder) BuildFromMarker(id interface{}) {
	panic("TODO: sliceBuilder.Marker")
}

func (this *sliceBuilder) BuildFromReference(id interface{}) {
	panic("TODO: sliceBuilder.Reference")
}

func (this *sliceBuilder) PrepareForListContents() {
	this.root.setCurrentBuilder(this)
}

func (this *sliceBuilder) PrepareForMapContents() {
	builderPanicBadEvent(this, this.dstType, "PrepareForMapContents")
}

func (this *sliceBuilder) NotifyChildContainerFinished(value reflect.Value) {
	this.root.setCurrentBuilder(this)
	this.storeValue(value)
}
