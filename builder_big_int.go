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

type bigIntBuilder struct {
}

func newBigIntBuilder() ObjectBuilder {
	return &bigIntBuilder{}
}

func (_this *bigIntBuilder) IsContainerOnly() bool {
	return false
}

func (_this *bigIntBuilder) PostCacheInitBuilder() {
}

func (_this *bigIntBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder, options *BuilderOptions) ObjectBuilder {
	return _this
}

func (_this *bigIntBuilder) SetParent(parent ObjectBuilder) {
}

func (_this *bigIntBuilder) BuildFromNil(dst reflect.Value) {
	builderPanicBadEvent(_this, typeBigInt, "Nil")
}

func (_this *bigIntBuilder) BuildFromBool(value bool, dst reflect.Value) {
	builderPanicBadEvent(_this, typeBigInt, "Bool")
}

func (_this *bigIntBuilder) BuildFromInt(value int64, dst reflect.Value) {
	setBigIntFromInt(value, dst)
}

func (_this *bigIntBuilder) BuildFromUint(value uint64, dst reflect.Value) {
	setBigIntFromUint(value, dst)
}

func (_this *bigIntBuilder) BuildFromBigInt(value *big.Int, dst reflect.Value) {
	dst.Set(reflect.ValueOf(*value))
}

func (_this *bigIntBuilder) BuildFromFloat(value float64, dst reflect.Value) {
	setBigIntFromFloat(value, dst)
}

func (_this *bigIntBuilder) BuildFromBigFloat(value *big.Float, dst reflect.Value) {
	setBigIntFromBigFloat(value, dst)
}

func (_this *bigIntBuilder) BuildFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	setBigIntFromDecimalFloat(value, dst)
}

func (_this *bigIntBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	setBigIntFromBigDecimalFloat(value, dst)
}

func (_this *bigIntBuilder) BuildFromUUID(value []byte, dst reflect.Value) {
	builderPanicBadEvent(_this, typeBigInt, "UUID")
}

func (_this *bigIntBuilder) BuildFromString(value string, dst reflect.Value) {
	builderPanicBadEvent(_this, typeBigInt, "String")
}

func (_this *bigIntBuilder) BuildFromBytes(value []byte, dst reflect.Value) {
	builderPanicBadEvent(_this, typeBigInt, "Bytes")
}

func (_this *bigIntBuilder) BuildFromURI(value *url.URL, dst reflect.Value) {
	builderPanicBadEvent(_this, typeBigInt, "URI")
}

func (_this *bigIntBuilder) BuildFromTime(value time.Time, dst reflect.Value) {
	builderPanicBadEvent(_this, typeBigInt, "Time")
}

func (_this *bigIntBuilder) BuildFromCompactTime(value *compact_time.Time, dst reflect.Value) {
	builderPanicBadEvent(_this, typeBigInt, "CompactTime")
}

func (_this *bigIntBuilder) BuildBeginList() {
	builderPanicBadEvent(_this, typeBigInt, "List")
}

func (_this *bigIntBuilder) BuildBeginMap() {
	builderPanicBadEvent(_this, typeBigInt, "Map")
}

func (_this *bigIntBuilder) BuildEndContainer() {
	builderPanicBadEvent(_this, typeBigInt, "ContainerEnd")
}

func (_this *bigIntBuilder) BuildFromMarker(id interface{}) {
	panic("TODO: bigIntBuilder.Marker")
}

func (_this *bigIntBuilder) BuildFromReference(id interface{}) {
	panic("TODO: bigIntBuilder.Reference")
}

func (_this *bigIntBuilder) PrepareForListContents() {
	builderPanicBadEvent(_this, typeBigInt, "PrepareForListContents")
}

func (_this *bigIntBuilder) PrepareForMapContents() {
	builderPanicBadEvent(_this, typeBigInt, "PrepareForMapContents")
}

func (_this *bigIntBuilder) NotifyChildContainerFinished(value reflect.Value) {
	builderPanicBadEvent(_this, typeBigInt, "NotifyChildContainerFinished")
}

