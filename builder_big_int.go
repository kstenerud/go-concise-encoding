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

type bigIntBuilder struct {
}

func newBigIntBuilder() ObjectBuilder {
	return &bigIntBuilder{}
}

func (this *bigIntBuilder) PostCacheInitBuilder() {
}

func (this *bigIntBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder) ObjectBuilder {
	return this
}

func (this *bigIntBuilder) BuildFromNil(dst reflect.Value) {
	builderPanicBadEvent(this, typeBigInt, "Nil")
}

func (this *bigIntBuilder) BuildFromBool(value bool, dst reflect.Value) {
	builderPanicBadEvent(this, typeBigInt, "Bool")
}

func (this *bigIntBuilder) BuildFromInt(value int64, dst reflect.Value) {
	setBigIntFromInt(value, dst)
}

func (this *bigIntBuilder) BuildFromUint(value uint64, dst reflect.Value) {
	setBigIntFromUint(value, dst)
}

func (this *bigIntBuilder) BuildFromBigInt(value *big.Int, dst reflect.Value) {
	dst.Set(reflect.ValueOf(*value))
}

func (this *bigIntBuilder) BuildFromFloat(value float64, dst reflect.Value) {
	setBigIntFromFloat(value, dst)
}

func (this *bigIntBuilder) BuildFromBigFloat(value *big.Float, dst reflect.Value) {
	setBigIntFromBigFloat(value, dst)
}

func (this *bigIntBuilder) BuildFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	setBigIntFromDecimalFloat(value, dst)
}

func (this *bigIntBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	setBigIntFromBigDecimalFloat(value, dst)
}

func (this *bigIntBuilder) BuildFromUUID(value []byte, dst reflect.Value) {
	builderPanicBadEvent(this, typeBigInt, "UUID")
}

func (this *bigIntBuilder) BuildFromString(value string, dst reflect.Value) {
	builderPanicBadEvent(this, typeBigInt, "String")
}

func (this *bigIntBuilder) BuildFromBytes(value []byte, dst reflect.Value) {
	builderPanicBadEvent(this, typeBigInt, "Bytes")
}

func (this *bigIntBuilder) BuildFromURI(value *url.URL, dst reflect.Value) {
	builderPanicBadEvent(this, typeBigInt, "URI")
}

func (this *bigIntBuilder) BuildFromTime(value time.Time, dst reflect.Value) {
	builderPanicBadEvent(this, typeBigInt, "Time")
}

func (this *bigIntBuilder) BuildBeginList() {
	builderPanicBadEvent(this, typeBigInt, "List")
}

func (this *bigIntBuilder) BuildBeginMap() {
	builderPanicBadEvent(this, typeBigInt, "Map")
}

func (this *bigIntBuilder) BuildEndContainer() {
	builderPanicBadEvent(this, typeBigInt, "ContainerEnd")
}

func (this *bigIntBuilder) BuildFromMarker(id interface{}) {
	panic("TODO: bigIntBuilder.Marker")
}

func (this *bigIntBuilder) BuildFromReference(id interface{}) {
	panic("TODO: bigIntBuilder.Reference")
}

func (this *bigIntBuilder) PrepareForListContents() {
}

func (this *bigIntBuilder) PrepareForMapContents() {
}

func (this *bigIntBuilder) NotifyChildContainerFinished(value reflect.Value) {
	builderPanicBadEvent(this, typeBigInt, "NotifyChildContainerFinished")
}

type pBigIntBuilder struct {
}

func newPBigIntBuilder() ObjectBuilder {
	return &pBigIntBuilder{}
}

func (this *pBigIntBuilder) PostCacheInitBuilder() {
}

func (this *pBigIntBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder) ObjectBuilder {
	return this
}

func (this *pBigIntBuilder) BuildFromNil(dst reflect.Value) {
	dst.Set(reflect.ValueOf((*big.Int)(nil)))
}

func (this *pBigIntBuilder) BuildFromBool(value bool, dst reflect.Value) {
	builderPanicBadEvent(this, typePBigInt, "Bool")
}

func (this *pBigIntBuilder) BuildFromInt(value int64, dst reflect.Value) {
	setPBigIntFromInt(value, dst)
}

func (this *pBigIntBuilder) BuildFromUint(value uint64, dst reflect.Value) {
	setPBigIntFromUint(value, dst)
}

func (this *pBigIntBuilder) BuildFromBigInt(value *big.Int, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (this *pBigIntBuilder) BuildFromFloat(value float64, dst reflect.Value) {
	setPBigIntFromFloat(value, dst)
}

func (this *pBigIntBuilder) BuildFromBigFloat(value *big.Float, dst reflect.Value) {
	setPBigIntFromBigFloat(value, dst)
}

func (this *pBigIntBuilder) BuildFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	setPBigIntFromDecimalFloat(value, dst)
}

func (this *pBigIntBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	setPBigIntFromBigDecimalFloat(value, dst)
}

func (this *pBigIntBuilder) BuildFromUUID(value []byte, dst reflect.Value) {
	builderPanicBadEvent(this, typePBigInt, "UUID")
}

func (this *pBigIntBuilder) BuildFromString(value string, dst reflect.Value) {
	builderPanicBadEvent(this, typePBigInt, "String")
}

func (this *pBigIntBuilder) BuildFromBytes(value []byte, dst reflect.Value) {
	builderPanicBadEvent(this, typePBigInt, "Bytes")
}

func (this *pBigIntBuilder) BuildFromURI(value *url.URL, dst reflect.Value) {
	builderPanicBadEvent(this, typePBigInt, "URI")
}

func (this *pBigIntBuilder) BuildFromTime(value time.Time, dst reflect.Value) {
	builderPanicBadEvent(this, typePBigInt, "Time")
}

func (this *pBigIntBuilder) BuildBeginList() {
	builderPanicBadEvent(this, typePBigInt, "List")
}

func (this *pBigIntBuilder) BuildBeginMap() {
	builderPanicBadEvent(this, typePBigInt, "Map")
}

func (this *pBigIntBuilder) BuildEndContainer() {
	builderPanicBadEvent(this, typePBigInt, "ContainerEnd")
}

func (this *pBigIntBuilder) BuildFromMarker(id interface{}) {
	panic("TODO: pBigIntBuilder.Marker")
}

func (this *pBigIntBuilder) BuildFromReference(id interface{}) {
	panic("TODO: pBigIntBuilder.Reference")
}

func (this *pBigIntBuilder) PrepareForListContents() {
}

func (this *pBigIntBuilder) PrepareForMapContents() {
}

func (this *pBigIntBuilder) NotifyChildContainerFinished(value reflect.Value) {
	builderPanicBadEvent(this, typePBigInt, "NotifyChildContainerFinished")
}
