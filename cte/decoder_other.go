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
	"math/big"

	"github.com/kstenerud/go-concise-encoding/internal/chars"

	"github.com/kstenerud/go-concise-encoding/internal/common"
)

type decodeInvalidChar struct{}

var global_decodeInvalidChar decodeInvalidChar

func (_this decodeInvalidChar) Run(ctx *DecoderContext) {
	ctx.Errorf("Unexpected [%v]", ctx.DescribeCurrentChar())
}

type decodeWhitespace struct{}

var global_decodeWhitespace decodeWhitespace

func (_this decodeWhitespace) Run(ctx *DecoderContext) {
	if ctx.Stream.PeekByteAllowEOF().HasProperty(chars.StructWS) {
		ctx.NotifyStructuralWS()
	}

	ctx.Stream.SkipWhitespace()
}

type decodeByFirstChar struct{}

var global_decodeByFirstChar decodeByFirstChar

func (_this decodeByFirstChar) Run(ctx *DecoderContext) {
	global_decodeWhitespace.Run(ctx)
	b := ctx.Stream.PeekByteAllowEOF()
	decoderOp := decoderOpsByFirstChar[b]
	decoderOp.Run(ctx)
}

type decodePostInvisible struct{}

var global_decodePostInvisible decodePostInvisible

func (_this decodePostInvisible) Run(ctx *DecoderContext) {
	ctx.UnstackDecoder()
	global_decodeByFirstChar.Run(ctx)
}

type decodeDocumentBegin struct{}

var global_decodeDocumentBegin decodeDocumentBegin

func (_this decodeDocumentBegin) Run(ctx *DecoderContext) {
	// Technically disallowed, but we'll support it anyway.
	global_decodeWhitespace.Run(ctx)

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

	global_decodeWhitespace.Run(ctx)
	ctx.AssertHasStructuralWS()

	ctx.EventReceiver.OnBeginDocument()
	ctx.EventReceiver.OnVersion(version)
	ctx.ChangeDecoder(global_decodeTopLevel)
}

type decodeTopLevel struct{}

var global_decodeTopLevel decodeTopLevel

func (_this decodeTopLevel) Run(ctx *DecoderContext) {
	ctx.AssertHasStructuralWS()
	ctx.ChangeDecoder(global_decodeEndDocument)
	global_decodeByFirstChar.Run(ctx)
}

type decodeEndDocument struct{}

var global_decodeEndDocument decodeEndDocument

func (_this decodeEndDocument) Run(ctx *DecoderContext) {
	ctx.EventReceiver.OnEndDocument()
	ctx.IsDocumentComplete = true
}

type decodeNumericPositive struct{}

var global_decodeNumericPositive decodeNumericPositive

func (_this decodeNumericPositive) Run(ctx *DecoderContext) {
	ctx.AssertHasStructuralWS()
	token := ctx.Stream.ReadToken()
	token.AssertNotEmpty(ctx.TextPos, "numeric")

	// 00000000-0000-0000-0000-000000000000
	if len(token) == 36 && token[8] == '-' {
		ctx.EventReceiver.OnUID(token.DecodeUID(ctx.TextPos))
		ctx.RequireStructuralWS()
		return
	}

	value, bigValue, digitCount, decodedCount := token.DecodeDecimalUint(ctx.TextPos)
	sign := 1

	// 123
	if token.IsAtEnd(decodedCount) {
		decodeTokenAsDecimalInt(ctx, token, decodedCount, value, bigValue, sign)
		ctx.RequireStructuralWS()
		return
	}

	switch token[decodedCount] {
	case '.':
		// 1.23
		decodeTokenAsDecimalFloat(ctx, token, decodedCount, digitCount, value, bigValue, sign)
	case '-':
		// 2000-01-01
		// TODO: Check for overflow
		decodeTokenAsDate(ctx, token, decodedCount, int(value))
	case ':':
		// 10:23:45
		// TODO: Check for overflow
		decodeTokenAsTime(ctx, token, decodedCount, int(value))
	default:
		token.UnexpectedChar(ctx.TextPos, decodedCount, "numeric")
	}
	ctx.RequireStructuralWS()
}

