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

type stringBuilder struct {
}

var globalStringBuilder = &stringBuilder{}

func newStringBuilder() ObjectBuilder {
	return globalStringBuilder
}

func (this *stringBuilder) IsContainerOnly() bool {
	return false
}

func (this *stringBuilder) PostCacheInitBuilder() {
}

func (this *stringBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder) ObjectBuilder {
	return this
}

func (this *stringBuilder) BuildFromNil(dst reflect.Value) {
	// Go doesn't have the concept of a nil string.
	dst.SetString("")
}

func (this *stringBuilder) BuildFromBool(value bool, dst reflect.Value) {
	builderPanicBadEvent(this, typeString, "Bool")
}

func (this *stringBuilder) BuildFromInt(value int64, dst reflect.Value) {
	builderPanicBadEvent(this, typeString, "Int")
}

func (this *stringBuilder) BuildFromUint(value uint64, dst reflect.Value) {
	builderPanicBadEvent(this, typeString, "Uint")
}

func (this *stringBuilder) BuildFromBigInt(value *big.Int, dst reflect.Value) {
	builderPanicBadEvent(this, typeString, "BigInt")
}

func (this *stringBuilder) BuildFromFloat(value float64, dst reflect.Value) {
	builderPanicBadEvent(this, typeString, "Float")
}

func (this *stringBuilder) BuildFromBigFloat(value *big.Float, dst reflect.Value) {
	builderPanicBadEvent(this, typeString, "BigFloat")
}

func (this *stringBuilder) BuildFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	builderPanicBadEvent(this, typeString, "DecimalFloat")
}

func (this *stringBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	builderPanicBadEvent(this, typeString, "BigDecimalFloat")
}

func (this *stringBuilder) BuildFromUUID(value []byte, dst reflect.Value) {
	builderPanicBadEvent(this, typeString, "UUID")
}

func (this *stringBuilder) BuildFromString(value string, dst reflect.Value) {
	dst.SetString(value)
}

func (this *stringBuilder) BuildFromBytes(value []byte, dst reflect.Value) {
	builderPanicBadEvent(this, typeString, "Bytes")
}

func (this *stringBuilder) BuildFromURI(value *url.URL, dst reflect.Value) {
	builderPanicBadEvent(this, typeString, "URI")
}

func (this *stringBuilder) BuildFromTime(value time.Time, dst reflect.Value) {
	builderPanicBadEvent(this, typeString, "Time")
}

func (this *stringBuilder) BuildBeginList() {
	builderPanicBadEvent(this, typeString, "List")
}

func (this *stringBuilder) BuildBeginMap() {
	builderPanicBadEvent(this, typeString, "Map")
}

func (this *stringBuilder) BuildEndContainer() {
	builderPanicBadEvent(this, typeString, "ContainerEnd")
}

func (this *stringBuilder) BuildFromMarker(id interface{}) {
	panic("TODO: stringBuilder.Marker")
}

func (this *stringBuilder) BuildFromReference(id interface{}) {
	panic("TODO: stringBuilder.Reference")
}

func (this *stringBuilder) PrepareForListContents() {
	builderPanicBadEvent(this, typeString, "PrepareForListContents")
}

func (this *stringBuilder) PrepareForMapContents() {
	builderPanicBadEvent(this, typeString, "PrepareForMapContents")
}

func (this *stringBuilder) NotifyChildContainerFinished(value reflect.Value) {
	builderPanicBadEvent(this, typeString, "NotifyChildContainerFinished")
}
