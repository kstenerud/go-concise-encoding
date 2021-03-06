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
	"fmt"
	"io"
	"math"
	"math/big"
	"unicode/utf8"

	"github.com/kstenerud/go-concise-encoding/internal/chars"
	"github.com/kstenerud/go-concise-encoding/internal/common"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
)

type DecodeBuffer struct {
	reader    io.Reader
	lastByte  chars.ByteWithEOF
	hasUnread bool
	isEOF     bool
	byteBuff  [1]byte

	readPos       int
	lineCount     int
	colCount      int
	lastLineCount int
	lastColCount  int

	token            []byte
	verbatimSentinel []byte
}

func NewReadBuffer(reader io.Reader) *DecodeBuffer {
	_this := &DecodeBuffer{}
	_this.Init(reader)
	return _this
}

// Init the read buffer. You may call this again to re-initialize the buffer.
func (_this *DecodeBuffer) Init(reader io.Reader) {
	_this.reader = reader
	if cap(_this.token) == 0 {
		_this.token = make([]byte, 0, 16)
	}
	_this.hasUnread = false
	_this.isEOF = false
	_this.readPos = 0
	_this.lineCount = 0
	_this.colCount = 0
}

// Bytes

func (_this *DecodeBuffer) advanceLineColCount(b chars.ByteWithEOF) {
	_this.lastLineCount = _this.lineCount
	_this.lastColCount = _this.colCount

	switch b {
	case '\n':
		_this.lineCount++
		_this.colCount = 1
	case chars.EndOfDocumentMarker:
		// Do nothing
	default:
		_this.colCount++
	}

	_this.readPos++
}

func (_this *DecodeBuffer) retreatLineColCount() {
	_this.readPos--
	_this.lineCount = _this.lastLineCount
	_this.colCount = _this.lastColCount
}

func (_this *DecodeBuffer) readNext() {
	// Can't create local [1]byte here because it mallocs? WTF???
	if _, err := _this.reader.Read(_this.byteBuff[:]); err == nil {
		_this.lastByte = chars.ByteWithEOF(_this.byteBuff[0])
	} else {
		if err == io.EOF {
			_this.isEOF = true
			_this.lastByte = chars.EndOfDocumentMarker
		} else {
			panic(err)
		}
	}
}

func (_this *DecodeBuffer) UnreadByte() {
	if _this.hasUnread {
		panic("Cannot unread twice")
	}
	_this.hasUnread = true
	_this.retreatLineColCount()
}

func (_this *DecodeBuffer) PeekByteAllowEOD() chars.ByteWithEOF {
	if !_this.hasUnread {
		_this.readNext()
		_this.hasUnread = true
	}
	return _this.lastByte
}

func (_this *DecodeBuffer) ReadByteAllowEOD() chars.ByteWithEOF {
	if !_this.hasUnread {
		_this.readNext()
	}

	_this.hasUnread = false
	_this.advanceLineColCount(_this.lastByte)
	return _this.lastByte
}

func (_this *DecodeBuffer) PeekByteNoEOD() byte {
	b := _this.PeekByteAllowEOD()
	if b == chars.EndOfDocumentMarker {
		_this.UnexpectedEOD()
	}
	return byte(b)
}

func (_this *DecodeBuffer) ReadByteNoEOD() byte {
	b := _this.ReadByteAllowEOD()
	if b == chars.EndOfDocumentMarker {
		_this.UnexpectedEOD()
	}
	return byte(b)
}

func (_this *DecodeBuffer) AdvanceByte() {
	_this.ReadByteNoEOD()
}

func (_this *DecodeBuffer) SkipWhileProperty(property chars.CharProperty) {
	for _this.ReadByteAllowEOD().HasProperty(property) {
	}
	_this.UnreadByte()
}

// Tokens

func (_this *DecodeBuffer) TokenBegin() {
	_this.token = _this.token[:0]
}

func (_this *DecodeBuffer) TokenStripLastByte() {
	_this.token = _this.token[:len(_this.token)-1]
}

func (_this *DecodeBuffer) TokenStripLastBytes(count int) {
	_this.token = _this.token[:len(_this.token)-count]
}

func (_this *DecodeBuffer) TokenAppendByte(b byte) {
	_this.token = append(_this.token, b)
}

func (_this *DecodeBuffer) TokenAppendBytes(b []byte) {
	_this.token = append(_this.token, b...)
}

