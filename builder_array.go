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

type arrayBuilder struct {
	// Const data
	dstType reflect.Type

	// Cloned data (must be populated)
	elemBuilder ObjectBuilder

	// Clone inserted data
	root   *RootBuilder
	parent ObjectBuilder

	// Variable data (must be reset)
	container reflect.Value
	index     int
}

func newArrayBuilder(dstType reflect.Type) ObjectBuilder {
	return &arrayBuilder{
		dstType: dstType,
	}
}

func (this *arrayBuilder) PostCacheInitBuilder() {
	this.elemBuilder = getBuilderForType(this.dstType.Elem())
}

func (this *arrayBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder) ObjectBuilder {
	that := &arrayBuilder{
		dstType: this.dstType,
		parent:  parent,
		root:    root,
	}
	that.elemBuilder = this.elemBuilder.CloneFromTemplate(root, that)
	that.reset()
	return that
}

func (this *arrayBuilder) reset() {
	this.container = reflect.New(this.dstType).Elem()
	this.index = 0
}

func (this *arrayBuilder) currentElem() reflect.Value {
	return this.container.Index(this.index)
}

func (this *arrayBuilder) BuildFromNil(ignored reflect.Value) {
	this.elemBuilder.BuildFromNil(this.currentElem())
	this.index++
}

func (this *arrayBuilder) BuildFromBool(value bool, ignored reflect.Value) {
	this.elemBuilder.BuildFromBool(value, this.currentElem())
	this.index++
}

func (this *arrayBuilder) BuildFromInt(value int64, ignored reflect.Value) {
	this.elemBuilder.BuildFromInt(value, this.currentElem())
	this.index++
}

func (this *arrayBuilder) BuildFromUint(value uint64, ignored reflect.Value) {
	this.elemBuilder.BuildFromUint(value, this.currentElem())
	this.index++
}

func (this *arrayBuilder) BuildFromBigInt(value *big.Int, ignored reflect.Value) {
	this.elemBuilder.BuildFromBigInt(value, this.currentElem())
	this.index++
}

func (this *arrayBuilder) BuildFromFloat(value float64, ignored reflect.Value) {
	this.elemBuilder.BuildFromFloat(value, this.currentElem())
	this.index++
}

func (this *arrayBuilder) BuildFromBigFloat(value *big.Float, ignored reflect.Value) {
	this.elemBuilder.BuildFromBigFloat(value, this.currentElem())
	this.index++
}

func (this *arrayBuilder) BuildFromDecimalFloat(value compact_float.DFloat, ignored reflect.Value) {
	this.elemBuilder.BuildFromDecimalFloat(value, this.currentElem())
	this.index++
}

func (this *arrayBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, ignored reflect.Value) {
	this.elemBuilder.BuildFromBigDecimalFloat(value, this.currentElem())
	this.index++
}

func (this *arrayBuilder) BuildFromUUID(value []byte, ignored reflect.Value) {
	this.elemBuilder.BuildFromUUID(value, this.currentElem())
	this.index++
}

func (this *arrayBuilder) BuildFromString(value string, ignored reflect.Value) {
	this.elemBuilder.BuildFromString(value, this.currentElem())
	this.index++
}

func (this *arrayBuilder) BuildFromBytes(value []byte, dst reflect.Value) {
	// TODO: Is this the right way to do this?
	for i := 0; i < len(value); i++ {
		elem := dst.Index(i + this.index)
		elem.SetUint(uint64(value[i]))
	}
}

func (this *arrayBuilder) BuildFromURI(value *url.URL, ignored reflect.Value) {
	this.elemBuilder.BuildFromURI(value, this.currentElem())
	this.index++
}

func (this *arrayBuilder) BuildFromTime(value time.Time, ignored reflect.Value) {
	this.elemBuilder.BuildFromTime(value, this.currentElem())
	this.index++
}

func (this *arrayBuilder) BuildBeginList() {
	this.elemBuilder.PrepareForListContents()
}

func (this *arrayBuilder) BuildBeginMap() {
	this.elemBuilder.PrepareForMapContents()
}

func (this *arrayBuilder) BuildEndContainer() {
	object := this.container
	this.reset()
	this.parent.NotifyChildContainerFinished(object)
}

func (this *arrayBuilder) BuildFromMarker(id interface{}) {
	panic("TODO: arrayBuilder.BuildFromMarker")
}

func (this *arrayBuilder) BuildFromReference(id interface{}) {
	panic("TODO: arrayBuilder.BuildFromReference")
}

func (this *arrayBuilder) PrepareForListContents() {
	this.root.setCurrentBuilder(this)
}

func (this *arrayBuilder) PrepareForMapContents() {
	builderPanicBadEvent(this, this.dstType, "PrepareForMapContents")
}

func (this *arrayBuilder) NotifyChildContainerFinished(value reflect.Value) {
	this.root.setCurrentBuilder(this)
	this.currentElem().Set(value)
	this.index++
}
