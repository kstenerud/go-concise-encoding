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
	"math/big"

	compact_float "github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-concise-encoding/internal/chars"
	"github.com/kstenerud/go-concise-encoding/internal/common"
)

func decodeInvalidChar(ctx *DecoderContext) {
	ctx.Errorf("Unexpected [%v]", ctx.DescribeCurrentChar())
}

func decodeWhitespace(ctx *DecoderContext) {
	if ctx.Stream.PeekByteAllowEOF().HasProperty(chars.StructWS) {
		ctx.NotifyStructuralWS()
	}

	ctx.Stream.SkipWhitespace()
}

func decodeByFirstChar(ctx *DecoderContext) {
	decodeWhitespace(ctx)
	b := ctx.Stream.PeekByteAllowEOF()
	decoderOp := decoderOpsByFirstChar[b]
	decoderOp(ctx)
}

func decodePostInvisible(ctx *DecoderContext) {
	ctx.UnstackDecoder()
	decodeByFirstChar(ctx)
}

func decodeDocumentBegin(ctx *DecoderContext) {
	// Technically disallowed, but we'll support it anyway.
	decodeWhitespace(ctx)

	if b := ctx.Stream.ReadByteNoEOF(); b != 'c' && b != 'C' {
		ctx.Errorf(`Expected document to begin with "c" but got [%v]`, ctx.DescribeCurrentChar())
	}

	token := ctx.Stream.ReadToken()
	version, digitCount, _ := token.DecodeSmallDecimalUint(ctx.TextPos)
	if digitCount == 0 {
		ctx.UnexpectedChar("version number")
	}

	// TODO: Remove this when releasing V1
	if version == 1 {
		version = 0
	}

	decodeWhitespace(ctx)
	ctx.AssertHasStructuralWS()

	ctx.EventReceiver.OnBeginDocument()
	ctx.EventReceiver.OnVersion(version)
	ctx.ChangeDecoder(decodeTopLevel)
}

func decodeTopLevel(ctx *DecoderContext) {
	ctx.AssertHasStructuralWS()
	ctx.ChangeDecoder(decodeEndDocument)
	decodeByFirstChar(ctx)
}

func decodeEndDocument(ctx *DecoderContext) {
	ctx.EventReceiver.OnEndDocument()
	ctx.IsDocumentComplete = true
}

func decodeNumericPositive(ctx *DecoderContext) {
	ctx.AssertHasStructuralWS()
	token := ctx.Stream.ReadToken()
	token.AssertNotEmpty(ctx.TextPos, "numeric")

	// 00000000-0000-0000-0000-000000000000
	if len(token) == 36 && token[8] == '-' {
		ctx.EventReceiver.OnUID(token.DecodeUID(ctx.TextPos))
		ctx.RequireStructuralWS()
		return
	}

	value, bigValue, _, decodedCount := token.DecodeDecimalUint(ctx.TextPos)
	sign := 1

	// 123
	if token.IsAtEnd(decodedCount) {
		continueDecodingAsDecimalInt(ctx, token, decodedCount, value, bigValue, sign)
		ctx.RequireStructuralWS()
		return
	}

	switch token[decodedCount] {
	case '.', ',':
		// 1.23
		continueDecodingAsDecimalFloat(ctx, token, decodedCount, value, bigValue, sign)
	case 'e', 'E':
		// 1e+10
		continueDecodingAsDecimalExponent(ctx, token, decodedCount, value, bigValue, sign)
	case '-':
		// 2000-01-01
		// TODO: Check for overflow
		continueDecodingAsDate(ctx, token, decodedCount, int(value))
	case ':':
		// 10:23:45
		// TODO: Check for overflow
		continueDecodingAsTime(ctx, token, decodedCount, int(value))
	default:
		token.UnexpectedChar(ctx.TextPos, decodedCount, "numeric")
	}
	ctx.RequireStructuralWS()
}

