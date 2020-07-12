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

	"github.com/kstenerud/go-concise-encoding/options"

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
	root    *RootBuilder
	parent  ObjectBuilder
	options *options.BuilderOptions

	// Variable data (must be reset)
	ppContainer **reflect.Value
}

func newSliceBuilder(dstType reflect.Type) ObjectBuilder {
	return &sliceBuilder{
		dstType: dstType,
	}
}

func (_this *sliceBuilder) String() string {
	return fmt.Sprintf("%v<%v>", reflect.TypeOf(_this), _this.elemBuilder)
}

func (_this *sliceBuilder) IsContainerOnly() bool {
	return true
}

func (_this *sliceBuilder) PostCacheInitBuilder() {
	_this.elemBuilder = getBuilderForType(_this.dstType.Elem())
}

func (_this *sliceBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder, options *options.BuilderOptions) ObjectBuilder {
	that := &sliceBuilder{
		dstType:     _this.dstType,
		parent:      parent,
		root:        root,
		elemBuilder: _this.elemBuilder,
		options:     options,
	}
	that.reset()
	return that
}

func (_this *sliceBuilder) SetParent(parent ObjectBuilder) {
	_this.parent = parent
}

func (_this *sliceBuilder) reset() {
	container := reflect.MakeSlice(_this.dstType, 0, defaultSliceCap)
	_this.ppContainer = new(*reflect.Value)
	*_this.ppContainer = &container
}

func (_this *sliceBuilder) newElem() reflect.Value {
	return reflect.New(_this.dstType.Elem()).Elem()
}

func (_this *sliceBuilder) storeValue(value reflect.Value) {
	**_this.ppContainer = reflect.Append(**_this.ppContainer, value)
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
	object := **_this.ppContainer
	_this.reset()
	_this.parent.NotifyChildContainerFinished(object)
}

func (_this *sliceBuilder) BuildBeginMarker(id interface{}) {
	origBuilder := _this.elemBuilder
	_this.elemBuilder = newMarkerObjectBuilder(_this, origBuilder, func(object reflect.Value) {
		_this.elemBuilder = origBuilder
		_this.root.GetMarkerRegistry().NotifyMarker(id, object)
	})
}

func (_this *sliceBuilder) BuildFromReference(id interface{}) {
	ppContainer := _this.ppContainer
	index := (**ppContainer).Len()
	elem := _this.newElem()
	_this.storeValue(elem)
	_this.root.GetMarkerRegistry().NotifyReference(id, func(object reflect.Value) {
		setAnythingFromAnything(object, (**ppContainer).Index(index))
	})
}

func (_this *sliceBuilder) PrepareForListContents() {
	_this.elemBuilder = _this.elemBuilder.CloneFromTemplate(_this.root, _this, _this.options)
	_this.root.SetCurrentBuilder(_this)
}

func (_this *sliceBuilder) PrepareForMapContents() {
	builderPanicBadEventType(_this, _this.dstType, "PrepareForMapContents")
}

func (_this *sliceBuilder) NotifyChildContainerFinished(value reflect.Value) {
	_this.root.SetCurrentBuilder(_this)
	_this.storeValue(value)
}
