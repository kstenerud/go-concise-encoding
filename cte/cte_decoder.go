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
	"io"
	"math"
	"strings"

	"github.com/kstenerud/go-concise-encoding/debug"
	"github.com/kstenerud/go-concise-encoding/events"
	"github.com/kstenerud/go-concise-encoding/internal/common"
	"github.com/kstenerud/go-concise-encoding/options"
)

// Decode a CTE document from reader, sending all data events to eventReceiver.
// If options is nil, default options will be used.
func Decode(reader io.Reader, eventReceiver events.DataEventReceiver, options *options.CTEDecoderOptions) (err error) {
	return NewDecoder(reader, eventReceiver, options).Decode()
}

// Decodes CTE documents.
type Decoder struct {
	eventReceiver  events.DataEventReceiver
	buffer         CTEReadBuffer
	containerState []cteDecoderState
	currentState   cteDecoderState
	options        options.CTEDecoderOptions
}

// Create a new CTE decoder, which will read from reader and send data events
// to nextReceiver. If options is nil, default options will be used.
func NewDecoder(reader io.Reader, eventReceiver events.DataEventReceiver, options *options.CTEDecoderOptions) *Decoder {
	_this := &Decoder{}
	_this.Init(reader, eventReceiver, options)
	return _this
}

// Initialize this decoder, which will read from reader and send data events
// to nextReceiver. If options is nil, default options will be used.
func (_this *Decoder) Init(reader io.Reader, eventReceiver events.DataEventReceiver, options *options.CTEDecoderOptions) {
	_this.options = *options.WithDefaultsApplied()
	_this.buffer.Init(reader, _this.options.BufferSize, chooseLowWater(_this.options.BufferSize))
	_this.eventReceiver = eventReceiver
}

