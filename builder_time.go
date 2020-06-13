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

func (this *timeBuilder) IsContainerOnly() bool {
	return false
}

func (this *timeBuilder) PostCacheInitBuilder() {
}

func (this *timeBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder) ObjectBuilder {
	return this
}

func (this *timeBuilder) BuildFromNil(dst reflect.Value) {
	builderPanicBadEvent(this, typeTime, "Nil")
}

func (this *timeBuilder) BuildFromBool(value bool, dst reflect.Value) {
	builderPanicBadEvent(this, typeTime, "Bool")
}

func (this *timeBuilder) BuildFromInt(value int64, dst reflect.Value) {
	builderPanicBadEvent(this, typeTime, "Int")
}

func (this *timeBuilder) BuildFromUint(value uint64, dst reflect.Value) {
	builderPanicBadEvent(this, typeTime, "Uint")
}

func (this *timeBuilder) BuildFromBigInt(value *big.Int, dst reflect.Value) {
	builderPanicBadEvent(this, typeTime, "BigInt")
}

func (this *timeBuilder) BuildFromFloat(value float64, dst reflect.Value) {
	builderPanicBadEvent(this, typeTime, "Float")
}

func (this *timeBuilder) BuildFromBigFloat(value *big.Float, dst reflect.Value) {
	builderPanicBadEvent(this, typeTime, "BigFloat")
}

func (this *timeBuilder) BuildFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	builderPanicBadEvent(this, typeTime, "DecimalFloat")
}

func (this *timeBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	builderPanicBadEvent(this, typeTime, "BigDecimalFloat")
}

func (this *timeBuilder) BuildFromUUID(value []byte, dst reflect.Value) {
	builderPanicBadEvent(this, typeTime, "UUID")
}

func (this *timeBuilder) BuildFromString(value string, dst reflect.Value) {
	builderPanicBadEvent(this, typeTime, "String")
}

func (this *timeBuilder) BuildFromBytes(value []byte, dst reflect.Value) {
	builderPanicBadEvent(this, typeTime, "Bytes")
}

func (this *timeBuilder) BuildFromURI(value *url.URL, dst reflect.Value) {
	builderPanicBadEvent(this, typeTime, "URI")
}

func (this *timeBuilder) BuildFromTime(value time.Time, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (this *timeBuilder) BuildFromCompactTime(value *compact_time.Time, dst reflect.Value) {
	v, err := value.AsGoTime()
	if err != nil {
		panic(err)
	}
	dst.Set(reflect.ValueOf(v))
}

func (this *timeBuilder) BuildBeginList() {
	builderPanicBadEvent(this, typeTime, "List")
}

func (this *timeBuilder) BuildBeginMap() {
	builderPanicBadEvent(this, typeTime, "Map")
}

func (this *timeBuilder) BuildEndContainer() {
	builderPanicBadEvent(this, typeTime, "ContainerEnd")
}

func (this *timeBuilder) BuildFromMarker(id interface{}) {
	panic("TODO: timeBuilder.Marker")
}

func (this *timeBuilder) BuildFromReference(id interface{}) {
	panic("TODO: timeBuilder.Reference")
}

func (this *timeBuilder) PrepareForListContents() {
	builderPanicBadEvent(this, typeTime, "PrepareForListContents")
}

func (this *timeBuilder) PrepareForMapContents() {
	builderPanicBadEvent(this, typeTime, "PrepareForMapContents")
}

func (this *timeBuilder) NotifyChildContainerFinished(value reflect.Value) {
	builderPanicBadEvent(this, typeTime, "NotifyChildContainerFinished")
}

// Compact Time

type compactTimeBuilder struct {
}

var globalCompactTimeBuilder compactTimeBuilder

func newCompactTimeBuilder() ObjectBuilder {
	return &globalCompactTimeBuilder
}

func (this *compactTimeBuilder) IsContainerOnly() bool {
	return false
}

func (this *compactTimeBuilder) PostCacheInitBuilder() {
}

func (this *compactTimeBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder) ObjectBuilder {
	return this
}

func (this *compactTimeBuilder) BuildFromNil(dst reflect.Value) {
	builderPanicBadEvent(this, typeCompactTime, "Nil")
}

func (this *compactTimeBuilder) BuildFromBool(value bool, dst reflect.Value) {
	builderPanicBadEvent(this, typeCompactTime, "Bool")
}

func (this *compactTimeBuilder) BuildFromInt(value int64, dst reflect.Value) {
	builderPanicBadEvent(this, typeCompactTime, "Int")
}

func (this *compactTimeBuilder) BuildFromUint(value uint64, dst reflect.Value) {
	builderPanicBadEvent(this, typeCompactTime, "Uint")
}

func (this *compactTimeBuilder) BuildFromBigInt(value *big.Int, dst reflect.Value) {
	builderPanicBadEvent(this, typeCompactTime, "BigInt")
}

func (this *compactTimeBuilder) BuildFromFloat(value float64, dst reflect.Value) {
	builderPanicBadEvent(this, typeCompactTime, "Float")
}

func (this *compactTimeBuilder) BuildFromBigFloat(value *big.Float, dst reflect.Value) {
	builderPanicBadEvent(this, typeCompactTime, "BigFloat")
}

func (this *compactTimeBuilder) BuildFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	builderPanicBadEvent(this, typeCompactTime, "DecimalFloat")
}

