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

type bigIntBuilder struct {
	// Template Data
	dstType reflect.Type

	// Instance Data
	opts *options.BuilderOptions
}

func newBigIntBuilder(dstType reflect.Type) ObjectBuilder {
	return &bigIntBuilder{
		dstType: dstType,
	}
}
func (_this *bigIntBuilder) String() string {
	return fmt.Sprintf("%v", reflect.TypeOf(_this))
}
func (_this *bigIntBuilder) panicBadEvent(name string, args ...interface{}) {
	PanicBadEventWithType(_this, _this.dstType, name, args...)
}
func (_this *bigIntBuilder) InitTemplate(_ *Session) {}
func (_this *bigIntBuilder) NewInstance(_ *RootBuilder, _ ObjectBuilder, opts *options.BuilderOptions) ObjectBuilder {
	return &bigIntBuilder{
		opts: opts,
	}
}
func (_this *bigIntBuilder) SetParent(_ ObjectBuilder) {}
func (_this *bigIntBuilder) BuildFromNil(_ reflect.Value) {
	_this.panicBadEvent("BuildFromNil")
}
func (_this *bigIntBuilder) BuildFromBool(_ bool, _ reflect.Value) {
	_this.panicBadEvent("BuildFromBool")
}
func (_this *bigIntBuilder) BuildFromInt(value int64, dst reflect.Value) {
	setBigIntFromInt(value, dst)
}
func (_this *bigIntBuilder) BuildFromUint(value uint64, dst reflect.Value) {
	setBigIntFromUint(value, dst)
}
func (_this *bigIntBuilder) BuildFromBigInt(value *big.Int, dst reflect.Value) {
	setBigIntFromBigInt(value, dst)
}
func (_this *bigIntBuilder) BuildFromFloat(value float64, dst reflect.Value) {
	setBigIntFromFloat(value, dst, _this.opts.FloatToBigIntMaxBase2Exponent)
}
func (_this *bigIntBuilder) BuildFromBigFloat(value *big.Float, dst reflect.Value) {
	setBigIntFromBigFloat(value, dst, _this.opts.FloatToBigIntMaxBase2Exponent)
}
func (_this *bigIntBuilder) BuildFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	setBigIntFromDecimalFloat(value, dst, _this.opts.FloatToBigIntMaxBase10Exponent)
}
func (_this *bigIntBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	setBigIntFromBigDecimalFloat(value, dst, _this.opts.FloatToBigIntMaxBase10Exponent)
}
func (_this *bigIntBuilder) BuildFromUUID(_ []byte, _ reflect.Value) {
	_this.panicBadEvent("BuildFromUUID")
}
func (_this *bigIntBuilder) BuildFromArray(arrayType events.ArrayType, value []byte, dst reflect.Value) {
	_this.panicBadEvent("TypedArray(%v)", arrayType)
}
func (_this *bigIntBuilder) BuildFromTime(_ time.Time, _ reflect.Value) {
	_this.panicBadEvent("BuildFromTime")
}
func (_this *bigIntBuilder) BuildFromCompactTime(_ *compact_time.Time, _ reflect.Value) {
	_this.panicBadEvent("BuildFromCompactTime")
}
func (_this *bigIntBuilder) BuildBeginList() {
	_this.panicBadEvent("BuildBeginList")
}
func (_this *bigIntBuilder) BuildBeginMap() {
	_this.panicBadEvent("BuildBeginMap")
}
func (_this *bigIntBuilder) BuildEndContainer() {
	_this.panicBadEvent("BuildEndContainer")
}
func (_this *bigIntBuilder) BuildBeginMarker(_ interface{}) {
	_this.panicBadEvent("BuildBeginMarker")
}
func (_this *bigIntBuilder) BuildFromReference(_ interface{}) {
	_this.panicBadEvent("BuildFromReference")
}
func (_this *bigIntBuilder) PrepareForListContents() {
	_this.panicBadEvent("PrepareForListContents")
}
func (_this *bigIntBuilder) PrepareForMapContents() {
	_this.panicBadEvent("PrepareForMapContents")
}
func (_this *bigIntBuilder) NotifyChildContainerFinished(_ reflect.Value) {
	_this.panicBadEvent("NotifyChildContainerFinished")
}

// ============================================================================

