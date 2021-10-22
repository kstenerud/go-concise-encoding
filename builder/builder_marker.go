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

type markerObjectBuilder struct {
	isContainer bool
	id          []byte
	child       Builder
}

func newMarkerObjectBuilder(id []byte, child Builder) *markerObjectBuilder {
	return &markerObjectBuilder{
		isContainer: false,
		id:          id,
		child:       child,
	}
}

func (_this *markerObjectBuilder) String() string {
	return fmt.Sprintf("%v<%v>", reflect.TypeOf(_this), _this.child)
}

func (_this *markerObjectBuilder) onObjectFinished(ctx *Context, dst reflect.Value) {
	if !_this.isContainer {
		ctx.UnstackBuilder()
		ctx.NotifyMarker(_this.id, dst)
	}
}

func (_this *markerObjectBuilder) BuildFromNil(ctx *Context, dst reflect.Value) reflect.Value {
	object := _this.child.BuildFromNil(ctx, dst)
	_this.onObjectFinished(ctx, object)
	return object
}

func (_this *markerObjectBuilder) BuildFromBool(ctx *Context, value bool, dst reflect.Value) reflect.Value {
	object := _this.child.BuildFromBool(ctx, value, dst)
	_this.onObjectFinished(ctx, object)
	return object
}

func (_this *markerObjectBuilder) BuildFromInt(ctx *Context, value int64, dst reflect.Value) reflect.Value {
	object := _this.child.BuildFromInt(ctx, value, dst)
	_this.onObjectFinished(ctx, object)
	return object
}

func (_this *markerObjectBuilder) BuildFromUint(ctx *Context, value uint64, dst reflect.Value) reflect.Value {
	object := _this.child.BuildFromUint(ctx, value, dst)
	_this.onObjectFinished(ctx, object)
	return object
}

func (_this *markerObjectBuilder) BuildFromBigInt(ctx *Context, value *big.Int, dst reflect.Value) reflect.Value {
	object := _this.child.BuildFromBigInt(ctx, value, dst)
	_this.onObjectFinished(ctx, object)
	return object
}

func (_this *markerObjectBuilder) BuildFromFloat(ctx *Context, value float64, dst reflect.Value) reflect.Value {
	object := _this.child.BuildFromFloat(ctx, value, dst)
	_this.onObjectFinished(ctx, object)
	return object
}

func (_this *markerObjectBuilder) BuildFromBigFloat(ctx *Context, value *big.Float, dst reflect.Value) reflect.Value {
	object := _this.child.BuildFromBigFloat(ctx, value, dst)
	_this.onObjectFinished(ctx, object)
	return object
}

func (_this *markerObjectBuilder) BuildFromDecimalFloat(ctx *Context, value compact_float.DFloat, dst reflect.Value) reflect.Value {
	object := _this.child.BuildFromDecimalFloat(ctx, value, dst)
	_this.onObjectFinished(ctx, object)
	return object
}

func (_this *markerObjectBuilder) BuildFromBigDecimalFloat(ctx *Context, value *apd.Decimal, dst reflect.Value) reflect.Value {
	object := _this.child.BuildFromBigDecimalFloat(ctx, value, dst)
	_this.onObjectFinished(ctx, object)
	return object
}

func (_this *markerObjectBuilder) BuildFromUID(ctx *Context, value []byte, dst reflect.Value) reflect.Value {
	object := _this.child.BuildFromUID(ctx, value, dst)
	_this.onObjectFinished(ctx, object)
	return object
}

func (_this *markerObjectBuilder) BuildFromArray(ctx *Context, arrayType events.ArrayType, value []byte, dst reflect.Value) reflect.Value {
	object := _this.child.BuildFromArray(ctx, arrayType, value, dst)
	_this.onObjectFinished(ctx, object)
	return object
}

func (_this *markerObjectBuilder) BuildFromStringlikeArray(ctx *Context, arrayType events.ArrayType, value string, dst reflect.Value) reflect.Value {
	object := _this.child.BuildFromStringlikeArray(ctx, arrayType, value, dst)
	_this.onObjectFinished(ctx, object)
	return object
}

func (_this *markerObjectBuilder) BuildFromMedia(ctx *Context, mediaType string, data []byte, dst reflect.Value) reflect.Value {
	object := _this.child.BuildFromMedia(ctx, mediaType, data, dst)
	_this.onObjectFinished(ctx, object)
	return dst
}

func (_this *markerObjectBuilder) BuildFromTime(ctx *Context, value time.Time, dst reflect.Value) reflect.Value {
	object := _this.child.BuildFromTime(ctx, value, dst)
	_this.onObjectFinished(ctx, object)
	return object
}

func (_this *markerObjectBuilder) BuildFromCompactTime(ctx *Context, value compact_time.Time, dst reflect.Value) reflect.Value {
	object := _this.child.BuildFromCompactTime(ctx, value, dst)
	_this.onObjectFinished(ctx, object)
	return object
}

func (_this *markerObjectBuilder) BuildNewList(ctx *Context) {
	_this.isContainer = true
	_this.child.BuildNewList(ctx)
}

func (_this *markerObjectBuilder) BuildNewMap(ctx *Context) {
	_this.isContainer = true
	_this.child.BuildNewMap(ctx)
}

func (_this *markerObjectBuilder) BuildNewMarkup(ctx *Context, name []byte) {
	_this.isContainer = true
	_this.child.BuildNewMarkup(ctx, name)
}

func (_this *markerObjectBuilder) BuildNewEdge(ctx *Context) {
	_this.isContainer = true
	_this.child.BuildNewEdge(ctx)
}

func (_this *markerObjectBuilder) BuildNewNode(ctx *Context) {
	_this.isContainer = true
	_this.child.BuildNewNode(ctx)
}

func (_this *markerObjectBuilder) BuildEndContainer(ctx *Context) {
	_this.child.BuildEndContainer(ctx)
}

func (_this *markerObjectBuilder) NotifyChildContainerFinished(ctx *Context, value reflect.Value) {
	if !_this.isContainer {
		PanicBadEvent(_this, "NotifyChildContainerFinished (isContainer = false)")
	}
	ctx.NotifyMarker(_this.id, value)
	ctx.UnstackBuilderAndNotifyChildFinished(value)
}
