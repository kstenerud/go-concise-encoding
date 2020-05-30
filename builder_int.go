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

type intBuilder struct {
	// Const data
	dstType reflect.Type
}

func newIntBuilder(dstType reflect.Type) ObjectBuilder {
	return &intBuilder{
		dstType: dstType,
	}
}

func (this *intBuilder) PostCacheInitBuilder() {
}

func (this *intBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder) ObjectBuilder {
	return this
}

func (this *intBuilder) BuildFromNil(dst reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "Nil")
}

func (this *intBuilder) BuildFromBool(value bool, dst reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "Bool")
}

func (this *intBuilder) BuildFromInt(value int64, dst reflect.Value) {
	setIntFromInt(value, dst)
}

func (this *intBuilder) BuildFromUint(value uint64, dst reflect.Value) {
	setIntFromUint(value, dst)
}

func (this *intBuilder) BuildFromBigInt(value *big.Int, dst reflect.Value) {
	setIntFromBigInt(value, dst)
}

func (this *intBuilder) BuildFromFloat(value float64, dst reflect.Value) {
	setIntFromFloat(value, dst)
}

func (this *intBuilder) BuildFromBigFloat(value *big.Float, dst reflect.Value) {
	setIntFromBigFloat(value, dst)
}

func (this *intBuilder) BuildFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	setIntFromDecimalFloat(value, dst)
}

func (this *intBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	setIntFromBigDecimalFloat(value, dst)
}

func (this *intBuilder) BuildFromUUID(value []byte, dst reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "UUID")
}

func (this *intBuilder) BuildFromString(value string, dst reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "String")
}

func (this *intBuilder) BuildFromBytes(value []byte, dst reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "Bytes")
}

func (this *intBuilder) BuildFromURI(value *url.URL, dst reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "URI")
}

func (this *intBuilder) BuildFromTime(value time.Time, dst reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "Time")
}

func (this *intBuilder) BuildBeginList() {
	builderPanicBadEvent(this, this.dstType, "List")
}

func (this *intBuilder) BuildBeginMap() {
	builderPanicBadEvent(this, this.dstType, "Map")
}

func (this *intBuilder) BuildEndContainer() {
	builderPanicBadEvent(this, this.dstType, "ContainerEnd")
}

func (this *intBuilder) BuildFromMarker(id interface{}) {
	panic("TODO: intBuilder.Marker")
}

func (this *intBuilder) BuildFromReference(id interface{}) {
	panic("TODO: intBuilder.Reference")
}

func (this *intBuilder) PrepareForListContents() {
}

func (this *intBuilder) PrepareForMapContents() {
}

func (this *intBuilder) NotifyChildContainerFinished(value reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "NotifyChildContainerFinished")
}
