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

	version, bigVersion, digitCount := ctx.Stream.ReadDecimalUint(0, nil)
	if digitCount == 0 {
		ctx.UnexpectedChar("version number")
	}
	if bigVersion != nil {
		ctx.Errorf("Version too big")
	}
	// TODO: Remove this when releasing V1
	if version == 1 {
		version = 0
	}

	b := ctx.Stream.PeekByteNoEOF()
	if !chars.ByteHasProperty(b, chars.StructWS) {
		ctx.UnexpectedChar("whitespace after version")
	}
	global_decodeWhitespace.Run(ctx)

	ctx.EventReceiver.OnBeginDocument()
	ctx.EventReceiver.OnVersion(version)
	ctx.ChangeDecoder(global_decodeTopLevel)
}

type decodeTopLevel struct{}

var global_decodeTopLevel decodeTopLevel

func (_this decodeTopLevel) Run(ctx *DecoderContext) {
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
	coefficient, bigCoefficient, digitCount := ctx.Stream.ReadDecimalUint(0, nil)
	b := ctx.Stream.ReadByteAllowEOF()
	switch {
	case b == '-':
		// TODO: Could be UUID  (followed by 0-9a-fA-F *4)
		v := ctx.Stream.ReadDate(int64(coefficient))
		ctx.Stream.AssertAtObjectEnd("date")
		ctx.EventReceiver.OnCompactTime(v)
		return
	case b == ':':
		v := ctx.Stream.ReadTime(int(coefficient))
		ctx.Stream.AssertAtObjectEnd("time")
		ctx.EventReceiver.OnCompactTime(v)
		return
	case b == '.':
		value, bigValue, _ := ctx.Stream.ReadDecimalFloat(1, coefficient, bigCoefficient, digitCount)
		ctx.Stream.AssertAtObjectEnd("float")
		if bigValue != nil {
			ctx.EventReceiver.OnBigDecimalFloat(bigValue)
		} else {
			ctx.EventReceiver.OnDecimalFloat(value)
		}
		return
	case b.HasProperty(chars.LowerAF | chars.UpperAF):
		ctx.Stream.UnreadByte()
		ctx.EventReceiver.OnUUID(ctx.Stream.ReadUUIDWithDecimalDecoded(coefficient, digitCount))
		return
	default:
		if b.HasProperty(chars.ObjectEnd) {
			ctx.Stream.UnreadByte()
			if bigCoefficient != nil {
				ctx.EventReceiver.OnBigInt(bigCoefficient)
			} else {
				ctx.EventReceiver.OnPositiveInt(coefficient)
			}
			return
		}
	}
	ctx.UnexpectedChar("numeric")
}

func reinterpretDecAsHex(v uint64) uint64 {
	var result uint64
	for position := 0; v != 0; position += 4 {
		result |= (v % 10) << position
		v /= 10
	}
	return result
}

type decodeUUID struct{}

var global_decodeUUID decodeUUID

func (_this decodeUUID) Run(ctx *DecoderContext) {
	ctx.EventReceiver.OnUUID(ctx.Stream.ReadUUIDWithDecimalDecoded(0, 0))
}

type advanceAndDecodeNumericNegative struct{}

var global_advanceAndDecodeNumericNegative advanceAndDecodeNumericNegative

