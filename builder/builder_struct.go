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
	"strings"
	"time"

	"github.com/kstenerud/go-concise-encoding/events"
	"github.com/kstenerud/go-concise-encoding/internal/common"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
)

type structBuilderField struct {
	Name      string
	Index     int
	Omit      bool
	OmitEmpty bool
	OmitValue string
}

func (_this *structBuilderField) applyTags(tags string) {
	if tags == "" {
		return
	}

	requiresValue := func(kv []string, key string) {
		if len(kv) != 2 {
			panic(fmt.Errorf(`tag key "%s" requires a value`, key))
		}
	}

	for _, entry := range strings.Split(tags, ",") {
		kv := strings.Split(entry, "=")
		switch strings.TrimSpace(kv[0]) {
		// TODO: lossy/nolossy
		// TODO: lowercase/origcase
		case "-":
			_this.Omit = true
		case "omit":
			if len(kv) == 1 {
				_this.Omit = true
			} else {
				_this.OmitValue = strings.TrimSpace(kv[1])
			}
		case "omitempty":
			// TODO: Implement omitempty
			_this.OmitEmpty = true
		case "name":
			requiresValue(kv, "name")
			_this.Name = strings.TrimSpace(kv[1])
		default:
			panic(fmt.Errorf("%v: Unknown Concise Encoding struct tag field", entry))
		}
	}
}

type structBuilder struct {
	dstType                reflect.Type
	generatorDescs         map[string]*structBuilderGeneratorDesc
	nameBuilderGenerator   BuilderGenerator
	ignoreBuilderGenerator BuilderGenerator
	nextBuilderGenerator   BuilderGenerator
	container              reflect.Value
	nextValue              reflect.Value
	nextIsKey              bool
	nextIsIgnored          bool
}

type structBuilderGeneratorDesc struct {
	field            *structBuilderField
	builderGenerator BuilderGenerator
}

func newStructBuilderGenerator(getBuilderGeneratorForType BuilderGeneratorGetter, dstType reflect.Type) BuilderGenerator {
	nameBuilderGenerator := getBuilderGeneratorForType(reflect.TypeOf(""))
	ignoreBuilderGenerator := generateIgnoreBuilder
	generatorDescs := make(map[string]*structBuilderGeneratorDesc)

	for i := 0; i < dstType.NumField(); i++ {
		reflectField := dstType.Field(i)
		if reflectField.PkgPath == "" {
			builderGenerator := getBuilderGeneratorForType(reflectField.Type)
			structField := &structBuilderField{
				Name:  reflectField.Name,
				Index: i,
			}
			structField.applyTags(reflectField.Tag.Get("ce"))
			generatorDescs[structField.Name] = &structBuilderGeneratorDesc{
				field:            structField,
				builderGenerator: builderGenerator,
			}
		}
	}

	// Make lowercase mappings as well in case we later do case-insensitive field name matching
	for _, desc := range generatorDescs {
		lowerName := common.ASCIIToLower(desc.field.Name)
		if _, exists := generatorDescs[lowerName]; !exists {
			generatorDescs[lowerName] = desc
		}
	}

	return func(ctx *Context) ObjectBuilder {
		builder := &structBuilder{
			dstType:                dstType,
			generatorDescs:         generatorDescs,
			nameBuilderGenerator:   nameBuilderGenerator,
			ignoreBuilderGenerator: ignoreBuilderGenerator,
		}
		builder.reset()
		return builder
	}
}

func (_this *structBuilder) String() string {
	return fmt.Sprintf("%v<%v>", reflect.TypeOf(_this), _this.dstType)
}

func (_this *structBuilder) reset() {
	_this.nextBuilderGenerator = _this.nameBuilderGenerator
	_this.container = reflect.New(_this.dstType).Elem()
	_this.nextValue = reflect.Value{}
	_this.nextIsKey = true
	_this.nextIsIgnored = false
}

func (_this *structBuilder) swapKeyValue() {
	_this.nextIsKey = !_this.nextIsKey
}

func (_this *structBuilder) BuildFromNil(ctx *Context, _ reflect.Value) reflect.Value {
	_this.nextBuilderGenerator(ctx).BuildFromNil(ctx, _this.nextValue)
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

func (_this *structBuilder) BuildFromUUID(ctx *Context, value []byte, _ reflect.Value) reflect.Value {
	_this.nextBuilderGenerator(ctx).BuildFromUUID(ctx, value, _this.nextValue)
	object := _this.nextValue
	_this.swapKeyValue()
	return object
}

func (_this *structBuilder) BuildFromArray(ctx *Context, arrayType events.ArrayType, value []byte, _ reflect.Value) reflect.Value {
	switch arrayType {
	case events.ArrayTypeString:
		if _this.nextIsKey {
			if ctx.Options.CaseInsensitiveStructFieldNames {
				common.ASCIIBytesToLower(value)
			}
			name := string(value)

			if generatorDesc, ok := _this.generatorDescs[name]; ok {
				_this.nextBuilderGenerator = generatorDesc.builderGenerator
				_this.nextValue = _this.container.Field(generatorDesc.field.Index)
			} else {
				_this.nextBuilderGenerator = _this.ignoreBuilderGenerator
				_this.nextIsIgnored = true
				break
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

func (_this *structBuilder) BuildFromTime(ctx *Context, value time.Time, _ reflect.Value) reflect.Value {
	_this.nextBuilderGenerator(ctx).BuildFromTime(ctx, value, _this.nextValue)
	object := _this.nextValue
	_this.swapKeyValue()
	return object
}

func (_this *structBuilder) BuildFromCompactTime(ctx *Context, value *compact_time.Time, _ reflect.Value) reflect.Value {
	_this.nextBuilderGenerator(ctx).BuildFromCompactTime(ctx, value, _this.nextValue)
	object := _this.nextValue
	_this.swapKeyValue()
	return object
}

func (_this *structBuilder) BuildInitiateList(ctx *Context) {
	_this.nextBuilderGenerator(ctx).BuildBeginListContents(ctx)
}

func (_this *structBuilder) BuildInitiateMap(ctx *Context) {
	_this.nextBuilderGenerator(ctx).BuildBeginMapContents(ctx)
}

func (_this *structBuilder) BuildEndContainer(ctx *Context) {
	object := _this.container
	_this.reset()
	ctx.UnstackBuilderAndNotifyChildFinished(object)
}

func (_this *structBuilder) BuildBeginMapContents(ctx *Context) {
	ctx.StackBuilder(_this)
}

func (_this *structBuilder) BuildFromReference(ctx *Context, id interface{}) {
	nextValue := _this.nextValue
	_this.swapKeyValue()
	ctx.NotifyReference(id, func(object reflect.Value) {
		setAnythingFromAnything(object, nextValue)
	})
}

func (_this *structBuilder) NotifyChildContainerFinished(ctx *Context, value reflect.Value) {
	if _this.nextIsIgnored {
		_this.nextIsIgnored = false
		return
	}

	_this.nextValue.Set(value)
	_this.swapKeyValue()
}
