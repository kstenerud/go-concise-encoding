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
	"bytes"
	"math"
	"math/big"
	"unicode/utf8"

	"github.com/kstenerud/go-concise-encoding/internal/chars"
	"github.com/kstenerud/go-concise-encoding/internal/common"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
)

type Token []byte

func (_this Token) IsAtEnd(tokenOffset int) bool {
	return tokenOffset == len(_this)
}

func (_this Token) UnexpectedChar(textPos *TextPositionCounter, tokenOffset int, decoding string) {
	_this.adjustTextPosition(textPos, tokenOffset)
	textPos.UnexpectedChar(decoding)
}

func (_this Token) AssertNotEmpty(textPos *TextPositionCounter, decoding string) {
	_this.assertNotEnd(textPos, 0, decoding)
}

func (_this Token) AssertAtEnd(textPos *TextPositionCounter, decoding string) {
	_this.assertPosIsEnd(textPos, 0, decoding)
}

func (_this Token) adjustTextPosition(textPos *TextPositionCounter, tokenOffset int) {
	textPos.Retreat(len(_this)-tokenOffset, chars.ByteWithEOF(_this[tokenOffset]))
}

func (_this Token) errorf(textPos *TextPositionCounter, tokenOffset int, format string, args ...interface{}) {
	_this.adjustTextPosition(textPos, tokenOffset)
	textPos.Errorf(format, args...)
}

func (_this Token) unexpectedError(textPos *TextPositionCounter, err error, decoding string) {
	_this.errorf(textPos, 0, "unexpected error [%v] while decoding %v from [%s]", err, decoding, string(_this))
}

func (_this Token) expectCharAtOffset(textPos *TextPositionCounter, tokenOffset int, ch byte, decoding string) {
	if _this[tokenOffset] != ch {
		_this.UnexpectedChar(textPos, tokenOffset, decoding)
	}
}

func (_this Token) assertNotEnd(textPos *TextPositionCounter, tokenOffset int, decoding string) {
	if tokenOffset >= len(_this) {
		textPos.UnexpectedEOF(decoding)
	}
}

func (_this Token) assertPosIsEnd(textPos *TextPositionCounter, tokenOffset int, decoding string) {
	if tokenOffset < len(_this) {
		r, _ := utf8.DecodeRune(_this[tokenOffset:])
		_this.errorf(textPos, tokenOffset, "unexpected character %c at end of %v", r, decoding)
	}
}

func (_this Token) isEmpty() bool {
	return len(_this) == 0
}

// ----------------------------------------------------------------------------

func (_this Token) DecodeNamedValue(textPos *TextPositionCounter) string {
	_this.assertNotEnd(textPos, 0, "named value")

	common.ASCIIBytesToLower(_this)
	return string(_this)
}

// ----------------------------------------------------------------------------

func (_this Token) DecodeBinaryUint(textPos *TextPositionCounter) (value uint64, bigValue *big.Int, digitCount int, decodedCount int) {
	_this.assertNotEnd(textPos, 0, "binary uint")

	const maxPreShiftBinary = uint64(0x7fffffffffffffff)
	pos := 0

	for ; pos < len(_this); pos++ {
		b := _this[pos]
		if chars.ByteHasProperty(b, chars.DigitBase2) {
			if value > maxPreShiftBinary {
				bigValue = new(big.Int).SetUint64(value)
				break
			}
			nextDigitValue := b - '0'
			value = value<<1 + uint64(nextDigitValue)
			digitCount++
		} else {
			_this.expectCharAtOffset(textPos, pos, charNumericWhitespace, "binary int")
		}
	}

	if bigValue == nil {
		decodedCount = pos
		return
	}

	for ; pos < len(_this); pos++ {
		b := _this[pos]
		if chars.ByteHasProperty(b, chars.DigitBase2) {
			nextDigitValue := b - '0'
			bigValue = bigValue.Mul(bigValue, common.BigInt2)
			bigValue = bigValue.Add(bigValue, big.NewInt(int64(nextDigitValue)))
			digitCount++
		} else {
			_this.expectCharAtOffset(textPos, pos, charNumericWhitespace, "binary int")
		}
	}

	decodedCount = pos
	return
}

func (_this Token) DecodeSmallBinaryUint(textPos *TextPositionCounter) (value uint64, digitCount int, decodedCount int) {
	return _this.decodeSmallUintWrapper(textPos, Token.DecodeBinaryUint)
}

