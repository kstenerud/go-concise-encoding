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

const (
	kvBuilderKey   = 0
	kvBuilderValue = 1
)

type mapBuilder struct {
	mapType         reflect.Type
	kvTypes         [2]reflect.Type
	kvGenerators    [2]BuilderGenerator
	container       reflect.Value
	key             reflect.Value
	builderIndex    int
	nextGenerator   BuilderGenerator
	nextStoreMethod func(*mapBuilder, reflect.Value)
}

func newMapBuilderGenerator(getBuilderGeneratorForType BuilderGeneratorGetter, mapType reflect.Type) BuilderGenerator {
	kvTypes := [2]reflect.Type{mapType.Key(), mapType.Elem()}
	kvGenerators := [2]BuilderGenerator{getBuilderGeneratorForType(kvTypes[0]), getBuilderGeneratorForType(kvTypes[1])}

	return func(ctx *Context) Builder {
		builder := &mapBuilder{
			mapType:      mapType,
			kvTypes:      kvTypes,
			kvGenerators: kvGenerators,
		}
		return builder
	}
}

func (_this *mapBuilder) String() string {
	return fmt.Sprintf("%v<%v:%v>", reflect.TypeOf(_this), _this.kvGenerators[0], _this.kvGenerators[1])
}

func (_this *mapBuilder) reset() {
	_this.container = reflect.MakeMap(_this.mapType)
	_this.key = reflect.Value{}
	_this.builderIndex = kvBuilderKey
	_this.nextGenerator = _this.kvGenerators[_this.builderIndex]
	_this.nextStoreMethod = mapBuilderKVStoreMethods[_this.builderIndex]
}

func (_this *mapBuilder) storeKey(value reflect.Value) {
	_this.key = value
}

func (_this *mapBuilder) storeValue(value reflect.Value) {
	_this.container.SetMapIndex(_this.key, value)
}

var mapBuilderKVStoreMethods = []func(*mapBuilder, reflect.Value){
	(*mapBuilder).storeKey,
	(*mapBuilder).storeValue,
}

func (_this *mapBuilder) store(value reflect.Value) {
	_this.nextStoreMethod(_this, value)
	_this.swapKeyValue()
}

func (_this *mapBuilder) swapKeyValue() {
	_this.builderIndex = (_this.builderIndex + 1) & 1
	_this.nextGenerator = _this.kvGenerators[_this.builderIndex]
	_this.nextStoreMethod = mapBuilderKVStoreMethods[_this.builderIndex]
}

func (_this *mapBuilder) newElem() reflect.Value {
	return reflect.New(_this.kvTypes[_this.builderIndex]).Elem()
}

func (_this *mapBuilder) BuildFromNil(ctx *Context, _ reflect.Value) reflect.Value {
	object := _this.newElem()
	if _this.container.IsValid() {
		_this.nextGenerator(ctx).BuildFromNil(ctx, object)
		_this.store(object)
	}
	return object
}

func (_this *mapBuilder) BuildFromBool(ctx *Context, value bool, _ reflect.Value) reflect.Value {
	object := _this.newElem()
	_this.nextGenerator(ctx).BuildFromBool(ctx, value, object)
	_this.store(object)
	return object
}

func (_this *mapBuilder) BuildFromInt(ctx *Context, value int64, _ reflect.Value) reflect.Value {
	object := _this.newElem()
	_this.nextGenerator(ctx).BuildFromInt(ctx, value, object)
	_this.store(object)
	return object
}

func (_this *mapBuilder) BuildFromUint(ctx *Context, value uint64, _ reflect.Value) reflect.Value {
	object := _this.newElem()
	_this.nextGenerator(ctx).BuildFromUint(ctx, value, object)
	_this.store(object)
	return object
}

func (_this *mapBuilder) BuildFromBigInt(ctx *Context, value *big.Int, _ reflect.Value) reflect.Value {
	object := _this.newElem()
	_this.nextGenerator(ctx).BuildFromBigInt(ctx, value, object)
	_this.store(object)
	return object
}

