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
	"math/big"
	"reflect"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
)

// ============================================================================

type bigDecimalFloatBuilder struct{}

var globalBigDecimalFloatBuilder = &bigDecimalFloatBuilder{}

func generateBigDecimalFloatBuilder(ctx *Context) ObjectBuilder { return globalBigDecimalFloatBuilder }
func (_this *bigDecimalFloatBuilder) String() string            { return reflect.TypeOf(_this).String() }

func (_this *bigDecimalFloatBuilder) BuildFromInt(ctx *Context, value int64, dst reflect.Value) reflect.Value {
	setBigDecimalFloatFromInt(value, dst)
	return dst
}
func (_this *bigDecimalFloatBuilder) BuildFromUint(ctx *Context, value uint64, dst reflect.Value) reflect.Value {
	setBigDecimalFloatFromUint(value, dst)
	return dst
}
func (_this *bigDecimalFloatBuilder) BuildFromBigInt(ctx *Context, value *big.Int, dst reflect.Value) reflect.Value {
	setBigDecimalFloatFromBigInt(value, dst)
	return dst
}
func (_this *bigDecimalFloatBuilder) BuildFromFloat(ctx *Context, value float64, dst reflect.Value) reflect.Value {
	setBigDecimalFloatFromFloat(value, dst)
	return dst
}
func (_this *bigDecimalFloatBuilder) BuildFromDecimalFloat(ctx *Context, value compact_float.DFloat, dst reflect.Value) reflect.Value {
	setBigDecimalFloatFromDecimalFloat(value, dst)
	return dst
}
func (_this *bigDecimalFloatBuilder) BuildFromBigFloat(ctx *Context, value *big.Float, dst reflect.Value) reflect.Value {
	setBigDecimalFloatFromBigFloat(value, dst)
	return dst
}
func (_this *bigDecimalFloatBuilder) BuildFromBigDecimalFloat(ctx *Context, value *apd.Decimal, dst reflect.Value) reflect.Value {
	setBigDecimalFloatFromBigDecimalFloat(value, dst)
	return dst
}

// ============================================================================

type bigFloatBuilder struct{}

var globalBigFloatBuilder = &bigFloatBuilder{}

func generateBigFloatBuilder(ctx *Context) ObjectBuilder { return globalBigFloatBuilder }
func (_this *bigFloatBuilder) String() string            { return reflect.TypeOf(_this).String() }

func (_this *bigFloatBuilder) BuildFromInt(ctx *Context, value int64, dst reflect.Value) reflect.Value {
	setBigFloatFromInt(value, dst)
	return dst
}
func (_this *bigFloatBuilder) BuildFromUint(ctx *Context, value uint64, dst reflect.Value) reflect.Value {
	setBigFloatFromUint(value, dst)
	return dst
}
func (_this *bigFloatBuilder) BuildFromBigInt(ctx *Context, value *big.Int, dst reflect.Value) reflect.Value {
	setBigFloatFromBigInt(value, dst)
	return dst
}
func (_this *bigFloatBuilder) BuildFromFloat(ctx *Context, value float64, dst reflect.Value) reflect.Value {
	setBigFloatFromFloat(value, dst)
	return dst
}
func (_this *bigFloatBuilder) BuildFromDecimalFloat(ctx *Context, value compact_float.DFloat, dst reflect.Value) reflect.Value {
	setBigFloatFromDecimalFloat(value, dst)
	return dst
}
func (_this *bigFloatBuilder) BuildFromBigFloat(ctx *Context, value *big.Float, dst reflect.Value) reflect.Value {
	setBigFloatFromBigFloat(value, dst)
	return dst
}
func (_this *bigFloatBuilder) BuildFromBigDecimalFloat(ctx *Context, value *apd.Decimal, dst reflect.Value) reflect.Value {
	setBigFloatFromBigDecimalFloat(value, dst)
	return dst
}

// ============================================================================

