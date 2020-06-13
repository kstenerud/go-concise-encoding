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

func CTEDecode(document []byte, eventReceiver DataEventReceiver) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	decoder := NewCTEDecoder([]byte(document), eventReceiver)
	decoder.Decode()
	return
}

type CTEDecoder struct {
	eventReceiver  DataEventReceiver
	document       []byte
	tokenStart     int
	tokenPos       int
	endPos         int
	containerState []cteDecoderState
	currentState   cteDecoderState
}

func NewCTEDecoder(document []byte, eventReceiver DataEventReceiver) *CTEDecoder {
	this := &CTEDecoder{}
	this.Init(document, eventReceiver)
	return this
}

func (this *CTEDecoder) Init(document []byte, eventReceiver DataEventReceiver) {
	this.document = document
	this.eventReceiver = eventReceiver
	this.endPos = len(document) - 1
}

func (this *CTEDecoder) Decode() (err error) {
	this.currentState = cteDecoderStateAwaitObject

	// Forgive initial whitespace even though it's technically not allowed
	this.decodeWhitespace()

	// TODO: Inline containers etc
	this.handleVersion()

	for !this.isEndOfDocument() {
		this.handleNextState()
	}
	this.eventReceiver.OnEndDocument()
	return
}

// ============================================================================

// Bytes

func (this *CTEDecoder) getByteAt(index int) cteByte {
	return cteByte(this.document[index])
}

func (this *CTEDecoder) peekByteAt(offset int) cteByte {
	return this.getByteAt(this.tokenPos + offset)
}

func (this *CTEDecoder) peekByte() cteByte {
	if this.isEndOfDocument() {
		return cteByteEndOfDocument
	}
	return this.getByteAt(this.tokenPos)
}

func (this *CTEDecoder) readByte() (b cteByte) {
	if this.isEndOfDocument() {
		return cteByteEndOfDocument
	}
	b = this.getByteAt(this.tokenPos)
	this.advanceByte()
	return
}

func (this *CTEDecoder) advanceByte() {
	this.tokenPos++
}

func (this *CTEDecoder) readUntilByte(b byte) {
	i := this.tokenPos
	for ; i <= this.endPos && this.document[i] != b; i++ {
	}
	this.tokenPos = i
}

func (this *CTEDecoder) readUntilProperty(property cteByteProprty) {
	i := this.tokenPos
	for ; i <= this.endPos && !hasProperty(this.document[i], property); i++ {
	}
	this.tokenPos = i
}

func (this *CTEDecoder) readWhileProperty(property cteByteProprty) {
	i := this.tokenPos
	for ; i <= this.endPos && hasProperty(this.document[i], property); i++ {
	}
	this.tokenPos = i
}

func (this *CTEDecoder) getCharBeginIndex(index int) int {
	i := index
	// UTF-8 continuation characters have the form 10xxxxxx
	for b := this.document[i]; b >= 0x80 && b <= 0xc0; b = this.document[i] {
		i--
	}
	return i
}

func (this *CTEDecoder) getCharEndIndex(index int) int {
	i := index
	// UTF-8 continuation characters have the form 10xxxxxx
	for b := this.document[i]; b >= 0x80 && b <= 0xc0; b = this.document[i] {
		i++
	}
	return i
}

func (this *CTEDecoder) ungetByte() {
	this.tokenPos--
}
func (this *CTEDecoder) ungetAll() {
	this.tokenPos = this.tokenStart
}

func (this *CTEDecoder) isEndOfDocument() bool {
	return this.tokenPos > this.endPos
}

// Tokens

func (this *CTEDecoder) endToken() {
	this.tokenStart = this.tokenPos
}

func (this *CTEDecoder) endObject() {
	this.endToken()
	this.transitionToNextState()
}

// State

func (this *CTEDecoder) stackContainer(newState cteDecoderState) {
	this.containerState = append(this.containerState, this.currentState)
	this.currentState = newState
}

func (this *CTEDecoder) unstackContainer() {
	index := len(this.containerState) - 1
	this.currentState = this.containerState[index]
	this.containerState = this.containerState[:index]
}

func (this *CTEDecoder) changeState(newState cteDecoderState) {
	this.currentState = newState
}

func (this *CTEDecoder) transitionToNextState() {
	this.currentState = cteDecoderStateTransitions[this.currentState]
}

// Errors

func (this *CTEDecoder) assertNotEOD() {
	if this.isEndOfDocument() {
		this.unexpectedEOD()
	}
}

