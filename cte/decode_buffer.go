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

	"github.com/kstenerud/go-concise-encoding/buffer"
	"github.com/kstenerud/go-concise-encoding/internal/chars"
	"github.com/kstenerud/go-concise-encoding/internal/common"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
)

type DecodeBuffer struct {
	buffer          buffer.StreamingReadBuffer
	readPos         int
	tokenStart      int
	subtokenStart   int
	lineCount       int
	colCount        int
	workBufferStart int
	workBufferPos   int
}

// Create a new CTE read buffer. The buffer will be empty until RefillIfNecessary() is
// called.
//
// readBufferSize determines the initial size of the buffer, and
// loWaterByteCount determines when RefillIfNecessary() refills the buffer from
// the reader.
func NewReadBuffer(reader io.Reader, readBufferSize int, loWaterByteCount int) *DecodeBuffer {
	_this := &DecodeBuffer{}
	_this.Init(reader, readBufferSize, loWaterByteCount)
	return _this
}

// Init the read buffer. You may call this again to re-initialize the buffer.
//
// readBufferSize determines the initial size of the buffer, and
// loWaterByteCount determines when RefillIfNecessary() refills the buffer from
// the reader.
func (_this *DecodeBuffer) Init(reader io.Reader, readBufferSize int, loWaterByteCount int) {
	_this.buffer.Init(reader, readBufferSize, loWaterByteCount)
	_this.readPos = 0
	_this.tokenStart = 0
	_this.subtokenStart = 0
	_this.workBufferStart = 0
	_this.workBufferPos = 0
	_this.lineCount = 0
	_this.colCount = 0
}

func (_this *DecodeBuffer) RefillIfNecessary() {
	offset := _this.buffer.RefillIfNecessary(_this.tokenStart, _this.readPos)
	if _this.tokenStart-offset < 0 {
		_this.Errorf("BUG: TokenStart was %v, tokenPos was %v", _this.tokenStart-offset, _this.readPos)
	}
	_this.readPos += offset
	_this.tokenStart += offset
	_this.subtokenStart += offset
	_this.workBufferStart += offset
	_this.workBufferPos += offset
}

func (_this *DecodeBuffer) IsEndOfDocument() bool {
	return _this.buffer.IsEOF() && !_this.hasUnreadByte()
}

// Bytes

func (_this *DecodeBuffer) PeekByteNoEOD() byte {
	_this.RefillIfNecessary()
	if _this.IsEndOfDocument() {
		_this.UnexpectedEOD()
	}
	return _this.buffer.ByteAtOffset(_this.readPos)
}

func (_this *DecodeBuffer) PeekByteAllowEOD() chars.ByteWithEOF {
	_this.RefillIfNecessary()
	if _this.IsEndOfDocument() {
		return chars.EndOfDocumentMarker
	}
	return chars.ByteWithEOF(_this.buffer.ByteAtOffset(_this.readPos))
}

func (_this *DecodeBuffer) ReadByte() byte {
	b := _this.PeekByteNoEOD()
	_this.AdvanceByte()
	return b
}

func (_this *DecodeBuffer) AdvanceByte() {
	val := lineColCounter[_this.buffer.Buffer[_this.readPos]]
	_this.colCount += val & 1
	_this.colCount &= ^(val >> 1)
	_this.lineCount += (val >> 1) & 1

	_this.readPos++
}

// Warning: Do not use this to unget a linefeed!
func (_this *DecodeBuffer) UngetByte() {
	_this.readPos--
	_this.colCount--
}

// Warning: Do not use this to unget a linefeed!
func (_this *DecodeBuffer) UngetBytes(count int) {
	_this.readPos -= count
	_this.colCount -= count
}

func (_this *DecodeBuffer) skipBytes(byteCount int) {
	_this.buffer.RequireBytes(_this.readPos, byteCount)
	_this.readPos += byteCount
	_this.colCount += len(string(_this.buffer.Buffer[_this.readPos : _this.readPos+byteCount]))
}

func (_this *DecodeBuffer) ReadUntilPropertyNoEOD(property chars.CharProperty) {
	for !chars.ByteHasProperty(_this.PeekByteNoEOD(), property) {
		_this.AdvanceByte()
	}
}

