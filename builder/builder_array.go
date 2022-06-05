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

type arrayBuilder struct {
	containerType reflect.Type
	elemGenerator BuilderGenerator
	container     reflect.Value
	elemIndex     int
}

func newArrayBuilderGenerator(getBuilderGeneratorForType BuilderGeneratorGetter, containerType reflect.Type) BuilderGenerator {
	elemBuilderGenerator := getBuilderGeneratorForType(containerType.Elem())
	return func(ctx *Context) Builder {
		return &arrayBuilder{
			containerType: containerType,
			elemGenerator: elemBuilderGenerator,
			container:     reflect.New(containerType).Elem(),
			elemIndex:     0,
		}
	}
}

func (_this *arrayBuilder) String() string {
	return fmt.Sprintf("%v<%v>", reflect.TypeOf(_this), _this.containerType.Elem())
}

func (_this *arrayBuilder) advanceElem() reflect.Value {
	elem := _this.container.Index(_this.elemIndex)
	_this.elemIndex++
	return elem
}

func (_this *arrayBuilder) BuildFromNull(ctx *Context, _ reflect.Value) reflect.Value {
	object := _this.advanceElem()
	_this.elemGenerator(ctx).BuildFromNull(ctx, object)
	return object
}

func (_this *arrayBuilder) BuildFromBool(ctx *Context, value bool, _ reflect.Value) reflect.Value {
	object := _this.advanceElem()
	_this.elemGenerator(ctx).BuildFromBool(ctx, value, object)
	return object
}

func (_this *arrayBuilder) BuildFromInt(ctx *Context, value int64, _ reflect.Value) reflect.Value {
	object := _this.advanceElem()
	_this.elemGenerator(ctx).BuildFromInt(ctx, value, object)
	return object
}

func (_this *arrayBuilder) BuildFromUint(ctx *Context, value uint64, _ reflect.Value) reflect.Value {
	object := _this.advanceElem()
	_this.elemGenerator(ctx).BuildFromUint(ctx, value, object)
	return object
}

func (_this *arrayBuilder) BuildFromBigInt(ctx *Context, value *big.Int, _ reflect.Value) reflect.Value {
	object := _this.advanceElem()
	_this.elemGenerator(ctx).BuildFromBigInt(ctx, value, object)
	return object
}

func (_this *arrayBuilder) BuildFromFloat(ctx *Context, value float64, _ reflect.Value) reflect.Value {
	object := _this.advanceElem()
	_this.elemGenerator(ctx).BuildFromFloat(ctx, value, object)
	return object
}

func (_this *arrayBuilder) BuildFromBigFloat(ctx *Context, value *big.Float, _ reflect.Value) reflect.Value {
	object := _this.advanceElem()
	_this.elemGenerator(ctx).BuildFromBigFloat(ctx, value, object)
	return object
}

func (_this *arrayBuilder) BuildFromDecimalFloat(ctx *Context, value compact_float.DFloat, _ reflect.Value) reflect.Value {
	object := _this.advanceElem()
	_this.elemGenerator(ctx).BuildFromDecimalFloat(ctx, value, object)
	return object
}

func (_this *arrayBuilder) BuildFromBigDecimalFloat(ctx *Context, value *apd.Decimal, _ reflect.Value) reflect.Value {
	object := _this.advanceElem()
	_this.elemGenerator(ctx).BuildFromBigDecimalFloat(ctx, value, object)
	return object
}

func (_this *arrayBuilder) BuildFromUID(ctx *Context, value []byte, _ reflect.Value) reflect.Value {
	object := _this.advanceElem()
	_this.elemGenerator(ctx).BuildFromUID(ctx, value, object)
	return object
}

func (_this *arrayBuilder) BuildFromArray(ctx *Context, arrayType events.ArrayType, value []byte, _ reflect.Value) reflect.Value {
	object := _this.advanceElem()
	_this.elemGenerator(ctx).BuildFromArray(ctx, arrayType, value, object)
	return object
}

func (_this *arrayBuilder) BuildFromStringlikeArray(ctx *Context, arrayType events.ArrayType, value string, _ reflect.Value) reflect.Value {
	object := _this.advanceElem()
	_this.elemGenerator(ctx).BuildFromStringlikeArray(ctx, arrayType, value, object)
	return object
}

func (_this *arrayBuilder) BuildFromMedia(ctx *Context, mediaType string, data []byte, _ reflect.Value) reflect.Value {
	object := _this.advanceElem()
	_this.elemGenerator(ctx).BuildFromMedia(ctx, mediaType, data, object)
	return object
}

func (_this *arrayBuilder) BuildFromTime(ctx *Context, value time.Time, _ reflect.Value) reflect.Value {
	object := _this.advanceElem()
	_this.elemGenerator(ctx).BuildFromTime(ctx, value, object)
	return object
}

func (_this *arrayBuilder) BuildFromCompactTime(ctx *Context, value compact_time.Time, _ reflect.Value) reflect.Value {
	object := _this.advanceElem()
	_this.elemGenerator(ctx).BuildFromCompactTime(ctx, value, object)
	return object
}

func (_this *arrayBuilder) BuildNewList(ctx *Context) {
	_this.elemGenerator(ctx).BuildBeginListContents(ctx)
}

func (_this *arrayBuilder) BuildNewMap(ctx *Context) {
	_this.elemGenerator(ctx).BuildBeginMapContents(ctx)
}

func (_this *arrayBuilder) BuildNewNode(ctx *Context) {
	_this.elemGenerator(ctx).BuildBeginNodeContents(ctx)
}

func (_this *arrayBuilder) BuildNewEdge(ctx *Context) {
	_this.elemGenerator(ctx).BuildBeginEdgeContents(ctx)
}

func (_this *arrayBuilder) BuildEndContainer(ctx *Context) {
	object := _this.container
	ctx.UnstackBuilderAndNotifyChildFinished(object)
}

func (_this *arrayBuilder) BuildBeginListContents(ctx *Context) {
	ctx.StackBuilder(_this)
}

func (_this *arrayBuilder) BuildFromReference(ctx *Context, id []byte) {
	container := _this.container
	index := _this.elemIndex
	_this.elemIndex++
	ctx.NotifyReference(id, func(object reflect.Value) {
		setAnythingFromAnything(object, container.Index(index))
	})
}

func (_this *arrayBuilder) NotifyChildContainerFinished(ctx *Context, value reflect.Value) {
	_this.advanceElem().Set(value)
}
