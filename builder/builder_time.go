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
	"github.com/kstenerud/go-concise-encoding/internal/common"
	"github.com/kstenerud/go-concise-encoding/options"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
)

// Go Time

type timeBuilder struct {
	// Template Data
	session *Session
}

func newTimeBuilder() ObjectBuilder {
	return &timeBuilder{}
}

func (_this *timeBuilder) String() string {
	return fmt.Sprintf("%v", reflect.TypeOf(_this))
}

func (_this *timeBuilder) panicBadEvent(name string, args ...interface{}) {
	PanicBadEventWithType(_this, common.TypeTime, name, args...)
}

func (_this *timeBuilder) InitTemplate(session *Session) {
	_this.session = session
}

func (_this *timeBuilder) NewInstance(_ *RootBuilder, _ ObjectBuilder, _ *options.BuilderOptions) ObjectBuilder {
	return _this
}

func (_this *timeBuilder) SetParent(_ ObjectBuilder) {
}

func (_this *timeBuilder) BuildFromNil(_ reflect.Value) {
	_this.panicBadEvent("Nil")
}

func (_this *timeBuilder) BuildFromBool(_ bool, _ reflect.Value) {
	_this.panicBadEvent("Bool")
}

func (_this *timeBuilder) BuildFromInt(_ int64, _ reflect.Value) {
	_this.panicBadEvent("Int")
}

func (_this *timeBuilder) BuildFromUint(_ uint64, _ reflect.Value) {
	_this.panicBadEvent("Uint")
}

func (_this *timeBuilder) BuildFromBigInt(_ *big.Int, _ reflect.Value) {
	_this.panicBadEvent("BigInt")
}

func (_this *timeBuilder) BuildFromFloat(_ float64, _ reflect.Value) {
	_this.panicBadEvent("Float")
}

func (_this *timeBuilder) BuildFromBigFloat(_ *big.Float, _ reflect.Value) {
	_this.panicBadEvent("BigFloat")
}

func (_this *timeBuilder) BuildFromDecimalFloat(_ compact_float.DFloat, _ reflect.Value) {
	_this.panicBadEvent("DecimalFloat")
}

func (_this *timeBuilder) BuildFromBigDecimalFloat(_ *apd.Decimal, _ reflect.Value) {
	_this.panicBadEvent("BigDecimalFloat")
}

func (_this *timeBuilder) BuildFromUUID(_ []byte, _ reflect.Value) {
	_this.panicBadEvent("UUID")
}

func (_this *timeBuilder) BuildFromArray(arrayType events.ArrayType, value []byte, dst reflect.Value) {
	if !_this.session.TryBuildFromCustom(_this, arrayType, value, dst) {
		_this.panicBadEvent("TypedArray(%v)", arrayType)
	}
}

