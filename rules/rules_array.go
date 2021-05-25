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

// Imposes the structural rules that enforce a well-formed concise encoding
// document.
package rules

type ArrayRule struct{}

func (_this *ArrayRule) String() string { return "Array Rule" }
func (_this *ArrayRule) OnArrayChunk(ctx *Context, length uint64, moreChunksFollow bool) {
	if length == 0 {
		ctx.tryEndArray(moreChunksFollow, nil)
		return
	}

	ctx.BeginChunkAnyType(length, moreChunksFollow)
}

// =============================================================================

type ArrayChunkRule struct{}

func (_this *ArrayChunkRule) String() string { return "Array Chunk Rule" }
func (_this *ArrayChunkRule) OnArrayData(ctx *Context, data []byte) {
	ctx.MarkCompletedChunkByteCount(uint64(len(data)))
	if ctx.chunkActualByteCount == ctx.chunkExpectedByteCount {
		ctx.EndChunkAnyType()
	}
}

// =============================================================================

type StringRule struct{}

func (_this *StringRule) String() string { return "String Rule" }
func (_this *StringRule) OnArrayChunk(ctx *Context, length uint64, moreChunksFollow bool) {
	if length == 0 {
		ctx.tryEndArray(moreChunksFollow, nil)
		return
	}

	ctx.BeginChunkString(length, moreChunksFollow)
}

// =============================================================================

type StringChunkRule struct{}

func (_this *StringChunkRule) String() string { return "String Chunk Rule" }
func (_this *StringChunkRule) OnArrayData(ctx *Context, data []byte) {
	ctx.MarkCompletedChunkByteCount(uint64(len(data)))
	firstRuneBytes, nextRunesBytes := ctx.StreamStringData(data)

	ctx.ValidateArrayDataFunc(firstRuneBytes)
	ctx.ValidateArrayDataFunc(nextRunesBytes)

	if ctx.chunkActualByteCount == ctx.chunkExpectedByteCount {
		ctx.EndChunkString()
	}
}

// =============================================================================

type StringBuilderRule struct{}

func (_this *StringBuilderRule) String() string { return "String Builder Rule" }
func (_this *StringBuilderRule) OnArrayChunk(ctx *Context, length uint64, moreChunksFollow bool) {
	if length == 0 {
		ctx.tryEndArray(moreChunksFollow, nil)
		return
	}

	ctx.BeginChunkStringBuilder(length, moreChunksFollow)
}

// =============================================================================

type StringBuilderChunkRule struct{}

func (_this *StringBuilderChunkRule) String() string { return "String Builder Chunk Rule" }
func (_this *StringBuilderChunkRule) OnArrayData(ctx *Context, data []byte) {
	ctx.MarkCompletedChunkByteCount(uint64(len(data)))
	ctx.AddBuiltArrayBytes(data)
	if ctx.chunkActualByteCount == ctx.chunkExpectedByteCount {
		ctx.EndChunkString()
	}
}
