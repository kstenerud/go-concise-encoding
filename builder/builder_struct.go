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
)

type structBuilderField struct {
	Name      string
	IndexPath []int
}

func (_this *structBuilderField) GetField(container reflect.Value) (value reflect.Value) {
	value = container
	for _, index := range _this.IndexPath {
		value = value.Field(index)
	}
	return
}

func (_this *structBuilderField) applyTags(field reflect.StructField) {
	tags := common.DecodeGoTags(field)
	_this.Name = tags.Name
}

type structBuilder struct {
	dstType              reflect.Type
	generatorDescs       map[string]*structBuilderGeneratorDesc
	nameBuilderGenerator BuilderGenerator
	nextBuilderGenerator BuilderGenerator
	container            reflect.Value
	nextValue            reflect.Value
	nextIsKey            bool
}

type structBuilderGeneratorDesc struct {
	field            *structBuilderField
	builderGenerator BuilderGenerator
}

func makeGeneratorDescs(
	getBuilderGeneratorForType BuilderGeneratorGetter,
	dstType reflect.Type,
	currentPath []int,
	generatorDescs map[string]*structBuilderGeneratorDesc) {

	for i := 0; i < dstType.NumField(); i++ {
		reflectField := dstType.Field(i)
		if reflectField.IsExported() {
			path := make([]int, len(currentPath)+1)
			copy(path, currentPath)
			path[len(currentPath)] = i

			if reflectField.Anonymous {
				// Treat the embedded struct as if it were part of the parent struct
				makeGeneratorDescs(getBuilderGeneratorForType, reflectField.Type, path, generatorDescs)
			} else {
				builderGenerator := getBuilderGeneratorForType(reflectField.Type)
				structField := &structBuilderField{
					Name:      reflectField.Name,
					IndexPath: path,
				}
				structField.applyTags(reflectField)
				generatorDescs[structField.Name] = &structBuilderGeneratorDesc{
					field:            structField,
					builderGenerator: builderGenerator,
				}
			}
		}
	}
}

func newStructBuilderGenerator(getBuilderGeneratorForType BuilderGeneratorGetter, dstType reflect.Type) BuilderGenerator {
	nameBuilderGenerator := getBuilderGeneratorForType(reflect.TypeOf(""))
	generatorDescs := make(map[string]*structBuilderGeneratorDesc)

	makeGeneratorDescs(getBuilderGeneratorForType, dstType, []int{}, generatorDescs)

	// Make lowercase mappings as well in case we later do case-insensitive field name matching
	for _, desc := range generatorDescs {
		lowerName := common.ToStructFieldIdentifier(desc.field.Name)
		if _, exists := generatorDescs[lowerName]; !exists {
			generatorDescs[lowerName] = desc
		}
	}

	return func(ctx *Context) Builder {
		return &structBuilder{
			dstType:              dstType,
			generatorDescs:       generatorDescs,
			nameBuilderGenerator: nameBuilderGenerator,
			nextBuilderGenerator: nameBuilderGenerator,
			container:            reflect.New(dstType).Elem(),
			nextValue:            reflect.Value{},
			nextIsKey:            true,
		}
	}
}

func (_this *structBuilder) String() string {
	return fmt.Sprintf("%v<%v>", reflect.TypeOf(_this), _this.dstType)
}

func (_this *structBuilder) swapKeyValue() {
	_this.nextIsKey = !_this.nextIsKey
}

func (_this *structBuilder) BuildFromNull(ctx *Context, _ reflect.Value) reflect.Value {
	_this.nextBuilderGenerator(ctx).BuildFromNull(ctx, _this.nextValue)
	object := _this.nextValue
	_this.swapKeyValue()
	return object
}

func (_this *structBuilder) BuildFromBool(ctx *Context, value bool, _ reflect.Value) reflect.Value {
	_this.nextBuilderGenerator(ctx).BuildFromBool(ctx, value, _this.nextValue)
	object := _this.nextValue
	_this.swapKeyValue()
	return object
}

func (_this *structBuilder) BuildFromInt(ctx *Context, value int64, _ reflect.Value) reflect.Value {
	_this.nextBuilderGenerator(ctx).BuildFromInt(ctx, value, _this.nextValue)
	object := _this.nextValue
	_this.swapKeyValue()
	return object
}

func (_this *structBuilder) BuildFromUint(ctx *Context, value uint64, _ reflect.Value) reflect.Value {
	_this.nextBuilderGenerator(ctx).BuildFromUint(ctx, value, _this.nextValue)
	object := _this.nextValue
	_this.swapKeyValue()
	return object
}

func (_this *structBuilder) BuildFromBigInt(ctx *Context, value *big.Int, _ reflect.Value) reflect.Value {
	_this.nextBuilderGenerator(ctx).BuildFromBigInt(ctx, value, _this.nextValue)
	object := _this.nextValue
	_this.swapKeyValue()
	return object
}

func (_this *structBuilder) BuildFromFloat(ctx *Context, value float64, _ reflect.Value) reflect.Value {
	_this.nextBuilderGenerator(ctx).BuildFromFloat(ctx, value, _this.nextValue)
	object := _this.nextValue
	_this.swapKeyValue()
	return object
}

func (_this *structBuilder) BuildFromBigFloat(ctx *Context, value *big.Float, _ reflect.Value) reflect.Value {
	_this.nextBuilderGenerator(ctx).BuildFromBigFloat(ctx, value, _this.nextValue)
	object := _this.nextValue
	_this.swapKeyValue()
	return object
}

func (_this *structBuilder) BuildFromDecimalFloat(ctx *Context, value compact_float.DFloat, _ reflect.Value) reflect.Value {
	_this.nextBuilderGenerator(ctx).BuildFromDecimalFloat(ctx, value, _this.nextValue)
	object := _this.nextValue
	_this.swapKeyValue()
	return object
}

