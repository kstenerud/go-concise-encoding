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
	Options         options.BuilderOptions
	dstType         reflect.Type
	referenceFiller ReferenceFiller

	CustomBinaryBuildFunction  options.CustomBuildFunction
	CustomTextBuildFunction    options.CustomBuildFunction
	CurrentBuilder             Builder
	GetBuilderGeneratorForType func(dstType reflect.Type) BuilderGenerator
	builderStack               []Builder

	chunkedData             []byte
	chunkRemainingLength    uint64
	moreChunksFollow        bool
	arrayCompletionCallback func([]byte)
}

func (_this *Context) Init(opts *options.BuilderOptions,
	dstType reflect.Type,
	customBinaryBuildFunction options.CustomBuildFunction,
	customTextBuildFunction options.CustomBuildFunction,
	getBuilderGeneratorForType func(dstType reflect.Type) BuilderGenerator,
) {
	opts = opts.WithDefaultsApplied()
	_this.Options = *opts
	_this.dstType = dstType
	_this.CustomBinaryBuildFunction = customBinaryBuildFunction
	_this.CustomTextBuildFunction = customTextBuildFunction
	_this.GetBuilderGeneratorForType = getBuilderGeneratorForType
	_this.builderStack = make([]Builder, 0, 16)

	_this.referenceFiller.Init()
}

func (_this *Context) updateCurrentBuilder() {
	_this.CurrentBuilder = _this.builderStack[len(_this.builderStack)-1]
}

func (_this *Context) StackBuilder(builder Builder) {
	_this.builderStack = append(_this.builderStack, builder)
	_this.updateCurrentBuilder()
}

func (_this *Context) UnstackBuilder() Builder {
	oldTop := _this.CurrentBuilder
	_this.builderStack = _this.builderStack[:len(_this.builderStack)-1]
	_this.updateCurrentBuilder()
	return oldTop
}

func (_this *Context) UnstackBuilderAndNotifyChildFinished(container reflect.Value) Builder {
	oldTop := _this.UnstackBuilder()
	_this.CurrentBuilder.NotifyChildContainerFinished(_this, container)
	return oldTop
}

func (_this *Context) IgnoreNext() {
	_this.StackBuilder(globalIgnoreBuilder)
}

func (_this *Context) NANext() {
	_this.StackBuilder(globalNABuilder)
}

func (_this *Context) NotifyMarker(id interface{}, value reflect.Value) {
	_this.referenceFiller.NotifyMarker(id, value)
}
func (_this *Context) NotifyReference(lookingForID interface{}, valueSetter func(value reflect.Value)) {
	_this.referenceFiller.NotifyReference(lookingForID, valueSetter)
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

func (_this *Context) TryBuildFromCustom(builder Builder, arrayType events.ArrayType, value []byte, dst reflect.Value) bool {
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

func (_this *Context) BeginArray(arrayCompletionCallback func([]byte)) {
	_this.arrayCompletionCallback = arrayCompletionCallback
	_this.chunkedData = _this.chunkedData[:0]
}
func (_this *Context) BeginArrayChunk(length uint64, moreChunksFollow bool) {
	_this.chunkRemainingLength = length
	_this.moreChunksFollow = moreChunksFollow
	if !_this.moreChunksFollow && _this.chunkRemainingLength == 0 {
		_this.arrayCompletionCallback(_this.chunkedData)
		_this.chunkedData = _this.chunkedData[:0]
	}
}
func (_this *Context) AddArrayData(data []byte) {
	_this.chunkedData = append(_this.chunkedData, data...)
	_this.chunkRemainingLength -= uint64(len(data))
	if !_this.moreChunksFollow && _this.chunkRemainingLength == 0 {
		_this.arrayCompletionCallback(_this.chunkedData)
		_this.chunkedData = _this.chunkedData[:0]
	}
}
