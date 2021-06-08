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

type relationshipBuilder struct {
	builders     []Builder
	components   []reflect.Value
	index        int
	relationship types.Relationship
}

func generateRelationshipBuilder(ctx *Context) Builder {
	return &relationshipBuilder{
		builders: []Builder{
			globalInterfaceBuilder,
			globalPRidBuilder,
			globalInterfaceBuilder,
		},
		components: []reflect.Value{
			reflect.New(common.TypeInterface).Elem(),
			reflect.New(common.TypePURL).Elem(),
			reflect.New(common.TypeInterface).Elem(),
		},
	}
}

func (_this *relationshipBuilder) String() string { return fmt.Sprintf("%v", reflect.TypeOf(_this)) }

func (_this *relationshipBuilder) tryFinish(ctx *Context) {
	_this.index++
	if _this.index >= 3 {
		// TODO: Fighting the compiler here to stop it from making a stack allocation
		obj := reflect.New(reflect.TypeOf(types.Relationship{})).Elem()
		obj.FieldByName("Subject").Set(_this.components[0])
		obj.FieldByName("Predicate").Set(_this.components[1])
		obj.FieldByName("Object").Set(_this.components[2])
		ctx.UnstackBuilder()
		ctx.CurrentBuilder.NotifyChildContainerFinished(ctx, obj)
	}
}

func (_this *relationshipBuilder) BuildFromNil(ctx *Context, _ reflect.Value) reflect.Value {
	result := _this.builders[_this.index].BuildFromNil(ctx, _this.components[_this.index])
	_this.tryFinish(ctx)
	return result
}

func (_this *relationshipBuilder) BuildFromBool(ctx *Context, value bool, _ reflect.Value) reflect.Value {
	result := _this.builders[_this.index].BuildFromBool(ctx, value, _this.components[_this.index])
	_this.tryFinish(ctx)
	return result
}

func (_this *relationshipBuilder) BuildFromInt(ctx *Context, value int64, _ reflect.Value) reflect.Value {
	result := _this.builders[_this.index].BuildFromInt(ctx, value, _this.components[_this.index])
	_this.tryFinish(ctx)
	return result
}

func (_this *relationshipBuilder) BuildFromUint(ctx *Context, value uint64, _ reflect.Value) reflect.Value {
	result := _this.builders[_this.index].BuildFromUint(ctx, value, _this.components[_this.index])
	_this.tryFinish(ctx)
	return result
}

func (_this *relationshipBuilder) BuildFromBigInt(ctx *Context, value *big.Int, _ reflect.Value) reflect.Value {
	result := _this.builders[_this.index].BuildFromBigInt(ctx, value, _this.components[_this.index])
	_this.tryFinish(ctx)
	return result
}

func (_this *relationshipBuilder) BuildFromFloat(ctx *Context, value float64, _ reflect.Value) reflect.Value {
	result := _this.builders[_this.index].BuildFromFloat(ctx, value, _this.components[_this.index])
	_this.tryFinish(ctx)
	return result
}

func (_this *relationshipBuilder) BuildFromBigFloat(ctx *Context, value *big.Float, _ reflect.Value) reflect.Value {
	result := _this.builders[_this.index].BuildFromBigFloat(ctx, value, _this.components[_this.index])
	_this.tryFinish(ctx)
	return result
}

func (_this *relationshipBuilder) BuildFromDecimalFloat(ctx *Context, value compact_float.DFloat, _ reflect.Value) reflect.Value {
	result := _this.builders[_this.index].BuildFromDecimalFloat(ctx, value, _this.components[_this.index])
	_this.tryFinish(ctx)
	return result
}

func (_this *relationshipBuilder) BuildFromBigDecimalFloat(ctx *Context, value *apd.Decimal, _ reflect.Value) reflect.Value {
	result := _this.builders[_this.index].BuildFromBigDecimalFloat(ctx, value, _this.components[_this.index])
	_this.tryFinish(ctx)
	return result
}

func (_this *relationshipBuilder) BuildFromUID(ctx *Context, value []byte, _ reflect.Value) reflect.Value {
	result := _this.builders[_this.index].BuildFromUID(ctx, value, _this.components[_this.index])
	_this.tryFinish(ctx)
	return result
}

func (_this *relationshipBuilder) BuildFromArray(ctx *Context, arrayType events.ArrayType, value []byte, dst reflect.Value) reflect.Value {
	result := _this.builders[_this.index].BuildFromArray(ctx, arrayType, value, _this.components[_this.index])
	_this.tryFinish(ctx)
	return result
}

func (_this *relationshipBuilder) BuildFromStringlikeArray(ctx *Context, arrayType events.ArrayType, value string, dst reflect.Value) reflect.Value {
	result := _this.builders[_this.index].BuildFromStringlikeArray(ctx, arrayType, value, _this.components[_this.index])
	_this.tryFinish(ctx)
	return result
}

func (_this *relationshipBuilder) BuildFromMedia(ctx *Context, mediaType string, data []byte, _ reflect.Value) reflect.Value {
	result := _this.builders[_this.index].BuildFromMedia(ctx, mediaType, data, _this.components[_this.index])
	_this.tryFinish(ctx)
	return result
}

func (_this *relationshipBuilder) BuildFromTime(ctx *Context, value time.Time, _ reflect.Value) reflect.Value {
	result := _this.builders[_this.index].BuildFromTime(ctx, value, _this.components[_this.index])
	_this.tryFinish(ctx)
	return result
}

func (_this *relationshipBuilder) BuildFromCompactTime(ctx *Context, value compact_time.Time, _ reflect.Value) reflect.Value {
	result := _this.builders[_this.index].BuildFromCompactTime(ctx, value, _this.components[_this.index])
	_this.tryFinish(ctx)
	return result
}

func (_this *relationshipBuilder) BuildInitiateList(ctx *Context) {
	_this.builders[_this.index].BuildBeginListContents(ctx)
}

func (_this *relationshipBuilder) BuildInitiateMap(ctx *Context) {
	_this.builders[_this.index].BuildBeginMapContents(ctx)
}

func (_this *relationshipBuilder) BuildFromReference(ctx *Context, id []byte) {
	_this.builders[_this.index].BuildFromReference(ctx, id)
	_this.tryFinish(ctx)
}

func (_this *relationshipBuilder) BuildInitiateMarkup(ctx *Context, name []byte) {
	_this.builders[_this.index].BuildBeginMarkupContents(ctx, name)
	_this.tryFinish(ctx)
}

func (_this *relationshipBuilder) NotifyChildContainerFinished(ctx *Context, value reflect.Value) {
	_this.components[_this.index] = value
	_this.tryFinish(ctx)
}