func (_this Token) DecodeSmallBinaryInt(textPos *TextPositionCounter) (value int64, digitCount int, decodedCount int) {
	return _this.decodeSmallIntWrapper(textPos, Token.DecodeBinaryUint)
}

// ----------------------------------------------------------------------------

func (_this Token) DecodeOctalUint(textPos *TextPositionCounter) (value uint64, bigValue *big.Int, digitCount int, decodedCount int) {
	_this.assertNotEnd(textPos, 0, "octal uint")

	const maxPreShiftOctal = uint64(0x1fffffffffffffff)
	pos := 0

	for ; pos < len(_this); pos++ {
		b := _this[pos]
		if chars.ByteHasProperty(b, chars.DigitBase8) {
			if value > maxPreShiftOctal {
				bigValue = new(big.Int).SetUint64(value)
				break
			}
			nextDigitValue := b - '0'
			value = value<<3 + uint64(nextDigitValue)
			digitCount++
		} else {
			_this.expectCharAtOffset(textPos, pos, charNumericWhitespace, "octal int")
		}
	}

	if bigValue == nil {
		decodedCount = pos
		return
	}

	for ; pos < len(_this); pos++ {
		b := _this[pos]
		if chars.ByteHasProperty(b, chars.DigitBase8) {
			nextDigitValue := b - '0'
			bigValue = bigValue.Mul(bigValue, common.BigInt8)
			bigValue = bigValue.Add(bigValue, big.NewInt(int64(nextDigitValue)))
			digitCount++
		} else {
			_this.expectCharAtOffset(textPos, pos, charNumericWhitespace, "octal int")
		}
	}

	decodedCount = pos
	return
}

func (_this Token) DecodeSmallOctalUint(textPos *TextPositionCounter) (value uint64, digitCount int, decodedCount int) {
	return _this.decodeSmallUintWrapper(textPos, Token.DecodeOctalUint)
}

func (_this Token) DecodeSmallOctalInt(textPos *TextPositionCounter) (value int64, digitCount int, decodedCount int) {
	return _this.decodeSmallIntWrapper(textPos, Token.DecodeOctalUint)
}

// ----------------------------------------------------------------------------

func (_this Token) CompleteDecimalUint(textPos *TextPositionCounter, startValue uint64, bigStartValue *big.Int) (value uint64, bigValue *big.Int, digitCount int, decodedCount int) {
	const maxPreShiftDecimal = uint64(1844674407370955161)
	const maxLastDigitDecimal = 5

	// Note: Don't "assert at end" in this function because caller may be reading a potential float

	pos := 0

	if bigStartValue == nil {
		value = startValue
		for ; pos < len(_this); pos++ {
			b := _this[pos]
			if b == charNumericWhitespace {
				continue
			}
			if !chars.ByteHasProperty(b, chars.DigitBase10) {
				decodedCount = pos
				return
			}

			nextDigitValue := b - '0'
			if value > maxPreShiftDecimal || (value == maxPreShiftDecimal && nextDigitValue > maxLastDigitDecimal) {
				bigStartValue = new(big.Int).SetUint64(value)
				break
			}
			value = value*10 + uint64(nextDigitValue)
			digitCount++
		}

		if bigStartValue == nil {
			decodedCount = pos
			return
		}
	}

	bigValue = bigStartValue
	for ; pos < len(_this); pos++ {
		b := _this[pos]
		if b == charNumericWhitespace {
			continue
		}
		if !chars.ByteHasProperty(b, chars.DigitBase10) {
			break
		}

		nextDigitValue := b - '0'
		bigValue = bigValue.Mul(bigValue, common.BigInt10)
		bigValue = bigValue.Add(bigValue, big.NewInt(int64(nextDigitValue)))
		digitCount++
	}

	decodedCount = pos
	return
}

func (_this Token) DecodeDecimalUint(textPos *TextPositionCounter) (value uint64, bigValue *big.Int, digitCount int, decodedCount int) {
	_this.assertNotEnd(textPos, 0, "decimal uint")
	return _this.CompleteDecimalUint(textPos, 0, nil)
}

func (_this Token) DecodeSmallDecimalUint(textPos *TextPositionCounter) (value uint64, digitCount int, decodedCount int) {
	return _this.decodeSmallUintWrapper(textPos, Token.DecodeDecimalUint)
}

func (_this Token) DecodeSmallDecimalInt(textPos *TextPositionCounter) (value int64, digitCount int, decodedCount int) {
	return _this.decodeSmallIntWrapper(textPos, Token.DecodeDecimalUint)
}

// ----------------------------------------------------------------------------