func decodeUID(ctx *DecoderContext) {
	ctx.AssertHasStructuralWS()
	token := ctx.Stream.ReadToken()
	ctx.EventReceiver.OnUID(token.DecodeUID(ctx.TextPos))
	ctx.RequireStructuralWS()
}

func decodeTokenAsNegative0Based(ctx *DecoderContext, token Token) {
	sign := -1

	// 0
	if len(token) == 1 {
		ctx.EventReceiver.OnNegativeInt(0)
		return
	}

	switch token[1] {
	case 'b', 'B':
		// 0b1010
		continueDecodingAsBinaryInt(ctx, token, sign)
		return
	case 'o', 'O':
		// 0o1234
		continueDecodingAsOctalInt(ctx, token, sign)
		return
	case 'x', 'X':
		// 0x1234
		continueDecodingAsHexNumber(ctx, token, sign)
		return
	}

	value, bigValue, _, decodedCount := token.DecodeDecimalUint(ctx.TextPos)

	// 0123
	if token.IsAtEnd(decodedCount) {
		continueDecodingAsDecimalInt(ctx, token, decodedCount, value, bigValue, sign)
		return
	}

	switch token[decodedCount] {
	case '.', ',':
		// 0.123
		continueDecodingAsDecimalFloat(ctx, token, decodedCount, value, bigValue, sign)
		return
	case 'e', 'E':
		// 1e+10
		continueDecodingAsDecimalExponent(ctx, token, decodedCount, value, bigValue, sign)
	case '-':
		// -2000-01-01
		// TODO: Check for overflow
		continueDecodingAsDate(ctx, token, decodedCount, int(value)*sign)
	default:
		token.UnexpectedChar(ctx.TextPos, decodedCount, "negative 0-based numeric")
		return
	}
}

func advanceAndDecodeNumericNegative(ctx *DecoderContext) {
	ctx.AssertHasStructuralWS()
	ctx.Stream.AdvanceByte() // Advance past '-'
	sign := -1

	token := ctx.Stream.ReadToken()
	token.AssertNotEmpty(ctx.TextPos, "negative numeric")

	switch token[0] {
	case '0':
		decodeTokenAsNegative0Based(ctx, token)
		ctx.RequireStructuralWS()
		return
	case 'i', 'I':
		common.ASCIIBytesToLower(token)
		namedValue := string(token)
		if namedValue != "inf" {
			ctx.Errorf("Unknown named value: %v", namedValue)
		}
		ctx.EventReceiver.OnDecimalFloat(compact_float.NegativeInfinity())
		ctx.RequireStructuralWS()
		return
	}

	value, bigValue, _, decodedCount := token.DecodeDecimalUint(ctx.TextPos)

	// 123
	if token.IsAtEnd(decodedCount) {
		continueDecodingAsDecimalInt(ctx, token, decodedCount, value, bigValue, sign)
		ctx.RequireStructuralWS()
		return
	}

	switch token[decodedCount] {
	case '.', ',':
		// 1.23
		continueDecodingAsDecimalFloat(ctx, token, decodedCount, value, bigValue, sign)
	case 'e', 'E':
		// 1e+10
		continueDecodingAsDecimalExponent(ctx, token, decodedCount, value, bigValue, sign)
	case '-':
		// 2000-01-01
		// TODO: Check for overflow
		continueDecodingAsDate(ctx, token, decodedCount, int(value)*sign)
	default:
		token.UnexpectedChar(ctx.TextPos, decodedCount, "numeric")
	}
	ctx.RequireStructuralWS()
}

