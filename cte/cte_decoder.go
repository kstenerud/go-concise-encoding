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

// Performs encoding and decoding of Concise Text Encoding documents
// (https://github.com/kstenerud/concise-encoding/blob/master/cte-specification.md).
//
// The decoder decodes a document to produce data events, and the encoder
// consumes data events to produce a document.
package cte

import (
	"bytes"
	"fmt"
	"io"
	"math"

	"github.com/kstenerud/go-concise-encoding/internal/chars"

	"github.com/kstenerud/go-concise-encoding/debug"
	"github.com/kstenerud/go-concise-encoding/events"
	"github.com/kstenerud/go-concise-encoding/internal/common"
	"github.com/kstenerud/go-concise-encoding/options"
)

// Decodes CTE documents.
type OldDecoder struct {
	buffer         DecodeBuffer
	eventReceiver  events.DataEventReceiver
	containerState []cteDecoderState
	currentState   cteDecoderState
	opts           options.CTEDecoderOptions
}

// Create a new CTE decoder, which will read from reader and send data events
// to nextReceiver. If opts is nil, default options will be used.
func NewOldDecoder(opts *options.CTEDecoderOptions) *OldDecoder {
	_this := &OldDecoder{}
	_this.Init(opts)
	return _this
}

// Initialize this decoder, which will read from reader and send data events
// to nextReceiver. If opts is nil, default options will be used.
func (_this *OldDecoder) Init(opts *options.CTEDecoderOptions) {
	opts = opts.WithDefaultsApplied()
	_this.opts = *opts
}

func (_this *OldDecoder) reset() {
	_this.buffer.Reset()
	_this.eventReceiver = nil
	_this.containerState = _this.containerState[:0]
	_this.currentState = 0
}

// Run the complete decode process. The document and data receiver specified
// when initializing the decoder will be used.
func (_this *OldDecoder) Decode(reader io.Reader, eventReceiver events.DataEventReceiver) (err error) {
	defer func() {
		_this.reset()
		if !debug.DebugOptions.PassThroughPanics {
			if r := recover(); r != nil {
				switch v := r.(type) {
				case error:
					err = v
				default:
					err = fmt.Errorf("%v", r)
				}
			}
		}
	}()

	_this.buffer.Init(reader, _this.opts.BufferSize, chooseLowWater(_this.opts.BufferSize))
	_this.eventReceiver = eventReceiver

	_this.eventReceiver.OnBeginDocument()

	_this.buffer.RefillIfNecessary()

	_this.currentState = cteDecoderStateAwaitObject

	// Forgive initial whitespace even though it's technically not allowed
	_this.buffer.SkipWhitespace()

	_this.handleVersion()

	for !_this.buffer.IsEndOfDocument() {
		_this.handleNextState()
		_this.buffer.RefillIfNecessary()
	}

	_this.eventReceiver.OnEndDocument()
	return
}

func (_this *OldDecoder) DecodeDocument(document []byte, eventReceiver events.DataEventReceiver) (err error) {
	return _this.Decode(bytes.NewBuffer(document), eventReceiver)
}

// ============================================================================

// Internal

func chooseLowWater(bufferSize int) int {
	lowWater := bufferSize / 50
	if lowWater < 30 {
		lowWater = 30
	}
	return lowWater
}

// State

func (_this *OldDecoder) changeState(newState cteDecoderState) {
	_this.currentState = newState
}

func (_this *OldDecoder) stackContainer(nextState cteDecoderState) {
	_this.containerState = append(_this.containerState, _this.currentState)
	_this.changeState(nextState)
}

func (_this *OldDecoder) unstackContainer() {
	index := len(_this.containerState) - 1
	_this.changeState(_this.containerState[index])
	_this.containerState = _this.containerState[:index]
}

func (_this *OldDecoder) endObject() {
	_this.changeState(cteDecoderStateTransitions[_this.currentState])
}

// Handlers

type cteDecoderHandlerFunction func(*OldDecoder)

func (_this *OldDecoder) handleNothing() {
}

func (_this *OldDecoder) handleNextState() {
	cteDecoderStateHandlers[_this.currentState](_this)
}

func (_this *OldDecoder) handleObject() {
	charBasedHandlers[_this.buffer.PeekByteAllowEOD()](_this)
}

func (_this *OldDecoder) handleInvalidChar() {
	_this.buffer.Errorf("Unexpected [%v]", _this.buffer.DescribeCurrentChar())
}

func (_this *OldDecoder) handleKVSeparator() {
	_this.buffer.SkipWhitespace()
	if _this.buffer.PeekByteNoEOD() != '=' {
		_this.buffer.Errorf("Expected map separator (=) but got [%v]", _this.buffer.DescribeCurrentChar())
	}
	_this.buffer.AdvanceByte()
	_this.buffer.SkipWhitespace()
	_this.endObject()
}

func (_this *OldDecoder) handleWhitespace() {
	_this.buffer.SkipWhitespace()
}

func (_this *OldDecoder) handleVersion() {
	if b := _this.buffer.PeekByteNoEOD(); b != 'c' && b != 'C' {
		_this.buffer.Errorf(`Expected document to begin with "c" but got [%v]`, _this.buffer.DescribeCurrentChar())
	}

	_this.buffer.AdvanceByte()

	version, bigVersion, digitCount := _this.buffer.DecodeDecimalUint(0, nil)
	if digitCount == 0 {
		_this.buffer.UnexpectedChar("version number")
	}
	if bigVersion != nil {
		_this.buffer.Errorf("Version too big")
	}

	b := _this.buffer.PeekByteAllowEOD()
	if !b.HasProperty(chars.CharIsWhitespace) && b != chars.EndOfDocumentMarker {
		_this.buffer.UnexpectedChar("whitespace after version")
	}

	_this.eventReceiver.OnVersion(version)
}