func (_this Token) CompleteHexUint(textPos *TextPositionCounter, startValue uint64, bigStartValue *big.Int) (value uint64, bigValue *big.Int, digitCount int, decodedCount int) {
	const maxPreShiftHex = uint64(0x0fffffffffffffff)
	pos := 0

	if bigStartValue == nil {
		value = startValue
		for ; pos < len(_this); pos++ {
			b := _this[pos]
			var nextNybble byte
			switch {
			case b == charNumericWhitespace:
				continue
			case chars.ByteHasProperty(b, chars.DigitBase10):
				nextNybble = b - '0'
			case chars.ByteHasProperty(b, chars.LowerAF):
				nextNybble = b - 'a' + 10
			case chars.ByteHasProperty(b, chars.UpperAF):
				nextNybble = b - 'A' + 10
			default:
				decodedCount = pos
				return
			}

			if value > maxPreShiftHex {
				bigStartValue = new(big.Int).SetUint64(value)
				break
			}
			value = value<<4 + uint64(nextNybble)
			digitCount++
		}

		if bigStartValue == nil {
			_this.assertPosIsEnd(textPos, pos, "hexadecimal int")
			decodedCount = pos
			return
		}
	}

	bigValue = bigStartValue
	for ; pos < len(_this); pos++ {
		b := _this[pos]
		var nextNybble byte
		switch {
		case b == charNumericWhitespace:
			continue
		case chars.ByteHasProperty(b, chars.DigitBase10):
			nextNybble = b - '0'
		case chars.ByteHasProperty(b, chars.LowerAF):
			nextNybble = b - 'a' + 10
		case chars.ByteHasProperty(b, chars.UpperAF):
			nextNybble = b - 'A' + 10
		default:
			decodedCount = pos
			return
		}

		bigValue = bigValue.Mul(bigValue, common.BigInt16)
		bigValue = bigValue.Add(bigValue, big.NewInt(int64(nextNybble)))
		digitCount++
	}

	decodedCount = pos
	return
}

func (_this Token) DecodeHexUint(textPos *TextPositionCounter) (value uint64, bigValue *big.Int, digitCount int, decodedCount int) {
	_this.assertNotEnd(textPos, 0, "hex uint")
	return _this.CompleteHexUint(textPos, 0, nil)
}

func (_this Token) DecodeSmallHexUint(textPos *TextPositionCounter) (value uint64, digitCount int, decodedCount int) {
	return _this.decodeSmallUintWrapper(textPos, Token.DecodeHexUint)
}

func (_this Token) DecodeSmallHexInt(textPos *TextPositionCounter) (value int64, digitCount int, decodedCount int) {
	return _this.decodeSmallIntWrapper(textPos, Token.DecodeHexUint)
}

// ----------------------------------------------------------------------------

func (_this Token) DecodeUint(textPos *TextPositionCounter) (value uint64, bigValue *big.Int, digitCount int, decodedCount int) {
	_this.assertNotEnd(textPos, 0, "uint")
	if _this[0] == '0' {
		if len(_this) == 1 {
			return 0, nil, 1, 1
		}
		switch _this[1] {
		case 'b', 'B':
			value, bigValue, digitCount, decodedCount = _this[2:].DecodeBinaryUint(textPos)
			decodedCount += 2
			return
		case 'o', 'O':
			value, bigValue, digitCount, decodedCount = _this[2:].DecodeOctalUint(textPos)
			decodedCount += 2
			return
		case 'x', 'X':
			value, bigValue, digitCount, decodedCount = _this[2:].DecodeHexUint(textPos)
			decodedCount += 2
			return
		}
	}

	return _this.DecodeDecimalUint(textPos)
}

func (_this Token) DecodeSmallUint(textPos *TextPositionCounter) (value uint64, digitCount int, decodedCount int) {
	return _this.decodeSmallUintWrapper(textPos, Token.DecodeUint)
}

func (_this Token) DecodeSmallInt(textPos *TextPositionCounter) (value int64, digitCount int, decodedCount int) {
	return _this.decodeSmallIntWrapper(textPos, Token.DecodeUint)
}

// ----------------------------------------------------------------------------

