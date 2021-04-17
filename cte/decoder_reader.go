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

type Reader struct {
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
	_this.readPos = 0
	_this.lineCount = 0
	_this.colCount = 0
}

// Bytes

func (_this *Reader) advanceLineColCount(b chars.ByteWithEOF) {
	// TODO: Handle more line terminators:
	// LF:    Line Feed, U+000A
	// VT:    Vertical Tab, U+000B
	// FF:    Form Feed, U+000C
	// CR:    Carriage Return, U+000D
	// NEL:   Next Line, U+0085                      : C2 85
	// LS:    Line Separator, U+2028                 : E2 80 A8
	// PS:    Paragraph Separator, U+2029            : E2 80 A9

	_this.lastLineCount = _this.lineCount
	_this.lastColCount = _this.colCount

	switch b {
	case '\n':
		_this.lineCount++
		_this.colCount = 1
	case chars.EOFMarker:
		// Do nothing
	default:
		_this.colCount++
	}

	_this.readPos++
}

func (_this *Reader) retreatLineColCount() {
	_this.readPos--
	_this.lineCount = _this.lastLineCount
	_this.colCount = _this.lastColCount
}

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
	_this.retreatLineColCount()
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
	_this.advanceLineColCount(_this.lastByte)
	return _this.lastByte
}

func (_this *Reader) PeekByteNoEOF() byte {
	b := _this.PeekByteAllowEOF()
	if b == chars.EOFMarker {
		_this.UnexpectedEOF()
	}
	return byte(b)
}

