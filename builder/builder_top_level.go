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

// topLevelContainerBuilder proxies the first build instruction to make sure containers
// are properly built. See BuildBeginList and BuildBeginMap.
type topLevelContainerBuilder struct {
	builder ObjectBuilder
	root    *RootBuilder
}

func newTopLevelContainerBuilder(root *RootBuilder, builder ObjectBuilder) ObjectBuilder {
	return &topLevelContainerBuilder{
		builder: builder,
		root:    root,
	}
}

func (_this *topLevelContainerBuilder) IsContainerOnly() bool {
	builderPanicBadEvent(_this, "IsContainerOnly")
	return false
}

func (_this *topLevelContainerBuilder) PostCacheInitBuilder() {
	builderPanicBadEvent(_this, "PostCacheInitBuilder")
}

func (_this *topLevelContainerBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder, options *BuilderOptions) ObjectBuilder {
	builderPanicBadEvent(_this, "CloneFromTemplate")
	return nil
}

func (_this *topLevelContainerBuilder) SetParent(parent ObjectBuilder) {
	builderPanicBadEvent(_this, "SetParent")
}

func (_this *topLevelContainerBuilder) BuildFromNil(dst reflect.Value) {
	builderPanicBadEvent(_this, "Nil")
}

func (_this *topLevelContainerBuilder) BuildFromBool(value bool, dst reflect.Value) {
	builderPanicBadEvent(_this, "Bool")
}

func (_this *topLevelContainerBuilder) BuildFromInt(value int64, dst reflect.Value) {
	builderPanicBadEvent(_this, "Int")
}

func (_this *topLevelContainerBuilder) BuildFromUint(value uint64, dst reflect.Value) {
	builderPanicBadEvent(_this, "Uint")
}

func (_this *topLevelContainerBuilder) BuildFromBigInt(value *big.Int, dst reflect.Value) {
	builderPanicBadEvent(_this, "BigInt")
}

func (_this *topLevelContainerBuilder) BuildFromFloat(value float64, dst reflect.Value) {
	builderPanicBadEvent(_this, "Float")
}

func (_this *topLevelContainerBuilder) BuildFromBigFloat(value *big.Float, dst reflect.Value) {
	builderPanicBadEvent(_this, "BigFloat")
}

func (_this *topLevelContainerBuilder) BuildFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	builderPanicBadEvent(_this, "DecimalFloat")
}

func (_this *topLevelContainerBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	builderPanicBadEvent(_this, "BigDecimalFloat")
}

func (_this *topLevelContainerBuilder) BuildFromUUID(value []byte, dst reflect.Value) {
	builderPanicBadEvent(_this, "UUID")
}

func (_this *topLevelContainerBuilder) BuildFromString(value string, dst reflect.Value) {
	builderPanicBadEvent(_this, "String")
}

func (_this *topLevelContainerBuilder) BuildFromBytes(value []byte, dst reflect.Value) {
	builderPanicBadEvent(_this, "Bytes")
}

func (_this *topLevelContainerBuilder) BuildFromURI(value *url.URL, dst reflect.Value) {
	builderPanicBadEvent(_this, "URI")
}

func (_this *topLevelContainerBuilder) BuildFromTime(value time.Time, dst reflect.Value) {
	builderPanicBadEvent(_this, "Time")
}

func (_this *topLevelContainerBuilder) BuildFromCompactTime(value *compact_time.Time, dst reflect.Value) {
	builderPanicBadEvent(_this, "CompactTime")
}

func (_this *topLevelContainerBuilder) BuildBeginList() {
	_this.root.SetCurrentBuilder(_this.builder)
	_this.builder.PrepareForListContents()
}

func (_this *topLevelContainerBuilder) BuildBeginMap() {
	_this.root.SetCurrentBuilder(_this.builder)
	_this.builder.PrepareForMapContents()
}

func (_this *topLevelContainerBuilder) BuildEndContainer() {
	builderPanicBadEvent(_this, "End")
}

func (_this *topLevelContainerBuilder) PrepareForListContents() {
	builderPanicBadEvent(_this, "PrepareForListContents")
}

func (_this *topLevelContainerBuilder) PrepareForMapContents() {
	builderPanicBadEvent(_this, "PrepareForMapContents")
}

func (_this *topLevelContainerBuilder) NotifyChildContainerFinished(value reflect.Value) {
	builderPanicBadEvent(_this, "NotifyChildContainerFinished")
}

func (_this *topLevelContainerBuilder) BuildBeginMarker(id interface{}) {
	panic("TODO: BuildBeginMarker")
}

func (_this *topLevelContainerBuilder) BuildFromReference(id interface{}) {
	panic("TODO: BuildFromReference")
}
