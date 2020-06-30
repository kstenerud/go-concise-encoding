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

type bigDecimalFloatBuilder struct {
}

func newBigDecimalFloatBuilder() ObjectBuilder {
	return &bigDecimalFloatBuilder{}
}

func (_this *bigDecimalFloatBuilder) IsContainerOnly() bool {
	return false
}

func (_this *bigDecimalFloatBuilder) PostCacheInitBuilder() {
}

func (_this *bigDecimalFloatBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder, options *BuilderOptions) ObjectBuilder {
	return _this
}

func (_this *bigDecimalFloatBuilder) SetParent(parent ObjectBuilder) {
}

func (_this *bigDecimalFloatBuilder) BuildFromNil(dst reflect.Value) {
	builderPanicBadEvent(_this, typeBigDecimalFloat, "Nil")
}

func (_this *bigDecimalFloatBuilder) BuildFromBool(value bool, dst reflect.Value) {
	builderPanicBadEvent(_this, typeBigDecimalFloat, "Bool")
}

func (_this *bigDecimalFloatBuilder) BuildFromInt(value int64, dst reflect.Value) {
	setBigDecimalFloatFromInt(value, dst)
}

func (_this *bigDecimalFloatBuilder) BuildFromUint(value uint64, dst reflect.Value) {
	setBigDecimalFloatFromUint(value, dst)
}

func (_this *bigDecimalFloatBuilder) BuildFromBigInt(value *big.Int, dst reflect.Value) {
	setBigDecimalFloatFromBigInt(value, dst)
}

func (_this *bigDecimalFloatBuilder) BuildFromFloat(value float64, dst reflect.Value) {
	setBigDecimalFloatFromFloat(value, dst)
}

func (_this *bigDecimalFloatBuilder) BuildFromBigFloat(value *big.Float, dst reflect.Value) {
	setBigDecimalFloatFromBigFloat(value, dst)
}

func (_this *bigDecimalFloatBuilder) BuildFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	setBigDecimalFloatFromDecimalFloat(value, dst)
}

func (_this *bigDecimalFloatBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	dst.Set(reflect.ValueOf(*value))
}

func (_this *bigDecimalFloatBuilder) BuildFromUUID(value []byte, dst reflect.Value) {
	builderPanicBadEvent(_this, typeBigDecimalFloat, "UUID")
}

func (_this *bigDecimalFloatBuilder) BuildFromString(value string, dst reflect.Value) {
	builderPanicBadEvent(_this, typeBigDecimalFloat, "String")
}

func (_this *bigDecimalFloatBuilder) BuildFromBytes(value []byte, dst reflect.Value) {
	builderPanicBadEvent(_this, typeBigDecimalFloat, "Bytes")
}

func (_this *bigDecimalFloatBuilder) BuildFromURI(value *url.URL, dst reflect.Value) {
	builderPanicBadEvent(_this, typeBigDecimalFloat, "URI")
}

func (_this *bigDecimalFloatBuilder) BuildFromTime(value time.Time, dst reflect.Value) {
	builderPanicBadEvent(_this, typeBigDecimalFloat, "Time")
}

func (_this *bigDecimalFloatBuilder) BuildFromCompactTime(value *compact_time.Time, dst reflect.Value) {
	builderPanicBadEvent(_this, typeBigDecimalFloat, "CompactTime")
}

func (_this *bigDecimalFloatBuilder) BuildBeginList() {
	builderPanicBadEvent(_this, typeBigDecimalFloat, "List")
}

func (_this *bigDecimalFloatBuilder) BuildBeginMap() {
	builderPanicBadEvent(_this, typeBigDecimalFloat, "Map")
}

func (_this *bigDecimalFloatBuilder) BuildEndContainer() {
	builderPanicBadEvent(_this, typeBigDecimalFloat, "ContainerEnd")
}

func (_this *bigDecimalFloatBuilder) BuildFromMarker(id interface{}) {
	panic("TODO: bigDecimalFloatBuilder.Marker")
}

func (_this *bigDecimalFloatBuilder) BuildFromReference(id interface{}) {
	panic("TODO: bigDecimalFloatBuilder.Reference")
}

func (_this *bigDecimalFloatBuilder) PrepareForListContents() {
	builderPanicBadEvent(_this, typeBigDecimalFloat, "PrepareForListContents")
}

func (_this *bigDecimalFloatBuilder) PrepareForMapContents() {
	builderPanicBadEvent(_this, typeBigDecimalFloat, "PrepareForMapContents")
}

func (_this *bigDecimalFloatBuilder) NotifyChildContainerFinished(value reflect.Value) {
	builderPanicBadEvent(_this, typeBigDecimalFloat, "NotifyChildContainerFinished")
}