func (_this Token) CompleteDecimalFloat(textPos *TextPositionCounter, sign int64, coefficient uint64, bigCoefficient *big.Int, wholePortionDigitCount int) (value compact_float.DFloat, bigValue *apd.Decimal, decodedCount int) {
	// Assumption: First byte is '.'

	pos := 1
	var exponent int32
	var fractionalDigitCount int
	coefficient, bigCoefficient, fractionalDigitCount, decodedCount = _this[pos:].CompleteDecimalUint(textPos, coefficient, bigCoefficient)
	if fractionalDigitCount == 0 {
		_this.UnexpectedChar(textPos, pos, "decimal float fractional")
	}
	pos += decodedCount

	if !_this.IsAtEnd(pos) && (_this[pos] == 'e' || _this[pos] == 'E') {
		pos++
		_this.assertNotEnd(textPos, pos, "decimal float")
		exponentSign := int32(1)
		switch _this[pos] {
		case '+':
			pos++
			_this.assertNotEnd(textPos, pos, "decimal float")
		case '-':
			exponentSign = -1
			pos++
			_this.assertNotEnd(textPos, pos, "decimal float")
		}
		exp, bigExp, expDigitCount, expDecodedCount := _this[pos:].DecodeDecimalUint(textPos)
		if expDigitCount == 0 {
			_this.UnexpectedChar(textPos, pos, "decimal float exponent")
		}
		if bigExp != nil {
			_this.errorf(textPos, pos, "Exponent %v is too big", bigExp)
		}
		if exp > 0x7fffffff {
			_this.errorf(textPos, pos, "Exponent %v is too big", exp)
		}
		exponent = int32(exp) * exponentSign
		pos += expDecodedCount
	}

	decodedCount = pos
	exponent -= int32(fractionalDigitCount)

	if coefficient == 0 && bigCoefficient == nil {
		if sign < 0 {
			value = compact_float.NegativeZero()
		}
		return
	}

	if bigCoefficient != nil {
		bigValue = apd.NewWithBigInt(bigCoefficient, exponent)
		if sign < 0 {
			bigValue.Negative = true
		}
		return
	}
	if coefficient > 0x7fffffffffffffff {
		bigCoefficient = new(big.Int).SetUint64(coefficient)
		bigValue = apd.NewWithBigInt(bigCoefficient, exponent)
		if sign < 0 {
			bigValue.Negative = true
		}
		return
	}

	value = compact_float.DFloatValue(exponent, int64(coefficient)*sign)
	return
}

func (_this Token) CompleteHexFloat(textPos *TextPositionCounter, sign int64, coefficient uint64, bigCoefficient *big.Int, coefficientDigitCount int) (value float64, bigValue *big.Float, decodedCount int) {
	// Assumption: First byte is '.'

	pos := 1
	var exponent int
	var fractionalDigitCount int
	coefficient, bigCoefficient, fractionalDigitCount, decodedCount = _this[pos:].CompleteHexUint(textPos, coefficient, bigCoefficient)
	if fractionalDigitCount == 0 {
		_this.UnexpectedChar(textPos, pos, "float fractional")
	}
	pos += decodedCount

	if !_this.IsAtEnd(pos) && (_this[pos] == 'p' || _this[pos] == 'P') {
		pos++
		_this.assertNotEnd(textPos, pos, "hex float")
		exponentSign := 1
		switch _this[pos] {
		case '+':
			pos++
			_this.assertNotEnd(textPos, pos, "decimal float")
		case '-':
			exponentSign = -1
			pos++
			_this.assertNotEnd(textPos, pos, "decimal float")
		}
		exp, bigExp, expDigitCount, expDecodedCount := _this[pos:].DecodeDecimalUint(textPos)
		if expDigitCount == 0 {
			_this.UnexpectedChar(textPos, pos, "hex float exponent")
		}
		if bigExp != nil {
			_this.errorf(textPos, pos, "Exponent %v is too big", bigExp)
		}
		// TODO: What is max exponent size?
		exponent = int(exp) * exponentSign
		pos += expDecodedCount
	}

	decodedCount = pos

	adjustedExponent := exponent - fractionalDigitCount*4

	if bigCoefficient != nil {
		bigValue = &big.Float{}
		bigValue = bigValue.SetInt(bigCoefficient)
		if sign < 0 {
			bigValue = bigValue.Neg(bigValue)
		}
		bigValue = bigValue.SetMantExp(bigValue, adjustedExponent)
		return
	}

	if coefficient > common.Float64CoefficientMax {
		bigCoefficient = &big.Int{}
		bigCoefficient = bigCoefficient.SetUint64(coefficient)
		bigValue = &big.Float{}
		bigValue = bigValue.SetInt(bigCoefficient)
		if sign < 0 {
			bigValue = bigValue.Neg(bigValue)
		}
		bigValue = bigValue.SetMantExp(bigValue, adjustedExponent)
		return
	}

	normalizedExponent := exponent - (coefficientDigitCount-1)*4
	if normalizedExponent > common.Float64ExponentMax || normalizedExponent < common.Float64ExponentMin {
		bigValue = &big.Float{}
		bigValue = bigValue.SetInt64(int64(coefficient))
		if sign < 0 {
			bigValue = bigValue.Neg(bigValue)
		}
		bigValue = bigValue.SetMantExp(bigValue, adjustedExponent)
		return
	}

	value = float64(sign) * float64(coefficient) * math.Pow(float64(2), float64(adjustedExponent))
	return
}

