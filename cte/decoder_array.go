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
	"math"

	"github.com/kstenerud/go-concise-encoding/events"
	"github.com/kstenerud/go-concise-encoding/internal/common"
)

type advanceAndDecodeQuotedString struct{}

var global_advanceAndDecodeQuotedString advanceAndDecodeQuotedString

func (_this advanceAndDecodeQuotedString) Run(ctx *DecoderContext) {
	ctx.Stream.AdvanceByte() // Advance past '"'

	bytes := ctx.Stream.ReadQuotedString()
	ctx.EventReceiver.OnArray(events.ArrayTypeString, uint64(len(bytes)), bytes)
}

type advanceAndDecodeResourceID struct{}

var global_advanceAndDecodeResourceID advanceAndDecodeResourceID

func (_this advanceAndDecodeResourceID) Run(ctx *DecoderContext) {
	ctx.Stream.AdvanceByte() // Advance past '@'
	if ctx.Stream.ReadByteNoEOF() != '"' {
		ctx.Stream.UnreadByte()
		ctx.Stream.unexpectedChar("resource ID")
	}

	bytes := ctx.Stream.ReadQuotedString()
	if ctx.Stream.PeekByteAllowEOF() == ':' {
		ctx.EventReceiver.OnArray(events.ArrayTypeResourceIDConcat, uint64(len(bytes)), bytes)
		ctx.Stream.AdvanceByte()
		if ctx.Stream.PeekByteNoEOF() != '"' {
			ctx.Errorf("Only strings may be appended to a resource ID")
		}
		global_advanceAndDecodeQuotedString.Run(ctx)
		return
	}
	ctx.EventReceiver.OnArray(events.ArrayTypeResourceID, uint64(len(bytes)), bytes)
}

func decodeArrayType(ctx *DecoderContext) string {
	arrayType := ctx.Stream.ReadToken()
	if len(arrayType) > 0 && arrayType[len(arrayType)-1] == '|' {
		arrayType = arrayType[:len(arrayType)-1]
		ctx.Stream.UnreadByte()
	}
	common.ASCIIBytesToLower(arrayType)
	return string(arrayType)
}

func finishTypedArray(ctx *DecoderContext, arrayType events.ArrayType, digitType string, bytesPerElement int, data []byte) {
	switch ctx.Stream.ReadByteNoEOF() {
	case '|':
		ctx.EventReceiver.OnArray(arrayType, uint64(len(data)/bytesPerElement), data)
		return
	default:
		ctx.Errorf("Expected %v digits", digitType)
	}
}

type advanceAndDecodeTypedArrayBegin struct{}

var global_advanceAndDecodeTypedArrayBegin advanceAndDecodeTypedArrayBegin

type uintTokenDecoder func(Token, *TextPositionCounter) (v uint64, digitCount int, decodedCount int)
type intTokenDecoder func(Token, *TextPositionCounter) (v int64, digitCount int, decodedCount int)
type floatTokenDecoder func(Token, *TextPositionCounter) (v float64, decodedCount int)

func (_this advanceAndDecodeTypedArrayBegin) decodeElementIntAnyType(ctx *DecoderContext) (v int64, success bool) {
	token := ctx.Stream.ReadToken()
	// TODO: This needs to check for other bases
	value, digitCount, decodedCount := token.DecodeSmallDecimalInt(ctx.TextPos)
	token[decodedCount:].AssertAtEnd(ctx.TextPos, "integer")
	return value, digitCount > 0
}

func (_this advanceAndDecodeTypedArrayBegin) decodeElementIntOctal(ctx *DecoderContext) (v int64, success bool) {
	token := ctx.Stream.ReadToken()
	value, digitCount, decodedCount := token.DecodeSmallOctalInt(ctx.TextPos)
	token[decodedCount:].AssertAtEnd(ctx.TextPos, "integer")
	return value, digitCount > 0
}

