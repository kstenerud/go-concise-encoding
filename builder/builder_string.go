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

	"github.com/kstenerud/go-concise-encoding/events"
	"github.com/kstenerud/go-concise-encoding/internal/common"
	"github.com/kstenerud/go-concise-encoding/options"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
)

type stringBuilder struct {
	// Template Data
	session *Session
}

func newStringBuilder() ObjectBuilder {
	return &stringBuilder{}
}

func (_this *stringBuilder) String() string {
	return fmt.Sprintf("%v", reflect.TypeOf(_this))
}

func (_this *stringBuilder) InitTemplate(session *Session) {
	_this.session = session
}

func (_this *stringBuilder) NewInstance(_ *RootBuilder, _ ObjectBuilder, _ *options.BuilderOptions) ObjectBuilder {
	return _this
}

func (_this *stringBuilder) SetParent(_ ObjectBuilder) {
}

func (_this *stringBuilder) BuildFromNil(dst reflect.Value) {
	// Go doesn't have the concept of a nil string.
	dst.SetString("")
}

func (_this *stringBuilder) BuildFromBool(_ bool, _ reflect.Value) {
	PanicBadEventWithType(_this, common.TypeString, "Bool")
}

func (_this *stringBuilder) BuildFromInt(_ int64, _ reflect.Value) {
	PanicBadEventWithType(_this, common.TypeString, "Int")
}

func (_this *stringBuilder) BuildFromUint(_ uint64, _ reflect.Value) {
	PanicBadEventWithType(_this, common.TypeString, "Uint")
}

func (_this *stringBuilder) BuildFromBigInt(_ *big.Int, _ reflect.Value) {
	PanicBadEventWithType(_this, common.TypeString, "BigInt")
}

func (_this *stringBuilder) BuildFromFloat(_ float64, _ reflect.Value) {
	PanicBadEventWithType(_this, common.TypeString, "Float")
}

func (_this *stringBuilder) BuildFromBigFloat(_ *big.Float, _ reflect.Value) {
	PanicBadEventWithType(_this, common.TypeString, "BigFloat")
}

func (_this *stringBuilder) BuildFromDecimalFloat(_ compact_float.DFloat, _ reflect.Value) {
	PanicBadEventWithType(_this, common.TypeString, "DecimalFloat")
}

func (_this *stringBuilder) BuildFromBigDecimalFloat(_ *apd.Decimal, _ reflect.Value) {
	PanicBadEventWithType(_this, common.TypeString, "BigDecimalFloat")
}

func (_this *stringBuilder) BuildFromUUID(_ []byte, _ reflect.Value) {
	PanicBadEventWithType(_this, common.TypeString, "UUID")
}

func (_this *stringBuilder) BuildFromString(value []byte, dst reflect.Value) {
	dst.SetString(string(value))
}

func (_this *stringBuilder) BuildFromVerbatimString(value []byte, dst reflect.Value) {
	dst.SetString(string(value))
}

func (_this *stringBuilder) BuildFromURI(_ *url.URL, _ reflect.Value) {
	PanicBadEventWithType(_this, common.TypeString, "URI")
}

func (_this *stringBuilder) BuildFromCustomBinary(value []byte, dst reflect.Value) {
	if err := _this.session.GetCustomBinaryBuildFunction()(value, dst); err != nil {
		PanicBuildFromCustomBinary(_this, value, dst.Type(), err)
	}
}

func (_this *stringBuilder) BuildFromCustomText(value []byte, dst reflect.Value) {
	if err := _this.session.GetCustomTextBuildFunction()(value, dst); err != nil {
		PanicBuildFromCustomText(_this, value, dst.Type(), err)
	}
}

func (_this *stringBuilder) BuildFromTypedArray(arrayType events.ArrayType, _ []byte, _ reflect.Value) {
	PanicBadEventWithType(_this, common.TypeString, "TypedArray(%v)", arrayType)
}

func (_this *stringBuilder) BuildFromTime(_ time.Time, _ reflect.Value) {
	PanicBadEventWithType(_this, common.TypeString, "Time")
}

func (_this *stringBuilder) BuildFromCompactTime(_ *compact_time.Time, _ reflect.Value) {
	PanicBadEventWithType(_this, common.TypeString, "CompactTime")
}

func (_this *stringBuilder) BuildBeginList() {
	PanicBadEventWithType(_this, common.TypeString, "List")
}

func (_this *stringBuilder) BuildBeginMap() {
	PanicBadEventWithType(_this, common.TypeString, "Map")
}

func (_this *stringBuilder) BuildEndContainer() {
	PanicBadEventWithType(_this, common.TypeString, "ContainerEnd")
}

func (_this *stringBuilder) BuildBeginMarker(_ interface{}) {
	PanicBadEventWithType(_this, common.TypeString, "Marker")
}

func (_this *stringBuilder) BuildFromReference(_ interface{}) {
	PanicBadEventWithType(_this, common.TypeString, "Reference")
}

func (_this *stringBuilder) PrepareForListContents() {
	PanicBadEventWithType(_this, common.TypeString, "PrepareForListContents")
}

func (_this *stringBuilder) PrepareForMapContents() {
	PanicBadEventWithType(_this, common.TypeString, "PrepareForMapContents")
}

func (_this *stringBuilder) NotifyChildContainerFinished(_ reflect.Value) {
	PanicBadEventWithType(_this, common.TypeString, "NotifyChildContainerFinished")
}
