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

type floatBuilder struct {
	// Const data
	dstType reflect.Type
}

func newFloatBuilder(dstType reflect.Type) ObjectBuilder {
	return &floatBuilder{
		dstType: dstType,
	}
}

func (this *floatBuilder) IsContainerOnly() bool {
	return false
}

func (this *floatBuilder) PostCacheInitBuilder() {
}

func (this *floatBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder) ObjectBuilder {
	return this
}

func (this *floatBuilder) BuildFromNil(dst reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "Nil")
}

func (this *floatBuilder) BuildFromBool(value bool, dst reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "Bool")
}

func (this *floatBuilder) BuildFromInt(value int64, dst reflect.Value) {
	setFloatFromInt(value, dst)
}

func (this *floatBuilder) BuildFromUint(value uint64, dst reflect.Value) {
	setFloatFromUint(value, dst)
}

func (this *floatBuilder) BuildFromBigInt(value *big.Int, dst reflect.Value) {
	setFloatFromBigInt(value, dst)
}

func (this *floatBuilder) BuildFromFloat(value float64, dst reflect.Value) {
	dst.SetFloat(value)
}

func (this *floatBuilder) BuildFromBigFloat(value *big.Float, dst reflect.Value) {
	setFloatFromBigFloat(value, dst)
}

func (this *floatBuilder) BuildFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	dst.SetFloat(value.Float())
}

func (this *floatBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	setFloatFromBigDecimalFloat(value, dst)
}

func (this *floatBuilder) BuildFromUUID(value []byte, dst reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "UUID")
}

func (this *floatBuilder) BuildFromString(value string, dst reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "String")
}

func (this *floatBuilder) BuildFromBytes(value []byte, dst reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "Bytes")
}

func (this *floatBuilder) BuildFromURI(value *url.URL, dst reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "URI")
}

func (this *floatBuilder) BuildFromTime(value time.Time, dst reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "Time")
}

func (this *floatBuilder) BuildBeginList() {
	builderPanicBadEvent(this, this.dstType, "ListBegin")
}

func (this *floatBuilder) BuildBeginMap() {
	builderPanicBadEvent(this, this.dstType, "MapBegin")
}

func (this *floatBuilder) BuildEndContainer() {
	builderPanicBadEvent(this, this.dstType, "ContainerEnd")
}

func (this *floatBuilder) BuildFromMarker(id interface{}) {
	panic("TODO: floatBuilder.BuildFromMarker")
}

func (this *floatBuilder) BuildFromReference(id interface{}) {
	panic("TODO: floatBuilder.BuildFromReference")
}

func (this *floatBuilder) PrepareForListContents() {
	builderPanicBadEvent(this, this.dstType, "PrepareForListContents")
}

func (this *floatBuilder) PrepareForMapContents() {
	builderPanicBadEvent(this, this.dstType, "PrepareForMapContents")
}

func (this *floatBuilder) NotifyChildContainerFinished(value reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "NotifyChildContainerFinished")
}