func (_this advanceAndDecodeTypedArrayBegin) Run(ctx *DecoderContext) {
	ctx.Stream.AdvanceByte() // Advance past '|'

	arrayType := decodeArrayType(ctx)
	ctx.Stream.SkipWhitespace()
	switch arrayType {
	case "cb":
		decodeCustomBinary(ctx)
	case "ct":
		decodeCustomText(ctx)
	case "u":
		decodeArrayUUID(ctx)
	case "b":
		decodeArrayBoolean(ctx)
	case "u8":
		decodeArrayU8(ctx, "uint", Token.DecodeSmallUint)
	case "u8b":
		decodeArrayU8(ctx, "binary", Token.DecodeSmallBinaryUint)
	case "u8o":
		decodeArrayU8(ctx, "octal", Token.DecodeSmallOctalUint)
	case "u8x":
		decodeArrayU8(ctx, "hex", Token.DecodeSmallHexUint)
	case "u16":
		decodeArrayU16(ctx, "uint", Token.DecodeSmallUint)
	case "u16b":
		decodeArrayU16(ctx, "binary", Token.DecodeSmallBinaryUint)
	case "u16o":
		decodeArrayU16(ctx, "octal", Token.DecodeSmallOctalUint)
	case "u16x":
		decodeArrayU16(ctx, "hex", Token.DecodeSmallHexUint)
	case "u32":
		decodeArrayU32(ctx, "uint", Token.DecodeSmallUint)
	case "u32b":
		decodeArrayU32(ctx, "binary", Token.DecodeSmallBinaryUint)
	case "u32o":
		decodeArrayU32(ctx, "octal", Token.DecodeSmallOctalUint)
	case "u32x":
		decodeArrayU32(ctx, "hex", Token.DecodeSmallHexUint)
	case "u64":
		decodeArrayU64(ctx, "uint", Token.DecodeSmallUint)
	case "u64b":
		decodeArrayU64(ctx, "binary", Token.DecodeSmallBinaryUint)
	case "u64o":
		decodeArrayU64(ctx, "octal", Token.DecodeSmallOctalUint)
	case "u64x":
		decodeArrayU64(ctx, "hex", Token.DecodeSmallHexUint)
	case "i8":
		decodeArrayI8(ctx, "int", Token.DecodeSmallInt)
	case "i8b":
		decodeArrayI8(ctx, "binary", Token.DecodeSmallBinaryInt)
	case "i8o":
		decodeArrayI8(ctx, "octal", Token.DecodeSmallOctalInt)
	case "i8x":
		decodeArrayI8(ctx, "hex", Token.DecodeSmallHexInt)
	case "i16":
		decodeArrayI16(ctx, "integer", Token.DecodeSmallInt)
	case "i16b":
		decodeArrayI16(ctx, "binary", Token.DecodeSmallBinaryInt)
	case "i16o":
		decodeArrayI16(ctx, "octal", Token.DecodeSmallOctalInt)
	case "i16x":
		decodeArrayI16(ctx, "hex", Token.DecodeSmallHexInt)
	case "i32":
		decodeArrayI32(ctx, "integer", Token.DecodeSmallInt)
	case "i32b":
		decodeArrayI32(ctx, "binary", Token.DecodeSmallBinaryInt)
	case "i32o":
		decodeArrayI32(ctx, "octal", Token.DecodeSmallOctalInt)
	case "i32x":
		decodeArrayI32(ctx, "hex", Token.DecodeSmallHexInt)
	case "i64":
		decodeArrayI64(ctx, "integer", Token.DecodeSmallInt)
	case "i64b":
		decodeArrayI64(ctx, "binary", Token.DecodeSmallBinaryInt)
	case "i64o":
		decodeArrayI64(ctx, "octal", Token.DecodeSmallOctalInt)
	case "i64x":
		decodeArrayI64(ctx, "hex", Token.DecodeSmallHexInt)
	case "f16":
		decodeArrayF16(ctx, "float", Token.DecodeSmallFloat)
	case "f16x":
		decodeArrayF16(ctx, "hex", Token.DecodeSmallHexFloat)
	case "f32":
		decodeArrayF32(ctx, "float", Token.DecodeSmallFloat)
	case "f32x":
		decodeArrayF32(ctx, "hex", Token.DecodeSmallHexFloat)
	case "f64":
		decodeArrayF64(ctx, "float", Token.DecodeSmallFloat)
	case "f64x":
		decodeArrayF64(ctx, "hex", Token.DecodeSmallHexFloat)
	default:
		ctx.Errorf("%s: Unhandled array type", arrayType)
	}
}

