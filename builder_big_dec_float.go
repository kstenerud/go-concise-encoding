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

type bigDecimalFloatBuilder struct {
}

func newBigDecimalFloatBuilder() ObjectBuilder {
	return &bigDecimalFloatBuilder{}
}

func (this *bigDecimalFloatBuilder) PostCacheInitBuilder() {
}

func (this *bigDecimalFloatBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder) ObjectBuilder {
	return this
}

func (this *bigDecimalFloatBuilder) BuildFromNil(dst reflect.Value) {
	builderPanicBadEvent(this, typeBigDecimalFloat, "Nil")
}

func (this *bigDecimalFloatBuilder) BuildFromBool(value bool, dst reflect.Value) {
	builderPanicBadEvent(this, typeBigDecimalFloat, "Bool")
}

func (this *bigDecimalFloatBuilder) BuildFromInt(value int64, dst reflect.Value) {
	setBigDecimalFloatFromInt(value, dst)
}

func (this *bigDecimalFloatBuilder) BuildFromUint(value uint64, dst reflect.Value) {
	setBigDecimalFloatFromUint(value, dst)
}

func (this *bigDecimalFloatBuilder) BuildFromBigInt(value *big.Int, dst reflect.Value) {
	setBigDecimalFloatFromBigInt(value, dst)
}

func (this *bigDecimalFloatBuilder) BuildFromFloat(value float64, dst reflect.Value) {
	setBigDecimalFloatFromFloat(value, dst)
}

func (this *bigDecimalFloatBuilder) BuildFromBigFloat(value *big.Float, dst reflect.Value) {
	setBigDecimalFloatFromBigFloat(value, dst)
}

func (this *bigDecimalFloatBuilder) BuildFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	setBigDecimalFloatFromDecimalFloat(value, dst)
}

func (this *bigDecimalFloatBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	dst.Set(reflect.ValueOf(*value))
}

func (this *bigDecimalFloatBuilder) BuildFromUUID(value []byte, dst reflect.Value) {
	builderPanicBadEvent(this, typeBigDecimalFloat, "UUID")
}

func (this *bigDecimalFloatBuilder) BuildFromString(value string, dst reflect.Value) {
	builderPanicBadEvent(this, typeBigDecimalFloat, "String")
}

func (this *bigDecimalFloatBuilder) BuildFromBytes(value []byte, dst reflect.Value) {
	builderPanicBadEvent(this, typeBigDecimalFloat, "Bytes")
}

func (this *bigDecimalFloatBuilder) BuildFromURI(value *url.URL, dst reflect.Value) {
	builderPanicBadEvent(this, typeBigDecimalFloat, "URI")
}

func (this *bigDecimalFloatBuilder) BuildFromTime(value time.Time, dst reflect.Value) {
	builderPanicBadEvent(this, typeBigDecimalFloat, "Time")
}

func (this *bigDecimalFloatBuilder) BuildBeginList() {
	builderPanicBadEvent(this, typeBigDecimalFloat, "List")
}

func (this *bigDecimalFloatBuilder) BuildBeginMap() {
	builderPanicBadEvent(this, typeBigDecimalFloat, "Map")
}

func (this *bigDecimalFloatBuilder) BuildEndContainer() {
	builderPanicBadEvent(this, typeBigDecimalFloat, "ContainerEnd")
}

func (this *bigDecimalFloatBuilder) BuildFromMarker(id interface{}) {
	panic("TODO: bigDecimalFloatBuilder.Marker")
}

func (this *bigDecimalFloatBuilder) BuildFromReference(id interface{}) {
	panic("TODO: bigDecimalFloatBuilder.Reference")
}

func (this *bigDecimalFloatBuilder) PrepareForListContents() {
}

func (this *bigDecimalFloatBuilder) PrepareForMapContents() {
}

func (this *bigDecimalFloatBuilder) NotifyChildContainerFinished(value reflect.Value) {
	builderPanicBadEvent(this, typeBigDecimalFloat, "NotifyChildContainerFinished")
}

