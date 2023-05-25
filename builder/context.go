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

	"github.com/kstenerud/go-concise-encoding/configuration"
)

type recordTypeKey func(*Context, Builder)

var unusedValue reflect.Value

type Context struct {
	config          *configuration.BuilderConfiguration
	dstType         reflect.Type
	referenceFiller ReferenceFiller

	CustomBinaryBuildFunction  configuration.CustomBinaryBuildFunction
	CustomTextBuildFunction    configuration.CustomTextBuildFunction
	CurrentBuilder             Builder
	GetBuilderGeneratorForType func(dstType reflect.Type) BuilderGenerator
	builderStack               []Builder

	chunkedData             []byte
	chunkRemainingLength    uint64
	moreChunksFollow        bool
	arrayCompletionCallback func(*Context)

	recordTypes    map[string][]recordTypeKey
	recordType     []recordTypeKey
	recordTypeName string
}

func (_this *Context) Init(config *configuration.BuilderConfiguration,
	dstType reflect.Type,
	customBinaryBuildFunction configuration.CustomBinaryBuildFunction,
	customTextBuildFunction configuration.CustomTextBuildFunction,
	getBuilderGeneratorForType func(dstType reflect.Type) BuilderGenerator,
) {
	if config == nil {
		defaultConfig := configuration.DefaultBuilderConfiguration()
		config = &defaultConfig
	} else {
		config.ApplyDefaults()
	}
	_this.config = config
	_this.dstType = dstType
	_this.CustomBinaryBuildFunction = customBinaryBuildFunction
	_this.CustomTextBuildFunction = customTextBuildFunction
	_this.GetBuilderGeneratorForType = getBuilderGeneratorForType
	_this.builderStack = make([]Builder, 0, 16)
	_this.recordTypes = make(map[string][]recordTypeKey)

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

func (_this *Context) SwapBuilder(builder Builder) Builder {
	oldTop := _this.CurrentBuilder
	_this.builderStack = _this.builderStack[:len(_this.builderStack)-1]
	_this.builderStack = append(_this.builderStack, builder)
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

func (_this *Context) BeginRecordType(id []byte) {
	_this.StackBuilder(generateRecordTypeBuilder(_this))
	_this.recordTypeName = string(id)
	_this.recordType = _this.recordType[:0]
}

func (_this *Context) BeginRecord(id []byte) {
	keys := _this.recordTypes[string(id)]
	_this.CurrentBuilder.BuildNewMap(_this)
	builder := _this.UnstackBuilder()
	_this.StackBuilder(generateRecordBuilder(_this, keys, builder))
}

func (_this *Context) AddRecordTypeKey(key recordTypeKey) {
	_this.recordType = append(_this.recordType, key)
}

func (_this *Context) EndRecordType() {
	_this.recordTypes[_this.recordTypeName] = _this.recordType
	_this.UnstackBuilder()
}

func (_this *Context) NotifyMarker(id []byte, value reflect.Value) {
	_this.referenceFiller.NotifyMarker(id, value)
}
func (_this *Context) NotifyLocalReference(lookingForID []byte, valueSetter func(value reflect.Value)) {
	_this.referenceFiller.NotifyLocalReference(lookingForID, valueSetter)
}

func (_this *Context) BeginMarkerObject(id []byte) {
	marker := newMarkerObjectBuilder(id, _this.CurrentBuilder)
	_this.StackBuilder(marker)
}

func (_this *Context) TryBuildFromCustomBinary(builder Builder, customType uint64, data []byte, dst reflect.Value) {
	if err := _this.CustomBinaryBuildFunction(customType, data, dst); err != nil {
		PanicBuildFromCustomBinary(builder, customType, data, dst.Type(), err)
	}
}

func (_this *Context) TryBuildFromCustomText(builder Builder, customType uint64, data string, dst reflect.Value) {
	if err := _this.CustomTextBuildFunction(customType, data, dst); err != nil {
		PanicBuildFromCustomText(builder, customType, data, dst.Type(), err)
	}
}

func (_this *Context) BeginArray(arrayCompletionCallback func(*Context)) {
	_this.arrayCompletionCallback = arrayCompletionCallback
	_this.chunkedData = _this.chunkedData[:0]
}
func (_this *Context) ContinueMultiComponentArray(arrayCompletionCallback func(*Context)) {
	_this.arrayCompletionCallback = arrayCompletionCallback
}
func (_this *Context) BeginArrayChunk(length uint64, moreChunksFollow bool) {
	_this.chunkRemainingLength = length
	_this.moreChunksFollow = moreChunksFollow
	if !_this.moreChunksFollow && _this.chunkRemainingLength == 0 {
		_this.arrayCompletionCallback(_this)
	}
}
func (_this *Context) AddArrayData(data []byte) {
	_this.chunkedData = append(_this.chunkedData, data...)
	_this.chunkRemainingLength -= uint64(len(data))
	if !_this.moreChunksFollow && _this.chunkRemainingLength == 0 {
		_this.arrayCompletionCallback(_this)
	}
}