func (this *CTEDecoder) assertAtObjectEnd(decoding string) {
	if !this.peekByte().HasProperty(ctePropertyObjectEnd) {
		this.unexpectedChar(decoding)
	}
}

func (this *CTEDecoder) unexpectedEOD() {
	this.errorf("Unexpected end of document")
}

func (this *CTEDecoder) errorAt(index int, format string, args ...interface{}) {
	lineNumber := 1
	lineStart := 0
	for i := 0; i < index; i++ {
		if this.document[i] == '\n' {
			lineNumber++
			lineStart = i
		}
	}

	colNumber := 1
	for i := lineStart; i < index; i++ {
		b := this.getByteAt(i)
		if b < 0x80 || b > 0xc0 {
			colNumber++
		}
	}

	msg := fmt.Sprintf(format, args...)
	panic(fmt.Errorf("Offset %v (line %v, col %v): %v", index, lineNumber, colNumber, msg))
	this.errorf(format, args...)
}

func (this *CTEDecoder) errorf(format string, args ...interface{}) {
	this.errorAt(this.tokenPos, format, args...)
}

func (this *CTEDecoder) unexpectedCharAt(index int, decoding string) {
	this.errorAt(index, "Unexpected [%v] while decoding %v", this.describeCharAt(index), decoding)
}

func (this *CTEDecoder) unexpectedChar(decoding string) {
	this.unexpectedCharAt(this.tokenPos, decoding)
}

func (this *CTEDecoder) describeCharAt(index int) string {
	if index > this.endPos {
		return "EOD"
	}

	charStart := this.getCharBeginIndex(index)
	charEnd := this.getCharEndIndex(index)
	if charEnd-charStart > 1 {
		return string(this.document[charStart:charEnd])
	}

	b := this.document[charStart]
	if b > ' ' && b <= '~' {
		return string(b)
	}
	if b == ' ' {
		return "SP"
	}
	return fmt.Sprintf("0x%02x", b)
}

func (this *CTEDecoder) describeCurrentChar() string {
	return this.describeCharAt(this.tokenPos)
}

// Decoders

func (this *CTEDecoder) decodeWhitespace() {
	endPos := this.endPos
	i := this.tokenPos
	for ; i <= endPos; i++ {
		if !hasProperty(this.document[i], ctePropertyWhitespace) {
			break
		}
	}
	this.tokenPos = i
	this.endToken()
}

func (this *CTEDecoder) decodeBinaryInteger() (value uint64) {
	endPos := this.endPos
	i := this.tokenPos
	for ; i <= endPos; i++ {
		b := this.document[i]
		if !hasProperty(b, ctePropertyBinaryDigit) {
			break
		}
		oldValue := value
		value = value<<1 + uint64(b-'0')
		if value < oldValue {
			this.errorf("Overflow reading binary integer")
		}
	}
	this.tokenPos = i
	return
}

func (this *CTEDecoder) decodeOctalInteger() (value uint64) {
	endPos := this.endPos
	i := this.tokenPos
	for ; i <= endPos; i++ {
		b := this.document[i]
		if !hasProperty(b, ctePropertyOctalDigit) {
			break
		}
		oldValue := value
		value = value<<3 + uint64(b-'0')
		if value < oldValue {
			this.errorf("Overflow reading octal value")
		}
	}
	this.tokenPos = i
	return
}

// startValue is only used if bigStartValue is nil
// bigValue will be nil unless the value was too big for a uint64
func (this *CTEDecoder) decodeDecimalInteger(startValue uint64, bigStartValue *big.Int) (value uint64, bigValue *big.Int, digitCount int) {
	endPos := this.endPos
	i := this.tokenPos

	if bigStartValue == nil {
		value = startValue
		for ; i <= endPos; i++ {
			b := this.document[i]
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
			b := this.document[i]
			if !hasProperty(b, cteProperty09) {
				break
			}
			bigValue = bigValue.Mul(bigValue, bigInt10)
			bigValue = bigValue.Add(bigValue, big.NewInt(int64(b-'0')))
		}
	}

	digitCount = i - this.tokenPos
	this.tokenPos = i
	return
}

func (this *CTEDecoder) decodeHexInteger(startValue uint64) (value uint64, digitCount int) {
	value = startValue
	endPos := this.endPos
	i := this.tokenPos
Loop:
	for ; i <= endPos; i++ {
		b := this.document[i]
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
			this.errorf("Overflow reading hex value")
		}
	}
	digitCount = i - this.tokenPos
	this.tokenPos = i
	return
}

