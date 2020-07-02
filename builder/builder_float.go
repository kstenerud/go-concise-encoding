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

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
)

type floatBuilder struct {
	// Const data
	dstType reflect.Type
}

func newFloatBuilder(dstType reflect.Type) ObjectBuilder {
	return &floatBuilder{
		dstType: dstType,
	}
}

func (_this *floatBuilder) IsContainerOnly() bool {
	return false
}

func (_this *floatBuilder) PostCacheInitBuilder() {
}

func (_this *floatBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder, options *BuilderOptions) ObjectBuilder {
	return _this
}

func (_this *floatBuilder) SetParent(parent ObjectBuilder) {
}

func (_this *floatBuilder) BuildFromNil(dst reflect.Value) {
	builderPanicBadEventType(_this, _this.dstType, "Nil")
}

func (_this *floatBuilder) BuildFromBool(value bool, dst reflect.Value) {
	builderPanicBadEventType(_this, _this.dstType, "Bool")
}

func (_this *floatBuilder) BuildFromInt(value int64, dst reflect.Value) {
	setFloatFromInt(value, dst)
}

func (_this *floatBuilder) BuildFromUint(value uint64, dst reflect.Value) {
	setFloatFromUint(value, dst)
}

func (_this *floatBuilder) BuildFromBigInt(value *big.Int, dst reflect.Value) {
	setFloatFromBigInt(value, dst)
}

func (_this *floatBuilder) BuildFromFloat(value float64, dst reflect.Value) {
	setFloatFromFloat(value, dst)
}

func (_this *floatBuilder) BuildFromBigFloat(value *big.Float, dst reflect.Value) {
	setFloatFromBigFloat(value, dst)
}

func (_this *floatBuilder) BuildFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	setFloatFromDecimalFloat(value, dst)
}

func (_this *floatBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	setFloatFromBigDecimalFloat(value, dst)
}

func (_this *floatBuilder) BuildFromUUID(value []byte, dst reflect.Value) {
	builderPanicBadEventType(_this, _this.dstType, "UUID")
}

func (_this *floatBuilder) BuildFromString(value string, dst reflect.Value) {
	builderPanicBadEventType(_this, _this.dstType, "String")
}

func (_this *floatBuilder) BuildFromBytes(value []byte, dst reflect.Value) {
	builderPanicBadEventType(_this, _this.dstType, "Bytes")
}

func (_this *floatBuilder) BuildFromURI(value *url.URL, dst reflect.Value) {
	builderPanicBadEventType(_this, _this.dstType, "URI")
}

func (_this *floatBuilder) BuildFromTime(value time.Time, dst reflect.Value) {
	builderPanicBadEventType(_this, _this.dstType, "Time")
}

func (_this *floatBuilder) BuildFromCompactTime(value *compact_time.Time, dst reflect.Value) {
	builderPanicBadEventType(_this, _this.dstType, "CompactTime")
}

func (_this *floatBuilder) BuildBeginList() {
	builderPanicBadEventType(_this, _this.dstType, "ListBegin")
}

func (_this *floatBuilder) BuildBeginMap() {
	builderPanicBadEventType(_this, _this.dstType, "MapBegin")
}

func (_this *floatBuilder) BuildEndContainer() {
	builderPanicBadEventType(_this, _this.dstType, "ContainerEnd")
}

func (_this *floatBuilder) BuildBeginMarker(id interface{}) {
	panic("TODO: floatBuilder.BuildBeginMarker")
}

func (_this *floatBuilder) BuildFromReference(id interface{}) {
	panic("TODO: floatBuilder.BuildFromReference")
}

func (_this *floatBuilder) PrepareForListContents() {
	builderPanicBadEventType(_this, _this.dstType, "PrepareForListContents")
}

func (_this *floatBuilder) PrepareForMapContents() {
	builderPanicBadEventType(_this, _this.dstType, "PrepareForMapContents")
}

func (_this *floatBuilder) NotifyChildContainerFinished(value reflect.Value) {
	builderPanicBadEventType(_this, _this.dstType, "NotifyChildContainerFinished")
}
