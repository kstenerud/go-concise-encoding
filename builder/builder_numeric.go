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

	"github.com/kstenerud/go-concise-encoding/options"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
)

// ============================================================================

type bigDecimalFloatBuilder struct {
	// Template Data
	dstType reflect.Type
}

func newBigDecimalFloatBuilder(dstType reflect.Type) ObjectBuilder {
	return &bigDecimalFloatBuilder{
		dstType: dstType,
	}
}
func (_this *bigDecimalFloatBuilder) String() string {
	return fmt.Sprintf("%v", reflect.TypeOf(_this))
}
func (_this *bigDecimalFloatBuilder) panicBadEvent(name string, args ...interface{}) {
	PanicBadEventWithType(_this, _this.dstType, name, args...)
}
func (_this *bigDecimalFloatBuilder) InitTemplate(_ *Session) {}
func (_this *bigDecimalFloatBuilder) NewInstance(_ *RootBuilder, _ ObjectBuilder, _ *options.BuilderOptions) ObjectBuilder {
	return _this
}
func (_this *bigDecimalFloatBuilder) SetParent(_ ObjectBuilder) {}
func (_this *bigDecimalFloatBuilder) BuildFromInt(value int64, dst reflect.Value) {
	setBigDecimalFloatFromInt(value, dst)
}
func (_this *bigDecimalFloatBuilder) BuildFromUint(value uint64, dst reflect.Value) {
	setBigDecimalFloatFromUint(value, dst)
}
func (_this *bigDecimalFloatBuilder) BuildFromBigInt(value *big.Int, dst reflect.Value) {
	setBigDecimalFloatFromBigInt(value, dst)
}
func (_this *bigDecimalFloatBuilder) BuildFromFloat(value float64, dst reflect.Value) {
	setBigDecimalFloatFromFloat(value, dst)
}
func (_this *bigDecimalFloatBuilder) BuildFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	setBigDecimalFloatFromDecimalFloat(value, dst)
}
func (_this *bigDecimalFloatBuilder) BuildFromBigFloat(value *big.Float, dst reflect.Value) {
	setBigDecimalFloatFromBigFloat(value, dst)
}
func (_this *bigDecimalFloatBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	setBigDecimalFloatFromBigDecimalFloat(value, dst)
}

// ============================================================================

type bigFloatBuilder struct {
	// Template Data
	dstType reflect.Type
}

func newBigFloatBuilder(dstType reflect.Type) ObjectBuilder {
	return &bigFloatBuilder{
		dstType: dstType,
	}
}
func (_this *bigFloatBuilder) String() string {
	return fmt.Sprintf("%v", reflect.TypeOf(_this))
}
func (_this *bigFloatBuilder) panicBadEvent(name string, args ...interface{}) {
	PanicBadEventWithType(_this, _this.dstType, name, args...)
}
func (_this *bigFloatBuilder) InitTemplate(_ *Session) {}
func (_this *bigFloatBuilder) NewInstance(_ *RootBuilder, _ ObjectBuilder, _ *options.BuilderOptions) ObjectBuilder {
	return _this
}
func (_this *bigFloatBuilder) SetParent(_ ObjectBuilder) {}
func (_this *bigFloatBuilder) BuildFromInt(value int64, dst reflect.Value) {
	setBigFloatFromInt(value, dst)
}
func (_this *bigFloatBuilder) BuildFromUint(value uint64, dst reflect.Value) {
	setBigFloatFromUint(value, dst)
}
func (_this *bigFloatBuilder) BuildFromBigInt(value *big.Int, dst reflect.Value) {
	setBigFloatFromBigInt(value, dst)
}
func (_this *bigFloatBuilder) BuildFromFloat(value float64, dst reflect.Value) {
	setBigFloatFromFloat(value, dst)
}
func (_this *bigFloatBuilder) BuildFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	setBigFloatFromDecimalFloat(value, dst)
}
func (_this *bigFloatBuilder) BuildFromBigFloat(value *big.Float, dst reflect.Value) {
	setBigFloatFromBigFloat(value, dst)
}
func (_this *bigFloatBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	setBigFloatFromBigDecimalFloat(value, dst)
}

// ============================================================================

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

// ============================================================================

type decimalFloatBuilder struct {
	// Template Data
	dstType reflect.Type
}

