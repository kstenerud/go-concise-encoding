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
	"strings"

	"github.com/kstenerud/go-concise-encoding/ce/events"
	"github.com/kstenerud/go-concise-encoding/internal/common"
)

func advanceAndDecodeQuotedString(ctx *DecoderContext) {
	ctx.AssertHasStructuralWS()
	ctx.Stream.AdvanceByte() // Advance past '"'

	bytes := ctx.Stream.ReadQuotedString()
	ctx.EventReceiver.OnArray(events.ArrayTypeString, uint64(len(bytes)), bytes)
	ctx.AwaitStructuralWS()
}

func decodeResourceID(ctx *DecoderContext) {
	bytes := ctx.Stream.ReadQuotedString()
	ctx.EventReceiver.OnArray(events.ArrayTypeResourceID, uint64(len(bytes)), bytes)
	ctx.AwaitStructuralWS()
}

func decodeRemoteReference(ctx *DecoderContext) {
	bytes := ctx.Stream.ReadQuotedString()
	ctx.EventReceiver.OnArray(events.ArrayTypeReferenceRemote, uint64(len(bytes)), bytes)
	ctx.AwaitStructuralWS()
}

func decodeArrayType(ctx *DecoderContext) []byte {
	arrayType := ctx.Stream.ReadToken()
	if len(arrayType) > 0 && arrayType[len(arrayType)-1] == '|' {
		arrayType = arrayType[:len(arrayType)-1]
		ctx.Stream.UnreadLastByte()
	}
	common.ASCIIBytesToLower(arrayType)
	return arrayType
}

func finishTypedArray(ctx *DecoderContext) {
	switch ctx.Stream.ReadByteNoEOF() {
	case '|':
		data := ctx.Scratch
		elemCount := uint64(len(data) / ctx.ArrayBytesPerElement)
		if ctx.ArrayContainsComments {
			ctx.EventReceiver.OnArrayChunk(elemCount, false)
			ctx.EventReceiver.OnArrayData(data)
		} else {
			ctx.EventReceiver.OnArray(ctx.ArrayType, elemCount, data)
		}
	default:
		ctx.Errorf("Expected %v digits", ctx.ArrayDigitType)
	}
	ctx.ArrayContainsComments = false
}

type uintTokenDecoder func(Token, *TextPositionCounter) (v uint64, digitCount int, decodedCount int)
type intTokenDecoder func(Token, *TextPositionCounter) (v int64, digitCount int, decodedCount int)
type floatTokenDecoder func(Token, *TextPositionCounter) (v float64, decodedCount int)

func advanceAndDecodeTypedArrayBegin(ctx *DecoderContext) {
	ctx.AssertHasStructuralWS()
	ctx.Stream.AdvanceByte() // Advance past '|'

	arrayType := decodeArrayType(ctx)
	arrayTypeAsString := string(arrayType)
	ctx.Stream.SkipWhitespace()
	switch arrayTypeAsString {
	case "u":
		decodeArrayUID(ctx)
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
	case "", "c":
		ctx.Errorf("%v: Malformed array type", arrayType)
	default:
		if arrayTypeAsString[0] == 'c' {
			decodeCustomType(ctx, arrayTypeAsString)
			break
		}
		if strings.ContainsRune(arrayTypeAsString, '/') {
			decodeMedia(ctx, arrayTypeAsString)
			break
		}
		ctx.Errorf("%v: Unknown array type", arrayType)
	}
	ctx.AwaitStructuralWS()
}

func decodeMedia(ctx *DecoderContext, mediaType string) {
	// TODO: Better media type validation

	ctx.BeginArray("string", events.ArrayTypeMedia, 1)
	if contents, ok := tryReadStringlikeArray(ctx, "media"); ok {
		ctx.EventReceiver.OnMedia(mediaType, contents)
		return
	}

	ctx.BeginArray("hex", events.ArrayTypeMedia, 1)
	contents := readHexByteArray(ctx, events.ArrayTypeMedia, "media")
	ctx.EventReceiver.OnMedia(mediaType, contents)
}

