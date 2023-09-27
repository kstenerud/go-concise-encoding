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
	compact_float "github.com/kstenerud/go-compact-float"
	compact_time "github.com/kstenerud/go-compact-time"
	"github.com/kstenerud/go-concise-encoding/ce/events"
)

type recordBuilder struct {
	builder Builder
	keys    []recordTypeKey
	index   int
}

func generateRecordBuilder(ctx *Context, keys []recordTypeKey, builder Builder) Builder {
	return &recordBuilder{
		builder: builder,
		keys:    keys,
	}
}
func (_this *recordBuilder) String() string { return reflect.TypeOf(_this).String() }

func (_this *recordBuilder) sendKey(ctx *Context) {
	_this.keys[_this.index](ctx, _this.builder)
	_this.index++
}

func (_this *recordBuilder) BuildFromNull(ctx *Context, dst reflect.Value) reflect.Value {
	_this.sendKey(ctx)
	return _this.builder.BuildFromNull(ctx, dst)
}

func (_this *recordBuilder) BuildFromBool(ctx *Context, value bool, dst reflect.Value) reflect.Value {
	_this.sendKey(ctx)
	return _this.builder.BuildFromBool(ctx, value, dst)
}

func (_this *recordBuilder) BuildFromInt(ctx *Context, value int64, dst reflect.Value) reflect.Value {
	_this.sendKey(ctx)
	return _this.builder.BuildFromInt(ctx, value, dst)
}

func (_this *recordBuilder) BuildFromUint(ctx *Context, value uint64, dst reflect.Value) reflect.Value {
	_this.sendKey(ctx)
	return _this.builder.BuildFromUint(ctx, value, dst)
}

func (_this *recordBuilder) BuildFromBigInt(ctx *Context, value *big.Int, dst reflect.Value) reflect.Value {
	_this.sendKey(ctx)
	return _this.builder.BuildFromBigInt(ctx, value, dst)
}

func (_this *recordBuilder) BuildFromFloat(ctx *Context, value float64, dst reflect.Value) reflect.Value {
	_this.sendKey(ctx)
	return _this.builder.BuildFromFloat(ctx, value, dst)
}

func (_this *recordBuilder) BuildFromBigFloat(ctx *Context, value *big.Float, dst reflect.Value) reflect.Value {
	_this.sendKey(ctx)
	return _this.builder.BuildFromBigFloat(ctx, value, dst)
}

func (_this *recordBuilder) BuildFromDecimalFloat(ctx *Context, value compact_float.DFloat, dst reflect.Value) reflect.Value {
	_this.sendKey(ctx)
	return _this.builder.BuildFromDecimalFloat(ctx, value, dst)
}

func (_this *recordBuilder) BuildFromBigDecimalFloat(ctx *Context, value *apd.Decimal, dst reflect.Value) reflect.Value {
	_this.sendKey(ctx)
	return _this.builder.BuildFromBigDecimalFloat(ctx, value, dst)
}

func (_this *recordBuilder) BuildFromUID(ctx *Context, value []byte, dst reflect.Value) reflect.Value {
	_this.sendKey(ctx)
	return _this.builder.BuildFromUID(ctx, value, dst)
}

func (_this *recordBuilder) BuildFromArray(ctx *Context, arrayType events.ArrayType, value []byte, dst reflect.Value) reflect.Value {
	_this.sendKey(ctx)
	return _this.builder.BuildFromArray(ctx, arrayType, value, dst)
}

func (_this *recordBuilder) BuildFromStringlikeArray(ctx *Context, arrayType events.ArrayType, value string, dst reflect.Value) reflect.Value {
	_this.sendKey(ctx)
	return _this.builder.BuildFromStringlikeArray(ctx, arrayType, value, dst)
}

func (_this *recordBuilder) BuildFromCustomBinary(ctx *Context, customType uint64, value []byte, dst reflect.Value) reflect.Value {
	_this.sendKey(ctx)
	return _this.builder.BuildFromCustomBinary(ctx, customType, value, dst)
}

func (_this *recordBuilder) BuildFromCustomText(ctx *Context, customType uint64, value string, dst reflect.Value) reflect.Value {
	_this.sendKey(ctx)
	return _this.builder.BuildFromCustomText(ctx, customType, value, dst)
}

func (_this *recordBuilder) BuildFromMedia(ctx *Context, mediaType string, data []byte, dst reflect.Value) reflect.Value {
	_this.sendKey(ctx)
	return _this.builder.BuildFromMedia(ctx, mediaType, data, dst)
}

func (_this *recordBuilder) BuildFromTime(ctx *Context, value compact_time.Time, dst reflect.Value) reflect.Value {
	_this.sendKey(ctx)
	return _this.builder.BuildFromTime(ctx, value, dst)
}

func (_this *recordBuilder) BuildNewList(ctx *Context) {
	_this.sendKey(ctx)
	_this.builder.BuildNewList(ctx)
}

func (_this *recordBuilder) BuildNewMap(ctx *Context) {
	_this.sendKey(ctx)
	_this.builder.BuildNewMap(ctx)
}

func (_this *recordBuilder) BuildNewEdge(ctx *Context) {
	_this.sendKey(ctx)
	_this.builder.BuildNewEdge(ctx)
}

func (_this *recordBuilder) BuildNewNode(ctx *Context) {
	_this.sendKey(ctx)
	_this.builder.BuildNewNode(ctx)
}

func (_this *recordBuilder) BuildEndContainer(ctx *Context) {
	_this.builder.BuildEndContainer(ctx)
}

func (_this *recordBuilder) BuildArtificiallyEndContainer(ctx *Context) {
	_this.builder.BuildArtificiallyEndContainer(ctx)
}

func (_this *recordBuilder) BuildFromLocalReference(ctx *Context, id []byte) {
	_this.sendKey(ctx)
	_this.builder.BuildFromLocalReference(ctx, id)
}

func (_this *recordBuilder) NotifyChildContainerFinished(ctx *Context, value reflect.Value) {
	_this.builder.NotifyChildContainerFinished(ctx, value)
}
