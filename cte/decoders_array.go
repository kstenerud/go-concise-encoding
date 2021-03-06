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

	"github.com/kstenerud/go-concise-encoding/internal/common"

	"github.com/kstenerud/go-concise-encoding/events"
	"github.com/kstenerud/go-concise-encoding/internal/chars"
)

func advanceAndDecodeQuotedString(ctx *DecoderContext) {
	ctx.Stream.AdvanceByte() // Advance past '"'

	bytes := ctx.Stream.DecodeQuotedString()
	ctx.EventReceiver.OnArray(events.ArrayTypeString, uint64(len(bytes)), bytes)
}

func decodeUnquotedString(ctx *DecoderContext) {
	bytes := ctx.Stream.DecodeUnquotedString()
	ctx.EventReceiver.OnArray(events.ArrayTypeString, uint64(len(bytes)), bytes)
}

func decodeArrayType(ctx *DecoderContext) string {
	ctx.Stream.TokenBegin()
	ctx.Stream.TokenReadUntilPropertyNoEOD(chars.CharIsObjectEnd)
	arrayType := ctx.Stream.TokenGet()
	if len(arrayType) > 0 && arrayType[len(arrayType)-1] == '|' {
		arrayType = arrayType[:len(arrayType)-1]
		ctx.Stream.UnreadByte()
	}
	common.ASCIIBytesToLower(arrayType)
	return string(arrayType)
}

func finishTypedArray(ctx *DecoderContext, arrayType events.ArrayType, digitType string, bytesPerElement int, data []byte) {
	switch ctx.Stream.ReadByteNoEOD() {
	case '|':
		ctx.EventReceiver.OnArray(arrayType, uint64(len(data)/bytesPerElement), data)
		return
	default:
		ctx.Stream.Errorf("Expected %v digits", digitType)
	}
}

func decodeTypedArrayBegin(ctx *DecoderContext) {
	ctx.Stream.AdvanceByte() // Advance past '|'

	arrayType := decodeArrayType(ctx)
	ctx.Stream.SkipWhitespace()
	switch arrayType {
	case "cb":
		decodeCustomBinary(ctx)
	case "ct":
		decodeCustomText(ctx)
	case "r":
		decodeRID(ctx)
	case "u":
		decodeArrayUUID(ctx)
	case "u8":
		decodeArrayU8(ctx, "integer", ctx.Stream.DecodeSmallUint)
	case "u8b":
		decodeArrayU8(ctx, "binary", ctx.Stream.DecodeSmallBinaryUint)
	case "u8o":
		decodeArrayU8(ctx, "octal", ctx.Stream.DecodeSmallOctalUint)
	case "u8x":
		decodeArrayU8(ctx, "hex", ctx.Stream.DecodeSmallHexUint)
	case "u16":
		decodeArrayU16(ctx, "integer", ctx.Stream.DecodeSmallUint)
	case "u16b":
		decodeArrayU16(ctx, "binary", ctx.Stream.DecodeSmallBinaryUint)
	case "u16o":
		decodeArrayU16(ctx, "octal", ctx.Stream.DecodeSmallOctalUint)
	case "u16x":
		decodeArrayU16(ctx, "hex", ctx.Stream.DecodeSmallHexUint)
	case "u32":
		decodeArrayU32(ctx, "integer", ctx.Stream.DecodeSmallUint)
	case "u32b":
		decodeArrayU32(ctx, "binary", ctx.Stream.DecodeSmallBinaryUint)
	case "u32o":
		decodeArrayU32(ctx, "octal", ctx.Stream.DecodeSmallOctalUint)
	case "u32x":
		decodeArrayU32(ctx, "hex", ctx.Stream.DecodeSmallHexUint)
	case "u64":
		decodeArrayU64(ctx, "integer", ctx.Stream.DecodeSmallUint)
	case "u64b":
		decodeArrayU64(ctx, "binary", ctx.Stream.DecodeSmallBinaryUint)
	case "u64o":
		decodeArrayU64(ctx, "octal", ctx.Stream.DecodeSmallOctalUint)
	case "u64x":
		decodeArrayU64(ctx, "hex", ctx.Stream.DecodeSmallHexUint)
	case "i8":
		decodeArrayI8(ctx, "integer", ctx.Stream.DecodeSmallInt)
	case "i8b":
		decodeArrayI8(ctx, "binary", ctx.Stream.DecodeSmallBinaryInt)
	case "i8o":
		decodeArrayI8(ctx, "octal", ctx.Stream.DecodeSmallOctalInt)
	case "i8x":
		decodeArrayI8(ctx, "hex", ctx.Stream.DecodeSmallHexInt)
	case "i16":
		decodeArrayI16(ctx, "integer", ctx.Stream.DecodeSmallInt)
	case "i16b":
		decodeArrayI16(ctx, "binary", ctx.Stream.DecodeSmallBinaryInt)
	case "i16o":
		decodeArrayI16(ctx, "octal", ctx.Stream.DecodeSmallOctalInt)
	case "i16x":
		decodeArrayI16(ctx, "hex", ctx.Stream.DecodeSmallHexInt)
	case "i32":
		decodeArrayI32(ctx, "integer", ctx.Stream.DecodeSmallInt)
	case "i32b":
		decodeArrayI32(ctx, "binary", ctx.Stream.DecodeSmallBinaryInt)
	case "i32o":
		decodeArrayI32(ctx, "octal", ctx.Stream.DecodeSmallOctalInt)
	case "i32x":
		decodeArrayI32(ctx, "hex", ctx.Stream.DecodeSmallHexInt)
	case "i64":
		decodeArrayI64(ctx, "integer", ctx.Stream.DecodeSmallInt)
	case "i64b":
		decodeArrayI64(ctx, "binary", ctx.Stream.DecodeSmallBinaryInt)
	case "i64o":
		decodeArrayI64(ctx, "octal", ctx.Stream.DecodeSmallOctalInt)
	case "i64x":
		decodeArrayI64(ctx, "hex", ctx.Stream.DecodeSmallHexInt)
	case "f16":
		decodeArrayF16(ctx, "float", ctx.Stream.DecodeSmallFloat)
	case "f16x":
		decodeArrayF16(ctx, "hex float", ctx.Stream.DecodeSmallHexFloat)
	case "f32":
		decodeArrayF32(ctx, "float", ctx.Stream.DecodeSmallFloat)
	case "f32x":
		decodeArrayF32(ctx, "hex float", ctx.Stream.DecodeSmallHexFloat)
	case "f64":
		decodeArrayF64(ctx, "float", ctx.Stream.DecodeSmallFloat)
	case "f64x":
		decodeArrayF64(ctx, "hex float", ctx.Stream.DecodeSmallHexFloat)
	default:
		ctx.Stream.Errorf("%s: Unhandled array type", arrayType)
	}
}