type bigIntBuilder struct{}

var globalBigIntBuilder = &bigIntBuilder{}

func generateBigIntBuilder(ctx *Context) ObjectBuilder { return globalBigIntBuilder }
func (_this *bigIntBuilder) String() string            { return reflect.TypeOf(_this).String() }

func (_this *bigIntBuilder) BuildFromInt(ctx *Context, value int64, dst reflect.Value) reflect.Value {
	setBigIntFromInt(value, dst)
	return dst
}
func (_this *bigIntBuilder) BuildFromUint(ctx *Context, value uint64, dst reflect.Value) reflect.Value {
	setBigIntFromUint(value, dst)
	return dst
}
func (_this *bigIntBuilder) BuildFromBigInt(ctx *Context, value *big.Int, dst reflect.Value) reflect.Value {
	setBigIntFromBigInt(value, dst)
	return dst
}
func (_this *bigIntBuilder) BuildFromFloat(ctx *Context, value float64, dst reflect.Value) reflect.Value {
	setBigIntFromFloat(value, dst, ctx.Options.FloatToBigIntMaxBase2Exponent)
	return dst
}
func (_this *bigIntBuilder) BuildFromBigFloat(ctx *Context, value *big.Float, dst reflect.Value) reflect.Value {
	setBigIntFromBigFloat(value, dst, ctx.Options.FloatToBigIntMaxBase2Exponent)
	return dst
}
func (_this *bigIntBuilder) BuildFromDecimalFloat(ctx *Context, value compact_float.DFloat, dst reflect.Value) reflect.Value {
	setBigIntFromDecimalFloat(value, dst, ctx.Options.FloatToBigIntMaxBase10Exponent)
	return dst
}
func (_this *bigIntBuilder) BuildFromBigDecimalFloat(ctx *Context, value *apd.Decimal, dst reflect.Value) reflect.Value {
	setBigIntFromBigDecimalFloat(value, dst, ctx.Options.FloatToBigIntMaxBase10Exponent)
	return dst
}

// ============================================================================

type decimalFloatBuilder struct{}

var globalDecimalFloatBuilder = &decimalFloatBuilder{}

func generateDecimalFloatBuilder(ctx *Context) ObjectBuilder { return globalDecimalFloatBuilder }
func (_this *decimalFloatBuilder) String() string            { return reflect.TypeOf(_this).String() }

func (_this *decimalFloatBuilder) BuildFromInt(ctx *Context, value int64, dst reflect.Value) reflect.Value {
	setDecimalFloatFromInt(value, dst)
	return dst
}
func (_this *decimalFloatBuilder) BuildFromUint(ctx *Context, value uint64, dst reflect.Value) reflect.Value {
	setDecimalFloatFromUint(value, dst)
	return dst
}
func (_this *decimalFloatBuilder) BuildFromBigInt(ctx *Context, value *big.Int, dst reflect.Value) reflect.Value {
	setDecimalFloatFromBigInt(value, dst)
	return dst
}
func (_this *decimalFloatBuilder) BuildFromFloat(ctx *Context, value float64, dst reflect.Value) reflect.Value {
	setDecimalFloatFromFloat(value, dst)
	return dst
}
func (_this *decimalFloatBuilder) BuildFromDecimalFloat(ctx *Context, value compact_float.DFloat, dst reflect.Value) reflect.Value {
	setDecimalFloatFromDecimalFloat(value, dst)
	return dst
}
func (_this *decimalFloatBuilder) BuildFromBigFloat(ctx *Context, value *big.Float, dst reflect.Value) reflect.Value {
	setDecimalFloatFromBigFloat(value, dst)
	return dst
}
func (_this *decimalFloatBuilder) BuildFromBigDecimalFloat(ctx *Context, value *apd.Decimal, dst reflect.Value) reflect.Value {
	setDecimalFloatFromBigDecimalFloat(value, dst)
	return dst
}

// ============================================================================