func reinterpretDecAsHex(v uint64) uint64 {
	var result uint64
	for position := 0; v != 0; position += 4 {
		result |= (v % 10) << position
		v /= 10
	}
	return result
}

type decodeUID struct{}

var global_decodeUID decodeUID

func (_this decodeUID) Run(ctx *DecoderContext) {
	ctx.AssertHasStructuralWS()
	token := ctx.Stream.ReadToken()
	ctx.EventReceiver.OnUID(token.DecodeUID(ctx.TextPos))
	ctx.RequireStructuralWS()
}

type advanceAndDecodeNumericNegative struct{}

var global_advanceAndDecodeNumericNegative advanceAndDecodeNumericNegative

func (_this advanceAndDecodeNumericNegative) decode0Based(ctx *DecoderContext, token Token) {
	sign := -1

	// 0
	if len(token) == 1 {
		ctx.EventReceiver.OnNegativeInt(0)
		return
	}

	switch token[1] {
	case 'b', 'B':
		// 0b1010
		decodeTokenAsBinaryInt(ctx, token, sign)
		return
	case 'o', 'O':
		// 0o1234
		decodeTokenAsOctalInt(ctx, token, sign)
		return
	case 'x', 'X':
		// 0x1234
		decodeTokenAsHexNumber(ctx, token, sign)
		return
	}

	value, bigValue, digitCount, decodedCount := token.DecodeDecimalUint(ctx.TextPos)

	// 0123
	if token.IsAtEnd(decodedCount) {
		decodeTokenAsDecimalInt(ctx, token, decodedCount, value, bigValue, sign)
		return
	}

	switch token[decodedCount] {
	case '.':
		// 0.123
		decodeTokenAsDecimalFloat(ctx, token, decodedCount, digitCount, value, bigValue, sign)
		return
	default:
		token.UnexpectedChar(ctx.TextPos, decodedCount, "negative 0-based numeric")
		return
	}
}

func (_this advanceAndDecodeNumericNegative) Run(ctx *DecoderContext) {
	ctx.AssertHasStructuralWS()
	ctx.Stream.AdvanceByte() // Advance past '-'
	sign := -1

	token := ctx.Stream.ReadToken()
	token.AssertNotEmpty(ctx.TextPos, "negative numeric")

	switch token[0] {
	case '0':
		_this.decode0Based(ctx, token)
		ctx.RequireStructuralWS()
		return
	case 'i', 'I':
		common.ASCIIBytesToLower(token)
		namedValue := string(token)
		if namedValue != "inf" {
			ctx.Errorf("Unknown named value: %v", namedValue)
		}
		ctx.EventReceiver.OnFloat(math.Inf(-1))
		ctx.RequireStructuralWS()
		return
	}

	value, bigValue, digitCount, decodedCount := token.DecodeDecimalUint(ctx.TextPos)

	// 123
	if token.IsAtEnd(decodedCount) {
		decodeTokenAsDecimalInt(ctx, token, decodedCount, value, bigValue, sign)
		ctx.RequireStructuralWS()
		return
	}

	switch token[decodedCount] {
	case '.':
		// 1.23
		decodeTokenAsDecimalFloat(ctx, token, decodedCount, digitCount, value, bigValue, sign)
	case '-':
		// 2000-01-01
		// TODO: Check for overflow
		decodeTokenAsDate(ctx, token, decodedCount, int(value)*sign)
	default:
		token.UnexpectedChar(ctx.TextPos, decodedCount, "numeric")
	}
	ctx.RequireStructuralWS()
}

type decode0Based struct{}

var global_decode0Based decode0Based