func (_this *Reader) ReadByteNoEOF() byte {
	b := _this.ReadByteAllowEOF()
	if b == chars.EOFMarker {
		_this.UnexpectedEOF()
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

func (_this *Reader) TokenUnreadByte() {
	_this.UnreadByte()
	_this.token = _this.token[:len(_this.token)-1]
}

func (_this *Reader) TokenReadUntilByte(untilByte byte) {
	for {
		b := _this.ReadByteNoEOF()
		if b == untilByte {
			_this.UnreadByte()
			break
		}
		_this.TokenAppendByte(b)
	}
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

func (_this *Reader) TokenReadWhilePropertyNoEOF(property chars.Properties) {
	for {
		b := _this.ReadByteNoEOF()
		if !chars.ByteHasProperty(b, property) {
			_this.UnreadByte()
			break
		}
		_this.TokenAppendByte(b)
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

func (_this *Reader) TokenReadUntilByteNoEOF(untilByte byte) {
	for {
		b := _this.ReadByteNoEOF()
		if b == untilByte {
			_this.UnreadByte()
			break
		}
		_this.TokenAppendByte(b)
	}
}

func (_this *Reader) TokenReadUntilOneOfBytesNoEOF(untilBytes ...byte) {
	for {
		b := _this.ReadByteNoEOF()
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

func (_this *Reader) AssertAtObjectEnd(decoding string) {
	if !_this.PeekByteAllowEOF().HasProperty(chars.ObjectEnd) {
		_this.UnexpectedChar(decoding)
	}
}

func (_this *Reader) Errorf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	panic(fmt.Errorf("offset %v (line %v, col %v): %v", _this.readPos, _this.lineCount+1, _this.colCount+1, msg))
}

func (_this *Reader) UnexpectedEOF() {
	_this.Errorf("unexpected end of document")
}

func (_this *Reader) UnexpectedError(err error, decoding string) {
	_this.Errorf("unexpected error [%v] while decoding %v", err, decoding)
}

func (_this *Reader) UnexpectedChar(decoding string) {
	_this.Errorf("unexpected [%v] while decoding %v", _this.DescribeCurrentChar(), decoding)
}

func (_this *Reader) DescribeCurrentChar() string {
	b := _this.PeekByteAllowEOF()
	switch {
	case b == chars.EOFMarker:
		return "EOF"
	case b == ' ':
		return "SP"
	case b > ' ' && b <= '~':
		return fmt.Sprintf("%c", b)
	default:
		return fmt.Sprintf("0x%02x", b)
	}
}

// Decoders

func (_this *Reader) SkipWhitespace() {
	_this.SkipWhileProperty(chars.StructWS)
}

func (_this *Reader) ReadToken() []byte {
	_this.TokenBegin()
	_this.TokenReadUntilPropertyAllowEOF(chars.ObjectEnd)
	return _this.TokenGet()
}

const maxPreShiftBinary = uint64(0x7fffffffffffffff)

func (_this *Reader) ReadSmallBinaryUint() (value uint64, digitCount int) {
	v, vBig, count := _this.ReadBinaryUint()
	if vBig != nil {
		_this.Errorf("Value cannot be > 64 bits")
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
		_this.Errorf("Integer value too big for element")
	}
	return int64(v) * sign, count
}

func (_this *Reader) ReadBinaryUint() (value uint64, bigValue *big.Int, digitCount int) {
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
		_this.Errorf("Value cannot be > 64 bits")
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
		_this.Errorf("Integer value too big for element")
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

const maxPreShiftDecimal = uint64(1844674407370955161)
const maxLastDigitDecimal = 5

// startValue is only used if bigStartValue is nil
// bigValue will be nil unless the value was too big for a uint64
func (_this *Reader) ReadDecimalUint(startValue uint64, bigStartValue *big.Int) (value uint64, bigValue *big.Int, digitCount int) {
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
		_this.Errorf("Value cannot be > 64 bits")
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
		_this.Errorf("Integer value too big for element")
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
		_this.Errorf("Integer value too big for element")
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
		_this.Errorf("Integer value too big for element")
	}
	return int64(v) * sign, count
}

func (_this *Reader) ReadDecimalFloat(sign int64, coefficient uint64, bigCoefficient *big.Int, coefficientDigitCount int) (value compact_float.DFloat, bigValue *apd.Decimal, digitCount int) {
	exponent := int32(0)
	fractionalDigitCount := 0
	coefficient, bigCoefficient, fractionalDigitCount = _this.ReadDecimalUint(coefficient, bigCoefficient)
	if fractionalDigitCount == 0 {
		_this.UnexpectedChar("float fractional")
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

func (_this *Reader) ReadHexFloat(sign int64, coefficient uint64, bigCoefficient *big.Int, coefficientDigitCount int) (value float64, bigValue *big.Float, digitCount int) {

	exponent := 0
	fractionalDigitCount := 0
	coefficient, bigCoefficient, fractionalDigitCount = _this.ReadHexUint(coefficient, bigCoefficient)
	b := _this.PeekByteAllowEOF()
	if fractionalDigitCount == 0 && b != 'p' && b != 'P' {
		_this.UnexpectedChar("hex float fractional")
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
		_this.UnexpectedChar("hex float coefficient")
	}
	if bigU != nil || u > maxFloat64Coefficient {
		_this.Errorf("Value too big for element")
	}
	b = _this.PeekByteAllowEOF()
	switch {
	case b == '.':
		_this.AdvanceByte()
		f, bigF, digitCount := _this.ReadHexFloat(sign, u, nil, coefficientDigitCount)
		if bigF != nil {
			_this.Errorf("Value too big for element")
		}
		return f, digitCount
	case b == 'p' || b == 'P':
		f, bigF, digitCount := _this.ReadHexFloat(sign, u, nil, coefficientDigitCount)
		if bigF != nil {
			_this.Errorf("Value too big for element")
		}
		return f, digitCount
	case b.HasProperty(chars.ObjectEnd):
		return float64(u) * float64(sign), coefficientDigitCount
	default:
		_this.UnexpectedChar("hex float")
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
				_this.UnexpectedChar("float")
			}
			if bigU != nil || u > maxFloat64Coefficient {
				_this.Errorf("Value too big for element")
			}

			b = _this.ReadByteAllowEOF()
			switch {
			case b == '.':
				f, bigF, digitCount := _this.ReadHexFloat(sign, u, nil, coefficientDigitCount)
				if bigF != nil {
					_this.Errorf("Value too big for element")
				}
				return f, digitCount
			case b.HasProperty(chars.StructWS):
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

	u, bigU, coefficientDigitCount := _this.ReadDecimalUint(0, nil)
	if initialZero {
		coefficientDigitCount++
	}
	if coefficientDigitCount == 0 {
		_this.UnexpectedChar("float")
	}
	if bigU != nil || u > maxFloat64Coefficient {
		_this.Errorf("Value too big for element")
	}

	b = _this.ReadByteAllowEOF()
	switch {
	case b == '.':
		f, bigF, digitCount := _this.ReadDecimalFloat(sign, u, nil, coefficientDigitCount)
		if bigF != nil {
			_this.Errorf("Value too big for element")
		}
		normalizedExponent := int(f.Exponent) + digitCount - 1
		if normalizedExponent < minFloat64DecimalExponent || normalizedExponent > maxFloat64DecimalExponent {
			_this.Errorf("Value too big for element")
		}

		return f.Float(), digitCount
	case b.HasProperty(chars.StructWS):
		_this.UnreadByte()
		return float64(u) * float64(sign), coefficientDigitCount
	default:
		_this.UnexpectedChar("float")
		return 0, 0
	}

}

func (_this *Reader) ReadNamedValue() []byte {
	_this.TokenBegin()
	_this.TokenReadUntilPropertyAllowEOF(chars.ObjectEnd)
	namedValue := _this.TokenGet()
	if len(namedValue) == 0 {
		_this.UnexpectedChar("name")
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
				_this.UnexpectedChar("unicode escape")
			}
		}

		_this.TokenAppendRune(codepoint)
	case '.':
		_this.TokenReadVerbatimSequence()
	default:
		_this.UnexpectedChar("escape sequence")
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
	_this.TokenReadUntilPropertyAllowEOF(chars.ObjectEnd)
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

var maxDayByMonth = []int{0, 31, 29, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

func (_this *Reader) ReadDate(year int64) compact_time.Time {
	month, _, digitCount := _this.ReadDecimalUint(0, nil)
	if digitCount > 2 {
		_this.Errorf("Month field is too long")
	}
	if month < 1 || month > 12 {
		_this.Errorf("Month %v is invalid", month)
	}
	if _this.ReadByteNoEOF() != '-' {
		_this.UnexpectedChar("month")
	}

	var day uint64
	day, _, digitCount = _this.ReadDecimalUint(0, nil)
	if digitCount == 0 {
		_this.UnexpectedChar("day")
	}
	if digitCount > 2 {
		_this.Errorf("Day field is too long")
	}
	if day < 1 || int(day) > maxDayByMonth[month] {
		_this.Errorf("Day %v is invalid", day)
	}
	if _this.ReadByteAllowEOF() != '/' {
		_this.UnreadByte()
		t, err := compact_time.NewDate(int(year), int(month), int(day))
		if err != nil {
			_this.UnexpectedError(err, "date")
		}
		return t
	}

	var hour uint64
	hour, _, digitCount = _this.ReadDecimalUint(0, nil)
	if digitCount == 0 {
		_this.UnexpectedChar("hour")
	}
	if digitCount > 2 {
		_this.Errorf("Hour field is too long")
	}
	if _this.ReadByteNoEOF() != ':' {
		_this.UnreadByte()
		_this.UnexpectedChar("hour")
	}
	t := _this.ReadTime(int(hour))
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

func (_this *Reader) ReadTime(hour int) compact_time.Time {
	if hour < 0 || hour > 23 {
		_this.Errorf("Hour %v is invalid", hour)
	}
	minute, _, digitCount := _this.ReadDecimalUint(0, nil)
	if digitCount > 2 {
		_this.Errorf("Minute field is too long")
	}
	if minute < 0 || minute > 59 {
		_this.Errorf("Minute %v is invalid", minute)
	}
	if _this.ReadByteNoEOF() != ':' {
		_this.UnreadByte()
		_this.UnexpectedChar("minute")
	}

	var second uint64
	second, _, digitCount = _this.ReadDecimalUint(0, nil)
	if digitCount > 2 {
		_this.Errorf("Second field is too long")
	}
	if second < 0 || second > 60 {
		_this.Errorf("Second %v is invalid", second)
	}
	var nanosecond int

	b := _this.ReadByteAllowEOF()

	if b == '.' {
		v, _, digitCount := _this.ReadDecimalUint(0, nil)
		if digitCount == 0 {
			_this.UnexpectedChar("nanosecond")
		}
		if digitCount > 9 {
			_this.Errorf("Nanosecond field is too long")
		}
		nanosecond = int(v)
		nanosecond *= subsecondMagnitudes[digitCount]
		b = _this.ReadByteAllowEOF()
	}

	if b == '/' {
		next := _this.PeekByteNoEOF()
		if chars.ByteHasProperty(next, chars.DigitBase10) || next == '-' {
			lat, long := _this.ReadLatLong()
			t, err := compact_time.NewTimeLatLong(hour, int(minute), int(second), nanosecond, lat, long)
			if err != nil {
				_this.UnexpectedError(err, "time lat/long")
			}
			return t
		}

		_this.TokenBegin()
		_this.TokenReadWhilePropertyAllowEOF(chars.AreaLocation)
		areaLocation := string(_this.TokenGet())
		t, err := compact_time.NewTime(hour, int(minute), int(second), nanosecond, areaLocation)
		if err != nil {
			_this.UnexpectedError(err, "time area/loc")
		}
		return t
	}

	if b.HasProperty(chars.ObjectEnd) {
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

func (_this *Reader) ReadLatLongPortion(name string) (value int) {
	sign := 1
	if _this.PeekByteNoEOF() == '-' {
		_this.AdvanceByte()
		sign = -1
	}
	whole, _, digitCount := _this.ReadDecimalUint(0, nil)
	switch digitCount {
	case 1, 2, 3:
		// 1-3 digits are allowed
	case 0:
		_this.UnexpectedChar(name)
	default:
		_this.Errorf("Too many digits decoding %v", name)
	}

	var fractional uint64
	if _this.ReadByteAllowEOF() == '.' {
		fractional, _, digitCount = _this.ReadDecimalUint(0, nil)
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

	return sign * int(whole*100+fractional)
}

func (_this *Reader) ReadLatLong() (latitudeHundredths, longitudeHundredths int) {
	latitudeHundredths = _this.ReadLatLongPortion("latitude")

	if _this.ReadByteNoEOF() != '/' {
		_this.UnexpectedChar("latitude/longitude")
	}

	longitudeHundredths = _this.ReadLatLongPortion("longitude")

	return
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

func (_this *Reader) ExtractUUID(data []byte) []byte {
	if len(data) != 36 ||
		data[8] != '-' ||
		data[13] != '-' ||
		data[18] != '-' ||
		data[23] != '-' {
		_this.Errorf("Malformed UUID or unknown named value: [%s]", string(data))
	}

	decodeHex := func(b byte) byte {
		switch {
		case chars.ByteHasProperty(b, chars.DigitBase10):
			return byte(b - '0')
		case chars.ByteHasProperty(b, chars.LowerAF):
			return byte(b - 'a' + 10)
		case chars.ByteHasProperty(b, chars.UpperAF):
			return byte(b - 'A' + 10)
		default:
			_this.Errorf("Unexpected char [%c] in UUID [%s]", b, string(data))
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

	decodeSection(data[:8], data)
	decodeSection(data[9:13], data[4:])
	decodeSection(data[14:18], data[6:])
	decodeSection(data[19:23], data[8:])
	decodeSection(data[24:36], data[10:])

	return data[:16]
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