type pBigIntBuilder struct {
	// Template Data
	dstType reflect.Type

	// Instance Data
	opts *options.BuilderOptions
}

func newPBigIntBuilder(dstType reflect.Type) ObjectBuilder {
	return &pBigIntBuilder{
		dstType: dstType,
	}
}
func (_this *pBigIntBuilder) String() string {
	return fmt.Sprintf("%v", reflect.TypeOf(_this))
}
func (_this *pBigIntBuilder) panicBadEvent(name string, args ...interface{}) {
	PanicBadEventWithType(_this, _this.dstType, name, args...)
}
func (_this *pBigIntBuilder) InitTemplate(_ *Session) {}
func (_this *pBigIntBuilder) NewInstance(_ *RootBuilder, _ ObjectBuilder, opts *options.BuilderOptions) ObjectBuilder {
	return &pBigIntBuilder{
		opts: opts,
	}
}
func (_this *pBigIntBuilder) SetParent(_ ObjectBuilder) {}
func (_this *pBigIntBuilder) BuildFromNil(dst reflect.Value) {
	dst.Set(reflect.ValueOf((*big.Int)(nil)))
}
func (_this *pBigIntBuilder) BuildFromBool(_ bool, _ reflect.Value) {
	_this.panicBadEvent("BuildFromBool")
}
func (_this *pBigIntBuilder) BuildFromInt(value int64, dst reflect.Value) {
	setPBigIntFromInt(value, dst)
}
func (_this *pBigIntBuilder) BuildFromUint(value uint64, dst reflect.Value) {
	setPBigIntFromUint(value, dst)
}
func (_this *pBigIntBuilder) BuildFromBigInt(value *big.Int, dst reflect.Value) {
	setPBigIntFromBigInt(value, dst)
}
func (_this *pBigIntBuilder) BuildFromFloat(value float64, dst reflect.Value) {
	setPBigIntFromFloat(value, dst, _this.opts.FloatToBigIntMaxBase2Exponent)
}
func (_this *pBigIntBuilder) BuildFromBigFloat(value *big.Float, dst reflect.Value) {
	setPBigIntFromBigFloat(value, dst, _this.opts.FloatToBigIntMaxBase2Exponent)
}
func (_this *pBigIntBuilder) BuildFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	setPBigIntFromDecimalFloat(value, dst, _this.opts.FloatToBigIntMaxBase10Exponent)
}
func (_this *pBigIntBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	setPBigIntFromBigDecimalFloat(value, dst, _this.opts.FloatToBigIntMaxBase10Exponent)
}
func (_this *pBigIntBuilder) BuildFromUUID(_ []byte, _ reflect.Value) {
	_this.panicBadEvent("BuildFromUUID")
}
func (_this *pBigIntBuilder) BuildFromArray(arrayType events.ArrayType, value []byte, dst reflect.Value) {
	_this.panicBadEvent("TypedArray(%v)", arrayType)
}
func (_this *pBigIntBuilder) BuildFromTime(_ time.Time, _ reflect.Value) {
	_this.panicBadEvent("BuildFromTime")
}
func (_this *pBigIntBuilder) BuildFromCompactTime(_ *compact_time.Time, _ reflect.Value) {
	_this.panicBadEvent("BuildFromCompactTime")
}
func (_this *pBigIntBuilder) BuildBeginList() {
	_this.panicBadEvent("BuildBeginList")
}
func (_this *pBigIntBuilder) BuildBeginMap() {
	_this.panicBadEvent("BuildBeginMap")
}
func (_this *pBigIntBuilder) BuildEndContainer() {
	_this.panicBadEvent("BuildEndContainer")
}
func (_this *pBigIntBuilder) BuildBeginMarker(_ interface{}) {
	_this.panicBadEvent("BuildBeginMarker")
}
func (_this *pBigIntBuilder) BuildFromReference(_ interface{}) {
	_this.panicBadEvent("BuildFromReference")
}
func (_this *pBigIntBuilder) PrepareForListContents() {
	_this.panicBadEvent("PrepareForListContents")
}
func (_this *pBigIntBuilder) PrepareForMapContents() {
	_this.panicBadEvent("PrepareForMapContents")
}
func (_this *pBigIntBuilder) NotifyChildContainerFinished(_ reflect.Value) {
	_this.panicBadEvent("NotifyChildContainerFinished")
}