func (_this Token) DecodeSmallFloat(textPos *TextPositionCounter) (value float64, decodedCount int) {
	_this.AssertNotEmpty(textPos, "float")

	pos := 0

	sign := int64(1)
	if _this[0] == '-' {
		sign = -1
		pos++
	}

	if len(_this) >= pos+3 && _this[pos] == '0' && (_this[pos+1] == 'x' || _this[pos+1] == 'X') {
		pos += 2
		v, bytesDecoded := _this[pos:].DecodeSmallHexFloat(textPos)
		value = v * float64(sign)
		pos += bytesDecoded
		decodedCount = pos
		return
	}

	coefficient, bigCoefficient, coefficientDigitCount, bydesDecoded := _this[pos:].DecodeUint(textPos)
	pos += bydesDecoded
	if bydesDecoded == 0 {
		common.ASCIIBytesToLower(_this)
		switch {
		case bytes.Equal(_this, byteStringNan):
			value = common.QuietNan
			decodedCount = pos + 3
			return
		case bytes.Equal(_this, byteStringSnan):
			value = common.SignalingNan
			decodedCount = pos + 4
			return
		case bytes.Equal(_this[pos:], byteStringInf):
			value = math.Inf(int(sign))
			decodedCount = pos + 3
			return
		}
		_this.UnexpectedChar(textPos, pos, "float")
	}

	if _this.IsAtEnd(pos) {
		if bigCoefficient != nil {
			// Just disallow this edge case.
			_this.errorf(textPos, pos, "coefficient is too big; use 0x1.1p1 form instead.")
		}
		// Just directly convert. This may cause rounding but it's the best way.
		value = float64(coefficient)
		decodedCount = pos
		return
	}

	if _this[pos] != '.' {
		_this.UnexpectedChar(textPos, pos, "float")
	}
	// Note: Do not advance past '.' because CompleteDecimalFloat expects it
	var dfloatValue compact_float.DFloat
	var bigValue *apd.Decimal
	dfloatValue, bigValue, bydesDecoded = _this[pos:].CompleteDecimalFloat(textPos, sign, coefficient, bigCoefficient, coefficientDigitCount)
	pos += bydesDecoded
	if bydesDecoded == 0 {
		_this.UnexpectedChar(textPos, pos, "float")
	}
	if bigValue != nil {
		_this.errorf(textPos, 0, "float value is too big")
	}
	_this.assertPosIsEnd(textPos, pos, "float")

	// TODO: Check exponent sizes
	value = dfloatValue.Float()
	decodedCount = pos
	return
}

func (_this Token) DecodeSmallHexFloat(textPos *TextPositionCounter) (value float64, decodedCount int) {
	_this.AssertNotEmpty(textPos, "hex float")
	pos := 0

	sign := int64(1)
	if _this[pos] == '-' {
		sign = -1
		pos++
	}

	// Note: The "0x" is implied, and not actually present in the text.

	coefficient, bigCoefficient, coefficientDigitCount, bytesDecoded := _this[pos:].DecodeHexUint(textPos)
	pos += bytesDecoded
	if bytesDecoded == 0 {
		common.ASCIIBytesToLower(_this)
		switch {
		case bytes.Equal(_this, byteStringNan):
			value = common.QuietNan
			decodedCount = pos + 3
			return
		case bytes.Equal(_this, byteStringSnan):
			value = common.SignalingNan
			decodedCount = pos + 4
			return
		case bytes.Equal(_this[pos:], byteStringInf):
			value = math.Inf(int(sign))
			decodedCount = pos + 3
			return
		}
		_this.UnexpectedChar(textPos, pos, "hex float")
	}

	if _this.IsAtEnd(pos) {
		if bigCoefficient != nil {
			// Just disallow this edge case.
			_this.errorf(textPos, pos, "coefficient is too big; use 0x1.1p1 form instead.")
		}
		// Just directly convert. This may cause rounding but it's the best way.
		value = float64(coefficient)
		decodedCount = pos
		return
	}

	if _this[pos] != '.' {
		_this.UnexpectedChar(textPos, pos, "hex float")
	}
	// Note: Do not advance past '.' because CompleteHexFloat expects it
	var bigValue *big.Float
	value, bigValue, bytesDecoded = _this[pos:].CompleteHexFloat(textPos, sign, coefficient, bigCoefficient, coefficientDigitCount)
	pos += bytesDecoded
	if bytesDecoded == 0 {
		_this.UnexpectedChar(textPos, pos, "hex float")
	}
	if bigValue != nil {
		_this.errorf(textPos, 0, "float value is too big")
	}
	_this.assertPosIsEnd(textPos, pos, "hex float")

	decodedCount = pos
	return
}

