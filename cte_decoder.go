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

package concise_encoding

import (
	"bytes"
	"fmt"
	"math"
	"math/big"
	"strings"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
)

type CTEDecoderOptions struct {
	// TODO: ShouldZeroCopy option
	ShouldZeroCopy bool
	// TODO: ImpliedVersion option
	ImpliedVersion uint
	// TODO: ImpliedTLContainer option
	ImpliedTLContainer TLContainerType
}

// Decode a CTE document, sending all data events to the specified event receiver.
func CTEDecode(document []byte, eventReceiver DataEventReceiver, options *CTEDecoderOptions) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	decoder := NewCTEDecoder([]byte(document), eventReceiver, options)
	decoder.Decode()
	return
}

// Decodes CTE documents
type CTEDecoder struct {
	eventReceiver  DataEventReceiver
	document       []byte
	tokenStart     int
	tokenPos       int
	endPos         int
	containerState []cteDecoderState
	currentState   cteDecoderState
	options        CTEDecoderOptions
}

func NewCTEDecoder(document []byte, eventReceiver DataEventReceiver, options *CTEDecoderOptions) *CTEDecoder {
	_this := &CTEDecoder{}
	_this.Init(document, eventReceiver, options)
	return _this
}

func (_this *CTEDecoder) Init(document []byte, eventReceiver DataEventReceiver, options *CTEDecoderOptions) {
	_this.document = document
	_this.eventReceiver = eventReceiver
	if options != nil {
		_this.options = *options
	}
	_this.endPos = len(document) - 1
}

// Run the complete decode process. The document and data receiver specified
// when initializing the decoder will be used.
func (_this *CTEDecoder) Decode() (err error) {
	_this.currentState = cteDecoderStateAwaitObject

	// Forgive initial whitespace even though it's technically not allowed
	_this.decodeWhitespace()

	// TODO: Inline containers etc
	_this.handleVersion()

	for !_this.isEndOfDocument() {
		_this.handleNextState()
	}
	_this.eventReceiver.OnEndDocument()
	return
}

// ============================================================================

// Bytes

func (_this *CTEDecoder) getByteAt(index int) cteByte {
	return cteByte(_this.document[index])
}

func (_this *CTEDecoder) peekByteAt(offset int) cteByte {
	return _this.getByteAt(_this.tokenPos + offset)
}

func (_this *CTEDecoder) peekByte() cteByte {
	if _this.isEndOfDocument() {
		return cteByteEndOfDocument
	}
	return _this.getByteAt(_this.tokenPos)
}

func (_this *CTEDecoder) readByte() (b cteByte) {
	if _this.isEndOfDocument() {
		return cteByteEndOfDocument
	}
	b = _this.getByteAt(_this.tokenPos)
	_this.advanceByte()
	return
}

func (_this *CTEDecoder) advanceByte() {
	_this.tokenPos++
}

func (_this *CTEDecoder) readUntilByte(b byte) {
	i := _this.tokenPos
	for ; i <= _this.endPos && _this.document[i] != b; i++ {
	}
	_this.tokenPos = i
}

func (_this *CTEDecoder) readUntilProperty(property cteByteProprty) {
	i := _this.tokenPos
	for ; i <= _this.endPos && !hasProperty(_this.document[i], property); i++ {
	}
	_this.tokenPos = i
}

func (_this *CTEDecoder) readWhileProperty(property cteByteProprty) {
	i := _this.tokenPos
	for ; i <= _this.endPos && hasProperty(_this.document[i], property); i++ {
	}
	_this.tokenPos = i
}

func (_this *CTEDecoder) getCharBeginIndex(index int) int {
	i := index
	// UTF-8 continuation characters have the form 10xxxxxx
	for b := _this.document[i]; b >= 0x80 && b <= 0xc0; b = _this.document[i] {
		i--
	}
	return i
}

func (_this *CTEDecoder) getCharEndIndex(index int) int {
	i := index
	// UTF-8 continuation characters have the form 10xxxxxx
	for b := _this.document[i]; b >= 0x80 && b <= 0xc0; b = _this.document[i] {
		i++
	}
	return i
}

func (_this *CTEDecoder) ungetByte() {
	_this.tokenPos--
}
func (_this *CTEDecoder) ungetAll() {
	_this.tokenPos = _this.tokenStart
}

func (_this *CTEDecoder) isEndOfDocument() bool {
	return _this.tokenPos > _this.endPos
}

// Tokens

func (_this *CTEDecoder) endToken() {
	_this.tokenStart = _this.tokenPos
}

func (_this *CTEDecoder) endObject() {
	_this.endToken()
	_this.transitionToNextState()
}

// State

func (_this *CTEDecoder) stackContainer(newState cteDecoderState) {
	_this.containerState = append(_this.containerState, _this.currentState)
	_this.currentState = newState
}

func (_this *CTEDecoder) unstackContainer() {
	index := len(_this.containerState) - 1
	_this.currentState = _this.containerState[index]
	_this.containerState = _this.containerState[:index]
}

func (_this *CTEDecoder) changeState(newState cteDecoderState) {
	_this.currentState = newState
}

func (_this *CTEDecoder) transitionToNextState() {
	_this.currentState = cteDecoderStateTransitions[_this.currentState]
}

// Errors

func (_this *CTEDecoder) assertNotEOD() {
	if _this.isEndOfDocument() {
		_this.unexpectedEOD()
	}
}

func (_this *CTEDecoder) assertAtObjectEnd(decoding string) {
	if !_this.peekByte().HasProperty(ctePropertyObjectEnd) {
		_this.unexpectedChar(decoding)
	}
}

func (_this *CTEDecoder) unexpectedEOD() {
	_this.errorf("Unexpected end of document")
}

func (_this *CTEDecoder) errorAt(index int, format string, args ...interface{}) {
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
		b := _this.getByteAt(i)
		if b < 0x80 || b > 0xc0 {
			colNumber++
		}
	}

	msg := fmt.Sprintf(format, args...)
	panic(fmt.Errorf("Offset %v (line %v, col %v): %v", index, lineNumber, colNumber, msg))
	_this.errorf(format, args...)
}