type floatBuilder struct{}

var globalFloatBuilder = &floatBuilder{}

func generateFloatBuilder(ctx *Context) ObjectBuilder { return globalFloatBuilder }
func (_this *floatBuilder) String() string            { return reflect.TypeOf(_this).String() }

func (_this *floatBuilder) BuildFromInt(ctx *Context, value int64, dst reflect.Value) reflect.Value {
	setFloatFromInt(value, dst)
	return dst
}
func (_this *floatBuilder) BuildFromUint(ctx *Context, value uint64, dst reflect.Value) reflect.Value {
	setFloatFromUint(value, dst)
	return dst
}
func (_this *floatBuilder) BuildFromBigInt(ctx *Context, value *big.Int, dst reflect.Value) reflect.Value {
	setFloatFromBigInt(value, dst)
	return dst
}
func (_this *floatBuilder) BuildFromFloat(ctx *Context, value float64, dst reflect.Value) reflect.Value {
	setFloatFromFloat(value, dst)
	return dst
}
func (_this *floatBuilder) BuildFromDecimalFloat(ctx *Context, value compact_float.DFloat, dst reflect.Value) reflect.Value {
	setFloatFromDecimalFloat(value, dst)
	return dst
}
func (_this *floatBuilder) BuildFromBigFloat(ctx *Context, value *big.Float, dst reflect.Value) reflect.Value {
	setFloatFromBigFloat(value, dst)
	return dst
}
func (_this *floatBuilder) BuildFromBigDecimalFloat(ctx *Context, value *apd.Decimal, dst reflect.Value) reflect.Value {
	setFloatFromBigDecimalFloat(value, dst)
	return dst
}

// ============================================================================

type intBuilder struct{}

var globalIntBuilder = &intBuilder{}

func generateIntBuilder(ctx *Context) ObjectBuilder { return globalIntBuilder }
func (_this *intBuilder) String() string            { return reflect.TypeOf(_this).String() }

func (_this *intBuilder) BuildFromInt(ctx *Context, value int64, dst reflect.Value) reflect.Value {
	setIntFromInt(value, dst)
	return dst
}
func (_this *intBuilder) BuildFromUint(ctx *Context, value uint64, dst reflect.Value) reflect.Value {
	setIntFromUint(value, dst)
	return dst
}
func (_this *intBuilder) BuildFromBigInt(ctx *Context, value *big.Int, dst reflect.Value) reflect.Value {
	setIntFromBigInt(value, dst)
	return dst
}
func (_this *intBuilder) BuildFromFloat(ctx *Context, value float64, dst reflect.Value) reflect.Value {
	setIntFromFloat(value, dst)
	return dst
}
func (_this *intBuilder) BuildFromDecimalFloat(ctx *Context, value compact_float.DFloat, dst reflect.Value) reflect.Value {
	setIntFromDecimalFloat(value, dst)
	return dst
}
func (_this *intBuilder) BuildFromBigFloat(ctx *Context, value *big.Float, dst reflect.Value) reflect.Value {
	setIntFromBigFloat(value, dst)
	return dst
}
func (_this *intBuilder) BuildFromBigDecimalFloat(ctx *Context, value *apd.Decimal, dst reflect.Value) reflect.Value {
	setIntFromBigDecimalFloat(value, dst)
	return dst
}

// ============================================================================

type pBigDecimalFloatBuilder struct{}

var globalPBigDecimalFloatBuilder = &pBigDecimalFloatBuilder{}

func generatePBigDecimalFloatBuilder(ctx *Context) ObjectBuilder { return globalPBigDecimalFloatBuilder }
func (_this *pBigDecimalFloatBuilder) String() string            { return reflect.TypeOf(_this).String() }

