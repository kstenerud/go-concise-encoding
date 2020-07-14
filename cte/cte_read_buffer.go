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
	"fmt"
	"io"
	"math"
	"math/big"

	"github.com/kstenerud/go-concise-encoding/buffer"

	"github.com/kstenerud/go-concise-encoding/internal/common"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
)

type CTEReadBuffer struct {
	buffer        buffer.StreamingReadBuffer
	tokenStart    int
	subtokenStart int
	tokenPos      int
	tempBuffer    bytes.Buffer
	lineCount     int
	colCount      int
}

// Create a new CTE read buffer. The buffer will be empty until RefillIfNecessary() is
// called.
//
// readBufferSize determines the initial size of the buffer, and
// loWaterByteCount determines when RefillIfNecessary() refills the buffer from
// the reader.
func NewReadBuffer(reader io.Reader, readBufferSize int, loWaterByteCount int) *CTEReadBuffer {
	_this := &CTEReadBuffer{}
	_this.Init(reader, readBufferSize, loWaterByteCount)
	return _this
}

// Init the read buffer. The buffer will be empty until RefillIfNecessary() is
// called.
//
// readBufferSize determines the initial size of the buffer, and
// loWaterByteCount determines when RefillIfNecessary() refills the buffer from
// the reader.
func (_this *CTEReadBuffer) Init(reader io.Reader, readBufferSize int, loWaterByteCount int) {
	_this.buffer.Init(reader, readBufferSize, loWaterByteCount)
}

// Bytes

func (_this *CTEReadBuffer) PeekByteNoEOD() byte {
	_this.RefillIfNecessary()
	if _this.IsEndOfDocument() {
		_this.UnexpectedEOD()
	}
	return _this.buffer.ByteAtOffset(_this.tokenPos)
}

func (_this *CTEReadBuffer) PeekByteAllowEOD() cteByte {
	_this.RefillIfNecessary()
	if _this.IsEndOfDocument() {
		return cteByteEndOfDocument
	}
	return cteByte(_this.buffer.ByteAtOffset(_this.tokenPos))
}

func (_this *CTEReadBuffer) ReadByte() byte {
	b := _this.PeekByteNoEOD()
	_this.AdvanceByte()
	return b
}

func (_this *CTEReadBuffer) SkipBytes(byteCount int) {
	_this.buffer.RequireBytes(_this.tokenPos, _this.tokenPos+byteCount)
	_this.tokenPos += byteCount
}

func (_this *CTEReadBuffer) AdvanceByte() {
	_this.tokenPos++
}

func (_this *CTEReadBuffer) RefillIfNecessary() {
	offset := _this.buffer.RefillIfNecessary(_this.tokenStart, _this.tokenPos)
	if _this.tokenStart-offset < 0 {
		panic(fmt.Errorf("TokenStart was %v, tokenPos was %v", _this.tokenStart-offset, _this.tokenPos))
	}
	_this.tokenPos += offset
	_this.tokenStart += offset
	_this.subtokenStart += offset
}

func (_this *CTEReadBuffer) ReadUntilPropertyNoEOD(property cteByteProprty) {
	for !hasProperty(_this.PeekByteNoEOD(), property) {
		_this.AdvanceByte()
	}
}

func (_this *CTEReadBuffer) ReadWhilePropertyAllowEOD(property cteByteProprty) {
	for _this.PeekByteAllowEOD().HasProperty(property) {
		_this.AdvanceByte()
	}
}

