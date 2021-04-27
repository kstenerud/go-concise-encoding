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
	"io"
	"math"
	"math/big"
	"unicode/utf8"

	"github.com/kstenerud/go-concise-encoding/internal/chars"
	"github.com/kstenerud/go-concise-encoding/internal/common"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
)

type Reader struct {
	reader    io.Reader
	lastByte  chars.ByteWithEOF
	hasUnread bool
	isEOF     bool
	byteBuff  [1]byte

	TextPos TextPositionCounter

	token            []byte
	verbatimSentinel []byte
}

func NewReader(reader io.Reader) *Reader {
	_this := &Reader{}
	_this.Init(reader)
	return _this
}

// Init the read buffer. You may call this again to re-initialize the buffer.
func (_this *Reader) Init(reader io.Reader) {
	_this.reader = reader
	if cap(_this.token) == 0 {
		_this.token = make([]byte, 0, 16)
	}
	_this.hasUnread = false
	_this.isEOF = false
}

// Bytes

func (_this *Reader) readNext() {
	// Can't create local [1]byte here because it mallocs? WTF???
	if _, err := _this.reader.Read(_this.byteBuff[:]); err == nil {
		_this.lastByte = chars.ByteWithEOF(_this.byteBuff[0])
	} else {
		if err == io.EOF {
			_this.isEOF = true
			_this.lastByte = chars.EOFMarker
		} else {
			panic(err)
		}
	}
}

func (_this *Reader) UnreadByte() {
	if _this.hasUnread {
		panic("Cannot unread twice")
	}
	_this.hasUnread = true
	_this.TextPos.RetreatOneChar()
}

func (_this *Reader) PeekByteAllowEOF() chars.ByteWithEOF {
	if !_this.hasUnread {
		_this.readNext()
		_this.hasUnread = true
	}
	return _this.lastByte
}

func (_this *Reader) ReadByteAllowEOF() chars.ByteWithEOF {
	if !_this.hasUnread {
		_this.readNext()
	}

	_this.hasUnread = false
	b := _this.lastByte
	_this.TextPos.Advance(b)
	return b
}

func (_this *Reader) PeekByteNoEOF() byte {
	b := _this.PeekByteAllowEOF()
	if b == chars.EOFMarker {
		_this.unexpectedEOF()
	}
	return byte(b)
}

func (_this *Reader) ReadByteNoEOF() byte {
	b := _this.ReadByteAllowEOF()
	if b == chars.EOFMarker {
		_this.unexpectedEOF()
	}
	return byte(b)
}

func (_this *Reader) AdvanceByte() {
	_this.ReadByteNoEOF()
}

func (_this *Reader) SkipWhileProperty(property chars.Properties) {
	for _this.ReadByteAllowEOF().HasProperty(property) {
	}
	_this.UnreadByte()
}

// Tokens

func (_this *Reader) TokenBegin() {
	_this.token = _this.token[:0]
}

func (_this *Reader) TokenStripLastByte() {
	_this.token = _this.token[:len(_this.token)-1]
}

func (_this *Reader) TokenStripLastBytes(count int) {
	_this.token = _this.token[:len(_this.token)-count]
}

func (_this *Reader) TokenAppendByte(b byte) {
	_this.token = append(_this.token, b)
}

func (_this *Reader) TokenAppendBytes(b []byte) {
	_this.token = append(_this.token, b...)
}

func (_this *Reader) TokenAppendRune(r rune) {
	if r < utf8.RuneSelf {
		_this.TokenAppendByte(byte(r))
	} else {
		pos := len(_this.token)
		_this.TokenAppendBytes([]byte{0, 0, 0, 0, 0})
		length := utf8.EncodeRune(_this.token[pos:], r)
		_this.token = _this.token[:pos+length]
	}
}

func (_this *Reader) TokenGet() []byte {
	return _this.token
}

func (_this *Reader) TokenReadByteNoEOF() byte {
	b := _this.ReadByteNoEOF()
	_this.TokenAppendByte(b)
	return b
}

