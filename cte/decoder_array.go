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
		decodeArrayU8(ctx, "integer", ctx.Stream.ReadSmallUint)
	case "u8b":
		decodeArrayU8(ctx, "binary", ctx.Stream.ReadSmallBinaryUint)
	case "u8o":
		decodeArrayU8(ctx, "octal", ctx.Stream.ReadSmallOctalUint)
	case "u8x":
		decodeArrayU8(ctx, "hex", ctx.Stream.ReadSmallHexUint)
	case "u16":
		decodeArrayU16(ctx, "integer", ctx.Stream.ReadSmallUint)
	case "u16b":
		decodeArrayU16(ctx, "binary", ctx.Stream.ReadSmallBinaryUint)
	case "u16o":
		decodeArrayU16(ctx, "octal", ctx.Stream.ReadSmallOctalUint)
	case "u16x":
		decodeArrayU16(ctx, "hex", ctx.Stream.ReadSmallHexUint)
	case "u32":
		decodeArrayU32(ctx, "integer", ctx.Stream.ReadSmallUint)
	case "u32b":
		decodeArrayU32(ctx, "binary", ctx.Stream.ReadSmallBinaryUint)
	case "u32o":
		decodeArrayU32(ctx, "octal", ctx.Stream.ReadSmallOctalUint)
	case "u32x":
		decodeArrayU32(ctx, "hex", ctx.Stream.ReadSmallHexUint)
	case "u64":
		decodeArrayU64(ctx, "integer", ctx.Stream.ReadSmallUint)
	case "u64b":
		decodeArrayU64(ctx, "binary", ctx.Stream.ReadSmallBinaryUint)
	case "u64o":
		decodeArrayU64(ctx, "octal", ctx.Stream.ReadSmallOctalUint)
	case "u64x":
		decodeArrayU64(ctx, "hex", ctx.Stream.ReadSmallHexUint)
	case "i8":
		decodeArrayI8(ctx, "integer", ctx.Stream.ReadSmallInt)
	case "i8b":
		decodeArrayI8(ctx, "binary", ctx.Stream.ReadSmallBinaryInt)
	case "i8o":
		decodeArrayI8(ctx, "octal", ctx.Stream.ReadSmallOctalInt)
	case "i8x":
		decodeArrayI8(ctx, "hex", ctx.Stream.ReadSmallHexInt)
	case "i16":
		decodeArrayI16(ctx, "integer", ctx.Stream.ReadSmallInt)
	case "i16b":
		decodeArrayI16(ctx, "binary", ctx.Stream.ReadSmallBinaryInt)
	case "i16o":
		decodeArrayI16(ctx, "octal", ctx.Stream.ReadSmallOctalInt)
	case "i16x":
		decodeArrayI16(ctx, "hex", ctx.Stream.ReadSmallHexInt)
	case "i32":
		decodeArrayI32(ctx, "integer", ctx.Stream.ReadSmallInt)
	case "i32b":
		decodeArrayI32(ctx, "binary", ctx.Stream.ReadSmallBinaryInt)
	case "i32o":
		decodeArrayI32(ctx, "octal", ctx.Stream.ReadSmallOctalInt)
	case "i32x":
		decodeArrayI32(ctx, "hex", ctx.Stream.ReadSmallHexInt)
	case "i64":
		decodeArrayI64(ctx, "integer", ctx.Stream.ReadSmallInt)
	case "i64b":
		decodeArrayI64(ctx, "binary", ctx.Stream.ReadSmallBinaryInt)
	case "i64o":
		decodeArrayI64(ctx, "octal", ctx.Stream.ReadSmallOctalInt)
	case "i64x":
		decodeArrayI64(ctx, "hex", ctx.Stream.ReadSmallHexInt)
	case "f16":
		decodeArrayF16(ctx, "float", ctx.Stream.ReadSmallFloat)
	case "f16x":
		decodeArrayF16(ctx, "hex float", ctx.Stream.ReadSmallHexFloat)
	case "f32":
		decodeArrayF32(ctx, "float", ctx.Stream.ReadSmallFloat)
	case "f32x":
		decodeArrayF32(ctx, "hex float", ctx.Stream.ReadSmallHexFloat)
	case "f64":
		decodeArrayF64(ctx, "float", ctx.Stream.ReadSmallFloat)
	case "f64x":
		decodeArrayF64(ctx, "hex float", ctx.Stream.ReadSmallHexFloat)
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
	var data []uint8
	for {
		ctx.Stream.SkipWhitespace()
		v, count := ctx.Stream.ReadSmallHexUint()
		if count == 0 {
			break
		}
		if v > maxUint8Value {
			ctx.Errorf("%v value too big for array type", digitType)
		}
		data = append(data, uint8(v))
	}
	finishTypedArray(ctx, events.ArrayTypeCustomBinary, digitType, 1, data)
}

