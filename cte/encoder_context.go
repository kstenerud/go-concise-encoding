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

func (_this *indenter) Reset() {
	_this.indent = _this.indent[:0]
	_this.indent = append(_this.indent, '\n')
}

func (_this *indenter) increase() {
	_this.indent = append(_this.indent, "    "...)
}

func (_this *indenter) decrease() {
	_this.indent = _this.indent[:len(_this.indent)-4]
}

func (_this *indenter) Get() []byte {
	return _this.indent
}

type EncoderContext struct {
	opts                options.CTEEncoderOptions
	indenter            indenter
	encoderStack        []Encoder
	CurrentEncoder      Encoder
	ContainerHasObjects bool
	currentPrefix       []byte
	Stream              Writer
	ArrayEngine         arrayEncoderEngine
}

func (_this *EncoderContext) Init(opts *options.CTEEncoderOptions) {
	_this.opts = *opts
	_this.ArrayEngine.Init(&_this.Stream, &_this.opts)
	_this.Stream.Init()
	_this.Reset()
}

func (_this *EncoderContext) Reset() {
	_this.indenter.Reset()
	_this.encoderStack = _this.encoderStack[:0]
	_this.CurrentEncoder = nil
	_this.Stack(&globalTopLevelEncoder)
	_this.SetStandardIndentPrefix()
}

func (_this *EncoderContext) Stack(encoder Encoder) {
	_this.encoderStack = append(_this.encoderStack, encoder)
	_this.CurrentEncoder = encoder
}

func (_this *EncoderContext) Unstack() {
	_this.encoderStack = _this.encoderStack[:len(_this.encoderStack)-1]
	_this.CurrentEncoder = _this.encoderStack[len(_this.encoderStack)-1]
}

func (_this *EncoderContext) ChangeEncoder(encoder Encoder) {
	_this.encoderStack[len(_this.encoderStack)-1] = encoder
	_this.CurrentEncoder = encoder
}

func (_this *EncoderContext) IncreaseIndent() {
	_this.indenter.increase()
}

func (_this *EncoderContext) DecreaseIndent() {
	_this.indenter.decrease()
}

func (_this *EncoderContext) WriteBasicIndent() {
	_this.Stream.WriteBytes(_this.indenter.Get())
}

func (_this *EncoderContext) SetIndentPrefix(value []byte) {
	_this.currentPrefix = value
}

func (_this *EncoderContext) SetStandardIndentPrefix() {
	_this.SetIndentPrefix(_this.indenter.Get())
}

func (_this *EncoderContext) SetStandardMapKeyPrefix() {
	_this.SetIndentPrefix(_this.indenter.Get())
}

func (_this *EncoderContext) SetStandardMapValuePrefix() {
	_this.SetIndentPrefix(mapValuePrefix)
}

func (_this *EncoderContext) SetMarkupAttributeKeyPrefix() {
	_this.SetIndentPrefix(markupAttributeKeyPrefix)
}

func (_this *EncoderContext) SetMarkupAttributeValuePrefix() {
	_this.SetIndentPrefix(markupAttributeValuePrefix)
}

func (_this *EncoderContext) ClearPrefix() {
	_this.SetIndentPrefix(emptyPrefix)
}

func (_this *EncoderContext) WriteCurrentPrefix() {
	_this.Stream.WriteBytes(_this.currentPrefix)
}

func (_this *EncoderContext) BeginStandardList() {
	_this.Stack(&globalListEncoder)
	_this.CurrentEncoder.Begin(_this)
}

func (_this *EncoderContext) BeginStandardMap() {
	_this.Stack(&globalMapKeyEncoder)
	_this.CurrentEncoder.Begin(_this)
}

func (_this *EncoderContext) BeginStandardMarkup() {
	_this.Stack(&globalMarkupNameEncoder)
	_this.CurrentEncoder.Begin(_this)
}

func (_this *EncoderContext) SwitchToMarkupAttributes() {
	_this.ChangeEncoder(&globalMarkupKeyEncoder)
}

func (_this *EncoderContext) SwitchToMarkupContents() {
	_this.ChangeEncoder(&globalMarkupContentsEncoder)
}

func (_this *EncoderContext) BeginStandardComment() {
	_this.Stack(&globalCommentEncoder)
	_this.CurrentEncoder.Begin(_this)
}

func (_this *EncoderContext) BeginStandardMarker() {
	_this.Stack(&globalMarkerIDEncoder)
	_this.CurrentEncoder.Begin(_this)
}

func (_this *EncoderContext) BeginStandardReference() {
	_this.Stack(&globalReferenceEncoder)
	_this.CurrentEncoder.Begin(_this)
}

func (_this *EncoderContext) BeginStandardConstant(name []byte, explicitValue bool) {
	_this.Stream.WriteByte('#')
	_this.Stream.WriteBytes(name)
	_this.Stack(&globalConstantEncoder)
}

func (_this *EncoderContext) BeginNA() {
	_this.Stream.WriteNA()
	_this.Stream.WriteConcat()
	_this.Stack(&globalPostInvisibleEncoder)
}

func (_this *EncoderContext) BeginStandardArray(arrayType events.ArrayType) {
	if arrayType == events.ArrayTypeResourceIDConcat {
		_this.Stack(&globalPostStreamRIDCatEncoder)
	}
	_this.Stack((&globalArrayEncoder))
	_this.ArrayEngine.BeginArray(arrayType, func() {
		_this.Unstack()
		_this.CurrentEncoder.ChildContainerFinished(_this, true)
	})
}

func (_this *EncoderContext) BeginPotentialRIDCat(arrayType events.ArrayType) {
	if arrayType == events.ArrayTypeResourceIDConcat {
		_this.Stack(&globalPostRIDCatEncoder)
	}
}

func (_this *EncoderContext) WriteStringlikeArray(arrayType events.ArrayType, data string) {
	// TODO: avoid string-to-bytes conversion?
	_this.ArrayEngine.EncodeArray(arrayType, uint64(len(data)), []byte(data))
	_this.BeginPotentialRIDCat(arrayType)
}

func (_this *EncoderContext) WriteArray(arrayType events.ArrayType, elementCount uint64, data []uint8) {
	_this.ArrayEngine.EncodeArray(arrayType, elementCount, data)
	_this.BeginPotentialRIDCat(arrayType)
}

func (_this *EncoderContext) WriteCommentString(data string) {
	// TODO: Not this
	_this.WriteMarkupContentStringData([]byte(data))
}

func (_this *EncoderContext) WriteCommentStringData(data []uint8) {
	// TODO: Need anything else?
	_this.Stream.WriteBytes(data)
}

func (_this *EncoderContext) WriteMarkupContentString(data string) {
	// TODO: Not this
	_this.WriteMarkupContentStringData([]byte(data))
}

func (_this *EncoderContext) WriteMarkupContentStringData(data []uint8) {
	_this.Stream.WritePotentiallyEscapedMarkupContents(data)
}

// ============================================================================

var (
	emptyPrefix                = []byte{}
	mapValuePrefix             = []byte{' ', '=', ' '}
	markupAttributeKeyPrefix   = []byte{' '}
	markupAttributeValuePrefix = []byte{'='}
)
