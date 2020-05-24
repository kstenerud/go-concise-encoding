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

var (
	builderIntfIntfMapType = reflect.TypeOf(map[interface{}]interface{}{})
	builderIntfSliceType   = reflect.TypeOf([]interface{}{})
	builderIntfType        = builderIntfSliceType.Elem()

	globalIntfBuilder        = &intfBuilder{}
	globalIntfSliceBuilder   = &intfSliceBuilder{}
	globalIntfIntfMapBuilder = &intfIntfMapBuilder{}
)

type intfBuilder struct {
	// Clone inserted data
	root   *RootBuilder
	parent ObjectBuilder
}

func newInterfaceBuilder() ObjectBuilder {
	return globalIntfBuilder
}

func (this *intfBuilder) PostCacheInitBuilder() {
}

func (this *intfBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder) ObjectBuilder {
	return &intfBuilder{
		parent: parent,
		root:   root,
	}
}

func (this *intfBuilder) BuildFromNil(dst reflect.Value) {
	dst.Set(reflect.Zero(builderIntfType))
}

func (this *intfBuilder) BuildFromBool(value bool, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (this *intfBuilder) BuildFromInt(value int64, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (this *intfBuilder) BuildFromUint(value uint64, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (this *intfBuilder) BuildFromBigInt(value *big.Int, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (this *intfBuilder) BuildFromFloat(value float64, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (this *intfBuilder) BuildFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (this *intfBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (this *intfBuilder) BuildFromUUID(value []byte, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (this *intfBuilder) BuildFromString(value string, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (this *intfBuilder) BuildFromBytes(value []byte, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (this *intfBuilder) BuildFromURI(value *url.URL, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (this *intfBuilder) BuildFromTime(value time.Time, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (this *intfBuilder) BuildBeginList() {
	builderPanicBadEvent(this, builderIntfType, "List")
}

func (this *intfBuilder) BuildBeginMap() {
	builderPanicBadEvent(this, builderIntfType, "Map")
}

func (this *intfBuilder) BuildEndContainer() {
	builderPanicBadEvent(this, builderIntfType, "ContainerEnd")
}

func (this *intfBuilder) BuildFromMarker(id interface{}) {
	panic("TODO: intfBuilder.Marker")
}

func (this *intfBuilder) BuildFromReference(id interface{}) {
	panic("TODO: intfBuilder.Reference")
}

func (this *intfBuilder) PrepareForListContents() {
	builder := globalIntfSliceBuilder.CloneFromTemplate(this.root, this.parent)
	builder.PrepareForListContents()
}

func (this *intfBuilder) PrepareForMapContents() {
	builder := globalIntfIntfMapBuilder.CloneFromTemplate(this.root, this.parent)
	builder.PrepareForMapContents()
}

func (this *intfBuilder) NotifyChildContainerFinished(value reflect.Value) {
	this.parent.NotifyChildContainerFinished(value)
}

// -----
// Slice
// -----

type intfSliceBuilder struct {
	// Clone inserted data
	root   *RootBuilder
	parent ObjectBuilder

	// Variable data (must be reset)
	container reflect.Value
}

func newIntfSliceBuilder() ObjectBuilder {
	return globalIntfSliceBuilder
}

func (this *intfSliceBuilder) PostCacheInitBuilder() {
}

func (this *intfSliceBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder) ObjectBuilder {
	that := &intfSliceBuilder{
		parent: parent,
		root:   root,
	}
	that.reset()
	return that
}

func (this *intfSliceBuilder) reset() {
	this.container = reflect.MakeSlice(builderIntfSliceType, 0, defaultSliceCap)
}

func (this *intfSliceBuilder) storeRValue(value reflect.Value) {
	this.container = reflect.Append(this.container, value)
}

func (this *intfSliceBuilder) storeValue(value interface{}) {
	this.storeRValue(reflect.ValueOf(value))
}

func (this *intfSliceBuilder) BuildFromNil(ignored reflect.Value) {
	this.storeRValue(reflect.New(builderIntfSliceType).Elem())
}

func (this *intfSliceBuilder) BuildFromBool(value bool, ignored reflect.Value) {
	this.storeValue(value)
}

func (this *intfSliceBuilder) BuildFromInt(value int64, ignored reflect.Value) {
	this.storeValue(value)
}

func (this *intfSliceBuilder) BuildFromUint(value uint64, ignored reflect.Value) {
	this.storeValue(value)
}

func (this *intfSliceBuilder) BuildFromBigInt(value *big.Int, ignored reflect.Value) {
	this.storeValue(value)
}

func (this *intfSliceBuilder) BuildFromFloat(value float64, ignored reflect.Value) {
	this.storeValue(value)
}

func (this *intfSliceBuilder) BuildFromDecimalFloat(value compact_float.DFloat, ignored reflect.Value) {
	this.storeValue(value)
}

func (this *intfSliceBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, ignored reflect.Value) {
	this.storeValue(value)
}

func (this *intfSliceBuilder) BuildFromUUID(value []byte, ignored reflect.Value) {
	this.storeValue(value)
}

func (this *intfSliceBuilder) BuildFromString(value string, ignored reflect.Value) {
	this.storeValue(value)
}

func (this *intfSliceBuilder) BuildFromBytes(value []byte, ignored reflect.Value) {
	this.storeValue(value)
}

func (this *intfSliceBuilder) BuildFromURI(value *url.URL, ignored reflect.Value) {
	this.storeValue(value)
}

func (this *intfSliceBuilder) BuildFromTime(value time.Time, ignored reflect.Value) {
	this.storeValue(value)
}

func (this *intfSliceBuilder) BuildBeginList() {
	builder := globalIntfSliceBuilder.CloneFromTemplate(this.root, this)
	builder.PrepareForListContents()
}

func (this *intfSliceBuilder) BuildBeginMap() {
	builder := globalIntfIntfMapBuilder.CloneFromTemplate(this.root, this)
	builder.PrepareForMapContents()
}

func (this *intfSliceBuilder) BuildEndContainer() {
	object := this.container
	this.reset()
	this.parent.NotifyChildContainerFinished(object)
}

func (this *intfSliceBuilder) BuildFromMarker(id interface{}) {
	panic("TODO: intfSliceBuilder.Marker")
}

func (this *intfSliceBuilder) BuildFromReference(id interface{}) {
	panic("TODO: intfSliceBuilder.Reference")
}

func (this *intfSliceBuilder) PrepareForListContents() {
	this.root.setCurrentBuilder(this)
}

func (this *intfSliceBuilder) PrepareForMapContents() {
	builderPanicBadEvent(this, builderIntfType, "PrepareForMapContents")
}

func (this *intfSliceBuilder) NotifyChildContainerFinished(value reflect.Value) {
	this.root.setCurrentBuilder(this)
	this.storeRValue(value)
}

// ---
// Map
// ---

type intfIntfMapBuilder struct {
	// Clone inserted data
	root   *RootBuilder
	parent ObjectBuilder

	// Variable data (must be reset)
	container reflect.Value
	key       reflect.Value
	nextIsKey bool
}

func newIntfIntfMapBuilder() ObjectBuilder {
	return globalIntfIntfMapBuilder
}

func (this *intfIntfMapBuilder) PostCacheInitBuilder() {
}

func (this *intfIntfMapBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder) ObjectBuilder {
	that := &intfIntfMapBuilder{
		parent: parent,
		root:   root,
	}
	that.reset()
	return that
}

func (this *intfIntfMapBuilder) reset() {
	this.container = reflect.MakeMap(builderIntfIntfMapType)
	this.key = reflect.Value{}
	this.nextIsKey = true
}

func (this *intfIntfMapBuilder) storeValue(value reflect.Value) {
	if this.nextIsKey {
		this.key = value
	} else {
		this.container.SetMapIndex(this.key, value)
	}
	this.nextIsKey = !this.nextIsKey
}

func (this *intfIntfMapBuilder) BuildFromNil(ignored reflect.Value) {
	this.storeValue(reflect.Zero(builderIntfType))
}

func (this *intfIntfMapBuilder) BuildFromBool(value bool, ignored reflect.Value) {
	this.storeValue(reflect.ValueOf(value))
}

func (this *intfIntfMapBuilder) BuildFromInt(value int64, ignored reflect.Value) {
	this.storeValue(reflect.ValueOf(value))
}

func (this *intfIntfMapBuilder) BuildFromUint(value uint64, ignored reflect.Value) {
	this.storeValue(reflect.ValueOf(value))
}

func (this *intfIntfMapBuilder) BuildFromBigInt(value *big.Int, ignored reflect.Value) {
	this.storeValue(reflect.ValueOf(value))
}

func (this *intfIntfMapBuilder) BuildFromFloat(value float64, ignored reflect.Value) {
	this.storeValue(reflect.ValueOf(value))
}

func (this *intfIntfMapBuilder) BuildFromDecimalFloat(value compact_float.DFloat, ignored reflect.Value) {
	this.storeValue(reflect.ValueOf(value))
}

func (this *intfIntfMapBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, ignored reflect.Value) {
	this.storeValue(reflect.ValueOf(value))
}

func (this *intfIntfMapBuilder) BuildFromUUID(value []byte, ignored reflect.Value) {
	this.storeValue(reflect.ValueOf(value))
}

func (this *intfIntfMapBuilder) BuildFromString(value string, ignored reflect.Value) {
	this.storeValue(reflect.ValueOf(value))
}

func (this *intfIntfMapBuilder) BuildFromBytes(value []byte, ignored reflect.Value) {
	this.storeValue(reflect.ValueOf(value))
}

func (this *intfIntfMapBuilder) BuildFromURI(value *url.URL, ignored reflect.Value) {
	this.storeValue(reflect.ValueOf(value))
}

func (this *intfIntfMapBuilder) BuildFromTime(value time.Time, ignored reflect.Value) {
	this.storeValue(reflect.ValueOf(value))
}

func (this *intfIntfMapBuilder) BuildBeginList() {
	builder := globalIntfSliceBuilder.CloneFromTemplate(this.root, this)
	builder.PrepareForListContents()
}

func (this *intfIntfMapBuilder) BuildBeginMap() {
	builder := globalIntfIntfMapBuilder.CloneFromTemplate(this.root, this)
	builder.PrepareForMapContents()
}

func (this *intfIntfMapBuilder) BuildEndContainer() {
	object := this.container
	this.reset()
	this.parent.NotifyChildContainerFinished(object)
}

func (this *intfIntfMapBuilder) BuildFromMarker(id interface{}) {
	panic("TODO: intfIntfMapBuilder.Marker")
}

func (this *intfIntfMapBuilder) BuildFromReference(id interface{}) {
	panic("TODO: intfIntfMapBuilder.Reference")
}

func (this *intfIntfMapBuilder) PrepareForListContents() {
	builderPanicBadEvent(this, builderIntfType, "PrepareForListContents")
}

func (this *intfIntfMapBuilder) PrepareForMapContents() {
	this.root.setCurrentBuilder(this)
}

func (this *intfIntfMapBuilder) NotifyChildContainerFinished(value reflect.Value) {
	this.root.setCurrentBuilder(this)
	this.storeValue(value)
}
