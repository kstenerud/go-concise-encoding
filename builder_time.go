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

func (_this *timeBuilder) BuildFromNil(dst reflect.Value) {
	builderPanicBadEvent(_this, typeTime, "Nil")
}

func (_this *timeBuilder) BuildFromBool(value bool, dst reflect.Value) {
	builderPanicBadEvent(_this, typeTime, "Bool")
}

func (_this *timeBuilder) BuildFromInt(value int64, dst reflect.Value) {
	builderPanicBadEvent(_this, typeTime, "Int")
}

func (_this *timeBuilder) BuildFromUint(value uint64, dst reflect.Value) {
	builderPanicBadEvent(_this, typeTime, "Uint")
}

func (_this *timeBuilder) BuildFromBigInt(value *big.Int, dst reflect.Value) {
	builderPanicBadEvent(_this, typeTime, "BigInt")
}

func (_this *timeBuilder) BuildFromFloat(value float64, dst reflect.Value) {
	builderPanicBadEvent(_this, typeTime, "Float")
}

func (_this *timeBuilder) BuildFromBigFloat(value *big.Float, dst reflect.Value) {
	builderPanicBadEvent(_this, typeTime, "BigFloat")
}

func (_this *timeBuilder) BuildFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	builderPanicBadEvent(_this, typeTime, "DecimalFloat")
}

func (_this *timeBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	builderPanicBadEvent(_this, typeTime, "BigDecimalFloat")
}

func (_this *timeBuilder) BuildFromUUID(value []byte, dst reflect.Value) {
	builderPanicBadEvent(_this, typeTime, "UUID")
}

func (_this *timeBuilder) BuildFromString(value string, dst reflect.Value) {
	builderPanicBadEvent(_this, typeTime, "String")
}

func (_this *timeBuilder) BuildFromBytes(value []byte, dst reflect.Value) {
	builderPanicBadEvent(_this, typeTime, "Bytes")
}

func (_this *timeBuilder) BuildFromURI(value *url.URL, dst reflect.Value) {
	builderPanicBadEvent(_this, typeTime, "URI")
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
	builderPanicBadEvent(_this, typeTime, "List")
}

func (_this *timeBuilder) BuildBeginMap() {
	builderPanicBadEvent(_this, typeTime, "Map")
}

func (_this *timeBuilder) BuildEndContainer() {
	builderPanicBadEvent(_this, typeTime, "ContainerEnd")
}

func (_this *timeBuilder) BuildFromMarker(id interface{}) {
	panic("TODO: timeBuilder.Marker")
}

func (_this *timeBuilder) BuildFromReference(id interface{}) {
	panic("TODO: timeBuilder.Reference")
}

func (_this *timeBuilder) PrepareForListContents() {
	builderPanicBadEvent(_this, typeTime, "PrepareForListContents")
}

func (_this *timeBuilder) PrepareForMapContents() {
	builderPanicBadEvent(_this, typeTime, "PrepareForMapContents")
}

func (_this *timeBuilder) NotifyChildContainerFinished(value reflect.Value) {
	builderPanicBadEvent(_this, typeTime, "NotifyChildContainerFinished")
}

// Compact Time

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

func (_this *compactTimeBuilder) BuildFromNil(dst reflect.Value) {
	builderPanicBadEvent(_this, typeCompactTime, "Nil")
}

func (_this *compactTimeBuilder) BuildFromBool(value bool, dst reflect.Value) {
	builderPanicBadEvent(_this, typeCompactTime, "Bool")
}

func (_this *compactTimeBuilder) BuildFromInt(value int64, dst reflect.Value) {
	builderPanicBadEvent(_this, typeCompactTime, "Int")
}

func (_this *compactTimeBuilder) BuildFromUint(value uint64, dst reflect.Value) {
	builderPanicBadEvent(_this, typeCompactTime, "Uint")
}

func (_this *compactTimeBuilder) BuildFromBigInt(value *big.Int, dst reflect.Value) {
	builderPanicBadEvent(_this, typeCompactTime, "BigInt")
}

func (_this *compactTimeBuilder) BuildFromFloat(value float64, dst reflect.Value) {
	builderPanicBadEvent(_this, typeCompactTime, "Float")
}

func (_this *compactTimeBuilder) BuildFromBigFloat(value *big.Float, dst reflect.Value) {
	builderPanicBadEvent(_this, typeCompactTime, "BigFloat")
}

func (_this *compactTimeBuilder) BuildFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	builderPanicBadEvent(_this, typeCompactTime, "DecimalFloat")
}

func (_this *compactTimeBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	builderPanicBadEvent(_this, typeCompactTime, "BigDecimalFloat")
}

func (_this *compactTimeBuilder) BuildFromUUID(value []byte, dst reflect.Value) {
	builderPanicBadEvent(_this, typeCompactTime, "UUID")
}

func (_this *compactTimeBuilder) BuildFromString(value string, dst reflect.Value) {
	builderPanicBadEvent(_this, typeCompactTime, "String")
}

func (_this *compactTimeBuilder) BuildFromBytes(value []byte, dst reflect.Value) {
	builderPanicBadEvent(_this, typeCompactTime, "Bytes")
}

func (_this *compactTimeBuilder) BuildFromURI(value *url.URL, dst reflect.Value) {
	builderPanicBadEvent(_this, typeCompactTime, "URI")
}

