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
	"math"
	"math/big"

	"github.com/kstenerud/go-concise-encoding/internal/common"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
)

type CTEReadBuffer struct {
	document      []byte
	tokenStart    int
	subtokenStart int
	tokenPos      int
	endPos        int
}

func NewReadBuffer(document []byte) *CTEReadBuffer {
	_this := &CTEReadBuffer{}
	_this.Init(document)
	return _this
}

func (_this *CTEReadBuffer) Init(document []byte) {
	_this.document = document
	_this.endPos = len(document) - 1
}

// Bytes

func (_this *CTEReadBuffer) GetByteAt(index int) cteByte {
	return cteByte(_this.document[index])
}

func (_this *CTEReadBuffer) PeekByte() cteByte {
	if _this.IsEndOfDocument() {
		return cteByteEndOfDocument
	}
	return _this.GetByteAt(_this.tokenPos)
}

func (_this *CTEReadBuffer) ReadByte() (b cteByte) {
	if _this.IsEndOfDocument() {
		return cteByteEndOfDocument
	}
	b = _this.GetByteAt(_this.tokenPos)
	_this.AdvanceByte()
	return
}

func (_this *CTEReadBuffer) AdvanceByte() {
	_this.tokenPos++
}

func (_this *CTEReadBuffer) ReadUntilByte(b byte) {
	i := _this.tokenPos
	for ; i <= _this.endPos && _this.document[i] != b; i++ {
	}
	_this.tokenPos = i
}

func (_this *CTEReadBuffer) ReadUntilProperty(property cteByteProprty) {
	i := _this.tokenPos
	for ; i <= _this.endPos && !hasProperty(_this.document[i], property); i++ {
	}
	_this.tokenPos = i
}

func (_this *CTEReadBuffer) ReadWhileProperty(property cteByteProprty) {
	i := _this.tokenPos
	for ; i <= _this.endPos && hasProperty(_this.document[i], property); i++ {
	}
	_this.tokenPos = i
}

func (_this *CTEReadBuffer) GetCharBeginIndex(index int) int {
	i := index
	// UTF-8 continuation characters have the form 10xxxxxx
	for b := _this.document[i]; b >= 0x80 && b <= 0xc0; b = _this.document[i] {
		i--
	}
	return i
}