// ----------------------------------------------------------------------------

type readUintFunc func(Token, *TextPositionCounter) (value uint64, bigValue *big.Int, digitCount int, decodedCount int)

func (_this Token) decodeSmallUintWrapper(textPos *TextPositionCounter, readUint readUintFunc) (value uint64, digitCount int, decodedCount int) {
	var bigValue *big.Int
	value, bigValue, digitCount, decodedCount = readUint(_this, textPos)
	if bigValue != nil {
		_this.errorf(textPos, 0, "Value cannot be > 64 bits")
	}
	return
}

func (_this Token) decodeSmallIntWrapper(textPos *TextPositionCounter, readUint readUintFunc) (value int64, digitCount int, decodedCount int) {
	_this.assertNotEnd(textPos, 0, "integer")
	additionalCount := 0
	sign := int64(1)
	if _this[0] == '-' {
		sign = -sign
		_this = _this[1:]
		additionalCount++
	}
	var bigValue *big.Int
	var uvalue uint64
	uvalue, bigValue, digitCount, decodedCount = readUint(_this, textPos)
	decodedCount += additionalCount
	if bigValue != nil {
		_this.errorf(textPos, 0, "Value cannot be > 64 bits")
	}

	if uvalue > 0x7fffffffffffffff && !(sign < 0 && uvalue == 0x8000000000000000) {
		_this.errorf(textPos, 0, "Integer value too big for element")
	}
	value = int64(uvalue) * sign
	return
}

// ----------------------------------------------------------------------------

func (_this Token) DecodeUID(textPos *TextPositionCounter) (uid []byte) {
	if len(_this) != 36 {
		_this.errorf(textPos, 0, "Expected a UID (xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx)")
	}

	decodeByte := func(src Token) byte {
		nybble1 := chars.HexCharValues[src[0]]
		if nybble1 == chars.InvalidHexChar {
			src.UnexpectedChar(textPos, 0, "UID")
		}
		nybble2 := chars.HexCharValues[src[1]]
		if nybble2 == chars.InvalidHexChar {
			src.UnexpectedChar(textPos, 1, "UID")
		}
		return (nybble1 << 4) | nybble2
	}

	decodeSection := func(src, dst Token, byteCount int) {
		for i := 0; i < byteCount; i++ {
			dst[i] = decodeByte(src)
			src = src[2:]
		}
	}

	expectDash := func(src Token, offset int) {
		_this.expectCharAtOffset(textPos, offset, '-', "UID")
	}

	decodeSection(_this, _this, 4)
	expectDash(_this, 8)
	decodeSection(_this[9:], _this[4:], 2)
	expectDash(_this, 13)
	decodeSection(_this[14:], _this[6:], 2)
	expectDash(_this, 18)
	decodeSection(_this[19:], _this[8:], 2)
	expectDash(_this, 23)
	decodeSection(_this[24:], _this[10:], 6)
	return []byte(_this[:16])
}

// ----------------------------------------------------------------------------

func (_this Token) DecodeUintNoWhitespace(textPos *TextPositionCounter) (value uint64, decodedCount int) {
	// More conservative max to make the calculation easier
	const maxPreShiftDecimal = uint64(1844674407370955161)

	pos := 0
	for ; pos < len(_this); pos++ {
		b := _this[pos]
		if !chars.ByteHasProperty(b, chars.DigitBase10) {
			break
		}

		if value > maxPreShiftDecimal {
			_this.errorf(textPos, pos, "integer value is too big")
		}

		value = value*10 + uint64(b-'0')
	}

	if pos == 0 {
		_this.errorf(textPos, 0, "expected an integer")
	}

	decodedCount = pos
	return
}