func (_this *timeBuilder) BuildFromTime(value time.Time, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (_this *timeBuilder) BuildFromCompactTime(value *compact_time.Time, dst reflect.Value) {
	v, err := value.AsGoTime()
	if err != nil {
		panic(err)
	}
	dst.Set(reflect.ValueOf(v))
}

func (_this *timeBuilder) BuildBeginList() {
	_this.panicBadEvent("List")
}

func (_this *timeBuilder) BuildBeginMap() {
	_this.panicBadEvent("Map")
}

func (_this *timeBuilder) BuildEndContainer() {
	_this.panicBadEvent("ContainerEnd")
}

func (_this *timeBuilder) BuildBeginMarker(_ interface{}) {
	_this.panicBadEvent("Marker")
}

func (_this *timeBuilder) BuildFromReference(_ interface{}) {
	_this.panicBadEvent("Reference")
}

func (_this *timeBuilder) PrepareForListContents() {
	_this.panicBadEvent("PrepareForListContents")
}

func (_this *timeBuilder) PrepareForMapContents() {
	_this.panicBadEvent("PrepareForMapContents")
}

func (_this *timeBuilder) NotifyChildContainerFinished(_ reflect.Value) {
	_this.panicBadEvent("NotifyChildContainerFinished")
}

// ============================================================================

type compactTimeBuilder struct {
	// Template Data
	session *Session
}

func newCompactTimeBuilder() ObjectBuilder {
	return &compactTimeBuilder{}
}

func (_this *compactTimeBuilder) String() string {
	return fmt.Sprintf("%v", reflect.TypeOf(_this))
}

func (_this *compactTimeBuilder) panicBadEvent(name string, args ...interface{}) {
	PanicBadEventWithType(_this, common.TypeCompactTime, name, args...)
}

func (_this *compactTimeBuilder) InitTemplate(session *Session) {
	_this.session = session
}

func (_this *compactTimeBuilder) NewInstance(_ *RootBuilder, _ ObjectBuilder, _ *options.BuilderOptions) ObjectBuilder {
	return _this
}

func (_this *compactTimeBuilder) SetParent(_ ObjectBuilder) {
}

func (_this *compactTimeBuilder) BuildFromNil(_ reflect.Value) {
	_this.panicBadEvent("Nil")
}

func (_this *compactTimeBuilder) BuildFromBool(_ bool, _ reflect.Value) {
	_this.panicBadEvent("Bool")
}

func (_this *compactTimeBuilder) BuildFromInt(_ int64, _ reflect.Value) {
	_this.panicBadEvent("Int")
}

func (_this *compactTimeBuilder) BuildFromUint(_ uint64, _ reflect.Value) {
	_this.panicBadEvent("Uint")
}

func (_this *compactTimeBuilder) BuildFromBigInt(_ *big.Int, _ reflect.Value) {
	_this.panicBadEvent("BigInt")
}

func (_this *compactTimeBuilder) BuildFromFloat(_ float64, _ reflect.Value) {
	_this.panicBadEvent("Float")
}

func (_this *compactTimeBuilder) BuildFromBigFloat(_ *big.Float, _ reflect.Value) {
	_this.panicBadEvent("BigFloat")
}

func (_this *compactTimeBuilder) BuildFromDecimalFloat(_ compact_float.DFloat, _ reflect.Value) {
	_this.panicBadEvent("DecimalFloat")
}

func (_this *compactTimeBuilder) BuildFromBigDecimalFloat(_ *apd.Decimal, _ reflect.Value) {
	_this.panicBadEvent("BigDecimalFloat")
}

func (_this *compactTimeBuilder) BuildFromUUID(_ []byte, _ reflect.Value) {
	_this.panicBadEvent("UUID")
}

func (_this *compactTimeBuilder) BuildFromArray(arrayType events.ArrayType, value []byte, dst reflect.Value) {
	if !_this.session.TryBuildFromCustom(_this, arrayType, value, dst) {
		_this.panicBadEvent("TypedArray(%v)", arrayType)
	}
}

func (_this *compactTimeBuilder) BuildFromTime(value time.Time, dst reflect.Value) {
	t, err := compact_time.AsCompactTime(value)
	if err != nil {
		panic(err)
	}
	dst.Set(reflect.ValueOf(*t))
}

func (_this *compactTimeBuilder) BuildFromCompactTime(value *compact_time.Time, dst reflect.Value) {
	dst.Set(reflect.ValueOf(*value))
}

func (_this *compactTimeBuilder) BuildBeginList() {
	_this.panicBadEvent("List")
}

func (_this *compactTimeBuilder) BuildBeginMap() {
	_this.panicBadEvent("Map")
}

func (_this *compactTimeBuilder) BuildEndContainer() {
	_this.panicBadEvent("ContainerEnd")
}

func (_this *compactTimeBuilder) BuildBeginMarker(_ interface{}) {
	_this.panicBadEvent("Marker")
}

func (_this *compactTimeBuilder) BuildFromReference(_ interface{}) {
	_this.panicBadEvent("Reference")
}

func (_this *compactTimeBuilder) PrepareForListContents() {
	_this.panicBadEvent("PrepareForListContents")
}

func (_this *compactTimeBuilder) PrepareForMapContents() {
	_this.panicBadEvent("PrepareForMapContents")
}

func (_this *compactTimeBuilder) NotifyChildContainerFinished(_ reflect.Value) {
	_this.panicBadEvent("NotifyChildContainerFinished")
}

// ============================================================================

type pCompactTimeBuilder struct {
	// Template Data
	session *Session
}

func newPCompactTimeBuilder() ObjectBuilder {
	return &pCompactTimeBuilder{}
}

func (_this *pCompactTimeBuilder) String() string {
	return fmt.Sprintf("%v", reflect.TypeOf(_this))
}

func (_this *pCompactTimeBuilder) panicBadEvent(name string, args ...interface{}) {
	PanicBadEventWithType(_this, common.TypePCompactTime, name, args...)
}

func (_this *pCompactTimeBuilder) InitTemplate(session *Session) {
	_this.session = session
}

func (_this *pCompactTimeBuilder) NewInstance(_ *RootBuilder, _ ObjectBuilder, _ *options.BuilderOptions) ObjectBuilder {
	return _this
}

func (_this *pCompactTimeBuilder) SetParent(_ ObjectBuilder) {
}

func (_this *pCompactTimeBuilder) BuildFromNil(dst reflect.Value) {
	dst.Set(reflect.ValueOf((*compact_time.Time)(nil)))
}

func (_this *pCompactTimeBuilder) BuildFromBool(_ bool, _ reflect.Value) {
	_this.panicBadEvent("Bool")
}

func (_this *pCompactTimeBuilder) BuildFromInt(_ int64, _ reflect.Value) {
	_this.panicBadEvent("Int")
}

func (_this *pCompactTimeBuilder) BuildFromUint(_ uint64, _ reflect.Value) {
	_this.panicBadEvent("Uint")
}

func (_this *pCompactTimeBuilder) BuildFromBigInt(_ *big.Int, _ reflect.Value) {
	_this.panicBadEvent("BigInt")
}

func (_this *pCompactTimeBuilder) BuildFromFloat(_ float64, _ reflect.Value) {
	_this.panicBadEvent("Float")
}

func (_this *pCompactTimeBuilder) BuildFromBigFloat(_ *big.Float, _ reflect.Value) {
	_this.panicBadEvent("BigFloat")
}

func (_this *pCompactTimeBuilder) BuildFromDecimalFloat(_ compact_float.DFloat, _ reflect.Value) {
	_this.panicBadEvent("DecimalFloat")
}

func (_this *pCompactTimeBuilder) BuildFromBigDecimalFloat(_ *apd.Decimal, _ reflect.Value) {
	_this.panicBadEvent("BigDecimalFloat")
}

func (_this *pCompactTimeBuilder) BuildFromUUID(_ []byte, _ reflect.Value) {
	_this.panicBadEvent("UUID")
}

func (_this *pCompactTimeBuilder) BuildFromArray(arrayType events.ArrayType, value []byte, dst reflect.Value) {
	if !_this.session.TryBuildFromCustom(_this, arrayType, value, dst) {
		_this.panicBadEvent("TypedArray(%v)", arrayType)
	}
}

func (_this *pCompactTimeBuilder) BuildFromTime(value time.Time, dst reflect.Value) {
	t, err := compact_time.AsCompactTime(value)
	if err != nil {
		panic(err)
	}
	dst.Set(reflect.ValueOf(t))
}

func (_this *pCompactTimeBuilder) BuildFromCompactTime(value *compact_time.Time, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (_this *pCompactTimeBuilder) BuildBeginList() {
	_this.panicBadEvent("List")
}

func (_this *pCompactTimeBuilder) BuildBeginMap() {
	_this.panicBadEvent("Map")
}

func (_this *pCompactTimeBuilder) BuildEndContainer() {
	_this.panicBadEvent("ContainerEnd")
}

func (_this *pCompactTimeBuilder) BuildBeginMarker(_ interface{}) {
	_this.panicBadEvent("Marker")
}

func (_this *pCompactTimeBuilder) BuildFromReference(_ interface{}) {
	_this.panicBadEvent("Reference")
}

func (_this *pCompactTimeBuilder) PrepareForListContents() {
	_this.panicBadEvent("PrepareForListContents")
}

func (_this *pCompactTimeBuilder) PrepareForMapContents() {
	_this.panicBadEvent("PrepareForMapContents")
}

func (_this *pCompactTimeBuilder) NotifyChildContainerFinished(_ reflect.Value) {
	_this.panicBadEvent("NotifyChildContainerFinished")
}