func (_this *DecodeBuffer) ReadUntilPropertyAllowEOD(property chars.CharProperty) {
	for !_this.PeekByteAllowEOD().HasProperty(property) {
		_this.AdvanceByte()
	}
}

func (_this *DecodeBuffer) ReadWhilePropertyNoEOD(property chars.CharProperty) {
	for chars.ByteHasProperty(_this.PeekByteNoEOD(), property) {
		_this.AdvanceByte()
	}
}

func (_this *DecodeBuffer) ReadWhilePropertyAllowEOD(property chars.CharProperty) {
	for _this.PeekByteAllowEOD().HasProperty(property) {
		_this.AdvanceByte()
	}
}

func (_this *DecodeBuffer) getCharBeginIndex(index int) int {
	i := index
	// UTF-8 continuation characters have the form 10xxxxxx
	for b := _this.buffer.Buffer[i]; b >= 0x80 && b <= 0xc0 && i >= 0; b = _this.buffer.Buffer[i] {
		i--
	}
	if i < 0 {
		i = 0
	}
	return i
}

func (_this *DecodeBuffer) getCharEndIndex(index int) int {
	i := index
	// UTF-8 continuation characters have the form 10xxxxxx
	for b := _this.buffer.Buffer[i]; b >= 0x80 && b <= 0xc0 && i < len(_this.buffer.Buffer); b = _this.buffer.Buffer[i] {
		i++
	}
	return i
}

func (_this *DecodeBuffer) hasUnreadByte() bool {
	return _this.buffer.HasByteAtOffset(_this.readPos)
}

func (_this *DecodeBuffer) readUntilByte(b byte) {
	for _this.PeekByteNoEOD() != b {
		_this.AdvanceByte()
	}
}

// Tokens

func (_this *DecodeBuffer) BeginToken() {
	_this.tokenStart = _this.readPos
}

func (_this *DecodeBuffer) GetToken() []byte {
	return _this.buffer.Buffer[_this.tokenStart:_this.readPos]
}

func (_this *DecodeBuffer) GetTokenWithEndOffset(offset int) []byte {
	return _this.buffer.Buffer[_this.tokenStart : _this.readPos+offset]
}

func (_this *DecodeBuffer) GetTokenLength() int {
	return _this.readPos - _this.tokenStart
}

func (_this *DecodeBuffer) GetTokenFirstByte() byte {
	return _this.buffer.Buffer[_this.tokenStart]
}

func (_this *DecodeBuffer) BeginSubtoken() {
	_this.subtokenStart = _this.readPos
}

func (_this *DecodeBuffer) GetSubtoken() []byte {
	return _this.buffer.Buffer[_this.subtokenStart:_this.readPos]
}

// Work Buffer

func (_this *DecodeBuffer) beginWorkBufferAtOffset(offset int) {
	_this.workBufferStart = _this.readPos + offset
	_this.workBufferPos = _this.workBufferStart
}

func (_this *DecodeBuffer) getTokenAndWorkBuffer() []byte {
	return _this.buffer.Buffer[_this.tokenStart:_this.workBufferPos]
}

func (_this *DecodeBuffer) writeWorkBufferByte(b byte) {
	_this.buffer.Buffer[_this.workBufferPos] = b
	_this.workBufferPos++
}

func (_this *DecodeBuffer) writeWorkBufferRune(r rune) {
	if r < utf8.RuneSelf {
		_this.writeWorkBufferByte(byte(r))
	} else {
		_this.workBufferPos += utf8.EncodeRune(_this.buffer.Buffer[_this.workBufferPos:], r)
	}
}

func (_this *DecodeBuffer) writeWorkBufferBytes(b []byte) {
	copy(_this.buffer.Buffer[_this.workBufferPos:], b)
	_this.workBufferPos += len(b)
}

func (_this *DecodeBuffer) endWorkBuffer() {
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
	index := _this.readPos
	if index >= len(_this.buffer.Buffer) {
		return "EOD"
	}

	charStart := _this.getCharBeginIndex(index)
	charEnd := _this.getCharEndIndex(index)
	if charEnd-charStart > 1 {
		return string(_this.buffer.Buffer[charStart:charEnd])
	}

	b := _this.buffer.Buffer[charStart]
	if b > ' ' && b <= '~' {
		return string(b)
	}
	if b == ' ' {
		return "SP"
	}
	return fmt.Sprintf("0x%02x", b)
}

