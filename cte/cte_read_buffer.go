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
	"strings"

	"github.com/kstenerud/go-concise-encoding/buffer"

	"github.com/kstenerud/go-concise-encoding/internal/common"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
)

type ReadBuffer struct {
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
func NewReadBuffer(reader io.Reader, readBufferSize int, loWaterByteCount int) *ReadBuffer {
	_this := &ReadBuffer{}
	_this.Init(reader, readBufferSize, loWaterByteCount)
	return _this
}

// Init the read buffer. The buffer will be empty until RefillIfNecessary() is
// called.
//
// readBufferSize determines the initial size of the buffer, and
// loWaterByteCount determines when RefillIfNecessary() refills the buffer from
// the reader.
func (_this *ReadBuffer) Init(reader io.Reader, readBufferSize int, loWaterByteCount int) {
	_this.buffer.Init(reader, readBufferSize, loWaterByteCount)
}

func (_this *ReadBuffer) Reset() {
	_this.buffer.Reset()
	_this.tokenStart = 0
	_this.subtokenStart = 0
	_this.tokenPos = 0
	_this.tempBuffer.Reset()
	_this.lineCount = 0
	_this.colCount = 0
}

// Bytes

func (_this *ReadBuffer) PeekByteNoEOD() byte {
	_this.RefillIfNecessary()
	if _this.IsEndOfDocument() {
		_this.UnexpectedEOD()
	}
	return _this.buffer.ByteAtOffset(_this.tokenPos)
}

func (_this *ReadBuffer) PeekByteAllowEOD() cteByte {
	_this.RefillIfNecessary()
	if _this.IsEndOfDocument() {
		return cteByteEndOfDocument
	}
	return cteByte(_this.buffer.ByteAtOffset(_this.tokenPos))
}

func (_this *ReadBuffer) ReadByte() byte {
	b := _this.PeekByteNoEOD()
	_this.AdvanceByte()
	return b
}

func (_this *ReadBuffer) SkipBytes(byteCount int) {
	_this.buffer.RequireBytes(_this.tokenPos, byteCount)
	_this.tokenPos += byteCount
}

func (_this *ReadBuffer) AdvanceByte() {
	_this.tokenPos++
}

func (_this *ReadBuffer) RefillIfNecessary() {
	offset := _this.buffer.RefillIfNecessary(_this.tokenStart, _this.tokenPos)
	if _this.tokenStart-offset < 0 {
		panic(fmt.Errorf("TokenStart was %v, tokenPos was %v", _this.tokenStart-offset, _this.tokenPos))
	}
	_this.tokenPos += offset
	_this.tokenStart += offset
	_this.subtokenStart += offset
}

func (_this *ReadBuffer) ReadUntilPropertyNoEOD(property cteByteProprty) {
	for !hasProperty(_this.PeekByteNoEOD(), property) {
		_this.AdvanceByte()
	}
}

func (_this *ReadBuffer) ReadWhilePropertyNoEOD(property cteByteProprty) {
	for hasProperty(_this.PeekByteNoEOD(), property) {
		_this.AdvanceByte()
	}
}

func (_this *ReadBuffer) ReadWhilePropertyAllowEOD(property cteByteProprty) {
	for _this.PeekByteAllowEOD().HasProperty(property) {
		_this.AdvanceByte()
	}
}

func (_this *ReadBuffer) GetCharBeginIndex(index int) int {
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

func (_this *ReadBuffer) GetCharEndIndex(index int) int {
	i := index
	// UTF-8 continuation characters have the form 10xxxxxx
	for b := _this.buffer.Buffer[i]; b >= 0x80 && b <= 0xc0 && i < len(_this.buffer.Buffer); b = _this.buffer.Buffer[i] {
		i++
	}
	return i
}

func (_this *ReadBuffer) UngetByte() {
	_this.tokenPos--
}

func (_this *ReadBuffer) UngetBytes(byteCount int) {
	_this.tokenPos -= byteCount
}

func (_this *ReadBuffer) UngetAll() {
	_this.tokenPos = _this.tokenStart
}

func (_this *ReadBuffer) IsEndOfDocument() bool {
	return _this.buffer.IsEOF() && !_this.hasUnreadByte()
}

func (_this *ReadBuffer) hasUnreadByte() bool {
	return _this.buffer.HasByteAtOffset(_this.tokenPos)
}

func (_this *ReadBuffer) readUntilByte(b byte) {
	for _this.PeekByteNoEOD() != b {
		_this.AdvanceByte()
	}
}

// Tokens

func (_this *ReadBuffer) GetToken() []byte {
	return _this.buffer.Buffer[_this.tokenStart:_this.tokenPos]
}

func (_this *ReadBuffer) GetTokenLength() int {
	return _this.tokenPos - _this.tokenStart
}

func (_this *ReadBuffer) GetTokenFirstByte() byte {
	return _this.buffer.Buffer[_this.tokenStart]
}

func (_this *ReadBuffer) EndToken() {
	_this.lineCount, _this.colCount = _this.countLinesColsToIndex(_this.tokenPos)
	_this.tokenStart = _this.tokenPos
}

func (_this *ReadBuffer) BeginSubtoken() {
	_this.subtokenStart = _this.tokenPos
}

func (_this *ReadBuffer) GetSubtoken() []byte {
	return _this.buffer.Buffer[_this.subtokenStart:_this.tokenPos]
}

// Errors

func (_this *ReadBuffer) AssertAtObjectEnd(decoding string) {
	if !_this.PeekByteAllowEOD().HasProperty(ctePropertyObjectEnd) {
		_this.UnexpectedChar(decoding)
	}
}

func (_this *ReadBuffer) UnexpectedEOD() {
	_this.Errorf("unexpected end of document")
}

func (_this *ReadBuffer) ErrorAt(index int, format string, args ...interface{}) {
	lineCount, colCount := _this.countLinesColsToIndex(index)

	msg := fmt.Sprintf(format, args...)
	panic(fmt.Errorf("offset %v (line %v, col %v): %v", index, lineCount+1, colCount+1, msg))
}

func (_this *ReadBuffer) Errorf(format string, args ...interface{}) {
	_this.ErrorAt(_this.tokenPos, format, args...)
}

func (_this *ReadBuffer) UnexpectedCharAt(index int, decoding string) {
	_this.ErrorAt(index, "unexpected [%v] while decoding %v", _this.DescribeCharAt(index), decoding)
}

func (_this *ReadBuffer) UnexpectedChar(decoding string) {
	_this.UnexpectedCharAt(_this.tokenPos, decoding)
}

func (_this *ReadBuffer) DescribeCharAt(index int) string {
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

func (_this *ReadBuffer) DescribeCurrentChar() string {
	return _this.DescribeCharAt(_this.tokenPos)
}

func (_this *ReadBuffer) countLinesColsToIndex(index int) (lineCount, colCount int) {
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

func (_this *ReadBuffer) SkipWhitespace() {
	_this.ReadWhilePropertyAllowEOD(ctePropertyWhitespace)
}

const maxPreShiftBinary = uint64(0x7fffffffffffffff)

func (_this *ReadBuffer) DecodeBinaryInteger() (value uint64) {
	for {
		b := _this.PeekByteAllowEOD()
		if b.HasProperty(ctePropertyNumericWhitespace) {
			_this.AdvanceByte()
			continue
		}
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

func (_this *ReadBuffer) DecodeOctalInteger() (value uint64) {
	for {
		b := _this.PeekByteAllowEOD()
		if b.HasProperty(ctePropertyNumericWhitespace) {
			_this.AdvanceByte()
			continue
		}
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
func (_this *ReadBuffer) DecodeDecimalInteger(startValue uint64, bigStartValue *big.Int) (value uint64, bigValue *big.Int, digitCount int) {
	if bigStartValue == nil {
		value = startValue
		for {
			b := _this.PeekByteAllowEOD()
			if b.HasProperty(ctePropertyNumericWhitespace) {
				_this.AdvanceByte()
				continue
			}
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

func (_this *ReadBuffer) DecodeHexInteger(startValue uint64) (value uint64, digitCount int) {
	value = startValue
	for {
		b := _this.PeekByteAllowEOD()
		nextNybble := uint64(0)
		switch {
		case b.HasProperty(ctePropertyNumericWhitespace):
			_this.AdvanceByte()
			continue
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

func (_this *ReadBuffer) DecodeQuotedStringWithEscapes() []byte {
	sb := strings.Builder{}
	sb.Write(_this.GetSubtoken())
	for {
		b := _this.PeekByteNoEOD()
		switch b {
		case '"':
			_this.AdvanceByte()
			return []byte(sb.String())
		case '\\':
			_this.AdvanceByte()
			escape := _this.PeekByteNoEOD()
			switch escape {
			case 'r':
				sb.WriteByte('\r')
			case 'n':
				sb.WriteByte('\n')
			case 't':
				sb.WriteByte('\t')
			case '\r':
				_this.AdvanceByte()
				if _this.PeekByteNoEOD() != '\n' {
					_this.UnexpectedChar("\\r\\n continuation")
				}
			case '\n':
				// Nothing to do
			case '"':
				sb.WriteByte('"')
			case '*':
				sb.WriteByte('*')
			case '/':
				sb.WriteByte('/')
			case '\\':
				sb.WriteByte('\\')
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				sb.WriteRune(_this.decodeUnicodeEscapeSequence(int(escape - '0')))
			default:
				_this.UnexpectedChar("quoted string escape sequence")
			}
			_this.AdvanceByte()
		default:
			sb.WriteByte(b)
			_this.AdvanceByte()
		}
	}
}

func (_this *ReadBuffer) decodeUnicodeEscapeSequence(length int) (codepoint rune) {
	for i := 0; i < length; i++ {
		_this.AdvanceByte()
		b := _this.PeekByteNoEOD()
		switch {
		case hasProperty(b, cteProperty09):
			codepoint = (codepoint << 4) | (rune(b) - '0')
		case hasProperty(b, ctePropertyLowercaseAF):
			codepoint = (codepoint << 4) | (rune(b) - 'a' + 10)
		case hasProperty(b, ctePropertyUppercaseAF):
			codepoint = (codepoint << 4) | (rune(b) - 'A' + 10)
		default:
			_this.UnexpectedChar("unicode sequence")
		}
	}
	return
}

func (_this *ReadBuffer) DecodeQuotedString() []byte {
	_this.BeginSubtoken()
	for {
		b := _this.PeekByteNoEOD()
		switch b {
		case '"':
			token := _this.GetSubtoken()
			_this.AdvanceByte()
			return token
		case '\\':
			return _this.DecodeQuotedStringWithEscapes()
		default:
			_this.AdvanceByte()
		}
	}
}

func (_this *ReadBuffer) DecodeURI() []byte {
	_this.BeginSubtoken()
	for {
		b := _this.PeekByteNoEOD()
		switch b {
		case '"':
			token := _this.GetSubtoken()
			_this.AdvanceByte()
			return token
		default:
			_this.AdvanceByte()
		}
	}
}

func (_this *ReadBuffer) DecodeVerbatimString() []byte {
	_this.BeginSubtoken()
	_this.ReadUntilPropertyNoEOD(ctePropertyWhitespace)
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
	_this.SkipBytes(sentinelLen - 1)

Outer:
	for {
		_this.ReadByte()
		sentinelStart := _this.tokenStart + sentinelOffset
		compareStart := _this.tokenPos - sentinelLen
		for i := sentinelLen - 1; i >= 0; i-- {
			if _this.buffer.Buffer[sentinelStart+i] != _this.buffer.Buffer[compareStart+i] {
				continue Outer
			}
		}
		_this.AssertAtObjectEnd("verbatim string")
		subtoken := _this.GetSubtoken()
		return subtoken[:len(subtoken)-sentinelLen]
	}
}

func (_this *ReadBuffer) DecodeCustomText() []byte {
	_this.BeginSubtoken()
	for {
		b := _this.PeekByteNoEOD()
		switch b {
		case '"':
			token := _this.GetSubtoken()
			_this.AdvanceByte()
			return token
		case '\\':
			return _this.DecodeCustomTextWithEscapes()
		default:
			_this.AdvanceByte()
		}
	}
}

func (_this *ReadBuffer) DecodeCustomTextWithEscapes() []byte {
	sb := strings.Builder{}
	sb.Write(_this.GetSubtoken())
	for {
		b := _this.PeekByteNoEOD()
		switch b {
		case '"':
			_this.AdvanceByte()
			return []byte(sb.String())
		case '\\':
			_this.AdvanceByte()
			escape := _this.PeekByteNoEOD()
			switch escape {
			case '"':
				sb.WriteByte('"')
			case '\\':
				sb.WriteByte('\\')
			case 'r':
				sb.WriteByte('\r')
			case 'n':
				sb.WriteByte('\n')
			case 't':
				sb.WriteByte('\t')
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				sb.WriteRune(_this.decodeUnicodeEscapeSequence(int(escape - '0')))
			default:
				_this.UnexpectedChar("custom text escape sequence")
			}
			_this.AdvanceByte()
		default:
			sb.WriteByte(b)
			_this.AdvanceByte()
		}
	}
}
func (_this *ReadBuffer) DecodeHexBytes() []byte {
	data := make([]byte, 0, 8)
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
			return data
		default:
			_this.UnexpectedCharAt(_this.tokenPos, "hex encoding")
		}
		if !firstNybble {
			data = append(data, nextByte)
			nextByte = 0
		}
		nextByte <<= 4
		firstNybble = !firstNybble
		_this.AdvanceByte()
	}
}

func (_this *ReadBuffer) DecodeUUID() []byte {
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

func (_this *ReadBuffer) DecodeDate(year int64) *compact_time.Time {
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

func (_this *ReadBuffer) DecodeTime(hour int) *compact_time.Time {
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

func (_this *ReadBuffer) DecodeLatLongPortion(name string) (value int) {
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

func (_this *ReadBuffer) DecodeLatLong() (latitudeHundredths, longitudeHundredths int) {
	latitudeHundredths = _this.DecodeLatLongPortion("latitude")

	if _this.PeekByteNoEOD() != '/' {
		_this.UnexpectedChar("latitude/longitude")
	}
	_this.AdvanceByte()

	longitudeHundredths = _this.DecodeLatLongPortion("longitude")

	return
}

func (_this *ReadBuffer) DecodeDecimalFloat(sign int64, coefficient uint64, bigCoefficient *big.Int, coefficientDigitCount int) (value compact_float.DFloat, bigValue *apd.Decimal) {
	exponent := int32(0)
	fractionalDigitCount := 0
	coefficient, bigCoefficient, fractionalDigitCount = _this.DecodeDecimalInteger(coefficient, bigCoefficient)
	if fractionalDigitCount == 0 {
		_this.UnexpectedChar("float fractional")
	}

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

	value = compact_float.DFloatValue(exponent, int64(coefficient)*sign)
	return
}

func (_this *ReadBuffer) DecodeHexFloat(sign int64, coefficient uint64, coefficientDigitCount int) float64 {
	exponent := 0
	fractionalDigitCount := 0
	coefficient, fractionalDigitCount = _this.DecodeHexInteger(coefficient)
	if fractionalDigitCount == 0 {
		_this.UnexpectedChar("float fractional")
	}

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

func (_this *ReadBuffer) DecodeSingleLineComment() []byte {
	// We're already past the "//" by this point
	_this.BeginSubtoken()
	_this.readUntilByte('\n')
	subtoken := _this.GetSubtoken()
	if len(subtoken) > 0 && subtoken[len(subtoken)-1] == '\r' {
		subtoken = subtoken[:len(subtoken)-1]
	}
	_this.AdvanceByte()

	return subtoken
}

func (_this *ReadBuffer) DecodeMultilineComment() ([]byte, nextType) {
	lastByte := _this.ReadByte()

	for {
		firstByte := lastByte
		lastByte = _this.ReadByte()

		if firstByte == '*' && lastByte == '/' {
			token := _this.GetToken()
			token = token[:len(token)-2]
			_this.EndToken()
			return token, nextIsCommentEnd
		}

		if firstByte == '/' && lastByte == '*' {
			token := _this.GetToken()
			token = token[:len(token)-2]
			_this.EndToken()
			return token, nextIsCommentBegin
		}
	}
}

func (_this *ReadBuffer) DecodeMarkupContent() ([]byte, nextType) {
	isInitialWS := true
	wsCount := 0
	isCommentInitiator := false
	sb := strings.Builder{}

	completeStringPortion := func() {
		if wsCount > 0 {
			if !isInitialWS {
				sb.WriteByte(' ')
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
			sb.WriteByte('/')
			isCommentInitiator = false
		}
	}

	for {
		currentByte := _this.ReadByte()
		switch currentByte {
		case '\r', '\n', '\t', ' ':
			wsCount++
		case '\\':
			completeStringPortion()
			escape := _this.PeekByteNoEOD()
			switch escape {
			case 'r':
				sb.WriteByte('\r')
			case 'n':
				sb.WriteByte('\n')
			case 't':
				sb.WriteByte('\t')
			case '*':
				sb.WriteByte('*')
			case '/':
				sb.WriteByte('/')
			case '<':
				sb.WriteByte('<')
			case '>':
				sb.WriteByte('>')
			case '`':
				sb.WriteByte('`')
			case '\\':
				sb.WriteByte('\\')
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				sb.WriteRune(_this.decodeUnicodeEscapeSequence(int(escape - '0')))
			default:
				_this.UnexpectedChar("quoted string escape sequence")
			}
			_this.AdvanceByte()
		case '<':
			completeStringPortion()
			addSlashIfNeeded()
			_this.EndToken()
			return []byte(sb.String()), nextIsMarkupBegin
		case '>':
			completeContentsPortion()
			addSlashIfNeeded()
			_this.EndToken()
			return []byte(sb.String()), nextIsMarkupEnd
		case '/':
			if isCommentInitiator {
				completeStringPortion()
				_this.EndToken()
				return []byte(sb.String()), nextIsSingleLineComment
			} else {
				isCommentInitiator = true
			}
		case '*':
			if isCommentInitiator {
				completeStringPortion()
				_this.EndToken()
				return []byte(sb.String()), nextIsCommentBegin
			} else {
				sb.WriteByte(currentByte)
			}
		case '`':
			completeStringPortion()
			addSlashIfNeeded()
			return []byte(sb.String()), nextIsVerbatimString
		default:
			completeStringPortion()
			addSlashIfNeeded()
			sb.WriteByte(currentByte)
		}
	}
}

// Decode a marker ID. asString will be empty if the result is an integer.
func (_this *ReadBuffer) DecodeMarkerID() (asString []byte, asUint uint64) {
	isInteger := true
	_this.BeginSubtoken()
Loop:
	for {
		b := _this.PeekByteAllowEOD()
		switch {
		case b.HasProperty(cteProperty09):
			if isInteger {
				if asUint > maxPreShiftDecimal {
					_this.Errorf("Integer marker ID is too big")
				} else {
					asUint = asUint*10 + uint64(b-'0')
				}
			}
		case b.HasProperty(ctePropertyMarkerID):
			isInteger = false
		default:
			break Loop
		}
		_this.AdvanceByte()
	}

	subtoken := _this.GetSubtoken()
	if len(subtoken) == 0 {
		_this.Errorf("Missing marker ID")
	}

	if !isInteger {
		asString = subtoken
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
	nextIsVerbatimString
	nextIsMarkupBegin
	nextIsMarkupEnd
)