func (_this *DecodeBuffer) TokenAppendRune(r rune) {
	if r < utf8.RuneSelf {
		_this.TokenAppendByte(byte(r))
	} else {
		pos := len(_this.token)
		_this.TokenAppendBytes([]byte{0, 0, 0, 0, 0})
		length := utf8.EncodeRune(_this.token[pos:], r)
		_this.token = _this.token[:pos+length]
	}
}

func (_this *DecodeBuffer) TokenGet() []byte {
	return _this.token
}

func (_this *DecodeBuffer) TokenReadByteNoEOD() byte {
	b := _this.ReadByteNoEOD()
	_this.TokenAppendByte(b)
	return b
}

func (_this *DecodeBuffer) TokenReadByteAllowEOD() chars.ByteWithEOF {
	b := _this.ReadByteAllowEOD()
	if b != chars.EndOfDocumentMarker {
		_this.TokenAppendByte(byte(b))
	}
	return b
}

func (_this *DecodeBuffer) TokenUnreadByte() {
	_this.UnreadByte()
	_this.token = _this.token[:len(_this.token)-1]
}

func (_this *DecodeBuffer) TokenReadUntilByte(untilByte byte) {
	for {
		b := _this.ReadByteNoEOD()
		if b == untilByte {
			_this.UnreadByte()
			break
		}
		_this.TokenAppendByte(b)
	}
}

func (_this *DecodeBuffer) TokenReadUntilAndIncludingByte(untilByte byte) {
	for {
		b := _this.ReadByteNoEOD()
		if b == untilByte {
			break
		}
		_this.TokenAppendByte(b)
	}
}

func (_this *DecodeBuffer) TokenReadUntilPropertyNoEOD(property chars.CharProperty) {
	for {
		b := _this.ReadByteNoEOD()
		if chars.ByteHasProperty(b, property) {
			_this.UnreadByte()
			break
		}
		_this.TokenAppendByte(b)
	}
}

func (_this *DecodeBuffer) TokenReadUntilPropertyAllowEOD(property chars.CharProperty) {
	for {
		b := _this.ReadByteAllowEOD()
		if b.HasProperty(property) {
			_this.UnreadByte()
			break
		}
		_this.TokenAppendByte(byte(b))
	}
}

func (_this *DecodeBuffer) TokenReadWhilePropertyNoEOD(property chars.CharProperty) {
	for {
		b := _this.ReadByteNoEOD()
		if !chars.ByteHasProperty(b, property) {
			_this.UnreadByte()
			break
		}
		_this.TokenAppendByte(b)
	}
}

func (_this *DecodeBuffer) TokenReadWhilePropertyAllowEOD(property chars.CharProperty) {
	for {
		b := _this.ReadByteAllowEOD()
		if !b.HasProperty(property) {
			_this.UnreadByte()
			break
		}
		_this.TokenAppendByte(byte(b))
	}
}

func (_this *DecodeBuffer) TokenReadUntilByteNoEOD(untilByte byte) {
	for {
		b := _this.ReadByteNoEOD()
		if b == untilByte {
			_this.UnreadByte()
			break
		}
		_this.TokenAppendByte(b)
	}
}

func (_this *DecodeBuffer) TokenReadUntilOneOfBytesNoEOD(untilBytes ...byte) {
	for {
		b := _this.ReadByteNoEOD()
		for _, check := range untilBytes {
			if b == check {
				_this.UnreadByte()
				return
			}
		}
		_this.TokenAppendByte(b)
	}
}

// Errors

func (_this *DecodeBuffer) AssertAtObjectEnd(decoding string) {
	if !_this.PeekByteAllowEOD().HasProperty(chars.CharIsObjectEnd) {
		_this.UnexpectedChar(decoding)
	}
}

func (_this *DecodeBuffer) Errorf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	panic(fmt.Errorf("offset %v (line %v, col %v): %v", _this.readPos, _this.lineCount+1, _this.colCount+1, msg))
}

func (_this *DecodeBuffer) UnexpectedEOD() {
	_this.Errorf("unexpected end of document")
}

func (_this *DecodeBuffer) UnexpectedError(err error, decoding string) {
	_this.Errorf("unexpected error [%v] while decoding %v", err, decoding)
}

