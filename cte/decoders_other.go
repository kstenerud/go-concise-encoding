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

	"github.com/kstenerud/go-concise-encoding/internal/chars"
	"github.com/kstenerud/go-concise-encoding/internal/common"
)

func decodeInvalidChar(ctx *DecoderContext) {
	ctx.Stream.Errorf("Unexpected [%v]", ctx.Stream.DescribeCurrentChar())
}

func decodeWhitespace(ctx *DecoderContext) {
	ctx.Stream.SkipWhitespace()
	return
}

func decodeByFirstChar(ctx *DecoderContext) {
	decodeWhitespace(ctx)
	b := ctx.Stream.PeekByteAllowEOD()
	decoderFunc := decoderFuncsByFirstChar[b]
	decoderFunc(ctx)
}

func decodePostInvisible(ctx *DecoderContext) {
	ctx.UnstackDecoder()
	decodeByFirstChar(ctx)
}

func decodeDocumentBegin(ctx *DecoderContext) {
	// Technically disallowed, but we'll support it anyway.
	decodeWhitespace(ctx)

	if b := ctx.Stream.ReadByteNoEOD(); b != 'c' && b != 'C' {
		ctx.Stream.Errorf(`Expected document to begin with "c" but got [%v]`, ctx.Stream.DescribeCurrentChar())
	}

	version, bigVersion, digitCount := ctx.Stream.DecodeDecimalUint(0, nil)
	if digitCount == 0 {
		ctx.Stream.UnexpectedChar("version number")
	}
	if bigVersion != nil {
		ctx.Stream.Errorf("Version too big")
	}

	b := ctx.Stream.PeekByteNoEOD()
	if !chars.ByteHasProperty(b, chars.CharIsWhitespace) {
		ctx.Stream.UnexpectedChar("whitespace after version")
	}
	decodeWhitespace(ctx)

	ctx.EventReceiver.OnBeginDocument()
	ctx.EventReceiver.OnVersion(version)
	ctx.ChangeDecoder(decodeTopLevel)
}

func decodeTopLevel(ctx *DecoderContext) {
	ctx.ChangeDecoder(decodeEndDocument)
	decodeByFirstChar(ctx)
}

func decodeEndDocument(ctx *DecoderContext) {
	ctx.EventReceiver.OnEndDocument()
	ctx.IsDocumentComplete = true
}

func decodeNumericPositive(ctx *DecoderContext) {
	coefficient, bigCoefficient, digitCount := ctx.Stream.DecodeDecimalUint(0, nil)
	b := ctx.Stream.ReadByteAllowEOD()
	switch b {
	case '-':
		v := ctx.Stream.DecodeDate(int64(coefficient))
		ctx.Stream.AssertAtObjectEnd("date")
		ctx.EventReceiver.OnCompactTime(v)
		return
	case ':':
		v := ctx.Stream.DecodeTime(int(coefficient))
		ctx.Stream.AssertAtObjectEnd("time")
		ctx.EventReceiver.OnCompactTime(v)
		return
	case '.':
		value, bigValue, _ := ctx.Stream.DecodeDecimalFloat(1, coefficient, bigCoefficient, digitCount)
		ctx.Stream.AssertAtObjectEnd("float")
		if bigValue != nil {
			ctx.EventReceiver.OnBigDecimalFloat(bigValue)
		} else {
			ctx.EventReceiver.OnDecimalFloat(value)
		}
		return
	default:
		if b.HasProperty(chars.CharIsObjectEnd) {
			ctx.Stream.UnreadByte()
			if bigCoefficient != nil {
				ctx.EventReceiver.OnBigInt(bigCoefficient)
			} else {
				ctx.EventReceiver.OnPositiveInt(coefficient)
			}
			return
		}
	}
	ctx.Stream.UnexpectedChar("numeric")
}

func advanceAndDecodeNumericNegative(ctx *DecoderContext) {
	ctx.Stream.AdvanceByte() // Advance past '-'

	switch ctx.Stream.ReadByteNoEOD() {
	case '0':
		decodeOtherBaseNegative(ctx)
		return
	case '@':
		namedValue := string(ctx.Stream.DecodeNamedValue())
		if namedValue != "inf" {
			ctx.Stream.Errorf("Unknown named value: %v", namedValue)
		}
		ctx.EventReceiver.OnFloat(math.Inf(-1))
		return
	default:
		ctx.Stream.UnreadByte()
	}

	coefficient, bigCoefficient, digitCount := ctx.Stream.DecodeDecimalUint(0, nil)
	b := ctx.Stream.ReadByteAllowEOD()
	switch b {
	case '-':
		v := ctx.Stream.DecodeDate(-int64(coefficient))
		ctx.Stream.AssertAtObjectEnd("time")
		ctx.EventReceiver.OnCompactTime(v)
		return
	case '.':
		value, bigValue, _ := ctx.Stream.DecodeDecimalFloat(-1, coefficient, bigCoefficient, digitCount)
		ctx.Stream.AssertAtObjectEnd("float")
		if bigValue != nil {
			ctx.EventReceiver.OnBigDecimalFloat(bigValue)
		} else {
			ctx.EventReceiver.OnDecimalFloat(value)
		}
		return
	default:
		if b.HasProperty(chars.CharIsObjectEnd) {
			ctx.Stream.UnreadByte()
			if bigCoefficient != nil {
				// TODO: More efficient way to negate?
				bigCoefficient = bigCoefficient.Mul(bigCoefficient, common.BigIntN1)
				ctx.EventReceiver.OnBigInt(bigCoefficient)
			} else {
				ctx.EventReceiver.OnNegativeInt(coefficient)
			}
			return
		}
	}
	ctx.Stream.UnexpectedChar("numeric")
}