func (_this *Reader) TokenReadByteAllowEOF() chars.ByteWithEOF {
	b := _this.ReadByteAllowEOF()
	if b != chars.EOFMarker {
		_this.TokenAppendByte(byte(b))
	}
	return b
}

func (_this *Reader) TokenReadUntilAndIncludingByte(untilByte byte) {
	for {
		b := _this.ReadByteNoEOF()
		if b == untilByte {
			break
		}
		_this.TokenAppendByte(b)
	}
}

func (_this *Reader) TokenReadUntilPropertyNoEOF(property chars.Properties) {
	for {
		b := _this.ReadByteNoEOF()
		if chars.ByteHasProperty(b, property) {
			_this.UnreadByte()
			break
		}
		_this.TokenAppendByte(b)
	}
}

func (_this *Reader) TokenReadUntilPropertyAllowEOF(property chars.Properties) {
	for {
		b := _this.ReadByteAllowEOF()
		if b.HasProperty(property) {
			_this.UnreadByte()
			break
		}
		_this.TokenAppendByte(byte(b))
	}
}

func (_this *Reader) TokenReadWhilePropertyAllowEOF(property chars.Properties) {
	for {
		b := _this.ReadByteAllowEOF()
		if !b.HasProperty(property) {
			_this.UnreadByte()
			break
		}
		_this.TokenAppendByte(byte(b))
	}
}

// Errors

func (_this *Reader) AssertAtObjectEnd(decoding string) {
	if !_this.PeekByteAllowEOF().HasProperty(chars.ObjectEnd) {
		_this.unexpectedChar(decoding)
	}
}

func (_this *Reader) errorf(format string, args ...interface{}) {
	_this.TextPos.Errorf(format, args...)
}

func (_this *Reader) unexpectedEOF() {
	_this.errorf("unexpected end of document")
}

func (_this *Reader) unexpectedError(err error, decoding string) {
	_this.errorf("unexpected error [%v] while decoding %v", err, decoding)
}

func (_this *Reader) unexpectedChar(decoding string) {
	_this.TextPos.UnexpectedChar(decoding)
}

// Decoders

func (_this *Reader) SkipWhitespace() {
	_this.SkipWhileProperty(chars.StructWS)
}

func (_this *Reader) ReadToken() Token {
	_this.TokenBegin()
	_this.TokenReadUntilPropertyAllowEOF(chars.ObjectEnd)
	return _this.TokenGet()
}

func (_this *Reader) ReadSmallBinaryUint() (value uint64, digitCount int) {
	v, vBig, count := _this.ReadBinaryUint()
	if vBig != nil {
		_this.errorf("Value cannot be > 64 bits")
	}
	return v, count
}

func (_this *Reader) ReadSmallBinaryInt() (value int64, digitCount int) {
	sign := int64(1)
	if _this.ReadByteAllowEOF() == '-' {
		sign = -sign
	} else {
		_this.UnreadByte()
	}
	v, count := _this.ReadSmallBinaryUint()

	if v > 0x7fffffffffffffff && !(sign < 0 && v == 0x8000000000000000) {
		_this.errorf("Integer value too big for element")
	}
	return int64(v) * sign, count
}

func (_this *Reader) ReadBinaryUint() (value uint64, bigValue *big.Int, digitCount int) {
	const maxPreShiftBinary = uint64(0x7fffffffffffffff)

	for {
		b := _this.ReadByteAllowEOF()
		switch {
		case b == charNumericWhitespace:
			continue
		case b.HasProperty(chars.DigitBase2):
			// Nothing to do
		default:
			_this.UnreadByte()
			return
		}

		if value > maxPreShiftBinary {
			bigValue = new(big.Int).SetUint64(value)
			_this.UnreadByte()
			break
		}
		nextDigitValue := b - '0'
		value = value<<1 + uint64(nextDigitValue)
		digitCount++
	}

	if bigValue == nil {
		return
	}

	for {
		b := _this.ReadByteAllowEOF()
		switch {
		case b == charNumericWhitespace:
			continue
		case b.HasProperty(chars.DigitBase2):
			// Nothing to do
		default:
			_this.UnreadByte()
			return
		}

		nextDigitValue := b - '0'
		bigValue = bigValue.Mul(bigValue, common.BigInt2)
		bigValue = bigValue.Add(bigValue, big.NewInt(int64(nextDigitValue)))
		digitCount++
	}

	return
}