func (_this *CTEDecoder) errorf(format string, args ...interface{}) {
	_this.errorAt(_this.tokenPos, format, args...)
}

func (_this *CTEDecoder) unexpectedCharAt(index int, decoding string) {
	_this.errorAt(index, "Unexpected [%v] while decoding %v", _this.describeCharAt(index), decoding)
}

func (_this *CTEDecoder) unexpectedChar(decoding string) {
	_this.unexpectedCharAt(_this.tokenPos, decoding)
}

func (_this *CTEDecoder) describeCharAt(index int) string {
	if index > _this.endPos {
		return "EOD"
	}

	charStart := _this.getCharBeginIndex(index)
	charEnd := _this.getCharEndIndex(index)
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

func (_this *CTEDecoder) describeCurrentChar() string {
	return _this.describeCharAt(_this.tokenPos)
}

// Decoders

func (_this *CTEDecoder) decodeWhitespace() {
	endPos := _this.endPos
	i := _this.tokenPos
	for ; i <= endPos; i++ {
		if !hasProperty(_this.document[i], ctePropertyWhitespace) {
			break
		}
	}
	_this.tokenPos = i
	_this.endToken()
}

func (_this *CTEDecoder) decodeBinaryInteger() (value uint64) {
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
			_this.errorf("Overflow reading binary integer")
		}
	}
	_this.tokenPos = i
	return
}

func (_this *CTEDecoder) decodeOctalInteger() (value uint64) {
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
			_this.errorf("Overflow reading octal value")
		}
	}
	_this.tokenPos = i
	return
}

// startValue is only used if bigStartValue is nil
// bigValue will be nil unless the value was too big for a uint64
func (_this *CTEDecoder) decodeDecimalInteger(startValue uint64, bigStartValue *big.Int) (value uint64, bigValue *big.Int, digitCount int) {
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
			bigValue = bigValue.Mul(bigValue, bigInt10)
			bigValue = bigValue.Add(bigValue, big.NewInt(int64(b-'0')))
		}
	}

	digitCount = i - _this.tokenPos
	_this.tokenPos = i
	return
}

func (_this *CTEDecoder) decodeHexInteger(startValue uint64) (value uint64, digitCount int) {
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
			_this.errorf("Overflow reading hex value")
		}
	}
	digitCount = i - _this.tokenPos
	_this.tokenPos = i
	return
}

func (_this *CTEDecoder) decodeQuotedStringWithEscapes(startPos, firstEscape int) string {
	_this.errorf("TODO: CTEDecoder: Escape sequences")
	return ""
}

func (_this *CTEDecoder) decodeQuotedString() string {
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
			return _this.decodeQuotedStringWithEscapes(startPos, i)
		}
	}
	_this.unexpectedEOD()
	return ""
}

func (_this *CTEDecoder) decodeHexBytes() []byte {
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
				_this.errorAt(i, "Missing last hex digit")
			}
			_this.tokenPos = i + 1
			return bytes
		default:
			_this.unexpectedCharAt(i, "hex encoding")
		}
		if !firstNybble {
			bytes = append(bytes, nextByte)
			nextByte = 0
		}
		nextByte <<= 4
		firstNybble = !firstNybble
	}
	_this.unexpectedEOD()
	return nil
}

func (_this *CTEDecoder) decodeUUID() []byte {
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
			_this.unexpectedCharAt(i, "UUID")
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
		_this.errorAt(i, "Unrecognized named value or malformed UUID")
	}

	_this.tokenPos = i
	return uuid
}

func (_this *CTEDecoder) decodeDate(year int64) *compact_time.Time {
	month, _, digitCount := _this.decodeDecimalInteger(0, nil)
	if digitCount > 2 {
		_this.errorf("Month field is too long")
	}
	if _this.peekByte() != '-' {
		_this.unexpectedChar("month")
	}
	_this.advanceByte()

	var day uint64
	day, _, digitCount = _this.decodeDecimalInteger(0, nil)
	if digitCount == 0 {
		_this.unexpectedChar("day")
	}
	if digitCount > 2 {
		_this.errorf("Day field is too long")
	}
	if _this.peekByte() != '/' {
		return compact_time.NewDate(int(year), int(month), int(day))
	}

	_this.advanceByte()
	var hour uint64
	hour, _, digitCount = _this.decodeDecimalInteger(0, nil)
	if digitCount == 0 {
		_this.unexpectedChar("hour")
	}
	if digitCount > 2 {
		_this.errorf("Hour field is too long")
	}
	if _this.readByte() != ':' {
		_this.ungetByte()
		_this.unexpectedChar("hour")
	}
	t := _this.decodeTime(int(hour))
	if t.TimezoneIs == compact_time.TypeLatitudeLongitude {
		return compact_time.NewTimestampLatLong(int(year), int(month), int(day),
			int(t.Hour), int(t.Minute), int(t.Second), int(t.Nanosecond),
			int(t.LatitudeHundredths), int(t.LongitudeHundredths))
	}
	return compact_time.NewTimestamp(int(year), int(month), int(day),
		int(t.Hour), int(t.Minute), int(t.Second), int(t.Nanosecond),
		t.AreaLocation)
}

