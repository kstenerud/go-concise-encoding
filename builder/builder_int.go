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

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
)

type intBuilder struct {
	// Const data
	dstType reflect.Type
}

func newIntBuilder(dstType reflect.Type) ObjectBuilder {
	return &intBuilder{
		dstType: dstType,
	}
}

func (_this *intBuilder) String() string {
	return fmt.Sprintf("%v<%v>", reflect.TypeOf(_this), _this.dstType)
}

func (_this *intBuilder) IsContainerOnly() bool {
	return false
}

func (_this *intBuilder) PostCacheInitBuilder() {
}

func (_this *intBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder, options *BuilderOptions) ObjectBuilder {
	return _this
}

func (_this *intBuilder) SetParent(parent ObjectBuilder) {
}

func (_this *intBuilder) BuildFromNil(dst reflect.Value) {
	builderPanicBadEventType(_this, _this.dstType, "Nil")
}

func (_this *intBuilder) BuildFromBool(value bool, dst reflect.Value) {
	builderPanicBadEventType(_this, _this.dstType, "Bool")
}

func (_this *intBuilder) BuildFromInt(value int64, dst reflect.Value) {
	setIntFromInt(value, dst)
}

func (_this *intBuilder) BuildFromUint(value uint64, dst reflect.Value) {
	setIntFromUint(value, dst)
}

func (_this *intBuilder) BuildFromBigInt(value *big.Int, dst reflect.Value) {
	setIntFromBigInt(value, dst)
}

func (_this *intBuilder) BuildFromFloat(value float64, dst reflect.Value) {
	setIntFromFloat(value, dst)
}

func (_this *intBuilder) BuildFromBigFloat(value *big.Float, dst reflect.Value) {
	setIntFromBigFloat(value, dst)
}

func (_this *intBuilder) BuildFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	setIntFromDecimalFloat(value, dst)
}

func (_this *intBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	setIntFromBigDecimalFloat(value, dst)
}

func (_this *intBuilder) BuildFromUUID(value []byte, dst reflect.Value) {
	builderPanicBadEventType(_this, _this.dstType, "UUID")
}

func (_this *intBuilder) BuildFromString(value string, dst reflect.Value) {
	builderPanicBadEventType(_this, _this.dstType, "String")
}

func (_this *intBuilder) BuildFromBytes(value []byte, dst reflect.Value) {
	builderPanicBadEventType(_this, _this.dstType, "Bytes")
}

func (_this *intBuilder) BuildFromURI(value *url.URL, dst reflect.Value) {
	builderPanicBadEventType(_this, _this.dstType, "URI")
}

func (_this *intBuilder) BuildFromTime(value time.Time, dst reflect.Value) {
	builderPanicBadEventType(_this, _this.dstType, "Time")
}

func (_this *intBuilder) BuildFromCompactTime(value *compact_time.Time, dst reflect.Value) {
	builderPanicBadEventType(_this, _this.dstType, "CompactTime")
}

func (_this *intBuilder) BuildBeginList() {
	builderPanicBadEventType(_this, _this.dstType, "List")
}

func (_this *intBuilder) BuildBeginMap() {
	builderPanicBadEventType(_this, _this.dstType, "Map")
}

func (_this *intBuilder) BuildEndContainer() {
	builderPanicBadEventType(_this, _this.dstType, "ContainerEnd")
}

func (_this *intBuilder) BuildBeginMarker(id interface{}) {
	panic("TODO: intBuilder.Marker")
}

func (_this *intBuilder) BuildFromReference(id interface{}) {
	panic("TODO: intBuilder.Reference")
}

func (_this *intBuilder) PrepareForListContents() {
	builderPanicBadEventType(_this, _this.dstType, "PrepareForListContents")
}

func (_this *intBuilder) PrepareForMapContents() {
	builderPanicBadEventType(_this, _this.dstType, "PrepareForMapContents")
}

func (_this *intBuilder) NotifyChildContainerFinished(value reflect.Value) {
	builderPanicBadEventType(_this, _this.dstType, "NotifyChildContainerFinished")
}
