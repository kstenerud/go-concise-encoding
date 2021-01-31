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

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
)

type referenceIDBuilder struct{}

var globalReferenceIDBuilder = &referenceIDBuilder{}

func newReferenceIDBuilder() *referenceIDBuilder { return globalReferenceIDBuilder }
func (_this *referenceIDBuilder) String() string { return reflect.TypeOf(_this).String() }

func (_this *referenceIDBuilder) BuildFromInt(ctx *Context, value int64, _ reflect.Value) reflect.Value {
	if value < 0 {
		PanicBadEvent(_this, "Negative Int")
	}
	ctx.StoreReferencedObject(value)
	return reflect.ValueOf(value)
}

func (_this *referenceIDBuilder) BuildFromUint(ctx *Context, value uint64, _ reflect.Value) reflect.Value {
	ctx.StoreReferencedObject(value)
	return reflect.ValueOf(value)
}

func (_this *referenceIDBuilder) BuildFromBigInt(ctx *Context, value *big.Int, _ reflect.Value) reflect.Value {
	if common.IsBigIntNegative(value) || !value.IsUint64() {
		PanicBadEvent(_this, "BigInt")
	}
	ctx.StoreReferencedObject(value.Uint64())
	return reflect.ValueOf(value)
}

// ============================================================================

type markerIDBuilder struct{}

var globalMarkerIDBuilder = &markerIDBuilder{}

func newMarkerIDBuilder() *markerIDBuilder    { return globalMarkerIDBuilder }
func (_this *markerIDBuilder) String() string { return reflect.TypeOf(_this).String() }

func (_this *markerIDBuilder) BuildFromInt(ctx *Context, value int64, _ reflect.Value) reflect.Value {
	if value < 0 {
		PanicBadEvent(_this, "Negative Int")
	}
	ctx.BeginMarkerObject(value)
	return reflect.ValueOf(value)
}

func (_this *markerIDBuilder) BuildFromUint(ctx *Context, value uint64, _ reflect.Value) reflect.Value {
	ctx.BeginMarkerObject(value)
	return reflect.ValueOf(value)
}

func (_this *markerIDBuilder) BuildFromBigInt(ctx *Context, value *big.Int, _ reflect.Value) reflect.Value {
	if common.IsBigIntNegative(value) || !value.IsUint64() {
		PanicBadEvent(_this, "BigInt")
	}
	ctx.BeginMarkerObject(value.Uint64())
	return reflect.ValueOf(value)
}

// ============================================================================

type markerObjectBuilder struct {
	isContainer bool
	id          interface{}
	child       ObjectBuilder
}

func newMarkerObjectBuilder(id interface{}, child ObjectBuilder) *markerObjectBuilder {
	return &markerObjectBuilder{
		id:    id,
		child: child,
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

func (_this *markerObjectBuilder) BuildFromUUID(ctx *Context, value []byte, dst reflect.Value) reflect.Value {
	object := _this.child.BuildFromUUID(ctx, value, dst)
	_this.onObjectFinished(ctx, object)
	return object
}

func (_this *markerObjectBuilder) BuildFromArray(ctx *Context, arrayType events.ArrayType, value []byte, dst reflect.Value) reflect.Value {
	object := _this.child.BuildFromArray(ctx, arrayType, value, dst)
	_this.onObjectFinished(ctx, object)
	return object
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

func (_this *markerObjectBuilder) BuildInitiateList(ctx *Context) {
	_this.isContainer = true
	_this.child.BuildInitiateList(ctx)
}

func (_this *markerObjectBuilder) BuildInitiateMap(ctx *Context) {
	_this.isContainer = true
	_this.child.BuildInitiateMap(ctx)
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