func (_this *CTEDecoder) decodeTime(hour int) *compact_time.Time {
	minute, _, digitCount := _this.decodeDecimalInteger(0, nil)
	if digitCount > 2 {
		_this.errorf("Minute field is too long")
	}
	if _this.peekByte() != ':' {
		_this.unexpectedChar("minute")
	}
	_this.advanceByte()
	var second uint64
	second, _, digitCount = _this.decodeDecimalInteger(0, nil)
	if digitCount > 2 {
		_this.errorf("Second field is too long")
	}
	var nanosecond int

	if _this.peekByte() == '.' {
		_this.advanceByte()
		v, _, digitCount := _this.decodeDecimalInteger(0, nil)
		if digitCount == 0 {
			_this.unexpectedChar("nanosecond")
		}
		if digitCount > 9 {
			_this.errorf("Nanosecond field is too long")
		}
		nanosecond = int(v)
		nanosecond *= subsecondMagnitudes[digitCount]
	}

	b := _this.peekByte()
	if b.HasProperty(ctePropertyObjectEnd) {
		return compact_time.NewTime(hour, int(minute), int(second), nanosecond, "")
	}

	if b != '/' {
		_this.unexpectedChar("time")
	}
	_this.advanceByte()

	if _this.peekByte().HasProperty(ctePropertyAZ) {
		areaLocationStart := _this.tokenPos
		_this.readWhileProperty(ctePropertyAreaLocation)
		if _this.peekByte() == '/' {
			_this.advanceByte()
			_this.readWhileProperty(ctePropertyAreaLocation)
		}
		areaLocation := string(_this.document[areaLocationStart:_this.tokenPos])
		return compact_time.NewTime(hour, int(minute), int(second), nanosecond, areaLocation)
	}

	lat, long := _this.decodeLatLong()
	return compact_time.NewTimeLatLong(hour, int(minute), int(second), nanosecond, lat, long)
}

func (_this *CTEDecoder) decodeLatLongPortion(name string) (value int) {
	whole, _, digitCount := _this.decodeDecimalInteger(0, nil)
	switch digitCount {
	case 1, 2, 3:
	// Nothing to do
	case 0:
		_this.unexpectedChar(name)
	default:
		_this.errorf("Too many digits decoding %v", name)
	}

	var fractional uint64
	b := _this.peekByte()
	if b == '.' {
		_this.advanceByte()
		fractional, _, digitCount = _this.decodeDecimalInteger(0, nil)
		switch digitCount {
		case 1:
			fractional *= 10
		case 2:
			// Nothing to do
		case 0:
			_this.unexpectedChar(name)
		default:
			_this.errorf("Too many digits decoding %v", name)
		}
	}
	return int(whole*100 + fractional)
}

func (_this *CTEDecoder) decodeLatLong() (latitudeHundredths, longitudeHundredths int) {
	latitudeHundredths = _this.decodeLatLongPortion("latitude")

	if _this.peekByte() != '/' {
		_this.unexpectedChar("latitude/longitude")
	}
	_this.advanceByte()

	longitudeHundredths = _this.decodeLatLongPortion("longitude")

	return
}

