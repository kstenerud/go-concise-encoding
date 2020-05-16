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

type bigFloatBuilder struct {
}

func newBigFloatBuilder() ObjectBuilder {
	return &bigFloatBuilder{}
}

func (this *bigFloatBuilder) PostCacheInitBuilder() {
}

func (this *bigFloatBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder) ObjectBuilder {
	return this
}

func (this *bigFloatBuilder) BuildFromNil(dst reflect.Value) {
	builderPanicBadEvent(this, typeBigFloat, "Nil")
}

func (this *bigFloatBuilder) BuildFromBool(value bool, dst reflect.Value) {
	builderPanicBadEvent(this, typeBigFloat, "Bool")
}

func (this *bigFloatBuilder) BuildFromInt(value int64, dst reflect.Value) {
	setBigFloatFromInt(value, dst)
}

func (this *bigFloatBuilder) BuildFromUint(value uint64, dst reflect.Value) {
	setBigFloatFromUint(value, dst)
}

func (this *bigFloatBuilder) BuildFromBigInt(value *big.Int, dst reflect.Value) {
	setBigFloatFromBigInt(value, dst)
}

func (this *bigFloatBuilder) BuildFromFloat(value float64, dst reflect.Value) {
	setBigFloatFromFloat(value, dst)
}

func (this *bigFloatBuilder) BuildFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	setBigFloatFromDecimalFloat(value, dst)
}

func (this *bigFloatBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	dst.Set(reflect.ValueOf(*value))
}

func (this *bigFloatBuilder) BuildFromString(value string, dst reflect.Value) {
	builderPanicBadEvent(this, typeBigFloat, "String")
}

func (this *bigFloatBuilder) BuildFromBytes(value []byte, dst reflect.Value) {
	builderPanicBadEvent(this, typeBigFloat, "Bytes")
}

func (this *bigFloatBuilder) BuildFromURI(value *url.URL, dst reflect.Value) {
	builderPanicBadEvent(this, typeBigFloat, "URI")
}

func (this *bigFloatBuilder) BuildFromTime(value time.Time, dst reflect.Value) {
	builderPanicBadEvent(this, typeBigFloat, "Time")
}

func (this *bigFloatBuilder) BuildBeginList() {
	builderPanicBadEvent(this, typeBigFloat, "List")
}

func (this *bigFloatBuilder) BuildBeginMap() {
	builderPanicBadEvent(this, typeBigFloat, "Map")
}

func (this *bigFloatBuilder) BuildEndContainer() {
	builderPanicBadEvent(this, typeBigFloat, "ContainerEnd")
}

func (this *bigFloatBuilder) BuildFromMarker(id interface{}) {
	panic("TODO: bigFloatBuilder.Marker")
}

func (this *bigFloatBuilder) BuildFromReference(id interface{}) {
	panic("TODO: bigFloatBuilder.Reference")
}

func (this *bigFloatBuilder) PrepareForListContents() {
}

func (this *bigFloatBuilder) PrepareForMapContents() {
}

func (this *bigFloatBuilder) NotifyChildContainerFinished(value reflect.Value) {
	builderPanicBadEvent(this, typeBigFloat, "NotifyChildContainerFinished")
}

type pBigFloatBuilder struct {
}

func newPBigFloatBuilder() ObjectBuilder {
	return &pBigFloatBuilder{}
}

func (this *pBigFloatBuilder) PostCacheInitBuilder() {
}

func (this *pBigFloatBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder) ObjectBuilder {
	return this
}

func (this *pBigFloatBuilder) BuildFromNil(dst reflect.Value) {
	builderPanicBadEvent(this, typePBigFloat, "Nil")
}

func (this *pBigFloatBuilder) BuildFromBool(value bool, dst reflect.Value) {
	builderPanicBadEvent(this, typePBigFloat, "Bool")
}

func (this *pBigFloatBuilder) BuildFromInt(value int64, dst reflect.Value) {
	setPBigFloatFromInt(value, dst)
}

func (this *pBigFloatBuilder) BuildFromUint(value uint64, dst reflect.Value) {
	setPBigFloatFromUint(value, dst)
}

func (this *pBigFloatBuilder) BuildFromBigInt(value *big.Int, dst reflect.Value) {
	setPBigFloatFromBigInt(value, dst)
}

func (this *pBigFloatBuilder) BuildFromFloat(value float64, dst reflect.Value) {
	setPBigFloatFromFloat(value, dst)
}

func (this *pBigFloatBuilder) BuildFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value.APD()))
}

func (this *pBigFloatBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (this *pBigFloatBuilder) BuildFromString(value string, dst reflect.Value) {
	builderPanicBadEvent(this, typePBigFloat, "String")
}

func (this *pBigFloatBuilder) BuildFromBytes(value []byte, dst reflect.Value) {
	builderPanicBadEvent(this, typePBigFloat, "Bytes")
}

func (this *pBigFloatBuilder) BuildFromURI(value *url.URL, dst reflect.Value) {
	builderPanicBadEvent(this, typePBigFloat, "URI")
}

func (this *pBigFloatBuilder) BuildFromTime(value time.Time, dst reflect.Value) {
	builderPanicBadEvent(this, typePBigFloat, "Time")
}

func (this *pBigFloatBuilder) BuildBeginList() {
	builderPanicBadEvent(this, typePBigFloat, "List")
}

func (this *pBigFloatBuilder) BuildBeginMap() {
	builderPanicBadEvent(this, typePBigFloat, "Map")
}

func (this *pBigFloatBuilder) BuildEndContainer() {
	builderPanicBadEvent(this, typePBigFloat, "ContainerEnd")
}

func (this *pBigFloatBuilder) BuildFromMarker(id interface{}) {
	panic("TODO: pBigFloatBuilder.Marker")
}

func (this *pBigFloatBuilder) BuildFromReference(id interface{}) {
	panic("TODO: pBigFloatBuilder.Reference")
}

func (this *pBigFloatBuilder) PrepareForListContents() {
}

func (this *pBigFloatBuilder) PrepareForMapContents() {
}

func (this *pBigFloatBuilder) NotifyChildContainerFinished(value reflect.Value) {
	builderPanicBadEvent(this, typePBigFloat, "NotifyChildContainerFinished")
}