func decode0Based(ctx *DecoderContext) {
	ctx.AssertHasStructuralWS()
	// Assumption: First character is 0

	token := ctx.Stream.ReadToken()
	token.AssertNotEmpty(ctx.TextPos, "0-based numeric")

	// 0
	if len(token) == 1 {
		ctx.EventReceiver.OnPositiveInt(0)
		ctx.RequireStructuralWS()
		return
	}

	// 00000000-0000-0000-0000-000000000000
	if len(token) == 36 && token[8] == '-' {
		ctx.EventReceiver.OnUID(token.DecodeUID(ctx.TextPos))
		ctx.RequireStructuralWS()
		return
	}

	sign := 1

	switch token[1] {
	case 'b', 'B':
		// 0b1010
		continueDecodingAsBinaryInt(ctx, token, sign)
		ctx.RequireStructuralWS()
		return
	case 'o', 'O':
		// 0o1234
		continueDecodingAsOctalInt(ctx, token, sign)
		ctx.RequireStructuralWS()
		return
	case 'x', 'X':
		// 0x1234
		continueDecodingAsHexNumber(ctx, token, sign)
		ctx.RequireStructuralWS()
		return
	}

	value, bigValue, _, decodedCount := token.DecodeDecimalUint(ctx.TextPos)

	// 0123
	if token.IsAtEnd(decodedCount) {
		continueDecodingAsDecimalInt(ctx, token, decodedCount, value, bigValue, sign)
		ctx.RequireStructuralWS()
		return
	}

	switch token[decodedCount] {
	case '.', ',':
		// 0.123
		continueDecodingAsDecimalFloat(ctx, token, decodedCount, value, bigValue, sign)
	case 'e', 'E':
		// 1e+10
		continueDecodingAsDecimalExponent(ctx, token, decodedCount, value, bigValue, sign)
	case '-':
		// 2000-01-01
		// TODO: Check for overflow
		continueDecodingAsDate(ctx, token, decodedCount, int(value)*sign)
	case ':':
		// 01:23:45
		continueDecodingAsTime(ctx, token, decodedCount, int(value))
	default:
		token.UnexpectedChar(ctx.TextPos, decodedCount, "0-based numeric")
	}
	ctx.RequireStructuralWS()
}

func decodeFalseOrUID(ctx *DecoderContext) {
	ctx.AssertHasStructuralWS()
	token := ctx.Stream.ReadToken()

	// 00000000-0000-0000-0000-000000000000
	if len(token) == 36 && token[8] == '-' {
		ctx.EventReceiver.OnUID(token.DecodeUID(ctx.TextPos))
		ctx.RequireStructuralWS()
		return
	}

	common.ASCIIBytesToLower(token)
	named := string(token)
	switch named {
	case "false":
		ctx.EventReceiver.OnFalse()
	default:
		ctx.Errorf("%v: Unknown named value", named)
	}
	ctx.RequireStructuralWS()
}

func decodeNamedValueI(ctx *DecoderContext) {
	ctx.AssertHasStructuralWS()
	namedValue := ctx.Stream.ReadNamedValue()
	switch string(namedValue) {
	case "inf":
		ctx.EventReceiver.OnDecimalFloat(compact_float.Infinity())
	default:
		ctx.Errorf("%v: Unknown named value", string(namedValue))
	}
	ctx.RequireStructuralWS()
}

func decodeNamedValueN(ctx *DecoderContext) {
	ctx.AssertHasStructuralWS()
	namedValue := ctx.Stream.ReadNamedValue()
	switch string(namedValue) {
	case "nan":
		ctx.EventReceiver.OnNan(false)
	case "null":
		ctx.EventReceiver.OnNull()
	default:
		ctx.Errorf("%v: Unknown named value", string(namedValue))
	}
	ctx.RequireStructuralWS()
}

func decodeNamedValueS(ctx *DecoderContext) {
	ctx.AssertHasStructuralWS()
	namedValue := ctx.Stream.ReadNamedValue()
	switch string(namedValue) {
	case "snan":
		ctx.EventReceiver.OnNan(true)
	default:
		ctx.Errorf("%v: Unknown named value", string(namedValue))
	}
	ctx.RequireStructuralWS()
}

func decodeNamedValueT(ctx *DecoderContext) {
	ctx.AssertHasStructuralWS()
	namedValue := ctx.Stream.ReadNamedValue()
	switch string(namedValue) {
	case "true":
		ctx.EventReceiver.OnTrue()
	default:
		ctx.Errorf("%v: Unknown named value", string(namedValue))
	}
	ctx.RequireStructuralWS()
}

