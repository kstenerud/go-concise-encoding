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
	"time"

	"github.com/kstenerud/go-concise-encoding/events"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
)

type urlBuilder struct{}

var globalUrlBuilder = &urlBuilder{}

func newUrlBuilderGenerator() BuilderGenerator {
	return func(ctx *Context) ObjectBuilder {
		builder := globalUrlBuilder
		ctx.StackBuilder(builder)
		return builder
	}
}

func (_this *urlBuilder) String() string { return reflect.TypeOf(_this).String() }

func (_this *urlBuilder) BuildFromArray(ctx *Context, arrayType events.ArrayType, value []byte, dst reflect.Value) reflect.Value {
	switch arrayType {
	case events.ArrayTypeResourceID:
		setRIDFromString(string(value), dst)
	default:
		ctx.UnstackBuilder()
		return ctx.CurrentBuilder.BuildFromArray(ctx, arrayType, value, dst)
	}
	return dst
}

func (_this *urlBuilder) BuildFromNil(ctx *Context, dst reflect.Value) reflect.Value {
	ctx.UnstackBuilder()
	return ctx.CurrentBuilder.BuildFromNil(ctx, dst)
}
func (_this *urlBuilder) BuildFromBool(ctx *Context, value bool, dst reflect.Value) reflect.Value {
	ctx.UnstackBuilder()
	return ctx.CurrentBuilder.BuildFromBool(ctx, value, dst)
}
func (_this *urlBuilder) BuildFromInt(ctx *Context, value int64, dst reflect.Value) reflect.Value {
	ctx.UnstackBuilder()
	return ctx.CurrentBuilder.BuildFromInt(ctx, value, dst)
}
func (_this *urlBuilder) BuildFromUint(ctx *Context, value uint64, dst reflect.Value) reflect.Value {
	ctx.UnstackBuilder()
	return ctx.CurrentBuilder.BuildFromUint(ctx, value, dst)
}
func (_this *urlBuilder) BuildFromBigInt(ctx *Context, value *big.Int, dst reflect.Value) reflect.Value {
	ctx.UnstackBuilder()
	return ctx.CurrentBuilder.BuildFromBigInt(ctx, value, dst)
}
func (_this *urlBuilder) BuildFromFloat(ctx *Context, value float64, dst reflect.Value) reflect.Value {
	ctx.UnstackBuilder()
	return ctx.CurrentBuilder.BuildFromFloat(ctx, value, dst)
}
func (_this *urlBuilder) BuildFromBigFloat(ctx *Context, value *big.Float, dst reflect.Value) reflect.Value {
	ctx.UnstackBuilder()
	return ctx.CurrentBuilder.BuildFromBigFloat(ctx, value, dst)
}
func (_this *urlBuilder) BuildFromDecimalFloat(ctx *Context, value compact_float.DFloat, dst reflect.Value) reflect.Value {
	ctx.UnstackBuilder()
	return ctx.CurrentBuilder.BuildFromDecimalFloat(ctx, value, dst)
}
func (_this *urlBuilder) BuildFromBigDecimalFloat(ctx *Context, value *apd.Decimal, dst reflect.Value) reflect.Value {
	ctx.UnstackBuilder()
	return ctx.CurrentBuilder.BuildFromBigDecimalFloat(ctx, value, dst)
}
func (_this *urlBuilder) BuildFromUUID(ctx *Context, value []byte, dst reflect.Value) reflect.Value {
	ctx.UnstackBuilder()
	return ctx.CurrentBuilder.BuildFromUUID(ctx, value, dst)
}
func (_this *urlBuilder) BuildFromTime(ctx *Context, value time.Time, dst reflect.Value) reflect.Value {
	ctx.UnstackBuilder()
	return ctx.CurrentBuilder.BuildFromTime(ctx, value, dst)
}
func (_this *urlBuilder) BuildFromCompactTime(ctx *Context, value compact_time.Time, dst reflect.Value) reflect.Value {
	ctx.UnstackBuilder()
	return ctx.CurrentBuilder.BuildFromCompactTime(ctx, value, dst)
}
func (_this *urlBuilder) BuildInitiateList(ctx *Context) {
	ctx.UnstackBuilder()
	ctx.CurrentBuilder.BuildInitiateList(ctx)
}
func (_this *urlBuilder) BuildInitiateMap(ctx *Context) {
	ctx.UnstackBuilder()
	ctx.CurrentBuilder.BuildInitiateMap(ctx)
}
func (_this *urlBuilder) BuildEndContainer(ctx *Context) {
	ctx.UnstackBuilder()
	ctx.CurrentBuilder.BuildEndContainer(ctx)
}
func (_this *urlBuilder) BuildFromReference(ctx *Context, id interface{}) {
	ctx.UnstackBuilder()
	ctx.CurrentBuilder.BuildFromReference(ctx, id)
}

