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

// topLevelContainerBuilder proxies the first build instruction to make sure containers
// are properly built. See BuildInitiateList and BuildInitiateMap.
type topLevelBuilder struct {
	builder ObjectBuilder
	root    *RootBuilder
}

func newTopLevelBuilder(root *RootBuilder, builder ObjectBuilder) ObjectBuilder {
	return &topLevelBuilder{
		builder: builder,
		root:    root,
	}
}

func (_this *topLevelBuilder) String() string { return reflect.TypeOf(_this).String() }

func (_this *topLevelBuilder) BuildFromNil(ctx *Context, dst reflect.Value) reflect.Value {
	_this.builder.BuildFromNil(ctx, dst)
	return dst
}

func (_this *topLevelBuilder) BuildFromBool(ctx *Context, value bool, dst reflect.Value) reflect.Value {
	_this.builder.BuildFromBool(ctx, value, dst)
	return dst
}

func (_this *topLevelBuilder) BuildFromInt(ctx *Context, value int64, dst reflect.Value) reflect.Value {
	_this.builder.BuildFromInt(ctx, value, dst)
	return dst
}

func (_this *topLevelBuilder) BuildFromUint(ctx *Context, value uint64, dst reflect.Value) reflect.Value {
	_this.builder.BuildFromUint(ctx, value, dst)
	return dst
}

func (_this *topLevelBuilder) BuildFromBigInt(ctx *Context, value *big.Int, dst reflect.Value) reflect.Value {
	_this.builder.BuildFromBigInt(ctx, value, dst)
	return dst
}

func (_this *topLevelBuilder) BuildFromFloat(ctx *Context, value float64, dst reflect.Value) reflect.Value {
	_this.builder.BuildFromFloat(ctx, value, dst)
	return dst
}

func (_this *topLevelBuilder) BuildFromBigFloat(ctx *Context, value *big.Float, dst reflect.Value) reflect.Value {
	_this.builder.BuildFromBigFloat(ctx, value, dst)
	return dst
}

func (_this *topLevelBuilder) BuildFromDecimalFloat(ctx *Context, value compact_float.DFloat, dst reflect.Value) reflect.Value {
	_this.builder.BuildFromDecimalFloat(ctx, value, dst)
	return dst
}

func (_this *topLevelBuilder) BuildFromBigDecimalFloat(ctx *Context, value *apd.Decimal, dst reflect.Value) reflect.Value {
	_this.builder.BuildFromBigDecimalFloat(ctx, value, dst)
	return dst
}

func (_this *topLevelBuilder) BuildFromUUID(ctx *Context, value []byte, dst reflect.Value) reflect.Value {
	_this.builder.BuildFromUUID(ctx, value, dst)
	return dst
}

func (_this *topLevelBuilder) BuildFromArray(ctx *Context, arrayType events.ArrayType, value []byte, dst reflect.Value) reflect.Value {
	_this.builder.BuildFromArray(ctx, arrayType, value, dst)
	return dst
}

func (_this *topLevelBuilder) BuildFromTime(ctx *Context, value time.Time, dst reflect.Value) reflect.Value {
	_this.builder.BuildFromTime(ctx, value, dst)
	return dst
}

func (_this *topLevelBuilder) BuildFromCompactTime(ctx *Context, value *compact_time.Time, dst reflect.Value) reflect.Value {
	_this.builder.BuildFromCompactTime(ctx, value, dst)
	return dst
}

func (_this *topLevelBuilder) BuildInitiateList(ctx *Context) {
	if reflect.TypeOf(_this.builder) == reflect.TypeOf((*interfaceBuilder)(nil)) {
		_this.builder = interfaceSliceBuilderGenerator()
	}
	_this.builder.BuildBeginListContents(ctx)
}

func (_this *topLevelBuilder) BuildInitiateMap(ctx *Context) {
	if reflect.TypeOf(_this.builder) == reflect.TypeOf(interfaceBuilder{}) {
		_this.builder = interfaceMapBuilderGenerator()
	}
	_this.builder.BuildBeginMapContents(ctx)
}

func (_this *topLevelBuilder) NotifyChildContainerFinished(ctx *Context, value reflect.Value) {
	_this.root.NotifyChildContainerFinished(ctx, value)
}
