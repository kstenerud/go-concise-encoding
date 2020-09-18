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

	"github.com/kstenerud/go-concise-encoding/debug"
	"github.com/kstenerud/go-concise-encoding/events"
	"github.com/kstenerud/go-concise-encoding/internal/common"
	"github.com/kstenerud/go-concise-encoding/options"
)

// Decodes CTE documents.
type Decoder struct {
	buffer         ReadBuffer
	eventReceiver  events.DataEventReceiver
	containerState []cteDecoderState
	currentState   cteDecoderState
	opts           options.CTEDecoderOptions
}

// Create a new CTE decoder, which will read from reader and send data events
// to nextReceiver. If opts is nil, default options will be used.
func NewDecoder(opts *options.CTEDecoderOptions) *Decoder {
	_this := &Decoder{}
	_this.Init(opts)
	return _this
}

// Initialize this decoder, which will read from reader and send data events
// to nextReceiver. If opts is nil, default options will be used.
func (_this *Decoder) Init(opts *options.CTEDecoderOptions) {
	opts = opts.WithDefaultsApplied()
	_this.opts = *opts
}

func (_this *Decoder) reset() {
	_this.buffer.Reset()
	_this.eventReceiver = nil
	_this.containerState = _this.containerState[:0]
	_this.currentState = 0
}