func (_this *OldDecoder) handleUnquotedString() {
	_this.buffer.BeginToken()
	_this.buffer.ReadUntilPropertyAllowEOD(chars.CharNeedsQuote)

	if _this.buffer.PeekByteAllowEOD().HasProperty(chars.CharIsObjectEnd) || _this.isOnCommentInitiator() {
		bytes := _this.buffer.GetToken()
		_this.eventReceiver.OnArray(events.ArrayTypeString, uint64(len(bytes)), bytes)
		_this.endObject()
		return
	}

	_this.buffer.UnexpectedChar("unquoted string")
}

func (_this *OldDecoder) isOnCommentInitiator() bool {
	if _this.buffer.PeekByteAllowEOD() == '/' {
		_this.buffer.AdvanceByte()
		b := _this.buffer.PeekByteAllowEOD()
		_this.buffer.UngetByte()
		switch b {
		case '/', '*':
			return true
		}
	}
	return false
}

func (_this *OldDecoder) handleQuotedString() {
	_this.buffer.AdvanceByte()
	bytes := _this.buffer.DecodeQuotedString()
	_this.eventReceiver.OnArray(events.ArrayTypeString, uint64(len(bytes)), bytes)
	_this.endObject()
}

func (_this *OldDecoder) handlePositiveNumeric() {
	coefficient, bigCoefficient, digitCount := _this.buffer.DecodeDecimalUint(0, nil)
	b := _this.buffer.PeekByteAllowEOD()
	switch b {
	case '-':
		_this.buffer.AdvanceByte()
		v := _this.buffer.DecodeDate(int64(coefficient))
		_this.buffer.AssertAtObjectEnd("date")
		_this.eventReceiver.OnCompactTime(v)
		_this.endObject()
		return
	case ':':
		_this.buffer.AdvanceByte()
		v := _this.buffer.DecodeTime(int(coefficient))
		_this.buffer.AssertAtObjectEnd("time")
		_this.eventReceiver.OnCompactTime(v)
		_this.endObject()
		return
	case '.':
		_this.buffer.AdvanceByte()
		value, bigValue, _ := _this.buffer.DecodeDecimalFloat(1, coefficient, bigCoefficient, digitCount)
		_this.buffer.AssertAtObjectEnd("float")
		if bigValue != nil {
			_this.eventReceiver.OnBigDecimalFloat(bigValue)
		} else {
			_this.eventReceiver.OnDecimalFloat(value)
		}
		_this.endObject()
		return
	default:
		if b.HasProperty(chars.CharIsObjectEnd) {
			if bigCoefficient != nil {
				_this.eventReceiver.OnBigInt(bigCoefficient)
			} else {
				_this.eventReceiver.OnPositiveInt(coefficient)
			}
			_this.endObject()
			return
		}
	}
	_this.buffer.UnexpectedChar("numeric")
}

func (_this *OldDecoder) handleNegativeNumeric() {
	_this.buffer.AdvanceByte()
	switch _this.buffer.PeekByteNoEOD() {
	case '0':
		_this.handleOtherBaseNegative()
		return
	case '@':
		_this.buffer.AdvanceByte()
		namedValue := string(_this.buffer.DecodeNamedValue())
		if namedValue != "inf" {
			_this.buffer.Errorf("Unknown named value: %v", namedValue)
		}
		_this.eventReceiver.OnFloat(math.Inf(-1))
		_this.endObject()
		return
	}

	coefficient, bigCoefficient, digitCount := _this.buffer.DecodeDecimalUint(0, nil)
	b := _this.buffer.PeekByteAllowEOD()
	switch b {
	case '-':
		_this.buffer.AdvanceByte()
		v := _this.buffer.DecodeDate(-int64(coefficient))
		_this.buffer.AssertAtObjectEnd("time")
		_this.eventReceiver.OnCompactTime(v)
		_this.endObject()
		return
	case '.':
		_this.buffer.AdvanceByte()
		value, bigValue, _ := _this.buffer.DecodeDecimalFloat(-1, coefficient, bigCoefficient, digitCount)
		_this.buffer.AssertAtObjectEnd("float")
		if bigValue != nil {
			_this.eventReceiver.OnBigDecimalFloat(bigValue)
		} else {
			_this.eventReceiver.OnDecimalFloat(value)
		}
		_this.endObject()
		return
	default:
		if b.HasProperty(chars.CharIsObjectEnd) {
			if bigCoefficient != nil {
				// TODO: More efficient way to negate?
				bigCoefficient = bigCoefficient.Mul(bigCoefficient, common.BigIntN1)
				_this.eventReceiver.OnBigInt(bigCoefficient)
			} else {
				_this.eventReceiver.OnNegativeInt(coefficient)
			}
			_this.endObject()
			return
		}
	}
	_this.buffer.UnexpectedChar("numeric")
}

