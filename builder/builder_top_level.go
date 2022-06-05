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

	"github.com/cockroachdb/apd/v2"
	compact_float "github.com/kstenerud/go-compact-float"
	compact_time "github.com/kstenerud/go-compact-time"
	"github.com/kstenerud/go-concise-encoding/events"
)

// topLevelContainerBuilder proxies the first build instruction to make sure containers
// are properly built. See BuildNewList and BuildNewMap.
type topLevelBuilder struct {
	builderGenerator          BuilderGenerator
	containerFinishedCallback func(value reflect.Value)
}

func newTopLevelBuilder(builderGenerator BuilderGenerator, containerFinishedCallback func(value reflect.Value)) Builder {
	return &topLevelBuilder{
		builderGenerator:          builderGenerator,
		containerFinishedCallback: containerFinishedCallback,
	}
}

func (_this *topLevelBuilder) String() string { return reflect.TypeOf(_this).String() }

func (_this *topLevelBuilder) BuildFromNull(ctx *Context, dst reflect.Value) reflect.Value {
	_this.builderGenerator(ctx).BuildFromNull(ctx, dst)
	return dst
}

func (_this *topLevelBuilder) BuildFromBool(ctx *Context, value bool, dst reflect.Value) reflect.Value {
	_this.builderGenerator(ctx).BuildFromBool(ctx, value, dst)
	return dst
}

func (_this *topLevelBuilder) BuildFromInt(ctx *Context, value int64, dst reflect.Value) reflect.Value {
	_this.builderGenerator(ctx).BuildFromInt(ctx, value, dst)
	return dst
}

func (_this *topLevelBuilder) BuildFromUint(ctx *Context, value uint64, dst reflect.Value) reflect.Value {
	_this.builderGenerator(ctx).BuildFromUint(ctx, value, dst)
	return dst
}

func (_this *topLevelBuilder) BuildFromBigInt(ctx *Context, value *big.Int, dst reflect.Value) reflect.Value {
	_this.builderGenerator(ctx).BuildFromBigInt(ctx, value, dst)
	return dst
}

func (_this *topLevelBuilder) BuildFromFloat(ctx *Context, value float64, dst reflect.Value) reflect.Value {
	_this.builderGenerator(ctx).BuildFromFloat(ctx, value, dst)
	return dst
}

func (_this *topLevelBuilder) BuildFromBigFloat(ctx *Context, value *big.Float, dst reflect.Value) reflect.Value {
	_this.builderGenerator(ctx).BuildFromBigFloat(ctx, value, dst)
	return dst
}

func (_this *topLevelBuilder) BuildFromDecimalFloat(ctx *Context, value compact_float.DFloat, dst reflect.Value) reflect.Value {
	_this.builderGenerator(ctx).BuildFromDecimalFloat(ctx, value, dst)
	return dst
}

func (_this *topLevelBuilder) BuildFromBigDecimalFloat(ctx *Context, value *apd.Decimal, dst reflect.Value) reflect.Value {
	_this.builderGenerator(ctx).BuildFromBigDecimalFloat(ctx, value, dst)
	return dst
}

func (_this *topLevelBuilder) BuildFromUID(ctx *Context, value []byte, dst reflect.Value) reflect.Value {
	_this.builderGenerator(ctx).BuildFromUID(ctx, value, dst)
	return dst
}

func (_this *topLevelBuilder) BuildFromArray(ctx *Context, arrayType events.ArrayType, value []byte, dst reflect.Value) reflect.Value {
	_this.builderGenerator(ctx).BuildFromArray(ctx, arrayType, value, dst)
	return dst
}

func (_this *topLevelBuilder) BuildFromStringlikeArray(ctx *Context, arrayType events.ArrayType, value string, dst reflect.Value) reflect.Value {
	_this.builderGenerator(ctx).BuildFromStringlikeArray(ctx, arrayType, value, dst)
	return dst
}

func (_this *topLevelBuilder) BuildFromMedia(ctx *Context, mediaType string, data []byte, dst reflect.Value) reflect.Value {
	_this.builderGenerator(ctx).BuildFromMedia(ctx, mediaType, data, dst)
	return dst
}

func (_this *topLevelBuilder) BuildFromTime(ctx *Context, value time.Time, dst reflect.Value) reflect.Value {
	_this.builderGenerator(ctx).BuildFromTime(ctx, value, dst)
	return dst
}

func (_this *topLevelBuilder) BuildFromCompactTime(ctx *Context, value compact_time.Time, dst reflect.Value) reflect.Value {
	_this.builderGenerator(ctx).BuildFromCompactTime(ctx, value, dst)
	return dst
}

func (_this *topLevelBuilder) BuildNewList(ctx *Context) {
	if reflect.TypeOf(_this.builderGenerator) == reflect.TypeOf((*interfaceBuilder)(nil)) {
		_this.builderGenerator = interfaceSliceBuilderGenerator
	}
	builder := _this.builderGenerator(ctx)
	builder.BuildBeginListContents(ctx)
}

func (_this *topLevelBuilder) BuildNewMap(ctx *Context) {
	if reflect.TypeOf(_this.builderGenerator) == reflect.TypeOf(interfaceBuilder{}) {
		_this.builderGenerator = interfaceMapBuilderGenerator
	}
	_this.builderGenerator(ctx).BuildBeginMapContents(ctx)
}

func (_this *topLevelBuilder) BuildNewNode(ctx *Context) {
	if reflect.TypeOf(_this.builderGenerator) == reflect.TypeOf((*interfaceBuilder)(nil)) {
		_this.builderGenerator = interfaceNodeBuilderGenerator
	}
	builder := _this.builderGenerator(ctx)
	builder.BuildBeginNodeContents(ctx)
}

func (_this *topLevelBuilder) BuildNewEdge(ctx *Context) {
	if reflect.TypeOf(_this.builderGenerator) == reflect.TypeOf((*interfaceBuilder)(nil)) {
		_this.builderGenerator = interfaceEdgeBuilderGenerator
	}
	builder := _this.builderGenerator(ctx)
	builder.BuildBeginEdgeContents(ctx)
}

func (_this *topLevelBuilder) NotifyChildContainerFinished(ctx *Context, value reflect.Value) {
	_this.containerFinishedCallback(value)
}