func advanceAndDecodeOtherBasePositive(ctx *DecoderContext) {
	ctx.Stream.AdvanceByte() // Advance past '0'

	b := ctx.Stream.ReadByteAllowEOD()
	if b.HasProperty(chars.CharIsObjectEnd) {
		ctx.Stream.UnreadByte()
		ctx.EventReceiver.OnPositiveInt(0)
		return
	}

	switch b {
	case 'b':
		v, bigV, _ := ctx.Stream.DecodeBinaryUint()
		ctx.Stream.AssertAtObjectEnd("binary integer")
		if bigV != nil {
			ctx.EventReceiver.OnBigInt(bigV)
		} else {
			ctx.EventReceiver.OnPositiveInt(v)
		}
	case 'o':
		v, bigV, _ := ctx.Stream.DecodeOctalUint()
		ctx.Stream.AssertAtObjectEnd("octal integer")
		if bigV != nil {
			ctx.EventReceiver.OnBigInt(bigV)
		} else {
			ctx.EventReceiver.OnPositiveInt(v)
		}
	case 'x':
		v, bigV, digitCount := ctx.Stream.DecodeHexUint(0, nil)
		if ctx.Stream.PeekByteAllowEOD() == '.' {
			ctx.Stream.AdvanceByte()
			fv, bigFV, _ := ctx.Stream.DecodeHexFloat(1, v, bigV, digitCount)
			ctx.Stream.AssertAtObjectEnd("hex float")
			if bigFV != nil {
				ctx.EventReceiver.OnBigFloat(bigFV)
			} else {
				ctx.EventReceiver.OnFloat(fv)
			}
		} else {
			ctx.Stream.AssertAtObjectEnd("hex integer")
			if bigV != nil {
				ctx.EventReceiver.OnBigInt(bigV)
			} else {
				ctx.EventReceiver.OnPositiveInt(v)
			}
		}
	case '.':
		value, bigValue, _ := ctx.Stream.DecodeDecimalFloat(1, 0, nil, 0)
		ctx.Stream.AssertAtObjectEnd("float")
		if bigValue != nil {
			ctx.EventReceiver.OnBigDecimalFloat(bigValue)
		} else {
			ctx.EventReceiver.OnDecimalFloat(value)
		}
	default:
		if b.HasProperty(chars.CharIsDigitBase10) && ctx.Stream.PeekByteNoEOD() == ':' {
			ctx.Stream.AdvanceByte()
			v := ctx.Stream.DecodeTime(int(b - '0'))
			ctx.Stream.AssertAtObjectEnd("time")
			ctx.EventReceiver.OnCompactTime(v)
			return
		}
		ctx.Stream.UnreadByte()
		ctx.Stream.UnexpectedChar("numeric base")
	}
}

func advanceAndDecodeOtherBaseNegative(ctx *DecoderContext) {
	ctx.Stream.AdvanceByte() // Advance past '0'
	decodeOtherBaseNegative(ctx)
}

