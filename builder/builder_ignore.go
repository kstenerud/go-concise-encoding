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

	"github.com/cockroachdb/apd/v2"
	compact_float "github.com/kstenerud/go-compact-float"
	compact_time "github.com/kstenerud/go-compact-time"
	"github.com/kstenerud/go-concise-encoding/events"
)

type ignoreBuilder struct{}

var globalIgnoreBuilder = &ignoreBuilder{}

func (_this *ignoreBuilder) String() string { return reflect.TypeOf(_this).String() }

func (_this *ignoreBuilder) BuildFromNull(ctx *Context, dst reflect.Value) reflect.Value {
	ctx.UnstackBuilder()
	return dst

}

func (_this *ignoreBuilder) BuildFromBool(ctx *Context, _ bool, dst reflect.Value) reflect.Value {
	ctx.UnstackBuilder()
	return dst
}

func (_this *ignoreBuilder) BuildFromInt(ctx *Context, _ int64, dst reflect.Value) reflect.Value {
	ctx.UnstackBuilder()
	return dst
}

func (_this *ignoreBuilder) BuildFromUint(ctx *Context, _ uint64, dst reflect.Value) reflect.Value {
	ctx.UnstackBuilder()
	return dst
}

func (_this *ignoreBuilder) BuildFromBigInt(ctx *Context, _ *big.Int, dst reflect.Value) reflect.Value {
	ctx.UnstackBuilder()
	return dst
}

func (_this *ignoreBuilder) BuildFromFloat(ctx *Context, _ float64, dst reflect.Value) reflect.Value {
	ctx.UnstackBuilder()
	return dst
}

func (_this *ignoreBuilder) BuildFromBigFloat(ctx *Context, _ *big.Float, dst reflect.Value) reflect.Value {
	ctx.UnstackBuilder()
	return dst
}

func (_this *ignoreBuilder) BuildFromDecimalFloat(ctx *Context, _ compact_float.DFloat, dst reflect.Value) reflect.Value {
	ctx.UnstackBuilder()
	return dst
}

func (_this *ignoreBuilder) BuildFromBigDecimalFloat(ctx *Context, _ *apd.Decimal, dst reflect.Value) reflect.Value {
	ctx.UnstackBuilder()
	return dst
}

func (_this *ignoreBuilder) BuildFromUID(ctx *Context, _ []byte, dst reflect.Value) reflect.Value {
	ctx.UnstackBuilder()
	return dst
}

func (_this *ignoreBuilder) BuildFromArray(ctx *Context, _ events.ArrayType, _ []byte, dst reflect.Value) reflect.Value {
	ctx.UnstackBuilder()
	return dst
}

func (_this *ignoreBuilder) BuildFromStringlikeArray(ctx *Context, _ events.ArrayType, _ string, dst reflect.Value) reflect.Value {
	ctx.UnstackBuilder()
	return dst
}

func (_this *ignoreBuilder) BuildFromMedia(ctx *Context, mediaType string, data []byte, dst reflect.Value) reflect.Value {
	ctx.UnstackBuilder()
	return dst
}

func (_this *ignoreBuilder) BuildFromTime(ctx *Context, _ time.Time, dst reflect.Value) reflect.Value {
	ctx.UnstackBuilder()
	return dst
}

func (_this *ignoreBuilder) BuildFromCompactTime(ctx *Context, _ compact_time.Time, dst reflect.Value) reflect.Value {
	ctx.UnstackBuilder()
	return dst
}

func (_this *ignoreBuilder) BuildNewList(ctx *Context) {
	ctx.StackBuilder(globalIgnoreContainerBuilder)
}

func (_this *ignoreBuilder) BuildNewMap(ctx *Context) {
	ctx.StackBuilder(globalIgnoreContainerBuilder)
}

func (_this *ignoreBuilder) BuildNewEdge(ctx *Context) {
	ctx.StackBuilder(generateIgnoreEdgeBuilder(ctx))
}

func (_this *ignoreBuilder) BuildNewNode(ctx *Context) {
	ctx.StackBuilder(globalIgnoreContainerBuilder)
}

func (_this *ignoreBuilder) BuildFromReference(ctx *Context, _ []byte) {
	ctx.UnstackBuilder()
}

func (_this *ignoreBuilder) NotifyChildContainerFinished(ctx *Context, _ reflect.Value) {
	ctx.UnstackBuilder()
}

// IgnoreXTimesBuilder

type ignoreXTimesBuilder struct {
	maxIndex int
	index    int
}

