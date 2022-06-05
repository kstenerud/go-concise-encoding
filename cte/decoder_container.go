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

func advanceAndDecodeMapBegin(ctx *DecoderContext) {
	ctx.AssertHasStructuralWS()
	ctx.Stream.AdvanceByte() // Advance past '{'

	ctx.EventReceiver.OnMap()
	ctx.StackDecoder(decodeMapKey)
	ctx.SetContainerType(ContainerTypeMap)
}

func decodeMapKey(ctx *DecoderContext) {
	ctx.ChangeDecoder(decodeMapValue)
	decodeByFirstChar(ctx)
}

func decodeMapValue(ctx *DecoderContext) {
	decodeWhitespace(ctx)
	if ctx.Stream.ReadByteNoEOF() != '=' {
		ctx.Errorf("Expected map separator (=) but got [%v]", ctx.DescribeCurrentChar())
	}
	ctx.NoNeedForWS()
	decodeWhitespace(ctx)
	ctx.ChangeDecoder(decodeMapKey)
	decodeByFirstChar(ctx)
}

func advanceAndDecodeMapEnd(ctx *DecoderContext) {
	ctx.Stream.AdvanceByte() // Advance past '}'

	ctx.AssertIsInMap()
	ctx.EventReceiver.OnEnd()
	ctx.UnstackDecoder()
	ctx.RequireStructuralWS()
}

func advanceAndDecodeListBegin(ctx *DecoderContext) {
	ctx.AssertHasStructuralWS()
	ctx.Stream.AdvanceByte() // Advance past '['

	ctx.EventReceiver.OnList()
	ctx.StackDecoder(decodeByFirstChar)
	ctx.SetContainerType(ContainerTypeList)
}

func advanceAndDecodeListEnd(ctx *DecoderContext) {
	ctx.Stream.AdvanceByte() // Advance past ']'

	ctx.AssertIsInList()
	ctx.EventReceiver.OnEnd()
	ctx.UnstackDecoder()
	ctx.RequireStructuralWS()
}

func decodeEdgeBegin(ctx *DecoderContext) {
	ctx.EventReceiver.OnEdge()
	ctx.StackDecoder(decodeEdgeEnd)
	ctx.StackDecoder(decodeEdgeComponent)
	ctx.StackDecoder(decodeEdgeComponent)
	ctx.StackDecoder(decodeEdgeComponent)
}

func decodeEdgeComponent(ctx *DecoderContext) {
	ctx.UnstackDecoder()
	decodeByFirstChar(ctx)
}

func decodeEdgeEnd(ctx *DecoderContext) {
	ctx.Stream.SkipWhitespace()
	if ctx.Stream.ReadByteNoEOF() != ')' {
		ctx.Stream.UnreadLastByte()
		ctx.Errorf("Expected ')' at end of edge structure")
	}
	ctx.UnstackDecoder()
	ctx.RequireStructuralWS()
}

func advanceAndDecodeNodeBegin(ctx *DecoderContext) {
	ctx.AssertHasStructuralWS()
	ctx.Stream.AdvanceByte() // Advance past '('

	ctx.EventReceiver.OnNode()
	ctx.StackDecoder(decodeByFirstChar)
	ctx.SetContainerType(ContainerTypeNode)
}

func advanceAndDecodeNodeEnd(ctx *DecoderContext) {
	ctx.Stream.AdvanceByte() // Advance past ')'

	ctx.AssertIsInNode()
	ctx.EventReceiver.OnEnd()
	ctx.UnstackDecoder()
	ctx.RequireStructuralWS()
}
