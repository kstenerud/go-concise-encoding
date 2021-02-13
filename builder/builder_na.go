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

type naBuilder struct{}

var globalNABuilder = &naBuilder{}

func generateNABuilder(ctx *Context) Builder { return globalNABuilder }
func (_this *naBuilder) String() string            { return reflect.TypeOf(_this).String() }

func (_this *naBuilder) BuildFromNil(ctx *Context, dst reflect.Value) reflect.Value {
	ctx.UnstackBuilder()
	return dst
}

func (_this *naBuilder) BuildFromBool(ctx *Context, _ bool, dst reflect.Value) reflect.Value {
	ctx.UnstackBuilder()
	return dst
}

func (_this *naBuilder) BuildFromInt(ctx *Context, _ int64, dst reflect.Value) reflect.Value {
	ctx.UnstackBuilder()
	return dst
}

func (_this *naBuilder) BuildFromUint(ctx *Context, _ uint64, dst reflect.Value) reflect.Value {
	ctx.UnstackBuilder()
	return dst
}

func (_this *naBuilder) BuildFromBigInt(ctx *Context, _ *big.Int, dst reflect.Value) reflect.Value {
	ctx.UnstackBuilder()
	return dst
}

func (_this *naBuilder) BuildFromFloat(ctx *Context, _ float64, dst reflect.Value) reflect.Value {
	ctx.UnstackBuilder()
	return dst
}

func (_this *naBuilder) BuildFromBigFloat(ctx *Context, _ *big.Float, dst reflect.Value) reflect.Value {
	ctx.UnstackBuilder()
	return dst
}

func (_this *naBuilder) BuildFromDecimalFloat(ctx *Context, _ compact_float.DFloat, dst reflect.Value) reflect.Value {
	ctx.UnstackBuilder()
	return dst
}

func (_this *naBuilder) BuildFromBigDecimalFloat(ctx *Context, _ *apd.Decimal, dst reflect.Value) reflect.Value {
	ctx.UnstackBuilder()
	return dst
}

func (_this *naBuilder) BuildFromUUID(ctx *Context, _ []byte, dst reflect.Value) reflect.Value {
	ctx.UnstackBuilder()
	return dst
}

func (_this *naBuilder) BuildFromArray(ctx *Context, _ events.ArrayType, _ []byte, dst reflect.Value) reflect.Value {
	ctx.UnstackBuilder()
	return dst
}

func (_this *naBuilder) BuildFromStringlikeArray(ctx *Context, _ events.ArrayType, _ string, dst reflect.Value) reflect.Value {
	ctx.UnstackBuilder()
	return dst
}

func (_this *naBuilder) BuildFromTime(ctx *Context, _ time.Time, dst reflect.Value) reflect.Value {
	ctx.UnstackBuilder()
	return dst
}

func (_this *naBuilder) BuildFromCompactTime(ctx *Context, _ compact_time.Time, dst reflect.Value) reflect.Value {
	ctx.UnstackBuilder()
	return dst
}

func (_this *naBuilder) BuildInitiateList(ctx *Context) {
	ctx.StackBuilder(_this)
}

func (_this *naBuilder) BuildInitiateMap(ctx *Context) {
	ctx.StackBuilder(_this)
}

func (_this *naBuilder) BuildEndContainer(ctx *Context) {
	ctx.UnstackBuilder()
}

func (_this *naBuilder) BuildBeginListContents(ctx *Context) {
	ctx.StackBuilder(_this)
}

func (_this *naBuilder) BuildBeginMapContents(ctx *Context) {
	ctx.StackBuilder(_this)
}

func (_this *naBuilder) BuildFromReference(ctx *Context, _ interface{}) {
	ctx.UnstackBuilder()
}

func (_this *naBuilder) NotifyChildContainerFinished(ctx *Context, _ reflect.Value) {
	ctx.UnstackBuilder()
}

func (_this *naBuilder) BuildConcatenate(ctx *Context) {
	panic("TODO: naBuilder.BuildConcatenate")
}