const maxPreShiftOctal = uint64(0x1fffffffffffffff)

func (_this *Reader) ReadSmallOctalUint() (value uint64, digitCount int) {
	v, vBig, count := _this.ReadOctalUint()
	if vBig != nil {
		_this.errorf("Value cannot be > 64 bits")
	}
	return v, count
}

func (_this *Reader) ReadSmallOctalInt() (value int64, digitCount int) {
	sign := int64(1)
	if _this.ReadByteAllowEOF() == '-' {
		sign = -sign
	} else {
		_this.UnreadByte()
	}
	v, count := _this.ReadSmallOctalUint()

	if v > 0x7fffffffffffffff && !(sign < 0 && v == 0x8000000000000000) {
		_this.errorf("Integer value too big for element")
	}
	return int64(v) * sign, count
}

func (_this *Reader) ReadOctalUint() (value uint64, bigValue *big.Int, digitCount int) {
	for {
		b := _this.ReadByteAllowEOF()
		switch {
		case b == charNumericWhitespace:
			continue
		case b.HasProperty(chars.DigitBase8):
			// Nothing to do
		default:
			_this.UnreadByte()
			return
		}

		if value > maxPreShiftOctal {
			bigValue = new(big.Int).SetUint64(value)
			_this.UnreadByte()
			break
		}
		nextDigitValue := b - '0'
		value = value<<3 + uint64(nextDigitValue)
		digitCount++
	}

	if bigValue == nil {
		return
	}

	for {
		b := _this.ReadByteAllowEOF()
		switch {
		case b == charNumericWhitespace:
			continue
		case b.HasProperty(chars.DigitBase8):
			// Nothing to do
		default:
			_this.UnreadByte()
			return
		}

		nextDigitValue := b - '0'
		bigValue = bigValue.Mul(bigValue, common.BigInt8)
		bigValue = bigValue.Add(bigValue, big.NewInt(int64(nextDigitValue)))
		digitCount++
	}

	return
}

// startValue is only used if bigStartValue is nil
// bigValue will be nil unless the value was too big for a uint64
func (_this *Reader) ReadDecimalUint(startValue uint64, bigStartValue *big.Int) (value uint64, bigValue *big.Int, digitCount int) {
	const maxPreShiftDecimal = uint64(1844674407370955161)
	const maxLastDigitDecimal = 5

	if bigStartValue == nil {
		value = startValue
		for {
			b := _this.ReadByteAllowEOF()
			switch {
			case b == charNumericWhitespace:
				continue
			case b.HasProperty(chars.DigitBase10):
				// Nothing to do
			default:
				_this.UnreadByte()
				return
			}

			nextDigitValue := b - '0'
			if value > maxPreShiftDecimal || (value == maxPreShiftDecimal && nextDigitValue > maxLastDigitDecimal) {
				bigStartValue = new(big.Int).SetUint64(value)
				_this.UnreadByte()
				break
			}
			value = value*10 + uint64(nextDigitValue)
			digitCount++
		}

		if bigStartValue == nil {
			return
		}
	}

	bigValue = bigStartValue
	for {
		b := _this.ReadByteAllowEOF()
		switch {
		case b == charNumericWhitespace:
			continue
		case b.HasProperty(chars.DigitBase10):
			// Nothing to do
		default:
			_this.UnreadByte()
			return
		}

		nextDigitValue := b - '0'
		bigValue = bigValue.Mul(bigValue, common.BigInt10)
		bigValue = bigValue.Add(bigValue, big.NewInt(int64(nextDigitValue)))
		digitCount++
	}

	return
}

const maxPreShiftHex = uint64(0x0fffffffffffffff)

func (_this *Reader) ReadSmallHexUint() (value uint64, digitCount int) {
	v, vBig, count := _this.ReadHexUint(0, nil)
	if vBig != nil {
		_this.errorf("Value cannot be > 64 bits")
	}
	return v, count
}

