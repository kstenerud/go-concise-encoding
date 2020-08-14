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
	// Template Data
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

func (_this *customBuilder) InitTemplate(_ *Session) {
	PanicBadEvent(_this, "InitTemplate")
}

func (_this *customBuilder) NewInstance(_ *RootBuilder, _ ObjectBuilder, _ *options.BuilderOptions) ObjectBuilder {
	return _this
}

func (_this *customBuilder) SetParent(_ ObjectBuilder) {
}

func (_this *customBuilder) BuildFromNil(_ reflect.Value) {
	PanicBadEvent(_this, "Nil")
}

func (_this *customBuilder) BuildFromBool(_ bool, _ reflect.Value) {
	PanicBadEvent(_this, "Bool")
}

func (_this *customBuilder) BuildFromInt(_ int64, _ reflect.Value) {
	PanicBadEvent(_this, "Int")
}

func (_this *customBuilder) BuildFromUint(_ uint64, _ reflect.Value) {
	PanicBadEvent(_this, "Uint")
}

func (_this *customBuilder) BuildFromBigInt(_ *big.Int, _ reflect.Value) {
	PanicBadEvent(_this, "BigInt")
}

func (_this *customBuilder) BuildFromFloat(_ float64, _ reflect.Value) {
	PanicBadEvent(_this, "Float")
}

func (_this *customBuilder) BuildFromBigFloat(_ *big.Float, _ reflect.Value) {
	PanicBadEvent(_this, "BigFloat")
}

func (_this *customBuilder) BuildFromDecimalFloat(_ compact_float.DFloat, _ reflect.Value) {
	PanicBadEvent(_this, "DecimalFloat")
}

func (_this *customBuilder) BuildFromBigDecimalFloat(_ *apd.Decimal, _ reflect.Value) {
	PanicBadEvent(_this, "BigDecimalFloat")
}

func (_this *customBuilder) BuildFromUUID(_ []byte, _ reflect.Value) {
	PanicBadEvent(_this, "UUID")
}

func (_this *customBuilder) BuildFromString(_ []byte, _ reflect.Value) {
	PanicBadEvent(_this, "String")
}

func (_this *customBuilder) BuildFromVerbatimString(_ []byte, _ reflect.Value) {
	PanicBadEvent(_this, "VerbatimString")
}

func (_this *customBuilder) BuildFromBytes(_ []byte, _ reflect.Value) {
	PanicBadEvent(_this, "Bytes")
}

func (_this *customBuilder) BuildFromCustomBinary(value []byte, dst reflect.Value) {
	if err := _this.session.GetCustomBinaryBuildFunction()(value, dst); err != nil {
		PanicBuildFromCustomBinary(_this, value, dst.Type(), err)
	}
}

func (_this *customBuilder) BuildFromCustomText(value []byte, dst reflect.Value) {
	if err := _this.session.GetCustomTextBuildFunction()(value, dst); err != nil {
		PanicBuildFromCustomText(_this, value, dst.Type(), err)
	}
}

func (_this *customBuilder) BuildFromURI(_ *url.URL, _ reflect.Value) {
	PanicBadEvent(_this, "URI")
}

func (_this *customBuilder) BuildFromTime(_ time.Time, _ reflect.Value) {
	PanicBadEvent(_this, "Time")
}

func (_this *customBuilder) BuildFromCompactTime(_ *compact_time.Time, _ reflect.Value) {
	PanicBadEvent(_this, "CompactTime")
}

func (_this *customBuilder) BuildBeginList() {
	PanicBadEvent(_this, "List")
}

func (_this *customBuilder) BuildBeginMap() {
	PanicBadEvent(_this, "Map")
}

func (_this *customBuilder) BuildEndContainer() {
	PanicBadEvent(_this, "ContainerEnd")
}

func (_this *customBuilder) BuildBeginMarker(_ interface{}) {
	PanicBadEvent(_this, "Marker")
}

func (_this *customBuilder) BuildFromReference(_ interface{}) {
	PanicBadEvent(_this, "Reference")
}

func (_this *customBuilder) PrepareForListContents() {
	PanicBadEvent(_this, "PrepareForListContents")
}

func (_this *customBuilder) PrepareForMapContents() {
	PanicBadEvent(_this, "PrepareForMapContents")
}

func (_this *customBuilder) NotifyChildContainerFinished(_ reflect.Value) {
	PanicBadEvent(_this, "NotifyChildContainerFinished")
}