func decodeOtherBaseNegative(ctx *DecoderContext) {

	b := ctx.Stream.PeekByteAllowEOD()
	if b.HasProperty(chars.CharIsObjectEnd) {
		// -0 has no decimal point (thus type int), so report it as positive 0.
		ctx.EventReceiver.OnPositiveInt(0)
		return
	}

	ctx.Stream.AdvanceByte() // Advance past the value now stored in b

	switch b {
	case 'b':
		v, bigV, _ := ctx.Stream.DecodeBinaryUint()
		ctx.Stream.AssertAtObjectEnd("binary integer")
		if bigV != nil {
			bigV = bigV.Neg(bigV)
			ctx.EventReceiver.OnBigInt(bigV)
		} else {
			ctx.EventReceiver.OnNegativeInt(v)
		}
	case 'o':
		v, bigV, _ := ctx.Stream.DecodeOctalUint()
		ctx.Stream.AssertAtObjectEnd("octal integer")
		if bigV != nil {
			bigV = bigV.Neg(bigV)
			ctx.EventReceiver.OnBigInt(bigV)
		} else {
			ctx.EventReceiver.OnNegativeInt(v)
		}
	case 'x':
		v, bigV, digitCount := ctx.Stream.DecodeHexUint(0, nil)
		if ctx.Stream.PeekByteAllowEOD() == '.' {
			ctx.Stream.AdvanceByte()
			fv, bigFV, _ := ctx.Stream.DecodeHexFloat(-1, v, bigV, digitCount)
			ctx.Stream.AssertAtObjectEnd("hex float")
			if bigFV != nil {
				ctx.EventReceiver.OnBigFloat(bigFV)
			} else {
				ctx.EventReceiver.OnFloat(fv)
			}
		} else {
			ctx.Stream.AssertAtObjectEnd("hex integer")
			if bigV != nil {
				bigV = bigV.Neg(bigV)
				ctx.EventReceiver.OnBigInt(bigV)
			} else {
				ctx.EventReceiver.OnNegativeInt(v)
			}
		}
	case '.':
		value, bigValue, _ := ctx.Stream.DecodeDecimalFloat(-1, 0, nil, 0)
		ctx.Stream.AssertAtObjectEnd("float")
		if bigValue != nil {
			ctx.EventReceiver.OnBigDecimalFloat(bigValue)
		} else {
			ctx.EventReceiver.OnDecimalFloat(value)
		}
	default:
		ctx.Stream.UnreadByte()
		ctx.Stream.UnexpectedChar("numeric base")
	}
}

func advanceAndDecodeNamedValueOrUUID(ctx *DecoderContext) {
	ctx.Stream.AdvanceByte() // Advance past '@'

	namedValue := ctx.Stream.DecodeNamedValue()
	switch string(namedValue) {
	case "na":
		ctx.EventReceiver.OnNA()
		if ctx.Stream.PeekByteAllowEOD() == ':' {
			ctx.Stream.AdvanceByte()
			decodeByFirstChar(ctx)
		} else {
			ctx.EventReceiver.OnNA()
		}
	case "nan":
		ctx.EventReceiver.OnNan(false)
	case "snan":
		ctx.EventReceiver.OnNan(true)
	case "inf":
		ctx.EventReceiver.OnFloat(math.Inf(1))
	case "false":
		ctx.EventReceiver.OnFalse()
	case "true":
		ctx.EventReceiver.OnTrue()
	default:
		ctx.EventReceiver.OnUUID(ctx.Stream.ExtractUUID(namedValue))
	}
}

func advanceAndDecodeConstant(ctx *DecoderContext) {
	ctx.Stream.AdvanceByte() // Advance past '#'

	name := ctx.Stream.DecodeUnquotedString()
	if ctx.Stream.PeekByteAllowEOD() == ':' {
		ctx.EventReceiver.OnConstant(name, true)
		ctx.Stream.AdvanceByte()
		decodeByFirstChar(ctx)
	} else {
		ctx.EventReceiver.OnConstant(name, false)
	}
}

func advanceAndDecodeMarker(ctx *DecoderContext) {
	ctx.Stream.AdvanceByte() // Advance past '&'

	ctx.EventReceiver.OnMarker()

	asString, asUint := ctx.Stream.DecodeMarkerID()
	if len(asString) > 0 {
		ctx.EventReceiver.OnArray(events.ArrayTypeString, uint64(len(asString)), asString)
	} else {
		ctx.EventReceiver.OnPositiveInt(asUint)
	}
	if ctx.Stream.PeekByteNoEOD() != ':' {
		ctx.Stream.Errorf("Missing colon between marker ID and marked value")
	}
	ctx.Stream.AdvanceByte()
	decodeByFirstChar(ctx)
}

func advanceAndDecodeReference(ctx *DecoderContext) {
	ctx.Stream.AdvanceByte() // Advance past '$'

	ctx.EventReceiver.OnReference()

	if ctx.Stream.PeekByteNoEOD() == '|' {
		ctx.Stream.AdvanceByte()
		arrayType := decodeArrayType(ctx)
		ctx.Stream.SkipWhitespace()
		if arrayType != "r" {
			ctx.Stream.Errorf("%s: Invalid array type for reference ID", arrayType)
		}
		decodeRID(ctx)
		return
	}

	asString, asUint := ctx.Stream.DecodeMarkerID()
	if len(asString) > 0 {
		ctx.EventReceiver.OnArray(events.ArrayTypeString, uint64(len(asString)), asString)
	} else {
		ctx.EventReceiver.OnPositiveInt(asUint)
	}
}

func advanceAndDecodeSuffix(ctx *DecoderContext) {
	ctx.Stream.AdvanceByte() // Advance past ':'

	panic("TODO: decodeSuffix")
}