// Decoders

func (_this *DecodeBuffer) SkipWhitespace() {
	_this.ReadWhilePropertyAllowEOD(chars.CharIsWhitespace)
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
	if _this.PeekByteAllowEOD() == '-' {
		sign = -sign
		_this.AdvanceByte()
	}
	v, count := _this.DecodeSmallBinaryUint()

	if v > 0x7fffffffffffffff && !(sign < 0 && v == 0x8000000000000000) {
		_this.Errorf("Integer value too big for element")
	}
	return int64(v) * sign, count
}

func (_this *DecodeBuffer) DecodeBinaryUint() (value uint64, bigValue *big.Int, digitCount int) {
	for {
		b := _this.PeekByteAllowEOD()
		nextDigitValue := uint64(0)
		switch {
		case b == charNumericWhitespace:
			_this.AdvanceByte()
			continue
		case b.HasProperty(chars.CharIsDigitBase2):
			nextDigitValue = uint64(b - '0')
		default:
			return
		}

		if value > maxPreShiftBinary {
			bigValue = new(big.Int).SetUint64(value)
			break
		}
		value = value<<1 + nextDigitValue
		digitCount++
		_this.AdvanceByte()
	}

	if bigValue == nil {
		return
	}

	for {
		b := _this.PeekByteAllowEOD()
		nextDigitValue := int64(0)
		switch {
		case b == charNumericWhitespace:
			_this.AdvanceByte()
			continue
		case b.HasProperty(chars.CharIsDigitBase2):
			nextDigitValue = int64(b - '0')
		default:
			return
		}

		bigValue = bigValue.Mul(bigValue, common.BigInt2)
		bigValue = bigValue.Add(bigValue, big.NewInt(nextDigitValue))
		digitCount++
		_this.AdvanceByte()
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
	if _this.PeekByteAllowEOD() == '-' {
		sign = -sign
		_this.AdvanceByte()
	}
	v, count := _this.DecodeSmallOctalUint()

	if v > 0x7fffffffffffffff && !(sign < 0 && v == 0x8000000000000000) {
		_this.Errorf("Integer value too big for element")
	}
	return int64(v) * sign, count
}

func (_this *DecodeBuffer) DecodeOctalUint() (value uint64, bigValue *big.Int, digitCount int) {
	for {
		b := _this.PeekByteAllowEOD()
		nextDigitValue := uint64(0)
		switch {
		case b == charNumericWhitespace:
			_this.AdvanceByte()
			continue
		case b.HasProperty(chars.CharIsDigitBase8):
			nextDigitValue = uint64(b - '0')
		default:
			return
		}

		if value > maxPreShiftOctal {
			bigValue = new(big.Int).SetUint64(value)
			break
		}
		value = value<<3 + nextDigitValue
		digitCount++
		_this.AdvanceByte()
	}

	if bigValue == nil {
		return
	}

	for {
		b := _this.PeekByteAllowEOD()
		nextDigitValue := int64(0)
		switch {
		case b == charNumericWhitespace:
			_this.AdvanceByte()
			continue
		case b.HasProperty(chars.CharIsDigitBase8):
			nextDigitValue = int64(b - '0')
		default:
			return
		}

		bigValue = bigValue.Mul(bigValue, common.BigInt8)
		bigValue = bigValue.Add(bigValue, big.NewInt(nextDigitValue))
		digitCount++
		_this.AdvanceByte()
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
			b := _this.PeekByteAllowEOD()
			if b == charNumericWhitespace {
				_this.AdvanceByte()
				continue
			}
			if !b.HasProperty(chars.CharIsDigitBase10) {
				return
			}
			bValue := uint64(b - '0')
			if value > maxPreShiftDecimal || (value == maxPreShiftDecimal && bValue > maxLastDigitDecimal) {
				bigStartValue = new(big.Int).SetUint64(value)
				break
			}
			value = value*10 + bValue
			digitCount++
			_this.AdvanceByte()
		}

		if bigStartValue == nil {
			return
		}
	}

	bigValue = bigStartValue
	for {
		b := _this.PeekByteAllowEOD()
		if b == charNumericWhitespace {
			_this.AdvanceByte()
			continue
		}
		if !b.HasProperty(chars.CharIsDigitBase10) {
			return
		}
		bigValue = bigValue.Mul(bigValue, common.BigInt10)
		bigValue = bigValue.Add(bigValue, big.NewInt(int64(b-'0')))
		digitCount++
		_this.AdvanceByte()
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
	if _this.PeekByteAllowEOD() == '-' {
		sign = -sign
		_this.AdvanceByte()
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
			b := _this.PeekByteAllowEOD()
			nextNybble := uint64(0)
			switch {
			case b == charNumericWhitespace:
				_this.AdvanceByte()
				continue
			case b.HasProperty(chars.CharIsDigitBase10):
				nextNybble = uint64(b - '0')
			case b.HasProperty(chars.CharIsLowerAF):
				nextNybble = uint64(b-'a') + 10
			case b.HasProperty(chars.CharIsUpperAF):
				nextNybble = uint64(b-'A') + 10
			default:
				return
			}

			if value > maxPreShiftHex {
				bigStartValue = new(big.Int).SetUint64(value)
				break
			}
			value = value<<4 + nextNybble
			digitCount++
			_this.AdvanceByte()
		}

		if bigStartValue == nil {
			return
		}
	}

	bigValue = bigStartValue
	for {
		b := _this.PeekByteAllowEOD()
		nextNybble := uint64(0)
		switch {
		case b == charNumericWhitespace:
			_this.AdvanceByte()
			continue
		case b.HasProperty(chars.CharIsDigitBase10):
			nextNybble = uint64(b - '0')
		case b.HasProperty(chars.CharIsLowerAF):
			nextNybble = uint64(b-'a') + 10
		case b.HasProperty(chars.CharIsUpperAF):
			nextNybble = uint64(b-'A') + 10
		default:
			return
		}

		bigValue = bigValue.Mul(bigValue, common.BigInt16)
		bigValue = bigValue.Add(bigValue, big.NewInt(int64(nextNybble)))
		digitCount++
		_this.AdvanceByte()
	}

	return
}