func decodeArrayBoolean(ctx *DecoderContext) {
	var data []uint8
	var elemCount uint64

	for {
		nextByte := byte(0)
		for i := 0; i < 8; i++ {
			ctx.Stream.SkipWhitespace()
			b := ctx.Stream.ReadByteNoEOF()
			switch b {
			case '|':
				if i > 0 {
					data = append(data, nextByte)
				}
				ctx.EventReceiver.OnArray(events.ArrayTypeBoolean, elemCount, data)
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
		data = append(data, nextByte)
	}
}

func decodeArrayU8(ctx *DecoderContext, digitType string, decodeElement func() (v uint64, digitCount int)) {
	var data []uint8
	for {
		ctx.Stream.SkipWhitespace()
		v, count := decodeElement()
		if count == 0 {
			break
		}
		if v > maxUint8Value {
			ctx.Errorf("%v value too big for array type", digitType)
		}
		data = append(data, uint8(v))
	}
	finishTypedArray(ctx, events.ArrayTypeUint8, digitType, 1, data)
}

func decodeArrayU16(ctx *DecoderContext, digitType string, decodeElement func() (v uint64, digitCount int)) {
	var data []uint8
	for {
		ctx.Stream.SkipWhitespace()
		v, count := decodeElement()
		if count == 0 {
			break
		}
		if v > maxUint16Value {
			ctx.Errorf("%v value too big for array type", digitType)
		}
		data = append(data, uint8(v), uint8(v>>8))
	}
	finishTypedArray(ctx, events.ArrayTypeUint16, digitType, 2, data)
}

func decodeArrayU32(ctx *DecoderContext, digitType string, decodeElement func() (v uint64, digitCount int)) {
	var data []uint8
	for {
		ctx.Stream.SkipWhitespace()
		v, count := decodeElement()
		if count == 0 {
			break
		}
		if v > maxUint32Value {
			ctx.Errorf("%v value too big for array type", digitType)
		}
		data = append(data, uint8(v), uint8(v>>8), uint8(v>>16), uint8(v>>24))
	}
	finishTypedArray(ctx, events.ArrayTypeUint32, digitType, 4, data)
}

func decodeArrayU64(ctx *DecoderContext, digitType string, decodeElement func() (v uint64, digitCount int)) {
	var data []uint8
	for {
		ctx.Stream.SkipWhitespace()
		v, count := decodeElement()
		if count == 0 {
			break
		}
		data = append(data, uint8(v), uint8(v>>8), uint8(v>>16), uint8(v>>24),
			uint8(v>>32), uint8(v>>40), uint8(v>>48), uint8(v>>56))
	}
	finishTypedArray(ctx, events.ArrayTypeUint64, digitType, 8, data)
}

func decodeArrayI8(ctx *DecoderContext, digitType string, decodeElement func() (v int64, digitCount int)) {
	var data []uint8
	for {
		ctx.Stream.SkipWhitespace()
		v, count := decodeElement()
		if count == 0 {
			break
		}
		if v < minInt8Value || v > maxInt8Value {
			ctx.Errorf("%v value too big for array type", digitType)
		}
		data = append(data, uint8(v))
	}
	finishTypedArray(ctx, events.ArrayTypeInt8, digitType, 1, data)
}

func decodeArrayI16(ctx *DecoderContext, digitType string, decodeElement func() (v int64, digitCount int)) {
	var data []uint8
	for {
		ctx.Stream.SkipWhitespace()
		v, count := decodeElement()
		if count == 0 {
			break
		}
		if v < minInt16Value || v > maxInt16Value {
			ctx.Errorf("%v value too big for array type", digitType)
		}
		data = append(data, uint8(v), uint8(v>>8))
	}
	finishTypedArray(ctx, events.ArrayTypeInt16, digitType, 2, data)
}

func decodeArrayI32(ctx *DecoderContext, digitType string, decodeElement func() (v int64, digitCount int)) {
	var data []uint8
	for {
		ctx.Stream.SkipWhitespace()
		v, count := decodeElement()
		if count == 0 {
			break
		}
		if v < minInt32Value || v > maxInt32Value {
			ctx.Errorf("%v value too big for array type", digitType)
		}
		data = append(data, uint8(v), uint8(v>>8), uint8(v>>16), uint8(v>>24))
	}
	finishTypedArray(ctx, events.ArrayTypeInt32, digitType, 4, data)
}

func decodeArrayI64(ctx *DecoderContext, digitType string, decodeElement func() (v int64, digitCount int)) {
	var data []uint8
	for {
		ctx.Stream.SkipWhitespace()
		v, count := decodeElement()
		if count == 0 {
			break
		}
		data = append(data, uint8(v), uint8(v>>8), uint8(v>>16), uint8(v>>24),
			uint8(v>>32), uint8(v>>40), uint8(v>>48), uint8(v>>56))
	}
	finishTypedArray(ctx, events.ArrayTypeInt64, digitType, 8, data)
}

func decodeArrayF16(ctx *DecoderContext, digitType string, decodeElement func() (v float64, digitCount int)) {
	var data []uint8
	for {
		ctx.Stream.SkipWhitespace()
		v, count := decodeElement()
		if count == 0 {
			break
		}

		exp := extractFloat64Exponent(v)
		if exp < minFloat32Exponent || exp > maxFloat32Exponent {
			ctx.Errorf("Exponent too big for bfloat16 type")
		}
		bits := math.Float32bits(float32(v))
		data = append(data, uint8(bits>>16), uint8(bits>>24))
	}
	finishTypedArray(ctx, events.ArrayTypeFloat16, digitType, 2, data)
}

func decodeArrayF32(ctx *DecoderContext, digitType string, decodeElement func() (v float64, digitCount int)) {
	var data []uint8
	for {
		ctx.Stream.SkipWhitespace()
		v, count := decodeElement()
		if count == 0 {
			break
		}

		exp := extractFloat64Exponent(v)
		if exp < minFloat32Exponent || exp > maxFloat32Exponent {
			ctx.Errorf("Exponent too big for float32 type")
		}
		bits := math.Float32bits(float32(v))
		data = append(data, uint8(bits), uint8(bits>>8), uint8(bits>>16), uint8(bits>>24))
	}
	finishTypedArray(ctx, events.ArrayTypeFloat32, digitType, 4, data)
}

func decodeArrayF64(ctx *DecoderContext, digitType string, decodeElement func() (v float64, digitCount int)) {
	var data []uint8
	for {
		ctx.Stream.SkipWhitespace()
		v, count := decodeElement()
		if count == 0 {
			break
		}
		bits := math.Float64bits(v)
		data = append(data, uint8(bits), uint8(bits>>8), uint8(bits>>16), uint8(bits>>24),
			uint8(bits>>32), uint8(bits>>40), uint8(bits>>48), uint8(bits>>56))
	}
	finishTypedArray(ctx, events.ArrayTypeFloat64, digitType, 8, data)
}

func decodeArrayUUID(ctx *DecoderContext) {
	var data []uint8
	ctx.Stream.SkipWhitespace()
	for ctx.Stream.PeekByteAllowEOF() != '|' {
		data = append(data, ctx.Stream.ReadUUIDWithDecimalDecoded(0, 0)...)
		ctx.Stream.SkipWhitespace()
	}
	finishTypedArray(ctx, events.ArrayTypeUUID, "uuid", 16, data)
}
