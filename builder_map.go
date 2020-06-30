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

const (
	kvBuilderKey   = 0
	kvBuilderValue = 1
)

type mapBuilder struct {
	// Const data
	dstType reflect.Type
	kvTypes [2]reflect.Type

	// Cloned data (must be populated)
	kvBuilders [2]ObjectBuilder

	// Clone inserted data
	root   *RootBuilder
	parent ObjectBuilder

	// Variable data (must be reset)
	container    reflect.Value
	builderIndex int
	key          reflect.Value
}

func newMapBuilder(dstType reflect.Type) ObjectBuilder {
	return &mapBuilder{
		dstType: dstType,
		kvTypes: [2]reflect.Type{dstType.Key(), dstType.Elem()},
	}
}

func (_this *mapBuilder) IsContainerOnly() bool {
	return true
}

func (_this *mapBuilder) PostCacheInitBuilder() {
	_this.kvBuilders[kvBuilderKey] = getBuilderForType(_this.dstType.Key())
	_this.kvBuilders[kvBuilderValue] = getBuilderForType(_this.dstType.Elem())
}

func (_this *mapBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder, options *BuilderOptions) ObjectBuilder {
	that := &mapBuilder{
		dstType: _this.dstType,
		kvTypes: _this.kvTypes,
		parent:  parent,
		root:    root,
	}
	that.kvBuilders[kvBuilderKey] = _this.kvBuilders[kvBuilderKey].CloneFromTemplate(root, that, options)
	that.kvBuilders[kvBuilderValue] = _this.kvBuilders[kvBuilderValue].CloneFromTemplate(root, that, options)
	that.reset()
	return that
}

func (_this *mapBuilder) SetParent(parent ObjectBuilder) {
	_this.parent = parent
}

func (_this *mapBuilder) reset() {
	_this.container = reflect.MakeMap(_this.dstType)
	_this.builderIndex = kvBuilderKey
	_this.key = reflect.Value{}
}

func (_this *mapBuilder) getBuilder() ObjectBuilder {
	return _this.kvBuilders[_this.builderIndex]
}

func (_this *mapBuilder) storeValue(value reflect.Value) {
	if _this.builderIndex == kvBuilderKey {
		_this.key = value
	} else {
		_this.container.SetMapIndex(_this.key, value)
	}
	_this.builderIndex = (_this.builderIndex + 1) & 1
}

func (_this *mapBuilder) newElem() reflect.Value {
	return reflect.New(_this.kvTypes[_this.builderIndex]).Elem()
}

func (_this *mapBuilder) BuildFromNil(ignored reflect.Value) {
	object := _this.newElem()
	_this.getBuilder().BuildFromNil(object)
	_this.storeValue(object)
}

func (_this *mapBuilder) BuildFromBool(value bool, ignored reflect.Value) {
	object := _this.newElem()
	_this.getBuilder().BuildFromBool(value, object)
	_this.storeValue(object)
}

func (_this *mapBuilder) BuildFromInt(value int64, ignored reflect.Value) {
	object := _this.newElem()
	_this.getBuilder().BuildFromInt(value, object)
	_this.storeValue(object)
}

func (_this *mapBuilder) BuildFromUint(value uint64, ignored reflect.Value) {
	object := _this.newElem()
	_this.getBuilder().BuildFromUint(value, object)
	_this.storeValue(object)
}

func (_this *mapBuilder) BuildFromBigInt(value *big.Int, ignored reflect.Value) {
	object := _this.newElem()
	_this.getBuilder().BuildFromBigInt(value, object)
	_this.storeValue(object)
}

func (_this *mapBuilder) BuildFromFloat(value float64, ignored reflect.Value) {
	object := _this.newElem()
	_this.getBuilder().BuildFromFloat(value, object)
	_this.storeValue(object)
}

func (_this *mapBuilder) BuildFromBigFloat(value *big.Float, ignored reflect.Value) {
	object := _this.newElem()
	_this.getBuilder().BuildFromBigFloat(value, object)
	_this.storeValue(object)
}

func (_this *mapBuilder) BuildFromDecimalFloat(value compact_float.DFloat, ignored reflect.Value) {
	object := _this.newElem()
	_this.getBuilder().BuildFromDecimalFloat(value, object)
	_this.storeValue(object)
}

func (_this *mapBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, ignored reflect.Value) {
	object := _this.newElem()
	_this.getBuilder().BuildFromBigDecimalFloat(value, object)
	_this.storeValue(object)
}

func (_this *mapBuilder) BuildFromUUID(value []byte, ignored reflect.Value) {
	object := _this.newElem()
	_this.getBuilder().BuildFromUUID(value, object)
	_this.storeValue(object)
}

func (_this *mapBuilder) BuildFromString(value string, ignored reflect.Value) {
	object := _this.newElem()
	_this.getBuilder().BuildFromString(value, object)
	_this.storeValue(object)
}

func (_this *mapBuilder) BuildFromBytes(value []byte, ignored reflect.Value) {
	object := _this.newElem()
	_this.getBuilder().BuildFromBytes(value, object)
	_this.storeValue(object)
}

func (_this *mapBuilder) BuildFromURI(value *url.URL, ignored reflect.Value) {
	_this.storeValue(reflect.ValueOf(value))
}

func (_this *mapBuilder) BuildFromTime(value time.Time, ignored reflect.Value) {
	_this.storeValue(reflect.ValueOf(value))
}

func (_this *mapBuilder) BuildFromCompactTime(value *compact_time.Time, ignored reflect.Value) {
	_this.storeValue(reflect.ValueOf(value))
}

func (_this *mapBuilder) BuildBeginList() {
	_this.getBuilder().PrepareForListContents()
}

func (_this *mapBuilder) BuildBeginMap() {
	_this.getBuilder().PrepareForMapContents()
}

func (_this *mapBuilder) BuildEndContainer() {
	object := _this.container
	_this.reset()
	_this.parent.NotifyChildContainerFinished(object)
}

func (_this *mapBuilder) BuildFromMarker(id interface{}) {
	panic("TODO: mapBuilder.Marker")
}

func (_this *mapBuilder) BuildFromReference(id interface{}) {
	panic("TODO: mapBuilder.Reference")
}

func (_this *mapBuilder) PrepareForListContents() {
	builderPanicBadEvent(_this, builderIntfType, "PrepareForListContents")
}

func (_this *mapBuilder) PrepareForMapContents() {
	_this.root.SetCurrentBuilder(_this)
}

func (_this *mapBuilder) NotifyChildContainerFinished(value reflect.Value) {
	_this.root.SetCurrentBuilder(_this)
	_this.storeValue(value)
}