// Returns a latitude or longitude value in hundredths of degrees
func (_this Token) DecodeLatOrLong(textPos *TextPositionCounter) (value int, decodedCount int) {
	pos := 0
	_this.assertNotEnd(textPos, pos, "latitude/longitude")

	sign := 1
	if _this[0] == '-' {
		sign = -1
		pos++
	}

	maxCount := len(_this)
	if maxCount > 3 {
		maxCount = 3
	}
	var digitCount int

	for digitCount = 0; digitCount < maxCount; digitCount++ {
		b := _this[pos]
		if !chars.ByteHasProperty(b, chars.DigitBase10) {
			break
		}
		value = value*10 + int(b-'0')
		pos++
	}
	if digitCount == 0 {
		_this.errorf(textPos, 0, "Empty lat/long field")
	}

	if !_this.IsAtEnd(pos) && _this[pos] == '.' {
		pos++
		_this.assertNotEnd(textPos, pos, "latitude/longitude")
		maxCount := len(_this) - pos
		if maxCount > 2 {
			maxCount = 2
		}
		for digitCount = 0; digitCount < maxCount; digitCount++ {
			b := _this[pos]
			if !chars.ByteHasProperty(b, chars.DigitBase10) {
				break
			}
			value = value*10 + int(b-'0')
			pos++
		}
		switch digitCount {
		case 0:
			_this.errorf(textPos, pos, "Missing hundredths portion")
		case 1:
			// Compensate for missing digit
			value *= 10
		}
	} else {
		// Compensate for missing hundredths
		value *= 100
	}

	value *= sign
	decodedCount = pos
	return
}

func (_this Token) CompleteDate(textPos *TextPositionCounter, year int) (t compact_time.Time, decodedCount int) {
	// Assumption: first byte is '-'
	pos := 1
	_this.assertNotEnd(textPos, pos, "date")

	// Month
	month, digitCount := _this[pos:].DecodeUintNoWhitespace(textPos)
	if digitCount == 0 {
		_this.UnexpectedChar(textPos, pos, "month")
	}
	if digitCount > 2 {
		_this.errorf(textPos, pos, "Month field is too long")
	}
	if month < 1 || month > 12 {
		_this.errorf(textPos, pos, "Month %v is invalid", month)
	}
	pos += digitCount
	_this.assertNotEnd(textPos, pos, "date")
	_this.expectCharAtOffset(textPos, pos, '-', "date")
	pos++

	// Day
	var day uint64
	day, digitCount = _this[pos:].DecodeUintNoWhitespace(textPos)
	if digitCount == 0 {
		_this.UnexpectedChar(textPos, pos, "day")
	}
	if digitCount > 2 {
		_this.errorf(textPos, pos, "Day field is too long")
	}
	if day < 1 || int(day) > common.MaxDayByMonth[month] {
		_this.errorf(textPos, pos, "Day %v is invalid", day)
	}
	pos += digitCount

	if _this.IsAtEnd(pos) {
		t = compact_time.NewDate(int(year), int(month), int(day))
		if err := t.Validate(); err != nil {
			_this.unexpectedError(textPos, err, "date")
		}
		decodedCount = pos
		return
	}

	// Timestamp
	_this.expectCharAtOffset(textPos, pos, '/', "date")
	pos++
	_this.assertNotEnd(textPos, pos, "date")

	var hour uint64
	hour, digitCount = _this[pos:].DecodeUintNoWhitespace(textPos)
	if digitCount < 1 || digitCount > 2 {
		_this.errorf(textPos, pos, "invalid hour value")
	}
	pos += digitCount
	t, decodedCount = _this[pos:].CompleteTime(textPos, year, int(month), int(day), int(hour))
	decodedCount += pos
	return
}

func (_this Token) decodeUTCOffset(textPos *TextPositionCounter, sign int) (minutesOffset int, decodedCount int) {
	if len(_this) != 4 {
		_this.errorf(textPos, 0, "%v: invalid UTC offset (must be in the format -1234 or +1234)", string(_this))
	}
	hhmm, decodedCount := _this.DecodeUintNoWhitespace(textPos)
	if decodedCount != 4 {
		_this.errorf(textPos, 0, "%v: invalid UTC offset (must be in the format -1234 or +1234)", string(_this))
	}
	hour := int(hhmm / 100)
	minute := int(hhmm % 100)
	if hour > 23 {
		_this.errorf(textPos, 0, "%v: invalid UTC offset hour (max is 23)", hour)
	}
	if minute > 59 {
		_this.errorf(textPos, 0, "%v: invalid UTC offset minute (max is 59)", minute)
	}
	minutesOffset = (hour*60 + minute) * sign
	return
}