func decodeCustomText(ctx *DecoderContext) {
	bytes := ctx.Stream.ReadStringArray()
	ctx.EventReceiver.OnArray(events.ArrayTypeCustomText, uint64(len(bytes)), bytes)
}

func decodeCustomBinary(ctx *DecoderContext) {
	digitType := "hex"
	ctx.Scratch = ctx.Scratch[:0]
	for {
		ctx.Stream.SkipWhitespace()
		token := ctx.Stream.ReadToken()
		if len(token) == 0 {
			break
		}
		v, _, decodedCount := token.DecodeSmallHexUint(ctx.TextPos)
		token[decodedCount:].AssertAtEnd(ctx.TextPos, digitType)
		if v > maxUint8Value {
			ctx.Errorf("%v value too big for array type", digitType)
		}
		ctx.Scratch = append(ctx.Scratch, uint8(v))
	}
	finishTypedArray(ctx, events.ArrayTypeCustomBinary, digitType, 1, ctx.Scratch)
}

func decodeArrayBoolean(ctx *DecoderContext) {
	ctx.Scratch = ctx.Scratch[:0]
	var elemCount uint64

	for {
		nextByte := byte(0)
		for i := 0; i < 8; i++ {
			ctx.Stream.SkipWhitespace()
			b := ctx.Stream.ReadByteNoEOF()
			switch b {
			case '|':
				if i > 0 {
					ctx.Scratch = append(ctx.Scratch, nextByte)
				}
				ctx.EventReceiver.OnArray(events.ArrayTypeBit, elemCount, ctx.Scratch)
				return
			case '0':
				// Nothing to do
			case '1':
				nextByte |= byte(1 << i)
			default:
				ctx.UnexpectedChar("boolean array")
			}
			elemCount++
		}
		ctx.Scratch = append(ctx.Scratch, nextByte)
	}
}

func decodeArrayU8(ctx *DecoderContext, digitType string, decodeElement uintTokenDecoder) {
	ctx.Scratch = ctx.Scratch[:0]
	for {
		ctx.Stream.SkipWhitespace()
		token := ctx.Stream.ReadToken()
		if len(token) == 0 {
			break
		}
		v, _, decodedCount := decodeElement(token, ctx.TextPos)
		token[decodedCount:].AssertAtEnd(ctx.TextPos, digitType)
		if v > maxUint8Value {
			ctx.Errorf("%v value too big for array type", digitType)
		}
		ctx.Scratch = append(ctx.Scratch, uint8(v))
	}
	finishTypedArray(ctx, events.ArrayTypeUint8, digitType, 1, ctx.Scratch)
}

func decodeArrayU16(ctx *DecoderContext, digitType string, decodeElement uintTokenDecoder) {
	ctx.Scratch = ctx.Scratch[:0]
	for {
		ctx.Stream.SkipWhitespace()
		token := ctx.Stream.ReadToken()
		if len(token) == 0 {
			break
		}
		v, _, decodedCount := decodeElement(token, ctx.TextPos)
		token[decodedCount:].AssertAtEnd(ctx.TextPos, digitType)
		if v > maxUint16Value {
			ctx.Errorf("%v value too big for array type", digitType)
		}
		ctx.Scratch = append(ctx.Scratch, uint8(v), uint8(v>>8))
	}
	finishTypedArray(ctx, events.ArrayTypeUint16, digitType, 2, ctx.Scratch)
}

