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
	"io"

	"github.com/kstenerud/go-concise-encoding/events"
	"github.com/kstenerud/go-concise-encoding/options"
)

type DecoderStackEntry struct {
	DecoderFunc      DecoderOp
	IsMarkupContents bool
}

type DecoderContext struct {
	opts               options.CTEDecoderOptions
	Stream             Reader
	TextPos            *TextPositionCounter
	EventReceiver      events.DataEventReceiver
	stack              []DecoderStackEntry
	IsDocumentComplete bool
}

func (_this *DecoderContext) Init(opts *options.CTEDecoderOptions, reader io.Reader, eventReceiver events.DataEventReceiver) {
	_this.opts = *opts
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

func (_this *DecoderContext) DecodeNext() {
	entry := _this.stack[len(_this.stack)-1]
	entry.DecoderFunc.Run(_this)
}

func (_this *DecoderContext) ChangeDecoder(decoder DecoderOp) {
	_this.stack[len(_this.stack)-1] = DecoderStackEntry{
		DecoderFunc:      decoder,
		IsMarkupContents: false,
	}
}

func (_this *DecoderContext) StackDecoder(decoder DecoderOp) {
	_this.stack = append(_this.stack, DecoderStackEntry{
		DecoderFunc:      decoder,
		IsMarkupContents: false,
	})
}

func (_this *DecoderContext) UnstackDecoder() DecoderStackEntry {
	_this.stack = _this.stack[:len(_this.stack)-1]
	return _this.stack[len(_this.stack)-1]
}

func (_this *DecoderContext) BeginComment() {
	_this.EventReceiver.OnComment()
	_this.StackDecoder(global_decodeCommentContents)
}

func (_this *DecoderContext) EndComment() {
	_this.EventReceiver.OnEnd()
	entry := _this.UnstackDecoder()
	if entry.DecoderFunc != global_decodeCommentContents && entry.DecoderFunc != global_decodeMarkupContents {
		_this.StackDecoder(global_decodePostInvisible)
	}
}

func (_this *DecoderContext) EndMarkup() {
	if !_this.stack[len(_this.stack)-1].IsMarkupContents {
		// Add second end message
		_this.EventReceiver.OnEnd()
	}
	_this.UnstackDecoder()
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
