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
)

type advanceAndDecodeMapBegin struct{}

var global_advanceAndDecodeMapBegin advanceAndDecodeMapBegin

func (_this advanceAndDecodeMapBegin) Run(ctx *DecoderContext) {
	ctx.Stream.AdvanceByte() // Advance past '{'

	ctx.EventReceiver.OnMap()
	ctx.StackDecoder(global_decodeMapKey)
}

type decodeMapKey struct{}

var global_decodeMapKey decodeMapKey

func (_this decodeMapKey) Run(ctx *DecoderContext) {
	ctx.ChangeDecoder(global_decodeMapValue)
	global_decodeByFirstChar.Run(ctx)
}

type decodeMapValue struct{}

var global_decodeMapValue decodeMapValue

func (_this decodeMapValue) Run(ctx *DecoderContext) {
	global_decodeWhitespace.Run(ctx)
	if ctx.Stream.ReadByteNoEOD() != '=' {
		ctx.Stream.Errorf("Expected map separator (=) but got [%v]", ctx.Stream.DescribeCurrentChar())
	}
	global_decodeWhitespace.Run(ctx)
	ctx.ChangeDecoder(global_decodeMapKey)
	global_decodeByFirstChar.Run(ctx)
}

type advanceAndDecodeMapEnd struct{}

var global_advanceAndDecodeMapEnd advanceAndDecodeMapEnd

func (_this advanceAndDecodeMapEnd) Run(ctx *DecoderContext) {
	ctx.Stream.AdvanceByte() // Advance past '}'

	ctx.EventReceiver.OnEnd()
	ctx.UnstackDecoder()
}

type advanceAndDecodeListBegin struct{}

var global_advanceAndDecodeListBegin advanceAndDecodeListBegin

func (_this advanceAndDecodeListBegin) Run(ctx *DecoderContext) {
	ctx.Stream.AdvanceByte() // Advance past '['

	ctx.EventReceiver.OnList()
	ctx.StackDecoder(global_decodeByFirstChar)
}

type advanceAndDecodeListEnd struct{}

var global_advanceAndDecodeListEnd advanceAndDecodeListEnd

func (_this advanceAndDecodeListEnd) Run(ctx *DecoderContext) {
	ctx.Stream.AdvanceByte() // Advance past ']'

	ctx.EventReceiver.OnEnd()
	ctx.UnstackDecoder()
}

type advanceAndDecodeMarkupBegin struct{}

var global_advanceAndDecodeMarkupBegin advanceAndDecodeMarkupBegin

func (_this advanceAndDecodeMarkupBegin) Run(ctx *DecoderContext) {
	ctx.Stream.AdvanceByte() // Advance past '<'
	global_decodeMarkupBegin.Run(ctx)
}

type decodeMarkupBegin struct{}

var global_decodeMarkupBegin decodeMarkupBegin

func (_this decodeMarkupBegin) Run(ctx *DecoderContext) {
	ctx.EventReceiver.OnMarkup()
	ctx.StackDecoder(global_decodeMarkupName)
}

type decodeMarkupName struct{}

var global_decodeMarkupName decodeMarkupName

func (_this decodeMarkupName) Run(ctx *DecoderContext) {
	global_decodeIdentifier.Run(ctx)
	ctx.ChangeDecoder(global_decodeMapKey)
}

type advanceAndDecodeMarkupContentBegin struct{}

var global_advanceAndDecodeMarkupContentBegin advanceAndDecodeMarkupContentBegin

func (_this advanceAndDecodeMarkupContentBegin) Run(ctx *DecoderContext) {
	ctx.Stream.AdvanceByte() // Advance past ','

	ctx.EventReceiver.OnEnd()
	ctx.ChangeDecoder(global_decodeMarkupContents)
}

type decodeMarkupContents struct{}

var global_decodeMarkupContents decodeMarkupContents

func (_this decodeMarkupContents) Run(ctx *DecoderContext) {
	ctx.stack[len(ctx.stack)-1].IsMarkupContents = true
	str, next := ctx.Stream.ReadMarkupContent()
	if len(str) > 0 {
		ctx.EventReceiver.OnArray(events.ArrayTypeString, uint64(len(str)), str)
	}
	switch next {
	case nextIsCommentBegin:
		ctx.BeginComment()
	case nextIsCommentEnd:
		ctx.EndComment()
	case nextIsSingleLineComment:
		ctx.EventReceiver.OnComment()
		contents := ctx.Stream.ReadSingleLineComment()
		if len(contents) > 0 {
			ctx.EventReceiver.OnArray(events.ArrayTypeString, uint64(len(contents)), contents)
		}
		ctx.EventReceiver.OnEnd()
	case nextIsMarkupBegin:
		global_decodeMarkupBegin.Run(ctx)
	case nextIsMarkupEnd:
		ctx.EventReceiver.OnEnd()
		ctx.UnstackDecoder()
	}
}

type advanceAndDecodeMarkupEnd struct{}

var global_advanceAndDecodeMarkupEnd advanceAndDecodeMarkupEnd

func (_this advanceAndDecodeMarkupEnd) Run(ctx *DecoderContext) {
	ctx.Stream.AdvanceByte() // Advance past '>'

	ctx.EventReceiver.OnEnd()
	ctx.EndMarkup()
}

type advanceAndDecodeComment struct{}

var global_advanceAndDecodeComment advanceAndDecodeComment

func (_this advanceAndDecodeComment) Run(ctx *DecoderContext) {
	ctx.Stream.AdvanceByte() // Advance past '/'

	b := ctx.Stream.ReadByteNoEOD()
	switch b {
	case '/':
		ctx.EventReceiver.OnComment()
		contents := ctx.Stream.ReadSingleLineComment()
		if len(contents) > 0 {
			ctx.EventReceiver.OnArray(events.ArrayTypeString, uint64(len(contents)), contents)
		}
		ctx.EventReceiver.OnEnd()
		ctx.StackDecoder(global_decodePostInvisible)
	case '*':
		ctx.BeginComment()
	default:
		ctx.Stream.Errorf("Unexpected comment initiator: [%c]", b)
	}
}

type decodeCommentContents struct{}

var global_decodeCommentContents decodeCommentContents

func (_this decodeCommentContents) Run(ctx *DecoderContext) {
	str, next := ctx.Stream.ReadMultilineComment()
	if len(str) > 0 {
		ctx.EventReceiver.OnArray(events.ArrayTypeString, uint64(len(str)), str)
	}
	switch next {
	case nextIsCommentBegin:
		ctx.BeginComment()
	case nextIsCommentEnd:
		ctx.EndComment()
	}
}

type advanceAndDecodeCommentEnd struct{}

var global_advanceAndDecodeCommentEnd advanceAndDecodeCommentEnd

func (_this advanceAndDecodeCommentEnd) Run(ctx *DecoderContext) {
	ctx.Stream.AdvanceByte() // Advance past '*'

	b := ctx.Stream.ReadByteNoEOD()
	switch b {
	case '/':
		ctx.UnstackDecoder()
		ctx.EventReceiver.OnEnd()
		global_decodeByFirstChar.Run(ctx)
	default:
		ctx.Stream.Errorf("Unexpected comment end char: [%c]", b)
	}
}