func decodeArrayU32(ctx *DecoderContext, digitType string, decodeElement uintTokenDecoder) {
	ctx.Scratch = ctx.Scratch[:0]
	for {
		ctx.Stream.SkipWhitespace()
		token := ctx.Stream.ReadToken()
		if len(token) == 0 {
			break
		}
		v, _, decodedCount := decodeElement(token, ctx.TextPos)
		token[decodedCount:].AssertAtEnd(ctx.TextPos, digitType)
		if v > maxUint32Value {
			ctx.Errorf("%v value too big for array type", digitType)
		}
		ctx.Scratch = append(ctx.Scratch, uint8(v), uint8(v>>8), uint8(v>>16), uint8(v>>24))
	}
	finishTypedArray(ctx, events.ArrayTypeUint32, digitType, 4, ctx.Scratch)
}

func decodeArrayU64(ctx *DecoderContext, digitType string, decodeElement uintTokenDecoder) {
	ctx.Scratch = ctx.Scratch[:0]
	for {
		ctx.Stream.SkipWhitespace()
		token := ctx.Stream.ReadToken()
		if len(token) == 0 {
			break
		}
		v, _, decodedCount := decodeElement(token, ctx.TextPos)
		token[decodedCount:].AssertAtEnd(ctx.TextPos, digitType)
		ctx.Scratch = append(ctx.Scratch, uint8(v), uint8(v>>8), uint8(v>>16), uint8(v>>24),
			uint8(v>>32), uint8(v>>40), uint8(v>>48), uint8(v>>56))
	}
	finishTypedArray(ctx, events.ArrayTypeUint64, digitType, 8, ctx.Scratch)
}

func decodeArrayI8(ctx *DecoderContext, digitType string, decodeElement intTokenDecoder) {
	ctx.Scratch = ctx.Scratch[:0]
	for {
		ctx.Stream.SkipWhitespace()
		token := ctx.Stream.ReadToken()
		if len(token) == 0 {
			break
		}
		v, _, decodedCount := decodeElement(token, ctx.TextPos)
		token[decodedCount:].AssertAtEnd(ctx.TextPos, digitType)
		if v < minInt8Value || v > maxInt8Value {
			ctx.Errorf("%v value too big for array type", digitType)
		}
		ctx.Scratch = append(ctx.Scratch, uint8(v))
	}
	finishTypedArray(ctx, events.ArrayTypeInt8, digitType, 1, ctx.Scratch)
}

func decodeArrayI16(ctx *DecoderContext, digitType string, decodeElement intTokenDecoder) {
	ctx.Scratch = ctx.Scratch[:0]
	for {
		ctx.Stream.SkipWhitespace()
		token := ctx.Stream.ReadToken()
		if len(token) == 0 {
			break
		}
		v, _, decodedCount := decodeElement(token, ctx.TextPos)
		token[decodedCount:].AssertAtEnd(ctx.TextPos, digitType)
		if v < minInt16Value || v > maxInt16Value {
			ctx.Errorf("%v value too big for array type", digitType)
		}
		ctx.Scratch = append(ctx.Scratch, uint8(v), uint8(v>>8))
	}
	finishTypedArray(ctx, events.ArrayTypeInt16, digitType, 2, ctx.Scratch)
}

func decodeArrayI32(ctx *DecoderContext, digitType string, decodeElement intTokenDecoder) {
	ctx.Scratch = ctx.Scratch[:0]
	for {
		ctx.Stream.SkipWhitespace()
		token := ctx.Stream.ReadToken()
		if len(token) == 0 {
			break
		}
		v, _, decodedCount := decodeElement(token, ctx.TextPos)
		token[decodedCount:].AssertAtEnd(ctx.TextPos, digitType)
		if v < minInt32Value || v > maxInt32Value {
			ctx.Errorf("%v value too big for array type", digitType)
		}
		ctx.Scratch = append(ctx.Scratch, uint8(v), uint8(v>>8), uint8(v>>16), uint8(v>>24))
	}
	finishTypedArray(ctx, events.ArrayTypeInt32, digitType, 4, ctx.Scratch)
}