// ============================================================================

type pBigIntBuilder struct {
}

func newPBigIntBuilder() ObjectBuilder {
	return &pBigIntBuilder{}
}

func (_this *pBigIntBuilder) IsContainerOnly() bool {
	return false
}

func (_this *pBigIntBuilder) PostCacheInitBuilder() {
}

func (_this *pBigIntBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder, options *BuilderOptions) ObjectBuilder {
	return _this
}

func (_this *pBigIntBuilder) SetParent(parent ObjectBuilder) {
}

func (_this *pBigIntBuilder) BuildFromNil(dst reflect.Value) {
	dst.Set(reflect.ValueOf((*big.Int)(nil)))
}

func (_this *pBigIntBuilder) BuildFromBool(value bool, dst reflect.Value) {
	builderPanicBadEvent(_this, typePBigInt, "Bool")
}

func (_this *pBigIntBuilder) BuildFromInt(value int64, dst reflect.Value) {
	setPBigIntFromInt(value, dst)
}

func (_this *pBigIntBuilder) BuildFromUint(value uint64, dst reflect.Value) {
	setPBigIntFromUint(value, dst)
}

func (_this *pBigIntBuilder) BuildFromBigInt(value *big.Int, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (_this *pBigIntBuilder) BuildFromFloat(value float64, dst reflect.Value) {
	setPBigIntFromFloat(value, dst)
}

func (_this *pBigIntBuilder) BuildFromBigFloat(value *big.Float, dst reflect.Value) {
	setPBigIntFromBigFloat(value, dst)
}

func (_this *pBigIntBuilder) BuildFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	setPBigIntFromDecimalFloat(value, dst)
}

func (_this *pBigIntBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	setPBigIntFromBigDecimalFloat(value, dst)
}

func (_this *pBigIntBuilder) BuildFromUUID(value []byte, dst reflect.Value) {
	builderPanicBadEvent(_this, typePBigInt, "UUID")
}

func (_this *pBigIntBuilder) BuildFromString(value string, dst reflect.Value) {
	builderPanicBadEvent(_this, typePBigInt, "String")
}

func (_this *pBigIntBuilder) BuildFromBytes(value []byte, dst reflect.Value) {
	builderPanicBadEvent(_this, typePBigInt, "Bytes")
}

func (_this *pBigIntBuilder) BuildFromURI(value *url.URL, dst reflect.Value) {
	builderPanicBadEvent(_this, typePBigInt, "URI")
}

func (_this *pBigIntBuilder) BuildFromTime(value time.Time, dst reflect.Value) {
	builderPanicBadEvent(_this, typePBigInt, "Time")
}

func (_this *pBigIntBuilder) BuildFromCompactTime(value *compact_time.Time, dst reflect.Value) {
	builderPanicBadEvent(_this, typePBigInt, "CompactTime")
}

func (_this *pBigIntBuilder) BuildBeginList() {
	builderPanicBadEvent(_this, typePBigInt, "List")
}

func (_this *pBigIntBuilder) BuildBeginMap() {
	builderPanicBadEvent(_this, typePBigInt, "Map")
}

func (_this *pBigIntBuilder) BuildEndContainer() {
	builderPanicBadEvent(_this, typePBigInt, "ContainerEnd")
}

func (_this *pBigIntBuilder) BuildFromMarker(id interface{}) {
	panic("TODO: pBigIntBuilder.Marker")
}

func (_this *pBigIntBuilder) BuildFromReference(id interface{}) {
	panic("TODO: pBigIntBuilder.Reference")
}

func (_this *pBigIntBuilder) PrepareForListContents() {
	builderPanicBadEvent(_this, typePBigInt, "PrepareForListContents")
}

func (_this *pBigIntBuilder) PrepareForMapContents() {
	builderPanicBadEvent(_this, typePBigInt, "PrepareForMapContents")
}

func (_this *pBigIntBuilder) NotifyChildContainerFinished(value reflect.Value) {
	builderPanicBadEvent(_this, typePBigInt, "NotifyChildContainerFinished")
}
