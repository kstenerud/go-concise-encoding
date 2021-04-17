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

	"github.com/kstenerud/go-concise-encoding/internal/chars"
	"github.com/kstenerud/go-concise-encoding/internal/common"
)

type decodeInvalidChar struct{}

var global_decodeInvalidChar decodeInvalidChar

func (_this decodeInvalidChar) Run(ctx *DecoderContext) {
	ctx.Stream.Errorf("Unexpected [%v]", ctx.Stream.DescribeCurrentChar())
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
	b := ctx.Stream.PeekByteAllowEOD()
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

	if b := ctx.Stream.ReadByteNoEOD(); b != 'c' && b != 'C' {
		ctx.Stream.Errorf(`Expected document to begin with "c" but got [%v]`, ctx.Stream.DescribeCurrentChar())
	}

	version, bigVersion, digitCount := ctx.Stream.ReadDecimalUint(0, nil)
	if digitCount == 0 {
		ctx.Stream.UnexpectedChar("version number")
	}
	if bigVersion != nil {
		ctx.Stream.Errorf("Version too big")
	}
	// TODO: Remove this when releasing V1
	if version == 1 {
		version = 0
	}

	b := ctx.Stream.PeekByteNoEOD()
	if !chars.ByteHasProperty(b, chars.StructWS) {
		ctx.Stream.UnexpectedChar("whitespace after version")
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
	b := ctx.Stream.ReadByteAllowEOD()
	switch b {
	case '-':
		v := ctx.Stream.ReadDate(int64(coefficient))
		ctx.Stream.AssertAtObjectEnd("date")
		ctx.EventReceiver.OnCompactTime(v)
		return
	case ':':
		v := ctx.Stream.ReadTime(int(coefficient))
		ctx.Stream.AssertAtObjectEnd("time")
		ctx.EventReceiver.OnCompactTime(v)
		return
	case '.':
		value, bigValue, _ := ctx.Stream.ReadDecimalFloat(1, coefficient, bigCoefficient, digitCount)
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
				ctx.EventReceiver.OnBigInt(bigCoefficient)
			} else {
				ctx.EventReceiver.OnPositiveInt(coefficient)
			}
			return
		}
	}
	ctx.Stream.UnexpectedChar("numeric")
}

type advanceAndDecodeNumericNegative struct{}

var global_advanceAndDecodeNumericNegative advanceAndDecodeNumericNegative

func (_this advanceAndDecodeNumericNegative) Run(ctx *DecoderContext) {
	ctx.Stream.AdvanceByte() // Advance past '-'

	switch ctx.Stream.PeekByteNoEOD() {
	case '0':
		ctx.Stream.AdvanceByte() // Advance past '0'
		global_decodeOtherBaseNegative.Run(ctx)
		return
	case 'i':
		namedValue := string(ctx.Stream.ReadNamedValue())
		if namedValue != "inf" {
			ctx.Stream.Errorf("Unknown named value: %v", namedValue)
		}
		ctx.EventReceiver.OnFloat(math.Inf(-1))
		return
	}

	coefficient, bigCoefficient, digitCount := ctx.Stream.ReadDecimalUint(0, nil)
	b := ctx.Stream.ReadByteAllowEOD()
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
	ctx.Stream.UnexpectedChar("numeric")
}

type advanceAndDecodeOtherBasePositive struct{}

var global_advanceAndDecodeOtherBasePositive advanceAndDecodeOtherBasePositive