func (this *compactTimeBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	builderPanicBadEvent(this, typeCompactTime, "BigDecimalFloat")
}

func (this *compactTimeBuilder) BuildFromUUID(value []byte, dst reflect.Value) {
	builderPanicBadEvent(this, typeCompactTime, "UUID")
}

func (this *compactTimeBuilder) BuildFromString(value string, dst reflect.Value) {
	builderPanicBadEvent(this, typeCompactTime, "String")
}

func (this *compactTimeBuilder) BuildFromBytes(value []byte, dst reflect.Value) {
	builderPanicBadEvent(this, typeCompactTime, "Bytes")
}

func (this *compactTimeBuilder) BuildFromURI(value *url.URL, dst reflect.Value) {
	builderPanicBadEvent(this, typeCompactTime, "URI")
}

func (this *compactTimeBuilder) BuildFromTime(value time.Time, dst reflect.Value) {
	dst.Set(reflect.ValueOf(*compact_time.AsCompactTime(value)))
}

func (this *compactTimeBuilder) BuildFromCompactTime(value *compact_time.Time, dst reflect.Value) {
	dst.Set(reflect.ValueOf(*value))
}

func (this *compactTimeBuilder) BuildBeginList() {
	builderPanicBadEvent(this, typeCompactTime, "List")
}

func (this *compactTimeBuilder) BuildBeginMap() {
	builderPanicBadEvent(this, typeCompactTime, "Map")
}

func (this *compactTimeBuilder) BuildEndContainer() {
	builderPanicBadEvent(this, typeCompactTime, "ContainerEnd")
}

func (this *compactTimeBuilder) BuildFromMarker(id interface{}) {
	panic("TODO: compactTimeBuilder.Marker")
}

func (this *compactTimeBuilder) BuildFromReference(id interface{}) {
	panic("TODO: compactTimeBuilder.Reference")
}

func (this *compactTimeBuilder) PrepareForListContents() {
	builderPanicBadEvent(this, typeCompactTime, "PrepareForListContents")
}

func (this *compactTimeBuilder) PrepareForMapContents() {
	builderPanicBadEvent(this, typeCompactTime, "PrepareForMapContents")
}

func (this *compactTimeBuilder) NotifyChildContainerFinished(value reflect.Value) {
	builderPanicBadEvent(this, typeCompactTime, "NotifyChildContainerFinished")
}

// PCompact Time

type pCompactTimeBuilder struct {
}

var globalPCompactTimeBuilder pCompactTimeBuilder

func newPCompactTimeBuilder() ObjectBuilder {
	return &globalPCompactTimeBuilder
}

