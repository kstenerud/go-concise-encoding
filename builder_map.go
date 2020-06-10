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

func (this *mapBuilder) IsContainerOnly() bool {
	return true
}

func (this *mapBuilder) PostCacheInitBuilder() {
	this.kvBuilders[kvBuilderKey] = getBuilderForType(this.dstType.Key())
	this.kvBuilders[kvBuilderValue] = getBuilderForType(this.dstType.Elem())
}

func (this *mapBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder) ObjectBuilder {
	that := &mapBuilder{
		dstType: this.dstType,
		kvTypes: this.kvTypes,
		parent:  parent,
		root:    root,
	}
	that.kvBuilders[kvBuilderKey] = this.kvBuilders[kvBuilderKey].CloneFromTemplate(root, that)
	that.kvBuilders[kvBuilderValue] = this.kvBuilders[kvBuilderValue].CloneFromTemplate(root, that)
	that.reset()
	return that
}

func (this *mapBuilder) reset() {
	this.container = reflect.MakeMap(this.dstType)
	this.builderIndex = kvBuilderKey
	this.key = reflect.Value{}
}

func (this *mapBuilder) getBuilder() ObjectBuilder {
	return this.kvBuilders[this.builderIndex]
}

func (this *mapBuilder) storeValue(value reflect.Value) {
	if this.builderIndex == kvBuilderKey {
		this.key = value
	} else {
		this.container.SetMapIndex(this.key, value)
	}
	this.builderIndex = (this.builderIndex + 1) & 1
}

func (this *mapBuilder) newElem() reflect.Value {
	return reflect.New(this.kvTypes[this.builderIndex]).Elem()
}

func (this *mapBuilder) BuildFromNil(ignored reflect.Value) {
	object := this.newElem()
	this.getBuilder().BuildFromNil(object)
	this.storeValue(object)
}

func (this *mapBuilder) BuildFromBool(value bool, ignored reflect.Value) {
	object := this.newElem()
	this.getBuilder().BuildFromBool(value, object)
	this.storeValue(object)
}

func (this *mapBuilder) BuildFromInt(value int64, ignored reflect.Value) {
	object := this.newElem()
	this.getBuilder().BuildFromInt(value, object)
	this.storeValue(object)
}

func (this *mapBuilder) BuildFromUint(value uint64, ignored reflect.Value) {
	object := this.newElem()
	this.getBuilder().BuildFromUint(value, object)
	this.storeValue(object)
}

func (this *mapBuilder) BuildFromBigInt(value *big.Int, ignored reflect.Value) {
	object := this.newElem()
	this.getBuilder().BuildFromBigInt(value, object)
	this.storeValue(object)
}

func (this *mapBuilder) BuildFromFloat(value float64, ignored reflect.Value) {
	object := this.newElem()
	this.getBuilder().BuildFromFloat(value, object)
	this.storeValue(object)
}

func (this *mapBuilder) BuildFromBigFloat(value *big.Float, ignored reflect.Value) {
	object := this.newElem()
	this.getBuilder().BuildFromBigFloat(value, object)
	this.storeValue(object)
}

func (this *mapBuilder) BuildFromDecimalFloat(value compact_float.DFloat, ignored reflect.Value) {
	object := this.newElem()
	this.getBuilder().BuildFromDecimalFloat(value, object)
	this.storeValue(object)
}

func (this *mapBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, ignored reflect.Value) {
	object := this.newElem()
	this.getBuilder().BuildFromBigDecimalFloat(value, object)
	this.storeValue(object)
}

func (this *mapBuilder) BuildFromUUID(value []byte, ignored reflect.Value) {
	object := this.newElem()
	this.getBuilder().BuildFromUUID(value, object)
	this.storeValue(object)
}

func (this *mapBuilder) BuildFromString(value string, ignored reflect.Value) {
	object := this.newElem()
	this.getBuilder().BuildFromString(value, object)
	this.storeValue(object)
}

func (this *mapBuilder) BuildFromBytes(value []byte, ignored reflect.Value) {
	object := this.newElem()
	this.getBuilder().BuildFromBytes(value, object)
	this.storeValue(object)
}

func (this *mapBuilder) BuildFromURI(value *url.URL, ignored reflect.Value) {
	this.storeValue(reflect.ValueOf(value))
}

func (this *mapBuilder) BuildFromTime(value time.Time, ignored reflect.Value) {
	this.storeValue(reflect.ValueOf(value))
}

func (this *mapBuilder) BuildBeginList() {
	this.getBuilder().PrepareForListContents()
}

func (this *mapBuilder) BuildBeginMap() {
	this.getBuilder().PrepareForMapContents()
}

func (this *mapBuilder) BuildEndContainer() {
	object := this.container
	this.reset()
	this.parent.NotifyChildContainerFinished(object)
}

func (this *mapBuilder) BuildFromMarker(id interface{}) {
	panic("TODO: mapBuilder.Marker")
}

func (this *mapBuilder) BuildFromReference(id interface{}) {
	panic("TODO: mapBuilder.Reference")
}

func (this *mapBuilder) PrepareForListContents() {
	builderPanicBadEvent(this, builderIntfType, "PrepareForListContents")
}

func (this *mapBuilder) PrepareForMapContents() {
	this.root.setCurrentBuilder(this)
}

func (this *mapBuilder) NotifyChildContainerFinished(value reflect.Value) {
	this.root.setCurrentBuilder(this)
	this.storeValue(value)
}