func (_this *OldDecoder) handleOtherBasePositive() {
	_this.buffer.AdvanceByte()
	b := _this.buffer.PeekByteAllowEOD()

	if b.HasProperty(chars.CharIsObjectEnd) && b != '.' {
		_this.eventReceiver.OnPositiveInt(0)
		_this.endObject()
		return
	}
	_this.buffer.AdvanceByte()

	switch b {
	case 'b':
		v, bigV, _ := _this.buffer.DecodeBinaryUint()
		_this.buffer.AssertAtObjectEnd("binary integer")
		if bigV != nil {
			_this.eventReceiver.OnBigInt(bigV)
		} else {
			_this.eventReceiver.OnPositiveInt(v)
		}
		_this.endObject()
	case 'o':
		v, bigV, _ := _this.buffer.DecodeOctalUint()
		_this.buffer.AssertAtObjectEnd("octal integer")
		if bigV != nil {
			_this.eventReceiver.OnBigInt(bigV)
		} else {
			_this.eventReceiver.OnPositiveInt(v)
		}
		_this.endObject()
	case 'x':
		v, bigV, digitCount := _this.buffer.DecodeHexUint(0, nil)
		if _this.buffer.PeekByteAllowEOD() == '.' {
			_this.buffer.AdvanceByte()
			fv, bigFV, _ := _this.buffer.DecodeHexFloat(1, v, bigV, digitCount)
			_this.buffer.AssertAtObjectEnd("hex float")
			if bigFV != nil {
				_this.eventReceiver.OnBigFloat(bigFV)
			} else {
				_this.eventReceiver.OnFloat(fv)
			}
			_this.endObject()
		} else {
			_this.buffer.AssertAtObjectEnd("hex integer")
			if bigV != nil {
				_this.eventReceiver.OnBigInt(bigV)
			} else {
				_this.eventReceiver.OnPositiveInt(v)
			}
			_this.endObject()
		}
	case '.':
		value, bigValue, _ := _this.buffer.DecodeDecimalFloat(1, 0, nil, 0)
		_this.buffer.AssertAtObjectEnd("float")
		if bigValue != nil {
			_this.eventReceiver.OnBigDecimalFloat(bigValue)
		} else {
			_this.eventReceiver.OnDecimalFloat(value)
		}
		_this.endObject()
	default:
		if b.HasProperty(chars.CharIsDigitBase10) && _this.buffer.PeekByteNoEOD() == ':' {
			_this.buffer.AdvanceByte()
			v := _this.buffer.DecodeTime(int(b - '0'))
			_this.buffer.AssertAtObjectEnd("time")
			_this.eventReceiver.OnCompactTime(v)
			_this.endObject()
			return
		}
		_this.buffer.UngetByte()
		_this.buffer.UnexpectedChar("numeric base")
	}
}

func (_this *OldDecoder) handleOtherBaseNegative() {
	_this.buffer.AdvanceByte()
	b := _this.buffer.ReadByte()
	switch b {
	case 'b':
		v, bigV, _ := _this.buffer.DecodeBinaryUint()
		_this.buffer.AssertAtObjectEnd("binary integer")
		if bigV != nil {
			bigV = bigV.Neg(bigV)
			_this.eventReceiver.OnBigInt(bigV)
		} else {
			_this.eventReceiver.OnNegativeInt(v)
		}
		_this.endObject()
	case 'o':
		v, bigV, _ := _this.buffer.DecodeOctalUint()
		_this.buffer.AssertAtObjectEnd("octal integer")
		if bigV != nil {
			bigV = bigV.Neg(bigV)
			_this.eventReceiver.OnBigInt(bigV)
		} else {
			_this.eventReceiver.OnNegativeInt(v)
		}
		_this.endObject()
	case 'x':
		v, bigV, digitCount := _this.buffer.DecodeHexUint(0, nil)
		if _this.buffer.PeekByteAllowEOD() == '.' {
			_this.buffer.AdvanceByte()
			fv, bigFV, _ := _this.buffer.DecodeHexFloat(-1, v, bigV, digitCount)
			_this.buffer.AssertAtObjectEnd("hex float")
			if bigFV != nil {
				_this.eventReceiver.OnBigFloat(bigFV)
			} else {
				_this.eventReceiver.OnFloat(fv)
			}
			_this.endObject()
		} else {
			_this.buffer.AssertAtObjectEnd("hex integer")
			if bigV != nil {
				bigV = bigV.Neg(bigV)
				_this.eventReceiver.OnBigInt(bigV)
			} else {
				_this.eventReceiver.OnNegativeInt(v)
			}
			_this.endObject()
		}
	case '.':
		value, bigValue, _ := _this.buffer.DecodeDecimalFloat(-1, 0, nil, 0)
		_this.buffer.AssertAtObjectEnd("float")
		if bigValue != nil {
			_this.eventReceiver.OnBigDecimalFloat(bigValue)
		} else {
			_this.eventReceiver.OnDecimalFloat(value)
		}
		_this.endObject()
	default:
		_this.buffer.UngetByte()
		_this.buffer.UnexpectedChar("numeric base")
	}
}

func (_this *OldDecoder) handleListBegin() {
	_this.buffer.AdvanceByte()
	_this.eventReceiver.OnList()
	_this.stackContainer(cteDecoderStateAwaitListItem)
}

func (_this *OldDecoder) handleListEnd() {
	_this.buffer.AdvanceByte()
	_this.eventReceiver.OnEnd()
	_this.unstackContainer()
	_this.endObject()
}

func (_this *OldDecoder) handleMapBegin() {
	_this.buffer.AdvanceByte()
	_this.eventReceiver.OnMap()
	_this.stackContainer(cteDecoderStateAwaitMapKey)
}

func (_this *OldDecoder) handleMapEnd() {
	_this.buffer.AdvanceByte()
	_this.eventReceiver.OnEnd()
	_this.unstackContainer()
	_this.endObject()
}

func (_this *OldDecoder) handleMetadataBegin() {
	_this.buffer.AdvanceByte()
	_this.eventReceiver.OnMetadata()
	_this.stackContainer(cteDecoderStateAwaitMetaKey)
}

func (_this *OldDecoder) handleMetadataEnd() {
	_this.buffer.AdvanceByte()
	_this.eventReceiver.OnEnd()
	_this.unstackContainer()
	// Don't transition state because metadata is a pseudo-object
}