func generateIgnoreXTimesBuilder(ctx *Context, maxIndex int) Builder {
	return &ignoreXTimesBuilder{
		maxIndex: maxIndex,
		index:    0,
	}
}

func generateIgnoreEdgeBuilder(ctx *Context) Builder {
	return generateIgnoreXTimesBuilder(ctx, 3)
}

func (_this *ignoreXTimesBuilder) String() string {
	return fmt.Sprintf("%v(%v)", reflect.TypeOf(_this), _this.maxIndex)
}

func (_this *ignoreXTimesBuilder) tryFinish(ctx *Context) {
	_this.index++
	if _this.index >= _this.maxIndex {
		ctx.UnstackBuilder()
		var obj = reflect.ValueOf(nil)
		ctx.CurrentBuilder.NotifyChildContainerFinished(ctx, obj)
	}
}

func (_this *ignoreXTimesBuilder) BuildFromNull(ctx *Context, dst reflect.Value) reflect.Value {
	_this.tryFinish(ctx)
	return dst
}

func (_this *ignoreXTimesBuilder) BuildFromBool(ctx *Context, value bool, dst reflect.Value) reflect.Value {
	_this.tryFinish(ctx)
	return dst
}

func (_this *ignoreXTimesBuilder) BuildFromInt(ctx *Context, value int64, dst reflect.Value) reflect.Value {
	_this.tryFinish(ctx)
	return dst
}

func (_this *ignoreXTimesBuilder) BuildFromUint(ctx *Context, value uint64, dst reflect.Value) reflect.Value {
	_this.tryFinish(ctx)
	return dst
}

func (_this *ignoreXTimesBuilder) BuildFromBigInt(ctx *Context, value *big.Int, dst reflect.Value) reflect.Value {
	_this.tryFinish(ctx)
	return dst
}

func (_this *ignoreXTimesBuilder) BuildFromFloat(ctx *Context, value float64, dst reflect.Value) reflect.Value {
	_this.tryFinish(ctx)
	return dst
}

func (_this *ignoreXTimesBuilder) BuildFromBigFloat(ctx *Context, value *big.Float, dst reflect.Value) reflect.Value {
	_this.tryFinish(ctx)
	return dst
}

func (_this *ignoreXTimesBuilder) BuildFromDecimalFloat(ctx *Context, value compact_float.DFloat, dst reflect.Value) reflect.Value {
	_this.tryFinish(ctx)
	return dst
}

func (_this *ignoreXTimesBuilder) BuildFromBigDecimalFloat(ctx *Context, value *apd.Decimal, dst reflect.Value) reflect.Value {
	_this.tryFinish(ctx)
	return dst
}

func (_this *ignoreXTimesBuilder) BuildFromUID(ctx *Context, value []byte, dst reflect.Value) reflect.Value {
	_this.tryFinish(ctx)
	return dst
}

func (_this *ignoreXTimesBuilder) BuildFromArray(ctx *Context, arrayType events.ArrayType, value []byte, dst reflect.Value) reflect.Value {
	_this.tryFinish(ctx)
	return dst
}

func (_this *ignoreXTimesBuilder) BuildFromStringlikeArray(ctx *Context, arrayType events.ArrayType, value string, dst reflect.Value) reflect.Value {
	_this.tryFinish(ctx)
	return dst
}

func (_this *ignoreXTimesBuilder) BuildFromMedia(ctx *Context, mediaType string, data []byte, dst reflect.Value) reflect.Value {
	_this.tryFinish(ctx)
	return dst
}

func (_this *ignoreXTimesBuilder) BuildFromTime(ctx *Context, value time.Time, dst reflect.Value) reflect.Value {
	_this.tryFinish(ctx)
	return dst
}

func (_this *ignoreXTimesBuilder) BuildFromCompactTime(ctx *Context, value compact_time.Time, dst reflect.Value) reflect.Value {
	_this.tryFinish(ctx)
	return dst
}

func (_this *ignoreXTimesBuilder) BuildNewList(ctx *Context) {
	ctx.StackBuilder(globalIgnoreContainerBuilder)
}

func (_this *ignoreXTimesBuilder) BuildNewMap(ctx *Context) {
	ctx.StackBuilder(globalIgnoreContainerBuilder)
}

func (_this *ignoreXTimesBuilder) BuildNewEdge(ctx *Context) {
	ctx.StackBuilder(globalIgnoreContainerBuilder)
}

func (_this *ignoreXTimesBuilder) BuildNewNode(ctx *Context) {
	ctx.StackBuilder(globalIgnoreContainerBuilder)
}

