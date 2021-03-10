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
	"fmt"

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
	currentPrefix       string
	Stream              EncodeBuffer
	ArrayEngine         arrayEncoderEngine
}

func (_this *EncoderContext) Init(opts *options.CTEEncoderOptions) {
	_this.opts = *opts
	_this.ArrayEngine.Init(&_this.Stream, &_this.opts)
	_this.Reset()
}

func (_this *EncoderContext) Reset() {
	_this.indenter.Reset()
	_this.Stream.Reset()
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
	_this.Stream.AddBytes(_this.indenter.Get())
}

func (_this *EncoderContext) SetIndentPrefix(value string) {
	_this.currentPrefix = value
}

func (_this *EncoderContext) SetStandardIndentPrefix() {
	_this.SetIndentPrefix(string(_this.indenter.Get()))
}

func (_this *EncoderContext) SetStandardMapKeyPrefix() {
	_this.SetIndentPrefix(string(_this.indenter.Get()))
}

func (_this *EncoderContext) SetStandardMapValuePrefix() {
	_this.SetIndentPrefix(" = ")
}

func (_this *EncoderContext) SetMarkupAttributeKeyPrefix() {
	_this.SetIndentPrefix(" ")
}

func (_this *EncoderContext) SetMarkupAttributeValuePrefix() {
	_this.SetIndentPrefix("=")
}

func (_this *EncoderContext) ClearPrefix() {
	_this.SetIndentPrefix("")
}

func (_this *EncoderContext) WriteCurrentPrefix() {
	_this.Stream.AddString(_this.currentPrefix)
	// TODO: Need to do this?
	// _this.ClearMapPrefix()
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

func (_this *EncoderContext) BeginStandardMetadata() {
	_this.Stack(&globalMetadataKeyEncoder)
	_this.CurrentEncoder.Begin(_this)
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

func (_this *EncoderContext) BeginStandardConcatenate() {
	panic(fmt.Errorf("TODO: EncoderContext.BeginConcatenate"))
}

func (_this *EncoderContext) BeginStandardConstant(name []byte, explicitValue bool) {
	_this.Stream.AddByte('#')
	_this.Stream.AddBytes(name)
	_this.Stack(&globalConstantEncoder)
}

func (_this *EncoderContext) BeginConcatenatedNA() {
	panic("TODO: EncoderContext.BeginConcatenatedNA")
	// _this.Stack(&globalNACatEncoder)
	// _this.CurrentEncoder.Begin(_this)
}

func (_this *EncoderContext) BeginStandardArray(arrayType events.ArrayType) {
	_this.Stack((&globalArrayEncoder))
	_this.ArrayEngine.BeginArray(arrayType, func() {
		_this.Unstack()
		_this.CurrentEncoder.ChildContainerFinished(_this, true)
	})
}

// pre-write (indent)
// post-write (lf?)
// list-type
// map-type (key section, value section)
// array-type
// metadata follow
// comment follow

// string in comment
// string in markup contents
// string as ID
// int as ID
// constant name
// custom string
// NA stuff

// Types that can be printed differently:
// - string (quoted, unquoted) (escape/noescape) (trim/notrim)
// - int (pos or neg) (bin, oct, dec, hex) (with/without prefix)
// - binary float (dec, hex) (with/without prefix)