func newDecimalFloatBuilder(dstType reflect.Type) ObjectBuilder {
	return &decimalFloatBuilder{
		dstType: dstType,
	}
}
func (_this *decimalFloatBuilder) String() string {
	return fmt.Sprintf("%v", reflect.TypeOf(_this))
}
func (_this *decimalFloatBuilder) panicBadEvent(name string, args ...interface{}) {
	PanicBadEventWithType(_this, _this.dstType, name, args...)
}
func (_this *decimalFloatBuilder) InitTemplate(_ *Session) {}
func (_this *decimalFloatBuilder) NewInstance(_ *RootBuilder, _ ObjectBuilder, _ *options.BuilderOptions) ObjectBuilder {
	return _this
}
func (_this *decimalFloatBuilder) SetParent(_ ObjectBuilder) {}
func (_this *decimalFloatBuilder) BuildFromInt(value int64, dst reflect.Value) {
	setDecimalFloatFromInt(value, dst)
}
func (_this *decimalFloatBuilder) BuildFromUint(value uint64, dst reflect.Value) {
	setDecimalFloatFromUint(value, dst)
}
func (_this *decimalFloatBuilder) BuildFromBigInt(value *big.Int, dst reflect.Value) {
	setDecimalFloatFromBigInt(value, dst)
}
func (_this *decimalFloatBuilder) BuildFromFloat(value float64, dst reflect.Value) {
	setDecimalFloatFromFloat(value, dst)
}
func (_this *decimalFloatBuilder) BuildFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	setDecimalFloatFromDecimalFloat(value, dst)
}
func (_this *decimalFloatBuilder) BuildFromBigFloat(value *big.Float, dst reflect.Value) {
	setDecimalFloatFromBigFloat(value, dst)
}
func (_this *decimalFloatBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	setDecimalFloatFromBigDecimalFloat(value, dst)
}

// ============================================================================

type floatBuilder struct {
	// Template Data
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
func (_this *floatBuilder) panicBadEvent(name string, args ...interface{}) {
	PanicBadEventWithType(_this, _this.dstType, name, args...)
}
func (_this *floatBuilder) InitTemplate(_ *Session) {}
func (_this *floatBuilder) NewInstance(_ *RootBuilder, _ ObjectBuilder, _ *options.BuilderOptions) ObjectBuilder {
	return _this
}
func (_this *floatBuilder) SetParent(_ ObjectBuilder) {}
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
func (_this *floatBuilder) BuildFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	setFloatFromDecimalFloat(value, dst)
}
func (_this *floatBuilder) BuildFromBigFloat(value *big.Float, dst reflect.Value) {
	setFloatFromBigFloat(value, dst)
}
func (_this *floatBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	setFloatFromBigDecimalFloat(value, dst)
}

// ============================================================================

type intBuilder struct {
	// Template Data
	dstType reflect.Type
}

func newIntBuilder(dstType reflect.Type) ObjectBuilder {
	return &intBuilder{
		dstType: dstType,
	}
}
func (_this *intBuilder) String() string {
	return fmt.Sprintf("%v", reflect.TypeOf(_this))
}
func (_this *intBuilder) panicBadEvent(name string, args ...interface{}) {
	PanicBadEventWithType(_this, _this.dstType, name, args...)
}
func (_this *intBuilder) InitTemplate(_ *Session) {}
func (_this *intBuilder) NewInstance(_ *RootBuilder, _ ObjectBuilder, _ *options.BuilderOptions) ObjectBuilder {
	return _this
}
func (_this *intBuilder) SetParent(_ ObjectBuilder) {}
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
func (_this *intBuilder) BuildFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	setIntFromDecimalFloat(value, dst)
}
func (_this *intBuilder) BuildFromBigFloat(value *big.Float, dst reflect.Value) {
	setIntFromBigFloat(value, dst)
}
func (_this *intBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	setIntFromBigDecimalFloat(value, dst)
}

// ============================================================================

type pBigDecimalFloatBuilder struct {
	// Template Data
	dstType reflect.Type
}

func newPBigDecimalFloatBuilder(dstType reflect.Type) ObjectBuilder {
	return &pBigDecimalFloatBuilder{
		dstType: dstType,
	}
}
func (_this *pBigDecimalFloatBuilder) String() string {
	return fmt.Sprintf("%v", reflect.TypeOf(_this))
}
func (_this *pBigDecimalFloatBuilder) panicBadEvent(name string, args ...interface{}) {
	PanicBadEventWithType(_this, _this.dstType, name, args...)
}
func (_this *pBigDecimalFloatBuilder) InitTemplate(_ *Session) {}
func (_this *pBigDecimalFloatBuilder) NewInstance(_ *RootBuilder, _ ObjectBuilder, _ *options.BuilderOptions) ObjectBuilder {
	return _this
}
func (_this *pBigDecimalFloatBuilder) SetParent(_ ObjectBuilder) {}
func (_this *pBigDecimalFloatBuilder) BuildFromNil(dst reflect.Value) {
	dst.Set(reflect.ValueOf((*apd.Decimal)(nil)))
}
func (_this *pBigDecimalFloatBuilder) BuildFromInt(value int64, dst reflect.Value) {
	setPBigDecimalFloatFromInt(value, dst)
}
func (_this *pBigDecimalFloatBuilder) BuildFromUint(value uint64, dst reflect.Value) {
	setPBigDecimalFloatFromUint(value, dst)
}
func (_this *pBigDecimalFloatBuilder) BuildFromBigInt(value *big.Int, dst reflect.Value) {
	setPBigDecimalFloatFromBigInt(value, dst)
}
func (_this *pBigDecimalFloatBuilder) BuildFromFloat(value float64, dst reflect.Value) {
	setPBigDecimalFloatFromFloat(value, dst)
}
func (_this *pBigDecimalFloatBuilder) BuildFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	setPBigDecimalFloatFromDecimalFloat(value, dst)
}
func (_this *pBigDecimalFloatBuilder) BuildFromBigFloat(value *big.Float, dst reflect.Value) {
	setPBigDecimalFloatFromBigFloat(value, dst)
}
func (_this *pBigDecimalFloatBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	setPBigDecimalFloatFromBigDecimalFloat(value, dst)
}