func (_this decode0Based) Run(ctx *DecoderContext) {
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
		decodeTokenAsBinaryInt(ctx, token, sign)
		ctx.RequireStructuralWS()
		return
	case 'o', 'O':
		// 0o1234
		decodeTokenAsOctalInt(ctx, token, sign)
		ctx.RequireStructuralWS()
		return
	case 'x', 'X':
		// 0x1234
		decodeTokenAsHexNumber(ctx, token, sign)
		ctx.RequireStructuralWS()
		return
	}

	value, bigValue, digitCount, decodedCount := token.DecodeDecimalUint(ctx.TextPos)

	// 0123
	if token.IsAtEnd(decodedCount) {
		decodeTokenAsDecimalInt(ctx, token, decodedCount, value, bigValue, sign)
		ctx.RequireStructuralWS()
		return
	}

	switch token[decodedCount] {
	case '.':
		// 0.123
		decodeTokenAsDecimalFloat(ctx, token, decodedCount, digitCount, value, bigValue, sign)
	case ':':
		// 01:23:45
		decodeTokenAsTime(ctx, token, decodedCount, int(value))
	default:
		token.UnexpectedChar(ctx.TextPos, decodedCount, "0-based numeric")
	}
	ctx.RequireStructuralWS()
}

type decodeFalseOrUID struct{}

var global_decodeFalseOrUID decodeFalseOrUID