func (_this *DecodeBuffer) DecodeSmallUint() (value uint64, digitCount int) {
	var bigV *big.Int

	if _this.PeekByteAllowEOD() == '0' {
		_this.AdvanceByte()
		switch _this.PeekByteAllowEOD() {
		case 'b', 'B':
			_this.AdvanceByte()
			value, bigV, digitCount = _this.DecodeBinaryUint()
		case 'o', 'O':
			_this.AdvanceByte()
			value, bigV, digitCount = _this.DecodeOctalUint()
		case 'x', 'X':
			_this.AdvanceByte()
			value, bigV, digitCount = _this.DecodeHexUint(0, nil)
		default:
			value, bigV, digitCount = _this.DecodeDecimalUint(0, nil)
			digitCount++
		}
	} else {
		value, bigV, digitCount = _this.DecodeDecimalUint(0, nil)
	}

	if bigV != nil {
		_this.Errorf("Integer value too big for element")
	}
	return
}

func (_this *DecodeBuffer) DecodeSmallInt() (value int64, digitCount int) {
	sign := int64(1)
	if _this.PeekByteAllowEOD() == '-' {
		sign = -sign
		_this.AdvanceByte()
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
	if fractionalDigitCount == 0 {
		_this.UnexpectedChar("hex float fractional")
	}
	digitCount = coefficientDigitCount + fractionalDigitCount

	b := _this.PeekByteAllowEOD()
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
	case b.HasProperty(chars.CharIsWhitespace):
		return float64(u) * float64(sign), coefficientDigitCount
	default:
		_this.UnexpectedChar("hex float")
		return 0, 0
	}
}

