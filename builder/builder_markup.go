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
	"reflect"

	"github.com/kstenerud/go-concise-encoding/events"
	"github.com/kstenerud/go-concise-encoding/types"
)

type markupBuilder struct {
	name string
}

func generateMarkupBuilder(ctx *Context) Builder { return &markupBuilder{} }

func (_this *markupBuilder) String() string {
	return fmt.Sprintf("%v", reflect.TypeOf(_this))
}

func (_this *markupBuilder) BuildBeginMarkupContents(ctx *Context, name []byte) {
	_this.name = string(name)
	ctx.StackBuilder(_this)
	interfaceMapBuilderGenerator(ctx).BuildBeginMapContents(ctx)
}

func (_this *markupBuilder) NotifyChildContainerFinished(ctx *Context, value reflect.Value) {
	ctx.SwapBuilder(newMarkupContentsBuilder(_this.name, value))
}

type markupContentsBuilder struct {
	markup *types.Markup
}

func newMarkupContentsBuilder(name string, attributes reflect.Value) *markupContentsBuilder {
	return &markupContentsBuilder{
		markup: &types.Markup{
			Name:       name,
			Attributes: attributes.Interface().(map[interface{}]interface{}),
		},
	}
}

func (_this *markupContentsBuilder) BuildFromArray(ctx *Context, arrayType events.ArrayType, value []byte, _ reflect.Value) reflect.Value {
	str := string(value)
	_this.markup.AddString(str)
	return reflect.ValueOf(str)
}

func (_this *markupContentsBuilder) BuildFromStringlikeArray(ctx *Context, arrayType events.ArrayType, value string, _ reflect.Value) reflect.Value {
	_this.markup.AddString(value)
	return reflect.ValueOf(value)
}

func (_this *markupContentsBuilder) BuildNewMarkup(ctx *Context, name []byte) {
	generator := ctx.GetBuilderGeneratorForType(reflect.TypeOf((*types.Markup)(nil)))
	generator(ctx).BuildBeginMarkupContents(ctx, name)
}

func (_this *markupContentsBuilder) NotifyChildContainerFinished(ctx *Context, value reflect.Value) {
	v := value.Interface().(*types.Markup)
	_this.markup.AddMarkup(v)
}

func (_this *markupContentsBuilder) BuildEndContainer(ctx *Context) {
	ctx.UnstackBuilderAndNotifyChildFinished(reflect.ValueOf(_this.markup).Elem())
}