func (_this advanceAndDecodeOtherBasePositive) Run(ctx *DecoderContext) {
	ctx.Stream.AdvanceByte() // Advance past '0'

	b := ctx.Stream.ReadByteAllowEOD()
	if b.HasProperty(chars.ObjectEnd) {
		ctx.Stream.UnreadByte()
		ctx.EventReceiver.OnPositiveInt(0)
		return
	}

	switch b {
	case 'b':
		v, bigV, _ := ctx.Stream.ReadBinaryUint()
		ctx.Stream.AssertAtObjectEnd("binary integer")
		if bigV != nil {
			ctx.EventReceiver.OnBigInt(bigV)
		} else {
			ctx.EventReceiver.OnPositiveInt(v)
		}
	case 'o':
		v, bigV, _ := ctx.Stream.ReadOctalUint()
		ctx.Stream.AssertAtObjectEnd("octal integer")
		if bigV != nil {
			ctx.EventReceiver.OnBigInt(bigV)
		} else {
			ctx.EventReceiver.OnPositiveInt(v)
		}
	case 'x':
		v, bigV, digitCount := ctx.Stream.ReadHexUint(0, nil)
		if ctx.Stream.PeekByteAllowEOD() == '.' {
			ctx.Stream.AdvanceByte()
			fv, bigFV, _ := ctx.Stream.ReadHexFloat(1, v, bigV, digitCount)
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
		value, bigValue, _ := ctx.Stream.ReadDecimalFloat(1, 0, nil, 0)
		ctx.Stream.AssertAtObjectEnd("float")
		if bigValue != nil {
			ctx.EventReceiver.OnBigDecimalFloat(bigValue)
		} else {
			ctx.EventReceiver.OnDecimalFloat(value)
		}
	default:
		if b.HasProperty(chars.DigitBase10) && ctx.Stream.PeekByteNoEOD() == ':' {
			ctx.Stream.AdvanceByte()
			v := ctx.Stream.ReadTime(int(b - '0'))
			ctx.Stream.AssertAtObjectEnd("time")
			ctx.EventReceiver.OnCompactTime(v)
			return
		}
		ctx.Stream.UnreadByte()
		ctx.Stream.UnexpectedChar("numeric base")
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

	b := ctx.Stream.PeekByteAllowEOD()
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
		if ctx.Stream.PeekByteAllowEOD() == '.' {
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
		ctx.Stream.UnexpectedChar("numeric base")
	}
}

type decodeNamedValueF struct{}

var global_decodeNamedValueF decodeNamedValueF

func (_this decodeNamedValueF) Run(ctx *DecoderContext) {
	namedValue := ctx.Stream.ReadNamedValue()
	switch string(namedValue) {
	case "false":
		ctx.EventReceiver.OnFalse()
	default:
		ctx.Stream.Errorf("%v: Unknown named value", string(namedValue))
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
		ctx.Stream.Errorf("%v: Unknown named value", string(namedValue))
	}
}

type decodeNamedValueN struct{}

var global_decodeNamedValueN decodeNamedValueN

func (_this decodeNamedValueN) Run(ctx *DecoderContext) {
	namedValue := ctx.Stream.ReadNamedValue()
	switch string(namedValue) {
	case "na":
		if ctx.Stream.ReadByteNoEOD() != ':' {
			ctx.Stream.UnreadByte()
			ctx.Stream.UnexpectedChar("NA")
		}
		ctx.EventReceiver.OnNA()
		global_decodeByFirstChar.Run(ctx)
	case "nan":
		ctx.EventReceiver.OnNan(false)
	case "nil":
		ctx.EventReceiver.OnNil()
	default:
		ctx.Stream.Errorf("%v: Unknown named value", string(namedValue))
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
		ctx.Stream.Errorf("%v: Unknown named value", string(namedValue))
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
		ctx.Stream.Errorf("%v: Unknown named value", string(namedValue))
	}
}

type advanceAndDecodeUUID struct{}

var global_advanceAndDecodeUUID advanceAndDecodeUUID

func (_this advanceAndDecodeUUID) Run(ctx *DecoderContext) {
	ctx.Stream.AdvanceByte() // Advance past '@'
	ctx.EventReceiver.OnUUID(ctx.Stream.ExtractUUID(ctx.Stream.ReadNamedValue()))
}

type advanceAndDecodeConstant struct{}

var global_advanceAndDecodeConstant advanceAndDecodeConstant

func (_this advanceAndDecodeConstant) Run(ctx *DecoderContext) {
	ctx.Stream.AdvanceByte() // Advance past '#'

	name := ctx.Stream.ReadIdentifier()
	if ctx.Stream.PeekByteAllowEOD() == ':' {
		ctx.EventReceiver.OnConstant(name, true)
		ctx.Stream.AdvanceByte()
		global_decodeByFirstChar.Run(ctx)
	} else {
		ctx.EventReceiver.OnConstant(name, false)
	}
}

type advanceAndDecodeMarker struct{}

var global_advanceAndDecodeMarker advanceAndDecodeMarker

func (_this advanceAndDecodeMarker) Run(ctx *DecoderContext) {
	ctx.Stream.AdvanceByte() // Advance past '&'

	ctx.EventReceiver.OnMarker(ctx.Stream.ReadIdentifier())
	if ctx.Stream.PeekByteNoEOD() != ':' {
		ctx.Stream.Errorf("Missing colon between marker ID and marked value")
	}
	ctx.Stream.AdvanceByte()
	global_decodeByFirstChar.Run(ctx)
}

type advanceAndDecodeReference struct{}

var global_advanceAndDecodeReference advanceAndDecodeReference

func (_this advanceAndDecodeReference) Run(ctx *DecoderContext) {
	ctx.Stream.AdvanceByte() // Advance past '$'

	if ctx.Stream.PeekByteNoEOD() == '|' {
		ctx.EventReceiver.OnRIDReference()
		ctx.Stream.AdvanceByte()
		arrayType := decodeArrayType(ctx)
		ctx.Stream.SkipWhitespace()
		if arrayType != "r" {
			ctx.Stream.Errorf("%s: Invalid array type for reference ID", arrayType)
		}
		decodeRID(ctx)
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
