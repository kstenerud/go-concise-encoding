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

type ptrBuilder struct {
	// Const data
	dstType reflect.Type

	// Cloned data (must be populated)
	elemBuilder ObjectBuilder

	// Clone inserted data
	root   *RootBuilder
	parent ObjectBuilder
}

func newPtrBuilder(dstType reflect.Type) ObjectBuilder {
	return &ptrBuilder{
		dstType: dstType,
	}
}

func (this *ptrBuilder) PostCacheInitBuilder() {
	this.elemBuilder = getBuilderForType(this.dstType.Elem())
}

func (this *ptrBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder) ObjectBuilder {
	that := &ptrBuilder{
		dstType: this.dstType,
		parent:  parent,
		root:    root,
	}
	that.elemBuilder = this.elemBuilder.CloneFromTemplate(root, that)
	return that
}

func (this *ptrBuilder) newElem() reflect.Value {
	return reflect.New(this.dstType.Elem())
}

func (this *ptrBuilder) BuildFromNil(dst reflect.Value) {
	dst.Set(reflect.Zero(this.dstType))
}

func (this *ptrBuilder) BuildFromBool(value bool, dst reflect.Value) {
	ptr := this.newElem()
	this.elemBuilder.BuildFromBool(value, ptr.Elem())
	dst.Set(ptr)
}

func (this *ptrBuilder) BuildFromInt(value int64, dst reflect.Value) {
	ptr := this.newElem()
	this.elemBuilder.BuildFromInt(value, ptr.Elem())
	dst.Set(ptr)
}

func (this *ptrBuilder) BuildFromUint(value uint64, dst reflect.Value) {
	ptr := this.newElem()
	this.elemBuilder.BuildFromUint(value, ptr.Elem())
	dst.Set(ptr)
}

func (this *ptrBuilder) BuildFromBigInt(value *big.Int, dst reflect.Value) {
	ptr := this.newElem()
	this.elemBuilder.BuildFromBigInt(value, ptr.Elem())
	dst.Set(ptr)
}

func (this *ptrBuilder) BuildFromFloat(value float64, dst reflect.Value) {
	ptr := this.newElem()
	this.elemBuilder.BuildFromFloat(value, ptr.Elem())
	dst.Set(ptr)
}

func (this *ptrBuilder) BuildFromBigFloat(value *big.Float, dst reflect.Value) {
	ptr := this.newElem()
	this.elemBuilder.BuildFromBigFloat(value, ptr.Elem())
	dst.Set(ptr)
}

func (this *ptrBuilder) BuildFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	ptr := this.newElem()
	this.elemBuilder.BuildFromDecimalFloat(value, ptr.Elem())
	dst.Set(ptr)
}

func (this *ptrBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	ptr := this.newElem()
	this.elemBuilder.BuildFromBigDecimalFloat(value, ptr.Elem())
	dst.Set(ptr)
}

func (this *ptrBuilder) BuildFromUUID(value []byte, dst reflect.Value) {
	ptr := this.newElem()
	this.elemBuilder.BuildFromUUID(value, ptr.Elem())
	dst.Set(ptr)
}

func (this *ptrBuilder) BuildFromString(value string, dst reflect.Value) {
	ptr := this.newElem()
	this.elemBuilder.BuildFromString(value, ptr.Elem())
	dst.Set(ptr)
}

func (this *ptrBuilder) BuildFromBytes(value []byte, dst reflect.Value) {
	ptr := this.newElem()
	this.elemBuilder.BuildFromBytes(value, ptr.Elem())
	dst.Set(ptr)
}

func (this *ptrBuilder) BuildFromURI(value *url.URL, dst reflect.Value) {
	ptr := this.newElem()
	this.elemBuilder.BuildFromURI(value, ptr.Elem())
	dst.Set(ptr)
}

func (this *ptrBuilder) BuildFromTime(value time.Time, dst reflect.Value) {
	ptr := this.newElem()
	this.elemBuilder.BuildFromTime(value, ptr.Elem())
	dst.Set(ptr)
}

func (this *ptrBuilder) BuildBeginList() {
	builderPanicBadEvent(this, this.dstType, "List")
}

func (this *ptrBuilder) BuildBeginMap() {
	builderPanicBadEvent(this, this.dstType, "Map")
}

func (this *ptrBuilder) BuildEndContainer() {
	builderPanicBadEvent(this, this.dstType, "ContainerEnd")
}

func (this *ptrBuilder) BuildFromMarker(id interface{}) {
	panic("TODO: ptrBuilder.Marker")
}

func (this *ptrBuilder) BuildFromReference(id interface{}) {
	panic("TODO: ptrBuilder.Reference")
}

func (this *ptrBuilder) PrepareForListContents() {
	this.elemBuilder.PrepareForListContents()
}

func (this *ptrBuilder) PrepareForMapContents() {
	this.elemBuilder.PrepareForMapContents()
}

func (this *ptrBuilder) NotifyChildContainerFinished(value reflect.Value) {
	this.parent.NotifyChildContainerFinished(value.Addr())
}