func (_this *OldDecoder) handleComment() {
	_this.buffer.AdvanceByte()
	switch _this.buffer.ReadByte() {
	case '/':
		_this.eventReceiver.OnComment()
		contents := _this.buffer.DecodeSingleLineComment()
		if len(contents) > 0 {
			_this.eventReceiver.OnArray(events.ArrayTypeString, uint64(len(contents)), contents)
		}
		_this.eventReceiver.OnEnd()
	case '*':
		_this.eventReceiver.OnComment()
		_this.stackContainer(cteDecoderStateAwaitCommentItem)
	default:
		_this.buffer.UngetByte()
		_this.buffer.UnexpectedChar("comment")
	}
}

func (_this *OldDecoder) handleCommentContent() {
	str, next := _this.buffer.DecodeMultilineComment()
	if len(str) > 0 {
		_this.eventReceiver.OnArray(events.ArrayTypeString, uint64(len(str)), str)
	}
	switch next {
	case nextIsCommentBegin:
		_this.eventReceiver.OnComment()
		_this.stackContainer(cteDecoderStateAwaitCommentItem)
	case nextIsCommentEnd:
		_this.eventReceiver.OnEnd()
		_this.unstackContainer()
	}
}

func (_this *OldDecoder) handleMarkupBegin() {
	_this.buffer.AdvanceByte()
	_this.eventReceiver.OnMarkup()
	_this.stackContainer(cteDecoderStateAwaitMarkupName)
}

func (_this *OldDecoder) handleMarkupContentBegin() {
	_this.buffer.AdvanceByte()
	_this.eventReceiver.OnEnd()
	_this.changeState(cteDecoderStateAwaitMarkupItem)
}

func (_this *OldDecoder) handleMarkupContent() {
	str, next := _this.buffer.DecodeMarkupContent()
	if len(str) > 0 {
		_this.eventReceiver.OnArray(events.ArrayTypeString, uint64(len(str)), str)
	}
	switch next {
	case nextIsCommentBegin:
		_this.eventReceiver.OnComment()
		_this.stackContainer(cteDecoderStateAwaitCommentItem)
	case nextIsCommentEnd:
		_this.eventReceiver.OnEnd()
		_this.unstackContainer()
	case nextIsSingleLineComment:
		_this.eventReceiver.OnComment()
		contents := _this.buffer.DecodeSingleLineComment()
		if len(contents) > 0 {
			_this.eventReceiver.OnArray(events.ArrayTypeString, uint64(len(contents)), contents)
		}
		_this.eventReceiver.OnEnd()
	case nextIsMarkupBegin:
		_this.eventReceiver.OnMarkup()
		_this.stackContainer(cteDecoderStateAwaitMarkupName)
	case nextIsMarkupEnd:
		_this.eventReceiver.OnEnd()
		_this.unstackContainer()
		_this.endObject()
	}
}

func (_this *OldDecoder) handleMarkupEnd() {
	_this.buffer.AdvanceByte()
	_this.eventReceiver.OnEnd()
	_this.eventReceiver.OnEnd()
	_this.unstackContainer()
	_this.endObject()
}