// Complete a time value. Pass 0 as year, month, and day to indicate no date portion.
func (_this Token) CompleteTime(textPos *TextPositionCounter, year, month, day, hour int) (t compact_time.Time, decodedCount int) {
	// Assumption: first byte is ':'
	pos := 1
	_this.assertNotEnd(textPos, pos, "time")

	timeType := compact_time.TimeTypeTime

	if day != 0 {
		timeType = compact_time.TimeTypeTimestamp
	}

	// Minute
	minute, digitCount := _this[pos:].DecodeUintNoWhitespace(textPos)
	if digitCount != 2 {
		_this.errorf(textPos, pos, "Minute field must be 2 characters long")
	}
	if minute > 59 {
		_this.errorf(textPos, pos, "Minute %v is invalid", minute)
	}
	pos += digitCount
	_this.assertNotEnd(textPos, pos, "time")
	_this.expectCharAtOffset(textPos, pos, ':', "time")
	pos++

	// Second
	var second uint64
	second, digitCount = _this[pos:].DecodeUintNoWhitespace(textPos)
	if digitCount != 2 {
		_this.errorf(textPos, pos, "Second field must be 2 characters long")
	}
	if second > 60 {
		_this.errorf(textPos, pos, "Second %v is invalid", second)
	}
	pos += digitCount

	// Nanosecond
	var nsec uint64
	if !_this.IsAtEnd(pos) && _this[pos] == '.' {
		pos++
		nsec, digitCount = _this[pos:].DecodeUintNoWhitespace(textPos)
		// TODO: Check for overflow
		nsec *= uint64(subsecondMagnitudes[digitCount])
		if digitCount == 0 {
			_this.UnexpectedChar(textPos, pos, "nanosecond")
		}
		if nsec > 999999999 {
			_this.errorf(textPos, pos, "nanosecond %v is invalid", nsec)
		}
		pos += digitCount
	}

	// Timezone
	var tz compact_time.Timezone
	if _this.IsAtEnd(pos) {
		tz = compact_time.TZAtUTC()
	} else {
		switch _this[pos] {
		case '+':
			pos++
			sign := 1
			minutes, decodedCount := _this[pos:].decodeUTCOffset(textPos, sign)
			tz.InitWithMinutesOffsetFromUTC(minutes)
			pos += decodedCount
		case '-':
			pos++
			sign := -1
			minutes, decodedCount := _this[pos:].decodeUTCOffset(textPos, sign)
			tz.InitWithMinutesOffsetFromUTC(minutes)
			pos += decodedCount
		case '/':
			pos++
			_this.assertNotEnd(textPos, pos, "time zone")

			if chars.ByteHasProperty(_this[pos], chars.DigitBase10) || _this[pos] == '-' {
				latitude, decodedCount := _this[pos:].DecodeLatOrLong(textPos)
				if latitude < -9000 || latitude > 9000 {
					_this.errorf(textPos, pos, "Latitude %v is invalid", float64(latitude)/100)
				}
				pos += decodedCount
				_this.expectCharAtOffset(textPos, pos, '/', "time zone")
				pos++
				var longitude int
				longitude, decodedCount = _this[pos:].DecodeLatOrLong(textPos)
				if longitude < -18000 || longitude > 18000 {
					_this.errorf(textPos, pos, "Longitude %v is invalid", float64(longitude)/100)
				}
				pos += decodedCount
				tz.InitWithLatLong(latitude, longitude)
				_this.assertPosIsEnd(textPos, pos, "time zone")
			} else if chars.ByteHasProperty(_this[pos], chars.AZ) {
				tz.InitWithAreaLocation(string(_this[pos:]))
				pos = len(_this)
				_this.assertPosIsEnd(textPos, pos, "time zone")
			} else {
				_this.UnexpectedChar(textPos, pos, "time zone")
			}
		default:
			_this.UnexpectedChar(textPos, pos, "time zone")
		}
	}

	switch timeType {
	case compact_time.TimeTypeTime:
		t.InitTime(hour, int(minute), int(second), int(nsec), tz)
	case compact_time.TimeTypeTimestamp:
		t.InitTimestamp(year, month, day, hour, int(minute), int(second), int(nsec), tz)
	}

	if err := t.Validate(); err != nil {
		_this.unexpectedError(textPos, err, "time")
	}
	decodedCount = pos
	return
}
