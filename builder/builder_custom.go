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

	"github.com/kstenerud/go-concise-encoding/options"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
)

type customBuilder struct {
	// Static data
	session *Session
}

func newCustomBuilder(session *Session) ObjectBuilder {
	return &customBuilder{
		session: session,
	}
}

func (_this *customBuilder) String() string {
	return fmt.Sprintf("%v", reflect.TypeOf(_this))
}

func (_this *customBuilder) PostCacheInitBuilder(session *Session) {
	BuilderPanicBadEvent(_this, "PostCacheInitBuilder")
}

func (_this *customBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder, options *options.BuilderOptions) ObjectBuilder {
	return _this
}

func (_this *customBuilder) SetParent(parent ObjectBuilder) {
}

func (_this *customBuilder) BuildFromNil(dst reflect.Value) {
	BuilderPanicBadEvent(_this, "Nil")
}

func (_this *customBuilder) BuildFromBool(value bool, dst reflect.Value) {
	BuilderPanicBadEvent(_this, "Bool")
}

func (_this *customBuilder) BuildFromInt(value int64, dst reflect.Value) {
	BuilderPanicBadEvent(_this, "Int")
}

func (_this *customBuilder) BuildFromUint(value uint64, dst reflect.Value) {
	BuilderPanicBadEvent(_this, "Uint")
}

func (_this *customBuilder) BuildFromBigInt(value *big.Int, dst reflect.Value) {
	BuilderPanicBadEvent(_this, "BigInt")
}

func (_this *customBuilder) BuildFromFloat(value float64, dst reflect.Value) {
	BuilderPanicBadEvent(_this, "Float")
}

func (_this *customBuilder) BuildFromBigFloat(value *big.Float, dst reflect.Value) {
	BuilderPanicBadEvent(_this, "BigFloat")
}

func (_this *customBuilder) BuildFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	BuilderPanicBadEvent(_this, "DecimalFloat")
}

func (_this *customBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	BuilderPanicBadEvent(_this, "BigDecimalFloat")
}

func (_this *customBuilder) BuildFromUUID(value []byte, dst reflect.Value) {
	BuilderPanicBadEvent(_this, "UUID")
}

func (_this *customBuilder) BuildFromString(value []byte, dst reflect.Value) {
	BuilderPanicBadEvent(_this, "String")
}

func (_this *customBuilder) BuildFromVerbatimString(value []byte, dst reflect.Value) {
	BuilderPanicBadEvent(_this, "VerbatimString")
}

func (_this *customBuilder) BuildFromBytes(value []byte, dst reflect.Value) {
	BuilderPanicBadEvent(_this, "Bytes")
}

func (_this *customBuilder) BuildFromCustomBinary(value []byte, dst reflect.Value) {
	if err := _this.session.GetCustomBinaryBuildFunction()(value, dst); err != nil {
		BuilderPanicBuildFromCustomBinary(_this, value, dst.Type(), err)
	}
}

func (_this *customBuilder) BuildFromCustomText(value []byte, dst reflect.Value) {
	if err := _this.session.GetCustomTextBuildFunction()(value, dst); err != nil {
		BuilderPanicBuildFromCustomText(_this, value, dst.Type(), err)
	}
}

func (_this *customBuilder) BuildFromURI(value *url.URL, dst reflect.Value) {
	BuilderPanicBadEvent(_this, "URI")
}

func (_this *customBuilder) BuildFromTime(value time.Time, dst reflect.Value) {
	BuilderPanicBadEvent(_this, "Time")
}

func (_this *customBuilder) BuildFromCompactTime(value *compact_time.Time, dst reflect.Value) {
	BuilderPanicBadEvent(_this, "CompactTime")
}

func (_this *customBuilder) BuildBeginList() {
	BuilderPanicBadEvent(_this, "List")
}

func (_this *customBuilder) BuildBeginMap() {
	BuilderPanicBadEvent(_this, "Map")
}

func (_this *customBuilder) BuildEndContainer() {
	BuilderPanicBadEvent(_this, "ContainerEnd")
}

func (_this *customBuilder) BuildBeginMarker(id interface{}) {
	BuilderPanicBadEvent(_this, "Marker")
}

func (_this *customBuilder) BuildFromReference(id interface{}) {
	BuilderPanicBadEvent(_this, "Reference")
}

func (_this *customBuilder) PrepareForListContents() {
	BuilderPanicBadEvent(_this, "PrepareForListContents")
}

func (_this *customBuilder) PrepareForMapContents() {
	BuilderPanicBadEvent(_this, "PrepareForMapContents")
}

func (_this *customBuilder) NotifyChildContainerFinished(value reflect.Value) {
	BuilderPanicBadEvent(_this, "NotifyChildContainerFinished")
}
