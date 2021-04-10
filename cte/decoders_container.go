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
)

func advanceAndDecodeMapBegin(ctx *DecoderContext) {
	ctx.Stream.AdvanceByte() // Advance past '{'

	ctx.EventReceiver.OnMap()
	ctx.StackDecoder(decodeMapKey)
}

func decodeMapKey(ctx *DecoderContext) {
	ctx.ChangeDecoder(decodeMapValue)
	decodeByFirstChar(ctx)
}

func decodeMapValue(ctx *DecoderContext) {
	decodeWhitespace(ctx)
	if ctx.Stream.ReadByteNoEOD() != '=' {
		ctx.Stream.Errorf("Expected map separator (=) but got [%v]", ctx.Stream.DescribeCurrentChar())
	}
	decodeWhitespace(ctx)
	ctx.ChangeDecoder(decodeMapKey)
	decodeByFirstChar(ctx)
}

func advanceAndDecodeMapEnd(ctx *DecoderContext) {
	ctx.Stream.AdvanceByte() // Advance past '}'

	ctx.EventReceiver.OnEnd()
	ctx.UnstackDecoder()
}

func advanceAndDecodeListBegin(ctx *DecoderContext) {
	ctx.Stream.AdvanceByte() // Advance past '['

	ctx.EventReceiver.OnList()
	ctx.StackDecoder(decodeByFirstChar)
}

func advanceAndDecodeListEnd(ctx *DecoderContext) {
	ctx.Stream.AdvanceByte() // Advance past ']'

	ctx.EventReceiver.OnEnd()
	ctx.UnstackDecoder()
}

func advanceAndDecodeMarkupBegin(ctx *DecoderContext) {
	ctx.Stream.AdvanceByte() // Advance past '<'
	decodeMarkupBegin(ctx)
}

func decodeMarkupBegin(ctx *DecoderContext) {
	ctx.EventReceiver.OnMarkup()
	ctx.StackDecoder(decodeMarkupName)
}

func decodeMarkupName(ctx *DecoderContext) {
	decodeIdentifier(ctx)
	ctx.ChangeDecoder(decodeMapKey)
}

func advanceAndDecodeMarkupContentBegin(ctx *DecoderContext) {
	ctx.Stream.AdvanceByte() // Advance past ','

	ctx.EventReceiver.OnEnd()
	ctx.ChangeDecoder(decodeMarkupContents)
}

func decodeMarkupContents(ctx *DecoderContext) {
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
		decodeMarkupBegin(ctx)
	case nextIsMarkupEnd:
		ctx.EventReceiver.OnEnd()
		ctx.UnstackDecoder()
	}
}

func advanceAndDecodeMarkupEnd(ctx *DecoderContext) {
	ctx.Stream.AdvanceByte() // Advance past '>'

	ctx.EventReceiver.OnEnd()
	ctx.EndMarkup()
}

func advanceAndDecodeComment(ctx *DecoderContext) {
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
		ctx.StackDecoder(decodePostInvisible)
	case '*':
		ctx.BeginComment()
	default:
		ctx.Stream.Errorf("Unexpected comment initiator: [%c]", b)
	}
}

func decodeCommentContents(ctx *DecoderContext) {
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

func advanceAndDecodeCommentEnd(ctx *DecoderContext) {
	ctx.Stream.AdvanceByte() // Advance past '*'

	b := ctx.Stream.ReadByteNoEOD()
	switch b {
	case '/':
		ctx.UnstackDecoder()
		ctx.EventReceiver.OnEnd()
		decodeByFirstChar(ctx)
	default:
		ctx.Stream.Errorf("Unexpected comment end char: [%c]", b)
	}
}
