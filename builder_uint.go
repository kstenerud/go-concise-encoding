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

type uintBuilder struct {
	// Const data
	dstType reflect.Type
}

func newUintBuilder(dstType reflect.Type) ObjectBuilder {
	return &uintBuilder{
		dstType: dstType,
	}
}

func (this *uintBuilder) IsContainerOnly() bool {
	return false
}

func (this *uintBuilder) PostCacheInitBuilder() {
}

func (this *uintBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder, options *BuilderOptions) ObjectBuilder {
	return this
}

func (this *uintBuilder) BuildFromNil(dst reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "Nil")
}

func (this *uintBuilder) BuildFromBool(value bool, dst reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "Bool")
}

func (this *uintBuilder) BuildFromInt(value int64, dst reflect.Value) {
	setUintFromInt(value, dst)
}

func (this *uintBuilder) BuildFromUint(value uint64, dst reflect.Value) {
	setUintFromUint(value, dst)
}

func (this *uintBuilder) BuildFromBigInt(value *big.Int, dst reflect.Value) {
	setUintFromBigInt(value, dst)
}

func (this *uintBuilder) BuildFromFloat(value float64, dst reflect.Value) {
	setUintFromFloat(value, dst)
}

func (this *uintBuilder) BuildFromBigFloat(value *big.Float, dst reflect.Value) {
	setUintFromBigFloat(value, dst)
}

func (this *uintBuilder) BuildFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	setUintFromDecimalFloat(value, dst)
}

func (this *uintBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	setUintFromBigDecimalFloat(value, dst)
}

func (this *uintBuilder) BuildFromUUID(value []byte, dst reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "UUID")
}

func (this *uintBuilder) BuildFromString(value string, dst reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "String")
}

func (this *uintBuilder) BuildFromBytes(value []byte, dst reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "Bytes")
}

func (this *uintBuilder) BuildFromURI(value *url.URL, dst reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "URI")
}

func (this *uintBuilder) BuildFromTime(value time.Time, dst reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "Time")
}

func (this *uintBuilder) BuildFromCompactTime(value *compact_time.Time, dst reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "CompactTime")
}

func (this *uintBuilder) BuildBeginList() {
	builderPanicBadEvent(this, this.dstType, "List")
}

func (this *uintBuilder) BuildBeginMap() {
	builderPanicBadEvent(this, this.dstType, "Map")
}

func (this *uintBuilder) BuildEndContainer() {
	builderPanicBadEvent(this, this.dstType, "ContainerEnd")
}

func (this *uintBuilder) BuildFromMarker(id interface{}) {
	panic("TODO: uintBuilder.Marker")
}

func (this *uintBuilder) BuildFromReference(id interface{}) {
	panic("TODO: uintBuilder.Reference")
}

func (this *uintBuilder) PrepareForListContents() {
	builderPanicBadEvent(this, this.dstType, "PrepareForListContents")
}

func (this *uintBuilder) PrepareForMapContents() {
	builderPanicBadEvent(this, this.dstType, "PrepareForMapContents")
}

func (this *uintBuilder) NotifyChildContainerFinished(value reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "NotifyChildContainerFinished")
}
