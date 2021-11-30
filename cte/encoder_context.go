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

package cte

import (
	"github.com/kstenerud/go-concise-encoding/events"
	"github.com/kstenerud/go-concise-encoding/options"
)

type indenter struct {
	indent []byte
}

var g_indent = []byte("    ")

const g_indent_length = 4

func (_this *indenter) Reset() {
	_this.indent = _this.indent[:0]
}

func (_this *indenter) increase() {
	_this.indent = append(_this.indent, g_indent...)
}

func (_this *indenter) decrease() {
	_this.indent = _this.indent[:len(_this.indent)-4]
}

func (_this *indenter) GetOriginPos() int {
	return len(_this.indent) - g_indent_length
}

func (_this *indenter) GetOrigin() []byte {
	originPos := _this.GetOriginPos()
	if originPos >= 0 {
		return _this.indent[:originPos]
	}
	return _this.indent
}

func (_this *indenter) GetIndent() []byte {
	return g_indent
}

func (_this *indenter) GetOriginAndIndent() []byte {
	return _this.indent
}

type EncoderContext struct {
	opts                options.CTEEncoderOptions
	indenter            indenter
	stack               []EncoderDecorator
	Decorator           EncoderDecorator
	ContainerHasObjects bool
	Stream              Writer
	ArrayEngine         arrayEncoderEngine
}

func (_this *EncoderContext) Init(opts *options.CTEEncoderOptions) {
	_this.opts = *opts
	_this.ArrayEngine.Init(&_this.Stream, &_this.opts)
	_this.Stream.Init()
}

func (_this *EncoderContext) Begin() {
	_this.indenter.Reset()
	_this.stack = _this.stack[:0]
	_this.Decorator = nil
	_this.Stack(&topLevelDecorator)
}

func (_this *EncoderContext) BeforeValue() {
	_this.Decorator.BeforeValue(_this)
}

func (_this *EncoderContext) AfterValue() {
	_this.Decorator.AfterValue(_this)
	_this.ContainerHasObjects = true
}

func (_this *EncoderContext) BeforeComment() {
	_this.Decorator.BeforeComment(_this)
}

func (_this *EncoderContext) AfterComment(isMultiline bool) {
	_this.Decorator.AfterComment(_this)
	_this.ContainerHasObjects = true
}

func (_this *EncoderContext) BeginContainer() {
	_this.ContainerHasObjects = false
}

func (_this *EncoderContext) EndContainer() {
	_this.Decorator.EndContainer(_this)
}

func (_this *EncoderContext) Stack(decorator EncoderDecorator) {
	_this.stack = append(_this.stack, decorator)
	_this.Decorator = decorator
}

func (_this *EncoderContext) Unstack() {
	_this.stack = _this.stack[:len(_this.stack)-1]
	_this.Decorator = _this.stack[len(_this.stack)-1]
}

func (_this *EncoderContext) Switch(decorator EncoderDecorator) {
	_this.stack[len(_this.stack)-1] = decorator
	_this.Decorator = decorator
}

func (_this *EncoderContext) IsAtOrigin() bool {
	return _this.Stream.Column == _this.indenter.GetOriginPos()
}

func (_this *EncoderContext) Indent() {
	_this.indenter.increase()
}

func (_this *EncoderContext) Unindent() {
	_this.indenter.decrease()
}

func (_this *EncoderContext) WriteReturnToOrigin() {
	if !_this.IsAtOrigin() {
		_this.Stream.WriteLF()
		_this.Stream.WriteBytesNotLF(_this.indenter.GetOrigin())
	}
}

func (_this *EncoderContext) WriteIndentIfOrigin() {
	if _this.IsAtOrigin() {
		_this.Stream.WriteBytesNotLF(_this.indenter.GetIndent())
	}
}

func (_this *EncoderContext) WriteIndentOrSpace() {
	if _this.IsAtOrigin() {
		_this.Stream.WriteBytesNotLF(_this.indenter.GetIndent())
	} else {
		_this.Stream.WriteByteNotLF(' ')
	}
}

func (_this *EncoderContext) WriteNewlineAndOrigin() {
	_this.Stream.WriteLF()
	_this.Stream.WriteBytesNotLF(_this.indenter.GetOrigin())
}

func (_this *EncoderContext) WriteIndent() {
	_this.Stream.WriteBytesNotLF(_this.indenter.GetIndent())
}

func (_this *EncoderContext) WriteNewlineAndOriginAndIndent() {
	_this.Stream.WriteLF()
	_this.Stream.WriteBytesNotLF(_this.indenter.GetOriginAndIndent())
}

func (_this *EncoderContext) WriteElementSeparator() {
	_this.Stream.WriteByteNotLF(' ')
}

func (_this *EncoderContext) WriteIdentifier(data []byte) {
	_this.Stream.WriteBytesNotLF(data)
}

func (_this *EncoderContext) EncodeStringlikeArray(arrayType events.ArrayType, data string) {
	_this.ArrayEngine.EncodeStringlikeArray(_this.Decorator.GetStringContext(), arrayType, data)
}

func (_this *EncoderContext) EncodeArray(arrayType events.ArrayType, elementCount uint64, data []uint8) {
	_this.ArrayEngine.EncodeArray(_this.Decorator.GetStringContext(), arrayType, elementCount, data)
}

func (_this *EncoderContext) BeginArray(arrayType events.ArrayType, completion func()) {
	finalCompletion := completion
	switch arrayType {
	case events.ArrayTypeCustomText, events.ArrayTypeRemoteRef, events.ArrayTypeResourceID, events.ArrayTypeString:
		// Do nothing
	default:
		_this.Stack(nonStringArrayDecorator)
		finalCompletion = func() {
			_this.Unstack()
			completion()
		}
	}

	_this.ArrayEngine.BeginArray(_this.Decorator.GetStringContext(), arrayType, finalCompletion)
}

func (_this *EncoderContext) BeginArrayChunk(elementCount uint64, moreChunksFollow bool) {
	_this.ArrayEngine.BeginChunk(elementCount, moreChunksFollow)
}

func (_this *EncoderContext) EncodeArrayData(data []byte) {
	_this.ArrayEngine.AddArrayData(data)
}

func (_this *EncoderContext) WriteCommentString(data string) {
	// TODO: Not this
	_this.WriteCommentStringData([]byte(data))
}

func (_this *EncoderContext) WriteCommentStringData(data []uint8) {
	// TODO: Need anything else?
	_this.Stream.WriteBytesPossibleLF(data)
}

func (_this *EncoderContext) WriteMarkupContentString(data string) {
	// TODO: Not this
	_this.WriteMarkupContentStringData([]byte(data))
}

func (_this *EncoderContext) WriteMarkupContentStringData(data []uint8) {
	_this.Stream.WritePotentiallyEscapedMarkupContentsBytes(data)
}