func decodeStringArray(ctx *DecoderContext, arrayType events.ArrayType) {
	bytes := ctx.Stream.DecodeStringArray()
	ctx.EventReceiver.OnArray(arrayType, uint64(len(bytes)), bytes)
}

func decodeCustomText(ctx *DecoderContext) {
	decodeStringArray(ctx, events.ArrayTypeCustomText)
}

func decodeRID(ctx *DecoderContext) {
	decodeStringArray(ctx, events.ArrayTypeResourceID)
}

func decodeCustomBinary(ctx *DecoderContext) {
	digitType := "hex"
	var data []uint8
	for {
		ctx.Stream.SkipWhitespace()
		v, count := ctx.Stream.DecodeSmallHexUint()
		if count == 0 {
			break
		}
		if v > maxUint8Value {
			ctx.Stream.Errorf("%v value too big for array type", digitType)
		}
		data = append(data, uint8(v))
	}
	finishTypedArray(ctx, events.ArrayTypeCustomBinary, digitType, 1, data)
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
			ctx.Stream.Errorf("%v value too big for array type", digitType)
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
			ctx.Stream.Errorf("%v value too big for array type", digitType)
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
			ctx.Stream.Errorf("%v value too big for array type", digitType)
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
			ctx.Stream.Errorf("%v value too big for array type", digitType)
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
			ctx.Stream.Errorf("%v value too big for array type", digitType)
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
			ctx.Stream.Errorf("%v value too big for array type", digitType)
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
			ctx.Stream.Errorf("Exponent too big for bfloat16 type")
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
			ctx.Stream.Errorf("Exponent too big for float32 type")
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
	for ctx.Stream.PeekByteAllowEOD() != '|' {
		token := ctx.Stream.DecodeToken()
		if len(token) == 0 {
			panic("Error")
		}
		if token[0] == '@' {
			token = token[1:]
		}
		data = append(data, ctx.Stream.ExtractUUID(token)...)
		ctx.Stream.SkipWhitespace()
	}
	finishTypedArray(ctx, events.ArrayTypeUUID, "uuid", 16, data)
}