// ============================================================================

type pUrlBuilder struct{}

var globalPUrlBuilder = &pUrlBuilder{}

func newPUrlBuilderGenerator() BuilderGenerator {
	return func(ctx *Context) ObjectBuilder {
		builder := globalPUrlBuilder
		ctx.StackBuilder(builder)
		return builder
	}
}

func (_this *pUrlBuilder) String() string { return reflect.TypeOf(_this).String() }

func (_this *pUrlBuilder) BuildFromNil(ctx *Context, dst reflect.Value) reflect.Value {
	dst.Set(reflect.Zero(dst.Type()))
	ctx.UnstackBuilder()
	ctx.NANext()
	return dst
}

func (_this *pUrlBuilder) BuildFromArray(ctx *Context, arrayType events.ArrayType, value []byte, dst reflect.Value) reflect.Value {
	switch arrayType {
	case events.ArrayTypeResourceID:
		setPRIDFromString(string(value), dst)
	default:
		ctx.UnstackBuilder()
		return ctx.CurrentBuilder.BuildFromArray(ctx, arrayType, value, dst)
	}
	return dst
}

func (_this *pUrlBuilder) BuildFromBool(ctx *Context, value bool, dst reflect.Value) reflect.Value {
	ctx.UnstackBuilder()
	return ctx.CurrentBuilder.BuildFromBool(ctx, value, dst)
}
func (_this *pUrlBuilder) BuildFromInt(ctx *Context, value int64, dst reflect.Value) reflect.Value {
	ctx.UnstackBuilder()
	return ctx.CurrentBuilder.BuildFromInt(ctx, value, dst)
}
func (_this *pUrlBuilder) BuildFromUint(ctx *Context, value uint64, dst reflect.Value) reflect.Value {
	ctx.UnstackBuilder()
	return ctx.CurrentBuilder.BuildFromUint(ctx, value, dst)
}
func (_this *pUrlBuilder) BuildFromBigInt(ctx *Context, value *big.Int, dst reflect.Value) reflect.Value {
	ctx.UnstackBuilder()
	return ctx.CurrentBuilder.BuildFromBigInt(ctx, value, dst)
}
func (_this *pUrlBuilder) BuildFromFloat(ctx *Context, value float64, dst reflect.Value) reflect.Value {
	ctx.UnstackBuilder()
	return ctx.CurrentBuilder.BuildFromFloat(ctx, value, dst)
}
func (_this *pUrlBuilder) BuildFromBigFloat(ctx *Context, value *big.Float, dst reflect.Value) reflect.Value {
	ctx.UnstackBuilder()
	return ctx.CurrentBuilder.BuildFromBigFloat(ctx, value, dst)
}
func (_this *pUrlBuilder) BuildFromDecimalFloat(ctx *Context, value compact_float.DFloat, dst reflect.Value) reflect.Value {
	ctx.UnstackBuilder()
	return ctx.CurrentBuilder.BuildFromDecimalFloat(ctx, value, dst)
}
func (_this *pUrlBuilder) BuildFromBigDecimalFloat(ctx *Context, value *apd.Decimal, dst reflect.Value) reflect.Value {
	ctx.UnstackBuilder()
	return ctx.CurrentBuilder.BuildFromBigDecimalFloat(ctx, value, dst)
}
func (_this *pUrlBuilder) BuildFromUUID(ctx *Context, value []byte, dst reflect.Value) reflect.Value {
	ctx.UnstackBuilder()
	return ctx.CurrentBuilder.BuildFromUUID(ctx, value, dst)
}
func (_this *pUrlBuilder) BuildFromTime(ctx *Context, value time.Time, dst reflect.Value) reflect.Value {
	ctx.UnstackBuilder()
	return ctx.CurrentBuilder.BuildFromTime(ctx, value, dst)
}
func (_this *pUrlBuilder) BuildFromCompactTime(ctx *Context, value compact_time.Time, dst reflect.Value) reflect.Value {
	ctx.UnstackBuilder()
	return ctx.CurrentBuilder.BuildFromCompactTime(ctx, value, dst)
}
func (_this *pUrlBuilder) BuildInitiateList(ctx *Context) {
	ctx.UnstackBuilder()
	ctx.CurrentBuilder.BuildInitiateList(ctx)
}
func (_this *pUrlBuilder) BuildInitiateMap(ctx *Context) {
	ctx.UnstackBuilder()
	ctx.CurrentBuilder.BuildInitiateMap(ctx)
}
func (_this *pUrlBuilder) BuildEndContainer(ctx *Context) {
	ctx.UnstackBuilder()
	ctx.CurrentBuilder.BuildEndContainer(ctx)
}
func (_this *pUrlBuilder) BuildFromReference(ctx *Context, id interface{}) {
	ctx.UnstackBuilder()
	ctx.CurrentBuilder.BuildFromReference(ctx, id)
}