type pBigDecimalFloatBuilder struct {
}

func newPBigDecimalFloatBuilder() ObjectBuilder {
	return &pBigDecimalFloatBuilder{}
}

func (this *pBigDecimalFloatBuilder) PostCacheInitBuilder() {
}

func (this *pBigDecimalFloatBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder) ObjectBuilder {
	return this
}

func (this *pBigDecimalFloatBuilder) BuildFromNil(dst reflect.Value) {
	dst.Set(reflect.ValueOf((*apd.Decimal)(nil)))
}

func (this *pBigDecimalFloatBuilder) BuildFromBool(value bool, dst reflect.Value) {
	builderPanicBadEvent(this, typePBigDecimalFloat, "Bool")
}

func (this *pBigDecimalFloatBuilder) BuildFromInt(value int64, dst reflect.Value) {
	setPBigDecimalFloatFromInt(value, dst)
}

func (this *pBigDecimalFloatBuilder) BuildFromUint(value uint64, dst reflect.Value) {
	setPBigDecimalFloatFromUint(value, dst)
}

func (this *pBigDecimalFloatBuilder) BuildFromBigInt(value *big.Int, dst reflect.Value) {
	setPBigDecimalFloatFromBigInt(value, dst)
}

func (this *pBigDecimalFloatBuilder) BuildFromFloat(value float64, dst reflect.Value) {
	setPBigDecimalFloatFromFloat(value, dst)
}

func (this *pBigDecimalFloatBuilder) BuildFromBigFloat(value *big.Float, dst reflect.Value) {
	setPBigDecimalFloatFromBigFloat(value, dst)
}

func (this *pBigDecimalFloatBuilder) BuildFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value.APD()))
}

func (this *pBigDecimalFloatBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (this *pBigDecimalFloatBuilder) BuildFromUUID(value []byte, dst reflect.Value) {
	builderPanicBadEvent(this, typePBigDecimalFloat, "UUID")
}

func (this *pBigDecimalFloatBuilder) BuildFromString(value string, dst reflect.Value) {
	builderPanicBadEvent(this, typePBigDecimalFloat, "String")
}

func (this *pBigDecimalFloatBuilder) BuildFromBytes(value []byte, dst reflect.Value) {
	builderPanicBadEvent(this, typePBigDecimalFloat, "Bytes")
}

func (this *pBigDecimalFloatBuilder) BuildFromURI(value *url.URL, dst reflect.Value) {
	builderPanicBadEvent(this, typePBigDecimalFloat, "URI")
}

func (this *pBigDecimalFloatBuilder) BuildFromTime(value time.Time, dst reflect.Value) {
	builderPanicBadEvent(this, typePBigDecimalFloat, "Time")
}

func (this *pBigDecimalFloatBuilder) BuildBeginList() {
	builderPanicBadEvent(this, typePBigDecimalFloat, "List")
}

func (this *pBigDecimalFloatBuilder) BuildBeginMap() {
	builderPanicBadEvent(this, typePBigDecimalFloat, "Map")
}

func (this *pBigDecimalFloatBuilder) BuildEndContainer() {
	builderPanicBadEvent(this, typePBigDecimalFloat, "ContainerEnd")
}

func (this *pBigDecimalFloatBuilder) BuildFromMarker(id interface{}) {
	panic("TODO: pBigDecimalFloatBuilder.Marker")
}

func (this *pBigDecimalFloatBuilder) BuildFromReference(id interface{}) {
	panic("TODO: pBigDecimalFloatBuilder.Reference")
}

func (this *pBigDecimalFloatBuilder) PrepareForListContents() {
}

func (this *pBigDecimalFloatBuilder) PrepareForMapContents() {
}

func (this *pBigDecimalFloatBuilder) NotifyChildContainerFinished(value reflect.Value) {
	builderPanicBadEvent(this, typePBigDecimalFloat, "NotifyChildContainerFinished")
}
