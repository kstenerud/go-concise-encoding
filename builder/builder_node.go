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
	"github.com/kstenerud/go-concise-encoding/internal/common"
	"github.com/kstenerud/go-concise-encoding/types"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
)

type nodeBuilder struct {
	isBuildingChildren bool
	node               reflect.Value
	value              reflect.Value
}

func generateNodeBuilder(ctx *Context) Builder {
	node := reflect.New(common.TypeNode).Elem()
	return &nodeBuilder{
		isBuildingChildren: false,
		node:               node,
		value:              node.Field(types.NodeFieldIndexValue),
	}
}

func (_this *nodeBuilder) String() string {
	return fmt.Sprintf("%v", reflect.TypeOf(_this))
}

func (_this *nodeBuilder) stackChildrenBuilder(ctx *Context) {
	_this.isBuildingChildren = true
	interfaceSliceBuilderGenerator(ctx).BuildBeginListContents(ctx)
}

func (_this *nodeBuilder) BuildFromNil(ctx *Context, _ reflect.Value) reflect.Value {
	globalInterfaceBuilder.BuildFromNil(ctx, _this.value)
	_this.stackChildrenBuilder(ctx)
	return _this.value
}

func (_this *nodeBuilder) BuildFromBool(ctx *Context, value bool, _ reflect.Value) reflect.Value {
	globalInterfaceBuilder.BuildFromBool(ctx, value, _this.value)
	_this.stackChildrenBuilder(ctx)
	return _this.value
}

func (_this *nodeBuilder) BuildFromInt(ctx *Context, value int64, _ reflect.Value) reflect.Value {
	globalInterfaceBuilder.BuildFromInt(ctx, value, _this.value)
	_this.stackChildrenBuilder(ctx)
	return _this.value
}

func (_this *nodeBuilder) BuildFromUint(ctx *Context, value uint64, _ reflect.Value) reflect.Value {
	globalInterfaceBuilder.BuildFromUint(ctx, value, _this.value)
	_this.stackChildrenBuilder(ctx)
	return _this.value
}

func (_this *nodeBuilder) BuildFromBigInt(ctx *Context, value *big.Int, _ reflect.Value) reflect.Value {
	globalInterfaceBuilder.BuildFromBigInt(ctx, value, _this.value)
	_this.stackChildrenBuilder(ctx)
	return _this.value
}

func (_this *nodeBuilder) BuildFromFloat(ctx *Context, value float64, _ reflect.Value) reflect.Value {
	globalInterfaceBuilder.BuildFromFloat(ctx, value, _this.value)
	_this.stackChildrenBuilder(ctx)
	return _this.value
}

func (_this *nodeBuilder) BuildFromBigFloat(ctx *Context, value *big.Float, _ reflect.Value) reflect.Value {
	globalInterfaceBuilder.BuildFromBigFloat(ctx, value, _this.value)
	_this.stackChildrenBuilder(ctx)
	return _this.value
}

func (_this *nodeBuilder) BuildFromDecimalFloat(ctx *Context, value compact_float.DFloat, _ reflect.Value) reflect.Value {
	globalInterfaceBuilder.BuildFromDecimalFloat(ctx, value, _this.value)
	_this.stackChildrenBuilder(ctx)
	return _this.value
}

func (_this *nodeBuilder) BuildFromBigDecimalFloat(ctx *Context, value *apd.Decimal, _ reflect.Value) reflect.Value {
	globalInterfaceBuilder.BuildFromBigDecimalFloat(ctx, value, _this.value)
	_this.stackChildrenBuilder(ctx)
	return _this.value
}

func (_this *nodeBuilder) BuildFromUID(ctx *Context, value []byte, _ reflect.Value) reflect.Value {
	globalInterfaceBuilder.BuildFromUID(ctx, value, _this.value)
	_this.stackChildrenBuilder(ctx)
	return _this.value
}

func (_this *nodeBuilder) BuildFromArray(ctx *Context, arrayType events.ArrayType, value []byte, _ reflect.Value) reflect.Value {
	globalInterfaceBuilder.BuildFromArray(ctx, arrayType, value, _this.value)
	_this.stackChildrenBuilder(ctx)
	return _this.value
}

func (_this *nodeBuilder) BuildFromStringlikeArray(ctx *Context, arrayType events.ArrayType, value string, _ reflect.Value) reflect.Value {
	globalInterfaceBuilder.BuildFromStringlikeArray(ctx, arrayType, value, _this.value)
	_this.stackChildrenBuilder(ctx)
	return _this.value
}

func (_this *nodeBuilder) BuildFromMedia(ctx *Context, mediaType string, data []byte, _ reflect.Value) reflect.Value {
	globalInterfaceBuilder.BuildFromMedia(ctx, mediaType, data, _this.value)
	_this.stackChildrenBuilder(ctx)
	return _this.value
}

func (_this *nodeBuilder) BuildFromTime(ctx *Context, value time.Time, _ reflect.Value) reflect.Value {
	globalInterfaceBuilder.BuildFromTime(ctx, value, _this.value)
	_this.stackChildrenBuilder(ctx)
	return _this.value
}

func (_this *nodeBuilder) BuildFromCompactTime(ctx *Context, value compact_time.Time, _ reflect.Value) reflect.Value {
	globalInterfaceBuilder.BuildFromCompactTime(ctx, value, _this.value)
	_this.stackChildrenBuilder(ctx)
	return _this.value
}

func (_this *nodeBuilder) BuildNewList(ctx *Context) {
	globalInterfaceBuilder.BuildBeginListContents(ctx)
}

func (_this *nodeBuilder) BuildNewMap(ctx *Context) {
	globalInterfaceBuilder.BuildBeginMapContents(ctx)
}

func (_this *nodeBuilder) BuildNewNode(ctx *Context) {
	globalInterfaceBuilder.BuildBeginNodeContents(ctx)
}

func (_this *nodeBuilder) BuildNewEdge(ctx *Context) {
	globalInterfaceBuilder.BuildBeginEdgeContents(ctx)
}

func (_this *nodeBuilder) BuildNewMarkup(ctx *Context, name []byte) {
	globalInterfaceBuilder.BuildBeginMarkupContents(ctx, name)
}

func (_this *nodeBuilder) BuildBeginNodeContents(ctx *Context) {
	ctx.StackBuilder(_this)
}

func (_this *nodeBuilder) BuildFromReference(ctx *Context, id []byte) {
	ctx.NotifyReference(id, func(object reflect.Value) {
		setAnythingFromAnything(object, _this.value)
	})
	_this.stackChildrenBuilder(ctx)
}

func (_this *nodeBuilder) NotifyChildContainerFinished(ctx *Context, container reflect.Value) {
	if _this.isBuildingChildren {
		_this.node.Field(types.NodeFieldIndexChildren).Set(container)
		ctx.UnstackBuilderAndNotifyChildFinished(_this.node)
	} else {
		_this.value.Set(container)
		_this.stackChildrenBuilder(ctx)
	}
}