func (_this *ignoreXTimesBuilder) BuildFromReference(ctx *Context, id []byte) {
	_this.tryFinish(ctx)
}

func (_this *ignoreXTimesBuilder) NotifyChildContainerFinished(ctx *Context, value reflect.Value) {
	_this.tryFinish(ctx)
}

// Ignore container

type ignoreContainerBuilder struct{}

var globalIgnoreContainerBuilder = &ignoreContainerBuilder{}

func (_this *ignoreContainerBuilder) String() string { return reflect.TypeOf(_this).String() }

func (_this *ignoreContainerBuilder) BuildFromNull(ctx *Context, dst reflect.Value) reflect.Value {
	return dst

}

func (_this *ignoreContainerBuilder) BuildFromBool(ctx *Context, _ bool, dst reflect.Value) reflect.Value {
	return dst
}

func (_this *ignoreContainerBuilder) BuildFromInt(ctx *Context, _ int64, dst reflect.Value) reflect.Value {
	return dst
}

func (_this *ignoreContainerBuilder) BuildFromUint(ctx *Context, _ uint64, dst reflect.Value) reflect.Value {
	return dst
}

func (_this *ignoreContainerBuilder) BuildFromBigInt(ctx *Context, _ *big.Int, dst reflect.Value) reflect.Value {
	return dst
}

func (_this *ignoreContainerBuilder) BuildFromFloat(ctx *Context, _ float64, dst reflect.Value) reflect.Value {
	return dst
}

func (_this *ignoreContainerBuilder) BuildFromBigFloat(ctx *Context, _ *big.Float, dst reflect.Value) reflect.Value {
	return dst
}

func (_this *ignoreContainerBuilder) BuildFromDecimalFloat(ctx *Context, _ compact_float.DFloat, dst reflect.Value) reflect.Value {
	return dst
}

func (_this *ignoreContainerBuilder) BuildFromBigDecimalFloat(ctx *Context, _ *apd.Decimal, dst reflect.Value) reflect.Value {
	return dst
}

func (_this *ignoreContainerBuilder) BuildFromUID(ctx *Context, _ []byte, dst reflect.Value) reflect.Value {
	return dst
}

func (_this *ignoreContainerBuilder) BuildFromArray(ctx *Context, _ events.ArrayType, _ []byte, dst reflect.Value) reflect.Value {
	return dst
}

func (_this *ignoreContainerBuilder) BuildFromStringlikeArray(ctx *Context, _ events.ArrayType, _ string, dst reflect.Value) reflect.Value {
	return dst
}

func (_this *ignoreContainerBuilder) BuildFromMedia(ctx *Context, mediaType string, data []byte, dst reflect.Value) reflect.Value {
	return dst
}

func (_this *ignoreContainerBuilder) BuildFromTime(ctx *Context, _ time.Time, dst reflect.Value) reflect.Value {
	return dst
}

func (_this *ignoreContainerBuilder) BuildFromCompactTime(ctx *Context, _ compact_time.Time, dst reflect.Value) reflect.Value {
	return dst
}

func (_this *ignoreContainerBuilder) BuildNewList(ctx *Context) {
	ctx.StackBuilder(_this)
}

func (_this *ignoreContainerBuilder) BuildNewMap(ctx *Context) {
	ctx.StackBuilder(_this)
}

func (_this *ignoreContainerBuilder) BuildNewEdge(ctx *Context) {
	ctx.StackBuilder(generateIgnoreEdgeBuilder(ctx))
}

func (_this *ignoreContainerBuilder) BuildNewNode(ctx *Context) {
	ctx.StackBuilder(_this)
}

func (_this *ignoreContainerBuilder) BuildEndContainer(ctx *Context) {
	ctx.UnstackBuilderAndNotifyChildFinished(reflect.Value{})
}

func (_this *ignoreContainerBuilder) BuildBeginListContents(ctx *Context) {
}

func (_this *ignoreContainerBuilder) BuildBeginMapContents(ctx *Context) {
}

func (_this *ignoreContainerBuilder) BuildBeginEdgeContents(ctx *Context) {
}

func (_this *ignoreContainerBuilder) BuildBeginNodeContents(ctx *Context) {
}

func (_this *ignoreContainerBuilder) BuildFromReference(ctx *Context, _ []byte) {
}

func (_this *ignoreContainerBuilder) NotifyChildContainerFinished(ctx *Context, _ reflect.Value) {
}