// ============================================================================

type pBigDecimalFloatBuilder struct {
}

func newPBigDecimalFloatBuilder() ObjectBuilder {
	return &pBigDecimalFloatBuilder{}
}

func (_this *pBigDecimalFloatBuilder) IsContainerOnly() bool {
	return false
}

func (_this *pBigDecimalFloatBuilder) PostCacheInitBuilder() {
}

func (_this *pBigDecimalFloatBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder, options *BuilderOptions) ObjectBuilder {
	return _this
}

func (_this *pBigDecimalFloatBuilder) SetParent(parent ObjectBuilder) {
}

func (_this *pBigDecimalFloatBuilder) BuildFromNil(dst reflect.Value) {
	dst.Set(reflect.ValueOf((*apd.Decimal)(nil)))
}

func (_this *pBigDecimalFloatBuilder) BuildFromBool(value bool, dst reflect.Value) {
	builderPanicBadEvent(_this, typePBigDecimalFloat, "Bool")
}

func (_this *pBigDecimalFloatBuilder) BuildFromInt(value int64, dst reflect.Value) {
	setPBigDecimalFloatFromInt(value, dst)
}

func (_this *pBigDecimalFloatBuilder) BuildFromUint(value uint64, dst reflect.Value) {
	setPBigDecimalFloatFromUint(value, dst)
}

func (_this *pBigDecimalFloatBuilder) BuildFromBigInt(value *big.Int, dst reflect.Value) {
	setPBigDecimalFloatFromBigInt(value, dst)
}

func (_this *pBigDecimalFloatBuilder) BuildFromFloat(value float64, dst reflect.Value) {
	setPBigDecimalFloatFromFloat(value, dst)
}

func (_this *pBigDecimalFloatBuilder) BuildFromBigFloat(value *big.Float, dst reflect.Value) {
	setPBigDecimalFloatFromBigFloat(value, dst)
}

func (_this *pBigDecimalFloatBuilder) BuildFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value.APD()))
}

func (_this *pBigDecimalFloatBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (_this *pBigDecimalFloatBuilder) BuildFromUUID(value []byte, dst reflect.Value) {
	builderPanicBadEvent(_this, typePBigDecimalFloat, "UUID")
}

func (_this *pBigDecimalFloatBuilder) BuildFromString(value string, dst reflect.Value) {
	builderPanicBadEvent(_this, typePBigDecimalFloat, "String")
}

func (_this *pBigDecimalFloatBuilder) BuildFromBytes(value []byte, dst reflect.Value) {
	builderPanicBadEvent(_this, typePBigDecimalFloat, "Bytes")
}

func (_this *pBigDecimalFloatBuilder) BuildFromURI(value *url.URL, dst reflect.Value) {
	builderPanicBadEvent(_this, typePBigDecimalFloat, "URI")
}

func (_this *pBigDecimalFloatBuilder) BuildFromTime(value time.Time, dst reflect.Value) {
	builderPanicBadEvent(_this, typePBigDecimalFloat, "Time")
}

func (_this *pBigDecimalFloatBuilder) BuildFromCompactTime(value *compact_time.Time, dst reflect.Value) {
	builderPanicBadEvent(_this, typeBigDecimalFloat, "CompactTime")
}

func (_this *pBigDecimalFloatBuilder) BuildBeginList() {
	builderPanicBadEvent(_this, typePBigDecimalFloat, "List")
}

func (_this *pBigDecimalFloatBuilder) BuildBeginMap() {
	builderPanicBadEvent(_this, typePBigDecimalFloat, "Map")
}

func (_this *pBigDecimalFloatBuilder) BuildEndContainer() {
	builderPanicBadEvent(_this, typePBigDecimalFloat, "ContainerEnd")
}

func (_this *pBigDecimalFloatBuilder) BuildFromMarker(id interface{}) {
	panic("TODO: pBigDecimalFloatBuilder.Marker")
}

func (_this *pBigDecimalFloatBuilder) BuildFromReference(id interface{}) {
	panic("TODO: pBigDecimalFloatBuilder.Reference")
}

func (_this *pBigDecimalFloatBuilder) PrepareForListContents() {
	builderPanicBadEvent(_this, typePBigDecimalFloat, "PrepareForListContents")
}

func (_this *pBigDecimalFloatBuilder) PrepareForMapContents() {
	builderPanicBadEvent(_this, typePBigDecimalFloat, "PrepareForMapContents")
}

func (_this *pBigDecimalFloatBuilder) NotifyChildContainerFinished(value reflect.Value) {
	builderPanicBadEvent(_this, typePBigDecimalFloat, "NotifyChildContainerFinished")
}