func (_this *compactTimeBuilder) BuildFromTime(value time.Time, dst reflect.Value) {
	dst.Set(reflect.ValueOf(*compact_time.AsCompactTime(value)))
}

func (_this *compactTimeBuilder) BuildFromCompactTime(value *compact_time.Time, dst reflect.Value) {
	dst.Set(reflect.ValueOf(*value))
}

func (_this *compactTimeBuilder) BuildBeginList() {
	builderPanicBadEvent(_this, typeCompactTime, "List")
}

func (_this *compactTimeBuilder) BuildBeginMap() {
	builderPanicBadEvent(_this, typeCompactTime, "Map")
}

func (_this *compactTimeBuilder) BuildEndContainer() {
	builderPanicBadEvent(_this, typeCompactTime, "ContainerEnd")
}

func (_this *compactTimeBuilder) BuildFromMarker(id interface{}) {
	panic("TODO: compactTimeBuilder.Marker")
}

func (_this *compactTimeBuilder) BuildFromReference(id interface{}) {
	panic("TODO: compactTimeBuilder.Reference")
}

func (_this *compactTimeBuilder) PrepareForListContents() {
	builderPanicBadEvent(_this, typeCompactTime, "PrepareForListContents")
}

func (_this *compactTimeBuilder) PrepareForMapContents() {
	builderPanicBadEvent(_this, typeCompactTime, "PrepareForMapContents")
}

func (_this *compactTimeBuilder) NotifyChildContainerFinished(value reflect.Value) {
	builderPanicBadEvent(_this, typeCompactTime, "NotifyChildContainerFinished")
}

// PCompact Time

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

func (_this *pCompactTimeBuilder) BuildFromNil(dst reflect.Value) {
	dst.Set(reflect.ValueOf((*compact_time.Time)(nil)))
}

func (_this *pCompactTimeBuilder) BuildFromBool(value bool, dst reflect.Value) {
	builderPanicBadEvent(_this, typePCompactTime, "Bool")
}

func (_this *pCompactTimeBuilder) BuildFromInt(value int64, dst reflect.Value) {
	builderPanicBadEvent(_this, typePCompactTime, "Int")
}

func (_this *pCompactTimeBuilder) BuildFromUint(value uint64, dst reflect.Value) {
	builderPanicBadEvent(_this, typePCompactTime, "Uint")
}

func (_this *pCompactTimeBuilder) BuildFromBigInt(value *big.Int, dst reflect.Value) {
	builderPanicBadEvent(_this, typePCompactTime, "BigInt")
}

func (_this *pCompactTimeBuilder) BuildFromFloat(value float64, dst reflect.Value) {
	builderPanicBadEvent(_this, typePCompactTime, "Float")
}

func (_this *pCompactTimeBuilder) BuildFromBigFloat(value *big.Float, dst reflect.Value) {
	builderPanicBadEvent(_this, typePCompactTime, "BigFloat")
}

func (_this *pCompactTimeBuilder) BuildFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	builderPanicBadEvent(_this, typePCompactTime, "DecimalFloat")
}

func (_this *pCompactTimeBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	builderPanicBadEvent(_this, typePCompactTime, "BigDecimalFloat")
}

func (_this *pCompactTimeBuilder) BuildFromUUID(value []byte, dst reflect.Value) {
	builderPanicBadEvent(_this, typePCompactTime, "UUID")
}

func (_this *pCompactTimeBuilder) BuildFromString(value string, dst reflect.Value) {
	builderPanicBadEvent(_this, typePCompactTime, "String")
}

func (_this *pCompactTimeBuilder) BuildFromBytes(value []byte, dst reflect.Value) {
	builderPanicBadEvent(_this, typePCompactTime, "Bytes")
}

func (_this *pCompactTimeBuilder) BuildFromURI(value *url.URL, dst reflect.Value) {
	builderPanicBadEvent(_this, typePCompactTime, "URI")
}

func (_this *pCompactTimeBuilder) BuildFromTime(value time.Time, dst reflect.Value) {
	dst.Set(reflect.ValueOf(compact_time.AsCompactTime(value)))
}

func (_this *pCompactTimeBuilder) BuildFromCompactTime(value *compact_time.Time, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (_this *pCompactTimeBuilder) BuildBeginList() {
	builderPanicBadEvent(_this, typePCompactTime, "List")
}

func (_this *pCompactTimeBuilder) BuildBeginMap() {
	builderPanicBadEvent(_this, typePCompactTime, "Map")
}

func (_this *pCompactTimeBuilder) BuildEndContainer() {
	builderPanicBadEvent(_this, typePCompactTime, "ContainerEnd")
}

func (_this *pCompactTimeBuilder) BuildFromMarker(id interface{}) {
	panic("TODO: pCompactTimeBuilder.Marker")
}

func (_this *pCompactTimeBuilder) BuildFromReference(id interface{}) {
	panic("TODO: pCompactTimeBuilder.Reference")
}

func (_this *pCompactTimeBuilder) PrepareForListContents() {
	builderPanicBadEvent(_this, typePCompactTime, "PrepareForListContents")
}

func (_this *pCompactTimeBuilder) PrepareForMapContents() {
	builderPanicBadEvent(_this, typePCompactTime, "PrepareForMapContents")
}

func (_this *pCompactTimeBuilder) NotifyChildContainerFinished(value reflect.Value) {
	builderPanicBadEvent(_this, typePCompactTime, "NotifyChildContainerFinished")
}