// ============================================================================

type pBigFloatBuilder struct {
	// Template Data
	dstType reflect.Type
}

func newPBigFloatBuilder(dstType reflect.Type) ObjectBuilder {
	return &pBigFloatBuilder{
		dstType: dstType,
	}
}
func (_this *pBigFloatBuilder) String() string {
	return fmt.Sprintf("%v", reflect.TypeOf(_this))
}
func (_this *pBigFloatBuilder) panicBadEvent(name string, args ...interface{}) {
	PanicBadEventWithType(_this, _this.dstType, name, args...)
}
func (_this *pBigFloatBuilder) InitTemplate(_ *Session) {}
func (_this *pBigFloatBuilder) NewInstance(_ *RootBuilder, _ ObjectBuilder, _ *options.BuilderOptions) ObjectBuilder {
	return _this
}
func (_this *pBigFloatBuilder) SetParent(_ ObjectBuilder) {}
func (_this *pBigFloatBuilder) BuildFromNil(dst reflect.Value) {
	dst.Set(reflect.ValueOf((*big.Float)(nil)))
}
func (_this *pBigFloatBuilder) BuildFromInt(value int64, dst reflect.Value) {
	setPBigFloatFromInt(value, dst)
}
func (_this *pBigFloatBuilder) BuildFromUint(value uint64, dst reflect.Value) {
	setPBigFloatFromUint(value, dst)
}
func (_this *pBigFloatBuilder) BuildFromBigInt(value *big.Int, dst reflect.Value) {
	setPBigFloatFromBigInt(value, dst)
}
func (_this *pBigFloatBuilder) BuildFromFloat(value float64, dst reflect.Value) {
	setPBigFloatFromFloat(value, dst)
}
func (_this *pBigFloatBuilder) BuildFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	setPBigFloatFromDecimalFloat(value, dst)
}
func (_this *pBigFloatBuilder) BuildFromBigFloat(value *big.Float, dst reflect.Value) {
	setPBigFloatFromBigFloat(value, dst)
}
func (_this *pBigFloatBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	setPBigFloatFromBigDecimalFloat(value, dst)
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

// ============================================================================

type uintBuilder struct {
	// Template Data
	dstType reflect.Type
}

func newUintBuilder(dstType reflect.Type) ObjectBuilder {
	return &uintBuilder{
		dstType: dstType,
	}
}
func (_this *uintBuilder) String() string {
	return fmt.Sprintf("%v", reflect.TypeOf(_this))
}
func (_this *uintBuilder) panicBadEvent(name string, args ...interface{}) {
	PanicBadEventWithType(_this, _this.dstType, name, args...)
}
func (_this *uintBuilder) InitTemplate(_ *Session) {}
func (_this *uintBuilder) NewInstance(_ *RootBuilder, _ ObjectBuilder, _ *options.BuilderOptions) ObjectBuilder {
	return _this
}
func (_this *uintBuilder) SetParent(_ ObjectBuilder) {}
func (_this *uintBuilder) BuildFromInt(value int64, dst reflect.Value) {
	setUintFromInt(value, dst)
}
func (_this *uintBuilder) BuildFromUint(value uint64, dst reflect.Value) {
	setUintFromUint(value, dst)
}
func (_this *uintBuilder) BuildFromBigInt(value *big.Int, dst reflect.Value) {
	setUintFromBigInt(value, dst)
}
func (_this *uintBuilder) BuildFromFloat(value float64, dst reflect.Value) {
	setUintFromFloat(value, dst)
}
func (_this *uintBuilder) BuildFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	setUintFromDecimalFloat(value, dst)
}
func (_this *uintBuilder) BuildFromBigFloat(value *big.Float, dst reflect.Value) {
	setUintFromBigFloat(value, dst)
}
func (_this *uintBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	setUintFromBigDecimalFloat(value, dst)
}