func (_this *pBigDecimalFloatBuilder) BuildFromNil(ctx *Context, dst reflect.Value) reflect.Value {
	dst.Set(reflect.ValueOf((*apd.Decimal)(nil)))
	ctx.NANext()
	return dst
}
func (_this *pBigDecimalFloatBuilder) BuildFromInt(ctx *Context, value int64, dst reflect.Value) reflect.Value {
	setPBigDecimalFloatFromInt(value, dst)
	return dst
}
func (_this *pBigDecimalFloatBuilder) BuildFromUint(ctx *Context, value uint64, dst reflect.Value) reflect.Value {
	setPBigDecimalFloatFromUint(value, dst)
	return dst
}
func (_this *pBigDecimalFloatBuilder) BuildFromBigInt(ctx *Context, value *big.Int, dst reflect.Value) reflect.Value {
	setPBigDecimalFloatFromBigInt(value, dst)
	return dst
}
func (_this *pBigDecimalFloatBuilder) BuildFromFloat(ctx *Context, value float64, dst reflect.Value) reflect.Value {
	setPBigDecimalFloatFromFloat(value, dst)
	return dst
}
func (_this *pBigDecimalFloatBuilder) BuildFromDecimalFloat(ctx *Context, value compact_float.DFloat, dst reflect.Value) reflect.Value {
	setPBigDecimalFloatFromDecimalFloat(value, dst)
	return dst
}
func (_this *pBigDecimalFloatBuilder) BuildFromBigFloat(ctx *Context, value *big.Float, dst reflect.Value) reflect.Value {
	setPBigDecimalFloatFromBigFloat(value, dst)
	return dst
}
func (_this *pBigDecimalFloatBuilder) BuildFromBigDecimalFloat(ctx *Context, value *apd.Decimal, dst reflect.Value) reflect.Value {
	setPBigDecimalFloatFromBigDecimalFloat(value, dst)
	return dst
}

// ============================================================================

type pBigFloatBuilder struct{}

var globalPBigFloatBuilder = &pBigFloatBuilder{}

func generatePBigFloatBuilder(ctx *Context) ObjectBuilder { return globalPBigFloatBuilder }
func (_this *pBigFloatBuilder) String() string            { return reflect.TypeOf(_this).String() }

func (_this *pBigFloatBuilder) BuildFromNil(ctx *Context, dst reflect.Value) reflect.Value {
	dst.Set(reflect.ValueOf((*big.Float)(nil)))
	ctx.NANext()
	return dst
}
func (_this *pBigFloatBuilder) BuildFromInt(ctx *Context, value int64, dst reflect.Value) reflect.Value {
	setPBigFloatFromInt(value, dst)
	return dst
}
func (_this *pBigFloatBuilder) BuildFromUint(ctx *Context, value uint64, dst reflect.Value) reflect.Value {
	setPBigFloatFromUint(value, dst)
	return dst
}
func (_this *pBigFloatBuilder) BuildFromBigInt(ctx *Context, value *big.Int, dst reflect.Value) reflect.Value {
	setPBigFloatFromBigInt(value, dst)
	return dst
}
func (_this *pBigFloatBuilder) BuildFromFloat(ctx *Context, value float64, dst reflect.Value) reflect.Value {
	setPBigFloatFromFloat(value, dst)
	return dst
}
func (_this *pBigFloatBuilder) BuildFromDecimalFloat(ctx *Context, value compact_float.DFloat, dst reflect.Value) reflect.Value {
	setPBigFloatFromDecimalFloat(value, dst)
	return dst
}
func (_this *pBigFloatBuilder) BuildFromBigFloat(ctx *Context, value *big.Float, dst reflect.Value) reflect.Value {
	setPBigFloatFromBigFloat(value, dst)
	return dst
}
func (_this *pBigFloatBuilder) BuildFromBigDecimalFloat(ctx *Context, value *apd.Decimal, dst reflect.Value) reflect.Value {
	setPBigFloatFromBigDecimalFloat(value, dst)
	return dst
}

// ============================================================================

type pBigIntBuilder struct{}

var globalPBigIntBuilder = &pBigIntBuilder{}

