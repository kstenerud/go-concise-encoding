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
	"math/big"
	"net/url"
	"reflect"
	"time"

	"github.com/kstenerud/go-concise-encoding/internal/common"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
)

// Go Time

type timeBuilder struct {
}

var globalTimeBuilder timeBuilder

func newTimeBuilder() ObjectBuilder {
	return &globalTimeBuilder
}

func (_this *timeBuilder) IsContainerOnly() bool {
	return false
}

func (_this *timeBuilder) PostCacheInitBuilder() {
}

func (_this *timeBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder, options *BuilderOptions) ObjectBuilder {
	return _this
}

func (_this *timeBuilder) SetParent(parent ObjectBuilder) {
}

func (_this *timeBuilder) BuildFromNil(dst reflect.Value) {
	builderPanicBadEventType(_this, common.TypeTime, "Nil")
}

func (_this *timeBuilder) BuildFromBool(value bool, dst reflect.Value) {
	builderPanicBadEventType(_this, common.TypeTime, "Bool")
}

func (_this *timeBuilder) BuildFromInt(value int64, dst reflect.Value) {
	builderPanicBadEventType(_this, common.TypeTime, "Int")
}

func (_this *timeBuilder) BuildFromUint(value uint64, dst reflect.Value) {
	builderPanicBadEventType(_this, common.TypeTime, "Uint")
}

func (_this *timeBuilder) BuildFromBigInt(value *big.Int, dst reflect.Value) {
	builderPanicBadEventType(_this, common.TypeTime, "BigInt")
}

func (_this *timeBuilder) BuildFromFloat(value float64, dst reflect.Value) {
	builderPanicBadEventType(_this, common.TypeTime, "Float")
}

func (_this *timeBuilder) BuildFromBigFloat(value *big.Float, dst reflect.Value) {
	builderPanicBadEventType(_this, common.TypeTime, "BigFloat")
}

func (_this *timeBuilder) BuildFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	builderPanicBadEventType(_this, common.TypeTime, "DecimalFloat")
}

func (_this *timeBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	builderPanicBadEventType(_this, common.TypeTime, "BigDecimalFloat")
}

func (_this *timeBuilder) BuildFromUUID(value []byte, dst reflect.Value) {
	builderPanicBadEventType(_this, common.TypeTime, "UUID")
}

func (_this *timeBuilder) BuildFromString(value string, dst reflect.Value) {
	builderPanicBadEventType(_this, common.TypeTime, "String")
}

func (_this *timeBuilder) BuildFromBytes(value []byte, dst reflect.Value) {
	builderPanicBadEventType(_this, common.TypeTime, "Bytes")
}

func (_this *timeBuilder) BuildFromURI(value *url.URL, dst reflect.Value) {
	builderPanicBadEventType(_this, common.TypeTime, "URI")
}