func (_this decodeFalseOrUID) Run(ctx *DecoderContext) {
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

type decodeNamedValueI struct{}

var global_decodeNamedValueI decodeNamedValueI

func (_this decodeNamedValueI) Run(ctx *DecoderContext) {
	ctx.AssertHasStructuralWS()
	namedValue := ctx.Stream.ReadNamedValue()
	switch string(namedValue) {
	case "inf":
		ctx.EventReceiver.OnFloat(math.Inf(1))
	default:
		ctx.Errorf("%v: Unknown named value", string(namedValue))
	}
	ctx.RequireStructuralWS()
}

type decodeNamedValueN struct{}

var global_decodeNamedValueN decodeNamedValueN

func (_this decodeNamedValueN) Run(ctx *DecoderContext) {
	ctx.AssertHasStructuralWS()
	namedValue := ctx.Stream.ReadNamedValue()
	switch string(namedValue) {
	case "na":
		if ctx.Stream.ReadByteNoEOF() != ':' {
			ctx.Stream.UnreadByte()
			ctx.UnexpectedChar("NA")
		}
		ctx.EventReceiver.OnNA()
		global_decodeByFirstChar.Run(ctx)
	case "nan":
		ctx.EventReceiver.OnNan(false)
	case "nil":
		ctx.EventReceiver.OnNil()
	default:
		ctx.Errorf("%v: Unknown named value", string(namedValue))
	}
	ctx.RequireStructuralWS()
}

type decodeNamedValueS struct{}

var global_decodeNamedValueS decodeNamedValueS

func (_this decodeNamedValueS) Run(ctx *DecoderContext) {
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

type decodeNamedValueT struct{}

var global_decodeNamedValueT decodeNamedValueT

func (_this decodeNamedValueT) Run(ctx *DecoderContext) {
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

type advanceAndDecodeConstant struct{}

var global_advanceAndDecodeConstant advanceAndDecodeConstant

func (_this advanceAndDecodeConstant) Run(ctx *DecoderContext) {
	ctx.AssertHasStructuralWS()
	ctx.Stream.AdvanceByte() // Advance past '#'

	ctx.EventReceiver.OnConstant(ctx.Stream.ReadIdentifier())
	ctx.RequireStructuralWS()
}

type advanceAndDecodeMarker struct{}

var global_advanceAndDecodeMarker advanceAndDecodeMarker

func (_this advanceAndDecodeMarker) Run(ctx *DecoderContext) {
	ctx.AssertHasStructuralWS()
	ctx.Stream.AdvanceByte() // Advance past '&'

	ctx.EventReceiver.OnMarker(ctx.Stream.ReadMarkerIdentifier())
	if ctx.Stream.PeekByteNoEOF() != ':' {
		ctx.Errorf("Missing colon between marker ID and marked value")
	}
	ctx.Stream.AdvanceByte()
	global_decodeByFirstChar.Run(ctx)
}

type advanceAndDecodeReference struct{}

var global_advanceAndDecodeReference advanceAndDecodeReference

func (_this advanceAndDecodeReference) Run(ctx *DecoderContext) {
	ctx.AssertHasStructuralWS()
	ctx.Stream.AdvanceByte() // Advance past '$'

	if ctx.Stream.PeekByteNoEOF() == '@' {
		ctx.EventReceiver.OnRIDReference()
		global_advanceAndDecodeResourceID.Run(ctx)
		return
	}

	ctx.EventReceiver.OnReference(ctx.Stream.ReadMarkerIdentifier())
}

// ===========================================================================

func decodeTokenAsBinaryInt(ctx *DecoderContext, token Token, sign int) {
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

func decodeTokenAsOctalInt(ctx *DecoderContext, token Token, sign int) {
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

func decodeTokenAsHexNumber(ctx *DecoderContext, token Token, sign int) {
	token = token[2:]
	value, bigValue, digitCount, decodedCount := token.DecodeHexUint(ctx.TextPos)
	token = token[decodedCount:]

	if token.IsAtEnd(0) {
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

	if token[0] != '.' {
		token.UnexpectedChar(ctx.TextPos, 0, "hexadecimal int")
	}

	var fvalue float64
	var bigFValue *big.Float
	fvalue, bigFValue, decodedCount = token.CompleteHexFloat(ctx.TextPos, int64(sign), value, bigValue, digitCount)
	token = token[decodedCount:]
	if bigFValue != nil {
		ctx.EventReceiver.OnBigFloat(bigFValue)
	} else {
		ctx.EventReceiver.OnFloat(fvalue)
	}

	token.AssertAtEnd(ctx.TextPos, "hexadecimal number")
}

func decodeTokenAsDecimalInt(ctx *DecoderContext, token Token, decodedCount int, value uint64, bigValue *big.Int, sign int) {
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

func decodeTokenAsDecimalFloat(ctx *DecoderContext, token Token, decodedCount int, digitCount int, value uint64, bigValue *big.Int, sign int) {
	token = token[decodedCount:]
	fvalue, bigFValue, decodedCount := token.CompleteDecimalFloat(ctx.TextPos, int64(sign), value, bigValue, digitCount)
	token = token[decodedCount:]
	if bigFValue != nil {
		ctx.EventReceiver.OnBigDecimalFloat(bigFValue)
	} else {
		ctx.EventReceiver.OnDecimalFloat(fvalue)
	}
	token.AssertAtEnd(ctx.TextPos, "decimal float")
}

func decodeTokenAsDate(ctx *DecoderContext, token Token, decodedCount int, year int) {
	token = token[decodedCount:]
	tvalue, decodedCount := token.CompleteDate(ctx.TextPos, year)
	token = token[decodedCount:]
	ctx.EventReceiver.OnCompactTime(tvalue)
	token.AssertAtEnd(ctx.TextPos, "date")
}

func decodeTokenAsTime(ctx *DecoderContext, token Token, decodedCount int, hour int) {
	if decodedCount < 1 || decodedCount > 2 {
		token.UnexpectedChar(ctx.TextPos, decodedCount, "time")
	}
	token = token[decodedCount:]
	tvalue, decodedCount := token.CompleteTime(ctx.TextPos, 0, 0, 0, hour)
	token = token[decodedCount:]
	ctx.EventReceiver.OnCompactTime(tvalue)
	token.AssertAtEnd(ctx.TextPos, "time")
}