func (this *CTEDecoder) decodeQuotedStringWithEscapes(startPos, firstEscape int) string {
	this.errorf("TODO: CTEDecoder: Escape sequences")
	return ""
}

func (this *CTEDecoder) decodeQuotedString() string {
	startPos := this.tokenPos
	endPos := this.endPos
	i := startPos
	for ; i <= endPos; i++ {
		switch this.document[i] {
		case '"':
			str := string(this.document[startPos:i])
			this.tokenPos = i + 1
			return str
		case '\\':
			return this.decodeQuotedStringWithEscapes(startPos, i)
		}
	}
	this.unexpectedEOD()
	return ""
}

func (this *CTEDecoder) decodeHexBytes() []byte {
	endPos := this.endPos
	i := this.tokenPos
	bytes := make([]byte, 0, 8)
	firstNybble := true
	nextByte := byte(0)
	for ; i <= endPos; i++ {
		b := this.document[i]
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
				this.errorAt(i, "Missing last hex digit")
			}
			this.tokenPos = i + 1
			return bytes
		default:
			this.unexpectedCharAt(i, "hex encoding")
		}
		if !firstNybble {
			bytes = append(bytes, nextByte)
			nextByte = 0
		}
		nextByte <<= 4
		firstNybble = !firstNybble
	}
	this.unexpectedEOD()
	return nil
}

