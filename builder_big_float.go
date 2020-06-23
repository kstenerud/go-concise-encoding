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

type bigFloatBuilder struct {
}

func newBigFloatBuilder() ObjectBuilder {
	return &bigFloatBuilder{}
}

func (_this *bigFloatBuilder) IsContainerOnly() bool {
	return false
}

func (_this *bigFloatBuilder) PostCacheInitBuilder() {
}

func (_this *bigFloatBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder, options *BuilderOptions) ObjectBuilder {
	return _this
}

func (_this *bigFloatBuilder) BuildFromNil(dst reflect.Value) {
	builderPanicBadEvent(_this, typeBigFloat, "Nil")
}

func (_this *bigFloatBuilder) BuildFromBool(value bool, dst reflect.Value) {
	builderPanicBadEvent(_this, typeBigFloat, "Bool")
}

func (_this *bigFloatBuilder) BuildFromInt(value int64, dst reflect.Value) {
	setBigFloatFromInt(value, dst)
}

func (_this *bigFloatBuilder) BuildFromUint(value uint64, dst reflect.Value) {
	setBigFloatFromUint(value, dst)
}

func (_this *bigFloatBuilder) BuildFromBigInt(value *big.Int, dst reflect.Value) {
	setBigFloatFromBigInt(value, dst)
}

func (_this *bigFloatBuilder) BuildFromFloat(value float64, dst reflect.Value) {
	setBigFloatFromFloat(value, dst)
}

func (_this *bigFloatBuilder) BuildFromBigFloat(value *big.Float, dst reflect.Value) {
	dst.Set(reflect.ValueOf(*value))
}

func (_this *bigFloatBuilder) BuildFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	setBigFloatFromDecimalFloat(value, dst)
}

func (_this *bigFloatBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	setBigFloatFromBigDecimalFloat(value, dst)
}

func (_this *bigFloatBuilder) BuildFromUUID(value []byte, dst reflect.Value) {
	builderPanicBadEvent(_this, typeBigFloat, "UUID")
}

func (_this *bigFloatBuilder) BuildFromString(value string, dst reflect.Value) {
	builderPanicBadEvent(_this, typeBigFloat, "String")
}

func (_this *bigFloatBuilder) BuildFromBytes(value []byte, dst reflect.Value) {
	builderPanicBadEvent(_this, typeBigFloat, "Bytes")
}

func (_this *bigFloatBuilder) BuildFromURI(value *url.URL, dst reflect.Value) {
	builderPanicBadEvent(_this, typeBigFloat, "URI")
}

func (_this *bigFloatBuilder) BuildFromTime(value time.Time, dst reflect.Value) {
	builderPanicBadEvent(_this, typeBigFloat, "Time")
}

func (_this *bigFloatBuilder) BuildFromCompactTime(value *compact_time.Time, dst reflect.Value) {
	builderPanicBadEvent(_this, typeBigFloat, "CompactTime")
}

func (_this *bigFloatBuilder) BuildBeginList() {
	builderPanicBadEvent(_this, typeBigFloat, "List")
}

func (_this *bigFloatBuilder) BuildBeginMap() {
	builderPanicBadEvent(_this, typeBigFloat, "Map")
}

func (_this *bigFloatBuilder) BuildEndContainer() {
	builderPanicBadEvent(_this, typeBigFloat, "ContainerEnd")
}

func (_this *bigFloatBuilder) BuildFromMarker(id interface{}) {
	panic("TODO: bigFloatBuilder.Marker")
}

func (_this *bigFloatBuilder) BuildFromReference(id interface{}) {
	panic("TODO: bigFloatBuilder.Reference")
}

func (_this *bigFloatBuilder) PrepareForListContents() {
	builderPanicBadEvent(_this, typeBigFloat, "PrepareForListContents")
}

func (_this *bigFloatBuilder) PrepareForMapContents() {
	builderPanicBadEvent(_this, typeBigFloat, "PrepareForMapContents")
}

func (_this *bigFloatBuilder) NotifyChildContainerFinished(value reflect.Value) {
	builderPanicBadEvent(_this, typeBigFloat, "NotifyChildContainerFinished")
}

