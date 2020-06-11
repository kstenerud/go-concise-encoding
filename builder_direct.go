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

// The direct builder has an unambiguous direct mapping from build event to
// a non-pointer destination type (for example, a bool is always a bool).
type directBuilder struct {
	// Const data
	dstType reflect.Type
}

func newDirectBuilder(dstType reflect.Type) ObjectBuilder {
	return &directBuilder{
		dstType: dstType,
	}
}

func (this *directBuilder) IsContainerOnly() bool {
	return false
}

func (this *directBuilder) PostCacheInitBuilder() {
}

func (this *directBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder) ObjectBuilder {
	return this
}

func (this *directBuilder) BuildFromNil(dst reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "Nil")
}

func (this *directBuilder) BuildFromBool(value bool, dst reflect.Value) {
	dst.SetBool(value)
}

func (this *directBuilder) BuildFromInt(value int64, dst reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "Int")
}

func (this *directBuilder) BuildFromUint(value uint64, dst reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "Uint")
}

func (this *directBuilder) BuildFromBigInt(value *big.Int, dst reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "BigInt")
}

func (this *directBuilder) BuildFromFloat(value float64, dst reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "Float")
}

func (this *directBuilder) BuildFromBigFloat(value *big.Float, dst reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "BigFloat")
}

func (this *directBuilder) BuildFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "DecimalFloat")
}

func (this *directBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "BigDecimalFloat")
}

func (this *directBuilder) BuildFromUUID(value []byte, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (this *directBuilder) BuildFromString(value string, dst reflect.Value) {
	dst.SetString(value)
}

func (this *directBuilder) BuildFromBytes(value []byte, dst reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "Bytes")
}

func (this *directBuilder) BuildFromURI(value *url.URL, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value).Elem())
}

func (this *directBuilder) BuildFromTime(value time.Time, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (this *directBuilder) BuildFromCompactTime(value *compact_time.Time, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (this *directBuilder) BuildBeginList() {
	builderPanicBadEvent(this, this.dstType, "List")
}

func (this *directBuilder) BuildBeginMap() {
	builderPanicBadEvent(this, this.dstType, "Map")
}

func (this *directBuilder) BuildEndContainer() {
	builderPanicBadEvent(this, this.dstType, "ContainerEnd")
}

func (this *directBuilder) BuildFromMarker(id interface{}) {
	panic("TODO: directBuilder.BuildFromMarker")
}

func (this *directBuilder) BuildFromReference(id interface{}) {
	panic("TODO: directBuilder.BuildFromReference")
}

func (this *directBuilder) PrepareForListContents() {
	builderPanicBadEvent(this, this.dstType, "PrepareForListContents")
}

func (this *directBuilder) PrepareForMapContents() {
	builderPanicBadEvent(this, this.dstType, "PrepareForMapContents")
}

func (this *directBuilder) NotifyChildContainerFinished(value reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "NotifyChildContainerFinished")
}

// The direct builder has an unambiguous direct mapping from build event to
// a pointer destination type (for example, a *url is always a *url).
type directPtrBuilder struct {
	// Const data
	dstType reflect.Type
}

func newDirectPtrBuilder(dstType reflect.Type) ObjectBuilder {
	return &directPtrBuilder{
		dstType: dstType,
	}
}

func (this *directPtrBuilder) IsContainerOnly() bool {
	return false
}

func (this *directPtrBuilder) PostCacheInitBuilder() {
}

func (this *directPtrBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder) ObjectBuilder {
	return this
}

func (this *directPtrBuilder) BuildFromNil(dst reflect.Value) {
	dst.Set(reflect.Zero(dst.Type()))
}

func (this *directPtrBuilder) BuildFromBool(value bool, dst reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "Bool")
}

func (this *directPtrBuilder) BuildFromInt(value int64, dst reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "Int")
}

func (this *directPtrBuilder) BuildFromUint(value uint64, dst reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "Uint")
}

func (this *directPtrBuilder) BuildFromBigInt(value *big.Int, dst reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "BigInt")
}

func (this *directPtrBuilder) BuildFromFloat(value float64, dst reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "Float")
}

func (this *directPtrBuilder) BuildFromBigFloat(value *big.Float, dst reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "BigFloat")
}

func (this *directPtrBuilder) BuildFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "DecimalFloat")
}

func (this *directPtrBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "BigDecimalFloat")
}

func (this *directPtrBuilder) BuildFromUUID(value []byte, dst reflect.Value) {
	dst.SetBytes(value)
}

func (this *directPtrBuilder) BuildFromString(value string, dst reflect.Value) {
	// String needs special handling since there's no such thing as a nil string
	// in go.
	builderPanicBadEvent(this, this.dstType, "String")
}

func (this *directPtrBuilder) BuildFromBytes(value []byte, dst reflect.Value) {
	dst.SetBytes(value)
}

func (this *directPtrBuilder) BuildFromURI(value *url.URL, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (this *directPtrBuilder) BuildFromTime(value time.Time, dst reflect.Value) {
	// TODO: Should non-pointer stuff be here?
	dst.Set(reflect.ValueOf(value))
}

func (this *directPtrBuilder) BuildFromCompactTime(value *compact_time.Time, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (this *directPtrBuilder) BuildBeginList() {
	builderPanicBadEvent(this, this.dstType, "List")
}

func (this *directPtrBuilder) BuildBeginMap() {
	builderPanicBadEvent(this, this.dstType, "Map")
}

func (this *directPtrBuilder) BuildEndContainer() {
	builderPanicBadEvent(this, this.dstType, "ContainerEnd")
}

func (this *directPtrBuilder) BuildFromMarker(id interface{}) {
	panic("TODO: directPtrBuilder.BuildFromMarker")
}

func (this *directPtrBuilder) BuildFromReference(id interface{}) {
	panic("TODO: directPtrBuilder.BuildFromReference")
}

func (this *directPtrBuilder) PrepareForListContents() {
	builderPanicBadEvent(this, this.dstType, "PrepareForListContents")
}

func (this *directPtrBuilder) PrepareForMapContents() {
	builderPanicBadEvent(this, this.dstType, "PrepareForMapContents")
}

func (this *directPtrBuilder) NotifyChildContainerFinished(value reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "NotifyChildContainerFinished")
}
