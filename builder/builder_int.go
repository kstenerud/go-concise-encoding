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

type intBuilder struct {
	// Template Data
	session *Session
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

func (_this *intBuilder) InitTemplate(session *Session) {
	_this.session = session
}

func (_this *intBuilder) NewInstance(_ *RootBuilder, _ ObjectBuilder, _ *options.BuilderOptions) ObjectBuilder {
	return _this
}

func (_this *intBuilder) SetParent(_ ObjectBuilder) {
}

func (_this *intBuilder) BuildFromNil(_ reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, _this.dstType, "Nil")
}

func (_this *intBuilder) BuildFromBool(_ bool, _ reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, _this.dstType, "Bool")
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

func (_this *intBuilder) BuildFromUUID(_ []byte, _ reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, _this.dstType, "UUID")
}

func (_this *intBuilder) BuildFromString(_ []byte, _ reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, _this.dstType, "String")
}

func (_this *intBuilder) BuildFromVerbatimString(_ []byte, _ reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, _this.dstType, "VerbatimString")
}

func (_this *intBuilder) BuildFromBytes(_ []byte, _ reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, _this.dstType, "Bytes")
}

func (_this *intBuilder) BuildFromCustomBinary(value []byte, dst reflect.Value) {
	if err := _this.session.GetCustomBinaryBuildFunction()(value, dst); err != nil {
		BuilderPanicBuildFromCustomBinary(_this, value, dst.Type(), err)
	}
}

func (_this *intBuilder) BuildFromCustomText(value []byte, dst reflect.Value) {
	if err := _this.session.GetCustomTextBuildFunction()(value, dst); err != nil {
		BuilderPanicBuildFromCustomText(_this, value, dst.Type(), err)
	}
}

func (_this *intBuilder) BuildFromURI(_ *url.URL, _ reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, _this.dstType, "URI")
}

func (_this *intBuilder) BuildFromTime(_ time.Time, _ reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, _this.dstType, "Time")
}

func (_this *intBuilder) BuildFromCompactTime(_ *compact_time.Time, _ reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, _this.dstType, "CompactTime")
}

func (_this *intBuilder) BuildBeginList() {
	BuilderWithTypePanicBadEvent(_this, _this.dstType, "List")
}

func (_this *intBuilder) BuildBeginMap() {
	BuilderWithTypePanicBadEvent(_this, _this.dstType, "Map")
}

func (_this *intBuilder) BuildEndContainer() {
	BuilderWithTypePanicBadEvent(_this, _this.dstType, "ContainerEnd")
}

func (_this *intBuilder) BuildBeginMarker(_ interface{}) {
	BuilderWithTypePanicBadEvent(_this, _this.dstType, "Marker")
}

func (_this *intBuilder) BuildFromReference(_ interface{}) {
	BuilderWithTypePanicBadEvent(_this, _this.dstType, "Reference")
}

func (_this *intBuilder) PrepareForListContents() {
	BuilderWithTypePanicBadEvent(_this, _this.dstType, "PrepareForListContents")
}

func (_this *intBuilder) PrepareForMapContents() {
	BuilderWithTypePanicBadEvent(_this, _this.dstType, "PrepareForMapContents")
}

func (_this *intBuilder) NotifyChildContainerFinished(_ reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, _this.dstType, "NotifyChildContainerFinished")
}