func (_this *DecodeBuffer) UnexpectedChar(decoding string) {
	_this.Errorf("unexpected [%v] while decoding %v", _this.DescribeCurrentChar(), decoding)
}

func (_this *DecodeBuffer) DescribeCurrentChar() string {
	b := _this.PeekByteAllowEOD()
	switch {
	case b == chars.EndOfDocumentMarker:
		return "EOD"
	case b == ' ':
		return "SP"
	case b > ' ' && b <= '~':
		return fmt.Sprintf("%c", b)
	default:
		return fmt.Sprintf("0x%02x", b)
	}
}

// Decoders

func (_this *DecodeBuffer) SkipWhitespace() {
	_this.SkipWhileProperty(chars.CharIsWhitespace)
}

const maxPreShiftBinary = uint64(0x7fffffffffffffff)

func (_this *DecodeBuffer) DecodeSmallBinaryUint() (value uint64, digitCount int) {
	v, vBig, count := _this.DecodeBinaryUint()
	if vBig != nil {
		_this.Errorf("Value cannot be > 64 bits")
	}
	return v, count
}

func (_this *DecodeBuffer) DecodeSmallBinaryInt() (value int64, digitCount int) {
	sign := int64(1)
	if _this.ReadByteAllowEOD() == '-' {
		sign = -sign
	} else {
		_this.UnreadByte()
	}
	v, count := _this.DecodeSmallBinaryUint()

	if v > 0x7fffffffffffffff && !(sign < 0 && v == 0x8000000000000000) {
		_this.Errorf("Integer value too big for element")
	}
	return int64(v) * sign, count
}