func decodeCustomType(ctx *DecoderContext, arrayTypeAsString string) {
	if len(arrayTypeAsString) == 1 {
		ctx.Errorf("%v: Missing custom type", arrayTypeAsString)
	}
	customType := stringToPlainUint(ctx, arrayTypeAsString[1:])

	if stringContents, ok := tryReadStringlikeArray(ctx, "custom text"); ok {
		ctx.EventReceiver.OnCustomText(customType, string(stringContents))
		return
	}

	byteContents := readHexByteArray(ctx, events.ArrayTypeCustomBinary, "custom binary")
	ctx.EventReceiver.OnCustomBinary(customType, byteContents)
}

func stringToPlainUint(ctx *DecoderContext, str string) uint64 {
	accum := uint64(0)
	for _, c := range str {
		if c < '0' || c > '9' {
			ctx.Errorf("%v: Non-numeric digit encountered", str)
		}
		accum = accum*10 + uint64(c-'0')
		if accum > 0xffffffff {
			ctx.Errorf("%v: Value out of range", str)
		}
	}
	return accum
}

func tryReadStringlikeArray(ctx *DecoderContext, decoding string) (contents []byte, ok bool) {
	ctx.Stream.SkipWhitespace()
	if ctx.Stream.PeekByteNoEOF() != '"' {
		return []byte{}, false
	}
	ctx.Stream.AdvanceByte()
	return readStringlikeArray(ctx, decoding), true
}

func readStringlikeArray(ctx *DecoderContext, decoding string) []byte {
	// Assumption: Last byte read was the opening double-quote
	bytes := ctx.Stream.ReadQuotedString()
	ctx.Stream.SkipWhitespace()
	if ctx.Stream.ReadByteNoEOF() != '|' {
		ctx.UnexpectedChar(decoding)
	}
	return bytes
}

func readHexByteArray(ctx *DecoderContext, arrayType events.ArrayType, decoding string) []byte {
	ctx.BeginArray("hex", arrayType, 1)
	for {
		ctx.Stream.SkipWhitespace()
		if tryDecodeCommentInArray(ctx) {
			continue
		}
		token := ctx.Stream.ReadToken()
		if len(token) == 0 {
			break
		}
		v, _, decodedCount := token.DecodeSmallHexUint(ctx.TextPos)
		token[decodedCount:].AssertAtEnd(ctx.TextPos, ctx.ArrayDigitType)
		if v > common.Uint8Max {
			ctx.Errorf("%v value too big for array type", ctx.ArrayDigitType)
		}
		ctx.Scratch = append(ctx.Scratch, uint8(v))
	}
	if ctx.Stream.ReadByteNoEOF() != '|' {
		ctx.UnexpectedChar(decoding)
	}
	return ctx.Scratch
}

