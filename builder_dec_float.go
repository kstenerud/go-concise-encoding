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

type dfloatBuilder struct {
}

func newDFloatBuilder() ObjectBuilder {
	return &dfloatBuilder{}
}

func (this *dfloatBuilder) PostCacheInitBuilder() {
}

func (this *dfloatBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder) ObjectBuilder {
	return this
}

func (this *dfloatBuilder) BuildFromNil(dst reflect.Value) {
	builderPanicBadEvent(this, typeDFloat, "Nil")
}

func (this *dfloatBuilder) BuildFromBool(value bool, dst reflect.Value) {
	builderPanicBadEvent(this, typeDFloat, "Bool")
}

func (this *dfloatBuilder) BuildFromInt(value int64, dst reflect.Value) {
	dst.Set(reflect.ValueOf(compact_float.DFloatValue(0, value)))
}

func (this *dfloatBuilder) BuildFromUint(value uint64, dst reflect.Value) {
	dst.Set(reflect.ValueOf(compact_float.DFloatFromUInt(value)))
}

func (this *dfloatBuilder) BuildFromBigInt(value *big.Int, dst reflect.Value) {
	dst.Set(reflect.ValueOf(compact_float.DFloatFromBigInt(value)))
}

func (this *dfloatBuilder) BuildFromFloat(value float64, dst reflect.Value) {
	dst.Set(reflect.ValueOf(compact_float.DFloatFromFloat64(value, 0)))
}

func (this *dfloatBuilder) BuildFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (this *dfloatBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	dst.Set(reflect.ValueOf(compact_float.DFloatFromAPD(value)))
}

func (this *dfloatBuilder) BuildFromUUID(value []byte, dst reflect.Value) {
	builderPanicBadEvent(this, typeDFloat, "UUID")
}

func (this *dfloatBuilder) BuildFromString(value string, dst reflect.Value) {
	builderPanicBadEvent(this, typeDFloat, "String")
}

func (this *dfloatBuilder) BuildFromBytes(value []byte, dst reflect.Value) {
	builderPanicBadEvent(this, typeDFloat, "Bytes")
}

func (this *dfloatBuilder) BuildFromURI(value *url.URL, dst reflect.Value) {
	builderPanicBadEvent(this, typeDFloat, "URI")
}

func (this *dfloatBuilder) BuildFromTime(value time.Time, dst reflect.Value) {
	builderPanicBadEvent(this, typeDFloat, "Time")
}

func (this *dfloatBuilder) BuildBeginList() {
	builderPanicBadEvent(this, typeDFloat, "List")
}

func (this *dfloatBuilder) BuildBeginMap() {
	builderPanicBadEvent(this, typeDFloat, "Map")
}

func (this *dfloatBuilder) BuildEndContainer() {
	builderPanicBadEvent(this, typeDFloat, "ContainerEnd")
}

func (this *dfloatBuilder) BuildFromMarker(id interface{}) {
	panic("TODO: dfloatBuilder.Marker")
}

func (this *dfloatBuilder) BuildFromReference(id interface{}) {
	panic("TODO: dfloatBuilder.Reference")
}

func (this *dfloatBuilder) PrepareForListContents() {
}

func (this *dfloatBuilder) PrepareForMapContents() {
}

func (this *dfloatBuilder) NotifyChildContainerFinished(value reflect.Value) {
	builderPanicBadEvent(this, typeDFloat, "NotifyChildContainerFinished")
}
