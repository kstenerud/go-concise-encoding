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

type floatBuilder struct {
	// Template Data
	session *Session
	dstType reflect.Type
}

func newFloatBuilder(dstType reflect.Type) ObjectBuilder {
	return &floatBuilder{
		dstType: dstType,
	}
}

func (_this *floatBuilder) String() string {
	return fmt.Sprintf("%v", reflect.TypeOf(_this))
}

func (_this *floatBuilder) InitTemplate(session *Session) {
	_this.session = session
}

func (_this *floatBuilder) NewInstance(_ *RootBuilder, _ ObjectBuilder, _ *options.BuilderOptions) ObjectBuilder {
	return _this
}

func (_this *floatBuilder) SetParent(_ ObjectBuilder) {
}

func (_this *floatBuilder) BuildFromNil(_ reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, _this.dstType, "Nil")
}

func (_this *floatBuilder) BuildFromBool(_ bool, _ reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, _this.dstType, "Bool")
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

func (_this *floatBuilder) BuildFromUUID(_ []byte, _ reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, _this.dstType, "UUID")
}

func (_this *floatBuilder) BuildFromString(_ []byte, _ reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, _this.dstType, "String")
}

func (_this *floatBuilder) BuildFromVerbatimString(_ []byte, _ reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, _this.dstType, "VerbatimString")
}

func (_this *floatBuilder) BuildFromBytes(_ []byte, _ reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, _this.dstType, "Bytes")
}

func (_this *floatBuilder) BuildFromCustomBinary(value []byte, dst reflect.Value) {
	if err := _this.session.GetCustomBinaryBuildFunction()(value, dst); err != nil {
		BuilderPanicBuildFromCustomBinary(_this, value, dst.Type(), err)
	}
}

func (_this *floatBuilder) BuildFromCustomText(value []byte, dst reflect.Value) {
	if err := _this.session.GetCustomTextBuildFunction()(value, dst); err != nil {
		BuilderPanicBuildFromCustomText(_this, value, dst.Type(), err)
	}
}

func (_this *floatBuilder) BuildFromURI(_ *url.URL, _ reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, _this.dstType, "URI")
}

func (_this *floatBuilder) BuildFromTime(_ time.Time, _ reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, _this.dstType, "Time")
}

func (_this *floatBuilder) BuildFromCompactTime(_ *compact_time.Time, _ reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, _this.dstType, "CompactTime")
}

func (_this *floatBuilder) BuildBeginList() {
	BuilderWithTypePanicBadEvent(_this, _this.dstType, "ListBegin")
}

func (_this *floatBuilder) BuildBeginMap() {
	BuilderWithTypePanicBadEvent(_this, _this.dstType, "MapBegin")
}

func (_this *floatBuilder) BuildEndContainer() {
	BuilderWithTypePanicBadEvent(_this, _this.dstType, "ContainerEnd")
}

func (_this *floatBuilder) BuildBeginMarker(_ interface{}) {
	BuilderWithTypePanicBadEvent(_this, _this.dstType, "Marker")
}

func (_this *floatBuilder) BuildFromReference(_ interface{}) {
	BuilderWithTypePanicBadEvent(_this, _this.dstType, "Reference")
}

func (_this *floatBuilder) PrepareForListContents() {
	BuilderWithTypePanicBadEvent(_this, _this.dstType, "PrepareForListContents")
}

func (_this *floatBuilder) PrepareForMapContents() {
	BuilderWithTypePanicBadEvent(_this, _this.dstType, "PrepareForMapContents")
}

func (_this *floatBuilder) NotifyChildContainerFinished(_ reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, _this.dstType, "NotifyChildContainerFinished")
}
