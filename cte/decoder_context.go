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
	"io"

	"github.com/kstenerud/go-concise-encoding/events"
	"github.com/kstenerud/go-concise-encoding/options"
)

type ContainerType uint16

const (
	ContainerTypeUnset ContainerType = iota
	ContainerTypeList
	ContainerTypeMap
	ContainerTypeNodeOrEdgeOrStructInstance
	ContainerTypeStructTemplate
	ContainerTypeEND
)

func (_this ContainerType) String() string {
	if _this >= ContainerTypeEND {
		panic(fmt.Errorf("BUG: ContainerType out of range: %x", uint(_this)))
	}
	switch _this {
	case ContainerTypeUnset:
		return "unset"
	case ContainerTypeList:
		return "list"
	case ContainerTypeMap:
		return "map"
	case ContainerTypeNodeOrEdgeOrStructInstance:
		return "node, edge, struct instance"
	case ContainerTypeStructTemplate:
		return "struct template"
	default:
		panic(fmt.Errorf("BUG: ContainerType.String(): Unknown container type %d (%x)", uint(_this), uint(_this)))
	}
}

type DecoderStackEntry struct {
	DecoderFunc   DecoderOp
	ContainerType ContainerType
}

type DecoderContext struct {
	opts                 *options.CEDecoderOptions
	Stream               Reader
	TextPos              *TextPositionCounter
	EventReceiver        events.DataEventReceiver
	stack                []DecoderStackEntry
	awaitingStructuralWS bool
	IsDocumentComplete   bool

	ArrayContainsComments bool
	ArrayType             events.ArrayType
	ArrayBytesPerElement  int
	ArrayDigitType        string
	Scratch               []byte
}

func (_this *DecoderContext) BeginArray(digitType string, arrayType events.ArrayType, elementWidth int) {
	_this.Scratch = _this.Scratch[:0]
	_this.ArrayDigitType = digitType
	_this.ArrayType = arrayType
	_this.ArrayBytesPerElement = elementWidth
	_this.ArrayContainsComments = false
}

func (_this *DecoderContext) BeginContainer(decoder DecoderOp, containerType ContainerType) {
	_this.stack = append(_this.stack, DecoderStackEntry{
		DecoderFunc:   decoder,
		ContainerType: containerType,
	})
}

func (_this *DecoderContext) EndContainer(allowedType ContainerType) {
	entry := _this.topOfStack()
	if entry.ContainerType != allowedType {
		_this.Errorf("Container type %v cannot be closed by container for type (%v)", entry.ContainerType, allowedType)
	}

	_this.UnstackDecoder()
}

func (_this *DecoderContext) topOfStack() *DecoderStackEntry {
	return &_this.stack[len(_this.stack)-1]
}

func (_this *DecoderContext) Init(opts *options.CEDecoderOptions, reader io.Reader, eventReceiver events.DataEventReceiver) {
	_this.opts = opts
	_this.Stream.Init(reader)
	_this.TextPos = &_this.Stream.TextPos
	_this.EventReceiver = eventReceiver
	if cap(_this.stack) > 0 {
		_this.stack = _this.stack[:0]
	} else {
		_this.stack = make([]DecoderStackEntry, 0, 16)
	}
	_this.IsDocumentComplete = false
}

func (_this *DecoderContext) SetEventReceiver(eventReceiver events.DataEventReceiver) {
	_this.EventReceiver = eventReceiver
}

func (_this *DecoderContext) AssertHasStructuralWS() {
	if _this.awaitingStructuralWS {
		_this.Errorf("Expected structural whitespace")
	}
}

func (_this *DecoderContext) AwaitStructuralWS() {
	_this.awaitingStructuralWS = true
}

func (_this *DecoderContext) NoNeedForWS() {
	_this.awaitingStructuralWS = false
}

func (_this *DecoderContext) NotifyStructuralWS() {
	_this.awaitingStructuralWS = false
}

func (_this *DecoderContext) DecodeNext() {
	_this.topOfStack().DecoderFunc(_this)
}

func (_this *DecoderContext) ChangeDecoder(decoder DecoderOp) {
	_this.topOfStack().DecoderFunc = decoder
}

func (_this *DecoderContext) StackDecoder(decoder DecoderOp) {
	_this.stack = append(_this.stack, DecoderStackEntry{
		DecoderFunc: decoder,
	})
}

func (_this *DecoderContext) UnstackDecoder() {
	_this.stack = _this.stack[:len(_this.stack)-1]
}

func (_this *DecoderContext) SetContainerType(containerType ContainerType) {
	_this.topOfStack().ContainerType = containerType
}

func (_this *DecoderContext) Errorf(format string, args ...interface{}) {
	_this.Stream.TextPos.Errorf(format, args...)
}

func (_this *DecoderContext) UnexpectedChar(decoding string) {
	_this.Stream.TextPos.Errorf("unexpected [%v] while decoding %v", _this.DescribeCurrentChar(), decoding)
}

func (_this *DecoderContext) DescribeCurrentChar() string {
	return _this.Stream.TextPos.DescribeCurrentChar()
}