func (_this advanceAndDecodeNumericNegative) Run(ctx *DecoderContext) {
	ctx.Stream.AdvanceByte() // Advance past '-'

	switch ctx.Stream.PeekByteNoEOF() {
	case '0':
		ctx.Stream.AdvanceByte() // Advance past '0'
		global_decodeOtherBaseNegative.Run(ctx)
		return
	case 'i':
		namedValue := string(ctx.Stream.ReadNamedValue())
		if namedValue != "inf" {
			ctx.Errorf("Unknown named value: %v", namedValue)
		}
		ctx.EventReceiver.OnFloat(math.Inf(-1))
		return
	}

	coefficient, bigCoefficient, digitCount := ctx.Stream.ReadDecimalUint(0, nil)
	b := ctx.Stream.ReadByteAllowEOF()
	switch b {
	case '-':
		v := ctx.Stream.ReadDate(-int64(coefficient))
		ctx.Stream.AssertAtObjectEnd("time")
		ctx.EventReceiver.OnCompactTime(v)
		return
	case '.':
		value, bigValue, _ := ctx.Stream.ReadDecimalFloat(-1, coefficient, bigCoefficient, digitCount)
		ctx.Stream.AssertAtObjectEnd("float")
		if bigValue != nil {
			ctx.EventReceiver.OnBigDecimalFloat(bigValue)
		} else {
			ctx.EventReceiver.OnDecimalFloat(value)
		}
		return
	default:
		if b.HasProperty(chars.ObjectEnd) {
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
	ctx.UnexpectedChar("numeric")
}

type decode0Based struct{}

var global_decode0Based decode0Based

func (_this decode0Based) Run(ctx *DecoderContext) {
	// Assumption: First character is 0

	token := ctx.Stream.ReadToken()
	token.AssertNotEmpty(ctx.TextPos, "0-based numeric")

	// 0
	if len(token) == 1 {
		ctx.EventReceiver.OnPositiveInt(0)
		return
	}

	// 00000000-0000-0000-0000-000000000000
	if len(token) == 36 && token[8] == '-' {
		decodeTokenAsUUID(ctx, token)
		return
	}

	switch token[1] {
	case 'b':
		// 0b1010
		decodeTokenAsBinaryInt(ctx, token)
		return
	case 'o':
		// 0o1234
		decodeTokenAsOctalInt(ctx, token)
		return
	case 'x':
		// 0x1234
		decodeTokenAsHexNumber(ctx, token)
		return
	}

	value, bigValue, digitCount, decodedCount := token.DecodeDecimalUint(ctx.TextPos)

	// 0123
	if token.IsAtEnd(decodedCount) {
		decodeTokenAsDecimalInt(ctx, token, decodedCount, value, bigValue)
		return
	}

	switch token[decodedCount] {
	case '.':
		// 0.123
		decodeTokenAsDecimalFloat(ctx, token, decodedCount, digitCount, value, bigValue)
		return
	case ':':
		// 01:23:45
		decodeTokenAsTime(ctx, token, decodedCount, int(value))
		return
	default:
		token.UnexpectedChar(ctx.TextPos, decodedCount, "0-based numeric")
		return
	}
}

type advanceAndDecodeOtherBaseNegative struct{}

var global_advanceAndDecodeOtherBaseNegative advanceAndDecodeOtherBaseNegative

func (_this advanceAndDecodeOtherBaseNegative) Run(ctx *DecoderContext) {
	ctx.Stream.AdvanceByte() // Advance past '0'
	global_decodeOtherBaseNegative.Run(ctx)
}

type decodeOtherBaseNegative struct{}

var global_decodeOtherBaseNegative decodeOtherBaseNegative

func (_this decodeOtherBaseNegative) Run(ctx *DecoderContext) {

	b := ctx.Stream.PeekByteAllowEOF()
	if b.HasProperty(chars.ObjectEnd) {
		// -0 has no decimal point (thus type int), so report it as positive 0.
		ctx.EventReceiver.OnPositiveInt(0)
		return
	}

	ctx.Stream.AdvanceByte() // Advance past the value now stored in b

	switch b {
	case 'b':
		v, bigV, _ := ctx.Stream.ReadBinaryUint()
		ctx.Stream.AssertAtObjectEnd("binary integer")
		if bigV != nil {
			bigV = bigV.Neg(bigV)
			ctx.EventReceiver.OnBigInt(bigV)
		} else {
			ctx.EventReceiver.OnNegativeInt(v)
		}
	case 'o':
		v, bigV, _ := ctx.Stream.ReadOctalUint()
		ctx.Stream.AssertAtObjectEnd("octal integer")
		if bigV != nil {
			bigV = bigV.Neg(bigV)
			ctx.EventReceiver.OnBigInt(bigV)
		} else {
			ctx.EventReceiver.OnNegativeInt(v)
		}
	case 'x':
		v, bigV, digitCount := ctx.Stream.ReadHexUint(0, nil)
		if ctx.Stream.PeekByteAllowEOF() == '.' {
			ctx.Stream.AdvanceByte()
			fv, bigFV, _ := ctx.Stream.ReadHexFloat(-1, v, bigV, digitCount)
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
		value, bigValue, _ := ctx.Stream.ReadDecimalFloat(-1, 0, nil, 0)
		ctx.Stream.AssertAtObjectEnd("float")
		if bigValue != nil {
			ctx.EventReceiver.OnBigDecimalFloat(bigValue)
		} else {
			ctx.EventReceiver.OnDecimalFloat(value)
		}
	default:
		ctx.Stream.UnreadByte()
		ctx.UnexpectedChar("numeric base")
	}
}

type decodeFalseOrUUID struct{}

var global_decodeFalseOrUUID decodeFalseOrUUID

func (_this decodeFalseOrUUID) Run(ctx *DecoderContext) {
	token := ctx.Stream.ReadToken()

	// 00000000-0000-0000-0000-000000000000
	if len(token) == 36 && token[8] == '-' {
		decodeTokenAsUUID(ctx, token)
		return
	}

	named := string(token)
	switch named {
	case "false":
		ctx.EventReceiver.OnFalse()
	default:
		ctx.Errorf("%v: Unknown named value", named)
	}
}

type decodeNamedValueI struct{}

var global_decodeNamedValueI decodeNamedValueI

func (_this decodeNamedValueI) Run(ctx *DecoderContext) {
	namedValue := ctx.Stream.ReadNamedValue()
	switch string(namedValue) {
	case "inf":
		ctx.EventReceiver.OnFloat(math.Inf(1))
	default:
		ctx.Errorf("%v: Unknown named value", string(namedValue))
	}
}

type decodeNamedValueN struct{}

var global_decodeNamedValueN decodeNamedValueN

func (_this decodeNamedValueN) Run(ctx *DecoderContext) {
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
}

type decodeNamedValueS struct{}

var global_decodeNamedValueS decodeNamedValueS

func (_this decodeNamedValueS) Run(ctx *DecoderContext) {
	namedValue := ctx.Stream.ReadNamedValue()
	switch string(namedValue) {
	case "snan":
		ctx.EventReceiver.OnNan(true)
	default:
		ctx.Errorf("%v: Unknown named value", string(namedValue))
	}
}

type decodeNamedValueT struct{}

var global_decodeNamedValueT decodeNamedValueT

func (_this decodeNamedValueT) Run(ctx *DecoderContext) {
	namedValue := ctx.Stream.ReadNamedValue()
	switch string(namedValue) {
	case "true":
		ctx.EventReceiver.OnTrue()
	default:
		ctx.Errorf("%v: Unknown named value", string(namedValue))
	}
}

type advanceAndDecodeConstant struct{}

var global_advanceAndDecodeConstant advanceAndDecodeConstant

func (_this advanceAndDecodeConstant) Run(ctx *DecoderContext) {
	ctx.Stream.AdvanceByte() // Advance past '#'

	ctx.EventReceiver.OnConstant(ctx.Stream.ReadIdentifier())
}

type advanceAndDecodeMarker struct{}

var global_advanceAndDecodeMarker advanceAndDecodeMarker

func (_this advanceAndDecodeMarker) Run(ctx *DecoderContext) {
	ctx.Stream.AdvanceByte() // Advance past '&'

	ctx.EventReceiver.OnMarker(ctx.Stream.ReadIdentifier())
	if ctx.Stream.PeekByteNoEOF() != ':' {
		ctx.Errorf("Missing colon between marker ID and marked value")
	}
	ctx.Stream.AdvanceByte()
	global_decodeByFirstChar.Run(ctx)
}

type advanceAndDecodeReference struct{}

var global_advanceAndDecodeReference advanceAndDecodeReference

func (_this advanceAndDecodeReference) Run(ctx *DecoderContext) {
	ctx.Stream.AdvanceByte() // Advance past '$'

	if ctx.Stream.PeekByteNoEOF() == '@' {
		ctx.EventReceiver.OnRIDReference()
		global_advanceAndDecodeResourceID.Run(ctx)
		return
	}

	ctx.EventReceiver.OnReference(ctx.Stream.ReadIdentifier())
}

type advanceAndDecodeSuffix struct{}

var global_advanceAndDecodeSuffix advanceAndDecodeSuffix

func (_this advanceAndDecodeSuffix) Run(ctx *DecoderContext) {
	ctx.Stream.AdvanceByte() // Advance past ':'

	panic("TODO: decodeSuffix")
}

// ===========================================================================

func decodeTokenAsBinaryInt(ctx *DecoderContext, token Token) {
	token = token[2:]
	value, bigValue, _, decodedCount := token.DecodeBinaryUint(ctx.TextPos)
	token = token[decodedCount:]
	if bigValue != nil {
		ctx.EventReceiver.OnBigInt(bigValue)
	} else {
		ctx.EventReceiver.OnPositiveInt(value)
	}
	token.AssertAtEnd(ctx.TextPos, "binary integer")
}

func decodeTokenAsOctalInt(ctx *DecoderContext, token Token) {
	token = token[2:]
	value, bigValue, _, decodedCount := token.DecodeOctalUint(ctx.TextPos)
	token = token[decodedCount:]
	if bigValue != nil {
		ctx.EventReceiver.OnBigInt(bigValue)
	} else {
		ctx.EventReceiver.OnPositiveInt(value)
	}
	token.AssertAtEnd(ctx.TextPos, "octal integer")
}

func decodeTokenAsHexNumber(ctx *DecoderContext, token Token) {
	token = token[2:]
	value, bigValue, digitCount, decodedCount := token.DecodeHexUint(ctx.TextPos)
	token = token[decodedCount:]

	if token.IsAtEnd(0) {
		if bigValue != nil {
			ctx.EventReceiver.OnBigInt(bigValue)
		} else {
			ctx.EventReceiver.OnPositiveInt(value)
		}
		return
	}

	if token[0] != '.' {
		token.UnexpectedChar(ctx.TextPos, 0, "hexadecimal int")
	}

	sign := int64(1)
	var fvalue float64
	var bigFValue *big.Float
	fvalue, bigFValue, decodedCount = token.CompleteHexFloat(ctx.TextPos, sign, value, bigValue, digitCount)
	token = token[decodedCount:]
	if bigFValue != nil {
		ctx.EventReceiver.OnBigFloat(bigFValue)
	} else {
		ctx.EventReceiver.OnFloat(fvalue)
	}

	token.AssertAtEnd(ctx.TextPos, "hexadecimal number")
}

func decodeTokenAsUUID(ctx *DecoderContext, token Token) {
	ctx.EventReceiver.OnUUID(token.DecodeUUID(ctx.TextPos))
}

func decodeTokenAsDecimalInt(ctx *DecoderContext, token Token, decodedCount int, value uint64, bigValue *big.Int) {
	token = token[decodedCount:]
	if bigValue != nil {
		ctx.EventReceiver.OnBigInt(bigValue)
	} else {
		ctx.EventReceiver.OnPositiveInt(value)
	}
	token.AssertAtEnd(ctx.TextPos, "decimal integer")
}

func decodeTokenAsDecimalFloat(ctx *DecoderContext, token Token, decodedCount int, digitCount int, value uint64, bigValue *big.Int) {
	sign := int64(1)
	token = token[decodedCount:]
	fvalue, bigFValue, decodedCount := token.CompleteDecimalFloat(ctx.TextPos, sign, value, bigValue, digitCount)
	token = token[decodedCount:]
	if bigFValue != nil {
		ctx.EventReceiver.OnBigDecimalFloat(bigFValue)
	} else {
		ctx.EventReceiver.OnDecimalFloat(fvalue)
	}
	token.AssertAtEnd(ctx.TextPos, "decimal float")
}

func decodeTokenAsTime(ctx *DecoderContext, token Token, decodedCount int, hour int) {
	if decodedCount != 2 {
		token.UnexpectedChar(ctx.TextPos, decodedCount, "time")
	}
	token = token[decodedCount:]
	tvalue, decodedCount := token.CompleteTime(ctx.TextPos, 0, 0, 0, hour)
	token = token[decodedCount:]
	ctx.EventReceiver.OnCompactTime(tvalue)
	token.AssertAtEnd(ctx.TextPos, "time")
}
