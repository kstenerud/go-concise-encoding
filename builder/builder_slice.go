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

const defaultSliceCap = 4

type sliceBuilder struct {
	dstType       reflect.Type
	elemGenerator BuilderGenerator
	ppContainer   **reflect.Value
}

func newSliceBuilderGenerator(getBuilderGeneratorForType BuilderGeneratorGetter, dstType reflect.Type) BuilderGenerator {
	builderGenerator := getBuilderGeneratorForType(dstType.Elem())

	return func(ctx *Context) Builder {
		container := reflect.MakeSlice(dstType, 0, defaultSliceCap)
		ppContainer := new(*reflect.Value)
		*ppContainer = &container

		return &sliceBuilder{
			dstType:       dstType,
			elemGenerator: builderGenerator,
			ppContainer:   ppContainer,
		}
	}
}

func (_this *sliceBuilder) String() string {
	return fmt.Sprintf("%v<%v>", reflect.TypeOf(_this), _this.elemGenerator)
}

func (_this *sliceBuilder) newElem() reflect.Value {
	return reflect.New(_this.dstType.Elem()).Elem()
}

func (_this *sliceBuilder) storeValue(value reflect.Value) {
	**_this.ppContainer = reflect.Append(**_this.ppContainer, value)
}

func (_this *sliceBuilder) BuildFromNull(ctx *Context, _ reflect.Value) reflect.Value {
	object := _this.newElem()
	_this.elemGenerator(ctx).BuildFromNull(ctx, object)
	if _this.ppContainer != nil {
		_this.storeValue(object)
	}
	return object
}

func (_this *sliceBuilder) BuildFromBool(ctx *Context, value bool, _ reflect.Value) reflect.Value {
	object := _this.newElem()
	_this.elemGenerator(ctx).BuildFromBool(ctx, value, object)
	_this.storeValue(object)
	return object
}

func (_this *sliceBuilder) BuildFromInt(ctx *Context, value int64, _ reflect.Value) reflect.Value {
	object := _this.newElem()
	_this.elemGenerator(ctx).BuildFromInt(ctx, value, object)
	_this.storeValue(object)
	return object
}

func (_this *sliceBuilder) BuildFromUint(ctx *Context, value uint64, _ reflect.Value) reflect.Value {
	object := _this.newElem()
	_this.elemGenerator(ctx).BuildFromUint(ctx, value, object)
	_this.storeValue(object)
	return object
}

func (_this *sliceBuilder) BuildFromBigInt(ctx *Context, value *big.Int, _ reflect.Value) reflect.Value {
	object := _this.newElem()
	_this.elemGenerator(ctx).BuildFromBigInt(ctx, value, object)
	_this.storeValue(object)
	return object
}

func (_this *sliceBuilder) BuildFromFloat(ctx *Context, value float64, _ reflect.Value) reflect.Value {
	object := _this.newElem()
	_this.elemGenerator(ctx).BuildFromFloat(ctx, value, object)
	_this.storeValue(object)
	return object
}

func (_this *sliceBuilder) BuildFromBigFloat(ctx *Context, value *big.Float, _ reflect.Value) reflect.Value {
	object := _this.newElem()
	_this.elemGenerator(ctx).BuildFromBigFloat(ctx, value, object)
	_this.storeValue(object)
	return object
}

func (_this *sliceBuilder) BuildFromDecimalFloat(ctx *Context, value compact_float.DFloat, _ reflect.Value) reflect.Value {
	object := _this.newElem()
	_this.elemGenerator(ctx).BuildFromDecimalFloat(ctx, value, object)
	_this.storeValue(object)
	return object
}

func (_this *sliceBuilder) BuildFromBigDecimalFloat(ctx *Context, value *apd.Decimal, _ reflect.Value) reflect.Value {
	object := _this.newElem()
	_this.elemGenerator(ctx).BuildFromBigDecimalFloat(ctx, value, object)
	_this.storeValue(object)
	return object
}

func (_this *sliceBuilder) BuildFromUID(ctx *Context, value []byte, _ reflect.Value) reflect.Value {
	object := _this.newElem()
	_this.elemGenerator(ctx).BuildFromUID(ctx, value, object)
	_this.storeValue(object)
	return object
}

func (_this *sliceBuilder) BuildFromArray(ctx *Context, arrayType events.ArrayType, value []byte, _ reflect.Value) reflect.Value {
	object := _this.newElem()
	_this.elemGenerator(ctx).BuildFromArray(ctx, arrayType, value, object)
	_this.storeValue(object)
	return object
}

func (_this *sliceBuilder) BuildFromStringlikeArray(ctx *Context, arrayType events.ArrayType, value string, _ reflect.Value) reflect.Value {
	object := _this.newElem()
	_this.elemGenerator(ctx).BuildFromStringlikeArray(ctx, arrayType, value, object)
	_this.storeValue(object)
	return object
}

func (_this *sliceBuilder) BuildFromMedia(ctx *Context, mediaType string, data []byte, _ reflect.Value) reflect.Value {
	object := _this.newElem()
	_this.elemGenerator(ctx).BuildFromMedia(ctx, mediaType, data, object)
	_this.storeValue(object)
	return object
}

func (_this *sliceBuilder) BuildFromTime(ctx *Context, value time.Time, _ reflect.Value) reflect.Value {
	object := _this.newElem()
	_this.elemGenerator(ctx).BuildFromTime(ctx, value, object)
	_this.storeValue(object)
	return object
}

func (_this *sliceBuilder) BuildFromCompactTime(ctx *Context, value compact_time.Time, _ reflect.Value) reflect.Value {
	object := _this.newElem()
	_this.elemGenerator(ctx).BuildFromCompactTime(ctx, value, object)
	_this.storeValue(object)
	return object
}

func (_this *sliceBuilder) BuildNewList(ctx *Context) {
	_this.elemGenerator(ctx).BuildBeginListContents(ctx)
}

func (_this *sliceBuilder) BuildNewMap(ctx *Context) {
	_this.elemGenerator(ctx).BuildBeginMapContents(ctx)
}

func (_this *sliceBuilder) BuildNewMarkup(ctx *Context, name []byte) {
	_this.elemGenerator(ctx).BuildBeginMarkupContents(ctx, name)
}

func (_this *sliceBuilder) BuildNewNode(ctx *Context) {
	_this.elemGenerator(ctx).BuildBeginNodeContents(ctx)
}

func (_this *sliceBuilder) BuildNewEdge(ctx *Context) {
	_this.elemGenerator(ctx).BuildBeginEdgeContents(ctx)
}

func (_this *sliceBuilder) BuildEndContainer(ctx *Context) {
	object := **_this.ppContainer
	ctx.UnstackBuilderAndNotifyChildFinished(object)
}

func (_this *sliceBuilder) BuildBeginListContents(ctx *Context) {
	ctx.StackBuilder(_this)
}

func (_this *sliceBuilder) BuildFromReference(ctx *Context, id []byte) {
	ppContainer := _this.ppContainer
	index := (**ppContainer).Len()
	elem := _this.newElem()
	_this.storeValue(elem)
	ctx.NotifyReference(id, func(object reflect.Value) {
		setAnythingFromAnything(object, (**ppContainer).Index(index))
	})
}

func (_this *sliceBuilder) NotifyChildContainerFinished(ctx *Context, value reflect.Value) {
	_this.storeValue(value)
}
