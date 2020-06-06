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

// This builder is used if the top-level object is a container or struct.
type tlContainerBuilder struct {
	// Const data
	dstType reflect.Type

	// Cloned data (must be populated)
	builder ObjectBuilder

	// Clone inserted data
	root   *RootBuilder
	parent ObjectBuilder
}

func newTLContainerBuilder(dstType reflect.Type) ObjectBuilder {
	return &tlContainerBuilder{
		dstType: dstType,
		builder: getBuilderForType(dstType),
	}
}

func (this *tlContainerBuilder) PostCacheInitBuilder() {
	// TLContainer is not cached, so we must init on creation
}

func (this *tlContainerBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder) ObjectBuilder {
	that := &tlContainerBuilder{
		dstType: this.dstType,
		builder: this.builder.CloneFromTemplate(root, parent),
		parent:  parent,
		root:    root,
	}
	return that
}

func (this *tlContainerBuilder) BuildFromNil(ignored reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "Nil")
}

func (this *tlContainerBuilder) BuildFromBool(value bool, ignored reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "Bool")
}

func (this *tlContainerBuilder) BuildFromInt(value int64, ignored reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "Int")
}

func (this *tlContainerBuilder) BuildFromUint(value uint64, ignored reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "Uint")
}

func (this *tlContainerBuilder) BuildFromBigInt(value *big.Int, ignored reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "BigInt")
}

func (this *tlContainerBuilder) BuildFromFloat(value float64, ignored reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "Float")
}

func (this *tlContainerBuilder) BuildFromBigFloat(value *big.Float, ignored reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "BigFloat")
}

func (this *tlContainerBuilder) BuildFromDecimalFloat(value compact_float.DFloat, ignored reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "DecimalFloat")
}

func (this *tlContainerBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, ignored reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "BigDecimalFloat")
}

func (this *tlContainerBuilder) BuildFromUUID(value []byte, ignored reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "UUID")
}

func (this *tlContainerBuilder) BuildFromString(value string, ignored reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "String")
}

func (this *tlContainerBuilder) BuildFromBytes(value []byte, dst reflect.Value) {
	this.root.setCurrentBuilder(this.builder)
	this.builder.BuildFromBytes(value, dst)
}

func (this *tlContainerBuilder) BuildFromURI(value *url.URL, ignored reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "URI")
}

func (this *tlContainerBuilder) BuildFromTime(value time.Time, ignored reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "Time")
}

func (this *tlContainerBuilder) BuildBeginList() {
	this.root.setCurrentBuilder(this.builder)
	this.builder.PrepareForListContents()
}

func (this *tlContainerBuilder) BuildBeginMap() {
	this.root.setCurrentBuilder(this.builder)
	this.builder.PrepareForMapContents()
}

func (this *tlContainerBuilder) BuildEndContainer() {
	builderPanicBadEvent(this, this.dstType, "End")
}

func (this *tlContainerBuilder) PrepareForListContents() {
	builderPanicBadEvent(this, this.dstType, "PrepareForListContents")
}

func (this *tlContainerBuilder) PrepareForMapContents() {
	builderPanicBadEvent(this, this.dstType, "PrepareForMapContents")
}

func (this *tlContainerBuilder) NotifyChildContainerFinished(value reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "NotifyChildContainerFinished")
}

func (this *tlContainerBuilder) BuildFromMarker(id interface{}) {
	panic("TODO: tlContainerBuilder.Marker")
}

func (this *tlContainerBuilder) BuildFromReference(id interface{}) {
	panic("TODO: tlContainerBuilder.Reference")
}