func (_this *structBuilder) BuildFromBigDecimalFloat(ctx *Context, value *apd.Decimal, _ reflect.Value) reflect.Value {
	_this.nextBuilderGenerator(ctx).BuildFromBigDecimalFloat(ctx, value, _this.nextValue)
	object := _this.nextValue
	_this.swapKeyValue()
	return object
}

func (_this *structBuilder) BuildFromUID(ctx *Context, value []byte, _ reflect.Value) reflect.Value {
	_this.nextBuilderGenerator(ctx).BuildFromUID(ctx, value, _this.nextValue)
	object := _this.nextValue
	_this.swapKeyValue()
	return object
}

func (_this *structBuilder) BuildFromArray(ctx *Context, arrayType events.ArrayType, value []byte, rv reflect.Value) reflect.Value {
	switch arrayType {
	case events.ArrayTypeString:
		if _this.nextIsKey {
			if ctx.config.Builder.CaseInsensitiveStructFieldNames {
				value = []byte(common.ToStructFieldIdentifier(string(value)))
			}

			if generatorDesc, ok := _this.generatorDescs[string(value)]; ok {
				_this.nextBuilderGenerator = generatorDesc.builderGenerator
				_this.nextValue = generatorDesc.field.GetField(_this.container)
			} else {
				ctx.StackBuilder(globalIgnoreBuilder)
				return rv
			}
		} else {
			_this.nextBuilderGenerator(ctx).BuildFromArray(ctx, arrayType, value, _this.nextValue)
		}
	default:
		_this.nextBuilderGenerator(ctx).BuildFromArray(ctx, arrayType, value, _this.nextValue)
	}
	object := _this.nextValue
	_this.swapKeyValue()
	return object
}

func (_this *structBuilder) BuildFromStringlikeArray(ctx *Context, arrayType events.ArrayType, value string, rv reflect.Value) reflect.Value {
	switch arrayType {
	case events.ArrayTypeString:
		if _this.nextIsKey {
			if ctx.config.Builder.CaseInsensitiveStructFieldNames {
				value = common.ToStructFieldIdentifier(value)
			}

			if generatorDesc, ok := _this.generatorDescs[value]; ok {
				_this.nextBuilderGenerator = generatorDesc.builderGenerator
				_this.nextValue = generatorDesc.field.GetField(_this.container)
			} else {
				ctx.StackBuilder(globalIgnoreBuilder)
				return rv
			}
		} else {
			_this.nextBuilderGenerator(ctx).BuildFromStringlikeArray(ctx, arrayType, value, _this.nextValue)
		}
	default:
		_this.nextBuilderGenerator(ctx).BuildFromStringlikeArray(ctx, arrayType, value, _this.nextValue)
	}
	object := _this.nextValue
	_this.swapKeyValue()
	return object
}

func (_this *structBuilder) BuildFromCustomBinary(ctx *Context, customType uint64, value []byte, dst reflect.Value) reflect.Value {
	_this.nextBuilderGenerator(ctx).BuildFromCustomBinary(ctx, customType, value, _this.nextValue)
	object := _this.nextValue
	_this.swapKeyValue()
	return object
}

func (_this *structBuilder) BuildFromCustomText(ctx *Context, customType uint64, value string, dst reflect.Value) reflect.Value {
	_this.nextBuilderGenerator(ctx).BuildFromCustomText(ctx, customType, value, _this.nextValue)
	object := _this.nextValue
	_this.swapKeyValue()
	return object
}

func (_this *structBuilder) BuildFromMedia(ctx *Context, mediaType string, data []byte, _ reflect.Value) reflect.Value {
	_this.nextBuilderGenerator(ctx).BuildFromMedia(ctx, mediaType, data, _this.nextValue)
	object := _this.nextValue
	_this.swapKeyValue()
	return object
}

func (_this *structBuilder) BuildFromTime(ctx *Context, value compact_time.Time, _ reflect.Value) reflect.Value {
	_this.nextBuilderGenerator(ctx).BuildFromTime(ctx, value, _this.nextValue)
	object := _this.nextValue
	_this.swapKeyValue()
	return object
}

func (_this *structBuilder) BuildNewList(ctx *Context) {
	_this.nextBuilderGenerator(ctx).BuildBeginListContents(ctx)
}

func (_this *structBuilder) BuildNewMap(ctx *Context) {
	_this.nextBuilderGenerator(ctx).BuildBeginMapContents(ctx)
}

func (_this *structBuilder) BuildNewNode(ctx *Context) {
	_this.nextBuilderGenerator(ctx).BuildBeginNodeContents(ctx)
}

func (_this *structBuilder) BuildNewEdge(ctx *Context) {
	_this.nextBuilderGenerator(ctx).BuildBeginEdgeContents(ctx)
}

func (_this *structBuilder) BuildEndContainer(ctx *Context) {
	object := _this.container
	ctx.UnstackBuilderAndNotifyChildFinished(object)
}

func (_this *structBuilder) BuildArtificiallyEndContainer(ctx *Context) {
	_this.BuildEndContainer(ctx)
}

func (_this *structBuilder) BuildBeginMapContents(ctx *Context) {
	ctx.StackBuilder(_this)
}

func (_this *structBuilder) BuildFromLocalReference(ctx *Context, id []byte) {
	nextValue := _this.nextValue
	_this.swapKeyValue()
	ctx.NotifyLocalReference(id, func(object reflect.Value) {
		setAnythingFromAnything(object, nextValue)
	})
}

func (_this *structBuilder) NotifyChildContainerFinished(ctx *Context, value reflect.Value) {
	_this.nextValue.Set(value)
	_this.swapKeyValue()
}
