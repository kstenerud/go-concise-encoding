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

package builder

import (
	"fmt"
	"math/big"
	"net/url"
	"reflect"
	"time"

	"github.com/kstenerud/go-concise-encoding/internal/common"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
)

type bigFloatBuilder struct {
}

func newBigFloatBuilder() ObjectBuilder {
	return &bigFloatBuilder{}
}

func (_this *bigFloatBuilder) String() string {
	return fmt.Sprintf("%v", reflect.TypeOf(_this))
}

func (_this *bigFloatBuilder) IsContainerOnly() bool {
	return false
}

func (_this *bigFloatBuilder) PostCacheInitBuilder() {
}

func (_this *bigFloatBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder, options *BuilderOptions) ObjectBuilder {
	return _this
}

func (_this *bigFloatBuilder) SetParent(parent ObjectBuilder) {
}

func (_this *bigFloatBuilder) BuildFromNil(dst reflect.Value) {
	builderPanicBadEventType(_this, common.TypeBigFloat, "Nil")
}

func (_this *bigFloatBuilder) BuildFromBool(value bool, dst reflect.Value) {
	builderPanicBadEventType(_this, common.TypeBigFloat, "Bool")
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
	builderPanicBadEventType(_this, common.TypeBigFloat, "UUID")
}

func (_this *bigFloatBuilder) BuildFromString(value string, dst reflect.Value) {
	builderPanicBadEventType(_this, common.TypeBigFloat, "String")
}

func (_this *bigFloatBuilder) BuildFromBytes(value []byte, dst reflect.Value) {
	builderPanicBadEventType(_this, common.TypeBigFloat, "Bytes")
}

func (_this *bigFloatBuilder) BuildFromURI(value *url.URL, dst reflect.Value) {
	builderPanicBadEventType(_this, common.TypeBigFloat, "URI")
}

func (_this *bigFloatBuilder) BuildFromTime(value time.Time, dst reflect.Value) {
	builderPanicBadEventType(_this, common.TypeBigFloat, "Time")
}

func (_this *bigFloatBuilder) BuildFromCompactTime(value *compact_time.Time, dst reflect.Value) {
	builderPanicBadEventType(_this, common.TypeBigFloat, "CompactTime")
}

func (_this *bigFloatBuilder) BuildBeginList() {
	builderPanicBadEventType(_this, common.TypeBigFloat, "List")
}

func (_this *bigFloatBuilder) BuildBeginMap() {
	builderPanicBadEventType(_this, common.TypeBigFloat, "Map")
}

func (_this *bigFloatBuilder) BuildEndContainer() {
	builderPanicBadEventType(_this, common.TypeBigFloat, "ContainerEnd")
}

func (_this *bigFloatBuilder) BuildBeginMarker(id interface{}) {
	panic("TODO: bigFloatBuilder.Marker")
}

func (_this *bigFloatBuilder) BuildFromReference(id interface{}) {
	panic("TODO: bigFloatBuilder.Reference")
}

func (_this *bigFloatBuilder) PrepareForListContents() {
	builderPanicBadEventType(_this, common.TypeBigFloat, "PrepareForListContents")
}

func (_this *bigFloatBuilder) PrepareForMapContents() {
	builderPanicBadEventType(_this, common.TypeBigFloat, "PrepareForMapContents")
}

func (_this *bigFloatBuilder) NotifyChildContainerFinished(value reflect.Value) {
	builderPanicBadEventType(_this, common.TypeBigFloat, "NotifyChildContainerFinished")
}

// ============================================================================

type pBigFloatBuilder struct {
}

func newPBigFloatBuilder() ObjectBuilder {
	return &pBigFloatBuilder{}
}

func (_this *pBigFloatBuilder) String() string {
	return fmt.Sprintf("%v", reflect.TypeOf(_this))
}

func (_this *pBigFloatBuilder) IsContainerOnly() bool {
	return false
}

func (_this *pBigFloatBuilder) PostCacheInitBuilder() {
}

func (_this *pBigFloatBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder, options *BuilderOptions) ObjectBuilder {
	return _this
}

func (_this *pBigFloatBuilder) SetParent(parent ObjectBuilder) {
}

func (_this *pBigFloatBuilder) BuildFromNil(dst reflect.Value) {
	dst.Set(reflect.ValueOf((*big.Float)(nil)))
}

func (_this *pBigFloatBuilder) BuildFromBool(value bool, dst reflect.Value) {
	builderPanicBadEventType(_this, common.TypePBigFloat, "Bool")
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
	builderPanicBadEventType(_this, common.TypePBigFloat, "UUID")
}

func (_this *pBigFloatBuilder) BuildFromString(value string, dst reflect.Value) {
	builderPanicBadEventType(_this, common.TypePBigFloat, "String")
}

func (_this *pBigFloatBuilder) BuildFromBytes(value []byte, dst reflect.Value) {
	builderPanicBadEventType(_this, common.TypePBigFloat, "Bytes")
}

func (_this *pBigFloatBuilder) BuildFromURI(value *url.URL, dst reflect.Value) {
	builderPanicBadEventType(_this, common.TypePBigFloat, "URI")
}

func (_this *pBigFloatBuilder) BuildFromTime(value time.Time, dst reflect.Value) {
	builderPanicBadEventType(_this, common.TypePBigFloat, "Time")
}

func (_this *pBigFloatBuilder) BuildFromCompactTime(value *compact_time.Time, dst reflect.Value) {
	builderPanicBadEventType(_this, common.TypePBigFloat, "CompactTime")
}

func (_this *pBigFloatBuilder) BuildBeginList() {
	builderPanicBadEventType(_this, common.TypePBigFloat, "List")
}

func (_this *pBigFloatBuilder) BuildBeginMap() {
	builderPanicBadEventType(_this, common.TypePBigFloat, "Map")
}

func (_this *pBigFloatBuilder) BuildEndContainer() {
	builderPanicBadEventType(_this, common.TypePBigFloat, "ContainerEnd")
}

func (_this *pBigFloatBuilder) BuildBeginMarker(id interface{}) {
	panic("TODO: pBigFloatBuilder.Marker")
}

func (_this *pBigFloatBuilder) BuildFromReference(id interface{}) {
	panic("TODO: pBigFloatBuilder.Reference")
}

func (_this *pBigFloatBuilder) PrepareForListContents() {
	builderPanicBadEventType(_this, common.TypePBigFloat, "PrepareForListContents")
}

func (_this *pBigFloatBuilder) PrepareForMapContents() {
	builderPanicBadEventType(_this, common.TypePBigFloat, "PrepareForMapContents")
}

func (_this *pBigFloatBuilder) NotifyChildContainerFinished(value reflect.Value) {
	builderPanicBadEventType(_this, common.TypePBigFloat, "NotifyChildContainerFinished")
}
