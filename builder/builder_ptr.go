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

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
)

type ptrBuilder struct {
	dstType       reflect.Type
	elemGenerator BuilderGenerator
}

func newPtrBuilderGenerator(getBuilderGeneratorForType BuilderGeneratorGetter, dstType reflect.Type) BuilderGenerator {
	builderGenerator := getBuilderGeneratorForType(dstType.Elem())

	return func(ctx *Context) Builder {
		return &ptrBuilder{
			dstType:       dstType,
			elemGenerator: builderGenerator,
		}
	}
}

func (_this *ptrBuilder) String() string {
	return fmt.Sprintf("%v<%v>", reflect.TypeOf(_this), _this.elemGenerator(nil))
}

func (_this *ptrBuilder) newElem() reflect.Value {
	return reflect.New(_this.dstType.Elem())
}

func (_this *ptrBuilder) BuildFromNil(ctx *Context, dst reflect.Value) reflect.Value {
	dst.Set(reflect.Zero(_this.dstType))
	return dst
}

func (_this *ptrBuilder) BuildFromBool(ctx *Context, value bool, dst reflect.Value) reflect.Value {
	ptr := _this.newElem()
	_this.elemGenerator(ctx).BuildFromBool(ctx, value, ptr.Elem())
	dst.Set(ptr)
	return dst
}

func (_this *ptrBuilder) BuildFromInt(ctx *Context, value int64, dst reflect.Value) reflect.Value {
	ptr := _this.newElem()
	_this.elemGenerator(ctx).BuildFromInt(ctx, value, ptr.Elem())
	dst.Set(ptr)
	return dst
}

func (_this *ptrBuilder) BuildFromUint(ctx *Context, value uint64, dst reflect.Value) reflect.Value {
	ptr := _this.newElem()
	_this.elemGenerator(ctx).BuildFromUint(ctx, value, ptr.Elem())
	dst.Set(ptr)
	return dst
}

func (_this *ptrBuilder) BuildFromBigInt(ctx *Context, value *big.Int, dst reflect.Value) reflect.Value {
	ptr := _this.newElem()
	_this.elemGenerator(ctx).BuildFromBigInt(ctx, value, ptr.Elem())
	dst.Set(ptr)
	return dst
}

func (_this *ptrBuilder) BuildFromFloat(ctx *Context, value float64, dst reflect.Value) reflect.Value {
	ptr := _this.newElem()
	_this.elemGenerator(ctx).BuildFromFloat(ctx, value, ptr.Elem())
	dst.Set(ptr)
	return dst
}

func (_this *ptrBuilder) BuildFromBigFloat(ctx *Context, value *big.Float, dst reflect.Value) reflect.Value {
	ptr := _this.newElem()
	_this.elemGenerator(ctx).BuildFromBigFloat(ctx, value, ptr.Elem())
	dst.Set(ptr)
	return dst
}

func (_this *ptrBuilder) BuildFromDecimalFloat(ctx *Context, value compact_float.DFloat, dst reflect.Value) reflect.Value {
	ptr := _this.newElem()
	_this.elemGenerator(ctx).BuildFromDecimalFloat(ctx, value, ptr.Elem())
	dst.Set(ptr)
	return dst
}

func (_this *ptrBuilder) BuildFromBigDecimalFloat(ctx *Context, value *apd.Decimal, dst reflect.Value) reflect.Value {
	ptr := _this.newElem()
	_this.elemGenerator(ctx).BuildFromBigDecimalFloat(ctx, value, ptr.Elem())
	dst.Set(ptr)
	return dst
}

func (_this *ptrBuilder) BuildFromUID(ctx *Context, value []byte, dst reflect.Value) reflect.Value {
	ptr := _this.newElem()
	_this.elemGenerator(ctx).BuildFromUID(ctx, value, ptr.Elem())
	dst.Set(ptr)
	return dst
}

func (_this *ptrBuilder) BuildFromArray(ctx *Context, arrayType events.ArrayType, value []byte, dst reflect.Value) reflect.Value {
	ptr := _this.newElem()
	_this.elemGenerator(ctx).BuildFromArray(ctx, arrayType, value, ptr.Elem())
	dst.Set(ptr)
	return dst
}

func (_this *ptrBuilder) BuildFromStringlikeArray(ctx *Context, arrayType events.ArrayType, value string, dst reflect.Value) reflect.Value {
	ptr := _this.newElem()
	_this.elemGenerator(ctx).BuildFromStringlikeArray(ctx, arrayType, value, ptr.Elem())
	dst.Set(ptr)
	return dst
}

func (_this *ptrBuilder) BuildFromMedia(ctx *Context, mediaType string, data []byte, dst reflect.Value) reflect.Value {
	ptr := _this.newElem()
	_this.elemGenerator(ctx).BuildFromMedia(ctx, mediaType, data, dst)
	dst.Set(ptr)
	return dst
}

func (_this *ptrBuilder) BuildFromTime(ctx *Context, value time.Time, dst reflect.Value) reflect.Value {
	ptr := _this.newElem()
	_this.elemGenerator(ctx).BuildFromTime(ctx, value, ptr.Elem())
	dst.Set(ptr)
	return dst
}

func (_this *ptrBuilder) BuildFromCompactTime(ctx *Context, value compact_time.Time, dst reflect.Value) reflect.Value {
	ptr := _this.newElem()
	_this.elemGenerator(ctx).BuildFromCompactTime(ctx, value, ptr.Elem())
	dst.Set(ptr)
	return dst
}

func (_this *ptrBuilder) BuildBeginListContents(ctx *Context) {
	ctx.StackBuilder(_this)
	_this.elemGenerator(ctx).BuildBeginListContents(ctx)
}

func (_this *ptrBuilder) BuildBeginMapContents(ctx *Context) {
	ctx.StackBuilder(_this)
	_this.elemGenerator(ctx).BuildBeginMapContents(ctx)
}

func (_this *ptrBuilder) BuildBeginMarkupContents(ctx *Context, name []byte) {
	ctx.StackBuilder(_this)
	_this.elemGenerator(ctx).BuildBeginMarkupContents(ctx, name)
}

func (_this *ptrBuilder) BuildBeginEdgeContents(ctx *Context) {
	ctx.StackBuilder(_this)
	_this.elemGenerator(ctx).BuildBeginEdgeContents(ctx)
}

func (_this *ptrBuilder) BuildBeginNodeContents(ctx *Context) {
	ctx.StackBuilder(_this)
	_this.elemGenerator(ctx).BuildBeginNodeContents(ctx)
}

func (_this *ptrBuilder) NotifyChildContainerFinished(ctx *Context, value reflect.Value) {
	ctx.UnstackBuilderAndNotifyChildFinished(value.Addr())
}