func (_this *OldDecoder) handleNamedValueOrUUID() {
	_this.buffer.AdvanceByte()
	namedValue := _this.buffer.DecodeNamedValue()
	switch string(namedValue) {
	case "na":
		_this.eventReceiver.OnNA()
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

	// UUID
	if len(namedValue) != 36 ||
		namedValue[8] != '-' ||
		namedValue[13] != '-' ||
		namedValue[18] != '-' ||
		namedValue[23] != '-' {
		_this.buffer.UngetBytes(len(namedValue) + 1)
		_this.buffer.Errorf("Malformed UUID or unknown named value: [%s]", string(namedValue))
	}

	decodeHex := func(b byte) byte {
		switch {
		case chars.ByteHasProperty(b, chars.CharIsDigitBase10):
			return byte(b - '0')
		case chars.ByteHasProperty(b, chars.CharIsLowerAF):
			return byte(b - 'a' + 10)
		case chars.ByteHasProperty(b, chars.CharIsUpperAF):
			return byte(b - 'A' + 10)
		default:
			_this.buffer.UngetBytes(len(namedValue) + 1)
			_this.buffer.Errorf("Unexpected char [%c] in UUID [%s]", b, string(namedValue))
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

	decodeSection(namedValue[:8], namedValue)
	decodeSection(namedValue[9:13], namedValue[4:])
	decodeSection(namedValue[14:18], namedValue[6:])
	decodeSection(namedValue[19:23], namedValue[8:])
	decodeSection(namedValue[24:36], namedValue[10:])

	_this.eventReceiver.OnUUID(namedValue[:16])
	_this.endObject()
}

func (_this *OldDecoder) handleConstant() {
	_this.buffer.AdvanceByte()
	name := _this.buffer.DecodeNamedValue()
	explicitValue := false
	if _this.buffer.PeekByteAllowEOD() == ':' {
		_this.buffer.AdvanceByte()
		explicitValue = true
	} else {
		_this.endObject()
	}
	_this.eventReceiver.OnConstant(name, explicitValue)
}

func (_this *OldDecoder) decodeStringArray(arrayType events.ArrayType) {
	bytes := _this.buffer.DecodeStringArray()
	_this.eventReceiver.OnArray(arrayType, uint64(len(bytes)), bytes)
	_this.endObject()
}

func (_this *OldDecoder) decodeCustomText() {
	_this.decodeStringArray(events.ArrayTypeCustomText)
}

func (_this *OldDecoder) decodeRID() {
	_this.decodeStringArray(events.ArrayTypeResourceID)
}

func (_this *OldDecoder) finishTypedArray(arrayType events.ArrayType, digitType string, bytesPerElement int, data []byte) {
	switch _this.buffer.PeekByteNoEOD() {
	case '|':
		_this.buffer.AdvanceByte()
		_this.eventReceiver.OnArray(arrayType, uint64(len(data)/bytesPerElement), data)
		_this.endObject()
		return
	default:
		_this.buffer.Errorf("Expected %v digits", digitType)
	}
}

func (_this *OldDecoder) decodeCustomBinary() {
	digitType := "hex"
	var data []uint8
	for {
		_this.buffer.ReadWhilePropertyNoEOD(chars.CharIsWhitespace)
		v, count := _this.buffer.DecodeSmallHexUint()
		if count == 0 {
			break
		}
		if v > maxUint8Value {
			_this.buffer.Errorf("%v value too big for array type", digitType)
		}
		data = append(data, uint8(v))
	}
	_this.finishTypedArray(events.ArrayTypeCustomBinary, digitType, 1, data)
}

func (_this *OldDecoder) decodeArrayU8(digitType string, decodeElement func() (v uint64, digitCount int)) {
	var data []uint8
	for {
		_this.buffer.ReadWhilePropertyNoEOD(chars.CharIsWhitespace)
		v, count := decodeElement()
		if count == 0 {
			break
		}
		if v > maxUint8Value {
			_this.buffer.Errorf("%v value too big for array type", digitType)
		}
		data = append(data, uint8(v))
	}
	_this.finishTypedArray(events.ArrayTypeUint8, digitType, 1, data)
}

func (_this *OldDecoder) decodeArrayU16(digitType string, decodeElement func() (v uint64, digitCount int)) {
	var data []uint8
	for {
		_this.buffer.ReadWhilePropertyNoEOD(chars.CharIsWhitespace)
		v, count := decodeElement()
		if count == 0 {
			break
		}
		if v > maxUint16Value {
			_this.buffer.Errorf("%v value too big for array type", digitType)
		}
		data = append(data, uint8(v), uint8(v>>8))
	}
	_this.finishTypedArray(events.ArrayTypeUint16, digitType, 2, data)
}

func (_this *OldDecoder) decodeArrayU32(digitType string, decodeElement func() (v uint64, digitCount int)) {
	var data []uint8
	for {
		_this.buffer.ReadWhilePropertyNoEOD(chars.CharIsWhitespace)
		v, count := decodeElement()
		if count == 0 {
			break
		}
		if v > maxUint32Value {
			_this.buffer.Errorf("%v value too big for array type", digitType)
		}
		data = append(data, uint8(v), uint8(v>>8), uint8(v>>16), uint8(v>>24))
	}
	_this.finishTypedArray(events.ArrayTypeUint32, digitType, 4, data)
}

func (_this *OldDecoder) decodeArrayU64(digitType string, decodeElement func() (v uint64, digitCount int)) {
	var data []uint8
	for {
		_this.buffer.ReadWhilePropertyNoEOD(chars.CharIsWhitespace)
		v, count := decodeElement()
		if count == 0 {
			break
		}
		data = append(data, uint8(v), uint8(v>>8), uint8(v>>16), uint8(v>>24),
			uint8(v>>32), uint8(v>>40), uint8(v>>48), uint8(v>>56))
	}
	_this.finishTypedArray(events.ArrayTypeUint64, digitType, 8, data)
}

func (_this *OldDecoder) decodeArrayI8(digitType string, decodeElement func() (v int64, digitCount int)) {
	var data []uint8
	for {
		_this.buffer.ReadWhilePropertyNoEOD(chars.CharIsWhitespace)
		v, count := decodeElement()
		if count == 0 {
			break
		}
		if v < minInt8Value || v > maxInt8Value {
			_this.buffer.Errorf("%v value too big for array type", digitType)
		}
		data = append(data, uint8(v))
	}
	_this.finishTypedArray(events.ArrayTypeInt8, digitType, 1, data)
}

func (_this *OldDecoder) decodeArrayI16(digitType string, decodeElement func() (v int64, digitCount int)) {
	var data []uint8
	for {
		_this.buffer.ReadWhilePropertyNoEOD(chars.CharIsWhitespace)
		v, count := decodeElement()
		if count == 0 {
			break
		}
		if v < minInt16Value || v > maxInt16Value {
			_this.buffer.Errorf("%v value too big for array type", digitType)
		}
		data = append(data, uint8(v), uint8(v>>8))
	}
	_this.finishTypedArray(events.ArrayTypeInt16, digitType, 2, data)
}

func (_this *OldDecoder) decodeArrayI32(digitType string, decodeElement func() (v int64, digitCount int)) {
	var data []uint8
	for {
		_this.buffer.ReadWhilePropertyNoEOD(chars.CharIsWhitespace)
		v, count := decodeElement()
		if count == 0 {
			break
		}
		if v < minInt32Value || v > maxInt32Value {
			_this.buffer.Errorf("%v value too big for array type", digitType)
		}
		data = append(data, uint8(v), uint8(v>>8), uint8(v>>16), uint8(v>>24))
	}
	_this.finishTypedArray(events.ArrayTypeInt32, digitType, 4, data)
}

func (_this *OldDecoder) decodeArrayI64(digitType string, decodeElement func() (v int64, digitCount int)) {
	var data []uint8
	for {
		_this.buffer.ReadWhilePropertyNoEOD(chars.CharIsWhitespace)
		v, count := decodeElement()
		if count == 0 {
			break
		}
		data = append(data, uint8(v), uint8(v>>8), uint8(v>>16), uint8(v>>24),
			uint8(v>>32), uint8(v>>40), uint8(v>>48), uint8(v>>56))
	}
	_this.finishTypedArray(events.ArrayTypeInt64, digitType, 8, data)
}

func (_this *OldDecoder) decodeArrayF16(digitType string, decodeElement func() (v float64, digitCount int)) {
	var data []uint8
	for {
		_this.buffer.ReadWhilePropertyNoEOD(chars.CharIsWhitespace)
		v, count := decodeElement()
		if count == 0 {
			break
		}

		exp := extractFloat64Exponent(v)
		if exp < minFloat32Exponent || exp > maxFloat32Exponent {
			_this.buffer.Errorf("Exponent too big for bfloat16 type")
		}
		bits := math.Float32bits(float32(v))
		data = append(data, uint8(bits>>16), uint8(bits>>24))
	}
	_this.finishTypedArray(events.ArrayTypeFloat16, digitType, 2, data)
}

func (_this *OldDecoder) decodeArrayF32(digitType string, decodeElement func() (v float64, digitCount int)) {
	var data []uint8
	for {
		_this.buffer.ReadWhilePropertyNoEOD(chars.CharIsWhitespace)
		v, count := decodeElement()
		if count == 0 {
			break
		}

		exp := extractFloat64Exponent(v)
		if exp < minFloat32Exponent || exp > maxFloat32Exponent {
			_this.buffer.Errorf("Exponent too big for float32 type")
		}
		bits := math.Float32bits(float32(v))
		data = append(data, uint8(bits), uint8(bits>>8), uint8(bits>>16), uint8(bits>>24))
	}
	_this.finishTypedArray(events.ArrayTypeFloat32, digitType, 4, data)
}

func (_this *OldDecoder) decodeArrayF64(digitType string, decodeElement func() (v float64, digitCount int)) {
	var data []uint8
	for {
		_this.buffer.ReadWhilePropertyNoEOD(chars.CharIsWhitespace)
		v, count := decodeElement()
		if count == 0 {
			break
		}
		bits := math.Float64bits(v)
		data = append(data, uint8(bits), uint8(bits>>8), uint8(bits>>16), uint8(bits>>24),
			uint8(bits>>32), uint8(bits>>40), uint8(bits>>48), uint8(bits>>56))
	}
	_this.finishTypedArray(events.ArrayTypeFloat64, digitType, 8, data)
}

func (_this *OldDecoder) readArrayType() string {
	_this.buffer.BeginToken()
	_this.buffer.ReadUntilPropertyNoEOD(chars.CharIsObjectEnd)
	arrayType := _this.buffer.GetToken()
	if len(arrayType) > 0 && arrayType[len(arrayType)-1] == '|' {
		arrayType = arrayType[:len(arrayType)-1]
		_this.buffer.UngetByte()
	}
	common.ASCIIBytesToLower(arrayType)
	return string(arrayType)
}

func (_this *OldDecoder) handleTypedArrayBegin() {
	_this.buffer.AdvanceByte()
	arrayType := _this.readArrayType()
	switch arrayType {
	case "cb":
		_this.decodeCustomBinary()
	case "ct":
		_this.decodeCustomText()
	case "r":
		_this.decodeRID()
	case "u":
		panic("TODO: CTEDecoder: UUID array")
	case "u8":
		_this.decodeArrayU8("integer", _this.buffer.DecodeSmallUint)
	case "u8b":
		_this.decodeArrayU8("binary", _this.buffer.DecodeSmallBinaryUint)
	case "u8o":
		_this.decodeArrayU8("octal", _this.buffer.DecodeSmallOctalUint)
	case "u8x":
		_this.decodeArrayU8("hex", _this.buffer.DecodeSmallHexUint)
	case "u16":
		_this.decodeArrayU16("integer", _this.buffer.DecodeSmallUint)
	case "u16b":
		_this.decodeArrayU16("binary", _this.buffer.DecodeSmallBinaryUint)
	case "u16o":
		_this.decodeArrayU16("octal", _this.buffer.DecodeSmallOctalUint)
	case "u16x":
		_this.decodeArrayU16("hex", _this.buffer.DecodeSmallHexUint)
	case "u32":
		_this.decodeArrayU32("integer", _this.buffer.DecodeSmallUint)
	case "u32b":
		_this.decodeArrayU32("binary", _this.buffer.DecodeSmallBinaryUint)
	case "u32o":
		_this.decodeArrayU32("octal", _this.buffer.DecodeSmallOctalUint)
	case "u32x":
		_this.decodeArrayU32("hex", _this.buffer.DecodeSmallHexUint)
	case "u64":
		_this.decodeArrayU64("integer", _this.buffer.DecodeSmallUint)
	case "u64b":
		_this.decodeArrayU64("binary", _this.buffer.DecodeSmallBinaryUint)
	case "u64o":
		_this.decodeArrayU64("octal", _this.buffer.DecodeSmallOctalUint)
	case "u64x":
		_this.decodeArrayU64("hex", _this.buffer.DecodeSmallHexUint)
	case "i8":
		_this.decodeArrayI8("integer", _this.buffer.DecodeSmallInt)
	case "i8b":
		_this.decodeArrayI8("binary", _this.buffer.DecodeSmallBinaryInt)
	case "i8o":
		_this.decodeArrayI8("octal", _this.buffer.DecodeSmallOctalInt)
	case "i8x":
		_this.decodeArrayI8("hex", _this.buffer.DecodeSmallHexInt)
	case "i16":
		_this.decodeArrayI16("integer", _this.buffer.DecodeSmallInt)
	case "i16b":
		_this.decodeArrayI16("binary", _this.buffer.DecodeSmallBinaryInt)
	case "i16o":
		_this.decodeArrayI16("octal", _this.buffer.DecodeSmallOctalInt)
	case "i16x":
		_this.decodeArrayI16("hex", _this.buffer.DecodeSmallHexInt)
	case "i32":
		_this.decodeArrayI32("integer", _this.buffer.DecodeSmallInt)
	case "i32b":
		_this.decodeArrayI32("binary", _this.buffer.DecodeSmallBinaryInt)
	case "i32o":
		_this.decodeArrayI32("octal", _this.buffer.DecodeSmallOctalInt)
	case "i32x":
		_this.decodeArrayI32("hex", _this.buffer.DecodeSmallHexInt)
	case "i64":
		_this.decodeArrayI64("integer", _this.buffer.DecodeSmallInt)
	case "i64b":
		_this.decodeArrayI64("binary", _this.buffer.DecodeSmallBinaryInt)
	case "i64o":
		_this.decodeArrayI64("octal", _this.buffer.DecodeSmallOctalInt)
	case "i64x":
		_this.decodeArrayI64("hex", _this.buffer.DecodeSmallHexInt)
	case "f16":
		_this.decodeArrayF16("float", _this.buffer.DecodeSmallFloat)
	case "f16x":
		_this.decodeArrayF16("hex float", _this.buffer.DecodeSmallHexFloat)
	case "f32":
		_this.decodeArrayF32("float", _this.buffer.DecodeSmallFloat)
	case "f32x":
		_this.decodeArrayF32("hex float", _this.buffer.DecodeSmallHexFloat)
	case "f64":
		_this.decodeArrayF64("float", _this.buffer.DecodeSmallFloat)
	case "f64x":
		_this.decodeArrayF64("hex float", _this.buffer.DecodeSmallHexFloat)
	default:
		_this.buffer.Errorf("%s: Unhandled array type", arrayType)
	}
}

func (_this *OldDecoder) handleReference() {
	_this.buffer.AdvanceByte()
	_this.eventReceiver.OnReference()
	if _this.buffer.PeekByteNoEOD() == '|' {
		_this.buffer.AdvanceByte()
		arrayType := _this.readArrayType()
		if arrayType != "r" {
			_this.buffer.Errorf("%s: Invalid array type for reference ID", arrayType)
		}
		_this.decodeRID()
		return
	}

	asString, asUint := _this.buffer.DecodeMarkerID()
	if len(asString) > 0 {
		_this.eventReceiver.OnArray(events.ArrayTypeString, uint64(len(asString)), asString)
	} else {
		_this.eventReceiver.OnPositiveInt(asUint)
	}
	_this.endObject()
}

func (_this *OldDecoder) handleMarker() {
	_this.buffer.AdvanceByte()
	_this.eventReceiver.OnMarker()
	asString, asUint := _this.buffer.DecodeMarkerID()
	if b := _this.buffer.PeekByteNoEOD(); b != ':' {
		_this.buffer.UnexpectedChar("marker ID")
	}
	_this.buffer.AdvanceByte()
	if len(asString) > 0 {
		_this.eventReceiver.OnArray(events.ArrayTypeString, uint64(len(asString)), asString)
	} else {
		_this.eventReceiver.OnPositiveInt(asUint)
	}
	// Don't end object here because the real object follows the marker ID
}

type cteDecoderState int

func (_this cteDecoderState) String() string {
	return cteDecoderStateNames[_this]
}

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
	cteDecoderStateAwaitReferenceID
	cteDecoderStateCount
)

var cteDecoderStateNames = []string{
	cteDecoderStateAwaitObject:            "cteDecoderStateAwaitObject",
	cteDecoderStateAwaitListItem:          "cteDecoderStateAwaitListItem",
	cteDecoderStateAwaitCommentItem:       "cteDecoderStateAwaitCommentItem",
	cteDecoderStateAwaitMapKey:            "cteDecoderStateAwaitMapKey",
	cteDecoderStateAwaitMapKVSeparator:    "cteDecoderStateAwaitMapKVSeparator",
	cteDecoderStateAwaitMapValue:          "cteDecoderStateAwaitMapValue",
	cteDecoderStateAwaitMetaKey:           "cteDecoderStateAwaitMetaKey",
	cteDecoderStateAwaitMetaKVSeparator:   "cteDecoderStateAwaitMetaKVSeparator",
	cteDecoderStateAwaitMetaValue:         "cteDecoderStateAwaitMetaValue",
	cteDecoderStateAwaitMarkupName:        "cteDecoderStateAwaitMarkupName",
	cteDecoderStateAwaitMarkupKey:         "cteDecoderStateAwaitMarkupKey",
	cteDecoderStateAwaitMarkupKVSeparator: "cteDecoderStateAwaitMarkupKVSeparator",
	cteDecoderStateAwaitMarkupValue:       "cteDecoderStateAwaitMarkupValue",
	cteDecoderStateAwaitMarkupItem:        "cteDecoderStateAwaitMarkupItem",
	cteDecoderStateAwaitReferenceID:       "cteDecoderStateAwaitReferenceID",
}

var cteDecoderStateTransitions = [cteDecoderStateCount]cteDecoderState{
	cteDecoderStateAwaitObject:            cteDecoderStateAwaitObject,
	cteDecoderStateAwaitListItem:          cteDecoderStateAwaitListItem,
	cteDecoderStateAwaitCommentItem:       cteDecoderStateAwaitCommentItem,
	cteDecoderStateAwaitMapKey:            cteDecoderStateAwaitMapKVSeparator,
	cteDecoderStateAwaitMapKVSeparator:    cteDecoderStateAwaitMapValue,
	cteDecoderStateAwaitMapValue:          cteDecoderStateAwaitMapKey,
	cteDecoderStateAwaitMetaKey:           cteDecoderStateAwaitMetaKVSeparator,
	cteDecoderStateAwaitMetaKVSeparator:   cteDecoderStateAwaitMetaValue,
	cteDecoderStateAwaitMetaValue:         cteDecoderStateAwaitMetaKey,
	cteDecoderStateAwaitMarkupName:        cteDecoderStateAwaitMarkupKey,
	cteDecoderStateAwaitMarkupKey:         cteDecoderStateAwaitMarkupKVSeparator,
	cteDecoderStateAwaitMarkupKVSeparator: cteDecoderStateAwaitMarkupValue,
	cteDecoderStateAwaitMarkupValue:       cteDecoderStateAwaitMarkupKey,
	cteDecoderStateAwaitMarkupItem:        cteDecoderStateAwaitMarkupItem,
	cteDecoderStateAwaitReferenceID:       cteDecoderStateAwaitObject,
}

var cteDecoderStateHandlers = [cteDecoderStateCount]cteDecoderHandlerFunction{
	cteDecoderStateAwaitObject:            (*OldDecoder).handleObject,
	cteDecoderStateAwaitListItem:          (*OldDecoder).handleObject,
	cteDecoderStateAwaitCommentItem:       (*OldDecoder).handleCommentContent,
	cteDecoderStateAwaitMapKey:            (*OldDecoder).handleObject,
	cteDecoderStateAwaitMapKVSeparator:    (*OldDecoder).handleKVSeparator,
	cteDecoderStateAwaitMapValue:          (*OldDecoder).handleObject,
	cteDecoderStateAwaitMetaKey:           (*OldDecoder).handleObject,
	cteDecoderStateAwaitMetaKVSeparator:   (*OldDecoder).handleKVSeparator,
	cteDecoderStateAwaitMetaValue:         (*OldDecoder).handleObject,
	cteDecoderStateAwaitMarkupName:        (*OldDecoder).handleObject,
	cteDecoderStateAwaitMarkupKey:         (*OldDecoder).handleObject,
	cteDecoderStateAwaitMarkupKVSeparator: (*OldDecoder).handleKVSeparator,
	cteDecoderStateAwaitMarkupValue:       (*OldDecoder).handleObject,
	cteDecoderStateAwaitMarkupItem:        (*OldDecoder).handleMarkupContent,
	cteDecoderStateAwaitReferenceID:       (*OldDecoder).handleObject,
}

var charBasedHandlers [0x101]cteDecoderHandlerFunction

func init() {
	for i := 0; i < 0x100; i++ {
		charBasedHandlers[i] = (*OldDecoder).handleInvalidChar
	}

	charBasedHandlers['\r'] = (*OldDecoder).handleWhitespace
	charBasedHandlers['\n'] = (*OldDecoder).handleWhitespace
	charBasedHandlers['\t'] = (*OldDecoder).handleWhitespace
	charBasedHandlers[' '] = (*OldDecoder).handleWhitespace

	charBasedHandlers['!'] = (*OldDecoder).handleInvalidChar
	charBasedHandlers['"'] = (*OldDecoder).handleQuotedString
	charBasedHandlers['#'] = (*OldDecoder).handleConstant
	charBasedHandlers['$'] = (*OldDecoder).handleReference
	charBasedHandlers['%'] = (*OldDecoder).handleInvalidChar
	charBasedHandlers['&'] = (*OldDecoder).handleMarker
	charBasedHandlers['\''] = (*OldDecoder).handleInvalidChar
	charBasedHandlers['('] = (*OldDecoder).handleMetadataBegin
	charBasedHandlers[')'] = (*OldDecoder).handleMetadataEnd
	charBasedHandlers['+'] = (*OldDecoder).handleInvalidChar
	charBasedHandlers[','] = (*OldDecoder).handleMarkupContentBegin
	charBasedHandlers['-'] = (*OldDecoder).handleNegativeNumeric
	charBasedHandlers['.'] = (*OldDecoder).handleInvalidChar
	charBasedHandlers['/'] = (*OldDecoder).handleComment

	charBasedHandlers['0'] = (*OldDecoder).handleOtherBasePositive
	for i := '1'; i <= '9'; i++ {
		charBasedHandlers[i] = (*OldDecoder).handlePositiveNumeric
	}

	charBasedHandlers[':'] = (*OldDecoder).handleInvalidChar
	charBasedHandlers[';'] = (*OldDecoder).handleInvalidChar
	charBasedHandlers['<'] = (*OldDecoder).handleMarkupBegin
	charBasedHandlers['>'] = (*OldDecoder).handleMarkupEnd
	charBasedHandlers['?'] = (*OldDecoder).handleInvalidChar
	charBasedHandlers['@'] = (*OldDecoder).handleNamedValueOrUUID

	for i := 'A'; i <= 'Z'; i++ {
		charBasedHandlers[i] = (*OldDecoder).handleUnquotedString
	}

	charBasedHandlers['['] = (*OldDecoder).handleListBegin
	charBasedHandlers['\\'] = (*OldDecoder).handleInvalidChar
	charBasedHandlers[']'] = (*OldDecoder).handleListEnd
	charBasedHandlers['^'] = (*OldDecoder).handleInvalidChar
	charBasedHandlers['_'] = (*OldDecoder).handleUnquotedString

	for i := 'a'; i <= 'z'; i++ {
		charBasedHandlers[i] = (*OldDecoder).handleUnquotedString
	}

	charBasedHandlers['{'] = (*OldDecoder).handleMapBegin
	charBasedHandlers['|'] = (*OldDecoder).handleTypedArrayBegin
	charBasedHandlers['}'] = (*OldDecoder).handleMapEnd
	charBasedHandlers['~'] = (*OldDecoder).handleInvalidChar

	for i := 0xc0; i < 0xf8; i++ {
		charBasedHandlers[i] = (*OldDecoder).handleUnquotedString
	}

	charBasedHandlers[chars.EndOfDocumentMarker] = (*OldDecoder).handleNothing
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