func (_this *timeBuilder) BuildFromTime(value time.Time, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (_this *timeBuilder) BuildFromCompactTime(value *compact_time.Time, dst reflect.Value) {
	v, err := value.AsGoTime()
	if err != nil {
		panic(err)
	}
	dst.Set(reflect.ValueOf(v))
}

func (_this *timeBuilder) BuildBeginList() {
	builderPanicBadEventType(_this, common.TypeTime, "List")
}

func (_this *timeBuilder) BuildBeginMap() {
	builderPanicBadEventType(_this, common.TypeTime, "Map")
}

func (_this *timeBuilder) BuildEndContainer() {
	builderPanicBadEventType(_this, common.TypeTime, "ContainerEnd")
}

func (_this *timeBuilder) BuildBeginMarker(id interface{}) {
	panic("TODO: timeBuilder.Marker")
}

func (_this *timeBuilder) BuildFromReference(id interface{}) {
	panic("TODO: timeBuilder.Reference")
}

func (_this *timeBuilder) PrepareForListContents() {
	builderPanicBadEventType(_this, common.TypeTime, "PrepareForListContents")
}

func (_this *timeBuilder) PrepareForMapContents() {
	builderPanicBadEventType(_this, common.TypeTime, "PrepareForMapContents")
}

func (_this *timeBuilder) NotifyChildContainerFinished(value reflect.Value) {
	builderPanicBadEventType(_this, common.TypeTime, "NotifyChildContainerFinished")
}

// ============================================================================

type compactTimeBuilder struct {
}

var globalCompactTimeBuilder compactTimeBuilder

func newCompactTimeBuilder() ObjectBuilder {
	return &globalCompactTimeBuilder
}

func (_this *compactTimeBuilder) IsContainerOnly() bool {
	return false
}

func (_this *compactTimeBuilder) PostCacheInitBuilder() {
}

func (_this *compactTimeBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder, options *BuilderOptions) ObjectBuilder {
	return _this
}

func (_this *compactTimeBuilder) SetParent(parent ObjectBuilder) {
}

func (_this *compactTimeBuilder) BuildFromNil(dst reflect.Value) {
	builderPanicBadEventType(_this, common.TypeCompactTime, "Nil")
}

func (_this *compactTimeBuilder) BuildFromBool(value bool, dst reflect.Value) {
	builderPanicBadEventType(_this, common.TypeCompactTime, "Bool")
}

func (_this *compactTimeBuilder) BuildFromInt(value int64, dst reflect.Value) {
	builderPanicBadEventType(_this, common.TypeCompactTime, "Int")
}

func (_this *compactTimeBuilder) BuildFromUint(value uint64, dst reflect.Value) {
	builderPanicBadEventType(_this, common.TypeCompactTime, "Uint")
}

func (_this *compactTimeBuilder) BuildFromBigInt(value *big.Int, dst reflect.Value) {
	builderPanicBadEventType(_this, common.TypeCompactTime, "BigInt")
}

func (_this *compactTimeBuilder) BuildFromFloat(value float64, dst reflect.Value) {
	builderPanicBadEventType(_this, common.TypeCompactTime, "Float")
}

func (_this *compactTimeBuilder) BuildFromBigFloat(value *big.Float, dst reflect.Value) {
	builderPanicBadEventType(_this, common.TypeCompactTime, "BigFloat")
}

func (_this *compactTimeBuilder) BuildFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	builderPanicBadEventType(_this, common.TypeCompactTime, "DecimalFloat")
}

func (_this *compactTimeBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	builderPanicBadEventType(_this, common.TypeCompactTime, "BigDecimalFloat")
}

func (_this *compactTimeBuilder) BuildFromUUID(value []byte, dst reflect.Value) {
	builderPanicBadEventType(_this, common.TypeCompactTime, "UUID")
}

func (_this *compactTimeBuilder) BuildFromString(value string, dst reflect.Value) {
	builderPanicBadEventType(_this, common.TypeCompactTime, "String")
}

func (_this *compactTimeBuilder) BuildFromBytes(value []byte, dst reflect.Value) {
	builderPanicBadEventType(_this, common.TypeCompactTime, "Bytes")
}

func (_this *compactTimeBuilder) BuildFromURI(value *url.URL, dst reflect.Value) {
	builderPanicBadEventType(_this, common.TypeCompactTime, "URI")
}

func (_this *compactTimeBuilder) BuildFromTime(value time.Time, dst reflect.Value) {
	dst.Set(reflect.ValueOf(*compact_time.AsCompactTime(value)))
}

func (_this *compactTimeBuilder) BuildFromCompactTime(value *compact_time.Time, dst reflect.Value) {
	dst.Set(reflect.ValueOf(*value))
}

func (_this *compactTimeBuilder) BuildBeginList() {
	builderPanicBadEventType(_this, common.TypeCompactTime, "List")
}

func (_this *compactTimeBuilder) BuildBeginMap() {
	builderPanicBadEventType(_this, common.TypeCompactTime, "Map")
}

func (_this *compactTimeBuilder) BuildEndContainer() {
	builderPanicBadEventType(_this, common.TypeCompactTime, "ContainerEnd")
}

func (_this *compactTimeBuilder) BuildBeginMarker(id interface{}) {
	panic("TODO: compactTimeBuilder.Marker")
}

func (_this *compactTimeBuilder) BuildFromReference(id interface{}) {
	panic("TODO: compactTimeBuilder.Reference")
}

func (_this *compactTimeBuilder) PrepareForListContents() {
	builderPanicBadEventType(_this, common.TypeCompactTime, "PrepareForListContents")
}

func (_this *compactTimeBuilder) PrepareForMapContents() {
	builderPanicBadEventType(_this, common.TypeCompactTime, "PrepareForMapContents")
}

func (_this *compactTimeBuilder) NotifyChildContainerFinished(value reflect.Value) {
	builderPanicBadEventType(_this, common.TypeCompactTime, "NotifyChildContainerFinished")
}

// ============================================================================

type pCompactTimeBuilder struct {
}

var globalPCompactTimeBuilder pCompactTimeBuilder