func (_this *mapBuilder) BuildFromFloat(ctx *Context, value float64, _ reflect.Value) reflect.Value {
	object := _this.newElem()
	_this.nextGenerator(ctx).BuildFromFloat(ctx, value, object)
	_this.store(object)
	return object
}

func (_this *mapBuilder) BuildFromBigFloat(ctx *Context, value *big.Float, _ reflect.Value) reflect.Value {
	object := _this.newElem()
	_this.nextGenerator(ctx).BuildFromBigFloat(ctx, value, object)
	_this.store(object)
	return object
}

func (_this *mapBuilder) BuildFromDecimalFloat(ctx *Context, value compact_float.DFloat, _ reflect.Value) reflect.Value {
	object := _this.newElem()
	_this.nextGenerator(ctx).BuildFromDecimalFloat(ctx, value, object)
	_this.store(object)
	return object
}

func (_this *mapBuilder) BuildFromBigDecimalFloat(ctx *Context, value *apd.Decimal, _ reflect.Value) reflect.Value {
	object := _this.newElem()
	_this.nextGenerator(ctx).BuildFromBigDecimalFloat(ctx, value, object)
	_this.store(object)
	return object
}

func (_this *mapBuilder) BuildFromUUID(ctx *Context, value []byte, _ reflect.Value) reflect.Value {
	object := _this.newElem()
	_this.nextGenerator(ctx).BuildFromUUID(ctx, value, object)
	_this.store(object)
	return object
}

func (_this *mapBuilder) BuildFromArray(ctx *Context, arrayType events.ArrayType, value []byte, _ reflect.Value) reflect.Value {
	object := _this.newElem()
	_this.nextGenerator(ctx).BuildFromArray(ctx, arrayType, value, object)
	_this.store(object)
	return object
}

func (_this *mapBuilder) BuildFromStringlikeArray(ctx *Context, arrayType events.ArrayType, value string, _ reflect.Value) reflect.Value {
	object := _this.newElem()
	_this.nextGenerator(ctx).BuildFromStringlikeArray(ctx, arrayType, value, object)
	_this.store(object)
	return object
}

func (_this *mapBuilder) BuildFromTime(ctx *Context, value time.Time, _ reflect.Value) reflect.Value {
	// TODO: Why was it this way?
	// _this.store(reflect.ValueOf(value))
	object := _this.newElem()
	_this.nextGenerator(ctx).BuildFromTime(ctx, value, object)
	_this.store(object)
	return object
}

func (_this *mapBuilder) BuildFromCompactTime(ctx *Context, value compact_time.Time, _ reflect.Value) reflect.Value {
	// TODO: Why was it this way?
	// _this.store(reflect.ValueOf(value))
	object := _this.newElem()
	_this.nextGenerator(ctx).BuildFromCompactTime(ctx, value, object)
	_this.store(object)
	return object
}

func (_this *mapBuilder) BuildInitiateList(ctx *Context) {
	_this.nextGenerator(ctx).BuildBeginListContents(ctx)
}

func (_this *mapBuilder) BuildInitiateMap(ctx *Context) {
	_this.nextGenerator(ctx).BuildBeginMapContents(ctx)
}

func (_this *mapBuilder) BuildEndContainer(ctx *Context) {
	object := _this.container
	ctx.UnstackBuilderAndNotifyChildFinished(object)
}

func (_this *mapBuilder) BuildBeginMapContents(ctx *Context) {
	ctx.StackBuilder(_this)
	_this.reset()
}

func (_this *mapBuilder) BuildFromReference(ctx *Context, id []byte) {
	container := _this.container
	key := _this.key
	tempValue := _this.newElem()
	_this.swapKeyValue()
	ctx.NotifyReference(id, func(object reflect.Value) {
		if container.Type().Elem().Kind() == reflect.Interface || object.Type() == container.Type().Elem() {
			// In case of self-referencing pointers, we need to pass the original container, not a copy.
			container.SetMapIndex(key, object)
		} else {
			setAnythingFromAnything(object, tempValue)
			container.SetMapIndex(key, tempValue)
		}
	})
}

func (_this *mapBuilder) NotifyChildContainerFinished(ctx *Context, value reflect.Value) {
	_this.store(value)
}
