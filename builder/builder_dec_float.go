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

	"github.com/kstenerud/go-concise-encoding/internal/common"
	"github.com/kstenerud/go-concise-encoding/options"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
)

type dfloatBuilder struct {
	// Template Data
	session *Session
}

func newDFloatBuilder() ObjectBuilder {
	return &dfloatBuilder{}
}

func (_this *dfloatBuilder) String() string {
	return fmt.Sprintf("%v", reflect.TypeOf(_this))
}

func (_this *dfloatBuilder) InitTemplate(session *Session) {
	_this.session = session
}

func (_this *dfloatBuilder) NewInstance(_ *RootBuilder, _ ObjectBuilder, _ *options.BuilderOptions) ObjectBuilder {
	return _this
}

func (_this *dfloatBuilder) SetParent(_ ObjectBuilder) {
}

func (_this *dfloatBuilder) BuildFromNil(_ reflect.Value) {
	PanicBadEventWithType(_this, common.TypeDFloat, "Nil")
}

func (_this *dfloatBuilder) BuildFromBool(_ bool, _ reflect.Value) {
	PanicBadEventWithType(_this, common.TypeDFloat, "Bool")
}

func (_this *dfloatBuilder) BuildFromInt(value int64, dst reflect.Value) {
	dst.Set(reflect.ValueOf(compact_float.DFloatValue(0, value)))
}

func (_this *dfloatBuilder) BuildFromUint(value uint64, dst reflect.Value) {
	dst.Set(reflect.ValueOf(compact_float.DFloatFromUInt(value)))
}

func (_this *dfloatBuilder) BuildFromBigInt(value *big.Int, dst reflect.Value) {
	dst.Set(reflect.ValueOf(compact_float.DFloatFromBigInt(value)))
}

func (_this *dfloatBuilder) BuildFromFloat(value float64, dst reflect.Value) {
	dst.Set(reflect.ValueOf(compact_float.DFloatFromFloat64(value, 0)))
}

func (_this *dfloatBuilder) BuildFromBigFloat(value *big.Float, dst reflect.Value) {
	dst.Set(reflect.ValueOf(compact_float.DFloatFromBigFloat(value)))
}

func (_this *dfloatBuilder) BuildFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (_this *dfloatBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	dst.Set(reflect.ValueOf(compact_float.DFloatFromAPD(value)))
}

func (_this *dfloatBuilder) BuildFromUUID(_ []byte, _ reflect.Value) {
	PanicBadEventWithType(_this, common.TypeDFloat, "UUID")
}

func (_this *dfloatBuilder) BuildFromString(_ []byte, _ reflect.Value) {
	PanicBadEventWithType(_this, common.TypeDFloat, "String")
}

func (_this *dfloatBuilder) BuildFromVerbatimString(_ []byte, _ reflect.Value) {
	PanicBadEventWithType(_this, common.TypeDFloat, "VerbatimString")
}

func (_this *dfloatBuilder) BuildFromURI(_ *url.URL, _ reflect.Value) {
	PanicBadEventWithType(_this, common.TypeDFloat, "URI")
}

func (_this *dfloatBuilder) BuildFromCustomBinary(value []byte, dst reflect.Value) {
	if err := _this.session.GetCustomBinaryBuildFunction()(value, dst); err != nil {
		PanicBuildFromCustomBinary(_this, value, dst.Type(), err)
	}
}

func (_this *dfloatBuilder) BuildFromCustomText(value []byte, dst reflect.Value) {
	if err := _this.session.GetCustomTextBuildFunction()(value, dst); err != nil {
		PanicBuildFromCustomText(_this, value, dst.Type(), err)
	}
}

func (_this *dfloatBuilder) BuildFromTypedArray(elemType reflect.Type, _ []byte, _ reflect.Value) {
	PanicBadEventWithType(_this, common.TypeDFloat, "TypedArray(%v)", elemType)
}

func (_this *dfloatBuilder) BuildFromTime(_ time.Time, _ reflect.Value) {
	PanicBadEventWithType(_this, common.TypeDFloat, "Time")
}

func (_this *dfloatBuilder) BuildFromCompactTime(_ *compact_time.Time, _ reflect.Value) {
	PanicBadEventWithType(_this, common.TypeDFloat, "CompactTime")
}

func (_this *dfloatBuilder) BuildBeginList() {
	PanicBadEventWithType(_this, common.TypeDFloat, "List")
}

func (_this *dfloatBuilder) BuildBeginMap() {
	PanicBadEventWithType(_this, common.TypeDFloat, "Map")
}

func (_this *dfloatBuilder) BuildEndContainer() {
	PanicBadEventWithType(_this, common.TypeDFloat, "ContainerEnd")
}

func (_this *dfloatBuilder) BuildBeginMarker(_ interface{}) {
	PanicBadEventWithType(_this, common.TypeDFloat, "Marker")
}

func (_this *dfloatBuilder) BuildFromReference(_ interface{}) {
	PanicBadEventWithType(_this, common.TypeDFloat, "Reference")
}

func (_this *dfloatBuilder) PrepareForListContents() {
	PanicBadEventWithType(_this, common.TypeDFloat, "PrepareForListContents")
}

func (_this *dfloatBuilder) PrepareForMapContents() {
	PanicBadEventWithType(_this, common.TypeDFloat, "PrepareForMapContents")
}

func (_this *dfloatBuilder) NotifyChildContainerFinished(_ reflect.Value) {
	PanicBadEventWithType(_this, common.TypeDFloat, "NotifyChildContainerFinished")
}