func (_this *DecodeBuffer) DecodeSmallFloat() (value float64, digitCount int) {
	sign := int64(1)
	b := _this.PeekByteAllowEOD()
	if b == '-' {
		sign = -sign
		_this.AdvanceByte()
	} else if !b.HasProperty(chars.CharIsDigitBase10) {
		return
	}

	if _this.PeekByteAllowEOD() == '0' {
		_this.AdvanceByte()
		switch _this.PeekByteAllowEOD() {
		case 'x', 'X':
			_this.AdvanceByte()
			u, bigU, coefficientDigitCount := _this.DecodeHexUint(0, nil)
			if coefficientDigitCount == 0 {
				_this.UnexpectedChar("float")
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
			case b.HasProperty(chars.CharIsWhitespace):
				return float64(u) * float64(sign), coefficientDigitCount
			default:
				_this.UnexpectedChar("float")
				return 0, 0
			}
		}
		_this.UngetByte()
	}

	u, bigU, coefficientDigitCount := _this.DecodeDecimalUint(0, nil)
	if coefficientDigitCount == 0 {
		_this.UnexpectedChar("float")
	}
	if bigU != nil || u > maxFloat64Coefficient {
		_this.Errorf("Value too big for element")
	}

	b = _this.PeekByteAllowEOD()
	switch {
	case b == '.':
		_this.AdvanceByte()
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
		return float64(u) * float64(sign), coefficientDigitCount
	default:
		_this.UnexpectedChar("float")
		return 0, 0
	}

}

func (_this *DecodeBuffer) DecodeNamedValue() []byte {
	_this.BeginToken()
	_this.ReadUntilPropertyAllowEOD(chars.CharIsObjectEnd)
	namedValue := _this.GetToken()
	if len(namedValue) == 0 {
		_this.UnexpectedChar("name")
	}
	common.ASCIIBytesToLower(namedValue)
	return namedValue
}

func (_this *DecodeBuffer) decodeVerbatimString() []byte {
	_this.BeginSubtoken()
	_this.ReadUntilPropertyNoEOD(chars.CharIsWhitespace)
	sentinel := _this.GetSubtoken()
	sentinelOffset := _this.subtokenStart - _this.tokenStart
	sentinelLen := len(sentinel)
	wsByte := _this.ReadByte()
	if wsByte == '\r' {
		if _this.ReadByte() != '\n' {
			_this.UngetByte()
			_this.UnexpectedChar("verbatim sentinel")
		}
	}
	_this.BeginSubtoken()
	_this.skipBytes(sentinelLen - 1)

Outer:
	for {
		_this.ReadByte()
		sentinelStart := _this.tokenStart + sentinelOffset
		compareStart := _this.readPos - sentinelLen
		for i := sentinelLen - 1; i >= 0; i-- {
			if _this.buffer.Buffer[sentinelStart+i] != _this.buffer.Buffer[compareStart+i] {
				continue Outer
			}
		}
		subtoken := _this.GetSubtoken()
		return subtoken[:len(subtoken)-sentinelLen]
	}
}

func (_this *DecodeBuffer) decodeEscape() {
	escape := _this.ReadByte()
	switch escape {
	case 't':
		_this.writeWorkBufferByte('\t')
	case 'n':
		_this.writeWorkBufferByte('\n')
	case 'r':
		_this.writeWorkBufferByte('\r')
	case '"', '*', '/', '<', '>', '\\', '|':
		_this.writeWorkBufferByte(escape)
	case '_':
		// Non-breaking space
		_this.writeWorkBufferBytes([]byte{0xc0, 0xa0})
	case '-':
		// Soft hyphen
		_this.writeWorkBufferBytes([]byte{0xc0, 0xad})
	case '\r', '\n':
		// Continuation
		_this.ReadWhilePropertyNoEOD(chars.CharIsWhitespace)
	case '0':
		_this.writeWorkBufferByte(0)
	case '1', '2', '3', '4', '5', '6', '7', '8', '9':
		length := int(escape - '0')
		codepoint := rune(0)
		for i := 0; i < length; i++ {
			b := _this.PeekByteNoEOD()
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
			_this.AdvanceByte()
		}

		if codepoint < utf8.RuneSelf {
			_this.writeWorkBufferByte(uint8(codepoint))
		} else {
			_this.writeWorkBufferRune(codepoint)
		}
	case '.':
		_this.writeWorkBufferBytes(_this.decodeVerbatimString())
	default:
		_this.UnexpectedChar("escape sequence")
	}
}

func (_this *DecodeBuffer) DecodeQuotedString() []byte {
	_this.BeginToken()
outer:
	for {
		b := _this.ReadByte()
		switch b {
		case '"':
			return _this.GetTokenWithEndOffset(-1)
		case '\\':
			_this.beginWorkBufferAtOffset(-1)
			_this.decodeEscape()
			break outer
		}
	}

	for {
		b := _this.ReadByte()
		switch b {
		case '"':
			str := _this.getTokenAndWorkBuffer()
			_this.endWorkBuffer()
			return str
		case '\\':
			_this.decodeEscape()
		default:
			_this.writeWorkBufferByte(b)
		}
	}
}