func (_this *Reader) ReadSmallHexInt() (value int64, digitCount int) {
	sign := int64(1)
	if _this.ReadByteAllowEOF() == '-' {
		sign = -sign
	} else {
		_this.UnreadByte()
	}
	v, count := _this.ReadSmallHexUint()

	if v > 0x7fffffffffffffff && !(sign < 0 && v == 0x8000000000000000) {
		_this.errorf("Integer value too big for element")
	}
	return int64(v) * sign, count
}

func (_this *Reader) ReadHexUint(startValue uint64, bigStartValue *big.Int) (value uint64, bigValue *big.Int, digitCount int) {
	if bigStartValue == nil {
		value = startValue
		for {
			b := _this.ReadByteAllowEOF()
			nextNybble := chars.ByteWithEOF(0)
			switch {
			case b == charNumericWhitespace:
				continue
			case b.HasProperty(chars.DigitBase10):
				nextNybble = b - '0'
			case b.HasProperty(chars.LowerAF):
				nextNybble = b - 'a' + 10
			case b.HasProperty(chars.UpperAF):
				nextNybble = b - 'A' + 10
			default:
				_this.UnreadByte()
				return
			}

			if value > maxPreShiftHex {
				bigStartValue = new(big.Int).SetUint64(value)
				_this.UnreadByte()
				break
			}
			value = value<<4 + uint64(nextNybble)
			digitCount++
		}

		if bigStartValue == nil {
			return
		}
	}

	bigValue = bigStartValue
	for {
		b := _this.ReadByteAllowEOF()
		nextNybble := chars.ByteWithEOF(0)
		switch {
		case b == charNumericWhitespace:
			continue
		case b.HasProperty(chars.DigitBase10):
			nextNybble = b - '0'
		case b.HasProperty(chars.LowerAF):
			nextNybble = b - 'a' + 10
		case b.HasProperty(chars.UpperAF):
			nextNybble = b - 'A' + 10
		default:
			_this.UnreadByte()
			return
		}

		bigValue = bigValue.Mul(bigValue, common.BigInt16)
		bigValue = bigValue.Add(bigValue, big.NewInt(int64(nextNybble)))
		digitCount++
	}

	return
}

func (_this *Reader) ReadSmallUint() (value uint64, digitCount int) {
	var bigV *big.Int

	if _this.ReadByteAllowEOF() == '0' {
		switch _this.ReadByteAllowEOF() {
		case 'b', 'B':
			value, bigV, digitCount = _this.ReadBinaryUint()
		case 'o', 'O':
			value, bigV, digitCount = _this.ReadOctalUint()
		case 'x', 'X':
			value, bigV, digitCount = _this.ReadHexUint(0, nil)
		default:
			_this.UnreadByte()
			value, bigV, digitCount = _this.ReadDecimalUint(0, nil)
			digitCount++
		}
	} else {
		_this.UnreadByte()
		value, bigV, digitCount = _this.ReadDecimalUint(0, nil)
	}

	if bigV != nil {
		_this.errorf("Integer value too big for element")
	}
	return
}

func (_this *Reader) ReadSmallInt() (value int64, digitCount int) {
	sign := int64(1)
	if _this.ReadByteAllowEOF() == '-' {
		sign = -sign
	} else {
		_this.UnreadByte()
	}

	v, count := _this.ReadSmallUint()

	if v > 0x7fffffffffffffff && !(sign < 0 && v == 0x8000000000000000) {
		_this.errorf("Integer value too big for element")
	}
	return int64(v) * sign, count
}

