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

func decodeInvalidChar(ctx *DecoderContext) {
	ctx.Stream.Errorf("Unexpected [%v]", ctx.Stream.DescribeCurrentChar())
}

func decodeWhitespace(ctx *DecoderContext) {
	ctx.Stream.SkipWhitespace()
	return
}

func decodeByFirstChar(ctx *DecoderContext) {
	decodeWhitespace(ctx)
	decoderFunc := decoderFuncsByFirstChar[ctx.Stream.PeekByteAllowEOD()]
	decoderFunc(ctx)
}

func decodeDocumentBegin(ctx *DecoderContext) {
	// Technically disallowed, but we'll support it anyway.
	decodeWhitespace(ctx)

	if b := ctx.Stream.PeekByteNoEOD(); b != 'c' && b != 'C' {
		ctx.Stream.Errorf(`Expected document to begin with "c" but got [%v]`, ctx.Stream.DescribeCurrentChar())
	}

	ctx.Stream.AdvanceByte()

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
	b := ctx.Stream.PeekByteAllowEOD()
	switch b {
	case '-':
		ctx.Stream.AdvanceByte()
		v := ctx.Stream.DecodeDate(int64(coefficient))
		ctx.Stream.AssertAtObjectEnd("date")
		ctx.EventReceiver.OnCompactTime(v)
		return
	case ':':
		ctx.Stream.AdvanceByte()
		v := ctx.Stream.DecodeTime(int(coefficient))
		ctx.Stream.AssertAtObjectEnd("time")
		ctx.EventReceiver.OnCompactTime(v)
		return
	case '.':
		ctx.Stream.AdvanceByte()
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

func decodeNumericNegative(ctx *DecoderContext) {
	ctx.Stream.AdvanceByte() // Advance past '-'

	switch ctx.Stream.PeekByteNoEOD() {
	case '0':
		decodeOtherBaseNegative(ctx)
		return
	case '@':
		ctx.Stream.AdvanceByte()
		namedValue := string(ctx.Stream.DecodeNamedValue())
		if namedValue != "inf" {
			ctx.Stream.Errorf("Unknown named value: %v", namedValue)
		}
		ctx.EventReceiver.OnFloat(math.Inf(-1))
		return
	}

	coefficient, bigCoefficient, digitCount := ctx.Stream.DecodeDecimalUint(0, nil)
	b := ctx.Stream.PeekByteAllowEOD()
	switch b {
	case '-':
		ctx.Stream.AdvanceByte()
		v := ctx.Stream.DecodeDate(-int64(coefficient))
		ctx.Stream.AssertAtObjectEnd("time")
		ctx.EventReceiver.OnCompactTime(v)
		return
	case '.':
		ctx.Stream.AdvanceByte()
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

func decodeOtherBasePositive(ctx *DecoderContext) {
	ctx.Stream.AdvanceByte() // Advance past '0'

	b := ctx.Stream.PeekByteAllowEOD()
	if b.HasProperty(chars.CharIsObjectEnd) {
		ctx.EventReceiver.OnPositiveInt(0)
		return
	}

	ctx.Stream.AdvanceByte() // Advance past the value now stored in b

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
		ctx.Stream.UngetByte()
		ctx.Stream.UnexpectedChar("numeric base")
	}
}

func decodeOtherBaseNegative(ctx *DecoderContext) {
	ctx.Stream.AdvanceByte() // Advance past '0'

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
		ctx.Stream.UngetByte()
		ctx.Stream.UnexpectedChar("numeric base")
	}
}

func decodeNamedValueOrUUID(ctx *DecoderContext) {
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
		return
	case "nan":
		ctx.EventReceiver.OnNan(false)
		return
	case "snan":
		ctx.EventReceiver.OnNan(true)
		return
	case "inf":
		ctx.EventReceiver.OnFloat(math.Inf(1))
		return
	case "false":
		ctx.EventReceiver.OnFalse()
		return
	case "true":
		ctx.EventReceiver.OnTrue()
		return
	}

	// UUID
	if len(namedValue) != 36 ||
		namedValue[8] != '-' ||
		namedValue[13] != '-' ||
		namedValue[18] != '-' ||
		namedValue[23] != '-' {
		ctx.Stream.UngetBytes(len(namedValue) + 1)
		ctx.Stream.Errorf("Malformed UUID or unknown named value: [%s]", string(namedValue))
	}

	decodeHex := func(b byte) byte {
		switch {
		case chars.ByteHasProperty(b, chars.CharIsDigitBase10):
			return byte(b - '0')
		case chars.ByteHasProperty(b, chars.CharIsLowerAF):
			return byte(b - 'a' + 10)
		case chars.ByteHasProperty(b, chars.CharIsUpperAF):
			return byte(b - 'A' + 10)
		default:
			ctx.Stream.UngetBytes(len(namedValue) + 1)
			ctx.Stream.Errorf("Unexpected char [%c] in UUID [%s]", b, string(namedValue))
			return 0
		}
	}

	decodeSection := func(src []byte, dst []byte) {
		iSrc := 0
		iDst := 0
		for iSrc < len(src) {
			dst[iDst] = (decodeHex(src[iSrc]) << 4) | decodeHex(src[iSrc+1])
			iDst++
			iSrc += 2
		}
	}

	decodeSection(namedValue[:8], namedValue)
	decodeSection(namedValue[9:13], namedValue[4:])
	decodeSection(namedValue[14:18], namedValue[6:])
	decodeSection(namedValue[19:23], namedValue[8:])
	decodeSection(namedValue[24:36], namedValue[10:])

	ctx.EventReceiver.OnUUID(namedValue[:16])
}

func decodeQuotedString(ctx *DecoderContext) {
	ctx.Stream.AdvanceByte() // Advance past '"'

	bytes := ctx.Stream.DecodeQuotedString()
	ctx.EventReceiver.OnArray(events.ArrayTypeString, uint64(len(bytes)), bytes)
}

func decodeUnquotedString(ctx *DecoderContext) {
	bytes := ctx.Stream.DecodeUnquotedString()
	ctx.EventReceiver.OnArray(events.ArrayTypeString, uint64(len(bytes)), bytes)
}

func decodeMapBegin(ctx *DecoderContext) {
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

func decodeMapEnd(ctx *DecoderContext) {
	ctx.Stream.AdvanceByte() // Advance past '}'

	ctx.EventReceiver.OnEnd()
	ctx.UnstackDecoder()
}

func decodeListBegin(ctx *DecoderContext) {
	ctx.Stream.AdvanceByte() // Advance past '['

	ctx.EventReceiver.OnList()
	ctx.StackDecoder(decodeByFirstChar)
}

func decodeListEnd(ctx *DecoderContext) {
	ctx.Stream.AdvanceByte() // Advance past ']'

	ctx.EventReceiver.OnEnd()
	ctx.UnstackDecoder()
}

func decodeMarkupBegin(ctx *DecoderContext) {
	ctx.Stream.AdvanceByte() // Advance past '<'

	ctx.EventReceiver.OnMarkup()
	ctx.StackDecoder(decodeMarkupName)
}

func decodeMarkupName(ctx *DecoderContext) {
	decodeByFirstChar(ctx)
	ctx.ChangeDecoder(decodeMapKey)
}

func decodeMarkupContentBegin(ctx *DecoderContext) {
	ctx.Stream.AdvanceByte() // Advance past ','

	ctx.EventReceiver.OnEnd()
	ctx.BeginMarkupContents()
}

func decodeMarkupContents(ctx *DecoderContext) {
	str, next := ctx.Stream.DecodeMarkupContent()
	if len(str) > 0 {
		ctx.EventReceiver.OnArray(events.ArrayTypeString, uint64(len(str)), str)
	}
	switch next {
	case nextIsCommentBegin:
		ctx.EventReceiver.OnComment()
		ctx.StackDecoder(decodeComment)
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
		ctx.EventReceiver.OnMarkup()
		ctx.StackDecoder(decodeMarkupBegin)
	case nextIsMarkupEnd:
		ctx.EventReceiver.OnEnd()
		ctx.UnstackDecoder()
	}
}

func decodeMarkupEnd(ctx *DecoderContext) {
	ctx.Stream.AdvanceByte() // Advance past '>'

	ctx.EventReceiver.OnEnd()
	ctx.EndMarkup()
}

func decodeConstant(ctx *DecoderContext) {
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

func decodeReference(ctx *DecoderContext) {
	ctx.Stream.AdvanceByte() // Advance past '$'

	ctx.EventReceiver.OnReference()
	panic("TODO: decodeReference")
}

func decodeMarker(ctx *DecoderContext) {
	ctx.Stream.AdvanceByte() // Advance past '&'

	ctx.EventReceiver.OnMarker()
	panic("TODO: decodeMarker")
}

func decodeMetadataBegin(ctx *DecoderContext) {
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
		ctx.Stream.Errorf("Expected Metadata separator (=) but got [%v]", ctx.Stream.DescribeCurrentChar())
	}
	ctx.Stream.AdvanceByte()
	decodeWhitespace(ctx)
	ctx.ChangeDecoder(decodeMetadataKey)
	decodeByFirstChar(ctx)
}

func decodeMetadataEnd(ctx *DecoderContext) {
	ctx.Stream.AdvanceByte() // Advance past ')'

	ctx.EventReceiver.OnEnd()
	ctx.ChangeDecoder(decodeMetadataCompletion)
}

func decodeMetadataCompletion(ctx *DecoderContext) {
	ctx.UnstackDecoder()
	decodeByFirstChar(ctx)
}

func decodeComment(ctx *DecoderContext) {

	ctx.EventReceiver.OnComment()
	panic("TODO: decodeComment")
}

func decodeSuffix(ctx *DecoderContext) {
	ctx.Stream.AdvanceByte() // Advance past ':'

	panic("TODO: decodeSuffix")
}