func (_this *DecodeBuffer) DecodeUnquotedString() []byte {
	_this.BeginToken()
	_this.ReadUntilPropertyAllowEOD(chars.CharNeedsQuote)
	return _this.GetToken()
}

func (_this *DecodeBuffer) DecodeStringArray() []byte {
	_this.SkipWhitespace()
	_this.BeginToken()
outer:
	for {
		b := _this.ReadByte()
		switch b {
		case '|':
			str := _this.GetTokenWithEndOffset(-1)
			_this.endWorkBuffer()
			return str
		case '\\':
			_this.beginWorkBufferAtOffset(-1)
			_this.decodeEscape()
			break outer
		}
	}

	for {
		b := _this.ReadByte()
		switch b {
		case '|':
			return _this.getTokenAndWorkBuffer()
		case '\\':
			_this.decodeEscape()
		default:
			_this.writeWorkBufferByte(b)
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
	if _this.PeekByteNoEOD() != '-' {
		_this.UnexpectedChar("month")
	}
	_this.AdvanceByte()

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
	if _this.PeekByteAllowEOD() != '/' {
		t, err := compact_time.NewDate(int(year), int(month), int(day))
		if err != nil {
			_this.UnexpectedError(err, "date")
		}
		return t
	}

	_this.AdvanceByte()
	var hour uint64
	hour, _, digitCount = _this.DecodeDecimalUint(0, nil)
	if digitCount == 0 {
		_this.UnexpectedChar("hour")
	}
	if digitCount > 2 {
		_this.Errorf("Hour field is too long")
	}
	if _this.ReadByte() != ':' {
		_this.UngetByte()
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
	if _this.PeekByteNoEOD() != ':' {
		_this.UnexpectedChar("minute")
	}
	_this.AdvanceByte()
	var second uint64
	second, _, digitCount = _this.DecodeDecimalUint(0, nil)
	if digitCount > 2 {
		_this.Errorf("Second field is too long")
	}
	if second < 0 || second > 60 {
		_this.Errorf("Second %v is invalid", second)
	}
	var nanosecond int

	if _this.PeekByteAllowEOD() == '.' {
		_this.AdvanceByte()
		v, _, digitCount := _this.DecodeDecimalUint(0, nil)
		if digitCount == 0 {
			_this.UnexpectedChar("nanosecond")
		}
		if digitCount > 9 {
			_this.Errorf("Nanosecond field is too long")
		}
		nanosecond = int(v)
		nanosecond *= subsecondMagnitudes[digitCount]
	}

	b := _this.PeekByteAllowEOD()

	if b == '/' {
		_this.AdvanceByte()

		// TODO: Multiple levels of /
		if chars.ByteHasProperty(_this.PeekByteNoEOD(), chars.CharIsAZ) {
			_this.BeginToken()
			_this.ReadWhilePropertyAllowEOD(chars.CharIsAreaLocation)
			if _this.PeekByteAllowEOD() == '/' {
				_this.AdvanceByte()
				_this.ReadWhilePropertyAllowEOD(chars.CharIsAreaLocation)
			}

			areaLocation := string(_this.GetToken())
			t, err := compact_time.NewTime(hour, int(minute), int(second), nanosecond, areaLocation)
			if err != nil {
				_this.UnexpectedError(err, "time area/loc")
			}
			return t
		}

		lat, long := _this.DecodeLatLong()
		t, err := compact_time.NewTimeLatLong(hour, int(minute), int(second), nanosecond, lat, long)
		if err != nil {
			_this.UnexpectedError(err, "time lat/long")
		}
		return t
	}

	if b.HasProperty(chars.CharIsObjectEnd) {
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
	b := _this.PeekByteAllowEOD()
	if b == '.' {
		_this.AdvanceByte()
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
	}
	return int(whole*100 + fractional)
}

func (_this *DecodeBuffer) DecodeLatLong() (latitudeHundredths, longitudeHundredths int) {
	latitudeHundredths = _this.DecodeLatLongPortion("latitude")

	if _this.PeekByteNoEOD() != '/' {
		_this.UnexpectedChar("latitude/longitude")
	}
	_this.AdvanceByte()

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

func (_this *DecodeBuffer) DecodeSingleLineComment() []byte {
	_this.BeginToken()
	_this.readUntilByte('\n')
	contents := _this.GetToken()
	_this.AdvanceByte()

	return trimWhitespace(contents)
}

func (_this *DecodeBuffer) DecodeMultilineComment() ([]byte, nextType) {
	_this.BeginToken()
	lastByte := _this.ReadByte()

	for {
		firstByte := lastByte
		lastByte = _this.ReadByte()

		if firstByte == '*' && lastByte == '/' {
			contents := _this.GetTokenWithEndOffset(-2)
			return trimWhitespace(contents), nextIsCommentEnd
		}

		if firstByte == '/' && lastByte == '*' {
			contents := _this.GetTokenWithEndOffset(-2)
			return trimWhitespace(contents), nextIsCommentBegin
		}
	}
}

func (_this *DecodeBuffer) DecodeMarkupContent() ([]byte, nextType) {
	isInitialWS := true
	wsCount := 0
	isCommentInitiator := false

	completeStringPortion := func() {
		if wsCount > 0 {
			if !isInitialWS {
				_this.writeWorkBufferByte(' ')
			}
			wsCount = 0
		}
		isInitialWS = false
	}

	completeContentsPortion := func() {
		wsCount = 0
		isInitialWS = false
	}

	addSlashIfNeeded := func() {
		if isCommentInitiator {
			_this.writeWorkBufferByte('/')
			isCommentInitiator = false
		}
	}

	_this.BeginToken()
	_this.beginWorkBufferAtOffset(0)
	for {
		currentByte := _this.ReadByte()
		switch currentByte {
		case '\r', '\n', '\t', ' ':
			wsCount++
		case '\\':
			_this.decodeEscape()
		case '<':
			completeStringPortion()
			addSlashIfNeeded()
			str := _this.getTokenAndWorkBuffer()
			_this.endWorkBuffer()
			return str, nextIsMarkupBegin
		case '>':
			completeContentsPortion()
			addSlashIfNeeded()
			str := _this.getTokenAndWorkBuffer()
			_this.endWorkBuffer()
			return str, nextIsMarkupEnd
		case '/':
			if isCommentInitiator {
				completeStringPortion()
				str := _this.getTokenAndWorkBuffer()
				_this.endWorkBuffer()
				return str, nextIsSingleLineComment
			} else {
				isCommentInitiator = true
			}
		case '*':
			if isCommentInitiator {
				completeStringPortion()
				str := _this.getTokenAndWorkBuffer()
				_this.endWorkBuffer()
				return str, nextIsCommentBegin
			} else {
				_this.writeWorkBufferByte(currentByte)
			}
		default:
			completeStringPortion()
			addSlashIfNeeded()
			_this.writeWorkBufferByte(currentByte)
		}
	}
}

// Decode a marker ID. asString will be empty if the result is an integer.
func (_this *DecodeBuffer) DecodeMarkerID() (asString []byte, asUint uint64) {
	_this.BeginToken()

	b := _this.PeekByteNoEOD()
	switch {
	case chars.ByteHasProperty(b, chars.CharIsDigitBase10):
		asUint, _ = _this.DecodeSmallUint()
	case !chars.ByteHasProperty(b, chars.CharNeedsQuoteFirst):
		for !_this.PeekByteAllowEOD().HasProperty(chars.CharNeedsQuote) {
			_this.AdvanceByte()
		}
		asString = _this.GetToken()
	default:
		_this.Errorf("Missing marker ID")
	}
	return
}

// ============================================================================

// Internal

func chooseLowWater(bufferSize int) int {
	lowWater := bufferSize / 50
	if lowWater < 32 {
		lowWater = 32
	}
	return lowWater
}

type nextType int

const (
	nextIsCommentBegin nextType = iota
	nextIsCommentEnd
	nextIsSingleLineComment
	nextIsMarkupBegin
	nextIsMarkupEnd
)

var lineColCounter [256]int

func init() {
	for i := 0; i < 0x80; i++ {
		lineColCounter[i] = 1
	}
	for i := 0xc1; i < 0xff; i++ {
		lineColCounter[i] = 1
	}
	lineColCounter['\n'] = -1
}

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