func (this *pCompactTimeBuilder) IsContainerOnly() bool {
	return false
}

func (this *pCompactTimeBuilder) PostCacheInitBuilder() {
}

func (this *pCompactTimeBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder) ObjectBuilder {
	return this
}

func (this *pCompactTimeBuilder) BuildFromNil(dst reflect.Value) {
	dst.Set(reflect.ValueOf((*compact_time.Time)(nil)))
}

func (this *pCompactTimeBuilder) BuildFromBool(value bool, dst reflect.Value) {
	builderPanicBadEvent(this, typePCompactTime, "Bool")
}

func (this *pCompactTimeBuilder) BuildFromInt(value int64, dst reflect.Value) {
	builderPanicBadEvent(this, typePCompactTime, "Int")
}

func (this *pCompactTimeBuilder) BuildFromUint(value uint64, dst reflect.Value) {
	builderPanicBadEvent(this, typePCompactTime, "Uint")
}

func (this *pCompactTimeBuilder) BuildFromBigInt(value *big.Int, dst reflect.Value) {
	builderPanicBadEvent(this, typePCompactTime, "BigInt")
}

func (this *pCompactTimeBuilder) BuildFromFloat(value float64, dst reflect.Value) {
	builderPanicBadEvent(this, typePCompactTime, "Float")
}

func (this *pCompactTimeBuilder) BuildFromBigFloat(value *big.Float, dst reflect.Value) {
	builderPanicBadEvent(this, typePCompactTime, "BigFloat")
}

func (this *pCompactTimeBuilder) BuildFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	builderPanicBadEvent(this, typePCompactTime, "DecimalFloat")
}

func (this *pCompactTimeBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	builderPanicBadEvent(this, typePCompactTime, "BigDecimalFloat")
}

func (this *pCompactTimeBuilder) BuildFromUUID(value []byte, dst reflect.Value) {
	builderPanicBadEvent(this, typePCompactTime, "UUID")
}

func (this *pCompactTimeBuilder) BuildFromString(value string, dst reflect.Value) {
	builderPanicBadEvent(this, typePCompactTime, "String")
}

func (this *pCompactTimeBuilder) BuildFromBytes(value []byte, dst reflect.Value) {
	builderPanicBadEvent(this, typePCompactTime, "Bytes")
}

func (this *pCompactTimeBuilder) BuildFromURI(value *url.URL, dst reflect.Value) {
	builderPanicBadEvent(this, typePCompactTime, "URI")
}

func (this *pCompactTimeBuilder) BuildFromTime(value time.Time, dst reflect.Value) {
	dst.Set(reflect.ValueOf(compact_time.AsCompactTime(value)))
}

func (this *pCompactTimeBuilder) BuildFromCompactTime(value *compact_time.Time, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (this *pCompactTimeBuilder) BuildBeginList() {
	builderPanicBadEvent(this, typePCompactTime, "List")
}

func (this *pCompactTimeBuilder) BuildBeginMap() {
	builderPanicBadEvent(this, typePCompactTime, "Map")
}

func (this *pCompactTimeBuilder) BuildEndContainer() {
	builderPanicBadEvent(this, typePCompactTime, "ContainerEnd")
}

func (this *pCompactTimeBuilder) BuildFromMarker(id interface{}) {
	panic("TODO: pCompactTimeBuilder.Marker")
}

func (this *pCompactTimeBuilder) BuildFromReference(id interface{}) {
	panic("TODO: pCompactTimeBuilder.Reference")
}

func (this *pCompactTimeBuilder) PrepareForListContents() {
	builderPanicBadEvent(this, typePCompactTime, "PrepareForListContents")
}

func (this *pCompactTimeBuilder) PrepareForMapContents() {
	builderPanicBadEvent(this, typePCompactTime, "PrepareForMapContents")
}

func (this *pCompactTimeBuilder) NotifyChildContainerFinished(value reflect.Value) {
	builderPanicBadEvent(this, typePCompactTime, "NotifyChildContainerFinished")
}