func advanceAndDecodeMarker(ctx *DecoderContext) {
	ctx.AssertHasStructuralWS()
	ctx.Stream.AdvanceByte() // Advance past '&'

	ctx.EventReceiver.OnMarker(ctx.Stream.ReadMarkerIdentifier())
	if ctx.Stream.PeekByteNoEOF() != ':' {
		ctx.Errorf("Missing colon between marker ID and marked value")
	}
	ctx.Stream.AdvanceByte()
	if ctx.Stream.NextByteHasProperty(chars.StructWS) {
		ctx.Errorf("whitespace not allowed between marker ID and marked object")
	}
	decodeByFirstChar(ctx)
}

func advanceAndDecodeReference(ctx *DecoderContext) {
	ctx.AssertHasStructuralWS()
	ctx.Stream.AdvanceByte() // Advance past '$'

	if ctx.Stream.PeekByteNoEOF() == '"' {
		ctx.Stream.AdvanceByte() // Advance past '"'
		decodeRemoteReference(ctx)
		return
	}

	ctx.EventReceiver.OnReference(ctx.Stream.ReadMarkerIdentifier())
}

func advanceAndDecodeEdgeOrResourceID(ctx *DecoderContext) {
	ctx.AssertHasStructuralWS()
	ctx.Stream.AdvanceByte() // Advance past '@'
	switch ctx.Stream.ReadByteNoEOF() {
	case '"':
		decodeResourceID(ctx)
	case '(':
		decodeEdgeBegin(ctx)
	default:
		ctx.Stream.UnreadLastByte()
		ctx.UnexpectedChar("edge or resource ID")
	}
}

// ===========================================================================

func continueDecodingAsBinaryInt(ctx *DecoderContext, token Token, sign int) {
	// Assumption: First two chars are 0b

	token = token[2:]
	value, bigValue, _, decodedCount := token.DecodeBinaryUint(ctx.TextPos)
	token = token[decodedCount:]
	if bigValue != nil {
		if sign < 0 {
			bigValue = bigValue.Neg(bigValue)
		}
		ctx.EventReceiver.OnBigInt(bigValue)
	} else {
		if sign >= 0 {
			ctx.EventReceiver.OnPositiveInt(value)
		} else {
			ctx.EventReceiver.OnNegativeInt(value)
		}
	}
	token.AssertAtEnd(ctx.TextPos, "binary integer")
}

func continueDecodingAsOctalInt(ctx *DecoderContext, token Token, sign int) {
	// Assumption: First two chars are 0o

	token = token[2:]
	value, bigValue, _, decodedCount := token.DecodeOctalUint(ctx.TextPos)
	token = token[decodedCount:]
	if bigValue != nil {
		if sign < 0 {
			bigValue = bigValue.Neg(bigValue)
		}
		ctx.EventReceiver.OnBigInt(bigValue)
	} else {
		if sign >= 0 {
			ctx.EventReceiver.OnPositiveInt(value)
		} else {
			ctx.EventReceiver.OnNegativeInt(value)
		}
	}
	token.AssertAtEnd(ctx.TextPos, "octal integer")
}

