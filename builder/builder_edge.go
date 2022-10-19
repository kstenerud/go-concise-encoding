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

	"github.com/cockroachdb/apd/v2"
	compact_float "github.com/kstenerud/go-compact-float"
	compact_time "github.com/kstenerud/go-compact-time"
	"github.com/kstenerud/go-concise-encoding/ce/events"
	"github.com/kstenerud/go-concise-encoding/internal/common"
	"github.com/kstenerud/go-concise-encoding/types"
)

type edgeBuilder struct {
	components []reflect.Value
	index      int
}

func generateEdgeBuilder(ctx *Context) Builder {
	return &edgeBuilder{
		components: []reflect.Value{
			reflect.New(common.TypeInterface).Elem(),
			reflect.New(common.TypeInterface).Elem(),
			reflect.New(common.TypeInterface).Elem(),
		},
		index: 0,
	}
}

func (_this *edgeBuilder) String() string { return fmt.Sprintf("%v", reflect.TypeOf(_this)) }

func (_this *edgeBuilder) tryFinish(ctx *Context) {
	const maxIndex = 3
	_this.index++
	if _this.index >= maxIndex {
		obj := reflect.New(reflect.TypeOf(types.Edge{})).Elem()
		obj.Field(types.EdgeFieldIndexSource).Set(_this.components[0])
		obj.Field(types.EdgeFieldIndexDescription).Set(_this.components[1])
		obj.Field(types.EdgeFieldIndexDestination).Set(_this.components[2])
		ctx.UnstackBuilder()
		ctx.CurrentBuilder.NotifyChildContainerFinished(ctx, obj)
	}
}

func (_this *edgeBuilder) BuildFromNull(ctx *Context, _ reflect.Value) reflect.Value {
	result := globalInterfaceBuilder.BuildFromNull(ctx, _this.components[_this.index])
	_this.tryFinish(ctx)
	return result
}

func (_this *edgeBuilder) BuildFromBool(ctx *Context, value bool, _ reflect.Value) reflect.Value {
	result := globalInterfaceBuilder.BuildFromBool(ctx, value, _this.components[_this.index])
	_this.tryFinish(ctx)
	return result
}

func (_this *edgeBuilder) BuildFromInt(ctx *Context, value int64, _ reflect.Value) reflect.Value {
	result := globalInterfaceBuilder.BuildFromInt(ctx, value, _this.components[_this.index])
	_this.tryFinish(ctx)
	return result
}

func (_this *edgeBuilder) BuildFromUint(ctx *Context, value uint64, _ reflect.Value) reflect.Value {
	result := globalInterfaceBuilder.BuildFromUint(ctx, value, _this.components[_this.index])
	_this.tryFinish(ctx)
	return result
}

func (_this *edgeBuilder) BuildFromBigInt(ctx *Context, value *big.Int, _ reflect.Value) reflect.Value {
	result := globalInterfaceBuilder.BuildFromBigInt(ctx, value, _this.components[_this.index])
	_this.tryFinish(ctx)
	return result
}

func (_this *edgeBuilder) BuildFromFloat(ctx *Context, value float64, _ reflect.Value) reflect.Value {
	result := globalInterfaceBuilder.BuildFromFloat(ctx, value, _this.components[_this.index])
	_this.tryFinish(ctx)
	return result
}

func (_this *edgeBuilder) BuildFromBigFloat(ctx *Context, value *big.Float, _ reflect.Value) reflect.Value {
	result := globalInterfaceBuilder.BuildFromBigFloat(ctx, value, _this.components[_this.index])
	_this.tryFinish(ctx)
	return result
}

func (_this *edgeBuilder) BuildFromDecimalFloat(ctx *Context, value compact_float.DFloat, _ reflect.Value) reflect.Value {
	result := globalInterfaceBuilder.BuildFromDecimalFloat(ctx, value, _this.components[_this.index])
	_this.tryFinish(ctx)
	return result
}

func (_this *edgeBuilder) BuildFromBigDecimalFloat(ctx *Context, value *apd.Decimal, _ reflect.Value) reflect.Value {
	result := globalInterfaceBuilder.BuildFromBigDecimalFloat(ctx, value, _this.components[_this.index])
	_this.tryFinish(ctx)
	return result
}

func (_this *edgeBuilder) BuildFromUID(ctx *Context, value []byte, _ reflect.Value) reflect.Value {
	result := globalInterfaceBuilder.BuildFromUID(ctx, value, _this.components[_this.index])
	_this.tryFinish(ctx)
	return result
}

func (_this *edgeBuilder) BuildFromArray(ctx *Context, arrayType events.ArrayType, value []byte, dst reflect.Value) reflect.Value {
	result := globalInterfaceBuilder.BuildFromArray(ctx, arrayType, value, _this.components[_this.index])
	_this.tryFinish(ctx)
	return result
}

func (_this *edgeBuilder) BuildFromStringlikeArray(ctx *Context, arrayType events.ArrayType, value string, dst reflect.Value) reflect.Value {
	result := globalInterfaceBuilder.BuildFromStringlikeArray(ctx, arrayType, value, _this.components[_this.index])
	_this.tryFinish(ctx)
	return result
}

func (_this *edgeBuilder) BuildFromCustomBinary(ctx *Context, customType uint64, data []byte, dst reflect.Value) reflect.Value {
	result := globalInterfaceBuilder.BuildFromCustomBinary(ctx, customType, data, _this.components[_this.index])
	_this.tryFinish(ctx)
	return result
}

func (_this *edgeBuilder) BuildFromCustomText(ctx *Context, customType uint64, data string, dst reflect.Value) reflect.Value {
	result := globalInterfaceBuilder.BuildFromCustomText(ctx, customType, data, _this.components[_this.index])
	_this.tryFinish(ctx)
	return result
}

func (_this *edgeBuilder) BuildFromMedia(ctx *Context, mediaType string, data []byte, _ reflect.Value) reflect.Value {
	result := globalInterfaceBuilder.BuildFromMedia(ctx, mediaType, data, _this.components[_this.index])
	_this.tryFinish(ctx)
	return result
}

func (_this *edgeBuilder) BuildFromTime(ctx *Context, value compact_time.Time, _ reflect.Value) reflect.Value {
	result := globalInterfaceBuilder.BuildFromTime(ctx, value, _this.components[_this.index])
	_this.tryFinish(ctx)
	return result
}

func (_this *edgeBuilder) BuildNewList(ctx *Context) {
	globalInterfaceBuilder.BuildBeginListContents(ctx)
}

func (_this *edgeBuilder) BuildNewMap(ctx *Context) {
	globalInterfaceBuilder.BuildBeginMapContents(ctx)
}

func (_this *edgeBuilder) BuildNewEdge(ctx *Context) {
	globalInterfaceBuilder.BuildBeginEdgeContents(ctx)
}

func (_this *edgeBuilder) BuildNewNode(ctx *Context) {
	globalInterfaceBuilder.BuildBeginNodeContents(ctx)
}

func (_this *edgeBuilder) BuildFromLocalReference(ctx *Context, id []byte) {
	globalInterfaceBuilder.BuildFromLocalReference(ctx, id)
	_this.tryFinish(ctx)
}

func (_this *edgeBuilder) BuildBeginEdgeContents(ctx *Context) {
	ctx.StackBuilder(_this)
}

func (_this *edgeBuilder) NotifyChildContainerFinished(ctx *Context, value reflect.Value) {
	_this.components[_this.index] = value
	_this.tryFinish(ctx)
}