func (_this *CTEReadBuffer) GetCharEndIndex(index int) int {
	i := index
	// UTF-8 continuation characters have the form 10xxxxxx
	for b := _this.document[i]; b >= 0x80 && b <= 0xc0; b = _this.document[i] {
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
	return _this.tokenPos > _this.endPos
}

// Tokens

func (_this *CTEReadBuffer) GetToken() []byte {
	return _this.document[_this.tokenStart:_this.tokenPos]
}

func (_this *CTEReadBuffer) GetTokenLength() int {
	return _this.tokenPos - _this.tokenStart
}

func (_this *CTEReadBuffer) GetTokenFirstByte() byte {
	return _this.document[_this.tokenStart]
}

func (_this *CTEReadBuffer) EndToken() {
	_this.tokenStart = _this.tokenPos
}

func (_this *CTEReadBuffer) BeginSubtoken() {
	_this.subtokenStart = _this.tokenPos
}

func (_this *CTEReadBuffer) GetSubtoken() []byte {
	return _this.document[_this.subtokenStart:_this.tokenPos]
}

// Errors

func (_this *CTEReadBuffer) AssertNotEOD() {
	if _this.IsEndOfDocument() {
		_this.UnexpectedEOD()
	}
}

func (_this *CTEReadBuffer) AssertAtObjectEnd(decoding string) {
	if !_this.PeekByte().HasProperty(ctePropertyObjectEnd) {
		_this.UnexpectedChar(decoding)
	}
}

func (_this *CTEReadBuffer) UnexpectedEOD() {
	_this.Errorf("Unexpected end of document")
}

func (_this *CTEReadBuffer) ErrorAt(index int, format string, args ...interface{}) {
	lineNumber := 1
	lineStart := 0
	for i := 0; i < index; i++ {
		if _this.document[i] == '\n' {
			lineNumber++
			lineStart = i
		}
	}

	colNumber := 1
	for i := lineStart; i < index; i++ {
		b := _this.GetByteAt(i)
		if b < 0x80 || b > 0xc0 {
			colNumber++
		}
	}

	msg := fmt.Sprintf(format, args...)
	panic(fmt.Errorf("Offset %v (line %v, col %v): %v", index, lineNumber, colNumber, msg))
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
	if index > _this.endPos {
		return "EOD"
	}

	charStart := _this.GetCharBeginIndex(index)
	charEnd := _this.GetCharEndIndex(index)
	if charEnd-charStart > 1 {
		return string(_this.document[charStart:charEnd])
	}

	b := _this.document[charStart]
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

// Decoders

func (_this *CTEReadBuffer) SkipWhitespace() {
	endPos := _this.endPos
	i := _this.tokenPos
	for ; i <= endPos; i++ {
		if !hasProperty(_this.document[i], ctePropertyWhitespace) {
			break
		}
	}
	_this.tokenPos = i
	_this.EndToken()
}

func (_this *CTEReadBuffer) DecodeBinaryInteger() (value uint64) {
	endPos := _this.endPos
	i := _this.tokenPos
	for ; i <= endPos; i++ {
		b := _this.document[i]
		if !hasProperty(b, ctePropertyBinaryDigit) {
			break
		}
		oldValue := value
		value = value<<1 + uint64(b-'0')
		if value < oldValue {
			_this.Errorf("Overflow reading binary integer")
		}
	}
	_this.tokenPos = i
	return
}

func (_this *CTEReadBuffer) DecodeOctalInteger() (value uint64) {
	endPos := _this.endPos
	i := _this.tokenPos
	for ; i <= endPos; i++ {
		b := _this.document[i]
		if !hasProperty(b, ctePropertyOctalDigit) {
			break
		}
		oldValue := value
		value = value<<3 + uint64(b-'0')
		if value < oldValue {
			_this.Errorf("Overflow reading octal value")
		}
	}
	_this.tokenPos = i
	return
}

// startValue is only used if bigStartValue is nil
// bigValue will be nil unless the value was too big for a uint64
func (_this *CTEReadBuffer) DecodeDecimalInteger(startValue uint64, bigStartValue *big.Int) (value uint64, bigValue *big.Int, digitCount int) {
	endPos := _this.endPos
	i := _this.tokenPos

	if bigStartValue == nil {
		value = startValue
		for ; i <= endPos; i++ {
			b := _this.document[i]
			if !hasProperty(b, cteProperty09) {
				break
			}
			oldValue := value
			value = value*10 + uint64(b-'0')
			if value/10 != oldValue {
				bigStartValue = new(big.Int).SetUint64(oldValue)
				break
			}
		}
	}

	if bigStartValue != nil {
		bigValue = bigStartValue
		for ; i <= endPos; i++ {
			b := _this.document[i]
			if !hasProperty(b, cteProperty09) {
				break
			}
			bigValue = bigValue.Mul(bigValue, common.BigInt10)
			bigValue = bigValue.Add(bigValue, big.NewInt(int64(b-'0')))
		}
	}

	digitCount = i - _this.tokenPos
	_this.tokenPos = i
	return
}

func (_this *CTEReadBuffer) DecodeHexInteger(startValue uint64) (value uint64, digitCount int) {
	value = startValue
	endPos := _this.endPos
	i := _this.tokenPos
Loop:
	for ; i <= endPos; i++ {
		b := _this.document[i]
		cp := cteByteProperties[b]
		oldValue := value
		switch {
		case cp.HasProperty(cteProperty09):
			value = value<<4 + uint64(b-'0')
		case cp.HasProperty(ctePropertyLowercaseAF):
			value = value<<4 + uint64(b-'a') + 10
		case cp.HasProperty(ctePropertyUppercaseAF):
			value = value<<4 + uint64(b-'A') + 10
		default:
			break Loop
		}
		if value < oldValue {
			_this.Errorf("Overflow reading hex value")
		}
	}
	digitCount = i - _this.tokenPos
	_this.tokenPos = i
	return
}

func (_this *CTEReadBuffer) DecodeQuotedStringWithEscapes(startPos, firstEscape int) string {
	_this.Errorf("TODO: CTEDecoder: Escape sequences")
	return ""
}

func (_this *CTEReadBuffer) DecodeQuotedString() string {
	startPos := _this.tokenPos
	endPos := _this.endPos
	i := startPos
	for ; i <= endPos; i++ {
		switch _this.document[i] {
		case '"':
			str := string(_this.document[startPos:i])
			_this.tokenPos = i + 1
			return str
		case '\\':
			return _this.DecodeQuotedStringWithEscapes(startPos, i)
		}
	}
	_this.UnexpectedEOD()
	return ""
}

func (_this *CTEReadBuffer) DecodeVerbatimString() (value []byte) {
	_this.BeginSubtoken()
	_this.ReadUntilProperty(ctePropertyWhitespace)
	sentinel := _this.GetSubtoken()
	sentinelLen := len(sentinel)
	wsByte := _this.ReadByte()
	if wsByte == '\r' {
		if _this.ReadByte() != '\n' {
			_this.UnexpectedChar("verbatim sentinel")
		}
	}
	_this.BeginSubtoken()

Outer:
	for _this.tokenPos += sentinelLen; _this.tokenPos <= _this.endPos+1; _this.tokenPos++ {
		docStart := _this.tokenPos - sentinelLen
		for i := sentinelLen - 1; i >= 0; i-- {
			if sentinel[i] != _this.document[docStart+i] {
				continue Outer
			}
		}
		_this.AssertAtObjectEnd("verbatim string")
		subtoken := _this.GetSubtoken()
		return subtoken[:len(subtoken)-sentinelLen]
	}
	_this.Errorf("Verbatim sentinel sequence [%v] not found in document", string(sentinel))
	return nil
}

func (_this *CTEReadBuffer) DecodeHexBytes() []byte {
	endPos := _this.endPos
	i := _this.tokenPos
	bytes := make([]byte, 0, 8)
	firstNybble := true
	nextByte := byte(0)
	for ; i <= endPos; i++ {
		b := _this.document[i]
		cp := cteByteProperties[b]
		switch {
		case cp.HasProperty(cteProperty09):
			nextByte |= b - '0'
		case cp.HasProperty(ctePropertyLowercaseAF):
			nextByte |= b - 'a' + 10
		case cp.HasProperty(ctePropertyUppercaseAF):
			nextByte |= b - 'A' + 10
		case cp.HasProperty(ctePropertyWhitespace):
			continue
		case b == '"':
			if !firstNybble {
				_this.ErrorAt(i, "Missing last hex digit")
			}
			_this.tokenPos = i + 1
			return bytes
		default:
			_this.UnexpectedCharAt(i, "hex encoding")
		}
		if !firstNybble {
			bytes = append(bytes, nextByte)
			nextByte = 0
		}
		nextByte <<= 4
		firstNybble = !firstNybble
	}
	_this.UnexpectedEOD()
	return nil
}

func (_this *CTEReadBuffer) DecodeUUID() []byte {
	endPos := _this.endPos
	i := _this.tokenPos
	uuid := make([]byte, 0, 16)
	dashCount := 0
	firstNybble := true
	nextByte := byte(0)
Loop:
	for ; i <= endPos; i++ {
		b := _this.document[i]
		cp := cteByteProperties[b]
		switch {
		case cp.HasProperty(cteProperty09):
			nextByte |= b - '0'
		case cp.HasProperty(ctePropertyLowercaseAF):
			nextByte |= b - 'a' + 10
		case cp.HasProperty(ctePropertyUppercaseAF):
			nextByte |= b - 'A' + 10
		case b == '-':
			dashCount++
			continue
		case cp.HasProperty(ctePropertyObjectEnd):
			break Loop
		default:
			_this.UnexpectedCharAt(i, "UUID")
		}
		if !firstNybble {
			uuid = append(uuid, nextByte)
			nextByte = 0
		}
		nextByte <<= 4
		firstNybble = !firstNybble
	}

	if len(uuid) != 16 ||
		dashCount != 4 ||
		_this.document[_this.tokenStart+9] != '-' ||
		_this.document[_this.tokenStart+14] != '-' ||
		_this.document[_this.tokenStart+19] != '-' ||
		_this.document[_this.tokenStart+24] != '-' {
		_this.ErrorAt(i, "Unrecognized named value or malformed UUID")
	}

	_this.tokenPos = i
	return uuid
}

func (_this *CTEReadBuffer) DecodeDate(year int64) *compact_time.Time {
	month, _, digitCount := _this.DecodeDecimalInteger(0, nil)
	if digitCount > 2 {
		_this.Errorf("Month field is too long")
	}
	if _this.PeekByte() != '-' {
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
	if _this.PeekByte() != '/' {
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
	if _this.PeekByte() != ':' {
		_this.UnexpectedChar("minute")
	}
	_this.AdvanceByte()
	var second uint64
	second, _, digitCount = _this.DecodeDecimalInteger(0, nil)
	if digitCount > 2 {
		_this.Errorf("Second field is too long")
	}
	var nanosecond int

	if _this.PeekByte() == '.' {
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

	b := _this.PeekByte()
	if b.HasProperty(ctePropertyObjectEnd) {
		return compact_time.NewTime(hour, int(minute), int(second), nanosecond, "")
	}

	if b != '/' {
		_this.UnexpectedChar("time")
	}
	_this.AdvanceByte()

	if _this.PeekByte().HasProperty(ctePropertyAZ) {
		areaLocationStart := _this.tokenPos
		_this.ReadWhileProperty(ctePropertyAreaLocation)
		if _this.PeekByte() == '/' {
			_this.AdvanceByte()
			_this.ReadWhileProperty(ctePropertyAreaLocation)
		}
		areaLocation := string(_this.document[areaLocationStart:_this.tokenPos])
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
	b := _this.PeekByte()
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

	if _this.PeekByte() != '/' {
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

	if _this.PeekByte() == 'e' {
		_this.AdvanceByte()
		exponentSign := int32(1)
		switch _this.PeekByte() {
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

	if _this.PeekByte() == 'p' {
		_this.AdvanceByte()
		exponentSign := 1
		switch _this.PeekByte() {
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
	commentStart := _this.tokenPos
	_this.ReadUntilByte('\n')
	commentEnd := _this.tokenPos
	if _this.document[commentEnd-1] == '\r' {
		commentEnd--
	}
	_this.AdvanceByte()

	return string(_this.document[commentStart:commentEnd])
}

func (_this *CTEReadBuffer) DecodeMultilineComment() (string, nextType) {
	lastByte := _this.ReadByte()
	_this.AssertNotEOD()

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

		_this.AssertNotEOD()
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
		case cteByteEndOfDocument:
			_this.UnexpectedEOD()
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