// Run the complete decode process. The document and data receiver specified
// when initializing the decoder will be used.
func (_this *Decoder) Decode() (err error) {
	defer func() {
		if !debug.DebugOptions.PassThroughPanics {
			if r := recover(); r != nil {
				err = r.(error)
			}
		}
	}()

	_this.buffer.RefillIfNecessary()

	_this.currentState = cteDecoderStateAwaitObject

	// Forgive initial whitespace even though it's technically not allowed
	_this.buffer.SkipWhitespace()
	_this.buffer.EndToken()

	// TODO: Inline containers etc
	_this.handleVersion()

	for !_this.buffer.IsEndOfDocument() {
		_this.handleNextState()
		_this.buffer.RefillIfNecessary()
	}
	_this.eventReceiver.OnEndDocument()
	return
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

	version, bigVersion, digitCount := _this.buffer.DecodeDecimalInteger(0, nil)
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
		_this.eventReceiver.OnString(string(_this.buffer.GetToken()))
		_this.endObject()
		return
	}

	// Bytes, Custom, URI
	if _this.buffer.GetTokenLength() == 1 && _this.buffer.PeekByteNoEOD() == '"' {
		_this.buffer.AdvanceByte()
		initiator := _this.buffer.GetTokenFirstByte()
		switch initiator {
		case 'b':
			_this.eventReceiver.OnBytes(_this.buffer.DecodeHexBytes())
			_this.endObject()
			return
		case 'c':
			_this.eventReceiver.OnCustomBinary(_this.buffer.DecodeHexBytes())
			_this.endObject()
			return
		case 't':
			_this.eventReceiver.OnCustomText(string(_this.buffer.DecodeCustomText()))
			_this.endObject()
			return
		case 'u':
			_this.eventReceiver.OnURI(string(_this.buffer.DecodeURI()))
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
	_this.eventReceiver.OnString(string(_this.buffer.DecodeQuotedString()))
	_this.endObject()
}

func (_this *Decoder) handlePositiveNumeric() {
	coefficient, bigCoefficient, digitCount := _this.buffer.DecodeDecimalInteger(0, nil)

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
		value, bigValue := _this.buffer.DecodeDecimalFloat(1, coefficient, bigCoefficient, digitCount)
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
		token := strings.ToLower(string(_this.buffer.GetSubtoken()))
		if token != "inf" {
			_this.buffer.Errorf("Unknown named value: %v", token)
		}
		_this.eventReceiver.OnFloat(math.Inf(-1))
		_this.endObject()
		return
	}

	coefficient, bigCoefficient, digitCount := _this.buffer.DecodeDecimalInteger(0, nil)

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
		value, bigValue := _this.buffer.DecodeDecimalFloat(-1, coefficient, bigCoefficient, digitCount)
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
		v := _this.buffer.DecodeBinaryInteger()
		_this.buffer.AssertAtObjectEnd("binary integer")
		_this.eventReceiver.OnPositiveInt(v)
		_this.endObject()
	case 'o':
		v := _this.buffer.DecodeOctalInteger()
		_this.buffer.AssertAtObjectEnd("octal integer")
		_this.eventReceiver.OnPositiveInt(v)
		_this.endObject()
	case 'x':
		v, digitCount := _this.buffer.DecodeHexInteger(0)
		if _this.buffer.PeekByteAllowEOD() == '.' {
			_this.buffer.AdvanceByte()
			fv := _this.buffer.DecodeHexFloat(1, v, digitCount)
			_this.buffer.AssertAtObjectEnd("hex float")
			_this.eventReceiver.OnFloat(fv)
			_this.endObject()
		} else {
			_this.buffer.AssertAtObjectEnd("hex integer")
			_this.eventReceiver.OnPositiveInt(v)
			_this.endObject()
		}
	case '.':
		value, bigValue := _this.buffer.DecodeDecimalFloat(1, 0, nil, 0)
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
		v := _this.buffer.DecodeBinaryInteger()
		_this.buffer.AssertAtObjectEnd("binary integer")
		_this.eventReceiver.OnNegativeInt(v)
		_this.endObject()
	case 'o':
		v := _this.buffer.DecodeOctalInteger()
		_this.buffer.AssertAtObjectEnd("octal integer")
		_this.eventReceiver.OnNegativeInt(v)
		_this.endObject()
	case 'x':
		v, digitCount := _this.buffer.DecodeHexInteger(0)
		if _this.buffer.PeekByteAllowEOD() == '.' {
			_this.buffer.AdvanceByte()
			fv := _this.buffer.DecodeHexFloat(-1, v, digitCount)
			_this.buffer.AssertAtObjectEnd("hex float")
			_this.eventReceiver.OnFloat(fv)
			_this.endObject()
		} else {
			_this.buffer.AssertAtObjectEnd("hex integer")
			_this.eventReceiver.OnNegativeInt(v)
			_this.endObject()
		}
	case '.':
		value, bigValue := _this.buffer.DecodeDecimalFloat(-1, 0, nil, 0)
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
			_this.eventReceiver.OnString(contents)
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
		_this.eventReceiver.OnString(str)
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
		_this.eventReceiver.OnString(str)
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
			_this.eventReceiver.OnString(string(str))
		}
		_this.buffer.EndToken()
	case nextIsSingleLineComment:
		_this.eventReceiver.OnComment()
		contents := _this.buffer.DecodeSingleLineComment()
		if len(contents) > 0 {
			_this.eventReceiver.OnString(contents)
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
	token := strings.ToLower(string(_this.buffer.GetSubtoken()))
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
	_this.eventReceiver.OnVerbatimString(string(_this.buffer.DecodeVerbatimString()))
	_this.endObject()
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
		_this.eventReceiver.OnURI(string(_this.buffer.DecodeQuotedString()))
		_this.endObject()
		return
	}

	asString, asUint := _this.buffer.DecodeMarkerID()
	if len(asString) > 0 {
		_this.eventReceiver.OnString(asString)
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
		_this.eventReceiver.OnString(asString)
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
	cteDecoderStateTransitions[cteDecoderStateAwaitReferenceID] = cteDecoderStateAwaitObject

	cteDecoderStateHandlers[cteDecoderStateAwaitObject] = (*Decoder).handleObject
	cteDecoderStateHandlers[cteDecoderStateAwaitListItem] = (*Decoder).handleObject
	cteDecoderStateHandlers[cteDecoderStateAwaitCommentItem] = (*Decoder).handleCommentContent
	cteDecoderStateHandlers[cteDecoderStateAwaitMapKey] = (*Decoder).handleObject
	cteDecoderStateHandlers[cteDecoderStateAwaitMapKVSeparator] = (*Decoder).handleKVSeparator
	cteDecoderStateHandlers[cteDecoderStateAwaitMapValue] = (*Decoder).handleObject
	cteDecoderStateHandlers[cteDecoderStateAwaitMetaKey] = (*Decoder).handleObject
	cteDecoderStateHandlers[cteDecoderStateAwaitMetaKVSeparator] = (*Decoder).handleKVSeparator
	cteDecoderStateHandlers[cteDecoderStateAwaitMetaValue] = (*Decoder).handleObject
	cteDecoderStateHandlers[cteDecoderStateAwaitMarkupName] = (*Decoder).handleObject
	cteDecoderStateHandlers[cteDecoderStateAwaitMarkupKey] = (*Decoder).handleObject
	cteDecoderStateHandlers[cteDecoderStateAwaitMarkupKVSeparator] = (*Decoder).handleKVSeparator
	cteDecoderStateHandlers[cteDecoderStateAwaitMarkupValue] = (*Decoder).handleObject
	cteDecoderStateHandlers[cteDecoderStateAwaitMarkupItem] = (*Decoder).handleMarkupContent
	cteDecoderStateHandlers[cteDecoderStateAwaitReferenceID] = (*Decoder).handleObject
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
	charBasedHandlers['#'] = (*Decoder).handleReference
	charBasedHandlers['$'] = (*Decoder).handleInvalidChar
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
	charBasedHandlers[';'] = (*Decoder).handleInvalidChar
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
	charBasedHandlers['|'] = (*Decoder).handleMarkupContentBegin
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

	cteByteProperties['-'] |= ctePropertyUnquotedMid | ctePropertyAreaLocation
	cteByteProperties['+'] |= ctePropertyUnquotedMid | ctePropertyAreaLocation
	cteByteProperties['.'] |= ctePropertyUnquotedMid
	cteByteProperties[':'] |= ctePropertyUnquotedMid
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