func generatePBigIntBuilder(ctx *Context) ObjectBuilder { return globalPBigIntBuilder }
func (_this *pBigIntBuilder) String() string            { return reflect.TypeOf(_this).String() }

func (_this *pBigIntBuilder) BuildFromNil(ctx *Context, dst reflect.Value) reflect.Value {
	dst.Set(reflect.ValueOf((*big.Int)(nil)))
	ctx.NANext()
	return dst
}
func (_this *pBigIntBuilder) BuildFromInt(ctx *Context, value int64, dst reflect.Value) reflect.Value {
	setPBigIntFromInt(value, dst)
	return dst
}
func (_this *pBigIntBuilder) BuildFromUint(ctx *Context, value uint64, dst reflect.Value) reflect.Value {
	setPBigIntFromUint(value, dst)
	return dst
}
func (_this *pBigIntBuilder) BuildFromBigInt(ctx *Context, value *big.Int, dst reflect.Value) reflect.Value {
	setPBigIntFromBigInt(value, dst)
	return dst
}
func (_this *pBigIntBuilder) BuildFromFloat(ctx *Context, value float64, dst reflect.Value) reflect.Value {
	setPBigIntFromFloat(value, dst, ctx.Options.FloatToBigIntMaxBase2Exponent)
	return dst
}
func (_this *pBigIntBuilder) BuildFromBigFloat(ctx *Context, value *big.Float, dst reflect.Value) reflect.Value {
	setPBigIntFromBigFloat(value, dst, ctx.Options.FloatToBigIntMaxBase2Exponent)
	return dst
}
func (_this *pBigIntBuilder) BuildFromDecimalFloat(ctx *Context, value compact_float.DFloat, dst reflect.Value) reflect.Value {
	setPBigIntFromDecimalFloat(value, dst, ctx.Options.FloatToBigIntMaxBase10Exponent)
	return dst
}
func (_this *pBigIntBuilder) BuildFromBigDecimalFloat(ctx *Context, value *apd.Decimal, dst reflect.Value) reflect.Value {
	setPBigIntFromBigDecimalFloat(value, dst, ctx.Options.FloatToBigIntMaxBase10Exponent)
	return dst
}

// ============================================================================

type uintBuilder struct {
	// Template Data
	dstType reflect.Type
}

var globalUintBuilder = &uintBuilder{}

func generateUintBuilder(ctx *Context) ObjectBuilder { return globalUintBuilder }
func (_this *uintBuilder) String() string            { return reflect.TypeOf(_this).String() }

func (_this *uintBuilder) BuildFromInt(ctx *Context, value int64, dst reflect.Value) reflect.Value {
	setUintFromInt(value, dst)
	return dst
}
func (_this *uintBuilder) BuildFromUint(ctx *Context, value uint64, dst reflect.Value) reflect.Value {
	setUintFromUint(value, dst)
	return dst
}
func (_this *uintBuilder) BuildFromBigInt(ctx *Context, value *big.Int, dst reflect.Value) reflect.Value {
	setUintFromBigInt(value, dst)
	return dst
}
func (_this *uintBuilder) BuildFromFloat(ctx *Context, value float64, dst reflect.Value) reflect.Value {
	setUintFromFloat(value, dst)
	return dst
}
func (_this *uintBuilder) BuildFromDecimalFloat(ctx *Context, value compact_float.DFloat, dst reflect.Value) reflect.Value {
	setUintFromDecimalFloat(value, dst)
	return dst
}
func (_this *uintBuilder) BuildFromBigFloat(ctx *Context, value *big.Float, dst reflect.Value) reflect.Value {
	setUintFromBigFloat(value, dst)
	return dst
}
func (_this *uintBuilder) BuildFromBigDecimalFloat(ctx *Context, value *apd.Decimal, dst reflect.Value) reflect.Value {
	setUintFromBigDecimalFloat(value, dst)
	return dst
}
