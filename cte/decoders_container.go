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
	if ctx.Stream.PeekByteNoEOD() != '=' {
		ctx.Stream.Errorf("Expected map separator (=) but got [%v]", ctx.Stream.DescribeCurrentChar())
	}
	ctx.Stream.AdvanceByte()
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
	decodeByFirstChar(ctx)
	ctx.ChangeDecoder(decodeMapKey)
}

func advanceAndDecodeMarkupContentBegin(ctx *DecoderContext) {
	ctx.Stream.AdvanceByte() // Advance past ','

	ctx.EventReceiver.OnEnd()
	ctx.BeginMarkupContents()
}

func decodeMarkupContents(ctx *DecoderContext) {
	ctx.stack[len(ctx.stack)-1].IsMarkupContents = true
	str, next := ctx.Stream.DecodeMarkupContent()
	if len(str) > 0 {
		ctx.EventReceiver.OnArray(events.ArrayTypeString, uint64(len(str)), str)
	}
	switch next {
	case nextIsCommentBegin:
		ctx.EventReceiver.OnComment()
		ctx.StackDecoder(decodeCommentContents)
	case nextIsCommentEnd:
		ctx.EventReceiver.OnEnd()
		ctx.UnstackDecoder()
	case nextIsSingleLineComment:
		ctx.EventReceiver.OnComment()
		contents := ctx.Stream.DecodeSingleLineComment()
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

	b := ctx.Stream.ReadByte()
	switch b {
	case '/':
		ctx.EventReceiver.OnComment()
		contents := ctx.Stream.DecodeSingleLineComment()
		if len(contents) > 0 {
			ctx.EventReceiver.OnArray(events.ArrayTypeString, uint64(len(contents)), contents)
		}
		ctx.EventReceiver.OnEnd()
		ctx.StackDecoder(decodePostInvisible)
	case '*':
		ctx.EventReceiver.OnComment()
		ctx.StackDecoder(decodeCommentContents)
	default:
		ctx.Stream.Errorf("Unexpected comment initiator: [%c]", b)
	}
}

func decodeCommentContents(ctx *DecoderContext) {
	str, next := ctx.Stream.DecodeMultilineComment()
	if len(str) > 0 {
		ctx.EventReceiver.OnArray(events.ArrayTypeString, uint64(len(str)), str)
	}
	switch next {
	case nextIsCommentBegin:
		ctx.EventReceiver.OnComment()
		ctx.StackDecoder(decodeCommentContents)
	case nextIsCommentEnd:
		ctx.EventReceiver.OnEnd()
		ctx.UnstackDecoder()
		ctx.StackDecoder(decodePostInvisible)
	}
}

func advanceAndDecodeCommentEnd(ctx *DecoderContext) {
	ctx.Stream.AdvanceByte() // Advance past '*'

	b := ctx.Stream.ReadByte()
	switch b {
	case '/':
		ctx.UnstackDecoder()
		ctx.EventReceiver.OnEnd()
		decodeByFirstChar(ctx)
	default:
		ctx.Stream.Errorf("Unexpected comment end char: [%c]", b)
	}
}

func advanceAndDecodeMetadataBegin(ctx *DecoderContext) {
	ctx.Stream.AdvanceByte() // Advance past '('

	ctx.EventReceiver.OnMetadata()
	ctx.StackDecoder(decodeMetadataKey)
}

func decodeMetadataKey(ctx *DecoderContext) {
	ctx.ChangeDecoder(decodeMetadataValue)
	decodeByFirstChar(ctx)
}

func decodeMetadataValue(ctx *DecoderContext) {
	decodeWhitespace(ctx)
	if ctx.Stream.PeekByteNoEOD() != '=' {
		// TODO: Allow comments before the =
		ctx.Stream.Errorf("Expected Metadata separator (=) but got [%v]", ctx.Stream.DescribeCurrentChar())
	}
	ctx.Stream.AdvanceByte()
	decodeWhitespace(ctx)
	ctx.ChangeDecoder(decodeMetadataKey)
	decodeByFirstChar(ctx)
}

func advanceAndDecodeMetadataEnd(ctx *DecoderContext) {
	ctx.Stream.AdvanceByte() // Advance past ')'

	ctx.EventReceiver.OnEnd()
	ctx.ChangeDecoder(decodeMetadataCompletion)
}

func decodeMetadataCompletion(ctx *DecoderContext) {
	ctx.UnstackDecoder()
	decodeByFirstChar(ctx)
}