func (_this *Reader) ReadDecimalFloat(sign int64, coefficient uint64, bigCoefficient *big.Int, coefficientDigitCount int) (value compact_float.DFloat, bigValue *apd.Decimal, digitCount int) {
	exponent := int32(0)
	fractionalDigitCount := 0
	coefficient, bigCoefficient, fractionalDigitCount = _this.ReadDecimalUint(coefficient, bigCoefficient)
	if fractionalDigitCount == 0 {
		_this.unexpectedChar("float fractional")
	}
	digitCount = coefficientDigitCount + fractionalDigitCount

	b := _this.PeekByteAllowEOF()
	if b == 'e' || b == 'E' {
		_this.AdvanceByte()
		exponentSign := int32(1)
		switch _this.PeekByteNoEOF() {
		case '+':
			_this.AdvanceByte()
		case '-':
			exponentSign = -1
			_this.AdvanceByte()
		}
		exp, bigExp, expDigitCount := _this.ReadDecimalUint(0, nil)
		if expDigitCount == 0 {
			_this.unexpectedChar("float exponent")
		}
		if bigExp != nil || exp > 0x7fffffff {
			_this.errorf("Exponent too big")
		}
		exponent = int32(exp) * exponentSign
	}

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

func (_this *Reader) ReadHexFloat(sign int64, coefficient uint64, bigCoefficient *big.Int, coefficientDigitCount int) (value float64, bigValue *big.Float, digitCount int) {

	exponent := 0
	fractionalDigitCount := 0
	coefficient, bigCoefficient, fractionalDigitCount = _this.ReadHexUint(coefficient, bigCoefficient)
	b := _this.PeekByteAllowEOF()
	if fractionalDigitCount == 0 && b != 'p' && b != 'P' {
		_this.unexpectedChar("hex float fractional")
	}
	digitCount = coefficientDigitCount + fractionalDigitCount

	if b == 'p' || b == 'P' {
		_this.AdvanceByte()
		exponentSign := 1
		switch _this.PeekByteNoEOF() {
		case '+':
			_this.AdvanceByte()
		case '-':
			exponentSign = -1
			_this.AdvanceByte()
		}
		exp, bigExp, expDigitCount := _this.ReadDecimalUint(0, nil)
		if expDigitCount == 0 {
			_this.unexpectedChar("hex float exponent")
		}
		if bigExp != nil {
			_this.errorf("Exponent too big")
		}
		exponent = int(exp) * exponentSign
	}

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

	if coefficient > maxFloat64Coefficient {
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
	if normalizedExponent > maxFloat64Exponent || normalizedExponent < minFloat64Exponent {
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

func (_this *Reader) ReadSmallHexFloat() (value float64, digitCount int) {
	sign := int64(1)
	b := _this.PeekByteAllowEOF()
	if b == '-' {
		sign = -sign
		_this.AdvanceByte()
	} else if !b.HasProperty(chars.DigitBase10 | chars.LowerAF | chars.UpperAF) {
		return
	}

	u, bigU, coefficientDigitCount := _this.ReadHexUint(0, nil)
	if coefficientDigitCount == 0 {
		_this.unexpectedChar("hex float coefficient")
	}
	if bigU != nil || u > maxFloat64Coefficient {
		_this.errorf("Value too big for element")
	}
	b = _this.PeekByteAllowEOF()
	switch {
	case b == '.':
		_this.AdvanceByte()
		f, bigF, digitCount := _this.ReadHexFloat(sign, u, nil, coefficientDigitCount)
		if bigF != nil {
			_this.errorf("Value too big for element")
		}
		return f, digitCount
	case b == 'p' || b == 'P':
		f, bigF, digitCount := _this.ReadHexFloat(sign, u, nil, coefficientDigitCount)
		if bigF != nil {
			_this.errorf("Value too big for element")
		}
		return f, digitCount
	case b.HasProperty(chars.ObjectEnd):
		return float64(u) * float64(sign), coefficientDigitCount
	default:
		_this.unexpectedChar("hex float")
		return 0, 0
	}
}

func (_this *Reader) ReadSmallFloat() (value float64, digitCount int) {
	sign := int64(1)
	b := _this.ReadByteAllowEOF()
	if b == '-' {
		sign = -sign
	} else if !b.HasProperty(chars.DigitBase10) {
		_this.UnreadByte()
		return
	} else {
		_this.UnreadByte()
	}

	initialZero := false
	if _this.ReadByteAllowEOF() == '0' {
		initialZero = true
		b = _this.ReadByteAllowEOF()
		if b == 'x' || b == 'X' {
			u, bigU, coefficientDigitCount := _this.ReadHexUint(0, nil)
			if coefficientDigitCount == 0 {
				_this.unexpectedChar("float")
			}
			if bigU != nil || u > maxFloat64Coefficient {
				_this.errorf("Value too big for element")
			}

			b = _this.ReadByteAllowEOF()
			switch {
			case b == '.':
				f, bigF, digitCount := _this.ReadHexFloat(sign, u, nil, coefficientDigitCount)
				if bigF != nil {
					_this.errorf("Value too big for element")
				}
				return f, digitCount
			case b.HasProperty(chars.StructWS):
				_this.UnreadByte()
				return float64(u) * float64(sign), coefficientDigitCount
			default:
				_this.unexpectedChar("float")
				return 0, 0
			}
		} else {
			_this.UnreadByte()
		}
	} else {
		_this.UnreadByte()
	}

	u, bigU, coefficientDigitCount := _this.ReadDecimalUint(0, nil)
	if initialZero {
		coefficientDigitCount++
	}
	if coefficientDigitCount == 0 {
		_this.unexpectedChar("float")
	}
	if bigU != nil || u > maxFloat64Coefficient {
		_this.errorf("Value too big for element")
	}

	b = _this.ReadByteAllowEOF()
	switch {
	case b == '.':
		f, bigF, digitCount := _this.ReadDecimalFloat(sign, u, nil, coefficientDigitCount)
		if bigF != nil {
			_this.errorf("Value too big for element")
		}
		normalizedExponent := int(f.Exponent) + digitCount - 1
		if normalizedExponent < minFloat64DecimalExponent || normalizedExponent > maxFloat64DecimalExponent {
			_this.errorf("Value too big for element")
		}

		return f.Float(), digitCount
	case b.HasProperty(chars.StructWS):
		_this.UnreadByte()
		return float64(u) * float64(sign), coefficientDigitCount
	default:
		_this.unexpectedChar("float")
		return 0, 0
	}

}

func (_this *Reader) ReadNamedValue() []byte {
	_this.TokenBegin()
	_this.TokenReadWhilePropertyAllowEOF(chars.AZ)
	namedValue := _this.TokenGet()
	if len(namedValue) == 0 {
		_this.unexpectedChar("name")
	}
	common.ASCIIBytesToLower(namedValue)
	return namedValue
}

func (_this *Reader) TokenReadVerbatimSequence() {
	_this.verbatimSentinel = _this.verbatimSentinel[:0]
	for {
		b := _this.ReadByteNoEOF()
		if chars.ByteHasProperty(b, chars.StructWS) {
			if b == '\r' {
				if _this.ReadByteNoEOF() != '\n' {
					_this.unexpectedChar("verbatim sentinel")
				}
			}
			break
		}
		_this.verbatimSentinel = append(_this.verbatimSentinel, b)
	}

	sentinelLength := len(_this.verbatimSentinel)

Outer:
	for {
		_this.TokenReadByteNoEOF()
		for i := 1; i <= sentinelLength; i++ {
			if _this.token[len(_this.token)-i] != _this.verbatimSentinel[sentinelLength-i] {
				continue Outer
			}
		}

		_this.TokenStripLastBytes(sentinelLength)
		return
	}
}

func (_this *Reader) TokenReadEscape() {
	escape := _this.ReadByteNoEOF()
	switch escape {
	case 't':
		_this.TokenAppendByte('\t')
	case 'n':
		_this.TokenAppendByte('\n')
	case 'r':
		_this.TokenAppendByte('\r')
	case '"', '*', '/', '<', '>', '\\', '|':
		_this.TokenAppendByte(escape)
	case '_':
		// Non-breaking space
		_this.TokenAppendBytes([]byte{0xc0, 0xa0})
	case '-':
		// Soft hyphen
		_this.TokenAppendBytes([]byte{0xc0, 0xad})
	case '\r', '\n':
		// Continuation
		_this.SkipWhitespace()
	case '0':
		_this.TokenAppendByte(0)
	case '1', '2', '3', '4', '5', '6', '7', '8', '9':
		length := int(escape - '0')
		codepoint := rune(0)
		for i := 0; i < length; i++ {
			b := _this.ReadByteNoEOF()
			switch {
			case chars.ByteHasProperty(b, chars.DigitBase10):
				codepoint = (codepoint << 4) | (rune(b) - '0')
			case chars.ByteHasProperty(b, chars.LowerAF):
				codepoint = (codepoint << 4) | (rune(b) - 'a' + 10)
			case chars.ByteHasProperty(b, chars.UpperAF):
				codepoint = (codepoint << 4) | (rune(b) - 'A' + 10)
			default:
				_this.unexpectedChar("unicode escape")
			}
		}

		_this.TokenAppendRune(codepoint)
	case '.':
		_this.TokenReadVerbatimSequence()
	default:
		_this.unexpectedChar("escape sequence")
	}
}

func (_this *Reader) ReadQuotedString() []byte {
	_this.TokenBegin()
	for {
		b := _this.TokenReadByteNoEOF()
		switch b {
		case '"':
			_this.TokenStripLastByte()
			return _this.TokenGet()
		case '\\':
			_this.TokenStripLastByte()
			_this.TokenReadEscape()
		}
	}
}

func (_this *Reader) ReadIdentifier() []byte {
	_this.TokenBegin()
	for {
		b := _this.ReadByteAllowEOF()
		// Only do a per-byte check here. The rules will do a per-rune check.
		if b == chars.EOFMarker || (b < 0x80 && !chars.IsRuneValidIdentifier(rune(b))) {
			_this.UnreadByte()
			break
		}
		_this.TokenAppendByte(byte(b))
	}
	return _this.TokenGet()
}

func (_this *Reader) ReadStringArray() []byte {
	_this.TokenBegin()
	for {
		b := _this.TokenReadByteNoEOF()
		switch b {
		case '|':
			_this.TokenStripLastByte()
			return _this.TokenGet()
		case '\\':
			_this.TokenStripLastByte()
			_this.TokenReadEscape()
		}
	}
}

func trimWhitespace(str []byte) []byte {
	for len(str) > 0 && chars.ByteHasProperty(str[0], chars.StructWS) {
		str = str[1:]
	}
	for len(str) > 0 && chars.ByteHasProperty(str[len(str)-1], chars.StructWS) {
		str = str[:len(str)-1]
	}
	return str
}

func trimWhitespaceMarkupContent(str []byte) []byte {
	for len(str) > 0 && chars.ByteHasProperty(str[0], chars.StructWS) {
		str = str[1:]
	}
	hasTrailingWS := false
	for len(str) > 0 && chars.ByteHasProperty(str[len(str)-1], chars.StructWS) {
		str = str[:len(str)-1]
		hasTrailingWS = true
	}
	if hasTrailingWS {
		str = append(str, ' ')
	}
	return str
}

func trimWhitespaceMarkupEnd(str []byte) []byte {
	return trimWhitespace(str)
}

func (_this *Reader) ReadSingleLineComment() []byte {
	_this.TokenBegin()
	_this.TokenReadUntilAndIncludingByte('\n')
	contents := _this.TokenGet()

	return trimWhitespace(contents)
}

func (_this *Reader) ReadMultilineComment() ([]byte, nextType) {
	_this.TokenBegin()
	lastByte := _this.TokenReadByteNoEOF()

	for {
		firstByte := lastByte
		lastByte = _this.TokenReadByteNoEOF()

		if firstByte == '*' && lastByte == '/' {
			_this.TokenStripLastBytes(2)
			contents := _this.TokenGet()
			return trimWhitespace(contents), nextIsCommentEnd
		}

		if firstByte == '/' && lastByte == '*' {
			_this.TokenStripLastBytes(2)
			contents := _this.TokenGet()
			return trimWhitespace(contents), nextIsCommentBegin
		}
	}
}

func (_this *Reader) ReadMarkupContent() ([]byte, nextType) {
	_this.TokenBegin()
	for {
		b := _this.TokenReadByteNoEOF()
		switch b {
		case '<':
			_this.TokenStripLastByte()
			return trimWhitespaceMarkupContent(_this.TokenGet()), nextIsMarkupBegin
		case '>':
			_this.TokenStripLastByte()
			return trimWhitespaceMarkupEnd(_this.TokenGet()), nextIsMarkupEnd
		case '/':
			switch _this.TokenReadByteAllowEOF() {
			case '*':
				_this.TokenStripLastBytes(2)
				return trimWhitespaceMarkupContent(_this.TokenGet()), nextIsCommentBegin
			case '/':
				_this.TokenStripLastBytes(2)
				return trimWhitespaceMarkupContent(_this.TokenGet()), nextIsSingleLineComment
			}
		case '\\':
			_this.TokenStripLastByte()
			_this.TokenReadEscape()

		}
	}
}

func (_this *Reader) DecodeHexDigit(b byte) byte {
	switch {
	case chars.ByteHasProperty(b, chars.DigitBase10):
		return byte(b - '0')
	case chars.ByteHasProperty(b, chars.LowerAF):
		return byte(b - 'a' + 10)
	case chars.ByteHasProperty(b, chars.UpperAF):
		return byte(b - 'A' + 10)
	default:
		_this.errorf("Unexpected char [%c] in UUID", b)
		return 0
	}
}

func (_this *Reader) decodeNextHexDigit() byte {
	return _this.DecodeHexDigit(_this.ReadByteNoEOF())
}

func (_this *Reader) decodeUUIDSection(byteOffset int, byteCount int) {
	for i := 0; i < byteCount; i++ {
		_this.token[byteOffset+i] = (_this.decodeNextHexDigit() << 4) | _this.decodeNextHexDigit()
	}
}

func (_this *Reader) decodeDash() {
	b := _this.ReadByteNoEOF()
	if b != '-' {
		_this.errorf("Unexpected char [%c] where [-] expected in UUID", b)
	}
}

func (_this *Reader) ReadUUIDWithTokenOffset(byteOffset int) []byte {
	_this.token = _this.token[:16]
	_this.decodeUUIDSection(byteOffset, 4-byteOffset)
	_this.decodeDash()
	_this.decodeUUIDSection(4, 2)
	_this.decodeDash()
	_this.decodeUUIDSection(6, 2)
	_this.decodeDash()
	_this.decodeUUIDSection(8, 2)
	_this.decodeDash()
	_this.decodeUUIDSection(10, 6)
	return _this.token
}

func (_this *Reader) ReadUUIDWithPreDecoded(decodedBytes ...byte) []byte {
	_this.token = _this.token[:16]
	copy(_this.token, decodedBytes)
	return _this.ReadUUIDWithTokenOffset(len(decodedBytes))
}

func (_this *Reader) ReadUUIDWithDecimalDecoded(alreadyDecodedAsDecimal uint64, decodedDigitCount int) []byte {
	_this.token = _this.token[:16]

	if decodedDigitCount > 0 {
		iLastDigit := decodedDigitCount / 2

		if decodedDigitCount&1 == 1 {
			_this.token[iLastDigit] = (byte(alreadyDecodedAsDecimal%10) << 4) | _this.decodeNextHexDigit()
			alreadyDecodedAsDecimal /= 10
			decodedDigitCount++
		}

		for i := iLastDigit - 1; i >= 0; i-- {
			v := byte(alreadyDecodedAsDecimal % 10)
			alreadyDecodedAsDecimal /= 10
			v |= byte(alreadyDecodedAsDecimal%10) << 4
			alreadyDecodedAsDecimal /= 10
			_this.token[i] = v
		}
	}

	return _this.ReadUUIDWithTokenOffset(decodedDigitCount / 2)
}

// ============================================================================

// Internal

type nextType int

const (
	nextIsCommentBegin nextType = iota
	nextIsCommentEnd
	nextIsSingleLineComment
	nextIsMarkupBegin
	nextIsMarkupEnd
)

var subsecondMagnitudes = []int{
	1000000000,
	100000000,
	10000000,
	1000000,
	100000,
	10000,
	1000,
	100,
	10,
	1,
}