func decodeArrayBoolean(ctx *DecoderContext) {
	ctx.BeginArray("bit", events.ArrayTypeBit, 0)
	var elemCount uint64

	for {
		nextByte := byte(0)
		for i := 0; i < 8; i++ {
			ctx.Stream.SkipWhitespace()
			// TODO: Boolean array Comments
			// if tryDecodeCommentInArray(ctx) {
			// 	continue
			// }
			b := ctx.Stream.ReadByteNoEOF()
			switch b {
			case '|':
				if i > 0 {
					ctx.Scratch = append(ctx.Scratch, nextByte)
				}
				ctx.EventReceiver.OnArray(ctx.ArrayType, elemCount, ctx.Scratch)
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
	ctx.BeginArray(digitType, events.ArrayTypeUint8, 1)
	for {
		ctx.Stream.SkipWhitespace()
		if tryDecodeCommentInArray(ctx) {
			continue
		}
		token := ctx.Stream.ReadToken()
		if len(token) == 0 {
			break
		}
		v, _, decodedCount := decodeElement(token, ctx.TextPos)
		token[decodedCount:].AssertAtEnd(ctx.TextPos, ctx.ArrayDigitType)
		if v > common.Uint8Max {
			ctx.Errorf("%v value too big for array type", ctx.ArrayDigitType)
		}
		ctx.Scratch = append(ctx.Scratch, uint8(v))
	}
	finishTypedArray(ctx)
}

func decodeArrayU16(ctx *DecoderContext, digitType string, decodeElement uintTokenDecoder) {
	ctx.BeginArray(digitType, events.ArrayTypeUint16, 2)
	for {
		ctx.Stream.SkipWhitespace()
		if tryDecodeCommentInArray(ctx) {
			continue
		}
		token := ctx.Stream.ReadToken()
		if len(token) == 0 {
			break
		}
		v, _, decodedCount := decodeElement(token, ctx.TextPos)
		token[decodedCount:].AssertAtEnd(ctx.TextPos, ctx.ArrayDigitType)
		if v > common.Uint16Max {
			ctx.Errorf("%v value too big for array type", ctx.ArrayDigitType)
		}
		ctx.Scratch = append(ctx.Scratch, uint8(v), uint8(v>>8))
	}
	finishTypedArray(ctx)
}

func decodeArrayU32(ctx *DecoderContext, digitType string, decodeElement uintTokenDecoder) {
	ctx.BeginArray(digitType, events.ArrayTypeUint32, 4)
	for {
		ctx.Stream.SkipWhitespace()
		if tryDecodeCommentInArray(ctx) {
			continue
		}
		token := ctx.Stream.ReadToken()
		if len(token) == 0 {
			break
		}
		v, _, decodedCount := decodeElement(token, ctx.TextPos)
		token[decodedCount:].AssertAtEnd(ctx.TextPos, ctx.ArrayDigitType)
		if v > common.Uint32Max {
			ctx.Errorf("%v value too big for array type", ctx.ArrayDigitType)
		}
		ctx.Scratch = append(ctx.Scratch, uint8(v), uint8(v>>8), uint8(v>>16), uint8(v>>24))
	}
	finishTypedArray(ctx)
}

func decodeArrayU64(ctx *DecoderContext, digitType string, decodeElement uintTokenDecoder) {
	ctx.BeginArray(digitType, events.ArrayTypeUint64, 8)
	for {
		ctx.Stream.SkipWhitespace()
		if tryDecodeCommentInArray(ctx) {
			continue
		}
		token := ctx.Stream.ReadToken()
		if len(token) == 0 {
			break
		}
		v, _, decodedCount := decodeElement(token, ctx.TextPos)
		token[decodedCount:].AssertAtEnd(ctx.TextPos, ctx.ArrayDigitType)
		ctx.Scratch = append(ctx.Scratch, uint8(v), uint8(v>>8), uint8(v>>16), uint8(v>>24),
			uint8(v>>32), uint8(v>>40), uint8(v>>48), uint8(v>>56))
	}
	finishTypedArray(ctx)
}

func decodeArrayI8(ctx *DecoderContext, digitType string, decodeElement intTokenDecoder) {
	ctx.BeginArray(digitType, events.ArrayTypeInt8, 1)
	for {
		ctx.Stream.SkipWhitespace()
		if tryDecodeCommentInArray(ctx) {
			continue
		}
		token := ctx.Stream.ReadToken()
		if len(token) == 0 {
			break
		}
		v, _, decodedCount := decodeElement(token, ctx.TextPos)
		token[decodedCount:].AssertAtEnd(ctx.TextPos, ctx.ArrayDigitType)
		if v < common.Int8Min || v > common.Int8Max {
			ctx.Errorf("%v value too big for array type", ctx.ArrayDigitType)
		}
		ctx.Scratch = append(ctx.Scratch, uint8(v))
	}
	finishTypedArray(ctx)
}

func decodeArrayI16(ctx *DecoderContext, digitType string, decodeElement intTokenDecoder) {
	ctx.BeginArray(digitType, events.ArrayTypeInt16, 2)
	for {
		ctx.Stream.SkipWhitespace()
		if tryDecodeCommentInArray(ctx) {
			continue
		}
		token := ctx.Stream.ReadToken()
		if len(token) == 0 {
			break
		}
		v, _, decodedCount := decodeElement(token, ctx.TextPos)
		token[decodedCount:].AssertAtEnd(ctx.TextPos, ctx.ArrayDigitType)
		if v < common.Int16Min || v > common.Int16Max {
			ctx.Errorf("%v value too big for array type", ctx.ArrayDigitType)
		}
		ctx.Scratch = append(ctx.Scratch, uint8(v), uint8(v>>8))
	}
	finishTypedArray(ctx)
}

func decodeArrayI32(ctx *DecoderContext, digitType string, decodeElement intTokenDecoder) {
	ctx.BeginArray(digitType, events.ArrayTypeInt32, 4)
	for {
		ctx.Stream.SkipWhitespace()
		if tryDecodeCommentInArray(ctx) {
			continue
		}
		token := ctx.Stream.ReadToken()
		if len(token) == 0 {
			break
		}
		v, _, decodedCount := decodeElement(token, ctx.TextPos)
		token[decodedCount:].AssertAtEnd(ctx.TextPos, ctx.ArrayDigitType)
		if v < common.Int32Min || v > common.Int32Max {
			ctx.Errorf("%v value too big for array type", ctx.ArrayDigitType)
		}
		ctx.Scratch = append(ctx.Scratch, uint8(v), uint8(v>>8), uint8(v>>16), uint8(v>>24))
	}
	finishTypedArray(ctx)
}

func decodeArrayI64(ctx *DecoderContext, digitType string, decodeElement intTokenDecoder) {
	ctx.BeginArray(digitType, events.ArrayTypeInt64, 8)
	for {
		ctx.Stream.SkipWhitespace()
		if tryDecodeCommentInArray(ctx) {
			continue
		}
		token := ctx.Stream.ReadToken()
		if len(token) == 0 {
			break
		}
		v, _, decodedCount := decodeElement(token, ctx.TextPos)
		token[decodedCount:].AssertAtEnd(ctx.TextPos, ctx.ArrayDigitType)
		ctx.Scratch = append(ctx.Scratch, uint8(v), uint8(v>>8), uint8(v>>16), uint8(v>>24),
			uint8(v>>32), uint8(v>>40), uint8(v>>48), uint8(v>>56))
	}
	finishTypedArray(ctx)
}

func decodeArrayF16(ctx *DecoderContext, digitType string, decodeElement floatTokenDecoder) {
	ctx.BeginArray(digitType, events.ArrayTypeFloat16, 2)
	for {
		ctx.Stream.SkipWhitespace()
		if tryDecodeCommentInArray(ctx) {
			continue
		}
		token := ctx.Stream.ReadToken()
		if len(token) == 0 {
			break
		}
		v, decodedCount := decodeElement(token, ctx.TextPos)
		token[decodedCount:].AssertAtEnd(ctx.TextPos, ctx.ArrayDigitType)
		f32bits := math.Float32bits(float32(v))
		if f32bits&0xffff != 0 {
			ctx.Errorf("%x is too big to fit into float16", float32(v))
		}
		f16bits := f32bits >> 16
		exp := common.Float64GetExponent(v)
		if exp < common.Float32ExponentMin || exp > common.Float32ExponentMax {
			if math.IsNaN(v) {
				if common.HasQuietNanBitSet64(v) {
					f16bits = uint32(common.Bfloat16QuietNanBits)
				} else {
					f16bits = uint32(common.Bfloat16SignalingNanBits)
				}
			} else if !math.IsInf(v, 0) {
				ctx.Errorf("Exponent too big for float32 type")
			}
		}
		ctx.Scratch = append(ctx.Scratch, uint8(f16bits), uint8(f16bits>>8))
	}
	finishTypedArray(ctx)
}

func decodeArrayF32(ctx *DecoderContext, digitType string, decodeElement floatTokenDecoder) {
	ctx.BeginArray(digitType, events.ArrayTypeFloat32, 4)
	for {
		ctx.Stream.SkipWhitespace()
		if tryDecodeCommentInArray(ctx) {
			continue
		}
		token := ctx.Stream.ReadToken()
		if len(token) == 0 {
			break
		}
		v, decodedCount := decodeElement(token, ctx.TextPos)
		token[decodedCount:].AssertAtEnd(ctx.TextPos, ctx.ArrayDigitType)
		exp := common.Float64GetExponent(v)
		if exp < common.Float32ExponentMin || exp > common.Float32ExponentMax {
			if !math.IsNaN(v) && !math.IsInf(v, 0) {
				ctx.Errorf("Exponent too big for float32 type")
			}
		}
		bits := math.Float32bits(float32(v))
		if math.IsNaN(v) {
			if common.HasQuietNanBitSet64(v) {
				bits = common.Float32QuietNanBits
			} else {
				bits = common.Float32SignalingNanBits
			}
		}
		ctx.Scratch = append(ctx.Scratch, uint8(bits), uint8(bits>>8), uint8(bits>>16), uint8(bits>>24))
	}
	finishTypedArray(ctx)
}

func decodeArrayF64(ctx *DecoderContext, digitType string, decodeElement floatTokenDecoder) {
	ctx.BeginArray(digitType, events.ArrayTypeFloat64, 8)
	for {
		ctx.Stream.SkipWhitespace()
		if tryDecodeCommentInArray(ctx) {
			continue
		}
		token := ctx.Stream.ReadToken()
		if len(token) == 0 {
			break
		}
		v, decodedCount := decodeElement(token, ctx.TextPos)
		token[decodedCount:].AssertAtEnd(ctx.TextPos, ctx.ArrayDigitType)
		bits := math.Float64bits(v)
		ctx.Scratch = append(ctx.Scratch, uint8(bits), uint8(bits>>8), uint8(bits>>16), uint8(bits>>24),
			uint8(bits>>32), uint8(bits>>40), uint8(bits>>48), uint8(bits>>56))
	}
	finishTypedArray(ctx)
}

func decodeArrayUID(ctx *DecoderContext) {
	ctx.BeginArray("uid", events.ArrayTypeUID, 16)
	for {
		ctx.Stream.SkipWhitespace()
		if ctx.Stream.PeekByteAllowEOF() == '|' {
			break
		}
		if tryDecodeCommentInArray(ctx) {
			continue
		}
		token := ctx.Stream.ReadToken()
		ctx.Scratch = append(ctx.Scratch, token.DecodeUID(ctx.TextPos)...)
	}
	finishTypedArray(ctx)
}

func tryDecodeCommentInArray(ctx *DecoderContext) bool {
	if ctx.Stream.PeekByteAllowEOF() != '/' {
		return false
	}
	ctx.Stream.AdvanceByte() // Advance past '/'

	if !ctx.ArrayContainsComments {
		ctx.EventReceiver.OnArrayBegin(ctx.ArrayType)
	}
	if len(ctx.Scratch) > 0 {
		elemCount := uint64(len(ctx.Scratch) / ctx.ArrayBytesPerElement)
		ctx.EventReceiver.OnArrayChunk(elemCount, true)
		ctx.EventReceiver.OnArrayData(ctx.Scratch)
		ctx.Scratch = ctx.Scratch[:0]
	}
	ctx.ArrayContainsComments = true

	b := ctx.Stream.ReadByteNoEOF()
	switch b {
	case '/':
		contents := ctx.Stream.ReadSingleLineComment()
		ctx.EventReceiver.OnComment(false, contents)
	case '*':
		contents := ctx.Stream.ReadMultiLineComment()
		ctx.EventReceiver.OnComment(true, contents)
	default:
		ctx.Errorf("Unexpected comment initiator: [%c]", b)
	}
	return true
}

func advanceAndDecodeComment(ctx *DecoderContext) {
	ctx.AssertHasStructuralWS()
	ctx.Stream.AdvanceByte() // Advance past '/'

	b := ctx.Stream.ReadByteNoEOF()
	switch b {
	case '/':
		contents := ctx.Stream.ReadSingleLineComment()
		ctx.EventReceiver.OnComment(false, contents)
		ctx.StackDecoder(decodePostInvisible)
	case '*':
		contents := ctx.Stream.ReadMultiLineComment()
		ctx.EventReceiver.OnComment(true, contents)
		ctx.StackDecoder(decodePostInvisible)
	default:
		ctx.Errorf("Unexpected comment initiator: [%c]", b)
	}
}