func decodeArrayI64(ctx *DecoderContext, digitType string, decodeElement intTokenDecoder) {
	ctx.Scratch = ctx.Scratch[:0]
	for {
		ctx.Stream.SkipWhitespace()
		token := ctx.Stream.ReadToken()
		if len(token) == 0 {
			break
		}
		v, _, decodedCount := decodeElement(token, ctx.TextPos)
		token[decodedCount:].AssertAtEnd(ctx.TextPos, digitType)
		ctx.Scratch = append(ctx.Scratch, uint8(v), uint8(v>>8), uint8(v>>16), uint8(v>>24),
			uint8(v>>32), uint8(v>>40), uint8(v>>48), uint8(v>>56))
	}
	finishTypedArray(ctx, events.ArrayTypeInt64, digitType, 8, ctx.Scratch)
}

func decodeArrayF16(ctx *DecoderContext, digitType string, decodeElement floatTokenDecoder) {
	ctx.Scratch = ctx.Scratch[:0]
	for {
		ctx.Stream.SkipWhitespace()
		token := ctx.Stream.ReadToken()
		if len(token) == 0 {
			break
		}
		v, decodedCount := decodeElement(token, ctx.TextPos)
		token[decodedCount:].AssertAtEnd(ctx.TextPos, digitType)
		exp := extractFloat64Exponent(v)
		if exp < minFloat32Exponent || exp > maxFloat32Exponent {
			ctx.Errorf("Exponent too big for float32 type")
		}
		bits := math.Float32bits(float32(v))
		ctx.Scratch = append(ctx.Scratch, uint8(bits>>16), uint8(bits>>24))
	}
	finishTypedArray(ctx, events.ArrayTypeFloat16, digitType, 2, ctx.Scratch)
}

func decodeArrayF32(ctx *DecoderContext, digitType string, decodeElement floatTokenDecoder) {
	ctx.Scratch = ctx.Scratch[:0]
	for {
		ctx.Stream.SkipWhitespace()
		token := ctx.Stream.ReadToken()
		if len(token) == 0 {
			break
		}
		v, decodedCount := decodeElement(token, ctx.TextPos)
		token[decodedCount:].AssertAtEnd(ctx.TextPos, digitType)
		exp := extractFloat64Exponent(v)
		if exp < minFloat32Exponent || exp > maxFloat32Exponent {
			ctx.Errorf("Exponent too big for float32 type")
		}
		bits := math.Float32bits(float32(v))
		ctx.Scratch = append(ctx.Scratch, uint8(bits), uint8(bits>>8), uint8(bits>>16), uint8(bits>>24))
	}
	finishTypedArray(ctx, events.ArrayTypeFloat32, digitType, 4, ctx.Scratch)
}

func decodeArrayF64(ctx *DecoderContext, digitType string, decodeElement floatTokenDecoder) {
	ctx.Scratch = ctx.Scratch[:0]
	for {
		ctx.Stream.SkipWhitespace()
		token := ctx.Stream.ReadToken()
		if len(token) == 0 {
			break
		}
		v, decodedCount := decodeElement(token, ctx.TextPos)
		token[decodedCount:].AssertAtEnd(ctx.TextPos, digitType)
		bits := math.Float64bits(v)
		ctx.Scratch = append(ctx.Scratch, uint8(bits), uint8(bits>>8), uint8(bits>>16), uint8(bits>>24),
			uint8(bits>>32), uint8(bits>>40), uint8(bits>>48), uint8(bits>>56))
	}
	finishTypedArray(ctx, events.ArrayTypeFloat64, digitType, 8, ctx.Scratch)
}

func decodeArrayUUID(ctx *DecoderContext) {
	ctx.Scratch = ctx.Scratch[:0]
	ctx.Stream.SkipWhitespace()
	for ctx.Stream.PeekByteAllowEOF() != '|' {
		token := ctx.Stream.ReadToken()
		ctx.Scratch = append(ctx.Scratch, token.DecodeUUID(ctx.TextPos)...)
		ctx.Stream.SkipWhitespace()
	}
	finishTypedArray(ctx, events.ArrayTypeUUID, "uuid", 16, ctx.Scratch)
}