func (_this *CTEReadBuffer) GetCharBeginIndex(index int) int {
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

func (_this *CTEReadBuffer) GetCharEndIndex(index int) int {
	i := index
	// UTF-8 continuation characters have the form 10xxxxxx
	for b := _this.buffer.Buffer[i]; b >= 0x80 && b <= 0xc0 && i < len(_this.buffer.Buffer); b = _this.buffer.Buffer[i] {
		i++
	}
	return i
}

func (_this *CTEReadBuffer) UngetByte() {
	_this.tokenPos--
}

func (_this *CTEReadBuffer) UngetBytes(byteCount int) {
	_this.tokenPos -= byteCount
}

func (_this *CTEReadBuffer) UngetAll() {
	_this.tokenPos = _this.tokenStart
}

func (_this *CTEReadBuffer) IsEndOfDocument() bool {
	return _this.buffer.IsEOF() && !_this.hasUnreadByte()
}

func (_this *CTEReadBuffer) hasUnreadByte() bool {
	return _this.buffer.HasByteAtOffset(_this.tokenPos)
}

func (_this *CTEReadBuffer) readUntilByte(b byte) {
	for _this.PeekByteNoEOD() != b {
		_this.AdvanceByte()
	}
}

// Tokens

func (_this *CTEReadBuffer) GetToken() []byte {
	return _this.buffer.Buffer[_this.tokenStart:_this.tokenPos]
}

func (_this *CTEReadBuffer) GetTokenLength() int {
	return _this.tokenPos - _this.tokenStart
}

func (_this *CTEReadBuffer) GetTokenFirstByte() byte {
	return _this.buffer.Buffer[_this.tokenStart]
}

func (_this *CTEReadBuffer) EndToken() {
	_this.lineCount, _this.colCount = _this.countLinesColsToIndex(_this.tokenPos)
	_this.tokenStart = _this.tokenPos
}

func (_this *CTEReadBuffer) BeginSubtoken() {
	_this.subtokenStart = _this.tokenPos
}

func (_this *CTEReadBuffer) GetSubtoken() []byte {
	return _this.buffer.Buffer[_this.subtokenStart:_this.tokenPos]
}

// Errors

func (_this *CTEReadBuffer) AssertAtObjectEnd(decoding string) {
	if !_this.PeekByteAllowEOD().HasProperty(ctePropertyObjectEnd) {
		_this.UnexpectedChar(decoding)
	}
}

func (_this *CTEReadBuffer) UnexpectedEOD() {
	_this.Errorf("Unexpected end of document")
}

func (_this *CTEReadBuffer) ErrorAt(index int, format string, args ...interface{}) {
	lineCount, colCount := _this.countLinesColsToIndex(index)

	msg := fmt.Sprintf(format, args...)
	panic(fmt.Errorf("Offset %v (line %v, col %v): %v", index, lineCount+1, colCount+1, msg))
	_this.Errorf(format, args...)
}

func (_this *CTEReadBuffer) Errorf(format string, args ...interface{}) {
	_this.ErrorAt(_this.tokenPos, format, args...)
}

func (_this *CTEReadBuffer) UnexpectedCharAt(index int, decoding string) {
	_this.ErrorAt(index, "Unexpected [%v] while decoding %v", _this.DescribeCharAt(index), decoding)
}

func (_this *CTEReadBuffer) UnexpectedChar(decoding string) {
	_this.UnexpectedCharAt(_this.tokenPos, decoding)
}

func (_this *CTEReadBuffer) DescribeCharAt(index int) string {
	if index >= len(_this.buffer.Buffer) {
		return "EOD"
	}

	charStart := _this.GetCharBeginIndex(index)
	charEnd := _this.GetCharEndIndex(index)
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

func (_this *CTEReadBuffer) DescribeCurrentChar() string {
	return _this.DescribeCharAt(_this.tokenPos)
}

func (_this *CTEReadBuffer) countLinesColsToIndex(index int) (lineCount, colCount int) {
	lineCount = _this.lineCount
	colCount = _this.colCount
	for i := _this.tokenStart; i < index; i++ {
		b := _this.buffer.Buffer[i]
		if b < 0x80 || b > 0xc0 {
			colCount++
		}
		if b == '\n' {
			lineCount++
			colCount = 0
		}
	}

	return
}

// Decoders

func (_this *CTEReadBuffer) SkipWhitespace() {
	_this.ReadWhilePropertyAllowEOD(ctePropertyWhitespace)
}

const maxPreShiftBinary = uint64(0x7fffffffffffffff)

func (_this *CTEReadBuffer) DecodeBinaryInteger() (value uint64) {
	for {
		b := _this.PeekByteAllowEOD()
		if !b.HasProperty(ctePropertyBinaryDigit) {
			return
		}
		if value > maxPreShiftBinary {
			// TODO: Support BigInt?
			_this.Errorf("Overflow reading binary integer")
		}
		value = value<<1 + uint64(b-'0')
		_this.AdvanceByte()
	}
}

const maxPreShiftOctal = uint64(0x1fffffffffffffff)

func (_this *CTEReadBuffer) DecodeOctalInteger() (value uint64) {
	for {
		b := _this.PeekByteAllowEOD()
		if !b.HasProperty(ctePropertyOctalDigit) {
			return
		}
		if value > maxPreShiftOctal {
			// TODO: Support BigInt?
			_this.Errorf("Overflow reading octal integer")
		}
		value = value<<3 + uint64(b-'0')
		_this.AdvanceByte()
	}
}

const maxPreShiftDecimal = uint64(0x1999999999999999)

// startValue is only used if bigStartValue is nil
// bigValue will be nil unless the value was too big for a uint64
func (_this *CTEReadBuffer) DecodeDecimalInteger(startValue uint64, bigStartValue *big.Int) (value uint64, bigValue *big.Int, digitCount int) {
	if bigStartValue == nil {
		value = startValue
		for {
			b := _this.PeekByteAllowEOD()
			if !b.HasProperty(cteProperty09) {
				return
			}
			if value > maxPreShiftDecimal {
				bigStartValue = new(big.Int).SetUint64(value)
				break
			}
			value = value*10 + uint64(b-'0')
			digitCount++
			_this.AdvanceByte()
		}
	}

	if bigStartValue != nil {
		bigValue = bigStartValue
		for {
			b := _this.PeekByteAllowEOD()
			if !b.HasProperty(cteProperty09) {
				return
			}
			bigValue = bigValue.Mul(bigValue, common.BigInt10)
			bigValue = bigValue.Add(bigValue, big.NewInt(int64(b-'0')))
			digitCount++
			_this.AdvanceByte()
		}
	}

	return
}

const maxPreShiftHex = uint64(0x0fffffffffffffff)

func (_this *CTEReadBuffer) DecodeHexInteger(startValue uint64) (value uint64, digitCount int) {
	value = startValue
	for {
		b := _this.PeekByteAllowEOD()
		nextNybble := uint64(0)
		switch {
		case b.HasProperty(cteProperty09):
			nextNybble = uint64(b - '0')
		case b.HasProperty(ctePropertyLowercaseAF):
			nextNybble = uint64(b-'a') + 10
		case b.HasProperty(ctePropertyUppercaseAF):
			nextNybble = uint64(b-'A') + 10
		default:
			return
		}

		if value > maxPreShiftHex {
			// TODO: Support BigInt?
			_this.Errorf("Overflow reading hex integer")
		}
		value = value<<4 + nextNybble
		digitCount++
		_this.AdvanceByte()
	}
}

func (_this *CTEReadBuffer) DecodeQuotedStringWithEscapes() string {
	_this.Errorf("TODO: CTEDecoder: Escape sequences")
	return ""
}

// TODO: Return []byte instead
func (_this *CTEReadBuffer) DecodeQuotedString() string {
	_this.BeginSubtoken()
	for {
		b := _this.PeekByteNoEOD()
		switch b {
		case '"':
			token := _this.GetSubtoken()
			_this.AdvanceByte()
			return string(token)
		case '\\':
			return _this.DecodeQuotedStringWithEscapes()
		}
		_this.AdvanceByte()
	}
}

func (_this *CTEReadBuffer) DecodeVerbatimString() []byte {
	_this.BeginSubtoken()
	_this.ReadUntilPropertyNoEOD(ctePropertyWhitespace)
	sentinel := _this.GetSubtoken()
	sentinelLen := len(sentinel)
	wsByte := _this.ReadByte()
	if wsByte == '\r' {
		if _this.ReadByte() != '\n' {
			_this.UngetByte()
			_this.UnexpectedChar("verbatim sentinel")
		}
	}
	_this.BeginSubtoken()
	_this.SkipBytes(sentinelLen - 1)

Outer:
	for {
		_this.ReadByte()
		compareStart := _this.tokenPos - sentinelLen
		for i := sentinelLen - 1; i >= 0; i-- {
			if sentinel[i] != _this.buffer.Buffer[compareStart+i] {
				continue Outer
			}
		}
		_this.AssertAtObjectEnd("verbatim string")
		subtoken := _this.GetSubtoken()
		return subtoken[:len(subtoken)-sentinelLen]
	}
}

func (_this *CTEReadBuffer) DecodeHexBytes() []byte {
	bytes := make([]byte, 0, 8)
	firstNybble := true
	nextByte := byte(0)
	for {
		b := _this.PeekByteNoEOD()
		switch {
		case hasProperty(b, cteProperty09):
			nextByte |= b - '0'
		case hasProperty(b, ctePropertyLowercaseAF):
			nextByte |= b - 'a' + 10
		case hasProperty(b, ctePropertyUppercaseAF):
			nextByte |= b - 'A' + 10
		case hasProperty(b, ctePropertyWhitespace):
			_this.AdvanceByte()
			continue
		case b == '"':
			if !firstNybble {
				_this.ErrorAt(_this.tokenPos, "Missing last hex digit")
			}
			_this.AdvanceByte()
			return bytes
		default:
			_this.UnexpectedCharAt(_this.tokenPos, "hex encoding")
		}
		if !firstNybble {
			bytes = append(bytes, nextByte)
			nextByte = 0
		}
		nextByte <<= 4
		firstNybble = !firstNybble
		_this.AdvanceByte()
	}
}

func (_this *CTEReadBuffer) DecodeUUID() []byte {
	uuid := make([]byte, 0, 16)
	dashCount := 0
	firstNybble := true
	nextByte := byte(0)
Loop:
	for {
		b := _this.PeekByteAllowEOD()
		switch {
		case b.HasProperty(cteProperty09):
			nextByte |= byte(b - '0')
		case b.HasProperty(ctePropertyLowercaseAF):
			nextByte |= byte(b - 'a' + 10)
		case b.HasProperty(ctePropertyUppercaseAF):
			nextByte |= byte(b - 'A' + 10)
		case b == '-':
			dashCount++
			_this.AdvanceByte()
			continue
		case b.HasProperty(ctePropertyObjectEnd):
			break Loop
		default:
			_this.UnexpectedCharAt(_this.tokenPos, "UUID")
		}
		if !firstNybble {
			uuid = append(uuid, nextByte)
			nextByte = 0
		}
		nextByte <<= 4
		firstNybble = !firstNybble
		_this.AdvanceByte()
	}

	if len(uuid) != 16 ||
		dashCount != 4 ||
		_this.buffer.Buffer[_this.tokenStart+9] != '-' ||
		_this.buffer.Buffer[_this.tokenStart+14] != '-' ||
		_this.buffer.Buffer[_this.tokenStart+19] != '-' ||
		_this.buffer.Buffer[_this.tokenStart+24] != '-' {
		_this.ErrorAt(_this.tokenPos-1, "Unrecognized named value or malformed UUID")
	}

	return uuid
}

func (_this *CTEReadBuffer) DecodeDate(year int64) *compact_time.Time {
	month, _, digitCount := _this.DecodeDecimalInteger(0, nil)
	if digitCount > 2 {
		_this.Errorf("Month field is too long")
	}
	if _this.PeekByteNoEOD() != '-' {
		_this.UnexpectedChar("month")
	}
	_this.AdvanceByte()

	var day uint64
	day, _, digitCount = _this.DecodeDecimalInteger(0, nil)
	if digitCount == 0 {
		_this.UnexpectedChar("day")
	}
	if digitCount > 2 {
		_this.Errorf("Day field is too long")
	}
	if _this.PeekByteAllowEOD() != '/' {
		return compact_time.NewDate(int(year), int(month), int(day))
	}

	_this.AdvanceByte()
	var hour uint64
	hour, _, digitCount = _this.DecodeDecimalInteger(0, nil)
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
	if t.TimezoneIs == compact_time.TypeLatitudeLongitude {
		return compact_time.NewTimestampLatLong(int(year), int(month), int(day),
			int(t.Hour), int(t.Minute), int(t.Second), int(t.Nanosecond),
			int(t.LatitudeHundredths), int(t.LongitudeHundredths))
	}
	return compact_time.NewTimestamp(int(year), int(month), int(day),
		int(t.Hour), int(t.Minute), int(t.Second), int(t.Nanosecond),
		t.AreaLocation)
}

func (_this *CTEReadBuffer) DecodeTime(hour int) *compact_time.Time {
	minute, _, digitCount := _this.DecodeDecimalInteger(0, nil)
	if digitCount > 2 {
		_this.Errorf("Minute field is too long")
	}
	if _this.PeekByteNoEOD() != ':' {
		_this.UnexpectedChar("minute")
	}
	_this.AdvanceByte()
	var second uint64
	second, _, digitCount = _this.DecodeDecimalInteger(0, nil)
	if digitCount > 2 {
		_this.Errorf("Second field is too long")
	}
	var nanosecond int

	if _this.PeekByteAllowEOD() == '.' {
		_this.AdvanceByte()
		v, _, digitCount := _this.DecodeDecimalInteger(0, nil)
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
	if b.HasProperty(ctePropertyObjectEnd) {
		return compact_time.NewTime(hour, int(minute), int(second), nanosecond, "")
	}

	if b != '/' {
		_this.UnexpectedChar("time")
	}
	_this.AdvanceByte()

	if hasProperty(_this.PeekByteNoEOD(), ctePropertyAZ) {
		_this.BeginSubtoken()
		_this.ReadWhilePropertyAllowEOD(ctePropertyAreaLocation)
		if _this.PeekByteAllowEOD() == '/' {
			_this.AdvanceByte()
			_this.ReadWhilePropertyAllowEOD(ctePropertyAreaLocation)
		}

		areaLocation := string(_this.GetSubtoken())
		return compact_time.NewTime(hour, int(minute), int(second), nanosecond, areaLocation)
	}

	lat, long := _this.DecodeLatLong()
	return compact_time.NewTimeLatLong(hour, int(minute), int(second), nanosecond, lat, long)
}

func (_this *CTEReadBuffer) DecodeLatLongPortion(name string) (value int) {
	whole, _, digitCount := _this.DecodeDecimalInteger(0, nil)
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
		fractional, _, digitCount = _this.DecodeDecimalInteger(0, nil)
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

func (_this *CTEReadBuffer) DecodeLatLong() (latitudeHundredths, longitudeHundredths int) {
	latitudeHundredths = _this.DecodeLatLongPortion("latitude")

	if _this.PeekByteNoEOD() != '/' {
		_this.UnexpectedChar("latitude/longitude")
	}
	_this.AdvanceByte()

	longitudeHundredths = _this.DecodeLatLongPortion("longitude")

	return
}

func (_this *CTEReadBuffer) DecodeDecimalFloat(sign int64, coefficient uint64, bigCoefficient *big.Int, coefficientDigitCount int) (value compact_float.DFloat, bigValue *apd.Decimal) {
	exponent := int32(0)
	fractionalDigitCount := 0
	coefficient, bigCoefficient, fractionalDigitCount = _this.DecodeDecimalInteger(coefficient, bigCoefficient)
	if fractionalDigitCount == 0 {
		_this.UnexpectedChar("float fractional")
	}

	if _this.PeekByteAllowEOD() == 'e' {
		_this.AdvanceByte()
		exponentSign := int32(1)
		switch _this.PeekByteNoEOD() {
		case '+':
			_this.AdvanceByte()
		case '-':
			exponentSign = -1
			_this.AdvanceByte()
		}
		exp, bigExp, digitCount := _this.DecodeDecimalInteger(0, nil)
		if digitCount == 0 {
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

	value = compact_float.DFloatValue(int32(exponent), int64(coefficient)*int64(sign))
	return
}

func (_this *CTEReadBuffer) DecodeHexFloat(sign int64, coefficient uint64, coefficientDigitCount int) float64 {
	exponent := 0
	fractionalDigitCount := 0
	coefficient, fractionalDigitCount = _this.DecodeHexInteger(coefficient)
	if fractionalDigitCount == 0 {
		_this.UnexpectedChar("float fractional")
	}

	if _this.PeekByteAllowEOD() == 'p' {
		_this.AdvanceByte()
		exponentSign := 1
		switch _this.PeekByteNoEOD() {
		case '+':
			_this.AdvanceByte()
		case '-':
			exponentSign = -1
			_this.AdvanceByte()
		}
		exp, bigExp, digitCount := _this.DecodeDecimalInteger(0, nil)
		if digitCount == 0 {
			_this.UnexpectedChar("float exponent")
		}
		if bigExp != nil {
			_this.Errorf("Exponent too big")
		}
		exponent = int(exp) * exponentSign
	}

	exponent -= fractionalDigitCount * 4

	// TODO: Overflow

	return float64(sign) * float64(coefficient) * math.Pow(float64(2), float64(exponent))
}

func (_this *CTEReadBuffer) DecodeSingleLineComment() string {
	// We're already past the "//" by this point
	_this.BeginSubtoken()
	_this.readUntilByte('\n')
	subtoken := _this.GetSubtoken()
	if len(subtoken) > 0 && subtoken[len(subtoken)-1] == '\r' {
		subtoken = subtoken[:len(subtoken)-1]
	}
	_this.AdvanceByte()

	return string(subtoken)
}

func (_this *CTEReadBuffer) DecodeMultilineComment() (string, nextType) {
	lastByte := _this.ReadByte()

	for {
		firstByte := lastByte
		lastByte = _this.ReadByte()

		if firstByte == '*' && lastByte == '/' {
			token := _this.GetToken()
			token = token[:len(token)-2]
			_this.EndToken()
			return string(token), nextIsCommentEnd
		}

		if firstByte == '/' && lastByte == '*' {
			token := _this.GetToken()
			token = token[:len(token)-2]
			_this.EndToken()
			return string(token), nextIsCommentBegin
		}
	}
}

func (_this *CTEReadBuffer) DecodeMarkupContent() (string, nextType) {
	isCommentInitiator := false

	for {
		currentByte := _this.ReadByte()
		switch currentByte {
		case '\\':
			panic("TODO: Escape sequences")
		case '<':
			token := _this.GetToken()
			token = token[:len(token)-1]
			_this.EndToken()
			return string(token), nextIsMarkupBegin
		case '>':
			token := _this.GetToken()
			token = token[:len(token)-1]
			_this.EndToken()
			return string(token), nextIsMarkupEnd
		case '/':
			if isCommentInitiator {
				token := _this.GetToken()
				token = token[:len(token)-2]
				_this.EndToken()
				return string(token), nextIsSingleLineComment
			}
		case '*':
			if isCommentInitiator {
				token := _this.GetToken()
				token = token[:len(token)-2]
				_this.EndToken()
				return string(token), nextIsCommentBegin
			}
		}

		isCommentInitiator = currentByte == '/'
	}
}

type nextType int

const (
	nextIsCommentBegin nextType = iota
	nextIsCommentEnd
	nextIsSingleLineComment
	nextIsMarkupBegin
	nextIsMarkupEnd
)