func (this *CTEDecoder) decodeUUID() []byte {
	endPos := this.endPos
	i := this.tokenPos
	uuid := make([]byte, 0, 16)
	dashCount := 0
	firstNybble := true
	nextByte := byte(0)
Loop:
	for ; i <= endPos; i++ {
		b := this.document[i]
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
			this.unexpectedCharAt(i, "UUID")
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
		this.document[this.tokenStart+9] != '-' ||
		this.document[this.tokenStart+14] != '-' ||
		this.document[this.tokenStart+19] != '-' ||
		this.document[this.tokenStart+24] != '-' {
		this.errorAt(i, "Unrecognized named value or malformed UUID")
	}

	this.tokenPos = i
	return uuid
}

func (this *CTEDecoder) decodeDate(year int64) *compact_time.Time {
	month, _, digitCount := this.decodeDecimalInteger(0, nil)
	if digitCount > 2 {
		this.errorf("Month field is too long")
	}
	if this.peekByte() != '-' {
		this.unexpectedChar("month")
	}
	this.advanceByte()

	var day uint64
	day, _, digitCount = this.decodeDecimalInteger(0, nil)
	if digitCount == 0 {
		this.unexpectedChar("day")
	}
	if digitCount > 2 {
		this.errorf("Day field is too long")
	}
	if this.peekByte() != '/' {
		return compact_time.NewDate(int(year), int(month), int(day))
	}

	this.advanceByte()
	var hour uint64
	hour, _, digitCount = this.decodeDecimalInteger(0, nil)
	if digitCount == 0 {
		this.unexpectedChar("hour")
	}
	if digitCount > 2 {
		this.errorf("Hour field is too long")
	}
	if this.readByte() != ':' {
		this.ungetByte()
		this.unexpectedChar("hour")
	}
	t := this.decodeTime(int(hour))
	if t.TimezoneIs == compact_time.TypeLatitudeLongitude {
		return compact_time.NewTimestampLatLong(int(year), int(month), int(day),
			int(t.Hour), int(t.Minute), int(t.Second), int(t.Nanosecond),
			int(t.LatitudeHundredths), int(t.LongitudeHundredths))
	}
	return compact_time.NewTimestamp(int(year), int(month), int(day),
		int(t.Hour), int(t.Minute), int(t.Second), int(t.Nanosecond),
		t.AreaLocation)
}

func (this *CTEDecoder) decodeTime(hour int) *compact_time.Time {
	minute, _, digitCount := this.decodeDecimalInteger(0, nil)
	if digitCount > 2 {
		this.errorf("Minute field is too long")
	}
	if this.peekByte() != ':' {
		this.unexpectedChar("minute")
	}
	this.advanceByte()
	var second uint64
	second, _, digitCount = this.decodeDecimalInteger(0, nil)
	if digitCount > 2 {
		this.errorf("Second field is too long")
	}
	var nanosecond int

	if this.peekByte() == '.' {
		this.advanceByte()
		v, _, digitCount := this.decodeDecimalInteger(0, nil)
		if digitCount == 0 {
			this.unexpectedChar("nanosecond")
		}
		if digitCount > 9 {
			this.errorf("Nanosecond field is too long")
		}
		nanosecond = int(v)
		nanosecond *= subsecondMagnitudes[digitCount]
	}

	b := this.peekByte()
	if b.HasProperty(ctePropertyObjectEnd) {
		return compact_time.NewTime(hour, int(minute), int(second), nanosecond, "")
	}

	if b != '/' {
		this.unexpectedChar("time")
	}
	this.advanceByte()

	if this.peekByte().HasProperty(ctePropertyAZ) {
		areaLocationStart := this.tokenPos
		this.readWhileProperty(ctePropertyAreaLocation)
		if this.peekByte() == '/' {
			this.advanceByte()
			this.readWhileProperty(ctePropertyAreaLocation)
		}
		areaLocation := string(this.document[areaLocationStart:this.tokenPos])
		return compact_time.NewTime(hour, int(minute), int(second), nanosecond, areaLocation)
	}

	lat, long := this.decodeLatLong()
	return compact_time.NewTimeLatLong(hour, int(minute), int(second), nanosecond, lat, long)
}

func (this *CTEDecoder) decodeLatLongPortion(name string) (value int) {
	whole, _, digitCount := this.decodeDecimalInteger(0, nil)
	switch digitCount {
	case 1, 2, 3:
	// Nothing to do
	case 0:
		this.unexpectedChar(name)
	default:
		this.errorf("Too many digits decoding %v", name)
	}

	var fractional uint64
	b := this.peekByte()
	if b == '.' {
		this.advanceByte()
		fractional, _, digitCount = this.decodeDecimalInteger(0, nil)
		switch digitCount {
		case 1:
			fractional *= 10
		case 2:
			// Nothing to do
		case 0:
			this.unexpectedChar(name)
		default:
			this.errorf("Too many digits decoding %v", name)
		}
	}
	return int(whole*100 + fractional)
}

func (this *CTEDecoder) decodeLatLong() (latitudeHundredths, longitudeHundredths int) {
	latitudeHundredths = this.decodeLatLongPortion("latitude")

	if this.peekByte() != '/' {
		this.unexpectedChar("latitude/longitude")
	}
	this.advanceByte()

	longitudeHundredths = this.decodeLatLongPortion("longitude")

	return
}

func (this *CTEDecoder) decodeDecimalFloat(sign int64, coefficient uint64, bigCoefficient *big.Int, coefficientDigitCount int) (value compact_float.DFloat, bigValue *apd.Decimal) {
	exponent := int32(0)
	fractionalDigitCount := 0
	coefficient, bigCoefficient, fractionalDigitCount = this.decodeDecimalInteger(coefficient, bigCoefficient)
	if fractionalDigitCount == 0 {
		this.unexpectedChar("float fractional")
	}

	if this.peekByte() == 'e' {
		this.advanceByte()
		exponentSign := int32(1)
		switch this.peekByte() {
		case '+':
			this.advanceByte()
		case '-':
			exponentSign = -1
			this.advanceByte()
		}
		exp, bigExp, digitCount := this.decodeDecimalInteger(0, nil)
		if digitCount == 0 {
			this.unexpectedChar("float exponent")
		}
		if bigExp != nil || exp > 0x7fffffff {
			this.errorf("Exponent too big")
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

func (this *CTEDecoder) decodeHexFloat(sign int64, coefficient uint64, coefficientDigitCount int) float64 {
	exponent := 0
	fractionalDigitCount := 0
	coefficient, fractionalDigitCount = this.decodeHexInteger(coefficient)
	if fractionalDigitCount == 0 {
		this.unexpectedChar("float fractional")
	}

	if this.peekByte() == 'p' {
		this.advanceByte()
		exponentSign := 1
		switch this.peekByte() {
		case '+':
			this.advanceByte()
		case '-':
			exponentSign = -1
			this.advanceByte()
		}
		exp, bigExp, digitCount := this.decodeDecimalInteger(0, nil)
		if digitCount == 0 {
			this.unexpectedChar("float exponent")
		}
		if bigExp != nil {
			this.errorf("Exponent too big")
		}
		exponent = int(exp) * exponentSign
	}

	exponent -= fractionalDigitCount * 4

	// TODO: Overflow

	return float64(sign) * float64(coefficient) * math.Pow(float64(2), float64(exponent))
}

func (this *CTEDecoder) decodeSingleLineComment() string {
	commentStart := this.tokenPos
	this.readUntilByte('\n')
	commentEnd := this.tokenPos
	if this.document[commentEnd-1] == '\r' {
		commentEnd--
	}

	return string(this.document[commentStart:commentEnd])
}

// Handlers

type cteDecoderHandlerFunction func(*CTEDecoder)

func (this *CTEDecoder) handleNothing() {
}

func (this *CTEDecoder) handleNextState() {
	cteDecoderStateHandlers[this.currentState](this)
}

func (this *CTEDecoder) handleObject() {
	charBasedHandlers[this.peekByte()](this)
}

func (this *CTEDecoder) handleInvalidChar() {
	this.errorf("Unexpected [%v]", this.describeCurrentChar())
}

func (this *CTEDecoder) handleInvalidState() {
	this.errorf("BUG: Invalid state: %v", this.currentState)
}

func (this *CTEDecoder) handleKVSeparator() {
	this.decodeWhitespace()
	b := this.peekByte()
	if b != '=' {
		this.errorf("Expected map separator (=) but got [%v]", this.describeCurrentChar())
	}
	this.advanceByte()
	this.decodeWhitespace()
	this.endObject()
}

func (this *CTEDecoder) handleWhitespace() {
	this.decodeWhitespace()
	this.endToken()
}

func (this *CTEDecoder) handleVersion() {
	if b := this.peekByte(); b != 'c' && b != 'C' {
		this.errorf(`Expected document to begin with "c" but got [%v]`, this.describeCurrentChar())
	}
	this.advanceByte()

	version, bigVersion, digitCount := this.decodeDecimalInteger(0, nil)
	if digitCount == 0 {
		this.unexpectedChar("version number")
	}
	if bigVersion != nil {
		this.errorf("Version too big")
	}

	if !this.peekByte().HasProperty(ctePropertyWhitespace) {
		this.unexpectedChar("whitespace after version")
	}
	this.advanceByte()

	this.eventReceiver.OnVersion(version)
	this.endToken()
}

func (this *CTEDecoder) handleStringish() {
	startPos := this.tokenPos
	endPos := this.endPos
	i := startPos
	var b byte
	for ; i <= endPos; i++ {
		b = this.document[i]
		if !hasProperty(b, ctePropertyUnquotedMid) {
			break
		}
	}

	// Unquoted string
	if hasProperty(b, ctePropertyObjectEnd) || i > endPos {
		bytes := this.document[startPos:i]
		this.tokenPos = i
		this.eventReceiver.OnString(string(bytes))
		this.endObject()
		return
	}

	// Bytes, Custom, URI
	if b == '"' && i-startPos == 1 {
		initiator := this.document[startPos]
		switch initiator {
		case 'b':
			this.tokenPos = i + 1
			this.eventReceiver.OnBytes(this.decodeHexBytes())
			this.endObject()
			return
		case 'c':
			this.tokenPos = i + 1
			this.eventReceiver.OnCustom(this.decodeHexBytes())
			this.endObject()
			return
		case 'u':
			this.tokenPos = i + 1
			this.eventReceiver.OnURI(this.decodeQuotedString())
			this.endObject()
			return
		default:
			this.unexpectedChar("byte array initiator")
		}
	}

	this.unexpectedChar("unquoted string")
}

func (this *CTEDecoder) handleQuotedString() {
	this.advanceByte()
	this.eventReceiver.OnString(this.decodeQuotedString())
	this.endObject()
}

func (this *CTEDecoder) handlePositiveNumeric() {
	coefficient, bigCoefficient, digitCount := this.decodeDecimalInteger(0, nil)

	// Integer
	if this.peekByte().HasProperty(ctePropertyObjectEnd) {
		if bigCoefficient != nil {
			this.eventReceiver.OnBigInt(bigCoefficient)
		} else {
			this.eventReceiver.OnPositiveInt(coefficient)
		}
		this.endObject()
		return
	}

	switch this.peekByte() {
	case '-':
		this.advanceByte()
		v := this.decodeDate(int64(coefficient))
		this.assertAtObjectEnd("date")
		this.eventReceiver.OnCompactTime(v)
		this.endObject()
	case ':':
		this.advanceByte()
		v := this.decodeTime(int(coefficient))
		this.assertAtObjectEnd("time")
		this.eventReceiver.OnCompactTime(v)
		this.endObject()
	case '.':
		this.advanceByte()
		value, bigValue := this.decodeDecimalFloat(1, coefficient, bigCoefficient, digitCount)
		this.assertAtObjectEnd("float")
		if bigValue != nil {
			this.eventReceiver.OnBigDecimalFloat(bigValue)
		} else {
			this.eventReceiver.OnDecimalFloat(value)
		}
		this.endObject()
	default:
		this.unexpectedChar("numeric")
	}
}

func (this *CTEDecoder) handleNegativeNumeric() {
	this.advanceByte()
	switch this.peekByte() {
	case '0':
		this.handleOtherBaseNegative()
		return
	case '@':
		this.advanceByte()
		nameStart := this.tokenPos
		this.readWhileProperty(ctePropertyAZ)
		token := strings.ToLower(string(this.document[nameStart:this.tokenPos]))
		if token != "inf" {
			this.errorf("Unknown named value: %v", token)
		}
		this.eventReceiver.OnFloat(math.Inf(-1))
		return
	}

	coefficient, bigCoefficient, digitCount := this.decodeDecimalInteger(0, nil)

	// Integer
	if this.peekByte().HasProperty(ctePropertyObjectEnd) {
		if bigCoefficient != nil {
			// TODO: More efficient way to negate?
			bigCoefficient = bigCoefficient.Mul(bigCoefficient, bigIntN1)
			this.eventReceiver.OnBigInt(bigCoefficient)
		} else {
			this.eventReceiver.OnNegativeInt(coefficient)
		}
		this.endObject()
		return
	}

	switch this.peekByte() {
	case '-':
		this.advanceByte()
		v := this.decodeDate(-int64(coefficient))
		this.assertAtObjectEnd("time")
		this.eventReceiver.OnCompactTime(v)
		this.endObject()
	case '.':
		this.advanceByte()
		value, bigValue := this.decodeDecimalFloat(-1, coefficient, bigCoefficient, digitCount)
		this.assertAtObjectEnd("float")
		if bigValue != nil {
			this.eventReceiver.OnBigDecimalFloat(bigValue)
		} else {
			this.eventReceiver.OnDecimalFloat(value)
		}
		this.endObject()
	default:
		this.unexpectedChar("numeric")
	}
}

func (this *CTEDecoder) handleOtherBasePositive() {
	this.advanceByte()
	b := this.peekByte()

	if b.HasProperty(ctePropertyObjectEnd) {
		this.eventReceiver.OnPositiveInt(0)
		this.endObject()
		return
	}
	this.advanceByte()

	switch b {
	case 'b':
		v := this.decodeBinaryInteger()
		this.assertAtObjectEnd("binary integer")
		this.eventReceiver.OnPositiveInt(v)
		this.endObject()
	case 'o':
		v := this.decodeOctalInteger()
		this.assertAtObjectEnd("octal integer")
		this.eventReceiver.OnPositiveInt(v)
		this.endObject()
	case 'x':
		v, digitCount := this.decodeHexInteger(0)
		if this.peekByte() == '.' {
			this.advanceByte()
			fv := this.decodeHexFloat(1, v, digitCount)
			this.assertAtObjectEnd("hex float")
			this.eventReceiver.OnFloat(fv)
			this.endObject()
		} else {
			this.assertAtObjectEnd("hex integer")
			this.eventReceiver.OnPositiveInt(v)
			this.endObject()
		}
	case '.':
		value, bigValue := this.decodeDecimalFloat(1, 0, nil, 0)
		this.assertAtObjectEnd("float")
		if bigValue != nil {
			this.eventReceiver.OnBigDecimalFloat(bigValue)
		} else {
			this.eventReceiver.OnDecimalFloat(value)
		}
		this.endObject()
	default:
		if b.HasProperty(cteProperty09) && this.peekByte() == ':' {
			this.advanceByte()
			v := this.decodeTime(int(b - '0'))
			this.assertAtObjectEnd("time")
			this.eventReceiver.OnCompactTime(v)
			this.endObject()
			return
		}
		this.ungetByte()
		this.unexpectedChar("numeric base")
	}
}

func (this *CTEDecoder) handleOtherBaseNegative() {
	this.advanceByte()
	b := this.readByte()
	switch b {
	case 'b':
		v := this.decodeBinaryInteger()
		this.assertAtObjectEnd("binary integer")
		this.eventReceiver.OnNegativeInt(v)
		this.endObject()
	case 'o':
		v := this.decodeOctalInteger()
		this.assertAtObjectEnd("octal integer")
		this.eventReceiver.OnNegativeInt(v)
		this.endObject()
	case 'x':
		v, digitCount := this.decodeHexInteger(0)
		if this.peekByte() == '.' {
			this.advanceByte()
			fv := this.decodeHexFloat(-1, v, digitCount)
			this.assertAtObjectEnd("hex float")
			this.eventReceiver.OnFloat(fv)
			this.endObject()
		} else {
			this.assertAtObjectEnd("hex integer")
			this.eventReceiver.OnNegativeInt(v)
			this.endObject()
		}
	case '.':
		value, bigValue := this.decodeDecimalFloat(-1, 0, nil, 0)
		this.assertAtObjectEnd("float")
		if bigValue != nil {
			this.eventReceiver.OnBigDecimalFloat(bigValue)
		} else {
			this.eventReceiver.OnDecimalFloat(value)
		}
		this.endObject()
	default:
		this.ungetByte()
		this.unexpectedChar("numeric base")
	}
}

func (this *CTEDecoder) handleListBegin() {
	this.advanceByte()
	this.eventReceiver.OnList()
	this.stackContainer(cteDecoderStateAwaitListItem)
	this.endToken()
}

func (this *CTEDecoder) handleListEnd() {
	this.advanceByte()
	this.eventReceiver.OnEnd()
	this.unstackContainer()
	this.endObject()
}

func (this *CTEDecoder) handleMapBegin() {
	this.advanceByte()
	this.eventReceiver.OnMap()
	this.stackContainer(cteDecoderStateAwaitMapKey)
	this.endToken()
}

func (this *CTEDecoder) handleMapEnd() {
	this.advanceByte()
	this.eventReceiver.OnEnd()
	this.unstackContainer()
	this.endObject()
}

func (this *CTEDecoder) handleMetadataBegin() {
	this.advanceByte()
	this.eventReceiver.OnMetadata()
	this.stackContainer(cteDecoderStateAwaitMetaKey)
	this.endToken()
}

func (this *CTEDecoder) handleMetadataEnd() {
	this.advanceByte()
	this.eventReceiver.OnEnd()
	this.unstackContainer()
	this.endToken()
	// Don't transition state because metadata is a pseudo-object
}

func (this *CTEDecoder) handleComment() {
	this.readByte()
	switch this.readByte() {
	case '/':
		this.eventReceiver.OnComment()
		contents := this.decodeSingleLineComment()
		if len(contents) > 0 {
			this.eventReceiver.OnString(contents)
		}
		this.eventReceiver.OnEnd()
		this.endToken()
		// Don't transition state because a comment is a pseudo-object
	case '*':
		this.eventReceiver.OnComment()
		this.stackContainer(cteDecoderStateAwaitCommentItem)
		this.endToken()
	default:
		this.ungetByte()
		this.unexpectedChar("comment")
	}
}

func (this *CTEDecoder) handleCommentContent() {
	startPos := this.tokenPos
	endPos := this.endPos

	// Skip first byte so that index-1 is not comparing against part of the
	// original multiline comment initiator '*'
	this.advanceByte()

	for i := this.tokenPos; i <= endPos; i++ {
		bLast := this.document[i-1]
		bNow := this.document[i]

		if bLast == '*' && bNow == '/' {
			endOffset := i - 1
			if endOffset > startPos {
				str := string(this.document[startPos:endOffset])
				this.eventReceiver.OnString(str)
			}
			this.tokenPos = i + 1
			this.eventReceiver.OnEnd()
			this.unstackContainer()
			this.endToken()
			return
		}

		if bLast == '/' && bNow == '*' {
			endOffset := i - 1
			if endOffset > this.tokenStart {
				str := string(this.document[startPos:endOffset])
				this.eventReceiver.OnString(str)
			}
			this.tokenPos = i + 1
			this.eventReceiver.OnComment()
			this.stackContainer(cteDecoderStateAwaitCommentItem)
			this.endToken()
			return
		}
	}

	this.unexpectedEOD()
}

func (this *CTEDecoder) handleMarkupBegin() {
	this.advanceByte()
	this.eventReceiver.OnMarkup()
	this.stackContainer(cteDecoderStateAwaitMarkupValue)
	this.endToken()
}

func (this *CTEDecoder) handleMarkupContentBegin() {
	this.advanceByte()
	this.eventReceiver.OnEnd()
	this.changeState(cteDecoderStateAwaitMarkupItem)
	this.endToken()
}

func (this *CTEDecoder) handleMarkupWithEscapes(startPos, firstEscape int) string {
	this.errorf("TODO: CTEDecoder: Markup with escape sequences, entity refs")
	return ""
}

func (this *CTEDecoder) handleMarkupContent() {
	startPos := this.tokenPos
	endPos := this.endPos
	i := startPos
	for ; i <= endPos; i++ {
		switch this.document[i] {
		case '/':
			switch this.getByteAt(i + 1) {
			case '/':
				if i > startPos {
					this.eventReceiver.OnString(string(this.document[startPos:i]))
				}
				this.tokenPos = i + 2
				this.eventReceiver.OnComment()
				contents := this.decodeSingleLineComment()
				if len(contents) > 0 {
					this.eventReceiver.OnString(contents)
				}
				this.advanceByte()
				this.eventReceiver.OnEnd()
				this.endToken()
				// Don't transition state because a comment is a pseudo-object
				return
			case '*':
				if i > startPos {
					this.eventReceiver.OnString(string(this.document[startPos:i]))
				}
				this.tokenPos = i + 2
				this.eventReceiver.OnComment()
				this.stackContainer(cteDecoderStateAwaitCommentItem)
				this.endToken()
				return
			}
		case '<':
			str := string(this.document[startPos:i])
			this.tokenPos = i + 1
			if len(str) > 0 {
				this.eventReceiver.OnString(str)
			}

			this.tokenPos = i
			this.handleMarkupBegin()
			return
		case '>':
			str := string(this.document[startPos:i])
			this.tokenPos = i + 1
			if len(str) > 0 {
				this.eventReceiver.OnString(str)
			}
			this.eventReceiver.OnEnd()
			this.unstackContainer()
			this.endObject()
			return
		case '\\':
			this.handleMarkupWithEscapes(startPos, i)
			return
		}
	}
	this.unexpectedEOD()
}

func (this *CTEDecoder) handleMarkupEnd() {
	this.advanceByte()
	this.eventReceiver.OnEnd()
	this.eventReceiver.OnEnd()
	this.unstackContainer()
	this.endObject()
}

func (this *CTEDecoder) handleNamedValue() {
	this.advanceByte()
	nameStart := this.tokenPos
	this.readWhileProperty(ctePropertyAZ)
	token := strings.ToLower(string(this.document[nameStart:this.tokenPos]))
	switch token {
	case "nil":
		this.eventReceiver.OnNil()
		this.endObject()
		return
	case "nan":
		this.eventReceiver.OnNan(false)
		this.endObject()
		return
	case "snan":
		this.eventReceiver.OnNan(true)
		this.endObject()
		return
	case "inf":
		this.eventReceiver.OnFloat(math.Inf(1))
		this.endObject()
		return
	case "false":
		this.eventReceiver.OnFalse()
		this.endObject()
		return
	case "true":
		this.eventReceiver.OnTrue()
		this.endObject()
		return
	}

	this.ungetAll()
	this.advanceByte()
	this.eventReceiver.OnUUID(this.decodeUUID())
	this.endObject()
}

func (this *CTEDecoder) handleVerbatimString() {
	this.advanceByte()
	this.assertNotEOD()
	sentinelStart := this.tokenPos
	this.readUntilProperty(ctePropertyWhitespace)
	sentinelEnd := this.tokenPos
	wsByte := this.readByte()
	if wsByte == '\r' {
		if this.readByte() != '\n' {
			this.unexpectedChar("verbatim sentinel")
		}
	}
	sentinel := this.document[sentinelStart:sentinelEnd]
	searchSpace := this.document[this.tokenPos:]
	index := bytes.Index(searchSpace, sentinel)
	if index < 0 {
		this.errorf("Verbatim sentinel sequence [%v] not found in document", string(sentinel))
	}
	str := string(searchSpace[:index])
	this.tokenPos += index + len(sentinel)
	this.assertAtObjectEnd("verbatim string")
	this.eventReceiver.OnString(str)
	this.endObject()
}

func (this *CTEDecoder) handleReference() {
	this.advanceByte()
	this.eventReceiver.OnReference()
	if this.peekByte().HasProperty(ctePropertyWhitespace) {
		this.errorf("Whitespace not allowed between reference and tag name")
	}
	this.endToken()
}

func (this *CTEDecoder) handleMarker() {
	this.advanceByte()
	this.eventReceiver.OnMarker()
	if this.peekByte().HasProperty(ctePropertyWhitespace) {
		this.errorf("Whitespace not allowed between marker and tag name")
	}
	this.endToken()
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

func (this cteByte) HasProperty(property cteByteProprty) bool {
	return cteByteProperties[this].HasProperty(property)
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

func (this cteByteProprty) HasProperty(property cteByteProprty) bool {
	return this&property != 0
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
