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
	"reflect"
	"time"

	"github.com/kstenerud/go-concise-encoding/events"
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

func (_this *customBuilder) panicBadEvent(name string, args ...interface{}) {
	PanicBadEvent(_this, name, args...)
}

func (_this *customBuilder) InitTemplate(_ *Session) {
	_this.panicBadEvent("InitTemplate")
}

func (_this *customBuilder) NewInstance(_ *RootBuilder, _ ObjectBuilder, _ *options.BuilderOptions) ObjectBuilder {
	return _this
}

func (_this *customBuilder) SetParent(_ ObjectBuilder) {
}

func (_this *customBuilder) BuildFromNil(_ reflect.Value) {
	_this.panicBadEvent("Nil")
}

func (_this *customBuilder) BuildFromBool(_ bool, _ reflect.Value) {
	_this.panicBadEvent("Bool")
}

func (_this *customBuilder) BuildFromInt(_ int64, _ reflect.Value) {
	_this.panicBadEvent("Int")
}

func (_this *customBuilder) BuildFromUint(_ uint64, _ reflect.Value) {
	_this.panicBadEvent("Uint")
}

func (_this *customBuilder) BuildFromBigInt(_ *big.Int, _ reflect.Value) {
	_this.panicBadEvent("BigInt")
}

func (_this *customBuilder) BuildFromFloat(_ float64, _ reflect.Value) {
	_this.panicBadEvent("Float")
}

func (_this *customBuilder) BuildFromBigFloat(_ *big.Float, _ reflect.Value) {
	_this.panicBadEvent("BigFloat")
}

func (_this *customBuilder) BuildFromDecimalFloat(_ compact_float.DFloat, _ reflect.Value) {
	_this.panicBadEvent("DecimalFloat")
}

func (_this *customBuilder) BuildFromBigDecimalFloat(_ *apd.Decimal, _ reflect.Value) {
	_this.panicBadEvent("BigDecimalFloat")
}

func (_this *customBuilder) BuildFromUUID(_ []byte, _ reflect.Value) {
	_this.panicBadEvent("UUID")
}

func (_this *customBuilder) BuildFromArray(arrayType events.ArrayType, value []byte, dst reflect.Value) {
	if !_this.session.TryBuildFromCustom(_this, arrayType, value, dst) {
		_this.panicBadEvent("TypedArray(%v)", arrayType)
	}
}

func (_this *customBuilder) BuildFromTime(_ time.Time, _ reflect.Value) {
	_this.panicBadEvent("Time")
}

func (_this *customBuilder) BuildFromCompactTime(_ *compact_time.Time, _ reflect.Value) {
	_this.panicBadEvent("CompactTime")
}

func (_this *customBuilder) BuildBeginList() {
	_this.panicBadEvent("List")
}

func (_this *customBuilder) BuildBeginMap() {
	_this.panicBadEvent("Map")
}

func (_this *customBuilder) BuildEndContainer() {
	_this.panicBadEvent("ContainerEnd")
}

func (_this *customBuilder) BuildBeginMarker(_ interface{}) {
	_this.panicBadEvent("Marker")
}

func (_this *customBuilder) BuildFromReference(_ interface{}) {
	_this.panicBadEvent("Reference")
}

func (_this *customBuilder) PrepareForListContents() {
	_this.panicBadEvent("PrepareForListContents")
}

func (_this *customBuilder) PrepareForMapContents() {
	_this.panicBadEvent("PrepareForMapContents")
}

func (_this *customBuilder) NotifyChildContainerFinished(_ reflect.Value) {
	_this.panicBadEvent("NotifyChildContainerFinished")
}