func continueDecodingAsHexNumber(ctx *DecoderContext, token Token, sign int) {
	// Assumption: First two chars are 0x

	pos := 2
	value, bigValue, digitCount, decodedCount := token[pos:].DecodeHexUint(ctx.TextPos)
	if decodedCount == 0 {
		token.UnexpectedChar(ctx.TextPos, pos, "hexadecimal")
	}
	pos += decodedCount

	if token.IsAtEnd(pos) {
		if bigValue != nil {
			if sign < 0 {
				bigValue = bigValue.Neg(bigValue)
			}
			ctx.EventReceiver.OnBigInt(bigValue)
		} else {
			if sign >= 0 {
				ctx.EventReceiver.OnPositiveInt(value)
			} else {
				ctx.EventReceiver.OnNegativeInt(value)
			}
		}
		return
	}

	switch token[pos] {
	case '.', ',':
		// 0x3.1f
		fvalue, bigFValue, decodedCount := token[pos:].CompleteHexFloat(ctx.TextPos, int64(sign), value, bigValue, digitCount)
		token = token[decodedCount+pos:]
		if bigFValue != nil {
			ctx.EventReceiver.OnBigFloat(bigFValue)
		} else {
			ctx.EventReceiver.OnFloat(fvalue)
		}
	case 'p', 'P':
		// 1p+10
		fvalue, bigFValue, decodedCount := token[pos:].CompleteHexExponent(ctx.TextPos, int64(sign), value, bigValue, digitCount)
		token = token[decodedCount+pos:]
		if bigFValue != nil {
			ctx.EventReceiver.OnBigFloat(bigFValue)
		} else {
			ctx.EventReceiver.OnFloat(fvalue)
		}
	default:
		token.UnexpectedChar(ctx.TextPos, pos, "numeric")
	}

	token.AssertAtEnd(ctx.TextPos, "hexadecimal number")
}

func continueDecodingAsDecimalInt(ctx *DecoderContext, token Token, decodedCount int, value uint64, bigValue *big.Int, sign int) {
	token = token[decodedCount:]
	if bigValue != nil {
		if sign < 0 {
			bigValue = bigValue.Neg(bigValue)
		}
		ctx.EventReceiver.OnBigInt(bigValue)
	} else {
		if sign >= 0 {
			ctx.EventReceiver.OnPositiveInt(value)
		} else {
			ctx.EventReceiver.OnNegativeInt(value)
		}
	}
	token.AssertAtEnd(ctx.TextPos, "decimal integer")
}

func continueDecodingAsDecimalFloat(ctx *DecoderContext, token Token, decodedCount int, value uint64, bigValue *big.Int, sign int) {
	token = token[decodedCount:]
	fvalue, bigFValue, decodedCount := token.CompleteDecimalFloat(ctx.TextPos, int64(sign), value, bigValue)
	token = token[decodedCount:]
	if bigFValue != nil {
		ctx.EventReceiver.OnBigDecimalFloat(bigFValue)
	} else {
		ctx.EventReceiver.OnDecimalFloat(fvalue)
	}
	token.AssertAtEnd(ctx.TextPos, "decimal float")
}

func continueDecodingAsDecimalExponent(ctx *DecoderContext, token Token, decodedCount int, value uint64, bigValue *big.Int, sign int) {
	token = token[decodedCount:]
	fvalue, bigFValue, decodedCount := token.CompleteDecimalExponent(ctx.TextPos, int64(sign), value, bigValue, 0)
	token = token[decodedCount:]
	if bigFValue != nil {
		ctx.EventReceiver.OnBigDecimalFloat(bigFValue)
	} else {
		ctx.EventReceiver.OnDecimalFloat(fvalue)
	}
	token.AssertAtEnd(ctx.TextPos, "decimal float")
}

func continueDecodingAsDate(ctx *DecoderContext, token Token, decodedCount int, year int) {
	token = token[decodedCount:]
	tvalue, decodedCount := token.CompleteDate(ctx.TextPos, year)
	token = token[decodedCount:]
	ctx.EventReceiver.OnCompactTime(tvalue)
	token.AssertAtEnd(ctx.TextPos, "date")
}

func continueDecodingAsTime(ctx *DecoderContext, token Token, decodedCount int, hour int) {
	if decodedCount < 1 || decodedCount > 2 {
		token.UnexpectedChar(ctx.TextPos, decodedCount, "time")
	}
	token = token[decodedCount:]
	tvalue, decodedCount := token.CompleteTime(ctx.TextPos, 0, 0, 0, hour)
	token = token[decodedCount:]
	ctx.EventReceiver.OnCompactTime(tvalue)
	token.AssertAtEnd(ctx.TextPos, "time")
}