type pBigFloatBuilder struct {
}

func newPBigFloatBuilder() ObjectBuilder {
	return &pBigFloatBuilder{}
}

func (_this *pBigFloatBuilder) IsContainerOnly() bool {
	return false
}

func (_this *pBigFloatBuilder) PostCacheInitBuilder() {
}

func (_this *pBigFloatBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder, options *BuilderOptions) ObjectBuilder {
	return _this
}

func (_this *pBigFloatBuilder) BuildFromNil(dst reflect.Value) {
	dst.Set(reflect.ValueOf((*big.Float)(nil)))
}

func (_this *pBigFloatBuilder) BuildFromBool(value bool, dst reflect.Value) {
	builderPanicBadEvent(_this, typePBigFloat, "Bool")
}

func (_this *pBigFloatBuilder) BuildFromInt(value int64, dst reflect.Value) {
	setPBigFloatFromInt(value, dst)
}

func (_this *pBigFloatBuilder) BuildFromUint(value uint64, dst reflect.Value) {
	setPBigFloatFromUint(value, dst)
}

func (_this *pBigFloatBuilder) BuildFromBigInt(value *big.Int, dst reflect.Value) {
	setPBigFloatFromBigInt(value, dst)
}

func (_this *pBigFloatBuilder) BuildFromFloat(value float64, dst reflect.Value) {
	setPBigFloatFromFloat(value, dst)
}

func (_this *pBigFloatBuilder) BuildFromBigFloat(value *big.Float, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (_this *pBigFloatBuilder) BuildFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	setPBigFloatFromDecimalFloat(value, dst)
}

func (_this *pBigFloatBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	setPBigFloatFromBigDecimalFloat(value, dst)
}

func (_this *pBigFloatBuilder) BuildFromUUID(value []byte, dst reflect.Value) {
	builderPanicBadEvent(_this, typePBigFloat, "UUID")
}

func (_this *pBigFloatBuilder) BuildFromString(value string, dst reflect.Value) {
	builderPanicBadEvent(_this, typePBigFloat, "String")
}

func (_this *pBigFloatBuilder) BuildFromBytes(value []byte, dst reflect.Value) {
	builderPanicBadEvent(_this, typePBigFloat, "Bytes")
}

func (_this *pBigFloatBuilder) BuildFromURI(value *url.URL, dst reflect.Value) {
	builderPanicBadEvent(_this, typePBigFloat, "URI")
}

func (_this *pBigFloatBuilder) BuildFromTime(value time.Time, dst reflect.Value) {
	builderPanicBadEvent(_this, typePBigFloat, "Time")
}

func (_this *pBigFloatBuilder) BuildFromCompactTime(value *compact_time.Time, dst reflect.Value) {
	builderPanicBadEvent(_this, typePBigFloat, "CompactTime")
}

func (_this *pBigFloatBuilder) BuildBeginList() {
	builderPanicBadEvent(_this, typePBigFloat, "List")
}

func (_this *pBigFloatBuilder) BuildBeginMap() {
	builderPanicBadEvent(_this, typePBigFloat, "Map")
}

func (_this *pBigFloatBuilder) BuildEndContainer() {
	builderPanicBadEvent(_this, typePBigFloat, "ContainerEnd")
}

func (_this *pBigFloatBuilder) BuildFromMarker(id interface{}) {
	panic("TODO: pBigFloatBuilder.Marker")
}

func (_this *pBigFloatBuilder) BuildFromReference(id interface{}) {
	panic("TODO: pBigFloatBuilder.Reference")
}

func (_this *pBigFloatBuilder) PrepareForListContents() {
	builderPanicBadEvent(_this, typePBigFloat, "PrepareForListContents")
}

func (_this *pBigFloatBuilder) PrepareForMapContents() {
	builderPanicBadEvent(_this, typePBigFloat, "PrepareForMapContents")
}

func (_this *pBigFloatBuilder) NotifyChildContainerFinished(value reflect.Value) {
	builderPanicBadEvent(_this, typePBigFloat, "NotifyChildContainerFinished")
}