// Run the complete decode process. The document and data receiver specified
// when initializing the decoder will be used.
func (_this *Decoder) Decode(reader io.Reader, eventReceiver events.DataEventReceiver) (err error) {
	defer func() {
		_this.reset()
		if !debug.DebugOptions.PassThroughPanics {
			if r := recover(); r != nil {
				err = r.(error)
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
	_this.buffer.EndToken()

	switch _this.opts.ImpliedStructure {
	case options.ImpliedStructureVersion:
		_this.eventReceiver.OnVersion(_this.opts.ConciseEncodingVersion)
	case options.ImpliedStructureList:
		_this.eventReceiver.OnVersion(_this.opts.ConciseEncodingVersion)
		_this.eventReceiver.OnList()
		_this.stackContainer(cteDecoderStateAwaitListItem)
	case options.ImpliedStructureMap:
		_this.eventReceiver.OnVersion(_this.opts.ConciseEncodingVersion)
		_this.eventReceiver.OnMap()
		_this.stackContainer(cteDecoderStateAwaitMapKey)
	default:
		_this.handleVersion()
	}

	for !_this.buffer.IsEndOfDocument() {
		_this.handleNextState()
		_this.buffer.RefillIfNecessary()
	}

	switch _this.opts.ImpliedStructure {
	case options.ImpliedStructureList, options.ImpliedStructureMap:
		_this.eventReceiver.OnEnd()
	}

	_this.eventReceiver.OnEndDocument()
	return
}

func (_this *Decoder) DecodeDocument(document []byte, eventReceiver events.DataEventReceiver) (err error) {
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

// Tokens

func (_this *Decoder) endObject() {
	_this.buffer.EndToken()
	_this.transitionToNextState()
}

// State

func (_this *Decoder) stackContainer(nextState cteDecoderState) {
	_this.containerState = append(_this.containerState, _this.currentState)
	_this.currentState = nextState
}

func (_this *Decoder) unstackContainer() {
	index := len(_this.containerState) - 1
	_this.currentState = _this.containerState[index]
	_this.containerState = _this.containerState[:index]
}

func (_this *Decoder) changeState(newState cteDecoderState) {
	_this.currentState = newState
}

func (_this *Decoder) transitionToNextState() {
	_this.currentState = cteDecoderStateTransitions[_this.currentState]
}

// Handlers

type cteDecoderHandlerFunction func(*Decoder)

func (_this *Decoder) handleNothing() {
}

func (_this *Decoder) handleNextState() {
	cteDecoderStateHandlers[_this.currentState](_this)
}

func (_this *Decoder) handleObject() {
	charBasedHandlers[_this.buffer.PeekByteAllowEOD()](_this)
}

func (_this *Decoder) handleInvalidChar() {
	_this.buffer.Errorf("Unexpected [%v]", _this.buffer.DescribeCurrentChar())
}

func (_this *Decoder) handleInvalidState() {
	_this.buffer.Errorf("BUG: Invalid state: %v", _this.currentState)
}

func (_this *Decoder) handleKVSeparator() {
	_this.buffer.SkipWhitespace()
	if _this.buffer.PeekByteNoEOD() != '=' {
		_this.buffer.Errorf("Expected map separator (=) but got [%v]", _this.buffer.DescribeCurrentChar())
	}
	_this.buffer.AdvanceByte()
	_this.buffer.SkipWhitespace()
	_this.buffer.EndToken()
	_this.endObject()
}

func (_this *Decoder) handleWhitespace() {
	_this.buffer.SkipWhitespace()
	_this.buffer.EndToken()
}

func (_this *Decoder) handleVersion() {
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

	if !hasProperty(_this.buffer.PeekByteNoEOD(), ctePropertyWhitespace) {
		_this.buffer.UnexpectedChar("whitespace after version")
	}
	_this.buffer.AdvanceByte()

	_this.eventReceiver.OnVersion(version)
	_this.buffer.EndToken()
}

func (_this *Decoder) handleStringish() {
	_this.buffer.ReadWhilePropertyAllowEOD(ctePropertyUnquotedMid)

	// Unquoted string
	if _this.buffer.PeekByteAllowEOD().HasProperty(ctePropertyObjectEnd) {
		bytes := _this.buffer.GetToken()
		_this.eventReceiver.OnArray(events.ArrayTypeString, uint64(len(bytes)), bytes)
		_this.endObject()
		return
	}

	// Bytes, Custom, URI
	if _this.buffer.GetTokenLength() == 1 && _this.buffer.PeekByteNoEOD() == '"' {
		// TODO: array chunking on big data instead of building a big slice
		_this.buffer.AdvanceByte()
		initiator := _this.buffer.GetTokenFirstByte()
		switch initiator {
		case 'b':
			bytes := _this.buffer.DecodeHexBytes()
			_this.eventReceiver.OnArray(events.ArrayTypeCustomBinary, uint64(len(bytes)), bytes)
			_this.endObject()
			return
		case 't':
			bytes := _this.buffer.DecodeCustomText()
			_this.eventReceiver.OnArray(events.ArrayTypeCustomText, uint64(len(bytes)), bytes)
			_this.endObject()
			return
		case 'u':
			bytes := _this.buffer.DecodeURI()
			_this.eventReceiver.OnArray(events.ArrayTypeURI, uint64(len(bytes)), bytes)
			_this.endObject()
			return
		default:
			_this.buffer.UngetByte()
			_this.buffer.UnexpectedChar("byte array initiator")
		}
	}

	_this.buffer.UnexpectedChar("unquoted string")
}

func (_this *Decoder) handleQuotedString() {
	_this.buffer.AdvanceByte()
	bytes := _this.buffer.DecodeQuotedString()
	_this.eventReceiver.OnArray(events.ArrayTypeString, uint64(len(bytes)), bytes)
	_this.endObject()
}

func (_this *Decoder) handlePositiveNumeric() {
	coefficient, bigCoefficient, digitCount := _this.buffer.DecodeDecimalUint(0, nil)

	// Integer
	if _this.buffer.PeekByteAllowEOD().HasProperty(ctePropertyObjectEnd) {
		if bigCoefficient != nil {
			_this.eventReceiver.OnBigInt(bigCoefficient)
		} else {
			_this.eventReceiver.OnPositiveInt(coefficient)
		}
		_this.endObject()
		return
	}

	switch _this.buffer.ReadByte() {
	case '-':
		v := _this.buffer.DecodeDate(int64(coefficient))
		_this.buffer.AssertAtObjectEnd("date")
		_this.eventReceiver.OnCompactTime(v)
		_this.endObject()
	case ':':
		v := _this.buffer.DecodeTime(int(coefficient))
		_this.buffer.AssertAtObjectEnd("time")
		_this.eventReceiver.OnCompactTime(v)
		_this.endObject()
	case '.':
		value, bigValue, _ := _this.buffer.DecodeDecimalFloat(1, coefficient, bigCoefficient, digitCount)
		_this.buffer.AssertAtObjectEnd("float")
		if bigValue != nil {
			_this.eventReceiver.OnBigDecimalFloat(bigValue)
		} else {
			_this.eventReceiver.OnDecimalFloat(value)
		}
		_this.endObject()
	default:
		_this.buffer.UngetByte()
		_this.buffer.UnexpectedChar("numeric")
	}
}

func (_this *Decoder) handleNegativeNumeric() {
	_this.buffer.AdvanceByte()
	switch _this.buffer.PeekByteNoEOD() {
	case '0':
		_this.handleOtherBaseNegative()
		return
	case '@':
		_this.buffer.AdvanceByte()
		_this.buffer.BeginSubtoken()
		_this.buffer.ReadWhilePropertyAllowEOD(ctePropertyAZ)
		subtoken := _this.buffer.GetSubtoken()
		common.ASCIIBytesToLower(subtoken)
		token := string(subtoken)
		if token != "inf" {
			_this.buffer.Errorf("Unknown named value: %v", token)
		}
		_this.eventReceiver.OnFloat(math.Inf(-1))
		_this.endObject()
		return
	}

	coefficient, bigCoefficient, digitCount := _this.buffer.DecodeDecimalUint(0, nil)

	// Integer
	if _this.buffer.PeekByteAllowEOD().HasProperty(ctePropertyObjectEnd) {
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

	switch _this.buffer.ReadByte() {
	case '-':
		v := _this.buffer.DecodeDate(-int64(coefficient))
		_this.buffer.AssertAtObjectEnd("time")
		_this.eventReceiver.OnCompactTime(v)
		_this.endObject()
	case '.':
		value, bigValue, _ := _this.buffer.DecodeDecimalFloat(-1, coefficient, bigCoefficient, digitCount)
		_this.buffer.AssertAtObjectEnd("float")
		if bigValue != nil {
			_this.eventReceiver.OnBigDecimalFloat(bigValue)
		} else {
			_this.eventReceiver.OnDecimalFloat(value)
		}
		_this.endObject()
	default:
		_this.buffer.UngetByte()
		_this.buffer.UnexpectedChar("numeric")
	}
}

func (_this *Decoder) handleOtherBasePositive() {
	_this.buffer.AdvanceByte()
	b := _this.buffer.PeekByteAllowEOD()

	if b.HasProperty(ctePropertyObjectEnd) {
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
		if b.HasProperty(cteProperty09) && _this.buffer.PeekByteNoEOD() == ':' {
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

func (_this *Decoder) handleOtherBaseNegative() {
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

func (_this *Decoder) handleListBegin() {
	_this.buffer.AdvanceByte()
	_this.eventReceiver.OnList()
	_this.stackContainer(cteDecoderStateAwaitListItem)
	_this.buffer.EndToken()
}

func (_this *Decoder) handleListEnd() {
	_this.buffer.AdvanceByte()
	_this.eventReceiver.OnEnd()
	_this.unstackContainer()
	_this.endObject()
}

func (_this *Decoder) handleMapBegin() {
	_this.buffer.AdvanceByte()
	_this.eventReceiver.OnMap()
	_this.stackContainer(cteDecoderStateAwaitMapKey)
	_this.buffer.EndToken()
}

func (_this *Decoder) handleMapEnd() {
	_this.buffer.AdvanceByte()
	_this.eventReceiver.OnEnd()
	_this.unstackContainer()
	_this.endObject()
}

func (_this *Decoder) handleMetadataBegin() {
	_this.buffer.AdvanceByte()
	_this.eventReceiver.OnMetadata()
	_this.stackContainer(cteDecoderStateAwaitMetaKey)
	_this.buffer.EndToken()
}

func (_this *Decoder) handleMetadataEnd() {
	_this.buffer.AdvanceByte()
	_this.eventReceiver.OnEnd()
	_this.unstackContainer()
	_this.buffer.EndToken()
	// Don't transition state because metadata is a pseudo-object
}

func (_this *Decoder) handleComment() {
	_this.buffer.AdvanceByte()
	switch _this.buffer.ReadByte() {
	case '/':
		_this.eventReceiver.OnComment()
		contents := _this.buffer.DecodeSingleLineComment()
		if len(contents) > 0 {
			_this.eventReceiver.OnArray(events.ArrayTypeString, uint64(len(contents)), contents)
		}
		_this.eventReceiver.OnEnd()
		_this.buffer.EndToken()
	case '*':
		_this.eventReceiver.OnComment()
		_this.stackContainer(cteDecoderStateAwaitCommentItem)
		_this.buffer.EndToken()
	default:
		_this.buffer.UngetByte()
		_this.buffer.UnexpectedChar("comment")
	}
}

func (_this *Decoder) handleCommentContent() {
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

func (_this *Decoder) handleMarkupBegin() {
	_this.buffer.AdvanceByte()
	_this.eventReceiver.OnMarkup()
	_this.stackContainer(cteDecoderStateAwaitMarkupName)
	_this.buffer.EndToken()
}

func (_this *Decoder) handleMarkupContentBegin() {
	_this.buffer.AdvanceByte()
	_this.eventReceiver.OnEnd()
	_this.changeState(cteDecoderStateAwaitMarkupItem)
	_this.buffer.EndToken()
}

func (_this *Decoder) handleMarkupContent() {
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
	case nextIsVerbatimString:
		str := _this.buffer.DecodeVerbatimString()
		if len(str) > 0 {
			_this.eventReceiver.OnArray(events.ArrayTypeVerbatimString, uint64(len(str)), str)
		}
		_this.buffer.EndToken()
	case nextIsSingleLineComment:
		_this.eventReceiver.OnComment()
		contents := _this.buffer.DecodeSingleLineComment()
		if len(contents) > 0 {
			_this.eventReceiver.OnArray(events.ArrayTypeString, uint64(len(contents)), contents)
		}
		_this.eventReceiver.OnEnd()
		_this.buffer.EndToken()
	case nextIsMarkupBegin:
		_this.eventReceiver.OnMarkup()
		_this.stackContainer(cteDecoderStateAwaitMarkupValue)
		_this.buffer.EndToken()
	case nextIsMarkupEnd:
		_this.eventReceiver.OnEnd()
		_this.unstackContainer()
		_this.endObject()
	}
}

func (_this *Decoder) handleMarkupEnd() {
	_this.buffer.AdvanceByte()
	_this.eventReceiver.OnEnd()
	_this.eventReceiver.OnEnd()
	_this.unstackContainer()
	_this.endObject()
}

func (_this *Decoder) handleNamedValue() {
	_this.buffer.AdvanceByte()
	_this.buffer.BeginSubtoken()
	_this.buffer.ReadWhilePropertyAllowEOD(ctePropertyAZ)
	subtoken := _this.buffer.GetSubtoken()
	common.ASCIIBytesToLower(subtoken)
	token := string(subtoken)
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

	_this.buffer.UngetAll()
	_this.buffer.AdvanceByte()
	_this.eventReceiver.OnUUID(_this.buffer.DecodeUUID())
	_this.endObject()
}

func (_this *Decoder) handleVerbatimString() {
	_this.buffer.AdvanceByte()
	bytes := _this.buffer.DecodeVerbatimString()
	_this.eventReceiver.OnArray(events.ArrayTypeVerbatimString, uint64(len(bytes)), bytes)
	_this.endObject()
}

func (_this *Decoder) finishTypedArray(arrayType events.ArrayType, digitType string, bytesPerElement int, data []byte) {
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

func (_this *Decoder) decodeArrayU8(digitType string, decodeElement func() (v uint64, digitCount int)) {
	var data []uint8
	for {
		_this.buffer.ReadWhilePropertyNoEOD(ctePropertyWhitespace)
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

func (_this *Decoder) decodeArrayU16(digitType string, decodeElement func() (v uint64, digitCount int)) {
	var data []uint8
	for {
		_this.buffer.ReadWhilePropertyNoEOD(ctePropertyWhitespace)
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

func (_this *Decoder) decodeArrayU32(digitType string, decodeElement func() (v uint64, digitCount int)) {
	var data []uint8
	for {
		_this.buffer.ReadWhilePropertyNoEOD(ctePropertyWhitespace)
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

func (_this *Decoder) decodeArrayU64(digitType string, decodeElement func() (v uint64, digitCount int)) {
	var data []uint8
	for {
		_this.buffer.ReadWhilePropertyNoEOD(ctePropertyWhitespace)
		v, count := decodeElement()
		if count == 0 {
			break
		}
		data = append(data, uint8(v), uint8(v>>8), uint8(v>>16), uint8(v>>24),
			uint8(v>>32), uint8(v>>40), uint8(v>>48), uint8(v>>56))
	}
	_this.finishTypedArray(events.ArrayTypeUint64, digitType, 8, data)
}

func (_this *Decoder) decodeArrayI8(digitType string, decodeElement func() (v int64, digitCount int)) {
	var data []uint8
	for {
		_this.buffer.ReadWhilePropertyNoEOD(ctePropertyWhitespace)
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

func (_this *Decoder) decodeArrayI16(digitType string, decodeElement func() (v int64, digitCount int)) {
	var data []uint8
	for {
		_this.buffer.ReadWhilePropertyNoEOD(ctePropertyWhitespace)
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

func (_this *Decoder) decodeArrayI32(digitType string, decodeElement func() (v int64, digitCount int)) {
	var data []uint8
	for {
		_this.buffer.ReadWhilePropertyNoEOD(ctePropertyWhitespace)
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

func (_this *Decoder) decodeArrayI64(digitType string, decodeElement func() (v int64, digitCount int)) {
	var data []uint8
	for {
		_this.buffer.ReadWhilePropertyNoEOD(ctePropertyWhitespace)
		v, count := decodeElement()
		if count == 0 {
			break
		}
		data = append(data, uint8(v), uint8(v>>8), uint8(v>>16), uint8(v>>24),
			uint8(v>>32), uint8(v>>40), uint8(v>>48), uint8(v>>56))
	}
	_this.finishTypedArray(events.ArrayTypeInt64, digitType, 8, data)
}

func (_this *Decoder) decodeArrayF16(digitType string, decodeElement func() (v float64, digitCount int)) {
	var data []uint8
	for {
		_this.buffer.ReadWhilePropertyNoEOD(ctePropertyWhitespace)
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

func (_this *Decoder) decodeArrayF32(digitType string, decodeElement func() (v float64, digitCount int)) {
	var data []uint8
	for {
		_this.buffer.ReadWhilePropertyNoEOD(ctePropertyWhitespace)
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

func (_this *Decoder) decodeArrayF64(digitType string, decodeElement func() (v float64, digitCount int)) {
	var data []uint8
	for {
		_this.buffer.ReadWhilePropertyNoEOD(ctePropertyWhitespace)
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

func (_this *Decoder) handleTypedArrayBegin() {
	_this.buffer.AdvanceByte()
	_this.buffer.BeginSubtoken()
	_this.buffer.ReadUntilPropertyNoEOD(ctePropertyWhitespace)
	subtoken := _this.buffer.GetSubtoken()
	if len(subtoken) > 0 && subtoken[len(subtoken)-1] == '|' {
		subtoken = subtoken[:len(subtoken)-1]
		_this.buffer.UngetByte()
	}
	common.ASCIIBytesToLower(subtoken)
	token := string(subtoken)
	switch token {
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
		panic(fmt.Errorf("TODO: Typed array decoder support for %s", token))
		_this.buffer.Errorf("%s: Unhandled array type", token)
	}
}

func (_this *Decoder) handleReference() {
	_this.buffer.AdvanceByte()
	_this.eventReceiver.OnReference()
	if _this.buffer.PeekByteNoEOD() == 'u' {
		_this.buffer.AdvanceByte()
		if _this.buffer.PeekByteNoEOD() != '"' {
			_this.buffer.UnexpectedChar("reference (uri)")
		}
		_this.buffer.AdvanceByte()
		bytes := _this.buffer.DecodeQuotedString()
		_this.eventReceiver.OnArray(events.ArrayTypeURI, uint64(len(bytes)), bytes)
		_this.endObject()
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

func (_this *Decoder) handleMarker() {
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
	_this.buffer.EndToken()
	// Don't end object here because the real object follows the marker ID
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
	cteDecoderStateAwaitReferenceID
	cteDecoderStateCount
)

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
	cteDecoderStateAwaitObject:            (*Decoder).handleObject,
	cteDecoderStateAwaitListItem:          (*Decoder).handleObject,
	cteDecoderStateAwaitCommentItem:       (*Decoder).handleCommentContent,
	cteDecoderStateAwaitMapKey:            (*Decoder).handleObject,
	cteDecoderStateAwaitMapKVSeparator:    (*Decoder).handleKVSeparator,
	cteDecoderStateAwaitMapValue:          (*Decoder).handleObject,
	cteDecoderStateAwaitMetaKey:           (*Decoder).handleObject,
	cteDecoderStateAwaitMetaKVSeparator:   (*Decoder).handleKVSeparator,
	cteDecoderStateAwaitMetaValue:         (*Decoder).handleObject,
	cteDecoderStateAwaitMarkupName:        (*Decoder).handleObject,
	cteDecoderStateAwaitMarkupKey:         (*Decoder).handleObject,
	cteDecoderStateAwaitMarkupKVSeparator: (*Decoder).handleKVSeparator,
	cteDecoderStateAwaitMarkupValue:       (*Decoder).handleObject,
	cteDecoderStateAwaitMarkupItem:        (*Decoder).handleMarkupContent,
	cteDecoderStateAwaitReferenceID:       (*Decoder).handleObject,
}

var charBasedHandlers [cteByteEndOfDocument + 1]cteDecoderHandlerFunction

func init() {
	for i := 0; i < cteByteEndOfDocument; i++ {
		charBasedHandlers[i] = (*Decoder).handleInvalidChar
	}

	charBasedHandlers['\r'] = (*Decoder).handleWhitespace
	charBasedHandlers['\n'] = (*Decoder).handleWhitespace
	charBasedHandlers['\t'] = (*Decoder).handleWhitespace
	charBasedHandlers[' '] = (*Decoder).handleWhitespace

	charBasedHandlers['!'] = (*Decoder).handleInvalidChar
	charBasedHandlers['"'] = (*Decoder).handleQuotedString
	charBasedHandlers['#'] = (*Decoder).handleInvalidChar
	charBasedHandlers['$'] = (*Decoder).handleReference
	charBasedHandlers['%'] = (*Decoder).handleInvalidChar
	charBasedHandlers['&'] = (*Decoder).handleMarker
	charBasedHandlers['\''] = (*Decoder).handleInvalidChar
	charBasedHandlers['('] = (*Decoder).handleMetadataBegin
	charBasedHandlers[')'] = (*Decoder).handleMetadataEnd
	charBasedHandlers['+'] = (*Decoder).handleInvalidChar
	charBasedHandlers[','] = (*Decoder).handleInvalidChar
	charBasedHandlers['-'] = (*Decoder).handleNegativeNumeric
	charBasedHandlers['.'] = (*Decoder).handleInvalidChar
	charBasedHandlers['/'] = (*Decoder).handleComment

	charBasedHandlers['0'] = (*Decoder).handleOtherBasePositive
	for i := '1'; i <= '9'; i++ {
		charBasedHandlers[i] = (*Decoder).handlePositiveNumeric
	}

	charBasedHandlers[':'] = (*Decoder).handleInvalidChar
	charBasedHandlers[';'] = (*Decoder).handleMarkupContentBegin
	charBasedHandlers['<'] = (*Decoder).handleMarkupBegin
	charBasedHandlers['>'] = (*Decoder).handleMarkupEnd
	charBasedHandlers['?'] = (*Decoder).handleInvalidChar
	charBasedHandlers['@'] = (*Decoder).handleNamedValue

	for i := 'A'; i <= 'Z'; i++ {
		charBasedHandlers[i] = (*Decoder).handleStringish
	}

	charBasedHandlers['['] = (*Decoder).handleListBegin
	charBasedHandlers['\\'] = (*Decoder).handleInvalidChar
	charBasedHandlers[']'] = (*Decoder).handleListEnd
	charBasedHandlers['^'] = (*Decoder).handleInvalidChar
	charBasedHandlers['_'] = (*Decoder).handleStringish
	charBasedHandlers['`'] = (*Decoder).handleVerbatimString

	for i := 'a'; i <= 'z'; i++ {
		charBasedHandlers[i] = (*Decoder).handleStringish
	}

	charBasedHandlers['{'] = (*Decoder).handleMapBegin
	charBasedHandlers['|'] = (*Decoder).handleTypedArrayBegin
	charBasedHandlers['}'] = (*Decoder).handleMapEnd
	charBasedHandlers['~'] = (*Decoder).handleInvalidChar

	for i := 0xc0; i < 0xf8; i++ {
		charBasedHandlers[i] = (*Decoder).handleStringish
	}

	charBasedHandlers[cteByteEndOfDocument] = (*Decoder).handleNothing
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
	ctePropertyNumericWhitespace
	ctePropertyObjectEnd
	ctePropertyUnquotedStart
	ctePropertyUnquotedMid
	ctePropertyAZ
	cteProperty09
	ctePropertyLowercaseAF
	ctePropertyUppercaseAF
	ctePropertyBinaryDigit
	ctePropertyOctalDigit
	ctePropertyAreaLocation
	ctePropertyMarkerID
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

	cteByteProperties[':'] |= ctePropertyUnquotedMid
	cteByteProperties['-'] |= ctePropertyUnquotedMid | ctePropertyAreaLocation
	cteByteProperties['+'] |= ctePropertyAreaLocation
	cteByteProperties['.'] |= ctePropertyUnquotedMid
	cteByteProperties['_'] |= ctePropertyUnquotedMid | ctePropertyUnquotedStart | ctePropertyAreaLocation | ctePropertyMarkerID | ctePropertyNumericWhitespace
	for i := '0'; i <= '9'; i++ {
		cteByteProperties[i] |= ctePropertyUnquotedMid | ctePropertyAreaLocation | ctePropertyMarkerID
	}
	for i := 'a'; i <= 'z'; i++ {
		cteByteProperties[i] |= ctePropertyUnquotedMid | ctePropertyUnquotedStart | ctePropertyAZ | ctePropertyAreaLocation | ctePropertyMarkerID
	}
	for i := 'A'; i <= 'Z'; i++ {
		cteByteProperties[i] |= ctePropertyUnquotedMid | ctePropertyUnquotedStart | ctePropertyAZ | ctePropertyAreaLocation | ctePropertyMarkerID
	}
	for i := 0xc0; i < 0xf8; i++ {
		// UTF-8 initiator
		cteByteProperties[i] |= ctePropertyUnquotedMid | ctePropertyUnquotedStart
	}
	for i := 0x80; i < 0xc0; i++ {
		// UTF-8 continuation
		cteByteProperties[i] |= ctePropertyUnquotedMid
	}
	// TODO: Completely invalid bytes?

	cteByteProperties['='] |= ctePropertyObjectEnd
	cteByteProperties[';'] |= ctePropertyObjectEnd
	cteByteProperties['['] |= ctePropertyObjectEnd
	cteByteProperties[']'] |= ctePropertyObjectEnd
	cteByteProperties['{'] |= ctePropertyObjectEnd
	cteByteProperties['}'] |= ctePropertyObjectEnd
	cteByteProperties[')'] |= ctePropertyObjectEnd
	cteByteProperties['('] |= ctePropertyObjectEnd
	cteByteProperties['<'] |= ctePropertyObjectEnd
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
