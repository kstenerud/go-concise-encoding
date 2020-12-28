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

// Builders consume events to produce objects.
//
// Builders respond to builder events in order to build arbitrary objects.
// Generally, they take template types and generate objects of those types.
package builder

import (
	"reflect"

	"github.com/kstenerud/go-concise-encoding/events"

	"github.com/kstenerud/go-concise-encoding/options"
)

type Context struct {
	Options                   options.BuilderOptions
	CustomBinaryBuildFunction options.CustomBuildFunction
	CustomTextBuildFunction   options.CustomBuildFunction
	NotifyMarker              func(id interface{}, value reflect.Value)
	NotifyReference           func(lookingForID interface{}, valueSetter func(value reflect.Value))
	CurrentBuilder            ObjectBuilder
	builderStack              []ObjectBuilder
}

func context(opts *options.BuilderOptions,
	customBinaryBuildFunction options.CustomBuildFunction,
	customTextBuildFunction options.CustomBuildFunction,
	notifyMarker func(id interface{}, value reflect.Value),
	notifyReference func(lookingForID interface{}, valueSetter func(value reflect.Value)),
) Context {
	opts = opts.WithDefaultsApplied()
	return Context{
		Options:                   *opts,
		CustomBinaryBuildFunction: customBinaryBuildFunction,
		CustomTextBuildFunction:   customTextBuildFunction,
		NotifyMarker:              notifyMarker,
		NotifyReference:           notifyReference,
		builderStack:              make([]ObjectBuilder, 0, 16),
	}
}

func (_this *Context) updateCurrentBuilder() {
	_this.CurrentBuilder = _this.builderStack[len(_this.builderStack)-1]
}

func (_this *Context) StackBuilder(builder ObjectBuilder) {
	_this.builderStack = append(_this.builderStack, builder)
	_this.updateCurrentBuilder()
}

func (_this *Context) UnstackBuilder() ObjectBuilder {
	oldTop := _this.CurrentBuilder
	_this.builderStack = _this.builderStack[:len(_this.builderStack)-1]
	_this.updateCurrentBuilder()
	return oldTop
}

func (_this *Context) UnstackBuilderAndNotifyChildFinished(container reflect.Value) ObjectBuilder {
	oldTop := _this.UnstackBuilder()
	_this.CurrentBuilder.NotifyChildContainerFinished(_this, container)
	return oldTop
}

func (_this *Context) BeginMarkerObject(id interface{}) {
	_this.UnstackBuilder()
	marker := newMarkerObjectBuilder(id, _this.CurrentBuilder)
	_this.StackBuilder(marker)
}

func (_this *Context) StoreReferencedObject(id interface{}) {
	_this.UnstackBuilder()
	_this.CurrentBuilder.BuildFromReference(_this, id)
}

func (_this *Context) TryBuildFromCustom(builder ObjectBuilder, arrayType events.ArrayType, value []byte, dst reflect.Value) bool {
	switch arrayType {
	case events.ArrayTypeCustomBinary:
		if err := _this.CustomBinaryBuildFunction(value, dst); err != nil {
			PanicBuildFromCustomBinary(builder, value, dst.Type(), err)
		}
		return true
	case events.ArrayTypeCustomText:
		if err := _this.CustomTextBuildFunction(value, dst); err != nil {
			PanicBuildFromCustomText(builder, value, dst.Type(), err)
		}
		return true
	default:
		return false
	}
}

type BuilderGenerator func() ObjectBuilder
type BuilderGeneratorGetter func(reflect.Type) BuilderGenerator
