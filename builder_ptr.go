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

type ptrBuilder struct {
	// Const data
	dstType reflect.Type

	// Cloned data (must be populated)
	elemBuilder ObjectBuilder

	// Clone inserted data
	parent ObjectBuilder
}

func newPtrBuilder(dstType reflect.Type) ObjectBuilder {
	return &ptrBuilder{
		dstType: dstType,
	}
}

func (_this *ptrBuilder) IsContainerOnly() bool {
	return _this.elemBuilder.IsContainerOnly()
}

func (_this *ptrBuilder) PostCacheInitBuilder() {
	_this.elemBuilder = getBuilderForType(_this.dstType.Elem())
}

func (_this *ptrBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder, options *BuilderOptions) ObjectBuilder {
	that := &ptrBuilder{
		dstType: _this.dstType,
		parent:  parent,
	}
	that.elemBuilder = _this.elemBuilder.CloneFromTemplate(root, that, options)
	return that
}

func (_this *ptrBuilder) newElem() reflect.Value {
	return reflect.New(_this.dstType.Elem())
}

func (_this *ptrBuilder) BuildFromNil(dst reflect.Value) {
	dst.Set(reflect.Zero(_this.dstType))
}

func (_this *ptrBuilder) BuildFromBool(value bool, dst reflect.Value) {
	ptr := _this.newElem()
	_this.elemBuilder.BuildFromBool(value, ptr.Elem())
	dst.Set(ptr)
}

func (_this *ptrBuilder) BuildFromInt(value int64, dst reflect.Value) {
	ptr := _this.newElem()
	_this.elemBuilder.BuildFromInt(value, ptr.Elem())
	dst.Set(ptr)
}

func (_this *ptrBuilder) BuildFromUint(value uint64, dst reflect.Value) {
	ptr := _this.newElem()
	_this.elemBuilder.BuildFromUint(value, ptr.Elem())
	dst.Set(ptr)
}

func (_this *ptrBuilder) BuildFromBigInt(value *big.Int, dst reflect.Value) {
	ptr := _this.newElem()
	_this.elemBuilder.BuildFromBigInt(value, ptr.Elem())
	dst.Set(ptr)
}

func (_this *ptrBuilder) BuildFromFloat(value float64, dst reflect.Value) {
	ptr := _this.newElem()
	_this.elemBuilder.BuildFromFloat(value, ptr.Elem())
	dst.Set(ptr)
}

func (_this *ptrBuilder) BuildFromBigFloat(value *big.Float, dst reflect.Value) {
	ptr := _this.newElem()
	_this.elemBuilder.BuildFromBigFloat(value, ptr.Elem())
	dst.Set(ptr)
}

func (_this *ptrBuilder) BuildFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	ptr := _this.newElem()
	_this.elemBuilder.BuildFromDecimalFloat(value, ptr.Elem())
	dst.Set(ptr)
}

func (_this *ptrBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	ptr := _this.newElem()
	_this.elemBuilder.BuildFromBigDecimalFloat(value, ptr.Elem())
	dst.Set(ptr)
}

func (_this *ptrBuilder) BuildFromUUID(value []byte, dst reflect.Value) {
	ptr := _this.newElem()
	_this.elemBuilder.BuildFromUUID(value, ptr.Elem())
	dst.Set(ptr)
}

func (_this *ptrBuilder) BuildFromString(value string, dst reflect.Value) {
	ptr := _this.newElem()
	_this.elemBuilder.BuildFromString(value, ptr.Elem())
	dst.Set(ptr)
}

func (_this *ptrBuilder) BuildFromBytes(value []byte, dst reflect.Value) {
	ptr := _this.newElem()
	_this.elemBuilder.BuildFromBytes(value, ptr.Elem())
	dst.Set(ptr)
}

func (_this *ptrBuilder) BuildFromURI(value *url.URL, dst reflect.Value) {
	ptr := _this.newElem()
	_this.elemBuilder.BuildFromURI(value, ptr.Elem())
	dst.Set(ptr)
}

func (_this *ptrBuilder) BuildFromTime(value time.Time, dst reflect.Value) {
	ptr := _this.newElem()
	_this.elemBuilder.BuildFromTime(value, ptr.Elem())
	dst.Set(ptr)
}

func (_this *ptrBuilder) BuildFromCompactTime(value *compact_time.Time, dst reflect.Value) {
	ptr := _this.newElem()
	_this.elemBuilder.BuildFromCompactTime(value, ptr.Elem())
	dst.Set(ptr)
}

func (_this *ptrBuilder) BuildBeginList() {
	builderPanicBadEvent(_this, _this.dstType, "List")
}

func (_this *ptrBuilder) BuildBeginMap() {
	builderPanicBadEvent(_this, _this.dstType, "Map")
}

func (_this *ptrBuilder) BuildEndContainer() {
	builderPanicBadEvent(_this, _this.dstType, "ContainerEnd")
}

func (_this *ptrBuilder) BuildFromMarker(id interface{}) {
	panic("TODO: ptrBuilder.Marker")
}

func (_this *ptrBuilder) BuildFromReference(id interface{}) {
	panic("TODO: ptrBuilder.Reference")
}

func (_this *ptrBuilder) PrepareForListContents() {
	_this.elemBuilder.PrepareForListContents()
}

func (_this *ptrBuilder) PrepareForMapContents() {
	_this.elemBuilder.PrepareForMapContents()
}

func (_this *ptrBuilder) NotifyChildContainerFinished(value reflect.Value) {
	_this.parent.NotifyChildContainerFinished(value.Addr())
}