func (_this *CTEDecoder) decodeDecimalFloat(sign int64, coefficient uint64, bigCoefficient *big.Int, coefficientDigitCount int) (value compact_float.DFloat, bigValue *apd.Decimal) {
	exponent := int32(0)
	fractionalDigitCount := 0
	coefficient, bigCoefficient, fractionalDigitCount = _this.decodeDecimalInteger(coefficient, bigCoefficient)
	if fractionalDigitCount == 0 {
		_this.unexpectedChar("float fractional")
	}

	if _this.peekByte() == 'e' {
		_this.advanceByte()
		exponentSign := int32(1)
		switch _this.peekByte() {
		case '+':
			_this.advanceByte()
		case '-':
			exponentSign = -1
			_this.advanceByte()
		}
		exp, bigExp, digitCount := _this.decodeDecimalInteger(0, nil)
		if digitCount == 0 {
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

	value = compact_float.DFloatValue(int32(exponent), int64(coefficient)*int64(sign))
	return
}

func (_this *CTEDecoder) decodeHexFloat(sign int64, coefficient uint64, coefficientDigitCount int) float64 {
	exponent := 0
	fractionalDigitCount := 0
	coefficient, fractionalDigitCount = _this.decodeHexInteger(coefficient)
	if fractionalDigitCount == 0 {
		_this.unexpectedChar("float fractional")
	}

	if _this.peekByte() == 'p' {
		_this.advanceByte()
		exponentSign := 1
		switch _this.peekByte() {
		case '+':
			_this.advanceByte()
		case '-':
			exponentSign = -1
			_this.advanceByte()
		}
		exp, bigExp, digitCount := _this.decodeDecimalInteger(0, nil)
		if digitCount == 0 {
			_this.unexpectedChar("float exponent")
		}
		if bigExp != nil {
			_this.errorf("Exponent too big")
		}
		exponent = int(exp) * exponentSign
	}

	exponent -= fractionalDigitCount * 4

	// TODO: Overflow

	return float64(sign) * float64(coefficient) * math.Pow(float64(2), float64(exponent))
}

func (_this *CTEDecoder) decodeSingleLineComment() string {
	commentStart := _this.tokenPos
	_this.readUntilByte('\n')
	commentEnd := _this.tokenPos
	if _this.document[commentEnd-1] == '\r' {
		commentEnd--
	}

	return string(_this.document[commentStart:commentEnd])
}

// Handlers

type cteDecoderHandlerFunction func(*CTEDecoder)

func (_this *CTEDecoder) handleNothing() {
}

func (_this *CTEDecoder) handleNextState() {
	cteDecoderStateHandlers[_this.currentState](_this)
}

func (_this *CTEDecoder) handleObject() {
	charBasedHandlers[_this.peekByte()](_this)
}

func (_this *CTEDecoder) handleInvalidChar() {
	_this.errorf("Unexpected [%v]", _this.describeCurrentChar())
}

func (_this *CTEDecoder) handleInvalidState() {
	_this.errorf("BUG: Invalid state: %v", _this.currentState)
}

func (_this *CTEDecoder) handleKVSeparator() {
	_this.decodeWhitespace()
	b := _this.peekByte()
	if b != '=' {
		_this.errorf("Expected map separator (=) but got [%v]", _this.describeCurrentChar())
	}
	_this.advanceByte()
	_this.decodeWhitespace()
	_this.endObject()
}

func (_this *CTEDecoder) handleWhitespace() {
	_this.decodeWhitespace()
	_this.endToken()
}

func (_this *CTEDecoder) handleVersion() {
	if b := _this.peekByte(); b != 'c' && b != 'C' {
		_this.errorf(`Expected document to begin with "c" but got [%v]`, _this.describeCurrentChar())
	}
	_this.advanceByte()

	version, bigVersion, digitCount := _this.decodeDecimalInteger(0, nil)
	if digitCount == 0 {
		_this.unexpectedChar("version number")
	}
	if bigVersion != nil {
		_this.errorf("Version too big")
	}

	if !_this.peekByte().HasProperty(ctePropertyWhitespace) {
		_this.unexpectedChar("whitespace after version")
	}
	_this.advanceByte()

	_this.eventReceiver.OnVersion(version)
	_this.endToken()
}

func (_this *CTEDecoder) handleStringish() {
	startPos := _this.tokenPos
	endPos := _this.endPos
	i := startPos
	var b byte
	for ; i <= endPos; i++ {
		b = _this.document[i]
		if !hasProperty(b, ctePropertyUnquotedMid) {
			break
		}
	}

	// Unquoted string
	if hasProperty(b, ctePropertyObjectEnd) || i > endPos {
		bytes := _this.document[startPos:i]
		_this.tokenPos = i
		_this.eventReceiver.OnString(string(bytes))
		_this.endObject()
		return
	}

	// Bytes, Custom, URI
	if b == '"' && i-startPos == 1 {
		initiator := _this.document[startPos]
		switch initiator {
		case 'b':
			_this.tokenPos = i + 1
			_this.eventReceiver.OnBytes(_this.decodeHexBytes())
			_this.endObject()
			return
		case 'c':
			_this.tokenPos = i + 1
			_this.eventReceiver.OnCustom(_this.decodeHexBytes())
			_this.endObject()
			return
		case 'u':
			_this.tokenPos = i + 1
			_this.eventReceiver.OnURI(_this.decodeQuotedString())
			_this.endObject()
			return
		default:
			_this.unexpectedChar("byte array initiator")
		}
	}

	_this.unexpectedChar("unquoted string")
}

func (_this *CTEDecoder) handleQuotedString() {
	_this.advanceByte()
	_this.eventReceiver.OnString(_this.decodeQuotedString())
	_this.endObject()
}

func (_this *CTEDecoder) handlePositiveNumeric() {
	coefficient, bigCoefficient, digitCount := _this.decodeDecimalInteger(0, nil)

	// Integer
	if _this.peekByte().HasProperty(ctePropertyObjectEnd) {
		if bigCoefficient != nil {
			_this.eventReceiver.OnBigInt(bigCoefficient)
		} else {
			_this.eventReceiver.OnPositiveInt(coefficient)
		}
		_this.endObject()
		return
	}

	switch _this.peekByte() {
	case '-':
		_this.advanceByte()
		v := _this.decodeDate(int64(coefficient))
		_this.assertAtObjectEnd("date")
		_this.eventReceiver.OnCompactTime(v)
		_this.endObject()
	case ':':
		_this.advanceByte()
		v := _this.decodeTime(int(coefficient))
		_this.assertAtObjectEnd("time")
		_this.eventReceiver.OnCompactTime(v)
		_this.endObject()
	case '.':
		_this.advanceByte()
		value, bigValue := _this.decodeDecimalFloat(1, coefficient, bigCoefficient, digitCount)
		_this.assertAtObjectEnd("float")
		if bigValue != nil {
			_this.eventReceiver.OnBigDecimalFloat(bigValue)
		} else {
			_this.eventReceiver.OnDecimalFloat(value)
		}
		_this.endObject()
	default:
		_this.unexpectedChar("numeric")
	}
}

func (_this *CTEDecoder) handleNegativeNumeric() {
	_this.advanceByte()
	switch _this.peekByte() {
	case '0':
		_this.handleOtherBaseNegative()
		return
	case '@':
		_this.advanceByte()
		nameStart := _this.tokenPos
		_this.readWhileProperty(ctePropertyAZ)
		token := strings.ToLower(string(_this.document[nameStart:_this.tokenPos]))
		if token != "inf" {
			_this.errorf("Unknown named value: %v", token)
		}
		_this.eventReceiver.OnFloat(math.Inf(-1))
		return
	}

	coefficient, bigCoefficient, digitCount := _this.decodeDecimalInteger(0, nil)

	// Integer
	if _this.peekByte().HasProperty(ctePropertyObjectEnd) {
		if bigCoefficient != nil {
			// TODO: More efficient way to negate?
			bigCoefficient = bigCoefficient.Mul(bigCoefficient, bigIntN1)
			_this.eventReceiver.OnBigInt(bigCoefficient)
		} else {
			_this.eventReceiver.OnNegativeInt(coefficient)
		}
		_this.endObject()
		return
	}

	switch _this.peekByte() {
	case '-':
		_this.advanceByte()
		v := _this.decodeDate(-int64(coefficient))
		_this.assertAtObjectEnd("time")
		_this.eventReceiver.OnCompactTime(v)
		_this.endObject()
	case '.':
		_this.advanceByte()
		value, bigValue := _this.decodeDecimalFloat(-1, coefficient, bigCoefficient, digitCount)
		_this.assertAtObjectEnd("float")
		if bigValue != nil {
			_this.eventReceiver.OnBigDecimalFloat(bigValue)
		} else {
			_this.eventReceiver.OnDecimalFloat(value)
		}
		_this.endObject()
	default:
		_this.unexpectedChar("numeric")
	}
}

func (_this *CTEDecoder) handleOtherBasePositive() {
	_this.advanceByte()
	b := _this.peekByte()

	if b.HasProperty(ctePropertyObjectEnd) {
		_this.eventReceiver.OnPositiveInt(0)
		_this.endObject()
		return
	}
	_this.advanceByte()

	switch b {
	case 'b':
		v := _this.decodeBinaryInteger()
		_this.assertAtObjectEnd("binary integer")
		_this.eventReceiver.OnPositiveInt(v)
		_this.endObject()
	case 'o':
		v := _this.decodeOctalInteger()
		_this.assertAtObjectEnd("octal integer")
		_this.eventReceiver.OnPositiveInt(v)
		_this.endObject()
	case 'x':
		v, digitCount := _this.decodeHexInteger(0)
		if _this.peekByte() == '.' {
			_this.advanceByte()
			fv := _this.decodeHexFloat(1, v, digitCount)
			_this.assertAtObjectEnd("hex float")
			_this.eventReceiver.OnFloat(fv)
			_this.endObject()
		} else {
			_this.assertAtObjectEnd("hex integer")
			_this.eventReceiver.OnPositiveInt(v)
			_this.endObject()
		}
	case '.':
		value, bigValue := _this.decodeDecimalFloat(1, 0, nil, 0)
		_this.assertAtObjectEnd("float")
		if bigValue != nil {
			_this.eventReceiver.OnBigDecimalFloat(bigValue)
		} else {
			_this.eventReceiver.OnDecimalFloat(value)
		}
		_this.endObject()
	default:
		if b.HasProperty(cteProperty09) && _this.peekByte() == ':' {
			_this.advanceByte()
			v := _this.decodeTime(int(b - '0'))
			_this.assertAtObjectEnd("time")
			_this.eventReceiver.OnCompactTime(v)
			_this.endObject()
			return
		}
		_this.ungetByte()
		_this.unexpectedChar("numeric base")
	}
}

func (_this *CTEDecoder) handleOtherBaseNegative() {
	_this.advanceByte()
	b := _this.readByte()
	switch b {
	case 'b':
		v := _this.decodeBinaryInteger()
		_this.assertAtObjectEnd("binary integer")
		_this.eventReceiver.OnNegativeInt(v)
		_this.endObject()
	case 'o':
		v := _this.decodeOctalInteger()
		_this.assertAtObjectEnd("octal integer")
		_this.eventReceiver.OnNegativeInt(v)
		_this.endObject()
	case 'x':
		v, digitCount := _this.decodeHexInteger(0)
		if _this.peekByte() == '.' {
			_this.advanceByte()
			fv := _this.decodeHexFloat(-1, v, digitCount)
			_this.assertAtObjectEnd("hex float")
			_this.eventReceiver.OnFloat(fv)
			_this.endObject()
		} else {
			_this.assertAtObjectEnd("hex integer")
			_this.eventReceiver.OnNegativeInt(v)
			_this.endObject()
		}
	case '.':
		value, bigValue := _this.decodeDecimalFloat(-1, 0, nil, 0)
		_this.assertAtObjectEnd("float")
		if bigValue != nil {
			_this.eventReceiver.OnBigDecimalFloat(bigValue)
		} else {
			_this.eventReceiver.OnDecimalFloat(value)
		}
		_this.endObject()
	default:
		_this.ungetByte()
		_this.unexpectedChar("numeric base")
	}
}

func (_this *CTEDecoder) handleListBegin() {
	_this.advanceByte()
	_this.eventReceiver.OnList()
	_this.stackContainer(cteDecoderStateAwaitListItem)
	_this.endToken()
}

func (_this *CTEDecoder) handleListEnd() {
	_this.advanceByte()
	_this.eventReceiver.OnEnd()
	_this.unstackContainer()
	_this.endObject()
}

func (_this *CTEDecoder) handleMapBegin() {
	_this.advanceByte()
	_this.eventReceiver.OnMap()
	_this.stackContainer(cteDecoderStateAwaitMapKey)
	_this.endToken()
}

func (_this *CTEDecoder) handleMapEnd() {
	_this.advanceByte()
	_this.eventReceiver.OnEnd()
	_this.unstackContainer()
	_this.endObject()
}

func (_this *CTEDecoder) handleMetadataBegin() {
	_this.advanceByte()
	_this.eventReceiver.OnMetadata()
	_this.stackContainer(cteDecoderStateAwaitMetaKey)
	_this.endToken()
}

func (_this *CTEDecoder) handleMetadataEnd() {
	_this.advanceByte()
	_this.eventReceiver.OnEnd()
	_this.unstackContainer()
	_this.endToken()
	// Don't transition state because metadata is a pseudo-object
}

func (_this *CTEDecoder) handleComment() {
	_this.readByte()
	switch _this.readByte() {
	case '/':
		_this.eventReceiver.OnComment()
		contents := _this.decodeSingleLineComment()
		if len(contents) > 0 {
			_this.eventReceiver.OnString(contents)
		}
		_this.eventReceiver.OnEnd()
		_this.endToken()
		// Don't transition state because a comment is a pseudo-object
	case '*':
		_this.eventReceiver.OnComment()
		_this.stackContainer(cteDecoderStateAwaitCommentItem)
		_this.endToken()
	default:
		_this.ungetByte()
		_this.unexpectedChar("comment")
	}
}

func (_this *CTEDecoder) handleCommentContent() {
	startPos := _this.tokenPos
	endPos := _this.endPos

	// Skip first byte so that index-1 is not comparing against part of the
	// original multiline comment initiator '*'
	_this.advanceByte()

	for i := _this.tokenPos; i <= endPos; i++ {
		bLast := _this.document[i-1]
		bNow := _this.document[i]

		if bLast == '*' && bNow == '/' {
			endOffset := i - 1
			if endOffset > startPos {
				str := string(_this.document[startPos:endOffset])
				_this.eventReceiver.OnString(str)
			}
			_this.tokenPos = i + 1
			_this.eventReceiver.OnEnd()
			_this.unstackContainer()
			_this.endToken()
			return
		}

		if bLast == '/' && bNow == '*' {
			endOffset := i - 1
			if endOffset > _this.tokenStart {
				str := string(_this.document[startPos:endOffset])
				_this.eventReceiver.OnString(str)
			}
			_this.tokenPos = i + 1
			_this.eventReceiver.OnComment()
			_this.stackContainer(cteDecoderStateAwaitCommentItem)
			_this.endToken()
			return
		}
	}

	_this.unexpectedEOD()
}

func (_this *CTEDecoder) handleMarkupBegin() {
	_this.advanceByte()
	_this.eventReceiver.OnMarkup()
	_this.stackContainer(cteDecoderStateAwaitMarkupValue)
	_this.endToken()
}

func (_this *CTEDecoder) handleMarkupContentBegin() {
	_this.advanceByte()
	_this.eventReceiver.OnEnd()
	_this.changeState(cteDecoderStateAwaitMarkupItem)
	_this.endToken()
}

func (_this *CTEDecoder) handleMarkupWithEscapes(startPos, firstEscape int) string {
	_this.errorf("TODO: CTEDecoder: Markup with escape sequences, entity refs")
	return ""
}

func (_this *CTEDecoder) handleMarkupContent() {
	startPos := _this.tokenPos
	endPos := _this.endPos
	i := startPos
	for ; i <= endPos; i++ {
		switch _this.document[i] {
		case '/':
			switch _this.getByteAt(i + 1) {
			case '/':
				if i > startPos {
					_this.eventReceiver.OnString(string(_this.document[startPos:i]))
				}
				_this.tokenPos = i + 2
				_this.eventReceiver.OnComment()
				contents := _this.decodeSingleLineComment()
				if len(contents) > 0 {
					_this.eventReceiver.OnString(contents)
				}
				_this.advanceByte()
				_this.eventReceiver.OnEnd()
				_this.endToken()
				// Don't transition state because a comment is a pseudo-object
				return
			case '*':
				if i > startPos {
					_this.eventReceiver.OnString(string(_this.document[startPos:i]))
				}
				_this.tokenPos = i + 2
				_this.eventReceiver.OnComment()
				_this.stackContainer(cteDecoderStateAwaitCommentItem)
				_this.endToken()
				return
			}
		case '<':
			str := string(_this.document[startPos:i])
			_this.tokenPos = i + 1
			if len(str) > 0 {
				_this.eventReceiver.OnString(str)
			}

			_this.tokenPos = i
			_this.handleMarkupBegin()
			return
		case '>':
			str := string(_this.document[startPos:i])
			_this.tokenPos = i + 1
			if len(str) > 0 {
				_this.eventReceiver.OnString(str)
			}
			_this.eventReceiver.OnEnd()
			_this.unstackContainer()
			_this.endObject()
			return
		case '\\':
			_this.handleMarkupWithEscapes(startPos, i)
			return
		}
	}
	_this.unexpectedEOD()
}

func (_this *CTEDecoder) handleMarkupEnd() {
	_this.advanceByte()
	_this.eventReceiver.OnEnd()
	_this.eventReceiver.OnEnd()
	_this.unstackContainer()
	_this.endObject()
}

func (_this *CTEDecoder) handleNamedValue() {
	_this.advanceByte()
	nameStart := _this.tokenPos
	_this.readWhileProperty(ctePropertyAZ)
	token := strings.ToLower(string(_this.document[nameStart:_this.tokenPos]))
	switch token {
	case "nil":
		_this.eventReceiver.OnNil()
		_this.endObject()
		return
	case "nan":
		_this.eventReceiver.OnNan(false)
		_this.endObject()
		return
	case "snan":
		_this.eventReceiver.OnNan(true)
		_this.endObject()
		return
	case "inf":
		_this.eventReceiver.OnFloat(math.Inf(1))
		_this.endObject()
		return
	case "false":
		_this.eventReceiver.OnFalse()
		_this.endObject()
		return
	case "true":
		_this.eventReceiver.OnTrue()
		_this.endObject()
		return
	}

	_this.ungetAll()
	_this.advanceByte()
	_this.eventReceiver.OnUUID(_this.decodeUUID())
	_this.endObject()
}

func (_this *CTEDecoder) handleVerbatimString() {
	_this.advanceByte()
	_this.assertNotEOD()
	sentinelStart := _this.tokenPos
	_this.readUntilProperty(ctePropertyWhitespace)
	sentinelEnd := _this.tokenPos
	wsByte := _this.readByte()
	if wsByte == '\r' {
		if _this.readByte() != '\n' {
			_this.unexpectedChar("verbatim sentinel")
		}
	}
	sentinel := _this.document[sentinelStart:sentinelEnd]
	searchSpace := _this.document[_this.tokenPos:]
	index := bytes.Index(searchSpace, sentinel)
	if index < 0 {
		_this.errorf("Verbatim sentinel sequence [%v] not found in document", string(sentinel))
	}
	str := string(searchSpace[:index])
	_this.tokenPos += index + len(sentinel)
	_this.assertAtObjectEnd("verbatim string")
	_this.eventReceiver.OnString(str)
	_this.endObject()
}

func (_this *CTEDecoder) handleReference() {
	_this.advanceByte()
	_this.eventReceiver.OnReference()
	if _this.peekByte().HasProperty(ctePropertyWhitespace) {
		_this.errorf("Whitespace not allowed between reference and tag name")
	}
	_this.endToken()
}

func (_this *CTEDecoder) handleMarker() {
	_this.advanceByte()
	_this.eventReceiver.OnMarker()
	if _this.peekByte().HasProperty(ctePropertyWhitespace) {
		_this.errorf("Whitespace not allowed between marker and tag name")
	}
	_this.endToken()
}

type cteDecoderState int

const (
	cteDecoderStateAwaitObject cteDecoderState = iota
	cteDecoderStateAwaitListItem
	cteDecoderStateAwaitCommentItem
	cteDecoderStateAwaitMapKey
	cteDecoderStateAwaitMapKVSeparator
	cteDecoderStateAwaitMapValue
	cteDecoderStateAwaitMetaKey
	cteDecoderStateAwaitMetaKVSeparator
	cteDecoderStateAwaitMetaValue
	cteDecoderStateAwaitMarkupName
	cteDecoderStateAwaitMarkupKey
	cteDecoderStateAwaitMarkupKVSeparator
	cteDecoderStateAwaitMarkupValue
	cteDecoderStateAwaitMarkupItem
	cteDecoderStateAwaitMarkerID
	cteDecoderStateAwaitReferenceID
	cteDecoderStateCount
)

var cteDecoderStateTransitions [cteDecoderStateCount]cteDecoderState
var cteDecoderStateHandlers [cteDecoderStateCount]cteDecoderHandlerFunction

func init() {
	cteDecoderStateTransitions[cteDecoderStateAwaitObject] = cteDecoderStateAwaitObject
	cteDecoderStateTransitions[cteDecoderStateAwaitListItem] = cteDecoderStateAwaitListItem
	cteDecoderStateTransitions[cteDecoderStateAwaitCommentItem] = cteDecoderStateAwaitCommentItem
	cteDecoderStateTransitions[cteDecoderStateAwaitMapKey] = cteDecoderStateAwaitMapKVSeparator
	cteDecoderStateTransitions[cteDecoderStateAwaitMapKVSeparator] = cteDecoderStateAwaitMapValue
	cteDecoderStateTransitions[cteDecoderStateAwaitMapValue] = cteDecoderStateAwaitMapKey
	cteDecoderStateTransitions[cteDecoderStateAwaitMetaKey] = cteDecoderStateAwaitMetaKVSeparator
	cteDecoderStateTransitions[cteDecoderStateAwaitMetaKVSeparator] = cteDecoderStateAwaitMetaValue
	cteDecoderStateTransitions[cteDecoderStateAwaitMetaValue] = cteDecoderStateAwaitMetaKey
	cteDecoderStateTransitions[cteDecoderStateAwaitMarkupName] = cteDecoderStateAwaitMarkupKey
	cteDecoderStateTransitions[cteDecoderStateAwaitMarkupKey] = cteDecoderStateAwaitMarkupKVSeparator
	cteDecoderStateTransitions[cteDecoderStateAwaitMarkupKVSeparator] = cteDecoderStateAwaitMarkupValue
	cteDecoderStateTransitions[cteDecoderStateAwaitMarkupValue] = cteDecoderStateAwaitMarkupKey
	cteDecoderStateTransitions[cteDecoderStateAwaitMarkupItem] = cteDecoderStateAwaitMarkupItem
	cteDecoderStateTransitions[cteDecoderStateAwaitMarkerID] = cteDecoderStateAwaitObject
	cteDecoderStateTransitions[cteDecoderStateAwaitReferenceID] = cteDecoderStateAwaitObject

	cteDecoderStateHandlers[cteDecoderStateAwaitObject] = (*CTEDecoder).handleObject
	cteDecoderStateHandlers[cteDecoderStateAwaitListItem] = (*CTEDecoder).handleObject
	cteDecoderStateHandlers[cteDecoderStateAwaitCommentItem] = (*CTEDecoder).handleCommentContent
	cteDecoderStateHandlers[cteDecoderStateAwaitMapKey] = (*CTEDecoder).handleObject
	cteDecoderStateHandlers[cteDecoderStateAwaitMapKVSeparator] = (*CTEDecoder).handleKVSeparator
	cteDecoderStateHandlers[cteDecoderStateAwaitMapValue] = (*CTEDecoder).handleObject
	cteDecoderStateHandlers[cteDecoderStateAwaitMetaKey] = (*CTEDecoder).handleObject
	cteDecoderStateHandlers[cteDecoderStateAwaitMetaKVSeparator] = (*CTEDecoder).handleKVSeparator
	cteDecoderStateHandlers[cteDecoderStateAwaitMetaValue] = (*CTEDecoder).handleObject
	cteDecoderStateHandlers[cteDecoderStateAwaitMarkupName] = (*CTEDecoder).handleObject
	cteDecoderStateHandlers[cteDecoderStateAwaitMarkupKey] = (*CTEDecoder).handleObject
	cteDecoderStateHandlers[cteDecoderStateAwaitMarkupKVSeparator] = (*CTEDecoder).handleKVSeparator
	cteDecoderStateHandlers[cteDecoderStateAwaitMarkupValue] = (*CTEDecoder).handleObject
	cteDecoderStateHandlers[cteDecoderStateAwaitMarkupItem] = (*CTEDecoder).handleMarkupContent
	cteDecoderStateHandlers[cteDecoderStateAwaitMarkerID] = (*CTEDecoder).handleObject
	cteDecoderStateHandlers[cteDecoderStateAwaitReferenceID] = (*CTEDecoder).handleObject
}

var charBasedHandlers [cteByteEndOfDocument + 1]cteDecoderHandlerFunction

func init() {
	for i := 0; i < cteByteEndOfDocument; i++ {
		charBasedHandlers[i] = (*CTEDecoder).handleInvalidChar
	}

	charBasedHandlers['\r'] = (*CTEDecoder).handleWhitespace
	charBasedHandlers['\n'] = (*CTEDecoder).handleWhitespace
	charBasedHandlers['\t'] = (*CTEDecoder).handleWhitespace
	charBasedHandlers[' '] = (*CTEDecoder).handleWhitespace

	charBasedHandlers['!'] = (*CTEDecoder).handleInvalidChar
	charBasedHandlers['"'] = (*CTEDecoder).handleQuotedString
	charBasedHandlers['#'] = (*CTEDecoder).handleReference
	charBasedHandlers['$'] = (*CTEDecoder).handleInvalidChar
	charBasedHandlers['%'] = (*CTEDecoder).handleInvalidChar
	charBasedHandlers['&'] = (*CTEDecoder).handleMarker
	charBasedHandlers['\''] = (*CTEDecoder).handleInvalidChar
	charBasedHandlers['('] = (*CTEDecoder).handleMetadataBegin
	charBasedHandlers[')'] = (*CTEDecoder).handleMetadataEnd
	charBasedHandlers['+'] = (*CTEDecoder).handleInvalidChar
	charBasedHandlers[','] = (*CTEDecoder).handleInvalidChar
	charBasedHandlers['-'] = (*CTEDecoder).handleNegativeNumeric
	charBasedHandlers['.'] = (*CTEDecoder).handleInvalidChar
	charBasedHandlers['/'] = (*CTEDecoder).handleComment

	charBasedHandlers['0'] = (*CTEDecoder).handleOtherBasePositive
	for i := '1'; i <= '9'; i++ {
		charBasedHandlers[i] = (*CTEDecoder).handlePositiveNumeric
	}

	charBasedHandlers[':'] = (*CTEDecoder).handleInvalidChar
	charBasedHandlers[';'] = (*CTEDecoder).handleInvalidChar
	charBasedHandlers['<'] = (*CTEDecoder).handleMarkupBegin
	charBasedHandlers['>'] = (*CTEDecoder).handleMarkupEnd
	charBasedHandlers['?'] = (*CTEDecoder).handleInvalidChar
	charBasedHandlers['@'] = (*CTEDecoder).handleNamedValue

	for i := 'A'; i <= 'Z'; i++ {
		charBasedHandlers[i] = (*CTEDecoder).handleStringish
	}

	charBasedHandlers['['] = (*CTEDecoder).handleListBegin
	charBasedHandlers['\\'] = (*CTEDecoder).handleInvalidChar
	charBasedHandlers[']'] = (*CTEDecoder).handleListEnd
	charBasedHandlers['^'] = (*CTEDecoder).handleInvalidChar
	charBasedHandlers['_'] = (*CTEDecoder).handleStringish
	charBasedHandlers['`'] = (*CTEDecoder).handleVerbatimString

	for i := 'a'; i <= 'z'; i++ {
		charBasedHandlers[i] = (*CTEDecoder).handleStringish
	}

	charBasedHandlers['{'] = (*CTEDecoder).handleMapBegin
	charBasedHandlers['|'] = (*CTEDecoder).handleMarkupContentBegin
	charBasedHandlers['}'] = (*CTEDecoder).handleMapEnd
	charBasedHandlers['~'] = (*CTEDecoder).handleInvalidChar

	for i := 0xc0; i < 0xf8; i++ {
		charBasedHandlers[i] = (*CTEDecoder).handleStringish
	}

	charBasedHandlers[cteByteEndOfDocument] = (*CTEDecoder).handleNothing
}

// Byte Properties

type cteByte int

func (_this cteByte) HasProperty(property cteByteProprty) bool {
	return cteByteProperties[_this].HasProperty(property)
}

func hasProperty(b byte, property cteByteProprty) bool {
	return cteByteProperties[b].HasProperty(property)
}

const cteByteEndOfDocument = 0x100

type cteByteProprty uint16

const (
	ctePropertyWhitespace cteByteProprty = 1 << iota
	ctePropertyObjectEnd
	ctePropertyUnquotedStart
	ctePropertyUnquotedMid
	ctePropertyAZ
	cteProperty09
	ctePropertyLowercaseAF
	ctePropertyUppercaseAF
	ctePropertyMarkupInitiator
	ctePropertyBinaryDigit
	ctePropertyOctalDigit
	ctePropertyAreaLocation
)

func (_this cteByteProprty) HasProperty(property cteByteProprty) bool {
	return _this&property != 0
}

var cteByteProperties [cteByteEndOfDocument + 1]cteByteProprty

func init() {
	cteByteProperties[' '] |= ctePropertyWhitespace
	cteByteProperties['\r'] |= ctePropertyWhitespace
	cteByteProperties['\n'] |= ctePropertyWhitespace
	cteByteProperties['\t'] |= ctePropertyWhitespace

	cteByteProperties['-'] |= ctePropertyUnquotedMid | ctePropertyAreaLocation
	cteByteProperties['+'] |= ctePropertyUnquotedMid | ctePropertyAreaLocation
	cteByteProperties['.'] |= ctePropertyUnquotedMid
	cteByteProperties[':'] |= ctePropertyUnquotedMid
	cteByteProperties['/'] |= ctePropertyUnquotedMid
	cteByteProperties['_'] |= ctePropertyUnquotedMid | ctePropertyUnquotedStart | ctePropertyAreaLocation
	for i := '0'; i <= '9'; i++ {
		cteByteProperties[i] |= ctePropertyUnquotedMid | ctePropertyAreaLocation
	}
	for i := 'a'; i <= 'z'; i++ {
		cteByteProperties[i] |= ctePropertyUnquotedMid | ctePropertyUnquotedStart | ctePropertyAZ | ctePropertyAreaLocation
	}
	for i := 'A'; i <= 'Z'; i++ {
		cteByteProperties[i] |= ctePropertyUnquotedMid | ctePropertyUnquotedStart | ctePropertyAZ | ctePropertyAreaLocation
	}
	for i := 0xc0; i < 0xf8; i++ {
		// UTF-8 initiator
		cteByteProperties[i] |= ctePropertyUnquotedMid | ctePropertyUnquotedStart
	}
	for i := 0x80; i < 0xc0; i++ {
		// UTF-8 continuation
		cteByteProperties[i] |= ctePropertyUnquotedMid
	}

	cteByteProperties['='] |= ctePropertyObjectEnd
	cteByteProperties[']'] |= ctePropertyObjectEnd
	cteByteProperties['}'] |= ctePropertyObjectEnd
	cteByteProperties[')'] |= ctePropertyObjectEnd
	cteByteProperties['>'] |= ctePropertyObjectEnd
	cteByteProperties['|'] |= ctePropertyObjectEnd
	cteByteProperties[' '] |= ctePropertyObjectEnd
	cteByteProperties['\r'] |= ctePropertyObjectEnd
	cteByteProperties['\n'] |= ctePropertyObjectEnd
	cteByteProperties['\t'] |= ctePropertyObjectEnd
	cteByteProperties[cteByteEndOfDocument] |= ctePropertyObjectEnd

	for i := '0'; i <= '9'; i++ {
		cteByteProperties[i] |= cteProperty09
	}
	for i := 'a'; i <= 'f'; i++ {
		cteByteProperties[i] |= ctePropertyLowercaseAF
	}
	for i := 'A'; i <= 'F'; i++ {
		cteByteProperties[i] |= ctePropertyUppercaseAF
	}

	for i := '0'; i <= '7'; i++ {
		cteByteProperties[i] |= ctePropertyOctalDigit
	}

	for i := '0'; i <= '1'; i++ {
		cteByteProperties[i] |= ctePropertyBinaryDigit
	}

	cteByteProperties['/'] |= ctePropertyMarkupInitiator
	cteByteProperties['<'] |= ctePropertyMarkupInitiator
	cteByteProperties['>'] |= ctePropertyMarkupInitiator
	cteByteProperties['\\'] |= ctePropertyMarkupInitiator
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