func newPCompactTimeBuilder() ObjectBuilder {
	return &globalPCompactTimeBuilder
}

func (_this *pCompactTimeBuilder) IsContainerOnly() bool {
	return false
}

func (_this *pCompactTimeBuilder) PostCacheInitBuilder() {
}

func (_this *pCompactTimeBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder, options *BuilderOptions) ObjectBuilder {
	return _this
}

func (_this *pCompactTimeBuilder) SetParent(parent ObjectBuilder) {
}

func (_this *pCompactTimeBuilder) BuildFromNil(dst reflect.Value) {
	dst.Set(reflect.ValueOf((*compact_time.Time)(nil)))
}

func (_this *pCompactTimeBuilder) BuildFromBool(value bool, dst reflect.Value) {
	builderPanicBadEventType(_this, common.TypePCompactTime, "Bool")
}

func (_this *pCompactTimeBuilder) BuildFromInt(value int64, dst reflect.Value) {
	builderPanicBadEventType(_this, common.TypePCompactTime, "Int")
}

func (_this *pCompactTimeBuilder) BuildFromUint(value uint64, dst reflect.Value) {
	builderPanicBadEventType(_this, common.TypePCompactTime, "Uint")
}

func (_this *pCompactTimeBuilder) BuildFromBigInt(value *big.Int, dst reflect.Value) {
	builderPanicBadEventType(_this, common.TypePCompactTime, "BigInt")
}

func (_this *pCompactTimeBuilder) BuildFromFloat(value float64, dst reflect.Value) {
	builderPanicBadEventType(_this, common.TypePCompactTime, "Float")
}

func (_this *pCompactTimeBuilder) BuildFromBigFloat(value *big.Float, dst reflect.Value) {
	builderPanicBadEventType(_this, common.TypePCompactTime, "BigFloat")
}

func (_this *pCompactTimeBuilder) BuildFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	builderPanicBadEventType(_this, common.TypePCompactTime, "DecimalFloat")
}

func (_this *pCompactTimeBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	builderPanicBadEventType(_this, common.TypePCompactTime, "BigDecimalFloat")
}

func (_this *pCompactTimeBuilder) BuildFromUUID(value []byte, dst reflect.Value) {
	builderPanicBadEventType(_this, common.TypePCompactTime, "UUID")
}

func (_this *pCompactTimeBuilder) BuildFromString(value string, dst reflect.Value) {
	builderPanicBadEventType(_this, common.TypePCompactTime, "String")
}

func (_this *pCompactTimeBuilder) BuildFromBytes(value []byte, dst reflect.Value) {
	builderPanicBadEventType(_this, common.TypePCompactTime, "Bytes")
}

func (_this *pCompactTimeBuilder) BuildFromURI(value *url.URL, dst reflect.Value) {
	builderPanicBadEventType(_this, common.TypePCompactTime, "URI")
}

func (_this *pCompactTimeBuilder) BuildFromTime(value time.Time, dst reflect.Value) {
	dst.Set(reflect.ValueOf(compact_time.AsCompactTime(value)))
}

func (_this *pCompactTimeBuilder) BuildFromCompactTime(value *compact_time.Time, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (_this *pCompactTimeBuilder) BuildBeginList() {
	builderPanicBadEventType(_this, common.TypePCompactTime, "List")
}

func (_this *pCompactTimeBuilder) BuildBeginMap() {
	builderPanicBadEventType(_this, common.TypePCompactTime, "Map")
}

func (_this *pCompactTimeBuilder) BuildEndContainer() {
	builderPanicBadEventType(_this, common.TypePCompactTime, "ContainerEnd")
}

func (_this *pCompactTimeBuilder) BuildBeginMarker(id interface{}) {
	panic("TODO: pCompactTimeBuilder.Marker")
}

func (_this *pCompactTimeBuilder) BuildFromReference(id interface{}) {
	panic("TODO: pCompactTimeBuilder.Reference")
}

func (_this *pCompactTimeBuilder) PrepareForListContents() {
	builderPanicBadEventType(_this, common.TypePCompactTime, "PrepareForListContents")
}

func (_this *pCompactTimeBuilder) PrepareForMapContents() {
	builderPanicBadEventType(_this, common.TypePCompactTime, "PrepareForMapContents")
}

func (_this *pCompactTimeBuilder) NotifyChildContainerFinished(value reflect.Value) {
	builderPanicBadEventType(_this, common.TypePCompactTime, "NotifyChildContainerFinished")
}