func (_this *DecodeBuffer) DecodeBinaryUint() (value uint64, bigValue *big.Int, digitCount int) {
	for {
		b := _this.ReadByteAllowEOD()
		switch {
		case b == charNumericWhitespace:
			continue
		case b.HasProperty(chars.CharIsDigitBase2):
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
		b := _this.ReadByteAllowEOD()
		switch {
		case b == charNumericWhitespace:
			continue
		case b.HasProperty(chars.CharIsDigitBase2):
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

func (_this *DecodeBuffer) DecodeSmallOctalUint() (value uint64, digitCount int) {
	v, vBig, count := _this.DecodeOctalUint()
	if vBig != nil {
		_this.Errorf("Value cannot be > 64 bits")
	}
	return v, count
}

func (_this *DecodeBuffer) DecodeSmallOctalInt() (value int64, digitCount int) {
	sign := int64(1)
	if _this.ReadByteAllowEOD() == '-' {
		sign = -sign
	} else {
		_this.UnreadByte()
	}
	v, count := _this.DecodeSmallOctalUint()

	if v > 0x7fffffffffffffff && !(sign < 0 && v == 0x8000000000000000) {
		_this.Errorf("Integer value too big for element")
	}
	return int64(v) * sign, count
}

func (_this *DecodeBuffer) DecodeOctalUint() (value uint64, bigValue *big.Int, digitCount int) {
	for {
		b := _this.ReadByteAllowEOD()
		switch {
		case b == charNumericWhitespace:
			continue
		case b.HasProperty(chars.CharIsDigitBase8):
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
		b := _this.ReadByteAllowEOD()
		switch {
		case b == charNumericWhitespace:
			continue
		case b.HasProperty(chars.CharIsDigitBase8):
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

const maxPreShiftDecimal = uint64(1844674407370955161)
const maxLastDigitDecimal = 5

// startValue is only used if bigStartValue is nil
// bigValue will be nil unless the value was too big for a uint64
func (_this *DecodeBuffer) DecodeDecimalUint(startValue uint64, bigStartValue *big.Int) (value uint64, bigValue *big.Int, digitCount int) {
	if bigStartValue == nil {
		value = startValue
		for {
			b := _this.ReadByteAllowEOD()
			switch {
			case b == charNumericWhitespace:
				continue
			case b.HasProperty(chars.CharIsDigitBase10):
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
		b := _this.ReadByteAllowEOD()
		switch {
		case b == charNumericWhitespace:
			continue
		case b.HasProperty(chars.CharIsDigitBase10):
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

func (_this *DecodeBuffer) DecodeSmallHexUint() (value uint64, digitCount int) {
	v, vBig, count := _this.DecodeHexUint(0, nil)
	if vBig != nil {
		_this.Errorf("Value cannot be > 64 bits")
	}
	return v, count
}

func (_this *DecodeBuffer) DecodeSmallHexInt() (value int64, digitCount int) {
	sign := int64(1)
	if _this.ReadByteAllowEOD() == '-' {
		sign = -sign
	} else {
		_this.UnreadByte()
	}
	v, count := _this.DecodeSmallHexUint()

	if v > 0x7fffffffffffffff && !(sign < 0 && v == 0x8000000000000000) {
		_this.Errorf("Integer value too big for element")
	}
	return int64(v) * sign, count
}

func (_this *DecodeBuffer) DecodeHexUint(startValue uint64, bigStartValue *big.Int) (value uint64, bigValue *big.Int, digitCount int) {
	if bigStartValue == nil {
		value = startValue
		for {
			b := _this.ReadByteAllowEOD()
			nextNybble := chars.ByteWithEOF(0)
			switch {
			case b == charNumericWhitespace:
				continue
			case b.HasProperty(chars.CharIsDigitBase10):
				nextNybble = b - '0'
			case b.HasProperty(chars.CharIsLowerAF):
				nextNybble = b - 'a' + 10
			case b.HasProperty(chars.CharIsUpperAF):
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
		b := _this.ReadByteAllowEOD()
		nextNybble := chars.ByteWithEOF(0)
		switch {
		case b == charNumericWhitespace:
			continue
		case b.HasProperty(chars.CharIsDigitBase10):
			nextNybble = b - '0'
		case b.HasProperty(chars.CharIsLowerAF):
			nextNybble = b - 'a' + 10
		case b.HasProperty(chars.CharIsUpperAF):
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

func (_this *DecodeBuffer) DecodeSmallUint() (value uint64, digitCount int) {
	var bigV *big.Int

	if _this.ReadByteAllowEOD() == '0' {
		switch _this.ReadByteAllowEOD() {
		case 'b', 'B':
			value, bigV, digitCount = _this.DecodeBinaryUint()
		case 'o', 'O':
			value, bigV, digitCount = _this.DecodeOctalUint()
		case 'x', 'X':
			value, bigV, digitCount = _this.DecodeHexUint(0, nil)
		default:
			_this.UnreadByte()
			value, bigV, digitCount = _this.DecodeDecimalUint(0, nil)
			digitCount++
		}
	} else {
		_this.UnreadByte()
		value, bigV, digitCount = _this.DecodeDecimalUint(0, nil)
	}

	if bigV != nil {
		_this.Errorf("Integer value too big for element")
	}
	return
}

func (_this *DecodeBuffer) DecodeSmallInt() (value int64, digitCount int) {
	sign := int64(1)
	if _this.ReadByteAllowEOD() == '-' {
		sign = -sign
	} else {
		_this.UnreadByte()
	}

	v, count := _this.DecodeSmallUint()

	if v > 0x7fffffffffffffff && !(sign < 0 && v == 0x8000000000000000) {
		_this.Errorf("Integer value too big for element")
	}
	return int64(v) * sign, count
}

func (_this *DecodeBuffer) DecodeDecimalFloat(sign int64, coefficient uint64, bigCoefficient *big.Int, coefficientDigitCount int) (value compact_float.DFloat, bigValue *apd.Decimal, digitCount int) {
	exponent := int32(0)
	fractionalDigitCount := 0
	coefficient, bigCoefficient, fractionalDigitCount = _this.DecodeDecimalUint(coefficient, bigCoefficient)
	if fractionalDigitCount == 0 {
		_this.UnexpectedChar("float fractional")
	}
	digitCount = coefficientDigitCount + fractionalDigitCount

	b := _this.PeekByteAllowEOD()
	if b == 'e' || b == 'E' {
		_this.AdvanceByte()
		exponentSign := int32(1)
		switch _this.PeekByteNoEOD() {
		case '+':
			_this.AdvanceByte()
		case '-':
			exponentSign = -1
			_this.AdvanceByte()
		}
		exp, bigExp, expDigitCount := _this.DecodeDecimalUint(0, nil)
		if expDigitCount == 0 {
			_this.UnexpectedChar("float exponent")
		}
		if bigExp != nil || exp > 0x7fffffff {
			_this.Errorf("Exponent too big")
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

func (_this *DecodeBuffer) DecodeHexFloat(sign int64, coefficient uint64, bigCoefficient *big.Int, coefficientDigitCount int) (value float64, bigValue *big.Float, digitCount int) {

	exponent := 0
	fractionalDigitCount := 0
	coefficient, bigCoefficient, fractionalDigitCount = _this.DecodeHexUint(coefficient, bigCoefficient)
	b := _this.PeekByteAllowEOD()
	if fractionalDigitCount == 0 && b != 'p' && b != 'P' {
		_this.UnexpectedChar("hex float fractional")
	}
	digitCount = coefficientDigitCount + fractionalDigitCount

	if b == 'p' || b == 'P' {
		_this.AdvanceByte()
		exponentSign := 1
		switch _this.PeekByteNoEOD() {
		case '+':
			_this.AdvanceByte()
		case '-':
			exponentSign = -1
			_this.AdvanceByte()
		}
		exp, bigExp, expDigitCount := _this.DecodeDecimalUint(0, nil)
		if expDigitCount == 0 {
			_this.UnexpectedChar("hex float exponent")
		}
		if bigExp != nil {
			_this.Errorf("Exponent too big")
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

func (_this *DecodeBuffer) DecodeSmallHexFloat() (value float64, digitCount int) {
	sign := int64(1)
	b := _this.PeekByteAllowEOD()
	if b == '-' {
		sign = -sign
		_this.AdvanceByte()
	} else if !b.HasProperty(chars.CharIsDigitBase10 | chars.CharIsLowerAF | chars.CharIsUpperAF) {
		return
	}

	u, bigU, coefficientDigitCount := _this.DecodeHexUint(0, nil)
	if coefficientDigitCount == 0 {
		_this.UnexpectedChar("hex float coefficient")
	}
	if bigU != nil || u > maxFloat64Coefficient {
		_this.Errorf("Value too big for element")
	}
	b = _this.PeekByteAllowEOD()
	switch {
	case b == '.':
		_this.AdvanceByte()
		f, bigF, digitCount := _this.DecodeHexFloat(sign, u, nil, coefficientDigitCount)
		if bigF != nil {
			_this.Errorf("Value too big for element")
		}
		return f, digitCount
	case b == 'p' || b == 'P':
		f, bigF, digitCount := _this.DecodeHexFloat(sign, u, nil, coefficientDigitCount)
		if bigF != nil {
			_this.Errorf("Value too big for element")
		}
		return f, digitCount
	case b.HasProperty(chars.CharIsObjectEnd):
		return float64(u) * float64(sign), coefficientDigitCount
	default:
		_this.UnexpectedChar("hex float")
		return 0, 0
	}
}

func (_this *DecodeBuffer) DecodeSmallFloat() (value float64, digitCount int) {
	sign := int64(1)
	b := _this.ReadByteAllowEOD()
	if b == '-' {
		sign = -sign
	} else if !b.HasProperty(chars.CharIsDigitBase10) {
		_this.UnreadByte()
		return
	} else {
		_this.UnreadByte()
	}

	initialZero := false
	if _this.ReadByteAllowEOD() == '0' {
		initialZero = true
		b = _this.ReadByteAllowEOD()
		if b == 'x' || b == 'X' {
			u, bigU, coefficientDigitCount := _this.DecodeHexUint(0, nil)
			if coefficientDigitCount == 0 {
				_this.UnexpectedChar("float")
			}
			if bigU != nil || u > maxFloat64Coefficient {
				_this.Errorf("Value too big for element")
			}

			b = _this.ReadByteAllowEOD()
			switch {
			case b == '.':
				f, bigF, digitCount := _this.DecodeHexFloat(sign, u, nil, coefficientDigitCount)
				if bigF != nil {
					_this.Errorf("Value too big for element")
				}
				return f, digitCount
			case b.HasProperty(chars.CharIsWhitespace):
				_this.UnreadByte()
				return float64(u) * float64(sign), coefficientDigitCount
			default:
				_this.UnexpectedChar("float")
				return 0, 0
			}
		} else {
			_this.UnreadByte()
		}
	} else {
		_this.UnreadByte()
	}

	u, bigU, coefficientDigitCount := _this.DecodeDecimalUint(0, nil)
	if initialZero {
		coefficientDigitCount++
	}
	if coefficientDigitCount == 0 {
		_this.UnexpectedChar("float")
	}
	if bigU != nil || u > maxFloat64Coefficient {
		_this.Errorf("Value too big for element")
	}

	b = _this.ReadByteAllowEOD()
	switch {
	case b == '.':
		f, bigF, digitCount := _this.DecodeDecimalFloat(sign, u, nil, coefficientDigitCount)
		if bigF != nil {
			_this.Errorf("Value too big for element")
		}
		normalizedExponent := int(f.Exponent) + digitCount - 1
		if normalizedExponent < minFloat64DecimalExponent || normalizedExponent > maxFloat64DecimalExponent {
			_this.Errorf("Value too big for element")
		}

		return f.Float(), digitCount
	case b.HasProperty(chars.CharIsWhitespace):
		_this.UnreadByte()
		return float64(u) * float64(sign), coefficientDigitCount
	default:
		_this.UnexpectedChar("float")
		return 0, 0
	}

}

func (_this *DecodeBuffer) DecodeNamedValue() []byte {
	_this.TokenBegin()
	_this.TokenReadUntilPropertyAllowEOD(chars.CharIsObjectEnd)
	namedValue := _this.TokenGet()
	if len(namedValue) == 0 {
		_this.UnexpectedChar("name")
	}
	common.ASCIIBytesToLower(namedValue)
	return namedValue
}

func (_this *DecodeBuffer) TokenDecodeVerbatimSequence() {
	_this.verbatimSentinel = _this.verbatimSentinel[:0]
	for {
		b := _this.ReadByteNoEOD()
		if chars.ByteHasProperty(b, chars.CharIsWhitespace) {
			if b == '\r' {
				if _this.ReadByteNoEOD() != '\n' {
					_this.UnexpectedChar("verbatim sentinel")
				}
			}
			break
		}
		_this.verbatimSentinel = append(_this.verbatimSentinel, b)
	}

	sentinelLength := len(_this.verbatimSentinel)

Outer:
	for {
		_this.TokenReadByteNoEOD()
		for i := 1; i <= sentinelLength; i++ {
			if _this.token[len(_this.token)-i] != _this.verbatimSentinel[sentinelLength-i] {
				continue Outer
			}
		}

		_this.TokenStripLastBytes(sentinelLength)
		return
	}
}

func (_this *DecodeBuffer) TokenDecodeEscape() {
	escape := _this.ReadByteNoEOD()
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
			b := _this.ReadByteNoEOD()
			switch {
			case chars.ByteHasProperty(b, chars.CharIsDigitBase10):
				codepoint = (codepoint << 4) | (rune(b) - '0')
			case chars.ByteHasProperty(b, chars.CharIsLowerAF):
				codepoint = (codepoint << 4) | (rune(b) - 'a' + 10)
			case chars.ByteHasProperty(b, chars.CharIsUpperAF):
				codepoint = (codepoint << 4) | (rune(b) - 'A' + 10)
			default:
				_this.UnexpectedChar("unicode escape")
			}
		}

		_this.TokenAppendRune(codepoint)
	case '.':
		_this.TokenDecodeVerbatimSequence()
	default:
		_this.UnexpectedChar("escape sequence")
	}
}

func (_this *DecodeBuffer) DecodeQuotedString() []byte {
	_this.TokenBegin()
	for {
		b := _this.TokenReadByteNoEOD()
		switch b {
		case '"':
			_this.TokenStripLastByte()
			return _this.TokenGet()
		case '\\':
			_this.TokenStripLastByte()
			_this.TokenDecodeEscape()
		}
	}
}

func (_this *DecodeBuffer) DecodeUnquotedString() []byte {
	_this.TokenBegin()
	_this.TokenReadUntilPropertyAllowEOD(chars.CharNeedsQuote)
	return _this.TokenGet()
}

func (_this *DecodeBuffer) DecodeStringArray() []byte {
	_this.TokenBegin()
	for {
		b := _this.TokenReadByteNoEOD()
		switch b {
		case '|':
			_this.TokenStripLastByte()
			return _this.TokenGet()
		case '\\':
			_this.TokenStripLastByte()
			_this.TokenDecodeEscape()
		}
	}
}

var maxDayByMonth = []int{0, 31, 29, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

func (_this *DecodeBuffer) DecodeDate(year int64) compact_time.Time {
	month, _, digitCount := _this.DecodeDecimalUint(0, nil)
	if digitCount > 2 {
		_this.Errorf("Month field is too long")
	}
	if month < 1 || month > 12 {
		_this.Errorf("Month %v is invalid", month)
	}
	if _this.ReadByteNoEOD() != '-' {
		_this.UnexpectedChar("month")
	}

	var day uint64
	day, _, digitCount = _this.DecodeDecimalUint(0, nil)
	if digitCount == 0 {
		_this.UnexpectedChar("day")
	}
	if digitCount > 2 {
		_this.Errorf("Day field is too long")
	}
	if day < 1 || int(day) > maxDayByMonth[month] {
		_this.Errorf("Day %v is invalid", day)
	}
	if _this.ReadByteAllowEOD() != '/' {
		_this.UnreadByte()
		t, err := compact_time.NewDate(int(year), int(month), int(day))
		if err != nil {
			_this.UnexpectedError(err, "date")
		}
		return t
	}

	var hour uint64
	hour, _, digitCount = _this.DecodeDecimalUint(0, nil)
	if digitCount == 0 {
		_this.UnexpectedChar("hour")
	}
	if digitCount > 2 {
		_this.Errorf("Hour field is too long")
	}
	if _this.ReadByteNoEOD() != ':' {
		_this.UnreadByte()
		_this.UnexpectedChar("hour")
	}
	t := _this.DecodeTime(int(hour))
	if t.TimezoneType == compact_time.TypeLatitudeLongitude {
		ts, err := compact_time.NewTimestampLatLong(int(year), int(month), int(day),
			int(t.Hour), int(t.Minute), int(t.Second), int(t.Nanosecond),
			int(t.LatitudeHundredths), int(t.LongitudeHundredths))
		if err != nil {
			_this.UnexpectedError(err, "timestamp lat/long")
		}
		return ts
	}
	ts, err := compact_time.NewTimestamp(int(year), int(month), int(day),
		int(t.Hour), int(t.Minute), int(t.Second), int(t.Nanosecond),
		t.ShortAreaLocation)
	if err != nil {
		_this.UnexpectedError(err, "timestamp area/loc")
	}
	return ts
}

func (_this *DecodeBuffer) DecodeTime(hour int) compact_time.Time {
	if hour < 0 || hour > 23 {
		_this.Errorf("Hour %v is invalid", hour)
	}
	minute, _, digitCount := _this.DecodeDecimalUint(0, nil)
	if digitCount > 2 {
		_this.Errorf("Minute field is too long")
	}
	if minute < 0 || minute > 59 {
		_this.Errorf("Minute %v is invalid", minute)
	}
	if _this.ReadByteNoEOD() != ':' {
		_this.UnreadByte()
		_this.UnexpectedChar("minute")
	}

	var second uint64
	second, _, digitCount = _this.DecodeDecimalUint(0, nil)
	if digitCount > 2 {
		_this.Errorf("Second field is too long")
	}
	if second < 0 || second > 60 {
		_this.Errorf("Second %v is invalid", second)
	}
	var nanosecond int

	b := _this.ReadByteAllowEOD()

	if b == '.' {
		v, _, digitCount := _this.DecodeDecimalUint(0, nil)
		if digitCount == 0 {
			_this.UnexpectedChar("nanosecond")
		}
		if digitCount > 9 {
			_this.Errorf("Nanosecond field is too long")
		}
		nanosecond = int(v)
		nanosecond *= subsecondMagnitudes[digitCount]
		b = _this.ReadByteAllowEOD()
	}

	if b == '/' {
		if chars.ByteHasProperty(_this.PeekByteNoEOD(), chars.CharIsDigitBase10) {
			lat, long := _this.DecodeLatLong()
			t, err := compact_time.NewTimeLatLong(hour, int(minute), int(second), nanosecond, lat, long)
			if err != nil {
				_this.UnexpectedError(err, "time lat/long")
			}
			return t
		}

		_this.TokenBegin()
		_this.TokenReadWhilePropertyAllowEOD(chars.CharIsAreaLocation)
		areaLocation := string(_this.TokenGet())
		t, err := compact_time.NewTime(hour, int(minute), int(second), nanosecond, areaLocation)
		if err != nil {
			_this.UnexpectedError(err, "time area/loc")
		}
		return t
	}

	if b.HasProperty(chars.CharIsObjectEnd) {
		_this.UnreadByte()
		t, err := compact_time.NewTime(hour, int(minute), int(second), nanosecond, "")
		if err != nil {
			_this.UnexpectedError(err, "time zero")
		}
		return t
	}

	_this.UnexpectedChar("time")
	return compact_time.Time{}
}

func (_this *DecodeBuffer) DecodeLatLongPortion(name string) (value int) {
	whole, _, digitCount := _this.DecodeDecimalUint(0, nil)
	switch digitCount {
	case 1, 2, 3:
		// 1-3 digits are allowed
	case 0:
		_this.UnexpectedChar(name)
	default:
		_this.Errorf("Too many digits decoding %v", name)
	}

	var fractional uint64
	if _this.ReadByteAllowEOD() == '.' {
		fractional, _, digitCount = _this.DecodeDecimalUint(0, nil)
		switch digitCount {
		case 1:
			// 1 digit: multiply fractional by 10
			fractional *= 10
		case 2:
			// 2 digits are allowed
		case 0:
			_this.UnexpectedChar(name)
		default:
			_this.Errorf("Too many digits decoding %v", name)
		}
	} else {
		_this.UnreadByte()
	}

	return int(whole*100 + fractional)
}

func (_this *DecodeBuffer) DecodeLatLong() (latitudeHundredths, longitudeHundredths int) {
	latitudeHundredths = _this.DecodeLatLongPortion("latitude")

	if _this.ReadByteNoEOD() != '/' {
		_this.UnexpectedChar("latitude/longitude")
	}

	longitudeHundredths = _this.DecodeLatLongPortion("longitude")

	return
}

func trimWhitespace(str []byte) []byte {
	for len(str) > 0 && chars.ByteHasProperty(str[0], chars.CharIsWhitespace) {
		str = str[1:]
	}
	for len(str) > 0 && chars.ByteHasProperty(str[len(str)-1], chars.CharIsWhitespace) {
		str = str[:len(str)-1]
	}
	return str
}

func trimWhitespaceMarkupContent(str []byte) []byte {
	for len(str) > 0 && chars.ByteHasProperty(str[0], chars.CharIsWhitespace) {
		str = str[1:]
	}
	hasTrailingWS := false
	for len(str) > 0 && chars.ByteHasProperty(str[len(str)-1], chars.CharIsWhitespace) {
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

func (_this *DecodeBuffer) DecodeSingleLineComment() []byte {
	_this.TokenBegin()
	_this.TokenReadUntilAndIncludingByte('\n')
	contents := _this.TokenGet()

	return trimWhitespace(contents)
}

func (_this *DecodeBuffer) DecodeMultilineComment() ([]byte, nextType) {
	_this.TokenBegin()
	lastByte := _this.TokenReadByteNoEOD()

	for {
		firstByte := lastByte
		lastByte = _this.TokenReadByteNoEOD()

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

func (_this *DecodeBuffer) DecodeMarkupContent() ([]byte, nextType) {
	_this.TokenBegin()
	for {
		b := _this.TokenReadByteNoEOD()
		switch b {
		case '<':
			_this.TokenStripLastByte()
			return trimWhitespaceMarkupContent(_this.TokenGet()), nextIsMarkupBegin
		case '>':
			_this.TokenStripLastByte()
			return trimWhitespaceMarkupEnd(_this.TokenGet()), nextIsMarkupEnd
		case '/':
			switch _this.TokenReadByteAllowEOD() {
			case '*':
				_this.TokenStripLastBytes(2)
				return trimWhitespaceMarkupContent(_this.TokenGet()), nextIsCommentBegin
			case '/':
				_this.TokenStripLastBytes(2)
				return trimWhitespaceMarkupContent(_this.TokenGet()), nextIsSingleLineComment
			}
		case '\\':
			_this.TokenStripLastByte()
			_this.TokenDecodeEscape()

		}
	}
}

// Decode a marker ID. asString will be empty if the result is an integer.
func (_this *DecodeBuffer) DecodeMarkerID() (asString []byte, asUint uint64) {

	b := _this.PeekByteNoEOD()
	switch {
	case chars.ByteHasProperty(b, chars.CharIsDigitBase10):
		asUint, _ = _this.DecodeSmallUint()
	case !chars.ByteHasProperty(b, chars.CharNeedsQuoteFirst):
		_this.TokenBegin()
		_this.TokenReadUntilPropertyAllowEOD(chars.CharNeedsQuote)
		asString = _this.TokenGet()
	default:
		_this.Errorf("Missing marker ID")
	}
	return
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
